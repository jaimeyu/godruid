package inMemory

import (
	"errors"
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

// // CheckAvailablility - InMemory implementation of CheckAvailablility
// func (pdb *PouchDBPluginServiceDatastoreInMemory) CheckAvailablility() (*pb.DBAvailableResponse, error) {
// 	// Stub to implement
// 	return nil, errors.New("CheckAvailablility() not implemented for InMemory DB")
// }

// // StoreDBSyncCheckpoint - InMemory implementation of StoreDBSyncCheckpoint
// func (pdb *PouchDBPluginServiceDatastoreInMemory) StoreDBSyncCheckpoint(dbCheckpoint *pb.DBSyncCheckpoint) (*pb.DBSyncCheckpointPutResponse, error) {
// 	// Stub to implement
// 	return nil, errors.New("StoreDBSyncCheckpoint() not implemented for InMemory DB")
// }

// // GetDBSyncCheckpoint - InMemory implementation of GetDBSyncCheckpoint
// func (pdb *PouchDBPluginServiceDatastoreInMemory) GetDBSyncCheckpoint(dbCheckpointID *pb.DBSyncCheckpointId, appendPrefix bool) (*pb.DBSyncCheckpoint, error) {
// 	// Stub to implement
// 	return nil, errors.New("GetDBSyncCheckpoint() not implemented for InMemory DB")
// }
