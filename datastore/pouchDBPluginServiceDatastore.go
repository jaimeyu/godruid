package datastore

import (
	"net/url"
)

const (
	// ChangeFeedStr - common name to refer to a DB Change Feed. For use in logs.
	ChangeFeedStr = "Changes Feed"

	// DBSyncCheckpointStr - common name to refer to a DB Checkpoint. For use in logs.
	DBSyncCheckpointStr = "DB Sync Checkpoint"

	// DBRevDiffStr - common name to refer to a DB Revision Diff. For use in logs.
	DBRevDiffStr = "DB Revision Diff"

	// DBSyncCheckpointPrefixStr - prefix required for storing objects as "local documents" per database.
	// Used during pouch - couch db syncronization
	DBSyncCheckpointPrefixStr = "_local/"
)

// PouchDBPluginServiceDatastore - interface which provides the functionality
// of the PouchDBPluginService Datastore.
type PouchDBPluginServiceDatastore interface {
	// GetChanges(*pb.DBChangesRequest) (*pb.DBChangesResponse, error)
	GetChanges(dbname string, queryParams *url.Values) (map[string]interface{}, error)
	CheckAvailability() (map[string]interface{}, error)
	StoreDBSyncCheckpoint(dbname string, queryParams *url.Values, request map[string]interface{}) (map[string]interface{}, error)
	GetDBSyncCheckpoint(dbName string, documentID string) (map[string]interface{}, error)
	GetDBRevisionDiff(dbname string, request map[string]interface{}) (map[string]interface{}, error)
}
