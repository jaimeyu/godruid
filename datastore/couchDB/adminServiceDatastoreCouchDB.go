package couchDB

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/accedian/adh-gather/datastore"
	"github.com/accedian/adh-gather/gather"
	"github.com/accedian/adh-gather/logger"

	pb "github.com/accedian/adh-gather/gathergrpc"
)

const adminUserType string = "adminUser"
const tenantDescriptorType string = "tenantDescriptor"

// AdminServiceDatastoreCouchDB - struct responsible for handling
// database operations for the Admin Service when using CouchDB
// as the storage option.
type AdminServiceDatastoreCouchDB struct {
	server string
	dbName string
}

// CreateDAO - instantiates a CouchDB implementation of the
// AdminServiceDatastore.
func CreateDAO() *AdminServiceDatastoreCouchDB {
	result := new(AdminServiceDatastoreCouchDB)
	cfg, err := gather.GetActiveConfig()
	if err != nil {
		logger.Log.Errorf("Falied to instantiate AdminServiceDatastoreCouchDB: %v", err)
	}

	provDBURL := fmt.Sprintf("%s:%d",
		cfg.ServerConfig.Datastore.BindIP,
		cfg.ServerConfig.Datastore.BindPort)
	logger.Log.Debug("CouchDB URL is: ", provDBURL)
	result.server = provDBURL
	result.dbName = result.server + "/adh-admin"

	return result
}

// CreateAdminUser - CouchDB implementation of CreateAdminUser
func (couchDB *AdminServiceDatastoreCouchDB) CreateAdminUser(user *pb.AdminUser) (*pb.AdminUser, error) {
	db, err := couchDB.GetDatabase()
	if err != nil {
		return nil, err
	}

	// Give the user a known id, type, and timestamps:
	user.Id = user.Username
	user.CreatedTimestamp = time.Now().Unix()
	user.LastModifiedTimestamp = user.GetCreatedTimestamp()
	user.Datatype = adminUserType

	// Marshal the Admin and read the bytes as string.
	storeFormat, err := ConvertDataToCouchDbSupportedModel(user)
	if err != nil {
		return nil, err
	}

	// Store the user in CouchDB
	_, rev, err := StoreDataInCouchDB(storeFormat, datastore.AdminUserStr, db)
	if err != nil {
		return nil, err
	}

	// Add the revision number to the response
	user.Rev = rev

	// Return the provisioned user.
	logger.Log.Infof("Created %s: %v\n", datastore.AdminUserStr, user)
	return user, nil
}

// UpdateAdminUser - CouchDB implementation of UpdateAdminUser
func (couchDB *AdminServiceDatastoreCouchDB) UpdateAdminUser(user *pb.AdminUser) (*pb.AdminUser, error) {
	db, err := couchDB.GetDatabase()
	if err != nil {
		return nil, err
	}

	// Update timestamp and make sure the type is properly set:
	user.LastModifiedTimestamp = time.Now().Unix()
	user.Datatype = adminUserType

	// Marshal the Admin and read the bytes as string.
	storeFormat, err := ConvertDataToCouchDbSupportedModel(user)

	// Store the user in CouchDB
	id, rev, err := StoreDataInCouchDB(storeFormat, datastore.AdminUserStr, db)
	if err != nil {
		return nil, err
	}

	// Add the evision number to the response
	user.Rev = rev
	logger.Log.Debugf("Successfully updated %s %s with rev %s", datastore.AdminUserStr, id, rev)

	// Return the provisioned user.
	logger.Log.Infof("Updated %s: %v\n", datastore.AdminUserStr, user)
	return user, nil
}

// DeleteAdminUser - CouchDB implementation of DeleteAdminUser
func (couchDB *AdminServiceDatastoreCouchDB) DeleteAdminUser(userID string) (*pb.AdminUser, error) {
	// Obtain the value of the existing record for a return value.
	existingUser, err := couchDB.GetAdminUser(userID)
	if err != nil {
		logger.Log.Errorf("Unable to delete %s: %v\n", datastore.AdminUserStr, err)
		return nil, err
	}

	// Perform the delete operation on CouchDB
	db, err := couchDB.GetDatabase()
	if err != nil {
		return nil, err
	}

	if err = DeleteByDocID(userID, datastore.AdminUserStr, db); err != nil {
		return nil, err
	}

	return existingUser, nil
}

// GetAdminUser - CouchDB implementation of GetAdminUser
func (couchDB *AdminServiceDatastoreCouchDB) GetAdminUser(userID string) (*pb.AdminUser, error) {
	db, err := couchDB.GetDatabase()
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
	res, err := convertGenericObjectToAdminUser(fetchedUser)
	if err != nil {
		return nil, err
	}

	return res, nil
}

// GetAllAdminUsers - CouchDB implementation of GetAllAdminUsers
func (couchDB *AdminServiceDatastoreCouchDB) GetAllAdminUsers() (*pb.AdminUserList, error) {
	db, err := couchDB.GetDatabase()
	if err != nil {
		return nil, err
	}

	fetchedUserList, err := GetAllOfType(adminUserType, datastore.AdminUserStr, db)
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
func (couchDB *AdminServiceDatastoreCouchDB) CreateTenant(tenantDescriptor *pb.TenantDescriptor) (*pb.TenantDescriptor, error) {
	db, err := couchDB.GetDatabase()
	if err != nil {
		return nil, err
	}

	// Give the tenantDescriptor a known id, type, and timestamps:
	tenantDescriptor.Id = tenantDescriptor.GetUrlSubdomain()
	tenantDescriptor.CreatedTimestamp = time.Now().Unix()
	tenantDescriptor.LastModifiedTimestamp = tenantDescriptor.GetCreatedTimestamp()
	tenantDescriptor.Datatype = tenantDescriptorType

	// Marshal the tenantDescriptor and read the bytes as string.
	storeFormat, err := ConvertDataToCouchDbSupportedModel(tenantDescriptor)
	if err != nil {
		return nil, err
	}

	// Store the user in CouchDB
	_, rev, err := StoreDataInCouchDB(storeFormat, datastore.TenantDescriptorStr, db)
	if err != nil {
		return nil, err
	}

	// Add the revision number to the response
	tenantDescriptor.Rev = rev

	// Return the provisioned user.
	logger.Log.Infof("Created %s: %v\n", datastore.TenantDescriptorStr, tenantDescriptor)
	return tenantDescriptor, nil
}

// UpdateTenantDescriptor - CouchDB implementation of UpdateTenantDescriptor
func (couchDB *AdminServiceDatastoreCouchDB) UpdateTenantDescriptor(tenantDescriptor *pb.TenantDescriptor) (*pb.TenantDescriptor, error) {
	db, err := couchDB.GetDatabase()
	if err != nil {
		return nil, err
	}

	// Update timestamp and make sure the type is properly set:
	tenantDescriptor.LastModifiedTimestamp = time.Now().Unix()
	tenantDescriptor.Datatype = tenantDescriptorType

	// Marshal the tenantDescriptor and read the bytes as string.
	storeFormat, err := ConvertDataToCouchDbSupportedModel(tenantDescriptor)

	// Store the user in CouchDB
	id, rev, err := StoreDataInCouchDB(storeFormat, datastore.TenantDescriptorStr, db)
	if err != nil {
		return nil, err
	}

	// Add the revision number to the response
	tenantDescriptor.Rev = rev
	logger.Log.Debugf("Successfully updated %s %s with rev %s", datastore.TenantDescriptorStr, id, rev)

	// Return the provisioned user.
	logger.Log.Infof("Updated %s: %v\n", datastore.TenantDescriptorStr, tenantDescriptor)
	return tenantDescriptor, nil
}

// DeleteTenant - CouchDB implementation of DeleteTenant
func (couchDB *AdminServiceDatastoreCouchDB) DeleteTenant(tenantID string) (*pb.TenantDescriptor, error) {
	// Obtain the value of the existing record for a return value.
	existingTenant, err := couchDB.GetTenantDescriptor(tenantID)
	if err != nil {
		logger.Log.Errorf("Unable to delete %s: %v\n", tenantDescriptorType, err)
		return nil, err
	}

	// Perform the delete operation on CouchDB
	db, err := couchDB.GetDatabase()
	if err != nil {
		return nil, err
	}

	if err = DeleteByDocID(tenantID, datastore.TenantDescriptorStr, db); err != nil {
		return nil, err
	}

	return existingTenant, nil
}

// GetTenantDescriptor - CouchDB implementation of GetTenantDescriptor
func (couchDB *AdminServiceDatastoreCouchDB) GetTenantDescriptor(tenantID string) (*pb.TenantDescriptor, error) {
	db, err := couchDB.GetDatabase()
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
	res, err := convertGenericObjectToTenantDescriptor(fetchedTenant)
	if err != nil {
		return nil, err
	}

	return res, nil
}

// Takes the map[string]interface{} generic data returned by CouchDB and
// converts it to an AdminUser.
func convertGenericObjectToAdminUser(genericUser map[string]interface{}) (*pb.AdminUser, error) {
	genericUserInBytes, err := ConvertGenericObjectToBytesWithCouchDbFields(genericUser)
	if err != nil {
		return nil, err
	}

	res := pb.AdminUser{}
	err = json.Unmarshal(genericUserInBytes, &res)
	if err != nil {
		logger.Log.Errorf("Error converting generic user data to %s type: %v\n", datastore.AdminUserStr, err)
		return nil, err
	}

	logger.Log.Debugf("Converted generic data to %s: %v\n", datastore.AdminUserStr, res)

	return &res, nil
}

// Takes a set of generic data that contains a list of AdminUsers and converts it to
// and ADH AdminUserList object
func convertGenericObjectListToAdminUserList(genericUserList []map[string]interface{}) (*pb.AdminUserList, error) {
	res := new(pb.AdminUserList)
	for _, genericUserObject := range genericUserList {
		user, err := convertGenericObjectToAdminUser(genericUserObject)
		if err != nil {
			continue
		}
		res.List = append(res.List, user)
	}

	logger.Log.Debugf("Converted generic data to %s List: %v\n", datastore.AdminUserStr, res)

	return res, nil
}

// Takes the map[string]interface{} generic data returned by CouchDB and
// converts it to an TenantDescriptor.
func convertGenericObjectToTenantDescriptor(genericData map[string]interface{}) (*pb.TenantDescriptor, error) {
	genericDataInBytes, err := ConvertGenericObjectToBytesWithCouchDbFields(genericData)
	if err != nil {
		return nil, err
	}

	res := pb.TenantDescriptor{}
	err = json.Unmarshal(genericDataInBytes, &res)
	if err != nil {
		logger.Log.Errorf("Error converting generic data to %s type: %v\n", datastore.TenantDescriptorStr, err)
		return nil, err
	}

	logger.Log.Debugf("Converted generic data to %s: %v\n", datastore.TenantDescriptorStr, res)

	return &res, nil
}
