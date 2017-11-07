package handlers

import (
	"context"
	"errors"

	"github.com/accedian/adh-gather/gather"

	"github.com/accedian/adh-gather/datastore"
	"github.com/accedian/adh-gather/datastore/couchDB"
	"github.com/accedian/adh-gather/datastore/inMemory"

	db "github.com/accedian/adh-gather/datastore"
	pb "github.com/accedian/adh-gather/gathergrpc"
	"github.com/accedian/adh-gather/logger"
	emp "github.com/golang/protobuf/ptypes/empty"
	wr "github.com/golang/protobuf/ptypes/wrappers"
)

const dbName string = "adh-admin"

// AdminServiceHandler - implementation of the interface for the gRPC
// Admin service. Anytime the Admin service changes, the logic to handle the
// API will be modified here.
type AdminServiceHandler struct {
	adminDB db.AdminServiceDatastore
}

// CreateAdminServiceHandler - used to generate a handler for the Admin Service.
func CreateAdminServiceHandler() *AdminServiceHandler {
	result := new(AdminServiceHandler)

	// Seteup the DB implementation based on configuration
	db, err := getAdminServiceDatastore()
	if err != nil {
		logger.Log.Fatalf("Unable to instantiate AdminServiceHandler: %v", err)
	}
	result.adminDB = db

	return result
}

// CreateAdminUser - Create an Administrative User.
func (ash *AdminServiceHandler) CreateAdminUser(ctx context.Context, user *pb.AdminUser) (*pb.AdminUser, error) {
	// Perform any validation here:
	logger.Log.Infof("Creating %s: %s", datastore.AdminUserStr, user)

	// Issue request to DAO Layer to Create the Admin User
	result, err := ash.adminDB.CreateAdminUser(user)
	if err != nil {
		logger.Log.Errorf("Unable to store %s: %v\n", datastore.AdminUserStr, err)
		return nil, err
	}

	// Succesfully Created the User, return the result.
	logger.Log.Infof("Created %s: %v\n", datastore.AdminUserStr, result)
	return result, nil
}

// UpdateAdminUser - Update an Administrative User.
func (ash *AdminServiceHandler) UpdateAdminUser(ctx context.Context, user *pb.AdminUser) (*pb.AdminUser, error) {
	// Perform any validation here:
	logger.Log.Infof("Updating %s: %s", datastore.AdminUserStr, user)

	// Issue request to DAO Layer to Create the Admin User
	result, err := ash.adminDB.UpdateAdminUser(user)
	if err != nil {
		logger.Log.Errorf("Unable to update %s: %v\n", datastore.AdminUserStr, err)
		return nil, err
	}

	// Succesfully Updated the User, return the result.
	logger.Log.Infof("Updated %s: %v\n", datastore.AdminUserStr, result)
	return result, nil
}

// DeleteAdminUser - Delete an Administrative User.
func (ash *AdminServiceHandler) DeleteAdminUser(ctx context.Context, userID *wr.StringValue) (*pb.AdminUser, error) {
	// Perform any validation here:

	// Issue request to DAO Layer to Create the Admin User
	result, err := ash.adminDB.DeleteAdminUser(userID.Value)
	if err != nil {
		logger.Log.Errorf("Unable to delete %s: %v\n", datastore.AdminUserStr, err)
		return nil, err
	}

	// Succesfully Deleted the User, return the result.
	logger.Log.Infof("Deleted %s: %v\n", datastore.AdminUserStr, result)
	return result, nil
}

// GetAdminUser - Retrieve an Administrative User by the ID.
func (ash *AdminServiceHandler) GetAdminUser(ctx context.Context, userID *wr.StringValue) (*pb.AdminUser, error) {
	// Perform and validation here:
	logger.Log.Infof("Retrieving %s: %s", datastore.AdminUserStr, userID.Value)

	// Issue request to DAO Layer to Get the requested Admin User
	result, err := ash.adminDB.GetAdminUser(userID.Value)
	if err != nil {
		logger.Log.Errorf("Unable to retrieve %s: %v\n", datastore.AdminUserStr, err)
		return nil, err
	}

	// Succesfully found the User, return the result.
	logger.Log.Infof("Retrieved %s: %v\n", datastore.AdminUserStr, result)
	return result, nil
}

// GetAllAdminUsers -  Retrieve all Administrative Users.
func (ash *AdminServiceHandler) GetAllAdminUsers(ctx context.Context, noValue *emp.Empty) (*pb.AdminUserList, error) {
	// Perform and validation here:
	logger.Log.Infof("Retrieving all %ss", datastore.AdminUserStr)

	// Issue request to DAO Layer to Get the requested Admin User List
	result, err := ash.adminDB.GetAllAdminUsers()
	if err != nil {
		logger.Log.Errorf("Unable to retrieve %ss: %v\n", datastore.AdminUserStr, err)
		return nil, err
	}

	// Succesfully found the Users, return the result list.
	logger.Log.Infof("Retrieved %d %ss\n", len(result.GetList()), datastore.AdminUserStr)
	return result, nil
}

// CreateTenant - Create a Tenant. This will store the identification details for the Tenant,
// TenantDescriptor, as well as generate the Tenant Datastore for the
// Tenant data.
func (ash *AdminServiceHandler) CreateTenant(ctx context.Context, tenantMeta *pb.TenantDescriptor) (*pb.TenantDescriptor, error) {
	// Perform and validation here:
	logger.Log.Infof("Creating %s: %v", datastore.TenantStr, tenantMeta)

	// Issue request to AdminService DAO to create the metadata record:
	result, err := ash.adminDB.CreateTenant(tenantMeta)
	if err != nil {
		logger.Log.Errorf("Unable to create %s: %v\n", datastore.TenantStr, err)
		return nil, err
	}

	// TODO: Issue request to TenantService Datastore to create necessary Tenant
	// data objects

	// Succesfully Createds the Tenant, return the metadata result.
	logger.Log.Infof("Created %s: %v\n", datastore.TenantStr, result)
	return result, nil
}

// UpdateTenantDescriptor - Update the metadata for a Tenant.
func (ash *AdminServiceHandler) UpdateTenantDescriptor(ctx context.Context, tenantMeta *pb.TenantDescriptor) (*pb.TenantDescriptor, error) {
	// Perform and validation here:
	logger.Log.Infof("Updating %s: %v", datastore.TenantDescriptorStr, tenantMeta)

	// Issue request to AdminService DAO to update the metadata record:
	result, err := ash.adminDB.UpdateTenantDescriptor(tenantMeta)
	if err != nil {
		logger.Log.Errorf("Unable to update %s: %v\n", datastore.TenantDescriptorStr, err)
		return nil, err
	}

	// Succesfully Createds the Tenant, return the metadata result.
	logger.Log.Infof("Updated %s: %v\n", datastore.TenantDescriptorStr, result)
	return result, nil
}

// DeleteTenant - Delete a Tenant by the provided ID. This operation will remove the Tenant
// datastore as well as the TenantDescriptor metadata.
func (ash *AdminServiceHandler) DeleteTenant(ctx context.Context, tenantID *wr.StringValue) (*pb.TenantDescriptor, error) {
	// Perform and validation here:
	logger.Log.Infof("Attempting to delete %s: %s", datastore.TenantStr, tenantID.Value)

	// Issue request to DAO Layer to Delete the requested Tenant
	result, err := ash.adminDB.DeleteTenant(tenantID.Value)
	if err != nil {
		logger.Log.Errorf("Unable to delete %s: %v\n", datastore.TenantStr, err)
		return nil, err
	}

	// TODO: Add calls here to Tenant Service to delete any related
	// tenant data.

	// Succesfully removed the Tenant, return the metadata that identified
	// the now deleted tenant.
	logger.Log.Infof("Successfully deleted %s. Previous %s: %v\n", datastore.TenantStr, datastore.TenantDescriptorStr, result)
	return result, nil
}

//GetTenantDescriptor - retrieves Tenant metadata for the provided tenantID.
func (ash *AdminServiceHandler) GetTenantDescriptor(ctx context.Context, tenantID *wr.StringValue) (*pb.TenantDescriptor, error) {
	// Perform and validation here:
	logger.Log.Infof("Retrieving %s: %s", datastore.TenantDescriptorStr, tenantID.Value)

	// Issue request to DAO Layer to Get the requested Tenant Metadata
	result, err := ash.adminDB.GetTenantDescriptor(tenantID.Value)
	if err != nil {
		logger.Log.Errorf("Unable to retrieve %s: %v\n", datastore.TenantDescriptorStr, err)
		return nil, err
	}

	// Succesfully found the Tenant Metadata, return the result.
	logger.Log.Infof("Retrieved %s: %v\n", datastore.TenantDescriptorStr, result)
	return result, nil
}

func getAdminServiceDatastore() (datastore.AdminServiceDatastore, error) {
	cfg, err := gather.GetActiveConfig()
	if err != nil {
		logger.Log.Errorf("Falied to instantiate AdminServiceHandler: %v", err)
		return nil, err
	}

	dbType := cfg.ServerConfig.StartupArgs.AdminDB
	if dbType == gather.COUCH {
		logger.Log.Debug("AdminService DB is using CouchDB Implementation")
		return couchDB.CreateAdminServiceDAO(), nil
	} else if dbType == gather.MEM {
		logger.Log.Debug("AdminService DB is using InMemory Implementation")
		return inMemory.CreateAdminServiceDAO(), nil
	}

	return nil, errors.New("No DB implementation provided for Admin Service. Check configuration")
}
