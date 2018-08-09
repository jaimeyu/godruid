package couchDB

import (
	"errors"
	"fmt"
	"strings"

	"github.com/getlantern/deepcopy"
	"github.com/leesper/couchdb-golang"

	"github.com/accedian/adh-gather/config"
	ds "github.com/accedian/adh-gather/datastore"
	"github.com/accedian/adh-gather/gather"
	"github.com/accedian/adh-gather/logger"

	"github.com/accedian/adh-gather/models"
	admmod "github.com/accedian/adh-gather/models/admin"
)

// AdminServiceDatastoreCouchDB - struct responsible for handling
// database operations for the Admin Service when using CouchDB
// as the storage option.
type AdminServiceDatastoreCouchDB struct {
	couchHost   string
	dbName      string
	dbNameAlone string
	server      *couchdb.Server
	cfg         config.Provider
}

// CreateAdminServiceDAO - instantiates a CouchDB implementation of the
// AdminServiceDatastore.
func CreateAdminServiceDAO() (*AdminServiceDatastoreCouchDB, error) {
	result := new(AdminServiceDatastoreCouchDB)
	result.cfg = gather.GetConfig()

	// Couch Server Configuration
	provDBURL := fmt.Sprintf("%s:%d",
		result.cfg.GetString(gather.CK_server_datastore_ip.String()),
		result.cfg.GetInt(gather.CK_server_datastore_port.String()))
	logger.Log.Debugf("Admin Service CouchDB URL is: %s", provDBURL)
	result.couchHost = provDBURL

	// Couch DB name configuration
	result.dbNameAlone = result.cfg.GetString(gather.CK_args_admindb_name.String())
	result.dbName = result.couchHost + "/" + result.dbNameAlone
	server, err := couchdb.NewServer(result.couchHost)
	if err != nil {
		logger.Log.Debugf("Falied to instantiate AdminServiceDatastoreCouchDB: %s", err.Error())
		return nil, err
	}
	logger.Log.Debugf("Admin Database is: %s", result.dbName)

	result.server = server
	return result, nil
}

// CreateAdminUser - CouchDB implementation of CreateAdminUser
func (asd *AdminServiceDatastoreCouchDB) CreateAdminUser(user *admmod.User) (*admmod.User, error) {
	user.ID = ds.GenerateID(user, string(admmod.AdminUserType))

	dataContainer := &admmod.User{}
	if err := createDataInCouch(asd.dbName, user, dataContainer, string(admmod.AdminUserType), admmod.AdminUserStr); err != nil {
		return nil, err
	}

	return dataContainer, nil
}

// UpdateAdminUser - CouchDB implementation of UpdateAdminUser
func (asd *AdminServiceDatastoreCouchDB) UpdateAdminUser(user *admmod.User) (*admmod.User, error) {
	user.ID = ds.PrependToDataID(user.ID, string(admmod.AdminUserType))

	dataContainer := &admmod.User{}
	if err := updateDataInCouch(asd.dbName, user, dataContainer, string(admmod.AdminUserType), admmod.AdminUserStr); err != nil {
		return nil, err
	}

	return dataContainer, nil
}

// DeleteAdminUser - CouchDB implementation of DeleteAdminUser
func (asd *AdminServiceDatastoreCouchDB) DeleteAdminUser(userID string) (*admmod.User, error) {
	userID = ds.PrependToDataID(userID, string(admmod.AdminUserType))

	dataContainer := admmod.User{}
	if err := deleteDataFromCouch(asd.dbName, userID, &dataContainer, admmod.AdminUserStr); err != nil {
		return nil, err
	}

	dataContainer.ID = ds.GetDataIDFromFullID(dataContainer.ID)
	return &dataContainer, nil
}

// GetAdminUser - CouchDB implementation of GetAdminUser
func (asd *AdminServiceDatastoreCouchDB) GetAdminUser(userID string) (*admmod.User, error) {
	userID = ds.PrependToDataID(userID, string(admmod.AdminUserType))

	dataContainer := admmod.User{}
	if err := getDataFromCouch(asd.dbName, userID, &dataContainer, admmod.AdminUserStr); err != nil {
		return nil, err
	}

	dataContainer.ID = ds.GetDataIDFromFullID(dataContainer.ID)
	return &dataContainer, nil
}

// GetAllAdminUsers - CouchDB implementation of GetAllAdminUsers
func (asd *AdminServiceDatastoreCouchDB) GetAllAdminUsers() ([]*admmod.User, error) {
	res := []*admmod.User{}
	if err := getAllOfTypeFromCouchAndFlatten(asd.dbName, string(admmod.AdminUserType), admmod.AdminUserStr, &res); err != nil {
		return nil, err
	}

	return res, nil
}

// CreateTenant - CouchDB implementation of CreateTenant
func (asd *AdminServiceDatastoreCouchDB) CreateTenant(tenantDescriptor *admmod.Tenant) (*admmod.Tenant, error) {
	logger.Log.Debugf("Creating %s: %v\n", admmod.TenantStr, models.AsJSONString(tenantDescriptor))
	tenantDescriptor.ID = ds.GenerateID(tenantDescriptor, string(admmod.TenantType))

	dataContainer := &admmod.Tenant{}
	if err := createDataInCouch(asd.dbName, tenantDescriptor, dataContainer, string(admmod.TenantType), admmod.TenantStr); err != nil {
		return nil, err
	}

	// Create a CouchDB database to isolate the tenant data
	_, err := asd.CreateDatabase(tenantDescriptor.ID)
	if err != nil {
		logger.Log.Debugf("Unable to create database for Tenant %s: %s", tenantDescriptor.ID, err.Error())
		return nil, err
	}

	// Create a CouchDB database to isolate monitored object data for the tenant
	_, err = asd.CreateDatabase(fmt.Sprintf("%s%s", tenantDescriptor.ID, monitoredObjectDBSuffix))
	if err != nil {
		logger.Log.Debugf("Unable to create monitored object database for Tenant %s: %s", tenantDescriptor.ID, err.Error())
		return nil, err
	}

	// Create a CouchDB database to isolate reports for the tenant
	_, err = asd.CreateDatabase(fmt.Sprintf("%s%s", tenantDescriptor.ID, reportObjectDBSuffix))
	if err != nil {
		logger.Log.Debugf("Unable to create report database for Tenant %s: %s", tenantDescriptor.ID, err.Error())
		return nil, err
	}

	// Add in the views/indicies necessary for the db:
	if err = asd.addTenantViewsToDB(tenantDescriptor.ID); err != nil {
		logger.Log.Debugf("Unable to add Views to DB for Tenant %s: %s", tenantDescriptor.ID, err.Error())
		return nil, err
	}

	// Return the provisioned object.
	logger.Log.Debugf("Created %s: %v\n", admmod.TenantStr, models.AsJSONString(dataContainer))
	return dataContainer, nil
}

// UpdateTenantDescriptor - CouchDB implementation of UpdateTenantDescriptor
func (asd *AdminServiceDatastoreCouchDB) UpdateTenantDescriptor(tenantDescriptor *admmod.Tenant) (*admmod.Tenant, error) {
	logger.Log.Debugf("Updating %s: %v\n", admmod.TenantStr, models.AsJSONString(tenantDescriptor))
	tenantDescriptor.ID = ds.PrependToDataID(tenantDescriptor.ID, string(admmod.TenantType))

	dataContainer := &admmod.Tenant{}
	if err := updateDataInCouch(asd.dbName, tenantDescriptor, dataContainer, string(admmod.TenantType), admmod.TenantStr); err != nil {
		return nil, err
	}

	logger.Log.Debugf("Updated %s: %v\n", admmod.TenantStr, models.AsJSONString(dataContainer))
	return dataContainer, nil
}

// DeleteTenant - CouchDB implementation of DeleteTenant
func (asd *AdminServiceDatastoreCouchDB) DeleteTenant(tenantID string) (*admmod.Tenant, error) {
	logger.Log.Debugf("Deleting %s: %s\n", admmod.TenantStr, tenantID)
	tenantIDWithPrefix := ds.PrependToDataID(tenantID, string(admmod.TenantType))

	// Obtain the value of the existing record for a return value.
	existingTenant, err := asd.GetTenantDescriptor(tenantID)
	if err != nil {
		logger.Log.Debugf("Unable to delete %s: %s", admmod.TenantStr, err.Error())
		return nil, err
	}

	// Setup deletion of the report datastore
	reportDBID := fmt.Sprintf("%s%s", tenantIDWithPrefix, reportObjectDBSuffix)
	if err = purgeDB(createDBPathStr(asd.couchHost, reportDBID)); err != nil {
		logger.Log.Debugf("Unable to purge DB contents for %s %s: %s", admmod.TenantStr, "reports", err.Error())
		return nil, err
	}
	err = asd.DeleteDatabase(reportDBID)
	if err != nil {
		logger.Log.Debugf("Unable to delete report database for Tenant %s: %s", tenantID, err.Error())
		return nil, err
	}

	// Setup deletion of the monitored object datastore
	moDBID := fmt.Sprintf("%s%s", tenantIDWithPrefix, monitoredObjectDBSuffix)
	if err = purgeDB(createDBPathStr(asd.couchHost, moDBID)); err != nil {
		logger.Log.Debugf("Unable to purge DB contents for %s %s: %s", admmod.TenantStr, "monitored objects", err.Error())
		return nil, err
	}
	err = asd.DeleteDatabase(moDBID)
	if err != nil {
		logger.Log.Debugf("Unable to delete monitored object database for Tenant %s: %s", tenantID, err.Error())
		return nil, err
	}

	// Purge the DB of records:
	if err = purgeDB(createDBPathStr(asd.couchHost, tenantIDWithPrefix)); err != nil {
		logger.Log.Debugf("Unable to purge DB contents for %s: %s", admmod.TenantStr, err.Error())
		return nil, err
	}

	// Try to delete the DB for the tenant
	if err := asd.DeleteDatabase(tenantIDWithPrefix); err != nil {
		logger.Log.Debugf("Unable to delete %s: %s", admmod.TenantStr, err.Error())
		return nil, err
	}

	if err = deleteData(asd.dbName, tenantIDWithPrefix, admmod.TenantStr); err != nil {
		logger.Log.Debugf("Unable to delete %s: %s", admmod.TenantStr, err.Error())
		return nil, err
	}

	// Return the deleted object.
	logger.Log.Debugf("Deleted %s: %v\n", admmod.TenantStr, models.AsJSONString(existingTenant))
	return existingTenant, nil
}

// GetTenantDescriptor - CouchDB implementation of GetTenantDescriptor
func (asd *AdminServiceDatastoreCouchDB) GetTenantDescriptor(tenantID string) (*admmod.Tenant, error) {
	logger.Log.Debugf("Fetching %s: %s\n", admmod.TenantStr, tenantID)
	tenantID = ds.PrependToDataID(tenantID, string(admmod.TenantType))

	dataContainer := admmod.Tenant{}
	if err := getDataFromCouch(asd.dbName, tenantID, &dataContainer, admmod.TenantStr); err != nil {
		return nil, err
	}

	logger.Log.Debugf("Retrieved %s: %v\n", admmod.TenantStr, models.AsJSONString(dataContainer))
	return &dataContainer, nil
}

// GetAllTenantDescriptors - CouchDB implementation of GetAllTenantDescriptors
func (asd *AdminServiceDatastoreCouchDB) GetAllTenantDescriptors() ([]*admmod.Tenant, error) {
	logger.Log.Debugf("Fetching all %s\n", admmod.TenantStr)
	res := []*admmod.Tenant{}
	if err := getAllOfTypeFromCouchAndFlatten(asd.dbName, string(admmod.TenantType), admmod.TenantStr, &res); err != nil {
		return nil, err
	}

	logger.Log.Debugf("Retrieved %d %s\n", len(res), admmod.TenantStr)
	return res, nil
}

// CreateDatabase - creates a database in CouchDB identified by the provided name.
func (asd *AdminServiceDatastoreCouchDB) CreateDatabase(dbName string) (ds.Database, error) {
	if len(dbName) == 0 {
		return nil, errors.New("Unable to create database if no identifier is provided")
	}
	if asd.server.Contains(dbName) {
		return nil, errors.New("Unable to create database '" + dbName + "': database already exists")
	}

	fmt.Printf("Server is: %v\n", asd.server)
	db, err := asd.server.Create(dbName)
	if err != nil {
		return nil, err
	}

	logger.Log.Debugf("Created DB %s\n", dbName)

	return db, nil
}

func (asd *AdminServiceDatastoreCouchDB) addTenantViewsToDB(dbName string) error {
	if len(dbName) == 0 {
		return errors.New("Unable to add views to a database if no database name is provided")
	}
	if !asd.server.Contains(dbName) {
		return errors.New("Unable to add views to database '" + dbName + "': database does not exist")
	}

	logger.Log.Debugf("Adding Tenant Views to DB %s", dbName)

	moDB, err := getDatabase(createDBPathStr(asd.couchHost, fmt.Sprintf("%s%s", dbName, monitoredObjectDBSuffix)))
	if err != nil {
		return err
	}

	// Store the views related to Monitored Objects
	for _, viewPayload := range getTenantViews() {
		_, _, err = storeDataInCouchDBWithQueryParams(viewPayload, "TenantView", moDB, nil)
		if err != nil {
			return err
		}
	}

	logger.Log.Debugf("Added views to DB %s\n", dbName)
	return nil
}

// DeleteDatabase - deletes a database in CouchDB identified by the provided name.
func (asd *AdminServiceDatastoreCouchDB) DeleteDatabase(dbName string) error {
	if len(dbName) == 0 {
		logger.Log.Debug("No database identifier provided, nothing to delete")
		return nil
	}
	if !asd.server.Contains(dbName) {
		logger.Log.Debugf("Unable to delete database '" + dbName + "': database does not exist")
		return nil
	}

	logger.Log.Debugf("Deleted DB %s\n", dbName)
	return asd.server.Delete(dbName)
}

// CreateIngestionDictionary - CouchDB implementation of CreateIngestionDictionary
func (asd *AdminServiceDatastoreCouchDB) CreateIngestionDictionary(ingDictionary *admmod.IngestionDictionary) (*admmod.IngestionDictionary, error) {
	logger.Log.Debugf("Creating %s: %v\n", admmod.IngestionDictionaryStr, models.AsJSONString(ingDictionary))
	// Only create one if one does not already exist:
	existing, _ := asd.GetIngestionDictionary()
	if existing != nil {
		return nil, fmt.Errorf("Can't create %s, it already exists", admmod.IngestionDictionaryStr)
	}

	// No pre-existing dictionary, go ahead and create one.
	ingDictionary.ID = ds.GenerateID(ingDictionary, string(admmod.IngestionDictionaryType))

	dataType := string(admmod.IngestionDictionaryType)
	dataContainer := admmod.IngestionDictionary{}
	if err := storeData(asd.dbName, ingDictionary, dataType, admmod.IngestionDictionaryStr, &dataContainer); err != nil {
		return nil, err
	}

	// Return the provisioned object.
	logger.Log.Debugf("Created %s: %v\n", admmod.IngestionDictionaryStr, models.AsJSONString(dataContainer))
	return &dataContainer, nil
}

// UpdateIngestionDictionary - CouchDB implementation of UpdateIngestionDictionary
func (asd *AdminServiceDatastoreCouchDB) UpdateIngestionDictionary(ingDictionary *admmod.IngestionDictionary) (*admmod.IngestionDictionary, error) {
	logger.Log.Debugf("Updating %s: %v\n", admmod.IngestionDictionaryStr, models.AsJSONString(ingDictionary))
	ingDictionary.ID = ds.PrependToDataID(ingDictionary.ID, string(admmod.IngestionDictionaryType))

	dataType := string(admmod.IngestionDictionaryType)
	dataContainer := admmod.IngestionDictionary{}
	if err := updateData(asd.dbName, ingDictionary, dataType, admmod.IngestionDictionaryStr, &dataContainer); err != nil {
		return nil, err
	}

	// Return the provisioned object.
	logger.Log.Debugf("Updated %s: %v\n", admmod.IngestionDictionaryStr, models.AsJSONString(dataContainer))
	return &dataContainer, nil
}

// DeleteIngestionDictionary - CouchDB implementation of DeleteIngestionDictionary
func (asd *AdminServiceDatastoreCouchDB) DeleteIngestionDictionary() (*admmod.IngestionDictionary, error) {
	logger.Log.Debugf("Deleting %s\n", admmod.IngestionDictionaryStr)
	// Obtain the value of the existing record for a return value.
	existingDictionary, err := asd.GetIngestionDictionary()
	if err != nil {
		logger.Log.Debugf("Unable to delete %s: %s", admmod.IngestionDictionaryStr, err.Error())
		return nil, err
	}

	deleteID := ds.PrependToDataID(existingDictionary.ID, string(admmod.IngestionDictionaryType))
	if err = deleteData(asd.dbName, deleteID, admmod.IngestionDictionaryStr); err != nil {
		logger.Log.Debugf("Unable to delete %s: %s", admmod.IngestionDictionaryStr, err.Error())
		return nil, err
	}

	// Return the deleted object.
	logger.Log.Debugf("Deleted %s: %v\n", admmod.IngestionDictionaryStr, models.AsJSONString(existingDictionary))
	return existingDictionary, nil

}

// GetIngestionDictionary - CouchDB implementation of GetIngestionDictionary
func (asd *AdminServiceDatastoreCouchDB) GetIngestionDictionary() (*admmod.IngestionDictionary, error) {
	logger.Log.Debugf("Retrieving %s\n", admmod.IngestionDictionaryStr)
	db, err := getDatabase(asd.dbName)
	if err != nil {
		return nil, err
	}

	fetchedData, err := getAllOfTypeByIDPrefix(string(admmod.IngestionDictionaryType), admmod.IngestionDictionaryStr, db)
	if err != nil {
		return nil, err
	}

	// Populate the response
	res := admmod.IngestionDictionary{}
	if len(fetchedData) != 0 {
		if err = convertGenericCouchDataToObject(fetchedData[0], &res, admmod.IngestionDictionaryStr); err != nil {
			return nil, err
		}
	} else {
		return nil, fmt.Errorf("Unable to find %s", admmod.IngestionDictionaryStr)
	}

	logger.Log.Debugf("Retrieved %s: %v\n", admmod.IngestionDictionaryStr, models.AsJSONString(res))
	return &res, nil
}

// GetTenantIDByAlias - InMemory impl of GetTenantIDByAlias
func (asd *AdminServiceDatastoreCouchDB) GetTenantIDByAlias(name string) (string, error) {
	logger.Log.Debugf("Getting Tenant ID for Tenant %s\n", name)
	// Retrieve just the subset of values.
	requestBody := map[string]interface{}{}
	requestBody["keys"] = []string{strings.ToLower(name)}

	fetchResponse, err := fetchDesignDocumentResults(requestBody, asd.dbName, tenantIDByNameIndex)
	if err != nil {
		return "", err
	}

	rows := fetchResponse["rows"].([]interface{})
	if rows == nil || len(rows) == 0 {
		return "", nil
	}
	obj := rows[0].(map[string]interface{})

	response := ds.GetDataIDFromFullID(obj["value"].(string))
	logger.Log.Debugf("Returning Tenant ID: %vs\n")
	return response, nil
}

// AddAdminViews - Adds the admin views (indicies) to the Admin DB.
func (asd *AdminServiceDatastoreCouchDB) AddAdminViews() error {

	logger.Log.Debugf("Adding Admin Views to DB %s", asd.dbNameAlone)

	db, err := getDatabase(asd.dbName)
	if err != nil {
		return err
	}

	// Store the sync checkpoint in CouchDB
	for _, viewPayload := range generateAdminViews() {
		_, _, err = storeDataInCouchDBWithQueryParams(viewPayload, "AdminView", db, nil)
		if err != nil {
			return err
		}
	}

	logger.Log.Debugf("Added views to DB %s\n", asd.dbNameAlone)
	return nil
}

// CreateValidTypes - CouchDB implementation of CreateValidTypes
func (asd *AdminServiceDatastoreCouchDB) CreateValidTypes(value *admmod.ValidTypes) (*admmod.ValidTypes, error) {
	logger.Log.Debugf("Creating %s: %v\n", admmod.ValidTypesStr, models.AsJSONString(value))
	value.ID = ds.GenerateID(value, string(admmod.ValidTypesType))

	// Only create one if one does not already exist:
	existing, _ := asd.GetValidTypes()
	if existing != nil {
		return nil, fmt.Errorf("Can't create %s, it already exists", admmod.ValidTypesStr)
	}

	dataType := string(admmod.ValidTypesType)
	dataContainer := admmod.ValidTypes{}
	if err := storeData(asd.dbName, value, dataType, admmod.ValidTypesStr, &dataContainer); err != nil {
		return nil, err
	}

	// Return the provisioned object.
	logger.Log.Debugf("Created %s: %v\n", admmod.ValidTypesStr, models.AsJSONString(dataContainer))
	return &dataContainer, nil
}

// UpdateValidTypes - CouchDB implementation of UpdateValidTypes
func (asd *AdminServiceDatastoreCouchDB) UpdateValidTypes(value *admmod.ValidTypes) (*admmod.ValidTypes, error) {
	logger.Log.Debugf("Updating %s: %v\n", admmod.ValidTypesStr, models.AsJSONString(value))
	value.ID = ds.PrependToDataID(value.ID, string(admmod.ValidTypesType))

	dataType := string(admmod.ValidTypesType)
	dataContainer := admmod.ValidTypes{}
	if err := updateData(asd.dbName, value, dataType, admmod.ValidTypesStr, &dataContainer); err != nil {
		return nil, err
	}

	// Return the provisioned object.
	logger.Log.Debugf("Updated %s: %v\n", admmod.ValidTypesStr, models.AsJSONString(dataContainer))
	return &dataContainer, nil
}

// GetValidTypes - CouchDB implementation of GetValidTypes
func (asd *AdminServiceDatastoreCouchDB) GetValidTypes() (*admmod.ValidTypes, error) {
	logger.Log.Debugf("Fetching %s\n", admmod.ValidTypesStr)
	db, err := getDatabase(asd.dbName)
	if err != nil {
		return nil, err
	}

	fetchedData, err := getAllOfTypeByIDPrefix(string(admmod.ValidTypesType), admmod.ValidTypesStr, db)
	if err != nil {
		return nil, err
	}

	// Populate the response
	res := admmod.ValidTypes{}
	if len(fetchedData) != 0 {
		if err = convertGenericCouchDataToObject(fetchedData[0], &res, admmod.ValidTypesStr); err != nil {
			return nil, err
		}
	} else {
		return nil, fmt.Errorf("No %s found", admmod.ValidTypesStr)
	}

	logger.Log.Debugf("Retrieved %s: %v\n", admmod.ValidTypesStr, models.AsJSONString(res))
	return &res, nil
}

// GetSpecificValidTypes - CouchDB implementation of GetSpecificValidTypes
func (asd *AdminServiceDatastoreCouchDB) GetSpecificValidTypes(value *admmod.ValidTypesRequest) (*admmod.ValidTypes, error) {
	logger.Log.Debugf("Fetching %s: %v\n", admmod.ValidTypesStr, models.AsJSONString(value))
	currentValidValuesRecord, err := asd.GetValidTypes()
	if err != nil {
		return nil, err
	}

	result := &admmod.ValidTypes{}
	deepcopy.Copy(&result, currentValidValuesRecord)
	if !value.MonitoredObjectTypes {
		result.MonitoredObjectTypes = nil
	}
	if !value.MonitoredObjectDeviceTypes {
		result.MonitoredObjectDeviceTypes = nil
	}

	logger.Log.Debugf("Retrieved %s: %v\n", admmod.ValidTypesStr, models.AsJSONString(result))
	return result, nil
}

// DeleteValidTypes - CouchDB implementation of DeleteValidTypes
func (asd *AdminServiceDatastoreCouchDB) DeleteValidTypes() (*admmod.ValidTypes, error) {
	logger.Log.Debugf("Deleting %s\n", admmod.ValidTypesStr)
	// Obtain the value of the existing record for a return value.
	existing, err := asd.GetValidTypes()
	if err != nil {
		logger.Log.Debugf("Unable to delete %s: %s", admmod.ValidTypesStr, err.Error())
		return nil, err
	}

	deleteID := ds.PrependToDataID(existing.ID, string(admmod.ValidTypesType))
	if err = deleteData(asd.dbName, deleteID, admmod.ValidTypesStr); err != nil {
		logger.Log.Debugf("Unable to delete %s: %s", admmod.ValidTypesStr, err.Error())
		return nil, err
	}

	// Return the deleted object.
	logger.Log.Debugf("Deleted %s: %v\n", admmod.ValidTypesStr, models.AsJSONString(existing))
	return existing, nil

}

func generateAdminViews() []map[string]interface{} {
	result := make([]map[string]interface{}, 0)

	tenantIDByName := map[string]interface{}{}
	tenantIDByName["_id"] = "_design/tenant"
	tenantIDByName["language"] = "javascript"
	byName := map[string]interface{}{}
	byName["map"] = "function(doc) {\n    if (doc.data && doc.data.datatype && doc.data.datatype === 'tenant') {\n        emit(doc.data.name.toLowerCase(), doc._id);\n    }\n}"
	views := map[string]interface{}{}
	views["byAlias"] = byName
	tenantIDByName["views"] = views

	logger.Log.Debug("Adding view for monitoredObjectCountByDomain")
	return append(result, tenantIDByName)
}
