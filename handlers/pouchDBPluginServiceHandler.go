package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

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

	pouchStr           = "pouch"
	getChangesStr      = pouchStr + "_changes"
	serverHBStr        = pouchStr + "_available"
	dbHBStr            = pouchStr + "_db_available"
	getCheckpointStr   = pouchStr + "_local_get"
	storeCheckpointStr = pouchStr + "_local_put"
	dbDiffStr          = pouchStr + "_diff"
	createDbStr        = pouchStr + "_db_put"
	bulkUpdateStr      = pouchStr + "_bulk_docs_put"
	bulkGetStr         = pouchStr + "_bulk_docs_get"
	allDocsStr         = pouchStr + "_all_docs"
	createDocStr       = pouchStr + "_db_doc_put"
	getDBDocStr        = pouchStr + "_db_doc_get"

	heartbeatStr = "heartbeat"
	heartbeatMax = 25000
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
			Pattern:     "/pouchdb/",
			HandlerFunc: result.CheckAvailability,
		},

		server.Route{
			Name:        "GetChanges",
			Method:      "GET",
			Pattern:     "/pouchdb/{dbname}/_changes",
			HandlerFunc: result.GetChanges,
		},

		server.Route{
			Name:        "StoreDBSyncCheckpoint",
			Method:      "PUT",
			Pattern:     "/pouchdb/{dbname}/_local/{docid}",
			HandlerFunc: result.StoreDBSyncCheckpoint,
		},

		server.Route{
			Name:        "GetDBSyncCheckpoint",
			Method:      "GET",
			Pattern:     "/pouchdb/{dbname}/_local/{docid}",
			HandlerFunc: result.GetDBSyncCheckpoint,
		},

		server.Route{
			Name:        "GetDBRevisionDiff",
			Method:      "POST",
			Pattern:     "/pouchdb/{dbname}/_revs_diff",
			HandlerFunc: result.GetDBRevisionDiff,
		},

		server.Route{
			Name:        "BulkDBUpdate",
			Method:      "POST",
			Pattern:     "/pouchdb/{dbname}/_bulk_docs",
			HandlerFunc: result.BulkDBUpdate,
		},

		server.Route{
			Name:        "CheckDBAvailability",
			Method:      "GET",
			Pattern:     "/pouchdb/{dbname}/",
			HandlerFunc: result.CheckDBAvailability,
		},

		server.Route{
			Name:        "CreateDB",
			Method:      "PUT",
			Pattern:     "/pouchdb/{dbname}/",
			HandlerFunc: result.CreateDB,
		},

		server.Route{
			Name:        "GetAllDBDocs",
			Method:      "POST",
			Pattern:     "/pouchdb/{dbname}/_all_docs",
			HandlerFunc: result.GetAllDBDocs,
		},

		server.Route{
			Name:        "GetDBDoc",
			Method:      "GET",
			Pattern:     "/pouchdb/{dbname}/{docid}",
			HandlerFunc: result.GetDBDoc,
		},

		server.Route{
			Name:        "BulkDBGet",
			Method:      "POST",
			Pattern:     "/pouchdb/{dbname}/_bulk_get",
			HandlerFunc: result.BulkDBGet,
		},
	}

	return result
}

// RegisterAPIHandlers - will bind any REST API routes defined in this service
// to the passed in request multiplexor.
func (psh *PouchDBPluginServiceHandler) RegisterAPIHandlers(router *mux.Router) {
	for _, route := range psh.routes {
		logger.Log.Debugf("Registering endpoint: %v", route)
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
	startTime := time.Now()
	defer r.Body.Close()

	// Cast an object to be used to send periodic data to the client
	flusher, ok := w.(http.Flusher)
	if !ok {
		logger.Log.Debug("Error casting writer to flusher")
	}

	w.Header().Add("Cache-Control", "must-revalidate")
	w.Header().Add("Content-Type", "application/json")
	w.Header().Add("X-Content-Type-Options", "nosniff")

	dbName := getDBFieldFromRequest(r, 2)

	logger.Log.Infof("Looking for changes from DB %s", dbName)

	//Issue request to DAO Layer to access the Changes Feed
	queryParams := r.URL.Query()

	heartbeat := queryParams.Get(heartbeatStr)
	stopHeartbeat := false
	if heartbeat != "" {
		// This is a longpoll, need to handel the open connection
		heartbeatInterval, err := strconv.Atoi(heartbeat)
		if err != nil {
			msg := fmt.Sprintf("Could not establish hearbeat for %s response: %s", db.ChangeFeedStr, err.Error())
			reportError(w, startTime, "400", getChangesStr, msg, http.StatusBadRequest)
			return
		}
		if heartbeatInterval > heartbeatMax {
			heartbeatInterval = heartbeatMax
		}

		// This request needs to intermittently return blank lines to keep the connection open.
		// Start a go routine to handle that:
		go func() {
			// Setup an object which will notify us when the client closes the connection to this request.
			cn, ok := w.(http.CloseNotifier)
			if !ok {
				msg := fmt.Sprintf("Unable to setup streaming response %s: %s", db.ChangeFeedStr, err.Error())
				reportError(w, startTime, "500", getChangesStr, msg, http.StatusInternalServerError)
				return
			}

			// Run the loop which will keep the connection with the client alive
			for {
				if stopHeartbeat {
					break
				}

				select {
				case <-cn.CloseNotify():
					logger.Log.Debug("Client disconnected from changes stream")
					return
				default:
					logger.Log.Debugf("Changes Heartbeat: %v", time.Now())
					fmt.Fprintf(w, "\n")
					flusher.Flush()
					time.Sleep(time.Duration(time.Millisecond * time.Duration(heartbeatInterval)))
				}
			}
		}()
	}

	result, err := psh.pouchPluginDB.GetChanges(dbName, &queryParams)
	stopHeartbeat = true
	if err != nil {
		msg := fmt.Sprintf("Unable to retrieve %s: %s", db.ChangeFeedStr, err.Error())
		reportError(w, startTime, "400", getChangesStr, msg, http.StatusBadRequest)
		return
	}

	// Succesfully fetched the Changes Feed, return the result. See
	logger.Log.Infof("Successfully accessed %s changes from DB %s\n", db.ChangeFeedStr, dbName)

	response, err := json.Marshal(result)
	if err != nil {
		msg := fmt.Sprintf("Error generating %s response: %s", db.ChangeFeedStr, err.Error())
		reportError(w, startTime, "400", getChangesStr, msg, http.StatusBadRequest)
		return
	}

	trackAPIMetrics(startTime, "200", getChangesStr)
	fmt.Fprintf(w, string(response))
}

// CheckAvailability - used to check if the CouchDB server is available.
// See http://docs.couchdb.org/en/2.1.1/api/server/common.html for the
// CouchDB documentation on this API.
func (psh *PouchDBPluginServiceHandler) CheckAvailability(w http.ResponseWriter, r *http.Request) {
	startTime := time.Now()

	logger.Log.Info("Checking for CouchDB availability")

	//Issue request to DAO Layer to access check availability
	result, err := psh.pouchPluginDB.CheckAvailability()
	if err != nil {
		msg := fmt.Sprintf("Error checking CouchDB availability: %s", err.Error())
		reportError(w, startTime, "500", serverHBStr, msg, http.StatusInternalServerError)
		return
	}

	// Succesfully accessed the couch server, return the result
	response, err := json.Marshal(result)
	if err != nil {
		msg := fmt.Sprintf("Error generating response: %s", err.Error())
		reportError(w, startTime, "500", serverHBStr, msg, http.StatusInternalServerError)
		return
	}
	logger.Log.Info("CouchDB server is available.\n")

	trackAPIMetrics(startTime, "200", serverHBStr)
	fmt.Fprintf(w, string(response))
}

// StoreDBSyncCheckpoint - persists a checkpoint used during synchronization between pouch and
// couch DB. See https://pouchdb.com/guides/local-documents.html for more details on the concept
// of CouchDB local documents.
func (psh *PouchDBPluginServiceHandler) StoreDBSyncCheckpoint(w http.ResponseWriter, r *http.Request) {
	startTime := time.Now()

	dbName := getDBFieldFromRequest(r, 2)
	logger.Log.Infof("Attempting to store %s to DB %s", db.DBSyncCheckpointStr, dbName)

	//Issue request to DAO Layer to store the DB Checkpoint
	queryParams := r.URL.Query()
	requestBody, err := getRequestBodyAsGenericObject(r)
	if err != nil {
		msg := fmt.Sprintf("Unable to read %s content: %s", db.DBSyncCheckpointStr, err.Error())
		reportError(w, startTime, "500", storeCheckpointStr, msg, http.StatusInternalServerError)
		return
	}

	result, err := psh.pouchPluginDB.StoreDBSyncCheckpoint(dbName, &queryParams, requestBody)
	if err != nil {
		msg := fmt.Sprintf("Unable to store %s: %s", db.DBSyncCheckpointStr, err.Error())
		reportError(w, startTime, "500", storeCheckpointStr, msg, http.StatusInternalServerError)
		return
	}

	// Succesfully stored the DB Checkpoint, return the result.
	logger.Log.Infof("Successfully stored %s to DB %s\n", db.DBSyncCheckpointStr, dbName)
	response, err := json.Marshal(result)
	if err != nil {
		msg := fmt.Sprintf("Error generating %s response: %s", db.DBSyncCheckpointStr, err.Error())
		reportError(w, startTime, "500", storeCheckpointStr, msg, http.StatusInternalServerError)
		return
	}

	trackAPIMetrics(startTime, "200", storeCheckpointStr)
	fmt.Fprintf(w, string(response))
}

// GetDBSyncCheckpoint - retrieves a stored DB Checkpoint for use in pouch - couch synchronization.
// See https://pouchdb.com/guides/local-documents.html for more details on the concept
// of CouchDB local documents.
func (psh *PouchDBPluginServiceHandler) GetDBSyncCheckpoint(w http.ResponseWriter, r *http.Request) {
	startTime := time.Now()

	dbName := getDBFieldFromRequest(r, 2)
	dbMethod := getDBFieldFromRequest(r, 3)
	docID := getDBFieldFromRequest(r, 4)

	// Need to build up the full "_local/docID" format as URL parsing
	// separates this.
	documentID := dbMethod + "/" + docID

	logger.Log.Infof("Attempting to retrieve %s %s from DB %s", db.DBSyncCheckpointStr, documentID, dbName)

	//Issue request to DAO Layer to fetch the DB Checkpoint
	result, err := psh.pouchPluginDB.GetDBSyncCheckpoint(dbName, documentID)
	if err != nil {
		if checkError(err, notFound) {
			msg := fmt.Sprintf("%s %s does not exist", db.DBSyncCheckpointStr, documentID)
			reportError(w, startTime, "404", getCheckpointStr, msg, http.StatusNotFound)
			return
		}
		msg := fmt.Sprintf("Unable to retrieve %s: %s", db.DBSyncCheckpointStr, err.Error())
		reportError(w, startTime, "500", getCheckpointStr, msg, http.StatusInternalServerError)
		return
	}

	// Succesfully retrieved the DB Checkpoint, return the result.
	logger.Log.Infof("Successfully retrieved %s %s from DB %s\n", db.DBSyncCheckpointStr, documentID, dbName)
	response, err := json.Marshal(result)
	if err != nil {
		msg := fmt.Sprintf("Error generating %s response: %s", db.DBSyncCheckpointStr, err.Error())
		reportError(w, startTime, "500", getCheckpointStr, msg, http.StatusInternalServerError)
		return
	}

	trackAPIMetrics(startTime, "200", getCheckpointStr)
	fmt.Fprintf(w, string(response))
}

// GetDBRevisionDiff - provides ability to query the DB, with a list of revision tags map to
// a documentID, and have the DB respond with a list of which revisions it does not have.
// See http://docs.couchdb.org/en/2.1.1/api/database/misc.html#db-revs-diff for Couch documentation
// on the API.
func (psh *PouchDBPluginServiceHandler) GetDBRevisionDiff(w http.ResponseWriter, r *http.Request) {
	startTime := time.Now()

	dbName := getDBFieldFromRequest(r, 2)

	logger.Log.Infof("Attempting to retrieve %s from DB %s", db.DBRevDiffStr, dbName)

	requestBody, err := getRequestBodyAsGenericObject(r)
	if err != nil {
		msg := fmt.Sprintf("Unable to read %s content: %s", db.DBRevDiffStr, err.Error())
		reportError(w, startTime, "400", dbDiffStr, msg, http.StatusBadRequest)
		return
	}

	//Issue request to DAO Layer to fetch the Revision Diff
	result, err := psh.pouchPluginDB.GetDBRevisionDiff(dbName, requestBody)
	if err != nil {
		msg := fmt.Sprintf("Unable to retrieve %s: %s", db.DBRevDiffStr, err.Error())
		reportError(w, startTime, "400", dbDiffStr, msg, http.StatusBadRequest)
		return
	}

	// Succesfully retrieved the DB Revision Diff, return the result.
	logger.Log.Infof("Successfully retrieved %s from DB %s\n", db.DBRevDiffStr, dbName)
	response, err := json.Marshal(result)
	if err != nil {
		msg := fmt.Sprintf("Error generating %s response: %s", db.DBRevDiffStr, err.Error())
		reportError(w, startTime, "500", dbDiffStr, msg, http.StatusInternalServerError)
		return
	}

	trackAPIMetrics(startTime, "200", dbDiffStr)
	fmt.Fprintf(w, string(response))
}

// BulkDBUpdate - allows multiple DB changes in one operation. See
// http://docs.couchdb.org/en/2.1.1/api/database/bulk-api.html#db-bulk-docs for
// CouchDB documentation of the API.
func (psh *PouchDBPluginServiceHandler) BulkDBUpdate(w http.ResponseWriter, r *http.Request) {
	startTime := time.Now()

	dbName := getDBFieldFromRequest(r, 2)

	logger.Log.Infof("Attempting to perform %s on DB %s", db.DBBulkUpdateStr, dbName)

	requestBody, err := getRequestBodyAsGenericObject(r)
	if err != nil {
		msg := fmt.Sprintf("Unable to read %s content: %s", db.DBBulkUpdateStr, err.Error())
		reportError(w, startTime, "400", bulkUpdateStr, msg, http.StatusBadRequest)
		return
	}

	//Issue request to DAO Layer to perform the bulk update
	result, err := psh.pouchPluginDB.BulkDBUpdate(dbName, requestBody)
	if err != nil {
		msg := fmt.Sprintf("Unable to complete %s: %s", db.DBBulkUpdateStr, err.Error())
		reportError(w, startTime, "400", bulkUpdateStr, msg, http.StatusBadRequest)
		return
	}

	// Succesfully performed the bulk update, return the result.
	logger.Log.Infof("Successfully completed %s from DB %s\n", db.DBBulkUpdateStr, dbName)
	response, err := json.Marshal(result)
	if err != nil {
		msg := fmt.Sprintf("Error generating %s response: %s", db.DBBulkUpdateStr, err.Error())
		reportError(w, startTime, "500", bulkUpdateStr, msg, http.StatusInternalServerError)
		return
	}

	trackAPIMetrics(startTime, "200", bulkUpdateStr)
	fmt.Fprintf(w, string(response))
}

// CheckDBAvailability - heartbeat for the given database.
func (psh *PouchDBPluginServiceHandler) CheckDBAvailability(w http.ResponseWriter, r *http.Request) {
	startTime := time.Now()

	dbName := getDBFieldFromRequest(r, 2)
	logger.Log.Infof("Checking for availability of DB %s", dbName)

	//Issue request to DAO Layer to access check availability
	result, err := psh.IsDBAvailable(dbName)
	if err != nil {
		if checkError(err, notFound) {
			msg := fmt.Sprintf("DB %s does not exist", dbName)
			reportError(w, startTime, "404", dbHBStr, msg, http.StatusNotFound)
			return
		}
		msg := fmt.Sprintf("Error checking availability of DB %s: %s", dbName, err.Error())
		reportError(w, startTime, "500", dbHBStr, msg, http.StatusInternalServerError)
		return
	}

	// Succesfully accessed the couch server, return the result
	response, err := json.Marshal(result)
	if err != nil {
		msg := fmt.Sprintf("Error generating response: %s", err.Error())
		reportError(w, startTime, "500", dbHBStr, msg, http.StatusInternalServerError)
		return
	}
	logger.Log.Infof("DB %s is available.\n", dbName)

	trackAPIMetrics(startTime, "200", dbHBStr)
	fmt.Fprintf(w, string(response))
}

// GetAllDBDocs - provides metadata on all docs in a DB. See
// http://docs.couchdb.org/en/2.1.1/api/database/bulk-api.html for
// Couch documentation of this API.
func (psh *PouchDBPluginServiceHandler) GetAllDBDocs(w http.ResponseWriter, r *http.Request) {
	startTime := time.Now()

	dbName := getDBFieldFromRequest(r, 2)

	logger.Log.Infof("Attempting to fetch %s from DB %s", db.DBAllDocsStr, dbName)

	requestBody, err := getRequestBodyAsGenericObject(r)
	if err != nil {
		msg := fmt.Sprintf("Unable to read %s request content: %s", db.DBAllDocsStr, err.Error())
		reportError(w, startTime, "400", allDocsStr, msg, http.StatusBadRequest)
		return
	}

	//Issue request to DAO Layer to perform the bulk fetch
	result, err := psh.pouchPluginDB.GetAllDBDocs(dbName, requestBody)
	if err != nil {
		msg := fmt.Sprintf("Unable to fetch %s: %s", db.DBAllDocsStr, err.Error())
		reportError(w, startTime, "400", allDocsStr, msg, http.StatusBadRequest)
		return
	}

	// Succesfully performed the bulk fetch, return the result.
	logger.Log.Infof("Successfully retrieved %s from DB %s\n", db.DBAllDocsStr, dbName)
	response, err := json.Marshal(result)
	if err != nil {
		msg := fmt.Sprintf("Error generating %s response: %s", db.DBAllDocsStr, err.Error())
		reportError(w, startTime, "500", allDocsStr, msg, http.StatusInternalServerError)
		return
	}

	trackAPIMetrics(startTime, "200", allDocsStr)
	fmt.Fprintf(w, string(response))
}

// CreateDB - provides the ability for pouch to create a couchDB.
func (psh *PouchDBPluginServiceHandler) CreateDB(w http.ResponseWriter, r *http.Request) {
	startTime := time.Now()

	dbName := getDBFieldFromRequest(r, 2)

	logger.Log.Infof("Attempting to create DB %s", db.DBAllDocsStr, dbName)

	//Issue request to DAO Layer to perform the DB creation
	result, err := psh.AddDB(dbName)
	if err != nil {
		msg := fmt.Sprintf("Unable to create DB %s: %s", dbName, err.Error())
		reportError(w, startTime, "400", createDbStr, msg, http.StatusBadRequest)
		return
	}

	// Succesfully performed the DB creation, return the result.
	logger.Log.Infof("Successfully created DB %s\n", dbName)
	response, err := json.Marshal(result)
	if err != nil {
		msg := fmt.Sprintf("Error generating DB creation response: %s", err.Error())
		reportError(w, startTime, "500", createDbStr, msg, http.StatusInternalServerError)
		return
	}

	trackAPIMetrics(startTime, "200", createDbStr)
	fmt.Fprintf(w, string(response))
}

// GetDBDoc - returns a document plus optional metadate about the document from CouchDB.
// See http://docs.couchdb.org/en/2.1.1/api/document/common.html for documentation of the API
func (psh *PouchDBPluginServiceHandler) GetDBDoc(w http.ResponseWriter, r *http.Request) {
	startTime := time.Now()

	dbName := getDBFieldFromRequest(r, 2)
	docID := getDBFieldFromRequest(r, 3)

	logger.Log.Infof("Fetching %s %s from DB %s", db.DBDocStr, docID, dbName)

	//Issue request to DAO Layer to access the Document
	queryParams := r.URL.Query()
	result, err := psh.pouchPluginDB.GetDoc(dbName, docID, &queryParams, &r.Header)
	if err != nil {
		msg := fmt.Sprintf("Unable to retrieve %s %s: %s", db.DBDocStr, docID, err.Error())
		reportError(w, startTime, "400", getDBDocStr, msg, http.StatusBadRequest)
		return
	}

	// Succesfully fetched the Document, return the result. See
	logger.Log.Infof("Successfully accessed %s %s from DB %s\n", db.DBDocStr, docID, dbName)
	response, err := json.Marshal(result["data"]) // Only need the data portion of this wrapper object
	if err != nil {
		msg := fmt.Sprintf("Error generating %s response: %s", db.DBDocStr, err.Error())
		reportError(w, startTime, "500", getDBDocStr, msg, http.StatusInternalServerError)
		return
	}

	trackAPIMetrics(startTime, "200", getDBDocStr)
	fmt.Fprintf(w, string(response))
}

// BulkDBGet - allows fetching multiple DB Documenta in one operation.
// There is no CouchDB documentation of the API.
func (psh *PouchDBPluginServiceHandler) BulkDBGet(w http.ResponseWriter, r *http.Request) {
	startTime := time.Now()

	dbName := getDBFieldFromRequest(r, 2)

	logger.Log.Infof("Attempting to perform %s on DB %s", db.DBBulkGetStr, dbName)

	queryParams := r.URL.Query()
	requestBody, err := getRequestBodyAsGenericObject(r)
	if err != nil {
		msg := fmt.Sprintf("Unable to read %s content: %s", db.DBBulkGetStr, err.Error())
		reportError(w, startTime, "400", bulkGetStr, msg, http.StatusBadRequest)
		return
	}

	//Issue request to DAO Layer to perform the bulk update
	result, err := psh.pouchPluginDB.BulkDBGet(dbName, &queryParams, requestBody)
	if err != nil {
		msg := fmt.Sprintf("Unable to complete %s: %s", db.DBBulkGetStr, err.Error())
		reportError(w, startTime, "400", bulkGetStr, msg, http.StatusBadRequest)
		return
	}

	// Succesfully performed the bulk get, return the result.
	logger.Log.Infof("Successfully completed %s from DB %s\n", db.DBBulkGetStr, dbName)
	response, err := json.Marshal(result)
	if err != nil {
		msg := fmt.Sprintf("Error generating %s response: %s", db.DBBulkGetStr, err.Error())
		reportError(w, startTime, "500", bulkGetStr, msg, http.StatusInternalServerError)
		return
	}

	trackAPIMetrics(startTime, "200", bulkGetStr)
	fmt.Fprintf(w, string(response))
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

// IsDBAvailable - checks if a DB is available.
func (psh *PouchDBPluginServiceHandler) IsDBAvailable(dbName string) (map[string]interface{}, error) {
	return psh.pouchPluginDB.CheckDBAvailability(dbName)
}

// AddDB - creates a DB instance.
func (psh *PouchDBPluginServiceHandler) AddDB(dbName string) (map[string]interface{}, error) {
	return psh.pouchPluginDB.CreateDB(dbName)
}
