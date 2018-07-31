package restapi

import (
	"fmt"
	"net/http"

	"github.com/accedian/adh-gather/config"
	"github.com/accedian/adh-gather/datastore"
	"github.com/accedian/adh-gather/gather"
	"github.com/accedian/adh-gather/handlers"
	"github.com/accedian/adh-gather/logger"
	"github.com/accedian/adh-gather/monitoring"
	"github.com/accedian/adh-gather/profile"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

const (
	globalChangesDBName = "_global_changes"
	replicatorDBName    = "_replicator"
	metadataDBName      = "_metadata"
	usersDBName         = "_users"
	statsDBName         = "_stats"
)

func startMonitoring(cfg config.Provider) {
	restBindIP := cfg.GetString(gather.CK_server_rest_ip.String())
	monPort := cfg.GetInt(gather.CK_server_monitoring_port.String())

	monitoring.InitMetrics()
	promServerMux := http.NewServeMux()

	logger.Log.Infof("Starting Prometheus Server")
	promServerMux.Handle("/metrics", promhttp.Handler())
	addr := fmt.Sprintf("%s:%d", restBindIP, monPort)
	if err := http.ListenAndServe(addr, promServerMux); err != nil {
		logger.Log.Fatalf("Unable to start monitoring function: %s", err.Error())
	}
}

func startProfile(cfg config.Provider) {
	restBindIP := cfg.GetString(gather.CK_server_rest_ip.String())
	monPort := cfg.GetInt(gather.CK_server_profile_port.String())

	pprofServerMux := http.NewServeMux()

	profile.AttachProfiler(pprofServerMux)

	logger.Log.Infof("Starting Profile Server")
	addr := fmt.Sprintf("%s:%d", restBindIP, monPort)
	if err := http.ListenAndServe(addr, pprofServerMux); err != nil {
		logger.Log.Fatalf("Unable to start profile function: %s", err.Error())
	}
}

func provisionCouchData(pouchSH *handlers.PouchDBPluginServiceHandler, adminDB datastore.AdminServiceDatastore, adminDBStr string, cfg config.Provider) {
	ensureBaseCouchDBsExist(pouchSH)
	ensureAdminDBExists(pouchSH, adminDB, adminDBStr)
}

func createCouchDB(pouchSH *handlers.PouchDBPluginServiceHandler, dbName string) error {
	// Make sure global changes db exists
	_, err := pouchSH.IsDBAvailable(dbName)
	if err != nil {
		logger.Log.Infof("Database %s does not exist. %s DB will now be created.", dbName, dbName)
		// Try to create the DB
		_, err = pouchSH.AddDB(dbName)
		if err != nil {
			return err
		}
	}

	return nil
}

func ensureBaseCouchDBsExist(pouchSH *handlers.PouchDBPluginServiceHandler) {
	if err := createCouchDB(pouchSH, globalChangesDBName); err != nil {
		logger.Log.Fatalf("Unable to create DB %s: %s", globalChangesDBName, err.Error())
	}
	if err := createCouchDB(pouchSH, metadataDBName); err != nil {
		logger.Log.Fatalf("Unable to create DB %s: %s", metadataDBName, err.Error())
	}
	if err := createCouchDB(pouchSH, replicatorDBName); err != nil {
		logger.Log.Fatalf("Unable to create DB %s: %s", replicatorDBName, err.Error())
	}
	if err := createCouchDB(pouchSH, usersDBName); err != nil {
		logger.Log.Fatalf("Unable to create DB %s: %s", usersDBName, err.Error())
	}
}

func ensureAdminDBExists(pouchSH *handlers.PouchDBPluginServiceHandler, adminDB datastore.AdminServiceDatastore, adminDBStr string) {
	_, err := pouchSH.IsDBAvailable(adminDBStr)
	if err != nil {
		logger.Log.Infof("Database %s does not exist. %s DB will now be created.", adminDBStr, adminDBStr)

		// Try to create the DB:
		_, err = pouchSH.AddDB(adminDBStr)
		if err != nil {
			logger.Log.Fatalf("Unable to create DB %s: %s", adminDBStr, err.Error())
		}

		// Also add the Views for Admin DB.
		err = adminDB.AddAdminViews()
		if err != nil {
			logger.Log.Fatalf("Unable to Add Views to DB %s: %s", adminDBStr, err.Error())
		}
	}

	logger.Log.Infof("Using %s as Administrative Database", adminDBStr)
}

func areStringSlicesEqual(slice1 []string, slice2 []string) bool {
	if (slice1 == nil && slice2 != nil) || (slice1 != nil && slice2 == nil) {
		return false
	}

	if slice1 == nil && slice2 == nil {
		return true
	}

	if len(slice1) != len(slice2) {
		return false
	}

	for _, value := range slice1 {
		if !gather.DoesSliceContainString(slice2, value) {
			return false
		}
	}

	return true
}
