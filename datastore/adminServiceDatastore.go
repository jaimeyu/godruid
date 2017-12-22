package datastore

import (
	pb "github.com/accedian/adh-gather/gathergrpc"
)

// AdminDataType - data type descriptors for objects stored in the admin datastore
type AdminDataType string

const (
	// AdminUserType - datatype string used to identify an Admin User in the datastore record
	AdminUserType AdminDataType = "adminUser"

	// TenantDescriptorType - datatype string used to identify a Tenant Descriptor in the datastore record
	TenantDescriptorType AdminDataType = "tenant"

	// IngestionDictionaryType - datatype string used to identify an IngestionDictionary in the datastore record
	IngestionDictionaryType AdminDataType = "ingestionDictionary"
)

const (
	// AdminUserStr - common name of the AdminUser data type for use in logs.
	AdminUserStr = "Admin User"

	// TenantDescriptorStr - common name of the TenantDescriptor data type for use in logs.
	TenantDescriptorStr = "Tenant Descriptor"

	// TenantStr - common name of the TenantDescriptor data type for use in logs.
	TenantStr = "Tenant"

	// IngestionDictionaryStr - common name of the IngestionDictionary data type for use in logs.
	IngestionDictionaryStr = "Ingestion Dictionary"
)

// AdminServiceDatastore - interface which provides the functionality
// of the AdminService Datastore.
type AdminServiceDatastore interface {
	AddAdminViews() error

	CreateAdminUser(*pb.AdminUserRequest) (*pb.AdminUserResponse, error)
	UpdateAdminUser(*pb.AdminUserRequest) (*pb.AdminUserResponse, error)
	DeleteAdminUser(string) (*pb.AdminUserResponse, error)
	GetAdminUser(string) (*pb.AdminUserResponse, error)
	GetAllAdminUsers() (*pb.AdminUserListResponse, error)

	CreateTenant(*pb.TenantDescriptorRequest) (*pb.TenantDescriptorResponse, error)
	UpdateTenantDescriptor(*pb.TenantDescriptorRequest) (*pb.TenantDescriptorResponse, error)
	DeleteTenant(string) (*pb.TenantDescriptorResponse, error)
	GetTenantDescriptor(string) (*pb.TenantDescriptorResponse, error)
	GetAllTenantDescriptors() (*pb.TenantDescriptorListResponse, error)

	CreateIngestionDictionary(ingDictionary *pb.IngestionDictionary) (*pb.IngestionDictionary, error)
	UpdateIngestionDictionary(ingDictionary *pb.IngestionDictionary) (*pb.IngestionDictionary, error)
	DeleteIngestionDictionary() (*pb.IngestionDictionary, error)
	GetIngestionDictionary() (*pb.IngestionDictionary, error)

	GetTenantIDByAlias(name string) (string, error)
}
