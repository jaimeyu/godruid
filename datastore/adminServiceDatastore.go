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
	CreateAdminUser(*pb.AdminUser) (*pb.AdminUser, error)
	UpdateAdminUser(*pb.AdminUser) (*pb.AdminUser, error)
	DeleteAdminUser(string) (*pb.AdminUser, error)
	GetAdminUser(string) (*pb.AdminUser, error)
	GetAllAdminUsers() (*pb.AdminUserList, error)

	CreateTenant(*pb.TenantDescriptor) (*pb.TenantDescriptor, error)
	UpdateTenantDescriptor(*pb.TenantDescriptor) (*pb.TenantDescriptor, error)
	DeleteTenant(string) (*pb.TenantDescriptor, error)
	GetTenantDescriptor(string) (*pb.TenantDescriptor, error)
}
