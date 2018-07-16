package couchDB

import (
	"fmt"
	"strings"

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

// InsertTenantViews - CouchDB implementation of InsertTenantViews
func (testDB *TestDataServiceDatastoreCouchDB) InsertTenantViews(dbName string) string {
	if len(dbName) == 0 {
		return fmt.Sprintf("Unable to insert Tenant Views if no DB name is provided")
	}

	fullDBName := createDBPathStr(testDB.couchHost, dbName)
	moDB, err := getDatabase(fullDBName)
	if err != nil {
		return fmt.Sprintf("Unable to insert Tenant Views: %s", err.Error())
	}

	errorMessageContainer := []string{}
	// Store the views related to Monitored Objects
	for _, viewPayload := range getTenantViews() {
		// See if a view already exists
		existing, err := getByDocID(viewPayload["_id"].(string), "Tenant View", moDB)
		if existing == nil || existing["_rev"] == nil {
			errorMessageContainer = append(errorMessageContainer, fmt.Sprintf("View %s was not found: %s. Going to try to create it.", viewPayload["_id"].(string), err.Error()))

			// Data does not exist, try to insert it
			_, _, err = storeDataInCouchDBWithQueryParams(viewPayload, "TenantView", moDB, nil)
			if err != nil {
				errorMessageContainer = append(errorMessageContainer, fmt.Sprintf("Error tring to create View %s: %s", viewPayload["_id"].(string), err.Error()))
			} else {
				errorMessageContainer = append(errorMessageContainer, fmt.Sprintf("Successfully created View %s", viewPayload["_id"].(string)))
			}
			continue
		}

		// Record exists, let's update it
		viewPayload["_rev"] = existing["_rev"].(string)
		_, _, err = storeDataInCouchDBWithQueryParams(viewPayload, "TenantView", moDB, nil)
		if err != nil {
			errorMessageContainer = append(errorMessageContainer, fmt.Sprintf("Error tring to update View %s: %s", viewPayload["_id"].(string), err.Error()))
		} else {
			errorMessageContainer = append(errorMessageContainer, fmt.Sprintf("Successfully updated View %s", viewPayload["_id"].(string)))
		}
	}

	return strings.Join(errorMessageContainer, "\n")
}
