package datastore

import (
	pb "github.com/accedian/adh-gather/gathergrpc"
)

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
