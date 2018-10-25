package couchDB

import (
	"fmt"

	ds "github.com/accedian/adh-gather/datastore"
	"github.com/accedian/adh-gather/logger"
	"github.com/accedian/adh-gather/models"
	admmod "github.com/accedian/adh-gather/models/admin"
	tenmod "github.com/accedian/adh-gather/models/tenant"
)

// CreateTenantMetadataConfig - CouchDB implementation of CreateTenantMetadataConfig
func (tsd *TenantServiceDatastoreCouchDB) CreateTenantMetadataConfig(metadataConfigReq *tenmod.MetadataConfig) (*tenmod.MetadataConfig, error) {
	logger.Log.Debugf("Creating %s: %v\n", tenmod.TenantMetadataConfigStr, models.AsJSONString(metadataConfigReq))
	metadataConfigReq.ID = ds.GenerateID(metadataConfigReq, string(tenmod.TenantMetadataConfigType))
	tenantID := ds.PrependToDataID(metadataConfigReq.TenantID, string(admmod.TenantType))

	// Only create one if one does not already exist:
	existing, _ := tsd.GetActiveTenantMetadataConfig(metadataConfigReq.TenantID)
	if existing != nil {
		return nil, fmt.Errorf("Can't create %s, it already exists", tenmod.TenantMetadataConfigStr)
	}

	tenantDBName := createDBPathStr(tsd.server, tenantID)
	dataContainer := &tenmod.MetadataConfig{}
	if err := createDataInCouch(tenantDBName, metadataConfigReq, dataContainer, string(tenmod.TenantMetadataConfigType), tenmod.TenantMetadataConfigStr); err != nil {
		return nil, err
	}
	logger.Log.Debugf("Created %s: %v\n", tenmod.TenantMetadataConfigStr, models.AsJSONString(dataContainer))
	return dataContainer, nil
}

// UpdateTenantMetadataConfig - CouchDB implementation of UpdateTenantMetadataConfig
func (tsd *TenantServiceDatastoreCouchDB) UpdateTenantMetadataConfig(metadataConfigReq *tenmod.MetadataConfig) (*tenmod.MetadataConfig, error) {
	logger.Log.Debugf("Updating %s: %v\n", tenmod.TenantMetadataConfigStr, models.AsJSONString(metadataConfigReq))
	metadataConfigReq.ID = ds.PrependToDataID(metadataConfigReq.ID, string(tenmod.TenantMetadataConfigType))
	tenantID := ds.PrependToDataID(metadataConfigReq.TenantID, string(admmod.TenantType))

	tenantDBName := createDBPathStr(tsd.server, tenantID)
	dataContainer := &tenmod.MetadataConfig{}
	if err := updateDataInCouch(tenantDBName, metadataConfigReq, dataContainer, string(tenmod.TenantMetadataConfigType), tenmod.TenantMetadataConfigStr); err != nil {
		return nil, err
	}
	logger.Log.Debugf("Updated %s: %v\n", tenmod.TenantMetadataConfigStr, models.AsJSONString(dataContainer))
	return dataContainer, nil
}

// GetTenantMetadataConfig - CouchDB implementation of GetTenantMetadataConfig
func (tsd *TenantServiceDatastoreCouchDB) GetTenantMetadataConfig(tenantID string, dataID string) (*tenmod.MetadataConfig, error) {
	logger.Log.Debugf("Fetching %s: %s\n", tenmod.TenantMetadataConfigStr, dataID)
	dataID = ds.PrependToDataID(dataID, string(tenmod.TenantMetadataConfigType))
	tenantID = ds.PrependToDataID(tenantID, string(admmod.TenantType))

	tenantDBName := createDBPathStr(tsd.server, tenantID)
	dataContainer := tenmod.MetadataConfig{}
	if err := getDataFromCouch(tenantDBName, dataID, &dataContainer, tenmod.TenantMetadataConfigStr); err != nil {
		return nil, err
	}

	logger.Log.Debugf("Retrieved %s: %v\n", tenmod.TenantMetadataConfigStr, models.AsJSONString(dataContainer))
	return &dataContainer, nil
}

// DeleteTenantMetadataConfig - CouchDB implementation of DeleteTenantMetadataConfig
func (tsd *TenantServiceDatastoreCouchDB) DeleteTenantMetadataConfig(tenantID string, dataID string) (*tenmod.MetadataConfig, error) {
	logger.Log.Debugf("Deleting %s: %s\n", tenmod.TenantMetadataConfigStr, dataID)
	dataID = ds.PrependToDataID(dataID, string(tenmod.TenantMetadataConfigType))
	tenantID = ds.PrependToDataID(tenantID, string(admmod.TenantType))

	tenantDBName := createDBPathStr(tsd.server, tenantID)
	dataContainer := tenmod.MetadataConfig{}
	if err := deleteDataFromCouch(tenantDBName, dataID, &dataContainer, tenmod.TenantMetadataConfigStr); err != nil {
		return nil, err
	}

	logger.Log.Debugf("Deleted %s: %v\n", tenmod.TenantMetadataConfigStr, models.AsJSONString(dataContainer))
	return &dataContainer, nil
}

// GetActiveTenantMetadataConfig - CouchDB implementation of GetActiveTenantMetadataConfig
func (tsd *TenantServiceDatastoreCouchDB) GetActiveTenantMetadataConfig(tenantID string) (*tenmod.MetadataConfig, error) {
	logger.Log.Debugf("Fetching active %s for Tenant %s\n", tenmod.TenantMetadataConfigStr, tenantID)
	tenantID = ds.PrependToDataID(tenantID, string(admmod.TenantType))

	tenantDBName := createDBPathStr(tsd.server, tenantID)
	db, err := getDatabase(tenantDBName)
	if err != nil {
		return nil, err
	}

	fetchedData, err := getAllOfTypeByIDPrefix(string(tenmod.TenantMetadataConfigType), tenmod.TenantMetadataConfigStr, db)
	if err != nil {
		return nil, err
	}

	logger.Log.Debugf("Retrieved: %s", models.AsJSONString(fetchedData))

	// Populate the response
	res := tenmod.MetadataConfig{}
	if len(fetchedData) != 0 {
		if err = convertGenericCouchDataToObject(fetchedData[0], &res, tenmod.TenantMetadataConfigStr); err != nil {
			return nil, err
		}
	} else {
		return nil, fmt.Errorf(ds.NotFoundStr)
	}

	logger.Log.Debugf("Retrieved %s: %v\n", tenmod.TenantMetadataConfigStr, models.AsJSONString(res))
	return &res, nil
}
