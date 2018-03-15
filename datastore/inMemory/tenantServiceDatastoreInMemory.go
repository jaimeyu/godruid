package inMemory

import (
	"errors"
	"fmt"

	ds "github.com/accedian/adh-gather/datastore"
	"github.com/accedian/adh-gather/logger"
	"github.com/accedian/adh-gather/models"
	"github.com/getlantern/deepcopy"
	uuid "github.com/satori/go.uuid"

	"github.com/accedian/adh-gather/models/common"
	tenmod "github.com/accedian/adh-gather/models/tenant"
)

// TenantServiceDatastoreInMemory - struct responsible for handling
// database operations for the Admin Service when using local memory
// as the storage option. Useful for tests.
type TenantServiceDatastoreInMemory struct {
	tenantToIDtoTenantUserMap   map[string]map[string]*tenmod.User
	tenantToIDtoTenantDomainMap map[string]map[string]*tenmod.Domain
}

// CreateTenantServiceDAO - returns an in-memory implementation of the Tenant Service
// datastore.
func CreateTenantServiceDAO() (*TenantServiceDatastoreInMemory, error) {
	res := new(TenantServiceDatastoreInMemory)

	res.tenantToIDtoTenantUserMap = map[string]map[string]*tenmod.User{}
	res.tenantToIDtoTenantDomainMap = map[string]map[string]*tenmod.Domain{}

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
func (tsd *TenantServiceDatastoreInMemory) CreateTenantUser(tenantUserRequest *tenmod.User) (*tenmod.User, error) {
	if len(tenantUserRequest.ID) != 0 {
		return nil, fmt.Errorf("%s already exists", tenmod.TenantUserStr)
	}
	if err := tsd.doesTenantExist(tenantUserRequest.TenantID, tenmod.TenantUserType); err != nil {
		// Make a place for the tenant
		tsd.tenantToIDtoTenantUserMap[tenantUserRequest.TenantID] = map[string]*tenmod.User{}
	}

	userCopy := tenmod.User{}
	deepcopy.Copy(&userCopy, tenantUserRequest)
	userCopy.ID = uuid.NewV4().String()
	userCopy.REV = uuid.NewV4().String()
	userCopy.Datatype = string(tenmod.TenantUserType)
	userCopy.CreatedTimestamp = ds.MakeTimestamp()
	userCopy.LastModifiedTimestamp = userCopy.CreatedTimestamp

	tsd.tenantToIDtoTenantUserMap[tenantUserRequest.TenantID][userCopy.ID] = &userCopy

	return &userCopy, nil
}

// UpdateTenantUser - InMemory implementation of UpdateTenantUser
func (tsd *TenantServiceDatastoreInMemory) UpdateTenantUser(tenantUserRequest *tenmod.User) (*tenmod.User, error) {
	if len(tenantUserRequest.ID) == 0 {
		return nil, fmt.Errorf("%s must have an ID", tenmod.TenantUserStr)
	}
	if len(tenantUserRequest.REV) == 0 {
		return nil, fmt.Errorf("%s must have a revision", tenmod.TenantUserStr)
	}
	if err := tsd.doesTenantExist(tenantUserRequest.TenantID, tenmod.TenantUserType); err != nil {
		return nil, fmt.Errorf("%s does not exist", tenmod.TenantUserStr)
	}

	userCopy := tenmod.User{}
	deepcopy.Copy(&userCopy, tenantUserRequest)
	userCopy.REV = uuid.NewV4().String()
	userCopy.Datatype = string(tenmod.TenantUserType)
	userCopy.LastModifiedTimestamp = ds.MakeTimestamp()

	tsd.tenantToIDtoTenantUserMap[tenantUserRequest.TenantID][userCopy.ID] = &userCopy

	return &userCopy, nil
}

// DeleteTenantUser - InMemory implementation of DeleteTenantUser
func (tsd *TenantServiceDatastoreInMemory) DeleteTenantUser(tenantID string, userID string) (*tenmod.User, error) {
	if len(userID) == 0 {
		return nil, fmt.Errorf("%s must provide a User ID", tenmod.TenantUserStr)
	}
	if len(tenantID) == 0 {
		return nil, fmt.Errorf("%s must provide a Tenant ID", tenmod.TenantUserStr)
	}
	if err := tsd.doesTenantExist(tenantID, tenmod.TenantUserType); err != nil {
		return nil, fmt.Errorf("%s does not exist", tenmod.TenantUserStr)
	}

	user, ok := tsd.tenantToIDtoTenantUserMap[tenantID][userID]
	if ok {
		delete(tsd.tenantToIDtoTenantUserMap[tenantID], userID)

		// Delete the tenant user map if there are no more users.
		if len(tsd.tenantToIDtoTenantUserMap[tenantID]) == 0 {
			delete(tsd.tenantToIDtoTenantUserMap, tenantID)
		}
		return user, nil
	}

	return nil, fmt.Errorf("%s not found", tenmod.TenantUserStr)
}

// GetTenantUser - InMemory implementation of GetTenantUser
func (tsd *TenantServiceDatastoreInMemory) GetTenantUser(tenantID string, userID string) (*tenmod.User, error) {
	if len(userID) == 0 {
		return nil, fmt.Errorf("%s must provide a User ID", tenmod.TenantUserStr)
	}
	if len(tenantID) == 0 {
		return nil, fmt.Errorf("%s must provide a Tenant ID", tenmod.TenantUserStr)
	}
	if err := tsd.doesTenantExist(tenantID, tenmod.TenantUserType); err != nil {
		return nil, fmt.Errorf("%s does not exist", tenmod.TenantUserStr)
	}

	user, ok := tsd.tenantToIDtoTenantUserMap[tenantID][userID]
	if ok {
		return user, nil
	}

	return nil, fmt.Errorf("%s not found", tenmod.TenantUserStr)
}

// GetAllTenantUsers - InMemory implementation of GetAllTenantUsers
func (tsd *TenantServiceDatastoreInMemory) GetAllTenantUsers(tenantID string) ([]*tenmod.User, error) {
	err := tsd.doesTenantExist(tenantID, tenmod.TenantUserType)
	if err != nil {
		return []*tenmod.User{}, nil
	}

	tenantUserList := make([]*tenmod.User, 0)

	for _, user := range tsd.tenantToIDtoTenantUserMap[tenantID] {
		tenantUserList = append(tenantUserList, user)
	}

	return tenantUserList, nil
}

// CreateTenantDomain - InMemory implementation of CreateTenantDomain
func (tsd *TenantServiceDatastoreInMemory) CreateTenantDomain(tenantDomainRequest *tenmod.Domain) (*tenmod.Domain, error) {
	if len(tenantDomainRequest.ID) != 0 {
		return nil, fmt.Errorf("%s already exists", tenmod.TenantDomainStr)
	}
	if err := tsd.doesTenantExist(tenantDomainRequest.TenantID, tenmod.TenantDomainType); err != nil {
		// Make a place for the tenant
		tsd.tenantToIDtoTenantDomainMap[tenantDomainRequest.TenantID] = map[string]*tenmod.Domain{}
	}

	recCopy := tenmod.Domain{}
	deepcopy.Copy(&recCopy, tenantDomainRequest)
	recCopy.ID = uuid.NewV4().String()
	recCopy.REV = uuid.NewV4().String()
	recCopy.Datatype = string(tenmod.TenantDomainType)
	recCopy.CreatedTimestamp = ds.MakeTimestamp()
	recCopy.LastModifiedTimestamp = recCopy.CreatedTimestamp

	tsd.tenantToIDtoTenantDomainMap[tenantDomainRequest.TenantID][recCopy.ID] = &recCopy

	return &recCopy, nil
}

// UpdateTenantDomain - InMemory implementation of UpdateTenantDomain
func (tsd *TenantServiceDatastoreInMemory) UpdateTenantDomain(tenantDomainRequest *tenmod.Domain) (*tenmod.Domain, error) {
	if len(tenantDomainRequest.ID) == 0 {
		return nil, fmt.Errorf("%s must have an ID", tenmod.TenantDomainStr)
	}
	if len(tenantDomainRequest.REV) == 0 {
		return nil, fmt.Errorf("%s must have a revision", tenmod.TenantDomainStr)
	}
	if err := tsd.doesTenantExist(tenantDomainRequest.TenantID, tenmod.TenantDomainType); err != nil {
		return nil, fmt.Errorf("%s does not exist", tenmod.TenantDomainStr)
	}

	recCopy := tenmod.Domain{}
	deepcopy.Copy(&recCopy, tenantDomainRequest)
	recCopy.REV = uuid.NewV4().String()
	recCopy.Datatype = string(tenmod.TenantDomainType)
	recCopy.LastModifiedTimestamp = ds.MakeTimestamp()

	tsd.tenantToIDtoTenantDomainMap[tenantDomainRequest.TenantID][recCopy.ID] = &recCopy

	return &recCopy, nil
}

// DeleteTenantDomain - InMemory implementation of DeleteTenantDomain
func (tsd *TenantServiceDatastoreInMemory) DeleteTenantDomain(tenantID string, dataID string) (*tenmod.Domain, error) {
	if len(dataID) == 0 {
		return nil, fmt.Errorf("%s must provide a Domain ID", tenmod.TenantDomainStr)
	}
	if len(tenantID) == 0 {
		return nil, fmt.Errorf("%s must provide a Tenant ID", tenmod.TenantDomainStr)
	}
	if err := tsd.doesTenantExist(tenantID, tenmod.TenantDomainType); err != nil {
		return nil, fmt.Errorf("%s does not exist", tenmod.TenantDomainStr)
	}

	rec, ok := tsd.tenantToIDtoTenantDomainMap[tenantID][dataID]
	logger.Log.Debugf(models.AsJSONString(tsd.tenantToIDtoTenantDomainMap))
	if ok {
		delete(tsd.tenantToIDtoTenantDomainMap[tenantID], dataID)

		// Delete the tenant user map if there are no more users.
		if len(tsd.tenantToIDtoTenantDomainMap[tenantID]) == 0 {
			delete(tsd.tenantToIDtoTenantDomainMap, tenantID)
		}
		return rec, nil
	}

	return nil, fmt.Errorf("%s not found", tenmod.TenantDomainStr)
}

// GetTenantDomain - InMemory implementation of GetTenantDomain
func (tsd *TenantServiceDatastoreInMemory) GetTenantDomain(tenantID string, dataID string) (*tenmod.Domain, error) {
	if len(dataID) == 0 {
		return nil, fmt.Errorf("%s must provide a Domain ID", tenmod.TenantDomainStr)
	}
	if len(tenantID) == 0 {
		return nil, fmt.Errorf("%s must provide a Tenant ID", tenmod.TenantDomainStr)
	}
	if err := tsd.doesTenantExist(tenantID, tenmod.TenantDomainType); err != nil {
		return nil, fmt.Errorf("%s does not exist", tenmod.TenantDomainStr)
	}

	rec, ok := tsd.tenantToIDtoTenantDomainMap[tenantID][dataID]
	if ok {
		return rec, nil
	}

	return nil, fmt.Errorf("%s not found", tenmod.TenantDomainStr)
}

// GetAllTenantDomains - InMemory implementation of GetAllTenantDomains
func (tsd *TenantServiceDatastoreInMemory) GetAllTenantDomains(tenantID string) ([]*tenmod.Domain, error) {
	err := tsd.doesTenantExist(tenantID, tenmod.TenantDomainType)
	if err != nil {
		return []*tenmod.Domain{}, nil
	}

	recList := make([]*tenmod.Domain, 0)

	for _, rec := range tsd.tenantToIDtoTenantDomainMap[tenantID] {
		recList = append(recList, rec)
	}

	return recList, nil
}

// CreateTenantIngestionProfile - InMemory implementation of CreateTenantIngestionProfile
func (tsd *TenantServiceDatastoreInMemory) CreateTenantIngestionProfile(tenantIngPrfReq *tenmod.IngestionProfile) (*tenmod.IngestionProfile, error) {
	// Stub to implement
	return nil, errors.New("Unsupported operation: CreateTenantIngestionProfile not implemented")
}

// UpdateTenantIngestionProfile - InMemory implementation of UpdateTenantIngestionProfile
func (tsd *TenantServiceDatastoreInMemory) UpdateTenantIngestionProfile(tenantIngPrfReq *tenmod.IngestionProfile) (*tenmod.IngestionProfile, error) {
	// Stub to implement
	return nil, errors.New("Unsupported operation: UpdateTenantIngestionProfile not implemented")
}

// GetTenantIngestionProfile - InMemory implementation of GetTenantIngestionProfile
func (tsd *TenantServiceDatastoreInMemory) GetTenantIngestionProfile(tenantID string, dataID string) (*tenmod.IngestionProfile, error) {
	// Stub to implement
	return nil, errors.New("Unsupported operation: GetTenantIngestionProfile not implemented")
}

// DeleteTenantIngestionProfile - InMemory implementation of DeleteTenantIngestionProfile
func (tsd *TenantServiceDatastoreInMemory) DeleteTenantIngestionProfile(tenantID string, dataID string) (*tenmod.IngestionProfile, error) {
	// Stub to implement
	return nil, errors.New("Unsupported operation: DeleteTenantIngestionProfile not implemented")
}

// CreateTenantThresholdProfile - InMemory implementation of CreateTenantThresholdProfile
func (tsd *TenantServiceDatastoreInMemory) CreateTenantThresholdProfile(tenantThreshPrfReq *tenmod.ThresholdProfile) (*tenmod.ThresholdProfile, error) {
	// Stub to implement
	return nil, errors.New("Unsupported operation: CreateTenantThresholdProfile not implemented")
}

// UpdateTenantThresholdProfile - InMemory implementation of UpdateTenantThresholdProfile
func (tsd *TenantServiceDatastoreInMemory) UpdateTenantThresholdProfile(tenantThreshPrfReq *tenmod.ThresholdProfile) (*tenmod.ThresholdProfile, error) {
	// Stub to implement
	return nil, errors.New("Unsupported operation: UpdateTenantThresholdProfile not implemented")
}

// GetTenantThresholdProfile - InMemory implementation of GetTenantThresholdProfile
func (tsd *TenantServiceDatastoreInMemory) GetTenantThresholdProfile(tenantID string, dataID string) (*tenmod.ThresholdProfile, error) {
	// Stub to implement
	return nil, errors.New("Unsupported operation: GetTenantThresholdProfile not implemented")
}

// DeleteTenantThresholdProfile - InMemory implementation of DeleteTenantThresholdProfile
func (tsd *TenantServiceDatastoreInMemory) DeleteTenantThresholdProfile(tenantID string, dataID string) (*tenmod.ThresholdProfile, error) {
	// Stub to implement
	return nil, errors.New("Unsupported operation: DeleteTenantThresholdProfile not implemented")
}

// CreateMonitoredObject - InMemory implementation of CreateMonitoredObject
func (tsd *TenantServiceDatastoreInMemory) CreateMonitoredObject(monitoredObjectReq *tenmod.MonitoredObject) (*tenmod.MonitoredObject, error) {
	// Stub to implement
	return nil, errors.New("Unsupported operation: CreateMonitoredObject not implemented")
}

// UpdateMonitoredObject - InMemory implementation of UpdateMonitoredObject
func (tsd *TenantServiceDatastoreInMemory) UpdateMonitoredObject(monitoredObjectReq *tenmod.MonitoredObject) (*tenmod.MonitoredObject, error) {
	// Stub to implement
	return nil, errors.New("Unsupported operation: UpdateMonitoredObject not implemented")
}

// GetMonitoredObject - InMemory implementation of GetMonitoredObject
func (tsd *TenantServiceDatastoreInMemory) GetMonitoredObject(tenantID string, dataID string) (*tenmod.MonitoredObject, error) {
	// Stub to implement
	return nil, errors.New("Unsupported operation: GetMonitoredObject not implemented")
}

// DeleteMonitoredObject - InMemory implementation of DeleteMonitoredObject
func (tsd *TenantServiceDatastoreInMemory) DeleteMonitoredObject(tenantID string, dataID string) (*tenmod.MonitoredObject, error) {
	// Stub to implement
	return nil, errors.New("Unsupported operation: DeleteMonitoredObject not implemented")
}

// GetAllMonitoredObjects - InMemory implementation of GetAllMonitoredObjects
func (tsd *TenantServiceDatastoreInMemory) GetAllMonitoredObjects(tenantID string) ([]*tenmod.MonitoredObject, error) {
	// Stub to implement
	return nil, errors.New("Unsupported operation: GetAllMonitoredObjects not implemented")
}

// GetMonitoredObjectToDomainMap - InMemory implementation of GetMonitoredObjectToDomainMap
func (tsd *TenantServiceDatastoreInMemory) GetMonitoredObjectToDomainMap(moByDomReq *tenmod.MonitoredObjectCountByDomainRequest) (*tenmod.MonitoredObjectCountByDomainResponse, error) {
	// Stub to implement
	return nil, errors.New("Unsupported operation: GetMonitoredObjectToDomainMap not implemented")
}

// CreateTenantMeta - InMemory implementation of CreateTenantMeta
func (tsd *TenantServiceDatastoreInMemory) CreateTenantMeta(meta *tenmod.Metadata) (*tenmod.Metadata, error) {
	// Stub to implement
	return nil, errors.New("Unsupported operation: CreateTenantMeta not implemented")
}

// UpdateTenantMeta - InMemory implementation of UpdateTenantMeta
func (tsd *TenantServiceDatastoreInMemory) UpdateTenantMeta(meta *tenmod.Metadata) (*tenmod.Metadata, error) {
	// Stub to implement
	return nil, errors.New("Unsupported operation: UpdateTenantMeta not implemented")
}

// DeleteTenantMeta - InMemory implementation of DeleteTenantMeta
func (tsd *TenantServiceDatastoreInMemory) DeleteTenantMeta(tenantID string) (*tenmod.Metadata, error) {
	// Stub to implement
	return nil, errors.New("Unsupported operation: DeleteTenantMeta not implemented")
}

// GetTenantMeta - InMemory implementation of GetTenantMeta
func (tsd *TenantServiceDatastoreInMemory) GetTenantMeta(tenantID string) (*tenmod.Metadata, error) {
	// Stub to implement
	return nil, errors.New("Unsupported operation: GetTenantMeta not implemented")
}

// GetActiveTenantIngestionProfile - InMemory implementation of GetActiveTenantIngestionProfile
func (tsd *TenantServiceDatastoreInMemory) GetActiveTenantIngestionProfile(tenantID string) (*tenmod.IngestionProfile, error) {
	// Stub to implement
	return nil, errors.New("Unsupported operation: GetActiveTenantIngestionProfile not implemented")
}

// GetAllTenantThresholdProfile - InMemory implementation of GetAllTenantThresholdProfile
func (tsd *TenantServiceDatastoreInMemory) GetAllTenantThresholdProfile(tenantID string) ([]*tenmod.ThresholdProfile, error) {
	// Stub to implement
	return nil, errors.New("Unsupported operation: GetAllTenantThresholdProfile not implemented")
}

// BulkInsertMonitoredObjects - InMemory implementation of BulkInsertMonitoredObjects
func (tsd *TenantServiceDatastoreInMemory) BulkInsertMonitoredObjects(tenantID string, value []*tenmod.MonitoredObject) ([]*common.BulkOperationResult, error) {
	// Stub to implement
	return nil, errors.New("BulkInsertMonitoredObjects() not implemented for InMemory DB")
}
