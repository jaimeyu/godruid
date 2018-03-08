package inMemory

import (
	"errors"
	"fmt"

	ds "github.com/accedian/adh-gather/datastore"
	pb "github.com/accedian/adh-gather/gathergrpc"
	"github.com/getlantern/deepcopy"
	uuid "github.com/satori/go.uuid"

	tenmod "github.com/accedian/adh-gather/models/tenant"
)

// TenantServiceDatastoreInMemory - struct responsible for handling
// database operations for the Admin Service when using local memory
// as the storage option. Useful for tests.
type TenantServiceDatastoreInMemory struct {
	tenantToIDtoTenantUserMap   map[string]map[string]*pb.TenantUser
	tenantToIDtoTenantDomainMap map[string]map[string]*pb.TenantDomain
}

// CreateTenantServiceDAO - returns an in-memory implementation of the Tenant Service
// datastore.
func CreateTenantServiceDAO() (*TenantServiceDatastoreInMemory, error) {
	res := new(TenantServiceDatastoreInMemory)

	res.tenantToIDtoTenantUserMap = map[string]map[string]*pb.TenantUser{}
	res.tenantToIDtoTenantDomainMap = map[string]map[string]*pb.TenantDomain{}

	return res, nil
}

func (tsd *TenantServiceDatastoreInMemory) doesTenantExist(tenantID string, ctx tenmod.TenantDataType) error {
	if len(tenantID) == 0 {
		return fmt.Errorf("%s does not exist", tenantID)
	}
	switch ctx {
	case tenmod.TenantUserType:
		if tsd.tenantToIDtoTenantUserMap[tenantID] == nil {
			return fmt.Errorf("%s does not exist", tenantID)
		}
	case tenmod.TenantDomainType:
		if tsd.tenantToIDtoTenantDomainMap[tenantID] == nil {
			return fmt.Errorf("%s does not exist", tenantID)
		}
	}

	return nil
}

// CreateTenantUser - InMemory implementation of CreateTenantUser
func (tsd *TenantServiceDatastoreInMemory) CreateTenantUser(tenantUserRequest *pb.TenantUser) (*pb.TenantUser, error) {
	if len(tenantUserRequest.XId) != 0 {
		return nil, fmt.Errorf("%s already exists", tenmod.TenantUserStr)
	}
	if err := tsd.doesTenantExist(tenantUserRequest.Data.TenantId, tenmod.TenantUserType); err != nil {
		// Make a place for the tenant
		tsd.tenantToIDtoTenantUserMap[tenantUserRequest.Data.TenantId] = map[string]*pb.TenantUser{}
	}

	userCopy := pb.TenantUser{}
	deepcopy.Copy(&userCopy, tenantUserRequest)
	userCopy.XId = uuid.NewV4().String()
	userCopy.XRev = uuid.NewV4().String()
	userCopy.Data.Datatype = string(tenmod.TenantUserType)
	userCopy.Data.CreatedTimestamp = ds.MakeTimestamp()
	userCopy.Data.LastModifiedTimestamp = userCopy.Data.GetCreatedTimestamp()

	tsd.tenantToIDtoTenantUserMap[tenantUserRequest.Data.TenantId][userCopy.XId] = &userCopy

	return &userCopy, nil
}

// UpdateTenantUser - InMemory implementation of UpdateTenantUser
func (tsd *TenantServiceDatastoreInMemory) UpdateTenantUser(tenantUserRequest *pb.TenantUser) (*pb.TenantUser, error) {
	if len(tenantUserRequest.XId) == 0 {
		return nil, fmt.Errorf("%s must have an ID", tenmod.TenantUserStr)
	}
	if len(tenantUserRequest.XRev) == 0 {
		return nil, fmt.Errorf("%s must have a revision", tenmod.TenantUserStr)
	}
	if err := tsd.doesTenantExist(tenantUserRequest.Data.TenantId, tenmod.TenantUserType); err != nil {
		return nil, fmt.Errorf("%s does not exist", tenmod.TenantUserStr)
	}

	userCopy := pb.TenantUser{}
	deepcopy.Copy(&userCopy, tenantUserRequest)
	userCopy.XRev = uuid.NewV4().String()
	userCopy.Data.Datatype = string(tenmod.TenantUserType)
	userCopy.Data.LastModifiedTimestamp = ds.MakeTimestamp()

	tsd.tenantToIDtoTenantUserMap[tenantUserRequest.Data.TenantId][userCopy.XId] = &userCopy

	return &userCopy, nil
}

// DeleteTenantUser - InMemory implementation of DeleteTenantUser
func (tsd *TenantServiceDatastoreInMemory) DeleteTenantUser(tenantUserIDRequest *pb.TenantUserIdRequest) (*pb.TenantUser, error) {
	if len(tenantUserIDRequest.UserId) == 0 {
		return nil, fmt.Errorf("%s must provide a User ID", tenmod.TenantUserStr)
	}
	if len(tenantUserIDRequest.TenantId) == 0 {
		return nil, fmt.Errorf("%s must provide a Tenant ID", tenmod.TenantUserStr)
	}
	if err := tsd.doesTenantExist(tenantUserIDRequest.TenantId, tenmod.TenantUserType); err != nil {
		return nil, fmt.Errorf("%s does not exist", tenmod.TenantUserStr)
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

	return nil, fmt.Errorf("%s not found", tenmod.TenantUserStr)
}

// GetTenantUser - InMemory implementation of GetTenantUser
func (tsd *TenantServiceDatastoreInMemory) GetTenantUser(tenantUserIDRequest *pb.TenantUserIdRequest) (*pb.TenantUser, error) {
	if len(tenantUserIDRequest.UserId) == 0 {
		return nil, fmt.Errorf("%s must provide a User ID", tenmod.TenantUserStr)
	}
	if len(tenantUserIDRequest.TenantId) == 0 {
		return nil, fmt.Errorf("%s must provide a Tenant ID", tenmod.TenantUserStr)
	}
	if err := tsd.doesTenantExist(tenantUserIDRequest.TenantId, tenmod.TenantUserType); err != nil {
		return nil, fmt.Errorf("%s does not exist", tenmod.TenantUserStr)
	}

	user, ok := tsd.tenantToIDtoTenantUserMap[tenantUserIDRequest.TenantId][tenantUserIDRequest.UserId]
	if ok {
		return user, nil
	}

	return nil, fmt.Errorf("%s not found", tenmod.TenantUserStr)
}

// GetAllTenantUsers - InMemory implementation of GetAllTenantUsers
func (tsd *TenantServiceDatastoreInMemory) GetAllTenantUsers(tenantID string) (*pb.TenantUserList, error) {
	err := tsd.doesTenantExist(tenantID, tenmod.TenantUserType)
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
	if len(tenantDomainRequest.XId) != 0 {
		return nil, fmt.Errorf("%s already exists", tenmod.TenantDomainStr)
	}
	if err := tsd.doesTenantExist(tenantDomainRequest.Data.TenantId, tenmod.TenantDomainType); err != nil {
		// Make a place for the tenant
		tsd.tenantToIDtoTenantDomainMap[tenantDomainRequest.Data.TenantId] = map[string]*pb.TenantDomain{}
	}

	recCopy := pb.TenantDomain{}
	deepcopy.Copy(&recCopy, tenantDomainRequest)
	recCopy.XId = uuid.NewV4().String()
	recCopy.XRev = uuid.NewV4().String()
	recCopy.Data.Datatype = string(tenmod.TenantDomainType)
	recCopy.Data.CreatedTimestamp = ds.MakeTimestamp()
	recCopy.Data.LastModifiedTimestamp = recCopy.Data.GetCreatedTimestamp()

	tsd.tenantToIDtoTenantDomainMap[tenantDomainRequest.Data.TenantId][recCopy.XId] = &recCopy

	return &recCopy, nil
}

// UpdateTenantDomain - InMemory implementation of UpdateTenantDomain
func (tsd *TenantServiceDatastoreInMemory) UpdateTenantDomain(tenantDomainRequest *pb.TenantDomain) (*pb.TenantDomain, error) {
	if len(tenantDomainRequest.XId) == 0 {
		return nil, fmt.Errorf("%s must have an ID", tenmod.TenantDomainStr)
	}
	if len(tenantDomainRequest.XRev) == 0 {
		return nil, fmt.Errorf("%s must have a revision", tenmod.TenantDomainStr)
	}
	if err := tsd.doesTenantExist(tenantDomainRequest.Data.TenantId, tenmod.TenantDomainType); err != nil {
		return nil, fmt.Errorf("%s does not exist", tenmod.TenantDomainStr)
	}

	recCopy := pb.TenantDomain{}
	deepcopy.Copy(&recCopy, tenantDomainRequest)
	recCopy.XRev = uuid.NewV4().String()
	recCopy.Data.Datatype = string(tenmod.TenantDomainType)
	recCopy.Data.LastModifiedTimestamp = ds.MakeTimestamp()

	tsd.tenantToIDtoTenantDomainMap[tenantDomainRequest.Data.TenantId][recCopy.XId] = &recCopy

	return &recCopy, nil
}

// DeleteTenantDomain - InMemory implementation of DeleteTenantDomain
func (tsd *TenantServiceDatastoreInMemory) DeleteTenantDomain(tenantDomainIDRequest *pb.TenantDomainIdRequest) (*pb.TenantDomain, error) {
	if len(tenantDomainIDRequest.DomainId) == 0 {
		return nil, fmt.Errorf("%s must provide a Domain ID", tenmod.TenantDomainStr)
	}
	if len(tenantDomainIDRequest.TenantId) == 0 {
		return nil, fmt.Errorf("%s must provide a Tenant ID", tenmod.TenantDomainStr)
	}
	if err := tsd.doesTenantExist(tenantDomainIDRequest.TenantId, tenmod.TenantDomainType); err != nil {
		return nil, fmt.Errorf("%s does not exist", tenmod.TenantDomainStr)
	}

	rec, ok := tsd.tenantToIDtoTenantDomainMap[tenantDomainIDRequest.TenantId][tenantDomainIDRequest.DomainId]
	if ok {
		delete(tsd.tenantToIDtoTenantDomainMap[tenantDomainIDRequest.TenantId], tenantDomainIDRequest.DomainId)

		// Delete the tenant user map if there are no more users.
		if len(tsd.tenantToIDtoTenantDomainMap[tenantDomainIDRequest.TenantId]) == 0 {
			delete(tsd.tenantToIDtoTenantDomainMap, tenantDomainIDRequest.TenantId)
		}
		return rec, nil
	}

	return nil, fmt.Errorf("%s not found", tenmod.TenantDomainStr)
}

// GetTenantDomain - InMemory implementation of GetTenantDomain
func (tsd *TenantServiceDatastoreInMemory) GetTenantDomain(tenantDomainIDRequest *pb.TenantDomainIdRequest) (*pb.TenantDomain, error) {
	if len(tenantDomainIDRequest.DomainId) == 0 {
		return nil, fmt.Errorf("%s must provide a Domain ID", tenmod.TenantDomainStr)
	}
	if len(tenantDomainIDRequest.TenantId) == 0 {
		return nil, fmt.Errorf("%s must provide a Tenant ID", tenmod.TenantDomainStr)
	}
	if err := tsd.doesTenantExist(tenantDomainIDRequest.TenantId, tenmod.TenantDomainType); err != nil {
		return nil, fmt.Errorf("%s does not exist", tenmod.TenantDomainStr)
	}

	rec, ok := tsd.tenantToIDtoTenantDomainMap[tenantDomainIDRequest.TenantId][tenantDomainIDRequest.DomainId]
	if ok {
		return rec, nil
	}

	return nil, fmt.Errorf("%s not found", tenmod.TenantDomainStr)
}

// GetAllTenantDomains - InMemory implementation of GetAllTenantDomains
func (tsd *TenantServiceDatastoreInMemory) GetAllTenantDomains(tenantID string) (*pb.TenantDomainList, error) {
	err := tsd.doesTenantExist(tenantID, tenmod.TenantDomainType)
	if err != nil {
		return &pb.TenantDomainList{Data: []*pb.TenantDomain{}}, nil
	}

	recList := pb.TenantDomainList{}
	recList.Data = make([]*pb.TenantDomain, 0)

	for _, rec := range tsd.tenantToIDtoTenantDomainMap[tenantID] {
		recList.Data = append(recList.Data, rec)
	}

	return &recList, nil
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
