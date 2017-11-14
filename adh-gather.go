package main

import (
	"context"
	"flag"
	"fmt"
	"net"
	"net/http"

	"github.com/accedian/adh-gather/gather"
	"github.com/accedian/adh-gather/gathergrpc/handlers"
	"github.com/accedian/adh-gather/logger"
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
	gsh *handlers.GRPCServiceHandler
}

func newServer() *GatherServer {
	s := new(GatherServer)
	s.gsh = handlers.CreateCoordinator()

	return s
}

func gRPCHandlerStart() {
	grpcBindPort := getActiveConfigOrExit().ServerConfig.GRPC.BindPort
	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", grpcBindPort))
	if err != nil {
		logger.Log.Fatalf("failed to start gRPC Service: %s", err.Error())
	}
	var opts []grpc.ServerOption

	grpcServer := grpc.NewServer(opts...)
	grpcServiceHandler := newServer().gsh
	pb.RegisterAdminProvisioningServiceServer(grpcServer, grpcServiceHandler)
	pb.RegisterTenantProvisioningServiceServer(grpcServer, grpcServiceHandler)

	logger.Log.Infof("gRPC service intiated on port: %d", grpcBindPort)
	grpcServer.Serve(lis)
}

func restHandlerStart() {
	cfg := getActiveConfigOrExit()
	restBindPort := cfg.ServerConfig.REST.BindPort
	grpcBindPort := cfg.ServerConfig.GRPC.BindPort

	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	mux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithInsecure()}
	err := pb.RegisterAdminProvisioningServiceHandlerFromEndpoint(ctx, mux, fmt.Sprintf("localhost:%d", grpcBindPort), opts)
	if err != nil {
		logger.Log.Fatalf("failed to start REST service: %s", err.Error())
	}
	err = pb.RegisterTenantProvisioningServiceHandlerFromEndpoint(ctx, mux, fmt.Sprintf("localhost:%d", grpcBindPort), opts)
	if err != nil {
		logger.Log.Fatalf("failed to start REST service: %s", err.Error())
	}

	logger.Log.Infof("REST service intiated on port: %d", restBindPort)

	http.ListenAndServe(fmt.Sprintf(":%d", restBindPort), mux)

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
	go restHandlerStart()
	gRPCHandlerStart()
}
