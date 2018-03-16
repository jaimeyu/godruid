package datastore

import (
	"github.com/accedian/adh-gather/models/common"
	tenmod "github.com/accedian/adh-gather/models/tenant"
)

// TenantServiceDatastore - interface which provides the functionality
// of the TenantService Datastore.
type TenantServiceDatastore interface {
	CreateTenantUser(*tenmod.User) (*tenmod.User, error)
	UpdateTenantUser(*tenmod.User) (*tenmod.User, error)
	DeleteTenantUser(tenantID string, userID string) (*tenmod.User, error)
	GetTenantUser(tenantID string, userID string) (*tenmod.User, error)
	GetAllTenantUsers(string) ([]*tenmod.User, error)

	CreateTenantDomain(*tenmod.Domain) (*tenmod.Domain, error)
	UpdateTenantDomain(*tenmod.Domain) (*tenmod.Domain, error)
	DeleteTenantDomain(tenantID string, dataID string) (*tenmod.Domain, error)
	GetTenantDomain(tenantID string, dataID string) (*tenmod.Domain, error)
	GetAllTenantDomains(string) ([]*tenmod.Domain, error)

	CreateTenantIngestionProfile(*tenmod.IngestionProfile) (*tenmod.IngestionProfile, error)
	UpdateTenantIngestionProfile(*tenmod.IngestionProfile) (*tenmod.IngestionProfile, error)
	GetTenantIngestionProfile(tenantID string, dataID string) (*tenmod.IngestionProfile, error)
	DeleteTenantIngestionProfile(tenantID string, dataID string) (*tenmod.IngestionProfile, error)
	GetActiveTenantIngestionProfile(tenantID string) (*tenmod.IngestionProfile, error)

	CreateTenantThresholdProfile(*tenmod.ThresholdProfile) (*tenmod.ThresholdProfile, error)
	UpdateTenantThresholdProfile(*tenmod.ThresholdProfile) (*tenmod.ThresholdProfile, error)
	GetTenantThresholdProfile(tenantID string, dataID string) (*tenmod.ThresholdProfile, error)
	DeleteTenantThresholdProfile(tenantID string, dataID string) (*tenmod.ThresholdProfile, error)
	GetAllTenantThresholdProfile(tenantID string) ([]*tenmod.ThresholdProfile, error)

	CreateMonitoredObject(monitoredObjectReq *tenmod.MonitoredObject) (*tenmod.MonitoredObject, error)
	UpdateMonitoredObject(monitoredObjectReq *tenmod.MonitoredObject) (*tenmod.MonitoredObject, error)
	GetMonitoredObject(tenantID string, dataID string) (*tenmod.MonitoredObject, error)
	DeleteMonitoredObject(tenantID string, dataID string) (*tenmod.MonitoredObject, error)
	GetAllMonitoredObjects(tenantID string) ([]*tenmod.MonitoredObject, error)
	GetMonitoredObjectToDomainMap(moByDomReq *tenmod.MonitoredObjectCountByDomainRequest) (*tenmod.MonitoredObjectCountByDomainResponse, error)
	BulkInsertMonitoredObjects(tenantID string, value []*tenmod.MonitoredObject) ([]*common.BulkOperationResult, error)

	CreateTenantMeta(meta *tenmod.Metadata) (*tenmod.Metadata, error)
	UpdateTenantMeta(meta *tenmod.Metadata) (*tenmod.Metadata, error)
	DeleteTenantMeta(tenantID string) (*tenmod.Metadata, error)
	GetTenantMeta(tenantID string) (*tenmod.Metadata, error)
}
