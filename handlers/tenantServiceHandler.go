package handlers

import (
	"context"
	"errors"
	"fmt"

	db "github.com/accedian/adh-gather/datastore"
	"github.com/accedian/adh-gather/datastore/couchDB"
	"github.com/accedian/adh-gather/datastore/inMemory"
	"github.com/accedian/adh-gather/gather"
	pb "github.com/accedian/adh-gather/gathergrpc"
	"github.com/accedian/adh-gather/logger"
	wr "github.com/golang/protobuf/ptypes/wrappers"
)

// TenantServiceHandler - implementation of the interface for the gRPC
// Tenant service. Anytime the Tenant service changes, the logic to handle the
// API will be modified here.
type TenantServiceHandler struct {
	tenantDB db.TenantServiceDatastore
}

// CreateTenantServiceHandler - used to generate a handler for the Admin Service.
func CreateTenantServiceHandler() *TenantServiceHandler {
	result := new(TenantServiceHandler)

	// Seteup the DB implementation based on configuration
	db, err := getTenantServiceDatastore()
	if err != nil {
		logger.Log.Fatalf("Unable to instantiate TenantServiceHandler: %s", err.Error())
	}
	result.tenantDB = db

	return result
}

func getTenantServiceDatastore() (db.TenantServiceDatastore, error) {
	cfg := gather.GetConfig()
	dbType := gather.DBImpl(cfg.GetInt(gather.CK_args_tenantdb_impl.String()))
	switch dbType {
	case gather.COUCH:
		logger.Log.Debug("TenantService DB is using CouchDB Implementation")
		return couchDB.CreateTenantServiceDAO()
	case gather.MEM:
		logger.Log.Debug("TenantService DB is using InMemory Implementation")
		return inMemory.CreateTenantServiceDAO()
	}

	return nil, errors.New("No DB implementation provided for Admin Service. Check configuration")
}

// CreateTenantUser - creates a user scoped to a single Tenant.
func (tsh *TenantServiceHandler) CreateTenantUser(ctx context.Context, tenantUserReq *pb.TenantUserRequest) (*pb.TenantUserResponse, error) {
	// Validate the request to ensure no invalid data is stored:
	if err := validateTenantUserRequest(tenantUserReq, false); err != nil {
		return nil, err
	}

	logger.Log.Infof("Creating %s: %s", db.TenantUserStr, tenantUserReq)

	// Issue request to DAO Layer to Create the Tenant User
	result, err := tsh.tenantDB.CreateTenantUser(tenantUserReq)
	if err != nil {
		return nil, fmt.Errorf("Unable to store %s: %s", db.TenantUserStr, err.Error())
	}

	// Succesfully Created the User, return the result.
	logger.Log.Infof("Created %s: %s\n", db.TenantUserStr, result.GetXId())
	return result, nil
}

// UpdateTenantUser - updates a user scoped to a single Tenant.
func (tsh *TenantServiceHandler) UpdateTenantUser(ctx context.Context, tenantUserReq *pb.TenantUserRequest) (*pb.TenantUserResponse, error) {
	// Validate the request to ensure no invalid data is stored:
	if err := validateTenantUserRequest(tenantUserReq, true); err != nil {
		return nil, err
	}

	logger.Log.Infof("Updating %s: %s", db.TenantUserStr, tenantUserReq)

	// Issue request to DAO Layer to Update the Tenant User
	result, err := tsh.tenantDB.UpdateTenantUser(tenantUserReq)
	if err != nil {
		return nil, fmt.Errorf("Unable to store %s: %s", db.TenantUserStr, err.Error())
	}

	// Succesfully Updated the User, return the result.
	logger.Log.Infof("Updated %s: %s\n", db.TenantUserStr, result.GetXId())
	return result, nil
}

// DeleteTenantUser - deletes a user scoped to a single Tenant.
func (tsh *TenantServiceHandler) DeleteTenantUser(ctx context.Context, tenantUserIDReq *pb.TenantUserIdRequest) (*pb.TenantUserResponse, error) {
	// Validate the request to ensure this operation is valid:
	if err := validateTenantUserIDRequest(tenantUserIDReq); err != nil {
		return nil, err
	}

	logger.Log.Infof("Deleting %s: %s", db.TenantUserStr, tenantUserIDReq.GetUserId())

	// Issue request to DAO Layer to Delete the Tenant User
	result, err := tsh.tenantDB.DeleteTenantUser(tenantUserIDReq)
	if err != nil {
		return nil, fmt.Errorf("Unable to delete %s: %s", db.TenantUserStr, err.Error())
	}

	// Succesfully Deleted the User, return the result.
	logger.Log.Infof("Deleted %s: %s\n", db.TenantUserStr, result.GetXId())
	return result, nil
}

// GetTenantUser - retrieves a user scoped to a single Tenant.
func (tsh *TenantServiceHandler) GetTenantUser(ctx context.Context, tenantUserIDReq *pb.TenantUserIdRequest) (*pb.TenantUserResponse, error) {
	// Validate the request to ensure this operatin is valid:
	if err := validateTenantUserIDRequest(tenantUserIDReq); err != nil {
		return nil, err
	}

	logger.Log.Infof("Retrieving %s: %s", db.TenantUserStr, tenantUserIDReq.GetUserId())

	// Issue request to DAO Layer to fetch the Tenant User
	result, err := tsh.tenantDB.GetTenantUser(tenantUserIDReq)
	if err != nil {
		return nil, fmt.Errorf("Unable to retrieve %s: %s", db.TenantUserStr, err.Error())
	}

	// Succesfully fetched the User, return the result.
	logger.Log.Infof("Retrieved %s: %s\n", db.TenantUserStr, result.GetXId())
	return result, nil
}

// GetAllTenantUsers - retrieves all users scoped to a single Tenant.
func (tsh *TenantServiceHandler) GetAllTenantUsers(ctx context.Context, tenantID *wr.StringValue) (*pb.TenantUserListResponse, error) {
	// Validate the request to ensure this operatin is valid:

	logger.Log.Infof("Retrieving all %ss for Tenant: %s", db.TenantUserStr, tenantID.Value)

	// Issue request to DAO Layer to fetch the Tenant Users
	result, err := tsh.tenantDB.GetAllTenantUsers(tenantID.Value)
	if err != nil {
		return nil, fmt.Errorf("Unable to retrieve %ss: %s", db.TenantUserStr, err.Error())
	}

	// Succesfully fetched the Users, return the result.
	logger.Log.Infof("Retrieved %d %ss:\n", len(result.GetData()), db.TenantUserStr)
	return result, nil
}

// CreateTenantDomain - creates a Domain scoped to a single Tenant.
func (tsh *TenantServiceHandler) CreateTenantDomain(ctx context.Context, tenantDomainRequest *pb.TenantDomainRequest) (*pb.TenantDomainResponse, error) {
	// Validate the request to ensure no invalid data is stored:
	if err := validateTenantDomainRequest(tenantDomainRequest, false); err != nil {
		return nil, err
	}

	logger.Log.Infof("Creating %s: %s", db.TenantDomainStr, tenantDomainRequest)

	// Issue request to DAO Layer to Create the Tenant Domain
	result, err := tsh.tenantDB.CreateTenantDomain(tenantDomainRequest)
	if err != nil {
		return nil, fmt.Errorf("Unable to store %s: %s", db.TenantDomainStr, err.Error())
	}

	// Succesfully Created the Domain, return the result.
	logger.Log.Infof("Created %s: %s\n", db.TenantDomainStr, result.GetXId())
	return result, nil
}

// UpdateTenantDomain - updates a Domain scoped to a single Tenant.
func (tsh *TenantServiceHandler) UpdateTenantDomain(ctx context.Context, tenantDomainRequest *pb.TenantDomainRequest) (*pb.TenantDomainResponse, error) {
	// Validate the request to ensure no invalid data is stored:
	if err := validateTenantDomainRequest(tenantDomainRequest, true); err != nil {
		return nil, err
	}

	logger.Log.Infof("Updating %s: %s", db.TenantDomainStr, tenantDomainRequest)

	// Issue request to DAO Layer to Update the Tenant Domain
	result, err := tsh.tenantDB.UpdateTenantDomain(tenantDomainRequest)
	if err != nil {
		return nil, fmt.Errorf("Unable to store %s: %s", db.TenantDomainStr, err.Error())
	}

	// Succesfully Updated the Domain, return the result.
	logger.Log.Infof("Updated %s: %s\n", db.TenantDomainStr, result.GetXId())
	return result, nil
}

// DeleteTenantDomain - deletes a Domain scoped to a single Tenant.
func (tsh *TenantServiceHandler) DeleteTenantDomain(ctx context.Context, tenantDomainIDRequest *pb.TenantDomainIdRequest) (*pb.TenantDomainResponse, error) {
	// Validate the request to ensure this operation is valid:
	if err := validateTenantDomainIDRequest(tenantDomainIDRequest); err != nil {
		return nil, err
	}

	logger.Log.Infof("Deleting %s: %s", db.TenantDomainStr, tenantDomainIDRequest.GetDomainId())

	// Issue request to DAO Layer to Delete the Tenant Domain
	result, err := tsh.tenantDB.DeleteTenantDomain(tenantDomainIDRequest)
	if err != nil {
		return nil, fmt.Errorf("Unable to delete %s: %s", db.TenantDomainStr, err.Error())
	}

	// Succesfully Deleted the Domain, return the result.
	logger.Log.Infof("Deleted %s: %s\n", db.TenantDomainStr, result.GetXId())
	return result, nil
}

// GetTenantDomain - retrieves a Domain scoped to a single Tenant.
func (tsh *TenantServiceHandler) GetTenantDomain(ctx context.Context, tenantDomainIDRequest *pb.TenantDomainIdRequest) (*pb.TenantDomainResponse, error) {
	// Validate the request to ensure this operatin is valid:
	if err := validateTenantDomainIDRequest(tenantDomainIDRequest); err != nil {
		return nil, err
	}

	logger.Log.Infof("Retrieving %s: %s", db.TenantDomainStr, tenantDomainIDRequest.GetDomainId())

	// Issue request to DAO Layer to fetch the Tenant Domain
	result, err := tsh.tenantDB.GetTenantDomain(tenantDomainIDRequest)
	if err != nil {
		return nil, fmt.Errorf("Unable to retrieve %s: %s", db.TenantDomainStr, err.Error())
	}

	// Succesfully fetched the Domain, return the result.
	logger.Log.Infof("Retrieved %s: %s\n", db.TenantDomainStr, result.GetXId())
	return result, nil
}

// GetAllTenantDomains - retrieves all Domains scoped to a single Tenant.
func (tsh *TenantServiceHandler) GetAllTenantDomains(ctx context.Context, tenantID *wr.StringValue) (*pb.TenantDomainListResponse, error) {
	// Validate the request to ensure this operation is valid:

	logger.Log.Infof("Retrieving all %ss for Tenant: %s", db.TenantDomainStr, tenantID.Value)

	// Issue request to DAO Layer to fetch the Tenant Domains
	result, err := tsh.tenantDB.GetAllTenantDomains(tenantID.Value)
	if err != nil {
		return nil, fmt.Errorf("Unable to retrieve %ss: %s", db.TenantDomainStr, err.Error())
	}

	// Succesfully fetched the Domains, return the result.
	logger.Log.Infof("Retrieved %d %ss:\n", len(result.GetData()), db.TenantDomainStr)
	return result, nil
}

// CreateTenantIngestionProfile - creates an Ingestion Profile scoped to a specific Tenant.
func (tsh *TenantServiceHandler) CreateTenantIngestionProfile(ctx context.Context, tenantIngPrfReq *pb.TenantIngestionProfileRequest) (*pb.TenantIngestionProfileResponse, error) {
	// Validate the request to ensure no invalid data is stored:
	if err := validateTenantIngPrfRequest(tenantIngPrfReq, false); err != nil {
		return nil, err
	}

	logger.Log.Infof("Creating %s: %s", db.TenantIngestionProfileStr, tenantIngPrfReq)

	// Issue request to DAO Layer to Create the Tenant Ingestion Profile
	result, err := tsh.tenantDB.CreateTenantIngestionProfile(tenantIngPrfReq)
	if err != nil {
		return nil, fmt.Errorf("Unable to store %s: %s", db.TenantIngestionProfileStr, err.Error())
	}

	// Succesfully Created the Ingestion Profile, return the result.
	logger.Log.Infof("Created %s: %s\n", db.TenantIngestionProfileStr, result.GetXId())
	return result, nil
}

// UpdateTenantIngestionProfile - updates an Ingestion Profile scoped to a specific Tenant.
func (tsh *TenantServiceHandler) UpdateTenantIngestionProfile(ctx context.Context, tenantIngPrfReq *pb.TenantIngestionProfileRequest) (*pb.TenantIngestionProfileResponse, error) {
	// Validate the request to ensure no invalid data is stored:
	if err := validateTenantIngPrfRequest(tenantIngPrfReq, true); err != nil {
		return nil, err
	}

	logger.Log.Infof("Updating %s: %s", db.TenantIngestionProfileStr, tenantIngPrfReq)

	// Issue request to DAO Layer to Update the Tenant Ingestion Profile
	result, err := tsh.tenantDB.UpdateTenantIngestionProfile(tenantIngPrfReq)
	if err != nil {
		return nil, fmt.Errorf("Unable to store %s: %s", db.TenantIngestionProfileStr, err.Error())
	}

	// Succesfully Updated the Ingestion Profile, return the result.
	logger.Log.Infof("Updated %s: %s\n", db.TenantIngestionProfileStr, result.GetXId())
	return result, nil
}

// GetTenantIngestionProfile - retrieves the Ingestion Profile for a single Tenant.
func (tsh *TenantServiceHandler) GetTenantIngestionProfile(ctx context.Context, tenantIngPrfIDReq *pb.TenantIngestionProfileIdRequest) (*pb.TenantIngestionProfileResponse, error) {
	// Validate the request to ensure the operation is valid:
	if err := validateTenantIngPrfIDRequest(tenantIngPrfIDReq); err != nil {
		return nil, err
	}

	logger.Log.Infof("Retrieving %s for Tenant %s", db.TenantIngestionProfileStr, tenantIngPrfIDReq.GetTenantId())

	// Issue request to DAO Layer to fetch the Tenant Ingestion Profile
	result, err := tsh.tenantDB.GetTenantIngestionProfile(tenantIngPrfIDReq)
	if err != nil {
		return nil, fmt.Errorf("Unable to retrieve %s: %s", db.TenantIngestionProfileStr, err.Error())
	}

	// Succesfully fetched the Ingestion Profile, return the result.
	logger.Log.Infof("Retrieved %s: %s\n", db.TenantIngestionProfileStr, result.GetXId())
	return result, nil
}

// DeleteTenantIngestionProfile - deletes the Ingestion Profile for a single Tenant.
func (tsh *TenantServiceHandler) DeleteTenantIngestionProfile(ctx context.Context, tenantIngPrfIDReq *pb.TenantIngestionProfileIdRequest) (*pb.TenantIngestionProfileResponse, error) {
	// Validate the request to ensure the operation is valid:
	if err := validateTenantIngPrfIDRequest(tenantIngPrfIDReq); err != nil {
		return nil, err
	}

	logger.Log.Infof("Deleting %s for Tenant %s", db.TenantIngestionProfileStr, tenantIngPrfIDReq.GetTenantId())

	// Issue request to DAO Layer to delete the Tenant Ingestion Profile
	result, err := tsh.tenantDB.DeleteTenantIngestionProfile(tenantIngPrfIDReq)
	if err != nil {
		return nil, fmt.Errorf("Unable to delete %s: %s", db.TenantIngestionProfileStr, err.Error())
	}

	// Succesfully deleted the Ingestion Profile, return the result.
	logger.Log.Infof("Deleted %s: %s\n", db.TenantIngestionProfileStr, result.GetXId())
	return result, nil
}

// CreateTenantThresholdProfile - creates an Threshold Profile scoped to a specific Tenant.
func (tsh *TenantServiceHandler) CreateTenantThresholdProfile(ctx context.Context, tenantThreshPrfReq *pb.TenantThresholdProfileRequest) (*pb.TenantThresholdProfileResponse, error) {
	// Validate the request to ensure no invalid data is stored:
	if err := validateTenantThreshPrfRequest(tenantThreshPrfReq, false); err != nil {
		return nil, err
	}

	logger.Log.Infof("Creating %s: %s", db.TenantThresholdProfileStr, tenantThreshPrfReq)

	// Issue request to DAO Layer to Create the Tenant Threshold Profile
	result, err := tsh.tenantDB.CreateTenantThresholdProfile(tenantThreshPrfReq)
	if err != nil {
		return nil, fmt.Errorf("Unable to store %s: %s", db.TenantThresholdProfileStr, err.Error())
	}

	// Succesfully Created the Threshold Profile, return the result.
	logger.Log.Infof("Created %s: %s\n", db.TenantThresholdProfileStr, result.GetXId())
	return result, nil
}

// UpdateTenantThresholdProfile - updates an Threshold Profile scoped to a specific Tenant.
func (tsh *TenantServiceHandler) UpdateTenantThresholdProfile(ctx context.Context, tenantThreshPrfReq *pb.TenantThresholdProfileRequest) (*pb.TenantThresholdProfileResponse, error) {
	// Validate the request to ensure no invalid data is stored:
	if err := validateTenantThreshPrfRequest(tenantThreshPrfReq, true); err != nil {
		return nil, err
	}

	logger.Log.Infof("Updating %s: %s", db.TenantThresholdProfileStr, tenantThreshPrfReq)

	// Issue request to DAO Layer to Update the Tenant Threshold Profile
	result, err := tsh.tenantDB.UpdateTenantThresholdProfile(tenantThreshPrfReq)
	if err != nil {
		return nil, fmt.Errorf("Unable to store %s: %s", db.TenantThresholdProfileStr, err.Error())
	}

	// Succesfully Updated the Threshold Profile, return the result.
	logger.Log.Infof("Updated %s: %s\n", db.TenantThresholdProfileStr, result.GetXId())
	return result, nil
}

// GetTenantThresholdProfile - retrieves the Threshold Profile for a single Tenant.
func (tsh *TenantServiceHandler) GetTenantThresholdProfile(ctx context.Context, tenantThreshPrfIDReq *pb.TenantThresholdProfileIdRequest) (*pb.TenantThresholdProfileResponse, error) {
	// Validate the request to ensure the operation is valid:
	if err := validateTenantThreshPrfIDRequest(tenantThreshPrfIDReq); err != nil {
		return nil, err
	}

	logger.Log.Infof("Retrieving %s for Tenant %s", db.TenantThresholdProfileStr, tenantThreshPrfIDReq.GetTenantId())

	// Issue request to DAO Layer to fetch the Tenant Threshold Profile
	result, err := tsh.tenantDB.GetTenantThresholdProfile(tenantThreshPrfIDReq)
	if err != nil {
		return nil, fmt.Errorf("Unable to retrieve %s: %s", db.TenantThresholdProfileStr, err.Error())
	}

	// Succesfully fetched the Threshold Profile, return the result.
	logger.Log.Infof("Retrieved %s: %s\n", db.TenantThresholdProfileStr, result.GetXId())
	return result, nil
}

// DeleteTenantThresholdProfile - deletes the Threshold Profile for a single Tenant.
func (tsh *TenantServiceHandler) DeleteTenantThresholdProfile(ctx context.Context, tenantThreshPrfIDReq *pb.TenantThresholdProfileIdRequest) (*pb.TenantThresholdProfileResponse, error) {
	// Validate the request to ensure the operation is valid:
	if err := validateTenantThreshPrfIDRequest(tenantThreshPrfIDReq); err != nil {
		return nil, err
	}

	logger.Log.Infof("Deleting %s for Tenant %s", db.TenantThresholdProfileStr, tenantThreshPrfIDReq.GetTenantId())

	// Issue request to DAO Layer to delete the Tenant Threshold Profile
	result, err := tsh.tenantDB.DeleteTenantThresholdProfile(tenantThreshPrfIDReq)
	if err != nil {
		return nil, fmt.Errorf("Unable to delete %s: %s", db.TenantThresholdProfileStr, err.Error())
	}

	// Succesfully deleted the Threshold Profile, return the result.
	logger.Log.Infof("Deleted %s: %s\n", db.TenantThresholdProfileStr, result.GetXId())
	return result, nil
}

// CreateMonitoredObject - creates a Monitored Object scoped to a specific tenant
func (tsh *TenantServiceHandler) CreateMonitoredObject(ctx context.Context, monitoredObjectReq *pb.MonitoredObjectRequest) (*pb.MonitoredObjectResponse, error) {
	// Validate the request to ensure no invalid data is stored:
	if err := validateMonitoredObjectRequest(monitoredObjectReq, false); err != nil {
		return nil, err
	}

	logger.Log.Infof("Creating %s: %s", db.TenantMonitoredObjectStr, monitoredObjectReq)

	// If no id is provided for the Manged Object, generate one
	objectID := monitoredObjectReq.GetXId()
	if len(objectID) == 0 {
		objectID := db.GenerateID(monitoredObjectReq.GetData(), string(db.TenantMonitoredObjectType))
		monitoredObjectReq.XId = objectID
	}

	// Issue request to DAO Layer to Create the Tenant Monitored Object
	result, err := tsh.tenantDB.CreateMonitoredObject(monitoredObjectReq)
	if err != nil {
		return nil, fmt.Errorf("Unable to store %s: %s", db.TenantMonitoredObjectStr, err.Error())
	}

	// Succesfully Created the Monitored, return the result.
	logger.Log.Infof("Created %s: %s\n", db.TenantMonitoredObjectStr, result.GetXId())
	return result, nil
}

// UpdateMonitoredObject - updates an MonitoredObject scoped to a specific Tenant.
func (tsh *TenantServiceHandler) UpdateMonitoredObject(ctx context.Context, monitoredObjectReq *pb.MonitoredObjectRequest) (*pb.MonitoredObjectResponse, error) {
	// Validate the request to ensure no invalid data is stored:
	if err := validateMonitoredObjectRequest(monitoredObjectReq, true); err != nil {
		return nil, err
	}

	logger.Log.Infof("Updating %s: %s", db.TenantMonitoredObjectStr, monitoredObjectReq)

	// If no id is provided for the Manged Object, generate one
	objectID := monitoredObjectReq.GetXId()
	if len(objectID) == 0 {
		objectID := db.GenerateID(monitoredObjectReq.GetData(), string(db.TenantMonitoredObjectType))
		monitoredObjectReq.XId = objectID
	}

	// Issue request to DAO Layer to Update the Tenant Monitored Object
	result, err := tsh.tenantDB.UpdateMonitoredObject(monitoredObjectReq)
	if err != nil {
		return nil, fmt.Errorf("Unable to store %s: %s", db.TenantMonitoredObjectStr, err.Error())
	}

	// Succesfully Updated the Monitored Object, return the result.
	logger.Log.Infof("Updated %s: %s\n", db.TenantMonitoredObjectStr, result.GetXId())
	return result, nil
}

// GetMonitoredObject - retrieves the MonitoredObject for a singler Tenant.
func (tsh *TenantServiceHandler) GetMonitoredObject(ctx context.Context, monitoredObjectIDReq *pb.MonitoredObjectIdRequest) (*pb.MonitoredObjectResponse, error) {
	// Validate the request to ensure no invalid data is stored:
	if err := validateMonitoredObjectIDRequest(monitoredObjectIDReq); err != nil {
		return nil, err
	}

	// Issue request to DAO Layer to fetch the Tenant Monitored Object
	result, err := tsh.tenantDB.GetMonitoredObject(monitoredObjectIDReq)
	if err != nil {
		return nil, fmt.Errorf("Unable to retrieve %s: %s", db.TenantMonitoredObjectStr, err.Error())
	}

	// Succesfully fetched the Monitored Object, return the result.
	logger.Log.Infof("Retrieved %s: %s\n", db.TenantMonitoredObjectStr, result.GetXId())
	return result, nil
}

// DeleteMonitoredObject - deletes the MonitoredObject for a singler Tenant.
func (tsh *TenantServiceHandler) DeleteMonitoredObject(ctx context.Context, monitoredObjectIDReq *pb.MonitoredObjectIdRequest) (*pb.MonitoredObjectResponse, error) {
	// Validate the request to ensure the operation is valid:
	if err := validateMonitoredObjectIDRequest(monitoredObjectIDReq); err != nil {
		return nil, err
	}

	logger.Log.Infof("Deleting %s for Tenant %s", db.TenantMonitoredObjectStr, monitoredObjectIDReq.GetTenantId())

	// Issue request to DAO Layer to delete the Tenant Monitored Object
	result, err := tsh.tenantDB.DeleteMonitoredObject(monitoredObjectIDReq)
	if err != nil {
		return nil, fmt.Errorf("Unable to delete %s: %s", db.TenantMonitoredObjectStr, err.Error())
	}

	// Succesfully deleted the MonitoredObject, return the result.
	logger.Log.Infof("Deleted %s: %s\n", db.TenantMonitoredObjectStr, result.GetXId())
	return result, nil
}

// GetAllMonitoredObjects - retrieves all MonitoredObjects scoped to a single Tenant.
func (tsh *TenantServiceHandler) GetAllMonitoredObjects(ctx context.Context, tenantID *wr.StringValue) (*pb.MonitoredObjectListResponse, error) {
	// Validate the request to ensure this operation is valid:

	logger.Log.Infof("Retrieving all %ss for Tenant: %s", db.TenantMonitoredObjectStr, tenantID.Value)

	// Issue request to DAO Layer to fetch the Tenant Monitored Objects
	result, err := tsh.tenantDB.GetAllMonitoredObjects(tenantID.Value)
	if err != nil {
		return nil, fmt.Errorf("Unable to retrieve %ss: %s", db.TenantMonitoredObjectStr, err.Error())
	}

	// Succesfully fetched the Monitored Objects, return the result.
	logger.Log.Infof("Retrieved %d %ss:\n", len(result.GetData()), db.TenantMonitoredObjectStr)
	return result, nil
}

// GetMonitoredObjectToDomainMap - retrieves a mapping of MonitoredObjects to each Domain. Will retrieve the mapping either as a count, or as a set of all
// MonitoredObjects that use each Domain.
func (tsh *TenantServiceHandler) GetMonitoredObjectToDomainMap(ctx context.Context, moByDomReq *pb.MonitoredObjectCountByDomainRequest) (*pb.MonitoredObjectCountByDomainResponse, error) {
	// Validate the request:
	if err := validateMonitoredObjectToDomainMapRequest(moByDomReq); err != nil {
		return nil, err
	}

	// Issue request to DAO Layer to fetch the Tenant Monitored Object Map
	result, err := tsh.tenantDB.GetMonitoredObjectToDomainMap(moByDomReq)
	if err != nil {
		return nil, fmt.Errorf("Unable to retrieve %s: %s", db.MonitoredObjectToDomainMapStr, err.Error())
	}

	// Succesfully fetched the Monitored Object Map, return the result.
	logger.Log.Infof("Successfully retrieved %s: %s\n", db.MonitoredObjectToDomainMapStr)
	return result, nil
}

// CreateTenantMeta - Create TenantMeta scoped to a Single Tenant.
func (tsh *TenantServiceHandler) CreateTenantMeta(ctx context.Context, meta *pb.TenantMeta) (*pb.TenantMeta, error) {
	// Validate the request to ensure no invalid data is stored:
	if err := validateTenantMetaRequest(meta, false); err != nil {
		return nil, err
	}

	logger.Log.Infof("Creating %s for Tenant %s", db.TenantMetaStr, meta.GetData().GetTenantId())

	// Issue request to DAO Layer to Create the record
	result, err := tsh.tenantDB.CreateTenantMeta(meta)
	if err != nil {
		return nil, fmt.Errorf("Unable to store %s: %s", db.TenantMetaStr, err.Error())
	}

	// Succesfully Created the record, return the result.
	logger.Log.Infof("Created %s: %s\n", db.TenantMetaStr, result.GetXId())
	return result, nil
}

// UpdateTenantMeta - Update TenantMeta scoped to a single Tenant.
func (tsh *TenantServiceHandler) UpdateTenantMeta(ctx context.Context, meta *pb.TenantMeta) (*pb.TenantMeta, error) {
	// Validate the request to ensure no invalid data is stored:
	if err := validateTenantMetaRequest(meta, true); err != nil {
		return nil, err
	}

	logger.Log.Infof("Updating %s: %s", db.TenantMetaStr, meta)

	// Issue request to DAO Layer to Update the record
	result, err := tsh.tenantDB.UpdateTenantMeta(meta)
	if err != nil {
		return nil, fmt.Errorf("Unable to store %s: %s", db.TenantMetaStr, err.Error())
	}

	// Succesfully Updated the record, return the result.
	logger.Log.Infof("Updated %s: %s\n", db.TenantMetaStr, result.GetXId())
	return result, nil
}

// DeleteTenantMeta - Delete TenantMeta scoped to a single Tenant.
func (tsh *TenantServiceHandler) DeleteTenantMeta(ctx context.Context, tenantID *wr.StringValue) (*pb.TenantMeta, error) {

	logger.Log.Infof("Deleting %s for Tenant %s", db.TenantMetaStr, tenantID.GetValue())

	// Issue request to DAO Layer to delete the record
	result, err := tsh.tenantDB.DeleteTenantMeta(tenantID.GetValue())
	if err != nil {
		return nil, fmt.Errorf("Unable to delete %s: %s", db.TenantMetaStr, err.Error())
	}

	// Succesfully deleted the record, return the result.
	logger.Log.Infof("Deleted %s: %s\n", db.TenantMetaStr, result.GetXId())
	return result, nil
}

// GetTenantMeta - Retrieve a User scoped to a single Tenant.
func (tsh *TenantServiceHandler) GetTenantMeta(ctx context.Context, tenantID *wr.StringValue) (*pb.TenantMeta, error) {

	// Issue request to DAO Layer to fetch the record
	result, err := tsh.tenantDB.GetTenantMeta(tenantID.GetValue())
	if err != nil {
		return nil, fmt.Errorf("Unable to retrieve %s: %s", db.TenantMetaStr, err.Error())
	}

	// Succesfully fetched the record, return the result.
	logger.Log.Infof("Retrieved %s: %s\n", db.TenantMetaStr, result.GetXId())
	return result, nil
}
