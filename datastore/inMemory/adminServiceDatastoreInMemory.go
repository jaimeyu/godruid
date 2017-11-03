package inMemory

import (
	pb "github.com/accedian/adh-gather/gathergrpc"
)

// AdminServiceDatastoreInMemory - struct responsible for handling
// database operations for the Admin Service when using local memory
// as the storage option. Useful for tests.
type AdminServiceDatastoreInMemory struct {
}

// CreateAdminUser - InMemory implementation of CreateAdminUser
func (memDB *AdminServiceDatastoreInMemory) CreateAdminUser(*pb.AdminUser) (*pb.AdminUser, error) {
	// Stub to implement
	return nil, nil
}

// UpdateAdminUser - InMemory implementation of UpdateAdminUser
func (memDB *AdminServiceDatastoreInMemory) UpdateAdminUser(*pb.AdminUser) (*pb.AdminUser, error) {
	// Stub to implement
	return nil, nil
}

// DeleteAdminUser - InMemory implementation of DeleteAdminUser
func (memDB *AdminServiceDatastoreInMemory) DeleteAdminUser(string) (*pb.AdminUser, error) {
	// Stub to implement
	return nil, nil
}

// GetAdminUser - InMemory implementation of GetAdminUser
func (memDB *AdminServiceDatastoreInMemory) GetAdminUser(string) (*pb.AdminUser, error) {
	// Stub to implement
	return nil, nil
}

// GetAllAdminUsers - InMemory implementation of GetAllAdminUsers
func (memDB *AdminServiceDatastoreInMemory) GetAllAdminUsers() (*pb.AdminUserList, error) {
	// Stub to implement
	return nil, nil
}

// CreateTenant - InMemory implementation of CreateTenant
func (memDB *AdminServiceDatastoreInMemory) CreateTenant(*pb.TenantDescriptor) (*pb.TenantDescriptor, error) {
	// Stub to implement
	return nil, nil
}

// UpdateTenantDescriptor - InMemory implementation of UpdateTenantDescriptor
func (memDB *AdminServiceDatastoreInMemory) UpdateTenantDescriptor(*pb.TenantDescriptor) (*pb.TenantDescriptor, error) {
	// Stub to implement
	return nil, nil
}

// DeleteTenant - InMemory implementation of DeleteTenant
func (memDB *AdminServiceDatastoreInMemory) DeleteTenant(string) (*pb.TenantDescriptor, error) {
	// Stub to implement
	return nil, nil
}

// GetTenantDescriptor - InMemory implementation of GetTenantDescriptor
func (memDB *AdminServiceDatastoreInMemory) GetTenantDescriptor(string) (*pb.TenantDescriptor, error) {
	// Stub to implement
	return nil, nil
}
