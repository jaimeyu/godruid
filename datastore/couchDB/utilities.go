package couchDB

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strings"

	"github.com/accedian/adh-gather/logger"
	couchdb "github.com/leesper/couchdb-golang"
)

// GetDatabase - returns the object used to issue commands to a CouchDB database
// instance.
func (couchDB *AdminServiceDatastoreCouchDB) GetDatabase() (*couchdb.Database, error) {
	db, err := couchdb.NewDatabase(couchDB.dbName)
	if err != nil {
		logger.Log.Errorf("Unable to connect to CouchDB %s: %v\n", couchDB.server, err)
		return nil, err
	}

	return db, nil
}

// InsertField - Helper metho to add fields to data models during conversion from ADH data
// model to CouchDB data model and vice-versa. Useful for metadata fields
// like '_id' and '_rev' that are key fields for operations in CouchDB.
func InsertField(genericData map[string]interface{}, fieldName string) {
	if genericData[fieldName] != nil {
		logger.Log.Debugf("Adding '%s' field to data: %v", fieldName, genericData)

		if strings.HasPrefix(fieldName, "_") {
			genericData[fieldName[1:]] = genericData[fieldName]
		} else {
			genericData["_"+fieldName] = genericData[fieldName]
		}

	}
}

// ConvertDataToCouchDbSupportedModel - Turns any object into a CouchDB ready entry
// that can be stored. Changes the provided object into a map[string]interface{} generic
//  object and adds the neccesary CouchDB metadata fields '_id' and '_rev' as provided
//by to orginal data.
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

	// Add in the _id field and _rev fields that are necessary for CouchDB
	InsertField(genericFormat, "id")
	InsertField(genericFormat, "rev")

	// Successfully converted the User
	return genericFormat, nil
}

// ConvertGenericObjectToBytesWithCouchDbFields - takes a generic set of CouchDB data,
// adds the necessary data fields back into the data model and then converts the data
// to a []byte. Useful as a preparation step before unmarshalling the bytes into a known
// ADH data model object.
func ConvertGenericObjectToBytesWithCouchDbFields(genericObject map[string]interface{}) ([]byte, error) {
	// Add in the _id field and _rev fields that are necessary for CouchDB
	InsertField(genericObject, "_id")
	InsertField(genericObject, "_rev")
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
	selector := fmt.Sprintf(`datatype == "%s"`, dataType)
	fetchedData, err := db.Query(nil, selector, nil, nil, nil, nil)
	if err != nil {
		logger.Log.Errorf("Error retrieving all %ss: %v\n", dataTypeStrForLogging, err)
		return nil, err
	}

	return fetchedData, nil
}
