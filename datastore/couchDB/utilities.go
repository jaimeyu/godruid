package couchDB

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strings"
	"sync"

	ds "github.com/accedian/adh-gather/datastore"
	"github.com/accedian/adh-gather/logger"
	"github.com/accedian/adh-gather/models"
	tenmod "github.com/accedian/adh-gather/models/tenant"
	couchdb "github.com/leesper/couchdb-golang"
)

const defaultQueryResultsLimit = 1000

// When metadata is being updated in monitored objects, we want to issue a request to the view
// so couchdb would start to build/update the view. Since the builds function is asynchronous,
// we don't want overlapping calls to the build since it makes no sense to start/stop a build
// while it is functioning.
var couchdbViewBuilderBusyMap sync.Map

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

	// Now that the object is generic, need to structure it with Couch specific fields:
	result := map[string]interface{}{}
	if genericFormat["_id"] != nil {
		result["_id"] = genericFormat["_id"]
		delete(genericFormat, "_id")
	}
	if genericFormat["_rev"] != nil {
		if len(genericFormat["_rev"].(string)) != 0 {
			result["_rev"] = genericFormat["_rev"]
		}
		delete(genericFormat, "_rev")
	}
	result["data"] = genericFormat

	// Successfully converted the User
	return result, nil
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

func convertGenericObjectToBytes(genericObject []map[string]interface{}) ([]byte, error) {
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
	logger.Log.Debugf("Attempting to store %s: %v", dataTypeStrForLogging, models.AsJSONString(dataToStore))

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

// getAvailableConfigs - Returns all tenant connectors matching tenantID, zone, that aren't already being used
func getAvailableConfigs(dataType, attr, value, dataTypeStrForLogging string, db *couchdb.Database) ([]map[string]interface{}, error) {
	logger.Log.Debugf("Attempting to retrieve all available %ss in %s %s\n", dataTypeStrForLogging, attr, value)

	selector := fmt.Sprintf(`data.%s == "%s" && data.datatype == "%s" && data.connectorInstanceId == ""`, attr, value, dataType)
	fetchedData, err := db.Query(nil, selector, nil, defaultQueryResultsLimit, nil, nil)
	if err != nil {
		logger.Log.Debugf("Error retrieving all %ss: %s", dataTypeStrForLogging, err.Error())
		return nil, err
	}

	return fetchedData, nil
}

// GetAllOfAny - retrieves a list of data of the specified attribute and value, along with the specified dataType from the couchDB instance.
func getAllOfAny(dataType, attr, value, dataTypeStrForLogging string, db *couchdb.Database) ([]map[string]interface{}, error) {
	logger.Log.Debugf("Attempting to retrieve all %ss in %s %s\n", dataTypeStrForLogging, attr, value)

	selector := fmt.Sprintf(`data.%s == "%s" && data.datatype == "%s"`, attr, value, dataType)
	fetchedData, err := db.Query(nil, selector, nil, defaultQueryResultsLimit, nil, nil)
	if err != nil {
		logger.Log.Debugf("Error retrieving all %ss: %s", dataTypeStrForLogging, err.Error())
		return nil, err
	}

	return fetchedData, nil
}

// GetAllOfType - retrieves a list of data of the specified dataType from the couchDB instance.
func getAllOfType(dataType string, dataTypeStrForLogging string, db *couchdb.Database) ([]map[string]interface{}, error) {
	logger.Log.Debugf("Attempting to retrieve all %ss\n", dataTypeStrForLogging)

	// Get the Admin User from CouchDB
	selector := fmt.Sprintf(`data.datatype == "%s"`, dataType)
	fetchedData, err := db.Query(nil, selector, nil, defaultQueryResultsLimit, nil, nil)
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

	// Get the data from CouchDB
	selector := fmt.Sprintf(`regex(_id, "^%s_")`, dataType)
	fetchedData, err := db.Query(nil, selector, nil, defaultQueryResultsLimit, nil, nil)
	if err != nil {
		logger.Log.Debugf("Error retrieving all %ss: %s", dataTypeStrForLogging, err.Error())
		return nil, err
	}

	// Strip out the prefix on all the IDs
	for _, data := range fetchedData {
		stripPrefixFromID(data)
	}

	return fetchedData, nil
}

func getAllInListOfIDs(idList []string, dataType string, dataTypeStrForLogging string, db *couchdb.Database) ([]map[string]interface{}, error) {
	logger.Log.Debugf("Attempting to retrieve all %ss\n", dataTypeStrForLogging)

	// Get the data from CouchDB
	selector := fmt.Sprintf(`in(_id, %s)`, printStringSliceInObjectDeclarationFormat(idList, dataType))
	fetchedData, err := db.Query(nil, selector, nil, defaultQueryResultsLimit, nil, nil)
	if err != nil {
		logger.Log.Debugf("Error retrieving all %ss: %s", dataTypeStrForLogging, err.Error())
		return nil, err
	}

	// Strip out the prefix on all the IDs
	for _, data := range fetchedData {
		stripPrefixFromID(data)
	}

	return fetchedData, nil
}

func printStringSliceInObjectDeclarationFormat(stringSlice []string, dataType string) string {
	if stringSlice == nil || len(stringSlice) == 0 {
		return "[]string{}"
	}

	datatypeString := fmt.Sprintf(`%s_2_`, dataType)
	result := "[]string{"
	for _, s := range stringSlice {
		result += fmt.Sprintf(`"%s%s",`, datatypeString, s)
	}
	result = strings.TrimRight(result, ",")
	return strings.Join([]string{result, `}`}, "")
}

func stripPrefixFromID(data map[string]interface{}) {
	if data["_id"] != nil {
		idStr := data["_id"].(string)
		if len(idStr) != 0 {
			data["_id"] = ds.GetDataIDFromFullID(idStr)
		}
	}
}

// ConvertGenericCouchDataToObject - takes an empty object of a known type and populates
// that object with the generic data.
func convertGenericCouchDataToObject(genericData map[string]interface{}, dataContainer interface{}, dataTypeStr string) error {
	// Flatten the object for conversion
	flattened := genericData["data"].(map[string]interface{})
	flattened["_id"] = genericData["_id"]
	flattened["_rev"] = genericData["_rev"]

	genericDataInBytes, err := convertGenericObjectToBytesWithCouchDbFields(flattened)
	if err != nil {
		return err
	}

	err = json.Unmarshal(genericDataInBytes, &dataContainer)
	if err != nil {
		logger.Log.Debugf("Error converting generic data to %s type: %s", dataTypeStr, err.Error())
		return err
	}

	logger.Log.Debugf("Converted generic data to %s: %v\n", dataTypeStr, models.AsJSONString(dataContainer))

	return nil
}

func convertCouchDataArrayToFlattenedArray(genericData []map[string]interface{}, dataContainer interface{}, dataTypeStr string) error {
	flattenedContainer := []map[string]interface{}{}
	for _, obj := range genericData {
		flattened := obj["data"].(map[string]interface{})
		flattened["_id"] = obj["_id"]
		flattened["_rev"] = obj["_rev"]
		flattenedContainer = append(flattenedContainer, flattened)
	}

	flattenedDataInBytes, err := convertGenericObjectToBytes(flattenedContainer)
	if err != nil {
		return err
	}

	err = json.Unmarshal(flattenedDataInBytes, &dataContainer)
	if err != nil {
		logger.Log.Debugf("Error converting generic data to %s type: %s", dataTypeStr, err.Error())
		return err
	}
	logger.Log.Debugf("Converted generic data to %s: %v\n", dataTypeStr, models.AsJSONString(dataContainer))

	return nil
}

func convertGenericArrayToObject(genericData []map[string]interface{}, dataContainer interface{}, dataTypeStr string) error {
	genericDataInBytes, err := convertGenericObjectToBytes(genericData)
	if err != nil {
		return err
	}

	err = json.Unmarshal(genericDataInBytes, &dataContainer)
	if err != nil {
		logger.Log.Debugf("Error converting generic data to %s type: %s", dataTypeStr, err.Error())
		return err
	}

	logger.Log.Debugf("Converted generic data to %s: %v\n", dataTypeStr, models.AsJSONString(dataContainer))

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

// storeData - encapsulates logic required for basic data storage for objects that follow the basic data format.
func storeData(dbName string, data interface{}, dataType string, dataTypeLogStr string, dataContainer interface{}) error {
	db, err := getDatabase(dbName)
	if err != nil {
		return err
	}

	// Convert to generic object
	storeFormat, err := convertDataToCouchDbSupportedModel(data)
	if err != nil {
		return err
	}

	// Give the data a known type, and timestamps:
	objectData := storeFormat["data"].(map[string]interface{})
	objectData["datatype"] = dataType
	objectData["createdTimestamp"] = ds.MakeTimestamp()
	objectData["lastModifiedTimestamp"] = objectData["createdTimestamp"]

	// Store the object in CouchDB
	_, _, err = storeDataInCouchDB(storeFormat, dataTypeLogStr, db)
	if err != nil {
		return err
	}

	stripPrefixFromID(storeFormat)

	// Populate the response
	if err = convertGenericCouchDataToObject(storeFormat, &dataContainer, dataTypeLogStr); err != nil {
		return err
	}

	// Return the provisioned object.
	logger.Log.Debugf("Created %s: %v\n", dataTypeLogStr, models.AsJSONString(dataContainer))
	return nil
}

// updateCouchDBDocWithStringDoc - Updates a couchdb design document in coucdb.
func updateCouchDBDocWithStringDoc(dbName string, data string, dataType string, dataTypeLogStr string, dataContainer interface{}) error {
	db, err := getDatabase(dbName)
	if err != nil {
		return err
	}

	genericFormat := make(map[string]interface{})

	err = json.Unmarshal([]byte(data), &genericFormat)
	if err != nil {
		return err
	}

	// Store the object in CouchDB
	_, _, err = storeDataInCouchDB(genericFormat, dataTypeLogStr, db)
	if err != nil {
		return err
	}

	// Populate the response
	if err = convertCouchDesignDocumentToObject(genericFormat, &dataContainer, dataTypeLogStr); err != nil {
		return err
	}

	// Return the provisioned object.
	if logger.IsDebugEnabled() {
		logger.Log.Debugf("Updated %s: %+v\n", dataTypeLogStr, models.AsJSONString(dataContainer))
	}
	return nil
}

// updateDesignDoc - encapsulates logic required for basic data updates for objects that follow the basic data format.
func updateDesignDoc(dbName string, data interface{}, dataType string, dataTypeLogStr string, dataContainer interface{}) error {
	db, err := getDatabase(dbName)
	if err != nil {
		return err
	}

	genericFormat := make(map[string]interface{})

	var dataToBytes []byte

	dataToBytes, err = json.Marshal(data)
	if err != nil {
		return err
	}
	err = json.Unmarshal(dataToBytes, &genericFormat)
	if err != nil {
		return err
	}

	// Store the object in CouchDB
	_, _, err = storeDataInCouchDB(genericFormat, dataTypeLogStr, db)
	if err != nil {
		return err
	}

	// Populate the response
	if err = convertCouchDesignDocumentToObject(genericFormat, &dataContainer, dataTypeLogStr); err != nil {
		return err
	}

	// Return the provisioned object.
	if logger.IsDebugEnabled() {
		logger.Log.Debugf("Updated %s: %+v\n", dataTypeLogStr, models.AsJSONString(dataContainer))
	}
	return nil
}

// updateData - encapsulates logic required for basic data updates for objects that follow the basic data format.
func updateData(dbName string, data interface{}, dataType string, dataTypeLogStr string, dataContainer interface{}) error {
	db, err := getDatabase(dbName)
	if err != nil {
		return err
	}

	// Convert to generic object
	storeFormat, err := convertDataToCouchDbSupportedModel(data)
	if err != nil {
		return err
	}

	// Give the data a known type, and timestamps:
	objectData := storeFormat["data"].(map[string]interface{})
	objectData["datatype"] = dataType
	objectData["lastModifiedTimestamp"] = ds.MakeTimestamp()

	// Store the object in CouchDB
	_, _, err = storeDataInCouchDB(storeFormat, dataTypeLogStr, db)
	if err != nil {
		return err
	}

	stripPrefixFromID(storeFormat)

	// Populate the response
	if err = convertGenericCouchDataToObject(storeFormat, &dataContainer, dataTypeLogStr); err != nil {
		return err
	}

	// Return the provisioned object.
	logger.Log.Debugf("Updated %s: %v\n", dataTypeLogStr, models.AsJSONString(dataContainer))
	return nil
}

// ConvertGenericCouchDataToObject - takes an empty object of a known type and populates
// that object with the generic data.
func convertCouchDesignDocumentToObject(genericData map[string]interface{}, dataContainer interface{}, dataTypeStr string) error {

	genericDataInBytes, err := convertGenericObjectToBytesWithCouchDbFields(genericData)
	if err != nil {
		return err
	}

	err = json.Unmarshal(genericDataInBytes, &dataContainer)
	if err != nil {
		logger.Log.Debugf("Error converting generic data to %s type: %s", dataTypeStr, err.Error())
		return err
	}

	logger.Log.Debugf("Converted generic data to %s: %v\n", dataTypeStr, models.AsJSONString(dataContainer))

	return nil
}

// getData - encapsulates logic required for basic data retrieval for objects that follow the basic data format.
func getData(dbName string, dataID string, dataTypeLogStr string, dataContainer interface{}) error {
	db, err := getDatabase(dbName)
	if err != nil {
		return err
	}

	// Retrieve the object data from CouchDB
	fetchedObject, err := getByDocID(dataID, dataTypeLogStr, db)
	if err != nil {
		return err
	}

	if dataTypeLogStr == tenmod.TenantMonitoredObjectKeysStr {
		if err = convertCouchDesignDocumentToObject(fetchedObject, &dataContainer, dataTypeLogStr); err != nil {
			return err
		}
	} else {

		// Strip prefix from the ID
		stripPrefixFromID(fetchedObject)

		// Marshal the response from the datastore to bytes so that it
		// can be Marshalled back to the proper type.
		if err = convertGenericCouchDataToObject(fetchedObject, &dataContainer, dataTypeLogStr); err != nil {
			return err
		}
	}

	logger.Log.Debugf("Retrieved %s: %v\n", dataTypeLogStr, models.AsJSONString(dataContainer))
	return nil
}

// deleteData - encapsulates logic required for basic data deletion for objects that follow the basic data format.
func deleteData(dbName string, dataID string, dataTypeLogStr string) error {
	// Perform the delete operation on CouchDB
	db, err := getDatabase(dbName)
	if err != nil {
		return err
	}

	if err = deleteByDocID(dataID, dataTypeLogStr, db); err != nil {
		return err
	}

	logger.Log.Debugf("Deleted %s: %s\n", dataTypeLogStr)
	return nil
}

func createDataInCouch(dbName string, dataToStore interface{}, dataContainer interface{}, dataType string, loggingStr string) error {
	logger.Log.Debugf("Creating %s: %v\n", loggingStr, models.AsJSONString(dataToStore))

	if err := storeData(dbName, dataToStore, dataType, loggingStr, &dataContainer); err != nil {
		return err
	}

	// Return the provisioned object.
	logger.Log.Debugf("Created %s: %v\n", loggingStr, models.AsJSONString(dataContainer))
	return nil
}

// Retrieve IDs from a particular view based on a key criteria
func getIDsByView(dbName string, designDocName string, viewName string, key string) ([]string, error) {

	db, err := getDatabase(dbName)
	view := createDBPathStr("_design", designDocName, "_view", viewName)

	qp := url.Values{}
	qp.Set("key", fmt.Sprintf("[\"%s\"]", key))

	vr, err := db.Get(view, qp)

	if err != nil {
		return nil, err
	}

	// If no rows are returned then immediately return an empty list
	moList := make([]string, 0)
	rows, found := vr["rows"]
	if !found {
		return moList, nil
	}

	for _, r := range rows.([]interface{}) {
		rMap := r.(map[string]interface{})
		moList = append(moList, rMap["id"].(string))
	}

	return moList, nil
}

func updateDataInCouch(dbName string, dataToStore interface{}, dataContainer interface{}, dataType string, loggingStr string) error {
	logger.Log.Debugf("Updating %s: %v\n", loggingStr, models.AsJSONString(dataToStore))

	if err := updateData(dbName, dataToStore, dataType, loggingStr, &dataContainer); err != nil {
		return err
	}

	// Return the provisioned object.
	logger.Log.Debugf("Updated %s: %v\n", loggingStr, models.AsJSONString(dataContainer))
	return nil
}

func getDesignDocumentFromCouch(dbName string, idToRetrieve string, dataContainer interface{}, loggingStr string) error {
	logger.Log.Debugf("Retrieving %s for %s\n", loggingStr, idToRetrieve)

	if err := getData(dbName, idToRetrieve, loggingStr, &dataContainer); err != nil {
		return err
	}

	// Return the provisioned object.
	logger.Log.Debugf("Retrieved %s: %v\n", loggingStr, models.AsJSONString(dataContainer))
	return nil
}
func getDataFromCouch(dbName string, idToRetrieve string, dataContainer interface{}, loggingStr string) error {
	logger.Log.Debugf("Retrieving %s for %s\n", loggingStr, idToRetrieve)

	if err := getData(dbName, idToRetrieve, loggingStr, &dataContainer); err != nil {
		return err
	}

	// Return the provisioned object.
	logger.Log.Debugf("Retrieved %s: %v\n", loggingStr, models.AsJSONString(dataContainer))
	return nil
}

func deleteDataFromCouch(dbName string, idToDelete string, dataContainer interface{}, loggingStr string) error {
	logger.Log.Debugf("Deleting %s for %s\n", loggingStr, idToDelete)

	// Obtain the value of the existing record for a return value.
	if err := getDataFromCouch(dbName, idToDelete, &dataContainer, loggingStr); err != nil {
		logger.Log.Debugf("Unable to fetch %s to delete: %s", loggingStr, err.Error())
		return err
	}

	if err := deleteData(dbName, idToDelete, loggingStr); err != nil {
		logger.Log.Debugf("Unable to delete %s: %s", loggingStr, err.Error())
		return err
	}

	// Return the deleted object.
	logger.Log.Debugf("Deleted %s: %v\n", loggingStr, models.AsJSONString(dataContainer))
	return nil
}

func getAllOfTypeFromCouch(dbName string, dataType string, loggingStr string, dataContainer interface{}) error {
	db, err := getDatabase(dbName)
	if err != nil {
		return err
	}

	fetchedList, err := getAllOfTypeByIDPrefix(dataType, loggingStr, db)
	if err != nil {
		return err
	}

	// Marshal the response from the datastore to bytes so that it
	// can be Marshalled back to the proper type.
	return convertGenericArrayToObject(fetchedList, &dataContainer, loggingStr)
}

func getAllOfTypeFromCouchAndFlatten(dbName string, dataType string, loggingStr string, dataContainer interface{}) error {
	db, err := getDatabase(dbName)
	if err != nil {
		return err
	}

	fetchedList, err := getAllOfTypeByIDPrefix(dataType, loggingStr, db)
	if err != nil {
		return err
	}

	// Marshal the response from the datastore to bytes so that it
	// can be Marshalled back to the proper type.
	return convertCouchDataArrayToFlattenedArray(fetchedList, dataContainer, loggingStr)
}

func getAllInIDListFromCouchAndFlatten(dbName string, idList []string, dataType string, loggingStr string, dataContainer interface{}) error {
	db, err := getDatabase(dbName)
	if err != nil {
		return err
	}

	fetchedList, err := getAllInListOfIDs(idList, dataType, loggingStr, db)
	if err != nil {
		return err
	}

	// Marshal the response from the datastore to bytes so that it
	// can be Marshalled back to the proper type.
	return convertCouchDataArrayToFlattenedArray(fetchedList, dataContainer, loggingStr)
}

// GenerateMonitoredObjectURL - Generates a Monitored Object URL
func GenerateMonitoredObjectURL(tenantID string, uri string) string {
	dbName := createDBPathStr(uri, fmt.Sprintf("tenant_2_%s%s/", tenantID, monitoredObjectDBSuffix))
	return dbName
}

// updateTenantMetadataMetadata - Updates the metadata in the TenantMetadata object' metakeys
func updateTenantMetadataMetadata(meta map[string]string, tenantMeta *tenmod.Metadata) ([]string, error) {
	keys := make([]string, 0)
	// Go thru a list of KV pairs and add the keys to the Metadata.
	// The idea is to cache all the known monitored  Metadata keys so the UI can do word completion
	for key, _ := range meta {

		if len(tenantMeta.MonitorObjectMetaKeys[key]) == 0 {
			logger.Log.Debugf("\"%s\" is a new key", key)
			keys = append(keys, key)
			// Stop being meta, Ahbed
			tenantMeta.MonitorObjectMetaKeys[key] = key
		}
	}

	return keys, nil
}
