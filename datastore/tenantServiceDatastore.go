package datastore

import (
	pb "github.com/accedian/adh-gather/gathergrpc"
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

	CreateMonitoredObject(monitoredObjectReq *pb.MonitoredObject) (*pb.MonitoredObject, error)
	UpdateMonitoredObject(monitoredObjectReq *pb.MonitoredObject) (*pb.MonitoredObject, error)
	GetMonitoredObject(monitoredObjectIDReq *pb.MonitoredObjectIdRequest) (*pb.MonitoredObject, error)
	DeleteMonitoredObject(monitoredObjectIDReq *pb.MonitoredObjectIdRequest) (*pb.MonitoredObject, error)
	GetAllMonitoredObjects(tenantID string) (*pb.MonitoredObjectList, error)
	GetMonitoredObjectToDomainMap(moByDomReq *pb.MonitoredObjectCountByDomainRequest) (*pb.MonitoredObjectCountByDomainResponse, error)

	CreateTenantMeta(meta *pb.TenantMetadata) (*pb.TenantMetadata, error)
	UpdateTenantMeta(meta *pb.TenantMetadata) (*pb.TenantMetadata, error)
	DeleteTenantMeta(tenantID string) (*pb.TenantMetadata, error)
	GetTenantMeta(tenantID string) (*pb.TenantMetadata, error)

	BulkInsertMonitoredObjects(value *pb.TenantMonitoredObjectSet) (*pb.BulkOperationResponse, error)
}
