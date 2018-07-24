package inMemory

import (
	"errors"
	"fmt"
	"strings"

	"github.com/satori/go.uuid"

	ds "github.com/accedian/adh-gather/datastore"
	admmod "github.com/accedian/adh-gather/models/admin"
	"github.com/getlantern/deepcopy"
)

// AdminServiceDatastoreInMemory - struct responsible for handling
// database operations for the Admin Service when using local memory
// as the storage option. Useful for tests.
type AdminServiceDatastoreInMemory struct {
	idToAdminUserMap  map[string]*admmod.User
	idToTenantDescMap map[string]*admmod.Tenant
	ingDictSlice      []*admmod.IngestionDictionary
	validTypeSlice    []*admmod.ValidTypes
}

// CreateAdminServiceDAO - returns an in-memory implementation of the Admin Service
// datastore.
func CreateAdminServiceDAO() (*AdminServiceDatastoreInMemory, error) {
	res := new(AdminServiceDatastoreInMemory)

	res.idToAdminUserMap = make(map[string]*admmod.User, 0)
	res.idToTenantDescMap = make(map[string]*admmod.Tenant, 0)
	res.ingDictSlice = make([]*admmod.IngestionDictionary, 1)
	res.validTypeSlice = make([]*admmod.ValidTypes, 1)

	return res, nil
}

// CreateAdminUser - InMemory implementation of CreateAdminUser
func (memDB *AdminServiceDatastoreInMemory) CreateAdminUser(user *admmod.User) (*admmod.User, error) {
	if len(user.ID) != 0 {
		return nil, fmt.Errorf("%s already exists", admmod.AdminUserStr)
	}

	userCopy := admmod.User{}
	deepcopy.Copy(&userCopy, user)
	userCopy.ID = uuid.NewV4().String()
	userCopy.REV = uuid.NewV4().String()
	userCopy.Datatype = string(admmod.AdminUserType)
	userCopy.CreatedTimestamp = ds.MakeTimestamp()
	userCopy.LastModifiedTimestamp = userCopy.CreatedTimestamp

	memDB.idToAdminUserMap[userCopy.ID] = &userCopy

	return &userCopy, nil
}

// UpdateAdminUser - InMemory implementation of UpdateAdminUser
func (memDB *AdminServiceDatastoreInMemory) UpdateAdminUser(user *admmod.User) (*admmod.User, error) {
	if len(user.ID) == 0 {
		return nil, fmt.Errorf("%s must have an ID", admmod.AdminUserStr)
	}
	if len(user.REV) == 0 {
		return nil, fmt.Errorf("%s must have a revision", admmod.AdminUserStr)
	}

	userCopy := admmod.User{}
	deepcopy.Copy(&userCopy, user)
	userCopy.REV = uuid.NewV4().String()
	userCopy.Datatype = string(admmod.AdminUserType)
	userCopy.LastModifiedTimestamp = ds.MakeTimestamp()

	memDB.idToAdminUserMap[userCopy.ID] = &userCopy

	return &userCopy, nil
}

// DeleteAdminUser - InMemory implementation of DeleteAdminUser
func (memDB *AdminServiceDatastoreInMemory) DeleteAdminUser(userID string) (*admmod.User, error) {
	if len(userID) == 0 {
		return nil, fmt.Errorf("%s must provide an ID", admmod.AdminUserStr)
	}

	user, ok := memDB.idToAdminUserMap[userID]
	if ok {
		delete(memDB.idToAdminUserMap, userID)
		return user, nil
	}

	return nil, fmt.Errorf("%s not found", admmod.AdminUserStr)
}

// GetAdminUser - InMemory implementation of GetAdminUser
func (memDB *AdminServiceDatastoreInMemory) GetAdminUser(userID string) (*admmod.User, error) {
	if len(userID) == 0 {
		return nil, fmt.Errorf("%s must provide an ID", admmod.AdminUserStr)
	}

	user, ok := memDB.idToAdminUserMap[userID]
	if ok {
		return user, nil
	}

	return nil, fmt.Errorf("%s not found", admmod.AdminUserStr)
}

// GetAllAdminUsers - InMemory implementation of GetAllAdminUsers
func (memDB *AdminServiceDatastoreInMemory) GetAllAdminUsers() ([]*admmod.User, error) {
	adminUserList := []*admmod.User{}

	for _, user := range memDB.idToAdminUserMap {
		adminUserList = append(adminUserList, user)
	}

	return adminUserList, nil
}

// CreateTenant - InMemory implementation of CreateTenant
func (memDB *AdminServiceDatastoreInMemory) CreateTenant(tenantDescriptor *admmod.Tenant) (*admmod.Tenant, error) {
	if len(tenantDescriptor.ID) != 0 {
		return nil, fmt.Errorf("%s already exists", admmod.TenantStr)
	}

	tenantCopy := admmod.Tenant{}
	deepcopy.Copy(&tenantCopy, tenantDescriptor)
	tenantCopy.ID = uuid.NewV4().String()
	tenantCopy.REV = uuid.NewV4().String()
	tenantCopy.Datatype = string(admmod.TenantType)
	tenantCopy.CreatedTimestamp = ds.MakeTimestamp()
	tenantCopy.LastModifiedTimestamp = tenantCopy.CreatedTimestamp

	memDB.idToTenantDescMap[tenantCopy.ID] = &tenantCopy

	return &tenantCopy, nil
}

// UpdateTenantDescriptor - InMemory implementation of UpdateTenantDescriptor
func (memDB *AdminServiceDatastoreInMemory) UpdateTenantDescriptor(tenantDescriptor *admmod.Tenant) (*admmod.Tenant, error) {
	if len(tenantDescriptor.ID) == 0 {
		return nil, fmt.Errorf("%s must have an ID", admmod.TenantStr)
	}
	if len(tenantDescriptor.REV) == 0 {
		return nil, fmt.Errorf("%s must have a revision", admmod.TenantStr)
	}

	tenantCopy := admmod.Tenant{}
	deepcopy.Copy(&tenantCopy, tenantDescriptor)
	tenantCopy.REV = uuid.NewV4().String()
	tenantCopy.Datatype = string(admmod.TenantType)
	tenantCopy.LastModifiedTimestamp = ds.MakeTimestamp()

	memDB.idToTenantDescMap[tenantCopy.ID] = &tenantCopy

	return &tenantCopy, nil
}

// DeleteTenant - InMemory implementation of DeleteTenant
func (memDB *AdminServiceDatastoreInMemory) DeleteTenant(tenantID string) (*admmod.Tenant, error) {
	if len(tenantID) == 0 {
		return nil, fmt.Errorf("%s must provide an ID", admmod.TenantStr)
	}

	tenant, ok := memDB.idToTenantDescMap[tenantID]
	if ok {
		delete(memDB.idToTenantDescMap, tenantID)
		return tenant, nil
	}

	return nil, fmt.Errorf("%s not found", admmod.TenantStr)
}

// GetTenantDescriptor - InMemory implementation of GetTenantDescriptor
func (memDB *AdminServiceDatastoreInMemory) GetTenantDescriptor(tenantID string) (*admmod.Tenant, error) {
	if len(tenantID) == 0 {
		return nil, fmt.Errorf("%s must provide an ID", admmod.TenantStr)
	}

	tenant, ok := memDB.idToTenantDescMap[tenantID]
	if ok {
		return tenant, nil
	}

	return nil, fmt.Errorf("%s not found", admmod.TenantStr)
}

// GetAllTenantDescriptors - InMemory implementation of GetAllTenantDescriptors
func (memDB *AdminServiceDatastoreInMemory) GetAllTenantDescriptors() ([]*admmod.Tenant, error) {
	tenantDescList := []*admmod.Tenant{}

	for _, tenant := range memDB.idToTenantDescMap {
		tenantDescList = append(tenantDescList, tenant)
	}

	return tenantDescList, nil
}

// CreateIngestionDictionary - InMemory implementation of CreateIngestionDictionary
func (memDB *AdminServiceDatastoreInMemory) CreateIngestionDictionary(ingDictionary *admmod.IngestionDictionary) (*admmod.IngestionDictionary, error) {
	// Make sure one does not already exist
	existing, _ := memDB.GetIngestionDictionary()
	if existing != nil {
		return nil, fmt.Errorf("Can't create %s, it already exists", admmod.IngestionDictionaryStr)
	}

	if len(ingDictionary.ID) != 0 {
		return nil, fmt.Errorf("%s already exists", admmod.IngestionDictionaryStr)
	}

	dictCopy := admmod.IngestionDictionary{}
	deepcopy.Copy(&dictCopy, ingDictionary)
	dictCopy.ID = uuid.NewV4().String()
	dictCopy.REV = uuid.NewV4().String()
	dictCopy.Datatype = string(admmod.IngestionDictionaryType)
	dictCopy.CreatedTimestamp = ds.MakeTimestamp()
	dictCopy.LastModifiedTimestamp = dictCopy.CreatedTimestamp

	memDB.ingDictSlice[0] = &dictCopy

	return &dictCopy, nil
}

// UpdateIngestionDictionary - InMemory implementation of UpdateIngestionDictionary
func (memDB *AdminServiceDatastoreInMemory) UpdateIngestionDictionary(ingDictionary *admmod.IngestionDictionary) (*admmod.IngestionDictionary, error) {
	if len(ingDictionary.ID) == 0 {
		return nil, fmt.Errorf("%s must have an ID", admmod.IngestionDictionaryStr)
	}
	if len(ingDictionary.REV) == 0 {
		return nil, fmt.Errorf("%s must have a revision", admmod.IngestionDictionaryStr)
	}

	dictCopy := admmod.IngestionDictionary{}
	deepcopy.Copy(&dictCopy, ingDictionary)
	dictCopy.REV = uuid.NewV4().String()
	dictCopy.Datatype = string(admmod.IngestionDictionaryType)
	dictCopy.LastModifiedTimestamp = ds.MakeTimestamp()

	memDB.ingDictSlice[0] = &dictCopy

	return &dictCopy, nil
}

// DeleteIngestionDictionary - InMemory implementation of DeleteIngestionDictionary
func (memDB *AdminServiceDatastoreInMemory) DeleteIngestionDictionary() (*admmod.IngestionDictionary, error) {
	existing, err := memDB.GetIngestionDictionary()
	if err != nil {
		return nil, err
	}

	memDB.ingDictSlice[0] = nil
	return existing, nil
}

// GetIngestionDictionary - InMemory implementation of GetIngestionDictionary
func (memDB *AdminServiceDatastoreInMemory) GetIngestionDictionary() (*admmod.IngestionDictionary, error) {
	if len(memDB.ingDictSlice) == 0 || memDB.ingDictSlice[0] == nil {
		return nil, fmt.Errorf("%s not found", admmod.IngestionDictionaryStr)
	}

	return memDB.ingDictSlice[0], nil
}

// GetTenantIDByAlias - InMemory impl of GetTenantIDByAlias
func (memDB *AdminServiceDatastoreInMemory) GetTenantIDByAlias(name string) (string, error) {
	for _, value := range memDB.idToTenantDescMap {
		if strings.ToLower(value.Name) == strings.ToLower(name) {
			return value.ID, nil
		}
	}

	return "", fmt.Errorf("No tenant found for name %s", name)
}

// AddAdminViews - Adds the admin views (indicies) to the Admin DB.
func (memDB *AdminServiceDatastoreInMemory) AddAdminViews() error {
	// Stub to implement
	return errors.New("GetTenantIDByName() not implemented for InMemory DB")
}

// CreateValidTypes - InMemory implementation of CreateValidTypes
func (memDB *AdminServiceDatastoreInMemory) CreateValidTypes(value *admmod.ValidTypes) (*admmod.ValidTypes, error) {
	// Make sure one does not already exist
	existing, _ := memDB.GetValidTypes()
	if existing != nil {
		return nil, fmt.Errorf("Can't create %s, it already exists", admmod.ValidTypesStr)
	}

	if len(value.ID) != 0 {
		return nil, fmt.Errorf("%s already exists", admmod.ValidTypesStr)
	}

	vtCopy := admmod.ValidTypes{}
	deepcopy.Copy(&vtCopy, value)
	vtCopy.ID = uuid.NewV4().String()
	vtCopy.REV = uuid.NewV4().String()
	vtCopy.CreatedTimestamp = ds.MakeTimestamp()
	vtCopy.LastModifiedTimestamp = vtCopy.CreatedTimestamp

	memDB.validTypeSlice[0] = &vtCopy

	return &vtCopy, nil
}

// UpdateValidTypes - InMemory implementation of UpdateValidTypes
func (memDB *AdminServiceDatastoreInMemory) UpdateValidTypes(value *admmod.ValidTypes) (*admmod.ValidTypes, error) {
	if len(value.ID) == 0 {
		return nil, fmt.Errorf("%s must have an ID", admmod.ValidTypesStr)
	}
	if len(value.REV) == 0 {
		return nil, fmt.Errorf("%s must have a revision", admmod.ValidTypesStr)
	}

	vtCopy := admmod.ValidTypes{}
	deepcopy.Copy(&vtCopy, value)
	vtCopy.REV = uuid.NewV4().String()
	vtCopy.LastModifiedTimestamp = ds.MakeTimestamp()

	memDB.validTypeSlice[0] = &vtCopy

	return &vtCopy, nil
}

// GetValidTypes - InMemory implementation of GetValidTypes
func (memDB *AdminServiceDatastoreInMemory) GetValidTypes() (*admmod.ValidTypes, error) {
	if len(memDB.validTypeSlice) == 0 || memDB.validTypeSlice[0] == nil {
		return nil, fmt.Errorf("%s not found", admmod.ValidTypesStr)
	}

	return memDB.validTypeSlice[0], nil
}

// GetSpecificValidTypes - InMemory implementation of GetSpecificValidTypes
func (memDB *AdminServiceDatastoreInMemory) GetSpecificValidTypes(value *admmod.ValidTypesRequest) (*admmod.ValidTypes, error) {
	if len(memDB.validTypeSlice) == 0 || memDB.validTypeSlice[0] == nil {
		return nil, fmt.Errorf("%s not found", admmod.ValidTypesStr)
	}

	vtCopy := admmod.ValidTypes{}
	deepcopy.Copy(&vtCopy, memDB.validTypeSlice[0])

	if !value.MonitoredObjectDeviceTypes {
		vtCopy.MonitoredObjectDeviceTypes = nil
	}
	if !value.MonitoredObjectTypes {
		vtCopy.MonitoredObjectTypes = nil
	}

	return &vtCopy, nil
}

// DeleteValidTypes - InMemory implementation of DeleteValidTypes
func (memDB *AdminServiceDatastoreInMemory) DeleteValidTypes() (*admmod.ValidTypes, error) {
	existing, err := memDB.GetValidTypes()
	if err != nil {
		return nil, err
	}

	memDB.validTypeSlice[0] = nil
	return existing, nil
}

// CreateDatabase - InMemory implementation of CreateDatabase
func (memDB *AdminServiceDatastoreInMemory) CreateDatabase(dbName string) (ds.Database, error) {
	// Nothing to do for in memory
	return nil, nil
}

// DeleteDatabase - InMemory implementation of DeleteDatabase
func (memDB *AdminServiceDatastoreInMemory) DeleteDatabase(dbName string) error {
	// Nothing to do for in memory
	return nil
}
