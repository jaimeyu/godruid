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

// Used as enum for retrieving parts of the PouchDB plugin URL
type pouchPluginURLPart int32

const (
	dbNameInURL     pouchPluginURLPart = 1
	dbMethodInURL   pouchPluginURLPart = 2
	documentIDInURL pouchPluginURLPart = 3
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

		server.Route{
			Name:        "GetDBRevisionDiff",
			Method:      "POST",
			Pattern:     "/{dbname}/_revs_diff",
			HandlerFunc: result.GetDBRevisionDiff,
		},

		server.Route{
			Name:        "BulkDBUpdate",
			Method:      "POST",
			Pattern:     "/{dbname}/_bulk_docs",
			HandlerFunc: result.BulkDBUpdate,
		},
	}

	return result
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
	// TODO: Validate the request to ensure this operation is valid:

	dbName := getDBFieldFromRequest(r, dbNameInURL)

	logger.Log.Infof("Looking for changes from DB %s", dbName)

	//Issue request to DAO Layer to access the Changes Feed
	queryParams := r.URL.Query()
	result, err := psh.pouchPluginDB.GetChanges(dbName, &queryParams)
	if err != nil {
		http.Error(w, fmt.Sprintf("Unable to retrieve %s: %s", db.ChangeFeedStr, err.Error()), http.StatusBadRequest)
		return
	}

	// Succesfully fetched the Changes Feed, return the result. See
	logger.Log.Infof("Successfully accessed %s changes from DB %s\n", db.ChangeFeedStr, dbName)
	response, err := json.Marshal(result)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error generating %s response: %s", db.ChangeFeedStr, err.Error()), http.StatusBadRequest)
		return
	}

	fmt.Fprintf(w, string(response))
}

// CheckAvailability - used to check if the CouchDB server is available.
// See http://docs.couchdb.org/en/2.1.1/api/server/common.html for the
// CouchDB documentation on this API.
func (psh *PouchDBPluginServiceHandler) CheckAvailability(w http.ResponseWriter, r *http.Request) {
	// TODO: Validate the request to ensure this operation is valid:

	logger.Log.Info("Checking for CouchDB availability")

	//Issue request to DAO Layer to access check availability
	result, err := psh.pouchPluginDB.CheckAvailability()
	if err != nil {
		http.Error(w, fmt.Sprintf("Error checking CouchDB availability: %s", err.Error()), http.StatusInternalServerError)
		return
	}

	// Succesfully accessed the couch server, return the result
	response, err := json.Marshal(result)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error generating response: %s", err.Error()), http.StatusInternalServerError)
		return
	}
	logger.Log.Info("CouchDB server is available.\n")

	fmt.Fprintf(w, string(response))
}

// StoreDBSyncCheckpoint - persists a checkpoint used during synchronization between pouch and
// couch DB. See https://pouchdb.com/guides/local-documents.html for more details on the concept
// of CouchDB local documents.
func (psh *PouchDBPluginServiceHandler) StoreDBSyncCheckpoint(w http.ResponseWriter, r *http.Request) {
	// TODO: Validate the request to ensure this operation is valid:

	dbName := getDBFieldFromRequest(r, dbNameInURL)
	logger.Log.Infof("Attempting to store %s to DB %s", db.DBSyncCheckpointStr, dbName)

	//Issue request to DAO Layer to store the DB Checkpoint
	queryParams := r.URL.Query()
	requestBody, err := getRequestBodyAsGenericObject(r)
	if err != nil {
		http.Error(w, fmt.Sprintf("Unable to read %s content: %s", db.DBSyncCheckpointStr, err.Error()), http.StatusBadRequest)
		return
	}

	result, err := psh.pouchPluginDB.StoreDBSyncCheckpoint(dbName, &queryParams, requestBody)
	if err != nil {
		http.Error(w, fmt.Sprintf("Unable to store %s: %s", db.DBSyncCheckpointStr, err.Error()), http.StatusInternalServerError)
		return
	}

	// Succesfully stored the DB Checkpoint, return the result.
	logger.Log.Infof("Successfully stored %s to DB %s\n", db.DBSyncCheckpointStr, dbName)
	response, err := json.Marshal(result)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error generating %s response: %s", db.DBSyncCheckpointStr, err.Error()), http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, string(response))
}

// GetDBSyncCheckpoint - retrieves a stored DB Checkpoint for use in pouch - couch synchronization.
// See https://pouchdb.com/guides/local-documents.html for more details on the concept
// of CouchDB local documents.
func (psh *PouchDBPluginServiceHandler) GetDBSyncCheckpoint(w http.ResponseWriter, r *http.Request) {
	// TODO: Validate the request to ensure this operation is valid:

	dbName := getDBFieldFromRequest(r, dbNameInURL)
	dbMethod := getDBFieldFromRequest(r, dbMethodInURL)
	docID := getDBFieldFromRequest(r, documentIDInURL)

	// Need to build up the full "_local/docID" format as URL parsing
	// separates this.
	documentID := dbMethod + "/" + docID

	logger.Log.Infof("Attempting to retrieve %s %s from DB %s", db.DBSyncCheckpointStr, documentID, dbName)

	//Issue request to DAO Layer to fetch the DB Checkpoint
	result, err := psh.pouchPluginDB.GetDBSyncCheckpoint(dbName, documentID)
	if err != nil {
		http.Error(w, fmt.Sprintf("Unable to retrieve %s: %s", db.DBSyncCheckpointStr, err.Error()), http.StatusInternalServerError)
		return
	}

	// Succesfully retrieved the DB Checkpoint, return the result.
	logger.Log.Infof("Successfully retrieved %s %s from DB %s\n", db.DBSyncCheckpointStr, documentID, dbName)
	response, err := json.Marshal(result)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error generating %s response: %s", db.DBSyncCheckpointStr, err.Error()), http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, string(response))
}

// GetDBRevisionDiff - provides ability to query the DB, with a list of revision tags map to
// a documentID, and have the DB respond with a list of which revisions it does not have.
// See http://docs.couchdb.org/en/2.1.1/api/database/misc.html#db-revs-diff for Couch documentation
// on the API.
func (psh *PouchDBPluginServiceHandler) GetDBRevisionDiff(w http.ResponseWriter, r *http.Request) {
	// TODO: Validate the request to ensure this operation is valid:

	dbName := getDBFieldFromRequest(r, dbNameInURL)

	logger.Log.Infof("Attempting to retrieve %s from DB %s", db.DBRevDiffStr, dbName)

	requestBody, err := getRequestBodyAsGenericObject(r)
	if err != nil {
		http.Error(w, fmt.Sprintf("Unable to read %s content: %s", db.DBRevDiffStr, err.Error()), http.StatusBadRequest)
		return
	}

	//Issue request to DAO Layer to fetch the Revision Diff
	result, err := psh.pouchPluginDB.GetDBRevisionDiff(dbName, requestBody)
	if err != nil {
		http.Error(w, fmt.Sprintf("Unable to retrieve %s: %s", db.DBRevDiffStr, err.Error()), http.StatusBadRequest)
		return
	}

	// Succesfully retrieved the DB Revision Diff, return the result.
	logger.Log.Infof("Successfully retrieved %s from DB %s\n", db.DBRevDiffStr, dbName)
	response, err := json.Marshal(result)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error generating %s response: %s", db.DBRevDiffStr, err.Error()), http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, string(response))
}

// BulkDBUpdate - allows multiple DB changes in one operation. See
// http://docs.couchdb.org/en/2.1.1/api/database/bulk-api.html#db-bulk-docs for
// CouchDB documentation of the API.
func (psh *PouchDBPluginServiceHandler) BulkDBUpdate(w http.ResponseWriter, r *http.Request) {
	// TODO: Validate the request to ensure this operation is valid:

	dbName := getDBFieldFromRequest(r, dbNameInURL)

	logger.Log.Infof("Attempting to perform %s on DB %s", db.DBBulkUpdateStr, dbName)

	requestBody, err := getRequestBodyAsGenericObject(r)
	if err != nil {
		http.Error(w, fmt.Sprintf("Unable to read %s content: %s", db.DBBulkUpdateStr, err.Error()), http.StatusBadRequest)
		return
	}

	//Issue request to DAO Layer to perform the bulk update
	result, err := psh.pouchPluginDB.BulkDBUpdate(dbName, requestBody)
	if err != nil {
		http.Error(w, fmt.Sprintf("Unable to complete %s: %s", db.DBBulkUpdateStr, err.Error()), http.StatusBadRequest)
		return
	}

	// Succesfully performed the bulk update, return the result.
	logger.Log.Infof("Successfully completed %s from DB %s\n", db.DBBulkUpdateStr, dbName)
	response, err := json.Marshal(result)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error generating %s response: %s", db.DBBulkUpdateStr, err.Error()), http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, string(response))
}

func getDBFieldFromRequest(r *http.Request, field pouchPluginURLPart) string {
	urlParts := strings.Split(r.URL.Path, "/")
	return urlParts[field]
}

func getRequestBodyAsGenericObject(r *http.Request) (map[string]interface{}, error) {
	decoder := json.NewDecoder(r.Body)
	var result map[string]interface{}
	err := decoder.Decode(&result)
	if err != nil {
		return nil, err
	}
	defer r.Body.Close()

	return result, nil
}

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
