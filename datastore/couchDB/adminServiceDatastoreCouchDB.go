package couchDB

import (
	"errors"
	"fmt"

	"github.com/leesper/couchdb-golang"

	"github.com/accedian/adh-gather/config"
	ds "github.com/accedian/adh-gather/datastore"
	"github.com/accedian/adh-gather/gather"
	"github.com/accedian/adh-gather/logger"

	pb "github.com/accedian/adh-gather/gathergrpc"
)

// AdminServiceDatastoreCouchDB - struct responsible for handling
// database operations for the Admin Service when using CouchDB
// as the storage option.
type AdminServiceDatastoreCouchDB struct {
	couchHost string
	dbName    string
	server    *couchdb.Server
	cfg       config.Provider
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
	dbName := result.cfg.GetString(gather.CK_args_admindb_name.String())
	result.dbName = result.couchHost + "/" + dbName
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
	logger.Log.Debugf("Creating %s: %v\n", ds.AdminUserStr, user)

	dataType := string(ds.AdminUserType)
	dataContainer := pb.AdminUserResponse{}
	if err := storeData(asd.dbName, user, dataType, ds.AdminUserStr, &dataContainer); err != nil {
		return nil, err
	}

	// Return the provisioned object.
	logger.Log.Debugf("Created %s: %v\n", ds.AdminUserStr, dataContainer)
	return &dataContainer, nil
}

// UpdateAdminUser - CouchDB implementation of UpdateAdminUser
func (asd *AdminServiceDatastoreCouchDB) UpdateAdminUser(user *pb.AdminUserRequest) (*pb.AdminUserResponse, error) {
	logger.Log.Debugf("Updating %s: %v\n", ds.AdminUserStr, user)

	dataType := string(ds.AdminUserType)
	dataContainer := pb.AdminUserResponse{}
	if err := updateData(asd.dbName, user, dataType, ds.AdminUserStr, &dataContainer); err != nil {
		return nil, err
	}

	// Return the provisioned object.
	logger.Log.Debugf("Updated %s: %v\n", ds.AdminUserStr, dataContainer)
	return &dataContainer, nil
}

// DeleteAdminUser - CouchDB implementation of DeleteAdminUser
func (asd *AdminServiceDatastoreCouchDB) DeleteAdminUser(userID string) (*pb.AdminUserResponse, error) {
	logger.Log.Debugf("Deleting %s for %s\n", ds.AdminUserStr, userID)

	// Obtain the value of the existing record for a return value.
	existingUser, err := asd.GetAdminUser(userID)
	if err != nil {
		logger.Log.Debugf("Unable to fetch %s to delete: %s", ds.AdminUserStr, err.Error())
		return nil, err
	}

	if err := deleteData(asd.dbName, userID, ds.AdminUserStr); err != nil {
		logger.Log.Debugf("Unable to delete %s: %s", ds.AdminUserStr, err.Error())
		return nil, err
	}

	// Return the deleted object.
	logger.Log.Debugf("Deleted %s: %v\n", ds.AdminUserStr, existingUser)
	return existingUser, nil
}

// GetAdminUser - CouchDB implementation of GetAdminUser
func (asd *AdminServiceDatastoreCouchDB) GetAdminUser(userID string) (*pb.AdminUserResponse, error) {
	logger.Log.Debugf("Retrieving %s for %s\n", ds.AdminUserStr, userID)

	dataContainer := pb.AdminUserResponse{}
	if err := getData(asd.dbName, userID, ds.AdminUserStr, &dataContainer); err != nil {
		return nil, err
	}

	// Return the provisioned object.
	logger.Log.Debugf("Retrieved %s: %v\n", ds.AdminUserStr, dataContainer)
	return &dataContainer, nil
}

// GetAllAdminUsers - CouchDB implementation of GetAllAdminUsers
func (asd *AdminServiceDatastoreCouchDB) GetAllAdminUsers() (*pb.AdminUserListResponse, error) {
	db, err := getDatabase(asd.dbName)
	if err != nil {
		return nil, err
	}

	fetchedUserList, err := getAllOfTypeByIDPrefix(string(ds.AdminUserType), ds.AdminUserStr, db)
	if err != nil {
		return nil, err
	}

	// Marshal the response from the datastore to bytes so that it
	// can be Marshalled back to the proper type.
	res, err := convertGenericObjectListToAdminUserList(fetchedUserList)
	if err != nil {
		return nil, err
	}

	logger.Log.Debugf("Retrieved %d %ss\n", len(res.GetData()), ds.AdminUserStr)
	return res, nil
}

// CreateTenant - CouchDB implementation of CreateTenant
func (asd *AdminServiceDatastoreCouchDB) CreateTenant(tenantDescriptor *pb.TenantDescriptorRequest) (*pb.TenantDescriptorResponse, error) {
	logger.Log.Debugf("Creating %s: %v\n", ds.TenantStr, tenantDescriptor)

	dataType := string(ds.TenantDescriptorType)
	dataContainer := pb.TenantDescriptorResponse{}
	if err := storeData(asd.dbName, tenantDescriptor, dataType, ds.TenantStr, &dataContainer); err != nil {
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
	logger.Log.Debugf("Created %s: %v\n", ds.TenantStr, dataContainer)
	return &dataContainer, nil
}

// UpdateTenantDescriptor - CouchDB implementation of UpdateTenantDescriptor
func (asd *AdminServiceDatastoreCouchDB) UpdateTenantDescriptor(tenantDescriptor *pb.TenantDescriptorRequest) (*pb.TenantDescriptorResponse, error) {
	logger.Log.Debugf("Updating %s: %v\n", ds.TenantStr, tenantDescriptor)

	dataType := string(ds.TenantDescriptorType)
	dataContainer := pb.TenantDescriptorResponse{}
	if err := updateData(asd.dbName, tenantDescriptor, dataType, ds.TenantStr, &dataContainer); err != nil {
		return nil, err
	}

	// Return the provisioned object.
	logger.Log.Debugf("Updated %s: %v\n", ds.TenantStr, dataContainer)
	return &dataContainer, nil
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

	// Truy to delete the DB for the tenant
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
	logger.Log.Debugf("Retrieving %s for %s\n", ds.TenantDescriptorStr, tenantID)

	dataContainer := pb.TenantDescriptorResponse{}
	if err := getData(asd.dbName, tenantID, ds.TenantDescriptorStr, &dataContainer); err != nil {
		return nil, err
	}

	// Return the provisioned object.
	logger.Log.Debugf("Retrieved %s: %v\n", ds.TenantDescriptorStr, dataContainer)
	return &dataContainer, nil
}

// GetAllTenantDescriptors - CouchDB implementation of GetAllTenantDescriptors
func (asd *AdminServiceDatastoreCouchDB) GetAllTenantDescriptors() (*pb.TenantDescriptorListResponse, error) {
	db, err := getDatabase(asd.dbName)
	if err != nil {
		return nil, err
	}

	fetchedTenantList, err := getAllOfTypeByIDPrefix(string(ds.TenantDescriptorType), ds.TenantStr, db)
	if err != nil {
		return nil, err
	}

	// Marshal the response from the datastore to bytes so that it
	// can be Marshalled back to the proper type.
	res, err := convertGenericObjectListToTenantDescriptorList(fetchedTenantList)
	if err != nil {
		return nil, err
	}

	logger.Log.Debugf("Found %d %ss\n", len(res.GetData()), ds.TenantDescriptorStr)
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

	// resource, err := couchdb.NewResource(createDBPathStr(asd.couchHost, dbName), nil)
	// if err != nil {
	// 	logger.Log.Debugf("Unable to add views to database: %s", err.Error())
	// 	return err
	// }

	// tenantViews := generateTenantViews()
	// for _, viewPayload := range tenantViews {
	// 	_, err := addDesignDocument(viewPayload, resource)
	// 	if err != nil {
	// 		return err
	// 	}
	// }

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

// Takes a set of generic data that contains a list of AdminUsers and converts it to
// and ADH AdminUserList object
func convertGenericObjectListToAdminUserList(genericUserList []map[string]interface{}) (*pb.AdminUserListResponse, error) {
	res := new(pb.AdminUserListResponse)
	for _, genericUserObject := range genericUserList {
		user := pb.AdminUserResponse{}
		if err := convertGenericCouchDataToObject(genericUserObject, &user, ds.AdminUserStr); err != nil {
			continue
		}
		res.Data = append(res.GetData(), &user)
	}

	logger.Log.Debugf("Converted generic data to %s List: %v\n", ds.AdminUserStr, res)

	return res, nil
}

func convertGenericObjectListToTenantDescriptorList(genericTenantList []map[string]interface{}) (*pb.TenantDescriptorListResponse, error) {
	res := new(pb.TenantDescriptorListResponse)
	for _, genericTenantObject := range genericTenantList {
		tenant := pb.TenantDescriptorResponse{}
		if err := convertGenericCouchDataToObject(genericTenantObject, &tenant, ds.TenantDescriptorStr); err != nil {
			continue
		}
		res.Data = append(res.GetData(), &tenant)
	}

	logger.Log.Debugf("Converted generic data to %s List: %v\n", ds.TenantDescriptorStr, res)

	return res, nil
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
