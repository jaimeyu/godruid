package inMemory

import (
	"errors"
	"fmt"

	ds "github.com/accedian/adh-gather/datastore"
	pb "github.com/accedian/adh-gather/gathergrpc"
	"github.com/getlantern/deepcopy"
	uuid "github.com/satori/go.uuid"
)

// TenantServiceDatastoreInMemory - struct responsible for handling
// database operations for the Admin Service when using local memory
// as the storage option. Useful for tests.
type TenantServiceDatastoreInMemory struct {
	tenantToIDtoTenantUserMap map[string]map[string]*pb.TenantUser
}

// CreateTenantServiceDAO - returns an in-memory implementation of the Tenant Service
// datastore.
func CreateTenantServiceDAO() (*TenantServiceDatastoreInMemory, error) {
	res := new(TenantServiceDatastoreInMemory)

	res.tenantToIDtoTenantUserMap = map[string]map[string]*pb.TenantUser{}

	return res, nil
}

func (tsd *TenantServiceDatastoreInMemory) doesTenantExist(tenantID string) error {
	if len(tenantID) == 0 {
		return fmt.Errorf("%s does not exist", tenantID)
	}
	if tsd.tenantToIDtoTenantUserMap[tenantID] == nil {
		return fmt.Errorf("%s does not exist", tenantID)
	}

	return nil
}

// CreateTenantUser - InMemory implementation of CreateTenantUser
func (tsd *TenantServiceDatastoreInMemory) CreateTenantUser(tenantUserRequest *pb.TenantUser) (*pb.TenantUser, error) {
	if len(tenantUserRequest.XId) != 0 {
		return nil, fmt.Errorf("%s already exists", ds.TenantUserStr)
	}
	if err := tsd.doesTenantExist(tenantUserRequest.Data.TenantId); err != nil {
		// Make a place for the tenant
		tsd.tenantToIDtoTenantUserMap[tenantUserRequest.Data.TenantId] = map[string]*pb.TenantUser{}
	}

	userCopy := pb.TenantUser{}
	deepcopy.Copy(&userCopy, tenantUserRequest)
	userCopy.XId = uuid.NewV4().String()
	userCopy.XRev = uuid.NewV4().String()
	userCopy.Data.Datatype = string(ds.TenantUserType)
	userCopy.Data.CreatedTimestamp = ds.MakeTimestamp()
	userCopy.Data.LastModifiedTimestamp = userCopy.Data.GetCreatedTimestamp()

	tsd.tenantToIDtoTenantUserMap[tenantUserRequest.Data.TenantId][userCopy.XId] = &userCopy

	return &userCopy, nil
}

// UpdateTenantUser - InMemory implementation of UpdateTenantUser
func (tsd *TenantServiceDatastoreInMemory) UpdateTenantUser(tenantUserRequest *pb.TenantUser) (*pb.TenantUser, error) {
	if len(tenantUserRequest.XId) == 0 {
		return nil, fmt.Errorf("%s must have an ID", ds.AdminUserStr)
	}
	if len(tenantUserRequest.XRev) == 0 {
		return nil, fmt.Errorf("%s must have a revision", ds.AdminUserStr)
	}
	if err := tsd.doesTenantExist(tenantUserRequest.Data.TenantId); err != nil {
		return nil, err
	}

	userCopy := pb.TenantUser{}
	deepcopy.Copy(&userCopy, tenantUserRequest)
	userCopy.XRev = uuid.NewV4().String()
	userCopy.Data.Datatype = string(ds.TenantUserType)
	userCopy.Data.LastModifiedTimestamp = ds.MakeTimestamp()

	tsd.tenantToIDtoTenantUserMap[tenantUserRequest.Data.TenantId][userCopy.XId] = &userCopy

	return &userCopy, nil
}

// DeleteTenantUser - InMemory implementation of DeleteTenantUser
func (tsd *TenantServiceDatastoreInMemory) DeleteTenantUser(tenantUserIDRequest *pb.TenantUserIdRequest) (*pb.TenantUser, error) {
	if len(tenantUserIDRequest.UserId) == 0 {
		return nil, fmt.Errorf("%s must provide a User ID", ds.TenantUserStr)
	}
	if len(tenantUserIDRequest.TenantId) == 0 {
		return nil, fmt.Errorf("%s must provide a Tenant ID", ds.TenantUserStr)
	}
	if err := tsd.doesTenantExist(tenantUserIDRequest.TenantId); err != nil {
		return nil, err
	}

	user, ok := tsd.tenantToIDtoTenantUserMap[tenantUserIDRequest.TenantId][tenantUserIDRequest.UserId]
	if ok {
		delete(tsd.tenantToIDtoTenantUserMap[tenantUserIDRequest.TenantId], tenantUserIDRequest.UserId)

		// Delete the tenant user map if there are no more users.
		if len(tsd.tenantToIDtoTenantUserMap[tenantUserIDRequest.TenantId]) == 0 {
			delete(tsd.tenantToIDtoTenantUserMap, tenantUserIDRequest.TenantId)
		}
		return user, nil
	}

	return nil, fmt.Errorf("%s not found", ds.TenantUserStr)
}

// GetTenantUser - InMemory implementation of GetTenantUser
func (tsd *TenantServiceDatastoreInMemory) GetTenantUser(tenantUserIDRequest *pb.TenantUserIdRequest) (*pb.TenantUser, error) {
	if len(tenantUserIDRequest.UserId) == 0 {
		return nil, fmt.Errorf("%s must provide a User ID", ds.TenantUserStr)
	}
	if len(tenantUserIDRequest.TenantId) == 0 {
		return nil, fmt.Errorf("%s must provide a Tenant ID", ds.TenantUserStr)
	}
	if err := tsd.doesTenantExist(tenantUserIDRequest.TenantId); err != nil {
		return nil, err
	}

	user, ok := tsd.tenantToIDtoTenantUserMap[tenantUserIDRequest.TenantId][tenantUserIDRequest.UserId]
	if ok {
		return user, nil
	}

	return nil, fmt.Errorf("%s not found", ds.TenantUserStr)
}

// GetAllTenantUsers - InMemory implementation of GetAllTenantUsers
func (tsd *TenantServiceDatastoreInMemory) GetAllTenantUsers(tenantID string) (*pb.TenantUserList, error) {
	err := tsd.doesTenantExist(tenantID)
	if err != nil {
		return &pb.TenantUserList{Data: []*pb.TenantUser{}}, nil
	}

	tenantUserList := pb.TenantUserList{}
	tenantUserList.Data = make([]*pb.TenantUser, 0)

	for _, user := range tsd.tenantToIDtoTenantUserMap[tenantID] {
		tenantUserList.Data = append(tenantUserList.Data, user)
	}

	return &tenantUserList, nil
}

// CreateTenantDomain - InMemory implementation of CreateTenantDomain
func (tsd *TenantServiceDatastoreInMemory) CreateTenantDomain(tenantDomainRequest *pb.TenantDomain) (*pb.TenantDomain, error) {
	// Stub to implement
	return nil, errors.New("Unsupported operation: CreateTenantDomain not implemented")
}

// UpdateTenantDomain - InMemory implementation of UpdateTenantDomain
func (tsd *TenantServiceDatastoreInMemory) UpdateTenantDomain(tenantDomainRequest *pb.TenantDomain) (*pb.TenantDomain, error) {
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
func (tsd *TenantServiceDatastoreInMemory) CreateTenantIngestionProfile(tenantIngPrfReq *pb.TenantIngestionProfile) (*pb.TenantIngestionProfile, error) {
	// Stub to implement
	return nil, errors.New("Unsupported operation: CreateTenantIngestionProfile not implemented")
}

// UpdateTenantIngestionProfile - InMemory implementation of UpdateTenantIngestionProfile
func (tsd *TenantServiceDatastoreInMemory) UpdateTenantIngestionProfile(tenantIngPrfReq *pb.TenantIngestionProfile) (*pb.TenantIngestionProfile, error) {
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

// CreateTenantThresholdProfile - InMemory implementation of CreateTenantThresholdProfile
func (tsd *TenantServiceDatastoreInMemory) CreateTenantThresholdProfile(tenantThreshPrfReq *pb.TenantThresholdProfile) (*pb.TenantThresholdProfile, error) {
	// Stub to implement
	return nil, errors.New("Unsupported operation: CreateTenantThresholdProfile not implemented")
}

// UpdateTenantThresholdProfile - InMemory implementation of UpdateTenantThresholdProfile
func (tsd *TenantServiceDatastoreInMemory) UpdateTenantThresholdProfile(tenantThreshPrfReq *pb.TenantThresholdProfile) (*pb.TenantThresholdProfile, error) {
	// Stub to implement
	return nil, errors.New("Unsupported operation: UpdateTenantThresholdProfile not implemented")
}

// GetTenantThresholdProfile - InMemory implementation of GetTenantThresholdProfile
func (tsd *TenantServiceDatastoreInMemory) GetTenantThresholdProfile(tenantIngPrfIDReq *pb.TenantThresholdProfileIdRequest) (*pb.TenantThresholdProfile, error) {
	// Stub to implement
	return nil, errors.New("Unsupported operation: GetTenantThresholdProfile not implemented")
}

// DeleteTenantThresholdProfile - InMemory implementation of DeleteTenantThresholdProfile
func (tsd *TenantServiceDatastoreInMemory) DeleteTenantThresholdProfile(tenantIngPrfIDReq *pb.TenantThresholdProfileIdRequest) (*pb.TenantThresholdProfile, error) {
	// Stub to implement
	return nil, errors.New("Unsupported operation: DeleteTenantThresholdProfile not implemented")
}

// CreateMonitoredObject - InMemory implementation of CreateMonitoredObject
func (tsd *TenantServiceDatastoreInMemory) CreateMonitoredObject(monitoredObjectReq *pb.MonitoredObject) (*pb.MonitoredObject, error) {
	// Stub to implement
	return nil, errors.New("Unsupported operation: CreateMonitoredObject not implemented")
}

// UpdateMonitoredObject - InMemory implementation of UpdateMonitoredObject
func (tsd *TenantServiceDatastoreInMemory) UpdateMonitoredObject(monitoredObjectReq *pb.MonitoredObject) (*pb.MonitoredObject, error) {
	// Stub to implement
	return nil, errors.New("Unsupported operation: UpdateMonitoredObject not implemented")
}

// GetMonitoredObject - InMemory implementation of GetMonitoredObject
func (tsd *TenantServiceDatastoreInMemory) GetMonitoredObject(monitoredObjectIDReq *pb.MonitoredObjectIdRequest) (*pb.MonitoredObject, error) {
	// Stub to implement
	return nil, errors.New("Unsupported operation: GetMonitoredObject not implemented")
}

// DeleteMonitoredObject - InMemory implementation of DeleteMonitoredObject
func (tsd *TenantServiceDatastoreInMemory) DeleteMonitoredObject(monitoredObjectIDReq *pb.MonitoredObjectIdRequest) (*pb.MonitoredObject, error) {
	// Stub to implement
	return nil, errors.New("Unsupported operation: DeleteMonitoredObject not implemented")
}

// GetAllMonitoredObjects - InMemory implementation of GetAllMonitoredObjects
func (tsd *TenantServiceDatastoreInMemory) GetAllMonitoredObjects(tenantID string) (*pb.MonitoredObjectList, error) {
	// Stub to implement
	return nil, errors.New("Unsupported operation: GetAllMonitoredObjects not implemented")
}

// GetMonitoredObjectToDomainMap - InMemory implementation of GetMonitoredObjectToDomainMap
func (tsd *TenantServiceDatastoreInMemory) GetMonitoredObjectToDomainMap(moByDomReq *pb.MonitoredObjectCountByDomainRequest) (*pb.MonitoredObjectCountByDomainResponse, error) {
	// Stub to implement
	return nil, errors.New("Unsupported operation: GetMonitoredObjectToDomainMap not implemented")
}

// CreateTenantMeta - InMemory implementation of CreateTenantMeta
func (tsd *TenantServiceDatastoreInMemory) CreateTenantMeta(meta *pb.TenantMetadata) (*pb.TenantMetadata, error) {
	// Stub to implement
	return nil, errors.New("Unsupported operation: CreateTenantMeta not implemented")
}

// UpdateTenantMeta - InMemory implementation of UpdateTenantMeta
func (tsd *TenantServiceDatastoreInMemory) UpdateTenantMeta(meta *pb.TenantMetadata) (*pb.TenantMetadata, error) {
	// Stub to implement
	return nil, errors.New("Unsupported operation: UpdateTenantMeta not implemented")
}

// DeleteTenantMeta - InMemory implementation of DeleteTenantMeta
func (tsd *TenantServiceDatastoreInMemory) DeleteTenantMeta(tenantID string) (*pb.TenantMetadata, error) {
	// Stub to implement
	return nil, errors.New("Unsupported operation: DeleteTenantMeta not implemented")
}

// GetTenantMeta - InMemory implementation of GetTenantMeta
func (tsd *TenantServiceDatastoreInMemory) GetTenantMeta(tenantID string) (*pb.TenantMetadata, error) {
	// Stub to implement
	return nil, errors.New("Unsupported operation: GetTenantMeta not implemented")
}

// GetActiveTenantIngestionProfile - InMemory implementation of GetActiveTenantIngestionProfile
func (tsd *TenantServiceDatastoreInMemory) GetActiveTenantIngestionProfile(tenantID string) (*pb.TenantIngestionProfile, error) {
	// Stub to implement
	return nil, errors.New("Unsupported operation: GetActiveTenantIngestionProfile not implemented")
}

// GetAllTenantThresholdProfile - InMemory implementation of GetAllTenantThresholdProfile
func (tsd *TenantServiceDatastoreInMemory) GetAllTenantThresholdProfile(tenantID string) (*pb.TenantThresholdProfileList, error) {
	// Stub to implement
	return nil, errors.New("Unsupported operation: GetAllTenantThresholdProfile not implemented")
}

// BulkInsertMonitoredObjects - InMemory implementation of BulkInsertMonitoredObjects
func (tsd *TenantServiceDatastoreInMemory) BulkInsertMonitoredObjects(value *pb.TenantMonitoredObjectSet) (*pb.BulkOperationResponse, error) {
	// Stub to implement
	return nil, errors.New("BulkInsertMonitoredObjects() not implemented for InMemory DB")
}
