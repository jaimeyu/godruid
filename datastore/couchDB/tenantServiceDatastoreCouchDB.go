package couchDB

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"time"

	"github.com/accedian/adh-gather/config"
	ds "github.com/accedian/adh-gather/datastore"
	"github.com/accedian/adh-gather/datastore/druid"

	"github.com/accedian/adh-gather/gather"
	"github.com/accedian/adh-gather/logger"
	"github.com/accedian/adh-gather/models"

	admmod "github.com/accedian/adh-gather/models/admin"
	"github.com/accedian/adh-gather/models/common"
	metmod "github.com/accedian/adh-gather/models/metrics"
	tenmod "github.com/accedian/adh-gather/models/tenant"
	mon "github.com/accedian/adh-gather/monitoring"
	couchdb "github.com/leesper/couchdb-golang"
)

const (
	monitoredObjectsByDomainIndex = "_design/monitoredObjectCount/_view/byDomain"
	monitoredObjectsByNameIndex   = "_design/moIndex/_view/byName"
)

// TenantServiceDatastoreCouchDB - struct responsible for handling
// database operations for the Tenant Service when using CouchDB
// as the storage option.
type TenantServiceDatastoreCouchDB struct {
	server              string
	cfg                 config.Provider
	connectorUpdateChan chan *tenmod.ConnectorConfig
	metricsDB           ds.DruidDatastore
	batchSize           int64
}

// CreateTenantServiceDAO - instantiates a CouchDB implementation of the
// TenantServiceDatastore.
func CreateTenantServiceDAO() (*TenantServiceDatastoreCouchDB, error) {
	result := new(TenantServiceDatastoreCouchDB)

	result.metricsDB = druid.NewDruidDatasctoreClient()
	result.cfg = gather.GetConfig()
	result.connectorUpdateChan = make(chan *tenmod.ConnectorConfig)

	provDBURL := fmt.Sprintf("%s:%d",
		result.cfg.GetString(gather.CK_server_datastore_ip.String()),
		result.cfg.GetInt(gather.CK_server_datastore_port.String()))
	logger.Log.Debugf("Tenant Service CouchDB URL is: %s, %v", provDBURL, result.connectorUpdateChan)
	result.server = provDBURL

	result.batchSize = int64(result.cfg.GetInt(gather.CK_server_datastore_batchsize.String()))

	return result, nil
}

//GetConnectorConfigUpdateChan - Get the Couchdb connector channel
func (tsd *TenantServiceDatastoreCouchDB) GetConnectorConfigUpdateChan() chan *tenmod.ConnectorConfig {
	return tsd.connectorUpdateChan
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

// CreateTenantConnectorInstance - CouchDB implementation of CreateTenantConnectorInstance
func (tsd *TenantServiceDatastoreCouchDB) CreateTenantConnectorInstance(tenantConnectorInstanceRequest *tenmod.ConnectorInstance) (*tenmod.ConnectorInstance, error) {
	logger.Log.Debugf("Creating %s: %v\n", tenmod.TenantConnectorInstanceStr, models.AsJSONString(tenantConnectorInstanceRequest))
	id := tenantConnectorInstanceRequest.ID
	if id == "" {
		tenantConnectorInstanceRequest.ID = ds.GenerateID(tenantConnectorInstanceRequest, string(tenmod.TenantConnectorInstanceType))
	} else {
		tenantConnectorInstanceRequest.ID = ds.PrependToDataID(id, string(tenmod.TenantConnectorInstanceType))
	}

	tenantID := ds.PrependToDataID(tenantConnectorInstanceRequest.TenantID, string(admmod.TenantType))

	tenantDBName := createDBPathStr(tsd.server, tenantID)
	dataContainer := &tenmod.ConnectorInstance{}
	if err := createDataInCouch(tenantDBName, tenantConnectorInstanceRequest, dataContainer, string(tenmod.TenantConnectorInstanceType), tenmod.TenantConnectorInstanceStr); err != nil {
		return nil, err
	}
	logger.Log.Debugf("Created %s: %v\n", tenmod.TenantConnectorInstanceStr, models.AsJSONString(dataContainer))
	return dataContainer, nil
}

// UpdateTenantConnectorInstance - CouchDB implementation of UpdateTenantConnectorInstance
func (tsd *TenantServiceDatastoreCouchDB) UpdateTenantConnectorInstance(tenantConnectorInstanceRequest *tenmod.ConnectorInstance) (*tenmod.ConnectorInstance, error) {
	logger.Log.Debugf("Updating %s: %v\n", tenmod.TenantConnectorInstanceStr, models.AsJSONString(tenantConnectorInstanceRequest))
	tenantConnectorInstanceRequest.ID = ds.PrependToDataID(tenantConnectorInstanceRequest.ID, string(tenmod.TenantConnectorInstanceType))
	tenantID := ds.PrependToDataID(tenantConnectorInstanceRequest.TenantID, string(admmod.TenantType))

	tenantDBName := createDBPathStr(tsd.server, tenantID)
	dataContainer := &tenmod.ConnectorInstance{}
	if err := updateDataInCouch(tenantDBName, tenantConnectorInstanceRequest, dataContainer, string(tenmod.TenantConnectorInstanceType), tenmod.TenantConnectorInstanceStr); err != nil {
		return nil, err
	}
	logger.Log.Debugf("Updated %s: %v\n", tenmod.TenantConnectorInstanceStr, models.AsJSONString(dataContainer))
	return dataContainer, nil
}

// DeleteTenantConnectorInstance - CouchDB implementation of DeleteTenantConnectorInstance
func (tsd *TenantServiceDatastoreCouchDB) DeleteTenantConnectorInstance(tenantID string, dataID string) (*tenmod.ConnectorInstance, error) {
	logger.Log.Debugf("Deleting %s: %s\n", tenmod.TenantConnectorInstanceStr, dataID)
	dataID = ds.PrependToDataID(dataID, string(tenmod.TenantConnectorInstanceType))
	tenantID = ds.PrependToDataID(tenantID, string(admmod.TenantType))

	tenantDBName := createDBPathStr(tsd.server, tenantID)
	dataContainer := tenmod.ConnectorInstance{}
	if err := deleteDataFromCouch(tenantDBName, dataID, &dataContainer, tenmod.TenantConnectorInstanceStr); err != nil {
		return nil, err
	}
	logger.Log.Debugf("Deleted %s: %v\n", tenmod.TenantConnectorInstanceStr, models.AsJSONString(dataContainer))
	return &dataContainer, nil
}

// GetTenantConnectorInstance - CouchDB implementation of GetTenantConnectorInstance
func (tsd *TenantServiceDatastoreCouchDB) GetTenantConnectorInstance(tenantID string, dataID string) (*tenmod.ConnectorInstance, error) {
	logger.Log.Debugf("Fetching %s: %s\n", tenmod.TenantConnectorInstanceStr, dataID)
	dataID = ds.PrependToDataID(dataID, string(tenmod.TenantConnectorInstanceType))
	tenantID = ds.PrependToDataID(tenantID, string(admmod.TenantType))

	tenantDBName := createDBPathStr(tsd.server, tenantID)
	dataContainer := tenmod.ConnectorInstance{}
	if err := getDataFromCouch(tenantDBName, dataID, &dataContainer, tenmod.TenantConnectorInstanceStr); err != nil {
		return nil, err
	}
	logger.Log.Debugf("Retrieved %s: %v\n", tenmod.TenantConnectorInstanceStr, models.AsJSONString(dataContainer))
	return &dataContainer, nil
}

// GetAllTenantConnectorInstances - CouchDB implementation of GetAllTenantConnectorInstances
func (tsd *TenantServiceDatastoreCouchDB) GetAllTenantConnectorInstances(tenantID string) ([]*tenmod.ConnectorInstance, error) {
	logger.Log.Debugf("Fetching all %s\n", tenmod.TenantConnectorInstanceStr)
	tenantID = ds.PrependToDataID(tenantID, string(admmod.TenantType))

	tenantDBName := createDBPathStr(tsd.server, tenantID)
	res := make([]*tenmod.ConnectorInstance, 0)

	if err := getAllOfTypeFromCouchAndFlatten(tenantDBName, string(tenmod.TenantConnectorInstanceType), tenmod.TenantConnectorInstanceStr, &res); err != nil {
		return nil, err
	}

	logger.Log.Debugf("Retrieved %d %s\n", len(res), tenmod.TenantConnectorInstanceStr)
	return res, nil
}

// CreateTenantConnectorConfig - CouchDB implementation of CreateTenantConnectorConfig
func (tsd *TenantServiceDatastoreCouchDB) CreateTenantConnectorConfig(TenantConnectorConfigRequest *tenmod.ConnectorConfig) (*tenmod.ConnectorConfig, error) {
	logger.Log.Debugf("Creating %s: %v\n", tenmod.TenantConnectorConfigStr, models.AsJSONString(TenantConnectorConfigRequest))
	TenantConnectorConfigRequest.ID = ds.GenerateID(TenantConnectorConfigRequest, string(tenmod.TenantConnectorConfigType))
	tenantID := ds.PrependToDataID(TenantConnectorConfigRequest.TenantID, string(admmod.TenantType))

	tenantDBName := createDBPathStr(tsd.server, tenantID)
	dataContainer := &tenmod.ConnectorConfig{}
	if err := createDataInCouch(tenantDBName, TenantConnectorConfigRequest, dataContainer, string(tenmod.TenantConnectorConfigType), tenmod.TenantConnectorConfigStr); err != nil {
		return nil, err
	}
	logger.Log.Debugf("Created %s: %v\n", tenmod.TenantConnectorConfigStr, models.AsJSONString(dataContainer))
	return dataContainer, nil
}

// UpdateTenantConnectorConfig - CouchDB implementation of UpdateTenantConnectorConfig
func (tsd *TenantServiceDatastoreCouchDB) UpdateTenantConnectorConfig(TenantConnectorConfigRequest *tenmod.ConnectorConfig) (*tenmod.ConnectorConfig, error) {
	logger.Log.Debugf("Updating %s: %v\n", tenmod.TenantConnectorConfigStr, models.AsJSONString(TenantConnectorConfigRequest))
	TenantConnectorConfigRequest.ID = ds.PrependToDataID(TenantConnectorConfigRequest.ID, string(tenmod.TenantConnectorConfigType))
	tenantID := ds.PrependToDataID(TenantConnectorConfigRequest.TenantID, string(admmod.TenantType))

	tenantDBName := createDBPathStr(tsd.server, tenantID)
	dataContainer := &tenmod.ConnectorConfig{}
	if err := updateDataInCouch(tenantDBName, TenantConnectorConfigRequest, dataContainer, string(tenmod.TenantConnectorConfigType), tenmod.TenantConnectorConfigStr); err != nil {
		return nil, err
	}
	logger.Log.Debugf("Updated %s: %v\n", tenmod.TenantConnectorConfigStr, models.AsJSONString(dataContainer))

	// put update on channel - don't block
	select {
	case tsd.connectorUpdateChan <- dataContainer:
		logger.Log.Debugf("Sending %s: %v to websocket server\n", tenmod.TenantConnectorConfigStr, models.AsJSONString(dataContainer))
	default:
		logger.Log.Debugf("Websocket server not listening to %s: %v \n", tenmod.TenantConnectorConfigStr, models.AsJSONString(dataContainer))
		break
	}

	return dataContainer, nil
}

// DeleteTenantConnectorConfig - CouchDB implementation of DeleteTenantConnectorConfig
func (tsd *TenantServiceDatastoreCouchDB) DeleteTenantConnectorConfig(tenantID string, dataID string) (*tenmod.ConnectorConfig, error) {
	logger.Log.Debugf("Deleting %s: %s\n", tenmod.TenantConnectorConfigStr, dataID)
	dataID = ds.PrependToDataID(dataID, string(tenmod.TenantConnectorConfigType))
	tenantID = ds.PrependToDataID(tenantID, string(admmod.TenantType))

	tenantDBName := createDBPathStr(tsd.server, tenantID)
	dataContainer := tenmod.ConnectorConfig{}
	if err := deleteDataFromCouch(tenantDBName, dataID, &dataContainer, tenmod.TenantConnectorConfigStr); err != nil {
		return nil, err
	}
	logger.Log.Debugf("Deleted %s: %v\n", tenmod.TenantConnectorConfigStr, models.AsJSONString(dataContainer))
	return &dataContainer, nil
}

// GetTenantConnectorConfig - CouchDB implementation of GetTenantConnectorConfig
func (tsd *TenantServiceDatastoreCouchDB) GetTenantConnectorConfig(tenantID string, dataID string) (*tenmod.ConnectorConfig, error) {
	logger.Log.Debugf("Fetching %s: %s\n", tenmod.TenantConnectorConfigStr, dataID)
	dataID = ds.PrependToDataID(dataID, string(tenmod.TenantConnectorConfigType))
	tenantID = ds.PrependToDataID(tenantID, string(admmod.TenantType))

	tenantDBName := createDBPathStr(tsd.server, tenantID)
	dataContainer := tenmod.ConnectorConfig{}
	if err := getDataFromCouch(tenantDBName, dataID, &dataContainer, tenmod.TenantConnectorConfigStr); err != nil {
		return nil, err
	}
	logger.Log.Debugf("Retrieved %s: %v\n", tenmod.TenantConnectorConfigStr, models.AsJSONString(dataContainer))
	return &dataContainer, nil
}

// GetAllTenantConnectorConfigs - CouchDB implementation of GetAllTenantConnectorConfigs
func (tsd *TenantServiceDatastoreCouchDB) GetAllTenantConnectorConfigs(tenantID, zone string) ([]*tenmod.ConnectorConfig, error) {
	logger.Log.Debugf("Fetching all %s\n", tenmod.TenantConnectorConfigStr)
	tenantID = ds.PrependToDataID(tenantID, string(admmod.TenantType))

	tenantDBName := createDBPathStr(tsd.server, tenantID)
	res := make([]*tenmod.ConnectorConfig, 0)

	if zone == "" {
		if err := getAllOfTypeFromCouchAndFlatten(tenantDBName, string(tenmod.TenantConnectorConfigType), tenmod.TenantConnectorConfigStr, &res); err != nil {
			return nil, err
		}
	} else {
		db, err := getDatabase(tenantDBName)
		if err != nil {
			return nil, err
		}

		fetchedList, err := getAllOfAny(string(tenmod.TenantConnectorConfigType), "zone", zone, string(tenmod.TenantConnectorConfigStr), db)
		if err != nil {
			return nil, err
		}

		if err := convertCouchDataArrayToFlattenedArray(fetchedList, &res, tenmod.TenantConnectorConfigStr); err != nil {
			return nil, err
		}
	}

	logger.Log.Debugf("Retrieved %d %s\n", len(res), tenmod.TenantConnectorConfigStr)
	return res, nil
}

// GetAllAvailableTenantConnectorConfigs - Returns all tenant connectors matching tenantID, zone, that aren't already being used
func (tsd *TenantServiceDatastoreCouchDB) GetAllAvailableTenantConnectorConfigs(tenantID, zone string) ([]*tenmod.ConnectorConfig, error) {
	logger.Log.Debugf("Fetching all available %s\n", tenmod.TenantConnectorConfigStr)
	tenantID = ds.PrependToDataID(tenantID, string(admmod.TenantType))

	tenantDBName := createDBPathStr(tsd.server, tenantID)
	res := make([]*tenmod.ConnectorConfig, 0)

	if zone == "" {
		if err := getAllOfTypeFromCouchAndFlatten(tenantDBName, string(tenmod.TenantConnectorConfigType), tenmod.TenantConnectorConfigStr, &res); err != nil {
			return nil, err
		}
	} else {
		db, err := getDatabase(tenantDBName)
		if err != nil {
			return nil, err
		}

		fetchedList, err := getAvailableConfigs(string(tenmod.TenantConnectorConfigType), "zone", zone, string(tenmod.TenantConnectorConfigStr), db)
		if err != nil {
			return nil, err
		}

		if err := convertCouchDataArrayToFlattenedArray(fetchedList, &res, tenmod.TenantConnectorConfigStr); err != nil {
			return nil, err
		}
	}

	logger.Log.Debugf("Retrieved %d %s\n", len(res), tenmod.TenantConnectorConfigStr)
	return res, nil
}

// GetAllTenantConnectorConfigsByInstanceID - Returns the TenantConnectorConfigConfigs with the given instance ID
func (tsd *TenantServiceDatastoreCouchDB) GetAllTenantConnectorConfigsByInstanceID(tenantID, instanceID string) ([]*tenmod.ConnectorConfig, error) {
	logger.Log.Debugf("Fetching %s with instance ID %s\n", tenmod.TenantConnectorConfigStr, instanceID)
	tenantID = ds.PrependToDataID(tenantID, string(admmod.TenantType))

	tenantDBName := createDBPathStr(tsd.server, tenantID)
	res := make([]*tenmod.ConnectorConfig, 0)

	db, err := getDatabase(tenantDBName)
	if err != nil {
		return nil, err
	}

	fetchedList, err := getAllOfAny(string(tenmod.TenantConnectorConfigType), "connectorInstanceId", instanceID, string(tenmod.TenantConnectorConfigStr), db)
	if err != nil {
		return nil, err
	}

	if err := convertCouchDataArrayToFlattenedArray(fetchedList, &res, tenmod.TenantConnectorConfigStr); err != nil {
		return nil, err
	}
	logger.Log.Debugf("Retrieved %d %s\n", len(res), tenmod.TenantConnectorConfigStr)
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

	// dataContainer := tenmod.ThresholdProfile{}
	// TODO: TEMPORARY FIX until the UI portion of a threshold profile is removed or the UI conforms to the data model.
	dataContainer := map[string]interface{}{}
	if err := getDataFromCouch(tenantDBName, dataID, &dataContainer, tenmod.TenantThresholdProfileStr); err != nil {
		return nil, err
	}
	thresholds := dataContainer["thresholds"].(map[string]interface{})
	vendorMap := thresholds["vendorMap"].(map[string]interface{})
	for _, val := range vendorMap {
		valAsMap := val.(map[string]interface{})
		if valAsMap["metricMap"] != nil {
			delete(valAsMap, "metricMap")
		}
	}

	// Convert the generic object to a thresholdProfile
	finalDataContainer := tenmod.ThresholdProfile{}
	genericDataInBytes, err := convertGenericObjectToBytesWithCouchDbFields(dataContainer)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(genericDataInBytes, &finalDataContainer)
	if err != nil {
		logger.Log.Debugf("Error converting generic data to %s type: %s", tenmod.TenantThresholdProfileStr, err.Error())
		return nil, err
	}

	logger.Log.Debugf("Retrieved %s: %v\n", tenmod.TenantThresholdProfileStr, models.AsJSONString(dataContainer))
	// return &dataContainer, nil
	return &finalDataContainer, nil
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

	dbName := createDBPathStr(tsd.server, fmt.Sprintf("%s%s", tenantID, monitoredObjectDBSuffix))
	dataContainer := &tenmod.MonitoredObject{}

	// Add missing metadata
	err := tsd.CheckAndAddMetadataView(monitoredObjectReq.TenantID, monitoredObjectReq)
	if err != nil {
		return nil, err
	}

	if err := createDataInCouch(dbName, monitoredObjectReq, dataContainer, string(tenmod.TenantMonitoredObjectType), tenmod.TenantMonitoredObjectStr); err != nil {
		return nil, err
	}
	logger.Log.Debugf("Created %s: %v\n", tenmod.TenantMonitoredObjectStr, models.AsJSONString(dataContainer))

	// Update the metadata before updating the monitored object
	err = tsd.UpdateMonitoredObjectMetadataViews(monitoredObjectReq.TenantID, dataContainer)
	if err != nil {
		return nil, err
	}
	return dataContainer, nil
}

// UpdateMonitoredObject - CouchDB implementation of UpdateMonitoredObject
func (tsd *TenantServiceDatastoreCouchDB) UpdateMonitoredObject(monitoredObjectReq *tenmod.MonitoredObject) (*tenmod.MonitoredObject, error) {
	logger.Log.Debugf("Updating %s: %v\n", tenmod.TenantMonitoredObjectStr, models.AsJSONString(monitoredObjectReq))
	monitoredObjectReq.ID = ds.PrependToDataID(monitoredObjectReq.ID, string(tenmod.TenantMonitoredObjectType))
	tenantID := ds.PrependToDataID(monitoredObjectReq.TenantID, string(admmod.TenantType))

	dbName := createDBPathStr(tsd.server, fmt.Sprintf("%s%s", tenantID, monitoredObjectDBSuffix))
	dataContainer := &tenmod.MonitoredObject{}
	if err := updateDataInCouch(dbName, monitoredObjectReq, dataContainer, string(tenmod.TenantMonitoredObjectType), tenmod.TenantMonitoredObjectStr); err != nil {
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

	dbName := createDBPathStr(tsd.server, fmt.Sprintf("%s%s", tenantID, monitoredObjectDBSuffix))
	dataContainer := tenmod.MonitoredObject{}
	if err := getDataFromCouch(dbName, dataID, &dataContainer, tenmod.TenantMonitoredObjectStr); err != nil {
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

	dbName := createDBPathStr(tsd.server, fmt.Sprintf("%s%s", tenantID, monitoredObjectDBSuffix))
	dataContainer := tenmod.MonitoredObject{}
	if err := deleteDataFromCouch(dbName, dataID, &dataContainer, tenmod.TenantMonitoredObjectStr); err != nil {
		return nil, err
	}
	logger.Log.Debugf("Deleted %s: %v\n", tenmod.TenantMonitoredObjectStr, models.AsJSONString(dataContainer))
	return &dataContainer, nil
}

// GetAllMonitoredObjects - CouchDB implementation of GetAllMonitoredObjects
func (tsd *TenantServiceDatastoreCouchDB) GetAllMonitoredObjects(tenantID string) ([]*tenmod.MonitoredObject, error) {
	logger.Log.Debugf("Fetching all %s\n", tenmod.TenantMonitoredObjectStr)
	tenantID = ds.PrependToDataID(tenantID, string(admmod.TenantType))

	dbName := createDBPathStr(tsd.server, fmt.Sprintf("%s%s", tenantID, monitoredObjectDBSuffix))
	res := make([]*tenmod.MonitoredObject, 0)
	if err := getAllOfTypeFromCouchAndFlatten(dbName, string(tenmod.TenantMonitoredObjectType), tenmod.TenantMonitoredObjectStr, &res); err != nil {
		return nil, err
	}

	logger.Log.Debugf("Retrieved %d %s\n", len(res), tenmod.TenantMonitoredObjectStr)
	return res, nil
}

// GetAllMonitoredObjectsByPage - CouchDB implementation of GetAllMonitoredObjectsByPage
func (tsd *TenantServiceDatastoreCouchDB) GetAllMonitoredObjectsByPage(tenantID string, startKey string, limit int64) ([]*tenmod.MonitoredObject, *common.PaginationOffsets, error) {
	logger.Log.Debugf("Fetching next %d %ss from startKey %s\n", limit, tenmod.TenantMonitoredObjectStr, startKey)
	tenantID = ds.PrependToDataID(tenantID, string(admmod.TenantType))

	dbName := createDBPathStr(tsd.server, fmt.Sprintf("%s%s", tenantID, monitoredObjectDBSuffix))
	db, err := getDatabase(dbName)
	if err != nil {
		return nil, nil, err
	}

	res := make([]*tenmod.MonitoredObject, 0)

	// Need to retrieve 1 more than the asking size to be able to give back a startKey for the next page
	var batchSize int64
	if limit <= 0 || limit > int64(tsd.batchSize) {
		batchSize = tsd.batchSize
		logger.Log.Warningf("Provided limit %d is outside of range [1 - %d]. Using value %d in query", limit, batchSize, batchSize)
	} else {
		batchSize = limit
	}

	// Get 1 more object than the real response so that we can have the start key of the next page
	batchPlus1 := batchSize + 1

	params := generatePaginationQueryParams(startKey, batchPlus1, true, false)
	fetchResponse, err := getByDocIDWithQueryParams(monitoredObjectsByNameIndex, tenmod.TenantMonitoredObjectStr, db, &params)
	if err != nil {
		return nil, nil, err
	}

	if fetchResponse["rows"] == nil {
		return nil, nil, fmt.Errorf(ds.NotFoundStr)
	}

	castedRows := fetchResponse["rows"].([]interface{})
	if len(castedRows) == 0 {
		return nil, nil, fmt.Errorf(ds.NotFoundStr)
	}

	// Convert interface results to map results
	rows := []map[string]interface{}{}
	for _, obj := range castedRows {
		castedObj := obj.(map[string]interface{})
		genericDoc := castedObj["doc"].(map[string]interface{})
		rows = append(rows, genericDoc)
	}

	convertCouchDataArrayToFlattenedArray(rows, &res, tenmod.TenantMonitoredObjectStr)

	nextPageStartKey := ""
	if int64(len(res)) == batchPlus1 {
		// Have an extra item, need to remove it and store the key for the next page
		nextPageStartKey = res[batchSize].ObjectName
		res = res[:batchSize]
	}

	paginationOffsets := common.PaginationOffsets{
		Self: startKey,
		Next: nextPageStartKey,
	}

	// Try to retrieve the previous page as well to get the previous start key
	prevPageParams := generatePaginationQueryParams(res[0].ObjectName, batchPlus1, true, true)
	prevPageResponse, err := getByDocIDWithQueryParams(monitoredObjectsByNameIndex, tenmod.TenantMonitoredObjectStr, db, &prevPageParams)
	if err == nil {
		// Try to get previous page details
		if prevPageResponse["rows"] != nil {
			prevPageRows := prevPageResponse["rows"].([]interface{})

			// There will always be 1 result at this point for the start of the current page, only add previous page key if there are actually records on the prev page
			if len(prevPageRows) > 1 {
				lastRow := prevPageRows[len(prevPageRows)-1].(map[string]interface{})
				paginationOffsets.Prev = lastRow["key"].(string)
			}
		}
	}

	logger.Log.Debugf("Retrieved %d %ss from startKey %s\n", len(res), tenmod.TenantMonitoredObjectStr, startKey)
	return res, &paginationOffsets, nil
}

// GetAllMonitoredObjectsInIDList - couchdb implementation of GetAllMonitoredObjectsInIDList
func (tsd *TenantServiceDatastoreCouchDB) GetAllMonitoredObjectsInIDList(tenantID string, idList []string) ([]*tenmod.MonitoredObject, error) {
	logger.Log.Debugf("Fetching all %s\n from list of %d IDs", tenmod.TenantMonitoredObjectStr, len(idList))
	tenantID = ds.PrependToDataID(tenantID, string(admmod.TenantType))

	dbName := createDBPathStr(tsd.server, fmt.Sprintf("%s%s", tenantID, monitoredObjectDBSuffix))
	res := make([]*tenmod.MonitoredObject, 0)
	if err := getAllInIDListFromCouchAndFlatten(dbName, idList, string(tenmod.TenantMonitoredObjectType), tenmod.TenantMonitoredObjectStr, &res); err != nil {
		return nil, err
	}

	logger.Log.Debugf("Retrieved %d %ss for list of %d IDs\n", len(res), tenmod.TenantMonitoredObjectStr, len(idList))
	return res, nil
}

// GetMonitoredObjectToDomainMap - CouchDB implementation of GetMonitoredObjectToDomainMap
func (tsd *TenantServiceDatastoreCouchDB) GetMonitoredObjectToDomainMap(moByDomReq *tenmod.MonitoredObjectCountByDomainRequest) (*tenmod.MonitoredObjectCountByDomainResponse, error) {
	logger.Log.Debugf("Fetching %s: %v\n", tenmod.MonitoredObjectToDomainMapStr, models.AsJSONString(moByDomReq))
	tenantID := ds.PrependToDataID(moByDomReq.TenantID, string(admmod.TenantType))

	dbName := createDBPathStr(tsd.server, fmt.Sprintf("%s%s", tenantID, monitoredObjectDBSuffix))
	db, err := getDatabase(dbName)
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

		fetchResponse, err = fetchDesignDocumentResults(requestBody, dbName, monitoredObjectsByDomainIndex)
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

func createNewTenantMetadataViews(dbName string, key string) error {
	// Create an index based on metadata keys
	err := createCouchDBViewIndex(dbName, metaIndexTemplate, key, []string{key}, metaFieldPrefix)
	if err != nil {
		msg := fmt.Sprintf("Could not create metadata Index for tenant %s, key %s. Error: %s", dbName, key, err.Error())
		return errors.New(msg)
	}
	// Create a view based on unique values per new value
	err = createCouchDBViewIndex(dbName, metaUniqueValuesViewsDdocTemplate, key, []string{key}, metaFieldPrefix)
	if err != nil {
		msg := fmt.Sprintf("Could not create metadata View for tenant %s, key %s. Error: %s", dbName, key, err.Error())
		return errors.New(msg)
	}
	return nil
}

// CheckAndAddMetadataView - Check if we're missing a couchdb view for this new metadata
func (tsd *TenantServiceDatastoreCouchDB) CheckAndAddMetadataView(tenantID string, monitoredObject *tenmod.MonitoredObject) error {

	// _, err := tsd.GetMetadataKeys(tenantID)
	// if err != nil {
	// 	return err
	// }
	// Create the couchDB views
	dbNameKeys := GenerateMonitoredObjectURL(tenantID, tsd.server)
	for key := range monitoredObject.Meta {

		// Check if metadata key is old/new
		/*if _, ok := newKeys[key]; ok {
			// Already in database
			continue
		}*/
		//if err := tsd.CheckMetaDdocExist(tenantID, key); err != nil {
		//logger.Log.Debugf("DDoc for %s does not exist, creating it,%s", key, err.Error())

		// Create an index based on metadata keys
		err := createNewTenantMetadataViews(dbNameKeys, key)
		if err != nil {
			msg := fmt.Sprintf("Could not create metadata Index for tenant %s, key %s. Error: %s", tenantID, key, err.Error())
			//return errors.New(msg)
			// This isn't critical error but log it
			logger.Log.Debug(msg)
		}
		//}
	}
	return nil
}

// UpdateMonitoredObjectMetadataViews - Updates the Tenant's Metadata's Monitored Object Meta key list
func (tsd *TenantServiceDatastoreCouchDB) UpdateMonitoredObjectMetadataViews(tenantID string, monitoredObject *tenmod.MonitoredObject) error {

	// Create the couchDB views
	dbNameKeys := GenerateMonitoredObjectURL(tenantID, tsd.server)
	if monitoredObject != nil {
		tsd.CheckAndAddMetadataView(tenantID, monitoredObject)

		// Now force the indexer to crunch!
		// Do not wait for this to finish, it will certainly take tens of minutes
		// Create/update the couchDB views
		for key := range monitoredObject.Meta {
			go indexViewTriggerBuild(dbNameKeys, MetaKeyIndexOf+key, "by"+key)
			go indexViewTriggerBuild(dbNameKeys, MetaKeyViewOf+key, "by"+key)
		}
	}

	go indexViewTriggerBuild(dbNameKeys, metakeysViewDdocName, MetakeysViewUniqueKeysURI)
	go indexViewTriggerBuild(dbNameKeys, metakeysViewDdocName, metakeysViewUniqueValuessURI)
	go indexViewTriggerBuild(dbNameKeys, metakeysViewDdocName, metaViewSearchLookup)
	go indexViewTriggerBuild(dbNameKeys, metakeysViewDdocName, metaViewLookupWords)
	go indexViewTriggerBuild(dbNameKeys, metakeysViewDdocName, metaViewAllValuesPerKey)

	return nil
}

// GetMetadataKeys - Gets all the known metadata keys from the couchdb view
func (tsd *TenantServiceDatastoreCouchDB) GetMetadataKeys(tenantId string) (map[string]int, error) {
	//https://megatron.npav.accedian.net/couchdb/tenant_2_b4772641-c19b-45fb-ad0e-848de0cfb862_monitored-objects/_design/metaViews/_view/uniqueKeys?group=true

	// model
	type couchViewItem struct {
		Key   string `json:"key"`
		Value int    `json:"value"`
	}
	type couchView struct {
		Rows []couchViewItem `json:"rows"`
	}
	dbName := GenerateMonitoredObjectURL(tenantId, tsd.server)
	db, err := getDatabase(dbName)
	if err != nil {
		return nil, err
	}
	v := url.Values{}
	v.Set("group", "true")

	//doc, err := db.Get("_design/"+metakeysViewDdocName+"/_view/"+MetakeysViewUniqueKeysURI, v)
	url := "_design/metaViews/_view/uniqueKeys"
	doc, err := db.Get(url, v)
	if err != nil {
		return nil, fmt.Errorf("Could not get view %s, %s", dbName+url, err.Error())
	}

	if logger.IsDebugEnabled() {
		logger.Log.Debugf("Getting metadata keys from %s -> %s", dbName+url, models.AsJSONString(doc))
	}

	var resp couchView
	raw, err := json.Marshal(doc)
	if err != nil {
		return nil, fmt.Errorf("Could not marshal (%s) %s, %s", dbName+url, models.AsJSONString(doc), err.Error())
	}
	err = json.Unmarshal(raw, &resp)
	if err != nil {
		return nil, fmt.Errorf("Could not unmarshal (%s) %s into couchView, %s", dbName+url, string(raw), err.Error())
	}

	rows := make(map[string]int, 0)
	for i, item := range resp.Rows {
		rows[item.Key] = resp.Rows[i].Value
	}

	return rows, nil
}

// CheckMetaDdocExist - Gets all the known metadata keys from the couchdb view
func (tsd *TenantServiceDatastoreCouchDB) CheckMetaDdocExist(tenantID string, docname string) error {
	//https://megatron.npav.accedian.net/couchdb/tenant_2_b4772641-c19b-45fb-ad0e-848de0cfb862_monitored-objects/_design/metaViews/_view/uniqueKeys?group=true

	dbName := GenerateMonitoredObjectURL(tenantID, tsd.server)
	db, err := getDatabase(dbName)
	if err != nil {
		return err
	}

	//doc, err := db.Get("_design/"+metakeysViewDdocName+"/_view/"+MetakeysViewUniqueKeysURI, v)
	url := fmt.Sprintf("_design/indexOf%s/_view/by%s", docname, docname)
	err = db.Contains(url)
	if err != nil {
		return fmt.Errorf("Could not get view %s, %s", dbName+url, err.Error())
	}

	url = fmt.Sprintf("_design/viewOf%s/_view/by%s", docname, docname)
	err = db.Contains(url)
	if err != nil {
		return fmt.Errorf("Could not get view %s, %s", dbName+url, err.Error())
	}

	return nil
}

// GetMonitoredObjectByObjectName - Returns an Monitored based on its Object Name
// This is useful because the tool used to ingest data from the NID creates unique IDs but
// the clients may use a different mapping based on monitored object name.
// Assumption right now is that most clients monitored object objectName will be unique.
func (tsd *TenantServiceDatastoreCouchDB) GetMonitoredObjectByObjectName(name string, tenantID string) (*tenmod.MonitoredObject, error) {

	dbName := GenerateMonitoredObjectURL(tenantID, tsd.server)
	db, err := getDatabase(dbName)
	if err != nil {
		return nil, err
	}

	index := "indexOfobjectName"
	selector := fmt.Sprintf(`data.objectName == "%s"`, name)
	// Expect only 1 return
	const expectOnly1Result = 1
	fetchedData, err := db.Query([]string{"_id"}, selector, nil, expectOnly1Result, nil, index)

	if err != nil {
		return nil, err
	}

	id, found := fetchedData[0]["_id"]

	if !found {
		return nil, errors.New(fmt.Sprintf("Could not find mapping of monitored object with name %s to id %s", name, id))
	}

	if logger.IsDebugEnabled() {
		logger.Log.Debugf("Found mapping of monitored object with name %s to id %s", name, id)
	}

	mo := &tenmod.MonitoredObject{}

	fetchedMonObject, err := getByDocID(id.(string), "by objectname", db)

	// Flatten the map so that the unmarshaller can properly build the object
	flatMO := fetchedMonObject["data"].(map[string]interface{})
	flatMO["_id"] = id.(string)
	flatMO["_rev"] = fetchedMonObject["_rev"].(string)

	if err != nil {
		return nil, err
	}
	raw, _ := json.Marshal(flatMO)
	err = json.Unmarshal(raw, mo)
	if err != nil {
		return nil, err
	}

	return mo, nil
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
	if err := getAllOfTypeFromCouchAndFlatten(tenantDBName, string(tenmod.TenantThresholdProfileType), tenmod.TenantThresholdProfileStr, &res); err != nil {
		return nil, err
	}

	logger.Log.Debugf("Retrieved %d %s\n", tenmod.TenantThresholdProfileStr, models.AsJSONString(res))
	return res, nil
}

// BulkInsertMonitoredObjects - CouchDB implementation of BulkInsertMonitoredObjects
func (tsd *TenantServiceDatastoreCouchDB) BulkInsertMonitoredObjects(tenantID string, value []*tenmod.MonitoredObject) ([]*common.BulkOperationResult, error) {
	logger.Log.Debugf("Bulk creating %s: %v\n", tenmod.TenantMonitoredObjectStr, models.AsJSONString(value))
	origTenantID := tenantID
	tenantID = ds.PrependToDataID(tenantID, string(admmod.TenantType))
	dbName := createDBPathStr(tsd.server, fmt.Sprintf("%s%s", tenantID, monitoredObjectDBSuffix))
	resource, err := couchdb.NewResource(dbName, nil)
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

		// Check and add missing metadata views
		err = tsd.CheckAndAddMetadataView(origTenantID, mo)
		if err != nil {
			return nil, err
		}
	}
	body := map[string]interface{}{
		"docs": data}

	fetchedData, err := performBulkUpdate(body, resource)
	if err != nil {
		return nil, err
	}

	// Now that we've done the bulk update, refresh the views
	for _, mo := range value {
		err = tsd.UpdateMonitoredObjectMetadataViews(origTenantID, mo)
		if err != nil {
			return nil, err
		}
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

// BulkUpdateMonitoredObjects - CouchDB implementation of BulkUpdateMonitoredObjects
func (tsd *TenantServiceDatastoreCouchDB) BulkUpdateMonitoredObjects(tenantID string, value []*tenmod.MonitoredObject) ([]*common.BulkOperationResult, error) {
	logger.Log.Debugf("Bulk updating %s: %v\n", tenmod.TenantMonitoredObjectStr, models.AsJSONString(value))
	tenantID = ds.PrependToDataID(tenantID, string(admmod.TenantType))
	dbName := createDBPathStr(tsd.server, fmt.Sprintf("%s%s", tenantID, monitoredObjectDBSuffix))
	resource, err := couchdb.NewResource(dbName, nil)
	if err != nil {
		return nil, err
	}

	// Iterate over the collection and populate necessary fields
	data := make([]map[string]interface{}, 0)
	for _, mo := range value {
		mo.ID = ds.PrependToDataID(mo.ID, string(tenmod.TenantMonitoredObjectType))
		genericMO, err := convertDataToCouchDbSupportedModel(mo)
		if err != nil {
			return nil, err
		}

		dataType := string(tenmod.TenantMonitoredObjectType)
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

	logger.Log.Debugf("Bulk update of %s result: %v\n", tenmod.TenantMonitoredObjectStr, models.AsJSONString(res))
	return res, nil
}

func (tsd *TenantServiceDatastoreCouchDB) CreateReportScheduleConfig(slaConfig *metmod.ReportScheduleConfig) (*metmod.ReportScheduleConfig, error) {
	logger.Log.Debugf("Creating %s: %v\n", metmod.ReportScheduleConfigStr, models.AsJSONString(slaConfig))
	slaConfig.ID = ds.GenerateID(slaConfig, string(metmod.ReportScheduleConfigType))
	tenantID := ds.PrependToDataID(slaConfig.TenantID, string(admmod.TenantType))
	tenantDBName := createDBPathStr(tsd.server, tenantID)

	dataContainer := &metmod.ReportScheduleConfig{}
	if err := createDataInCouch(tenantDBName, slaConfig, dataContainer, string(metmod.ReportScheduleConfigType), metmod.ReportScheduleConfigStr); err != nil {
		return nil, err
	}
	logger.Log.Debugf("Updated %s: %v\n", metmod.ReportScheduleConfigStr, models.AsJSONString(slaConfig))
	return dataContainer, nil

}
func (tsd *TenantServiceDatastoreCouchDB) UpdateReportScheduleConfig(slaConfig *metmod.ReportScheduleConfig) (*metmod.ReportScheduleConfig, error) {
	logger.Log.Debugf("Updating %s: %v\n", metmod.ReportScheduleConfigStr, models.AsJSONString(slaConfig))
	slaConfig.ID = ds.PrependToDataID(slaConfig.ID, string(metmod.ReportScheduleConfigType))
	tenantID := ds.PrependToDataID(slaConfig.TenantID, string(admmod.TenantType))
	tenantDBName := createDBPathStr(tsd.server, tenantID)

	dataContainer := &metmod.ReportScheduleConfig{}
	if err := updateDataInCouch(tenantDBName, slaConfig, dataContainer, string(metmod.ReportScheduleConfigType), metmod.ReportScheduleConfigStr); err != nil {
		return nil, err
	}
	logger.Log.Debugf("Updated %s: %v\n", metmod.ReportScheduleConfigStr, models.AsJSONString(slaConfig))
	return dataContainer, nil
}
func (tsd *TenantServiceDatastoreCouchDB) DeleteReportScheduleConfig(tenantID string, configID string) (*metmod.ReportScheduleConfig, error) {
	logger.Log.Debugf("Deleting %s %s\n", metmod.ReportScheduleConfigStr, configID)
	configID = ds.PrependToDataID(configID, string(metmod.ReportScheduleConfigType))
	tenantID = ds.PrependToDataID(tenantID, string(admmod.TenantType))

	tenantDBName := createDBPathStr(tsd.server, tenantID)
	dataContainer := &metmod.ReportScheduleConfig{}
	if err := deleteDataFromCouch(tenantDBName, configID, &dataContainer, metmod.ReportScheduleConfigStr); err != nil {
		return nil, err
	}
	logger.Log.Debugf("Deleted %s: %v\n", metmod.ReportScheduleConfigStr, models.AsJSONString(dataContainer))
	return dataContainer, nil
}
func (tsd *TenantServiceDatastoreCouchDB) GetReportScheduleConfig(tenantID string, configID string) (*metmod.ReportScheduleConfig, error) {
	logger.Log.Debugf("Fetching %s: %s\n", metmod.ReportScheduleConfigStr, configID)
	configID = ds.PrependToDataID(configID, string(metmod.ReportScheduleConfigType))
	tenantID = ds.PrependToDataID(tenantID, string(admmod.TenantType))

	tenantDBName := createDBPathStr(tsd.server, tenantID)
	dataContainer := &metmod.ReportScheduleConfig{}
	if err := getDataFromCouch(tenantDBName, configID, &dataContainer, metmod.ReportScheduleConfigStr); err != nil {
		return nil, err
	}
	logger.Log.Debugf("Retrieved %s: %v\n", metmod.ReportScheduleConfigStr, models.AsJSONString(dataContainer))
	return dataContainer, nil
}
func (tsd *TenantServiceDatastoreCouchDB) GetAllReportScheduleConfigs(tenantID string) ([]*metmod.ReportScheduleConfig, error) {
	logger.Log.Debugf("Fetching all %s\n", metmod.ReportScheduleConfigStr)
	tenantID = ds.PrependToDataID(tenantID, string(admmod.TenantType))

	tenantDBName := createDBPathStr(tsd.server, tenantID)
	res := make([]*metmod.ReportScheduleConfig, 0)
	if err := getAllOfTypeFromCouchAndFlatten(tenantDBName, string(metmod.ReportScheduleConfigType), metmod.ReportScheduleConfigStr, &res); err != nil {
		return nil, err
	}

	logger.Log.Debugf("Retrieved %d %s\n", len(res), metmod.ReportScheduleConfigStr)
	return res, nil
}

func (tsd *TenantServiceDatastoreCouchDB) CreateSLAReport(slaReport *metmod.SLAReport) (*metmod.SLAReport, error) {
	logger.Log.Debugf("Creating %s: %v\n", tenmod.TenantSLAReportStr, models.AsJSONString(slaReport))
	slaReport.ID = ds.GenerateID(slaReport, string(tenmod.TenantReportType))
	tenantID := ds.PrependToDataID(slaReport.TenantID, string(admmod.TenantType))
	tenantDBName := createDBPathStr(tsd.server, fmt.Sprintf("%s%s", tenantID, reportObjectDBSuffix))

	dataContainer := &metmod.SLAReport{}
	if err := createDataInCouch(tenantDBName, slaReport, dataContainer, string(tenmod.TenantReportType), tenmod.TenantSLAReportStr); err != nil {
		return nil, err
	}
	logger.Log.Debugf("Created %s: %v\n", tenmod.TenantSLAReportStr, models.AsJSONString(dataContainer))
	return dataContainer, nil
}
func (tsd *TenantServiceDatastoreCouchDB) DeleteSLAReport(tenantID string, slaReportID string) (*metmod.SLAReport, error) {
	logger.Log.Debugf("Fetching %s: %s\n", tenmod.TenantSLAReportStr, slaReportID)
	slaReportID = ds.PrependToDataID(slaReportID, string(tenmod.TenantReportType))
	tenantID = ds.PrependToDataID(tenantID, string(admmod.TenantType))
	tenantDBName := createDBPathStr(tsd.server, fmt.Sprintf("%s%s", tenantID, reportObjectDBSuffix))

	dataContainer := &metmod.SLAReport{}
	if err := deleteDataFromCouch(tenantDBName, slaReportID, &dataContainer, tenmod.TenantSLAReportStr); err != nil {
		return nil, err
	}
	logger.Log.Debugf("Retrieved %s: %v\n", tenmod.TenantSLAReportStr, models.AsJSONString(dataContainer))
	return dataContainer, nil
}
func (tsd *TenantServiceDatastoreCouchDB) GetSLAReport(tenantID string, slaReportID string) (*metmod.SLAReport, error) {
	logger.Log.Debugf("Fetching %s: %s\n", tenmod.TenantSLAReportStr, slaReportID)
	slaReportID = ds.PrependToDataID(slaReportID, string(tenmod.TenantReportType))
	tenantID = ds.PrependToDataID(tenantID, string(admmod.TenantType))
	tenantDBName := createDBPathStr(tsd.server, fmt.Sprintf("%s%s", tenantID, reportObjectDBSuffix))

	dataContainer := &metmod.SLAReport{}
	if err := getDataFromCouch(tenantDBName, slaReportID, &dataContainer, tenmod.TenantSLAReportStr); err != nil {
		return nil, err
	}
	logger.Log.Debugf("Retrieved %s: %v\n", tenmod.TenantSLAReportStr, models.AsJSONString(dataContainer))
	return dataContainer, nil
}
func (tsd *TenantServiceDatastoreCouchDB) GetAllSLAReports(tenantID string) ([]*metmod.SLAReport, error) {
	logger.Log.Debugf("Fetching all %s\n", tenmod.TenantSLAReportStr)
	tenantID = ds.PrependToDataID(tenantID, string(admmod.TenantType))
	tenantDBName := createDBPathStr(tsd.server, fmt.Sprintf("%s%s", tenantID, reportObjectDBSuffix))

	res := make([]*metmod.SLAReport, 0)
	if err := getAllOfTypeFromCouchAndFlatten(tenantDBName, string(tenmod.TenantReportType), tenmod.TenantSLAReportStr, &res); err != nil {
		return nil, err
	}

	logger.Log.Debugf("Retrieved %d %s\n", len(res), tenmod.TenantSLAReportStr)
	return res, nil
}
func (tsd *TenantServiceDatastoreCouchDB) CreateDashboard(dashboard *tenmod.Dashboard) (*tenmod.Dashboard, error) {
	dashboard.ID = ds.GenerateID(dashboard, string(tenmod.TenantDashboardType))
	tenantID := ds.PrependToDataID(dashboard.TenantID, string(admmod.TenantType))

	tenantDBName := createDBPathStr(tsd.server, tenantID)
	db, err := getDatabase(tenantDBName)
	if err != nil {
		return nil, err
	}

	// Convert to generic object
	storeFormat, err := convertDataToCouchDbSupportedModel(dashboard)
	if err != nil {
		return nil, err
	}

	_, _, err = storeDataInCouchDB(storeFormat, tenmod.TenantDashboardStr, db)
	if err != nil {
		return nil, err
	}

	stripPrefixFromID(storeFormat)

	// Populate the response
	dataContainer := &tenmod.Dashboard{}
	if err = convertGenericCouchDataToObject(storeFormat, dataContainer, tenmod.TenantDashboardStr); err != nil {
		return nil, err
	}

	return dataContainer, nil

}

func (tsd *TenantServiceDatastoreCouchDB) DeleteDashboard(tenantID string, dataID string) (*tenmod.Dashboard, error) {
	logger.Log.Debugf("Deleting %s: %s\n", tenmod.TenantDashboardStr, dataID)
	dataID = ds.PrependToDataID(dataID, string(tenmod.TenantDashboardType))
	tenantID = ds.PrependToDataID(tenantID, string(admmod.TenantType))

	tenantDBName := createDBPathStr(tsd.server, tenantID)
	dataContainer := tenmod.Dashboard{}
	if err := deleteDataFromCouch(tenantDBName, dataID, &dataContainer, tenmod.TenantDashboardStr); err != nil {
		return nil, err
	}
	logger.Log.Debugf("Deleted %s: %v\n", tenmod.TenantDashboardStr, models.AsJSONString(dataContainer))
	return &dataContainer, nil
}

func (tsd *TenantServiceDatastoreCouchDB) HasDashboardsWithDomain(tenantID string, domainID string) (bool, error) {
	logger.Log.Debugf("Fetching %s: %s\n", tenmod.TenantMetaStr, tenantID)
	tenantID = ds.PrependToDataID(tenantID, string(admmod.TenantType))

	tenantDBName := createDBPathStr(tsd.server, tenantID)
	db, err := getDatabase(tenantDBName)
	if err != nil {
		return false, err
	}

	fetchedData, err := getAllOfTypeByIDPrefix(string(tenmod.TenantDashboardType), tenmod.TenantDashboardStr, db)
	if err != nil {
		return false, err
	}

	logger.Log.Debugf("fetched %s", models.AsJSONString(fetchedData))
	for _, d := range fetchedData {
		data, ok := d["data"]
		if !ok {
			continue
		}
		if _, hasDomainIDs := data.(map[string]interface{})["domainSet"]; hasDomainIDs {
			dataContainer := tenmod.Dashboard{}
			if err = convertGenericCouchDataToObject(d, &dataContainer, tenmod.TenantDashboardStr); err != nil {
				return false, err
			}

			for _, v := range dataContainer.DomainSet {
				if v == domainID {
					return true, nil
				}
			}

		}
	}
	return false, nil
}

// CreateTenantDataCleaningProfile - CouchDB implementation of CreateTenantDataCleaningProfile
func (tsd *TenantServiceDatastoreCouchDB) CreateTenantDataCleaningProfile(dcp *tenmod.DataCleaningProfile) (*tenmod.DataCleaningProfile, error) {
	logger.Log.Debugf("Creating %s: %v\n", tenmod.TenantDataCleaningProfileStr, models.AsJSONString(dcp))
	dcp.ID = ds.GenerateID(dcp, string(tenmod.TenantDataCleaningProfileType))
	tenantID := ds.PrependToDataID(dcp.TenantID, string(admmod.TenantType))

	// Only create one if one does not already exist:
	existing, _ := tsd.GetAllTenantDataCleaningProfiles(dcp.TenantID)
	if len(existing) != 0 {
		return nil, fmt.Errorf("Can't create %s, it already exists", tenmod.TenantDataCleaningProfileStr)
	}

	tenantDBName := createDBPathStr(tsd.server, tenantID)
	dataContainer := &tenmod.DataCleaningProfile{}
	if err := createDataInCouch(tenantDBName, dcp, dataContainer, string(tenmod.TenantDataCleaningProfileType), tenmod.TenantDataCleaningProfileStr); err != nil {
		return nil, err
	}
	logger.Log.Debugf("Created %s: %v\n", tenmod.TenantDataCleaningProfileStr, models.AsJSONString(dataContainer))
	return dataContainer, nil
}

// UpdateTenantDataCleaningProfile - CouchDB implementation of UpdateTenantDataCleaningProfile
func (tsd *TenantServiceDatastoreCouchDB) UpdateTenantDataCleaningProfile(dcp *tenmod.DataCleaningProfile) (*tenmod.DataCleaningProfile, error) {
	logger.Log.Debugf("Updating %s: %v\n", tenmod.TenantDataCleaningProfileStr, models.AsJSONString(dcp))
	dcp.ID = ds.PrependToDataID(dcp.ID, string(tenmod.TenantDataCleaningProfileType))
	tenantID := ds.PrependToDataID(dcp.TenantID, string(admmod.TenantType))

	tenantDBName := createDBPathStr(tsd.server, tenantID)
	dataContainer := &tenmod.DataCleaningProfile{}
	if err := updateDataInCouch(tenantDBName, dcp, dataContainer, string(tenmod.TenantDataCleaningProfileType), tenmod.TenantDataCleaningProfileStr); err != nil {
		return nil, err
	}
	logger.Log.Debugf("Updated %s: %v\n", tenmod.TenantDataCleaningProfileStr, models.AsJSONString(dataContainer))
	return dataContainer, nil
}

// DeleteTenantDataCleaningProfile - CouchDB implementation of DeleteTenantDataCleaningProfile
func (tsd *TenantServiceDatastoreCouchDB) DeleteTenantDataCleaningProfile(tenantID string, dataID string) (*tenmod.DataCleaningProfile, error) {
	logger.Log.Debugf("Deleting %s: %s\n", tenmod.TenantDataCleaningProfileStr, tenantID)

	// Obtain the value of the existing record for a return value.
	existing, err := tsd.GetAllTenantDataCleaningProfiles(tenantID)
	if err != nil || len(existing) == 0 {
		return nil, fmt.Errorf("Unable to fetch %s to delete: %s", tenmod.TenantDataCleaningProfileStr, ds.NotFoundStr)
	}

	existingObject := existing[0]

	if existingObject.ID != dataID {
		return nil, fmt.Errorf("%s %s: %s", tenmod.TenantDataCleaningProfileStr, dataID, ds.NotFoundStr)
	}

	tenantID = ds.PrependToDataID(tenantID, string(admmod.TenantType))
	tenantDBName := createDBPathStr(tsd.server, tenantID)
	objectID := ds.PrependToDataID(existingObject.ID, string(tenmod.TenantDataCleaningProfileType))
	if err := deleteData(tenantDBName, objectID, tenmod.TenantDataCleaningProfileStr); err != nil {
		logger.Log.Debugf("Unable to delete %s: %s", tenmod.TenantDataCleaningProfileStr, err.Error())
		return nil, err
	}

	// Return the deleted object.
	logger.Log.Debugf("Deleted %s: %v\n", tenmod.TenantDataCleaningProfileStr, models.AsJSONString(existingObject))
	return existingObject, nil
}

// GetTenantDataCleaningProfile - CouchDB implementation of GetTenantDataCleaningProfile
func (tsd *TenantServiceDatastoreCouchDB) GetTenantDataCleaningProfile(tenantID string, dataID string) (*tenmod.DataCleaningProfile, error) {
	logger.Log.Debugf("Fetching %s: %s\n", tenmod.TenantDataCleaningProfileStr, dataID)
	dataID = ds.PrependToDataID(dataID, string(tenmod.TenantDataCleaningProfileType))
	tenantID = ds.PrependToDataID(tenantID, string(admmod.TenantType))

	tenantDBName := createDBPathStr(tsd.server, tenantID)
	dataContainer := tenmod.DataCleaningProfile{}
	if err := getDataFromCouch(tenantDBName, dataID, &dataContainer, tenmod.TenantDataCleaningProfileStr); err != nil {
		return nil, err
	}
	logger.Log.Debugf("Retrieved %s: %v\n", tenmod.TenantDataCleaningProfileStr, models.AsJSONString(dataContainer))
	return &dataContainer, nil
}

// GetAllTenantDataCleaningProfiles - CouchDB implementation of GetAllTenantDataCleaningProfiles
func (tsd *TenantServiceDatastoreCouchDB) GetAllTenantDataCleaningProfiles(tenantID string) ([]*tenmod.DataCleaningProfile, error) {
	logger.Log.Debugf("Fetching all %ss for Tenant %s\n", tenmod.TenantDataCleaningProfileStr, tenantID)
	tenantID = ds.PrependToDataID(tenantID, string(admmod.TenantType))

	tenantDBName := createDBPathStr(tsd.server, tenantID)
	res := make([]*tenmod.DataCleaningProfile, 0)
	if err := getAllOfTypeFromCouchAndFlatten(tenantDBName, string(tenmod.TenantDataCleaningProfileType), tenmod.TenantDataCleaningProfileStr, &res); err != nil {
		return nil, err
	}

	if len(res) == 0 {
		return nil, fmt.Errorf("%ss: %s", tenmod.TenantDataCleaningProfileStr, ds.NotFoundStr)
	}

	logger.Log.Debugf("Retrieved %d %ss\n", len(res), tenmod.TenantDataCleaningProfileStr)
	return res, nil
}

// GetMonitoredObjectIDsToMetaEntry - CouchDB implementation to retrieve all monitored object Ids associated with a specific metadata key/value pair
func (tsd *TenantServiceDatastoreCouchDB) GetMonitoredObjectIDsToMetaEntry(tenantID string, metakey string, metavalue string) ([]string, error) {
	logger.Log.Debugf("Fetching all %ss for Tenant %s with meta kay %s and value %s\n", tenmod.TenantMonitoredObjectKeysStr, tenantID, metakey, metavalue)

	timeStart := time.Now()
	tenantMODB := createDBPathStr(tsd.server, fmt.Sprintf("%s_monitored-objects", ds.PrependToDataID(tenantID, string(admmod.TenantType))))

	res, err := getIDsByView(tenantMODB, fmt.Sprintf("indexOf%s", metakey), fmt.Sprintf("by%s", metakey), metavalue)
	if err != nil {
		return nil, err
	}
	mon.TrackAPITimeMetricInSeconds(timeStart, "200", mon.DbGetIDByViewStr)

	return res, nil
}

// GetAllMonitoredObjectsIDs - uses the paginated DB call to acquire all monitored objects' ID
// Becarefully calling this function, it may take a long time to process
func (tsd *TenantServiceDatastoreCouchDB) GetAllMonitoredObjectsIDs(tenantID string) ([]string, error) {
	timeStart := time.Now()
	//ogger.Log.Debugf("Fetching next %d %ss from startKey %s\n", limit, tenmod.TenantMonitoredObjectStr, startKey)
	tenantID = ds.PrependToDataID(tenantID, string(admmod.TenantType))

	dbName := createDBPathStr(tsd.server, fmt.Sprintf("%s%s", tenantID, monitoredObjectDBSuffix))
	db, err := getDatabase(dbName)
	if err != nil {
		return nil, err
	}

	fetchResponse, err := getByDocIDWithQueryParams(monitoredObjectsByNameIndex, tenmod.TenantMonitoredObjectStr, db, nil)
	if err != nil {
		return nil, err
	}

	if fetchResponse["rows"] == nil {
		return nil, fmt.Errorf(ds.NotFoundStr)
	}

	castedRows := fetchResponse["rows"].([]interface{})
	if len(castedRows) == 0 {
		return nil, fmt.Errorf(ds.NotFoundStr)
	}
	var ids []string
	moCount := 0
	// Convert interface results to map results
	for _, obj := range castedRows {
		castedObj := obj.(map[string]interface{})
		genericDoc := castedObj["id"].(string)
		moID := ds.GetDataIDFromFullID(genericDoc)

		ids = append(ids, moID)
		moCount = moCount + 1
	}

	// Update counters
	mon.MonitoredObjectCounter.Set(float64(moCount))

	logger.Log.Debugf("Retrieved %d items from %ss\n", len(ids), tenmod.TenantMonitoredObjectStr)
	mon.TrackAPITimeMetricInSeconds(timeStart, "200", mon.DbGetAllMoIDStr)

	return ids, nil

}
