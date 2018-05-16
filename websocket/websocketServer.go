package websocket

import (
	"encoding/json"

	"net/http"

	"github.com/accedian/adh-gather/gather"
	"github.com/accedian/adh-gather/logger"
	"github.com/gorilla/websocket"
)

// ServerStruct - Struct containing the connections for various websocket clients
type ServerStruct struct {
	ConnectionMap map[string]*websocket.Conn
	Upgrader      websocket.Upgrader
}

// ConnectorMessage is a format for communicating to connector instances
type ConnectorMessage struct {
	ConnectorID string
	MsgType     string
	ErrorMsg    string
	ErrorCode   int
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
		case "Heartbeat":
			{
				logger.Log.Infof("Received Heartbeat from Connector with ID: %s", msg.ConnectorID)
			}
		case "Data":
			{
				logger.Log.Infof("Received Heartbeat from Connector with ID: %s", msg.ConnectorID)
			}
		default:
			{
				logger.Log.Errorf("Error from connector: %v. Error: %v. Message: %s", msg.ConnectorID, msg.ErrorCode, msg.ErrorMsg)
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

	wsServer := &ServerStruct{
		ConnectionMap: make(map[string]*websocket.Conn),
		Upgrader: websocket.Upgrader{
			ReadBufferSize:    1024,
			WriteBufferSize:   1024,
			EnableCompression: true,
		},
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
