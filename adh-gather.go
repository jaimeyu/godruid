package main

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"os"
	"strings"
	"io/ioutil"
	"encoding/json"

	"github.com/golang/protobuf/proto"
	"github.com/spf13/viper"

	"github.com/accedian/adh-gather/config"
	"github.com/accedian/adh-gather/gather"
	adhh "github.com/accedian/adh-gather/handlers"
	"github.com/accedian/adh-gather/logger"
	gh "github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/spf13/pflag"
	"google.golang.org/grpc"

	pb "github.com/accedian/adh-gather/gathergrpc"
	emp "github.com/golang/protobuf/ptypes/empty"
)

const (
	defaultIngestionDictionaryPath = "files/defaultIngestionDictionary.json"
)

var (
	configFilePath string
	enableTLS      bool
	tlsKeyFile     string
	tlsCertFile    string
	ingDictFilePath string
)

func init() {
	pflag.StringVar(&configFilePath, "config", "config/adh-gather.yml", "Specify a configuration file to use")
	pflag.StringVar(&tlsKeyFile, "tlskey", "/run/secrets/tls_key", "Specify a TLS Key file")
	pflag.StringVar(&tlsCertFile, "tlscert", "/run/secrets/tls_crt", "Specify a TLS Cert file")
	pflag.BoolVar(&enableTLS, "tls", true, "Specify if TLS should be enabled")
	pflag.StringVar(&ingDictFilePath, "ingDict", defaultIngestionDictionaryPath, "Specify file path of default Ingestion Dictionary")
}

// GatherServer - Server which will implement the gRPC Services.
type GatherServer struct {
	gsh     *adhh.GRPCServiceHandler
	pouchSH *adhh.PouchDBPluginServiceHandler
	testSH  *adhh.TestDataServiceHandler

	mux        *mux.Router
	gwmux      *runtime.ServeMux
	jsonAPIMux *runtime.ServeMux
}

func newServer() *GatherServer {
	s := new(GatherServer)
	s.gsh = adhh.CreateCoordinator()
	s.pouchSH = adhh.CreatePouchDBPluginServiceHandler()
	s.testSH = adhh.CreateTestDataServiceHandler()

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
	pb.RegisterMetricsServiceServer(grpcServer, gatherServer.gsh)

	logger.Log.Infof("gRPC service intiated on: %s", gRPCAddress)
	grpcServer.Serve(lis)
}

func restHandlerStart(gatherServer *GatherServer, cfg config.Provider) {
	restBindIP := cfg.GetString(gather.CK_server_rest_ip.String())
	restBindPort := cfg.GetInt(gather.CK_server_rest_port.String())
	grpcBindIP := cfg.GetString(gather.CK_server_grpc_ip.String())
	grpcBindPort := cfg.GetInt(gather.CK_server_grpc_port.String())

	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	gatherServer.mux = mux.NewRouter().StrictSlash(true)

	gatherServer.gwmux = runtime.NewServeMux()

	gatherServer.jsonAPIMux = runtime.NewServeMux(
		runtime.WithForwardResponseOption(
			func(ctx context.Context, w http.ResponseWriter, _ proto.Message) error {
				w.Header().Set("Content-Type", "application/vnd.api+json")
				return nil
			},
		),
	)

	opts := []grpc.DialOption{grpc.WithInsecure()}

	// Register the Admin Service
	if err := pb.RegisterAdminProvisioningServiceHandlerFromEndpoint(ctx, gatherServer.gwmux, fmt.Sprintf("%s:%d", grpcBindIP, grpcBindPort), opts); err != nil {
		logger.Log.Fatalf("failed to start REST service: %s", err.Error())
	}

	// Register the Tenant Service
	if err := pb.RegisterTenantProvisioningServiceHandlerFromEndpoint(ctx, gatherServer.gwmux, fmt.Sprintf("%s:%d", grpcBindIP, grpcBindPort), opts); err != nil {
		logger.Log.Fatalf("failed to start REST service: %s", err.Error())
	}

	// Register the Metrics Service
	if err := pb.RegisterMetricsServiceHandlerFromEndpoint(ctx, gatherServer.jsonAPIMux, fmt.Sprintf("%s:%d", grpcBindIP, grpcBindPort), opts); err != nil {
		logger.Log.Fatalf("failed to start REST service: %s", err.Error())
	}

	// Add in handling for non protobuf generated API endpoints:
	gatherServer.pouchSH.RegisterAPIHandlers(gatherServer.mux)
	gatherServer.testSH.RegisterAPIHandlers(gatherServer.mux)

	allowedOrigins := cfg.GetStringSlice(gather.CK_server_cors_allowedorigins.String())
	logger.Log.Debugf("Allowed Origins: %v", allowedOrigins)
	originsOption := gh.AllowedOrigins(allowedOrigins)
	methodsOption := gh.AllowedMethods([]string{"GET", "HEAD", "POST", "PUT", "OPTIONS", "DELETE"})
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

	if strings.Compare("application/vnd.api+json", r.Header.Get("Content-Type")) == 0 {
		gs.jsonAPIMux.ServeHTTP(w, r)
	} else if strings.Index(r.URL.Path, "/api/v1/") == 0 {
		gs.gwmux.ServeHTTP(w, r)
	} else {
		gs.mux.ServeHTTP(w, r)
	}
}

func areValidTypesEquivalent(obj1 *pb.ValidTypesData, obj2 *pb.ValidTypesData) bool {
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
	ensureAdminDBExists(gatherServer, adminDB)
	ensureIngestionDictionaryExists(gatherServer, adminDB)
	ensureValidTypesExists(gatherServer, adminDB)
}

func ensureAdminDBExists(gatherServer *GatherServer, adminDB string) {
	// Make sure the admin DB exists:
	_, err := gatherServer.pouchSH.IsDBAvailable(adminDB)
	if err != nil {
		logger.Log.Infof("Database %s does not exist. %s DB will now be created.", adminDB, adminDB)

		// Try to create the DB:
		_, err = gatherServer.pouchSH.AddDB(adminDB)
		if err != nil {
			logger.Log.Fatalf("Unable to create DB %s: %s", adminDB, err.Error())
		}

		// Also add the Views for Admin DB.
		err = gatherServer.gsh.AddAdminViews()
		if err != nil {
			logger.Log.Fatalf("Unable to Add Views to DB %s: %s", adminDB, err.Error())
		}
	}

	logger.Log.Infof("Using %s as Administrative Database", adminDB)
}

func ensureIngestionDictionaryExists(gatherServer *GatherServer, adminDB string) {
	defaultDictionaryBytes, err := ioutil.ReadFile(ingDictFilePath)
	if err != nil {
		logger.Log.Fatalf("Unable to read Default Ingestion Profile from file: %s", err.Error())
	}
	
	defaultDictionaryData := &pb.IngestionDictionaryData{}
	if err = json.Unmarshal(defaultDictionaryBytes, &defaultDictionaryData); err != nil {
		logger.Log.Fatalf("Unable to construct Default Ingestion Profile from file: %s", err.Error())
	}

	existingDictionary, err := gatherServer.gsh.GetIngestionDictionary(nil, &emp.Empty{})
	if err != nil {
		logger.Log.Debugf("Unable to fetch Ingestion Dictionary from DB %s: %s", adminDB, err.Error())

		// Provision the default IngestionDictionary
		_, err = gatherServer.gsh.CreateIngestionDictionary(nil, &pb.IngestionDictionary{Data: defaultDictionaryData})
		if err != nil {
			logger.Log.Fatalf("Unable to store Default Ingestion Profile from file: %s", err.Error())
		}

		return
	}

	// There is an existing dictionary, make sure it matches the known values.
	if !areIngestionDictionariesEqual(defaultDictionaryData, existingDictionary.Data) {
		existingDictionary.Data.Metrics = defaultDictionaryData.Metrics

		_, err = gatherServer.gsh.UpdateIngestionDictionary(nil, existingDictionary)
		if err != nil {
			logger.Log.Fatalf("Unable to update Default Ingestion Profile from file: %s", err.Error())
		}

		return
	}
}

func areIngestionDictionariesEqual(dict1 *pb.IngestionDictionaryData, dict2 *pb.IngestionDictionaryData) bool {
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

			if !areUIPartsEqual(metricDef.Ui, dict2.Metrics[vendor].MetricMap[metric].Ui) {
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

func areUIPartsEqual(ui1 *pb.IngestionDictionaryData_UIData, ui2 *pb.IngestionDictionaryData_UIData) bool {
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

func areMonitoredObjectTypesEqual(mot1 *pb.IngestionDictionaryData_MonitoredObjectType, mot2 *pb.IngestionDictionaryData_MonitoredObjectType) bool {
	if (mot1 == nil && mot2 != nil) || (mot1 != nil && mot2 == nil) {
		return false
	}

	if mot1 == nil && mot2 == nil {
		return true
	}

	if mot1.Key != mot2.Key {
		return false
	}
	if mot1.RawMetricId != mot2.RawMetricId {
		return false
	}

	if !areStringSlicesEqual(mot1.Units, mot2.Units){
		return false
	}

	if !areStringSlicesEqual(mot1.Directions, mot2.Directions) {
		return false
	}

	return true
}

func doesSliceOfMonitoredObjectTypesContain(container []*pb.IngestionDictionaryData_MonitoredObjectType, value *pb.IngestionDictionaryData_MonitoredObjectType) bool {
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
	provisionedValidTypes, err := gatherServer.gsh.GetValidTypes(nil, &emp.Empty{})
	if err != nil {
		logger.Log.Debugf("Unable to fetch Valid Values from DB %s: %s", adminDB, err.Error())

		// Provision the default values as a new object.
		provisionedValidTypes, err = gatherServer.gsh.CreateValidTypes(nil, &pb.ValidTypes{Data: gatherServer.gsh.DefaultValidTypes})
		if err != nil {
			logger.Log.Fatalf("Unable to Add Valid Values object to DB %s: %s", adminDB, err.Error())
		}
		return
	}
	if !areValidTypesEquivalent(provisionedValidTypes.Data, gatherServer.gsh.DefaultValidTypes) {
		// Need to add the known default values to the data store
		provisionedValidTypes.Data.MonitoredObjectTypes = gatherServer.gsh.DefaultValidTypes.MonitoredObjectTypes
		provisionedValidTypes.Data.MonitoredObjectDeviceTypes = gatherServer.gsh.DefaultValidTypes.MonitoredObjectDeviceTypes
		provisionedValidTypes, err = gatherServer.gsh.UpdateValidTypes(nil, provisionedValidTypes)
		if err != nil {
			logger.Log.Fatalf("Unable to Update Valid Values object to DB %s: %s", adminDB, err.Error())
		}
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

	// Load Configuration
	cfg := gather.LoadConfig(configFilePath, v)

	debug := cfg.GetBool(gather.CK_args_debug.String())
	if debug {
		logger.SetDebugLevel(true)
	} else {
		logger.SetDebugLevel(false)
	}

	logger.Log.Infof("Starting adh-gather broker with config '%s'", configFilePath)

	// Start the REST and gRPC Services
	gatherServer := newServer()

	adminDB := cfg.GetString(gather.CK_args_admindb_name.String())
	provisionCouchData(gatherServer, adminDB)

	go restHandlerStart(gatherServer, cfg)
	gRPCHandlerStart(gatherServer, cfg)

}
