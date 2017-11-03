package main

import (
	"context"
	"flag"
	"fmt"
	"log"
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

	cfg := getActiveConfigOrExit()
	dbBindIP := cfg.ServerConfig.Datastore.BindIP
	dbBindPort := cfg.ServerConfig.Datastore.BindPort
	s.gsh = handlers.CreateCoordinator(fmt.Sprintf("%s:%d", dbBindIP, dbBindPort))

	return s
}

func gRPCHandlerStart() {
	grpcBindPort := getActiveConfigOrExit().ServerConfig.GRPC.BindPort
	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", grpcBindPort))
	if err != nil {
		log.Fatalf("failed to start gRPC Service: %v", err)
	}
	var opts []grpc.ServerOption

	grpcServer := grpc.NewServer(opts...)
	pb.RegisterAdminProvisioningServiceServer(grpcServer, newServer().gsh)

	log.Printf("gRPC service intiated on port: %d", grpcBindPort)
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
		log.Fatalf("failed to start REST service: %v", err)
	}

	log.Printf("REST service intiated on port: %d", restBindPort)

	http.ListenAndServe(fmt.Sprintf(":%d", restBindPort), mux)

}

func getActiveConfigOrExit() *gather.Config {
	cfg, err := gather.GetActiveConfig()
	if err != nil {
		log.Fatalf("failed to start Gather Service: %v", err)
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

	cfg := gather.LoadConfig(configFilePath)

	fmt.Printf("Your config is %+v \n", cfg)

	go restHandlerStart()
	gRPCHandlerStart()
}
