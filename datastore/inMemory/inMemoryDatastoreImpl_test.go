package inMemory

import (
	"log"
	"testing"

	ds "github.com/accedian/adh-gather/datastore"

	dstestAdmin "github.com/accedian/adh-gather/datastore/test/admin"
	dstestTenant "github.com/accedian/adh-gather/datastore/test/tenant"
)

var (
	adminDB  ds.AdminServiceDatastore
	tenantDB ds.TenantServiceDatastore
)

func setupInMemoryDB() {
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

	RunAdminServiceDatastoreTests(t)
	RunTenantServiceDatastoreTests(t)
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
}
