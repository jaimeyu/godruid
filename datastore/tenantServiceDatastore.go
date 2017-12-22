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
	TenantMetaType TenantDataType = "meta"
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
	CreateTenantUser(*pb.TenantUserRequest) (*pb.TenantUserResponse, error)
	UpdateTenantUser(*pb.TenantUserRequest) (*pb.TenantUserResponse, error)
	DeleteTenantUser(*pb.TenantUserIdRequest) (*pb.TenantUserResponse, error)
	GetTenantUser(*pb.TenantUserIdRequest) (*pb.TenantUserResponse, error)
	GetAllTenantUsers(string) (*pb.TenantUserListResponse, error)

	CreateTenantDomain(*pb.TenantDomainRequest) (*pb.TenantDomainResponse, error)
	UpdateTenantDomain(*pb.TenantDomainRequest) (*pb.TenantDomainResponse, error)
	DeleteTenantDomain(*pb.TenantDomainIdRequest) (*pb.TenantDomainResponse, error)
	GetTenantDomain(*pb.TenantDomainIdRequest) (*pb.TenantDomainResponse, error)
	GetAllTenantDomains(string) (*pb.TenantDomainListResponse, error)

	CreateTenantIngestionProfile(*pb.TenantIngestionProfileRequest) (*pb.TenantIngestionProfileResponse, error)
	UpdateTenantIngestionProfile(*pb.TenantIngestionProfileRequest) (*pb.TenantIngestionProfileResponse, error)
	GetTenantIngestionProfile(*pb.TenantIngestionProfileIdRequest) (*pb.TenantIngestionProfileResponse, error)
	DeleteTenantIngestionProfile(*pb.TenantIngestionProfileIdRequest) (*pb.TenantIngestionProfileResponse, error)

	CreateTenantThresholdProfile(*pb.TenantThresholdProfileRequest) (*pb.TenantThresholdProfileResponse, error)
	UpdateTenantThresholdProfile(*pb.TenantThresholdProfileRequest) (*pb.TenantThresholdProfileResponse, error)
	GetTenantThresholdProfile(*pb.TenantThresholdProfileIdRequest) (*pb.TenantThresholdProfileResponse, error)
	DeleteTenantThresholdProfile(*pb.TenantThresholdProfileIdRequest) (*pb.TenantThresholdProfileResponse, error)

	CreateMonitoredObject(monitoredObjectReq *pb.MonitoredObjectRequest) (*pb.MonitoredObjectResponse, error)
	UpdateMonitoredObject(monitoredObjectReq *pb.MonitoredObjectRequest) (*pb.MonitoredObjectResponse, error)
	GetMonitoredObject(monitoredObjectIDReq *pb.MonitoredObjectIdRequest) (*pb.MonitoredObjectResponse, error)
	DeleteMonitoredObject(monitoredObjectIDReq *pb.MonitoredObjectIdRequest) (*pb.MonitoredObjectResponse, error)
	GetAllMonitoredObjects(tenantID string) (*pb.MonitoredObjectListResponse, error)
	GetMonitoredObjectToDomainMap(moByDomReq *pb.MonitoredObjectCountByDomainRequest) (*pb.MonitoredObjectCountByDomainResponse, error)

	CreateTenantMeta(meta *pb.TenantMetadata) (*pb.TenantMetadata, error)
	UpdateTenantMeta(meta *pb.TenantMetadata) (*pb.TenantMetadata, error)
	DeleteTenantMeta(tenantID string) (*pb.TenantMetadata, error)
	GetTenantMeta(tenantID string) (*pb.TenantMetadata, error)
}
