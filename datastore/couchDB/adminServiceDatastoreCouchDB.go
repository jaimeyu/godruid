package couchDB

import (
	"encoding/json"
	"fmt"
	"net/url"
	"time"

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

	logger.Log.Infof("Using DB %s to Create Admin User %v \n", couchDB.dbName, user)

	// Give the user a known id, type, and timestamps:
	user.Id = user.Username
	user.CreatedTimestamp = time.Now().Unix()
	user.LastModifiedTimestamp = user.GetCreatedTimestamp()
	user.Datatype = adminUserType

	// Marshal the Admin and read the bytes as string.
	storeFormat, err := ConvertDataToCouchDbSupportedModel(user)

	logger.Log.Debugf("Attempting to create Admin User: %v", storeFormat)

	// Store the user in PROV DB
	options := new(url.Values)
	id, rev, err := db.Save(storeFormat, *options)
	if err != nil {
		logger.Log.Errorf("Unable to create Admin User: %v\n", err)
		return nil, err
	}

	// Add the evision number to the response
	user.Rev = rev
	logger.Log.Debugf("Successfully created Admin User %s with rev %s", id, rev)

	// Return the provisioned user.
	logger.Log.Infof("Created Admin User: %v\n", user)
	return user, nil
}

// UpdateAdminUser - CouchDB implementation of UpdateAdminUser
func (couchDB *AdminServiceDatastoreCouchDB) UpdateAdminUser(user *pb.AdminUser) (*pb.AdminUser, error) {
	db, err := couchDB.GetDatabase()
	if err != nil {
		return nil, err
	}

	logger.Log.Infof("Using DB %s to update Admin User %v \n", couchDB.dbName, user)

	// Update timestamp and make sure the type is properly set:
	user.LastModifiedTimestamp = time.Now().Unix()
	user.Datatype = adminUserType

	// Marshal the Admin and read the bytes as string.
	storeFormat, err := ConvertDataToCouchDbSupportedModel(user)

	// Add the _rev field required for CouchDB conflict resolution
	logger.Log.Debugf("Attempting to update Admin User: %v", storeFormat)

	// Store the user in PROV DB
	options := new(url.Values)
	id, rev, err := db.Save(storeFormat, *options)
	if err != nil {
		logger.Log.Errorf("Unable to update Admin User: %v\n", err)
		return nil, err
	}

	// Add the evision number to the response
	user.Rev = rev
	logger.Log.Debugf("Successfully updated Admin User %s with rev %s", id, rev)

	// Return the provisioned user.
	logger.Log.Infof("Updated Admin User: %v\n", user)
	return user, nil
}

// DeleteAdminUser - CouchDB implementation of DeleteAdminUser
func (couchDB *AdminServiceDatastoreCouchDB) DeleteAdminUser(userID string) (*pb.AdminUser, error) {
	// Obtain the value of the existing record for a return value.
	existingUser, err := couchDB.GetAdminUser(userID)
	if err != nil {
		logger.Log.Errorf("Unable to delete Admin User: %v\n", err)
		return nil, err
	}

	// Perform the delete operation
	db, err := couchDB.GetDatabase()
	if err != nil {
		return nil, err
	}

	logger.Log.Infof("Using db %s to delete Admin User %s\n", couchDB.dbName, userID)

	err = db.Delete(userID)
	if err != nil {
		logger.Log.Errorf("Error deleting Admin User %s: %v\n", userID, err)
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

	logger.Log.Infof("Using db %s to retrieve Admin User %s\n", couchDB.dbName, userID)

	// Get the Admin User from CouchDB
	options := new(url.Values)
	fetchedUser, err := db.Get(userID, *options)
	if err != nil {
		logger.Log.Errorf("Error retrieving Admin User %s: %v\n", userID, err)
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

	logger.Log.Infof("Using db %s to retrieve all Admin Users\n", couchDB.dbName)

	// Get the Admin User from CouchDB
	selector := fmt.Sprintf(`datatype == "%s"`, adminUserType)
	fetchedUserList, err := db.Query(nil, selector, nil, nil, nil, nil)
	if err != nil {
		logger.Log.Errorf("Error retrieving Admin User List: %v\n", err)
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
	// Stub to implement
	return nil, nil
}

// UpdateTenantDescriptor - CouchDB implementation of UpdateTenantDescriptor
func (couchDB *AdminServiceDatastoreCouchDB) UpdateTenantDescriptor(tenantDescriptor *pb.TenantDescriptor) (*pb.TenantDescriptor, error) {
	// Stub to implement
	return nil, nil
}

// DeleteTenant - CouchDB implementation of DeleteTenant
func (couchDB *AdminServiceDatastoreCouchDB) DeleteTenant(tenantID string) (*pb.TenantDescriptor, error) {
	// Stub to implement
	return nil, nil
}

// GetTenantDescriptor - CouchDB implementation of GetTenantDescriptor
func (couchDB *AdminServiceDatastoreCouchDB) GetTenantDescriptor(tenantID string) (*pb.TenantDescriptor, error) {
	// Stub to implement
	return nil, nil
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
		logger.Log.Errorf("Error converting generic user data to Admin User type: %v\n", err)
		return nil, err
	}

	logger.Log.Debugf("Converted generic data to AdminUser: %v\n", res)

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

	logger.Log.Debugf("Converted generic data to AdminUserList: %v\n", res)

	return res, nil
}
