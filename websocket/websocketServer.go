package websocket

import (
	"encoding/json"
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
	HeartbeatChan   chan string
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
	MsgID       int
	ErrorCode   int
	Data        []byte
	Zone        string
}

// Reader - Function which reads websocket messages coming through the websocket connection
func (wsServer *ServerStruct) Reader(ws *websocket.Conn, connectorID string) {

	for {
		_, p, err := ws.ReadMessage()
		if err != nil {
			logger.Log.Errorf("Lost connection to Connector with ID: %v. Error: %v", connectorID, err)

			tenantID := wsServer.ConnectionMeta[connectorID].TenantID
			connectorConfigs, _ := wsServer.TenantDB.GetAllTenantConnectorsByInstanceID(tenantID, connectorID)

			// Lost connection to the connector, so we need to remove connectorInstanceID from any configs that have it
			for _, c := range connectorConfigs {
				c.ID = datastore.GetDataIDFromFullID(c.ID)
				c.ConnectorInstanceID = ""
				_, err = wsServer.TenantDB.UpdateTenantConnector(c)
				if err != nil {
					logger.Log.Errorf("Unable to remove connectorInstanceID from ConnectorConfig: %v, for tenant: %v. Error: %v", c.ID, tenantID, err)
					break
				}
			}

			// delete from connection map
			delete(wsServer.ConnectionMeta, connectorID)

			break
		}
		msg := &ConnectorMessage{}

		json.Unmarshal(p, msg)

		switch msg.MsgType {
		case "Config":
			{
				tenantID := msg.Tenant
				zone := msg.Zone
				var configs []*tenant.Connector

				logger.Log.Infof("Received config request from Connector with ID: %s", connectorID)
				wsServer.Lock.Lock()
				wsServer.ConnectionMeta[connectorID].TenantID = tenantID
				wsServer.Lock.Unlock()

				// Check if ConnectorInstances has an entry for connectorID
				connectorInstance, err := wsServer.TenantDB.GetTenantConnectorInstance(tenantID, connectorID)
				if err != nil {
					logger.Log.Errorf("Unable to retrieve connector instance for tenant: %v and connectorID: %v. Error: %v", tenantID, connectorID, err)
				}

				// The following is logic for choosing which connection to give to an incoming connector:
				// if connector hasn't been added to connector instances, add it
				if connectorInstance == nil {
					connectorInstance = &tenant.ConnectorInstance{
						ID:       connectorID,
						Hostname: msg.Hostname,
						TenantID: tenantID,
					}

					wsServer.TenantDB.CreateTenantConnectorInstance(connectorInstance)
					if err != nil {
						logger.Log.Errorf("Unable to create TenantConnectorInstance for tenant: %v. Error: %v", tenantID, err)
						break
					}

					// get all available configs
					configs, err = wsServer.TenantDB.GetAllAvailableTenantConnectors(tenantID, zone)

				} else {
					// We have a connectorInstance for the incoming connectorID

					// find configs that have an instanceID that matches connectorID
					configs, err = wsServer.TenantDB.GetAllTenantConnectorsByInstanceID(tenantID, connectorID)

					// if none of the configs are used by the connector instances, get available configs
					if len(configs) == 0 {
						// get all available configs
						configs, err = wsServer.TenantDB.GetAllAvailableTenantConnectors(tenantID, zone)
					}

					// if there are no available configs, make sure that the used configs are being used by a valid connector
					// and not by any stale connectors (Could happen if gather crashes)

					if len(configs) == 0 {
						allConfigs, err := wsServer.TenantDB.GetAllTenantConnectors(tenantID, zone)
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
					}
				}

				if err != nil {
					logger.Log.Errorf("Unable to find connectors for tenant: %v and zone: %v", tenantID, zone)
					break
				}

				// take the first available config, and assign a connector instance ID to it
				if len(configs) == 0 {
					logger.Log.Errorf("No available configurations for Connector with ID: %v", connectorID)
					break
				}

				// pick the first available connector
				selectedConfig := configs[0]

				logger.Log.Debugf("Sending following config: %v to connector with ID: %s", selectedConfig, connectorID)

				// remove the couchDB type prefix from the ID
				selectedConfig.ID = datastore.GetDataIDFromFullID(selectedConfig.ID)
				selectedConfig.ConnectorInstanceID = connectorID

				// Convert our data to JSON
				configJSON, _ := json.Marshal(selectedConfig)

				returnMsg := &ConnectorMessage{
					MsgType: "Config",
					Data:    configJSON,
				}

				msgJSON, _ := json.Marshal(returnMsg)

				// Send the config to the connector
				err = wsServer.ConnectionMeta[connectorID].Connection.WriteMessage(websocket.BinaryMessage, msgJSON)
				if err != nil {
					logger.Log.Errorf("Error sending configuration to Connector with ID: %v", connectorID)
					break
				}

				// After successfully sending config to connector, update ConnectorConfig with the instance iD
				_, err = wsServer.TenantDB.UpdateTenantConnector(selectedConfig)
				if err != nil {
					logger.Log.Errorf("Unable to update TenantConnector: %v, for tenant: %v. Error: %v", selectedConfig.ID, tenantID, err)
					break
				}

			}
		case "Heartbeat":
			{
				logger.Log.Debugf("Received Heartbeat from Connector with ID: %s", connectorID)

				// Notify our hearbeat monitor that we've received a heartbeat
				wsServer.ConnectionMeta[connectorID].HeartbeatChan <- connectorID

			}
		default:
			{
				logger.Log.Errorf("Error from connector: %v. Error: %v. Message: %s", connectorID, msg.ErrorCode, msg.ErrorMsg)
			}
		}

		// Cleanup any connections marked for closing
		if wsServer.ConnectionMeta[connectorID].CloseConnection {
			wsServer.ConnectionMeta[connectorID].Connection.Close()
			break
		}
	}
}

func (wsServer *ServerStruct) serveWs(w http.ResponseWriter, r *http.Request) {
	ws, err := wsServer.Upgrader.Upgrade(w, r, nil)
	connectorID := r.Header["Connectorid"][0]
	heartbeatChan := make(chan string)

	if err != nil {
		if _, ok := err.(websocket.HandshakeError); !ok {
			logger.Log.Errorf("Could not upgrade websocket connection from connector %v. Error: %v", connectorID, err)
		}
		return
	}

	connectionMeta := &ConnectionMeta{
		Connection:    ws,
		HeartbeatChan: heartbeatChan,
	}

	wsServer.ConnectionMeta[connectorID] = connectionMeta

	logger.Log.Infof("Connector with ID: %v, successfully connected.", connectorID)

	// If we don't get a heartbeat message in the required amount of time, mark connection for closing
	go func() {

		maxTimeWithoutHeartbeat := wsServer.Config.GetInt(gather.CK_connector_maxSecondsWithoutHeartbeat.String())
		ticker := time.NewTicker(time.Duration(maxTimeWithoutHeartbeat) * time.Second)

		for {
			select {
			case <-heartbeatChan:
				ticker = time.NewTicker(time.Duration(maxTimeWithoutHeartbeat) * time.Second)
			case <-ticker.C:
				// If connection gets marked for closing, exit this loop. It will be recreated for the new connection
				logger.Log.Errorf("Haven't received heartbeat from Connector: %s for %v seconds. Resetting connection.", connectorID, maxTimeWithoutHeartbeat)
				if wsServer.ConnectionMeta[connectorID] != nil {
					wsServer.Lock.Lock()
					wsServer.ConnectionMeta[connectorID].CloseConnection = true
					wsServer.Lock.Unlock()
				}
				return
			}
		}

	}()

	wsServer.Reader(ws, connectorID)
}

// Listens to config changes and sends the new config to the correct connector
func (wsServer *ServerStruct) listenToConnectorChanges() {

	// if a connector config changes, push it to the connector
	for config := range wsServer.TenantDB.GetConnectorUpdateChan() {
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

	go wsServer.listenToConnectorChanges()

	return wsServer
}
