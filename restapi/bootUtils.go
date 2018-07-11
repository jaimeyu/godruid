package restapi

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/accedian/adh-gather/config"
	"github.com/accedian/adh-gather/datastore"
	"github.com/accedian/adh-gather/gather"
	"github.com/accedian/adh-gather/handlers"
	"github.com/accedian/adh-gather/logger"
	admmod "github.com/accedian/adh-gather/models/admin"
	tenmod "github.com/accedian/adh-gather/models/tenant"
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

var (
	// ValidMonitoredObjectTypes - known Monitored Object types in the system.
	ValidMonitoredObjectTypes = map[string]tenmod.MonitoredObjectType{
		"pe": tenmod.TwampPE,
		"sf": tenmod.TwampSF,
		"sl": tenmod.TwampSL,
		string(tenmod.TwampPE): tenmod.TwampPE,
		string(tenmod.TwampSF): tenmod.TwampSF,
		string(tenmod.TwampSL): tenmod.TwampSL}

	// ValidMonitoredObjectDeviceTypes - known Monitored Object Device types in the system.
	ValidMonitoredObjectDeviceTypes = map[string]tenmod.MonitoredObjectDeviceType{
		string(tenmod.AccedianNID):  tenmod.AccedianNID,
		string(tenmod.AccedianVNID): tenmod.AccedianVNID}

	// DefaultValidTypes - default values for the valid types supported by datahub
	DefaultValidTypes = &admmod.ValidTypes{}
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

func areValidTypesEquivalent(obj1 *admmod.ValidTypes, obj2 *admmod.ValidTypes) bool {
	if (obj1 == nil && obj2 != nil) || (obj1 != nil && obj2 == nil) {
		return false
	}

	if obj1 == nil && obj2 == nil {
		return true
	}

	// Have 2 valid objects, do parameter comparison.
	// MonitoredObjectTypes
	if len(obj1.MonitoredObjectTypes) != len(obj2.MonitoredObjectTypes) {
		return false
	}
	for key, val := range obj1.MonitoredObjectTypes {
		if obj2.MonitoredObjectTypes[key] != val {
			return false
		}
	}

	// MonitoredObjectDeviceTypes
	if len(obj1.MonitoredObjectDeviceTypes) != len(obj2.MonitoredObjectDeviceTypes) {
		return false
	}
	for key, val := range obj1.MonitoredObjectDeviceTypes {
		if obj2.MonitoredObjectDeviceTypes[key] != val {
			return false
		}
	}

	return true
}

func provisionCouchData(pouchSH *handlers.PouchDBPluginServiceHandler, adminDB datastore.AdminServiceDatastore, adminDBStr string, cfg config.Provider) {
	ensureBaseCouchDBsExist(pouchSH)
	ensureAdminDBExists(pouchSH, adminDB, adminDBStr)
	ensureIngestionDictionaryExists(adminDB, adminDBStr, cfg)
	ensureValidTypesExists(adminDB, adminDBStr)
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

func ensureIngestionDictionaryExists(adminDB datastore.AdminServiceDatastore, adminDBStr string, cfg config.Provider) {
	ingDictFilePath := cfg.GetString("ingDict")
	defaultDictionaryBytes, err := ioutil.ReadFile(ingDictFilePath)
	if err != nil {
		logger.Log.Fatalf("Unable to read Default Ingestion Dictionary from file: %s", err.Error())
	}

	defaultDictionaryData := &admmod.IngestionDictionary{}
	if err = json.Unmarshal(defaultDictionaryBytes, &defaultDictionaryData); err != nil {
		logger.Log.Fatalf("Unable to construct Default Ingestion Dictionary from file: %s", err.Error())
	}

	existingDictionary, err := adminDB.GetIngestionDictionary()
	if err != nil {
		logger.Log.Debugf("Unable to fetch Ingestion Dictionary from DB %s: %s", adminDB, err.Error())

		// Provision the default IngestionDictionary
		_, err = adminDB.CreateIngestionDictionary(defaultDictionaryData)
		if err != nil {
			logger.Log.Fatalf("Unable to store Default Ingestion Profile from file: %s", err.Error())
		}

		return
	}

	// There is an existing dictionary, make sure it matches the known values.
	if !areIngestionDictionariesEqual(defaultDictionaryData, existingDictionary) {
		existingDictionary.Metrics = defaultDictionaryData.Metrics

		_, err = adminDB.UpdateIngestionDictionary(existingDictionary)
		if err != nil {
			logger.Log.Fatalf("Unable to update Default Ingestion Profile from file: %s", err.Error())
		}

		return
	}
}

func areIngestionDictionariesEqual(dict1 *admmod.IngestionDictionary, dict2 *admmod.IngestionDictionary) bool {
	if (dict1 == nil && dict2 != nil) || (dict1 != nil && dict2 == nil) {
		return false
	}

	if dict1 == nil && dict2 == nil {
		return true
	}

	// Have 2 valid objects, do parameter comparison.
	for vendor, metricMap := range dict1.Metrics {
		if dict2.Metrics[vendor] == nil {
			return false
		}

		for metric, metricDef := range metricMap.MetricMap {
			if dict2.Metrics[vendor].MetricMap[metric] == nil {
				return false
			}

			if !areUIPartsEqual(metricDef.UIData, dict2.Metrics[vendor].MetricMap[metric].UIData) {
				return false
			}

			for _, monitoredObjectType := range metricDef.MonitoredObjectTypes {
				if !doesSliceOfMonitoredObjectTypesContain(dict2.Metrics[vendor].MetricMap[metric].MonitoredObjectTypes, monitoredObjectType) {
					return false
				}
			}
		}
	}

	return true
}

func areUIPartsEqual(ui1 *admmod.UIData, ui2 *admmod.UIData) bool {
	if (ui1 == nil && ui2 != nil) || (ui1 != nil && ui2 == nil) {
		return false
	}

	if ui1 == nil && ui2 == nil {
		return true
	}

	if ui1.Group != ui2.Group {
		return false
	}
	if ui1.Position != ui2.Position {
		return false
	}

	return true
}

func areMonitoredObjectTypesEqual(mot1 *admmod.MonitoredObjectType, mot2 *admmod.MonitoredObjectType) bool {
	if (mot1 == nil && mot2 != nil) || (mot1 != nil && mot2 == nil) {
		return false
	}

	if mot1 == nil && mot2 == nil {
		return true
	}

	if mot1.Key != mot2.Key {
		return false
	}
	if mot1.RawMetricID != mot2.RawMetricID {
		return false
	}

	if !areStringSlicesEqual(mot1.Units, mot2.Units) {
		return false
	}

	if !areStringSlicesEqual(mot1.Directions, mot2.Directions) {
		return false
	}

	return true
}

func doesSliceOfMonitoredObjectTypesContain(container []*admmod.MonitoredObjectType, value *admmod.MonitoredObjectType) bool {
	for _, s := range container {
		if areMonitoredObjectTypesEqual(s, value) {
			return true
		}
	}
	return false
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

func ensureValidTypesExists(adminDB datastore.AdminServiceDatastore, adminDBStr string) {
	validMonObjTypes := make(map[string]string, 0)
	validMonObjDevTypes := make(map[string]string, 0)

	for key, val := range ValidMonitoredObjectTypes {
		validMonObjTypes[key] = string(val)
	}
	for key, val := range ValidMonitoredObjectDeviceTypes {
		validMonObjDevTypes[key] = string(val)
	}

	DefaultValidTypes = &admmod.ValidTypes{
		MonitoredObjectTypes:       validMonObjTypes,
		MonitoredObjectDeviceTypes: validMonObjDevTypes}

	// Make sure the valid types are provisioned.
	provisionedValidTypes, err := adminDB.GetValidTypes()
	if err != nil {
		logger.Log.Debugf("Unable to fetch Valid Values from DB %s: %s", adminDBStr, err.Error())

		// Provision the default values as a new object.
		provisionedValidTypes, err = adminDB.CreateValidTypes(DefaultValidTypes)
		if err != nil {
			logger.Log.Fatalf("Unable to Add Valid Values object to DB %s: %s", adminDBStr, err.Error())
		}
		return
	}
	if !areValidTypesEquivalent(provisionedValidTypes, DefaultValidTypes) {
		// Need to add the known default values to the data store
		provisionedValidTypes.MonitoredObjectTypes = DefaultValidTypes.MonitoredObjectTypes
		provisionedValidTypes.MonitoredObjectDeviceTypes = DefaultValidTypes.MonitoredObjectDeviceTypes
		provisionedValidTypes, err = adminDB.UpdateValidTypes(provisionedValidTypes)
		if err != nil {
			logger.Log.Fatalf("Unable to Update Valid Values object to DB %s: %s", adminDBStr, err.Error())
		}
	}
}
