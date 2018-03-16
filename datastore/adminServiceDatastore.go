package datastore

import (
	admmod "github.com/accedian/adh-gather/models/admin"
)

// AdminServiceDatastore - interface which provides the functionality
// of the AdminService Datastore.
type AdminServiceDatastore interface {
	AddAdminViews() error
	CreateDatabase(dbName string) (Database, error)

	CreateAdminUser(*admmod.User) (*admmod.User, error)
	UpdateAdminUser(*admmod.User) (*admmod.User, error)
	DeleteAdminUser(string) (*admmod.User, error)
	GetAdminUser(string) (*admmod.User, error)
	GetAllAdminUsers() ([]*admmod.User, error)

	CreateTenant(*admmod.Tenant) (*admmod.Tenant, error)
	UpdateTenantDescriptor(*admmod.Tenant) (*admmod.Tenant, error)
	DeleteTenant(string) (*admmod.Tenant, error)
	GetTenantDescriptor(string) (*admmod.Tenant, error)
	GetAllTenantDescriptors() ([]*admmod.Tenant, error)

	CreateIngestionDictionary(ingDictionary *admmod.IngestionDictionary) (*admmod.IngestionDictionary, error)
	UpdateIngestionDictionary(ingDictionary *admmod.IngestionDictionary) (*admmod.IngestionDictionary, error)
	DeleteIngestionDictionary() (*admmod.IngestionDictionary, error)
	GetIngestionDictionary() (*admmod.IngestionDictionary, error)

	GetTenantIDByAlias(name string) (string, error)

	CreateValidTypes(value *admmod.ValidTypes) (*admmod.ValidTypes, error)
	UpdateValidTypes(value *admmod.ValidTypes) (*admmod.ValidTypes, error)
	GetValidTypes() (*admmod.ValidTypes, error)
	GetSpecificValidTypes(value *admmod.ValidTypesRequest) (*admmod.ValidTypes, error)
	DeleteValidTypes() (*admmod.ValidTypes, error)
}
