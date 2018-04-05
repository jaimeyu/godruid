package admin

import (
	"encoding/json"
	"io/ioutil"
	"testing"
	"time"

	"github.com/accedian/adh-gather/logger"
	"github.com/getlantern/deepcopy"

	"github.com/stretchr/testify/assert"

	ds "github.com/accedian/adh-gather/datastore"
	admmod "github.com/accedian/adh-gather/models/admin"
	"github.com/accedian/adh-gather/models/common"
)

// AdminServiceDatastoreTestRunner - object used to run tests for any iplementation
// of the AdminServiceDatastore interface
type AdminServiceDatastoreTestRunner struct {
	adminDB ds.AdminServiceDatastore
}

func InitTestRunner(db ds.AdminServiceDatastore) *AdminServiceDatastoreTestRunner {
	return &AdminServiceDatastoreTestRunner{
		adminDB: db,
	}
}

func (runner *AdminServiceDatastoreTestRunner) RunAdminUserCRUD(t *testing.T) {

	const USER1 = "test1"
	const USER2 = "test2"
	const PASS1 = "pass1"
	const PASS2 = "pass2"
	const PASS3 = "pass3"
	const TOKEN1 = "token1"
	const TOKEN2 = "token2"

	// Validate that there are currently no records
	adminUserList, err := runner.adminDB.GetAllAdminUsers()
	assert.Nil(t, err)
	assert.Empty(t, adminUserList)

	// Create a record
	adminUser := admmod.User{
		Username:        USER1,
		Password:        PASS1,
		OnboardingToken: TOKEN1,
		State:           string(common.UserActive)}
	created, err := runner.adminDB.CreateAdminUser(&adminUser)
	assert.Nil(t, err)
	assert.NotNil(t, created)
	assert.NotEmpty(t, created.ID)
	assert.NotEmpty(t, created.REV)
	assert.Equal(t, string(admmod.AdminUserType), created.Datatype)
	assert.Equal(t, created.Username, USER1, "Username not the same")
	assert.Equal(t, created.Password, PASS1, "Password not the same")
	assert.Equal(t, created.OnboardingToken, TOKEN1, "OnboardingToken not the same")
	assert.True(t, created.CreatedTimestamp > 0, "CreatedTimestamp was not set")
	assert.True(t, created.LastModifiedTimestamp > 0, "LastmodifiedTimestamp was not set")

	// Get a record
	fetched, err := runner.adminDB.GetAdminUser(created.ID)
	assert.Nil(t, err)
	assert.Equal(t, created, fetched, "The retrieved record should be the same as the updated record")

	// Try to fetch a record even though none exist:
	fail, err := runner.adminDB.GetAdminUser("someID")
	assert.NotNil(t, err)
	assert.Nil(t, fail)

	time.Sleep(time.Millisecond * 2)

	// Update a record
	updateRecord := admmod.User{}
	deepcopy.Copy(&updateRecord, fetched)
	updateRecord.Password = PASS2
	updated, err := runner.adminDB.UpdateAdminUser(&updateRecord)
	assert.Nil(t, err)
	assert.NotNil(t, updated)
	assert.Equal(t, updated.ID, fetched.ID)
	assert.NotEqual(t, updated.REV, fetched.REV)
	assert.Equal(t, string(admmod.AdminUserType), updated.Datatype)
	assert.Equal(t, updated.Username, USER1, "Username not the same")
	assert.Equal(t, updated.Password, PASS2, "Password was not updated")
	assert.Equal(t, updated.OnboardingToken, TOKEN1, "OnboardingToken not the same")
	assert.Equal(t, updated.CreatedTimestamp, fetched.CreatedTimestamp, "CreatedTimestamp should not be updated")
	assert.True(t, updated.LastModifiedTimestamp > fetched.LastModifiedTimestamp, "LastmodifiedTimestamp was not updated")

	// Add a second record.
	adminUser2 := admmod.User{
		Username:        USER2,
		Password:        PASS3,
		OnboardingToken: TOKEN2,
		State:           string(common.UserActive)}
	created2, err := runner.adminDB.CreateAdminUser(&adminUser2)
	assert.Nil(t, err)
	assert.NotNil(t, created2)
	assert.NotEmpty(t, created2.ID)
	assert.NotEmpty(t, created2.REV)
	assert.Equal(t, string(admmod.AdminUserType), created2.Datatype)
	assert.Equal(t, created2.Username, USER2, "Username not the same")
	assert.Equal(t, created2.Password, PASS3, "Password not the same")
	assert.Equal(t, created2.OnboardingToken, TOKEN2, "OnboardingToken not the same")
	assert.True(t, created2.CreatedTimestamp > 0, "CreatedTimestamp was not set")
	assert.True(t, created2.LastModifiedTimestamp > 0, "LastmodifiedTimestamp was not set")

	// Get all records
	fetchedList, err := runner.adminDB.GetAllAdminUsers()
	assert.Nil(t, err)
	assert.NotEmpty(t, fetchedList)
	assert.True(t, len(fetchedList) == 2)

	// Delete a record that does not exist.
	deleted, err := runner.adminDB.DeleteAdminUser(string(fetched.ID))
	assert.Nil(t, err)
	assert.NotNil(t, deleted)
	assert.NotEmpty(t, deleted.ID)
	assert.NotEmpty(t, deleted.REV)
	assert.Equal(t, deleted.Username, fetched.Username, "Deleted Username not the same")

	// Get all records - should be 1
	fetchedList, err = runner.adminDB.GetAllAdminUsers()
	assert.Nil(t, err)
	assert.NotEmpty(t, fetchedList)
	assert.True(t, len(fetchedList) == 1)

	// Get a record that does not exist
	dne, err := runner.adminDB.GetAdminUser(deleted.ID)
	assert.NotNil(t, err)
	assert.Nil(t, dne)

	// Delete a record that oes not exist
	deleteDNE, err := runner.adminDB.DeleteAdminUser(deleted.ID)
	assert.NotNil(t, err)
	assert.Nil(t, deleteDNE)

	// Delete the last record
	deleted, err = runner.adminDB.DeleteAdminUser(string(created2.ID))
	assert.Nil(t, err)
	assert.NotNil(t, deleted)
	assert.NotEmpty(t, deleted.ID)
	assert.NotEmpty(t, deleted.REV)
	assert.Equal(t, deleted.Username, created2.Username, "Deleted Username not the same")

	// Get all records - should be empty
	fetchedList, err = runner.adminDB.GetAllAdminUsers()
	assert.Nil(t, err)
	assert.NotNil(t, fetchedList)
	assert.Empty(t, fetchedList)
}

func (runner *AdminServiceDatastoreTestRunner) RunTenantDescCRUD(t *testing.T) {

	const COMPANY1 = "test1"
	const COMPANY2 = "test2"
	const SUBDOMAIN1 = "pass1"
	const SUBDOMAIN2 = "pass2"
	const SUBDOMAIN3 = "pass3"

	// Validate that there are currently no records
	tenantList, err := runner.adminDB.GetAllTenantDescriptors()
	assert.Nil(t, err)
	assert.Empty(t, tenantList)

	// Create a record
	tenant1 := admmod.Tenant{
		Name:         COMPANY1,
		URLSubdomain: SUBDOMAIN1,
		State:        string(common.UserActive)}
	created, err := runner.adminDB.CreateTenant(&tenant1)
	assert.Nil(t, err)
	assert.NotNil(t, created)
	assert.NotEmpty(t, created.ID)
	assert.NotEmpty(t, created.REV)
	assert.Equal(t, string(admmod.TenantType), created.Datatype)
	assert.Equal(t, created.Name, COMPANY1, "Name not the same")
	assert.Equal(t, created.URLSubdomain, SUBDOMAIN1, "Subdomain not the same")
	assert.True(t, created.CreatedTimestamp > 0, "CreatedTimestamp was not set")
	assert.True(t, created.LastModifiedTimestamp > 0, "LastmodifiedTimestamp was not set")

	// Get a record
	fetched, err := runner.adminDB.GetTenantDescriptor(created.ID)
	assert.Nil(t, err)
	assert.Equal(t, created, fetched, "The retrieved record should be the same as the updated record")

	// Try to fetch a record even though none exist:
	fail, err := runner.adminDB.GetTenantDescriptor("someID")
	assert.NotNil(t, err)
	assert.Nil(t, fail)

	time.Sleep(time.Millisecond * 2)

	// Update a record
	updateRecord := admmod.Tenant{}
	deepcopy.Copy(&updateRecord, fetched)
	updateRecord.URLSubdomain = SUBDOMAIN2
	updated, err := runner.adminDB.UpdateTenantDescriptor(&updateRecord)
	assert.Nil(t, err)
	assert.NotNil(t, updated)
	assert.Equal(t, updated.ID, fetched.ID)
	assert.NotEqual(t, updated.REV, fetched.REV)
	assert.Equal(t, string(admmod.TenantType), updated.Datatype)
	assert.Equal(t, updated.Name, COMPANY1, "Name not the same")
	assert.Equal(t, updated.URLSubdomain, SUBDOMAIN2, "Subdomain was not updated")
	assert.Equal(t, updated.CreatedTimestamp, fetched.CreatedTimestamp, "CreatedTimestamp should not be updated")
	assert.True(t, updated.LastModifiedTimestamp > fetched.LastModifiedTimestamp, "LastmodifiedTimestamp was not updated")

	// Add a second record.
	tenant2 := admmod.Tenant{
		Name:         COMPANY2,
		URLSubdomain: SUBDOMAIN3,
		State:        string(common.UserActive)}
	created2, err := runner.adminDB.CreateTenant(&tenant2)
	assert.Nil(t, err)
	assert.NotNil(t, created2)
	assert.NotEmpty(t, created2.ID)
	assert.NotEmpty(t, created2.REV)
	assert.Equal(t, string(admmod.TenantType), created2.Datatype)
	assert.Equal(t, created2.Name, COMPANY2, "Name not the same")
	assert.Equal(t, created2.URLSubdomain, SUBDOMAIN3, "Subdomain not the same")
	assert.True(t, created2.CreatedTimestamp > 0, "CreatedTimestamp was not set")
	assert.True(t, created2.LastModifiedTimestamp > 0, "LastmodifiedTimestamp was not set")

	// Get all records
	fetchedList, err := runner.adminDB.GetAllTenantDescriptors()
	assert.Nil(t, err)
	assert.NotEmpty(t, fetchedList)
	assert.True(t, len(fetchedList) == 2)

	// Fetch a Tenant ID by username
	tenantID, err := runner.adminDB.GetTenantIDByAlias(COMPANY1)
	assert.NotEmpty(t, tenantID)
	assert.Equal(t, updated.ID, tenantID)

	// Delete a record that does not exist.
	deleted, err := runner.adminDB.DeleteTenant(string(fetched.ID))
	assert.Nil(t, err)
	assert.NotNil(t, deleted)
	assert.NotEmpty(t, deleted.ID)
	assert.NotEmpty(t, deleted.REV)
	assert.Equal(t, deleted.Name, fetched.Name, "Deleted Name not the same")

	// Get all records - should be 1
	fetchedList, err = runner.adminDB.GetAllTenantDescriptors()
	assert.Nil(t, err)
	assert.NotEmpty(t, fetchedList)
	assert.True(t, len(fetchedList) == 1)

	// Get a record that does not exist
	dne, err := runner.adminDB.GetTenantDescriptor(deleted.ID)
	assert.NotNil(t, err)
	assert.Nil(t, dne)

	// Delete a record that oes not exist
	deleteDNE, err := runner.adminDB.DeleteTenant(deleted.ID)
	assert.NotNil(t, err)
	assert.Nil(t, deleteDNE)

	// Delete the last record
	deleted, err = runner.adminDB.DeleteTenant(string(created2.ID))
	assert.Nil(t, err)
	assert.NotNil(t, deleted)
	assert.NotEmpty(t, deleted.ID)
	assert.NotEmpty(t, deleted.REV)
	assert.Equal(t, deleted.Name, created2.Name, "Deleted Name not the same")

	// Get all records - should be empty
	fetchedList, err = runner.adminDB.GetAllTenantDescriptors()
	assert.Nil(t, err)
	assert.Empty(t, fetchedList)
}

func (runner *AdminServiceDatastoreTestRunner) RunIngDictCRUD(t *testing.T) {

	var accTWAMP = "accedian-twamp"
	var accFLOW = "accedian-flowmeter"

	var delayMin = "delayMin"
	var delayMax = "delayMax"
	var delayAvg = "delayAvg"

	var throughputAvg = "throughputAvg"
	var throughputMax = "throughputMax"
	var throughputMin = "throughputMin"

	// Validate that there are currently no records
	ingPrf, err := runner.adminDB.GetIngestionDictionary()
	assert.NotNil(t, err)
	assert.Nil(t, ingPrf)

	// Read in the test dictionary from file
	defaultDictionaryBytes, err := ioutil.ReadFile("../test/files/testIngestionDictionary.json")
	if err != nil {
		logger.Log.Fatalf("Unable to read Default Ingestion Profile from file: %s", err.Error())
	}

	defaultDictionaryData := &admmod.IngestionDictionary{}
	if err = json.Unmarshal(defaultDictionaryBytes, &defaultDictionaryData); err != nil {
		logger.Log.Fatalf("Unable to construct Default Ingestion Profile from file: %s", err.Error())
	}

	// Create a record
	created, err := runner.adminDB.CreateIngestionDictionary(defaultDictionaryData)
	assert.Nil(t, err)
	assert.NotNil(t, created)
	assert.NotEmpty(t, created.ID)
	assert.NotEmpty(t, created.REV)
	assert.Equal(t, string(admmod.IngestionDictionaryType), created.Datatype)
	assert.NotEmpty(t, created.Metrics, "There should be metrics")
	assert.NotNil(t, created.Metrics[accTWAMP], "There should be accedian-twamp metrics")
	assert.NotNil(t, created.Metrics[accFLOW], "There should be accedian-flowmeter metrics")
	assert.NotEmpty(t, created.Metrics[accTWAMP].MetricMap, "There should be accedian-twamp metric definitions")
	assert.NotNil(t, created.Metrics[accTWAMP].MetricMap[delayMin], "There should be delayMin metrics")
	assert.NotNil(t, created.Metrics[accTWAMP].MetricMap[delayMax], "There should be delayMax metrics")
	assert.NotNil(t, created.Metrics[accTWAMP].MetricMap[delayAvg], "There should be delayAvg metrics")
	assert.NotEmpty(t, created.Metrics[accTWAMP].MetricMap[delayMin].MonitoredObjectTypes, "There should be delayMin monitored object definitions")
	assert.True(t, len(created.Metrics[accTWAMP].MetricMap[delayMin].MonitoredObjectTypes) == 3, "There should be 3 delayMin monitored object definitions")
	assert.NotEmpty(t, created.Metrics[accFLOW].MetricMap, "There should be accedian-flowmeter metric definitions")
	assert.NotNil(t, created.Metrics[accFLOW].MetricMap[throughputAvg], "There should be throughputAvg metrics")
	assert.NotNil(t, created.Metrics[accFLOW].MetricMap[throughputMax], "There should be throughputMax metrics")
	assert.NotNil(t, created.Metrics[accFLOW].MetricMap[throughputMin], "There should be throughputMin metrics")
	assert.NotEmpty(t, created.Metrics[accFLOW].MetricMap[throughputMax].MonitoredObjectTypes, "There should be throughputMax monitored object definitions")
	assert.True(t, len(created.Metrics[accFLOW].MetricMap[throughputMax].MonitoredObjectTypes) == 1, "There should be 1 throughputMax monitored object definitions")

	// Get a record
	fetched, err := runner.adminDB.GetIngestionDictionary()
	assert.Nil(t, err)
	assert.Equal(t, created, fetched, "The retrieved record should be the same as the created record")

	// Try to create a record that already exists, should fail
	created, err = runner.adminDB.CreateIngestionDictionary(created)
	assert.NotNil(t, err)
	assert.Nil(t, created, "Created should now be nil")

	time.Sleep(time.Millisecond * 2)

	// Update a record
	updateRecord := admmod.IngestionDictionary{}
	deepcopy.Copy(&updateRecord, fetched)
	delete(updateRecord.Metrics, accFLOW)
	updated, err := runner.adminDB.UpdateIngestionDictionary(&updateRecord)
	assert.Nil(t, err)
	assert.NotNil(t, updated)
	assert.NotEmpty(t, updated.ID)
	assert.NotEmpty(t, updated.REV)
	assert.Equal(t, fetched.ID, updated.ID, "Id values should be the same")
	assert.NotEqual(t, fetched.REV, updated.REV)
	assert.Equal(t, string(admmod.IngestionDictionaryType), updated.Datatype)
	assert.NotEmpty(t, updated.Metrics, "There should be metrics")
	assert.NotNil(t, updated.Metrics[accTWAMP], "There should be accedian-twamp metrics")
	assert.Nil(t, updated.Metrics[accFLOW], "There should not be any accedian-flowmeter metrics")

	// Delete the record
	deleted, err := runner.adminDB.DeleteIngestionDictionary()
	assert.Nil(t, err)
	assert.NotNil(t, deleted)
	assert.NotEmpty(t, deleted.ID)
	assert.NotEmpty(t, deleted.REV)
	assert.Equal(t, deleted, updated, "Deleted record not the same as last known version")

	// Get record - should fail
	fetched, err = runner.adminDB.GetIngestionDictionary()
	assert.NotNil(t, err)
	assert.Nil(t, fetched)

	// Delete record - should fail as no record exists
	fetched, err = runner.adminDB.DeleteIngestionDictionary()
	assert.NotNil(t, err)
	assert.Nil(t, fetched)
}

func (runner *AdminServiceDatastoreTestRunner) RunValidTypesCRUD(t *testing.T) {

	objTypeKey1 := "objTypeKey1"
	objTypeKey2 := "objTypeKey2"
	devTypeKey1 := "devTypeKey1"
	devTypeKey2 := "devTypeKey2"
	devTypeKey3 := "devTypeKey3"

	objTypeVal1 := "objTypeVal1"
	objTypeVal2 := "objTypeVal2"
	devTypeVal1 := "devTypeVal1"
	devTypeVal2 := "devTypeVal2"
	devTypeVal3 := "devTypeVal3"

	// Validate that there are currently no records
	validTypes, err := runner.adminDB.GetValidTypes()
	assert.NotNil(t, err)
	assert.Nil(t, validTypes)

	validTypeData := admmod.ValidTypes{
		MonitoredObjectTypes:       map[string]string{objTypeKey1: objTypeVal1, objTypeKey2: objTypeVal2},
		MonitoredObjectDeviceTypes: map[string]string{devTypeKey1: devTypeVal1, devTypeKey2: devTypeVal2}}
	// Create a record
	created, err := runner.adminDB.CreateValidTypes(&validTypeData)
	assert.Nil(t, err)
	assert.NotNil(t, created)
	assert.NotEmpty(t, created.ID)
	assert.NotEmpty(t, created.REV)
	assert.NotEmpty(t, created.MonitoredObjectTypes, "There should be mon obj types")
	assert.Equal(t, 2, len(created.MonitoredObjectTypes), "There should be 2 mon obj types")
	assert.NotEmpty(t, created.MonitoredObjectDeviceTypes, "There should be mon obj dev types")
	assert.Equal(t, 2, len(created.MonitoredObjectDeviceTypes), "There should be 2 mon obj dev types")
	assert.Equal(t, objTypeVal1, created.MonitoredObjectTypes[objTypeKey1])
	assert.Equal(t, devTypeVal1, created.MonitoredObjectDeviceTypes[devTypeKey1])

	// Get a record
	fetched, err := runner.adminDB.GetValidTypes()
	assert.Nil(t, err)
	assert.Equal(t, created, fetched, "The retrieved record should be the same as the created record")

	// Try to create a record that already exists, should fail
	created, err = runner.adminDB.CreateValidTypes(created)
	assert.NotNil(t, err)
	assert.Nil(t, created, "Created should now be nil")

	time.Sleep(time.Millisecond * 2)

	// Update a record
	updateRecord := admmod.ValidTypes{}
	deepcopy.Copy(&updateRecord, fetched)
	updateRecord.MonitoredObjectDeviceTypes[devTypeKey3] = devTypeVal3
	updated, err := runner.adminDB.UpdateValidTypes(&updateRecord)
	assert.Nil(t, err)
	assert.NotNil(t, updated)
	assert.NotEmpty(t, updated.ID)
	assert.NotEmpty(t, updated.REV)
	assert.Equal(t, fetched.ID, updated.ID, "Id values should be the same")
	assert.NotEqual(t, fetched.REV, updated.REV)
	assert.NotEmpty(t, updated.MonitoredObjectTypes, "There should be mon obj types")
	assert.Equal(t, 2, len(updated.MonitoredObjectTypes), "There should be 2 mon obj types")
	assert.NotEmpty(t, updated.MonitoredObjectDeviceTypes, "There should be mon obj dev types")
	assert.Equal(t, 3, len(updated.MonitoredObjectDeviceTypes), "There should be 2 mon obj dev types")
	assert.Equal(t, devTypeVal3, updated.MonitoredObjectDeviceTypes[devTypeKey3])

	// Retrieve a specific valid type
	specific, err := runner.adminDB.GetSpecificValidTypes(&admmod.ValidTypesRequest{MonitoredObjectTypes: true})
	assert.Nil(t, err)
	assert.NotNil(t, specific)
	assert.NotNil(t, specific.MonitoredObjectTypes)
	assert.Nil(t, specific.MonitoredObjectDeviceTypes)
	assert.Equal(t, updated.MonitoredObjectTypes, specific.MonitoredObjectTypes)

	specific, err = runner.adminDB.GetSpecificValidTypes(&admmod.ValidTypesRequest{MonitoredObjectDeviceTypes: true})
	assert.Nil(t, err)
	assert.NotNil(t, specific)
	assert.NotNil(t, specific.MonitoredObjectDeviceTypes)
	assert.Nil(t, specific.MonitoredObjectTypes)
	assert.Equal(t, updated.MonitoredObjectDeviceTypes, specific.MonitoredObjectDeviceTypes)

	specific, err = runner.adminDB.GetSpecificValidTypes(&admmod.ValidTypesRequest{MonitoredObjectDeviceTypes: false, MonitoredObjectTypes: false})
	assert.Nil(t, err)
	assert.NotNil(t, specific)
	assert.Nil(t, specific.MonitoredObjectDeviceTypes)
	assert.Nil(t, specific.MonitoredObjectTypes)

	// Delete the record
	deleted, err := runner.adminDB.DeleteValidTypes()
	assert.Nil(t, err)
	assert.NotNil(t, deleted)
	assert.NotEmpty(t, deleted.ID)
	assert.NotEmpty(t, deleted.REV)
	assert.Equal(t, deleted, updated, "Deleted record not the same as last known version")

	// Get record - should fail
	fetched, err = runner.adminDB.GetValidTypes()
	assert.NotNil(t, err)
	assert.Nil(t, fetched)

	// Delete record - should fail as no record exists
	fetched, err = runner.adminDB.DeleteValidTypes()
	assert.NotNil(t, err)
	assert.Nil(t, fetched)
}
