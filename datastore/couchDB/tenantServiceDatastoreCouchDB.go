package couchDB

import (
	"encoding/json"
	"fmt"

	"github.com/accedian/adh-gather/config"
	ds "github.com/accedian/adh-gather/datastore"
	"github.com/accedian/adh-gather/gather"
	"github.com/accedian/adh-gather/logger"
	"github.com/accedian/adh-gather/models"

	admmod "github.com/accedian/adh-gather/models/admin"
	"github.com/accedian/adh-gather/models/common"
	metmod "github.com/accedian/adh-gather/models/metrics"
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
	server              string
	cfg                 config.Provider
	connectorUpdateChan chan *tenmod.ConnectorConfig
}

// CreateTenantServiceDAO - instantiates a CouchDB implementation of the
// TenantServiceDatastore.
func CreateTenantServiceDAO() (*TenantServiceDatastoreCouchDB, error) {
	result := new(TenantServiceDatastoreCouchDB)
	result.cfg = gather.GetConfig()
	result.connectorUpdateChan = make(chan *tenmod.ConnectorConfig)

	provDBURL := fmt.Sprintf("%s:%d",
		result.cfg.GetString(gather.CK_server_datastore_ip.String()),
		result.cfg.GetInt(gather.CK_server_datastore_port.String()))
	logger.Log.Debugf("Tenant Service CouchDB URL is: %s, %v", provDBURL, result.connectorUpdateChan)
	result.server = provDBURL

	return result, nil
}

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
	if err := createDataInCouch(dbName, monitoredObjectReq, dataContainer, string(tenmod.TenantMonitoredObjectType), tenmod.TenantMonitoredObjectStr); err != nil {
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
