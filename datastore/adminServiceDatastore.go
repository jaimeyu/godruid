package datastore

import (
	pb "github.com/accedian/adh-gather/gathergrpc"
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
