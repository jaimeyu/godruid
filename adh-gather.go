package main

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"os"
	"strings"

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
)

var (
	configFilePath string
	enableTLS      bool
	tlsKeyFile     string
	tlsCertFile    string
)

func init() {
	pflag.StringVar(&configFilePath, "config", "config/adh-gather.yml", "Specify a configuration file to use")
	pflag.StringVar(&tlsKeyFile, "tlskey", "/run/secrets/tls_key", "Specify a TLS Key file")
	pflag.StringVar(&tlsCertFile, "tlscert", "/run/secrets/tls_crt", "Specify a TLS Cert file")
	pflag.BoolVar(&enableTLS, "tls", true, "Specify if TLS should be enabled")
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

func main() {
	pflag.Parse()
	v := viper.New()

	v.BindPFlags(pflag.CommandLine)

	configFilePath = v.GetString("config")
	enableTLS = v.GetBool("tls")
	tlsCertFile = v.GetString("tlscert")
	tlsKeyFile = v.GetString("tlskey")

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

	// Make sure the admin DB exists:
	adminDB := cfg.GetString(gather.CK_args_admindb_name.String())
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
	go restHandlerStart(gatherServer, cfg)
	gRPCHandlerStart(gatherServer, cfg)

}
