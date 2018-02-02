package couchDB

import (
	"fmt"

	"github.com/accedian/adh-gather/config"
	"github.com/accedian/adh-gather/gather"
	"github.com/accedian/adh-gather/logger"
)

// TestDataServiceDatastoreCouchDB - couchDB test data datastore impl
type TestDataServiceDatastoreCouchDB struct {
	couchHost string
	cfg       config.Provider
}

// CreateTestDataServiceDAO - instantiates a CouchDB implementation of the
// TestDataServiceDatastore.
func CreateTestDataServiceDAO() (*TestDataServiceDatastoreCouchDB, error) {
	result := new(TestDataServiceDatastoreCouchDB)
	result.cfg = gather.GetConfig()

	// Couch Server Configuration
	provDBURL := fmt.Sprintf("%s:%d",
		result.cfg.GetString(gather.CK_server_datastore_ip.String()),
		result.cfg.GetInt(gather.CK_server_datastore_port.String()))
	logger.Log.Debugf("TestData Service CouchDB URL is: %s", provDBURL)
	result.couchHost = provDBURL

	return result, nil
}

// GetAllDocsByDatatype - CouchDB implementation of GetAllDocsByDatatype
func (testDB *TestDataServiceDatastoreCouchDB) GetAllDocsByDatatype(dbName string, datatype string) ([]map[string]interface{}, error) {
	if len(dbName) == 0 || len(datatype) == 0 {
		return nil, fmt.Errorf("Unable to retrieve documnets if DB name and datatype are provided")
	}

	fullDBName := createDBPathStr(testDB.couchHost, dbName)
	db, err := getDatabase(fullDBName)
	if err != nil {
		return nil, err
	}

	return getAllOfTypeByIDPrefix(datatype, datatype, db)
}
