package couchDB

import (
	"encoding/json"
	"fmt"
	"net/url"
	"time"

	"github.com/accedian/adh-gather/gather"
	"github.com/accedian/adh-gather/logger"

	pb "github.com/accedian/adh-gather/gathergrpc"
	couchdb "github.com/leesper/couchdb-golang"
)

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
	db, err := couchdb.NewDatabase(couchDB.dbName)
	if err != nil {
		logger.Log.Errorf("Unable to connect to Prov DB %s: %v\n", couchDB.server, err)
		return nil, err
	}

	logger.Log.Infof("Using DB %s to Create user %v \n", couchDB.dbName, user)

	// Give the user a known id and timestamps:
	user.XId = user.Username
	user.CreatedTimestamp = time.Now().Unix()
	user.LastModifiedTimestamp = user.GetCreatedTimestamp()

	// Marshal the Admin and read the bytes as string.
	storeFormat, err := convertAdminUserToGenericObject(user)

	logger.Log.Debugf("Attempting to store user: %v", storeFormat)

	// Store the user in PROV DB
	options := new(url.Values)
	id, rev, err := db.Save(storeFormat, *options)
	if err != nil {
		logger.Log.Errorf("Unable to store admin user: %v\n", err)
		return nil, err
	}

	// Leaving this here as a reminder that I need to add revision to the data model.
	logger.Log.Debugf("Successfully stored user %s with rev %s", id, rev)

	// Return the provisioned user.
	logger.Log.Infof("Created Admin User: %v\n", user)
	return user, nil
}

// UpdateAdminUser - CouchDB implementation of UpdateAdminUser
func (couchDB *AdminServiceDatastoreCouchDB) UpdateAdminUser(user *pb.AdminUser) (*pb.AdminUser, error) {
	// Stub to implement
	return nil, nil
}

// DeleteAdminUser - CouchDB implementation of DeleteAdminUser
func (couchDB *AdminServiceDatastoreCouchDB) DeleteAdminUser(userID string) (*pb.AdminUser, error) {
	// Stub to implement
	return nil, nil
}

// GetAdminUser - CouchDB implementation of GetAdminUser
func (couchDB *AdminServiceDatastoreCouchDB) GetAdminUser(userID string) (*pb.AdminUser, error) {
	// Connect to PROV DB
	db, err := couchdb.NewDatabase(couchDB.dbName)
	if err != nil {
		logger.Log.Errorf("Unable to connect to Prov DB %s: %v\n", couchDB.server, err)
		return nil, err
	}

	logger.Log.Infof("Using db %s to GET Admin User %s\n", couchDB.dbName, userID)

	// Get the Admion User from CouchDB
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
	// Stub to implement
	return nil, nil
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

// Turns an AdminUser object into a map[string]interface{} so that it
// can be stored in CouchDB.
func convertAdminUserToGenericObject(user *pb.AdminUser) (map[string]interface{}, error) {
	userToBytes, err := json.Marshal(user)
	if err != nil {
		logger.Log.Errorf("Unable to convert user to format to persist: %v\n", err)
		return nil, err
	}
	var genericFormat map[string]interface{}
	err = json.Unmarshal(userToBytes, &genericFormat)
	if err != nil {
		logger.Log.Errorf("Unable to convert user to format to persist: %v\n", err)
		return nil, err
	}

	// Successfully converted the User
	return genericFormat, nil
}

// Takes the map[string]interface{} generic data returned by CouchDB and
// converts it to an AdminUser.
func convertGenericObjectToAdminUser(genericUser map[string]interface{}) (*pb.AdminUser, error) {
	genericUserInBytes, err := json.Marshal(genericUser)
	if err != nil {
		fmt.Printf("Error converting generic user data to Admin User type: %v\n", err)
		return nil, err
	}

	res := pb.AdminUser{}
	err = json.Unmarshal(genericUserInBytes, &res)
	if err != nil {
		fmt.Printf("Error converting generic user data to Admin User type: %v\n", err)
		return nil, err
	}

	logger.Log.Debugf("Converted generic data to AdminUser: %v\n", res)

	return &res, nil
}

func (couchDB *AdminServiceDatastoreCouchDB) getCouchDB() (*couchdb.Database, error) {
	db, err := couchdb.NewDatabase(couchDB.dbName)
	if err != nil {
		logger.Log.Errorf("Unable to connect to CouchDB %s: %v\n", couchDB.server, err)
		return nil, err
	}

	return db, nil
}
