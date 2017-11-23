package couchDB

import (
	"fmt"

	"github.com/accedian/adh-gather/config"
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
	cfg    config.Provider
}

// CreateTenantServiceDAO - instantiates a CouchDB implementation of the
// TenantServiceDatastore.
func CreateTenantServiceDAO() (*TenantServiceDatastoreCouchDB, error) {
	result := new(TenantServiceDatastoreCouchDB)
	result.cfg = gather.GetConfig()

	provDBURL := fmt.Sprintf("%s:%d",
		result.cfg.GetString(gather.CK_server_datastore_ip.String()),
		result.cfg.GetInt(gather.CK_server_datastore_port.String()))
	logger.Log.Debug("Tenant Service CouchDB URL is: ", provDBURL)
	result.server = provDBURL

	return result, nil
}

// CreateTenantUser - CouchDB implementation of CreateTenantUser
func (tsd *TenantServiceDatastoreCouchDB) CreateTenantUser(tenantUserRequest *pb.TenantUserRequest) (*pb.TenantUserResponse, error) {
	logger.Log.Debugf("Creating %s: %v\n", ds.TenantUserStr, tenantUserRequest)

	tenantDBName := createDBPathStr(tsd.server, tenantUserRequest.GetData().GetTenantId())
	dataType := string(ds.TenantUserType)
	dataContainer := pb.TenantUserResponse{}
	if err := storeData(tenantDBName, tenantUserRequest, dataType, ds.TenantUserStr, &dataContainer); err != nil {
		return nil, err
	}

	// Return the provisioned object.
	logger.Log.Debugf("Created %s: %v\n", ds.TenantUserStr, dataContainer)
	return &dataContainer, nil
}

// UpdateTenantUser - CouchDB implementation of UpdateTenantUser
func (tsd *TenantServiceDatastoreCouchDB) UpdateTenantUser(tenantUserRequest *pb.TenantUserRequest) (*pb.TenantUserResponse, error) {
	logger.Log.Debugf("Creating %s: %v\n", ds.TenantUserStr, tenantUserRequest)

	tenantDBName := createDBPathStr(tsd.server, tenantUserRequest.GetData().GetTenantId())
	dataType := string(ds.TenantUserType)
	dataContainer := pb.TenantUserResponse{}
	if err := updateData(tenantDBName, tenantUserRequest, dataType, ds.TenantUserStr, &dataContainer); err != nil {
		return nil, err
	}

	// Return the provisioned object.
	logger.Log.Debugf("Created %s: %v\n", ds.TenantUserStr, dataContainer)
	return &dataContainer, nil
}

// DeleteTenantUser - CouchDB implementation of DeleteTenantUser
func (tsd *TenantServiceDatastoreCouchDB) DeleteTenantUser(tenantUserIDRequest *pb.TenantUserIdRequest) (*pb.TenantUserResponse, error) {

	logger.Log.Debugf("Deleting %s for %v\n", ds.TenantUserStr, tenantUserIDRequest)

	// Obtain the value of the existing record for a return value.
	existingObject, err := tsd.GetTenantUser(tenantUserIDRequest)
	if err != nil {
		logger.Log.Debugf("Unable to fetch %s to delete: %s", ds.TenantUserStr, err.Error())
		return nil, err
	}

	tenantDBName := createDBPathStr(tsd.server, tenantUserIDRequest.GetTenantId())
	if err := deleteData(tenantDBName, tenantUserIDRequest.GetUserId(), ds.TenantUserStr); err != nil {
		logger.Log.Debugf("Unable to delete %s: %s", ds.TenantUserStr, err.Error())
		return nil, err
	}

	// Return the deleted object.
	logger.Log.Debugf("Deleted %s: %v\n", ds.TenantUserStr, existingObject)
	return existingObject, nil
}

// GetTenantUser - CouchDB implementation of GetTenantUser
func (tsd *TenantServiceDatastoreCouchDB) GetTenantUser(tenantUserIDRequest *pb.TenantUserIdRequest) (*pb.TenantUserResponse, error) {
	logger.Log.Debugf("Retrieving %s for %v\n", ds.TenantUserStr, tenantUserIDRequest)

	tenantDBName := createDBPathStr(tsd.server, tenantUserIDRequest.GetTenantId())
	dataContainer := pb.TenantUserResponse{}
	if err := getData(tenantDBName, tenantUserIDRequest.GetUserId(), ds.TenantUserStr, &dataContainer); err != nil {
		return nil, err
	}

	// Return the provisioned object.
	logger.Log.Debugf("Retrieved %s: %v\n", ds.TenantUserStr, dataContainer)
	return &dataContainer, nil
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
	logger.Log.Debugf("Creating %s: %v\n", ds.TenantDomainStr, tenantDomainRequest)

	tenantDBName := createDBPathStr(tsd.server, tenantDomainRequest.GetData().GetTenantId())
	dataType := string(ds.TenantDomainType)
	dataContainer := pb.TenantDomainResponse{}
	if err := storeData(tenantDBName, tenantDomainRequest, dataType, ds.TenantDomainStr, &dataContainer); err != nil {
		return nil, err
	}

	// Return the provisioned object.
	logger.Log.Debugf("Created %s: %v\n", ds.TenantDomainStr, dataContainer)
	return &dataContainer, nil
}

// UpdateTenantDomain - CouchDB implementation of UpdateTenantDomain
func (tsd *TenantServiceDatastoreCouchDB) UpdateTenantDomain(tenantDomainRequest *pb.TenantDomainRequest) (*pb.TenantDomainResponse, error) {
	logger.Log.Debugf("Creating %s: %v\n", ds.TenantDomainStr, tenantDomainRequest)

	tenantDBName := createDBPathStr(tsd.server, tenantDomainRequest.GetData().GetTenantId())
	dataType := string(ds.TenantDomainType)
	dataContainer := pb.TenantDomainResponse{}
	if err := updateData(tenantDBName, tenantDomainRequest, dataType, ds.TenantDomainStr, &dataContainer); err != nil {
		return nil, err
	}

	// Return the provisioned object.
	logger.Log.Debugf("Created %s: %v\n", ds.TenantDomainStr, dataContainer)
	return &dataContainer, nil
}

// DeleteTenantDomain - CouchDB implementation of DeleteTenantDomain
func (tsd *TenantServiceDatastoreCouchDB) DeleteTenantDomain(tenantDomainIDRequest *pb.TenantDomainIdRequest) (*pb.TenantDomainResponse, error) {
	logger.Log.Debugf("Deleting %s for %v\n", ds.TenantDomainStr, tenantDomainIDRequest)

	// Obtain the value of the existing record for a return value.
	existingObject, err := tsd.GetTenantDomain(tenantDomainIDRequest)
	if err != nil {
		logger.Log.Debugf("Unable to fetch %s to delete: %s", ds.TenantDomainStr, err.Error())
		return nil, err
	}

	tenantDBName := createDBPathStr(tsd.server, tenantDomainIDRequest.GetTenantId())
	if err := deleteData(tenantDBName, tenantDomainIDRequest.GetDomainId(), ds.TenantDomainStr); err != nil {
		logger.Log.Debugf("Unable to delete %s: %s", ds.TenantDomainStr, err.Error())
		return nil, err
	}

	// Return the deleted object.
	logger.Log.Debugf("Deleted %s: %v\n", ds.TenantDomainStr, existingObject)
	return existingObject, nil
}

// GetTenantDomain - CouchDB implementation of GetTenantDomain
func (tsd *TenantServiceDatastoreCouchDB) GetTenantDomain(tenantDomainIDRequest *pb.TenantDomainIdRequest) (*pb.TenantDomainResponse, error) {
	logger.Log.Debugf("Retrieving %s for %v\n", ds.TenantDomainStr, tenantDomainIDRequest)

	tenantDBName := createDBPathStr(tsd.server, tenantDomainIDRequest.GetTenantId())
	dataContainer := pb.TenantDomainResponse{}
	if err := getData(tenantDBName, tenantDomainIDRequest.GetDomainId(), ds.TenantDomainStr, &dataContainer); err != nil {
		return nil, err
	}

	// Return the provisioned object.
	logger.Log.Debugf("Retrieved %s: %v\n", ds.TenantDomainStr, dataContainer)
	return &dataContainer, nil
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
	logger.Log.Debugf("Creating %s: %v\n", ds.TenantIngestionProfileStr, tenantIngPrfReq)

	tenantDBName := createDBPathStr(tsd.server, tenantIngPrfReq.GetData().GetTenantId())
	dataType := string(ds.TenantIngestionProfileType)
	dataContainer := pb.TenantIngestionProfileResponse{}
	if err := storeData(tenantDBName, tenantIngPrfReq, dataType, ds.TenantIngestionProfileStr, &dataContainer); err != nil {
		return nil, err
	}

	// Return the provisioned object.
	logger.Log.Debugf("Created %s: %v\n", ds.TenantIngestionProfileStr, dataContainer)
	return &dataContainer, nil
}

// UpdateTenantIngestionProfile - CouchDB implementation of UpdateTenantIngestionProfile
func (tsd *TenantServiceDatastoreCouchDB) UpdateTenantIngestionProfile(tenantIngPrfReq *pb.TenantIngestionProfileRequest) (*pb.TenantIngestionProfileResponse, error) {
	logger.Log.Debugf("Updating %s: %v\n", ds.TenantIngestionProfileStr, tenantIngPrfReq)

	tenantDBName := createDBPathStr(tsd.server, tenantIngPrfReq.GetData().GetTenantId())
	dataType := string(ds.TenantIngestionProfileType)
	dataContainer := pb.TenantIngestionProfileResponse{}
	if err := updateData(tenantDBName, tenantIngPrfReq, dataType, ds.TenantIngestionProfileStr, &dataContainer); err != nil {
		return nil, err
	}

	// Return the provisioned object.
	logger.Log.Debugf("Updated %s: %v\n", ds.TenantIngestionProfileStr, dataContainer)
	return &dataContainer, nil
}

// GetTenantIngestionProfile - CouchDB implementation of GetTenantIngestionProfile
func (tsd *TenantServiceDatastoreCouchDB) GetTenantIngestionProfile(tenantIngPrfReq *pb.TenantIngestionProfileIdRequest) (*pb.TenantIngestionProfileResponse, error) {
	logger.Log.Debugf("Retrieving %s for %v\n", ds.TenantIngestionProfileStr, tenantIngPrfReq)

	tenantDBName := createDBPathStr(tsd.server, tenantIngPrfReq.GetTenantId())
	dataContainer := pb.TenantIngestionProfileResponse{}
	if err := getData(tenantDBName, tenantIngPrfReq.GetIngestionProfileId(), ds.TenantIngestionProfileStr, &dataContainer); err != nil {
		return nil, err
	}

	// Return the provisioned object.
	logger.Log.Debugf("Retrieved %s: %v\n", ds.TenantIngestionProfileStr, dataContainer)
	return &dataContainer, nil
}

// DeleteTenantIngestionProfile - CouchDB implementation of DeleteTenantIngestionProfile
func (tsd *TenantServiceDatastoreCouchDB) DeleteTenantIngestionProfile(tenantIngPrfReq *pb.TenantIngestionProfileIdRequest) (*pb.TenantIngestionProfileResponse, error) {
	logger.Log.Debugf("Deleting %s for %v\n", ds.TenantIngestionProfileStr, tenantIngPrfReq)

	// Obtain the value of the existing record for a return value.
	existingObject, err := tsd.GetTenantIngestionProfile(tenantIngPrfReq)
	if err != nil {
		logger.Log.Debugf("Unable to fetch %s to delete: %s", ds.TenantIngestionProfileStr, err.Error())
		return nil, err
	}

	tenantDBName := createDBPathStr(tsd.server, tenantIngPrfReq.GetTenantId())
	if err := deleteData(tenantDBName, tenantIngPrfReq.GetIngestionProfileId(), ds.TenantIngestionProfileStr); err != nil {
		logger.Log.Debugf("Unable to delete %s: %s", ds.TenantIngestionProfileStr, err.Error())
		return nil, err
	}

	// Return the deleted object.
	logger.Log.Debugf("Deleted %s: %v\n", ds.TenantIngestionProfileStr, existingObject)
	return existingObject, nil
}

// CreateMonitoredObject - CouchDB implementation of CreateMonitoredObject
func (tsd *TenantServiceDatastoreCouchDB) CreateMonitoredObject(monitoredObjectReq *pb.MonitoredObjectRequest) (*pb.MonitoredObjectResponse, error) {
	logger.Log.Debugf("Creating %s: %v\n", ds.TenantMonitoredObjectStr, monitoredObjectReq)

	tenantDBName := createDBPathStr(tsd.server, monitoredObjectReq.GetData().GetTenantId())
	dataType := string(ds.TenantMonitoredObjectType)
	dataContainer := pb.MonitoredObjectResponse{}
	if err := storeData(tenantDBName, monitoredObjectReq, dataType, ds.TenantMonitoredObjectStr, &dataContainer); err != nil {
		return nil, err
	}

	// Return the provisioned object.
	logger.Log.Debugf("Created %s: %v\n", ds.TenantMonitoredObjectStr, dataContainer)
	return &dataContainer, nil
}

// UpdateMonitoredObject - CouchDB implementation of UpdateMonitoredObject
func (tsd *TenantServiceDatastoreCouchDB) UpdateMonitoredObject(monitoredObjectReq *pb.MonitoredObjectRequest) (*pb.MonitoredObjectResponse, error) {
	logger.Log.Debugf("Updating %s: %v\n", ds.TenantMonitoredObjectStr, monitoredObjectReq)

	tenantDBName := createDBPathStr(tsd.server, monitoredObjectReq.GetData().GetTenantId())
	dataType := string(ds.TenantMonitoredObjectType)
	dataContainer := pb.MonitoredObjectResponse{}
	if err := updateData(tenantDBName, monitoredObjectReq, dataType, ds.TenantMonitoredObjectStr, &dataContainer); err != nil {
		return nil, err
	}

	// Return the provisioned object.
	logger.Log.Debugf("Updated %s: %v\n", ds.TenantMonitoredObjectStr, dataContainer)
	return &dataContainer, nil
}

// GetMonitoredObject - CouchDB implementation of GetMonitoredObject
func (tsd *TenantServiceDatastoreCouchDB) GetMonitoredObject(monitoredObjectIDReq *pb.MonitoredObjectIdRequest) (*pb.MonitoredObjectResponse, error) {
	logger.Log.Debugf("Retrieving %s for %v\n", ds.TenantMonitoredObjectStr, monitoredObjectIDReq)

	tenantDBName := createDBPathStr(tsd.server, monitoredObjectIDReq.GetTenantId())
	dataContainer := pb.MonitoredObjectResponse{}
	if err := getData(tenantDBName, monitoredObjectIDReq.GetMonitoredObjectId(), ds.TenantMonitoredObjectStr, &dataContainer); err != nil {
		return nil, err
	}

	// Return the provisioned object.
	logger.Log.Debugf("Retrieved %s: %v\n", ds.TenantMonitoredObjectStr, dataContainer)
	return &dataContainer, nil
}

// DeleteMonitoredObject - CouchDB implementation of DeleteMonitoredObject
func (tsd *TenantServiceDatastoreCouchDB) DeleteMonitoredObject(monitoredObjectIDReq *pb.MonitoredObjectIdRequest) (*pb.MonitoredObjectResponse, error) {
	logger.Log.Debugf("Deleting %s for %v\n", ds.TenantMonitoredObjectStr, monitoredObjectIDReq)

	// Obtain the value of the existing record for a return value.
	existingObject, err := tsd.GetMonitoredObject(monitoredObjectIDReq)
	if err != nil {
		logger.Log.Debugf("Unable to fetch %s to delete: %s", ds.TenantMonitoredObjectStr, err.Error())
		return nil, err
	}

	tenantDBName := createDBPathStr(tsd.server, monitoredObjectIDReq.GetTenantId())
	if err := deleteData(tenantDBName, monitoredObjectIDReq.GetMonitoredObjectId(), ds.TenantMonitoredObjectStr); err != nil {
		logger.Log.Debugf("Unable to delete %s: %s", ds.TenantMonitoredObjectStr, err.Error())
		return nil, err
	}

	// Return the deleted object.
	logger.Log.Debugf("Deleted %s: %v\n", ds.TenantMonitoredObjectStr, existingObject)
	return existingObject, nil
}

// GetAllMonitoredObjects - CouchDB implementation of GetAllMonitoredObjects
func (tsd *TenantServiceDatastoreCouchDB) GetAllMonitoredObjects(tenantID string) (*pb.MonitoredObjectListResponse, error) {
	tenantDBName := createDBPathStr(tsd.server, tenantID)
	db, err := getDatabase(tenantDBName)
	if err != nil {
		return nil, err
	}

	fetchedObjectList, err := getAllOfTypeByIDPrefix(string(ds.TenantMonitoredObjectType), ds.TenantMonitoredObjectStr, db)
	if err != nil {
		return nil, err
	}

	// Marshal the response from the datastore to bytes so that it
	// can be Marshalled back to the proper type.
	res, err := convertGenericObjectListToMonitoredObjectList(fetchedObjectList)
	if err != nil {
		return nil, err
	}

	logger.Log.Debugf("Found %d %ss\n", len(res.GetData()), ds.TenantMonitoredObjectStr)
	return res, nil
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

func convertGenericObjectListToMonitoredObjectList(genericObjectList []map[string]interface{}) (*pb.MonitoredObjectListResponse, error) {
	res := pb.MonitoredObjectListResponse{}
	for _, genericDomainObject := range genericObjectList {
		object := pb.MonitoredObjectResponse{}
		if err := convertGenericCouchDataToObject(genericDomainObject, &object, ds.TenantMonitoredObjectStr); err != nil {
			continue
		}
		res.Data = append(res.GetData(), &object)
	}

	logger.Log.Debugf("Converted generic data to %s List: %v\n", ds.TenantMonitoredObjectStr, res)

	return &res, nil
}
