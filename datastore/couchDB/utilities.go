package couchDB

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strings"

	"github.com/accedian/adh-gather/logger"
	couchdb "github.com/leesper/couchdb-golang"
)

// ConvertDataToCouchDbSupportedModel - Turns any object into a CouchDB ready entry
// that can be stored. Changes the provided object into a map[string]interface{} generic
// object.
func convertDataToCouchDbSupportedModel(data interface{}) (map[string]interface{}, error) {
	dataToBytes, err := json.Marshal(data)
	if err != nil {
		logger.Log.Debugf("Unable to convert data to CouchDB format to persist: %s", err.Error())
		return nil, err
	}
	var genericFormat map[string]interface{}
	err = json.Unmarshal(dataToBytes, &genericFormat)
	if err != nil {
		logger.Log.Debugf("Unable to convert data to CouchDB format to persist: %s", err.Error())
		return nil, err
	}

	// Successfully converted the User
	return genericFormat, nil
}

// ConvertGenericObjectToBytesWithCouchDbFields - takes a generic set of CouchDB data,
// adds the necessary data fields back into the data model and then converts the data
// to a []byte. Useful as a preparation step before unmarshalling the bytes into a known
// ADH data model object.
func convertGenericObjectToBytesWithCouchDbFields(genericObject map[string]interface{}) ([]byte, error) {
	genericUserInBytes, err := json.Marshal(genericObject)
	if err != nil {
		logger.Log.Debugf("Error converting generic data to bytes: %s", err.Error())
		return nil, err
	}

	return genericUserInBytes, nil
}

// StoreDataInCouchDB - takes data that is already in a format ready to store in CouchDB
// and attempts to store it. Parameters are:
// dataToStore(the CouchDB ready data to be stored),
// dataTypeStrForLogging(human readable string of the type of data being stored),
// db (the CouchDB connector used to store the data.)
func storeDataInCouchDB(dataToStore map[string]interface{}, dataTypeStrForLogging string, db *couchdb.Database) (string, string, error) {
	return storeDataInCouchDBWithQueryParams(dataToStore, dataTypeStrForLogging, db, nil)
}

// StoreDataInCouchDB - takes data that is already in a format ready to store in CouchDB
// and attempts to store it. Parameters are:
// dataToStore(the CouchDB ready data to be stored),
// dataTypeStrForLogging(human readable string of the type of data being stored),
// db (the CouchDB connector used to store the data.)
// queryParams (the query parameters passed to the call to store the data)
func storeDataInCouchDBWithQueryParams(dataToStore map[string]interface{}, dataTypeStrForLogging string, db *couchdb.Database, queryParams *url.Values) (string, string, error) {
	logger.Log.Debugf("Attempting to store %s: %v", dataTypeStrForLogging, dataToStore)

	// Store the user in PROV DB
	if queryParams == nil {
		queryParams = new(url.Values)
	}
	id, rev, err := db.Save(dataToStore, *queryParams)
	if err != nil {
		logger.Log.Debugf("Unable to store %s: %s", dataTypeStrForLogging, err.Error())
		return "", "", err
	}

	logger.Log.Debugf("Successfully stored %s: id: %s, rev: %s", dataTypeStrForLogging, id, rev)
	return id, rev, nil
}

// DeleteByDocID - deletes a document (by the documentID) in specified CouchDB instance.
func deleteByDocID(docID string, dataTypeStrForLogging string, db *couchdb.Database) error {
	logger.Log.Debugf("Attempting to delete %s %s\n", dataTypeStrForLogging, docID)

	err := db.Delete(docID)
	if err != nil {
		logger.Log.Debugf("Error deleting %s %s: %s", dataTypeStrForLogging, docID, err.Error())
		return err
	}

	return nil
}

// GetByDocID - retrieves a document (by documentID) from the specified CouchDB instamnce.
func getByDocID(docID string, dataTypeStrForLogging string, db *couchdb.Database) (map[string]interface{}, error) {
	return getByDocIDWithQueryParams(docID, dataTypeStrForLogging, db, nil)
}

func getByDocIDWithQueryParams(docID string, dataTypeStrForLogging string, db *couchdb.Database, queryParams *url.Values) (map[string]interface{}, error) {
	logger.Log.Debugf("Attempting to retrieve %s %s\n", dataTypeStrForLogging, docID)

	// Get the Document from CouchDB
	if queryParams == nil {
		queryParams = new(url.Values)
	}
	fetchedData, err := db.Get(docID, *queryParams)
	if err != nil {
		logger.Log.Debugf("Error retrieving %s %s: %s", dataTypeStrForLogging, docID, err.Error())
		return nil, err
	}

	return fetchedData, nil
}

// GetAllOfType - retrieves a list of data of the specified dataType from the couchDB instance.
func getAllOfType(dataType string, dataTypeStrForLogging string, db *couchdb.Database) ([]map[string]interface{}, error) {
	logger.Log.Debugf("Attempting to retrieve all %ss\n", dataTypeStrForLogging)

	// Get the Admin User from CouchDB
	selector := fmt.Sprintf(`data.datatype == "%s"`, dataType)
	fetchedData, err := db.Query(nil, selector, nil, nil, nil, nil)
	if err != nil {
		logger.Log.Debugf("Error retrieving all %ss: %s", dataTypeStrForLogging, err.Error())
		return nil, err
	}

	return fetchedData, nil
}

// GetAllOfTypeByIDPrefix - retrieves a list of data whose ids start with
// the specified prefix from the couchDB instance.
func getAllOfTypeByIDPrefix(dataType string, dataTypeStrForLogging string, db *couchdb.Database) ([]map[string]interface{}, error) {
	logger.Log.Debugf("Attempting to retrieve all %ss\n", dataTypeStrForLogging)

	// Get the Admin User from CouchDB
	selector := fmt.Sprintf(`regex(_id, "^%s")`, dataType)
	fetchedData, err := db.Query(nil, selector, nil, nil, nil, nil)
	if err != nil {
		logger.Log.Debugf("Error retrieving all %ss: %s", dataTypeStrForLogging, err.Error())
		return nil, err
	}

	return fetchedData, nil
}

// ConvertGenericCouchDataToObject - takes an empty object of a known type and populates
// that object with the generic data.
func convertGenericCouchDataToObject(genericData map[string]interface{}, dataContainer interface{}, dataTypeStr string) error {
	genericDataInBytes, err := convertGenericObjectToBytesWithCouchDbFields(genericData)
	if err != nil {
		return err
	}

	err = json.Unmarshal(genericDataInBytes, &dataContainer)
	if err != nil {
		logger.Log.Debugf("Error converting generic data to %s type: %s", dataTypeStr, err.Error())
		return err
	}

	logger.Log.Debugf("Converted generic data to %s: %v\n", dataTypeStr, dataContainer)

	return nil
}

// GetDatabase - returns the object used to issue commands to a CouchDB database
// instance.
func getDatabase(dbConnectionName string) (*couchdb.Database, error) {
	db, err := couchdb.NewDatabase(dbConnectionName)
	if err != nil {
		logger.Log.Debugf("Unable to connect to CouchDB %s: %s", dbConnectionName, err.Error())
		return nil, err
	}

	return db, nil
}

// CreateDBPathStr - Helper method to handle logic specific to CouchDB for creating the
// URL to a database. Works by taking a server name (i.e. http://localhost:5894) and
// appending the path to the db.
func createDBPathStr(pathParts ...string) string {
	return strings.Join(pathParts, "/")
}

func accessDBChangesFeed(db *couchdb.Database, queryParams *url.Values) (map[string]interface{}, error) {
	// Get access to the Changes feed for the DB
	fetchedData, err := db.Changes(*queryParams)
	if err != nil {
		logger.Log.Debugf("Error accessing Changes feed: %s", err.Error())
		return nil, err
	}

	return fetchedData, nil
}
