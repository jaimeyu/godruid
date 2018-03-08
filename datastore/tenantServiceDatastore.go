package datastore

import (
	pb "github.com/accedian/adh-gather/gathergrpc"
)

// TenantServiceDatastore - interface which provides the functionality
// of the TenantService Datastore.
type TenantServiceDatastore interface {
	CreateTenantUser(*pb.TenantUser) (*pb.TenantUser, error)
	UpdateTenantUser(*pb.TenantUser) (*pb.TenantUser, error)
	DeleteTenantUser(*pb.TenantUserIdRequest) (*pb.TenantUser, error)
	GetTenantUser(*pb.TenantUserIdRequest) (*pb.TenantUser, error)
	GetAllTenantUsers(string) (*pb.TenantUserList, error)

	CreateTenantDomain(*pb.TenantDomain) (*pb.TenantDomain, error)
	UpdateTenantDomain(*pb.TenantDomain) (*pb.TenantDomain, error)
	DeleteTenantDomain(*pb.TenantDomainIdRequest) (*pb.TenantDomain, error)
	GetTenantDomain(*pb.TenantDomainIdRequest) (*pb.TenantDomain, error)
	GetAllTenantDomains(string) (*pb.TenantDomainList, error)

	CreateTenantIngestionProfile(*pb.TenantIngestionProfile) (*pb.TenantIngestionProfile, error)
	UpdateTenantIngestionProfile(*pb.TenantIngestionProfile) (*pb.TenantIngestionProfile, error)
	GetTenantIngestionProfile(*pb.TenantIngestionProfileIdRequest) (*pb.TenantIngestionProfile, error)
	DeleteTenantIngestionProfile(*pb.TenantIngestionProfileIdRequest) (*pb.TenantIngestionProfile, error)

	CreateTenantThresholdProfile(*pb.TenantThresholdProfile) (*pb.TenantThresholdProfile, error)
	UpdateTenantThresholdProfile(*pb.TenantThresholdProfile) (*pb.TenantThresholdProfile, error)
	GetTenantThresholdProfile(*pb.TenantThresholdProfileIdRequest) (*pb.TenantThresholdProfile, error)
	DeleteTenantThresholdProfile(*pb.TenantThresholdProfileIdRequest) (*pb.TenantThresholdProfile, error)

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

	GetActiveTenantIngestionProfile(tenantID string) (*pb.TenantIngestionProfile, error)
	GetAllTenantThresholdProfile(tenantID string) (*pb.TenantThresholdProfileList, error)

	BulkInsertMonitoredObjects(value *pb.TenantMonitoredObjectSet) (*pb.BulkOperationResponse, error)
}
