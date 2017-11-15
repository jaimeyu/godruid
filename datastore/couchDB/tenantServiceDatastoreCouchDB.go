package couchDB

import (
	"fmt"
	"time"

	ds "github.com/accedian/adh-gather/datastore"
	"github.com/accedian/adh-gather/gather"
	"github.com/accedian/adh-gather/logger"

	pb "github.com/accedian/adh-gather/gathergrpc"
)

// TenantServiceDatastoreCouchDB - struct responsible for handling
// database operations for the Tenant Service when using CouchDB
// as the storage option.
type TenantServiceDatastoreCouchDB struct {
	server string
}

// CreateTenantServiceDAO - instantiates a CouchDB implementation of the
// TenantServiceDatastore.
func CreateTenantServiceDAO() (*TenantServiceDatastoreCouchDB, error) {
	result := new(TenantServiceDatastoreCouchDB)
	cfg, err := gather.GetActiveConfig()
	if err != nil {
		logger.Log.Debugf("Falied to instantiate TenantServiceDatastoreCouchDB: %s", err.Error())
		return nil, err
	}

	provDBURL := fmt.Sprintf("%s:%d",
		cfg.ServerConfig.Datastore.BindIP,
		cfg.ServerConfig.Datastore.BindPort)
	logger.Log.Debug("Tenant Service CouchDB URL is: ", provDBURL)
	result.server = provDBURL

	return result, nil
}

// CreateTenantUser - CouchDB implementation of CreateTenantUser
func (tsd *TenantServiceDatastoreCouchDB) CreateTenantUser(tenantUserRequest *pb.TenantUserRequest) (*pb.TenantUserResponse, error) {
	tenantDBName := createDBPathStr(tsd.server, tenantUserRequest.GetData().GetTenantId())
	db, err := getDatabase(tenantDBName)
	if err != nil {
		return nil, err
	}

	// Give the user a known type, and timestamps:
	tenantUserRequest.Data.Datatype = string(ds.TenantUserType)
	tenantUserRequest.Data.CreatedTimestamp = time.Now().Unix()
	tenantUserRequest.Data.LastModifiedTimestamp = tenantUserRequest.GetData().GetCreatedTimestamp()

	// Marshal the user and read the bytes as string.
	storeFormat, err := convertDataToCouchDbSupportedModel(tenantUserRequest)
	if err != nil {
		return nil, err
	}

	// Store the user in CouchDB
	_, _, err = storeDataInCouchDB(storeFormat, ds.TenantUserStr, db)
	if err != nil {
		return nil, err
	}

	// Populate the response
	res := pb.TenantUserResponse{}
	if err = convertGenericCouchDataToObject(storeFormat, &res, ds.TenantUserStr); err != nil {
		return nil, err
	}

	// Return the provisioned user.
	logger.Log.Debugf("Created %s: %v\n", ds.TenantUserStr, res)
	return &res, nil
}

// UpdateTenantUser - CouchDB implementation of UpdateTenantUser
func (tsd *TenantServiceDatastoreCouchDB) UpdateTenantUser(tenantUserRequest *pb.TenantUserRequest) (*pb.TenantUserResponse, error) {
	tenantDBName := createDBPathStr(tsd.server, tenantUserRequest.GetData().GetTenantId())
	db, err := getDatabase(tenantDBName)
	if err != nil {
		return nil, err
	}

	// Update timestamp and make sure the type is properly set:
	tenantUserRequest.Data.Datatype = string(ds.TenantUserType)
	tenantUserRequest.Data.LastModifiedTimestamp = time.Now().Unix()

	// Marshal the user and read the bytes as string.
	storeFormat, err := convertDataToCouchDbSupportedModel(tenantUserRequest)
	if err != nil {
		return nil, err
	}

	// Store the user in CouchDB
	_, _, err = storeDataInCouchDB(storeFormat, ds.TenantUserStr, db)
	if err != nil {
		return nil, err
	}

	// Populate the response
	res := pb.TenantUserResponse{}
	if err = convertGenericCouchDataToObject(storeFormat, &res, ds.TenantUserStr); err != nil {
		return nil, err
	}

	// Return the provisioned user.
	logger.Log.Debugf("Updated %s: %v\n", ds.TenantUserStr, res)
	return &res, nil
}

// DeleteTenantUser - CouchDB implementation of DeleteTenantUser
func (tsd *TenantServiceDatastoreCouchDB) DeleteTenantUser(tenantUserIDRequest *pb.TenantUserIdRequest) (*pb.TenantUserResponse, error) {

	// Obtain the value of the existing record for a return value.
	existingUser, err := tsd.GetTenantUser(tenantUserIDRequest)
	if err != nil {
		logger.Log.Debugf("Unable to delete %s: %s", ds.TenantUserStr, err.Error())
		return nil, err
	}

	// Perform the delete operation on CouchDB
	tenantDBName := createDBPathStr(tsd.server, tenantUserIDRequest.GetTenantId())
	db, err := getDatabase(tenantDBName)
	if err != nil {
		return nil, err
	}

	if err = deleteByDocID(tenantUserIDRequest.GetUserId(), ds.TenantUserStr, db); err != nil {
		return nil, err
	}

	logger.Log.Debugf("Deleted %s: %v\n", ds.TenantUserStr, existingUser)
	return existingUser, nil
}

// GetTenantUser - CouchDB implementation of GetTenantUser
func (tsd *TenantServiceDatastoreCouchDB) GetTenantUser(tenantUserIDRequest *pb.TenantUserIdRequest) (*pb.TenantUserResponse, error) {
	tenantDBName := createDBPathStr(tsd.server, tenantUserIDRequest.GetTenantId())
	db, err := getDatabase(tenantDBName)
	if err != nil {
		return nil, err
	}

	// Retrieve the user data from CouchDB
	fetchedUser, err := getByDocID(tenantUserIDRequest.GetUserId(), ds.TenantUserStr, db)
	if err != nil {
		return nil, err
	}

	// Marshal the response from the datastore to bytes so that it
	// can be Marshalled back to the proper type.
	res := pb.TenantUserResponse{}
	if err = convertGenericCouchDataToObject(fetchedUser, &res, ds.TenantUserStr); err != nil {
		return nil, err
	}

	logger.Log.Debugf("Retrieved %s: %v\n", ds.TenantUserStr, res)
	return &res, nil
}

// GetAllTenantUsers - CouchDB implementation of GetAllTenantUsers
func (tsd *TenantServiceDatastoreCouchDB) GetAllTenantUsers(tenantID string) (*pb.TenantUserListResponse, error) {
	tenantDBName := createDBPathStr(tsd.server, tenantID)
	db, err := getDatabase(tenantDBName)
	if err != nil {
		return nil, err
	}

	fetchedUserList, err := getAllOfType(string(ds.TenantUserType), ds.TenantUserStr, db)
	if err != nil {
		return nil, err
	}

	// Marshal the response from the datastore to bytes so that it
	// can be Marshalled back to the proper type.
	res, err := convertGenericObjectListToTenantUserList(fetchedUserList)
	if err != nil {
		return nil, err
	}

	logger.Log.Debugf("Found %d %ss", len(res.GetData()), ds.TenantUserStr)
	return res, nil
}

// CreateTenantDomain - CouchDB implementation of CreateTenantDomain
func (tsd *TenantServiceDatastoreCouchDB) CreateTenantDomain(tenantDomainRequest *pb.TenantDomainRequest) (*pb.TenantDomainResponse, error) {
	tenantDBName := createDBPathStr(tsd.server, tenantDomainRequest.GetData().GetTenantId())
	db, err := getDatabase(tenantDBName)
	if err != nil {
		return nil, err
	}

	// Give the domain a known type, and timestamps:
	tenantDomainRequest.Data.Datatype = string(ds.TenantDomainType)
	tenantDomainRequest.Data.CreatedTimestamp = time.Now().Unix()
	tenantDomainRequest.Data.LastModifiedTimestamp = tenantDomainRequest.GetData().GetCreatedTimestamp()

	// Marshal the domain and read the bytes as string.
	storeFormat, err := convertDataToCouchDbSupportedModel(tenantDomainRequest)
	if err != nil {
		return nil, err
	}

	// Store the domain in CouchDB
	_, _, err = storeDataInCouchDB(storeFormat, ds.TenantDomainStr, db)
	if err != nil {
		return nil, err
	}

	// Populate the response
	res := pb.TenantDomainResponse{}
	if err = convertGenericCouchDataToObject(storeFormat, &res, ds.TenantUserStr); err != nil {
		return nil, err
	}

	// Return the provisioned user.
	logger.Log.Debugf("Created %s: %v\n", ds.TenantDomainStr, res)
	return &res, nil
}

// UpdateTenantDomain - CouchDB implementation of UpdateTenantDomain
func (tsd *TenantServiceDatastoreCouchDB) UpdateTenantDomain(tenantDomainRequest *pb.TenantDomainRequest) (*pb.TenantDomainResponse, error) {
	tenantDBName := createDBPathStr(tsd.server, tenantDomainRequest.GetData().GetTenantId())
	db, err := getDatabase(tenantDBName)
	if err != nil {
		return nil, err
	}

	// Update timestamp and make sure the type is properly set:
	tenantDomainRequest.Data.Datatype = string(ds.TenantDomainType)
	tenantDomainRequest.Data.LastModifiedTimestamp = time.Now().Unix()

	// Marshal the domain and read the bytes as string.
	storeFormat, err := convertDataToCouchDbSupportedModel(tenantDomainRequest)
	if err != nil {
		return nil, err
	}

	// Store the domain in CouchDB
	_, _, err = storeDataInCouchDB(storeFormat, ds.TenantDomainStr, db)
	if err != nil {
		return nil, err
	}

	// Populate the response
	res := pb.TenantDomainResponse{}
	if err = convertGenericCouchDataToObject(storeFormat, &res, ds.TenantUserStr); err != nil {
		return nil, err
	}

	// Return the provisioned user.
	logger.Log.Debugf("Updated %s: %v\n", ds.TenantDomainStr, res)
	return &res, nil
}

// DeleteTenantDomain - CouchDB implementation of DeleteTenantDomain
func (tsd *TenantServiceDatastoreCouchDB) DeleteTenantDomain(tenantDomainIDRequest *pb.TenantDomainIdRequest) (*pb.TenantDomainResponse, error) {

	// Obtain the value of the existing record for a return value.
	existingDomain, err := tsd.GetTenantDomain(tenantDomainIDRequest)
	if err != nil {
		logger.Log.Debugf("Unable to delete %s: %s", ds.TenantDomainStr, err.Error())
		return nil, err
	}

	// Perform the delete operation on CouchDB
	tenantDBName := createDBPathStr(tsd.server, tenantDomainIDRequest.GetTenantId())
	db, err := getDatabase(tenantDBName)
	if err != nil {
		return nil, err
	}

	if err = deleteByDocID(tenantDomainIDRequest.GetDomainId(), ds.TenantDomainStr, db); err != nil {
		return nil, err
	}

	logger.Log.Debugf("Deleted %s: %v\n", ds.TenantDomainStr, existingDomain)
	return existingDomain, nil
}

// GetTenantDomain - CouchDB implementation of GetTenantDomain
func (tsd *TenantServiceDatastoreCouchDB) GetTenantDomain(tenantDomainIDRequest *pb.TenantDomainIdRequest) (*pb.TenantDomainResponse, error) {
	tenantDBName := createDBPathStr(tsd.server, tenantDomainIDRequest.GetTenantId())
	db, err := getDatabase(tenantDBName)
	if err != nil {
		return nil, err
	}

	// Retrieve the domain data from CouchDB
	fetchedDomain, err := getByDocID(tenantDomainIDRequest.GetDomainId(), ds.TenantDomainStr, db)
	if err != nil {
		return nil, err
	}

	// Marshal the response from the datastore to bytes so that it
	// can be Marshalled back to the proper type.
	res := pb.TenantDomainResponse{}
	if err = convertGenericCouchDataToObject(fetchedDomain, &res, ds.TenantDomainStr); err != nil {
		return nil, err
	}

	logger.Log.Debugf("Retrieved %s: %v\n", ds.TenantDomainStr, res)
	return &res, nil
}

// GetAllTenantDomains - CouchDB implementation of GetAllTenantDomains
func (tsd *TenantServiceDatastoreCouchDB) GetAllTenantDomains(tenantID string) (*pb.TenantDomainListResponse, error) {
	tenantDBName := createDBPathStr(tsd.server, tenantID)
	db, err := getDatabase(tenantDBName)
	if err != nil {
		return nil, err
	}

	fetchedDomainList, err := getAllOfTypeByIDPrefix(string(ds.TenantDomainType), ds.TenantDomainStr, db)
	if err != nil {
		return nil, err
	}

	// Marshal the response from the datastore to bytes so that it
	// can be Marshalled back to the proper type.
	res, err := convertGenericObjectListToTenantDomainList(fetchedDomainList)
	if err != nil {
		return nil, err
	}

	logger.Log.Debugf("Found %d %ss\n", len(res.GetData()), ds.TenantDomainStr)
	return res, nil
}

// CreateTenantIngestionProfile - CouchDB implementation of CreateTenantIngestionProfile
func (tsd *TenantServiceDatastoreCouchDB) CreateTenantIngestionProfile(tenantIngPrfReq *pb.TenantIngestionProfileRequest) (*pb.TenantIngestionProfileResponse, error) {
	tenantDBName := createDBPathStr(tsd.server, tenantIngPrfReq.GetData().GetTenantId())
	db, err := getDatabase(tenantDBName)
	if err != nil {
		return nil, err
	}

	// Give the ingestion profile a known type, and timestamps:
	tenantIngPrfReq.Data.Datatype = string(ds.TenantIngestionProfileType)
	tenantIngPrfReq.Data.CreatedTimestamp = time.Now().Unix()
	tenantIngPrfReq.Data.LastModifiedTimestamp = tenantIngPrfReq.GetData().GetCreatedTimestamp()

	// Marshal the ingestion profile and read the bytes as string.
	storeFormat, err := convertDataToCouchDbSupportedModel(tenantIngPrfReq)
	if err != nil {
		return nil, err
	}

	// Store the ingestion profile in CouchDB
	_, _, err = storeDataInCouchDB(storeFormat, ds.TenantIngestionProfileStr, db)
	if err != nil {
		return nil, err
	}

	// Populate the response
	res := pb.TenantIngestionProfileResponse{}
	if err = convertGenericCouchDataToObject(storeFormat, &res, ds.TenantUserStr); err != nil {
		return nil, err
	}

	// Return the provisioned user.
	logger.Log.Debugf("Created %s: %v\n", ds.TenantIngestionProfileStr, res)
	return &res, nil
}

// UpdateTenantIngestionProfile - CouchDB implementation of UpdateTenantIngestionProfile
func (tsd *TenantServiceDatastoreCouchDB) UpdateTenantIngestionProfile(tenantIngPrfReq *pb.TenantIngestionProfileRequest) (*pb.TenantIngestionProfileResponse, error) {
	tenantDBName := createDBPathStr(tsd.server, tenantIngPrfReq.GetData().GetTenantId())
	db, err := getDatabase(tenantDBName)
	if err != nil {
		return nil, err
	}

	// Update timestamp and make sure the type is properly set:
	tenantIngPrfReq.Data.Datatype = string(ds.TenantIngestionProfileType)
	tenantIngPrfReq.Data.LastModifiedTimestamp = time.Now().Unix()

	// Marshal the ingestion profile and read the bytes as string.
	storeFormat, err := convertDataToCouchDbSupportedModel(tenantIngPrfReq)
	if err != nil {
		return nil, err
	}

	// Store the ingestion profile in CouchDB
	_, _, err = storeDataInCouchDB(storeFormat, ds.TenantIngestionProfileStr, db)
	if err != nil {
		return nil, err
	}

	// Populate the response
	res := pb.TenantIngestionProfileResponse{}
	if err = convertGenericCouchDataToObject(storeFormat, &res, ds.TenantUserStr); err != nil {
		return nil, err
	}

	// Return the provisioned user.
	logger.Log.Debugf("Updated %s: %s\n", ds.TenantIngestionProfileStr, res)
	return &res, nil
}

// GetTenantIngestionProfile - CouchDB implementation of GetTenantIngestionProfile
func (tsd *TenantServiceDatastoreCouchDB) GetTenantIngestionProfile(tenantIngPrfReq *pb.TenantIngestionProfileIdRequest) (*pb.TenantIngestionProfileResponse, error) {
	tenantDBName := createDBPathStr(tsd.server, tenantIngPrfReq.GetTenantId())
	db, err := getDatabase(tenantDBName)
	if err != nil {
		return nil, err
	}

	// Retrieve the ingestion profile data from CouchDB
	fetchedIngPrf, err := getByDocID(tenantIngPrfReq.GetIngestionProfileId(), ds.TenantIngestionProfileStr, db)
	if err != nil {
		return nil, err
	}

	// Marshal the response from the datastore to bytes so that it
	// can be Marshalled back to the proper type.
	res := pb.TenantIngestionProfileResponse{}
	if err = convertGenericCouchDataToObject(fetchedIngPrf, &res, ds.TenantIngestionProfileStr); err != nil {
		return nil, err
	}

	logger.Log.Debugf("Retrieved %s: %v\n", ds.TenantIngestionProfileStr, res)
	return &res, nil
}

// DeleteTenantIngestionProfile - CouchDB implementation of DeleteTenantIngestionProfile
func (tsd *TenantServiceDatastoreCouchDB) DeleteTenantIngestionProfile(tenantIngPrfReq *pb.TenantIngestionProfileIdRequest) (*pb.TenantIngestionProfileResponse, error) {

	// Obtain the value of the existing record for a return value.
	existingIngPrf, err := tsd.GetTenantIngestionProfile(tenantIngPrfReq)
	if err != nil {
		logger.Log.Debugf("Unable to delete %s: %s", ds.TenantIngestionProfileStr, err.Error())
		return nil, err
	}

	// Perform the delete operation on CouchDB
	tenantDBName := createDBPathStr(tsd.server, tenantIngPrfReq.GetTenantId())
	db, err := getDatabase(tenantDBName)
	if err != nil {
		return nil, err
	}

	if err = deleteByDocID(tenantIngPrfReq.GetIngestionProfileId(), ds.TenantIngestionProfileStr, db); err != nil {
		return nil, err
	}

	logger.Log.Debugf("Deleted %s: %v\n", ds.TenantIngestionProfileStr, existingIngPrf)
	return existingIngPrf, nil
}

// Takes a set of generic data that contains a list of TenantUsers and converts it to
// and ADH TenantUserList object
func convertGenericObjectListToTenantUserList(genericUserList []map[string]interface{}) (*pb.TenantUserListResponse, error) {
	res := pb.TenantUserListResponse{}
	for _, genericUserObject := range genericUserList {
		user := pb.TenantUserResponse{}
		if err := convertGenericCouchDataToObject(genericUserObject, &user, ds.TenantUserStr); err != nil {
			continue
		}
		res.Data = append(res.GetData(), &user)
	}

	logger.Log.Debugf("Converted generic data to %s List: %v\n", ds.TenantUserStr, res)

	return &res, nil
}

// Takes a set of generic data that contains a list of TenantDomains and converts it to
// and ADH TenantDomainList object
func convertGenericObjectListToTenantDomainList(genericDomainList []map[string]interface{}) (*pb.TenantDomainListResponse, error) {
	res := pb.TenantDomainListResponse{}
	for _, genericDomainObject := range genericDomainList {
		domain := pb.TenantDomainResponse{}
		if err := convertGenericCouchDataToObject(genericDomainObject, &domain, ds.TenantDomainStr); err != nil {
			continue
		}
		res.Data = append(res.GetData(), &domain)
	}

	logger.Log.Debugf("Converted generic data to %s List: %v\n", ds.TenantDomainStr, res)

	return &res, nil
}
