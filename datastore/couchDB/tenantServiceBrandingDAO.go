package couchDB

import (
	"fmt"

	ds "github.com/accedian/adh-gather/datastore"
	"github.com/accedian/adh-gather/logger"
	"github.com/accedian/adh-gather/models"
	admmod "github.com/accedian/adh-gather/models/admin"
	tenmod "github.com/accedian/adh-gather/models/tenant"
)

// CreateTenantBranding - CouchDB implementation of CreateTenantBranding
func (tsd *TenantServiceDatastoreCouchDB) CreateTenantBranding(tenantBrandingRequest *tenmod.Branding) (*tenmod.Branding, error) {
	logger.Log.Debugf("Creating %s: %v\n", tenmod.TenantBrandingStr, models.AsJSONString(tenantBrandingRequest))
	id := tenantBrandingRequest.ID
	if id == "" {
		tenantBrandingRequest.ID = ds.GenerateID(tenantBrandingRequest, string(tenmod.TenantBrandingType))
	} else {
		tenantBrandingRequest.ID = ds.PrependToDataID(id, string(tenmod.TenantBrandingType))
	}

	tenantID := ds.PrependToDataID(tenantBrandingRequest.TenantID, string(admmod.TenantType))

	tenantDBName := createDBPathStr(tsd.server, tenantID)
	dataContainer := &tenmod.Branding{}
	if err := createDataInCouch(tenantDBName, tenantBrandingRequest, dataContainer, string(tenmod.TenantBrandingType), tenmod.TenantBrandingStr); err != nil {
		return nil, err
	}
	logger.Log.Debugf("Created %s: %v\n", tenmod.TenantBrandingStr, models.AsJSONString(dataContainer))
	return dataContainer, nil
}

// UpdateTenantBranding - CouchDB implementation of UpdateTenantBranding
func (tsd *TenantServiceDatastoreCouchDB) UpdateTenantBranding(tenantBrandingRequest *tenmod.Branding) (*tenmod.Branding, error) {
	logger.Log.Debugf("Updating %s: %v\n", tenmod.TenantBrandingStr, models.AsJSONString(tenantBrandingRequest))
	tenantBrandingRequest.ID = ds.PrependToDataID(tenantBrandingRequest.ID, string(tenmod.TenantBrandingType))
	tenantID := ds.PrependToDataID(tenantBrandingRequest.TenantID, string(admmod.TenantType))

	tenantDBName := createDBPathStr(tsd.server, tenantID)
	dataContainer := &tenmod.Branding{}
	if err := updateDataInCouch(tenantDBName, tenantBrandingRequest, dataContainer, string(tenmod.TenantBrandingType), tenmod.TenantBrandingStr); err != nil {
		return nil, err
	}
	logger.Log.Debugf("Updated %s: %v\n", tenmod.TenantBrandingStr, models.AsJSONString(dataContainer))
	return dataContainer, nil
}

// DeleteTenantBranding - CouchDB implementation of DeleteTenantBranding
func (tsd *TenantServiceDatastoreCouchDB) DeleteTenantBranding(tenantID string, dataID string) (*tenmod.Branding, error) {
	logger.Log.Debugf("Deleting %s: %s\n", tenmod.TenantBrandingStr, dataID)
	dataID = ds.PrependToDataID(dataID, string(tenmod.TenantBrandingType))
	tenantID = ds.PrependToDataID(tenantID, string(admmod.TenantType))

	tenantDBName := createDBPathStr(tsd.server, tenantID)
	dataContainer := tenmod.Branding{}
	if err := deleteDataFromCouch(tenantDBName, dataID, &dataContainer, tenmod.TenantBrandingStr); err != nil {
		return nil, err
	}
	logger.Log.Debugf("Deleted %s: %v\n", tenmod.TenantBrandingStr, models.AsJSONString(dataContainer))
	return &dataContainer, nil
}

// GetTenantBranding - CouchDB implementation of GetTenantBranding
func (tsd *TenantServiceDatastoreCouchDB) GetTenantBranding(tenantID string, dataID string) (*tenmod.Branding, error) {
	logger.Log.Debugf("Fetching %s: %s\n", tenmod.TenantBrandingStr, dataID)
	dataID = ds.PrependToDataID(dataID, string(tenmod.TenantBrandingType))
	tenantID = ds.PrependToDataID(tenantID, string(admmod.TenantType))

	tenantDBName := createDBPathStr(tsd.server, tenantID)
	dataContainer := tenmod.Branding{}
	if err := getDataFromCouch(tenantDBName, dataID, &dataContainer, tenmod.TenantBrandingStr); err != nil {
		return nil, err
	}
	logger.Log.Debugf("Retrieved %s: %v\n", tenmod.TenantBrandingStr, models.AsJSONString(dataContainer))
	return &dataContainer, nil
}

// GetAllTenantBrandings - CouchDB implementation of GetAllTenantBrandings
func (tsd *TenantServiceDatastoreCouchDB) GetAllTenantBrandings(tenantID string) ([]*tenmod.Branding, error) {
	logger.Log.Debugf("Fetching all %s\n", tenmod.TenantBrandingStr)
	tenantID = ds.PrependToDataID(tenantID, string(admmod.TenantType))

	tenantDBName := createDBPathStr(tsd.server, tenantID)
	res := make([]*tenmod.Branding, 0)

	if err := getAllOfTypeFromCouchAndFlatten(tenantDBName, string(tenmod.TenantBrandingType), tenmod.TenantBrandingStr, &res); err != nil {
		return nil, err
	}

	if len(res) == 0 {
		return nil, fmt.Errorf(ds.NotFoundStr)
	}

	logger.Log.Debugf("Retrieved %d %s\n", len(res), tenmod.TenantBrandingStr)
	return res, nil
}
