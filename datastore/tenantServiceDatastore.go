package datastore

import (
	pb "github.com/accedian/adh-gather/gathergrpc"
)

// TenantUserStr - common name of the TenantUser data type for use in logs.
const TenantUserStr string = "Tenant User"

// TenantDomainStr - common name of the Tenant Domain data type for use in logs.
const TenantDomainStr string = "Tenant Domain"

// TenantIngestionProfileStr - common name of the Tenant Ingestion Profile data type for use in logs.
const TenantIngestionProfileStr string = "Tenant Ingestion Profile"

// TenantUserType - datatype string used to identify a Tenant User in the datastore record
const TenantUserType string = "user"

// TenantDomainType - datatype string used to identify a Tenant Domain in the datastore record
const TenantDomainType string = "domain"

// TenantIngPrfType - datatype string used to identify a Tenant Ingestion Profile in the datastore record
const TenantIngPrfType string = "ingestionProfile"

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
}
