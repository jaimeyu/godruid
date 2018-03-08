package inMemory

import (
	"errors"
	"fmt"
	"strings"

	"github.com/satori/go.uuid"

	ds "github.com/accedian/adh-gather/datastore"
	pb "github.com/accedian/adh-gather/gathergrpc"
	admmod "github.com/accedian/adh-gather/models/admin"
	"github.com/getlantern/deepcopy"
)

// AdminServiceDatastoreInMemory - struct responsible for handling
// database operations for the Admin Service when using local memory
// as the storage option. Useful for tests.
type AdminServiceDatastoreInMemory struct {
	idToAdminUserMap  map[string]*pb.AdminUser
	idToTenantDescMap map[string]*pb.TenantDescriptor
	ingDictSlice      []*pb.IngestionDictionary
	validTypeSlice    []*pb.ValidTypes
}

// CreateAdminServiceDAO - returns an in-memory implementation of the Admin Service
// datastore.
func CreateAdminServiceDAO() (*AdminServiceDatastoreInMemory, error) {
	res := new(AdminServiceDatastoreInMemory)

	res.idToAdminUserMap = make(map[string]*pb.AdminUser, 0)
	res.idToTenantDescMap = make(map[string]*pb.TenantDescriptor, 0)
	res.ingDictSlice = make([]*pb.IngestionDictionary, 1)
	res.validTypeSlice = make([]*pb.ValidTypes, 1)

	return res, nil
}

// CreateAdminUser - InMemory implementation of CreateAdminUser
func (memDB *AdminServiceDatastoreInMemory) CreateAdminUser(user *pb.AdminUser) (*pb.AdminUser, error) {
	if len(user.XId) != 0 {
		return nil, fmt.Errorf("%s already exists", admmod.AdminUserStr)
	}

	userCopy := pb.AdminUser{}
	deepcopy.Copy(&userCopy, user)
	userCopy.XId = uuid.NewV4().String()
	userCopy.XRev = uuid.NewV4().String()
	userCopy.Data.Datatype = string(admmod.AdminUserType)
	userCopy.Data.CreatedTimestamp = ds.MakeTimestamp()
	userCopy.Data.LastModifiedTimestamp = userCopy.Data.GetCreatedTimestamp()

	memDB.idToAdminUserMap[userCopy.XId] = &userCopy

	return &userCopy, nil
}

// UpdateAdminUser - InMemory implementation of UpdateAdminUser
func (memDB *AdminServiceDatastoreInMemory) UpdateAdminUser(user *pb.AdminUser) (*pb.AdminUser, error) {
	if len(user.XId) == 0 {
		return nil, fmt.Errorf("%s must have an ID", admmod.AdminUserStr)
	}
	if len(user.XRev) == 0 {
		return nil, fmt.Errorf("%s must have a revision", admmod.AdminUserStr)
	}

	userCopy := pb.AdminUser{}
	deepcopy.Copy(&userCopy, user)
	userCopy.XRev = uuid.NewV4().String()
	userCopy.Data.Datatype = string(admmod.AdminUserType)
	userCopy.Data.LastModifiedTimestamp = ds.MakeTimestamp()

	memDB.idToAdminUserMap[userCopy.XId] = &userCopy

	return &userCopy, nil
}

// DeleteAdminUser - InMemory implementation of DeleteAdminUser
func (memDB *AdminServiceDatastoreInMemory) DeleteAdminUser(userID string) (*pb.AdminUser, error) {
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
func (memDB *AdminServiceDatastoreInMemory) GetAdminUser(userID string) (*pb.AdminUser, error) {
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
func (memDB *AdminServiceDatastoreInMemory) GetAllAdminUsers() (*pb.AdminUserList, error) {
	adminUserList := pb.AdminUserList{}
	adminUserList.Data = make([]*pb.AdminUser, 0)

	for _, user := range memDB.idToAdminUserMap {
		adminUserList.Data = append(adminUserList.Data, user)
	}

	return &adminUserList, nil
}

// CreateTenant - InMemory implementation of CreateTenant
func (memDB *AdminServiceDatastoreInMemory) CreateTenant(tenantDescriptor *pb.TenantDescriptor) (*pb.TenantDescriptor, error) {
	if len(tenantDescriptor.XId) != 0 {
		return nil, fmt.Errorf("%s already exists", admmod.TenantStr)
	}

	tenantCopy := pb.TenantDescriptor{}
	deepcopy.Copy(&tenantCopy, tenantDescriptor)
	tenantCopy.XId = uuid.NewV4().String()
	tenantCopy.XRev = uuid.NewV4().String()
	tenantCopy.Data.Datatype = string(admmod.TenantType)
	tenantCopy.Data.CreatedTimestamp = ds.MakeTimestamp()
	tenantCopy.Data.LastModifiedTimestamp = tenantCopy.Data.GetCreatedTimestamp()

	memDB.idToTenantDescMap[tenantCopy.XId] = &tenantCopy

	return &tenantCopy, nil
}

// UpdateTenantDescriptor - InMemory implementation of UpdateTenantDescriptor
func (memDB *AdminServiceDatastoreInMemory) UpdateTenantDescriptor(tenantDescriptor *pb.TenantDescriptor) (*pb.TenantDescriptor, error) {
	if len(tenantDescriptor.XId) == 0 {
		return nil, fmt.Errorf("%s must have an ID", admmod.TenantStr)
	}
	if len(tenantDescriptor.XRev) == 0 {
		return nil, fmt.Errorf("%s must have a revision", admmod.TenantStr)
	}

	tenantCopy := pb.TenantDescriptor{}
	deepcopy.Copy(&tenantCopy, tenantDescriptor)
	tenantCopy.XRev = uuid.NewV4().String()
	tenantCopy.Data.Datatype = string(admmod.TenantType)
	tenantCopy.Data.LastModifiedTimestamp = ds.MakeTimestamp()

	memDB.idToTenantDescMap[tenantCopy.XId] = &tenantCopy

	return &tenantCopy, nil
}

// DeleteTenant - InMemory implementation of DeleteTenant
func (memDB *AdminServiceDatastoreInMemory) DeleteTenant(tenantID string) (*pb.TenantDescriptor, error) {
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
func (memDB *AdminServiceDatastoreInMemory) GetTenantDescriptor(tenantID string) (*pb.TenantDescriptor, error) {
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
func (memDB *AdminServiceDatastoreInMemory) GetAllTenantDescriptors() (*pb.TenantDescriptorList, error) {
	tenantDescList := pb.TenantDescriptorList{}
	tenantDescList.Data = make([]*pb.TenantDescriptor, 0)

	for _, tenant := range memDB.idToTenantDescMap {
		tenantDescList.Data = append(tenantDescList.Data, tenant)
	}

	return &tenantDescList, nil
}

// CreateIngestionDictionary - InMemory implementation of CreateIngestionDictionary
func (memDB *AdminServiceDatastoreInMemory) CreateIngestionDictionary(ingDictionary *pb.IngestionDictionary) (*pb.IngestionDictionary, error) {
	// Make sure one does not already exist
	existing, _ := memDB.GetIngestionDictionary()
	if existing != nil {
		return nil, fmt.Errorf("Can't create %s, it already exists", admmod.IngestionDictionaryStr)
	}

	if len(ingDictionary.XId) != 0 {
		return nil, fmt.Errorf("%s already exists", admmod.IngestionDictionaryStr)
	}

	dictCopy := pb.IngestionDictionary{}
	deepcopy.Copy(&dictCopy, ingDictionary)
	dictCopy.XId = uuid.NewV4().String()
	dictCopy.XRev = uuid.NewV4().String()
	dictCopy.Data.Datatype = string(admmod.IngestionDictionaryType)
	dictCopy.Data.CreatedTimestamp = ds.MakeTimestamp()
	dictCopy.Data.LastModifiedTimestamp = dictCopy.Data.GetCreatedTimestamp()

	memDB.ingDictSlice[0] = &dictCopy

	return &dictCopy, nil
}

// UpdateIngestionDictionary - InMemory implementation of UpdateIngestionDictionary
func (memDB *AdminServiceDatastoreInMemory) UpdateIngestionDictionary(ingDictionary *pb.IngestionDictionary) (*pb.IngestionDictionary, error) {
	if len(ingDictionary.XId) == 0 {
		return nil, fmt.Errorf("%s must have an ID", admmod.IngestionDictionaryStr)
	}
	if len(ingDictionary.XRev) == 0 {
		return nil, fmt.Errorf("%s must have a revision", admmod.IngestionDictionaryStr)
	}

	dictCopy := pb.IngestionDictionary{}
	deepcopy.Copy(&dictCopy, ingDictionary)
	dictCopy.XRev = uuid.NewV4().String()
	dictCopy.Data.Datatype = string(admmod.IngestionDictionaryType)
	dictCopy.Data.LastModifiedTimestamp = ds.MakeTimestamp()

	memDB.ingDictSlice[0] = &dictCopy

	return &dictCopy, nil
}

// DeleteIngestionDictionary - InMemory implementation of DeleteIngestionDictionary
func (memDB *AdminServiceDatastoreInMemory) DeleteIngestionDictionary() (*pb.IngestionDictionary, error) {
	existing, err := memDB.GetIngestionDictionary()
	if err != nil {
		return nil, err
	}

	memDB.ingDictSlice[0] = nil
	return existing, nil
}

// GetIngestionDictionary - InMemory implementation of GetIngestionDictionary
func (memDB *AdminServiceDatastoreInMemory) GetIngestionDictionary() (*pb.IngestionDictionary, error) {
	if len(memDB.ingDictSlice) == 0 || memDB.ingDictSlice[0] == nil {
		return nil, fmt.Errorf("%s not found", admmod.IngestionDictionaryStr)
	}

	return memDB.ingDictSlice[0], nil
}

// GetTenantIDByAlias - InMemory impl of GetTenantIDByAlias
func (memDB *AdminServiceDatastoreInMemory) GetTenantIDByAlias(name string) (string, error) {
	for _, value := range memDB.idToTenantDescMap {
		if strings.ToLower(value.Data.Name) == strings.ToLower(name) {
			return value.XId, nil
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
func (memDB *AdminServiceDatastoreInMemory) CreateValidTypes(value *pb.ValidTypes) (*pb.ValidTypes, error) {
	// Make sure one does not already exist
	existing, _ := memDB.GetValidTypes()
	if existing != nil {
		return nil, fmt.Errorf("Can't create %s, it already exists", admmod.ValidTypesStr)
	}

	if len(value.XId) != 0 {
		return nil, fmt.Errorf("%s already exists", admmod.ValidTypesStr)
	}

	vtCopy := pb.ValidTypes{}
	deepcopy.Copy(&vtCopy, value)
	vtCopy.XId = uuid.NewV4().String()
	vtCopy.XRev = uuid.NewV4().String()
	vtCopy.Data.CreatedTimestamp = ds.MakeTimestamp()
	vtCopy.Data.LastModifiedTimestamp = vtCopy.Data.GetCreatedTimestamp()

	memDB.validTypeSlice[0] = &vtCopy

	return &vtCopy, nil
}

// UpdateValidTypes - InMemory implementation of UpdateValidTypes
func (memDB *AdminServiceDatastoreInMemory) UpdateValidTypes(value *pb.ValidTypes) (*pb.ValidTypes, error) {
	if len(value.XId) == 0 {
		return nil, fmt.Errorf("%s must have an ID", admmod.ValidTypesStr)
	}
	if len(value.XRev) == 0 {
		return nil, fmt.Errorf("%s must have a revision", admmod.ValidTypesStr)
	}

	vtCopy := pb.ValidTypes{}
	deepcopy.Copy(&vtCopy, value)
	vtCopy.XRev = uuid.NewV4().String()
	vtCopy.Data.LastModifiedTimestamp = ds.MakeTimestamp()

	memDB.validTypeSlice[0] = &vtCopy

	return &vtCopy, nil
}

// GetValidTypes - InMemory implementation of GetValidTypes
func (memDB *AdminServiceDatastoreInMemory) GetValidTypes() (*pb.ValidTypes, error) {
	if len(memDB.validTypeSlice) == 0 || memDB.validTypeSlice[0] == nil {
		return nil, fmt.Errorf("%s not found", admmod.ValidTypesStr)
	}

	return memDB.validTypeSlice[0], nil
}

// GetSpecificValidTypes - InMemory implementation of GetSpecificValidTypes
func (memDB *AdminServiceDatastoreInMemory) GetSpecificValidTypes(value *pb.ValidTypesRequest) (*pb.ValidTypesData, error) {
	if len(memDB.validTypeSlice) == 0 || memDB.validTypeSlice[0] == nil {
		return nil, fmt.Errorf("%s not found", admmod.ValidTypesStr)
	}

	vtCopy := pb.ValidTypes{}
	deepcopy.Copy(&vtCopy, memDB.validTypeSlice[0])

	if !value.MonitoredObjectDeviceTypes {
		vtCopy.Data.MonitoredObjectDeviceTypes = nil
	}
	if !value.MonitoredObjectTypes {
		vtCopy.Data.MonitoredObjectTypes = nil
	}

	return vtCopy.Data, nil
}

// DeleteValidTypes - InMemory implementation of DeleteValidTypes
func (memDB *AdminServiceDatastoreInMemory) DeleteValidTypes() (*pb.ValidTypes, error) {
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
