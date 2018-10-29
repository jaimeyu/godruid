package couchDB

import (
	"fmt"

	ds "github.com/accedian/adh-gather/datastore"
	"github.com/accedian/adh-gather/logger"
	"github.com/accedian/adh-gather/models"
	admmod "github.com/accedian/adh-gather/models/admin"
	tenmod "github.com/accedian/adh-gather/models/tenant"
)

// CreateTenantLocale - CouchDB implementation of CreateTenantLocale
func (tsd *TenantServiceDatastoreCouchDB) CreateTenantLocale(tenantLocaleRequest *tenmod.Locale) (*tenmod.Locale, error) {
	logger.Log.Debugf("Creating %s: %v\n", tenmod.TenantLocaleStr, models.AsJSONString(tenantLocaleRequest))
	id := tenantLocaleRequest.ID
	if id == "" {
		tenantLocaleRequest.ID = ds.GenerateID(tenantLocaleRequest, string(tenmod.TenantLocaleType))
	} else {
		tenantLocaleRequest.ID = ds.PrependToDataID(id, string(tenmod.TenantLocaleType))
	}

	tenantID := ds.PrependToDataID(tenantLocaleRequest.TenantID, string(admmod.TenantType))

	tenantDBName := createDBPathStr(tsd.server, tenantID)
	dataContainer := &tenmod.Locale{}
	if err := createDataInCouch(tenantDBName, tenantLocaleRequest, dataContainer, string(tenmod.TenantLocaleType), tenmod.TenantLocaleStr); err != nil {
		return nil, err
	}
	logger.Log.Debugf("Created %s: %v\n", tenmod.TenantLocaleStr, models.AsJSONString(dataContainer))
	return dataContainer, nil
}

// UpdateTenantLocale - CouchDB implementation of UpdateTenantLocale
func (tsd *TenantServiceDatastoreCouchDB) UpdateTenantLocale(tenantLocaleRequest *tenmod.Locale) (*tenmod.Locale, error) {
	logger.Log.Debugf("Updating %s: %v\n", tenmod.TenantLocaleStr, models.AsJSONString(tenantLocaleRequest))
	tenantLocaleRequest.ID = ds.PrependToDataID(tenantLocaleRequest.ID, string(tenmod.TenantLocaleType))
	tenantID := ds.PrependToDataID(tenantLocaleRequest.TenantID, string(admmod.TenantType))

	tenantDBName := createDBPathStr(tsd.server, tenantID)
	dataContainer := &tenmod.Locale{}
	if err := updateDataInCouch(tenantDBName, tenantLocaleRequest, dataContainer, string(tenmod.TenantLocaleType), tenmod.TenantLocaleStr); err != nil {
		return nil, err
	}
	logger.Log.Debugf("Updated %s: %v\n", tenmod.TenantLocaleStr, models.AsJSONString(dataContainer))
	return dataContainer, nil
}

// DeleteTenantLocale - CouchDB implementation of DeleteTenantLocale
func (tsd *TenantServiceDatastoreCouchDB) DeleteTenantLocale(tenantID string, dataID string) (*tenmod.Locale, error) {
	logger.Log.Debugf("Deleting %s: %s\n", tenmod.TenantLocaleStr, dataID)
	dataID = ds.PrependToDataID(dataID, string(tenmod.TenantLocaleType))
	tenantID = ds.PrependToDataID(tenantID, string(admmod.TenantType))

	tenantDBName := createDBPathStr(tsd.server, tenantID)
	dataContainer := tenmod.Locale{}
	if err := deleteDataFromCouch(tenantDBName, dataID, &dataContainer, tenmod.TenantLocaleStr); err != nil {
		return nil, err
	}
	logger.Log.Debugf("Deleted %s: %v\n", tenmod.TenantLocaleStr, models.AsJSONString(dataContainer))
	return &dataContainer, nil
}

// GetTenantLocale - CouchDB implementation of GetTenantLocale
func (tsd *TenantServiceDatastoreCouchDB) GetTenantLocale(tenantID string, dataID string) (*tenmod.Locale, error) {
	logger.Log.Debugf("Fetching %s: %s\n", tenmod.TenantLocaleStr, dataID)
	dataID = ds.PrependToDataID(dataID, string(tenmod.TenantLocaleType))
	tenantID = ds.PrependToDataID(tenantID, string(admmod.TenantType))

	tenantDBName := createDBPathStr(tsd.server, tenantID)
	dataContainer := tenmod.Locale{}
	if err := getDataFromCouch(tenantDBName, dataID, &dataContainer, tenmod.TenantLocaleStr); err != nil {
		return nil, err
	}
	logger.Log.Debugf("Retrieved %s: %v\n", tenmod.TenantLocaleStr, models.AsJSONString(dataContainer))
	return &dataContainer, nil
}

// GetAllTenantLocales - CouchDB implementation of GetAllTenantLocales
func (tsd *TenantServiceDatastoreCouchDB) GetAllTenantLocales(tenantID string) ([]*tenmod.Locale, error) {
	logger.Log.Debugf("Fetching all %s\n", tenmod.TenantLocaleStr)
	tenantID = ds.PrependToDataID(tenantID, string(admmod.TenantType))

	tenantDBName := createDBPathStr(tsd.server, tenantID)
	res := make([]*tenmod.Locale, 0)

	if err := getAllOfTypeFromCouchAndFlatten(tenantDBName, string(tenmod.TenantLocaleType), tenmod.TenantLocaleStr, &res); err != nil {
		return nil, err
	}

	if len(res) == 0 {
		return nil, fmt.Errorf(ds.NotFoundStr)
	}

	logger.Log.Debugf("Retrieved %d %s\n", len(res), tenmod.TenantLocaleStr)
	return res, nil
}
