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

// AdminUserType - datatype string used to identify an Admin User in the datastore record
const AdminUserType string = "adminUser"

// TenantDescriptorType - datatype string used to identify an Tenant Descriptor in the datastore record
const TenantDescriptorType string = "tenant"

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
	GetAllTenantDescriptors() (*pb.TenantDescriptorListResponse, error)
}
