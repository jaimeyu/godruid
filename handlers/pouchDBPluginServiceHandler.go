package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"

	db "github.com/accedian/adh-gather/datastore"
	"github.com/accedian/adh-gather/datastore/couchDB"
	"github.com/accedian/adh-gather/datastore/inMemory"
	"github.com/accedian/adh-gather/gather"
	"github.com/accedian/adh-gather/logger"
	"github.com/accedian/adh-gather/server"
	"github.com/gorilla/mux"
)

// PouchDBPluginServiceHandler - handler of logic related to calls for the
// pass through PouchDB Plugin Service.
type PouchDBPluginServiceHandler struct {
	pouchPluginDB db.PouchDBPluginServiceDatastore
	routes        []server.Route
}

// CreatePouchDBPluginServiceHandler - used to create a PouchDB plugin service handler
// which handles calls to the PouchDB Plugin Service
func CreatePouchDBPluginServiceHandler() *PouchDBPluginServiceHandler {
	result := new(PouchDBPluginServiceHandler)

	// Seteup the DB implementation based on configuration
	db, err := getPouchDBPluginServiceDatastore()
	if err != nil {
		logger.Log.Fatalf("Unable to instantiate PouchDBPluginServiceHandler: %s", err.Error())
	}
	result.pouchPluginDB = db

	result.routes = []server.Route{
		server.Route{
			Name:        "CheckAvailability",
			Method:      "GET",
			Pattern:     "/",
			HandlerFunc: result.CheckAvailability,
		},

		server.Route{
			Name:        "GetChanges",
			Method:      "GET",
			Pattern:     "/{dbname}/_changes",
			HandlerFunc: result.GetChanges,
		},

		server.Route{
			Name:        "StoreDBSyncCheckpoint",
			Method:      "PUT",
			Pattern:     "/{dbname}/_local/{docid}",
			HandlerFunc: result.StoreDBSyncCheckpoint,
		},

		server.Route{
			Name:        "GetDBSyncCheckpoint",
			Method:      "GET",
			Pattern:     "/{dbname}/_local/{docid}",
			HandlerFunc: result.GetDBSyncCheckpoint,
		},
	}

	return result
}

func getPouchDBPluginServiceDatastore() (db.PouchDBPluginServiceDatastore, error) {
	cfg, err := gather.GetActiveConfig()
	if err != nil {
		return nil, fmt.Errorf("Falied to instantiate PouchDBPluginServiceHandler: %s", err.Error())
	}

	dbType := cfg.ServerConfig.StartupArgs.PouchPluginDB
	switch dbType {
	case gather.COUCH:
		logger.Log.Debug("PouchDBPluginService DB is using CouchDB Implementation")
		return couchDB.CreatePouchDBServiceDAO()
	case gather.MEM:
		logger.Log.Debug("PouchDBPluginService DB is using InMemory Implementation")
		return inMemory.CreatePouchDBPluginServiceDAO()
	}

	return nil, errors.New("No DB implementation provided for Admin Service. Check configuration")
}

// GetChanges - provides access to the Changes feed of the provided DB.
// See http://docs.couchdb.org/en/2.1.1/api/database/changes.html for details
// on the API format.
func (psh *PouchDBPluginServiceHandler) GetChanges(w http.ResponseWriter, r *http.Request) {
	// Validate the request to ensure this operation is valid:

	urlParts := strings.Split(r.URL.Path, "/")
	dbName := urlParts[1]

	logger.Log.Infof("Looking for changes from DB %s", dbName)

	//Issue request to DAO Layer to faccess the Changes Feed
	queryParams := r.URL.Query()
	result, err := psh.pouchPluginDB.GetChanges(dbName, &queryParams)
	if err != nil {
		http.Error(w, fmt.Sprintf("Unable to retrieve %s: %s", db.ChangeFeedStr, err.Error()), http.StatusBadRequest)
		return
	}

	// Succesfully fetched the Changes Feed, return the result.
	logger.Log.Infof("Successfully accessed %s changes from DB %s\n", db.ChangeFeedStr, dbName)
	response, err := json.Marshal(result)
	if err != nil {
		http.Error(w, fmt.Sprintf("Unable to retrieve %s: %s", db.ChangeFeedStr, err.Error()), http.StatusBadRequest)
		return
	}

	fmt.Fprintf(w, string(response))
}

func (psh *PouchDBPluginServiceHandler) CheckAvailability(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "CheckAvailability hit!")
}

func (psh *PouchDBPluginServiceHandler) StoreDBSyncCheckpoint(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "StoreDBSyncCheckpoint hit!")
}

func (psh *PouchDBPluginServiceHandler) GetDBSyncCheckpoint(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "GetDBSyncCheckpoint hit!")
}

// RegisterAPIHandlers - will bind any REST API routes defined in this service
// to the passed in request multiplexor.
func (psh *PouchDBPluginServiceHandler) RegisterAPIHandlers(router *mux.Router) {
	for _, route := range psh.routes {
		router.
			Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			Handler(route.HandlerFunc)
	}
}

// // GetChanges - used to subscribe to the changes feed from CouchDB.
// func (psh *PouchDBPluginServiceHandler) GetChanges(ctx context.Context, dbChangesRequest *pb.DBChangesRequest) (*pb.DBChangesResponse, error) {
// 	// Validate the request to ensure this operation is valid:

// 	logger.Log.Infof("Looking for changes from DB %s", dbChangesRequest.GetDbName())

// 	// Issue request to DAO Layer to faccess the Changes Feed
// 	result, err := psh.pouchPluginDB.GetChanges(dbChangesRequest)
// 	if err != nil {
// 		return nil, fmt.Errorf("Unable to retrieve %s: %s", db.ChangeFeedStr, err.Error())
// 	}

// 	// Succesfully fetched the Changes Feed, return the result.
// 	logger.Log.Infof("Retrieved %d changes from %s\n", len(result.GetResults()), db.ChangeFeedStr)
// 	return result, nil
// }

// // CheckAvailablility - ping the CouchDB server for availability.
// func (psh *PouchDBPluginServiceHandler) CheckAvailablility(ctx context.Context, noValue *emp.Empty) (*pb.DBAvailableResponse, error) {
// 	// Validate the request to ensure this operation is valid:

// 	logger.Log.Info("Checking availability of the CouchDB Server")

// 	// Issue request to DAO Layer to check DB Availability
// 	result, err := psh.pouchPluginDB.CheckAvailablility()
// 	if err != nil {
// 		return nil, fmt.Errorf("Unable to retrieve Couch DB availability: %s", err.Error())
// 	}

// 	// DB is available, send the response
// 	logger.Log.Info("CouchDB server is available\n")
// 	return result, nil
// }

// // StoreDBSyncCheckpoint - stores data used to keep track of sync position between pouch and couch DBs.
// func (psh *PouchDBPluginServiceHandler) StoreDBSyncCheckpoint(ctx context.Context, dbCheckpoint *pb.DBSyncCheckpoint) (*pb.DBSyncCheckpointPutResponse, error) {
// 	// Validate the request to ensure this operation is valid:

// 	logger.Log.Infof("Storing %s: %s", db.DBSyncCheckpointStr, dbCheckpoint.GetXId())

// 	// Issue request to DAO Layer to check DB Availability
// 	result, err := psh.pouchPluginDB.StoreDBSyncCheckpoint(dbCheckpoint)
// 	if err != nil {
// 		return nil, fmt.Errorf("Unable to store %s: %s", db.DBSyncCheckpointStr, err.Error())
// 	}

// 	// DB Sync Checkpoint stored, send the response
// 	logger.Log.Infof("%s %s stored.\n", db.DBSyncCheckpointStr, dbCheckpoint.GetXId())
// 	return result, nil
// }

// // GetDBSyncCheckpoint - retrieves a previously stored DB sync checkpoint between pouch and couch DB.
// func (psh *PouchDBPluginServiceHandler) GetDBSyncCheckpoint(ctx context.Context, dbCheckpointID *pb.DBSyncCheckpointId) (*pb.DBSyncCheckpoint, error) {
// 	// Validate the request to ensure this operation is valid:

// 	logger.Log.Infof("Retrieving %s: %s", db.DBSyncCheckpointStr, dbCheckpointID.GetXId())

// 	// Issue request to DAO Layer to check DB Availability
// 	result, err := psh.pouchPluginDB.GetDBSyncCheckpoint(dbCheckpointID, true)
// 	if err != nil {
// 		return nil, fmt.Errorf("Unable to retrieve %s: %s", db.DBSyncCheckpointStr, err.Error())
// 	}

// 	// DB Sync Checkpoint retrieved, send the response
// 	logger.Log.Infof("%s %s retrieved.\n", db.DBSyncCheckpointStr, dbCheckpointID.GetXId())
// 	return result, nil
// }
