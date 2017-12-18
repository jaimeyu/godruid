package inMemory

import (
	"errors"

	"github.com/accedian/adh-gather/datastore"
	pb "github.com/accedian/adh-gather/gathergrpc"
)

// AdminServiceDatastoreInMemory - struct responsible for handling
// database operations for the Admin Service when using local memory
// as the storage option. Useful for tests.
type AdminServiceDatastoreInMemory struct {
}

// CreateAdminServiceDAO - returns an in-memory implementation of the Admin Service
// datastore.
func CreateAdminServiceDAO() (datastore.AdminServiceDatastore, error) {
	res := new(AdminServiceDatastoreInMemory)

	return res, nil
}

// CreateAdminUser - InMemory implementation of CreateAdminUser
func (memDB *AdminServiceDatastoreInMemory) CreateAdminUser(user *pb.AdminUserRequest) (*pb.AdminUserResponse, error) {
	// Stub to implement
	return nil, errors.New("CreateAdminUser() not implemented for InMemory DB")
}

// UpdateAdminUser - InMemory implementation of UpdateAdminUser
func (memDB *AdminServiceDatastoreInMemory) UpdateAdminUser(user *pb.AdminUserRequest) (*pb.AdminUserResponse, error) {
	// Stub to implement
	return nil, errors.New("UpdateAdminUser() not implemented for InMemory DB")
}

// DeleteAdminUser - InMemory implementation of DeleteAdminUser
func (memDB *AdminServiceDatastoreInMemory) DeleteAdminUser(userID string) (*pb.AdminUserResponse, error) {
	// Stub to implement
	return nil, errors.New("DeleteAdminUser() not implemented for InMemory DB")
}

// GetAdminUser - InMemory implementation of GetAdminUser
func (memDB *AdminServiceDatastoreInMemory) GetAdminUser(userID string) (*pb.AdminUserResponse, error) {
	// Stub to implement
	return nil, errors.New("GetAdminUser() not implemented for InMemory DB")
}

// GetAllAdminUsers - InMemory implementation of GetAllAdminUsers
func (memDB *AdminServiceDatastoreInMemory) GetAllAdminUsers() (*pb.AdminUserListResponse, error) {
	// Stub to implement
	return nil, errors.New("GetAllAdminUsers() not implemented for InMemory DB")
}

// CreateTenant - InMemory implementation of CreateTenant
func (memDB *AdminServiceDatastoreInMemory) CreateTenant(tenantDescriptor *pb.TenantDescriptorRequest) (*pb.TenantDescriptorResponse, error) {
	// Stub to implement
	return nil, errors.New("CreateTenant() not implemented for InMemory DB")
}

// UpdateTenantDescriptor - InMemory implementation of UpdateTenantDescriptor
func (memDB *AdminServiceDatastoreInMemory) UpdateTenantDescriptor(tenantDescriptor *pb.TenantDescriptorRequest) (*pb.TenantDescriptorResponse, error) {
	// Stub to implement
	return nil, errors.New("UpdateTenantDescriptor() not implemented for InMemory DB")
}

// DeleteTenant - InMemory implementation of DeleteTenant
func (memDB *AdminServiceDatastoreInMemory) DeleteTenant(tenantID string) (*pb.TenantDescriptorResponse, error) {
	// Stub to implement
	return nil, errors.New("DeleteTenant() not implemented for InMemory DB")
}

// GetTenantDescriptor - InMemory implementation of GetTenantDescriptor
func (memDB *AdminServiceDatastoreInMemory) GetTenantDescriptor(tenantID string) (*pb.TenantDescriptorResponse, error) {
	// Stub to implement
	return nil, errors.New("GetTenantDescriptor() not implemented for InMemory DB")
}

// GetAllTenantDescriptors - InMemory implementation of GetAllTenantDescriptors
func (memDB *AdminServiceDatastoreInMemory) GetAllTenantDescriptors() (*pb.TenantDescriptorListResponse, error) {
	// Stub to implement
	return nil, errors.New("GetAllTenantDescriptors() not implemented for InMemory DB")
}

// CreateIngestionDictionary - InMemory implementation of CreateIngestionDictionary
func (memDB *AdminServiceDatastoreInMemory) CreateIngestionDictionary(ingDictionary *pb.IngestionDictionary) (*pb.IngestionDictionary, error) {
	// Stub to implement
	return nil, errors.New("CreateIngestionDictionary() not implemented for InMemory DB")
}

// UpdateIngestionDictionary - InMemory implementation of UpdateIngestionDictionary
func (memDB *AdminServiceDatastoreInMemory) UpdateIngestionDictionary(ingDictionary *pb.IngestionDictionary) (*pb.IngestionDictionary, error) {
	// Stub to implement
	return nil, errors.New("UpdateIngestionDictionary() not implemented for InMemory DB")
}

// DeleteIngestionDictionary - InMemory implementation of DeleteIngestionDictionary
func (memDB *AdminServiceDatastoreInMemory) DeleteIngestionDictionary() (*pb.IngestionDictionary, error) {
	// Stub to implement
	return nil, errors.New("DeleteIngestionDictionary() not implemented for InMemory DB")
}

// GetIngestionDictionary - InMemory implementation of GetIngestionDictionary
func (memDB *AdminServiceDatastoreInMemory) GetIngestionDictionary() (*pb.IngestionDictionary, error) {
	// Stub to implement
	return nil, errors.New("GetIngestionDictionary() not implemented for InMemory DB")
}
