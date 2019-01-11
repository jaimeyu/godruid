package datastore

import (
	"net/http"
	"net/url"
)

const (
	// ChangeFeedStr - common name to refer to a DB Change Feed. For use in logs.
	ChangeFeedStr = "Changes Feed"

	// DBSyncCheckpointStr - common name to refer to a DB Checkpoint. For use in logs.
	DBSyncCheckpointStr = "Sync Checkpoint"

	// DBRevDiffStr - common name to refer to a DB Revision Diff. For use in logs.
	DBRevDiffStr = "Revision Diff"

	// DBBulkUpdateStr - common name to refer to a DB Bulk Update. For use in logs.
	DBBulkUpdateStr = "Bulk Update"

	// DBBulkGetStr - common name to refer to a DB Bulk Update. For use in logs.
	DBBulkGetStr = "Bulk Get"

	// DBAllDocsStr - common name to refer to the DB metadata for all docs. For use in logs.
	DBAllDocsStr = "All Docs"

	// DBDocStr - common name to refer to a Document. For use in logs.
	DBDocStr = "Document"

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
	CheckDBAvailability(dbName string) (map[string]interface{}, error)
	StoreDBSyncCheckpoint(dbname string, queryParams *url.Values, request map[string]interface{}) (map[string]interface{}, error)
	GetDBSyncCheckpoint(dbName string, documentID string) (map[string]interface{}, error)
	GetDBRevisionDiff(dbname string, request map[string]interface{}) (map[string]interface{}, error)
	BulkDBUpdate(dbname string, request map[string]interface{}) ([]map[string]interface{}, error)
	BulkDBGet(dbname string, queryParams *url.Values, request map[string]interface{}) (map[string]interface{}, error)
	GetAllDBDocs(dbname string, request map[string]interface{}) (map[string]interface{}, error)
	CreateDB(dbname string) (map[string]interface{}, error)

	// Have to pass headers into this call as it changes the response type of the call.
	GetDoc(dbname string, docID string, queryParams *url.Values, headers *http.Header) (map[string]interface{}, error)

	GetByDesignDocument(dbName string, indexName string) (map[string]interface{}, error)
}
