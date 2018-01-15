package couchDB

import (
	"fmt"

	"github.com/accedian/adh-gather/config"
	ds "github.com/accedian/adh-gather/datastore"
	"github.com/accedian/adh-gather/gather"
	"github.com/accedian/adh-gather/logger"

	pb "github.com/accedian/adh-gather/gathergrpc"
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
func (tsd *TenantServiceDatastoreCouchDB) CreateTenantUser(tenantUserRequest *pb.TenantUserRequest) (*pb.TenantUserResponse, error) {
	tenantDBName := createDBPathStr(tsd.server, tenantUserRequest.GetData().GetTenantId())
	dataContainer := &pb.TenantUserResponse{}
	if err := createDataInCouch(tenantDBName, tenantUserRequest, dataContainer, string(ds.TenantUserType), ds.TenantUserStr); err != nil {
		return nil, err
	}
	return dataContainer, nil
}

// UpdateTenantUser - CouchDB implementation of UpdateTenantUser
func (tsd *TenantServiceDatastoreCouchDB) UpdateTenantUser(tenantUserRequest *pb.TenantUserRequest) (*pb.TenantUserResponse, error) {
	tenantDBName := createDBPathStr(tsd.server, tenantUserRequest.GetData().GetTenantId())
	dataContainer := &pb.TenantUserResponse{}
	if err := updateDataInCouch(tenantDBName, tenantUserRequest, dataContainer, string(ds.TenantUserType), ds.TenantUserStr); err != nil {
		return nil, err
	}
	return dataContainer, nil
}

// DeleteTenantUser - CouchDB implementation of DeleteTenantUser
func (tsd *TenantServiceDatastoreCouchDB) DeleteTenantUser(tenantUserIDRequest *pb.TenantUserIdRequest) (*pb.TenantUserResponse, error) {
	tenantDBName := createDBPathStr(tsd.server, tenantUserIDRequest.GetTenantId())
	dataContainer := pb.TenantUserResponse{}
	if err := deleteDataFromCouch(tenantDBName, tenantUserIDRequest.GetUserId(), &dataContainer, ds.TenantUserStr); err != nil {
		return nil, err
	}
	return &dataContainer, nil
}

// GetTenantUser - CouchDB implementation of GetTenantUser
func (tsd *TenantServiceDatastoreCouchDB) GetTenantUser(tenantUserIDRequest *pb.TenantUserIdRequest) (*pb.TenantUserResponse, error) {
	tenantDBName := createDBPathStr(tsd.server, tenantUserIDRequest.GetTenantId())
	dataContainer := pb.TenantUserResponse{}
	if err := getDataFromCouch(tenantDBName, tenantUserIDRequest.GetUserId(), &dataContainer, ds.TenantUserStr); err != nil {
		return nil, err
	}
	return &dataContainer, nil
}

// GetAllTenantUsers - CouchDB implementation of GetAllTenantUsers
func (tsd *TenantServiceDatastoreCouchDB) GetAllTenantUsers(tenantID string) (*pb.TenantUserListResponse, error) {
	tenantDBName := createDBPathStr(tsd.server, tenantID)
	res := &pb.TenantUserListResponse{}
	res.Data = make([]*pb.TenantUserResponse, 0)
	if err := getAllOfTypeFromCouch(tenantDBName, string(ds.TenantUserType), ds.TenantUserStr, &res.Data); err != nil {
		return nil, err
	}

	return res, nil
}

// CreateTenantDomain - CouchDB implementation of CreateTenantDomain
func (tsd *TenantServiceDatastoreCouchDB) CreateTenantDomain(tenantDomainRequest *pb.TenantDomainRequest) (*pb.TenantDomainResponse, error) {
	tenantDBName := createDBPathStr(tsd.server, tenantDomainRequest.GetData().GetTenantId())
	dataContainer := &pb.TenantDomainResponse{}
	if err := createDataInCouch(tenantDBName, tenantDomainRequest, dataContainer, string(ds.TenantDomainType), ds.TenantDomainStr); err != nil {
		return nil, err
	}
	return dataContainer, nil
}

// UpdateTenantDomain - CouchDB implementation of UpdateTenantDomain
func (tsd *TenantServiceDatastoreCouchDB) UpdateTenantDomain(tenantDomainRequest *pb.TenantDomainRequest) (*pb.TenantDomainResponse, error) {
	tenantDBName := createDBPathStr(tsd.server, tenantDomainRequest.GetData().GetTenantId())
	dataContainer := &pb.TenantDomainResponse{}
	if err := updateDataInCouch(tenantDBName, tenantDomainRequest, dataContainer, string(ds.TenantDomainType), ds.TenantDomainStr); err != nil {
		return nil, err
	}
	return dataContainer, nil
}

// DeleteTenantDomain - CouchDB implementation of DeleteTenantDomain
func (tsd *TenantServiceDatastoreCouchDB) DeleteTenantDomain(tenantDomainIDRequest *pb.TenantDomainIdRequest) (*pb.TenantDomainResponse, error) {
	tenantDBName := createDBPathStr(tsd.server, tenantDomainIDRequest.GetTenantId())
	dataContainer := pb.TenantDomainResponse{}
	if err := deleteDataFromCouch(tenantDBName, tenantDomainIDRequest.GetDomainId(), &dataContainer, ds.TenantDomainStr); err != nil {
		return nil, err
	}
	return &dataContainer, nil
}

// GetTenantDomain - CouchDB implementation of GetTenantDomain
func (tsd *TenantServiceDatastoreCouchDB) GetTenantDomain(tenantDomainIDRequest *pb.TenantDomainIdRequest) (*pb.TenantDomainResponse, error) {
	tenantDBName := createDBPathStr(tsd.server, tenantDomainIDRequest.GetTenantId())
	dataContainer := pb.TenantDomainResponse{}
	if err := getDataFromCouch(tenantDBName, tenantDomainIDRequest.GetDomainId(), &dataContainer, ds.TenantDomainStr); err != nil {
		return nil, err
	}
	return &dataContainer, nil
}

// GetAllTenantDomains - CouchDB implementation of GetAllTenantDomains
func (tsd *TenantServiceDatastoreCouchDB) GetAllTenantDomains(tenantID string) (*pb.TenantDomainListResponse, error) {
	tenantDBName := createDBPathStr(tsd.server, tenantID)
	res := &pb.TenantDomainListResponse{}
	res.Data = make([]*pb.TenantDomainResponse, 0)
	if err := getAllOfTypeFromCouch(tenantDBName, string(ds.TenantDomainType), ds.TenantDomainStr, &res.Data); err != nil {
		return nil, err
	}

	return res, nil
}

// CreateTenantIngestionProfile - CouchDB implementation of CreateTenantIngestionProfile
func (tsd *TenantServiceDatastoreCouchDB) CreateTenantIngestionProfile(tenantIngPrfReq *pb.TenantIngestionProfileRequest) (*pb.TenantIngestionProfileResponse, error) {
	tenantDBName := createDBPathStr(tsd.server, tenantIngPrfReq.GetData().GetTenantId())
	dataContainer := &pb.TenantIngestionProfileResponse{}
	if err := createDataInCouch(tenantDBName, tenantIngPrfReq, dataContainer, string(ds.TenantIngestionProfileType), ds.TenantIngestionProfileStr); err != nil {
		return nil, err
	}
	return dataContainer, nil
}

// UpdateTenantIngestionProfile - CouchDB implementation of UpdateTenantIngestionProfile
func (tsd *TenantServiceDatastoreCouchDB) UpdateTenantIngestionProfile(tenantIngPrfReq *pb.TenantIngestionProfileRequest) (*pb.TenantIngestionProfileResponse, error) {
	tenantDBName := createDBPathStr(tsd.server, tenantIngPrfReq.GetData().GetTenantId())
	dataContainer := &pb.TenantIngestionProfileResponse{}
	if err := updateDataInCouch(tenantDBName, tenantIngPrfReq, dataContainer, string(ds.TenantIngestionProfileType), ds.TenantIngestionProfileStr); err != nil {
		return nil, err
	}
	return dataContainer, nil
}

// GetTenantIngestionProfile - CouchDB implementation of GetTenantIngestionProfile
func (tsd *TenantServiceDatastoreCouchDB) GetTenantIngestionProfile(tenantIngPrfReq *pb.TenantIngestionProfileIdRequest) (*pb.TenantIngestionProfileResponse, error) {
	tenantDBName := createDBPathStr(tsd.server, tenantIngPrfReq.GetTenantId())
	dataContainer := pb.TenantIngestionProfileResponse{}
	if err := getDataFromCouch(tenantDBName, tenantIngPrfReq.GetIngestionProfileId(), &dataContainer, ds.TenantIngestionProfileStr); err != nil {
		return nil, err
	}
	return &dataContainer, nil
}

// DeleteTenantIngestionProfile - CouchDB implementation of DeleteTenantIngestionProfile
func (tsd *TenantServiceDatastoreCouchDB) DeleteTenantIngestionProfile(tenantIngPrfReq *pb.TenantIngestionProfileIdRequest) (*pb.TenantIngestionProfileResponse, error) {
	tenantDBName := createDBPathStr(tsd.server, tenantIngPrfReq.GetTenantId())
	dataContainer := pb.TenantIngestionProfileResponse{}
	if err := deleteDataFromCouch(tenantDBName, tenantIngPrfReq.GetIngestionProfileId(), &dataContainer, ds.TenantIngestionProfileStr); err != nil {
		return nil, err
	}
	return &dataContainer, nil
}

// CreateTenantThresholdProfile - CouchDB implementation of CreateTenantThresholdProfile
func (tsd *TenantServiceDatastoreCouchDB) CreateTenantThresholdProfile(tenantThreshPrfReq *pb.TenantThresholdProfileRequest) (*pb.TenantThresholdProfileResponse, error) {
	tenantDBName := createDBPathStr(tsd.server, tenantThreshPrfReq.GetData().GetTenantId())
	dataContainer := &pb.TenantThresholdProfileResponse{}
	if err := createDataInCouch(tenantDBName, tenantThreshPrfReq, dataContainer, string(ds.TenantThresholdProfileType), ds.TenantThresholdProfileStr); err != nil {
		return nil, err
	}
	return dataContainer, nil
}

// UpdateTenantThresholdProfile - CouchDB implementation of UpdateTenantThresholdProfile
func (tsd *TenantServiceDatastoreCouchDB) UpdateTenantThresholdProfile(tenantThreshPrfReq *pb.TenantThresholdProfileRequest) (*pb.TenantThresholdProfileResponse, error) {
	tenantDBName := createDBPathStr(tsd.server, tenantThreshPrfReq.GetData().GetTenantId())
	dataContainer := &pb.TenantThresholdProfileResponse{}
	if err := updateDataInCouch(tenantDBName, tenantThreshPrfReq, dataContainer, string(ds.TenantThresholdProfileType), ds.TenantThresholdProfileStr); err != nil {
		return nil, err
	}
	return dataContainer, nil
}

// GetTenantThresholdProfile - CouchDB implementation of GetTenantThresholdProfile
func (tsd *TenantServiceDatastoreCouchDB) GetTenantThresholdProfile(tenantThreshPrfReq *pb.TenantThresholdProfileIdRequest) (*pb.TenantThresholdProfileResponse, error) {
	tenantDBName := createDBPathStr(tsd.server, tenantThreshPrfReq.GetTenantId())
	dataContainer := pb.TenantThresholdProfileResponse{}
	if err := getDataFromCouch(tenantDBName, tenantThreshPrfReq.GetThresholdProfileId(), &dataContainer, ds.TenantThresholdProfileStr); err != nil {
		return nil, err
	}
	return &dataContainer, nil
}

// DeleteTenantThresholdProfile - CouchDB implementation of DeleteTenantThresholdProfile
func (tsd *TenantServiceDatastoreCouchDB) DeleteTenantThresholdProfile(tenantThreshPrfReq *pb.TenantThresholdProfileIdRequest) (*pb.TenantThresholdProfileResponse, error) {
	tenantDBName := createDBPathStr(tsd.server, tenantThreshPrfReq.GetTenantId())
	dataContainer := pb.TenantThresholdProfileResponse{}
	if err := deleteDataFromCouch(tenantDBName, tenantThreshPrfReq.GetThresholdProfileId(), &dataContainer, ds.TenantThresholdProfileStr); err != nil {
		return nil, err
	}
	return &dataContainer, nil
}

// CreateMonitoredObject - CouchDB implementation of CreateMonitoredObject
func (tsd *TenantServiceDatastoreCouchDB) CreateMonitoredObject(monitoredObjectReq *pb.MonitoredObjectRequest) (*pb.MonitoredObjectResponse, error) {
	tenantDBName := createDBPathStr(tsd.server, monitoredObjectReq.GetData().GetTenantId())
	dataContainer := &pb.MonitoredObjectResponse{}
	if err := createDataInCouch(tenantDBName, monitoredObjectReq, dataContainer, string(ds.TenantMonitoredObjectType), ds.TenantMonitoredObjectStr); err != nil {
		return nil, err
	}
	return dataContainer, nil
}

// UpdateMonitoredObject - CouchDB implementation of UpdateMonitoredObject
func (tsd *TenantServiceDatastoreCouchDB) UpdateMonitoredObject(monitoredObjectReq *pb.MonitoredObjectRequest) (*pb.MonitoredObjectResponse, error) {
	tenantDBName := createDBPathStr(tsd.server, monitoredObjectReq.GetData().GetTenantId())
	dataContainer := &pb.MonitoredObjectResponse{}
	if err := updateDataInCouch(tenantDBName, monitoredObjectReq, dataContainer, string(ds.TenantMonitoredObjectType), ds.TenantMonitoredObjectStr); err != nil {
		return nil, err
	}
	return dataContainer, nil
}

// GetMonitoredObject - CouchDB implementation of GetMonitoredObject
func (tsd *TenantServiceDatastoreCouchDB) GetMonitoredObject(monitoredObjectIDReq *pb.MonitoredObjectIdRequest) (*pb.MonitoredObjectResponse, error) {
	tenantDBName := createDBPathStr(tsd.server, monitoredObjectIDReq.GetTenantId())
	dataContainer := pb.MonitoredObjectResponse{}
	if err := getDataFromCouch(tenantDBName, monitoredObjectIDReq.GetMonitoredObjectId(), &dataContainer, ds.TenantMonitoredObjectStr); err != nil {
		return nil, err
	}
	return &dataContainer, nil
}

// DeleteMonitoredObject - CouchDB implementation of DeleteMonitoredObject
func (tsd *TenantServiceDatastoreCouchDB) DeleteMonitoredObject(monitoredObjectIDReq *pb.MonitoredObjectIdRequest) (*pb.MonitoredObjectResponse, error) {
	tenantDBName := createDBPathStr(tsd.server, monitoredObjectIDReq.GetTenantId())
	dataContainer := pb.MonitoredObjectResponse{}
	if err := deleteDataFromCouch(tenantDBName, monitoredObjectIDReq.GetMonitoredObjectId(), &dataContainer, ds.TenantMonitoredObjectStr); err != nil {
		return nil, err
	}
	return &dataContainer, nil
}

// GetAllMonitoredObjects - CouchDB implementation of GetAllMonitoredObjects
func (tsd *TenantServiceDatastoreCouchDB) GetAllMonitoredObjects(tenantID string) (*pb.MonitoredObjectListResponse, error) {
	tenantDBName := createDBPathStr(tsd.server, tenantID)
	res := &pb.MonitoredObjectListResponse{}
	res.Data = make([]*pb.MonitoredObjectResponse, 0)
	if err := getAllOfTypeFromCouch(tenantDBName, string(ds.TenantMonitoredObjectType), ds.TenantMonitoredObjectStr, &res.Data); err != nil {
		return nil, err
	}

	return res, nil
}

// GetMonitoredObjectToDomainMap - CouchDB implementation of GetMonitoredObjectToDomainMap
func (tsd *TenantServiceDatastoreCouchDB) GetMonitoredObjectToDomainMap(moByDomReq *pb.MonitoredObjectCountByDomainRequest) (*pb.MonitoredObjectCountByDomainResponse, error) {
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
		domainMap := map[string]*pb.MonitoredObjectList{}
		for _, row := range rows {
			obj := row.(map[string]interface{})
			key := obj["key"].(string)
			val := obj["value"].(string)
			if domainMap[key] == nil {
				domainMap[key] = &pb.MonitoredObjectList{}
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
	tenantDBName := createDBPathStr(tsd.server, meta.GetData().GetTenantId())
	dataContainer := &pb.TenantMetadata{}
	if err := createDataInCouch(tenantDBName, meta, dataContainer, string(ds.TenantMetaType), ds.TenantMetaStr); err != nil {
		return nil, err
	}
	return dataContainer, nil
}

// UpdateTenantMeta - CouchDB implementation of UpdateTenantMeta
func (tsd *TenantServiceDatastoreCouchDB) UpdateTenantMeta(meta *pb.TenantMetadata) (*pb.TenantMetadata, error) {
	tenantDBName := createDBPathStr(tsd.server, meta.GetData().GetTenantId())
	dataContainer := &pb.TenantMetadata{}
	if err := updateDataInCouch(tenantDBName, meta, dataContainer, string(ds.TenantMetaType), ds.TenantMetaStr); err != nil {
		return nil, err
	}
	return dataContainer, nil
}

// DeleteTenantMeta - CouchDB implementation of DeleteTenantMeta
func (tsd *TenantServiceDatastoreCouchDB) DeleteTenantMeta(tenantID string) (*pb.TenantMetadata, error) {
	logger.Log.Debugf("Deleting %s for %v\n", ds.TenantMetaStr, tenantID)

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
func (tsd *TenantServiceDatastoreCouchDB) GetActiveTenantIngestionProfile(tenantID string) (*pb.TenantIngestionProfileResponse, error) {
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
	res := pb.TenantIngestionProfileResponse{}
	if len(fetchedData) != 0 {
		if err = convertGenericCouchDataToObject(fetchedData[0], &res, ds.TenantIngestionProfileStr); err != nil {
			return nil, err
		}
	}

	logger.Log.Debugf("Found %s %v\n", ds.TenantIngestionProfileStr, res)
	return &res, nil
}

// GetAllTenantThresholdProfile - CouchDB implementation of GetAllTenantThresholdProfile
func (tsd *TenantServiceDatastoreCouchDB) GetAllTenantThresholdProfile(tenantID string) (*pb.TenantThresholdListResponse, error) {
	tenantDBName := createDBPathStr(tsd.server, tenantID)
	res := &pb.TenantThresholdListResponse{}
	res.Data = make([]*pb.TenantThresholdProfileResponse, 0)
	if err := getAllOfTypeFromCouch(tenantDBName, string(ds.TenantThresholdProfileType), ds.TenantThresholdProfileStr, &res.Data); err != nil {
		return nil, err
	}

	return res, nil
}
