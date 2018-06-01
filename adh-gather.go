package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
	"strings"
	"sync"
	"sync/atomic"

	"github.com/spf13/viper"

	"github.com/accedian/adh-gather/config"
	"github.com/accedian/adh-gather/gather"
	adhh "github.com/accedian/adh-gather/handlers"
	"github.com/accedian/adh-gather/logger"
	"github.com/accedian/adh-gather/monitoring"
	"github.com/accedian/adh-gather/profile"
	"github.com/accedian/adh-gather/websocket"
	gh "github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/spf13/pflag"
	"google.golang.org/grpc"

	pb "github.com/accedian/adh-gather/gathergrpc"
	admmod "github.com/accedian/adh-gather/models/admin"
	mon "github.com/accedian/adh-gather/monitoring"
	slasched "github.com/accedian/adh-gather/scheduler"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

const (
	defaultIngestionDictionaryPath = "files/defaultIngestionDictionary.json"
	defaultSwaggerFile             = "files/swagger.yml"
	globalChangesDBName            = "_global_changes"
	replicatorDBName               = "_replicator"
	metadataDBName                 = "_metadata"
	usersDBName                    = "_users"
	statsDBName                    = "_stats"
)

var (
	configFilePath  string
	enableTLS       bool
	tlsKeyFile      string
	tlsCertFile     string
	ingDictFilePath string
	swaggerFilePath string

	maxConcurrentMetricAPICalls uint64
	maxConcurrentProvAPICalls   uint64
	maxConcurrentPouchAPICalls  uint64

	concurrentMetricAPICounter uint64
	concurrentProvAPICounter   uint64
	concurrentPouchAPICounter  uint64

	metricAPIMutex = &sync.Mutex{}
	provAPIMutex   = &sync.Mutex{}
	pouchAPIMutex  = &sync.Mutex{}

	metricServiceEndpoints    []string
	enableChangeNotifications bool
	enableAuthorizationAAA    bool
)

func init() {
	pflag.StringVar(&configFilePath, "config", "config/adh-gather.yml", "Specify a configuration file to use")
	pflag.StringVar(&tlsKeyFile, "tlskey", "/run/secrets/tls_key", "Specify a TLS Key file")
	pflag.StringVar(&tlsCertFile, "tlscert", "/run/secrets/tls_crt", "Specify a TLS Cert file")
	pflag.BoolVar(&enableTLS, "tls", true, "Specify if TLS should be enabled")
	pflag.StringVar(&ingDictFilePath, "ingDict", defaultIngestionDictionaryPath, "Specify file path of default Ingestion Dictionary")
	pflag.StringVar(&swaggerFilePath, "swag", defaultSwaggerFile, "Specify file path of the Swagger documentation")

	pflag.BoolVar(&enableChangeNotifications, "changeNotifications", true, "Specify if Change Notifications should be enabled")

	pflag.BoolVar(&enableAuthorizationAAA, "enableAuthorizationAAA", true, "Specify if checking for Skylight AAA authorization is enabled")

	metricServiceEndpoints = []string{
		"/api/v1/histogram",
		"/api/v1/raw-metrics",
		"/api/v1/threshold-crossing-by-monitored-object",
		"/api/v1/threshold-crossing",
		"/api/v1/generate-sla-report",
		"/api/v1/threshold-crossing-by-monitored-object-top-n",
		"/api/v1/aggregated-metrics",
	}
}

// GatherServer - Server which will implement the gRPC Services.
type GatherServer struct {
	gsh         *adhh.GRPCServiceHandler
	pouchSH     *adhh.PouchDBPluginServiceHandler
	testSH      *adhh.TestDataServiceHandler
	msh         *adhh.MetricServiceHandler
	adminAPISH  *adhh.AdminServiceRESTHandler
	tenantAPISH *adhh.TenantServiceRESTHandler

	mux            *mux.Router
	jsonAPIMux     *mux.Router
	promServerMux  *http.ServeMux
	pprofServerMux *http.ServeMux
}

func newServer() *GatherServer {
	s := new(GatherServer)
	s.gsh = adhh.CreateCoordinator()
	s.pouchSH = adhh.CreatePouchDBPluginServiceHandler()
	s.testSH = adhh.CreateTestDataServiceHandler()

	s.msh = adhh.CreateMetricServiceHandler(s.gsh)
	s.adminAPISH = adhh.CreateAdminServiceRESTHandler()
	s.tenantAPISH = adhh.CreateTenantServiceRESTHandler()

	return s
}

func gRPCHandlerStart(gatherServer *GatherServer, cfg config.Provider) {
	gRPCAddress := fmt.Sprintf("%s:%d", cfg.GetString(gather.CK_server_grpc_ip.String()), cfg.GetInt(gather.CK_server_grpc_port.String()))

	lis, err := net.Listen("tcp", gRPCAddress)
	if err != nil {
		logger.Log.Fatalf("failed to start gRPC Service: %s", err.Error())
	}
	var opts []grpc.ServerOption

	grpcServer := grpc.NewServer(opts...)
	pb.RegisterAdminProvisioningServiceServer(grpcServer, gatherServer.gsh)
	pb.RegisterTenantProvisioningServiceServer(grpcServer, gatherServer.gsh)

	logger.Log.Infof("gRPC service intiated on: %s", gRPCAddress)
	grpcServer.Serve(lis)
}

func restHandlerStart(gatherServer *GatherServer, cfg config.Provider) {
	restBindIP := cfg.GetString(gather.CK_server_rest_ip.String())
	restBindPort := cfg.GetInt(gather.CK_server_rest_port.String())

	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	gatherServer.mux = mux.NewRouter().StrictSlash(true)
	gatherServer.jsonAPIMux = mux.NewRouter().StrictSlash(true)

	// Register all the API endpoints:
	gatherServer.pouchSH.RegisterAPIHandlers(gatherServer.mux)
	gatherServer.testSH.RegisterAPIHandlers(gatherServer.mux)
	gatherServer.msh.RegisterAPIHandlers(gatherServer.jsonAPIMux)
	gatherServer.adminAPISH.RegisterAPIHandlers(gatherServer.mux)
	gatherServer.tenantAPISH.RegisterAPIHandlers(gatherServer.mux)

	allowedOrigins := cfg.GetStringSlice(gather.CK_server_cors_allowedorigins.String())
	logger.Log.Debugf("Allowed Origins: %v", allowedOrigins)
	originsOption := gh.AllowedOrigins(allowedOrigins)
	methodsOption := gh.AllowedMethods([]string{"GET", "HEAD", "POST", "PUT", "PATCH", "OPTIONS", "DELETE"})
	headersOption := gh.AllowedHeaders([]string{"accept", "authorization", "content-type", "origin", "referer", "x-csrf-token"})
	logger.Log.Infof("REST service intiated on: %s:%d", restBindIP, restBindPort)

	// Enable TLS based on config
	handler := gh.CORS(originsOption, methodsOption, headersOption, gh.AllowCredentials())(gatherServer)
	addr := fmt.Sprintf("%s:%d", restBindIP, restBindPort)
	if enableTLS {
		if _, err := os.Stat(tlsCertFile); os.IsNotExist(err) {
			// No TLS cert file
			logger.Log.Fatalf("Failed to start Gather: TLS cert %s does not exist", tlsCertFile)
		}
		if _, err := os.Stat(tlsKeyFile); os.IsNotExist(err) {
			// No TLS cert file
			logger.Log.Fatalf("Failed to start Gather: TLS key %s does not exist", tlsKeyFile)
		}
		http.ListenAndServeTLS(addr, tlsCertFile, tlsKeyFile, handler)
	} else {
		http.ListenAndServe(addr, handler)
	}

}

// Handle requests based on the path provided. If it begins with the known
// gRPC REST GW handler prefix, then use that handler, use the default handler
// otherwise.
func (gs *GatherServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	logger.Log.Debugf("Received API call: %s", r.URL.Path)

	isPouch := strings.Index(r.URL.Path, "/pouchdb") == 0

	// Exclude Pouch API calls from the total count as it has expected failures.
	if !isPouch {
		mon.RecievedAPICalls.Inc()
	}

	if strings.Index(r.URL.Path, "/api/v1/") == 0 {
		if isMetricAPICall(r.URL.Path) {
			// Handle Metrics Calls
			if err := updateCounter(&concurrentMetricAPICounter, metricAPIMutex, true, maxConcurrentMetricAPICalls); err != nil {
				reportOverloaded(w, r, err.Error())
				mon.CompletedAPICalls.Inc()
				return
			}

			mon.IncrementCounter(mon.MetricAPIRecieved)

			gs.jsonAPIMux.ServeHTTP(w, r)

			updateCounter(&concurrentMetricAPICounter, metricAPIMutex, false, maxConcurrentMetricAPICalls)
			mon.IncrementCounter(mon.MetricAPICompleted)
			return
		}

		// Handle calls to our Admin and Tenant Services
		if err := updateCounter(&concurrentProvAPICounter, provAPIMutex, true, maxConcurrentProvAPICalls); err != nil {
			reportOverloaded(w, r, err.Error())
			mon.CompletedAPICalls.Inc()
			return
		}

		isTenant := strings.Index(r.URL.Path, "/api/v1/tenants") == 0
		if isTenant {
			mon.IncrementCounter(mon.TenantAPIRecieved)
		} else {
			mon.IncrementCounter(mon.AdminAPIRecieved)
		}

		gs.mux.ServeHTTP(w, r)

		updateCounter(&concurrentProvAPICounter, provAPIMutex, false, maxConcurrentProvAPICalls)
		if isTenant {
			mon.IncrementCounter(mon.TenantAPICompleted)
		} else {
			mon.IncrementCounter(mon.AdminAPICompleted)
		}
	} else if strings.Index(r.URL.Path, "/swagger.json") == 0 {
		// Handle requests for the swagger definition
		input, err := ioutil.ReadFile(swaggerFilePath)
		if err != nil {
			logger.Log.Fatalf("Unable to locate swagger definition: %s", err.Error())
		}
		w.Write(input)
	} else if isPouch && strings.Index(r.URL.Path, "_changes") > 0 {
		// This is a pouch _changes call which sometimes is cancelled by the client.
		// We do not want the same pouchDB load shedder used here as it will get hit for
		// incorrect reasons.
		mon.IncrementCounter(mon.PouchChangesAPIRecieved)
		gs.mux.ServeHTTP(w, r)
		mon.IncrementCounter(mon.PouchChangesAPICompleted)
	} else {
		// Handle all other endpoints (really just Pouch and Test Data right now.
		if err := updateCounter(&concurrentPouchAPICounter, pouchAPIMutex, true, maxConcurrentPouchAPICalls); err != nil {
			reportOverloaded(w, r, err.Error())
			mon.CompletedAPICalls.Inc()
			return
		}

		if isPouch {
			mon.IncrementCounter(mon.PouchAPIRecieved)
		}

		gs.mux.ServeHTTP(w, r)

		if isPouch {
			mon.IncrementCounter(mon.PouchAPICompleted)
			updateCounter(&concurrentPouchAPICounter, pouchAPIMutex, false, maxConcurrentPouchAPICalls)
		}

	}

	if !isPouch {
		mon.CompletedAPICalls.Inc()
	}
}

func isMetricAPICall(url string) bool {
	for _, val := range metricServiceEndpoints {
		if strings.Index(url, val) == 0 {
			return true
		}
	}

	return false
}

func reportOverloaded(w http.ResponseWriter, r *http.Request, errorStr string) {
	msg := fmt.Sprintf("Unable to complete %s API: %s", r.URL.Path, errorStr)
	logger.Log.Infof(msg)
	http.Error(w, msg, http.StatusServiceUnavailable)
}

// updateCounter - updates a counter only if the counter is within the valid range between 0 and maxOperations.
// returns an error if the counter could not be incremented due to reaching the max.
func updateCounter(counter *uint64, mutex *sync.Mutex, increment bool, maxOperations uint64) error {
	// increment, but ensure it stays below max.
	if increment {
		mutex.Lock()
		currentVal := atomic.LoadUint64(counter)
		if currentVal >= maxOperations {
			mutex.Unlock()
			return fmt.Errorf("Server has reached the maximum allowed operation of this type")
		}
		atomic.AddUint64(counter, 1)
		mutex.Unlock()
		return nil
	}

	// Decrement but keep at or above 0
	mutex.Lock()
	atomic.AddUint64(counter, ^uint64(0))
	newVal := atomic.LoadUint64(counter)
	if newVal > ^uint64(0)-1000 {
		atomic.StoreUint64(counter, 0)
	}
	mutex.Unlock()
	return nil
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

func doesSliceContainString(container []string, value string) bool {
	for _, s := range container {
		if s == value {
			return true
		}
	}
	return false
}

func provisionCouchData(gatherServer *GatherServer, adminDB string) {
	ensureBaseCouchDBsExist(gatherServer)
	ensureAdminDBExists(gatherServer, adminDB)
	ensureIngestionDictionaryExists(gatherServer, adminDB)
	ensureValidTypesExists(gatherServer, adminDB)
}

func createCouchDB(gatherServer *GatherServer, dbName string) error {
	// Make sure global changes db exists
	_, err := gatherServer.pouchSH.IsDBAvailable(dbName)
	if err != nil {
		logger.Log.Infof("Database %s does not exist. %s DB will now be created.", dbName, dbName)
		// Try to create the DB
		_, err = gatherServer.pouchSH.AddDB(dbName)
		if err != nil {
			return err
		}
	}

	return nil
}

func ensureBaseCouchDBsExist(gatherServer *GatherServer) {
	if err := createCouchDB(gatherServer, globalChangesDBName); err != nil {
		logger.Log.Fatalf("Unable to create DB %s: %s", globalChangesDBName, err.Error())
	}
	if err := createCouchDB(gatherServer, metadataDBName); err != nil {
		logger.Log.Fatalf("Unable to create DB %s: %s", metadataDBName, err.Error())
	}
	if err := createCouchDB(gatherServer, replicatorDBName); err != nil {
		logger.Log.Fatalf("Unable to create DB %s: %s", replicatorDBName, err.Error())
	}
	if err := createCouchDB(gatherServer, usersDBName); err != nil {
		logger.Log.Fatalf("Unable to create DB %s: %s", usersDBName, err.Error())
	}
}

func ensureAdminDBExists(gatherServer *GatherServer, adminDB string) {
	_, err := gatherServer.pouchSH.IsDBAvailable(adminDB)
	if err != nil {
		logger.Log.Infof("Database %s does not exist. %s DB will now be created.", adminDB, adminDB)

		// Try to create the DB:
		_, err = gatherServer.pouchSH.AddDB(adminDB)
		if err != nil {
			logger.Log.Fatalf("Unable to create DB %s: %s", adminDB, err.Error())
		}

		// Also add the Views for Admin DB.
		err = gatherServer.adminAPISH.AddAdminViews()
		if err != nil {
			logger.Log.Fatalf("Unable to Add Views to DB %s: %s", adminDB, err.Error())
		}
	}

	logger.Log.Infof("Using %s as Administrative Database", adminDB)
}

func ensureIngestionDictionaryExists(gatherServer *GatherServer, adminDB string) {
	defaultDictionaryBytes, err := ioutil.ReadFile(ingDictFilePath)
	if err != nil {
		logger.Log.Fatalf("Unable to read Default Ingestion Dictionary from file: %s", err.Error())
	}

	defaultDictionaryData := &admmod.IngestionDictionary{}
	if err = json.Unmarshal(defaultDictionaryBytes, &defaultDictionaryData); err != nil {
		logger.Log.Fatalf("Unable to construct Default Ingestion Dictionary from file: %s", err.Error())
	}

	existingDictionary, err := gatherServer.adminAPISH.GetIngestionDictionaryInternal()
	if err != nil {
		logger.Log.Debugf("Unable to fetch Ingestion Dictionary from DB %s: %s", adminDB, err.Error())

		// Provision the default IngestionDictionary
		_, err = gatherServer.adminAPISH.CreateIngestionDictionaryInternal(defaultDictionaryData)
		if err != nil {
			logger.Log.Fatalf("Unable to store Default Ingestion Profile from file: %s", err.Error())
		}

		return
	}

	// There is an existing dictionary, make sure it matches the known values.
	if !areIngestionDictionariesEqual(defaultDictionaryData, existingDictionary) {
		existingDictionary.Metrics = defaultDictionaryData.Metrics

		_, err = gatherServer.adminAPISH.UpdateIngestionDictionaryInternal(existingDictionary)
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
		if !doesSliceContainString(slice2, value) {
			return false
		}
	}

	return true
}

func ensureValidTypesExists(gatherServer *GatherServer, adminDB string) {
	// Make sure the valid types are provisioned.
	provisionedValidTypes, err := gatherServer.adminAPISH.GetValidTypesInternal()
	if err != nil {
		logger.Log.Debugf("Unable to fetch Valid Values from DB %s: %s", adminDB, err.Error())

		// Provision the default values as a new object.
		provisionedValidTypes, err = gatherServer.adminAPISH.CreateValidTypesInternal(gatherServer.gsh.DefaultValidTypes)
		if err != nil {
			logger.Log.Fatalf("Unable to Add Valid Values object to DB %s: %s", adminDB, err.Error())
		}
		return
	}
	if !areValidTypesEquivalent(provisionedValidTypes, gatherServer.gsh.DefaultValidTypes) {
		// Need to add the known default values to the data store
		provisionedValidTypes.MonitoredObjectTypes = gatherServer.gsh.DefaultValidTypes.MonitoredObjectTypes
		provisionedValidTypes.MonitoredObjectDeviceTypes = gatherServer.gsh.DefaultValidTypes.MonitoredObjectDeviceTypes
		provisionedValidTypes, err = gatherServer.adminAPISH.UpdateValidTypesInternal(provisionedValidTypes)
		if err != nil {
			logger.Log.Fatalf("Unable to Update Valid Values object to DB %s: %s", adminDB, err.Error())
		}
	}
}

func startMonitoring(gatherServer *GatherServer, cfg config.Provider) {
	restBindIP := cfg.GetString(gather.CK_server_rest_ip.String())
	monPort := cfg.GetInt(gather.CK_server_monitoring_port.String())

	monitoring.InitMetrics()
	gatherServer.promServerMux = http.NewServeMux()

	logger.Log.Infof("Starting Prometheus Server")
	gatherServer.promServerMux.Handle("/metrics", promhttp.Handler())
	addr := fmt.Sprintf("%s:%d", restBindIP, monPort)
	if err := http.ListenAndServe(addr, gatherServer.promServerMux); err != nil {
		logger.Log.Fatalf("Unable to start monitoring function: %s", err.Error())
	}
}

func startProfile(gatherServer *GatherServer, cfg config.Provider) {
	restBindIP := cfg.GetString(gather.CK_server_rest_ip.String())
	monPort := cfg.GetInt(gather.CK_server_profile_port.String())

	gatherServer.pprofServerMux = http.NewServeMux()

	profile.AttachProfiler(gatherServer.pprofServerMux)

	logger.Log.Infof("Starting Profile Server")
	addr := fmt.Sprintf("%s:%d", restBindIP, monPort)
	if err := http.ListenAndServe(addr, gatherServer.pprofServerMux); err != nil {
		logger.Log.Fatalf("Unable to start profile function: %s", err.Error())
	}
}

func startChangeNotificationHandler() {
	// Start monitoring changes and sending notifications
	hdlr := adhh.CreateChangeNotificationHandler()
	hdlr.SendChangeNotifications()
}

func modifySwagger(cfg config.Provider) {
	apiPort := cfg.GetInt(gather.CK_server_rest_port.String())

	var hostLine string
	var schemeLine string
	hostFromEnv := os.Getenv("API_TARGET")
	if len(hostFromEnv) == 0 {
		hostLine = fmt.Sprintf(`host: 'localhost:%d'`, apiPort)
		schemeLine = "  - http"
	} else {
		hostLine = fmt.Sprintf(`host: '%s'`, hostFromEnv)
		schemeLine = "  - https"
	}

	// Update the generated swagger file to contain the correct host
	input, err := ioutil.ReadFile(swaggerFilePath)
	if err != nil {
		logger.Log.Fatalf("Unable to locate swagger definition: %s", err.Error())
	}

	// Replace the host line
	containsHost := false
	lines := strings.Split(string(input), "\n")
	for i, line := range lines {
		if strings.Contains(line, `host:`) {
			containsHost = true
			lines[i] = hostLine
			break
		}
	}
	if !containsHost {
		// Insert the host into the swager file
		lines = append(lines[:2], append([]string{hostLine}, lines[2:]...)...)
	}

	// Append the appropriate scheme line:
	var index int
	for i, line := range lines {
		if strings.Contains(line, `schemes:`) {
			index = i + 1
			break
		}
	}
	lines[index] = schemeLine

	output := strings.Join(lines, "\n")
	err = ioutil.WriteFile(swaggerFilePath, []byte(output), 0644)
	if err != nil {
		log.Fatalln(err)
	}
}

func main() {
	pflag.Parse()
	v := viper.New()

	v.BindPFlags(pflag.CommandLine)

	configFilePath = v.GetString("config")
	enableTLS = v.GetBool("tls")
	tlsCertFile = v.GetString("tlscert")
	tlsKeyFile = v.GetString("tlskey")
	ingDictFilePath = v.GetString("ingDict")
	swaggerFilePath = v.GetString("swag")
	enableChangeNotifications = v.GetBool("changeNotifications")

	// Load Configuration
	cfg := gather.LoadConfig(configFilePath, v)
	enableAuthorizationAAA = v.GetBool("enableAuthorizationAAA")
	cfg.Set(gather.CK_args_authorizationAAA.String(), enableAuthorizationAAA)

	debug := cfg.GetBool(gather.CK_args_debug.String())
	if debug {
		logger.SetDebugLevel(true)
	} else {
		logger.SetDebugLevel(false)
	}

	logger.Log.Infof("Starting adh-gather broker with config '%s'", configFilePath)

	maxConcurrentMetricAPICalls = uint64(v.GetInt64(gather.CK_args_maxConcurrentMetricAPICalls.String()))
	maxConcurrentPouchAPICalls = uint64(v.GetInt64(gather.CK_args_maxConcurrentPouchAPICalls.String()))
	maxConcurrentProvAPICalls = uint64(v.GetInt64(gather.CK_args_maxConcurrentProvAPICalls.String()))
	logger.Log.Debugf("API caps are set to: Metric-%d Prov-%d Pouch-%d", maxConcurrentMetricAPICalls, maxConcurrentProvAPICalls, maxConcurrentPouchAPICalls)

	if enableChangeNotifications {
		// Start monitoring changes and sending notifications
		go startChangeNotificationHandler()
	}

	// Start the REST and gRPC Services
	gatherServer := newServer()

	// Register the metrics to be tracked in Gather
	go startMonitoring(gatherServer, cfg)

	// Start pprof profiler
	go startProfile(gatherServer, cfg)

	// Start websocket server
	websocket.Server(gatherServer.tenantAPISH.TenantDB)

	// modify the swagger for this deployment
	modifySwagger(cfg)

	adminDB := cfg.GetString(gather.CK_args_admindb_name.String())
	provisionCouchData(gatherServer, adminDB)

	slasched.Initialize(gatherServer.msh, nil, nil, 5)

	go restHandlerStart(gatherServer, cfg)
	gRPCHandlerStart(gatherServer, cfg)

}
