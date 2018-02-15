package test

import (
	"log"

	"github.com/accedian/adh-gather/logger"
	couchdb "github.com/leesper/couchdb-golang"
)

// ClearCouch - Helper method to clear all couch data.
func ClearCouch(couchServer *couchdb.Server) {
	dbs, err := couchServer.DBs()
	if err != nil {
		log.Fatalf("Could not delete DBs after test: %s", err.Error())
	}
	for _, dbname := range dbs {
		logger.Log.Debugf("Deleting DB %s", dbname)
		couchServer.Delete(dbname)
	}
}

// FailButContinue - Way to bypass the termination of the test executor before the DB can cleanup.
// Still exits the tests, but due to the os.Exit calls, will still stop execution.
func FailButContinue(testName string) {
	if r := recover(); r != nil {
		logger.Log.Debug("Failed Test %s", testName)
	}
}
