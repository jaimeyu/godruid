package inMemory

import (
	"errors"
	"net/http"
	"net/url"

	"github.com/accedian/adh-gather/datastore"
)

// PouchDBPluginServiceDatastoreInMemory - struct responsible for handling
// database operations for the PouchDBPlugin Service when using local memory
// as the storage option. Useful for tests.
type PouchDBPluginServiceDatastoreInMemory struct {
}

// CreatePouchDBPluginServiceDAO - returns an in-memory implementation of the PouchDB
// Plugin Service datastore.
func CreatePouchDBPluginServiceDAO() (datastore.PouchDBPluginServiceDatastore, error) {
	res := new(PouchDBPluginServiceDatastoreInMemory)

	return res, nil
}

// GetChanges - InMemory implementation of GetChanges
func (pdb *PouchDBPluginServiceDatastoreInMemory) GetChanges(dbname string, queryParams *url.Values) (map[string]interface{}, error) {
	// Stub to implement
	return nil, errors.New("GetChanges() not implemented for InMemory DB")
}

// CheckAvailability - InMemory implementation of CheckAvailability
func (pdb *PouchDBPluginServiceDatastoreInMemory) CheckAvailability() (map[string]interface{}, error) {
	// Stub to implement
	return nil, errors.New("CheckAvailablility() not implemented for InMemory DB")
}

// StoreDBSyncCheckpoint - InMemory implementation of StoreDBSyncCheckpoint
func (pdb *PouchDBPluginServiceDatastoreInMemory) StoreDBSyncCheckpoint(dbname string, queryParams *url.Values, request map[string]interface{}) (map[string]interface{}, error) {
	// Stub to implement
	return nil, errors.New("StoreDBSyncCheckpoint() not implemented for InMemory DB")
}

// GetDBSyncCheckpoint - InMemory implementation of GetDBSyncCheckpoint
func (pdb *PouchDBPluginServiceDatastoreInMemory) GetDBSyncCheckpoint(dbName string, documentID string) (map[string]interface{}, error) {
	// Stub to implement
	return nil, errors.New("GetDBSyncCheckpoint() not implemented for InMemory DB")
}

// GetDBRevisionDiff - InMemory implementation of GetDBRevisionDiff
func (pdb *PouchDBPluginServiceDatastoreInMemory) GetDBRevisionDiff(dbname string, request map[string]interface{}) (map[string]interface{}, error) {
	// Stub to implement
	return nil, errors.New("GetDBRevisionDiff() not implemented for InMemory DB")
}

// BulkDBUpdate - InMemory implementation of BulkDBUpdate
func (pdb *PouchDBPluginServiceDatastoreInMemory) BulkDBUpdate(dbname string, request map[string]interface{}) ([]map[string]interface{}, error) {
	// Stub to implement
	return nil, errors.New("BulkDBUpdate() not implemented for InMemory DB")
}

// BulkDBGet - InMemory implementation of BulkDBGet
func (pdb *PouchDBPluginServiceDatastoreInMemory) BulkDBGet(dbname string, queryParams *url.Values, request map[string]interface{}) (map[string]interface{}, error) {
	// Stub to implement
	return nil, errors.New("BulkDBGet() not implemented for InMemory DB")
}

// CheckDBAvailability - InMemory inmplementation of CheckDBAvailability
func (pdb *PouchDBPluginServiceDatastoreInMemory) CheckDBAvailability(dbName string) (map[string]interface{}, error) {
	// Stub to implement
	return nil, errors.New("CheckDBAvailability() not implemented for InMemory DB")
}

// GetAllDBDocs - InMemory inmplementation of GetAllDBDocs
func (pdb *PouchDBPluginServiceDatastoreInMemory) GetAllDBDocs(dbname string, request map[string]interface{}) (map[string]interface{}, error) {
	// Stub to implement
	return nil, errors.New("GetAllDBDocs() not implemented for InMemory DB")
}

// CreateDB - InMemory inmplementation of CreateDB
func (pdb *PouchDBPluginServiceDatastoreInMemory) CreateDB(dbname string) (map[string]interface{}, error) {
	// Stub to implement
	return nil, errors.New("CreateDB() not implemented for InMemory DB")
}

// GetDoc - InMemory inmplementation of GetDoc
func (pdb *PouchDBPluginServiceDatastoreInMemory) GetDoc(dbname string, docID string, queryParams *url.Values, headers *http.Header) (map[string]interface{}, error) {
	// Stub to implement
	return nil, errors.New("GetDoc() not implemented for InMemory DB")
}
