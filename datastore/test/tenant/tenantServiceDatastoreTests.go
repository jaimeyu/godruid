package tenant

import (
	"testing"
	"time"

	"github.com/getlantern/deepcopy"

	"github.com/stretchr/testify/assert"

	ds "github.com/accedian/adh-gather/datastore"
	pb "github.com/accedian/adh-gather/gathergrpc"
	admmod "github.com/accedian/adh-gather/models/admin"
	"github.com/accedian/adh-gather/models/common"
	tenmod "github.com/accedian/adh-gather/models/tenant"
)

// TenantServiceDatastoreTestRunner - object used to run tests for any iplementation
// of the TenantServiceDatastore interface
type TenantServiceDatastoreTestRunner struct {
	tenantDB ds.TenantServiceDatastore
	adminDB  ds.AdminServiceDatastore
}

func InitTestRunner(tdb ds.TenantServiceDatastore, adb ds.AdminServiceDatastore) *TenantServiceDatastoreTestRunner {
	return &TenantServiceDatastoreTestRunner{
		tenantDB: tdb,
		adminDB:  adb,
	}
}

func (runner *TenantServiceDatastoreTestRunner) RunTenantUserCRUD(t *testing.T) {
	const COMPANY1 = "UserCompany"
	const SUBDOMAIN1 = "subdom1"
	const USER1 = "test1"
	const USER2 = "test2"
	const PASS1 = "pass1"
	const PASS2 = "pass2"
	const PASS3 = "pass3"
	const TOKEN1 = "token1"
	const TOKEN2 = "token2"

	// Create a tenant
	data := admmod.Tenant{
		Name:         COMPANY1,
		URLSubdomain: SUBDOMAIN1,
		State:        string(common.UserActive)}
	tenantDescriptor, err := runner.adminDB.CreateTenant(&data)
	assert.Nil(t, err)
	assert.NotNil(t, tenantDescriptor)
	assert.Equal(t, COMPANY1, tenantDescriptor.Name)

	TENANT := ds.GetDataIDFromFullID(tenantDescriptor.ID)

	// Validate that there are currently no records
	tenantUserList, err := runner.tenantDB.GetAllTenantUsers(TENANT)
	assert.Nil(t, err)
	assert.NotNil(t, tenantUserList)

	// Create a record
	tenantUser := pb.TenantUserData{
		Username:        USER1,
		Password:        PASS1,
		OnboardingToken: TOKEN1,
		TenantId:        TENANT,
		State:           pb.UserState_INVITED}
	created, err := runner.tenantDB.CreateTenantUser(&pb.TenantUser{Data: &tenantUser})
	assert.Nil(t, err)
	assert.NotNil(t, created)
	assert.NotNil(t, created.XId)
	assert.NotNil(t, created.XRev)
	assert.NotEmpty(t, created.XId)
	assert.NotEmpty(t, created.XRev)
	assert.Equal(t, string(tenmod.TenantUserType), created.Data.Datatype)
	assert.Equal(t, created.Data.Username, USER1, "Username not the same")
	assert.Equal(t, created.Data.Password, PASS1, "Password not the same")
	assert.Equal(t, created.Data.TenantId, TENANT, "Tenant ID not the same")
	assert.Equal(t, created.Data.OnboardingToken, TOKEN1, "OnboardingToken not the same")
	assert.True(t, created.Data.CreatedTimestamp > 0, "CreatedTimestamp was not set")
	assert.True(t, created.Data.LastModifiedTimestamp > 0, "LastmodifiedTimestamp was not set")

	// Get a record
	fetched, err := runner.tenantDB.GetTenantUser(&pb.TenantUserIdRequest{TenantId: TENANT, UserId: created.XId})
	assert.Nil(t, err)
	assert.Equal(t, created, fetched, "The retrieved record should be the same as the created record")

	time.Sleep(time.Millisecond * 2)

	// Update a record
	updateRecord := pb.TenantUser{}
	deepcopy.Copy(&updateRecord, fetched)
	updateRecord.Data.Password = PASS2
	updated, err := runner.tenantDB.UpdateTenantUser(&updateRecord)
	assert.Nil(t, err)
	assert.NotNil(t, updated)
	assert.Equal(t, updated.XId, fetched.XId)
	assert.NotEqual(t, updated.XRev, fetched.XRev)
	assert.Equal(t, string(tenmod.TenantUserType), updated.Data.Datatype)
	assert.Equal(t, updated.Data.Username, USER1, "Username not the same")
	assert.Equal(t, updated.Data.Password, PASS2, "Password was not updated")
	assert.Equal(t, updated.Data.TenantId, TENANT, "Tenant ID not the same")
	assert.Equal(t, updated.Data.OnboardingToken, TOKEN1, "OnboardingToken not the same")
	assert.Equal(t, updated.Data.CreatedTimestamp, fetched.Data.CreatedTimestamp, "CreatedTimestamp should not be updated")
	assert.True(t, updated.Data.LastModifiedTimestamp > fetched.Data.LastModifiedTimestamp, "LastmodifiedTimestamp was not updated")

	// Add a second record.
	tenantUser2 := pb.TenantUserData{
		Username:        USER2,
		Password:        PASS3,
		OnboardingToken: TOKEN2,
		TenantId:        TENANT,
		State:           pb.UserState_INVITED}
	created2, err := runner.tenantDB.CreateTenantUser(&pb.TenantUser{Data: &tenantUser2})
	assert.Nil(t, err)
	assert.NotNil(t, created2)
	assert.NotNil(t, created2.XId)
	assert.NotNil(t, created2.XRev)
	assert.NotEmpty(t, created2.XId)
	assert.NotEmpty(t, created2.XRev)
	assert.Equal(t, string(tenmod.TenantUserType), created2.Data.Datatype)
	assert.Equal(t, created2.Data.Username, USER2, "Username not the same")
	assert.Equal(t, created2.Data.Password, PASS3, "Password not the same")
	assert.Equal(t, created2.Data.TenantId, TENANT, "Tenant ID not the same")
	assert.Equal(t, created2.Data.OnboardingToken, TOKEN2, "OnboardingToken not the same")
	assert.True(t, created2.Data.CreatedTimestamp > 0, "CreatedTimestamp was not set")
	assert.True(t, created2.Data.LastModifiedTimestamp > 0, "LastmodifiedTimestamp was not set")

	// Get all records
	fetchedList, err := runner.tenantDB.GetAllTenantUsers(TENANT)
	assert.Nil(t, err)
	assert.NotNil(t, fetchedList)
	assert.NotEmpty(t, fetchedList.Data)
	assert.True(t, len(fetchedList.Data) == 2)

	// Delete a record.
	deleted, err := runner.tenantDB.DeleteTenantUser(&pb.TenantUserIdRequest{TenantId: TENANT, UserId: fetched.XId})
	assert.Nil(t, err)
	assert.NotNil(t, deleted)
	assert.NotNil(t, deleted.XId)
	assert.NotNil(t, deleted.XRev)
	assert.NotEmpty(t, deleted.XId)
	assert.NotEmpty(t, deleted.XRev)
	assert.Equal(t, deleted.Data.Username, fetched.Data.Username, "Deleted Username not the same")

	// Get all records - should be 1
	fetchedList, err = runner.tenantDB.GetAllTenantUsers(TENANT)
	assert.Nil(t, err)
	assert.NotNil(t, fetchedList)
	assert.NotEmpty(t, fetchedList.Data)
	assert.True(t, len(fetchedList.Data) == 1)

	// Get a record that does not exist
	dne, err := runner.tenantDB.GetTenantUser(&pb.TenantUserIdRequest{TenantId: TENANT, UserId: deleted.XId})
	assert.NotNil(t, err)
	assert.Nil(t, dne)

	// Delete a record that oes not exist
	deleteDNE, err := runner.tenantDB.DeleteTenantUser(&pb.TenantUserIdRequest{TenantId: TENANT, UserId: deleted.XId})
	assert.NotNil(t, err)
	assert.Nil(t, deleteDNE)

	// Delete the last record
	deleted, err = runner.tenantDB.DeleteTenantUser(&pb.TenantUserIdRequest{TenantId: TENANT, UserId: created2.XId})
	assert.Nil(t, err)
	assert.NotNil(t, deleted)
	assert.NotNil(t, deleted.XId)
	assert.NotNil(t, deleted.XRev)
	assert.NotEmpty(t, deleted.XId)
	assert.NotEmpty(t, deleted.XRev)
	assert.Equal(t, deleted.Data.Username, created2.Data.Username, "Deleted Username not the same")

	// Get all records - should be empty
	fetchedList, err = runner.tenantDB.GetAllTenantUsers(TENANT)
	assert.Nil(t, err)
	assert.NotNil(t, tenantUserList)
}

func (runner *TenantServiceDatastoreTestRunner) RunTenantDomainCRUD(t *testing.T) {

	const COMPANY1 = "DomainCompany"
	const SUBDOMAIN1 = "subdom1"
	const DOM1 = "domain1"
	const DOM2 = "domain2"
	const DOM3 = "domain3"
	const COLOR1 = "color1"
	const COLOR2 = "color2"
	const THRPRF = "ThresholdPrf"

	// Create a tenant
	data := admmod.Tenant{
		Name:         COMPANY1,
		URLSubdomain: SUBDOMAIN1,
		State:        string(common.UserActive)}
	tenantDescriptor, err := runner.adminDB.CreateTenant(&data)
	assert.Nil(t, err)
	assert.NotNil(t, tenantDescriptor)
	assert.Equal(t, COMPANY1, tenantDescriptor.Name)

	TENANT := ds.GetDataIDFromFullID(tenantDescriptor.ID)

	// Validate that there are currently no records
	recList, err := runner.tenantDB.GetAllTenantDomains(TENANT)
	assert.Nil(t, err)
	assert.NotNil(t, recList)

	// Create a record
	tenantDomain := pb.TenantDomainData{
		Name:                DOM1,
		TenantId:            TENANT,
		Color:               COLOR1,
		ThresholdProfileSet: []string{THRPRF}}
	created, err := runner.tenantDB.CreateTenantDomain(&pb.TenantDomain{Data: &tenantDomain})
	assert.Nil(t, err)
	assert.NotNil(t, created)
	assert.NotNil(t, created.XId)
	assert.NotNil(t, created.XRev)
	assert.NotEmpty(t, created.XId)
	assert.NotEmpty(t, created.XRev)
	assert.Equal(t, string(tenmod.TenantDomainType), created.Data.Datatype)
	assert.Equal(t, created.Data.Name, DOM1, "Name not the same")
	assert.Equal(t, created.Data.Color, COLOR1, "Color not the same")
	assert.Equal(t, created.Data.ThresholdProfileSet[0], THRPRF, "Threshold Profile ID not the same")
	assert.Equal(t, created.Data.TenantId, TENANT, "Tenant ID not the same")
	assert.True(t, created.Data.CreatedTimestamp > 0, "CreatedTimestamp was not set")
	assert.True(t, created.Data.LastModifiedTimestamp > 0, "LastmodifiedTimestamp was not set")

	// Get a record
	fetched, err := runner.tenantDB.GetTenantDomain(&pb.TenantDomainIdRequest{TenantId: TENANT, DomainId: created.XId})
	assert.Nil(t, err)
	assert.Equal(t, created, fetched, "The retrieved record should be the same as the created record")

	time.Sleep(time.Millisecond * 2)

	// Update a record
	updateRecord := pb.TenantDomain{}
	deepcopy.Copy(&updateRecord, fetched)
	updateRecord.Data.Color = COLOR2
	updated, err := runner.tenantDB.UpdateTenantDomain(&updateRecord)
	assert.Nil(t, err)
	assert.NotNil(t, updated)
	assert.Equal(t, updated.XId, fetched.XId)
	assert.NotEqual(t, updated.XRev, fetched.XRev)
	assert.Equal(t, string(tenmod.TenantDomainType), updated.Data.Datatype)
	assert.Equal(t, updated.Data.Name, DOM1, "Name not the same")
	assert.Equal(t, updated.Data.Color, COLOR2, "Password was not updated")
	assert.Equal(t, updated.Data.TenantId, TENANT, "Tenant ID not the same")
	assert.Equal(t, updated.Data.ThresholdProfileSet[0], THRPRF, "Threshold Profile ID not the same")
	assert.Equal(t, updated.Data.CreatedTimestamp, fetched.Data.CreatedTimestamp, "CreatedTimestamp should not be updated")
	assert.True(t, updated.Data.LastModifiedTimestamp > fetched.Data.LastModifiedTimestamp, "LastmodifiedTimestamp was not updated")

	// Add a second record.
	tenantDomain2 := pb.TenantDomainData{
		Name:     DOM2,
		TenantId: TENANT,
		Color:    COLOR1}
	created2, err := runner.tenantDB.CreateTenantDomain(&pb.TenantDomain{Data: &tenantDomain2})
	assert.Nil(t, err)
	assert.NotNil(t, created2)
	assert.NotNil(t, created2.XId)
	assert.NotNil(t, created2.XRev)
	assert.NotEmpty(t, created2.XId)
	assert.NotEmpty(t, created2.XRev)
	assert.Equal(t, string(tenmod.TenantDomainType), created2.Data.Datatype)
	assert.Equal(t, created2.Data.Name, DOM2, "Name not the same")
	assert.Equal(t, created2.Data.Color, COLOR1, "Password not the same")
	assert.Equal(t, created2.Data.TenantId, TENANT, "Tenant ID not the same")
	assert.True(t, len(created2.Data.ThresholdProfileSet) == 0, "Should not be a Threshold Profile ID")
	assert.True(t, created2.Data.CreatedTimestamp > 0, "CreatedTimestamp was not set")
	assert.True(t, created2.Data.LastModifiedTimestamp > 0, "LastmodifiedTimestamp was not set")

	// Get all records
	fetchedList, err := runner.tenantDB.GetAllTenantDomains(TENANT)
	assert.Nil(t, err)
	assert.NotNil(t, fetchedList)
	assert.NotEmpty(t, fetchedList.Data)
	assert.True(t, len(fetchedList.Data) == 2)

	// Delete a record.
	deleted, err := runner.tenantDB.DeleteTenantDomain(&pb.TenantDomainIdRequest{TenantId: TENANT, DomainId: fetched.XId})
	assert.Nil(t, err)
	assert.NotNil(t, deleted)
	assert.NotNil(t, deleted.XId)
	assert.NotNil(t, deleted.XRev)
	assert.NotEmpty(t, deleted.XId)
	assert.NotEmpty(t, deleted.XRev)
	assert.Equal(t, deleted.Data.Name, fetched.Data.Name, "Deleted name not the same")

	// Get all records - should be 1
	fetchedList, err = runner.tenantDB.GetAllTenantDomains(TENANT)
	assert.Nil(t, err)
	assert.NotNil(t, fetchedList)
	assert.NotEmpty(t, fetchedList.Data)
	assert.True(t, len(fetchedList.Data) == 1)

	// Get a record that does not exist
	dne, err := runner.tenantDB.GetTenantDomain(&pb.TenantDomainIdRequest{TenantId: TENANT, DomainId: deleted.XId})
	assert.NotNil(t, err)
	assert.Nil(t, dne)

	// Delete a record that oes not exist
	deleteDNE, err := runner.tenantDB.DeleteTenantDomain(&pb.TenantDomainIdRequest{TenantId: TENANT, DomainId: deleted.XId})
	assert.NotNil(t, err)
	assert.Nil(t, deleteDNE)

	// Delete the last record
	deleted, err = runner.tenantDB.DeleteTenantDomain(&pb.TenantDomainIdRequest{TenantId: TENANT, DomainId: created2.XId})
	assert.Nil(t, err)
	assert.NotNil(t, deleted)
	assert.NotNil(t, deleted.XId)
	assert.NotNil(t, deleted.XRev)
	assert.NotEmpty(t, deleted.XId)
	assert.NotEmpty(t, deleted.XRev)
	assert.Equal(t, deleted.Data.Name, created2.Data.Name, "Deleted name not the same")

	// Get all records - should be empty
	fetchedList, err = runner.tenantDB.GetAllTenantDomains(TENANT)
	assert.Nil(t, err)
	assert.NotNil(t, recList)
}
