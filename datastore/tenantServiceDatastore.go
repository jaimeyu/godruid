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

// TenantServiceDatastore - interface which provides the functionality
// of the TenantService Datastore.
type TenantServiceDatastore interface {
	CreateTenantUser(*pb.TenantUserRequest) (*pb.TenantUser, error)
	UpdateTenantUser(*pb.TenantUserRequest) (*pb.TenantUser, error)
	DeleteTenantUser(*pb.TenantUserIdRequest) (*pb.TenantUser, error)
	GetTenantUser(*pb.TenantUserIdRequest) (*pb.TenantUser, error)
	GetAllTenantUsers(string) ([]*pb.TenantUser, error)

	CreateTenantDomain(*pb.TenantDomainRequest) (*pb.TenantDomain, error)
	UpdateTenantDomain(*pb.TenantDomainRequest) (*pb.TenantDomain, error)
	DeleteTenantDomain(*pb.TenantDomainIdRequest) (*pb.TenantDomain, error)
	GetTenantDomain(*pb.TenantDomainIdRequest) (*pb.TenantDomain, error)
	GetAllTenantDomains(string) ([]*pb.TenantDomain, error)

	CreateTenantIngestionProfile(*pb.TenantIngestionProfileRequest) (*pb.TenantIngestionProfile, error)
	UpdateTenantIngestionProfile(*pb.TenantIngestionProfileRequest) (*pb.TenantIngestionProfile, error)
	GetTenantIngestionProfile(*pb.TenantIngestionProfileIdRequest) (*pb.TenantIngestionProfile, error)
	DeleteTenantIngestionProfile(*pb.TenantIngestionProfileIdRequest) (*pb.TenantIngestionProfile, error)
}
