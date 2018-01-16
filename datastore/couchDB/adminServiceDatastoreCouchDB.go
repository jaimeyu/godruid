package couchDB

import (
	"errors"
	"fmt"
	"strings"

	"github.com/leesper/couchdb-golang"

	"github.com/accedian/adh-gather/config"
	ds "github.com/accedian/adh-gather/datastore"
	"github.com/accedian/adh-gather/gather"
	"github.com/accedian/adh-gather/logger"

	pb "github.com/accedian/adh-gather/gathergrpc"
)

const (
	tenantIDByNameIndex = "_design/tenant/_view/byAlias"
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
	logger.Log.Debug("Admin Service CouchDB URL is: ", provDBURL)
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
func (asd *AdminServiceDatastoreCouchDB) CreateAdminUser(user *pb.AdminUserRequest) (*pb.AdminUserResponse, error) {
	dataContainer := &pb.AdminUserResponse{}
	if err := createDataInCouch(asd.dbName, user, dataContainer, string(ds.AdminUserType), ds.AdminUserStr); err != nil {
		return nil, err
	}
	return dataContainer, nil
}

// UpdateAdminUser - CouchDB implementation of UpdateAdminUser
func (asd *AdminServiceDatastoreCouchDB) UpdateAdminUser(user *pb.AdminUserRequest) (*pb.AdminUserResponse, error) {
	dataContainer := &pb.AdminUserResponse{}
	if err := updateDataInCouch(asd.dbName, user, dataContainer, string(ds.AdminUserType), ds.AdminUserStr); err != nil {
		return nil, err
	}
	return dataContainer, nil
}

// DeleteAdminUser - CouchDB implementation of DeleteAdminUser
func (asd *AdminServiceDatastoreCouchDB) DeleteAdminUser(userID string) (*pb.AdminUserResponse, error) {
	dataContainer := pb.AdminUserResponse{}
	if err := deleteDataFromCouch(asd.dbName, userID, &dataContainer, ds.AdminUserStr); err != nil {
		return nil, err
	}
	return &dataContainer, nil
}

// GetAdminUser - CouchDB implementation of GetAdminUser
func (asd *AdminServiceDatastoreCouchDB) GetAdminUser(userID string) (*pb.AdminUserResponse, error) {
	dataContainer := pb.AdminUserResponse{}
	if err := getDataFromCouch(asd.dbName, userID, &dataContainer, ds.AdminUserStr); err != nil {
		return nil, err
	}
	return &dataContainer, nil
}

// GetAllAdminUsers - CouchDB implementation of GetAllAdminUsers
func (asd *AdminServiceDatastoreCouchDB) GetAllAdminUsers() (*pb.AdminUserListResponse, error) {
	res := &pb.AdminUserListResponse{}
	res.Data = make([]*pb.AdminUserResponse, 0)
	if err := getAllOfTypeFromCouch(asd.dbName, string(ds.AdminUserType), ds.AdminUserStr, &res.Data); err != nil {
		return nil, err
	}

	return res, nil
}

// CreateTenant - CouchDB implementation of CreateTenant
func (asd *AdminServiceDatastoreCouchDB) CreateTenant(tenantDescriptor *pb.TenantDescriptorRequest) (*pb.TenantDescriptorResponse, error) {
	dataContainer := &pb.TenantDescriptorResponse{}
	if err := createDataInCouch(asd.dbName, tenantDescriptor, dataContainer, string(ds.TenantDescriptorType), ds.TenantDescriptorStr); err != nil {
		return nil, err
	}

	// Create a CouchDB database to isolate the tenant data
	_, err := asd.createDatabase(tenantDescriptor.GetXId())
	if err != nil {
		logger.Log.Debugf("Unable to create database for Tenant %s: %s", tenantDescriptor.GetXId(), err.Error())
		return nil, err
	}

	// Add in the views/indicies necessary for the db:
	if err = asd.addTenantViewsToDB(tenantDescriptor.GetXId()); err != nil {
		logger.Log.Debugf("Unable to add Views to DB for Tenant %s: %s", tenantDescriptor.GetXId(), err.Error())
		return nil, err
	}

	// Return the provisioned object.
	return dataContainer, nil
}

// UpdateTenantDescriptor - CouchDB implementation of UpdateTenantDescriptor
func (asd *AdminServiceDatastoreCouchDB) UpdateTenantDescriptor(tenantDescriptor *pb.TenantDescriptorRequest) (*pb.TenantDescriptorResponse, error) {
	dataContainer := &pb.TenantDescriptorResponse{}
	if err := updateDataInCouch(asd.dbName, tenantDescriptor, dataContainer, string(ds.TenantDescriptorType), ds.TenantStr); err != nil {
		return nil, err
	}
	return dataContainer, nil
}

// DeleteTenant - CouchDB implementation of DeleteTenant
func (asd *AdminServiceDatastoreCouchDB) DeleteTenant(tenantID string) (*pb.TenantDescriptorResponse, error) {
	logger.Log.Debugf("Deleting %s for %s\n", ds.TenantStr, tenantID)

	// Obtain the value of the existing record for a return value.
	existingTenant, err := asd.GetTenantDescriptor(tenantID)
	if err != nil {
		logger.Log.Debugf("Unable to delete %s: %s", ds.TenantStr, err.Error())
		return nil, err
	}

	// Purge the DB of records:
	if err = purgeDB(createDBPathStr(asd.couchHost, tenantID)); err != nil {
		logger.Log.Debugf("Unable to purge DB contents for %s: %s", ds.TenantStr, err.Error())
		return nil, err
	}

	// Try to delete the DB for the tenant
	if err := asd.deleteDatabase(tenantID); err != nil {
		logger.Log.Debugf("Unable to delete %s: %s", ds.TenantStr, err.Error())
		return nil, err
	}

	if err = deleteData(asd.dbName, tenantID, ds.TenantStr); err != nil {
		logger.Log.Debugf("Unable to delete %s: %s", ds.TenantStr, err.Error())
		return nil, err
	}

	// Return the deleted object.
	logger.Log.Debugf("Deleted %s: %v\n", ds.TenantStr, existingTenant)
	return existingTenant, nil
}

// GetTenantDescriptor - CouchDB implementation of GetTenantDescriptor
func (asd *AdminServiceDatastoreCouchDB) GetTenantDescriptor(tenantID string) (*pb.TenantDescriptorResponse, error) {
	dataContainer := pb.TenantDescriptorResponse{}
	if err := getDataFromCouch(asd.dbName, tenantID, &dataContainer, ds.TenantDescriptorStr); err != nil {
		return nil, err
	}
	return &dataContainer, nil
}

// GetAllTenantDescriptors - CouchDB implementation of GetAllTenantDescriptors
func (asd *AdminServiceDatastoreCouchDB) GetAllTenantDescriptors() (*pb.TenantDescriptorListResponse, error) {
	res := &pb.TenantDescriptorListResponse{}
	res.Data = make([]*pb.TenantDescriptorResponse, 0)
	if err := getAllOfTypeFromCouch(asd.dbName, string(ds.TenantDescriptorType), ds.TenantDescriptorStr, &res.Data); err != nil {
		return nil, err
	}

	return res, nil
}

// createDatabase - creates a database in CouchDB identified by the provided name.
func (asd *AdminServiceDatastoreCouchDB) createDatabase(dbName string) (*couchdb.Database, error) {
	if len(dbName) == 0 {
		return nil, errors.New("Unable to create database if no identifier is provided")
	}
	if asd.server.Contains(dbName) {
		return nil, errors.New("Unable to create database '" + dbName + "': database already exists")
	}

	logger.Log.Debugf("Created DB %s\n", dbName)
	return asd.server.Create(dbName)
}

func (asd *AdminServiceDatastoreCouchDB) addTenantViewsToDB(dbName string) error {
	if len(dbName) == 0 {
		return errors.New("Unable to add views to a database if no database name is provided")
	}
	if !asd.server.Contains(dbName) {
		return errors.New("Unable to add views to database '" + dbName + "': database does not exist")
	}

	logger.Log.Debugf("Adding Tenant Views to DB %s", dbName)

	db, err := getDatabase(createDBPathStr(asd.couchHost, dbName))
	if err != nil {
		return err
	}

	// Store the sync checkpoint in CouchDB
	for _, viewPayload := range generateTenantViews() {
		_, _, err = storeDataInCouchDBWithQueryParams(viewPayload, "TenantView", db, nil)
		if err != nil {
			return err
		}
	}

	logger.Log.Debugf("Added views to DB %s\n", dbName)
	return nil
}

// deleteDatabase - deletes a database in CouchDB identified by the provided name.
func (asd *AdminServiceDatastoreCouchDB) deleteDatabase(dbName string) error {
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
func (asd *AdminServiceDatastoreCouchDB) CreateIngestionDictionary(ingDictionary *pb.IngestionDictionary) (*pb.IngestionDictionary, error) {
	logger.Log.Debugf("Creating %s: %v\n", ds.IngestionDictionaryStr, ingDictionary)

	dataType := string(ds.IngestionDictionaryType)
	dataContainer := pb.IngestionDictionary{}
	if err := storeData(asd.dbName, ingDictionary, dataType, ds.IngestionDictionaryStr, &dataContainer); err != nil {
		return nil, err
	}

	// Return the provisioned object.
	logger.Log.Debugf("Created %s: %v\n", ds.IngestionDictionaryStr, dataContainer)
	return &dataContainer, nil
}

// UpdateIngestionDictionary - CouchDB implementation of UpdateIngestionDictionary
func (asd *AdminServiceDatastoreCouchDB) UpdateIngestionDictionary(ingDictionary *pb.IngestionDictionary) (*pb.IngestionDictionary, error) {
	logger.Log.Debugf("Updating %s: %v\n", ds.IngestionDictionaryStr, ingDictionary)

	dataType := string(ds.IngestionDictionaryType)
	dataContainer := pb.IngestionDictionary{}
	if err := updateData(asd.dbName, ingDictionary, dataType, ds.IngestionDictionaryStr, &dataContainer); err != nil {
		return nil, err
	}

	// Return the provisioned object.
	logger.Log.Debugf("Updated %s: %v\n", ds.IngestionDictionaryStr, dataContainer)
	return &dataContainer, nil
}

// DeleteIngestionDictionary - CouchDB implementation of DeleteIngestionDictionary
func (asd *AdminServiceDatastoreCouchDB) DeleteIngestionDictionary() (*pb.IngestionDictionary, error) {
	logger.Log.Debugf("Deleting %s\n", ds.IngestionDictionaryStr)

	// Obtain the value of the existing record for a return value.
	existingDictionary, err := asd.GetIngestionDictionary()
	if err != nil {
		logger.Log.Debugf("Unable to delete %s: %s", ds.IngestionDictionaryStr, err.Error())
		return nil, err
	}

	if err = deleteData(asd.dbName, existingDictionary.GetXId(), ds.IngestionDictionaryStr); err != nil {
		logger.Log.Debugf("Unable to delete %s: %s", ds.IngestionDictionaryStr, err.Error())
		return nil, err
	}

	// Return the deleted object.
	logger.Log.Debugf("Deleted %s: %v\n", ds.IngestionDictionaryStr, existingDictionary)
	return existingDictionary, nil

}

// GetIngestionDictionary - CouchDB implementation of GetIngestionDictionary
func (asd *AdminServiceDatastoreCouchDB) GetIngestionDictionary() (*pb.IngestionDictionary, error) {
	db, err := getDatabase(asd.dbName)
	if err != nil {
		return nil, err
	}

	fetchedData, err := getAllOfTypeByIDPrefix(string(ds.IngestionDictionaryType), ds.IngestionDictionaryStr, db)
	if err != nil {
		return nil, err
	}

	// Populate the response
	res := pb.IngestionDictionary{}
	if len(fetchedData) != 0 {
		if err = convertGenericCouchDataToObject(fetchedData[0], &res, ds.IngestionDictionaryStr); err != nil {
			return nil, err
		}
	}

	logger.Log.Debugf("Found %s %v\n", ds.IngestionDictionaryStr, res)
	return &res, nil
}

// GetTenantIDByAlias - InMemory impl of GetTenantIDByAlias
func (asd *AdminServiceDatastoreCouchDB) GetTenantIDByAlias(name string) (string, error) {

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

	response := obj["value"].(string)
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

// Produces all of the views/indicies necessary for the Tenant DB
func generateTenantViews() []map[string]interface{} {
	result := make([]map[string]interface{}, 0)

	monitoredObjectCountByDomain := map[string]interface{}{}
	monitoredObjectCountByDomain["_id"] = "_design/monitoredObjectCount"
	monitoredObjectCountByDomain["language"] = "javascript"
	byDomain := map[string]interface{}{}
	byDomain["map"] = "function(doc) {\n    if (doc.data && doc.data.datatype && doc.data.datatype === 'monitoredObject' && doc.data.domainSet) {\n      for (var i in doc.data.domainSet) {\n        emit(doc.data.domainSet[i], doc._id);\n      }\n    }\n}"
	views := map[string]interface{}{}
	views["byDomain"] = byDomain
	monitoredObjectCountByDomain["views"] = views

	logger.Log.Debug("Adding view for monitoredObjectCountByDomain")
	return append(result, monitoredObjectCountByDomain)
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
