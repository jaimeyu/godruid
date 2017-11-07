package handlers

import (
	"context"

	pb "github.com/accedian/adh-gather/gathergrpc"
	"github.com/accedian/adh-gather/logger"
	emp "github.com/golang/protobuf/ptypes/empty"
	wr "github.com/golang/protobuf/ptypes/wrappers"
)

// GRPCServiceHandler - implementer of all gRPC Services. Offloads
// implementation details to each unique service handler. When new
// gRPC services are added, a new Service Handler should be created,
// and a pointer to that object should be added to this wrapper.
type GRPCServiceHandler struct {
	ash *AdminServiceHandler
	tsh *TenantServiceHandler
}

// CreateCoordinator - used to create a gRPC service handler wrapper
// that coordinates the logic to satisfy all gRPC service
// interfaces.
func CreateCoordinator() *GRPCServiceHandler {
	result := new(GRPCServiceHandler)

	result.ash = CreateAdminServiceHandler()
	result.tsh = CreateTenantServiceHandler()

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
	// Create the Tenant metadata record and reserve space to store isolated Tenant data
	result, err := gsh.ash.CreateTenant(ctx, tenantMeta)
	if err != nil {
		return nil, err
	}

	// Create a default Ingestion Profile for the Tenant.
	ingPrfReq := pb.TenantIngestionProfileRequest{TenantId: result.GetId(), IngestionProfile: createDefaultTenantIngPrf()}
	_, err = gsh.tsh.CreateTenantIngestionProfile(ctx, &ingPrfReq)
	if err != nil {
		logger.Log.Errorf("Unable to create Ingestion Profile for Tenant %s. The Tenant does exist though, so may need to create the Ingestion Profile manually", result.GetId())
	}

	return result, nil
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

// CreateTenantUser - creates a user scoped to a single Tenant.
func (gsh *GRPCServiceHandler) CreateTenantUser(ctx context.Context, tenantUserReq *pb.TenantUserRequest) (*pb.TenantUserResponse, error) {
	return gsh.tsh.CreateTenantUser(ctx, tenantUserReq)
}

// UpdateTenantUser - updates a user scoped to a single Tenant.
func (gsh *GRPCServiceHandler) UpdateTenantUser(ctx context.Context, tenantUserReq *pb.TenantUserRequest) (*pb.TenantUserResponse, error) {
	return gsh.tsh.UpdateTenantUser(ctx, tenantUserReq)
}

// DeleteTenantUser - deletes a user scoped to a single Tenant.
func (gsh *GRPCServiceHandler) DeleteTenantUser(ctx context.Context, tenantUserIdReq *pb.TenantUserIdRequest) (*pb.TenantUserResponse, error) {
	return gsh.tsh.DeleteTenantUser(ctx, tenantUserIdReq)
}

// GetTenantUser - retrieves a user scoped to a single Tenant.
func (gsh *GRPCServiceHandler) GetTenantUser(ctx context.Context, tenantUserIdReq *pb.TenantUserIdRequest) (*pb.TenantUserResponse, error) {
	return gsh.tsh.GetTenantUser(ctx, tenantUserIdReq)
}

// GetAllTenantUsers - retrieves all users scoped to a single Tenant.
func (gsh *GRPCServiceHandler) GetAllTenantUsers(ctx context.Context, tenantID *wr.StringValue) (*pb.TenantUserListResponse, error) {
	return gsh.tsh.GetAllTenantUsers(ctx, tenantID)
}

// CreateTenantDomain - creates a Domain scoped to a single Tenant.
func (gsh *GRPCServiceHandler) CreateTenantDomain(ctx context.Context, tenantDomainRequest *pb.TenantDomainRequest) (*pb.TenantDomainResponse, error) {
	return gsh.tsh.CreateTenantDomain(ctx, tenantDomainRequest)
}

// UpdateTenantDomain - updates a Domain scoped to a single Tenant.
func (gsh *GRPCServiceHandler) UpdateTenantDomain(ctx context.Context, tenantDomainRequest *pb.TenantDomainRequest) (*pb.TenantDomainResponse, error) {
	return gsh.tsh.UpdateTenantDomain(ctx, tenantDomainRequest)
}

// DeleteTenantDomain - deletes a Domain scoped to a single Tenant.
func (gsh *GRPCServiceHandler) DeleteTenantDomain(ctx context.Context, tenantDomainIDRequest *pb.TenantDomainIdRequest) (*pb.TenantDomainResponse, error) {
	return gsh.tsh.DeleteTenantDomain(ctx, tenantDomainIDRequest)
}

// GetTenantDomain - retrieves a Domain scoped to a single Tenant.
func (gsh *GRPCServiceHandler) GetTenantDomain(ctx context.Context, tenantDomainIDRequest *pb.TenantDomainIdRequest) (*pb.TenantDomainResponse, error) {
	return gsh.tsh.GetTenantDomain(ctx, tenantDomainIDRequest)
}

// GetAllTenantDomains - retrieves all Domains scoped to a single Tenant.
func (gsh *GRPCServiceHandler) GetAllTenantDomains(ctx context.Context, tenantID *wr.StringValue) (*pb.TenantDomainListResponse, error) {
	return gsh.tsh.GetAllTenantDomains(ctx, tenantID)
}

// CreateTenantIngestionProfile - updates an Ingestion Profile scoped to a specific Tenant.
func (gsh *GRPCServiceHandler) CreateTenantIngestionProfile(ctx context.Context, tenantIngPrfReq *pb.TenantIngestionProfileRequest) (*pb.TenantIngestionProfileResponse, error) {
	return gsh.tsh.CreateTenantIngestionProfile(ctx, tenantIngPrfReq)
}

// UpdateTenantIngestionProfile - updates an Ingestion Profile scoped to a specific Tenant.
func (gsh *GRPCServiceHandler) UpdateTenantIngestionProfile(ctx context.Context, tenantIngPrfReq *pb.TenantIngestionProfileRequest) (*pb.TenantIngestionProfileResponse, error) {
	return gsh.tsh.UpdateTenantIngestionProfile(ctx, tenantIngPrfReq)
}

// GetTenantIngestionProfile - retrieves the Ingestion Profile for a singler Tenant.
func (gsh *GRPCServiceHandler) GetTenantIngestionProfile(ctx context.Context, tenantID *pb.TenantIngestionProfileIdRequest) (*pb.TenantIngestionProfileResponse, error) {
	return gsh.tsh.GetTenantIngestionProfile(ctx, tenantID)
}

// DeleteTenantIngestionProfile - retrieves the Ingestion Profile for a singler Tenant.
func (gsh *GRPCServiceHandler) DeleteTenantIngestionProfile(ctx context.Context, tenantID *pb.TenantIngestionProfileIdRequest) (*pb.TenantIngestionProfileResponse, error) {
	return gsh.tsh.DeleteTenantIngestionProfile(ctx, tenantID)
}

func createDefaultTenantIngPrf() *pb.TenantIngestionProfile {
	ingPrf := pb.TenantIngestionProfile{}
	ingPrf.ScpUsername = "default"
	ingPrf.ScpPassword = "password"

	return &ingPrf
}
