package inMemory

import (
	"errors"

	"github.com/accedian/adh-gather/datastore"
	pb "github.com/accedian/adh-gather/gathergrpc"
)

// TenantServiceDatastoreInMemory - struct responsible for handling
// database operations for the Admin Service when using local memory
// as the storage option. Useful for tests.
type TenantServiceDatastoreInMemory struct {
}

// CreateTenantServiceDAO - returns an in-memory implementation of the Tenant Service
// datastore.
func CreateTenantServiceDAO() datastore.TenantServiceDatastore {
	res := new(TenantServiceDatastoreInMemory)

	return res
}

// CreateTenantUser - InMemory implementation of CreateTenantUser
func (tsd *TenantServiceDatastoreInMemory) CreateTenantUser(tenantUserRequest *pb.TenantUserRequest) (*pb.TenantUser, error) {
	// Stub to implement
	return nil, errors.New("Unsupported operation: CreateTenantUser not implemented")
}

// UpdateTenantUser - InMemory implementation of UpdateTenantUser
func (tsd *TenantServiceDatastoreInMemory) UpdateTenantUser(tenantUserRequest *pb.TenantUserRequest) (*pb.TenantUser, error) {
	// Stub to implement
	return nil, errors.New("Unsupported operation: UpdateTenantUser not implemented")
}

// DeleteTenantUser - InMemory implementation of DeleteTenantUser
func (tsd *TenantServiceDatastoreInMemory) DeleteTenantUser(tenantUserIDRequest *pb.TenantUserIdRequest) (*pb.TenantUser, error) {
	// Stub to implement
	return nil, errors.New("Unsupported operation: DeleteTenantUser not implemented")
}

// GetTenantUser - InMemory implementation of GetTenantUser
func (tsd *TenantServiceDatastoreInMemory) GetTenantUser(tenantUserIDRequest *pb.TenantUserIdRequest) (*pb.TenantUser, error) {
	// Stub to implement
	return nil, errors.New("Unsupported operation: GetTenantUser not implemented")
}

// GetAllTenantUsers - InMemory implementation of GetAllTenantUsers
func (tsd *TenantServiceDatastoreInMemory) GetAllTenantUsers(tenantID string) (*pb.TenantUserList, error) {
	// Stub to implement
	return nil, errors.New("Unsupported operation: GetAllTenantUsers not implemented")
}

// CreateTenantDomain - InMemory implementation of CreateTenantDomain
func (tsd *TenantServiceDatastoreInMemory) CreateTenantDomain(tenantDomainRequest *pb.TenantDomainRequest) (*pb.TenantDomain, error) {
	// Stub to implement
	return nil, errors.New("Unsupported operation: CreateTenantDomain not implemented")
}

// UpdateTenantDomain - InMemory implementation of UpdateTenantDomain
func (tsd *TenantServiceDatastoreInMemory) UpdateTenantDomain(tenantDomainRequest *pb.TenantDomainRequest) (*pb.TenantDomain, error) {
	// Stub to implement
	return nil, errors.New("Unsupported operation: UpdateTenantDomain not implemented")
}

// DeleteTenantDomain - InMemory implementation of DeleteTenantDomain
func (tsd *TenantServiceDatastoreInMemory) DeleteTenantDomain(tenantDomainIDRequest *pb.TenantDomainIdRequest) (*pb.TenantDomain, error) {
	// Stub to implement
	return nil, errors.New("Unsupported operation: DeleteTenantDomain not implemented")
}

// GetTenantDomain - InMemory implementation of GetTenantDomain
func (tsd *TenantServiceDatastoreInMemory) GetTenantDomain(tenantDomainIDRequest *pb.TenantDomainIdRequest) (*pb.TenantDomain, error) {
	// Stub to implement
	return nil, errors.New("Unsupported operation: GetTenantDomain not implemented")
}

// GetAllTenantDomains - InMemory implementation of GetAllTenantDomains
func (tsd *TenantServiceDatastoreInMemory) GetAllTenantDomains(tenantID string) (*pb.TenantDomainList, error) {
	// Stub to implement
	return nil, errors.New("Unsupported operation: GetAllTenantDomains not implemented")
}

// CreateTenantIngestionProfile - InMemory implementation of CreateTenantIngestionProfile
func (tsd *TenantServiceDatastoreInMemory) CreateTenantIngestionProfile(tenantIngPrfReq *pb.TenantIngestionProfileRequest) (*pb.TenantIngestionProfile, error) {
	// Stub to implement
	return nil, errors.New("Unsupported operation: CreateTenantIngestionProfile not implemented")
}

// UpdateTenantIngestionProfile - InMemory implementation of UpdateTenantIngestionProfile
func (tsd *TenantServiceDatastoreInMemory) UpdateTenantIngestionProfile(tenantIngPrfReq *pb.TenantIngestionProfileRequest) (*pb.TenantIngestionProfile, error) {
	// Stub to implement
	return nil, errors.New("Unsupported operation: UpdateTenantIngestionProfile not implemented")
}

// GetTenantIngestionProfile - InMemory implementation of GetTenantIngestionProfile
func (tsd *TenantServiceDatastoreInMemory) GetTenantIngestionProfile(tenantIngPrfIDReq *pb.TenantIngestionProfileIdRequest) (*pb.TenantIngestionProfile, error) {
	// Stub to implement
	return nil, errors.New("Unsupported operation: GetTenantIngestionProfile not implemented")
}

// DeleteTenantIngestionProfile - InMemory implementation of DeleteTenantIngestionProfile
func (tsd *TenantServiceDatastoreInMemory) DeleteTenantIngestionProfile(tenantIngPrfIDReq *pb.TenantIngestionProfileIdRequest) (*pb.TenantIngestionProfile, error) {
	// Stub to implement
	return nil, errors.New("Unsupported operation: DeleteTenantIngestionProfile not implemented")
}
