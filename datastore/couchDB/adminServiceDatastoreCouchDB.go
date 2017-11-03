package couchDB

import (
	pb "github.com/accedian/adh-gather/gathergrpc"
)

// AdminServiceDatastoreCouchDB - struct responsible for handling
// database operations for the Admin Service when using CouchDB
// as the storage option.
type AdminServiceDatastoreCouchDB struct {
}

// CreateAdminUser - CouchDB implementation of CreateAdminUser
func (couchDB *AdminServiceDatastoreCouchDB) CreateAdminUser(*pb.AdminUser) (*pb.AdminUser, error) {
	// Stub to implement
	return nil, nil
}

// UpdateAdminUser - CouchDB implementation of UpdateAdminUser
func (couchDB *AdminServiceDatastoreCouchDB) UpdateAdminUser(*pb.AdminUser) (*pb.AdminUser, error) {
	// Stub to implement
	return nil, nil
}

// DeleteAdminUser - CouchDB implementation of DeleteAdminUser
func (couchDB *AdminServiceDatastoreCouchDB) DeleteAdminUser(string) (*pb.AdminUser, error) {
	// Stub to implement
	return nil, nil
}

// GetAdminUser - CouchDB implementation of GetAdminUser
func (couchDB *AdminServiceDatastoreCouchDB) GetAdminUser(string) (*pb.AdminUser, error) {
	// Stub to implement
	return nil, nil
}

// GetAllAdminUsers - CouchDB implementation of GetAllAdminUsers
func (couchDB *AdminServiceDatastoreCouchDB) GetAllAdminUsers() (*pb.AdminUserList, error) {
	// Stub to implement
	return nil, nil
}

// CreateTenant - CouchDB implementation of CreateTenant
func (couchDB *AdminServiceDatastoreCouchDB) CreateTenant(*pb.TenantDescriptor) (*pb.TenantDescriptor, error) {
	// Stub to implement
	return nil, nil
}

// UpdateTenantDescriptor - CouchDB implementation of UpdateTenantDescriptor
func (couchDB *AdminServiceDatastoreCouchDB) UpdateTenantDescriptor(*pb.TenantDescriptor) (*pb.TenantDescriptor, error) {
	// Stub to implement
	return nil, nil
}

// DeleteTenant - CouchDB implementation of DeleteTenant
func (couchDB *AdminServiceDatastoreCouchDB) DeleteTenant(string) (*pb.TenantDescriptor, error) {
	// Stub to implement
	return nil, nil
}

// GetTenantDescriptor - CouchDB implementation of GetTenantDescriptor
func (couchDB *AdminServiceDatastoreCouchDB) GetTenantDescriptor(string) (*pb.TenantDescriptor, error) {
	// Stub to implement
	return nil, nil
}
