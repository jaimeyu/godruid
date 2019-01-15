package websocket

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"sync"
	"time"

	"net/http"

	"github.com/accedian/adh-gather/config"
	"github.com/accedian/adh-gather/datastore"
	"github.com/accedian/adh-gather/gather"
	"github.com/accedian/adh-gather/logger"
	"github.com/accedian/adh-gather/models/tenant"
	"github.com/gorilla/websocket"
)

// ConnectionInfo - Struct containing metadata about the websocket connection
type ConnectionMeta struct {
	Connection      *websocket.Conn
	TenantID        string
	LastHeartbeat   int64
	CloseConnection bool // mark connection to be closed as soon as possible
}

// ServerStruct - Struct containing the connections for various websocket clients
type ServerStruct struct {
	ConnectionMeta map[string]*ConnectionMeta
	Upgrader       websocket.Upgrader
	TenantDB       datastore.TenantServiceDatastore
	Config         config.Provider
	Lock           sync.Mutex
}

// ConnectorMessage is a format for communicating to connector instances

type ConnectorMessage struct {
	Filename    string
	Tenant      string
	Hostname    string
	ConnectorID string
	DataType    string
	MsgType     string
	ErrorMsg    string
	ObjectType  string
	MsgID       int
	ErrorCode   int
	Data        []byte
	Zone        string
}

type PtExport struct {
	XMLName  xml.Name `xml:"Ptexport"`
	Version  string   `xml:"version,attr"`
	Response Response
}

type Response struct {
	XMLName xml.Name `xml:"Response"`
	Sess    []Sess   `xml:"Sess"`
}
type Sess struct {
	XMLName xml.Name `xml:"Sess"`
	CID     string   `xml:"cid,attr"`
	SID     string   `xml:"sid,attr"`
	Name    string   `xml:"name,attr"`
	Type    string   `xml:"type,attr"`
}

var (
	batchSize = -1
	// TODO: PEYO this is temporary, this mapping will live outside of gather
	objectTypes = make(map[string]string)
)

func setBatchSize() {
	if batchSize < 0 {
		cfg := gather.GetConfig()
		batchSize = cfg.GetInt(gather.CK_server_datastore_batchsize.String())
		logger.Log.Debugf("Using BatchSize of %d", batchSize)
	}
}

// Reader - Function which reads websocket messages coming through the websocket connection
func (wsServer *ServerStruct) Reader(ws *websocket.Conn, connectorID string) {

	for ws != nil {
		_, p, err := ws.ReadMessage()

		if err != nil {
			logger.Log.Errorf("Lost connection to Connector with ID: %v. Error: %v", connectorID, err)

			if wsServer.ConnectionMeta[connectorID] == nil {
				logger.Log.Debugf("Connection to Connector: %s has already been deleted", connectorID)
				break
			}

			tenantID := wsServer.ConnectionMeta[connectorID].TenantID
			connectorConfigs, _ := wsServer.TenantDB.GetAllTenantConnectorConfigsByInstanceID(tenantID, connectorID)

			// Lost connection to the connector, so we need to remove connectorInstanceID from any configs that have it
			for _, c := range connectorConfigs {
				c.ID = datastore.GetDataIDFromFullID(c.ID)
				c.ConnectorInstanceID = ""

				_, err = wsServer.TenantDB.UpdateTenantConnectorConfig(c)
				if err != nil {
					logger.Log.Errorf("Unable to remove connectorInstanceID from ConnectorConfig: %v, for tenant: %v. Error: %v", c.ID, tenantID, err)
					break
				}
			}

			// Remove ConnectorInstance
			_, err = wsServer.TenantDB.DeleteTenantConnectorInstance(tenantID, connectorID)
			if err != nil {
				logger.Log.Errorf("Unable to delete connectorInstance with ID: %v, for tenant: %v. Error: %v", connectorID, tenantID, err)
				break
			}
			break
		}
		msg := &ConnectorMessage{}

		json.Unmarshal(p, msg)

		switch msg.MsgType {
		case "Config":
			{
				tenantID := msg.Tenant
				zone := msg.Zone
				var configs []*tenant.ConnectorConfig

				logger.Log.Infof("Received config request from Connector with ID: %s", connectorID)

				// Check if ConnectorInstances has an entry for connectorID
				connectorInstance, err := wsServer.TenantDB.GetTenantConnectorInstance(tenantID, connectorID)
				if err != nil {
					logger.Log.Errorf("Unable to retrieve connector instance for tenant: %v and connectorID: %v. Error: %v", tenantID, connectorID, err)
				}

				// The following is logic for choosing which configuration to give to an incoming connector:
				// if connector hasn't been added to connector instances, add it
				if connectorInstance == nil {
					connectorInstance = &tenant.ConnectorInstance{
						ID:       connectorID,
						Hostname: msg.Hostname,
						TenantID: tenantID,
						Status:   "connected",
					}

					_, err = wsServer.TenantDB.CreateTenantConnectorInstance(connectorInstance)
					if err != nil {
						logger.Log.Errorf("Unable to create TenantConnectorInstance for tenant: %v. Error: %v", tenantID, err)
						break
					}

					// get all available configs
					configs, err = wsServer.TenantDB.GetAllAvailableTenantConnectorConfigs(tenantID, zone)

				} else {
					// We have a connectorInstance for the incoming connectorID

					// find configs that have an instanceID that matches connectorID
					configs, err = wsServer.TenantDB.GetAllTenantConnectorConfigsByInstanceID(tenantID, connectorID)

					// if none of the configs are used by the connector instances, get available configs
					if len(configs) == 0 {
						// get all available configs
						configs, err = wsServer.TenantDB.GetAllAvailableTenantConnectorConfigs(tenantID, zone)
					}
				}

				if err != nil {
					logger.Log.Errorf("Unable to find connectors for tenant: %v and zone: %v", tenantID, zone)
					break
				}

				// if there are no available configs, make sure that the used configs are being used by a valid connector
				// and not by any stale connectors (Could happen if gather crashes)
				if len(configs) == 0 {
					allConfigs, err := wsServer.TenantDB.GetAllTenantConnectorConfigs(tenantID, zone)
					if err != nil {
						logger.Log.Errorf("Unable to find connectors for tenant: %v and zone: %v", tenantID, zone)
						break
					}
					for _, c := range allConfigs {
						if wsServer.ConnectionMeta[c.ConnectorInstanceID] == nil {
							c.ConnectorInstanceID = ""
							configs = append(configs, c)
						}
					}

					if len(configs) == 0 {
						errMsg := fmt.Sprintf("No available configurations for Connector with ID: %v", connectorID)
						logger.Log.Error(errMsg)

						returnMsg := &ConnectorMessage{
							MsgType:  "Config",
							ErrorMsg: errMsg,
						}

						msgJSON, _ := json.Marshal(returnMsg)

						err = wsServer.ConnectionMeta[connectorID].Connection.WriteMessage(websocket.BinaryMessage, msgJSON)

						if err != nil {
							logger.Log.Errorf("Error sending configuration to Connector with ID: %v", connectorID)
							break
						}
						break
					}

				}

				// pick the first available connector
				selectedID := datastore.GetDataIDFromFullID(configs[0].ID)
				selectedConfig, _ := wsServer.TenantDB.GetTenantConnectorConfig(tenantID, selectedID)

				logger.Log.Infof("Sending following config: %v to connector with ID: %s", selectedConfig, connectorID)

				// remove the couchDB type prefix from the ID
				selectedConfig.ID = datastore.GetDataIDFromFullID(selectedConfig.ID)
				selectedConfig.ConnectorInstanceID = connectorID

				wsServer.Lock.Lock()
				wsServer.ConnectionMeta[connectorID].TenantID = tenantID
				wsServer.ConnectionMeta[connectorID].LastHeartbeat = time.Now().Unix()
				wsServer.Lock.Unlock()

				// After successfully sending config to connector, update ConnectorConfig with the instance iD
				_, err = wsServer.TenantDB.UpdateTenantConnectorConfig(selectedConfig)
				if err != nil {
					logger.Log.Errorf("Unable to update TenantConnectorConfig: %v, for tenant: %v. Error: %v", selectedConfig.ID, tenantID, err)
					break
				}
			}
		case "SessionUpdate":
			{
				setBatchSize()

				logger.Log.Infof("Received Session Update from Connector with ID: %s", connectorID)
				monitoredObjectNames := PtExport{}

				if err := xml.Unmarshal(msg.Data, &monitoredObjectNames); err != nil {
					logger.Log.Errorf("Error unmarshalling session names from connector: %v. Error: %v. Message: %s", connectorID, msg.ErrorCode, msg.ErrorMsg)
				}

				// Create a mapping of MonitoredObjectID to Session data for fast lookup
				sessionDataMap := map[string]Sess{}
				tenantID := msg.Tenant

				// Make sure we only handle the bulk requests in batches of 1000
				moFetchBuffer := make([]string, 0, batchSize)
				moUpdateBuffer := make([]*tenant.MonitoredObject, 0)
				for i, m := range monitoredObjectNames.Response.Sess {

					monObjID := m.CID + "-" + m.SID

					// Log scenario where Monitored Object name is an empty string, but still make the update.
					if err != nil {
						logger.Log.Errorf("Unable to update name of MonitoredObject with ID: %v, for tenant: %v. Error: %v", monObjID, tenantID, err)
					}

					// Update the fast lookup mapping:
					sessionDataMap[monObjID] = m

					currentIndexInRange := i % batchSize

					// Send a fetch request if buffer is full
					if i != 0 && currentIndexInRange == 0 {
						logger.Log.Debugf("Retrieving batch of Monitored Objects for Tenant %s from %d IDs", tenantID, batchSize)
						objectsToUpdate, err := wsServer.TenantDB.GetAllMonitoredObjectsInIDList(tenantID, moFetchBuffer)
						if err != nil {
							logger.Log.Errorf("Unable to retrieve batch of Monitored Objects for tenant: %s. Error: %s", tenantID, err.Error())
						}
						moUpdateBuffer = append(moUpdateBuffer, objectsToUpdate...)
						moFetchBuffer = make([]string, 0, batchSize) // Reset the fetch buffer
					}

					// Add the current object to the fetch request list:
					moFetchBuffer = append(moFetchBuffer, monObjID)
				}

				// Issue request to get any remaining items:
				logger.Log.Debugf("Retrieving last batch of Monitored Objects for Tenant %s from %d IDs", tenantID, len(moFetchBuffer))
				lastBatchFromFetch, err := wsServer.TenantDB.GetAllMonitoredObjectsInIDList(tenantID, moFetchBuffer)
				if err != nil {
					logger.Log.Errorf("Unable to retrieve last batch of Monitored Objects for tenant: %s. Error: %s", tenantID, err.Error())
				}
				moUpdateBuffer = append(moUpdateBuffer, lastBatchFromFetch...)

				// Process the update requests in batches
				moUpdateBatch := make([]*tenant.MonitoredObject, 0, batchSize)
				for i, m := range moUpdateBuffer {

					currentIndexInRange := i % batchSize

					// Send a batch of Monitored Objects for update if it is time.
					if i != 0 && currentIndexInRange == 0 {
						logger.Log.Debugf("Updating batch of %d Monitored Objects for Tenant %s", batchSize, tenantID)
						_, err = wsServer.TenantDB.BulkUpdateMonitoredObjects(tenantID, moUpdateBatch)
						if err != nil {
							logger.Log.Errorf("Unable to update batch of MonitoredObjects for tenant: %s. Error: %s", tenantID, err.Error())
						}
						moUpdateBatch = make([]*tenant.MonitoredObject, 0, batchSize)
					}

					// Otherwise, just add the record to the batch
					updateProperties := sessionDataMap[m.ID]
					m.ObjectName = updateProperties.Name
					m.ObjectType = objectTypes[updateProperties.Type]

					logger.Log.Debugf("Updating Monitored object %s to have name %s as per properties %v", m.ID, m.ObjectName, updateProperties)
					moUpdateBatch = append(moUpdateBatch, m)
				}

				// Issue request to update any remaining items:
				logger.Log.Debugf("Updating batch of %d Monitored Objects for Tenant %s", len(moUpdateBatch), tenantID)
				_, err = wsServer.TenantDB.BulkUpdateMonitoredObjects(tenantID, moUpdateBatch)
				if err != nil {
					logger.Log.Errorf("Unable to update last batch of MonitoredObjects for tenant: %s. Error: %s", tenantID, err.Error())
				}

				returnMsg := &ConnectorMessage{
					MsgType: "Session",
				}

				msgJSON, _ := json.Marshal(returnMsg)

				err = wsServer.ConnectionMeta[connectorID].Connection.WriteMessage(websocket.BinaryMessage, msgJSON)

			}
		default:
			{
				logger.Log.Errorf("Error from connector: %v. Error: %v. Message: %s", connectorID, msg.ErrorCode, msg.ErrorMsg)
			}
		}
	}
}

// Create initial websocket connection
func (wsServer *ServerStruct) serveWs(w http.ResponseWriter, r *http.Request) {
	ws, err := wsServer.Upgrader.Upgrade(w, r, nil)
	connectorID := r.Header["Connectorid"][0]

	if err != nil {
		if _, ok := err.(websocket.HandshakeError); !ok {
			logger.Log.Errorf("Could not upgrade websocket connection from connector %v. Error: %v", connectorID, err)
		}
		return
	}

	connectionMeta := &ConnectionMeta{
		Connection:    ws,
		LastHeartbeat: time.Now().Unix(),
	}
	wsServer.Lock.Lock()
	wsServer.ConnectionMeta[connectorID] = connectionMeta
	wsServer.Lock.Unlock()

	logger.Log.Infof("Connector with ID: %v, successfully connected.", connectorID)

	wsServer.Reader(ws, connectorID)
}

// Listens to config changes and sends the new config to the correct connector
func (wsServer *ServerStruct) listenToConnectorChanges() {

	// if a connector config changes, push it to the connector
	for config := range wsServer.TenantDB.GetConnectorConfigUpdateChan() {

		instanceID := config.ConnectorInstanceID
		meta := wsServer.ConnectionMeta[instanceID]
		if instanceID != "" && meta != nil {
			wsConn := meta.Connection
			configJSON, _ := json.Marshal(config)

			returnMsg := &ConnectorMessage{
				MsgType: "Config",
				Data:    configJSON,
			}

			msgJSON, _ := json.Marshal(returnMsg)

			// Send the config to the connector
			err := wsConn.WriteMessage(websocket.BinaryMessage, msgJSON)
			if err != nil {
				logger.Log.Errorf("Error sending configuration to Connector with ID: %v", instanceID)
				break
			}
		}
	}
}

// Server server waiting to accept websocket connections from the connector
func Server(tenantDB datastore.TenantServiceDatastore) *ServerStruct {

	objectTypes["8"] = "eth-lb"
	objectTypes["9"] = "eth-dm"
	objectTypes["10"] = "twamp-sf"
	objectTypes["11"] = "echo-udp"
	objectTypes["12"] = "echo-icmp"
	objectTypes["13"] = "eth-vs"
	objectTypes["16"] = "twamp-sl"

	cfg := gather.GetConfig()

	wsServer := &ServerStruct{
		ConnectionMeta: make(map[string]*ConnectionMeta),
		Upgrader: websocket.Upgrader{
			ReadBufferSize:    1024,
			WriteBufferSize:   1024,
			EnableCompression: true,
		},
		TenantDB: tenantDB,
		Config:   cfg,
	}

	http.HandleFunc("/wsstatus", wsServer.serveWs)

	host := cfg.GetString(gather.CK_server_websocket_ip.String())
	port := cfg.GetString(gather.CK_server_websocket_port.String())
	addr := host + ":" + port

	go func() {
		logger.Log.Infof("Starting Websocket Server on: %v:%v", host, port)
		http.ListenAndServe(addr, nil)
	}()

	// If we don't see heartbeats for the maximum time allowed, close the websocket connection
	go func() {
		maxSecondsWithoutHeartbeat := int64(cfg.GetInt("connector.maxSecondsWithoutHeartbeat"))

		heartbeatTicker := time.NewTicker(time.Duration(maxSecondsWithoutHeartbeat) * time.Second)
		for range heartbeatTicker.C {
			now := time.Now().Unix()
			for cID, meta := range wsServer.ConnectionMeta {
				// No hearbeat has been received for this connector, so we need to clean out its connection,
				// and clear out the connectorInstance, as well as the connectorInstanceID from the connector config.
				if now-meta.LastHeartbeat > maxSecondsWithoutHeartbeat {
					logger.Log.Errorf("No Heartbeat messages have been received from Connector with ID: %v for %v seconds. Terminating connection.", cID, maxSecondsWithoutHeartbeat)
					wsServer.Lock.Lock()
					wsServer.ConnectionMeta[cID].Connection.Close()
					delete(wsServer.ConnectionMeta, cID)
					wsServer.Lock.Unlock()
				}
			}
		}
	}()

	go wsServer.listenToConnectorChanges()

	return wsServer
}
