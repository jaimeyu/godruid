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

	// ValidTypesType - datatype string used to identify a ValidTypes object in the datastore record
	ValidTypesType AdminDataType = "validTypes"
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

	// ValidTypesStr - common name of the ValidTypes data type for use in logs.
	ValidTypesStr = "Valid Types object"
)

// AdminServiceDatastore - interface which provides the functionality
// of the AdminService Datastore.
type AdminServiceDatastore interface {
	AddAdminViews() error
	CreateDatabase(dbName string) (Database, error)

	CreateAdminUser(*pb.AdminUser) (*pb.AdminUser, error)
	UpdateAdminUser(*pb.AdminUser) (*pb.AdminUser, error)
	DeleteAdminUser(string) (*pb.AdminUser, error)
	GetAdminUser(string) (*pb.AdminUser, error)
	GetAllAdminUsers() (*pb.AdminUserList, error)

	CreateTenant(*pb.TenantDescriptor) (*pb.TenantDescriptor, error)
	UpdateTenantDescriptor(*pb.TenantDescriptor) (*pb.TenantDescriptor, error)
	DeleteTenant(string) (*pb.TenantDescriptor, error)
	GetTenantDescriptor(string) (*pb.TenantDescriptor, error)
	GetAllTenantDescriptors() (*pb.TenantDescriptorList, error)

	CreateIngestionDictionary(ingDictionary *pb.IngestionDictionary) (*pb.IngestionDictionary, error)
	UpdateIngestionDictionary(ingDictionary *pb.IngestionDictionary) (*pb.IngestionDictionary, error)
	DeleteIngestionDictionary() (*pb.IngestionDictionary, error)
	GetIngestionDictionary() (*pb.IngestionDictionary, error)

	GetTenantIDByAlias(name string) (string, error)

	CreateValidTypes(value *pb.ValidTypes) (*pb.ValidTypes, error)
	UpdateValidTypes(value *pb.ValidTypes) (*pb.ValidTypes, error)
	GetValidTypes() (*pb.ValidTypes, error)
	GetSpecificValidTypes(value *pb.ValidTypesRequest) (*pb.ValidTypesData, error)
	DeleteValidTypes() (*pb.ValidTypes, error)
}
