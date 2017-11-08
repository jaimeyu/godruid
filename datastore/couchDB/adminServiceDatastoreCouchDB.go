package couchDB

import (
	"errors"
	"fmt"
	"time"

	"github.com/leesper/couchdb-golang"

	"github.com/accedian/adh-gather/datastore"
	"github.com/accedian/adh-gather/gather"
	"github.com/accedian/adh-gather/logger"

	pb "github.com/accedian/adh-gather/gathergrpc"
)

const adminUserType string = "adminUser"
const tenantDescriptorType string = "tenant"

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
func CreateAdminServiceDAO() *AdminServiceDatastoreCouchDB {
	result := new(AdminServiceDatastoreCouchDB)
	cfg, err := gather.GetActiveConfig()
	if err != nil {
		logger.Log.Errorf("Falied to instantiate AdminServiceDatastoreCouchDB: %v", err)
	}

	provDBURL := fmt.Sprintf("%s:%d",
		cfg.ServerConfig.Datastore.BindIP,
		cfg.ServerConfig.Datastore.BindPort)
	logger.Log.Debug("Admin Service CouchDB URL is: ", provDBURL)
	result.couchHost = provDBURL
	result.dbName = result.couchHost + "/adh-admin"
	server, err := couchdb.NewServer(result.couchHost)
	if err != nil {
		logger.Log.Errorf("Falied to instantiate AdminServiceDatastoreCouchDB: %v", err)
	}

	result.server = server
	return result
}

// CreateAdminUser - CouchDB implementation of CreateAdminUser
func (asd *AdminServiceDatastoreCouchDB) CreateAdminUser(user *pb.AdminUserRequest) (*pb.AdminUserResponse, error) {
	db, err := GetDatabase(asd.dbName)
	if err != nil {
		return nil, err
	}

	// Give the user a known type and timestamps:
	user.Data.Datatype = adminUserType
	user.Data.CreatedTimestamp = time.Now().Unix()
	user.Data.LastModifiedTimestamp = user.GetData().GetCreatedTimestamp()

	// Marshal the Admin and read the bytes as string.
	storeFormat, err := ConvertDataToCouchDbSupportedModel(user)
	if err != nil {
		return nil, err
	}

	// Store the user in CouchDB
	_, _, err = StoreDataInCouchDB(storeFormat, datastore.AdminUserStr, db)
	if err != nil {
		return nil, err
	}

	// Populate the response
	res := pb.AdminUserResponse{}
	err = ConvertGenericCouchDataToObject(storeFormat, &res, datastore.AdminUserStr)
	if err != nil {
		return nil, err
	}

	// err = ConvertGenericCouchDataToObject(storeFormat, &res, datastore.AdminUserStr)
	// if err != nil {
	// 	return nil, err
	// }

	// Return the provisioned user.
	logger.Log.Infof("Created %s: %v\n", datastore.AdminUserStr, res)
	return &res, nil
}

// UpdateAdminUser - CouchDB implementation of UpdateAdminUser
func (asd *AdminServiceDatastoreCouchDB) UpdateAdminUser(user *pb.AdminUserRequest) (*pb.AdminUserResponse, error) {
	db, err := GetDatabase(asd.dbName)
	if err != nil {
		return nil, err
	}

	// Update timestamp and make sure the type is properly set:
	user.Data.Datatype = adminUserType
	user.Data.LastModifiedTimestamp = time.Now().Unix()

	// Marshal the Admin and read the bytes as string.
	storeFormat, err := ConvertDataToCouchDbSupportedModel(user)
	if err != nil {
		return nil, err
	}

	// Store the user in CouchDB
	_, _, err = StoreDataInCouchDB(storeFormat, datastore.AdminUserStr, db)
	if err != nil {
		return nil, err
	}

	// Populate the response
	res := pb.AdminUserResponse{}
	err = ConvertGenericCouchDataToObject(storeFormat, &res, datastore.AdminUserStr)
	if err != nil {
		return nil, err
	}

	// Return the provisioned user.
	logger.Log.Infof("Updated %s: %v\n", datastore.AdminUserStr, res)
	return &res, nil
}

// DeleteAdminUser - CouchDB implementation of DeleteAdminUser
func (asd *AdminServiceDatastoreCouchDB) DeleteAdminUser(userID string) (*pb.AdminUserResponse, error) {
	// Obtain the value of the existing record for a return value.
	existingUser, err := asd.GetAdminUser(userID)
	if err != nil {
		logger.Log.Errorf("Unable to delete %s: %v\n", datastore.AdminUserStr, err)
		return nil, err
	}

	// Perform the delete operation on CouchDB
	db, err := GetDatabase(asd.dbName)
	if err != nil {
		return nil, err
	}

	if err = DeleteByDocID(userID, datastore.AdminUserStr, db); err != nil {
		return nil, err
	}

	return existingUser, nil
}

// GetAdminUser - CouchDB implementation of GetAdminUser
func (asd *AdminServiceDatastoreCouchDB) GetAdminUser(userID string) (*pb.AdminUserResponse, error) {
	db, err := GetDatabase(asd.dbName)
	if err != nil {
		return nil, err
	}

	// Retrieve the user data from CouchDB
	fetchedUser, err := GetByDocID(userID, datastore.AdminUserStr, db)
	if err != nil {
		return nil, err
	}

	// Marshal the response from the datastore to bytes so that it
	// can be Marshalled back to the proper type.
	res := pb.AdminUserResponse{}
	err = ConvertGenericCouchDataToObject(fetchedUser, &res, datastore.AdminUserStr)
	if err != nil {
		return nil, err
	}

	return &res, nil
}

// GetAllAdminUsers - CouchDB implementation of GetAllAdminUsers
func (asd *AdminServiceDatastoreCouchDB) GetAllAdminUsers() (*pb.AdminUserListResponse, error) {
	db, err := GetDatabase(asd.dbName)
	if err != nil {
		return nil, err
	}

	fetchedUserList, err := GetAllOfTypeByIDPrefix(adminUserType, datastore.AdminUserStr, db)
	if err != nil {
		return nil, err
	}

	// Marshal the response from the datastore to bytes so that it
	// can be Marshalled back to the proper type.
	res, err := convertGenericObjectListToAdminUserList(fetchedUserList)
	if err != nil {
		return nil, err
	}

	return res, nil
}

// CreateTenant - CouchDB implementation of CreateTenant
func (asd *AdminServiceDatastoreCouchDB) CreateTenant(tenantDescriptor *pb.TenantDescriptorRequest) (*pb.TenantDescriptorResponse, error) {
	db, err := GetDatabase(asd.dbName)
	if err != nil {
		return nil, err
	}

	// Give the tenantDescriptor a known type, and timestamps:
	tenantDescriptor.Data.Datatype = tenantDescriptorType
	tenantDescriptor.Data.CreatedTimestamp = time.Now().Unix()
	tenantDescriptor.Data.LastModifiedTimestamp = tenantDescriptor.GetData().GetCreatedTimestamp()

	// Marshal the tenantDescriptor and read the bytes as string.
	storeFormat, err := ConvertDataToCouchDbSupportedModel(tenantDescriptor)
	if err != nil {
		return nil, err
	}

	// Store the tenant metadata in CouchDB
	_, _, err = StoreDataInCouchDB(storeFormat, datastore.TenantDescriptorStr, db)
	if err != nil {
		return nil, err
	}

	// Create a CouchDB database to isolate the tenant data
	_, err = asd.createDatabase(tenantDescriptor.GetXId())
	if err != nil {
		logger.Log.Errorf("Unable to create database for Tenant %s: %v", tenantDescriptor.GetXId(), err)
		return nil, err
	}

	// Populate the response
	res := pb.TenantDescriptorResponse{}
	err = ConvertGenericCouchDataToObject(storeFormat, &res, datastore.TenantDescriptorStr)
	if err != nil {
		return nil, err
	}

	// Return the provisioned user.
	logger.Log.Infof("Created %s: %v\n", datastore.TenantDescriptorStr, res)
	return &res, nil
}

// UpdateTenantDescriptor - CouchDB implementation of UpdateTenantDescriptor
func (asd *AdminServiceDatastoreCouchDB) UpdateTenantDescriptor(tenantDescriptor *pb.TenantDescriptorRequest) (*pb.TenantDescriptorResponse, error) {
	db, err := GetDatabase(asd.dbName)
	if err != nil {
		return nil, err
	}

	// Update timestamp and make sure the type is properly set:
	tenantDescriptor.Data.Datatype = tenantDescriptorType
	tenantDescriptor.Data.LastModifiedTimestamp = time.Now().Unix()

	// Marshal the tenantDescriptor and read the bytes as string.
	storeFormat, err := ConvertDataToCouchDbSupportedModel(tenantDescriptor)

	// Store the user in CouchDB
	_, _, err = StoreDataInCouchDB(storeFormat, datastore.TenantDescriptorStr, db)
	if err != nil {
		return nil, err
	}

	// Populate the response
	res := pb.TenantDescriptorResponse{}
	err = ConvertGenericCouchDataToObject(storeFormat, &res, datastore.TenantDescriptorStr)
	if err != nil {
		return nil, err
	}

	// Return the provisioned user.
	logger.Log.Infof("Updated %s: %v\n", datastore.TenantDescriptorStr, res)
	return &res, nil
}

// DeleteTenant - CouchDB implementation of DeleteTenant
func (asd *AdminServiceDatastoreCouchDB) DeleteTenant(tenantID string) (*pb.TenantDescriptorResponse, error) {
	// Obtain the value of the existing record for a return value.
	existingTenant, err := asd.GetTenantDescriptor(tenantID)
	if err != nil {
		logger.Log.Errorf("Unable to delete %s: %v\n", tenantDescriptorType, err)
		return nil, err
	}

	// Perform the delete operation on CouchDB
	db, err := GetDatabase(asd.dbName)
	if err != nil {
		return nil, err
	}

	if err = DeleteByDocID(tenantID, datastore.TenantDescriptorStr, db); err != nil {
		return nil, err
	}

	return existingTenant, nil
}

// GetTenantDescriptor - CouchDB implementation of GetTenantDescriptor
func (asd *AdminServiceDatastoreCouchDB) GetTenantDescriptor(tenantID string) (*pb.TenantDescriptorResponse, error) {
	db, err := GetDatabase(asd.dbName)
	if err != nil {
		return nil, err
	}

	// Retrieve the tenant data from CouchDB
	fetchedTenant, err := GetByDocID(tenantID, datastore.TenantDescriptorStr, db)
	if err != nil {
		return nil, err
	}

	// Marshal the response from the datastore to bytes so that it
	// can be Marshalled back to the proper type.
	res := pb.TenantDescriptorResponse{}
	err = ConvertGenericCouchDataToObject(fetchedTenant, &res, datastore.TenantDescriptorStr)
	if err != nil {
		return nil, err
	}

	return &res, nil
}

// createDatabase - creates a database in CouchDB identified by the provided name.
func (asd *AdminServiceDatastoreCouchDB) createDatabase(dbName string) (*couchdb.Database, error) {
	if len(dbName) == 0 {
		return nil, errors.New("Unable to create database if no identifier is provided")
	}
	if asd.server.Contains(dbName) {
		return nil, errors.New("Unable to create database '" + dbName + "': database already exists")
	}

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

	return asd.server.Delete(dbName)
}

// Takes a set of generic data that contains a list of AdminUsers and converts it to
// and ADH AdminUserList object
func convertGenericObjectListToAdminUserList(genericUserList []map[string]interface{}) (*pb.AdminUserListResponse, error) {
	res := new(pb.AdminUserListResponse)
	for _, genericUserObject := range genericUserList {
		user := pb.AdminUserResponse{}
		err := ConvertGenericCouchDataToObject(genericUserObject, &user, datastore.AdminUserStr)
		if err != nil {
			continue
		}
		res.List = append(res.List, &user)
	}

	logger.Log.Debugf("Converted generic data to %s List: %v\n", datastore.AdminUserStr, res)

	return res, nil
}
