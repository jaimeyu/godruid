package test

import (
	"log"

	"github.com/accedian/adh-gather/logger"
	couchdb "github.com/leesper/couchdb-golang"
)

// ClearCouch - Helper method to clear all couch data.
func ClearCouch(couchServer *couchdb.Server) {
	if couchServer != nil {
		dbs, err := couchServer.DBs()
		if err != nil {
			log.Fatalf("Could not delete DBs after test: %s", err.Error())
		}
		for _, dbname := range dbs {
			logger.Log.Debugf("Deleting DB %s", dbname)
			couchServer.Delete(dbname)
		}
	}
}
