package couchDB

import (
	"fmt"
	"log"
	"testing"
	"time"

	"github.com/accedian/adh-gather/gather"
	"github.com/accedian/adh-gather/logger"
	"github.com/cenkalti/backoff"
	"github.com/spf13/viper"

	"github.com/leesper/couchdb-golang"

	dstest "github.com/accedian/adh-gather/datastore/test"
	dstestAdmin "github.com/accedian/adh-gather/datastore/test/admin"
	dstestTenant "github.com/accedian/adh-gather/datastore/test/tenant"
)

const (
	adminDBName = "adh-admin"
)

var (
	adminDB  *AdminServiceDatastoreCouchDB
	tenantDB *TenantServiceDatastoreCouchDB
)

func setupCouchDB() *couchdb.Server {
	// Configure the test AdminService DAO to use the newly started couch docker image
	cfg := gather.LoadConfig("../../config/adh-gather-test.yml", viper.New())

	// Before the tests run, setup the adh-admin db
	couchHost := cfg.GetString(gather.CK_server_datastore_ip.String())
	couchPort := cfg.GetString(gather.CK_server_datastore_port.String())

	couchServer, err := couchdb.NewServer(fmt.Sprintf("%s:%s", couchHost, couchPort))
	if err != nil {
		log.Fatalf("error connecting to couch server: %s", err.Error())
	}

	b := backoff.NewExponentialBackOff()
	b.MaxElapsedTime = 3 * time.Minute

	err = backoff.Retry(func() error {
		ver, err := couchServer.Version()
		logger.Log.Debugf("Test CouchDB version is: %s", ver)
		return err
	}, b)
	if err != nil {
		log.Fatalf("error connecting to couch server: %s", err.Error())
	}

	// Couch Run.
	dstest.ClearCouch(couchServer)
	adminDB, err = CreateAdminServiceDAO()
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

	tenantDB, err = CreateTenantServiceDAO()
	if err != nil {
		log.Fatalf("Could not create couchdb tenant DAO: %s", err.Error())
	}

	return couchServer
}

func TestCouchDBImplMain(t *testing.T) {
	couchServer := setupCouchDB()
	defer dstest.ClearCouch(couchServer)

	// RunAdminServiceDatastoreTests(t)
	RunTenantServiceDatastoreTests(t)
}

func RunAdminServiceDatastoreTests(t *testing.T) {
	// Issue the test to the AdminServiceDatastoreTestRunner
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
}
