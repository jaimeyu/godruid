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
func (tsd *TenantServiceDatastoreCouchDB) CreateTenantUser(tenantUserRequest *pb.TenantUserRequest) (*pb.TenantUserResponse, error) {
	tenantDBName := CreateDBPathStr(tsd.server, tenantUserRequest.GetData().GetTenantId())
	db, err := GetDatabase(tenantDBName)
	if err != nil {
		return nil, err
	}

	// Give the user a known type, and timestamps:
	tenantUserRequest.Data.Datatype = tenantUserType
	tenantUserRequest.Data.CreatedTimestamp = time.Now().Unix()
	tenantUserRequest.Data.LastModifiedTimestamp = tenantUserRequest.GetData().GetCreatedTimestamp()

	// Marshal the user and read the bytes as string.
	storeFormat, err := ConvertDataToCouchDbSupportedModel(tenantUserRequest)
	if err != nil {
		return nil, err
	}

	// Store the user in CouchDB
	_, _, err = StoreDataInCouchDB(storeFormat, datastore.TenantUserStr, db)
	if err != nil {
		return nil, err
	}

	// Populate the response
	res := pb.TenantUserResponse{}
	err = ConvertGenericCouchDataToObject(storeFormat, &res, datastore.TenantUserStr)
	if err != nil {
		return nil, err
	}

	// Return the provisioned user.
	logger.Log.Infof("Created %s: %v\n", datastore.TenantUserStr, res)
	return &res, nil
}

// UpdateTenantUser - CouchDB implementation of UpdateTenantUser
func (tsd *TenantServiceDatastoreCouchDB) UpdateTenantUser(tenantUserRequest *pb.TenantUserRequest) (*pb.TenantUserResponse, error) {
	tenantDBName := CreateDBPathStr(tsd.server, tenantUserRequest.GetData().GetTenantId())
	db, err := GetDatabase(tenantDBName)
	if err != nil {
		return nil, err
	}

	// Update timestamp and make sure the type is properly set:
	tenantUserRequest.Data.Datatype = tenantUserType
	tenantUserRequest.Data.LastModifiedTimestamp = time.Now().Unix()

	// Marshal the user and read the bytes as string.
	storeFormat, err := ConvertDataToCouchDbSupportedModel(tenantUserRequest)
	if err != nil {
		return nil, err
	}

	// Store the user in CouchDB
	_, _, err = StoreDataInCouchDB(storeFormat, datastore.TenantUserStr, db)
	if err != nil {
		return nil, err
	}

	// Populate the response
	res := pb.TenantUserResponse{}
	err = ConvertGenericCouchDataToObject(storeFormat, &res, datastore.TenantUserStr)
	if err != nil {
		return nil, err
	}

	// Return the provisioned user.
	logger.Log.Infof("Updated %s: %v\n", datastore.TenantUserStr, res)
	return &res, nil
}

// DeleteTenantUser - CouchDB implementation of DeleteTenantUser
func (tsd *TenantServiceDatastoreCouchDB) DeleteTenantUser(tenantUserIDRequest *pb.TenantUserIdRequest) (*pb.TenantUserResponse, error) {

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
func (tsd *TenantServiceDatastoreCouchDB) GetTenantUser(tenantUserIDRequest *pb.TenantUserIdRequest) (*pb.TenantUserResponse, error) {
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
	res := pb.TenantUserResponse{}
	err = ConvertGenericCouchDataToObject(fetchedUser, &res, datastore.TenantUserStr)
	if err != nil {
		return nil, err
	}

	return &res, nil
}

// GetAllTenantUsers - CouchDB implementation of GetAllTenantUsers
func (tsd *TenantServiceDatastoreCouchDB) GetAllTenantUsers(tenantID string) (*pb.TenantUserListResponse, error) {
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
func (tsd *TenantServiceDatastoreCouchDB) CreateTenantDomain(tenantDomainRequest *pb.TenantDomainRequest) (*pb.TenantDomainResponse, error) {
	tenantDBName := CreateDBPathStr(tsd.server, tenantDomainRequest.GetData().GetTenantId())
	db, err := GetDatabase(tenantDBName)
	if err != nil {
		return nil, err
	}

	// Give the domain a known type, and timestamps:
	tenantDomainRequest.Data.Datatype = tenantDomainType
	tenantDomainRequest.Data.CreatedTimestamp = time.Now().Unix()
	tenantDomainRequest.Data.LastModifiedTimestamp = tenantDomainRequest.GetData().GetCreatedTimestamp()

	// Marshal the domain and read the bytes as string.
	storeFormat, err := ConvertDataToCouchDbSupportedModel(tenantDomainRequest)
	if err != nil {
		return nil, err
	}

	// Store the domain in CouchDB
	_, _, err = StoreDataInCouchDB(storeFormat, datastore.TenantDomainStr, db)
	if err != nil {
		return nil, err
	}

	// Populate the response
	res := pb.TenantDomainResponse{}
	err = ConvertGenericCouchDataToObject(storeFormat, &res, datastore.TenantUserStr)
	if err != nil {
		return nil, err
	}

	// Return the provisioned user.
	logger.Log.Infof("Created %s: %v\n", datastore.TenantDomainStr, res)
	return &res, nil
}

// UpdateTenantDomain - CouchDB implementation of UpdateTenantDomain
func (tsd *TenantServiceDatastoreCouchDB) UpdateTenantDomain(tenantDomainRequest *pb.TenantDomainRequest) (*pb.TenantDomainResponse, error) {
	tenantDBName := CreateDBPathStr(tsd.server, tenantDomainRequest.GetData().GetTenantId())
	db, err := GetDatabase(tenantDBName)
	if err != nil {
		return nil, err
	}

	// Update timestamp and make sure the type is properly set:
	tenantDomainRequest.Data.Datatype = tenantDomainType
	tenantDomainRequest.Data.LastModifiedTimestamp = time.Now().Unix()

	// Marshal the domain and read the bytes as string.
	storeFormat, err := ConvertDataToCouchDbSupportedModel(tenantDomainRequest)
	if err != nil {
		return nil, err
	}

	// Store the domain in CouchDB
	_, _, err = StoreDataInCouchDB(storeFormat, datastore.TenantDomainStr, db)
	if err != nil {
		return nil, err
	}

	// Populate the response
	res := pb.TenantDomainResponse{}
	err = ConvertGenericCouchDataToObject(storeFormat, &res, datastore.TenantUserStr)
	if err != nil {
		return nil, err
	}

	// Return the provisioned user.
	logger.Log.Infof("Updated %s: %v\n", datastore.TenantDomainStr, res)
	return &res, nil
}

// DeleteTenantDomain - CouchDB implementation of DeleteTenantDomain
func (tsd *TenantServiceDatastoreCouchDB) DeleteTenantDomain(tenantDomainIDRequest *pb.TenantDomainIdRequest) (*pb.TenantDomainResponse, error) {

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
func (tsd *TenantServiceDatastoreCouchDB) GetTenantDomain(tenantDomainIDRequest *pb.TenantDomainIdRequest) (*pb.TenantDomainResponse, error) {
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
	res := pb.TenantDomainResponse{}
	err = ConvertGenericCouchDataToObject(fetchedDomain, &res, datastore.TenantDomainStr)
	if err != nil {
		return nil, err
	}

	return &res, nil
}

// GetAllTenantDomains - CouchDB implementation of GetAllTenantDomains
func (tsd *TenantServiceDatastoreCouchDB) GetAllTenantDomains(tenantID string) (*pb.TenantDomainListResponse, error) {
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
func (tsd *TenantServiceDatastoreCouchDB) CreateTenantIngestionProfile(tenantIngPrfReq *pb.TenantIngestionProfileRequest) (*pb.TenantIngestionProfileResponse, error) {
	tenantDBName := CreateDBPathStr(tsd.server, tenantIngPrfReq.GetData().GetTenantId())
	db, err := GetDatabase(tenantDBName)
	if err != nil {
		return nil, err
	}

	// Give the ingestion profile a known type, and timestamps:
	tenantIngPrfReq.Data.Datatype = tenantIngPrfType
	tenantIngPrfReq.Data.CreatedTimestamp = time.Now().Unix()
	tenantIngPrfReq.Data.LastModifiedTimestamp = tenantIngPrfReq.GetData().GetCreatedTimestamp()

	// Marshal the ingestion profile and read the bytes as string.
	storeFormat, err := ConvertDataToCouchDbSupportedModel(tenantIngPrfReq)
	if err != nil {
		return nil, err
	}

	// Store the ingestion profile in CouchDB
	_, _, err = StoreDataInCouchDB(storeFormat, datastore.TenantIngestionProfileStr, db)
	if err != nil {
		return nil, err
	}

	// Populate the response
	res := pb.TenantIngestionProfileResponse{}
	err = ConvertGenericCouchDataToObject(storeFormat, &res, datastore.TenantUserStr)
	if err != nil {
		return nil, err
	}

	// Return the provisioned user.
	logger.Log.Infof("Created %s: %v\n", datastore.TenantIngestionProfileStr, res)
	return &res, nil
}

// UpdateTenantIngestionProfile - CouchDB implementation of UpdateTenantIngestionProfile
func (tsd *TenantServiceDatastoreCouchDB) UpdateTenantIngestionProfile(tenantIngPrfReq *pb.TenantIngestionProfileRequest) (*pb.TenantIngestionProfileResponse, error) {
	tenantDBName := CreateDBPathStr(tsd.server, tenantIngPrfReq.GetData().GetTenantId())
	db, err := GetDatabase(tenantDBName)
	if err != nil {
		return nil, err
	}

	// Update timestamp and make sure the type is properly set:
	tenantIngPrfReq.Data.Datatype = tenantIngPrfType
	tenantIngPrfReq.Data.LastModifiedTimestamp = time.Now().Unix()

	// Marshal the ingestion profile and read the bytes as string.
	storeFormat, err := ConvertDataToCouchDbSupportedModel(tenantIngPrfReq)
	if err != nil {
		return nil, err
	}

	// Store the ingestion profile in CouchDB
	_, _, err = StoreDataInCouchDB(storeFormat, datastore.TenantIngestionProfileStr, db)
	if err != nil {
		return nil, err
	}

	// Populate the response
	res := pb.TenantIngestionProfileResponse{}
	err = ConvertGenericCouchDataToObject(storeFormat, &res, datastore.TenantUserStr)
	if err != nil {
		return nil, err
	}

	// Return the provisioned user.
	logger.Log.Infof("Updated %s: %v\n", datastore.TenantIngestionProfileStr, res)
	return &res, nil
}

// GetTenantIngestionProfile - CouchDB implementation of GetTenantIngestionProfile
func (tsd *TenantServiceDatastoreCouchDB) GetTenantIngestionProfile(tenantIngPrfReq *pb.TenantIngestionProfileIdRequest) (*pb.TenantIngestionProfileResponse, error) {
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
	res := pb.TenantIngestionProfileResponse{}
	err = ConvertGenericCouchDataToObject(fetchedIngPrf, &res, datastore.TenantIngestionProfileStr)
	if err != nil {
		return nil, err
	}

	return &res, nil
}

// DeleteTenantIngestionProfile - CouchDB implementation of DeleteTenantIngestionProfile
func (tsd *TenantServiceDatastoreCouchDB) DeleteTenantIngestionProfile(tenantIngPrfReq *pb.TenantIngestionProfileIdRequest) (*pb.TenantIngestionProfileResponse, error) {

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
func convertGenericObjectListToTenantUserList(genericUserList []map[string]interface{}) (*pb.TenantUserListResponse, error) {
	res := pb.TenantUserListResponse{}
	for _, genericUserObject := range genericUserList {
		user := pb.TenantUser{}
		err := ConvertGenericCouchDataToObject(genericUserObject, &user, datastore.TenantUserStr)
		if err != nil {
			continue
		}
		res.List = append(res.GetList(), &user)
	}

	logger.Log.Debugf("Converted generic data to %s List: %v\n", datastore.TenantUserStr, res)

	return &res, nil
}

// Takes a set of generic data that contains a list of TenantDomains and converts it to
// and ADH TenantDomainList object
func convertGenericObjectListToTenantDomainList(genericDomainList []map[string]interface{}) (*pb.TenantDomainListResponse, error) {
	res := pb.TenantDomainListResponse{}
	for _, genericDomainObject := range genericDomainList {
		domain := pb.TenantDomain{}
		err := ConvertGenericCouchDataToObject(genericDomainObject, &domain, datastore.TenantDomainStr)
		if err != nil {
			continue
		}
		res.List = append(res.GetList(), &domain)
	}

	logger.Log.Debugf("Converted generic data to %s List: %v\n", datastore.TenantDomainStr, res)

	return &res, nil
}
