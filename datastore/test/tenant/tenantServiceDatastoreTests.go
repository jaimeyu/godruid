package tenant

import (
	"testing"
	"time"

	"github.com/getlantern/deepcopy"

	"github.com/stretchr/testify/assert"

	ds "github.com/accedian/adh-gather/datastore"
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
	tenantUser := tenmod.User{
		Username:        USER1,
		Password:        PASS1,
		OnboardingToken: TOKEN1,
		TenantID:        TENANT,
		State:           string(common.UserActive)}
	created, err := runner.tenantDB.CreateTenantUser(&tenantUser)
	assert.Nil(t, err)
	assert.NotNil(t, created)
	assert.NotEmpty(t, created.ID)
	assert.NotEmpty(t, created.REV)
	assert.Equal(t, string(tenmod.TenantUserType), created.Datatype)
	assert.Equal(t, created.Username, USER1, "Username not the same")
	assert.Equal(t, created.Password, PASS1, "Password not the same")
	assert.Equal(t, created.TenantID, TENANT, "Tenant ID not the same")
	assert.Equal(t, created.OnboardingToken, TOKEN1, "OnboardingToken not the same")
	assert.True(t, created.CreatedTimestamp > 0, "CreatedTimestamp was not set")
	assert.True(t, created.LastModifiedTimestamp > 0, "LastmodifiedTimestamp was not set")

	// Get a record
	fetched, err := runner.tenantDB.GetTenantUser(TENANT, created.ID)
	assert.Nil(t, err)
	assert.Equal(t, created, fetched, "The retrieved record should be the same as the created record")

	time.Sleep(time.Millisecond * 2)

	// Update a record
	updateRecord := tenmod.User{}
	deepcopy.Copy(&updateRecord, fetched)
	updateRecord.Password = PASS2
	updated, err := runner.tenantDB.UpdateTenantUser(&updateRecord)
	assert.Nil(t, err)
	assert.NotNil(t, updated)
	assert.Equal(t, updated.ID, fetched.ID)
	assert.NotEqual(t, updated.REV, fetched.REV)
	assert.Equal(t, string(tenmod.TenantUserType), updated.Datatype)
	assert.Equal(t, updated.Username, USER1, "Username not the same")
	assert.Equal(t, updated.Password, PASS2, "Password was not updated")
	assert.Equal(t, updated.TenantID, TENANT, "Tenant ID not the same")
	assert.Equal(t, updated.OnboardingToken, TOKEN1, "OnboardingToken not the same")
	assert.Equal(t, updated.CreatedTimestamp, fetched.CreatedTimestamp, "CreatedTimestamp should not be updated")
	assert.True(t, updated.LastModifiedTimestamp > fetched.LastModifiedTimestamp, "LastmodifiedTimestamp was not updated")

	// Add a second record.
	tenantUser2 := tenmod.User{
		Username:        USER2,
		Password:        PASS3,
		OnboardingToken: TOKEN2,
		TenantID:        TENANT,
		State:           string(common.UserInvited)}
	created2, err := runner.tenantDB.CreateTenantUser(&tenantUser2)
	assert.Nil(t, err)
	assert.NotNil(t, created2)
	assert.NotEmpty(t, created2.ID)
	assert.NotEmpty(t, created2.REV)
	assert.Equal(t, string(tenmod.TenantUserType), created2.Datatype)
	assert.Equal(t, created2.Username, USER2, "Username not the same")
	assert.Equal(t, created2.Password, PASS3, "Password not the same")
	assert.Equal(t, created2.TenantID, TENANT, "Tenant ID not the same")
	assert.Equal(t, created2.OnboardingToken, TOKEN2, "OnboardingToken not the same")
	assert.True(t, created2.CreatedTimestamp > 0, "CreatedTimestamp was not set")
	assert.True(t, created2.LastModifiedTimestamp > 0, "LastmodifiedTimestamp was not set")

	// Get all records
	fetchedList, err := runner.tenantDB.GetAllTenantUsers(TENANT)
	assert.Nil(t, err)
	assert.NotEmpty(t, fetchedList)
	assert.True(t, len(fetchedList) == 2)

	// Delete a record.
	deleted, err := runner.tenantDB.DeleteTenantUser(TENANT, fetched.ID)
	assert.Nil(t, err)
	assert.NotNil(t, deleted)
	assert.NotEmpty(t, deleted.ID)
	assert.NotEmpty(t, deleted.REV)
	assert.Equal(t, deleted.Username, fetched.Username, "Deleted Username not the same")

	// Get all records - should be 1
	fetchedList, err = runner.tenantDB.GetAllTenantUsers(TENANT)
	assert.Nil(t, err)
	assert.NotEmpty(t, fetchedList)
	assert.True(t, len(fetchedList) == 1)

	// Get a record that does not exist
	dne, err := runner.tenantDB.GetTenantUser(TENANT, deleted.ID)
	assert.NotNil(t, err)
	assert.Nil(t, dne)

	// Delete a record that oes not exist
	deleteDNE, err := runner.tenantDB.DeleteTenantUser(TENANT, deleted.ID)
	assert.NotNil(t, err)
	assert.Nil(t, deleteDNE)

	// Delete the last record
	deleted, err = runner.tenantDB.DeleteTenantUser(TENANT, created2.ID)
	assert.Nil(t, err)
	assert.NotNil(t, deleted)
	assert.NotEmpty(t, deleted.ID)
	assert.NotEmpty(t, deleted.REV)
	assert.Equal(t, deleted.Username, created2.Username, "Deleted Username not the same")

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
	tenantDomain := tenmod.Domain{
		Name:                DOM1,
		TenantID:            TENANT,
		Color:               COLOR1,
		ThresholdProfileSet: []string{THRPRF}}
	created, err := runner.tenantDB.CreateTenantDomain(&tenantDomain)
	assert.Nil(t, err)
	assert.NotNil(t, created)
	assert.NotEmpty(t, created.ID)
	assert.NotEmpty(t, created.REV)
	assert.Equal(t, string(tenmod.TenantDomainType), created.Datatype)
	assert.Equal(t, created.Name, DOM1, "Name not the same")
	assert.Equal(t, created.Color, COLOR1, "Color not the same")
	assert.Equal(t, created.ThresholdProfileSet[0], THRPRF, "Threshold Profile ID not the same")
	assert.Equal(t, created.TenantID, TENANT, "Tenant ID not the same")
	assert.True(t, created.CreatedTimestamp > 0, "CreatedTimestamp was not set")
	assert.True(t, created.LastModifiedTimestamp > 0, "LastmodifiedTimestamp was not set")

	// Get a record
	fetched, err := runner.tenantDB.GetTenantDomain(TENANT, created.ID)
	assert.Nil(t, err)
	assert.Equal(t, created, fetched, "The retrieved record should be the same as the created record")

	time.Sleep(time.Millisecond * 2)

	// Update a record
	updateRecord := tenmod.Domain{}
	deepcopy.Copy(&updateRecord, fetched)
	updateRecord.Color = COLOR2
	updated, err := runner.tenantDB.UpdateTenantDomain(&updateRecord)
	assert.Nil(t, err)
	assert.NotNil(t, updated)
	assert.Equal(t, updated.ID, fetched.ID)
	assert.NotEqual(t, updated.REV, fetched.REV)
	assert.Equal(t, string(tenmod.TenantDomainType), updated.Datatype)
	assert.Equal(t, updated.Name, DOM1, "Name not the same")
	assert.Equal(t, updated.Color, COLOR2, "Password was not updated")
	assert.Equal(t, updated.TenantID, TENANT, "Tenant ID not the same")
	assert.Equal(t, updated.ThresholdProfileSet[0], THRPRF, "Threshold Profile ID not the same")
	assert.Equal(t, updated.CreatedTimestamp, fetched.CreatedTimestamp, "CreatedTimestamp should not be updated")
	assert.True(t, updated.LastModifiedTimestamp > fetched.LastModifiedTimestamp, "LastmodifiedTimestamp was not updated")

	// Add a second record.
	tenantDomain2 := tenmod.Domain{
		Name:     DOM2,
		TenantID: TENANT,
		Color:    COLOR1}
	created2, err := runner.tenantDB.CreateTenantDomain(&tenantDomain2)
	assert.Nil(t, err)
	assert.NotNil(t, created2)
	assert.NotEmpty(t, created2.ID)
	assert.NotEmpty(t, created2.REV)
	assert.Equal(t, string(tenmod.TenantDomainType), created2.Datatype)
	assert.Equal(t, created2.Name, DOM2, "Name not the same")
	assert.Equal(t, created2.Color, COLOR1, "Password not the same")
	assert.Equal(t, created2.TenantID, TENANT, "Tenant ID not the same")
	assert.True(t, len(created2.ThresholdProfileSet) == 0, "Should not be a Threshold Profile ID")
	assert.True(t, created2.CreatedTimestamp > 0, "CreatedTimestamp was not set")
	assert.True(t, created2.LastModifiedTimestamp > 0, "LastmodifiedTimestamp was not set")

	// Get all records
	fetchedList, err := runner.tenantDB.GetAllTenantDomains(TENANT)
	assert.Nil(t, err)
	assert.NotEmpty(t, fetchedList)
	assert.True(t, len(fetchedList) == 2)

	// Delete a record.
	deleted, err := runner.tenantDB.DeleteTenantDomain(TENANT, fetched.ID)
	assert.Nil(t, err)
	assert.NotNil(t, deleted)
	assert.NotEmpty(t, deleted.ID)
	assert.NotEmpty(t, deleted.REV)
	assert.Equal(t, deleted.Name, fetched.Name, "Deleted name not the same")

	// Get all records - should be 1
	fetchedList, err = runner.tenantDB.GetAllTenantDomains(TENANT)
	assert.Nil(t, err)
	assert.NotEmpty(t, fetchedList)
	assert.True(t, len(fetchedList) == 1)

	// Get a record that does not exist
	dne, err := runner.tenantDB.GetTenantDomain(TENANT, deleted.ID)
	assert.NotNil(t, err)
	assert.Nil(t, dne)

	// Delete a record that oes not exist
	deleteDNE, err := runner.tenantDB.DeleteTenantDomain(TENANT, deleted.ID)
	assert.NotNil(t, err)
	assert.Nil(t, deleteDNE)

	// Delete the last record
	deleted, err = runner.tenantDB.DeleteTenantDomain(TENANT, created2.ID)
	assert.Nil(t, err)
	assert.NotNil(t, deleted)
	assert.NotEmpty(t, deleted.ID)
	assert.NotEmpty(t, deleted.REV)
	assert.Equal(t, deleted.Name, created2.Name, "Deleted name not the same")

	// Get all records - should be empty
	fetchedList, err = runner.tenantDB.GetAllTenantDomains(TENANT)
	assert.Nil(t, err)
	assert.NotNil(t, recList)
}
