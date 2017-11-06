package couchDB

import (
	"encoding/json"
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
