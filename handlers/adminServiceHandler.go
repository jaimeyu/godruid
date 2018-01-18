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
func (ash *AdminServiceHandler) CreateAdminUser(ctx context.Context, user *pb.AdminUser) (*pb.AdminUser, error) {
	// Validate the request to ensure no invalid data is stored:
	if err := validateAdminUserRequest(user, false); err != nil {
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
func (ash *AdminServiceHandler) UpdateAdminUser(ctx context.Context, user *pb.AdminUser) (*pb.AdminUser, error) {
	// Validate the request to ensure no invalid data is stored:
	if err := validateAdminUserRequest(user, true); err != nil {
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
func (ash *AdminServiceHandler) DeleteAdminUser(ctx context.Context, userID *wr.StringValue) (*pb.AdminUser, error) {
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
func (ash *AdminServiceHandler) GetAdminUser(ctx context.Context, userID *wr.StringValue) (*pb.AdminUser, error) {
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
func (ash *AdminServiceHandler) GetAllAdminUsers(ctx context.Context, noValue *emp.Empty) (*pb.AdminUserList, error) {
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
func (ash *AdminServiceHandler) CreateTenant(ctx context.Context, tenantMeta *pb.TenantDescriptor) (*pb.TenantDescriptor, error) {
	// Validate the request to ensure no invalid data is stored:
	if err := validateTenantDescriptorRequest(tenantMeta, false); err != nil {
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
func (ash *AdminServiceHandler) UpdateTenantDescriptor(ctx context.Context, tenantMeta *pb.TenantDescriptor) (*pb.TenantDescriptor, error) {
	// Validate the request to ensure no invalid data is stored:
	if err := validateTenantDescriptorRequest(tenantMeta, true); err != nil {
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
func (ash *AdminServiceHandler) DeleteTenant(ctx context.Context, tenantID *wr.StringValue) (*pb.TenantDescriptor, error) {
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
func (ash *AdminServiceHandler) GetTenantDescriptor(ctx context.Context, tenantID *wr.StringValue) (*pb.TenantDescriptor, error) {
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
func (ash *AdminServiceHandler) GetAllTenantDescriptors(ctx context.Context, noValue *emp.Empty) (*pb.TenantDescriptorList, error) {
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

// CreateIngestionDictionary - Update an IngestionDictionary used for the entire deployment.
func (ash *AdminServiceHandler) CreateIngestionDictionary(ctx context.Context, ingDictionary *pb.IngestionDictionary) (*pb.IngestionDictionary, error) {
	// Validate the request to ensure no invalid data is stored:
	if err := validateIngestionDictionary(ingDictionary, false); err != nil {
		return nil, err
	}
	logger.Log.Infof("Creating %s: %s", datastore.IngestionDictionaryStr, ingDictionary.GetXId())

	// Issue request to AdminService DAO to create the record:
	result, err := ash.adminDB.CreateIngestionDictionary(ingDictionary)
	if err != nil {
		return nil, fmt.Errorf("Unable to create %s: %s", datastore.IngestionDictionaryStr, err.Error())
	}

	// Succesfully Created the record
	logger.Log.Infof("Created %s: %s\n", datastore.IngestionDictionaryStr, result.GetXId())
	return result, nil
}

// UpdateIngestionDictionary - Update an IngestionDictionary used for the entire deployment.
func (ash *AdminServiceHandler) UpdateIngestionDictionary(ctx context.Context, ingDictionary *pb.IngestionDictionary) (*pb.IngestionDictionary, error) {
	if err := validateIngestionDictionary(ingDictionary, true); err != nil {
		return nil, err
	}
	logger.Log.Infof("Updating %s: %s", datastore.IngestionDictionaryStr, ingDictionary.GetXId())

	// Issue request to AdminService DAO to update the record
	result, err := ash.adminDB.UpdateIngestionDictionary(ingDictionary)
	if err != nil {
		return nil, fmt.Errorf("Unable to update %s: %s", datastore.IngestionDictionaryStr, err.Error())
	}

	// Succesfully updated the record
	logger.Log.Infof("Updated %s: %s\n", datastore.IngestionDictionaryStr, result.GetXId())
	return result, nil
}

// DeleteIngestionDictionary - Delete an IngestionDictionary used for the entire deployment.
func (ash *AdminServiceHandler) DeleteIngestionDictionary(ctx context.Context, noValue *emp.Empty) (*pb.IngestionDictionary, error) {
	logger.Log.Infof("Attempting to delete %s", datastore.IngestionDictionaryStr)

	// Issue request to DAO Layer to Delete the record
	result, err := ash.adminDB.DeleteIngestionDictionary()
	if err != nil {
		return nil, fmt.Errorf("Unable to delete %s: %s", datastore.IngestionDictionaryStr, err.Error())
	}

	// Succesfully removed the record, return the previous record
	logger.Log.Infof("Successfully deleted %s. Previous %s: %s\n", datastore.IngestionDictionaryStr, datastore.IngestionDictionaryStr, result.GetXId())
	return result, nil
}

// GetIngestionDictionary - Retrieve an IngestionDictionary used for the entire deployment.
func (ash *AdminServiceHandler) GetIngestionDictionary(ctx context.Context, noValuie *emp.Empty) (*pb.IngestionDictionary, error) {
	logger.Log.Infof("Retrieving %s", datastore.IngestionDictionaryStr)

	// Issue request to DAO Layer to Get the requested record
	result, err := ash.adminDB.GetIngestionDictionary()
	if err != nil {
		return nil, fmt.Errorf("Unable to retrieve %s: %s", datastore.IngestionDictionaryStr, err.Error())
	}

	// Succesfully found the record, return the result.
	logger.Log.Infof("Retrieved %s: %s\n", datastore.IngestionDictionaryStr, result.GetXId())
	return result, nil
}

// GetTenantIDByAlias - retrieve a Tenant ID by the common name of the Tenant
func (ash *AdminServiceHandler) GetTenantIDByAlias(ctx context.Context, name *wr.StringValue) (*wr.StringValue, error) {

	logger.Log.Infof("Retrieving Tenant ID by name: %s", name)

	// Issue request to DAO Layer to Get the requested Tenant ID
	result, err := ash.adminDB.GetTenantIDByAlias(name.GetValue())
	if err != nil {
		return nil, fmt.Errorf("Unable to retrieve Tenant ID: %s", err.Error())
	}

	// Return Tenant ID
	logger.Log.Infof("Retrieved Tenant ID: %s\n", result)
	return &wr.StringValue{Value: result}, nil
}

// AddAdminViews - Add admin views to Admin DB.
func (ash *AdminServiceHandler) AddAdminViews() error {
	logger.Log.Info("Adding Views to Admin DB")

	// Issue request to DAO Layer to Get the requested Tenant ID
	return ash.adminDB.AddAdminViews()
}

func getAdminServiceDatastore() (datastore.AdminServiceDatastore, error) {
	cfg := gather.GetConfig()
	dbType := gather.DBImpl(cfg.GetInt(gather.CK_args_admindb_impl.String()))
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
