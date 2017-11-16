package couchDB

import (
	"encoding/json"
	"fmt"
	"net/url"

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
	resource  *couchdb.Resource
}

// CreatePouchDBServiceDAO - instantiates a CouchDB implementation of the
// PouchDBPluginServiceDatastore.
func CreatePouchDBServiceDAO() (*PouchDBServiceDatastoreCouchDB, error) {
	result := new(PouchDBServiceDatastoreCouchDB)
	cfg, err := gather.GetActiveConfig()
	if err != nil {
		logger.Log.Debugf("Falied to instantiate PouchDBServiceDatastoreCouchDB: %s", err.Error())
		return nil, err
	}

	provDBURL := fmt.Sprintf("%s:%d",
		cfg.ServerConfig.Datastore.BindIP,
		cfg.ServerConfig.Datastore.BindPort)
	logger.Log.Debug("Admin Service CouchDB URL is: ", provDBURL)
	result.couchHost = provDBURL

	resource, err := couchdb.NewResource(result.couchHost, nil)
	if err != nil {
		logger.Log.Debugf("Falied to instantiate PouchDBServiceDatastoreCouchDB: %s", err.Error())
		return nil, err
	}
	result.resource = resource

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

	// Retrieve the DB Changes Feed data from CouchDB
	fetchedData, err := psd.checkIfAvailable()
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
	// Validate the request to ensure this operation is valid:

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

// ************************ Extensions of CouchDB-GoLang functionality ************************ //

// checkIfAvailable - contacts the CouchDB server to ascertain availability
func (psd *PouchDBServiceDatastoreCouchDB) checkIfAvailable() (map[string]interface{}, error) {
	var jsonMap map[string]interface{}

	_, data, err := psd.resource.GetJSON("", nil, nil)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(data, &jsonMap)
	if err != nil {
		return nil, err
	}

	return jsonMap, nil
}

// ************************ End of CouchDB-GoLang functionality ************************ //
