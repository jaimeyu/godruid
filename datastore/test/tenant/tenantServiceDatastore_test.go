package tenant

// import (
// 	"fmt"
// 	"log"
// 	"os"
// 	"testing"
// 	"time"

// 	"github.com/accedian/adh-gather/gather"
// 	"github.com/accedian/adh-gather/logger"
// 	"github.com/cenkalti/backoff"
// 	"github.com/getlantern/deepcopy"
// 	"github.com/spf13/viper"

// 	"github.com/accedian/adh-gather/config"
// 	"github.com/leesper/couchdb-golang"
// 	"github.com/stretchr/testify/assert"

// 	ds "github.com/accedian/adh-gather/datastore"
// 	couchDB "github.com/accedian/adh-gather/datastore/couchDB"
// 	mem "github.com/accedian/adh-gather/datastore/inMemory"
// 	dstest "github.com/accedian/adh-gather/datastore/test"
// 	pb "github.com/accedian/adh-gather/gathergrpc"
// )

// const (
// 	adminDBName = "adh-admin"
// 	TENANT      = "tenant"
// )

// var (
// 	couchHost   string
// 	couchPort   string
// 	couchServer *couchdb.Server
// 	cfg         config.Provider
// 	adminDB     ds.AdminServiceDatastore
// 	tenantDB    ds.TenantServiceDatastore
// )

// func TestMain(m *testing.M) {
// 	// Configure the test AdminService DAO to use the newly started couch docker image
// 	cfg := gather.LoadConfig("../../../config/adh-gather-test.yml", viper.New())

// 	// Before the tests run, setup the adh-admin db
// 	couchHost = cfg.GetString(gather.CK_server_datastore_ip.String())
// 	couchPort = cfg.GetString(gather.CK_server_datastore_port.String())

// 	couchServer, err := couchdb.NewServer(fmt.Sprintf("%s:%s", couchHost, couchPort))
// 	if err != nil {
// 		log.Fatalf("error connecting to couch server: %s", err.Error())
// 	}

// 	b := backoff.NewExponentialBackOff()
// 	b.MaxElapsedTime = 3 * time.Minute

// 	err = backoff.Retry(func() error {
// 		ver, err := couchServer.Version()
// 		logger.Log.Debugf("Test CouchDB version is: %s", ver)
// 		return err
// 	}, b)
// 	if err != nil {
// 		log.Fatalf("error connecting to couch server: %s", err.Error())
// 	}

// 	// Couch Run.
// 	dstest.ClearCouch(couchServer)
// 	adminDB, err = couchDB.CreateAdminServiceDAO()
// 	if err != nil {
// 		log.Fatalf("Could not create couchdb admin DAO: %s", err.Error())
// 	}
// 	_, err = adminDB.CreateDatabase(adminDBName)
// 	if err != nil {
// 		log.Fatalf("Could not create admin DB: %s", err.Error())
// 	}
// 	err = adminDB.AddAdminViews()
// 	if err != nil {
// 		log.Fatalf("Could not populate admin indicies: %s", err.Error())
// 	}
// 	// Setup the Tenant DB:
// 	_, err = adminDB.CreateDatabase(ds.PrependToDataID(TENANT, string(ds.TenantDescriptorType)))
// 	if err != nil {
// 		log.Fatalf("Could not create Tenant DB: %s", err.Error())
// 	}

// 	tenantDB, err = couchDB.CreateTenantServiceDAO()
// 	if err != nil {
// 		log.Fatalf("Could not create couchdb tenant DAO: %s", err.Error())
// 	}

// 	code := m.Run()

// 	dstest.ClearCouch(couchServer)

// 	// If there were test failures, stop executing
// 	if code != 0 {
// 		os.Exit(code)
// 	}

// 	// InMemory Run:
// 	adminDB, err = mem.CreateAdminServiceDAO()
// 	if err != nil {
// 		log.Fatalf("Could not create in-mem admin DAO: %s", err.Error())
// 	}
// 	tenantDB, err = mem.CreateTenantServiceDAO()
// 	if err != nil {
// 		log.Fatalf("Could not create in-mem tenant DAO: %s", err.Error())
// 	}

// 	code = m.Run()
// 	os.Exit(code)

// }

// func TestTenantUserCRUD(t *testing.T) {
// 	defer dstest.FailButContinue("TestTenantUserCRUD")

// 	const USER1 = "test1"
// 	const USER2 = "test2"
// 	const PASS1 = "pass1"
// 	const PASS2 = "pass2"
// 	const PASS3 = "pass3"
// 	const TOKEN1 = "token1"
// 	const TOKEN2 = "token2"

// 	// Validate that there are currently no records
// 	tenantUserList, err := tenantDB.GetAllTenantUsers(TENANT)
// 	assert.Nil(t, err)
// 	assert.NotNil(t, tenantUserList)
// 	assert.Empty(t, tenantUserList.Data)

// 	// Create a record
// 	tenantUser := pb.TenantUserData{
// 		Username:        USER1,
// 		Password:        PASS1,
// 		OnboardingToken: TOKEN1,
// 		TenantId:        TENANT,
// 		State:           pb.UserState_INVITED}
// 	created, err := tenantDB.CreateTenantUser(&pb.TenantUser{Data: &tenantUser})
// 	assert.Nil(t, err)
// 	assert.NotNil(t, created)
// 	assert.NotNil(t, created.XId)
// 	assert.NotNil(t, created.XRev)
// 	assert.NotEmpty(t, created.XId)
// 	assert.NotEmpty(t, created.XRev)
// 	assert.Equal(t, string(ds.TenantUserType), created.Data.Datatype)
// 	assert.Equal(t, created.Data.Username, USER1, "Username not the same")
// 	assert.Equal(t, created.Data.Password, PASS1, "Password not the same")
// 	assert.Equal(t, created.Data.TenantId, TENANT, "Tenant ID not the same")
// 	assert.Equal(t, created.Data.OnboardingToken, TOKEN1, "OnboardingToken not the same")
// 	assert.True(t, created.Data.CreatedTimestamp > 0, "CreatedTimestamp was not set")
// 	assert.True(t, created.Data.LastModifiedTimestamp > 0, "LastmodifiedTimestamp was not set")

// 	// Get a record
// 	fetched, err := tenantDB.GetTenantUser(&pb.TenantUserIdRequest{TenantId: TENANT, UserId: created.XId})
// 	assert.Nil(t, err)
// 	assert.Equal(t, created, fetched, "The retrieved record should be the same as the created record")

// 	time.Sleep(time.Millisecond * 2)

// 	// Update a record
// 	updateRecord := pb.TenantUser{}
// 	deepcopy.Copy(&updateRecord, fetched)
// 	updateRecord.Data.Password = PASS2
// 	updated, err := tenantDB.UpdateTenantUser(&updateRecord)
// 	assert.Nil(t, err)
// 	assert.NotNil(t, updated)
// 	assert.Equal(t, updated.XId, fetched.XId)
// 	assert.NotEqual(t, updated.XRev, fetched.XRev)
// 	assert.Equal(t, string(ds.TenantUserType), updated.Data.Datatype)
// 	assert.Equal(t, updated.Data.Username, USER1, "Username not the same")
// 	assert.Equal(t, updated.Data.Password, PASS2, "Password was not updated")
// 	assert.Equal(t, updated.Data.TenantId, TENANT, "Tenant ID not the same")
// 	assert.Equal(t, updated.Data.OnboardingToken, TOKEN1, "OnboardingToken not the same")
// 	assert.Equal(t, updated.Data.CreatedTimestamp, fetched.Data.CreatedTimestamp, "CreatedTimestamp should not be updated")
// 	assert.True(t, updated.Data.LastModifiedTimestamp > fetched.Data.LastModifiedTimestamp, "LastmodifiedTimestamp was not updated")

// 	// Add a second record.
// 	tenantUser2 := pb.TenantUserData{
// 		Username:        USER2,
// 		Password:        PASS3,
// 		OnboardingToken: TOKEN2,
// 		TenantId:        TENANT,
// 		State:           pb.UserState_INVITED}
// 	created2, err := tenantDB.CreateTenantUser(&pb.TenantUser{Data: &tenantUser2})
// 	assert.Nil(t, err)
// 	assert.NotNil(t, created2)
// 	assert.NotNil(t, created2.XId)
// 	assert.NotNil(t, created2.XRev)
// 	assert.NotEmpty(t, created2.XId)
// 	assert.NotEmpty(t, created2.XRev)
// 	assert.Equal(t, string(ds.TenantUserType), created2.Data.Datatype)
// 	assert.Equal(t, created2.Data.Username, USER2, "Username not the same")
// 	assert.Equal(t, created2.Data.Password, PASS3, "Password not the same")
// 	assert.Equal(t, created2.Data.TenantId, TENANT, "Tenant ID not the same")
// 	assert.Equal(t, created2.Data.OnboardingToken, TOKEN2, "OnboardingToken not the same")
// 	assert.True(t, created2.Data.CreatedTimestamp > 0, "CreatedTimestamp was not set")
// 	assert.True(t, created2.Data.LastModifiedTimestamp > 0, "LastmodifiedTimestamp was not set")

// 	// Get all records
// 	fetchedList, err := tenantDB.GetAllTenantUsers(TENANT)
// 	assert.Nil(t, err)
// 	assert.NotNil(t, fetchedList)
// 	assert.NotEmpty(t, fetchedList.Data)
// 	assert.True(t, len(fetchedList.Data) == 2)

// 	// Delete a record.
// 	deleted, err := tenantDB.DeleteTenantUser(&pb.TenantUserIdRequest{TenantId: TENANT, UserId: fetched.XId})
// 	assert.Nil(t, err)
// 	assert.NotNil(t, deleted)
// 	assert.NotNil(t, deleted.XId)
// 	assert.NotNil(t, deleted.XRev)
// 	assert.NotEmpty(t, deleted.XId)
// 	assert.NotEmpty(t, deleted.XRev)
// 	assert.Equal(t, deleted.Data.Username, fetched.Data.Username, "Deleted Username not the same")

// 	// Get all records - should be 1
// 	fetchedList, err = tenantDB.GetAllTenantUsers(TENANT)
// 	assert.Nil(t, err)
// 	assert.NotNil(t, fetchedList)
// 	assert.NotEmpty(t, fetchedList.Data)
// 	assert.True(t, len(fetchedList.Data) == 1)

// 	// Get a record that does not exist
// 	dne, err := tenantDB.GetTenantUser(&pb.TenantUserIdRequest{TenantId: TENANT, UserId: deleted.XId})
// 	assert.NotNil(t, err)
// 	assert.Nil(t, dne)

// 	// Delete a record that oes not exist
// 	deleteDNE, err := tenantDB.DeleteTenantUser(&pb.TenantUserIdRequest{TenantId: TENANT, UserId: deleted.XId})
// 	assert.NotNil(t, err)
// 	assert.Nil(t, deleteDNE)

// 	// Delete the last record
// 	deleted, err = tenantDB.DeleteTenantUser(&pb.TenantUserIdRequest{TenantId: TENANT, UserId: created2.XId})
// 	assert.Nil(t, err)
// 	assert.NotNil(t, deleted)
// 	assert.NotNil(t, deleted.XId)
// 	assert.NotNil(t, deleted.XRev)
// 	assert.NotEmpty(t, deleted.XId)
// 	assert.NotEmpty(t, deleted.XRev)
// 	assert.Equal(t, deleted.Data.Username, created2.Data.Username, "Deleted Username not the same")

// 	// Get all records - should be empty
// 	fetchedList, err = tenantDB.GetAllTenantUsers(TENANT)
// 	assert.Nil(t, err)
// 	assert.NotNil(t, fetchedList)
// 	assert.Empty(t, fetchedList.Data)
// }

// func TestTenantDomainCRUD(t *testing.T) {
// 	defer dstest.FailButContinue("TestTenantDomainCRUD")

// 	const DOM1 = "domain1"
// 	const DOM2 = "domain2"
// 	const DOM3 = "domain3"
// 	const COLOR1 = "color1"
// 	const COLOR2 = "color2"
// 	const THRPRF = "ThresholdPrf"

// 	// Validate that there are currently no records
// 	recList, err := tenantDB.GetAllTenantDomains(TENANT)
// 	assert.Nil(t, err)
// 	assert.NotNil(t, recList)
// 	assert.Empty(t, recList.Data)

// 	// Create a record
// 	tenantDomain := pb.TenantDomainData{
// 		Name:                DOM1,
// 		TenantId:            TENANT,
// 		Color:               COLOR1,
// 		ThresholdProfileSet: []string{THRPRF}}
// 	created, err := tenantDB.CreateTenantDomain(&pb.TenantDomain{Data: &tenantDomain})
// 	assert.Nil(t, err)
// 	assert.NotNil(t, created)
// 	assert.NotNil(t, created.XId)
// 	assert.NotNil(t, created.XRev)
// 	assert.NotEmpty(t, created.XId)
// 	assert.NotEmpty(t, created.XRev)
// 	assert.Equal(t, string(ds.TenantDomainType), created.Data.Datatype)
// 	assert.Equal(t, created.Data.Name, DOM1, "Name not the same")
// 	assert.Equal(t, created.Data.Color, COLOR1, "Color not the same")
// 	assert.Equal(t, created.Data.ThresholdProfileSet[0], THRPRF, "Threshold Profile ID not the same")
// 	assert.Equal(t, created.Data.TenantId, TENANT, "Tenant ID not the same")
// 	assert.True(t, created.Data.CreatedTimestamp > 0, "CreatedTimestamp was not set")
// 	assert.True(t, created.Data.LastModifiedTimestamp > 0, "LastmodifiedTimestamp was not set")

// 	// Get a record
// 	fetched, err := tenantDB.GetTenantDomain(&pb.TenantDomainIdRequest{TenantId: TENANT, DomainId: created.XId})
// 	assert.Nil(t, err)
// 	assert.Equal(t, created, fetched, "The retrieved record should be the same as the created record")

// 	time.Sleep(time.Millisecond * 2)

// 	// Update a record
// 	updateRecord := pb.TenantDomain{}
// 	deepcopy.Copy(&updateRecord, fetched)
// 	updateRecord.Data.Color = COLOR2
// 	updated, err := tenantDB.UpdateTenantDomain(&updateRecord)
// 	assert.Nil(t, err)
// 	assert.NotNil(t, updated)
// 	assert.Equal(t, updated.XId, fetched.XId)
// 	assert.NotEqual(t, updated.XRev, fetched.XRev)
// 	assert.Equal(t, string(ds.TenantDomainType), updated.Data.Datatype)
// 	assert.Equal(t, updated.Data.Name, DOM1, "Name not the same")
// 	assert.Equal(t, updated.Data.Color, COLOR2, "Password was not updated")
// 	assert.Equal(t, updated.Data.TenantId, TENANT, "Tenant ID not the same")
// 	assert.Equal(t, updated.Data.ThresholdProfileSet[0], THRPRF, "Threshold Profile ID not the same")
// 	assert.Equal(t, updated.Data.CreatedTimestamp, fetched.Data.CreatedTimestamp, "CreatedTimestamp should not be updated")
// 	assert.True(t, updated.Data.LastModifiedTimestamp > fetched.Data.LastModifiedTimestamp, "LastmodifiedTimestamp was not updated")

// 	// Add a second record.
// 	tenantDomain2 := pb.TenantDomainData{
// 		Name:     DOM2,
// 		TenantId: TENANT,
// 		Color:    COLOR1}
// 	created2, err := tenantDB.CreateTenantDomain(&pb.TenantDomain{Data: &tenantDomain2})
// 	assert.Nil(t, err)
// 	assert.NotNil(t, created2)
// 	assert.NotNil(t, created2.XId)
// 	assert.NotNil(t, created2.XRev)
// 	assert.NotEmpty(t, created2.XId)
// 	assert.NotEmpty(t, created2.XRev)
// 	assert.Equal(t, string(ds.TenantDomainType), created2.Data.Datatype)
// 	assert.Equal(t, created2.Data.Name, DOM2, "Name not the same")
// 	assert.Equal(t, created2.Data.Color, COLOR1, "Password not the same")
// 	assert.Equal(t, created2.Data.TenantId, TENANT, "Tenant ID not the same")
// 	assert.True(t, len(created2.Data.ThresholdProfileSet) == 0, "Should not be a Threshold Profile ID")
// 	assert.True(t, created2.Data.CreatedTimestamp > 0, "CreatedTimestamp was not set")
// 	assert.True(t, created2.Data.LastModifiedTimestamp > 0, "LastmodifiedTimestamp was not set")

// 	// Get all records
// 	fetchedList, err := tenantDB.GetAllTenantDomains(TENANT)
// 	assert.Nil(t, err)
// 	assert.NotNil(t, fetchedList)
// 	assert.NotEmpty(t, fetchedList.Data)
// 	assert.True(t, len(fetchedList.Data) == 2)

// 	// Delete a record.
// 	deleted, err := tenantDB.DeleteTenantDomain(&pb.TenantDomainIdRequest{TenantId: TENANT, DomainId: fetched.XId})
// 	assert.Nil(t, err)
// 	assert.NotNil(t, deleted)
// 	assert.NotNil(t, deleted.XId)
// 	assert.NotNil(t, deleted.XRev)
// 	assert.NotEmpty(t, deleted.XId)
// 	assert.NotEmpty(t, deleted.XRev)
// 	assert.Equal(t, deleted.Data.Name, fetched.Data.Name, "Deleted name not the same")

// 	// Get all records - should be 1
// 	fetchedList, err = tenantDB.GetAllTenantDomains(TENANT)
// 	assert.Nil(t, err)
// 	assert.NotNil(t, fetchedList)
// 	assert.NotEmpty(t, fetchedList.Data)
// 	assert.True(t, len(fetchedList.Data) == 1)

// 	// Get a record that does not exist
// 	dne, err := tenantDB.GetTenantDomain(&pb.TenantDomainIdRequest{TenantId: TENANT, DomainId: deleted.XId})
// 	assert.NotNil(t, err)
// 	assert.Nil(t, dne)

// 	// Delete a record that oes not exist
// 	deleteDNE, err := tenantDB.DeleteTenantDomain(&pb.TenantDomainIdRequest{TenantId: TENANT, DomainId: deleted.XId})
// 	assert.NotNil(t, err)
// 	assert.Nil(t, deleteDNE)

// 	// Delete the last record
// 	deleted, err = tenantDB.DeleteTenantDomain(&pb.TenantDomainIdRequest{TenantId: TENANT, DomainId: created2.XId})
// 	assert.Nil(t, err)
// 	assert.NotNil(t, deleted)
// 	assert.NotNil(t, deleted.XId)
// 	assert.NotNil(t, deleted.XRev)
// 	assert.NotEmpty(t, deleted.XId)
// 	assert.NotEmpty(t, deleted.XRev)
// 	assert.Equal(t, deleted.Data.Name, created2.Data.Name, "Deleted name not the same")

// 	// Get all records - should be empty
// 	fetchedList, err = tenantDB.GetAllTenantDomains(TENANT)
// 	assert.Nil(t, err)
// 	assert.NotNil(t, fetchedList)
// 	assert.Empty(t, fetchedList.Data)
// }

// // func TestTenantDescCRUD(t *testing.T) {
// // 	defer failButContinue("TestTenantDescCRUD")

// // 	const COMPANY1 = "test1"
// // 	const COMPANY2 = "test2"
// // 	const SUBDOMAIN1 = "pass1"
// // 	const SUBDOMAIN2 = "pass2"
// // 	const SUBDOMAIN3 = "pass3"

// // 	// Validate that there are currently no records
// // 	tenantList, err := adminDB.GetAllTenantDescriptors()
// // 	assert.Nil(t, err)
// // 	assert.NotNil(t, tenantList)
// // 	assert.Empty(t, tenantList.Data)

// // 	// Create a record
// // 	tenant1 := pb.TenantDescriptorData{
// // 		Name:         COMPANY1,
// // 		UrlSubdomain: SUBDOMAIN1,
// // 		State:        pb.UserState_ACTIVE}
// // 	created, err := adminDB.CreateTenant(&pb.TenantDescriptor{Data: &tenant1})
// // 	assert.Nil(t, err)
// // 	assert.NotNil(t, created)
// // 	assert.NotNil(t, created.XId)
// // 	assert.NotNil(t, created.XRev)
// // 	assert.NotEmpty(t, created.XId)
// // 	assert.NotEmpty(t, created.XRev)
// // 	assert.Equal(t, string(ds.TenantDescriptorType), created.Data.Datatype)
// // 	assert.Equal(t, created.Data.Name, COMPANY1, "Name not the same")
// // 	assert.Equal(t, created.Data.UrlSubdomain, SUBDOMAIN1, "Subdomain not the same")
// // 	assert.True(t, created.Data.CreatedTimestamp > 0, "CreatedTimestamp was not set")
// // 	assert.True(t, created.Data.LastModifiedTimestamp > 0, "LastmodifiedTimestamp was not set")

// // 	// Get a record
// // 	fetched, err := adminDB.GetTenantDescriptor(created.XId)
// // 	assert.Nil(t, err)
// // 	assert.Equal(t, created, fetched, "The retrieved record should be the same as the updated record")

// // 	time.Sleep(time.Millisecond * 2)

// // 	// Update a record
// // 	updateRecord := pb.TenantDescriptor{}
// // 	deepcopy.Copy(&updateRecord, fetched)
// // 	updateRecord.Data.UrlSubdomain = SUBDOMAIN2
// // 	updated, err := adminDB.UpdateTenantDescriptor(&updateRecord)
// // 	assert.Nil(t, err)
// // 	assert.NotNil(t, updated)
// // 	assert.Equal(t, updated.XId, fetched.XId)
// // 	assert.NotEqual(t, updated.XRev, fetched.XRev)
// // 	assert.Equal(t, string(ds.TenantDescriptorType), updated.Data.Datatype)
// // 	assert.Equal(t, updated.Data.Name, COMPANY1, "Name not the same")
// // 	assert.Equal(t, updated.Data.UrlSubdomain, SUBDOMAIN2, "Subdomain was not updated")
// // 	assert.Equal(t, updated.Data.CreatedTimestamp, fetched.Data.CreatedTimestamp, "CreatedTimestamp should not be updated")
// // 	assert.True(t, updated.Data.LastModifiedTimestamp > fetched.Data.LastModifiedTimestamp, "LastmodifiedTimestamp was not updated")

// // 	// Add a second record.
// // 	tenant2 := pb.TenantDescriptorData{
// // 		Name:         COMPANY2,
// // 		UrlSubdomain: SUBDOMAIN3,
// // 		State:        pb.UserState_ACTIVE}
// // 	created2, err := adminDB.CreateTenant(&pb.TenantDescriptor{Data: &tenant2})
// // 	assert.Nil(t, err)
// // 	assert.NotNil(t, created2)
// // 	assert.NotNil(t, created2.XId)
// // 	assert.NotNil(t, created2.XRev)
// // 	assert.NotEmpty(t, created2.XId)
// // 	assert.NotEmpty(t, created2.XRev)
// // 	assert.Equal(t, string(ds.TenantDescriptorType), created2.Data.Datatype)
// // 	assert.Equal(t, created2.Data.Name, COMPANY2, "Name not the same")
// // 	assert.Equal(t, created2.Data.UrlSubdomain, SUBDOMAIN3, "Subdomain not the same")
// // 	assert.True(t, created2.Data.CreatedTimestamp > 0, "CreatedTimestamp was not set")
// // 	assert.True(t, created2.Data.LastModifiedTimestamp > 0, "LastmodifiedTimestamp was not set")

// // 	// Get all records
// // 	fetchedList, err := adminDB.GetAllTenantDescriptors()
// // 	assert.Nil(t, err)
// // 	assert.NotNil(t, fetchedList)
// // 	assert.NotEmpty(t, fetchedList.Data)
// // 	assert.True(t, len(fetchedList.Data) == 2)

// // 	// Fetch a Tenant ID by username
// // 	tenantID, err := adminDB.GetTenantIDByAlias(COMPANY1)
// // 	assert.Nil(t, err)
// // 	assert.NotEmpty(t, tenantID)
// // 	assert.Equal(t, updated.XId, tenantID)

// // 	// Delete a record that does not exist.
// // 	deleted, err := adminDB.DeleteTenant(string(fetched.XId))
// // 	assert.Nil(t, err)
// // 	assert.NotNil(t, deleted)
// // 	assert.NotNil(t, deleted.XId)
// // 	assert.NotNil(t, deleted.XRev)
// // 	assert.NotEmpty(t, deleted.XId)
// // 	assert.NotEmpty(t, deleted.XRev)
// // 	assert.Equal(t, deleted.Data.Name, fetched.Data.Name, "Deleted Name not the same")

// // 	// Get all records - should be 1
// // 	fetchedList, err = adminDB.GetAllTenantDescriptors()
// // 	assert.Nil(t, err)
// // 	assert.NotNil(t, fetchedList)
// // 	assert.NotEmpty(t, fetchedList.Data)
// // 	assert.True(t, len(fetchedList.Data) == 1)

// // 	// Get a record that does not exist
// // 	dne, err := adminDB.GetTenantDescriptor(deleted.XId)
// // 	assert.NotNil(t, err)
// // 	assert.Nil(t, dne)

// // 	// Delete a record that oes not exist
// // 	deleteDNE, err := adminDB.DeleteTenant(deleted.XId)
// // 	assert.NotNil(t, err)
// // 	assert.Nil(t, deleteDNE)

// // 	// Delete the last record
// // 	deleted, err = adminDB.DeleteTenant(string(created2.XId))
// // 	assert.Nil(t, err)
// // 	assert.NotNil(t, deleted)
// // 	assert.NotNil(t, deleted.XId)
// // 	assert.NotNil(t, deleted.XRev)
// // 	assert.NotEmpty(t, deleted.XId)
// // 	assert.NotEmpty(t, deleted.XRev)
// // 	assert.Equal(t, deleted.Data.Name, created2.Data.Name, "Deleted Name not the same")

// // 	// Get all records - should be empty
// // 	fetchedList, err = adminDB.GetAllTenantDescriptors()
// // 	assert.Nil(t, err)
// // 	assert.NotNil(t, fetchedList)
// // 	assert.Empty(t, fetchedList.Data)
// // }

// // func TestIngDictCRUD(t *testing.T) {
// // 	defer failButContinue("TestIngDictCRUD")

// // 	var accTWAMP = string(handlers.AccedianTwamp)
// // 	var accFLOW = string(handlers.AccedianFlowmeter)

// // 	var delayMin = "delayMin"
// // 	var delayMax = "delayMax"
// // 	var delayAvg = "delayAvg"

// // 	var throughputAvg = "throughputAvg"
// // 	var throughputMax = "throughputMax"
// // 	var throughputMin = "throughputMin"

// // 	// Validate that there are currently no records
// // 	ingPrf, err := adminDB.GetIngestionDictionary()
// // 	assert.NotNil(t, err)
// // 	assert.Nil(t, ingPrf)

// // 	// Read in the test dictionary from file
// // 	defaultDictionaryBytes, err := ioutil.ReadFile("./files/testIngestionDictionary.json")
// // 	if err != nil {
// // 		logger.Log.Fatalf("Unable to read Default Ingestion Profile from file: %s", err.Error())
// // 	}

// // 	defaultDictionaryData := &pb.IngestionDictionaryData{}
// // 	if err = json.Unmarshal(defaultDictionaryBytes, &defaultDictionaryData); err != nil {
// // 		logger.Log.Fatalf("Unable to construct Default Ingestion Profile from file: %s", err.Error())
// // 	}

// // 	// Create a record
// // 	created, err := adminDB.CreateIngestionDictionary(&pb.IngestionDictionary{Data: defaultDictionaryData})
// // 	assert.Nil(t, err)
// // 	assert.NotNil(t, created)
// // 	assert.NotNil(t, created.XId)
// // 	assert.NotNil(t, created.XRev)
// // 	assert.NotEmpty(t, created.XId)
// // 	assert.NotEmpty(t, created.XRev)
// // 	assert.Equal(t, string(ds.IngestionDictionaryType), created.Data.Datatype)
// // 	assert.NotEmpty(t, created.Data.Metrics, "There should be metrics")
// // 	assert.NotNil(t, created.Data.Metrics[accTWAMP], "There should be accedian-twamp metrics")
// // 	assert.NotNil(t, created.Data.Metrics[accFLOW], "There should be accedian-flowmeter metrics")
// // 	assert.NotEmpty(t, created.Data.Metrics[accTWAMP].MetricMap, "There should be accedian-twamp metric definitions")
// // 	assert.NotNil(t, created.Data.Metrics[accTWAMP].MetricMap[delayMin], "There should be delayMin metrics")
// // 	assert.NotNil(t, created.Data.Metrics[accTWAMP].MetricMap[delayMax], "There should be delayMax metrics")
// // 	assert.NotNil(t, created.Data.Metrics[accTWAMP].MetricMap[delayAvg], "There should be delayAvg metrics")
// // 	assert.NotEmpty(t, created.Data.Metrics[accTWAMP].MetricMap[delayMin].MonitoredObjectTypes, "There should be delayMin monitored object definitions")
// // 	assert.True(t, len(created.Data.Metrics[accTWAMP].MetricMap[delayMin].MonitoredObjectTypes) == 3, "There should be 3 delayMin monitored object definitions")
// // 	assert.NotEmpty(t, created.Data.Metrics[accFLOW].MetricMap, "There should be accedian-flowmeter metric definitions")
// // 	assert.NotNil(t, created.Data.Metrics[accFLOW].MetricMap[throughputAvg], "There should be throughputAvg metrics")
// // 	assert.NotNil(t, created.Data.Metrics[accFLOW].MetricMap[throughputMax], "There should be throughputMax metrics")
// // 	assert.NotNil(t, created.Data.Metrics[accFLOW].MetricMap[throughputMin], "There should be throughputMin metrics")
// // 	assert.NotEmpty(t, created.Data.Metrics[accFLOW].MetricMap[throughputMax].MonitoredObjectTypes, "There should be throughputMax monitored object definitions")
// // 	assert.True(t, len(created.Data.Metrics[accFLOW].MetricMap[throughputMax].MonitoredObjectTypes) == 1, "There should be 1 throughputMax monitored object definitions")

// // 	// Get a record
// // 	fetched, err := adminDB.GetIngestionDictionary()
// // 	assert.Nil(t, err)
// // 	assert.Equal(t, created, fetched, "The retrieved record should be the same as the created record")

// // 	// Try to create a record that already exists, should fail
// // 	created, err = adminDB.CreateIngestionDictionary(created)
// // 	assert.NotNil(t, err)
// // 	assert.Nil(t, created, "Created should now be nil")

// // 	time.Sleep(time.Millisecond * 2)

// // 	// Update a record
// // 	updateRecord := pb.IngestionDictionary{}
// // 	deepcopy.Copy(&updateRecord, fetched)
// // 	delete(updateRecord.Data.Metrics, accFLOW)
// // 	updated, err := adminDB.UpdateIngestionDictionary(&updateRecord)
// // 	assert.Nil(t, err)
// // 	assert.NotNil(t, updated)
// // 	assert.NotNil(t, updated.XId)
// // 	assert.NotNil(t, updated.XRev)
// // 	assert.NotEmpty(t, updated.XId)
// // 	assert.NotEmpty(t, updated.XRev)
// // 	assert.Equal(t, fetched.XId, updated.XId, "Id values should be the same")
// // 	assert.NotEqual(t, fetched.XRev, updated.XRev)
// // 	assert.Equal(t, string(ds.IngestionDictionaryType), updated.Data.Datatype)
// // 	assert.NotEmpty(t, updated.Data.Metrics, "There should be metrics")
// // 	assert.NotNil(t, updated.Data.Metrics[accTWAMP], "There should be accedian-twamp metrics")
// // 	assert.Nil(t, updated.Data.Metrics[accFLOW], "There should not be any accedian-flowmeter metrics")

// // 	// Delete the record
// // 	deleted, err := adminDB.DeleteIngestionDictionary()
// // 	assert.Nil(t, err)
// // 	assert.NotNil(t, deleted)
// // 	assert.NotNil(t, deleted.XId)
// // 	assert.NotNil(t, deleted.XRev)
// // 	assert.NotEmpty(t, deleted.XId)
// // 	assert.NotEmpty(t, deleted.XRev)
// // 	assert.Equal(t, deleted, updated, "Deleted record not the same as last known version")

// // 	// Get record - should fail
// // 	fetched, err = adminDB.GetIngestionDictionary()
// // 	assert.NotNil(t, err)
// // 	assert.Nil(t, fetched)

// // 	// Delete record - should fail as no record exists
// // 	fetched, err = adminDB.DeleteIngestionDictionary()
// // 	assert.NotNil(t, err)
// // 	assert.Nil(t, fetched)
// // }

// // func TestValidTypesCRUD(t *testing.T) {
// // 	defer failButContinue("TestValidTypesCRUD")

// // 	objTypeKey1 := "objTypeKey1"
// // 	objTypeKey2 := "objTypeKey2"
// // 	devTypeKey1 := "devTypeKey1"
// // 	devTypeKey2 := "devTypeKey2"
// // 	devTypeKey3 := "devTypeKey3"

// // 	objTypeVal1 := "objTypeVal1"
// // 	objTypeVal2 := "objTypeVal2"
// // 	devTypeVal1 := "devTypeVal1"
// // 	devTypeVal2 := "devTypeVal2"
// // 	devTypeVal3 := "devTypeVal3"

// // 	// Validate that there are currently no records
// // 	validTypes, err := adminDB.GetValidTypes()
// // 	assert.NotNil(t, err)
// // 	assert.Nil(t, validTypes)

// // 	validTypeData := pb.ValidTypesData{
// // 		MonitoredObjectTypes:       map[string]string{objTypeKey1: objTypeVal1, objTypeKey2: objTypeVal2},
// // 		MonitoredObjectDeviceTypes: map[string]string{devTypeKey1: devTypeVal1, devTypeKey2: devTypeVal2}}
// // 	// Create a record
// // 	created, err := adminDB.CreateValidTypes(&pb.ValidTypes{Data: &validTypeData})
// // 	assert.Nil(t, err)
// // 	assert.NotNil(t, created)
// // 	assert.NotNil(t, created.XId)
// // 	assert.NotNil(t, created.XRev)
// // 	assert.NotEmpty(t, created.XId)
// // 	assert.NotEmpty(t, created.XRev)
// // 	assert.NotEmpty(t, created.Data.MonitoredObjectTypes, "There should be mon obj types")
// // 	assert.Equal(t, 2, len(created.Data.MonitoredObjectTypes), "There should be 2 mon obj types")
// // 	assert.NotEmpty(t, created.Data.MonitoredObjectDeviceTypes, "There should be mon obj dev types")
// // 	assert.Equal(t, 2, len(created.Data.MonitoredObjectDeviceTypes), "There should be 2 mon obj dev types")
// // 	assert.Equal(t, objTypeVal1, created.Data.MonitoredObjectTypes[objTypeKey1])
// // 	assert.Equal(t, devTypeVal1, created.Data.MonitoredObjectDeviceTypes[devTypeKey1])

// // 	// Get a record
// // 	fetched, err := adminDB.GetValidTypes()
// // 	assert.Nil(t, err)
// // 	assert.Equal(t, created, fetched, "The retrieved record should be the same as the created record")

// // 	// Try to create a record that already exists, should fail
// // 	created, err = adminDB.CreateValidTypes(created)
// // 	assert.NotNil(t, err)
// // 	assert.Nil(t, created, "Created should now be nil")

// // 	time.Sleep(time.Millisecond * 2)

// // 	// Update a record
// // 	updateRecord := pb.ValidTypes{}
// // 	deepcopy.Copy(&updateRecord, fetched)
// // 	updateRecord.Data.MonitoredObjectDeviceTypes[devTypeKey3] = devTypeVal3
// // 	updated, err := adminDB.UpdateValidTypes(&updateRecord)
// // 	assert.Nil(t, err)
// // 	assert.NotNil(t, updated)
// // 	assert.NotNil(t, updated.XId)
// // 	assert.NotNil(t, updated.XRev)
// // 	assert.NotEmpty(t, updated.XId)
// // 	assert.NotEmpty(t, updated.XRev)
// // 	assert.Equal(t, fetched.XId, updated.XId, "Id values should be the same")
// // 	assert.NotEqual(t, fetched.XRev, updated.XRev)
// // 	assert.NotEmpty(t, updated.Data.MonitoredObjectTypes, "There should be mon obj types")
// // 	assert.Equal(t, 2, len(updated.Data.MonitoredObjectTypes), "There should be 2 mon obj types")
// // 	assert.NotEmpty(t, updated.Data.MonitoredObjectDeviceTypes, "There should be mon obj dev types")
// // 	assert.Equal(t, 3, len(updated.Data.MonitoredObjectDeviceTypes), "There should be 2 mon obj dev types")
// // 	assert.Equal(t, devTypeVal3, updated.Data.MonitoredObjectDeviceTypes[devTypeKey3])

// // 	// Retrieve a specific valid type
// // 	specific, err := adminDB.GetSpecificValidTypes(&pb.ValidTypesRequest{MonitoredObjectTypes: true})
// // 	assert.Nil(t, err)
// // 	assert.NotNil(t, specific)
// // 	assert.NotNil(t, specific.MonitoredObjectTypes)
// // 	assert.Nil(t, specific.MonitoredObjectDeviceTypes)
// // 	assert.Equal(t, updated.Data.MonitoredObjectTypes, specific.MonitoredObjectTypes)

// // 	specific, err = adminDB.GetSpecificValidTypes(&pb.ValidTypesRequest{MonitoredObjectDeviceTypes: true})
// // 	assert.Nil(t, err)
// // 	assert.NotNil(t, specific)
// // 	assert.NotNil(t, specific.MonitoredObjectDeviceTypes)
// // 	assert.Nil(t, specific.MonitoredObjectTypes)
// // 	assert.Equal(t, updated.Data.MonitoredObjectDeviceTypes, specific.MonitoredObjectDeviceTypes)

// // 	specific, err = adminDB.GetSpecificValidTypes(&pb.ValidTypesRequest{MonitoredObjectDeviceTypes: false, MonitoredObjectTypes: false})
// // 	assert.Nil(t, err)
// // 	assert.NotNil(t, specific)
// // 	assert.Nil(t, specific.MonitoredObjectDeviceTypes)
// // 	assert.Nil(t, specific.MonitoredObjectTypes)

// // 	// Delete the record
// // 	deleted, err := adminDB.DeleteValidTypes()
// // 	assert.Nil(t, err)
// // 	assert.NotNil(t, deleted)
// // 	assert.NotNil(t, deleted.XId)
// // 	assert.NotNil(t, deleted.XRev)
// // 	assert.NotEmpty(t, deleted.XId)
// // 	assert.NotEmpty(t, deleted.XRev)
// // 	assert.Equal(t, deleted, updated, "Deleted record not the same as last known version")

// // 	// Get record - should fail
// // 	fetched, err = adminDB.GetValidTypes()
// // 	assert.NotNil(t, err)
// // 	assert.Nil(t, fetched)

// // 	// Delete record - should fail as no record exists
// // 	fetched, err = adminDB.DeleteValidTypes()
// // 	assert.NotNil(t, err)
// // 	assert.Nil(t, fetched)
// // }
