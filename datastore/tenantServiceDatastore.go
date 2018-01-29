package datastore

import (
	pb "github.com/accedian/adh-gather/gathergrpc"
)

// TenantDataType - enumeration of the types of data stored in the Tenant Datastore
type TenantDataType string

const (
	// TenantUserType - datatype string used to identify a Tenant User in the datastore record
	TenantUserType TenantDataType = "user"

	// TenantDomainType - datatype string used to identify a Tenant Domain in the datastore record
	TenantDomainType TenantDataType = "domain"

	// TenantIngestionProfileType - datatype string used to identify a Tenant Ingestion Profile in the datastore record
	TenantIngestionProfileType TenantDataType = "ingestionProfile"

	// TenantMonitoredObjectType - datatype string used to identify a Tenant MonitoredObject in the datastore record
	TenantMonitoredObjectType TenantDataType = "monitoredObject"

	// TenantThresholdProfileType - datatype string used to identify a Tenant Ingestion Profile in the datastore record
	TenantThresholdProfileType TenantDataType = "thresholdProfile"

	// TenantMetaType - datatype string used to identify a Tenant Meta in the datastore record
	TenantMetaType TenantDataType = "tenantMetadata"
)

const (
	// TenantUserStr - common name of the TenantUser data type for use in logs.
	TenantUserStr = "Tenant User"

	// TenantDomainStr - common name of the Tenant Domain data type for use in logs.
	TenantDomainStr = "Tenant Domain"

	// TenantIngestionProfileStr - common name of the Tenant Ingestion Profile data type for use in logs.
	TenantIngestionProfileStr = "Tenant Ingestion Profile"

	// TenantMonitoredObjectStr - common name of the Tenant Monitored Object data type for use in logs.
	TenantMonitoredObjectStr = "Tenant Monitored Object"

	// TenantThresholdProfileStr - common name of the Tenant Ingestion Profile data type for use in logs.
	TenantThresholdProfileStr = "Tenant Threshold Profile"

	// MonitoredObjectToDomainMapStr - common name for the Monitored Object to Doamin Map for use in logs.
	MonitoredObjectToDomainMapStr = "Monitored Object to Doamin Map"

	// TenantMetaStr - common name for the Meta for use in logs.
	TenantMetaStr = "Tenant Meta"
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
