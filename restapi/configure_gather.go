// This file is safe to edit. Once it exists it will not be overwritten

package restapi

import (
	"crypto/tls"
	"net/http"
	"strings"

	errors "github.com/go-openapi/errors"
	runtime "github.com/go-openapi/runtime"

	"github.com/accedian/adh-gather/datastore"
	"github.com/accedian/adh-gather/datastore/druid"
	"github.com/accedian/adh-gather/gather"
	"github.com/accedian/adh-gather/handlers"
	"github.com/accedian/adh-gather/logger"
	"github.com/accedian/adh-gather/restapi/operations"
	slasched "github.com/accedian/adh-gather/scheduler"
	"github.com/accedian/adh-gather/websocket"
	mux "github.com/gorilla/mux"
)

//go:generate swagger generate server --target .. --name gather --spec ../files/swagger.yml --model-package swagmodels --exclude-main --exclude-spec

const (
	testDataAPIPrefix = "/test-data"
)

var (
	metricSH      *handlers.MetricServiceHandler
	testSH        *handlers.TestDataServiceHandler
	nonSwaggerMUX *mux.Router
	adminDB       datastore.AdminServiceDatastore
	tenantDB      datastore.TenantServiceDatastore
	druidDB       datastore.DruidDatastore

	metricServiceV1APIRouteRoots = []string{
		"/api/v1/threshold-crossing", "/api/v1/threshold-crossing-by-monitored-object", "/api/v1/threshold-crossing-by-monitored-object-top-n",
		"/api/v1/generate-sla-report", "/api/v1/histogram", "/api/v1/raw-metrics", "/api/v2/raw-metrics", "/api/v1/aggregated-metrics", "/api/v1/topn-metrics",
	}

	supportedOrigins = []string{}
	supportedMethods = "GET, HEAD, POST, PUT, PATCH, OPTIONS"
)

func configureFlags(api *operations.GatherAPI) {
	// api.CommandLineOptionsGroups = []swag.CommandLineOptionsGroup{ ... }
}

func configureAPI(api *operations.GatherAPI) http.Handler {
	// configure the api here
	api.ServeError = errors.ServeError

	// Set your custom logger if needed. Default one is log.Printf
	// Expected interface func(string, ...interface{})
	//
	// Example:
	// api.Logger = log.Printf

	api.JSONConsumer = runtime.JSONConsumer()

	api.JSONProducer = runtime.JSONProducer()

	api.TxtProducer = runtime.TextProducer()

	handlers.InitializeAuthHelper()

	// Create DAO objects to handle data retrieval as needed.
	var err error
	adminDB, err = handlers.GetAdminServiceDatastore()
	if err != nil {
		logger.Log.Fatalf("Unable to instantiate Admin Service DAO: %s", err.Error())
	}

	tenantDB, err = handlers.GetTenantServiceDatastore()
	if err != nil {
		logger.Log.Fatalf("Unable to instantiate Tenant Service DAO: %s", err.Error())
	}

	druidDB := druid.NewDruidDatasctoreClient()

	// Register the V1 APIs
	configureAdminServiceV1API(api, adminDB, tenantDB)
	configureTenantServiceV1API(api, tenantDB)
	// configureMetricServiceV1API(api, druidDB)  - Not using this implementation but leaving it in in case we change our minds
	configurev1APIThatWeMayRemove(api, tenantDB)

	// Register the V2 APIs
	configureAdminServiceV2API(api, adminDB, tenantDB)
	configureTenantServiceV2API(api, tenantDB, druidDB)
	// configureMetricServiceV2API(api, druidDB)

	api.ServerShutdown = func() {}

	// Setup handling of non-swagger generated APIs and other boot processes
	return setupGlobalMiddleware(api.Serve(setupMiddlewares))
}

// The TLS configuration before HTTPS server starts.
func configureTLS(tlsConfig *tls.Config) {
	// Make all necessary changes to the TLS configuration here.
}

// As soon as server is initialized but not run yet, this function will be called.
// If you need to modify a config, store server instance to stop it individually later, this is the place.
// This function can be called multiple times, depending on the number of serving schemes.
// scheme value will be set accordingly: "http", "https" or "unix"
func configureServer(s *http.Server, scheme, addr string) {
}

// The middleware configuration is for the handler executors. These do not apply to the swagger.json document.
// The middleware executes after routing but before authentication, binding and validation
func setupMiddlewares(handler http.Handler) http.Handler {
	return handler
}

// The middleware configuration happens before anything, this middleware also applies to serving the swagger.json document.
// So this is a good place to plug in a panic handling middleware, logging and metrics
func setupGlobalMiddleware(handler http.Handler) http.Handler {
	cfg := gather.GetConfig()

	// Setup non-swagger MUX
	metricSH = handlers.CreateMetricServiceHandler()
	testSH = handlers.CreateTestDataServiceHandler()
	nonSwaggerMUX = mux.NewRouter().StrictSlash(true)
	testSH.RegisterAPIHandlers(nonSwaggerMUX)
	metricSH.RegisterAPIHandlers(nonSwaggerMUX)

	supportedOrigins = cfg.GetStringSlice(gather.CK_server_cors_allowedorigins.String())

	// Register the metrics to be tracked in Gather
	go startMonitoring(cfg)

	// Start pprof profiler
	go startProfile(cfg)

	// Start websocket server
	websocket.Server(tenantDB)

	// Start the scheduler for handling SLA report generation
	slasched.Initialize(metricSH, nil, nil, 5)

	// Make sure necessary Couch data is present
	pouchSH := handlers.CreatePouchDBPluginServiceHandler()
	adminDBStr := cfg.GetString(gather.CK_args_admindb_name.String())
	provisionCouchData(pouchSH, adminDB, adminDBStr, cfg)

	return addNonSwaggerHandler(handler)
}

func addNonSwaggerHandler(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger.Log.Debugf("Received request: %s%s", r.Method, r.URL)

		// Handle CORS headers:
		reqHeader := r.Header.Get("Origin")
		logger.Log.Debugf("Recieved origin is: %s", reqHeader)
		if isOriginInSupportedList(reqHeader) {

			w.Header().Set("Access-Control-Allow-Origin", reqHeader)
			w.Header().Set("Access-Control-Allow-Methods", supportedMethods)
		}

		if r.Method == http.MethodOptions {
			optionsHandler(w, r)
			return
		}

		if strings.Index(r.URL.Path, testDataAPIPrefix) == 0 {
			// Test Data Call
			nonSwaggerMUX.ServeHTTP(w, r)
		} else if gather.DoesSliceContainString(metricServiceV1APIRouteRoots, r.URL.Path) {
			// Metric Service V1 call
			nonSwaggerMUX.ServeHTTP(w, r)
		} else {
			next.ServeHTTP(w, r)
		}

	})
}

func isOriginInSupportedList(requestOrigin string) bool {
	for _, originRegex := range supportedOrigins {
		if strings.HasSuffix(requestOrigin, originRegex) {
			return true
		}
	}

	return false
}

func optionsHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNoContent)
}
