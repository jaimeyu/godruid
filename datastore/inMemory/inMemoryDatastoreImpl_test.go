package inMemory

import (
	"log"
	"testing"

	ds "github.com/accedian/adh-gather/datastore"
	"github.com/accedian/adh-gather/gather"
	tenmod "github.com/accedian/adh-gather/models/tenant"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"

	dstestAdmin "github.com/accedian/adh-gather/datastore/test/admin"
	dstestTenant "github.com/accedian/adh-gather/datastore/test/tenant"
)

var (
	adminDB  ds.AdminServiceDatastore
	tenantDB ds.TenantServiceDatastore
)

const (
	InvalidTenantDataType tenmod.TenantDataType = "invalid"
)

func setupInMemoryDB() {
	cfg := gather.LoadConfig("../../config/adh-gather-test.yml", viper.New())
	cfg.Set("ingDict", "../../files/defaultIngestionDictionary.json")

	var err error
	adminDB, err = CreateAdminServiceDAO()
	if err != nil {
		log.Fatalf("Could not create in-mem admin DAO: %s", err.Error())
	}

	tenantDB, err = CreateTenantServiceDAO()
	if err != nil {
		log.Fatalf("Could not create in-mem tenant DAO: %s", err.Error())
	}
}

func TestInMemoryImplMain(t *testing.T) {
	setupInMemoryDB()
	TestPackageSpecificFunctions(t)

	RunAdminServiceDatastoreTests(t)
	RunTenantServiceDatastoreTests(t)
}

func TestPackageSpecificFunctions(t *testing.T) {
	// Test failure condition of DoesTenantExist
	testDB, err := CreateTenantServiceDAO()
	assert.Nil(t, err)
	assert.NotNil(t, testDB)

	err = testDB.DoesTenantExist("", tenmod.TenantUserType)
	assert.NotNil(t, err, "Should not be able to validate a Tenant with no Tenant data")

	err = nil

	err = testDB.DoesTenantExist("something", InvalidTenantDataType)
	assert.NotNil(t, err, "Should not be able to validate a Tenant with an invalid datatype")
}

func RunAdminServiceDatastoreTests(t *testing.T) {

	tester := dstestAdmin.InitTestRunner(adminDB)
	tester.RunAdminUserCRUD(t)
	tester.RunTenantDescCRUD(t)
	tester.RunIngDictCRUD(t)
	tester.RunValidTypesCRUD(t)
}

func RunTenantServiceDatastoreTests(t *testing.T) {
	// Issue the test to the TenantServiceDatastoreTestRunner
	tester := dstestTenant.InitTestRunner(tenantDB, adminDB)
	tester.RunTenantUserCRUD(t)
	tester.RunTenantDomainCRUD(t)
	tester.RunTenantMonitoredObjectCRUD(t)
	tester.RunTenantMetadataCRUD(t)
	tester.RunTenantIngestionProfileCRUD(t)
	//TODO: re-enable this test when the threshold profile model is fixed to no include UI details or UI uses the API to provision them
	//tester.RunTenantThresholdProfileCRUD(t)
	tester.RunGetMonitoredObjectByDomainMapTest(t)
	tester.RunTenantDataCleaningProfileCRUD(t)

	tester.RunDashboardCRUD(t)
	tester.RunCardCRUD(t)
}
