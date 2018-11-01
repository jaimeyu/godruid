package handlers

import (
	"context"
	"errors"
	"fmt"

	"github.com/accedian/adh-gather/datastore"
	"github.com/accedian/adh-gather/gather"
	admmod "github.com/accedian/adh-gather/models/admin"

	"github.com/accedian/adh-gather/datastore/couchDB"

	pb "github.com/accedian/adh-gather/gathergrpc"
	"github.com/accedian/adh-gather/logger"
	emp "github.com/golang/protobuf/ptypes/empty"
	wr "github.com/golang/protobuf/ptypes/wrappers"
)

// AdminServiceHandler - implementation of the interface for the gRPC
// Admin service. Anytime the Admin service changes, the logic to handle the
// API will be modified here.
type AdminServiceHandler struct {
	adminDB datastore.AdminServiceDatastore
}

// CreateAdminServiceHandler - used to generate a handler for the Admin Service.
func CreateAdminServiceHandler() *AdminServiceHandler {
	result := new(AdminServiceHandler)

	// Seteup the DB implementation based on configuration
	db, err := GetAdminServiceDatastore()
	if err != nil {
		logger.Log.Fatalf("Unable to instantiate AdminServiceHandler: %s", err.Error())
	}
	result.adminDB = db

	return result
}

// CreateTenant - Create a Tenant. This will store the identification details for the Tenant,
// TenantDescriptor, as well as generate the Tenant Datastore for the
// Tenant data.
func (ash *AdminServiceHandler) CreateTenant(ctx context.Context, tenantMeta *pb.TenantDescriptor) (*pb.TenantDescriptor, error) {
	// Validate the request to ensure no invalid data is stored:
	if err := validateTenantDescriptorRequest(tenantMeta, false); err != nil {
		msg := fmt.Sprintf("Unable to validate request to create %s: %s", admmod.TenantStr, err.Error())
		logger.Log.Error(msg)
		return nil, fmt.Errorf(msg)
	}
	logger.Log.Infof("Creating %s: %s", admmod.TenantStr, tenantMeta.GetXId())

	// Convert the protobuf object to the proper type:
	converted := admmod.Tenant{}
	if err := pb.ConvertFromPBObject(tenantMeta, &converted); err != nil {
		msg := fmt.Sprintf("Unable to convert request to store %s: %s", admmod.TenantStr, err.Error())
		logger.Log.Error(msg)
		return nil, fmt.Errorf(msg)
	}

	// Issue request to AdminService DAO to create the metadata record:
	result, err := ash.adminDB.CreateTenant(&converted)
	if err != nil {
		msg := fmt.Sprintf("Unable to create %s: %s", admmod.TenantStr, err.Error())
		logger.Log.Error(msg)
		return nil, fmt.Errorf(msg)
	}

	// Convert the result back to PB object
	response := pb.TenantDescriptor{}
	if err := pb.ConvertToPBObject(result, &response); err != nil {
		msg := fmt.Sprintf("Unable to convert response to store %s: %s", admmod.TenantStr, err.Error())
		logger.Log.Error(msg)
		return nil, fmt.Errorf(msg)
	}

	// Succesfully Created the Tenant, return the metadata result.
	logger.Log.Infof("Created %s: %s\n", admmod.TenantStr, response.GetXId())
	return &response, nil
}

// UpdateTenantDescriptor - Update the metadata for a Tenant.
func (ash *AdminServiceHandler) UpdateTenantDescriptor(ctx context.Context, tenantMeta *pb.TenantDescriptor) (*pb.TenantDescriptor, error) {
	// Validate the request to ensure no invalid data is stored:
	// if err := validateTenantDescriptorRequest(tenantMeta, true); err != nil {
	// 	msg := fmt.Sprintf("Unable to validate request to update %s: %s", admmod.TenantStr, err.Error())
	// 	logger.Log.Error(msg)
	// 	return nil, fmt.Errorf(msg)
	// }
	// logger.Log.Infof("Updating %s: %s", admmod.TenantStr, tenantMeta.GetXId())

	// // Issue request to AdminService DAO to update the metadata record:
	// result, err := ash.adminDB.UpdateTenantDescriptor(tenantMeta)
	// if err != nil {
	// 	msg := fmt.Sprintf("Unable to update %s: %s", admmod.TenantStr, err.Error())
	// 	logger.Log.Error(msg)
	// 	return nil, fmt.Errorf(msg)
	// }

	// // Succesfully Createds the Tenant, return the metadata result.
	// logger.Log.Infof("Updated %s: %s\n", admmod.TenantStr, result.GetXId())
	// return result, nil
	return nil, nil
}

// DeleteTenant - Delete a Tenant by the provided ID. This operation will remove the Tenant
// datastore as well as the TenantDescriptor metadata.
func (ash *AdminServiceHandler) DeleteTenant(ctx context.Context, tenantID *wr.StringValue) (*pb.TenantDescriptor, error) {
	// Perform and validation here:
	// logger.Log.Infof("Attempting to delete %s: %s", admmod.TenantStr, tenantID.Value)

	// // Issue request to DAO Layer to Delete the requested Tenant
	// result, err := ash.adminDB.DeleteTenant(tenantID.Value)
	// if err != nil {
	// 	msg := fmt.Sprintf("Unable to delete %s: %s", admmod.TenantStr, err.Error())
	// 	logger.Log.Error(msg)
	// 	return nil, fmt.Errorf(msg)
	// }

	// // Succesfully removed the Tenant, return the metadata that identified
	// // the now deleted tenant.
	// logger.Log.Infof("Successfully deleted %s. Previous %s: %s\n", admmod.TenantStr, admmod.TenantStr, result.GetXId())
	// return result, nil
	return nil, nil
}

//GetTenantDescriptor - retrieves Tenant metadata for the provided tenantID.
func (ash *AdminServiceHandler) GetTenantDescriptor(ctx context.Context, tenantID *wr.StringValue) (*pb.TenantDescriptor, error) {
	// Perform and validation here:
	// logger.Log.Infof("Retrieving %s: %s", admmod.TenantStr, tenantID.Value)

	// // Issue request to DAO Layer to Get the requested Tenant Metadata
	// result, err := ash.adminDB.GetTenantDescriptor(tenantID.Value)
	// if err != nil {
	// 	msg := fmt.Sprintf("Unable to retrieve %s: %s", admmod.TenantStr, err.Error())
	// 	logger.Log.Error(msg)
	// 	return nil, fmt.Errorf(msg)
	// }

	// // Succesfully found the Tenant Metadata, return the result.
	// logger.Log.Infof("Retrieved %s: %s\n", admmod.TenantStr, result.GetXId())
	// return result, nil
	return nil, nil
}

// GetAllTenantDescriptors -  Retrieve all Tenant Descriptors.
func (ash *AdminServiceHandler) GetAllTenantDescriptors(ctx context.Context, noValue *emp.Empty) (*pb.TenantDescriptorList, error) {
	// Perform and validation here:
	// logger.Log.Infof("Retrieving all %ss", admmod.TenantStr)

	// // Issue request to DAO Layer to Get the requested Tenant Descriptor List
	// result, err := ash.adminDB.GetAllTenantDescriptors()
	// if err != nil {
	// 	msg := fmt.Sprintf("Unable to retrieve %ss: %s", admmod.TenantStr, err.Error())
	// 	logger.Log.Error(msg)
	// 	return nil, fmt.Errorf(msg)
	// }

	// // Succesfully found the Tenant Descriptors, return the result list.
	// logger.Log.Infof("Retrieved %d %ss\n", len(result.GetData()), admmod.TenantStr)
	// return result, nil
	return nil, nil
}

// GetTenantIDByAlias - retrieve a Tenant ID by the common name of the Tenant
func (ash *AdminServiceHandler) GetTenantIDByAlias(ctx context.Context, name *wr.StringValue) (*wr.StringValue, error) {

	logger.Log.Infof("Retrieving Tenant ID by name: %s", name)

	// Issue request to DAO Layer to Get the requested Tenant ID
	result, err := ash.adminDB.GetTenantIDByAlias(name.GetValue())
	if err != nil {
		msg := fmt.Sprintf("Unable to retrieve Tenant ID: %s", err.Error())
		logger.Log.Error(msg)
		return nil, fmt.Errorf(msg)
	}

	// Return Tenant ID
	logger.Log.Infof("Retrieved Tenant ID: %s\n", result)
	return &wr.StringValue{Value: result}, nil
}

func GetAdminServiceDatastore() (datastore.AdminServiceDatastore, error) {
	cfg := gather.GetConfig()
	dbType := gather.DBImpl(cfg.GetInt(gather.CK_args_admindb_impl.String()))
	switch dbType {
	case gather.COUCH:
		logger.Log.Debug("AdminService DB is using CouchDB Implementation")
		return couchDB.CreateAdminServiceDAO()
	}

	return nil, errors.New("No DB implementation provided for Admin Service. Check configuration")
}
