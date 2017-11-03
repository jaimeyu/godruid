package handlers

import (
	"context"

	pb "github.com/accedian/adh-gather/gathergrpc"
	emp "github.com/golang/protobuf/ptypes/empty"
	wr "github.com/golang/protobuf/ptypes/wrappers"
)

// GRPCServiceHandler - implementer of all gRPC Services. Offloads
// implementation details to each unique service handler. When new
// gRPC services are added, a new Service Handler should be created,
// and a pointer to that object should be added to this wrapper.
type GRPCServiceHandler struct {
	ash *AdminServiceHandler
}

// CreateCoordinator - used to create a gRPC service handler wrapper
// that coordinates the logic to satisfy all gRPC service
// interfaces.
func CreateCoordinator(provDBURL string) *GRPCServiceHandler {
	result := new(GRPCServiceHandler)
	result.ash = CreateHandler(provDBURL)
	return result
}

// CreateAdminUser - Create an Administrative User.
func (gsh *GRPCServiceHandler) CreateAdminUser(ctx context.Context, user *pb.AdminUser) (*pb.AdminUser, error) {
	return gsh.ash.CreateAdminUser(ctx, user)
}

// UpdateAdminUser - Update an Administrative User.
func (gsh *GRPCServiceHandler) UpdateAdminUser(ctx context.Context, user *pb.AdminUser) (*pb.AdminUser, error) {
	return gsh.ash.UpdateAdminUser(ctx, user)
}

// DeleteAdminUser - Delete an Administrative User.
func (gsh *GRPCServiceHandler) DeleteAdminUser(ctx context.Context, userID *wr.StringValue) (*pb.AdminUser, error) {
	return gsh.ash.DeleteAdminUser(ctx, userID)
}

// GetAdminUser - Retrieve an Administrative User by the ID.
func (gsh *GRPCServiceHandler) GetAdminUser(ctx context.Context, userID *wr.StringValue) (*pb.AdminUser, error) {
	return gsh.ash.GetAdminUser(ctx, userID)
}

// GetAllAdminUsers -  Retrieve all Administrative Users.
func (gsh *GRPCServiceHandler) GetAllAdminUsers(ctx context.Context, noValue *emp.Empty) (*pb.AdminUserList, error) {
	return gsh.ash.GetAllAdminUsers(ctx, noValue)
}

// CreateTenant - Create a Tenant. This will store the identification details for the Tenant,
// TenantDescriptor, as well as generate the Tenant Datastore for the
// Tenant data.
func (gsh *GRPCServiceHandler) CreateTenant(ctx context.Context, tenantMeta *pb.TenantDescriptor) (*pb.TenantDescriptor, error) {
	return gsh.ash.CreateTenant(ctx, tenantMeta)
}

// UpdateTenantDescriptor - Update the metadata for a Tenant.
func (gsh *GRPCServiceHandler) UpdateTenantDescriptor(ctx context.Context, tenantMeta *pb.TenantDescriptor) (*pb.TenantDescriptor, error) {
	return gsh.ash.UpdateTenantDescriptor(ctx, tenantMeta)
}

// DeleteTenant - Delete a Tenant by the provided ID. This operation will remove the Tenant
// datastore as well as the TenantDescriptor metadata.
func (gsh *GRPCServiceHandler) DeleteTenant(ctx context.Context, tenantID *wr.StringValue) (*pb.TenantDescriptor, error) {
	return gsh.ash.DeleteTenant(ctx, tenantID)
}

//GetTenantDescriptor - retrieves Tenant metadata for the provided tenantID.
func (gsh *GRPCServiceHandler) GetTenantDescriptor(ctx context.Context, tenantID *wr.StringValue) (*pb.TenantDescriptor, error) {
	return gsh.ash.GetTenantDescriptor(ctx, tenantID)
}
