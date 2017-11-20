package couchDB

import (
	"errors"
	"fmt"
	"time"

	"github.com/leesper/couchdb-golang"

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
}

// CreateAdminServiceDAO - instantiates a CouchDB implementation of the
// AdminServiceDatastore.
func CreateAdminServiceDAO() (*AdminServiceDatastoreCouchDB, error) {
	result := new(AdminServiceDatastoreCouchDB)
	cfg, err := gather.GetActiveConfig()
	if err != nil {
		logger.Log.Debugf("Falied to instantiate AdminServiceDatastoreCouchDB: %s", err.Error())
		return nil, err
	}

	// Couch Server Configuration
	provDBURL := fmt.Sprintf("%s:%d",
		cfg.ServerConfig.Datastore.BindIP,
		cfg.ServerConfig.Datastore.BindPort)
	logger.Log.Debug("Admin Service CouchDB URL is: ", provDBURL)
	result.couchHost = provDBURL

	// Couch DB name configuration
	dbName := cfg.ServerConfig.StartupArgs.AdminDB.Name
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
	db, err := getDatabase(asd.dbName)
	if err != nil {
		return nil, err
	}

	// Give the user a known type and timestamps:
	user.Data.Datatype = string(ds.AdminUserType)
	user.Data.CreatedTimestamp = time.Now().Unix()
	user.Data.LastModifiedTimestamp = user.GetData().GetCreatedTimestamp()

	// Marshal the Admin and read the bytes as string.
	storeFormat, err := convertDataToCouchDbSupportedModel(user)
	if err != nil {
		return nil, err
	}

	// Store the user in CouchDB
	_, _, err = storeDataInCouchDB(storeFormat, ds.AdminUserStr, db)
	if err != nil {
		return nil, err
	}

	// Populate the response
	res := pb.AdminUserResponse{}
	if err = convertGenericCouchDataToObject(storeFormat, &res, ds.AdminUserStr); err != nil {
		return nil, err
	}

	// Return the provisioned user.
	logger.Log.Debugf("Created %s: %v\n", ds.AdminUserStr, res)
	return &res, nil
}

// UpdateAdminUser - CouchDB implementation of UpdateAdminUser
func (asd *AdminServiceDatastoreCouchDB) UpdateAdminUser(user *pb.AdminUserRequest) (*pb.AdminUserResponse, error) {
	db, err := getDatabase(asd.dbName)
	if err != nil {
		return nil, err
	}

	// Update timestamp and make sure the type is properly set:
	user.Data.Datatype = string(ds.AdminUserType)
	user.Data.LastModifiedTimestamp = time.Now().Unix()

	// Marshal the Admin and read the bytes as string.
	storeFormat, err := convertDataToCouchDbSupportedModel(user)
	if err != nil {
		return nil, err
	}

	// Store the user in CouchDB
	_, _, err = storeDataInCouchDB(storeFormat, ds.AdminUserStr, db)
	if err != nil {
		return nil, err
	}

	// Populate the response
	res := pb.AdminUserResponse{}
	if err = convertGenericCouchDataToObject(storeFormat, &res, ds.AdminUserStr); err != nil {
		return nil, err
	}

	// Return the provisioned user.
	logger.Log.Debugf("Updated %s: %v\n", ds.AdminUserStr, res)
	return &res, nil
}

// DeleteAdminUser - CouchDB implementation of DeleteAdminUser
func (asd *AdminServiceDatastoreCouchDB) DeleteAdminUser(userID string) (*pb.AdminUserResponse, error) {
	// Obtain the value of the existing record for a return value.
	existingUser, err := asd.GetAdminUser(userID)
	if err != nil {
		logger.Log.Debugf("Unable to delete %s: %s", ds.AdminUserStr, err.Error())
		return nil, err
	}

	// Perform the delete operation on CouchDB
	db, err := getDatabase(asd.dbName)
	if err != nil {
		return nil, err
	}

	if err = deleteByDocID(userID, ds.AdminUserStr, db); err != nil {
		return nil, err
	}

	logger.Log.Debugf("Deleted %s: %v\n", ds.AdminUserStr, existingUser)
	return existingUser, nil
}

// GetAdminUser - CouchDB implementation of GetAdminUser
func (asd *AdminServiceDatastoreCouchDB) GetAdminUser(userID string) (*pb.AdminUserResponse, error) {
	db, err := getDatabase(asd.dbName)
	if err != nil {
		return nil, err
	}

	// Retrieve the user data from CouchDB
	fetchedUser, err := getByDocID(userID, ds.AdminUserStr, db)
	if err != nil {
		return nil, err
	}

	// Marshal the response from the datastore to bytes so that it
	// can be Marshalled back to the proper type.
	res := pb.AdminUserResponse{}
	if err = convertGenericCouchDataToObject(fetchedUser, &res, ds.AdminUserStr); err != nil {
		return nil, err
	}

	logger.Log.Debugf("Retrieved %s: %v\n", ds.AdminUserStr, res)
	return &res, nil
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
	db, err := getDatabase(asd.dbName)
	if err != nil {
		return nil, err
	}

	// Give the tenantDescriptor a known type, and timestamps:
	tenantDescriptor.Data.Datatype = string(ds.TenantDescriptorType)
	tenantDescriptor.Data.CreatedTimestamp = time.Now().Unix()
	tenantDescriptor.Data.LastModifiedTimestamp = tenantDescriptor.GetData().GetCreatedTimestamp()

	// Marshal the tenantDescriptor and read the bytes as string.
	storeFormat, err := convertDataToCouchDbSupportedModel(tenantDescriptor)
	if err != nil {
		return nil, err
	}

	// Store the tenant metadata in CouchDB
	_, _, err = storeDataInCouchDB(storeFormat, ds.TenantDescriptorStr, db)
	if err != nil {
		return nil, err
	}

	// Create a CouchDB database to isolate the tenant data
	_, err = asd.createDatabase(tenantDescriptor.GetXId())
	if err != nil {
		logger.Log.Debugf("Unable to create database for Tenant %s: %s", tenantDescriptor.GetXId(), err.Error())
		return nil, err
	}

	// Populate the response
	res := pb.TenantDescriptorResponse{}
	if err = convertGenericCouchDataToObject(storeFormat, &res, ds.TenantDescriptorStr); err != nil {
		return nil, err
	}

	// Return the provisioned user.
	logger.Log.Debugf("Created %s: %v\n", ds.TenantDescriptorStr, res)
	return &res, nil
}

// UpdateTenantDescriptor - CouchDB implementation of UpdateTenantDescriptor
func (asd *AdminServiceDatastoreCouchDB) UpdateTenantDescriptor(tenantDescriptor *pb.TenantDescriptorRequest) (*pb.TenantDescriptorResponse, error) {
	db, err := getDatabase(asd.dbName)
	if err != nil {
		return nil, err
	}

	// Update timestamp and make sure the type is properly set:
	tenantDescriptor.Data.Datatype = string(ds.TenantDescriptorType)
	tenantDescriptor.Data.LastModifiedTimestamp = time.Now().Unix()

	// Marshal the tenantDescriptor and read the bytes as string.
	storeFormat, err := convertDataToCouchDbSupportedModel(tenantDescriptor)
	if err != nil {
		return nil, err
	}

	// Store the user in CouchDB
	_, _, err = storeDataInCouchDB(storeFormat, ds.TenantDescriptorStr, db)
	if err != nil {
		return nil, err
	}

	// Populate the response
	res := pb.TenantDescriptorResponse{}
	if err = convertGenericCouchDataToObject(storeFormat, &res, ds.TenantDescriptorStr); err != nil {
		return nil, err
	}

	// Return the provisioned user.
	logger.Log.Debugf("Updated %s: %s\n", ds.TenantDescriptorStr, res)
	return &res, nil
}

// DeleteTenant - CouchDB implementation of DeleteTenant
func (asd *AdminServiceDatastoreCouchDB) DeleteTenant(tenantID string) (*pb.TenantDescriptorResponse, error) {
	// Obtain the value of the existing record for a return value.
	existingTenant, err := asd.GetTenantDescriptor(tenantID)
	if err != nil {
		logger.Log.Debugf("Unable to delete %s: %s", ds.TenantDescriptorType, err.Error())
		return nil, err
	}

	// Perform the delete operation on CouchDB
	db, err := getDatabase(asd.dbName)
	if err != nil {
		return nil, err
	}

	if err = deleteByDocID(tenantID, ds.TenantDescriptorStr, db); err != nil {
		return nil, err
	}

	logger.Log.Debugf("Deleted %s: %v\n", ds.TenantDescriptorStr, existingTenant)
	return existingTenant, nil
}

// GetTenantDescriptor - CouchDB implementation of GetTenantDescriptor
func (asd *AdminServiceDatastoreCouchDB) GetTenantDescriptor(tenantID string) (*pb.TenantDescriptorResponse, error) {
	db, err := getDatabase(asd.dbName)
	if err != nil {
		return nil, err
	}

	// Retrieve the tenant data from CouchDB
	fetchedTenant, err := getByDocID(tenantID, ds.TenantDescriptorStr, db)
	if err != nil {
		return nil, err
	}

	// Marshal the response from the datastore to bytes so that it
	// can be Marshalled back to the proper type.
	res := pb.TenantDescriptorResponse{}
	if err = convertGenericCouchDataToObject(fetchedTenant, &res, ds.TenantDescriptorStr); err != nil {
		return nil, err
	}

	logger.Log.Debugf("Retrieved %s: %v\n", ds.TenantDescriptorStr, res)
	return &res, nil
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
