package main

import (
	"context"
	"flag"
	"fmt"
	"net"
	"net/http"
	"strings"

	"github.com/accedian/adh-gather/gather"
	"github.com/accedian/adh-gather/handlers"
	"github.com/accedian/adh-gather/logger"
	"github.com/gorilla/mux"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"google.golang.org/grpc"

	pb "github.com/accedian/adh-gather/gathergrpc"
)

var (
	configFilePath string
	debug          bool
)

func init() {
	flag.StringVar(&configFilePath, "config", "config/adh-gather.yml", "Specify a configuration file to use")
	flag.BoolVar(&debug, "debug", false, "Enable debug mode (and logs)")
}

// GatherServer - Server which will implement the gRPC Services.
type GatherServer struct {
	gsh     *handlers.GRPCServiceHandler
	pouchSH *handlers.PouchDBPluginServiceHandler

	mux   *mux.Router
	gwmux *runtime.ServeMux
}

func newServer() *GatherServer {
	s := new(GatherServer)
	s.gsh = handlers.CreateCoordinator()
	s.pouchSH = handlers.CreatePouchDBPluginServiceHandler()

	return s
}

func gRPCHandlerStart(gatherServer *GatherServer) {
	grpcBindPort := getActiveConfigOrExit().ServerConfig.GRPC.BindPort
	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", grpcBindPort))
	if err != nil {
		logger.Log.Fatalf("failed to start gRPC Service: %s", err.Error())
	}
	var opts []grpc.ServerOption

	grpcServer := grpc.NewServer(opts...)
	pb.RegisterAdminProvisioningServiceServer(grpcServer, gatherServer.gsh)
	pb.RegisterTenantProvisioningServiceServer(grpcServer, gatherServer.gsh)
	// pb.RegisterPouchDBPluginServiceServer(grpcServer, gatherServer.pouchSH)

	logger.Log.Infof("gRPC service intiated on port: %d", grpcBindPort)
	grpcServer.Serve(lis)
}

func restHandlerStart(gatherServer *GatherServer) {
	cfg := getActiveConfigOrExit()
	restBindPort := cfg.ServerConfig.REST.BindPort
	grpcBindPort := cfg.ServerConfig.GRPC.BindPort

	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	gatherServer.mux = mux.NewRouter().StrictSlash(true)

	gatherServer.gwmux = runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithInsecure()}

	// // Setup the overall handler:
	// mux := gatherServer.pouchSH.Router

	// Register handlers for the Non-gRPC REST GW calls.

	// Register the Admin Service
	if err := pb.RegisterAdminProvisioningServiceHandlerFromEndpoint(ctx, gatherServer.gwmux, fmt.Sprintf("localhost:%d", grpcBindPort), opts); err != nil {
		logger.Log.Fatalf("failed to start REST service: %s", err.Error())
	}

	// Register the Tenant Service
	if err := pb.RegisterTenantProvisioningServiceHandlerFromEndpoint(ctx, gatherServer.gwmux, fmt.Sprintf("localhost:%d", grpcBindPort), opts); err != nil {
		logger.Log.Fatalf("failed to start REST service: %s", err.Error())
	}

	// Add in handling for non protobuf generated API endpoints:
	gatherServer.pouchSH.RegisterAPIHandlers(gatherServer.mux)

	// // Register the PouchDBPlugin Service
	// if err := pb.RegisterPouchDBPluginServiceHandlerFromEndpoint(ctx, gwmux, fmt.Sprintf("localhost:%d", grpcBindPort), opts); err != nil {
	// 	logger.Log.Fatalf("failed to start REST service: %s", err.Error())
	// }

	// Handle all the generated gRPC REST GW calls from the same overall mux
	// mux.Handle("/api/v1/", gwmux)

	logger.Log.Infof("REST service intiated on port: %d", restBindPort)
	http.ListenAndServe(fmt.Sprintf(":%d", restBindPort), gatherServer)

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

func getActiveConfigOrExit() *gather.Config {
	cfg, err := gather.GetActiveConfig()
	if err != nil {
		logger.Log.Fatalf("failed to start Gather Service: %s", err.Error())
	}

	return cfg
}

func main() {
	flag.Parse()

	if debug {
		logger.SetDebugLevel(true)
	} else {
		logger.SetDebugLevel(false)
	}

	logger.Log.Infof("Starting adh-gather broker with config '%s'", configFilePath)

	// Load Configuration
	cfg := gather.LoadConfig(configFilePath)
	fmt.Printf("Your config is %+v \n", cfg)

	// Start the REST and gRPC Services
	gatherServer := newServer()
	go restHandlerStart(gatherServer)
	gRPCHandlerStart(gatherServer)
}
