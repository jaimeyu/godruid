package datastore

import (
	pb "github.com/accedian/adh-gather/gathergrpc"
)

// AdminUserStr - common name of the AdminUser data type for use in logs.
const AdminUserStr string = "Admin User"

// TenantDescriptorStr - common name of the TenantDescriptor data type for use in logs.
const TenantDescriptorStr string = "Tenant Descriptor"

// TenantStr - common name of the TenantDescriptor data type for use in logs.
const TenantStr string = "Tenant"

// AdminServiceDatastore - interface which provides the functionality
// of the AdminService Datastore.
type AdminServiceDatastore interface {
	CreateAdminUser(*pb.AdminUserRequest) (*pb.AdminUserResponse, error)
	UpdateAdminUser(*pb.AdminUserRequest) (*pb.AdminUserResponse, error)
	DeleteAdminUser(string) (*pb.AdminUserResponse, error)
	GetAdminUser(string) (*pb.AdminUserResponse, error)
	GetAllAdminUsers() (*pb.AdminUserListResponse, error)

	CreateTenant(*pb.TenantDescriptorRequest) (*pb.TenantDescriptorResponse, error)
	UpdateTenantDescriptor(*pb.TenantDescriptorRequest) (*pb.TenantDescriptorResponse, error)
	DeleteTenant(string) (*pb.TenantDescriptorResponse, error)
	GetTenantDescriptor(string) (*pb.TenantDescriptorResponse, error)
}
