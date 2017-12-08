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
func CreateTenantServiceDAO() (datastore.TenantServiceDatastore, error) {
	res := new(TenantServiceDatastoreInMemory)

	return res, nil
}

// CreateTenantUser - InMemory implementation of CreateTenantUser
func (tsd *TenantServiceDatastoreInMemory) CreateTenantUser(tenantUserRequest *pb.TenantUserRequest) (*pb.TenantUserResponse, error) {
	// Stub to implement
	return nil, errors.New("Unsupported operation: CreateTenantUser not implemented")
}

// UpdateTenantUser - InMemory implementation of UpdateTenantUser
func (tsd *TenantServiceDatastoreInMemory) UpdateTenantUser(tenantUserRequest *pb.TenantUserRequest) (*pb.TenantUserResponse, error) {
	// Stub to implement
	return nil, errors.New("Unsupported operation: UpdateTenantUser not implemented")
}

// DeleteTenantUser - InMemory implementation of DeleteTenantUser
func (tsd *TenantServiceDatastoreInMemory) DeleteTenantUser(tenantUserIDRequest *pb.TenantUserIdRequest) (*pb.TenantUserResponse, error) {
	// Stub to implement
	return nil, errors.New("Unsupported operation: DeleteTenantUser not implemented")
}

// GetTenantUser - InMemory implementation of GetTenantUser
func (tsd *TenantServiceDatastoreInMemory) GetTenantUser(tenantUserIDRequest *pb.TenantUserIdRequest) (*pb.TenantUserResponse, error) {
	// Stub to implement
	return nil, errors.New("Unsupported operation: GetTenantUser not implemented")
}

// GetAllTenantUsers - InMemory implementation of GetAllTenantUsers
func (tsd *TenantServiceDatastoreInMemory) GetAllTenantUsers(tenantID string) (*pb.TenantUserListResponse, error) {
	// Stub to implement
	return nil, errors.New("Unsupported operation: GetAllTenantUsers not implemented")
}

// CreateTenantDomain - InMemory implementation of CreateTenantDomain
func (tsd *TenantServiceDatastoreInMemory) CreateTenantDomain(tenantDomainRequest *pb.TenantDomainRequest) (*pb.TenantDomainResponse, error) {
	// Stub to implement
	return nil, errors.New("Unsupported operation: CreateTenantDomain not implemented")
}

// UpdateTenantDomain - InMemory implementation of UpdateTenantDomain
func (tsd *TenantServiceDatastoreInMemory) UpdateTenantDomain(tenantDomainRequest *pb.TenantDomainRequest) (*pb.TenantDomainResponse, error) {
	// Stub to implement
	return nil, errors.New("Unsupported operation: UpdateTenantDomain not implemented")
}

// DeleteTenantDomain - InMemory implementation of DeleteTenantDomain
func (tsd *TenantServiceDatastoreInMemory) DeleteTenantDomain(tenantDomainIDRequest *pb.TenantDomainIdRequest) (*pb.TenantDomainResponse, error) {
	// Stub to implement
	return nil, errors.New("Unsupported operation: DeleteTenantDomain not implemented")
}

// GetTenantDomain - InMemory implementation of GetTenantDomain
func (tsd *TenantServiceDatastoreInMemory) GetTenantDomain(tenantDomainIDRequest *pb.TenantDomainIdRequest) (*pb.TenantDomainResponse, error) {
	// Stub to implement
	return nil, errors.New("Unsupported operation: GetTenantDomain not implemented")
}

// GetAllTenantDomains - InMemory implementation of GetAllTenantDomains
func (tsd *TenantServiceDatastoreInMemory) GetAllTenantDomains(tenantID string) (*pb.TenantDomainListResponse, error) {
	// Stub to implement
	return nil, errors.New("Unsupported operation: GetAllTenantDomains not implemented")
}

// CreateTenantIngestionProfile - InMemory implementation of CreateTenantIngestionProfile
func (tsd *TenantServiceDatastoreInMemory) CreateTenantIngestionProfile(tenantIngPrfReq *pb.TenantIngestionProfileRequest) (*pb.TenantIngestionProfileResponse, error) {
	// Stub to implement
	return nil, errors.New("Unsupported operation: CreateTenantIngestionProfile not implemented")
}

// UpdateTenantIngestionProfile - InMemory implementation of UpdateTenantIngestionProfile
func (tsd *TenantServiceDatastoreInMemory) UpdateTenantIngestionProfile(tenantIngPrfReq *pb.TenantIngestionProfileRequest) (*pb.TenantIngestionProfileResponse, error) {
	// Stub to implement
	return nil, errors.New("Unsupported operation: UpdateTenantIngestionProfile not implemented")
}

// GetTenantIngestionProfile - InMemory implementation of GetTenantIngestionProfile
func (tsd *TenantServiceDatastoreInMemory) GetTenantIngestionProfile(tenantIngPrfIDReq *pb.TenantIngestionProfileIdRequest) (*pb.TenantIngestionProfileResponse, error) {
	// Stub to implement
	return nil, errors.New("Unsupported operation: GetTenantIngestionProfile not implemented")
}

// DeleteTenantIngestionProfile - InMemory implementation of DeleteTenantIngestionProfile
func (tsd *TenantServiceDatastoreInMemory) DeleteTenantIngestionProfile(tenantIngPrfIDReq *pb.TenantIngestionProfileIdRequest) (*pb.TenantIngestionProfileResponse, error) {
	// Stub to implement
	return nil, errors.New("Unsupported operation: DeleteTenantIngestionProfile not implemented")
}

// CreateTenantThresholdProfile - InMemory implementation of CreateTenantThresholdProfile
func (tsd *TenantServiceDatastoreInMemory) CreateTenantThresholdProfile(tenantThreshPrfReq *pb.TenantThresholdProfileRequest) (*pb.TenantThresholdProfileResponse, error) {
	// Stub to implement
	return nil, errors.New("Unsupported operation: CreateTenantThresholdProfile not implemented")
}

// UpdateTenantThresholdProfile - InMemory implementation of UpdateTenantThresholdProfile
func (tsd *TenantServiceDatastoreInMemory) UpdateTenantThresholdProfile(tenantThreshPrfReq *pb.TenantThresholdProfileRequest) (*pb.TenantThresholdProfileResponse, error) {
	// Stub to implement
	return nil, errors.New("Unsupported operation: UpdateTenantThresholdProfile not implemented")
}

// GetTenantThresholdProfile - InMemory implementation of GetTenantThresholdProfile
func (tsd *TenantServiceDatastoreInMemory) GetTenantThresholdProfile(tenantIngPrfIDReq *pb.TenantThresholdProfileIdRequest) (*pb.TenantThresholdProfileResponse, error) {
	// Stub to implement
	return nil, errors.New("Unsupported operation: GetTenantThresholdProfile not implemented")
}

// DeleteTenantThresholdProfile - InMemory implementation of DeleteTenantThresholdProfile
func (tsd *TenantServiceDatastoreInMemory) DeleteTenantThresholdProfile(tenantIngPrfIDReq *pb.TenantThresholdProfileIdRequest) (*pb.TenantThresholdProfileResponse, error) {
	// Stub to implement
	return nil, errors.New("Unsupported operation: DeleteTenantThresholdProfile not implemented")
}

// CreateMonitoredObject - InMemory implementation of CreateMonitoredObject
func (tsd *TenantServiceDatastoreInMemory) CreateMonitoredObject(monitoredObjectReq *pb.MonitoredObjectRequest) (*pb.MonitoredObjectResponse, error) {
	// Stub to implement
	return nil, errors.New("Unsupported operation: CreateMonitoredObject not implemented")
}

// UpdateMonitoredObject - InMemory implementation of UpdateMonitoredObject
func (tsd *TenantServiceDatastoreInMemory) UpdateMonitoredObject(monitoredObjectReq *pb.MonitoredObjectRequest) (*pb.MonitoredObjectResponse, error) {
	// Stub to implement
	return nil, errors.New("Unsupported operation: UpdateMonitoredObject not implemented")
}

// GetMonitoredObject - InMemory implementation of GetMonitoredObject
func (tsd *TenantServiceDatastoreInMemory) GetMonitoredObject(monitoredObjectIDReq *pb.MonitoredObjectIdRequest) (*pb.MonitoredObjectResponse, error) {
	// Stub to implement
	return nil, errors.New("Unsupported operation: GetMonitoredObject not implemented")
}

// DeleteMonitoredObject - InMemory implementation of DeleteMonitoredObject
func (tsd *TenantServiceDatastoreInMemory) DeleteMonitoredObject(monitoredObjectIDReq *pb.MonitoredObjectIdRequest) (*pb.MonitoredObjectResponse, error) {
	// Stub to implement
	return nil, errors.New("Unsupported operation: DeleteMonitoredObject not implemented")
}

// GetAllMonitoredObjects - InMemory implementation of GetAllMonitoredObjects
func (tsd *TenantServiceDatastoreInMemory) GetAllMonitoredObjects(tenantID string) (*pb.MonitoredObjectListResponse, error) {
	// Stub to implement
	return nil, errors.New("Unsupported operation: GetAllMonitoredObjects not implemented")
}
