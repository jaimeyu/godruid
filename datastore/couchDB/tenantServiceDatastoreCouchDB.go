package couchDB

import (
	"fmt"

	"github.com/accedian/adh-gather/config"
	ds "github.com/accedian/adh-gather/datastore"
	"github.com/accedian/adh-gather/gather"
	"github.com/accedian/adh-gather/logger"
	"github.com/accedian/adh-gather/models"

	admmod "github.com/accedian/adh-gather/models/admin"
	"github.com/accedian/adh-gather/models/common"
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
func (tsd *TenantServiceDatastoreCouchDB) CreateTenantUser(tenantUserRequest *tenmod.User) (*tenmod.User, error) {
	logger.Log.Debugf("Creating %s: %v\n", tenmod.TenantUserStr, models.AsJSONString(tenantUserRequest))
	tenantUserRequest.ID = ds.GenerateID(tenantUserRequest, string(tenmod.TenantUserType))
	tenantID := ds.PrependToDataID(tenantUserRequest.TenantID, string(admmod.TenantType))

	tenantDBName := createDBPathStr(tsd.server, tenantID)
	dataContainer := &tenmod.User{}
	if err := createDataInCouch(tenantDBName, tenantUserRequest, dataContainer, string(tenmod.TenantUserType), tenmod.TenantUserStr); err != nil {
		return nil, err
	}
	logger.Log.Debugf("Cresated %s: %v\n", tenmod.TenantUserStr, models.AsJSONString(dataContainer))
	return dataContainer, nil
}

// UpdateTenantUser - CouchDB implementation of UpdateTenantUser
func (tsd *TenantServiceDatastoreCouchDB) UpdateTenantUser(tenantUserRequest *tenmod.User) (*tenmod.User, error) {
	logger.Log.Debugf("Updating %s: %v\n", tenmod.TenantUserStr, models.AsJSONString(tenantUserRequest))
	tenantUserRequest.ID = ds.PrependToDataID(tenantUserRequest.ID, string(tenmod.TenantUserType))
	tenantID := ds.PrependToDataID(tenantUserRequest.TenantID, string(admmod.TenantType))

	tenantDBName := createDBPathStr(tsd.server, tenantID)
	dataContainer := &tenmod.User{}
	if err := updateDataInCouch(tenantDBName, tenantUserRequest, dataContainer, string(tenmod.TenantUserType), tenmod.TenantUserStr); err != nil {
		return nil, err
	}
	logger.Log.Debugf("Updated %s: %v\n", tenmod.TenantUserStr, models.AsJSONString(tenantUserRequest))
	return dataContainer, nil
}

// DeleteTenantUser - CouchDB implementation of DeleteTenantUser
func (tsd *TenantServiceDatastoreCouchDB) DeleteTenantUser(tenantID string, dataID string) (*tenmod.User, error) {
	logger.Log.Debugf("Deleting %s %s\n", tenmod.TenantUserStr, dataID)
	dataID = ds.PrependToDataID(dataID, string(tenmod.TenantUserType))
	tenantID = ds.PrependToDataID(tenantID, string(admmod.TenantType))

	tenantDBName := createDBPathStr(tsd.server, tenantID)
	dataContainer := tenmod.User{}
	if err := deleteDataFromCouch(tenantDBName, dataID, &dataContainer, tenmod.TenantUserStr); err != nil {
		return nil, err
	}
	logger.Log.Debugf("Deleted %s: %v\n", tenmod.TenantUserStr, models.AsJSONString(dataContainer))
	return &dataContainer, nil
}

// GetTenantUser - CouchDB implementation of GetTenantUser
func (tsd *TenantServiceDatastoreCouchDB) GetTenantUser(tenantID string, dataID string) (*tenmod.User, error) {
	logger.Log.Debugf("Fetching %s: %s\n", tenmod.TenantUserStr, dataID)
	dataID = ds.PrependToDataID(dataID, string(tenmod.TenantUserType))
	tenantID = ds.PrependToDataID(tenantID, string(admmod.TenantType))

	tenantDBName := createDBPathStr(tsd.server, tenantID)
	dataContainer := tenmod.User{}
	if err := getDataFromCouch(tenantDBName, dataID, &dataContainer, tenmod.TenantUserStr); err != nil {
		return nil, err
	}
	logger.Log.Debugf("Retrieved %s: %v\n", tenmod.TenantUserStr, models.AsJSONString(dataContainer))
	return &dataContainer, nil
}

// GetAllTenantUsers - CouchDB implementation of GetAllTenantUsers
func (tsd *TenantServiceDatastoreCouchDB) GetAllTenantUsers(tenantID string) ([]*tenmod.User, error) {
	logger.Log.Debugf("Fetching all %s\n", tenmod.TenantUserStr)
	tenantID = ds.PrependToDataID(tenantID, string(admmod.TenantType))

	tenantDBName := createDBPathStr(tsd.server, tenantID)
	res := make([]*tenmod.User, 0)
	if err := getAllOfTypeFromCouchAndFlatten(tenantDBName, string(tenmod.TenantUserType), tenmod.TenantUserStr, &res); err != nil {
		return nil, err
	}

	logger.Log.Debugf("Retrieved %d %s\n", len(res), tenmod.TenantUserStr)
	return res, nil
}

// CreateTenantDomain - CouchDB implementation of CreateTenantDomain
func (tsd *TenantServiceDatastoreCouchDB) CreateTenantDomain(tenantDomainRequest *tenmod.Domain) (*tenmod.Domain, error) {
	logger.Log.Debugf("Creating %s: %v\n", tenmod.TenantDomainStr, models.AsJSONString(tenantDomainRequest))
	tenantDomainRequest.ID = ds.GenerateID(tenantDomainRequest, string(tenmod.TenantDomainType))
	tenantID := ds.PrependToDataID(tenantDomainRequest.TenantID, string(admmod.TenantType))

	tenantDBName := createDBPathStr(tsd.server, tenantID)
	dataContainer := &tenmod.Domain{}
	if err := createDataInCouch(tenantDBName, tenantDomainRequest, dataContainer, string(tenmod.TenantDomainType), tenmod.TenantDomainStr); err != nil {
		return nil, err
	}
	logger.Log.Debugf("Created %s: %v\n", tenmod.TenantDomainStr, models.AsJSONString(dataContainer))
	return dataContainer, nil
}

// UpdateTenantDomain - CouchDB implementation of UpdateTenantDomain
func (tsd *TenantServiceDatastoreCouchDB) UpdateTenantDomain(tenantDomainRequest *tenmod.Domain) (*tenmod.Domain, error) {
	logger.Log.Debugf("Updating %s: %v\n", tenmod.TenantDomainStr, models.AsJSONString(tenantDomainRequest))
	tenantDomainRequest.ID = ds.PrependToDataID(tenantDomainRequest.ID, string(tenmod.TenantDomainType))
	tenantID := ds.PrependToDataID(tenantDomainRequest.TenantID, string(admmod.TenantType))

	tenantDBName := createDBPathStr(tsd.server, tenantID)
	dataContainer := &tenmod.Domain{}
	if err := updateDataInCouch(tenantDBName, tenantDomainRequest, dataContainer, string(tenmod.TenantDomainType), tenmod.TenantDomainStr); err != nil {
		return nil, err
	}
	logger.Log.Debugf("Updated %s: %v\n", tenmod.TenantDomainStr, models.AsJSONString(dataContainer))
	return dataContainer, nil
}

// DeleteTenantDomain - CouchDB implementation of DeleteTenantDomain
func (tsd *TenantServiceDatastoreCouchDB) DeleteTenantDomain(tenantID string, dataID string) (*tenmod.Domain, error) {
	logger.Log.Debugf("Deleting %s: %s\n", tenmod.TenantDomainStr, dataID)
	dataID = ds.PrependToDataID(dataID, string(tenmod.TenantDomainType))
	tenantID = ds.PrependToDataID(tenantID, string(admmod.TenantType))

	tenantDBName := createDBPathStr(tsd.server, tenantID)
	dataContainer := tenmod.Domain{}
	if err := deleteDataFromCouch(tenantDBName, dataID, &dataContainer, tenmod.TenantDomainStr); err != nil {
		return nil, err
	}
	logger.Log.Debugf("Deleted %s: %v\n", tenmod.TenantDomainStr, models.AsJSONString(dataContainer))
	return &dataContainer, nil
}

// GetTenantDomain - CouchDB implementation of GetTenantDomain
func (tsd *TenantServiceDatastoreCouchDB) GetTenantDomain(tenantID string, dataID string) (*tenmod.Domain, error) {
	logger.Log.Debugf("Fetching %s: %s\n", tenmod.TenantDomainStr, dataID)
	dataID = ds.PrependToDataID(dataID, string(tenmod.TenantDomainType))
	tenantID = ds.PrependToDataID(tenantID, string(admmod.TenantType))

	tenantDBName := createDBPathStr(tsd.server, tenantID)
	dataContainer := tenmod.Domain{}
	if err := getDataFromCouch(tenantDBName, dataID, &dataContainer, tenmod.TenantDomainStr); err != nil {
		return nil, err
	}
	logger.Log.Debugf("Retrieved %s: %v\n", tenmod.TenantDomainStr, models.AsJSONString(dataContainer))
	return &dataContainer, nil
}

// GetAllTenantDomains - CouchDB implementation of GetAllTenantDomains
func (tsd *TenantServiceDatastoreCouchDB) GetAllTenantDomains(tenantID string) ([]*tenmod.Domain, error) {
	logger.Log.Debugf("Fetching all %s\n", tenmod.TenantDomainStr)
	tenantID = ds.PrependToDataID(tenantID, string(admmod.TenantType))

	tenantDBName := createDBPathStr(tsd.server, tenantID)
	res := make([]*tenmod.Domain, 0)
	if err := getAllOfTypeFromCouchAndFlatten(tenantDBName, string(tenmod.TenantDomainType), tenmod.TenantDomainStr, &res); err != nil {
		return nil, err
	}

	logger.Log.Debugf("Retrieved %d %s\n", len(res), tenmod.TenantDomainStr)
	return res, nil
}

// CreateTenantIngestionProfile - CouchDB implementation of CreateTenantIngestionProfile
func (tsd *TenantServiceDatastoreCouchDB) CreateTenantIngestionProfile(tenantIngPrfReq *tenmod.IngestionProfile) (*tenmod.IngestionProfile, error) {
	logger.Log.Debugf("Creating %s: %v\n", tenmod.TenantIngestionProfileStr, models.AsJSONString(tenantIngPrfReq))
	tenantIngPrfReq.ID = ds.GenerateID(tenantIngPrfReq, string(tenmod.TenantIngestionProfileType))
	tenantID := ds.PrependToDataID(tenantIngPrfReq.TenantID, string(admmod.TenantType))

	// Only create one if one does not already exist:
	existing, _ := tsd.GetActiveTenantIngestionProfile(tenantIngPrfReq.TenantID)
	if existing != nil {
		return nil, fmt.Errorf("Can't create %s, it already exists", tenmod.TenantIngestionProfileStr)
	}

	tenantDBName := createDBPathStr(tsd.server, tenantID)
	dataContainer := &tenmod.IngestionProfile{}
	if err := createDataInCouch(tenantDBName, tenantIngPrfReq, dataContainer, string(tenmod.TenantIngestionProfileType), tenmod.TenantIngestionProfileStr); err != nil {
		return nil, err
	}
	logger.Log.Debugf("Created %s: %v\n", tenmod.TenantIngestionProfileStr, models.AsJSONString(dataContainer))
	return dataContainer, nil
}

// UpdateTenantIngestionProfile - CouchDB implementation of UpdateTenantIngestionProfile
func (tsd *TenantServiceDatastoreCouchDB) UpdateTenantIngestionProfile(tenantIngPrfReq *tenmod.IngestionProfile) (*tenmod.IngestionProfile, error) {
	logger.Log.Debugf("Updating %s: %v\n", tenmod.TenantIngestionProfileStr, models.AsJSONString(tenantIngPrfReq))
	tenantIngPrfReq.ID = ds.PrependToDataID(tenantIngPrfReq.ID, string(tenmod.TenantIngestionProfileType))
	tenantID := ds.PrependToDataID(tenantIngPrfReq.TenantID, string(admmod.TenantType))

	tenantDBName := createDBPathStr(tsd.server, tenantID)
	dataContainer := &tenmod.IngestionProfile{}
	if err := updateDataInCouch(tenantDBName, tenantIngPrfReq, dataContainer, string(tenmod.TenantIngestionProfileType), tenmod.TenantIngestionProfileStr); err != nil {
		return nil, err
	}
	logger.Log.Debugf("Updated %s: %v\n", tenmod.TenantIngestionProfileStr, models.AsJSONString(dataContainer))
	return dataContainer, nil
}

// GetTenantIngestionProfile - CouchDB implementation of GetTenantIngestionProfile
func (tsd *TenantServiceDatastoreCouchDB) GetTenantIngestionProfile(tenantID string, dataID string) (*tenmod.IngestionProfile, error) {
	logger.Log.Debugf("Fetching %s: %s\n", tenmod.TenantIngestionProfileStr, dataID)
	dataID = ds.PrependToDataID(dataID, string(tenmod.TenantIngestionProfileType))
	tenantID = ds.PrependToDataID(tenantID, string(admmod.TenantType))

	tenantDBName := createDBPathStr(tsd.server, tenantID)
	dataContainer := tenmod.IngestionProfile{}
	if err := getDataFromCouch(tenantDBName, dataID, &dataContainer, tenmod.TenantIngestionProfileStr); err != nil {
		return nil, err
	}
	logger.Log.Debugf("Retrieved %s: %v\n", tenmod.TenantIngestionProfileStr, models.AsJSONString(dataContainer))
	return &dataContainer, nil
}

// DeleteTenantIngestionProfile - CouchDB implementation of DeleteTenantIngestionProfile
func (tsd *TenantServiceDatastoreCouchDB) DeleteTenantIngestionProfile(tenantID string, dataID string) (*tenmod.IngestionProfile, error) {
	logger.Log.Debugf("Deleting %s: %s\n", tenmod.TenantIngestionProfileStr, dataID)
	dataID = ds.PrependToDataID(dataID, string(tenmod.TenantIngestionProfileType))
	tenantID = ds.PrependToDataID(tenantID, string(admmod.TenantType))

	tenantDBName := createDBPathStr(tsd.server, tenantID)
	dataContainer := tenmod.IngestionProfile{}
	if err := deleteDataFromCouch(tenantDBName, dataID, &dataContainer, tenmod.TenantIngestionProfileStr); err != nil {
		return nil, err
	}
	logger.Log.Debugf("Deleted %s: %v\n", tenmod.TenantIngestionProfileStr, models.AsJSONString(dataContainer))
	return &dataContainer, nil
}

// CreateTenantThresholdProfile - CouchDB implementation of CreateTenantThresholdProfile
func (tsd *TenantServiceDatastoreCouchDB) CreateTenantThresholdProfile(tenantThreshPrfReq *tenmod.ThresholdProfile) (*tenmod.ThresholdProfile, error) {
	logger.Log.Debugf("Creating %s: %v\n", tenmod.TenantThresholdProfileStr, models.AsJSONString(tenantThreshPrfReq))
	tenantThreshPrfReq.ID = ds.GenerateID(tenantThreshPrfReq, string(tenmod.TenantThresholdProfileType))
	tenantID := ds.PrependToDataID(tenantThreshPrfReq.TenantID, string(admmod.TenantType))

	tenantDBName := createDBPathStr(tsd.server, tenantID)
	dataContainer := &tenmod.ThresholdProfile{}
	if err := createDataInCouch(tenantDBName, tenantThreshPrfReq, dataContainer, string(tenmod.TenantThresholdProfileType), tenmod.TenantThresholdProfileStr); err != nil {
		return nil, err
	}
	logger.Log.Debugf("Created %s: %v\n", tenmod.TenantThresholdProfileStr, models.AsJSONString(dataContainer))
	return dataContainer, nil
}

// UpdateTenantThresholdProfile - CouchDB implementation of UpdateTenantThresholdProfile
func (tsd *TenantServiceDatastoreCouchDB) UpdateTenantThresholdProfile(tenantThreshPrfReq *tenmod.ThresholdProfile) (*tenmod.ThresholdProfile, error) {
	logger.Log.Debugf("Updating %s: %v\n", tenmod.TenantThresholdProfileStr, models.AsJSONString(tenantThreshPrfReq))
	tenantThreshPrfReq.ID = ds.PrependToDataID(tenantThreshPrfReq.ID, string(tenmod.TenantThresholdProfileType))
	tenantID := ds.PrependToDataID(tenantThreshPrfReq.TenantID, string(admmod.TenantType))

	tenantDBName := createDBPathStr(tsd.server, tenantID)
	dataContainer := &tenmod.ThresholdProfile{}
	if err := updateDataInCouch(tenantDBName, tenantThreshPrfReq, dataContainer, string(tenmod.TenantThresholdProfileType), tenmod.TenantThresholdProfileStr); err != nil {
		return nil, err
	}
	logger.Log.Debugf("Updated %s: %v\n", tenmod.TenantThresholdProfileStr, models.AsJSONString(dataContainer))
	return dataContainer, nil
}

// GetTenantThresholdProfile - CouchDB implementation of GetTenantThresholdProfile
func (tsd *TenantServiceDatastoreCouchDB) GetTenantThresholdProfile(tenantID string, dataID string) (*tenmod.ThresholdProfile, error) {
	logger.Log.Debugf("Fetching %s: %s\n", tenmod.TenantThresholdProfileStr, dataID)
	dataID = ds.PrependToDataID(dataID, string(tenmod.TenantThresholdProfileType))
	tenantID = ds.PrependToDataID(tenantID, string(admmod.TenantType))

	tenantDBName := createDBPathStr(tsd.server, tenantID)
	dataContainer := tenmod.ThresholdProfile{}
	if err := getDataFromCouch(tenantDBName, dataID, &dataContainer, tenmod.TenantThresholdProfileStr); err != nil {
		return nil, err
	}
	logger.Log.Debugf("Retrieved %s: %v\n", tenmod.TenantThresholdProfileStr, models.AsJSONString(dataContainer))
	return &dataContainer, nil
}

// DeleteTenantThresholdProfile - CouchDB implementation of DeleteTenantThresholdProfile
func (tsd *TenantServiceDatastoreCouchDB) DeleteTenantThresholdProfile(tenantID string, dataID string) (*tenmod.ThresholdProfile, error) {
	logger.Log.Debugf("Deleting %s: %s\n", tenmod.TenantThresholdProfileStr, dataID)
	dataID = ds.PrependToDataID(dataID, string(tenmod.TenantThresholdProfileType))
	tenantID = ds.PrependToDataID(tenantID, string(admmod.TenantType))

	tenantDBName := createDBPathStr(tsd.server, tenantID)
	dataContainer := tenmod.ThresholdProfile{}
	if err := deleteDataFromCouch(tenantDBName, dataID, &dataContainer, tenmod.TenantThresholdProfileStr); err != nil {
		return nil, err
	}
	logger.Log.Debugf("Deleted %s: %v\n", tenmod.TenantThresholdProfileStr, models.AsJSONString(dataContainer))
	return &dataContainer, nil
}

// CreateMonitoredObject - CouchDB implementation of CreateMonitoredObject
func (tsd *TenantServiceDatastoreCouchDB) CreateMonitoredObject(monitoredObjectReq *tenmod.MonitoredObject) (*tenmod.MonitoredObject, error) {
	logger.Log.Debugf("Creating %s: %v\n", tenmod.TenantMonitoredObjectStr, models.AsJSONString(monitoredObjectReq))
	monitoredObjectReq.ID = ds.GenerateID(monitoredObjectReq, string(tenmod.TenantMonitoredObjectType))
	tenantID := ds.PrependToDataID(monitoredObjectReq.TenantID, string(admmod.TenantType))

	tenantDBName := createDBPathStr(tsd.server, tenantID)
	dataContainer := &tenmod.MonitoredObject{}
	if err := createDataInCouch(tenantDBName, monitoredObjectReq, dataContainer, string(tenmod.TenantMonitoredObjectType), tenmod.TenantMonitoredObjectStr); err != nil {
		return nil, err
	}
	logger.Log.Debugf("Created %s: %v\n", tenmod.TenantMonitoredObjectStr, models.AsJSONString(dataContainer))
	return dataContainer, nil
}

// UpdateMonitoredObject - CouchDB implementation of UpdateMonitoredObject
func (tsd *TenantServiceDatastoreCouchDB) UpdateMonitoredObject(monitoredObjectReq *tenmod.MonitoredObject) (*tenmod.MonitoredObject, error) {
	logger.Log.Debugf("Updating %s: %v\n", tenmod.TenantMonitoredObjectStr, models.AsJSONString(monitoredObjectReq))
	monitoredObjectReq.ID = ds.PrependToDataID(monitoredObjectReq.ID, string(tenmod.TenantMonitoredObjectType))
	tenantID := ds.PrependToDataID(monitoredObjectReq.TenantID, string(admmod.TenantType))

	tenantDBName := createDBPathStr(tsd.server, tenantID)
	dataContainer := &tenmod.MonitoredObject{}
	if err := updateDataInCouch(tenantDBName, monitoredObjectReq, dataContainer, string(tenmod.TenantMonitoredObjectType), tenmod.TenantMonitoredObjectStr); err != nil {
		return nil, err
	}
	logger.Log.Debugf("Updated %s: %v\n", tenmod.TenantMonitoredObjectStr, models.AsJSONString(dataContainer))
	return dataContainer, nil
}

// GetMonitoredObject - CouchDB implementation of GetMonitoredObject
func (tsd *TenantServiceDatastoreCouchDB) GetMonitoredObject(tenantID string, dataID string) (*tenmod.MonitoredObject, error) {
	logger.Log.Debugf("Fetching %s: %s\n", tenmod.TenantMonitoredObjectStr, dataID)
	dataID = ds.PrependToDataID(dataID, string(tenmod.TenantMonitoredObjectType))
	tenantID = ds.PrependToDataID(tenantID, string(admmod.TenantType))

	tenantDBName := createDBPathStr(tsd.server, tenantID)
	dataContainer := tenmod.MonitoredObject{}
	if err := getDataFromCouch(tenantDBName, dataID, &dataContainer, tenmod.TenantMonitoredObjectStr); err != nil {
		return nil, err
	}
	logger.Log.Debugf("Retrieved %s: %v\n", tenmod.TenantMonitoredObjectStr, models.AsJSONString(dataContainer))
	return &dataContainer, nil
}

// DeleteMonitoredObject - CouchDB implementation of DeleteMonitoredObject
func (tsd *TenantServiceDatastoreCouchDB) DeleteMonitoredObject(tenantID string, dataID string) (*tenmod.MonitoredObject, error) {
	logger.Log.Debugf("Deleting %s: %s\n", tenmod.TenantMonitoredObjectStr, dataID)
	dataID = ds.PrependToDataID(dataID, string(tenmod.TenantMonitoredObjectType))
	tenantID = ds.PrependToDataID(tenantID, string(admmod.TenantType))

	tenantDBName := createDBPathStr(tsd.server, tenantID)
	dataContainer := tenmod.MonitoredObject{}
	if err := deleteDataFromCouch(tenantDBName, dataID, &dataContainer, tenmod.TenantMonitoredObjectStr); err != nil {
		return nil, err
	}
	logger.Log.Debugf("Deleted %s: %v\n", tenmod.TenantMonitoredObjectStr, models.AsJSONString(dataContainer))
	return &dataContainer, nil
}

// GetAllMonitoredObjects - CouchDB implementation of GetAllMonitoredObjects
func (tsd *TenantServiceDatastoreCouchDB) GetAllMonitoredObjects(tenantID string) ([]*tenmod.MonitoredObject, error) {
	logger.Log.Debugf("Fetching all %s\n", tenmod.TenantMonitoredObjectStr)
	tenantID = ds.PrependToDataID(tenantID, string(admmod.TenantType))

	tenantDBName := createDBPathStr(tsd.server, tenantID)
	res := make([]*tenmod.MonitoredObject, 0)
	if err := getAllOfTypeFromCouchAndFlatten(tenantDBName, string(tenmod.TenantMonitoredObjectType), tenmod.TenantMonitoredObjectStr, &res); err != nil {
		return nil, err
	}

	logger.Log.Debugf("Retrieved %d %s\n", len(res), tenmod.TenantMonitoredObjectStr)
	return res, nil
}

// GetMonitoredObjectToDomainMap - CouchDB implementation of GetMonitoredObjectToDomainMap
func (tsd *TenantServiceDatastoreCouchDB) GetMonitoredObjectToDomainMap(moByDomReq *tenmod.MonitoredObjectCountByDomainRequest) (*tenmod.MonitoredObjectCountByDomainResponse, error) {
	logger.Log.Debugf("Fetching %s: %v\n", tenmod.MonitoredObjectToDomainMapStr, models.AsJSONString(moByDomReq))
	tenantID := ds.PrependToDataID(moByDomReq.TenantID, string(admmod.TenantType))

	tenantDBName := createDBPathStr(tsd.server, tenantID)
	db, err := getDatabase(tenantDBName)
	if err != nil {
		return nil, err
	}

	// Get response data either by subset, or for all domains
	domainSet := moByDomReq.DomainSet
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
		requestBody["keys"] = moByDomReq.DomainSet

		fetchResponse, err = fetchDesignDocumentResults(requestBody, tenantDBName, monitoredObjectsByDomainIndex)
		if err != nil {
			return nil, err
		}

	}

	if fetchResponse["rows"] == nil {
		return &tenmod.MonitoredObjectCountByDomainResponse{}, nil
	}

	// Response will vary depending on if it is aggregated or not
	response := tenmod.MonitoredObjectCountByDomainResponse{}
	rows := fetchResponse["rows"].([]interface{})
	if moByDomReq.ByCount {
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
		domainMap := map[string][]string{}
		for _, row := range rows {
			obj := row.(map[string]interface{})
			key := obj["key"].(string)
			val := ds.GetDataIDFromFullID(obj["value"].(string))
			if domainMap[key] == nil {
				domainMap[key] = []string{}
			}
			domainMap[key] = append(domainMap[key], val)
		}
		response.DomainToMonitoredObjectSetMap = domainMap
	}

	logger.Log.Debugf("Retrieved %s: %v\n", tenmod.MonitoredObjectToDomainMapStr, models.AsJSONString(response))
	return &response, nil
}

// CreateTenantMeta - CouchDB implementation of CreateTenantMeta
func (tsd *TenantServiceDatastoreCouchDB) CreateTenantMeta(meta *tenmod.Metadata) (*tenmod.Metadata, error) {
	logger.Log.Debugf("Creating %s: %v\n", tenmod.TenantMetaStr, models.AsJSONString(meta))
	meta.ID = ds.GenerateID(meta, string(tenmod.TenantMetaType))
	tenantID := ds.PrependToDataID(meta.TenantID, string(admmod.TenantType))

	// Only create one if one does not already exist:
	existing, _ := tsd.GetTenantMeta(meta.TenantID)
	if existing != nil {
		return nil, fmt.Errorf("Can't create %s, it already exists", tenmod.TenantMetaStr)
	}

	tenantDBName := createDBPathStr(tsd.server, tenantID)
	dataContainer := &tenmod.Metadata{}
	if err := createDataInCouch(tenantDBName, meta, dataContainer, string(tenmod.TenantMetaType), tenmod.TenantMetaStr); err != nil {
		return nil, err
	}
	logger.Log.Debugf("Created %s: %v\n", tenmod.TenantMetaStr, models.AsJSONString(dataContainer))
	return dataContainer, nil
}

// UpdateTenantMeta - CouchDB implementation of UpdateTenantMeta
func (tsd *TenantServiceDatastoreCouchDB) UpdateTenantMeta(meta *tenmod.Metadata) (*tenmod.Metadata, error) {
	logger.Log.Debugf("Updating %s: %v\n", tenmod.TenantMetaStr, models.AsJSONString(meta))
	meta.ID = ds.PrependToDataID(meta.ID, string(tenmod.TenantMetaType))
	tenantID := ds.PrependToDataID(meta.TenantID, string(admmod.TenantType))

	tenantDBName := createDBPathStr(tsd.server, tenantID)
	dataContainer := &tenmod.Metadata{}
	if err := updateDataInCouch(tenantDBName, meta, dataContainer, string(tenmod.TenantMetaType), tenmod.TenantMetaStr); err != nil {
		return nil, err
	}
	logger.Log.Debugf("Updated %s: %v\n", tenmod.TenantMetaStr, models.AsJSONString(dataContainer))
	return dataContainer, nil
}

// DeleteTenantMeta - CouchDB implementation of DeleteTenantMeta
func (tsd *TenantServiceDatastoreCouchDB) DeleteTenantMeta(tenantID string) (*tenmod.Metadata, error) {
	logger.Log.Debugf("Deleting %s: %s\n", tenmod.TenantMetaStr, tenantID)

	// Obtain the value of the existing record for a return value.
	existingObject, err := tsd.GetTenantMeta(tenantID)
	if err != nil {
		logger.Log.Debugf("Unable to fetch %s to delete: %s", tenmod.TenantMetaStr, err.Error())
		return nil, err
	}

	tenantID = ds.PrependToDataID(tenantID, string(admmod.TenantType))
	tenantDBName := createDBPathStr(tsd.server, tenantID)
	objectID := ds.PrependToDataID(existingObject.ID, string(tenmod.TenantMetaType))
	if err := deleteData(tenantDBName, objectID, tenmod.TenantMetaStr); err != nil {
		logger.Log.Debugf("Unable to delete %s: %s", tenmod.TenantMetaStr, err.Error())
		return nil, err
	}

	// Return the deleted object.
	logger.Log.Debugf("Deleted %s: %v\n", tenmod.TenantMetaStr, models.AsJSONString(existingObject))
	return existingObject, nil
}

// GetTenantMeta - CouchDB implementation of GetTenantMeta
func (tsd *TenantServiceDatastoreCouchDB) GetTenantMeta(tenantID string) (*tenmod.Metadata, error) {
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
	res := tenmod.Metadata{}
	if len(fetchedData) != 0 {
		if err = convertGenericCouchDataToObject(fetchedData[0], &res, tenmod.TenantMetaStr); err != nil {
			return nil, err
		}
	} else {
		return nil, fmt.Errorf("Unable to find %s", tenmod.TenantMetaStr)
	}

	logger.Log.Debugf("Retrieved %s: %v\n", tenmod.TenantMetaStr, models.AsJSONString(res))
	return &res, nil
}

// GetActiveTenantIngestionProfile - CouchDB implementation of GetActiveTenantIngestionProfile
func (tsd *TenantServiceDatastoreCouchDB) GetActiveTenantIngestionProfile(tenantID string) (*tenmod.IngestionProfile, error) {
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

	logger.Log.Debugf("Retrieved: %s", models.AsJSONString(fetchedData))

	// Populate the response
	res := tenmod.IngestionProfile{}
	if len(fetchedData) != 0 {
		if err = convertGenericCouchDataToObject(fetchedData[0], &res, tenmod.TenantIngestionProfileStr); err != nil {
			return nil, err
		}
	} else {
		return nil, fmt.Errorf("%s not found", tenmod.TenantIngestionProfileStr)
	}

	logger.Log.Debugf("Retrieved %s: %v\n", tenmod.TenantIngestionProfileStr, models.AsJSONString(res))
	return &res, nil
}

// GetAllTenantThresholdProfile - CouchDB implementation of GetAllTenantThresholdProfile
func (tsd *TenantServiceDatastoreCouchDB) GetAllTenantThresholdProfile(tenantID string) ([]*tenmod.ThresholdProfile, error) {
	logger.Log.Debugf("Fetching all %s for Tenant %s\n", tenmod.TenantThresholdProfileStr, tenantID)
	tenantID = ds.PrependToDataID(tenantID, string(admmod.TenantType))

	tenantDBName := createDBPathStr(tsd.server, tenantID)
	res := make([]*tenmod.ThresholdProfile, 0)
	if err := getAllOfTypeFromCouch(tenantDBName, string(tenmod.TenantThresholdProfileType), tenmod.TenantThresholdProfileStr, &res); err != nil {
		return nil, err
	}

	logger.Log.Debugf("Retrieved %d %s\n", tenmod.TenantThresholdProfileStr, models.AsJSONString(res))
	return res, nil
}

// BulkInsertMonitoredObjects - CouchDB implementation of BulkInsertMonitoredObjects
func (tsd *TenantServiceDatastoreCouchDB) BulkInsertMonitoredObjects(tenantID string, value []*tenmod.MonitoredObject) ([]*common.BulkOperationResult, error) {
	logger.Log.Debugf("Bulk creating %s: %v\n", tenmod.TenantMonitoredObjectStr, models.AsJSONString(value))
	tenantID = ds.PrependToDataID(tenantID, string(admmod.TenantType))
	tenantDBName := createDBPathStr(tsd.server, tenantID)
	resource, err := couchdb.NewResource(tenantDBName, nil)
	if err != nil {
		return nil, err
	}

	// Iterate over the collection and populate necessary fields
	data := make([]map[string]interface{}, 0)
	for _, mo := range value {
		genericMO, err := convertDataToCouchDbSupportedModel(mo)
		if err != nil {
			return nil, err
		}

		dataType := string(tenmod.TenantMonitoredObjectType)
		genericMO["_id"] = ds.GenerateID(mo, dataType)

		dataProp := genericMO["data"].(map[string]interface{})
		dataProp["datatype"] = dataType
		dataProp["createdTimestamp"] = ds.MakeTimestamp()
		dataProp["lastModifiedTimestamp"] = genericMO["createdTimestamp"]

		data = append(data, genericMO)
	}
	body := map[string]interface{}{
		"docs": data}

	fetchedData, err := performBulkUpdate(body, resource)
	if err != nil {
		return nil, err
	}

	// Populate the response
	res := make([]*common.BulkOperationResult, 0)
	for _, fetched := range fetchedData {
		newObj := common.BulkOperationResult{}
		if fetched["id"] != nil {
			newObj.ID = fetched["id"].(string)
		}
		if fetched["rev"] != nil {
			newObj.REV = fetched["rev"].(string)
		}
		if fetched["reason"] != nil {
			newObj.REASON = fetched["reason"].(string)
		}
		if fetched["error"] != nil {
			newObj.ERROR = fetched["error"].(string)
		}
		if fetched["ok"] != nil {
			newObj.OK = fetched["ok"].(bool)
		}

		newObj.ID = ds.GetDataIDFromFullID(newObj.ID)
		res = append(res, &newObj)
	}

	logger.Log.Debugf("Bulk create of %s result: %v\n", tenmod.TenantMonitoredObjectStr, models.AsJSONString(res))
	return res, nil
}
