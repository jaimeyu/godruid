package websocket

import (
	"encoding/json"

	"net/http"

	"github.com/accedian/adh-gather/datastore"
	"github.com/accedian/adh-gather/datastore/couchDB"
	"github.com/accedian/adh-gather/gather"
	"github.com/accedian/adh-gather/logger"
	"github.com/accedian/adh-gather/models/tenant"
	"github.com/gorilla/websocket"
)

// ServerStruct - Struct containing the connections for various websocket clients
type ServerStruct struct {
	ConnectionMap map[string]*websocket.Conn
	Upgrader      websocket.Upgrader
	TenantDB      *couchDB.TenantServiceDatastoreCouchDB
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
			break
		}
		msg := &ConnectorMessage{}

		json.Unmarshal(p, msg)

		switch msg.MsgType {
		case "Config":
			{
				tenantID := msg.Tenant
				zone := msg.Zone

				logger.Log.Infof("Received config request from Connector with ID: %s", connectorID)

				// Check if ConnectorInstances has an entry for connectorID
				connectorInstance, err := wsServer.TenantDB.GetTenantConnectorInstance(tenantID, connectorID)
				if err != nil {
					logger.Log.Errorf("Unable to retrieve connector instance for tenant: %v and connectorID: %v. Error: %v", tenantID, connectorID, err)
				}

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
				}

				// get all available configs
				configs, err := wsServer.TenantDB.GetAllAvailableTenantConnectors(tenantID, zone)

				if err != nil {
					logger.Log.Errorf("Unable to find connectors for tenant: %v and zone: %v", tenantID, zone)
					break
				}

				// take the first available config, and assign a connector instance ID to it
				if len(configs) == 0 {
					logger.Log.Errorf("No available configurations for Connector with ID: %v", connectorID)
					break
				}

				selectedConfig := configs[0]

				selectedConfig.ID = datastore.GetDataIDFromFullID(selectedConfig.ID)
				selectedConfig.ConnectorInstanceID = connectorID

				// Send the config to the connector
				configJSON, _ := json.Marshal(selectedConfig)

				returnMsg := &ConnectorMessage{
					MsgType: "Config",
					Data:    configJSON,
				}

				msgJSON, _ := json.Marshal(returnMsg)

				err = wsServer.ConnectionMap[connectorID].WriteMessage(websocket.BinaryMessage, msgJSON)
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
				logger.Log.Infof("Received Heartbeat from Connector with ID: %s", connectorID)
			}
		default:
			{
				logger.Log.Errorf("Error from connector: %v. Error: %v. Message: %s", connectorID, msg.ErrorCode, msg.ErrorMsg)
			}
		}
	}
}

func (wsServer *ServerStruct) serveWs(w http.ResponseWriter, r *http.Request) {
	ws, err := wsServer.Upgrader.Upgrade(w, r, nil)
	connectorID := r.Header["Connectorid"][0]

	if err != nil {
		if _, ok := err.(websocket.HandshakeError); !ok {
			logger.Log.Errorf("Could not upgrade websocket connection from connector %v. Error: %v", connectorID, err)
		}
		return
	}

	wsServer.ConnectionMap[connectorID] = ws

	logger.Log.Infof("Connector with ID: %v, successfully connected.", connectorID)
	wsServer.Reader(ws, connectorID)
}

// Server server waiting to accept websocket connections from the connector
func Server() *ServerStruct {

	cfg := gather.GetConfig()

	tenantDB, err := couchDB.CreateTenantServiceDAO()
	if err != nil {
		logger.Log.Errorf("Could not create couchdb tenant DAO: %s", err.Error())
	}

	wsServer := &ServerStruct{
		ConnectionMap: make(map[string]*websocket.Conn),
		Upgrader: websocket.Upgrader{
			ReadBufferSize:    1024,
			WriteBufferSize:   1024,
			EnableCompression: true,
		},
		TenantDB: tenantDB,
	}

	http.HandleFunc("/wsstatus", wsServer.serveWs)

	host := cfg.GetString(gather.CK_server_websocket_ip.String())
	port := cfg.GetString(gather.CK_server_websocket_port.String())
	addr := host + ":" + port

	go func() {
		logger.Log.Infof("Starting Websocket Server on: %v:%v", host, port)
		http.ListenAndServe(addr, nil)
	}()

	return wsServer
}
