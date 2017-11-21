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

type httpErrorString string

const (
	notFound httpErrorString = "status 404 - not found"
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

		server.Route{
			Name:        "CheckDBAvailability",
			Method:      "GET",
			Pattern:     "/{dbname}/",
			HandlerFunc: result.CheckDBAvailability,
		},

		server.Route{
			Name:        "CreateDB",
			Method:      "PUT",
			Pattern:     "/{dbname}/",
			HandlerFunc: result.CreateDB,
		},

		server.Route{
			Name:        "GetAllDBDocs",
			Method:      "POST",
			Pattern:     "/{dbname}/_all_docs",
			HandlerFunc: result.GetAllDBDocs,
		},

		server.Route{
			Name:        "GetDBDoc",
			Method:      "GET",
			Pattern:     "/{dbname}/{docid}",
			HandlerFunc: result.GetDBDoc,
		},

		server.Route{
			Name:        "BulkDBGet",
			Method:      "POST",
			Pattern:     "/{dbname}/_bulk_get",
			HandlerFunc: result.BulkDBGet,
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
	cfg := gather.GetConfig()
	dbType := gather.DBImpl(cfg.GetInt(gather.CK_args_pouchplugindb_impl.String()))
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

	dbName := getDBFieldFromRequest(r, 1)

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

	dbName := getDBFieldFromRequest(r, 1)
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

	dbName := getDBFieldFromRequest(r, 1)
	dbMethod := getDBFieldFromRequest(r, 2)
	docID := getDBFieldFromRequest(r, 3)

	// Need to build up the full "_local/docID" format as URL parsing
	// separates this.
	documentID := dbMethod + "/" + docID

	logger.Log.Infof("Attempting to retrieve %s %s from DB %s", db.DBSyncCheckpointStr, documentID, dbName)

	//Issue request to DAO Layer to fetch the DB Checkpoint
	result, err := psh.pouchPluginDB.GetDBSyncCheckpoint(dbName, documentID)
	if err != nil {
		if checkError(err, notFound) {
			http.Error(w, fmt.Sprintf("%s %s does not exist", db.DBSyncCheckpointStr, documentID), http.StatusNotFound)
			return
		}
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

	dbName := getDBFieldFromRequest(r, 1)

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

	dbName := getDBFieldFromRequest(r, 1)

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

// CheckDBAvailability - heartbeat for the given database.
func (psh *PouchDBPluginServiceHandler) CheckDBAvailability(w http.ResponseWriter, r *http.Request) {
	// TODO: Validate the request to ensure this operation is valid:

	dbName := getDBFieldFromRequest(r, 1)
	logger.Log.Infof("Checking for availability of DB %s", dbName)

	//Issue request to DAO Layer to access check availability
	result, err := psh.pouchPluginDB.CheckDBAvailability(dbName)
	if err != nil {
		if checkError(err, notFound) {
			http.Error(w, fmt.Sprintf("DB %s does not exist", dbName), http.StatusNotFound)
			return
		}
		http.Error(w, fmt.Sprintf("Error checking availability of DB %s: %s", dbName, err.Error()), http.StatusInternalServerError)
		return
	}

	// Succesfully accessed the couch server, return the result
	response, err := json.Marshal(result)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error generating response: %s", err.Error()), http.StatusInternalServerError)
		return
	}
	logger.Log.Infof("DB %s is available.\n", dbName)

	fmt.Fprintf(w, string(response))
}

// GetAllDBDocs - provides metadata on all docs in a DB. See
// http://docs.couchdb.org/en/2.1.1/api/database/bulk-api.html for
// Couch documentation of this API.
func (psh *PouchDBPluginServiceHandler) GetAllDBDocs(w http.ResponseWriter, r *http.Request) {
	// TODO: Validate the request to ensure this operation is valid:

	dbName := getDBFieldFromRequest(r, 1)

	logger.Log.Infof("Attempting to fetch %s from DB %s", db.DBAllDocsStr, dbName)

	requestBody, err := getRequestBodyAsGenericObject(r)
	if err != nil {
		http.Error(w, fmt.Sprintf("Unable to read %s request content: %s", db.DBAllDocsStr, err.Error()), http.StatusBadRequest)
		return
	}

	//Issue request to DAO Layer to perform the bulk fetch
	result, err := psh.pouchPluginDB.GetAllDBDocs(dbName, requestBody)
	if err != nil {
		http.Error(w, fmt.Sprintf("Unable to fetch %s: %s", db.DBAllDocsStr, err.Error()), http.StatusBadRequest)
		return
	}

	// Succesfully performed the bulk fetch, return the result.
	logger.Log.Infof("Successfully retrieved %s from DB %s\n", db.DBAllDocsStr, dbName)
	response, err := json.Marshal(result)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error generating %s response: %s", db.DBAllDocsStr, err.Error()), http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, string(response))
}

// CreateDB - provides the ability for pouch to create a couchDB.
func (psh *PouchDBPluginServiceHandler) CreateDB(w http.ResponseWriter, r *http.Request) {
	// TODO: Validate the request to ensure this operation is valid:

	dbName := getDBFieldFromRequest(r, 1)

	logger.Log.Infof("Attempting to create DB %s", db.DBAllDocsStr, dbName)

	//Issue request to DAO Layer to perform the DB creation
	result, err := psh.pouchPluginDB.CreateDB(dbName)
	if err != nil {
		http.Error(w, fmt.Sprintf("Unable to create DB %s: %s", dbName, err.Error()), http.StatusBadRequest)
		return
	}

	// Succesfully performed the DB creation, return the result.
	logger.Log.Infof("Successfully created DB %s\n", dbName)
	response, err := json.Marshal(result)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error generating DB creation response: %s", err.Error()), http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, string(response))
}

// GetDBDoc - returns a document plus optional metadate about the document from CouchDB.
// See http://docs.couchdb.org/en/2.1.1/api/document/common.html for documentation of the API
func (psh *PouchDBPluginServiceHandler) GetDBDoc(w http.ResponseWriter, r *http.Request) {
	// TODO: Validate the request to ensure this operation is valid:

	dbName := getDBFieldFromRequest(r, 1)
	docID := getDBFieldFromRequest(r, 2)

	logger.Log.Infof("Fetching %s %s from DB %s", db.DBDocStr, docID, dbName)

	//Issue request to DAO Layer to access the Document
	queryParams := r.URL.Query()
	result, err := psh.pouchPluginDB.GetDoc(dbName, docID, &queryParams, &r.Header)
	if err != nil {
		http.Error(w, fmt.Sprintf("Unable to retrieve %s %s: %s", db.DBDocStr, docID, err.Error()), http.StatusBadRequest)
		return
	}

	// Succesfully fetched the Document, return the result. See
	logger.Log.Infof("Successfully accessed %s %s from DB %s\n", db.DBDocStr, docID, dbName)
	response, err := json.Marshal(result["data"]) // Only need the data portion of this wrapper object
	if err != nil {
		http.Error(w, fmt.Sprintf("Error generating %s response: %s", db.DBDocStr, err.Error()), http.StatusBadRequest)
		return
	}

	fmt.Fprintf(w, string(response))
}

// BulkDBGet - allows fetching multiple DB Documenta in one operation.
// There is no CouchDB documentation of the API.
func (psh *PouchDBPluginServiceHandler) BulkDBGet(w http.ResponseWriter, r *http.Request) {
	// TODO: Validate the request to ensure this operation is valid:

	dbName := getDBFieldFromRequest(r, 1)

	logger.Log.Infof("Attempting to perform %s on DB %s", db.DBBulkGetStr, dbName)

	queryParams := r.URL.Query()
	requestBody, err := getRequestBodyAsGenericObject(r)
	if err != nil {
		http.Error(w, fmt.Sprintf("Unable to read %s content: %s", db.DBBulkGetStr, err.Error()), http.StatusBadRequest)
		return
	}

	//Issue request to DAO Layer to perform the bulk update
	result, err := psh.pouchPluginDB.BulkDBGet(dbName, &queryParams, requestBody)
	if err != nil {
		http.Error(w, fmt.Sprintf("Unable to complete %s: %s", db.DBBulkGetStr, err.Error()), http.StatusBadRequest)
		return
	}

	// Succesfully performed the bulk get, return the result.
	logger.Log.Infof("Successfully completed %s from DB %s\n", db.DBBulkGetStr, dbName)
	response, err := json.Marshal(result)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error generating %s response: %s", db.DBBulkGetStr, err.Error()), http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, string(response))
}

func getDBFieldFromRequest(r *http.Request, urlPart int32) string {
	urlParts := strings.Split(r.URL.Path, "/")
	return urlParts[urlPart]
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

func checkError(err error, errorType httpErrorString) bool {
	if strings.Contains(err.Error(), string(errorType)) {
		return true
	}

	return false
}
