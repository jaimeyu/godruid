package main

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"strings"

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
)

func init() {
	pflag.StringVar(&configFilePath, "config", "config/adh-gather.yml", "Specify a configuration file to use")
}

// GatherServer - Server which will implement the gRPC Services.
type GatherServer struct {
	gsh     *adhh.GRPCServiceHandler
	pouchSH *adhh.PouchDBPluginServiceHandler

	mux   *mux.Router
	gwmux *runtime.ServeMux
}

func newServer() *GatherServer {
	s := new(GatherServer)
	s.gsh = adhh.CreateCoordinator()
	s.pouchSH = adhh.CreatePouchDBPluginServiceHandler()

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
	grpcBindIP := cfg.GetString(gather.CK_server_grpc_ip.String())
	grpcBindPort := cfg.GetInt(gather.CK_server_grpc_port.String())

	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	gatherServer.mux = mux.NewRouter().StrictSlash(true)

	gatherServer.gwmux = runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithInsecure()}

	// Register the Admin Service
	if err := pb.RegisterAdminProvisioningServiceHandlerFromEndpoint(ctx, gatherServer.gwmux, fmt.Sprintf("%s:%d", grpcBindIP, grpcBindPort), opts); err != nil {
		logger.Log.Fatalf("failed to start REST service: %s", err.Error())
	}

	// Register the Tenant Service
	if err := pb.RegisterTenantProvisioningServiceHandlerFromEndpoint(ctx, gatherServer.gwmux, fmt.Sprintf("%s:%d", grpcBindIP, grpcBindPort), opts); err != nil {
		logger.Log.Fatalf("failed to start REST service: %s", err.Error())
	}

	// Add in handling for non protobuf generated API endpoints:
	gatherServer.pouchSH.RegisterAPIHandlers(gatherServer.mux)

	logger.Log.Infof("REST service intiated on: %s:%d", restBindIP, restBindPort)
	originsOption := gh.AllowedOrigins([]string{"http://localhost:4200"})
	methodsOption := gh.AllowedMethods([]string{"GET", "HEAD", "POST", "PUT", "OPTIONS", "DELETE"})
	headersOption := gh.AllowedHeaders([]string{"accept", "authorization", "content-type", "origin", "referer", "x-csrf-token"})
	http.ListenAndServe(fmt.Sprintf(":%d", restBindPort), gh.CORS(originsOption, methodsOption, headersOption, gh.AllowCredentials())(gatherServer))

}

// Handle requests based on the path provided. If it begins with the known
// gRPC REST GW handler prefix, then use that handler, use the default handler
// otherwise.
func (gs *GatherServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if strings.Index(r.URL.Path, "/api/v1/") == 0 {
		gs.gwmux.ServeHTTP(w, r)
	} else {
		gs.mux.ServeHTTP(w, r)
	}
}

func main() {
	pflag.Parse()
	v := viper.New()

	v.BindPFlags(pflag.CommandLine)

	configFilePath := v.GetString("config")

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
	go restHandlerStart(gatherServer, cfg)
	gRPCHandlerStart(gatherServer, cfg)
}
