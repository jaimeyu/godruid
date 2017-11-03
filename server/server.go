package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"

	"golang.org/x/net/context"
	"google.golang.org/grpc"

	pb "github.com/accedian/adh-gather/gathergrpc"
	"github.com/accedian/adh-gather/gathergrpc/handlers"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
)

var (
	port       = flag.Int("port", 10000, "The server port")
	server     = flag.String("server", "", "Server interface to bind to")
	restPort   = flag.Int("restPort", 10001, "The port used for REST operations")
	provDB     = flag.String("provDB", "http://localhost", "The server to connect to for provisioning details.")
	provDBPort = flag.Int("provDBPort", 5984, "The port used for operations with Provisioning Database")
)

// GatherServer - Server which will implement the gRPC Services.
type GatherServer struct {
	gsh *handlers.GRPCServiceHandler
}

func newServer() *GatherServer {
	s := new(GatherServer)
	s.gsh = handlers.CreateCoordinator(fmt.Sprintf("%s:%d", *provDB, *provDBPort))

	// TODO: Load in config and stuff here.

	return s
}

func gRPCHandlerStart() {
	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", *port))
	if err != nil {
		log.Fatalf("failed to start gRPC Service: %v", err)
	}
	var opts []grpc.ServerOption

	grpcServer := grpc.NewServer(opts...)
	pb.RegisterAdminProvisioningServiceServer(grpcServer, newServer().gsh)

	log.Printf("gRPC service intiated on port: %d", *port)
	grpcServer.Serve(lis)
}

func restHandlerStart() {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	mux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithInsecure()}
	err := pb.RegisterAdminProvisioningServiceHandlerFromEndpoint(ctx, mux, fmt.Sprintf("localhost:%d", *port), opts)
	if err != nil {
		log.Fatalf("failed to start REST service: %v", err)
	}

	log.Printf("REST service intiated on port: %d", *restPort)

	http.ListenAndServe(fmt.Sprintf(":%d", *restPort), mux)

}

func main() {
	flag.Parse()

	go restHandlerStart()
	gRPCHandlerStart()
}
