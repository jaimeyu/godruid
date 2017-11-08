package couchDB

import (
	"fmt"
	"time"

	"github.com/accedian/adh-gather/datastore"
	"github.com/accedian/adh-gather/gather"
	"github.com/accedian/adh-gather/logger"

	pb "github.com/accedian/adh-gather/gathergrpc"
)

const tenantUserType string = "tenantUser"
const tenantDomainType string = "tenantDomain"
const tenantIngPrfType string = "tenantIngPrf"

// TenantServiceDatastoreCouchDB - struct responsible for handling
// database operations for the Tenant Service when using CouchDB
// as the storage option.
type TenantServiceDatastoreCouchDB struct {
	server string
}

// CreateTenantServiceDAO - instantiates a CouchDB implementation of the
// TenantServiceDatastore.
func CreateTenantServiceDAO() *TenantServiceDatastoreCouchDB {
	result := new(TenantServiceDatastoreCouchDB)
	cfg, err := gather.GetActiveConfig()
	if err != nil {
		logger.Log.Errorf("Falied to instantiate TenantServiceDatastoreCouchDB: %v", err)
	}

	provDBURL := fmt.Sprintf("%s:%d",
		cfg.ServerConfig.Datastore.BindIP,
		cfg.ServerConfig.Datastore.BindPort)
	logger.Log.Debug("Tenant Service CouchDB URL is: ", provDBURL)
	result.server = provDBURL

	return result
}

// CreateTenantUser - CouchDB implementation of CreateTenantUser
func (tsd *TenantServiceDatastoreCouchDB) CreateTenantUser(tenantUserRequest *pb.TenantUserRequest) (*pb.TenantUser, error) {
	tenantDBName := CreateDBPathStr(tsd.server, tenantUserRequest.GetTenantId())
	db, err := GetDatabase(tenantDBName)
	if err != nil {
		return nil, err
	}

	// Give the user a known id, type, and timestamps:
	user := tenantUserRequest.GetUser()
	user.Id = user.GetUsername()
	user.CreatedTimestamp = time.Now().Unix()
	user.LastModifiedTimestamp = user.GetCreatedTimestamp()
	user.Datatype = tenantUserType

	// Marshal the user and read the bytes as string.
	storeFormat, err := ConvertDataToCouchDbSupportedModel(user)
	if err != nil {
		return nil, err
	}

	// Store the user in CouchDB
	_, rev, err := StoreDataInCouchDB(storeFormat, datastore.TenantUserStr, db)
	if err != nil {
		return nil, err
	}

	// Add the revision number to the response
	user.Rev = rev

	// Return the provisioned user.
	logger.Log.Infof("Created %s: %v\n", datastore.TenantUserStr, user)
	return user, nil
}

// UpdateTenantUser - CouchDB implementation of UpdateTenantUser
func (tsd *TenantServiceDatastoreCouchDB) UpdateTenantUser(tenantUserRequest *pb.TenantUserRequest) (*pb.TenantUser, error) {
	tenantDBName := CreateDBPathStr(tsd.server, tenantUserRequest.GetTenantId())
	db, err := GetDatabase(tenantDBName)
	if err != nil {
		return nil, err
	}

	// Update timestamp and make sure the type is properly set:
	user := tenantUserRequest.GetUser()
	user.LastModifiedTimestamp = time.Now().Unix()
	user.Datatype = tenantUserType

	// Marshal the user and read the bytes as string.
	storeFormat, err := ConvertDataToCouchDbSupportedModel(user)
	if err != nil {
		return nil, err
	}

	// Store the user in CouchDB
	id, rev, err := StoreDataInCouchDB(storeFormat, datastore.TenantUserStr, db)
	if err != nil {
		return nil, err
	}

	// Add the evision number to the response
	user.Rev = rev
	logger.Log.Debugf("Successfully updated %s %s with rev %s", datastore.TenantUserStr, id, rev)

	// Return the provisioned user.
	logger.Log.Infof("Updated %s: %v\n", datastore.TenantUserStr, user)
	return user, nil
}

// DeleteTenantUser - CouchDB implementation of DeleteTenantUser
func (tsd *TenantServiceDatastoreCouchDB) DeleteTenantUser(tenantUserIDRequest *pb.TenantUserIdRequest) (*pb.TenantUser, error) {

	// Obtain the value of the existing record for a return value.
	existingUser, err := tsd.GetTenantUser(tenantUserIDRequest)
	if err != nil {
		logger.Log.Errorf("Unable to delete %s: %v\n", datastore.TenantUserStr, err)
		return nil, err
	}

	// Perform the delete operation on CouchDB
	tenantDBName := CreateDBPathStr(tsd.server, tenantUserIDRequest.GetTenantId())
	db, err := GetDatabase(tenantDBName)
	if err != nil {
		return nil, err
	}

	if err = DeleteByDocID(tenantUserIDRequest.GetUserId(), datastore.TenantUserStr, db); err != nil {
		return nil, err
	}

	return existingUser, nil
}

// GetTenantUser - CouchDB implementation of GetTenantUser
func (tsd *TenantServiceDatastoreCouchDB) GetTenantUser(tenantUserIDRequest *pb.TenantUserIdRequest) (*pb.TenantUser, error) {
	tenantDBName := CreateDBPathStr(tsd.server, tenantUserIDRequest.GetTenantId())
	db, err := GetDatabase(tenantDBName)
	if err != nil {
		return nil, err
	}

	// Retrieve the user data from CouchDB
	fetchedUser, err := GetByDocID(tenantUserIDRequest.GetUserId(), datastore.TenantUserStr, db)
	if err != nil {
		return nil, err
	}

	// Marshal the response from the datastore to bytes so that it
	// can be Marshalled back to the proper type.
	res := pb.TenantUser{}
	err = ConvertGenericCouchDataToObject(fetchedUser, &res, datastore.TenantUserStr)
	if err != nil {
		return nil, err
	}

	return &res, nil
}

// GetAllTenantUsers - CouchDB implementation of GetAllTenantUsers
func (tsd *TenantServiceDatastoreCouchDB) GetAllTenantUsers(tenantID string) ([]*pb.TenantUser, error) {
	tenantDBName := CreateDBPathStr(tsd.server, tenantID)
	db, err := GetDatabase(tenantDBName)
	if err != nil {
		return nil, err
	}

	fetchedUserList, err := GetAllOfType(tenantUserType, datastore.TenantUserStr, db)
	if err != nil {
		return nil, err
	}

	// Marshal the response from the datastore to bytes so that it
	// can be Marshalled back to the proper type.
	res, err := convertGenericObjectListToTenantUserList(fetchedUserList)
	if err != nil {
		return nil, err
	}

	return res, nil
}

// CreateTenantDomain - CouchDB implementation of CreateTenantDomain
func (tsd *TenantServiceDatastoreCouchDB) CreateTenantDomain(tenantDomainRequest *pb.TenantDomainRequest) (*pb.TenantDomain, error) {
	tenantDBName := CreateDBPathStr(tsd.server, tenantDomainRequest.GetTenantId())
	db, err := GetDatabase(tenantDBName)
	if err != nil {
		return nil, err
	}

	// Give the domain a known id, type, and timestamps:
	domain := tenantDomainRequest.GetDomain()
	domain.Id = domain.GetName()
	domain.CreatedTimestamp = time.Now().Unix()
	domain.LastModifiedTimestamp = domain.GetCreatedTimestamp()
	domain.Datatype = tenantDomainType

	// Marshal the domain and read the bytes as string.
	storeFormat, err := ConvertDataToCouchDbSupportedModel(domain)
	if err != nil {
		return nil, err
	}

	// Store the domain in CouchDB
	_, rev, err := StoreDataInCouchDB(storeFormat, datastore.TenantDomainStr, db)
	if err != nil {
		return nil, err
	}

	// Add the revision number to the response
	domain.Rev = rev

	// Return the provisioned user.
	logger.Log.Infof("Created %s: %v\n", datastore.TenantDomainStr, domain)
	return domain, nil
}

// UpdateTenantDomain - CouchDB implementation of UpdateTenantDomain
func (tsd *TenantServiceDatastoreCouchDB) UpdateTenantDomain(tenantDomainRequest *pb.TenantDomainRequest) (*pb.TenantDomain, error) {
	tenantDBName := CreateDBPathStr(tsd.server, tenantDomainRequest.GetTenantId())
	db, err := GetDatabase(tenantDBName)
	if err != nil {
		return nil, err
	}

	// Update timestamp and make sure the type is properly set:
	domain := tenantDomainRequest.GetDomain()
	domain.LastModifiedTimestamp = time.Now().Unix()
	domain.Datatype = tenantDomainType

	// Marshal the domain and read the bytes as string.
	storeFormat, err := ConvertDataToCouchDbSupportedModel(domain)
	if err != nil {
		return nil, err
	}

	// Store the domain in CouchDB
	id, rev, err := StoreDataInCouchDB(storeFormat, datastore.TenantDomainStr, db)
	if err != nil {
		return nil, err
	}

	// Add the revision number to the response
	domain.Rev = rev
	logger.Log.Debugf("Successfully updated %s %s with rev %s", datastore.TenantDomainStr, id, rev)

	// Return the provisioned user.
	logger.Log.Infof("Updated %s: %v\n", datastore.TenantDomainStr, domain)
	return domain, nil
}

// DeleteTenantDomain - CouchDB implementation of DeleteTenantDomain
func (tsd *TenantServiceDatastoreCouchDB) DeleteTenantDomain(tenantDomainIDRequest *pb.TenantDomainIdRequest) (*pb.TenantDomain, error) {

	// Obtain the value of the existing record for a return value.
	existingDomain, err := tsd.GetTenantDomain(tenantDomainIDRequest)
	if err != nil {
		logger.Log.Errorf("Unable to delete %s: %v\n", datastore.TenantDomainStr, err)
		return nil, err
	}

	// Perform the delete operation on CouchDB
	tenantDBName := CreateDBPathStr(tsd.server, tenantDomainIDRequest.GetTenantId())
	db, err := GetDatabase(tenantDBName)
	if err != nil {
		return nil, err
	}

	if err = DeleteByDocID(tenantDomainIDRequest.GetDomainId(), datastore.TenantDomainStr, db); err != nil {
		return nil, err
	}

	return existingDomain, nil
}

// GetTenantDomain - CouchDB implementation of GetTenantDomain
func (tsd *TenantServiceDatastoreCouchDB) GetTenantDomain(tenantDomainIDRequest *pb.TenantDomainIdRequest) (*pb.TenantDomain, error) {
	tenantDBName := CreateDBPathStr(tsd.server, tenantDomainIDRequest.GetTenantId())
	db, err := GetDatabase(tenantDBName)
	if err != nil {
		return nil, err
	}

	// Retrieve the domain data from CouchDB
	fetchedDomain, err := GetByDocID(tenantDomainIDRequest.GetDomainId(), datastore.TenantDomainStr, db)
	if err != nil {
		return nil, err
	}

	// Marshal the response from the datastore to bytes so that it
	// can be Marshalled back to the proper type.
	res := pb.TenantDomain{}
	err = ConvertGenericCouchDataToObject(fetchedDomain, &res, datastore.TenantDomainStr)
	if err != nil {
		return nil, err
	}

	return &res, nil
}

// GetAllTenantDomains - CouchDB implementation of GetAllTenantDomains
func (tsd *TenantServiceDatastoreCouchDB) GetAllTenantDomains(tenantID string) ([]*pb.TenantDomain, error) {
	tenantDBName := CreateDBPathStr(tsd.server, tenantID)
	db, err := GetDatabase(tenantDBName)
	if err != nil {
		return nil, err
	}

	fetchedDomainList, err := GetAllOfType(tenantDomainType, datastore.TenantDomainStr, db)
	if err != nil {
		return nil, err
	}

	// Marshal the response from the datastore to bytes so that it
	// can be Marshalled back to the proper type.
	res, err := convertGenericObjectListToTenantDomainList(fetchedDomainList)
	if err != nil {
		return nil, err
	}

	return res, nil
}

// CreateTenantIngestionProfile - CouchDB implementation of CreateTenantIngestionProfile
func (tsd *TenantServiceDatastoreCouchDB) CreateTenantIngestionProfile(tenantIngPrfReq *pb.TenantIngestionProfileRequest) (*pb.TenantIngestionProfile, error) {
	tenantDBName := CreateDBPathStr(tsd.server, tenantIngPrfReq.GetTenantId())
	db, err := GetDatabase(tenantDBName)
	if err != nil {
		return nil, err
	}

	// Give the ingestion profile a known id, type, and timestamps:
	ingPrf := tenantIngPrfReq.GetIngestionProfile()
	ingPrf.Id = tenantIngPrfType
	ingPrf.CreatedTimestamp = time.Now().Unix()
	ingPrf.LastModifiedTimestamp = ingPrf.GetCreatedTimestamp()
	ingPrf.Datatype = tenantIngPrfType

	// Marshal the ingestion profile and read the bytes as string.
	storeFormat, err := ConvertDataToCouchDbSupportedModel(ingPrf)
	if err != nil {
		return nil, err
	}

	// Store the ingestion profile in CouchDB
	_, rev, err := StoreDataInCouchDB(storeFormat, datastore.TenantIngestionProfileStr, db)
	if err != nil {
		return nil, err
	}

	// Add the revision number to the response
	ingPrf.Rev = rev

	// Return the provisioned user.
	logger.Log.Infof("Created %s: %v\n", datastore.TenantIngestionProfileStr, ingPrf)
	return ingPrf, nil
}

// UpdateTenantIngestionProfile - CouchDB implementation of UpdateTenantIngestionProfile
func (tsd *TenantServiceDatastoreCouchDB) UpdateTenantIngestionProfile(tenantIngPrfReq *pb.TenantIngestionProfileRequest) (*pb.TenantIngestionProfile, error) {
	tenantDBName := CreateDBPathStr(tsd.server, tenantIngPrfReq.GetTenantId())
	db, err := GetDatabase(tenantDBName)
	if err != nil {
		return nil, err
	}

	// Update timestamp and make sure the type is properly set:
	ingPrf := tenantIngPrfReq.GetIngestionProfile()
	ingPrf.LastModifiedTimestamp = time.Now().Unix()
	ingPrf.Datatype = tenantIngPrfType

	// Marshal the ingestion profile and read the bytes as string.
	storeFormat, err := ConvertDataToCouchDbSupportedModel(ingPrf)
	if err != nil {
		return nil, err
	}

	// Store the ingestion profile in CouchDB
	id, rev, err := StoreDataInCouchDB(storeFormat, datastore.TenantIngestionProfileStr, db)
	if err != nil {
		return nil, err
	}

	// Add the revision number to the response
	ingPrf.Rev = rev
	logger.Log.Debugf("Successfully updated %s %s with rev %s", datastore.TenantIngestionProfileStr, id, rev)

	// Return the provisioned user.
	logger.Log.Infof("Updated %s: %v\n", datastore.TenantIngestionProfileStr, ingPrf)
	return ingPrf, nil
}

// GetTenantIngestionProfile - CouchDB implementation of GetTenantIngestionProfile
func (tsd *TenantServiceDatastoreCouchDB) GetTenantIngestionProfile(tenantIngPrfReq *pb.TenantIngestionProfileIdRequest) (*pb.TenantIngestionProfile, error) {
	tenantDBName := CreateDBPathStr(tsd.server, tenantIngPrfReq.GetTenantId())
	db, err := GetDatabase(tenantDBName)
	if err != nil {
		return nil, err
	}

	// Retrieve the ingestion profile data from CouchDB
	fetchedIngPrf, err := GetByDocID(tenantIngPrfReq.GetIngestionProfileId(), datastore.TenantIngestionProfileStr, db)
	if err != nil {
		return nil, err
	}

	// Marshal the response from the datastore to bytes so that it
	// can be Marshalled back to the proper type.
	res := pb.TenantIngestionProfile{}
	err = ConvertGenericCouchDataToObject(fetchedIngPrf, &res, datastore.TenantIngestionProfileStr)
	if err != nil {
		return nil, err
	}

	return &res, nil
}

// DeleteTenantIngestionProfile - CouchDB implementation of DeleteTenantIngestionProfile
func (tsd *TenantServiceDatastoreCouchDB) DeleteTenantIngestionProfile(tenantIngPrfReq *pb.TenantIngestionProfileIdRequest) (*pb.TenantIngestionProfile, error) {

	// Obtain the value of the existing record for a return value.
	existingIngPrf, err := tsd.GetTenantIngestionProfile(tenantIngPrfReq)
	if err != nil {
		logger.Log.Errorf("Unable to delete %s: %v\n", datastore.TenantIngestionProfileStr, err)
		return nil, err
	}

	// Perform the delete operation on CouchDB
	tenantDBName := CreateDBPathStr(tsd.server, tenantIngPrfReq.GetTenantId())
	db, err := GetDatabase(tenantDBName)
	if err != nil {
		return nil, err
	}

	if err = DeleteByDocID(tenantIngPrfReq.GetIngestionProfileId(), datastore.TenantIngestionProfileStr, db); err != nil {
		return nil, err
	}

	return existingIngPrf, nil
}

// Takes a set of generic data that contains a list of TenantUsers and converts it to
// and ADH TenantUserList object
func convertGenericObjectListToTenantUserList(genericUserList []map[string]interface{}) ([]*pb.TenantUser, error) {
	res := make([]*pb.TenantUser, 0)
	for _, genericUserObject := range genericUserList {
		user := pb.TenantUser{}
		err := ConvertGenericCouchDataToObject(genericUserObject, &user, datastore.TenantUserStr)
		if err != nil {
			continue
		}
		res = append(res, &user)
	}

	logger.Log.Debugf("Converted generic data to %s List: %v\n", datastore.TenantUserStr, res)

	return res, nil
}

// Takes a set of generic data that contains a list of TenantDomains and converts it to
// and ADH TenantDomainList object
func convertGenericObjectListToTenantDomainList(genericDomainList []map[string]interface{}) ([]*pb.TenantDomain, error) {
	res := make([]*pb.TenantDomain, 0)
	for _, genericDomainObject := range genericDomainList {
		domain := pb.TenantDomain{}
		err := ConvertGenericCouchDataToObject(genericDomainObject, &domain, datastore.TenantDomainStr)
		if err != nil {
			continue
		}
		res = append(res, &domain)
	}

	logger.Log.Debugf("Converted generic data to %s List: %v\n", datastore.TenantDomainStr, res)

	return res, nil
}
