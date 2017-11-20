package handlers

import (
	"context"
	"errors"
	"fmt"

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
		logger.Log.Fatalf("Unable to instantiate AdminServiceHandler: %s", err.Error())
	}
	result.adminDB = db

	return result
}

// CreateAdminUser - Create an Administrative User.
func (ash *AdminServiceHandler) CreateAdminUser(ctx context.Context, user *pb.AdminUserRequest) (*pb.AdminUserResponse, error) {
	// Validate the request to ensure no invalid data is stored:
	if err := validateAdminUserRequest(user); err != nil {
		return nil, err
	}
	logger.Log.Infof("Creating %s: %s", datastore.AdminUserStr, user)

	// Issue request to DAO Layer to Create the Admin User
	result, err := ash.adminDB.CreateAdminUser(user)
	if err != nil {
		return nil, fmt.Errorf("Unable to store %s: %s", datastore.AdminUserStr, err.Error())
	}

	// Succesfully Created the User, return the result.
	logger.Log.Infof("Created %s: %s\n", datastore.AdminUserStr, result.GetXId())
	return result, nil
}

// UpdateAdminUser - Update an Administrative User.
func (ash *AdminServiceHandler) UpdateAdminUser(ctx context.Context, user *pb.AdminUserRequest) (*pb.AdminUserResponse, error) {
	// Validate the request to ensure no invalid data is stored:
	if err := validateAdminUserRequest(user); err != nil {
		return nil, err
	}
	logger.Log.Infof("Updating %s: %s", datastore.AdminUserStr, user)

	// Issue request to DAO Layer to Create the Admin User
	result, err := ash.adminDB.UpdateAdminUser(user)
	if err != nil {
		return nil, fmt.Errorf("Unable to update %s: %s", datastore.AdminUserStr, err.Error())
	}

	// Succesfully Updated the User, return the result.
	logger.Log.Infof("Updated %s: %s\n", datastore.AdminUserStr, result.GetXId())
	return result, nil
}

// DeleteAdminUser - Delete an Administrative User.
func (ash *AdminServiceHandler) DeleteAdminUser(ctx context.Context, userID *wr.StringValue) (*pb.AdminUserResponse, error) {
	// Perform any validation here:

	// Issue request to DAO Layer to Create the Admin User
	result, err := ash.adminDB.DeleteAdminUser(userID.Value)
	if err != nil {
		return nil, fmt.Errorf("Unable to delete %s: %s", datastore.AdminUserStr, err.Error())
	}

	// Succesfully Deleted the User, return the result.
	logger.Log.Infof("Deleted %s: %s\n", datastore.AdminUserStr, result.GetXId())
	return result, nil
}

// GetAdminUser - Retrieve an Administrative User by the ID.
func (ash *AdminServiceHandler) GetAdminUser(ctx context.Context, userID *wr.StringValue) (*pb.AdminUserResponse, error) {
	// Perform and validation here:
	logger.Log.Infof("Retrieving %s: %s", datastore.AdminUserStr, userID.Value)

	// Issue request to DAO Layer to Get the requested Admin User
	result, err := ash.adminDB.GetAdminUser(userID.Value)
	if err != nil {
		return nil, fmt.Errorf("Unable to retrieve %s: %s", datastore.AdminUserStr, err.Error())
	}

	// Succesfully found the User, return the result.
	logger.Log.Infof("Retrieved %s: %s\n", datastore.AdminUserStr, result.GetXId())
	return result, nil
}

// GetAllAdminUsers -  Retrieve all Administrative Users.
func (ash *AdminServiceHandler) GetAllAdminUsers(ctx context.Context, noValue *emp.Empty) (*pb.AdminUserListResponse, error) {
	// Perform and validation here:
	logger.Log.Infof("Retrieving all %ss", datastore.AdminUserStr)

	// Issue request to DAO Layer to Get the requested Admin User List
	result, err := ash.adminDB.GetAllAdminUsers()
	if err != nil {
		return nil, fmt.Errorf("Unable to retrieve %ss: %s", datastore.AdminUserStr, err.Error())
	}

	// Succesfully found the Users, return the result list.
	logger.Log.Infof("Retrieved %d %ss\n", len(result.GetData()), datastore.AdminUserStr)
	return result, nil
}

// CreateTenant - Create a Tenant. This will store the identification details for the Tenant,
// TenantDescriptor, as well as generate the Tenant Datastore for the
// Tenant data.
func (ash *AdminServiceHandler) CreateTenant(ctx context.Context, tenantMeta *pb.TenantDescriptorRequest) (*pb.TenantDescriptorResponse, error) {
	// Validate the request to ensure no invalid data is stored:
	if err := validateTenantDescriptorRequest(tenantMeta); err != nil {
		return nil, err
	}
	logger.Log.Infof("Creating %s: %s", datastore.TenantStr, tenantMeta.GetXId())

	// Issue request to AdminService DAO to create the metadata record:
	result, err := ash.adminDB.CreateTenant(tenantMeta)
	if err != nil {
		return nil, fmt.Errorf("Unable to create %s: %s", datastore.TenantStr, err.Error())
	}

	// Succesfully Created the Tenant, return the metadata result.
	logger.Log.Infof("Created %s: %s\n", datastore.TenantStr, result.GetXId())
	return result, nil
}

// UpdateTenantDescriptor - Update the metadata for a Tenant.
func (ash *AdminServiceHandler) UpdateTenantDescriptor(ctx context.Context, tenantMeta *pb.TenantDescriptorRequest) (*pb.TenantDescriptorResponse, error) {
	// Validate the request to ensure no invalid data is stored:
	if err := validateTenantDescriptorRequest(tenantMeta); err != nil {
		return nil, err
	}
	logger.Log.Infof("Updating %s: %s", datastore.TenantDescriptorStr, tenantMeta.GetXId())

	// Issue request to AdminService DAO to update the metadata record:
	result, err := ash.adminDB.UpdateTenantDescriptor(tenantMeta)
	if err != nil {
		return nil, fmt.Errorf("Unable to update %s: %s", datastore.TenantDescriptorStr, err.Error())
	}

	// Succesfully Createds the Tenant, return the metadata result.
	logger.Log.Infof("Updated %s: %s\n", datastore.TenantDescriptorStr, result.GetXId())
	return result, nil
}

// DeleteTenant - Delete a Tenant by the provided ID. This operation will remove the Tenant
// datastore as well as the TenantDescriptor metadata.
func (ash *AdminServiceHandler) DeleteTenant(ctx context.Context, tenantID *wr.StringValue) (*pb.TenantDescriptorResponse, error) {
	// Perform and validation here:
	logger.Log.Infof("Attempting to delete %s: %s", datastore.TenantStr, tenantID.Value)

	// Issue request to DAO Layer to Delete the requested Tenant
	result, err := ash.adminDB.DeleteTenant(tenantID.Value)
	if err != nil {
		return nil, fmt.Errorf("Unable to delete %s: %s", datastore.TenantStr, err.Error())
	}

	// Succesfully removed the Tenant, return the metadata that identified
	// the now deleted tenant.
	logger.Log.Infof("Successfully deleted %s. Previous %s: %s\n", datastore.TenantStr, datastore.TenantDescriptorStr, result.GetXId())
	return result, nil
}

//GetTenantDescriptor - retrieves Tenant metadata for the provided tenantID.
func (ash *AdminServiceHandler) GetTenantDescriptor(ctx context.Context, tenantID *wr.StringValue) (*pb.TenantDescriptorResponse, error) {
	// Perform and validation here:
	logger.Log.Infof("Retrieving %s: %s", datastore.TenantDescriptorStr, tenantID.Value)

	// Issue request to DAO Layer to Get the requested Tenant Metadata
	result, err := ash.adminDB.GetTenantDescriptor(tenantID.Value)
	if err != nil {
		return nil, fmt.Errorf("Unable to retrieve %s: %s", datastore.TenantDescriptorStr, err.Error())
	}

	// Succesfully found the Tenant Metadata, return the result.
	logger.Log.Infof("Retrieved %s: %s\n", datastore.TenantDescriptorStr, result.GetXId())
	return result, nil
}

// GetAllTenantDescriptors -  Retrieve all Tenant Descriptors.
func (ash *AdminServiceHandler) GetAllTenantDescriptors(ctx context.Context, noValue *emp.Empty) (*pb.TenantDescriptorListResponse, error) {
	// Perform and validation here:
	logger.Log.Infof("Retrieving all %ss", datastore.TenantStr)

	// Issue request to DAO Layer to Get the requested Tenant Descriptor List
	result, err := ash.adminDB.GetAllTenantDescriptors()
	if err != nil {
		return nil, fmt.Errorf("Unable to retrieve %ss: %s", datastore.TenantStr, err.Error())
	}

	// Succesfully found the Tenant Descriptors, return the result list.
	logger.Log.Infof("Retrieved %d %ss\n", len(result.GetData()), datastore.TenantStr)
	return result, nil
}

func getAdminServiceDatastore() (datastore.AdminServiceDatastore, error) {
	cfg, err := gather.GetActiveConfig()
	if err != nil {
		return nil, fmt.Errorf("Falied to instantiate AdminServiceHandler: %s", err.Error())
	}

	dbType := cfg.ServerConfig.StartupArgs.AdminDB.Impl
	switch dbType {
	case gather.COUCH:
		logger.Log.Debug("AdminService DB is using CouchDB Implementation")
		return couchDB.CreateAdminServiceDAO()
	case gather.MEM:
		logger.Log.Debug("AdminService DB is using InMemory Implementation")
		return inMemory.CreateAdminServiceDAO()
	}

	return nil, errors.New("No DB implementation provided for Admin Service. Check configuration")
}

func validateAdminUserRequest(request *pb.AdminUserRequest) error {
	if request == nil || request.GetData() == nil {
		return errors.New("Invalid AdminUserRequest: no Admin User data provided")
	}

	if len(request.GetXId()) == 0 {
		return errors.New("Invalid AdminUserRequest: no Admin User ID provided")
	}

	return nil
}

func validateTenantDescriptorRequest(request *pb.TenantDescriptorRequest) error {
	if request == nil || request.GetData() == nil {
		return errors.New("Invalid TenantDescriptorRequest: no Tenant Descriptor data provided")
	}

	if len(request.GetXId()) == 0 {
		return errors.New("Invalid TenantDescriptorRequest: no Tenant ID provided")
	}

	return nil
}
