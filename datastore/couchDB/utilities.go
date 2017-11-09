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
func ConvertDataToCouchDbSupportedModel(data interface{}) (map[string]interface{}, error) {
	dataToBytes, err := json.Marshal(data)
	if err != nil {
		logger.Log.Errorf("Unable to convert data to CouchDB format to persist: %v\n", err)
		return nil, err
	}
	var genericFormat map[string]interface{}
	err = json.Unmarshal(dataToBytes, &genericFormat)
	if err != nil {
		logger.Log.Errorf("Unable to convert data to CouchDB format to persist: %v\n", err)
		return nil, err
	}

	// Successfully converted the User
	return genericFormat, nil
}

// ConvertGenericObjectToBytesWithCouchDbFields - takes a generic set of CouchDB data,
// adds the necessary data fields back into the data model and then converts the data
// to a []byte. Useful as a preparation step before unmarshalling the bytes into a known
// ADH data model object.
func ConvertGenericObjectToBytesWithCouchDbFields(genericObject map[string]interface{}) ([]byte, error) {
	genericUserInBytes, err := json.Marshal(genericObject)
	if err != nil {
		logger.Log.Errorf("Error converting generic data to bytes: %v\n", err)
		return nil, err
	}

	return genericUserInBytes, nil
}

// StoreDataInCouchDB - takes data that is already in a format ready to store in CouchDB
// and attempts to store it. Parameters are:
// dataToStore(the CouchDB ready data to be stored),
// dataTypeStrForLogging(human readable string of the type of data being stored),
// db (the CouchDB connector used to store the data.)
func StoreDataInCouchDB(dataToStore map[string]interface{}, dataTypeStrForLogging string, db *couchdb.Database) (string, string, error) {
	logger.Log.Debugf("Attempting to store %s: %v", dataTypeStrForLogging, dataToStore)

	// Store the user in PROV DB
	options := new(url.Values)
	id, rev, err := db.Save(dataToStore, *options)
	if err != nil {
		logger.Log.Errorf("Unable to store %s: %v\n", dataTypeStrForLogging, err)
		return "", "", err
	}

	logger.Log.Debugf("Successfully stored %s: id: %s, rev: %s", dataTypeStrForLogging, id, rev)
	return id, rev, nil
}

// DeleteByDocID - deletes a document (by the documentID) in specified CouchDB instance.
func DeleteByDocID(docID string, dataTypeStrForLogging string, db *couchdb.Database) error {
	logger.Log.Infof("Attempting to delete %s %s\n", dataTypeStrForLogging, docID)

	err := db.Delete(docID)
	if err != nil {
		logger.Log.Errorf("Error deleting %s %s: %v\n", dataTypeStrForLogging, docID, err)
		return err
	}

	return nil
}

// GetByDocID - retrieves a document (by documentID) from the specified CouchDB instamnce.
func GetByDocID(docID string, dataTypeStrForLogging string, db *couchdb.Database) (map[string]interface{}, error) {
	logger.Log.Infof("Attempting to retrieve %s %s\n", dataTypeStrForLogging, docID)

	// Get the Admin User from CouchDB
	options := new(url.Values)
	fetchedData, err := db.Get(docID, *options)
	if err != nil {
		logger.Log.Errorf("Error retrieving %s %s: %v\n", dataTypeStrForLogging, docID, err)
		return nil, err
	}

	return fetchedData, nil
}

// GetAllOfType - retrieves a list of data of the specified dataType from the couchDB instance.
func GetAllOfType(dataType string, dataTypeStrForLogging string, db *couchdb.Database) ([]map[string]interface{}, error) {
	logger.Log.Infof("Attempting to retrieve all %ss\n", dataTypeStrForLogging)

	// Get the Admin User from CouchDB
	selector := fmt.Sprintf(`data.datatype == "%s"`, dataType)
	fetchedData, err := db.Query(nil, selector, nil, nil, nil, nil)
	if err != nil {
		logger.Log.Errorf("Error retrieving all %ss: %v\n", dataTypeStrForLogging, err)
		return nil, err
	}

	return fetchedData, nil
}

// GetAllOfTypeByIDPrefix - retrieves a list of data whose ids start with
// the specified prefix from the couchDB instance.
func GetAllOfTypeByIDPrefix(dataType string, dataTypeStrForLogging string, db *couchdb.Database) ([]map[string]interface{}, error) {
	logger.Log.Infof("Attempting to retrieve all %ss\n", dataTypeStrForLogging)

	// Get the Admin User from CouchDB
	selector := fmt.Sprintf(`regex(_id, "^%s")`, dataType)
	fetchedData, err := db.Query(nil, selector, nil, nil, nil, nil)
	if err != nil {
		logger.Log.Errorf("Error retrieving all %ss: %v\n", dataTypeStrForLogging, err)
		return nil, err
	}

	return fetchedData, nil
}

// ConvertGenericCouchDataToObject - takes an empty object of a known type and populates
// that object with the generic data.
func ConvertGenericCouchDataToObject(genericData map[string]interface{}, dataContainer interface{}, dataTypeStr string) error {
	genericDataInBytes, err := ConvertGenericObjectToBytesWithCouchDbFields(genericData)
	if err != nil {
		return err
	}

	err = json.Unmarshal(genericDataInBytes, &dataContainer)
	if err != nil {
		logger.Log.Errorf("Error converting generic data to %s type: %v\n", dataTypeStr, err)
		return err
	}

	logger.Log.Debugf("Converted generic data to %s: %v\n", dataTypeStr, dataContainer)

	return nil
}

// GetDatabase - returns the object used to issue commands to a CouchDB database
// instance.
func GetDatabase(dbConnectionName string) (*couchdb.Database, error) {
	db, err := couchdb.NewDatabase(dbConnectionName)
	if err != nil {
		logger.Log.Errorf("Unable to connect to CouchDB %s: %v\n", dbConnectionName, err)
		return nil, err
	}

	return db, nil
}

// CreateDBPathStr - Helper method to handle logic specific to CouchDB for creating the
// URL to a database. Works by taking a server name (i.e. http://localhost:5894) and
// appending the path to the db.
func CreateDBPathStr(dbServerStr string, dbPathStr string) string {
	return strings.Join([]string{dbServerStr, "/", dbPathStr}, "")
}
