package test

import (
	"log"
	"testing"
	"fmt"
	"strconv"
	"io/ioutil"
	"encoding/json"

	"github.com/accedian/adh-gather/gather"
	"github.com/accedian/adh-gather/logger"
	"github.com/accedian/adh-gather/handlers"
	"github.com/spf13/viper"

	"github.com/accedian/adh-gather/config"
	"github.com/leesper/couchdb-golang"
	"github.com/stretchr/testify/assert"

	dockertest "gopkg.in/ory-am/dockertest.v3"
	pb "github.com/accedian/adh-gather/gathergrpc"
	ds "github.com/accedian/adh-gather/datastore"
	mem "github.com/accedian/adh-gather/datastore/inMemory"
	couchDB "github.com/accedian/adh-gather/datastore/couchDB"
)

const (
	adminDBName = "adh-admin"
)

var (
	couchServer *couchdb.Server
	cfg         config.Provider
	adminDB 	ds.AdminServiceDatastore
)

func TestMain(m *testing.M) {
	pool, err := dockertest.NewPool("")
	if err != nil {
		log.Fatalf("Could not connect to docker: %s", err)
	}

	// pulls the couchdb image and start it in a mode that will self-cleanup.
	resource, err := pool.Run("apache/couchdb", "2.1.1", []string{"--rm", "-it"})
	if err != nil {
		log.Fatalf("Could not start resource: %s", err)
	}

	// exponential backoff-retry, because the application in the container might not be ready to accept connections yet
	if err := pool.Retry(func() error {
		var err error

		// Spin up the couch docker image and make sure we can reach it.
		couchHost := "http://0.0.0.0"
		couchPort := resource.GetPort("5984/tcp")
		couchServer, err = couchdb.NewServer(fmt.Sprintf("%s:%s", couchHost, couchPort))
		if err != nil {
			return err
		}

		// Configure the test AdminService DAO to use the newly started couch docker image
		cfg := gather.LoadConfig("../../config/adh-gather-debug.yml", viper.New())
		cfg.Set(gather.CK_server_datastore_ip.String(), couchHost)
		portAsInt, err := strconv.Atoi(couchPort)
		if err != nil {
			return err
		}
		cfg.Set(gather.CK_server_datastore_port.String(), portAsInt)

		ver, err := couchServer.Version()
		logger.Log.Debugf("Test CouchDB version is: %s", ver)

		return err
	}); err != nil {
		log.Fatalf("Could not connect to docker: %s", err.Error())
	}

	// Before the tests run, setup the adh-admin db

	// Couch Run.
	adminDB, err = couchDB.CreateAdminServiceDAO()
	if err != nil {
		log.Fatalf("Could not create couchdb admin DAO: %s", err.Error())
	}
	_, err = adminDB.CreateDatabase(adminDBName)
	if err != nil {
		log.Fatalf("Could not create admin DB: %s", err.Error())
	}
	err = adminDB.AddAdminViews()
	if err != nil {
		log.Fatalf("Could not populate admin indicies: %s", err.Error())
	}

	defer func(res *dockertest.Resource) {
		if err := pool.Purge(resource); err != nil {
			log.Fatalf("Could not purge resource: %s", err.Error())
		}
	}(resource)

	m.Run()


	// InMemory Run:
	adminDB, err = mem.CreateAdminServiceDAO()
	if err != nil {
		log.Fatalf("Could not create in-mem admin DAO: %s", err.Error())
	}

	m.Run()
}

func TestAdminUserCRUD(t *testing.T) {
	const USER1 = "test1"
	const USER2 = "test2"
	const PASS1 = "pass1"
	const PASS2 = "pass2"
	const PASS3 = "pass3"
	const TOKEN1 = "token1"
	const TOKEN2 = "token2"

	// Validate that there are currently no records
	adminUserList, err := adminDB.GetAllAdminUsers()
	assert.Nil(t, err)
	assert.NotNil(t, adminUserList)
	assert.Empty(t, adminUserList.Data)

	// Create a record
	adminUser := pb.AdminUserData{
		Username: USER1,
		Password: PASS1,
		OnboardingToken: TOKEN1,
		State: pb.UserState_INVITED}
	created, err := adminDB.CreateAdminUser(&pb.AdminUser{ Data: &adminUser})
	assert.Nil(t, err)
	assert.NotNil(t, created)
	assert.NotNil(t, created.XId)
	assert.NotNil(t, created.XRev)
	assert.NotEmpty(t, created.XId)
	assert.NotEmpty(t, created.XRev)
	assert.Equal(t, string(ds.AdminUserType), created.Data.Datatype)
	assert.Equal(t, created.Data.Username, USER1, "Username not the same")
	assert.Equal(t, created.Data.Password, PASS1,"Password not the same")
	assert.Equal(t, created.Data.OnboardingToken, TOKEN1, "OnboardingToken not the same")
	assert.True(t, created.Data.CreatedTimestamp > 0, "CreatedTimestamp was not set")
	assert.True(t, created.Data.LastModifiedTimestamp > 0, "LastmodifiedTimestamp was not set")

	// Get a record
	fetched, err := adminDB.GetAdminUser(created.XId)
	assert.Nil(t, err)
	assert.Equal(t, created, fetched, "The retrieved record should be the same as the updated record")

	// Try to create a record that already exists, should fail
	created, err = adminDB.CreateAdminUser(created)
	assert.NotNil(t, err)
	assert.Nil(t, created, "Created should now be nil")

	// Update a record
	updateRecord := *fetched
	updateRecord.Data.Password = PASS2
	updated, err := adminDB.UpdateAdminUser(&updateRecord)
	assert.Nil(t, err)
	assert.NotNil(t, updated)
	assert.Equal(t, updated.XId, fetched.XId)
	assert.NotEqual(t, updated.XRev, fetched.XRev)
	assert.Equal(t, string(ds.AdminUserType), updated.Data.Datatype)
	assert.Equal(t, updated.Data.Username, USER1, "Username not the same")
	assert.Equal(t, updated.Data.Password, PASS2,"Password was not updated")
	assert.Equal(t, updated.Data.OnboardingToken, TOKEN1, "OnboardingToken not the same")
	assert.Equal(t, updated.Data.CreatedTimestamp, fetched.Data.CreatedTimestamp, "CreatedTimestamp should not be updated")
	assert.True(t, updated.Data.LastModifiedTimestamp > fetched.Data.LastModifiedTimestamp, "LastmodifiedTimestamp was not updated")

	// Add a second record.
	adminUser2 := pb.AdminUserData{
		Username: USER2,
		Password: PASS3,
		OnboardingToken: TOKEN2,
		State: pb.UserState_INVITED}
	created2, err := adminDB.CreateAdminUser(&pb.AdminUser{ Data: &adminUser2})
	assert.Nil(t, err)
	assert.NotNil(t, created2)
	assert.NotNil(t, created2.XId)
	assert.NotNil(t, created2.XRev)
	assert.NotEmpty(t, created2.XId)
	assert.NotEmpty(t, created2.XRev)
	assert.Equal(t, string(ds.AdminUserType), created2.Data.Datatype)
	assert.Equal(t, created2.Data.Username, USER2, "Username not the same")
	assert.Equal(t, created2.Data.Password, PASS3,"Password not the same")
	assert.Equal(t, created2.Data.OnboardingToken, TOKEN2, "OnboardingToken not the same")
	assert.True(t, created2.Data.CreatedTimestamp > 0, "CreatedTimestamp was not set")
	assert.True(t, created2.Data.LastModifiedTimestamp > 0, "LastmodifiedTimestamp was not set")

	// Get all records 
	fetchedList, err := adminDB.GetAllAdminUsers()
	assert.Nil(t, err)
	assert.NotNil(t, fetchedList)
	assert.NotEmpty(t, fetchedList.Data)
	assert.True(t, len(fetchedList.Data) == 2)

	// Delete a record that does not exist.
	deleted, err := adminDB.DeleteAdminUser(string(fetched.XId))
	assert.Nil(t, err)
	assert.NotNil(t, deleted)
	assert.NotNil(t, deleted.XId)
	assert.NotNil(t, deleted.XRev)
	assert.NotEmpty(t, deleted.XId)
	assert.NotEmpty(t, deleted.XRev)
	assert.Equal(t, deleted.Data.Username, fetched.Data.Username, "Deleted Username not the same")

	// Get all records - should be 1
	fetchedList, err = adminDB.GetAllAdminUsers()
	assert.Nil(t, err)
	assert.NotNil(t, fetchedList)
	assert.NotEmpty(t, fetchedList.Data)
	assert.True(t, len(fetchedList.Data) == 1)

	// Get a record that does not exist
	dne, err := adminDB.GetAdminUser(deleted.XId)
	assert.NotNil(t, err)
	assert.Nil(t, dne)

	// Delete a record that oes not exist
	deleteDNE, err := adminDB.DeleteAdminUser(deleted.XId)
	assert.NotNil(t, err)
	assert.Nil(t, deleteDNE)

	// Delete the last record
	deleted, err = adminDB.DeleteAdminUser(string(created2.XId))
	assert.Nil(t, err)
	assert.NotNil(t, deleted)
	assert.NotNil(t, deleted.XId)
	assert.NotNil(t, deleted.XRev)
	assert.NotEmpty(t, deleted.XId)
	assert.NotEmpty(t, deleted.XRev)
	assert.Equal(t, deleted.Data.Username, created2.Data.Username, "Deleted Username not the same")

	// Get all records - should be empty
	fetchedList, err = adminDB.GetAllAdminUsers()
	assert.Nil(t, err)
	assert.NotNil(t, fetchedList)
	assert.Empty(t, fetchedList.Data)
}

func TestTenantDescCRUD(t *testing.T) {
	const COMPANY1 = "test1"
	const COMPANY2 = "test2"
	const SUBDOMAIN1 = "pass1"
	const SUBDOMAIN2 = "pass2"
	const SUBDOMAIN3 = "pass3"

	// Validate that there are currently no records
	tenantList, err := adminDB.GetAllTenantDescriptors()
	assert.Nil(t, err)
	assert.NotNil(t, tenantList)
	assert.Empty(t, tenantList.Data)

	// Create a record
	tenant1 := pb.TenantDescriptorData{
		Name: COMPANY1,
		UrlSubdomain: SUBDOMAIN1,
		State: pb.UserState_ACTIVE}
	created, err := adminDB.CreateTenant(&pb.TenantDescriptor{ Data: &tenant1})
	assert.Nil(t, err)
	assert.NotNil(t, created)
	assert.NotNil(t, created.XId)
	assert.NotNil(t, created.XRev)
	assert.NotEmpty(t, created.XId)
	assert.NotEmpty(t, created.XRev)
	assert.Equal(t, string(ds.TenantDescriptorType), created.Data.Datatype)
	assert.Equal(t, created.Data.Name, COMPANY1, "Name not the same")
	assert.Equal(t, created.Data.UrlSubdomain, SUBDOMAIN1,"Subdomain not the same")
	assert.True(t, created.Data.CreatedTimestamp > 0, "CreatedTimestamp was not set")
	assert.True(t, created.Data.LastModifiedTimestamp > 0, "LastmodifiedTimestamp was not set")

	// Get a record
	fetched, err := adminDB.GetTenantDescriptor(created.XId)
	assert.Nil(t, err)
	assert.Equal(t, created, fetched, "The retrieved record should be the same as the updated record")

	// Try to create a record that already exists, should fail
	created, err = adminDB.CreateTenant(created)
	assert.NotNil(t, err)
	assert.Nil(t, created, "Created should now be nil")

	// Update a record
	updateRecord := *fetched
	updateRecord.Data.UrlSubdomain = SUBDOMAIN2
	updated, err := adminDB.UpdateTenantDescriptor(&updateRecord)
	assert.Nil(t, err)
	assert.NotNil(t, updated)
	assert.Equal(t, updated.XId, fetched.XId)
	assert.NotEqual(t, updated.XRev, fetched.XRev)
	assert.Equal(t, string(ds.TenantDescriptorType), updated.Data.Datatype)
	assert.Equal(t, updated.Data.Name, COMPANY1, "Name not the same")
	assert.Equal(t, updated.Data.UrlSubdomain, SUBDOMAIN2,"Subdomain was not updated")
	assert.Equal(t, updated.Data.CreatedTimestamp, fetched.Data.CreatedTimestamp, "CreatedTimestamp should not be updated")
	assert.True(t, updated.Data.LastModifiedTimestamp > fetched.Data.LastModifiedTimestamp, "LastmodifiedTimestamp was not updated")

	// Add a second record.
	tenant2 := pb.TenantDescriptorData{
		Name: COMPANY2,
		UrlSubdomain: SUBDOMAIN3,
		State: pb.UserState_ACTIVE}
	created2, err := adminDB.CreateTenant(&pb.TenantDescriptor{ Data: &tenant2})
	assert.Nil(t, err)
	assert.NotNil(t, created2)
	assert.NotNil(t, created2.XId)
	assert.NotNil(t, created2.XRev)
	assert.NotEmpty(t, created2.XId)
	assert.NotEmpty(t, created2.XRev)
	assert.Equal(t, string(ds.TenantDescriptorType), created2.Data.Datatype)
	assert.Equal(t, created2.Data.Name, COMPANY2, "Name not the same")
	assert.Equal(t, created2.Data.UrlSubdomain, SUBDOMAIN3,"Subdomain not the same")
	assert.True(t, created2.Data.CreatedTimestamp > 0, "CreatedTimestamp was not set")
	assert.True(t, created2.Data.LastModifiedTimestamp > 0, "LastmodifiedTimestamp was not set")

	// Get all records 
	fetchedList, err := adminDB.GetAllTenantDescriptors()
	assert.Nil(t, err)
	assert.NotNil(t, fetchedList)
	assert.NotEmpty(t, fetchedList.Data)
	assert.True(t, len(fetchedList.Data) == 2)

	// Fetch a Tenant ID by username
	tenantID, err := adminDB.GetTenantIDByAlias(COMPANY1)
	assert.Nil(t, err)
	assert.NotEmpty(t, tenantID)
	assert.Equal(t, updated.XId, tenantID)

	// Delete a record that does not exist.
	deleted, err := adminDB.DeleteTenant(string(fetched.XId))
	assert.Nil(t, err)
	assert.NotNil(t, deleted)
	assert.NotNil(t, deleted.XId)
	assert.NotNil(t, deleted.XRev)
	assert.NotEmpty(t, deleted.XId)
	assert.NotEmpty(t, deleted.XRev)
	assert.Equal(t, deleted.Data.Name, fetched.Data.Name, "Deleted Name not the same")

	// Get all records - should be 1
	fetchedList, err = adminDB.GetAllTenantDescriptors()
	assert.Nil(t, err)
	assert.NotNil(t, fetchedList)
	assert.NotEmpty(t, fetchedList.Data)
	assert.True(t, len(fetchedList.Data) == 1)

	// Get a record that does not exist
	dne, err := adminDB.GetTenantDescriptor(deleted.XId)
	assert.NotNil(t, err)
	assert.Nil(t, dne)

	// Delete a record that oes not exist
	deleteDNE, err := adminDB.DeleteTenant(deleted.XId)
	assert.NotNil(t, err)
	assert.Nil(t, deleteDNE)

	// Delete the last record
	deleted, err = adminDB.DeleteTenant(string(created2.XId))
	assert.Nil(t, err)
	assert.NotNil(t, deleted)
	assert.NotNil(t, deleted.XId)
	assert.NotNil(t, deleted.XRev)
	assert.NotEmpty(t, deleted.XId)
	assert.NotEmpty(t, deleted.XRev)
	assert.Equal(t, deleted.Data.Name, created2.Data.Name, "Deleted Name not the same")

	// Get all records - should be empty
	fetchedList, err = adminDB.GetAllTenantDescriptors()
	assert.Nil(t, err)
	assert.NotNil(t, fetchedList)
	assert.Empty(t, fetchedList.Data)
}

func TestIngDictCRUD(t *testing.T) {
	var accTWAMP = string(handlers.AccedianTwamp)
	var accFLOW = string(handlers.AccedianFlowmeter)

	var delayMin = "delayMin"
	var delayMax = "delayMax"
	var delayAvg = "delayAvg"

	var throughputAvg = "throughputAvg"
	var throughputMax = "throughputMax"
	var throughputMin = "throughputMin"

	// Validate that there are currently no records
	ingPrf, err := adminDB.GetIngestionDictionary()
	assert.NotNil(t, err)
	assert.Nil(t, ingPrf)

	// Read in the test dictionary from file
	defaultDictionaryBytes, err := ioutil.ReadFile("./files/testIngestionDictionary.json")
	if err != nil {
		logger.Log.Fatalf("Unable to read Default Ingestion Profile from file: %s", err.Error())
	}
	
	defaultDictionaryData := &pb.IngestionDictionaryData{}
	if err = json.Unmarshal(defaultDictionaryBytes, &defaultDictionaryData); err != nil {
		logger.Log.Fatalf("Unable to construct Default Ingestion Profile from file: %s", err.Error())
	}

	// Create a record
	created, err := adminDB.CreateIngestionDictionary(&pb.IngestionDictionary{ Data: defaultDictionaryData})
	assert.Nil(t, err)
	assert.NotNil(t, created)
	assert.NotNil(t, created.XId)
	assert.NotNil(t, created.XRev)
	assert.NotEmpty(t, created.XId)
	assert.NotEmpty(t, created.XRev)
	assert.Equal(t, string(ds.IngestionDictionaryType), created.Data.Datatype)
	assert.NotEmpty(t, created.Data.Metrics, "There should be metrics")
	assert.NotNil(t, created.Data.Metrics[accTWAMP], "There should be accedian-twamp metrics")
	assert.NotNil(t, created.Data.Metrics[accFLOW], "There should be accedian-flowmeter metrics")
	assert.NotEmpty(t, created.Data.Metrics[accTWAMP].MetricMap, "There should be accedian-twamp metric definitions")
	assert.NotNil(t, created.Data.Metrics[accTWAMP].MetricMap[delayMin], "There should be delayMin metrics")
	assert.NotNil(t, created.Data.Metrics[accTWAMP].MetricMap[delayMax], "There should be delayMax metrics")
	assert.NotNil(t, created.Data.Metrics[accTWAMP].MetricMap[delayAvg], "There should be delayAvg metrics")
	assert.NotEmpty(t, created.Data.Metrics[accTWAMP].MetricMap[delayMin].MonitoredObjectTypes, "There should be delayMin monitored object definitions")
	assert.True(t, len(created.Data.Metrics[accTWAMP].MetricMap[delayMin].MonitoredObjectTypes) == 3 , "There should be 3 delayMin monitored object definitions")
	assert.NotEmpty(t, created.Data.Metrics[accFLOW].MetricMap, "There should be accedian-flowmeter metric definitions")
	assert.NotNil(t, created.Data.Metrics[accFLOW].MetricMap[throughputAvg], "There should be throughputAvg metrics")
	assert.NotNil(t, created.Data.Metrics[accFLOW].MetricMap[throughputMax], "There should be throughputMax metrics")
	assert.NotNil(t, created.Data.Metrics[accFLOW].MetricMap[throughputMin], "There should be throughputMin metrics")
	assert.NotEmpty(t, created.Data.Metrics[accFLOW].MetricMap[throughputMax].MonitoredObjectTypes, "There should be throughputMax monitored object definitions")
	assert.True(t, len(created.Data.Metrics[accFLOW].MetricMap[throughputMax].MonitoredObjectTypes) == 1 , "There should be 1 throughputMax monitored object definitions")

	// Get a record
	fetched, err := adminDB.GetIngestionDictionary()
	assert.Nil(t, err)
	assert.Equal(t, created, fetched, "The retrieved record should be the same as the created record")

	// Try to create a record that already exists, should fail
	created, err = adminDB.CreateIngestionDictionary(created)
	assert.NotNil(t, err)
	assert.Nil(t, created, "Created should now be nil")

	// Update a record
	updateRecord := *fetched
	delete(updateRecord.Data.Metrics, accFLOW)
	updated, err := adminDB.UpdateIngestionDictionary(&updateRecord)
	assert.Nil(t, err)
	assert.NotNil(t, updated)
	assert.NotNil(t, updated.XId)
	assert.NotNil(t, updated.XRev)
	assert.NotEmpty(t, updated.XId)
	assert.NotEmpty(t, updated.XRev)
	assert.Equal(t, fetched.XId, updated.XId, "Id values should be the same")
	assert.NotEqual(t, fetched.XRev, updated.XRev)
	assert.Equal(t, string(ds.IngestionDictionaryType), updated.Data.Datatype)
	assert.NotEmpty(t, updated.Data.Metrics, "There should be metrics")
	assert.NotNil(t, updated.Data.Metrics[accTWAMP], "There should be accedian-twamp metrics")
	assert.Nil(t, updated.Data.Metrics[accFLOW], "There should not be any accedian-flowmeter metrics")

	// Delete the record
	deleted, err := adminDB.DeleteIngestionDictionary()
	assert.Nil(t, err)
	assert.NotNil(t, deleted)
	assert.NotNil(t, deleted.XId)
	assert.NotNil(t, deleted.XRev)
	assert.NotEmpty(t, deleted.XId)
	assert.NotEmpty(t, deleted.XRev)
	assert.Equal(t, deleted, updated, "Deleted record not the same as last known version")

	// Get record - should fail
	fetched, err = adminDB.GetIngestionDictionary()
	assert.NotNil(t, err)
	assert.Nil(t, fetched)

	// Delete record - should fail as no record exists
	fetched, err = adminDB.DeleteIngestionDictionary()
	assert.NotNil(t, err)
	assert.Nil(t, fetched)
}