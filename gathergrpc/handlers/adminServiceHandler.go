package handlers

import (
	"context"

	pb "github.com/accedian/adh-gather/gathergrpc"
	emp "github.com/golang/protobuf/ptypes/empty"
	wr "github.com/golang/protobuf/ptypes/wrappers"
)

// AdminServiceHandler - implementation of the interface for the gRPC
// Admin service. Anytime the Admin service changes, the logic to handle the
// API will be modified here.
type AdminServiceHandler struct {
}

// CreateAdminUser - Create an Administrative User.
func (ash *AdminServiceHandler) CreateAdminUser(ctx context.Context, user *pb.AdminUser) (*pb.AdminUser, error) {
	// Stub to implement
	return nil, nil
}

// UpdateAdminUser - Update an Administrative User.
func (ash *AdminServiceHandler) UpdateAdminUser(ctx context.Context, user *pb.AdminUser) (*pb.AdminUser, error) {
	// Stub to implement
	return nil, nil
}

// DeleteAdminUser - Delete an Administrative User.
func (ash *AdminServiceHandler) DeleteAdminUser(ctx context.Context, userID *wr.StringValue) (*pb.AdminUser, error) {
	// Stub to implement
	return nil, nil
}

// GetAdminUser - Retrieve an Administrative User by the ID.
func (ash *AdminServiceHandler) GetAdminUser(ctx context.Context, userID *wr.StringValue) (*pb.AdminUser, error) {
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
func (ash *AdminServiceHandler) GetAllAdminUsers(ctx context.Context, noValue *emp.Empty) (*pb.AdminUserList, error) {
	// Stub to implement
	return nil, nil
}

// CreateTenant - Create a Tenant. This will store the identification details for the Tenant,
// TenantDescriptor, as well as generate the Tenant Datastore for the
// Tenant data.
func (ash *AdminServiceHandler) CreateTenant(ctx context.Context, tenantMeta *pb.TenantDescriptor) (*pb.TenantDescriptor, error) {
	// Stub to implement
	return nil, nil
}

// UpdateTenantDescriptor - Update the metadata for a Tenant.
func (ash *AdminServiceHandler) UpdateTenantDescriptor(ctx context.Context, tenantMeta *pb.TenantDescriptor) (*pb.TenantDescriptor, error) {
	// Stub to implement
	return nil, nil
}

// DeleteTenant - Delete a Tenant by the provided ID. This operation will remove the Tenant
// datastore as well as the TenantDescriptor metadata.
func (ash *AdminServiceHandler) DeleteTenant(ctx context.Context, tenantID *wr.StringValue) (*pb.TenantDescriptor, error) {
	// Stub to implement
	return nil, nil
}

//GetTenantDescriptor - retrieves Tenant metadata for the provided tenantID.
func (ash *AdminServiceHandler) GetTenantDescriptor(ctx context.Context, tenantID *wr.StringValue) (*pb.TenantDescriptor, error) {
	// Stub to implement
	return nil, nil
}
