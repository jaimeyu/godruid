package couchDB

import (
	"fmt"

	"github.com/accedian/adh-gather/config"
	ds "github.com/accedian/adh-gather/datastore"
	"github.com/accedian/adh-gather/gather"
	"github.com/accedian/adh-gather/logger"

	pb "github.com/accedian/adh-gather/gathergrpc"
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
	logger.Log.Debug("Tenant Service CouchDB URL is: ", provDBURL)
	result.server = provDBURL

	return result, nil
}

// CreateTenantUser - CouchDB implementation of CreateTenantUser
func (tsd *TenantServiceDatastoreCouchDB) CreateTenantUser(tenantUserRequest *pb.TenantUser) (*pb.TenantUser, error) {
	tenantUserRequest.XId = ds.GenerateID(tenantUserRequest.GetData(), string(ds.TenantUserType))
	tenantID := ds.PrependToDataID(tenantUserRequest.GetData().GetTenantId(), string(ds.TenantDescriptorType))

	tenantDBName := createDBPathStr(tsd.server, tenantID)
	dataContainer := &pb.TenantUser{}
	if err := createDataInCouch(tenantDBName, tenantUserRequest, dataContainer, string(ds.TenantUserType), ds.TenantUserStr); err != nil {
		return nil, err
	}
	return dataContainer, nil
}

// UpdateTenantUser - CouchDB implementation of UpdateTenantUser
func (tsd *TenantServiceDatastoreCouchDB) UpdateTenantUser(tenantUserRequest *pb.TenantUser) (*pb.TenantUser, error) {
	tenantUserRequest.XId = ds.PrependToDataID(tenantUserRequest.XId, string(ds.TenantUserType))
	tenantID := ds.PrependToDataID(tenantUserRequest.GetData().GetTenantId(), string(ds.TenantDescriptorType))

	tenantDBName := createDBPathStr(tsd.server, tenantID)
	dataContainer := &pb.TenantUser{}
	if err := updateDataInCouch(tenantDBName, tenantUserRequest, dataContainer, string(ds.TenantUserType), ds.TenantUserStr); err != nil {
		return nil, err
	}
	return dataContainer, nil
}

// DeleteTenantUser - CouchDB implementation of DeleteTenantUser
func (tsd *TenantServiceDatastoreCouchDB) DeleteTenantUser(tenantUserIDRequest *pb.TenantUserIdRequest) (*pb.TenantUser, error) {
	tenantUserIDRequest.UserId = ds.PrependToDataID(tenantUserIDRequest.UserId, string(ds.TenantUserType))
	tenantUserIDRequest.TenantId = ds.PrependToDataID(tenantUserIDRequest.TenantId, string(ds.TenantDescriptorType))

	tenantDBName := createDBPathStr(tsd.server, tenantUserIDRequest.GetTenantId())
	dataContainer := pb.TenantUser{}
	if err := deleteDataFromCouch(tenantDBName, tenantUserIDRequest.GetUserId(), &dataContainer, ds.TenantUserStr); err != nil {
		return nil, err
	}
	return &dataContainer, nil
}

// GetTenantUser - CouchDB implementation of GetTenantUser
func (tsd *TenantServiceDatastoreCouchDB) GetTenantUser(tenantUserIDRequest *pb.TenantUserIdRequest) (*pb.TenantUser, error) {
	tenantUserIDRequest.UserId = ds.PrependToDataID(tenantUserIDRequest.UserId, string(ds.TenantUserType))
	tenantUserIDRequest.TenantId = ds.PrependToDataID(tenantUserIDRequest.TenantId, string(ds.TenantDescriptorType))

	tenantDBName := createDBPathStr(tsd.server, tenantUserIDRequest.GetTenantId())
	dataContainer := pb.TenantUser{}
	if err := getDataFromCouch(tenantDBName, tenantUserIDRequest.GetUserId(), &dataContainer, ds.TenantUserStr); err != nil {
		return nil, err
	}
	return &dataContainer, nil
}

// GetAllTenantUsers - CouchDB implementation of GetAllTenantUsers
func (tsd *TenantServiceDatastoreCouchDB) GetAllTenantUsers(tenantID string) (*pb.TenantUserList, error) {
	tenantID = ds.PrependToDataID(tenantID, string(ds.TenantDescriptorType))

	tenantDBName := createDBPathStr(tsd.server, tenantID)
	res := &pb.TenantUserList{}
	res.Data = make([]*pb.TenantUser, 0)
	if err := getAllOfTypeFromCouch(tenantDBName, string(ds.TenantUserType), ds.TenantUserStr, &res.Data); err != nil {
		return nil, err
	}

	return res, nil
}

// CreateTenantDomain - CouchDB implementation of CreateTenantDomain
func (tsd *TenantServiceDatastoreCouchDB) CreateTenantDomain(tenantDomainRequest *pb.TenantDomain) (*pb.TenantDomain, error) {
	tenantDomainRequest.XId = ds.GenerateID(tenantDomainRequest.GetData(), string(ds.TenantDomainType))
	tenantID := ds.PrependToDataID(tenantDomainRequest.GetData().GetTenantId(), string(ds.TenantDescriptorType))

	tenantDBName := createDBPathStr(tsd.server, tenantID)
	dataContainer := &pb.TenantDomain{}
	if err := createDataInCouch(tenantDBName, tenantDomainRequest, dataContainer, string(ds.TenantDomainType), ds.TenantDomainStr); err != nil {
		return nil, err
	}
	return dataContainer, nil
}

// UpdateTenantDomain - CouchDB implementation of UpdateTenantDomain
func (tsd *TenantServiceDatastoreCouchDB) UpdateTenantDomain(tenantDomainRequest *pb.TenantDomain) (*pb.TenantDomain, error) {
	tenantDomainRequest.XId = ds.PrependToDataID(tenantDomainRequest.XId, string(ds.TenantDomainType))
	tenantID := ds.PrependToDataID(tenantDomainRequest.GetData().GetTenantId(), string(ds.TenantDescriptorType))

	tenantDBName := createDBPathStr(tsd.server, tenantID)
	dataContainer := &pb.TenantDomain{}
	if err := updateDataInCouch(tenantDBName, tenantDomainRequest, dataContainer, string(ds.TenantDomainType), ds.TenantDomainStr); err != nil {
		return nil, err
	}
	return dataContainer, nil
}

// DeleteTenantDomain - CouchDB implementation of DeleteTenantDomain
func (tsd *TenantServiceDatastoreCouchDB) DeleteTenantDomain(tenantDomainIDRequest *pb.TenantDomainIdRequest) (*pb.TenantDomain, error) {
	tenantDomainIDRequest.DomainId = ds.PrependToDataID(tenantDomainIDRequest.DomainId, string(ds.TenantDomainType))
	tenantDomainIDRequest.TenantId = ds.PrependToDataID(tenantDomainIDRequest.TenantId, string(ds.TenantDescriptorType))

	tenantDBName := createDBPathStr(tsd.server, tenantDomainIDRequest.GetTenantId())
	dataContainer := pb.TenantDomain{}
	if err := deleteDataFromCouch(tenantDBName, tenantDomainIDRequest.GetDomainId(), &dataContainer, ds.TenantDomainStr); err != nil {
		return nil, err
	}
	return &dataContainer, nil
}

// GetTenantDomain - CouchDB implementation of GetTenantDomain
func (tsd *TenantServiceDatastoreCouchDB) GetTenantDomain(tenantDomainIDRequest *pb.TenantDomainIdRequest) (*pb.TenantDomain, error) {
	tenantDomainIDRequest.DomainId = ds.PrependToDataID(tenantDomainIDRequest.DomainId, string(ds.TenantDomainType))
	tenantDomainIDRequest.TenantId = ds.PrependToDataID(tenantDomainIDRequest.TenantId, string(ds.TenantDescriptorType))

	tenantDBName := createDBPathStr(tsd.server, tenantDomainIDRequest.GetTenantId())
	dataContainer := pb.TenantDomain{}
	if err := getDataFromCouch(tenantDBName, tenantDomainIDRequest.GetDomainId(), &dataContainer, ds.TenantDomainStr); err != nil {
		return nil, err
	}
	return &dataContainer, nil
}

// GetAllTenantDomains - CouchDB implementation of GetAllTenantDomains
func (tsd *TenantServiceDatastoreCouchDB) GetAllTenantDomains(tenantID string) (*pb.TenantDomainList, error) {
	tenantID = ds.PrependToDataID(tenantID, string(ds.TenantDescriptorType))

	tenantDBName := createDBPathStr(tsd.server, tenantID)
	res := &pb.TenantDomainList{}
	res.Data = make([]*pb.TenantDomain, 0)
	if err := getAllOfTypeFromCouch(tenantDBName, string(ds.TenantDomainType), ds.TenantDomainStr, &res.Data); err != nil {
		return nil, err
	}

	return res, nil
}

// CreateTenantIngestionProfile - CouchDB implementation of CreateTenantIngestionProfile
func (tsd *TenantServiceDatastoreCouchDB) CreateTenantIngestionProfile(tenantIngPrfReq *pb.TenantIngestionProfile) (*pb.TenantIngestionProfile, error) {
	tenantIngPrfReq.XId = ds.GenerateID(tenantIngPrfReq.GetData(), string(ds.TenantIngestionProfileType))
	tenantID := ds.PrependToDataID(tenantIngPrfReq.GetData().GetTenantId(), string(ds.TenantDescriptorType))

	tenantDBName := createDBPathStr(tsd.server, tenantID)
	dataContainer := &pb.TenantIngestionProfile{}
	if err := createDataInCouch(tenantDBName, tenantIngPrfReq, dataContainer, string(ds.TenantIngestionProfileType), ds.TenantIngestionProfileStr); err != nil {
		return nil, err
	}
	return dataContainer, nil
}

// UpdateTenantIngestionProfile - CouchDB implementation of UpdateTenantIngestionProfile
func (tsd *TenantServiceDatastoreCouchDB) UpdateTenantIngestionProfile(tenantIngPrfReq *pb.TenantIngestionProfile) (*pb.TenantIngestionProfile, error) {
	tenantIngPrfReq.XId = ds.PrependToDataID(tenantIngPrfReq.XId, string(ds.TenantIngestionProfileType))
	tenantID := ds.PrependToDataID(tenantIngPrfReq.GetData().GetTenantId(), string(ds.TenantDescriptorType))

	tenantDBName := createDBPathStr(tsd.server, tenantID)
	dataContainer := &pb.TenantIngestionProfile{}
	if err := updateDataInCouch(tenantDBName, tenantIngPrfReq, dataContainer, string(ds.TenantIngestionProfileType), ds.TenantIngestionProfileStr); err != nil {
		return nil, err
	}
	return dataContainer, nil
}

// GetTenantIngestionProfile - CouchDB implementation of GetTenantIngestionProfile
func (tsd *TenantServiceDatastoreCouchDB) GetTenantIngestionProfile(tenantIngPrfReq *pb.TenantIngestionProfileIdRequest) (*pb.TenantIngestionProfile, error) {
	tenantIngPrfReq.IngestionProfileId = ds.PrependToDataID(tenantIngPrfReq.IngestionProfileId, string(ds.TenantIngestionProfileType))
	tenantIngPrfReq.TenantId = ds.PrependToDataID(tenantIngPrfReq.TenantId, string(ds.TenantDescriptorType))

	tenantDBName := createDBPathStr(tsd.server, tenantIngPrfReq.GetTenantId())
	dataContainer := pb.TenantIngestionProfile{}
	if err := getDataFromCouch(tenantDBName, tenantIngPrfReq.GetIngestionProfileId(), &dataContainer, ds.TenantIngestionProfileStr); err != nil {
		return nil, err
	}
	return &dataContainer, nil
}

// DeleteTenantIngestionProfile - CouchDB implementation of DeleteTenantIngestionProfile
func (tsd *TenantServiceDatastoreCouchDB) DeleteTenantIngestionProfile(tenantIngPrfReq *pb.TenantIngestionProfileIdRequest) (*pb.TenantIngestionProfile, error) {
	tenantIngPrfReq.IngestionProfileId = ds.PrependToDataID(tenantIngPrfReq.IngestionProfileId, string(ds.TenantIngestionProfileType))
	tenantIngPrfReq.TenantId = ds.PrependToDataID(tenantIngPrfReq.TenantId, string(ds.TenantDescriptorType))

	tenantDBName := createDBPathStr(tsd.server, tenantIngPrfReq.GetTenantId())
	dataContainer := pb.TenantIngestionProfile{}
	if err := deleteDataFromCouch(tenantDBName, tenantIngPrfReq.GetIngestionProfileId(), &dataContainer, ds.TenantIngestionProfileStr); err != nil {
		return nil, err
	}
	return &dataContainer, nil
}

// CreateTenantThresholdProfile - CouchDB implementation of CreateTenantThresholdProfile
func (tsd *TenantServiceDatastoreCouchDB) CreateTenantThresholdProfile(tenantThreshPrfReq *pb.TenantThresholdProfile) (*pb.TenantThresholdProfile, error) {
	tenantThreshPrfReq.XId = ds.GenerateID(tenantThreshPrfReq.GetData(), string(ds.TenantThresholdProfileType))
	tenantID := ds.PrependToDataID(tenantThreshPrfReq.GetData().GetTenantId(), string(ds.TenantDescriptorType))

	tenantDBName := createDBPathStr(tsd.server, tenantID)
	dataContainer := &pb.TenantThresholdProfile{}
	if err := createDataInCouch(tenantDBName, tenantThreshPrfReq, dataContainer, string(ds.TenantThresholdProfileType), ds.TenantThresholdProfileStr); err != nil {
		return nil, err
	}
	return dataContainer, nil
}

// UpdateTenantThresholdProfile - CouchDB implementation of UpdateTenantThresholdProfile
func (tsd *TenantServiceDatastoreCouchDB) UpdateTenantThresholdProfile(tenantThreshPrfReq *pb.TenantThresholdProfile) (*pb.TenantThresholdProfile, error) {
	tenantThreshPrfReq.XId = ds.PrependToDataID(tenantThreshPrfReq.XId, string(ds.TenantThresholdProfileType))
	tenantID := ds.PrependToDataID(tenantThreshPrfReq.GetData().GetTenantId(), string(ds.TenantDescriptorType))

	tenantDBName := createDBPathStr(tsd.server, tenantID)
	dataContainer := &pb.TenantThresholdProfile{}
	if err := updateDataInCouch(tenantDBName, tenantThreshPrfReq, dataContainer, string(ds.TenantThresholdProfileType), ds.TenantThresholdProfileStr); err != nil {
		return nil, err
	}
	return dataContainer, nil
}

// GetTenantThresholdProfile - CouchDB implementation of GetTenantThresholdProfile
func (tsd *TenantServiceDatastoreCouchDB) GetTenantThresholdProfile(tenantThreshPrfReq *pb.TenantThresholdProfileIdRequest) (*pb.TenantThresholdProfile, error) {
	tenantThreshPrfReq.ThresholdProfileId = ds.PrependToDataID(tenantThreshPrfReq.ThresholdProfileId, string(ds.TenantThresholdProfileType))
	tenantThreshPrfReq.TenantId = ds.PrependToDataID(tenantThreshPrfReq.TenantId, string(ds.TenantDescriptorType))

	tenantDBName := createDBPathStr(tsd.server, tenantThreshPrfReq.GetTenantId())
	dataContainer := pb.TenantThresholdProfile{}
	if err := getDataFromCouch(tenantDBName, tenantThreshPrfReq.GetThresholdProfileId(), &dataContainer, ds.TenantThresholdProfileStr); err != nil {
		return nil, err
	}
	return &dataContainer, nil
}

// DeleteTenantThresholdProfile - CouchDB implementation of DeleteTenantThresholdProfile
func (tsd *TenantServiceDatastoreCouchDB) DeleteTenantThresholdProfile(tenantThreshPrfReq *pb.TenantThresholdProfileIdRequest) (*pb.TenantThresholdProfile, error) {
	tenantThreshPrfReq.ThresholdProfileId = ds.PrependToDataID(tenantThreshPrfReq.ThresholdProfileId, string(ds.TenantThresholdProfileType))
	tenantThreshPrfReq.TenantId = ds.PrependToDataID(tenantThreshPrfReq.TenantId, string(ds.TenantDescriptorType))

	tenantDBName := createDBPathStr(tsd.server, tenantThreshPrfReq.GetTenantId())
	dataContainer := pb.TenantThresholdProfile{}
	if err := deleteDataFromCouch(tenantDBName, tenantThreshPrfReq.GetThresholdProfileId(), &dataContainer, ds.TenantThresholdProfileStr); err != nil {
		return nil, err
	}
	return &dataContainer, nil
}

// CreateMonitoredObject - CouchDB implementation of CreateMonitoredObject
func (tsd *TenantServiceDatastoreCouchDB) CreateMonitoredObject(monitoredObjectReq *pb.MonitoredObject) (*pb.MonitoredObject, error) {
	monitoredObjectReq.XId = ds.GenerateID(monitoredObjectReq.GetData(), string(ds.TenantMonitoredObjectType))
	tenantID := ds.PrependToDataID(monitoredObjectReq.GetData().GetTenantId(), string(ds.TenantDescriptorType))

	tenantDBName := createDBPathStr(tsd.server, tenantID)
	dataContainer := &pb.MonitoredObject{}
	if err := createDataInCouch(tenantDBName, monitoredObjectReq, dataContainer, string(ds.TenantMonitoredObjectType), ds.TenantMonitoredObjectStr); err != nil {
		return nil, err
	}
	return dataContainer, nil
}

// UpdateMonitoredObject - CouchDB implementation of UpdateMonitoredObject
func (tsd *TenantServiceDatastoreCouchDB) UpdateMonitoredObject(monitoredObjectReq *pb.MonitoredObject) (*pb.MonitoredObject, error) {
	monitoredObjectReq.XId = ds.PrependToDataID(monitoredObjectReq.XId, string(ds.TenantMonitoredObjectType))
	tenantID := ds.PrependToDataID(monitoredObjectReq.GetData().GetTenantId(), string(ds.TenantDescriptorType))

	tenantDBName := createDBPathStr(tsd.server, tenantID)
	dataContainer := &pb.MonitoredObject{}
	if err := updateDataInCouch(tenantDBName, monitoredObjectReq, dataContainer, string(ds.TenantMonitoredObjectType), ds.TenantMonitoredObjectStr); err != nil {
		return nil, err
	}
	return dataContainer, nil
}

// GetMonitoredObject - CouchDB implementation of GetMonitoredObject
func (tsd *TenantServiceDatastoreCouchDB) GetMonitoredObject(monitoredObjectIDReq *pb.MonitoredObjectIdRequest) (*pb.MonitoredObject, error) {
	monitoredObjectIDReq.MonitoredObjectId = ds.PrependToDataID(monitoredObjectIDReq.MonitoredObjectId, string(ds.TenantMonitoredObjectType))
	monitoredObjectIDReq.TenantId = ds.PrependToDataID(monitoredObjectIDReq.TenantId, string(ds.TenantDescriptorType))

	tenantDBName := createDBPathStr(tsd.server, monitoredObjectIDReq.GetTenantId())
	dataContainer := pb.MonitoredObject{}
	if err := getDataFromCouch(tenantDBName, monitoredObjectIDReq.GetMonitoredObjectId(), &dataContainer, ds.TenantMonitoredObjectStr); err != nil {
		return nil, err
	}
	return &dataContainer, nil
}

// DeleteMonitoredObject - CouchDB implementation of DeleteMonitoredObject
func (tsd *TenantServiceDatastoreCouchDB) DeleteMonitoredObject(monitoredObjectIDReq *pb.MonitoredObjectIdRequest) (*pb.MonitoredObject, error) {
	monitoredObjectIDReq.MonitoredObjectId = ds.PrependToDataID(monitoredObjectIDReq.MonitoredObjectId, string(ds.TenantMonitoredObjectType))
	monitoredObjectIDReq.TenantId = ds.PrependToDataID(monitoredObjectIDReq.TenantId, string(ds.TenantDescriptorType))

	tenantDBName := createDBPathStr(tsd.server, monitoredObjectIDReq.GetTenantId())
	dataContainer := pb.MonitoredObject{}
	if err := deleteDataFromCouch(tenantDBName, monitoredObjectIDReq.GetMonitoredObjectId(), &dataContainer, ds.TenantMonitoredObjectStr); err != nil {
		return nil, err
	}
	return &dataContainer, nil
}

// GetAllMonitoredObjects - CouchDB implementation of GetAllMonitoredObjects
func (tsd *TenantServiceDatastoreCouchDB) GetAllMonitoredObjects(tenantID string) (*pb.MonitoredObjectList, error) {
	tenantID = ds.PrependToDataID(tenantID, string(ds.TenantDescriptorType))

	tenantDBName := createDBPathStr(tsd.server, tenantID)
	res := &pb.MonitoredObjectList{}
	res.Data = make([]*pb.MonitoredObject, 0)
	if err := getAllOfTypeFromCouch(tenantDBName, string(ds.TenantMonitoredObjectType), ds.TenantMonitoredObjectStr, &res.Data); err != nil {
		return nil, err
	}

	return res, nil
}

// GetMonitoredObjectToDomainMap - CouchDB implementation of GetMonitoredObjectToDomainMap
func (tsd *TenantServiceDatastoreCouchDB) GetMonitoredObjectToDomainMap(moByDomReq *pb.MonitoredObjectCountByDomainRequest) (*pb.MonitoredObjectCountByDomainResponse, error) {
	moByDomReq.TenantId = ds.PrependToDataID(moByDomReq.TenantId, string(ds.TenantDescriptorType))

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
		fetchResponse, err = getByDocID(monitoredObjectsByDomainIndex, ds.MonitoredObjectToDomainMapStr, db)
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

	logger.Log.Debugf("Returning %s: %vs\n", ds.MonitoredObjectToDomainMapStr, response)
	return &response, nil
}

// CreateTenantMeta - CouchDB implementation of CreateTenantMeta
func (tsd *TenantServiceDatastoreCouchDB) CreateTenantMeta(meta *pb.TenantMetadata) (*pb.TenantMetadata, error) {
	meta.XId = ds.GenerateID(meta.GetData(), string(ds.TenantMetaType))
	tenantID := ds.PrependToDataID(meta.GetData().GetTenantId(), string(ds.TenantDescriptorType))

	tenantDBName := createDBPathStr(tsd.server, tenantID)
	dataContainer := &pb.TenantMetadata{}
	if err := createDataInCouch(tenantDBName, meta, dataContainer, string(ds.TenantMetaType), ds.TenantMetaStr); err != nil {
		return nil, err
	}
	return dataContainer, nil
}

// UpdateTenantMeta - CouchDB implementation of UpdateTenantMeta
func (tsd *TenantServiceDatastoreCouchDB) UpdateTenantMeta(meta *pb.TenantMetadata) (*pb.TenantMetadata, error) {
	meta.XId = ds.PrependToDataID(meta.XId, string(ds.TenantMetaType))
	tenantID := ds.PrependToDataID(meta.GetData().GetTenantId(), string(ds.TenantDescriptorType))

	tenantDBName := createDBPathStr(tsd.server, tenantID)
	dataContainer := &pb.TenantMetadata{}
	if err := updateDataInCouch(tenantDBName, meta, dataContainer, string(ds.TenantMetaType), ds.TenantMetaStr); err != nil {
		return nil, err
	}
	return dataContainer, nil
}

// DeleteTenantMeta - CouchDB implementation of DeleteTenantMeta
func (tsd *TenantServiceDatastoreCouchDB) DeleteTenantMeta(tenantID string) (*pb.TenantMetadata, error) {
	tenantID = ds.PrependToDataID(tenantID, string(ds.TenantDescriptorType))

	// Obtain the value of the existing record for a return value.
	existingObject, err := tsd.GetTenantMeta(tenantID)
	if err != nil {
		logger.Log.Debugf("Unable to fetch %s to delete: %s", ds.TenantMetaStr, err.Error())
		return nil, err
	}

	tenantDBName := createDBPathStr(tsd.server, tenantID)
	if err := deleteData(tenantDBName, existingObject.GetXId(), ds.TenantMetaStr); err != nil {
		logger.Log.Debugf("Unable to delete %s: %s", ds.TenantMetaStr, err.Error())
		return nil, err
	}

	// Return the deleted object.
	logger.Log.Debugf("Deleted %s: %v\n", ds.TenantMetaStr, existingObject)
	return existingObject, nil
}

// GetTenantMeta - CouchDB implementation of GetTenantMeta
func (tsd *TenantServiceDatastoreCouchDB) GetTenantMeta(tenantID string) (*pb.TenantMetadata, error) {
	tenantID = ds.PrependToDataID(tenantID, string(ds.TenantDescriptorType))

	tenantDBName := createDBPathStr(tsd.server, tenantID)
	db, err := getDatabase(tenantDBName)
	if err != nil {
		return nil, err
	}

	fetchedData, err := getAllOfTypeByIDPrefix(string(ds.TenantMetaType), ds.TenantMetaStr, db)
	if err != nil {
		return nil, err
	}

	// Populate the response
	res := pb.TenantMetadata{}
	if len(fetchedData) != 0 {
		if err = convertGenericCouchDataToObject(fetchedData[0], &res, ds.TenantMetaStr); err != nil {
			return nil, err
		}
	}

	logger.Log.Debugf("Found %s %v\n", ds.TenantMetaStr, res)
	return &res, nil
}

// GetActiveTenantIngestionProfile - CouchDB implementation of GetActiveTenantIngestionProfile
func (tsd *TenantServiceDatastoreCouchDB) GetActiveTenantIngestionProfile(tenantID string) (*pb.TenantIngestionProfile, error) {
	tenantID = ds.PrependToDataID(tenantID, string(ds.TenantDescriptorType))

	tenantDBName := createDBPathStr(tsd.server, tenantID)
	db, err := getDatabase(tenantDBName)
	if err != nil {
		return nil, err
	}

	fetchedData, err := getAllOfTypeByIDPrefix(string(ds.TenantIngestionProfileType), ds.TenantIngestionProfileStr, db)
	if err != nil {
		return nil, err
	}

	// Populate the response
	res := pb.TenantIngestionProfile{}
	if len(fetchedData) != 0 {
		if err = convertGenericCouchDataToObject(fetchedData[0], &res, ds.TenantIngestionProfileStr); err != nil {
			return nil, err
		}
	}

	logger.Log.Debugf("Found %s %v\n", ds.TenantIngestionProfileStr, res)
	return &res, nil
}

// GetAllTenantThresholdProfile - CouchDB implementation of GetAllTenantThresholdProfile
func (tsd *TenantServiceDatastoreCouchDB) GetAllTenantThresholdProfile(tenantID string) (*pb.TenantThresholdProfileList, error) {
	tenantID = ds.PrependToDataID(tenantID, string(ds.TenantDescriptorType))

	tenantDBName := createDBPathStr(tsd.server, tenantID)
	res := &pb.TenantThresholdProfileList{}
	res.Data = make([]*pb.TenantThresholdProfile, 0)
	if err := getAllOfTypeFromCouch(tenantDBName, string(ds.TenantThresholdProfileType), ds.TenantThresholdProfileStr, &res.Data); err != nil {
		return nil, err
	}

	return res, nil
}

// BulkInsertMonitoredObjects - CouchDB implementation of BulkInsertMonitoredObjects
func (tsd *TenantServiceDatastoreCouchDB) BulkInsertMonitoredObjects(value *pb.TenantMonitoredObjectSet) (*pb.BulkOperationResponse, error) {
	tenantID := ds.PrependToDataID(value.TenantId, string(ds.TenantDescriptorType))
	tenantDBName := createDBPathStr(tsd.server, tenantID)
	resource, err := couchdb.NewResource(tenantDBName, nil)
	if err != nil {
		return nil, err
	}

	// Iterate over the collection and populate necessary fields
	for _, mo := range value.MonitoredObjectSet {
		dataType := string(ds.TenantMonitoredObjectType)
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
	
	return &res, nil
}
