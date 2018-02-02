package inMemory

import (
	"github.com/satori/go.uuid"
	"fmt"
	"errors"
	"time"

	pb "github.com/accedian/adh-gather/gathergrpc"
	ds "github.com/accedian/adh-gather/datastore"
)

// AdminServiceDatastoreInMemory - struct responsible for handling
// database operations for the Admin Service when using local memory
// as the storage option. Useful for tests.
type AdminServiceDatastoreInMemory struct {
	idToAdminUserMap map[string]*pb.AdminUser 
	idToTenantDescMap map[string]*pb.TenantDescriptor 
	ingDictSlice []*pb.IngestionDictionary 
}

// CreateAdminServiceDAO - returns an in-memory implementation of the Admin Service
// datastore.
func CreateAdminServiceDAO() (*AdminServiceDatastoreInMemory, error) {
	res := new(AdminServiceDatastoreInMemory)

	res.idToAdminUserMap = make(map[string]*pb.AdminUser, 0)
	res.idToTenantDescMap = make(map[string]*pb.TenantDescriptor, 0)
	res.ingDictSlice = make([]*pb.IngestionDictionary, 1)

	return res, nil
}

// CreateAdminUser - InMemory implementation of CreateAdminUser
func (memDB *AdminServiceDatastoreInMemory) CreateAdminUser(user *pb.AdminUser) (*pb.AdminUser, error) {
	if len(user.XId) != 0 {
		return nil, fmt.Errorf("%s already exists", ds.AdminUserStr)
	}

	userCopy := *user
	userCopy.XId = uuid.NewV4().String()
	userCopy.XRev = uuid.NewV4().String()
	userCopy.Data.Datatype = string(ds.AdminUserType)
	userCopy.Data.CreatedTimestamp = time.Now().Unix()
	userCopy.Data.LastModifiedTimestamp = userCopy.Data.GetCreatedTimestamp()

	memDB.idToAdminUserMap[userCopy.XId] = &userCopy

	return &userCopy, nil
}

// UpdateAdminUser - InMemory implementation of UpdateAdminUser
func (memDB *AdminServiceDatastoreInMemory) UpdateAdminUser(user *pb.AdminUser) (*pb.AdminUser, error) {
	if len(user.XId) == 0 {
		return nil, fmt.Errorf("%s must have an ID", ds.AdminUserStr)
	}
	if len(user.XRev) == 0 {
		return nil, fmt.Errorf("%s must have a revision", ds.AdminUserStr)
	}

	userCopy := *user
	userCopy.XRev = uuid.NewV4().String()
	userCopy.Data.Datatype = string(ds.AdminUserType)
	userCopy.Data.LastModifiedTimestamp = time.Now().Unix()

	memDB.idToAdminUserMap[userCopy.XId] = &userCopy

	return &userCopy, nil
}

// DeleteAdminUser - InMemory implementation of DeleteAdminUser
func (memDB *AdminServiceDatastoreInMemory) DeleteAdminUser(userID string) (*pb.AdminUser, error) {
	if len(userID) == 0 {
		return nil, fmt.Errorf("%s must provide an ID", ds.AdminUserStr)
	}

	user, ok := memDB.idToAdminUserMap[userID];
    if ok {
		delete(memDB.idToAdminUserMap, userID);
		return user, nil
    }

	return nil, fmt.Errorf("%s not found", ds.AdminUserStr)
}

// GetAdminUser - InMemory implementation of GetAdminUser
func (memDB *AdminServiceDatastoreInMemory) GetAdminUser(userID string) (*pb.AdminUser, error) {
	if len(userID) == 0 {
		return nil, fmt.Errorf("%s must provide an ID", ds.AdminUserStr)
	}

	user, ok := memDB.idToAdminUserMap[userID];
    if ok {
		return user, nil
    }

	return nil, fmt.Errorf("%s not found", ds.AdminUserStr)
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
		return nil, fmt.Errorf("%s already exists", ds.TenantDescriptorStr)
	}

	tenantCopy := *tenantDescriptor
	tenantCopy.XId = uuid.NewV4().String()
	tenantCopy.XRev = uuid.NewV4().String()
	tenantCopy.Data.Datatype = string(ds.TenantDescriptorType)
	tenantCopy.Data.CreatedTimestamp = time.Now().Unix()
	tenantCopy.Data.LastModifiedTimestamp = tenantCopy.Data.GetCreatedTimestamp()

	memDB.idToTenantDescMap[tenantCopy.XId] = &tenantCopy

	return &tenantCopy, nil
}

// UpdateTenantDescriptor - InMemory implementation of UpdateTenantDescriptor
func (memDB *AdminServiceDatastoreInMemory) UpdateTenantDescriptor(tenantDescriptor *pb.TenantDescriptor) (*pb.TenantDescriptor, error) {
	if len(tenantDescriptor.XId) == 0 {
		return nil, fmt.Errorf("%s must have an ID", ds.TenantDescriptorStr)
	}
	if len(tenantDescriptor.XRev) == 0 {
		return nil, fmt.Errorf("%s must have a revision", ds.TenantDescriptorStr)
	}

	tenantCopy := *tenantDescriptor
	tenantCopy.XRev = uuid.NewV4().String()
	tenantCopy.Data.Datatype = string(ds.TenantDescriptorType)
	tenantCopy.Data.LastModifiedTimestamp = time.Now().Unix()

	memDB.idToTenantDescMap[tenantCopy.XId] = &tenantCopy

	return &tenantCopy, nil
}

// DeleteTenant - InMemory implementation of DeleteTenant
func (memDB *AdminServiceDatastoreInMemory) DeleteTenant(tenantID string) (*pb.TenantDescriptor, error) {
	if len(tenantID) == 0 {
		return nil, fmt.Errorf("%s must provide an ID", ds.TenantDescriptorStr)
	}

	tenant, ok := memDB.idToTenantDescMap[tenantID];
    if ok {
		delete(memDB.idToTenantDescMap, tenantID);
		return tenant, nil
    }

	return nil, fmt.Errorf("%s not found", ds.TenantDescriptorStr)
}

// GetTenantDescriptor - InMemory implementation of GetTenantDescriptor
func (memDB *AdminServiceDatastoreInMemory) GetTenantDescriptor(tenantID string) (*pb.TenantDescriptor, error) {
	if len(tenantID) == 0 {
		return nil, fmt.Errorf("%s must provide an ID", ds.TenantDescriptorStr)
	}

	tenant, ok := memDB.idToTenantDescMap[tenantID];
    if ok {
		return tenant, nil
    }

	return nil, fmt.Errorf("%s not found", ds.TenantDescriptorStr)
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
		return nil, fmt.Errorf("Can't create %s, it already exists", ds.IngestionDictionaryStr)
	}
	
	if len(ingDictionary.XId) != 0 {
		return nil, fmt.Errorf("%s already exists", ds.AdminUserStr)
	}

	dictCopy := *ingDictionary
	dictCopy.XId = uuid.NewV4().String()
	dictCopy.XRev = uuid.NewV4().String()
	dictCopy.Data.Datatype = string(ds.IngestionDictionaryType)
	dictCopy.Data.CreatedTimestamp = time.Now().Unix()
	dictCopy.Data.LastModifiedTimestamp = dictCopy.Data.GetCreatedTimestamp()

	memDB.ingDictSlice[0] = &dictCopy

	return &dictCopy, nil
}

// UpdateIngestionDictionary - InMemory implementation of UpdateIngestionDictionary
func (memDB *AdminServiceDatastoreInMemory) UpdateIngestionDictionary(ingDictionary *pb.IngestionDictionary) (*pb.IngestionDictionary, error) {
	if len(ingDictionary.XId) == 0 {
		return nil, fmt.Errorf("%s must have an ID", ds.IngestionDictionaryStr)
	}
	if len(ingDictionary.XRev) == 0 {
		return nil, fmt.Errorf("%s must have a revision", ds.IngestionDictionaryStr)
	}

	dictCopy := *ingDictionary
	dictCopy.XRev = uuid.NewV4().String()
	dictCopy.Data.Datatype = string(ds.TenantDescriptorType)
	dictCopy.Data.LastModifiedTimestamp = time.Now().Unix()

	memDB.ingDictSlice[0] = &dictCopy

	return &dictCopy, nil
}

// DeleteIngestionDictionary - InMemory implementation of DeleteIngestionDictionary
func (memDB *AdminServiceDatastoreInMemory) DeleteIngestionDictionary() (*pb.IngestionDictionary, error) {
	if len(memDB.ingDictSlice) == 0 {
		return nil, fmt.Errorf("%s not found", ds.TenantDescriptorStr)
	}

	res := memDB.ingDictSlice[0]
	memDB.ingDictSlice[0] = nil
	return res, nil
}

// GetIngestionDictionary - InMemory implementation of GetIngestionDictionary
func (memDB *AdminServiceDatastoreInMemory) GetIngestionDictionary() (*pb.IngestionDictionary, error) {
	if len(memDB.ingDictSlice) == 0 {
		return nil, fmt.Errorf("%s not found", ds.TenantDescriptorStr)
	}

	return memDB.ingDictSlice[0], nil
}

// GetTenantIDByAlias - InMemory impl of GetTenantIDByAlias
func (memDB *AdminServiceDatastoreInMemory) GetTenantIDByAlias(name string) (string, error) {
	// Stub to implement
	return "", errors.New("GetTenantIDByName() not implemented for InMemory DB")
}

// AddAdminViews - Adds the admin views (indicies) to the Admin DB.
func (memDB *AdminServiceDatastoreInMemory) AddAdminViews() error {
	// Stub to implement
	return errors.New("GetTenantIDByName() not implemented for InMemory DB")
}

// CreateValidTypes - InMemory implementation of CreateValidTypes
func (memDB *AdminServiceDatastoreInMemory) CreateValidTypes(value *pb.ValidTypes) (*pb.ValidTypes, error) {
	// Stub to implement
	return nil, errors.New("CreateValidTypes() not implemented for InMemory DB")
}

// UpdateValidTypes - InMemory implementation of UpdateValidTypes
func (memDB *AdminServiceDatastoreInMemory) UpdateValidTypes(value *pb.ValidTypes) (*pb.ValidTypes, error) {
	// Stub to implement
	return nil, errors.New("UpdateValidTypes() not implemented for InMemory DB")
}

// GetValidTypes - InMemory implementation of GetValidTypes
func (memDB *AdminServiceDatastoreInMemory) GetValidTypes() (*pb.ValidTypes, error) {
	// Stub to implement
	return nil, errors.New("GetValidTypes() not implemented for InMemory DB")
}

// GetSpecificValidTypes - InMemory implementation of GetSpecificValidTypes
func (memDB *AdminServiceDatastoreInMemory) GetSpecificValidTypes(value *pb.ValidTypesRequest) (*pb.ValidTypesData, error) {
	// Stub to implement
	return nil, errors.New("GetSpecificValidTypes() not implemented for InMemory DB")
}

// CreateDatabase - InMemory implementation of CreateDatabase
func (memDB *AdminServiceDatastoreInMemory) CreateDatabase(dbName string) (ds.Database, error) {
	// Stub to implement
	return nil, errors.New("CreateDatabase() not implemented for InMemory DB")
}