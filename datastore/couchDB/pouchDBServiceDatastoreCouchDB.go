package couchDB

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"

	"github.com/accedian/adh-gather/config"
	ds "github.com/accedian/adh-gather/datastore"
	"github.com/accedian/adh-gather/gather"
	"github.com/accedian/adh-gather/logger"
	couchdb "github.com/leesper/couchdb-golang"
)

// PouchDBServiceDatastoreCouchDB - struct responsible for handling
// database operations for the PouchDB Plugin Service when using CouchDB
// as the storage option.
type PouchDBServiceDatastoreCouchDB struct {
	couchHost string
	cfg       config.Provider
}

// CreatePouchDBServiceDAO - instantiates a CouchDB implementation of the
// PouchDBPluginServiceDatastore.
func CreatePouchDBServiceDAO() (*PouchDBServiceDatastoreCouchDB, error) {
	result := new(PouchDBServiceDatastoreCouchDB)
	result.cfg = gather.GetConfig()

	provDBURL := fmt.Sprintf("%s:%d",
		result.cfg.GetString(gather.CK_server_datastore_ip.String()),
		result.cfg.GetInt(gather.CK_server_datastore_port.String()))
	logger.Log.Debugf("Pouch Plugin Service CouchDB URL is: %s", provDBURL)
	result.couchHost = provDBURL

	return result, nil
}

// GetChanges - CouchDB implementation of GetChanges
func (psd *PouchDBServiceDatastoreCouchDB) GetChanges(dbname string, queryParams *url.Values) (map[string]interface{}, error) {
	logger.Log.Debugf("Attempting to access %s for DB %s with options: %v", ds.ChangeFeedStr, dbname, queryParams)

	db, err := getDatabase(createDBPathStr(psd.couchHost, dbname))
	if err != nil {
		return nil, err
	}

	// Retrieve the DB Changes Feed data from CouchDB
	fetchedData, err := accessDBChangesFeed(db, queryParams)
	if err != nil {
		return nil, err
	}

	logger.Log.Debugf("%s for DB %s returned: %v\n", ds.ChangeFeedStr, dbname, fetchedData)
	return fetchedData, nil
}

// CheckAvailability - CouchDB implementation of CheckAvailability
func (psd *PouchDBServiceDatastoreCouchDB) CheckAvailability() (map[string]interface{}, error) {

	resource, err := couchdb.NewResource(psd.couchHost, nil)
	if err != nil {
		logger.Log.Debugf("Falied to check availability of Couch Server %s: %s", ds.DBRevDiffStr, err.Error())
		return nil, err
	}

	// Retrieve the DB Changes Feed data from CouchDB
	fetchedData, err := checkIfAvailable(resource)
	if err != nil {
		return nil, err
	}

	logger.Log.Debugf("CheckAvailibility complete: %v\n", fetchedData)
	return fetchedData, nil
}

// StoreDBSyncCheckpoint - CouchDB implementation of StoreDBSyncCheckpoint
func (psd *PouchDBServiceDatastoreCouchDB) StoreDBSyncCheckpoint(dbname string, queryParams *url.Values, request map[string]interface{}) (map[string]interface{}, error) {
	// Validate the request to ensure this operation is valid:

	logger.Log.Debugf("Storing %s: %v - using options: %v", ds.DBSyncCheckpointStr, request, queryParams)

	db, err := getDatabase(createDBPathStr(psd.couchHost, dbname))
	if err != nil {
		return nil, err
	}

	// Store the sync checkpoint in CouchDB
	id, rev, err := storeDataInCouchDBWithQueryParams(request, ds.DBSyncCheckpointStr, db, queryParams)
	if err != nil {
		return nil, err
	}

	result := make(map[string]interface{})
	result["id"] = id
	result["ok"] = true
	result["rev"] = rev

	// DB Sync Checkpoint stored, send the response
	logger.Log.Debugf("%s %v stored.\n", ds.DBSyncCheckpointStr, result)
	return result, nil
}

// GetDBSyncCheckpoint - CoiuchDB implementation of GetDBSyncCheckpoint
func (psd *PouchDBServiceDatastoreCouchDB) GetDBSyncCheckpoint(dbname string, documentID string) (map[string]interface{}, error) {
	logger.Log.Debugf("Retrieving %s: %s", ds.DBSyncCheckpointStr, documentID)

	db, err := getDatabase(createDBPathStr(psd.couchHost, dbname))
	if err != nil {
		return nil, err
	}

	// Retrieve the checkpoint data from CouchDB
	fetchedData, err := getByDocID(documentID, ds.DBSyncCheckpointStr, db)
	if err != nil {
		return nil, err
	}

	// DB Sync Checkpoint retrieved, send the response
	logger.Log.Debugf("%s %v retrieved.\n", ds.DBSyncCheckpointStr, fetchedData)
	return fetchedData, nil
}

// GetDBRevisionDiff - CouchDB implementation of GetDBRevisionDiff
func (psd *PouchDBServiceDatastoreCouchDB) GetDBRevisionDiff(dbname string, request map[string]interface{}) (map[string]interface{}, error) {
	logger.Log.Debugf("Retrieving %s %v from DB %s", ds.DBRevDiffStr, request, dbname)

	// Create a resource that can make the revs diff call to Couch
	resource, err := couchdb.NewResource(createDBPathStr(psd.couchHost, dbname), nil)
	if err != nil {
		logger.Log.Debugf("Falied to retrieve %s: %s", ds.DBRevDiffStr, err.Error())
		return nil, err
	}

	// Retrieve the checkpoint data from CouchDB
	fetchedData, err := fetchRevDiff(request, resource)
	if err != nil {
		return nil, err
	}

	// DB Sync Checkpoint retrieved, send the response
	logger.Log.Debugf("Retrieved %s %v from DB %s\n", ds.DBRevDiffStr, fetchedData, dbname)
	return fetchedData, nil
}

// BulkDBUpdate - CouchDB implementation of BulkDBUpdate
func (psd *PouchDBServiceDatastoreCouchDB) BulkDBUpdate(dbname string, request map[string]interface{}) ([]map[string]interface{}, error) {
	return bulkUpdate(createDBPathStr(psd.couchHost, dbname), request)
}

// CheckDBAvailability - CouchDB inmplementation of CheckDBAvailability
func (psd *PouchDBServiceDatastoreCouchDB) CheckDBAvailability(dbName string) (map[string]interface{}, error) {
	resource, err := couchdb.NewResource(createDBPathStr(psd.couchHost, dbName), nil)
	if err != nil {
		logger.Log.Debugf("Falied to check availability of DB %s: %s", dbName, err.Error())
		return nil, err
	}

	// Retrieve the DB status data from CouchDB
	fetchedData, err := checkIfAvailable(resource)
	if err != nil {
		return nil, err
	}

	logger.Log.Debugf("CheckAvailibility complete: %v\n", fetchedData)
	return fetchedData, nil
}

// GetAllDBDocs - CouchDB inmplementation of GetAllDBDocs
func (psd *PouchDBServiceDatastoreCouchDB) GetAllDBDocs(dbname string, request map[string]interface{}) (map[string]interface{}, error) {
	return getAllDocsFromDB(createDBPathStr(psd.couchHost, dbname), request)
}

// CreateDB - Couch inmplementation of CreateDB
func (psd *PouchDBServiceDatastoreCouchDB) CreateDB(dbname string) (map[string]interface{}, error) {
	resource, err := couchdb.NewResource(psd.couchHost, nil)
	if err != nil {
		logger.Log.Debugf("Falied to create DB %s: %s", dbname, err.Error())
		return nil, err
	}

	// Issue request to create the DB
	fetchedData, err := createDB(dbname, resource)
	if err != nil {
		return nil, err
	}

	logger.Log.Debugf("DB %s created: %v\n", dbname, fetchedData)
	return fetchedData, nil
}

// GetDoc - CouchDB inmplementation of GetDoc
func (psd *PouchDBServiceDatastoreCouchDB) GetDoc(dbname string, docID string, queryParams *url.Values, headers *http.Header) (map[string]interface{}, error) {
	logger.Log.Debugf("Attempting to fetch %s %s for DB %s with options: %v", ds.DBDocStr, docID, dbname, queryParams)

	resource, err := couchdb.NewResource(createDBPathStr(psd.couchHost, dbname), *headers)
	if err != nil {
		logger.Log.Debugf("Falied to create DB %s: %s", dbname, err.Error())
		return nil, err
	}

	// db, err := getDatabase(createDBPathStr(psd.couchHost, dbname))
	// if err != nil {
	// 	return nil, err
	// }

	// Retrieve the DB Changes Feed data from CouchDB
	fetchedData, err := getDoc(docID, resource, queryParams, headers)
	if err != nil {
		return nil, err
	}

	logger.Log.Debugf("Fetched %s %s for DB %s returned: %v\n", ds.DBDocStr, docID, dbname, fetchedData)
	return fetchedData, nil
}

// BulkDBGet - Couch implementation of BulkDBGet
func (psd *PouchDBServiceDatastoreCouchDB) BulkDBGet(dbname string, queryParams *url.Values, request map[string]interface{}) (map[string]interface{}, error) {
	logger.Log.Debugf("Performing %s %v on DB %s", ds.DBBulkGetStr, request, dbname)

	// Create a resource that can make the bulk fetch call to Couch
	resource, err := couchdb.NewResource(createDBPathStr(psd.couchHost, dbname), nil)
	if err != nil {
		logger.Log.Debugf("Falied to perform %s: %s", ds.DBBulkGetStr, err.Error())
		return nil, err
	}

	// Retrieve the bulk data from CouchDB
	fetchedData, err := performBulkGet(queryParams, request, resource)
	if err != nil {
		return nil, err
	}

	// DB Sync Checkpoint retrieved, send the response
	logger.Log.Debugf("Completed %s on DB %s: %v\n", ds.DBBulkGetStr, dbname, fetchedData)
	return fetchedData, nil
}

// ************************ Extensions of CouchDB-GoLang functionality ************************ //

// checkIfAvailable - contacts the CouchDB server to ascertain availability
func checkIfAvailable(resource *couchdb.Resource) (map[string]interface{}, error) {
	var jsonMap map[string]interface{}

	_, data, err := resource.GetJSON("", nil, nil)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(data, &jsonMap)
	if err != nil {
		return nil, err
	}

	return jsonMap, nil
}

func getDoc(docID string, resource *couchdb.Resource, queryParams *url.Values, headers *http.Header) (map[string]interface{}, error) {

	_, data, err := resource.GetJSON(docID, *headers, *queryParams)
	if err != nil {
		return nil, err
	}

	isArrayResponse := queryParams != nil && queryParams.Get("open_revs") != ""
	if isArrayResponse {
		// Parse result as array and then add to object to return.
		var jsonMap []map[string]interface{}
		err = json.Unmarshal(data, &jsonMap)
		if err != nil {
			return nil, err
		}

		result := make(map[string]interface{})
		result["data"] = jsonMap
		return result, nil
	}

	// Not an array response, just parse as object and then return
	// the result wrapped in the result object
	result := make(map[string]interface{})
	val, err := parseData(data)
	if err != nil {
		return nil, err
	}
	result["data"] = val
	return result, nil
}

func createDB(dbname string, resource *couchdb.Resource) (map[string]interface{}, error) {
	_, data, err := resource.PutJSON(dbname, nil, nil, nil)
	if err != nil {
		return nil, err
	}

	return parseData(data)
}

// fetchRevDiff - retrieves the revision diff for the provided map to revision list request.
func fetchRevDiff(body map[string]interface{}, resource *couchdb.Resource) (map[string]interface{}, error) {
	_, data, err := resource.PostJSON("_revs_diff", nil, body, nil)
	if err != nil {
		return nil, err
	}

	return parseData(data)
}

func addDesignDocument(body map[string]interface{}, resource *couchdb.Resource) (map[string]interface{}, error) {
	_, data, err := resource.PutJSON("", nil, body, nil)
	if err != nil {
		return nil, err
	}

	return parseData(data)
}

func fetchDesignDocumentResults(body map[string]interface{}, dbName string, indexName string) (map[string]interface{}, error) {
	resource, err := couchdb.NewResource(dbName, nil)
	if err != nil {
		logger.Log.Debugf("Falied to create resource to DB %s: %s", dbName, err.Error())
		return nil, err
	}

	_, data, err := resource.PostJSON(indexName, nil, body, nil)
	if err != nil {
		return nil, err
	}

	return parseData(data)
}

func performBulkUpdate(body map[string]interface{}, resource *couchdb.Resource) ([]map[string]interface{}, error) {
	_, data, err := resource.PostJSON("_bulk_docs", nil, body, nil)
	if err != nil {
		return nil, err
	}

	return parseDataArray(data)
}

func performBulkGet(queryParams *url.Values, body map[string]interface{}, resource *couchdb.Resource) (map[string]interface{}, error) {
	_, data, err := resource.PostJSON("_bulk_get", nil, body, *queryParams)
	if err != nil {
		return nil, err
	}

	return parseData(data)
}

func fetchAllDocs(body map[string]interface{}, resource *couchdb.Resource) (map[string]interface{}, error) {
	_, data, err := resource.PostJSON("_all_docs", nil, body, nil)
	if err != nil {
		return nil, err
	}

	return parseData(data)
}

func parseData(data []byte) (map[string]interface{}, error) {
	result := map[string]interface{}{}
	err := json.Unmarshal(data, &result)
	if err != nil {
		return result, err
	}
	if _, ok := result["error"]; ok {
		reason := result["reason"].(string)
		return result, errors.New(reason)
	}
	return result, nil
}

func parseDataArray(data []byte) ([]map[string]interface{}, error) {
	result := []map[string]interface{}{}
	err := json.Unmarshal(data, &result)
	if err != nil {
		return result, err
	}

	return result, nil
}

func getAllDocsFromDB(dbName string, request map[string]interface{}) (map[string]interface{}, error) {
	logger.Log.Debugf("Performing %s %v on DB %s", ds.DBAllDocsStr, request, dbName)

	// Create a resource that can make the fetch all docs call to Couch
	resource, err := couchdb.NewResource(dbName, nil)
	if err != nil {
		logger.Log.Debugf("Falied to fetch %s: %s", ds.DBAllDocsStr, err.Error())
		return nil, err
	}

	// Retrieve the all doc metadata from CouchDB
	fetchedData, err := fetchAllDocs(request, resource)
	if err != nil {
		return nil, err
	}

	// All Docs data retrieved, send the response
	logger.Log.Debugf("Completed fetch of %s on DB %s\n", ds.DBAllDocsStr, dbName)
	return fetchedData, nil
}

func bulkUpdate(dbName string, request map[string]interface{}) ([]map[string]interface{}, error) {
	logger.Log.Debugf("Performing %s %v on DB %s", ds.DBBulkUpdateStr, request, dbName)

	// Create a resource that can make the bulk update call to Couch
	resource, err := couchdb.NewResource(dbName, nil)
	if err != nil {
		logger.Log.Debugf("Falied to perform %s: %s", ds.DBBulkUpdateStr, err.Error())
		return nil, err
	}

	// Retrieve the checkpoint data from CouchDB
	fetchedData, err := performBulkUpdate(request, resource)
	if err != nil {
		return nil, err
	}

	// DB Sync Checkpoint retrieved, send the response
	logger.Log.Debugf("Completed %s on DB %s\n", ds.DBBulkUpdateStr, dbName)
	return fetchedData, nil
}

func purgeDB(dbName string) error {
	// Get a list of all documents from the DB:
	docs, err := getAllDocsFromDB(dbName, map[string]interface{}{})
	if err != nil {
		return fmt.Errorf("Unable to purge DB %s: %s", dbName, err.Error())
	}

	if docs["rows"] != nil {
		// There are documents to delete, build up a bulk delete request body for all of them
		docList := docs["rows"].([]interface{})
		docsToDelete := make([]map[string]interface{}, 0)
		for _, doc := range docList {
			docObj := doc.(map[string]interface{})
			docID := docObj["id"].(string)
			docRev := docObj["value"].(map[string]interface{})["rev"].(string)
			docsToDelete = append(docsToDelete, map[string]interface{}{"_id": docID, "_rev": docRev, "_deleted": true})
		}

		deleteBody := map[string]interface{}{"docs": docsToDelete}
		logger.Log.Debugf("Attempting to delete the following from DB %s: %v", dbName, docsToDelete)

		_, err = bulkUpdate(dbName, deleteBody)
		if err != nil {
			return fmt.Errorf("Unable to purge DB %s: %s", dbName, err.Error())
		}
	}

	return nil
}

// ************************ End of CouchDB-GoLang functionality ************************ //
