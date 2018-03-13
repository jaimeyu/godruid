package couchDB

import (
	"fmt"

	"github.com/accedian/adh-gather/config"
	ds "github.com/accedian/adh-gather/datastore"
	"github.com/accedian/adh-gather/gather"
	"github.com/accedian/adh-gather/logger"
	"github.com/accedian/adh-gather/models"

	pb "github.com/accedian/adh-gather/gathergrpc"
	admmod "github.com/accedian/adh-gather/models/admin"
	tenmod "github.com/accedian/adh-gather/models/tenant"
	couchdb "github.com/leesper/couchdb-golang"
)

const (
	monitoredObjectsByDomainIndex = "_design/monitoredObjectCount/_view/byDomain"
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
	logger.Log.Debugf("Tenant Service CouchDB URL is: %s", provDBURL)
	result.server = provDBURL

	return result, nil
}

// CreateTenantUser - CouchDB implementation of CreateTenantUser
func (tsd *TenantServiceDatastoreCouchDB) CreateTenantUser(tenantUserRequest *pb.TenantUser) (*pb.TenantUser, error) {
	logger.Log.Debugf("Creating %s: %v\n", tenmod.TenantUserStr, models.AsJSONString(tenantUserRequest))
	tenantUserRequest.XId = ds.GenerateID(tenantUserRequest.GetData(), string(tenmod.TenantUserType))
	tenantID := ds.PrependToDataID(tenantUserRequest.GetData().GetTenantId(), string(admmod.TenantType))

	tenantDBName := createDBPathStr(tsd.server, tenantID)
	dataContainer := &pb.TenantUser{}
	if err := createDataInCouch(tenantDBName, tenantUserRequest, dataContainer, string(tenmod.TenantUserType), tenmod.TenantUserStr); err != nil {
		return nil, err
	}
	logger.Log.Debugf("Cresated %s: %v\n", tenmod.TenantUserStr, models.AsJSONString(dataContainer))
	return dataContainer, nil
}

// UpdateTenantUser - CouchDB implementation of UpdateTenantUser
func (tsd *TenantServiceDatastoreCouchDB) UpdateTenantUser(tenantUserRequest *pb.TenantUser) (*pb.TenantUser, error) {
	logger.Log.Debugf("Updating %s: %v\n", tenmod.TenantUserStr, models.AsJSONString(tenantUserRequest))
	tenantUserRequest.XId = ds.PrependToDataID(tenantUserRequest.XId, string(tenmod.TenantUserType))
	tenantID := ds.PrependToDataID(tenantUserRequest.GetData().GetTenantId(), string(admmod.TenantType))

	tenantDBName := createDBPathStr(tsd.server, tenantID)
	dataContainer := &pb.TenantUser{}
	if err := updateDataInCouch(tenantDBName, tenantUserRequest, dataContainer, string(tenmod.TenantUserType), tenmod.TenantUserStr); err != nil {
		return nil, err
	}
	logger.Log.Debugf("Updated %s: %v\n", tenmod.TenantUserStr, models.AsJSONString(tenantUserRequest))
	return dataContainer, nil
}

// DeleteTenantUser - CouchDB implementation of DeleteTenantUser
func (tsd *TenantServiceDatastoreCouchDB) DeleteTenantUser(tenantUserIDRequest *pb.TenantUserIdRequest) (*pb.TenantUser, error) {
	logger.Log.Debugf("Deleting %s: %v\n", tenmod.TenantUserStr, models.AsJSONString(tenantUserIDRequest))
	tenantUserIDRequest.UserId = ds.PrependToDataID(tenantUserIDRequest.UserId, string(tenmod.TenantUserType))
	tenantUserIDRequest.TenantId = ds.PrependToDataID(tenantUserIDRequest.TenantId, string(admmod.TenantType))

	tenantDBName := createDBPathStr(tsd.server, tenantUserIDRequest.GetTenantId())
	dataContainer := pb.TenantUser{}
	if err := deleteDataFromCouch(tenantDBName, tenantUserIDRequest.GetUserId(), &dataContainer, tenmod.TenantUserStr); err != nil {
		return nil, err
	}
	logger.Log.Debugf("Deleted %s: %v\n", tenmod.TenantUserStr, models.AsJSONString(dataContainer))
	return &dataContainer, nil
}

// GetTenantUser - CouchDB implementation of GetTenantUser
func (tsd *TenantServiceDatastoreCouchDB) GetTenantUser(tenantUserIDRequest *pb.TenantUserIdRequest) (*pb.TenantUser, error) {
	logger.Log.Debugf("Fetching %s: %v\n", tenmod.TenantUserStr, models.AsJSONString(tenantUserIDRequest))
	tenantUserIDRequest.UserId = ds.PrependToDataID(tenantUserIDRequest.UserId, string(tenmod.TenantUserType))
	tenantUserIDRequest.TenantId = ds.PrependToDataID(tenantUserIDRequest.TenantId, string(admmod.TenantType))

	tenantDBName := createDBPathStr(tsd.server, tenantUserIDRequest.GetTenantId())
	dataContainer := pb.TenantUser{}
	if err := getDataFromCouch(tenantDBName, tenantUserIDRequest.GetUserId(), &dataContainer, tenmod.TenantUserStr); err != nil {
		return nil, err
	}
	logger.Log.Debugf("Retrieved %s: %v\n", tenmod.TenantUserStr, models.AsJSONString(dataContainer))
	return &dataContainer, nil
}

// GetAllTenantUsers - CouchDB implementation of GetAllTenantUsers
func (tsd *TenantServiceDatastoreCouchDB) GetAllTenantUsers(tenantID string) (*pb.TenantUserList, error) {
	logger.Log.Debugf("Fetching all %s\n", tenmod.TenantUserStr)
	tenantID = ds.PrependToDataID(tenantID, string(admmod.TenantType))

	tenantDBName := createDBPathStr(tsd.server, tenantID)
	res := &pb.TenantUserList{}
	res.Data = make([]*pb.TenantUser, 0)
	if err := getAllOfTypeFromCouch(tenantDBName, string(tenmod.TenantUserType), tenmod.TenantUserStr, &res.Data); err != nil {
		return nil, err
	}

	logger.Log.Debugf("Retrieved %d %s\n", len(res.Data), tenmod.TenantUserStr)
	return res, nil
}

// CreateTenantDomain - CouchDB implementation of CreateTenantDomain
func (tsd *TenantServiceDatastoreCouchDB) CreateTenantDomain(tenantDomainRequest *pb.TenantDomain) (*pb.TenantDomain, error) {
	logger.Log.Debugf("Creating %s: %v\n", tenmod.TenantDomainStr, models.AsJSONString(tenantDomainRequest))
	tenantDomainRequest.XId = ds.GenerateID(tenantDomainRequest.GetData(), string(tenmod.TenantDomainType))
	tenantID := ds.PrependToDataID(tenantDomainRequest.GetData().GetTenantId(), string(admmod.TenantType))

	tenantDBName := createDBPathStr(tsd.server, tenantID)
	dataContainer := &pb.TenantDomain{}
	if err := createDataInCouch(tenantDBName, tenantDomainRequest, dataContainer, string(tenmod.TenantDomainType), tenmod.TenantDomainStr); err != nil {
		return nil, err
	}
	logger.Log.Debugf("Created %s: %v\n", tenmod.TenantDomainStr, models.AsJSONString(dataContainer))
	return dataContainer, nil
}

// UpdateTenantDomain - CouchDB implementation of UpdateTenantDomain
func (tsd *TenantServiceDatastoreCouchDB) UpdateTenantDomain(tenantDomainRequest *pb.TenantDomain) (*pb.TenantDomain, error) {
	logger.Log.Debugf("Updating %s: %v\n", tenmod.TenantDomainStr, models.AsJSONString(tenantDomainRequest))
	tenantDomainRequest.XId = ds.PrependToDataID(tenantDomainRequest.XId, string(tenmod.TenantDomainType))
	tenantID := ds.PrependToDataID(tenantDomainRequest.GetData().GetTenantId(), string(admmod.TenantType))

	tenantDBName := createDBPathStr(tsd.server, tenantID)
	dataContainer := &pb.TenantDomain{}
	if err := updateDataInCouch(tenantDBName, tenantDomainRequest, dataContainer, string(tenmod.TenantDomainType), tenmod.TenantDomainStr); err != nil {
		return nil, err
	}
	logger.Log.Debugf("Updated %s: %v\n", tenmod.TenantDomainStr, models.AsJSONString(dataContainer))
	return dataContainer, nil
}

// DeleteTenantDomain - CouchDB implementation of DeleteTenantDomain
func (tsd *TenantServiceDatastoreCouchDB) DeleteTenantDomain(tenantDomainIDRequest *pb.TenantDomainIdRequest) (*pb.TenantDomain, error) {
	logger.Log.Debugf("Deleting %s: %v\n", tenmod.TenantDomainStr, models.AsJSONString(tenantDomainIDRequest))
	tenantDomainIDRequest.DomainId = ds.PrependToDataID(tenantDomainIDRequest.DomainId, string(tenmod.TenantDomainType))
	tenantDomainIDRequest.TenantId = ds.PrependToDataID(tenantDomainIDRequest.TenantId, string(admmod.TenantType))

	tenantDBName := createDBPathStr(tsd.server, tenantDomainIDRequest.GetTenantId())
	dataContainer := pb.TenantDomain{}
	if err := deleteDataFromCouch(tenantDBName, tenantDomainIDRequest.GetDomainId(), &dataContainer, tenmod.TenantDomainStr); err != nil {
		return nil, err
	}
	logger.Log.Debugf("Deleted %s: %v\n", tenmod.TenantDomainStr, models.AsJSONString(dataContainer))
	return &dataContainer, nil
}

// GetTenantDomain - CouchDB implementation of GetTenantDomain
func (tsd *TenantServiceDatastoreCouchDB) GetTenantDomain(tenantDomainIDRequest *pb.TenantDomainIdRequest) (*pb.TenantDomain, error) {
	logger.Log.Debugf("Fetching %s: %v\n", tenmod.TenantDomainStr, models.AsJSONString(tenantDomainIDRequest))
	tenantDomainIDRequest.DomainId = ds.PrependToDataID(tenantDomainIDRequest.DomainId, string(tenmod.TenantDomainType))
	tenantDomainIDRequest.TenantId = ds.PrependToDataID(tenantDomainIDRequest.TenantId, string(admmod.TenantType))

	tenantDBName := createDBPathStr(tsd.server, tenantDomainIDRequest.GetTenantId())
	dataContainer := pb.TenantDomain{}
	if err := getDataFromCouch(tenantDBName, tenantDomainIDRequest.GetDomainId(), &dataContainer, tenmod.TenantDomainStr); err != nil {
		return nil, err
	}
	logger.Log.Debugf("Retrieved %s: %v\n", tenmod.TenantDomainStr, models.AsJSONString(dataContainer))
	return &dataContainer, nil
}

// GetAllTenantDomains - CouchDB implementation of GetAllTenantDomains
func (tsd *TenantServiceDatastoreCouchDB) GetAllTenantDomains(tenantID string) (*pb.TenantDomainList, error) {
	logger.Log.Debugf("Fetching all %s\n", tenmod.TenantDomainStr)
	tenantID = ds.PrependToDataID(tenantID, string(admmod.TenantType))

	tenantDBName := createDBPathStr(tsd.server, tenantID)
	res := &pb.TenantDomainList{}
	res.Data = make([]*pb.TenantDomain, 0)
	if err := getAllOfTypeFromCouch(tenantDBName, string(tenmod.TenantDomainType), tenmod.TenantDomainStr, &res.Data); err != nil {
		return nil, err
	}

	logger.Log.Debugf("Retrieved %d %s\n", len(res.Data), tenmod.TenantDomainStr)
	return res, nil
}

// CreateTenantIngestionProfile - CouchDB implementation of CreateTenantIngestionProfile
func (tsd *TenantServiceDatastoreCouchDB) CreateTenantIngestionProfile(tenantIngPrfReq *pb.TenantIngestionProfile) (*pb.TenantIngestionProfile, error) {
	logger.Log.Debugf("Creating %s: %v\n", tenmod.TenantIngestionProfileStr, models.AsJSONString(tenantIngPrfReq))
	tenantIngPrfReq.XId = ds.GenerateID(tenantIngPrfReq.GetData(), string(tenmod.TenantIngestionProfileType))
	tenantID := ds.PrependToDataID(tenantIngPrfReq.GetData().GetTenantId(), string(admmod.TenantType))

	tenantDBName := createDBPathStr(tsd.server, tenantID)
	dataContainer := &pb.TenantIngestionProfile{}
	if err := createDataInCouch(tenantDBName, tenantIngPrfReq, dataContainer, string(tenmod.TenantIngestionProfileType), tenmod.TenantIngestionProfileStr); err != nil {
		return nil, err
	}
	logger.Log.Debugf("Created %s: %v\n", tenmod.TenantIngestionProfileStr, models.AsJSONString(dataContainer))
	return dataContainer, nil
}

// UpdateTenantIngestionProfile - CouchDB implementation of UpdateTenantIngestionProfile
func (tsd *TenantServiceDatastoreCouchDB) UpdateTenantIngestionProfile(tenantIngPrfReq *pb.TenantIngestionProfile) (*pb.TenantIngestionProfile, error) {
	logger.Log.Debugf("Updating %s: %v\n", tenmod.TenantIngestionProfileStr, models.AsJSONString(tenantIngPrfReq))
	tenantIngPrfReq.XId = ds.PrependToDataID(tenantIngPrfReq.XId, string(tenmod.TenantIngestionProfileType))
	tenantID := ds.PrependToDataID(tenantIngPrfReq.GetData().GetTenantId(), string(admmod.TenantType))

	tenantDBName := createDBPathStr(tsd.server, tenantID)
	dataContainer := &pb.TenantIngestionProfile{}
	if err := updateDataInCouch(tenantDBName, tenantIngPrfReq, dataContainer, string(tenmod.TenantIngestionProfileType), tenmod.TenantIngestionProfileStr); err != nil {
		return nil, err
	}
	logger.Log.Debugf("Updated %s: %v\n", tenmod.TenantIngestionProfileStr, models.AsJSONString(dataContainer))
	return dataContainer, nil
}

// GetTenantIngestionProfile - CouchDB implementation of GetTenantIngestionProfile
func (tsd *TenantServiceDatastoreCouchDB) GetTenantIngestionProfile(tenantIngPrfReq *pb.TenantIngestionProfileIdRequest) (*pb.TenantIngestionProfile, error) {
	logger.Log.Debugf("Fetching %s: %v\n", tenmod.TenantIngestionProfileStr, models.AsJSONString(tenantIngPrfReq))
	tenantIngPrfReq.IngestionProfileId = ds.PrependToDataID(tenantIngPrfReq.IngestionProfileId, string(tenmod.TenantIngestionProfileType))
	tenantIngPrfReq.TenantId = ds.PrependToDataID(tenantIngPrfReq.TenantId, string(admmod.TenantType))

	tenantDBName := createDBPathStr(tsd.server, tenantIngPrfReq.GetTenantId())
	dataContainer := pb.TenantIngestionProfile{}
	if err := getDataFromCouch(tenantDBName, tenantIngPrfReq.GetIngestionProfileId(), &dataContainer, tenmod.TenantIngestionProfileStr); err != nil {
		return nil, err
	}
	logger.Log.Debugf("Retrieved %s: %v\n", tenmod.TenantIngestionProfileStr, models.AsJSONString(dataContainer))
	return &dataContainer, nil
}

// DeleteTenantIngestionProfile - CouchDB implementation of DeleteTenantIngestionProfile
func (tsd *TenantServiceDatastoreCouchDB) DeleteTenantIngestionProfile(tenantIngPrfReq *pb.TenantIngestionProfileIdRequest) (*pb.TenantIngestionProfile, error) {
	logger.Log.Debugf("Deleting %s: %v\n", tenmod.TenantIngestionProfileStr, models.AsJSONString(tenantIngPrfReq))
	tenantIngPrfReq.IngestionProfileId = ds.PrependToDataID(tenantIngPrfReq.IngestionProfileId, string(tenmod.TenantIngestionProfileType))
	tenantIngPrfReq.TenantId = ds.PrependToDataID(tenantIngPrfReq.TenantId, string(admmod.TenantType))

	tenantDBName := createDBPathStr(tsd.server, tenantIngPrfReq.GetTenantId())
	dataContainer := pb.TenantIngestionProfile{}
	if err := deleteDataFromCouch(tenantDBName, tenantIngPrfReq.GetIngestionProfileId(), &dataContainer, tenmod.TenantIngestionProfileStr); err != nil {
		return nil, err
	}
	logger.Log.Debugf("Deleted %s: %v\n", tenmod.TenantIngestionProfileStr, models.AsJSONString(dataContainer))
	return &dataContainer, nil
}

// CreateTenantThresholdProfile - CouchDB implementation of CreateTenantThresholdProfile
func (tsd *TenantServiceDatastoreCouchDB) CreateTenantThresholdProfile(tenantThreshPrfReq *pb.TenantThresholdProfile) (*pb.TenantThresholdProfile, error) {
	logger.Log.Debugf("Creating %s: %v\n", tenmod.TenantThresholdProfileStr, models.AsJSONString(tenantThreshPrfReq))
	tenantThreshPrfReq.XId = ds.GenerateID(tenantThreshPrfReq.GetData(), string(tenmod.TenantThresholdProfileType))
	tenantID := ds.PrependToDataID(tenantThreshPrfReq.GetData().GetTenantId(), string(admmod.TenantType))

	tenantDBName := createDBPathStr(tsd.server, tenantID)
	dataContainer := &pb.TenantThresholdProfile{}
	if err := createDataInCouch(tenantDBName, tenantThreshPrfReq, dataContainer, string(tenmod.TenantThresholdProfileType), tenmod.TenantThresholdProfileStr); err != nil {
		return nil, err
	}
	logger.Log.Debugf("Created %s: %v\n", tenmod.TenantThresholdProfileStr, models.AsJSONString(dataContainer))
	return dataContainer, nil
}

// UpdateTenantThresholdProfile - CouchDB implementation of UpdateTenantThresholdProfile
func (tsd *TenantServiceDatastoreCouchDB) UpdateTenantThresholdProfile(tenantThreshPrfReq *pb.TenantThresholdProfile) (*pb.TenantThresholdProfile, error) {
	logger.Log.Debugf("Updating %s: %v\n", tenmod.TenantThresholdProfileStr, models.AsJSONString(tenantThreshPrfReq))
	tenantThreshPrfReq.XId = ds.PrependToDataID(tenantThreshPrfReq.XId, string(tenmod.TenantThresholdProfileType))
	tenantID := ds.PrependToDataID(tenantThreshPrfReq.GetData().GetTenantId(), string(admmod.TenantType))

	tenantDBName := createDBPathStr(tsd.server, tenantID)
	dataContainer := &pb.TenantThresholdProfile{}
	if err := updateDataInCouch(tenantDBName, tenantThreshPrfReq, dataContainer, string(tenmod.TenantThresholdProfileType), tenmod.TenantThresholdProfileStr); err != nil {
		return nil, err
	}
	logger.Log.Debugf("Updated %s: %v\n", tenmod.TenantThresholdProfileStr, models.AsJSONString(dataContainer))
	return dataContainer, nil
}

// GetTenantThresholdProfile - CouchDB implementation of GetTenantThresholdProfile
func (tsd *TenantServiceDatastoreCouchDB) GetTenantThresholdProfile(tenantThreshPrfReq *pb.TenantThresholdProfileIdRequest) (*pb.TenantThresholdProfile, error) {
	logger.Log.Debugf("Fetching %s: %v\n", tenmod.TenantThresholdProfileStr, models.AsJSONString(tenantThreshPrfReq))
	tenantThreshPrfReq.ThresholdProfileId = ds.PrependToDataID(tenantThreshPrfReq.ThresholdProfileId, string(tenmod.TenantThresholdProfileType))
	tenantThreshPrfReq.TenantId = ds.PrependToDataID(tenantThreshPrfReq.TenantId, string(admmod.TenantType))

	tenantDBName := createDBPathStr(tsd.server, tenantThreshPrfReq.GetTenantId())
	dataContainer := pb.TenantThresholdProfile{}
	if err := getDataFromCouch(tenantDBName, tenantThreshPrfReq.GetThresholdProfileId(), &dataContainer, tenmod.TenantThresholdProfileStr); err != nil {
		return nil, err
	}
	logger.Log.Debugf("Retrieved %s: %v\n", tenmod.TenantThresholdProfileStr, models.AsJSONString(dataContainer))
	return &dataContainer, nil
}

// DeleteTenantThresholdProfile - CouchDB implementation of DeleteTenantThresholdProfile
func (tsd *TenantServiceDatastoreCouchDB) DeleteTenantThresholdProfile(tenantThreshPrfReq *pb.TenantThresholdProfileIdRequest) (*pb.TenantThresholdProfile, error) {
	logger.Log.Debugf("Deleting %s: %v\n", tenmod.TenantThresholdProfileStr, models.AsJSONString(tenantThreshPrfReq))
	tenantThreshPrfReq.ThresholdProfileId = ds.PrependToDataID(tenantThreshPrfReq.ThresholdProfileId, string(tenmod.TenantThresholdProfileType))
	tenantThreshPrfReq.TenantId = ds.PrependToDataID(tenantThreshPrfReq.TenantId, string(admmod.TenantType))

	tenantDBName := createDBPathStr(tsd.server, tenantThreshPrfReq.GetTenantId())
	dataContainer := pb.TenantThresholdProfile{}
	if err := deleteDataFromCouch(tenantDBName, tenantThreshPrfReq.GetThresholdProfileId(), &dataContainer, tenmod.TenantThresholdProfileStr); err != nil {
		return nil, err
	}
	logger.Log.Debugf("Deleted %s: %v\n", tenmod.TenantThresholdProfileStr, models.AsJSONString(dataContainer))
	return &dataContainer, nil
}

// CreateMonitoredObject - CouchDB implementation of CreateMonitoredObject
func (tsd *TenantServiceDatastoreCouchDB) CreateMonitoredObject(monitoredObjectReq *pb.MonitoredObject) (*pb.MonitoredObject, error) {
	logger.Log.Debugf("Creating %s: %v\n", tenmod.TenantMonitoredObjectStr, models.AsJSONString(monitoredObjectReq))
	monitoredObjectReq.XId = ds.GenerateID(monitoredObjectReq.GetData(), string(tenmod.TenantMonitoredObjectType))
	tenantID := ds.PrependToDataID(monitoredObjectReq.GetData().GetTenantId(), string(admmod.TenantType))

	tenantDBName := createDBPathStr(tsd.server, tenantID)
	dataContainer := &pb.MonitoredObject{}
	if err := createDataInCouch(tenantDBName, monitoredObjectReq, dataContainer, string(tenmod.TenantMonitoredObjectType), tenmod.TenantMonitoredObjectStr); err != nil {
		return nil, err
	}
	logger.Log.Debugf("Created %s: %v\n", tenmod.TenantMonitoredObjectStr, models.AsJSONString(dataContainer))
	return dataContainer, nil
}

// UpdateMonitoredObject - CouchDB implementation of UpdateMonitoredObject
func (tsd *TenantServiceDatastoreCouchDB) UpdateMonitoredObject(monitoredObjectReq *pb.MonitoredObject) (*pb.MonitoredObject, error) {
	logger.Log.Debugf("Updating %s: %v\n", tenmod.TenantMonitoredObjectStr, models.AsJSONString(monitoredObjectReq))
	monitoredObjectReq.XId = ds.PrependToDataID(monitoredObjectReq.XId, string(tenmod.TenantMonitoredObjectType))
	tenantID := ds.PrependToDataID(monitoredObjectReq.GetData().GetTenantId(), string(admmod.TenantType))

	tenantDBName := createDBPathStr(tsd.server, tenantID)
	dataContainer := &pb.MonitoredObject{}
	if err := updateDataInCouch(tenantDBName, monitoredObjectReq, dataContainer, string(tenmod.TenantMonitoredObjectType), tenmod.TenantMonitoredObjectStr); err != nil {
		return nil, err
	}
	logger.Log.Debugf("Updated %s: %v\n", tenmod.TenantMonitoredObjectStr, models.AsJSONString(dataContainer))
	return dataContainer, nil
}

// GetMonitoredObject - CouchDB implementation of GetMonitoredObject
func (tsd *TenantServiceDatastoreCouchDB) GetMonitoredObject(monitoredObjectIDReq *pb.MonitoredObjectIdRequest) (*pb.MonitoredObject, error) {
	logger.Log.Debugf("Fetching %s: %v\n", tenmod.TenantMonitoredObjectStr, models.AsJSONString(monitoredObjectIDReq))
	monitoredObjectIDReq.MonitoredObjectId = ds.PrependToDataID(monitoredObjectIDReq.MonitoredObjectId, string(tenmod.TenantMonitoredObjectType))
	monitoredObjectIDReq.TenantId = ds.PrependToDataID(monitoredObjectIDReq.TenantId, string(admmod.TenantType))

	tenantDBName := createDBPathStr(tsd.server, monitoredObjectIDReq.GetTenantId())
	dataContainer := pb.MonitoredObject{}
	if err := getDataFromCouch(tenantDBName, monitoredObjectIDReq.GetMonitoredObjectId(), &dataContainer, tenmod.TenantMonitoredObjectStr); err != nil {
		return nil, err
	}
	logger.Log.Debugf("Retrieved %s: %v\n", tenmod.TenantMonitoredObjectStr, models.AsJSONString(dataContainer))
	return &dataContainer, nil
}

// DeleteMonitoredObject - CouchDB implementation of DeleteMonitoredObject
func (tsd *TenantServiceDatastoreCouchDB) DeleteMonitoredObject(monitoredObjectIDReq *pb.MonitoredObjectIdRequest) (*pb.MonitoredObject, error) {
	logger.Log.Debugf("Deleting %s: %v\n", tenmod.TenantMonitoredObjectStr, models.AsJSONString(monitoredObjectIDReq))
	monitoredObjectIDReq.MonitoredObjectId = ds.PrependToDataID(monitoredObjectIDReq.MonitoredObjectId, string(tenmod.TenantMonitoredObjectType))
	monitoredObjectIDReq.TenantId = ds.PrependToDataID(monitoredObjectIDReq.TenantId, string(admmod.TenantType))

	tenantDBName := createDBPathStr(tsd.server, monitoredObjectIDReq.GetTenantId())
	dataContainer := pb.MonitoredObject{}
	if err := deleteDataFromCouch(tenantDBName, monitoredObjectIDReq.GetMonitoredObjectId(), &dataContainer, tenmod.TenantMonitoredObjectStr); err != nil {
		return nil, err
	}
	logger.Log.Debugf("Deleted %s: %v\n", tenmod.TenantMonitoredObjectStr, models.AsJSONString(dataContainer))
	return &dataContainer, nil
}

// GetAllMonitoredObjects - CouchDB implementation of GetAllMonitoredObjects
func (tsd *TenantServiceDatastoreCouchDB) GetAllMonitoredObjects(tenantID string) (*pb.MonitoredObjectList, error) {
	logger.Log.Debugf("Fetching all %s\n", tenmod.TenantMonitoredObjectStr)
	tenantID = ds.PrependToDataID(tenantID, string(admmod.TenantType))

	tenantDBName := createDBPathStr(tsd.server, tenantID)
	res := &pb.MonitoredObjectList{}
	res.Data = make([]*pb.MonitoredObject, 0)
	if err := getAllOfTypeFromCouch(tenantDBName, string(tenmod.TenantMonitoredObjectType), tenmod.TenantMonitoredObjectStr, &res.Data); err != nil {
		return nil, err
	}

	logger.Log.Debugf("Retrieved %d %s\n", len(res.Data), tenmod.TenantMonitoredObjectStr)
	return res, nil
}

// GetMonitoredObjectToDomainMap - CouchDB implementation of GetMonitoredObjectToDomainMap
func (tsd *TenantServiceDatastoreCouchDB) GetMonitoredObjectToDomainMap(moByDomReq *pb.MonitoredObjectCountByDomainRequest) (*pb.MonitoredObjectCountByDomainResponse, error) {
	logger.Log.Debugf("Fetching %s: %v\n", tenmod.MonitoredObjectToDomainMapStr, models.AsJSONString(moByDomReq))
	moByDomReq.TenantId = ds.PrependToDataID(moByDomReq.TenantId, string(admmod.TenantType))

	tenantDBName := createDBPathStr(tsd.server, moByDomReq.GetTenantId())
	db, err := getDatabase(tenantDBName)
	if err != nil {
		return nil, err
	}

	// Get response data either by subset, or for all domains
	domainSet := moByDomReq.GetDomainSet()
	var fetchResponse map[string]interface{}
	if domainSet == nil || len(domainSet) == 0 {
		// Retrieve values for all domains
		fetchResponse, err = getByDocID(monitoredObjectsByDomainIndex, tenmod.MonitoredObjectToDomainMapStr, db)
		if err != nil {
			return nil, err
		}
	} else {
		// Retrieve just the subset of values.
		requestBody := map[string]interface{}{}
		requestBody["keys"] = moByDomReq.GetDomainSet()

		fetchResponse, err = fetchDesignDocumentResults(requestBody, tenantDBName, monitoredObjectsByDomainIndex)
		if err != nil {
			return nil, err
		}

	}

	if fetchResponse["rows"] == nil {
		return &pb.MonitoredObjectCountByDomainResponse{}, nil
	}

	// Response will vary depending on if it is aggregated or not
	response := pb.MonitoredObjectCountByDomainResponse{}
	rows := fetchResponse["rows"].([]interface{})
	if moByDomReq.GetByCount() {
		// Aggregate the data into a mapping of Domain ID to count.
		domainMap := map[string]int64{}
		for _, row := range rows {
			obj := row.(map[string]interface{})
			key := obj["key"].(string)
			domainMap[key] = domainMap[key] + 1
		}
		response.DomainToMonitoredObjectCountMap = domainMap
	} else {
		// Return the results as a map of Domain name to values.
		domainMap := map[string]*pb.MonitoredObjectSet{}
		for _, row := range rows {
			obj := row.(map[string]interface{})
			key := obj["key"].(string)
			val := ds.GetDataIDFromFullID(obj["value"].(string))
			if domainMap[key] == nil {
				domainMap[key] = &pb.MonitoredObjectSet{}
			}
			domainMap[key].MonitoredObjectSet = append(domainMap[key].GetMonitoredObjectSet(), val)
		}
		response.DomainToMonitoredObjectSetMap = domainMap
	}

	logger.Log.Debugf("Retrieved %s: %v\n", tenmod.MonitoredObjectToDomainMapStr, models.AsJSONString(response))
	return &response, nil
}

// CreateTenantMeta - CouchDB implementation of CreateTenantMeta
func (tsd *TenantServiceDatastoreCouchDB) CreateTenantMeta(meta *pb.TenantMetadata) (*pb.TenantMetadata, error) {
	logger.Log.Debugf("Creating %s: %v\n", tenmod.TenantMetaStr, models.AsJSONString(meta))
	meta.XId = ds.GenerateID(meta.GetData(), string(tenmod.TenantMetaType))
	tenantID := ds.PrependToDataID(meta.GetData().GetTenantId(), string(admmod.TenantType))

	tenantDBName := createDBPathStr(tsd.server, tenantID)
	dataContainer := &pb.TenantMetadata{}
	if err := createDataInCouch(tenantDBName, meta, dataContainer, string(tenmod.TenantMetaType), tenmod.TenantMetaStr); err != nil {
		return nil, err
	}
	logger.Log.Debugf("Created %s: %v\n", tenmod.TenantMetaStr, models.AsJSONString(dataContainer))
	return dataContainer, nil
}

// UpdateTenantMeta - CouchDB implementation of UpdateTenantMeta
func (tsd *TenantServiceDatastoreCouchDB) UpdateTenantMeta(meta *pb.TenantMetadata) (*pb.TenantMetadata, error) {
	logger.Log.Debugf("Updating %s: %v\n", tenmod.TenantMetaStr, models.AsJSONString(meta))
	meta.XId = ds.PrependToDataID(meta.XId, string(tenmod.TenantMetaType))
	tenantID := ds.PrependToDataID(meta.GetData().GetTenantId(), string(admmod.TenantType))

	tenantDBName := createDBPathStr(tsd.server, tenantID)
	dataContainer := &pb.TenantMetadata{}
	if err := updateDataInCouch(tenantDBName, meta, dataContainer, string(tenmod.TenantMetaType), tenmod.TenantMetaStr); err != nil {
		return nil, err
	}
	logger.Log.Debugf("Updated %s: %v\n", tenmod.TenantMetaStr, models.AsJSONString(dataContainer))
	return dataContainer, nil
}

// DeleteTenantMeta - CouchDB implementation of DeleteTenantMeta
func (tsd *TenantServiceDatastoreCouchDB) DeleteTenantMeta(tenantID string) (*pb.TenantMetadata, error) {
	logger.Log.Debugf("Deleting %s: %s\n", tenmod.TenantMetaStr, tenantID)
	tenantID = ds.PrependToDataID(tenantID, string(admmod.TenantType))

	// Obtain the value of the existing record for a return value.
	existingObject, err := tsd.GetTenantMeta(tenantID)
	if err != nil {
		logger.Log.Debugf("Unable to fetch %s to delete: %s", tenmod.TenantMetaStr, err.Error())
		return nil, err
	}

	tenantDBName := createDBPathStr(tsd.server, tenantID)
	if err := deleteData(tenantDBName, existingObject.GetXId(), tenmod.TenantMetaStr); err != nil {
		logger.Log.Debugf("Unable to delete %s: %s", tenmod.TenantMetaStr, err.Error())
		return nil, err
	}

	// Return the deleted object.
	logger.Log.Debugf("Deleted %s: %v\n", tenmod.TenantMetaStr, models.AsJSONString(existingObject))
	return existingObject, nil
}

// GetTenantMeta - CouchDB implementation of GetTenantMeta
func (tsd *TenantServiceDatastoreCouchDB) GetTenantMeta(tenantID string) (*pb.TenantMetadata, error) {
	logger.Log.Debugf("Fetching %s: %s\n", tenmod.TenantMetaStr, tenantID)
	tenantID = ds.PrependToDataID(tenantID, string(admmod.TenantType))

	tenantDBName := createDBPathStr(tsd.server, tenantID)
	db, err := getDatabase(tenantDBName)
	if err != nil {
		return nil, err
	}

	fetchedData, err := getAllOfTypeByIDPrefix(string(tenmod.TenantMetaType), tenmod.TenantMetaStr, db)
	if err != nil {
		return nil, err
	}

	// Populate the response
	res := pb.TenantMetadata{}
	if len(fetchedData) != 0 {
		if err = convertGenericCouchDataToObject(fetchedData[0], &res, tenmod.TenantMetaStr); err != nil {
			return nil, err
		}
	}

	logger.Log.Debugf("Retrieved %s: %v\n", tenmod.TenantMetaStr, models.AsJSONString(res))
	return &res, nil
}

// GetActiveTenantIngestionProfile - CouchDB implementation of GetActiveTenantIngestionProfile
func (tsd *TenantServiceDatastoreCouchDB) GetActiveTenantIngestionProfile(tenantID string) (*pb.TenantIngestionProfile, error) {
	logger.Log.Debugf("Fetching active %s for Tenant %s\n", tenmod.TenantIngestionProfileStr, tenantID)
	tenantID = ds.PrependToDataID(tenantID, string(admmod.TenantType))

	tenantDBName := createDBPathStr(tsd.server, tenantID)
	db, err := getDatabase(tenantDBName)
	if err != nil {
		return nil, err
	}

	fetchedData, err := getAllOfTypeByIDPrefix(string(tenmod.TenantIngestionProfileType), tenmod.TenantIngestionProfileStr, db)
	if err != nil {
		return nil, err
	}

	// Populate the response
	res := pb.TenantIngestionProfile{}
	if len(fetchedData) != 0 {
		if err = convertGenericCouchDataToObject(fetchedData[0], &res, tenmod.TenantIngestionProfileStr); err != nil {
			return nil, err
		}
	}

	logger.Log.Debugf("Retrieved %s: %v\n", tenmod.TenantIngestionProfileStr, models.AsJSONString(res))
	return &res, nil
}

// GetAllTenantThresholdProfile - CouchDB implementation of GetAllTenantThresholdProfile
func (tsd *TenantServiceDatastoreCouchDB) GetAllTenantThresholdProfile(tenantID string) (*pb.TenantThresholdProfileList, error) {
	logger.Log.Debugf("Fetching all %s for Tenant %s\n", tenmod.TenantThresholdProfileStr, tenantID)
	tenantID = ds.PrependToDataID(tenantID, string(admmod.TenantType))

	tenantDBName := createDBPathStr(tsd.server, tenantID)
	res := &pb.TenantThresholdProfileList{}
	res.Data = make([]*pb.TenantThresholdProfile, 0)
	if err := getAllOfTypeFromCouch(tenantDBName, string(tenmod.TenantThresholdProfileType), tenmod.TenantThresholdProfileStr, &res.Data); err != nil {
		return nil, err
	}

	logger.Log.Debugf("Retrieved %d %s\n", tenmod.TenantThresholdProfileStr, models.AsJSONString(res))
	return res, nil
}

// BulkInsertMonitoredObjects - CouchDB implementation of BulkInsertMonitoredObjects
func (tsd *TenantServiceDatastoreCouchDB) BulkInsertMonitoredObjects(value *pb.TenantMonitoredObjectSet) (*pb.BulkOperationResponse, error) {
	logger.Log.Debugf("Bulk creating %s: %v\n", tenmod.TenantMonitoredObjectStr, models.AsJSONString(value))
	tenantID := ds.PrependToDataID(value.TenantId, string(admmod.TenantType))
	tenantDBName := createDBPathStr(tsd.server, tenantID)
	resource, err := couchdb.NewResource(tenantDBName, nil)
	if err != nil {
		return nil, err
	}

	// Iterate over the collection and populate necessary fields
	for _, mo := range value.MonitoredObjectSet {
		dataType := string(tenmod.TenantMonitoredObjectType)
		mo.XId = ds.GenerateID(mo.Data, dataType)
		mo.Data.Datatype = dataType
		mo.Data.CreatedTimestamp = ds.MakeTimestamp()
		mo.Data.LastModifiedTimestamp = mo.Data.GetCreatedTimestamp()
	}
	body := map[string]interface{}{
		"docs": value.MonitoredObjectSet}

	fetchedData, err := performBulkUpdate(body, resource)
	if err != nil {
		return nil, err
	}

	// Populate the response
	res := pb.BulkOperationResponse{}
	res.Results = make([]*pb.BulkOperationResult, 0)
	for _, fetched := range fetchedData {
		newObj := pb.BulkOperationResult{}
		if err = convertGenericCouchDataToObject(fetched, &newObj, ds.DBBulkUpdateStr); err != nil {
			return nil, err
		}
		newObj.Id = ds.GetDataIDFromFullID(newObj.Id)
		res.Results = append(res.Results, &newObj)
	}

	logger.Log.Debugf("Bulk create of %s result: %v\n", tenmod.TenantMonitoredObjectStr, models.AsJSONString(res))
	return &res, nil
}
