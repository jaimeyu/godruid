package main

import (
	"flag"
	"fmt"
	"log"
	"net"

	"golang.org/x/net/context"
	"google.golang.org/grpc"

	pb "github.com/accedian/adh-gather/gathergrpc"
	emp "github.com/golang/protobuf/ptypes/empty"
	wr "github.com/golang/protobuf/ptypes/wrappers"
)

var (
	port   = flag.Int("port", 10000, "The server port")
	server = flag.String("server", "", "Server interface to bind to")
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
	// Stub to implement
	return nil, nil
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

func main() {
	flag.Parse()
	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	var opts []grpc.ServerOption

	grpcServer := grpc.NewServer(opts...)
	pb.RegisterAdminProvisioningServiceServer(grpcServer, newServer())
	grpcServer.Serve(lis)
}
