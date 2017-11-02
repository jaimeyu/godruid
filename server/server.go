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
	emp "github.com/golang/protobuf/ptypes/empty"
	wr "github.com/golang/protobuf/ptypes/wrappers"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
)

var (
	port     = flag.Int("port", 10000, "The server port")
	server   = flag.String("server", "", "Server interface to bind to")
	restPort = flag.Int("restPort", 10001, "The port used for REST operations")
)

// GatherServer - Server which will implement the gRPC Services.
type GatherServer struct {
}

// CreateAdminUser - Create an Administrative User.
func (s *GatherServer) CreateAdminUser(ctx context.Context, user *pb.AdminUser) (*pb.AdminUser, error) {
	// Stub to implement
	return nil, nil
}

// UpdateAdminUser - Update an Administrative User.
func (s *GatherServer) UpdateAdminUser(ctx context.Context, user *pb.AdminUser) (*pb.AdminUser, error) {
	// Stub to implement
	return nil, nil
}

// DeleteAdminUser - Delete an Administrative User.
func (s *GatherServer) DeleteAdminUser(ctx context.Context, userID *wr.StringValue) (*pb.AdminUser, error) {
	// Stub to implement
	return nil, nil
}

// GetAdminUser - Retrieve an Administrative User by the ID.
func (s *GatherServer) GetAdminUser(ctx context.Context, userID *wr.StringValue) (*pb.AdminUser, error) {
	// Just for test...return a fake AdminUser object.
	return &pb.AdminUser{
		XId:                   "fakeid",
		Username:              "best@admin.com",
		Password:              "fakey",
		SendOnboardingEmail:   true,
		OnboardingToken:       "somekeyvaluetoken",
		UserVerified:          false,
		State:                 1,
		CreatedTimestamp:      2346738246278,
		LastModifiedTimestamp: 689548964845}, nil
}

// GetAllAdminUsers -  Retrieve all Administrative Users.
func (s *GatherServer) GetAllAdminUsers(ctx context.Context, noValue *emp.Empty) (*pb.AdminUserList, error) {
	// Stub to implement
	return nil, nil
}

// CreateTenant - Create a Tenant. This will store the identification details for the Tenant,
// TenantDescriptor, as well as generate the Tenant Datastore for the
// Tenant data.
func (s *GatherServer) CreateTenant(ctx context.Context, tenantMeta *pb.TenantDescriptor) (*pb.TenantDescriptor, error) {
	// Stub to implement
	return nil, nil
}

// UpdateTenantDescriptor - Update the metadata for a Tenant.
func (s *GatherServer) UpdateTenantDescriptor(ctx context.Context, tenantMeta *pb.TenantDescriptor) (*pb.TenantDescriptor, error) {
	// Stub to implement
	return nil, nil
}

// DeleteTenant - Delete a Tenant by the provided ID. This operation will remove the Tenant
// datastore as well as the TenantDescriptor metadata.
func (s *GatherServer) DeleteTenant(ctx context.Context, tenantID *wr.StringValue) (*pb.TenantDescriptor, error) {
	// Stub to implement
	return nil, nil
}

//GetTenantDescriptor - retrieves Tenant metadata for the provided tenantID.
func (s *GatherServer) GetTenantDescriptor(ctx context.Context, tenantID *wr.StringValue) (*pb.TenantDescriptor, error) {
	// Stub to implement
	return nil, nil
}

func newServer() *GatherServer {
	s := new(GatherServer)

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
	pb.RegisterAdminProvisioningServiceServer(grpcServer, newServer())
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

	http.ListenAndServe(fmt.Sprintf(":%d", *restPort), mux)

}

func main() {
	flag.Parse()

	go restHandlerStart()
	gRPCHandlerStart()
}
