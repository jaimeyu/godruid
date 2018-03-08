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

	tenmod "github.com/accedian/adh-gather/models/tenant"
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
func (tsh *TenantServiceHandler) CreateTenantUser(ctx context.Context, tenantUserReq *pb.TenantUser) (*pb.TenantUser, error) {
	// Validate the request to ensure no invalid data is stored:
	if err := validateTenantUserRequest(tenantUserReq, false); err != nil {
		msg := fmt.Sprintf("Unable to validate request to store %s: %s", tenmod.TenantUserStr, err.Error())
		logger.Log.Error(msg)
		return nil, fmt.Errorf(msg)
	}

	logger.Log.Infof("Creating %s: %s", tenmod.TenantUserStr, tenantUserReq)

	// Issue request to DAO Layer to Create the Tenant User
	result, err := tsh.tenantDB.CreateTenantUser(tenantUserReq)
	if err != nil {
		msg := fmt.Sprintf("Unable to store %s: %s", tenmod.TenantUserStr, err.Error())
		logger.Log.Error(msg)
		return nil, fmt.Errorf(msg)
	}

	// Succesfully Created the User, return the result.
	logger.Log.Infof("Created %s: %s\n", tenmod.TenantUserStr, result.GetXId())
	return result, nil
}

// UpdateTenantUser - updates a user scoped to a single Tenant.
func (tsh *TenantServiceHandler) UpdateTenantUser(ctx context.Context, tenantUserReq *pb.TenantUser) (*pb.TenantUser, error) {
	// Validate the request to ensure no invalid data is stored:
	if err := validateTenantUserRequest(tenantUserReq, true); err != nil {
		msg := fmt.Sprintf("Unable to validate request to store %s: %s", tenmod.TenantUserStr, err.Error())
		logger.Log.Error(msg)
		return nil, fmt.Errorf(msg)
	}

	logger.Log.Infof("Updating %s: %s", tenmod.TenantUserStr, tenantUserReq)

	// Issue request to DAO Layer to Update the Tenant User
	result, err := tsh.tenantDB.UpdateTenantUser(tenantUserReq)
	if err != nil {
		msg := fmt.Sprintf("Unable to store %s: %s", tenmod.TenantUserStr, err.Error())
		logger.Log.Error(msg)
		return nil, fmt.Errorf(msg)
	}

	// Succesfully Updated the User, return the result.
	logger.Log.Infof("Updated %s: %s\n", tenmod.TenantUserStr, result.GetXId())
	return result, nil
}

// DeleteTenantUser - deletes a user scoped to a single Tenant.
func (tsh *TenantServiceHandler) DeleteTenantUser(ctx context.Context, tenantUserIDReq *pb.TenantUserIdRequest) (*pb.TenantUser, error) {
	// Validate the request to ensure this operation is valid:
	if err := validateTenantUserIDRequest(tenantUserIDReq); err != nil {
		msg := fmt.Sprintf("Unable to validate request to delete %s: %s", tenmod.TenantUserStr, err.Error())
		logger.Log.Error(msg)
		return nil, fmt.Errorf(msg)
	}

	logger.Log.Infof("Deleting %s: %s", tenmod.TenantUserStr, tenantUserIDReq.GetUserId())

	// Issue request to DAO Layer to Delete the Tenant User
	result, err := tsh.tenantDB.DeleteTenantUser(tenantUserIDReq)
	if err != nil {
		msg := fmt.Sprintf("Unable to delete %s: %s", tenmod.TenantUserStr, err.Error())
		logger.Log.Error(msg)
		return nil, fmt.Errorf(msg)
	}

	// Succesfully Deleted the User, return the result.
	logger.Log.Infof("Deleted %s: %s\n", tenmod.TenantUserStr, result.GetXId())
	return result, nil
}

// GetTenantUser - retrieves a user scoped to a single Tenant.
func (tsh *TenantServiceHandler) GetTenantUser(ctx context.Context, tenantUserIDReq *pb.TenantUserIdRequest) (*pb.TenantUser, error) {
	// Validate the request to ensure this operatin is valid:
	if err := validateTenantUserIDRequest(tenantUserIDReq); err != nil {
		msg := fmt.Sprintf("Unable to validate request to fetch %s: %s", tenmod.TenantUserStr, err.Error())
		logger.Log.Error(msg)
		return nil, fmt.Errorf(msg)
	}

	logger.Log.Infof("Retrieving %s: %s", tenmod.TenantUserStr, tenantUserIDReq.GetUserId())

	// Issue request to DAO Layer to fetch the Tenant User
	result, err := tsh.tenantDB.GetTenantUser(tenantUserIDReq)
	if err != nil {
		msg := fmt.Sprintf("Unable to fetch %s: %s", tenmod.TenantUserStr, err.Error())
		logger.Log.Error(msg)
		return nil, fmt.Errorf(msg)
	}

	// Succesfully fetched the User, return the result.
	logger.Log.Infof("Retrieved %s: %s\n", tenmod.TenantUserStr, result.GetXId())
	return result, nil
}

// GetAllTenantUsers - retrieves all users scoped to a single Tenant.
func (tsh *TenantServiceHandler) GetAllTenantUsers(ctx context.Context, tenantID *wr.StringValue) (*pb.TenantUserList, error) {
	// Validate the request to ensure this operatin is valid:

	logger.Log.Infof("Retrieving all %ss for Tenant: %s", tenmod.TenantUserStr, tenantID.Value)

	// Issue request to DAO Layer to fetch the Tenant Users
	result, err := tsh.tenantDB.GetAllTenantUsers(tenantID.Value)
	if err != nil {
		msg := fmt.Sprintf("Unable to retrieve %ss: %s", tenmod.TenantUserStr, err.Error())
		logger.Log.Error(msg)
		return nil, fmt.Errorf(msg)
	}

	// Succesfully fetched the Users, return the result.
	logger.Log.Infof("Retrieved %d %ss:\n", len(result.GetData()), tenmod.TenantUserStr)
	return result, nil
}

// CreateTenantDomain - creates a Domain scoped to a single Tenant.
func (tsh *TenantServiceHandler) CreateTenantDomain(ctx context.Context, tenantDomainRequest *pb.TenantDomain) (*pb.TenantDomain, error) {
	// Validate the request to ensure no invalid data is stored:
	if err := validateTenantDomainRequest(tenantDomainRequest, false); err != nil {
		msg := fmt.Sprintf("Unable to validate request to store %s: %s", tenmod.TenantDomainStr, err.Error())
		logger.Log.Error(msg)
		return nil, fmt.Errorf(msg)
	}

	logger.Log.Infof("Creating %s: %s", tenmod.TenantDomainStr, tenantDomainRequest)

	// Issue request to DAO Layer to Create the Tenant Domain
	result, err := tsh.tenantDB.CreateTenantDomain(tenantDomainRequest)
	if err != nil {
		msg := fmt.Sprintf("Unable to store %s: %s", tenmod.TenantDomainStr, err.Error())
		logger.Log.Error(msg)
		return nil, fmt.Errorf(msg)
	}

	// Succesfully Created the Domain, return the result.
	logger.Log.Infof("Created %s: %s\n", tenmod.TenantDomainStr, result.GetXId())
	return result, nil
}

// UpdateTenantDomain - updates a Domain scoped to a single Tenant.
func (tsh *TenantServiceHandler) UpdateTenantDomain(ctx context.Context, tenantDomainRequest *pb.TenantDomain) (*pb.TenantDomain, error) {
	// Validate the request to ensure no invalid data is stored:
	if err := validateTenantDomainRequest(tenantDomainRequest, true); err != nil {
		msg := fmt.Sprintf("Unable to validate request to store %s: %s", tenmod.TenantDomainStr, err.Error())
		logger.Log.Error(msg)
		return nil, fmt.Errorf(msg)
	}

	logger.Log.Infof("Updating %s: %s", tenmod.TenantDomainStr, tenantDomainRequest)

	// Issue request to DAO Layer to Update the Tenant Domain
	result, err := tsh.tenantDB.UpdateTenantDomain(tenantDomainRequest)
	if err != nil {
		msg := fmt.Sprintf("Unable to store %s: %s", tenmod.TenantDomainStr, err.Error())
		logger.Log.Error(msg)
		return nil, fmt.Errorf(msg)
	}

	// Succesfully Updated the Domain, return the result.
	logger.Log.Infof("Updated %s: %s\n", tenmod.TenantDomainStr, result.GetXId())
	return result, nil
}

// DeleteTenantDomain - deletes a Domain scoped to a single Tenant.
func (tsh *TenantServiceHandler) DeleteTenantDomain(ctx context.Context, tenantDomainIDRequest *pb.TenantDomainIdRequest) (*pb.TenantDomain, error) {
	// Validate the request to ensure this operation is valid:
	if err := validateTenantDomainIDRequest(tenantDomainIDRequest); err != nil {
		msg := fmt.Sprintf("Unable to validate request to delete %s: %s", tenmod.TenantDomainStr, err.Error())
		logger.Log.Error(msg)
		return nil, fmt.Errorf(msg)
	}

	logger.Log.Infof("Deleting %s: %s", tenmod.TenantDomainStr, tenantDomainIDRequest.GetDomainId())

	// Issue request to DAO Layer to Delete the Tenant Domain
	result, err := tsh.tenantDB.DeleteTenantDomain(tenantDomainIDRequest)
	if err != nil {
		msg := fmt.Sprintf("Unable to delete %s: %s", tenmod.TenantDomainStr, err.Error())
		logger.Log.Error(msg)
		return nil, fmt.Errorf(msg)
	}

	// Succesfully Deleted the Domain, return the result.
	logger.Log.Infof("Deleted %s: %s\n", tenmod.TenantDomainStr, result.GetXId())
	return result, nil
}

// GetTenantDomain - retrieves a Domain scoped to a single Tenant.
func (tsh *TenantServiceHandler) GetTenantDomain(ctx context.Context, tenantDomainIDRequest *pb.TenantDomainIdRequest) (*pb.TenantDomain, error) {
	// Validate the request to ensure this operatin is valid:
	if err := validateTenantDomainIDRequest(tenantDomainIDRequest); err != nil {
		msg := fmt.Sprintf("Unable to validate request to fetch %s: %s", tenmod.TenantDomainStr, err.Error())
		logger.Log.Error(msg)
		return nil, fmt.Errorf(msg)
	}

	logger.Log.Infof("Retrieving %s: %s", tenmod.TenantDomainStr, tenantDomainIDRequest.GetDomainId())

	// Issue request to DAO Layer to fetch the Tenant Domain
	result, err := tsh.tenantDB.GetTenantDomain(tenantDomainIDRequest)
	if err != nil {
		msg := fmt.Sprintf("Unable to fetch %s: %s", tenmod.TenantDomainStr, err.Error())
		logger.Log.Error(msg)
		return nil, fmt.Errorf(msg)
	}

	// Succesfully fetched the Domain, return the result.
	logger.Log.Infof("Retrieved %s: %s\n", tenmod.TenantDomainStr, result.GetXId())
	return result, nil
}

// GetAllTenantDomains - retrieves all Domains scoped to a single Tenant.
func (tsh *TenantServiceHandler) GetAllTenantDomains(ctx context.Context, tenantID *wr.StringValue) (*pb.TenantDomainList, error) {
	// Validate the request to ensure this operation is valid:

	logger.Log.Infof("Retrieving all %ss for Tenant: %s", tenmod.TenantDomainStr, tenantID.Value)

	// Issue request to DAO Layer to fetch the Tenant Domains
	result, err := tsh.tenantDB.GetAllTenantDomains(tenantID.Value)
	if err != nil {
		msg := fmt.Sprintf("Unable to retrieve %ss: %s", tenmod.TenantDomainStr, err.Error())
		logger.Log.Error(msg)
		return nil, fmt.Errorf(msg)
	}

	// Succesfully fetched the Domains, return the result.
	logger.Log.Infof("Retrieved %d %ss:\n", len(result.GetData()), tenmod.TenantDomainStr)
	return result, nil
}

// CreateTenantIngestionProfile - creates an Ingestion Profile scoped to a specific Tenant.
func (tsh *TenantServiceHandler) CreateTenantIngestionProfile(ctx context.Context, tenantIngPrfReq *pb.TenantIngestionProfile) (*pb.TenantIngestionProfile, error) {
	// Validate the request to ensure no invalid data is stored:
	if err := validateTenantIngPrfRequest(tenantIngPrfReq, false); err != nil {
		msg := fmt.Sprintf("Unable to validate request to store %s: %s", tenmod.TenantIngestionProfileStr, err.Error())
		logger.Log.Error(msg)
		return nil, fmt.Errorf(msg)
	}

	logger.Log.Infof("Creating %s: %s", tenmod.TenantIngestionProfileStr, tenantIngPrfReq)

	// Issue request to DAO Layer to Create the Tenant Ingestion Profile
	result, err := tsh.tenantDB.CreateTenantIngestionProfile(tenantIngPrfReq)
	if err != nil {
		msg := fmt.Sprintf("Unable to store %s: %s", tenmod.TenantIngestionProfileStr, err.Error())
		logger.Log.Error(msg)
		return nil, fmt.Errorf(msg)
	}

	// Succesfully Created the Ingestion Profile, return the result.
	logger.Log.Infof("Created %s: %s\n", tenmod.TenantIngestionProfileStr, result.GetXId())
	return result, nil
}

// UpdateTenantIngestionProfile - updates an Ingestion Profile scoped to a specific Tenant.
func (tsh *TenantServiceHandler) UpdateTenantIngestionProfile(ctx context.Context, tenantIngPrfReq *pb.TenantIngestionProfile) (*pb.TenantIngestionProfile, error) {
	// Validate the request to ensure no invalid data is stored:
	if err := validateTenantIngPrfRequest(tenantIngPrfReq, true); err != nil {
		msg := fmt.Sprintf("Unable to validate request to store %s: %s", tenmod.TenantIngestionProfileStr, err.Error())
		logger.Log.Error(msg)
		return nil, fmt.Errorf(msg)
	}

	logger.Log.Infof("Updating %s: %s", tenmod.TenantIngestionProfileStr, tenantIngPrfReq)

	// Issue request to DAO Layer to Update the Tenant Ingestion Profile
	result, err := tsh.tenantDB.UpdateTenantIngestionProfile(tenantIngPrfReq)
	if err != nil {
		msg := fmt.Sprintf("Unable to store %s: %s", tenmod.TenantIngestionProfileStr, err.Error())
		logger.Log.Error(msg)
		return nil, fmt.Errorf(msg)
	}

	// Succesfully Updated the Ingestion Profile, return the result.
	logger.Log.Infof("Updated %s: %s\n", tenmod.TenantIngestionProfileStr, result.GetXId())
	return result, nil
}

// GetTenantIngestionProfile - retrieves the Ingestion Profile for a single Tenant.
func (tsh *TenantServiceHandler) GetTenantIngestionProfile(ctx context.Context, tenantIngPrfIDReq *pb.TenantIngestionProfileIdRequest) (*pb.TenantIngestionProfile, error) {
	// Validate the request to ensure the operation is valid:
	if err := validateTenantIngPrfIDRequest(tenantIngPrfIDReq); err != nil {
		msg := fmt.Sprintf("Unable to validate request to fetch %s: %s", tenmod.TenantIngestionProfileStr, err.Error())
		logger.Log.Error(msg)
		return nil, fmt.Errorf(msg)
	}

	logger.Log.Infof("Retrieving %s for Tenant %s", tenmod.TenantIngestionProfileStr, tenantIngPrfIDReq.GetTenantId())

	// Issue request to DAO Layer to fetch the Tenant Ingestion Profile
	result, err := tsh.tenantDB.GetTenantIngestionProfile(tenantIngPrfIDReq)
	if err != nil {
		msg := fmt.Sprintf("Unable to fetch %s: %s", tenmod.TenantIngestionProfileStr, err.Error())
		logger.Log.Error(msg)
		return nil, fmt.Errorf(msg)
	}

	// Succesfully fetched the Ingestion Profile, return the result.
	logger.Log.Infof("Retrieved %s: %s\n", tenmod.TenantIngestionProfileStr, result.GetXId())
	return result, nil
}

// DeleteTenantIngestionProfile - deletes the Ingestion Profile for a single Tenant.
func (tsh *TenantServiceHandler) DeleteTenantIngestionProfile(ctx context.Context, tenantIngPrfIDReq *pb.TenantIngestionProfileIdRequest) (*pb.TenantIngestionProfile, error) {
	// Validate the request to ensure the operation is valid:
	if err := validateTenantIngPrfIDRequest(tenantIngPrfIDReq); err != nil {
		msg := fmt.Sprintf("Unable to validate request to delete %s: %s", tenmod.TenantIngestionProfileStr, err.Error())
		logger.Log.Error(msg)
		return nil, fmt.Errorf(msg)
	}

	logger.Log.Infof("Deleting %s for Tenant %s", tenmod.TenantIngestionProfileStr, tenantIngPrfIDReq.GetTenantId())

	// Issue request to DAO Layer to delete the Tenant Ingestion Profile
	result, err := tsh.tenantDB.DeleteTenantIngestionProfile(tenantIngPrfIDReq)
	if err != nil {
		msg := fmt.Sprintf("Unable to delete %s: %s", tenmod.TenantIngestionProfileStr, err.Error())
		logger.Log.Error(msg)
		return nil, fmt.Errorf(msg)
	}

	// Succesfully deleted the Ingestion Profile, return the result.
	logger.Log.Infof("Deleted %s: %s\n", tenmod.TenantIngestionProfileStr, result.GetXId())
	return result, nil
}

// CreateTenantThresholdProfile - creates an Threshold Profile scoped to a specific Tenant.
func (tsh *TenantServiceHandler) CreateTenantThresholdProfile(ctx context.Context, tenantThreshPrfReq *pb.TenantThresholdProfile) (*pb.TenantThresholdProfile, error) {
	// Validate the request to ensure no invalid data is stored:
	if err := validateTenantThreshPrfRequest(tenantThreshPrfReq, false); err != nil {
		msg := fmt.Sprintf("Unable to validate request to store %s: %s", tenmod.TenantThresholdProfileStr, err.Error())
		logger.Log.Error(msg)
		return nil, fmt.Errorf(msg)
	}

	logger.Log.Infof("Creating %s: %s", tenmod.TenantThresholdProfileStr, tenantThreshPrfReq)

	// Issue request to DAO Layer to Create the Tenant Threshold Profile
	result, err := tsh.tenantDB.CreateTenantThresholdProfile(tenantThreshPrfReq)
	if err != nil {
		msg := fmt.Sprintf("Unable to store %s: %s", tenmod.TenantThresholdProfileStr, err.Error())
		logger.Log.Error(msg)
		return nil, fmt.Errorf(msg)
	}

	// Succesfully Created the Threshold Profile, return the result.
	logger.Log.Infof("Created %s: %s\n", tenmod.TenantThresholdProfileStr, result.GetXId())
	return result, nil
}

// UpdateTenantThresholdProfile - updates an Threshold Profile scoped to a specific Tenant.
func (tsh *TenantServiceHandler) UpdateTenantThresholdProfile(ctx context.Context, tenantThreshPrfReq *pb.TenantThresholdProfile) (*pb.TenantThresholdProfile, error) {
	// Validate the request to ensure no invalid data is stored:
	if err := validateTenantThreshPrfRequest(tenantThreshPrfReq, true); err != nil {
		msg := fmt.Sprintf("Unable to validate request to store %s: %s", tenmod.TenantThresholdProfileStr, err.Error())
		logger.Log.Error(msg)
		return nil, fmt.Errorf(msg)
	}

	logger.Log.Infof("Updating %s: %s", tenmod.TenantThresholdProfileStr, tenantThreshPrfReq)

	// Issue request to DAO Layer to Update the Tenant Threshold Profile
	result, err := tsh.tenantDB.UpdateTenantThresholdProfile(tenantThreshPrfReq)
	if err != nil {
		msg := fmt.Sprintf("Unable to store %s: %s", tenmod.TenantThresholdProfileStr, err.Error())
		logger.Log.Error(msg)
		return nil, fmt.Errorf(msg)
	}

	// Succesfully Updated the Threshold Profile, return the result.
	logger.Log.Infof("Updated %s: %s\n", tenmod.TenantThresholdProfileStr, result.GetXId())
	return result, nil
}

// GetTenantThresholdProfile - retrieves the Threshold Profile for a single Tenant.
func (tsh *TenantServiceHandler) GetTenantThresholdProfile(ctx context.Context, tenantThreshPrfIDReq *pb.TenantThresholdProfileIdRequest) (*pb.TenantThresholdProfile, error) {
	// Validate the request to ensure the operation is valid:
	if err := validateTenantThreshPrfIDRequest(tenantThreshPrfIDReq); err != nil {
		msg := fmt.Sprintf("Unable to validate request to fetch %s: %s", tenmod.TenantThresholdProfileStr, err.Error())
		logger.Log.Error(msg)
		return nil, fmt.Errorf(msg)
	}

	logger.Log.Infof("Retrieving %s for Tenant %s", tenmod.TenantThresholdProfileStr, tenantThreshPrfIDReq.GetTenantId())

	// Issue request to DAO Layer to fetch the Tenant Threshold Profile
	result, err := tsh.tenantDB.GetTenantThresholdProfile(tenantThreshPrfIDReq)
	if err != nil {
		msg := fmt.Sprintf("Unable to fetch %s: %s", tenmod.TenantThresholdProfileStr, err.Error())
		logger.Log.Error(msg)
		return nil, fmt.Errorf(msg)
	}

	// Succesfully fetched the Threshold Profile, return the result.
	logger.Log.Infof("Retrieved %s: %s\n", tenmod.TenantThresholdProfileStr, result.GetXId())
	return result, nil
}

// DeleteTenantThresholdProfile - deletes the Threshold Profile for a single Tenant.
func (tsh *TenantServiceHandler) DeleteTenantThresholdProfile(ctx context.Context, tenantThreshPrfIDReq *pb.TenantThresholdProfileIdRequest) (*pb.TenantThresholdProfile, error) {
	// Validate the request to ensure the operation is valid:
	if err := validateTenantThreshPrfIDRequest(tenantThreshPrfIDReq); err != nil {
		msg := fmt.Sprintf("Unable to validate request to delete %s: %s", tenmod.TenantThresholdProfileStr, err.Error())
		logger.Log.Error(msg)
		return nil, fmt.Errorf(msg)
	}

	logger.Log.Infof("Deleting %s for Tenant %s", tenmod.TenantThresholdProfileStr, tenantThreshPrfIDReq.GetTenantId())

	// Issue request to DAO Layer to delete the Tenant Threshold Profile
	result, err := tsh.tenantDB.DeleteTenantThresholdProfile(tenantThreshPrfIDReq)
	if err != nil {
		msg := fmt.Sprintf("Unable to delete %s: %s", tenmod.TenantThresholdProfileStr, err.Error())
		logger.Log.Error(msg)
		return nil, fmt.Errorf(msg)
	}

	// Succesfully deleted the Threshold Profile, return the result.
	logger.Log.Infof("Deleted %s: %s\n", tenmod.TenantThresholdProfileStr, result.GetXId())
	return result, nil
}

// CreateMonitoredObject - creates a Monitored Object scoped to a specific tenant
func (tsh *TenantServiceHandler) CreateMonitoredObject(ctx context.Context, monitoredObjectReq *pb.MonitoredObject) (*pb.MonitoredObject, error) {
	// Validate the request to ensure no invalid data is stored:
	if err := validateMonitoredObjectRequest(monitoredObjectReq, false); err != nil {
		msg := fmt.Sprintf("Unable to validate request to store %s: %s", tenmod.TenantMonitoredObjectStr, err.Error())
		logger.Log.Error(msg)
		return nil, fmt.Errorf(msg)
	}

	logger.Log.Infof("Creating %s: %s", tenmod.TenantMonitoredObjectStr, monitoredObjectReq)

	// Issue request to DAO Layer to Create the Tenant Monitored Object
	result, err := tsh.tenantDB.CreateMonitoredObject(monitoredObjectReq)
	if err != nil {
		msg := fmt.Sprintf("Unable to store %s: %s", tenmod.TenantMonitoredObjectStr, err.Error())
		logger.Log.Error(msg)
		return nil, fmt.Errorf(msg)
	}

	// Succesfully Created the Monitored, return the result.
	logger.Log.Infof("Created %s: %s\n", tenmod.TenantMonitoredObjectStr, result.GetXId())
	return result, nil
}

// UpdateMonitoredObject - updates an MonitoredObject scoped to a specific Tenant.
func (tsh *TenantServiceHandler) UpdateMonitoredObject(ctx context.Context, monitoredObjectReq *pb.MonitoredObject) (*pb.MonitoredObject, error) {
	// Validate the request to ensure no invalid data is stored:
	if err := validateMonitoredObjectRequest(monitoredObjectReq, true); err != nil {
		msg := fmt.Sprintf("Unable to validate request to store %s: %s", tenmod.TenantMonitoredObjectStr, err.Error())
		logger.Log.Error(msg)
		return nil, fmt.Errorf(msg)
	}

	logger.Log.Infof("Updating %s: %s", tenmod.TenantMonitoredObjectStr, monitoredObjectReq)

	// Issue request to DAO Layer to Update the Tenant Monitored Object
	result, err := tsh.tenantDB.UpdateMonitoredObject(monitoredObjectReq)
	if err != nil {
		msg := fmt.Sprintf("Unable to store %s: %s", tenmod.TenantMonitoredObjectStr, err.Error())
		logger.Log.Error(msg)
		return nil, fmt.Errorf(msg)
	}

	// Succesfully Updated the Monitored Object, return the result.
	logger.Log.Infof("Updated %s: %s\n", tenmod.TenantMonitoredObjectStr, result.GetXId())
	return result, nil
}

// GetMonitoredObject - retrieves the MonitoredObject for a singler Tenant.
func (tsh *TenantServiceHandler) GetMonitoredObject(ctx context.Context, monitoredObjectIDReq *pb.MonitoredObjectIdRequest) (*pb.MonitoredObject, error) {
	// Validate the request to ensure no invalid data is stored:
	if err := validateMonitoredObjectIDRequest(monitoredObjectIDReq); err != nil {
		msg := fmt.Sprintf("Unable to validate request to fetch %s: %s", tenmod.TenantMonitoredObjectStr, err.Error())
		logger.Log.Error(msg)
		return nil, fmt.Errorf(msg)
	}

	// Issue request to DAO Layer to fetch the Tenant Monitored Object
	result, err := tsh.tenantDB.GetMonitoredObject(monitoredObjectIDReq)
	if err != nil {
		msg := fmt.Sprintf("Unable to fetch %s: %s", tenmod.TenantMonitoredObjectStr, err.Error())
		logger.Log.Error(msg)
		return nil, fmt.Errorf(msg)
	}

	// Succesfully fetched the Monitored Object, return the result.
	logger.Log.Infof("Retrieved %s: %s\n", tenmod.TenantMonitoredObjectStr, result.GetXId())
	return result, nil
}

// DeleteMonitoredObject - deletes the MonitoredObject for a singler Tenant.
func (tsh *TenantServiceHandler) DeleteMonitoredObject(ctx context.Context, monitoredObjectIDReq *pb.MonitoredObjectIdRequest) (*pb.MonitoredObject, error) {
	// Validate the request to ensure the operation is valid:
	if err := validateMonitoredObjectIDRequest(monitoredObjectIDReq); err != nil {
		msg := fmt.Sprintf("Unable to validate request to delete %s: %s", tenmod.TenantMonitoredObjectStr, err.Error())
		logger.Log.Error(msg)
		return nil, fmt.Errorf(msg)
	}

	logger.Log.Infof("Deleting %s for Tenant %s", tenmod.TenantMonitoredObjectStr, monitoredObjectIDReq.GetTenantId())

	// Issue request to DAO Layer to delete the Tenant Monitored Object
	result, err := tsh.tenantDB.DeleteMonitoredObject(monitoredObjectIDReq)
	if err != nil {
		msg := fmt.Sprintf("Unable to delete %s: %s", tenmod.TenantMonitoredObjectStr, err.Error())
		logger.Log.Error(msg)
		return nil, fmt.Errorf(msg)
	}

	// Succesfully deleted the MonitoredObject, return the result.
	logger.Log.Infof("Deleted %s: %s\n", tenmod.TenantMonitoredObjectStr, result.GetXId())
	return result, nil
}

// GetAllMonitoredObjects - retrieves all MonitoredObjects scoped to a single Tenant.
func (tsh *TenantServiceHandler) GetAllMonitoredObjects(ctx context.Context, tenantID *wr.StringValue) (*pb.MonitoredObjectList, error) {
	// Validate the request to ensure this operation is valid:

	logger.Log.Infof("Retrieving all %ss for Tenant: %s", tenmod.TenantMonitoredObjectStr, tenantID.Value)

	// Issue request to DAO Layer to fetch the Tenant Monitored Objects
	result, err := tsh.tenantDB.GetAllMonitoredObjects(tenantID.Value)
	if err != nil {
		msg := fmt.Sprintf("Unable to retrieve %ss: %s", tenmod.TenantMonitoredObjectStr, err.Error())
		logger.Log.Error(msg)
		return nil, fmt.Errorf(msg)
	}

	// Succesfully fetched the Monitored Objects, return the result.
	logger.Log.Infof("Retrieved %d %ss:\n", len(result.GetData()), tenmod.TenantMonitoredObjectStr)
	return result, nil
}

// GetMonitoredObjectToDomainMap - retrieves a mapping of MonitoredObjects to each Domain. Will retrieve the mapping either as a count, or as a set of all
// MonitoredObjects that use each Domain.
func (tsh *TenantServiceHandler) GetMonitoredObjectToDomainMap(ctx context.Context, moByDomReq *pb.MonitoredObjectCountByDomainRequest) (*pb.MonitoredObjectCountByDomainResponse, error) {
	// Validate the request:
	if err := validateMonitoredObjectToDomainMapRequest(moByDomReq); err != nil {
		msg := fmt.Sprintf("Unable to validate request to fetch %s: %s", tenmod.MonitoredObjectToDomainMapStr, err.Error())
		logger.Log.Error(msg)
		return nil, fmt.Errorf(msg)
	}

	// Issue request to DAO Layer to fetch the Tenant Monitored Object Map
	result, err := tsh.tenantDB.GetMonitoredObjectToDomainMap(moByDomReq)
	if err != nil {
		msg := fmt.Sprintf("Unable to fetch %s: %s", tenmod.MonitoredObjectToDomainMapStr, err.Error())
		logger.Log.Error(msg)
		return nil, fmt.Errorf(msg)
	}

	// Succesfully fetched the Monitored Object Map, return the result.
	logger.Log.Infof("Successfully retrieved %s: %s\n", tenmod.MonitoredObjectToDomainMapStr)
	return result, nil
}

// CreateTenantMeta - Create TenantMeta scoped to a Single Tenant.
func (tsh *TenantServiceHandler) CreateTenantMeta(ctx context.Context, meta *pb.TenantMetadata) (*pb.TenantMetadata, error) {
	// Validate the request to ensure no invalid data is stored:
	if err := validateTenantMetaRequest(meta, false); err != nil {
		msg := fmt.Sprintf("Unable to validate request to store %s: %s", tenmod.TenantMetaStr, err.Error())
		logger.Log.Error(msg)
		return nil, fmt.Errorf(msg)
	}

	logger.Log.Infof("Creating %s for Tenant %s", tenmod.TenantMetaStr, meta.GetData().GetTenantId())

	// Issue request to DAO Layer to Create the record
	result, err := tsh.tenantDB.CreateTenantMeta(meta)
	if err != nil {
		msg := fmt.Sprintf("Unable to store %s: %s", tenmod.TenantMetaStr, err.Error())
		logger.Log.Error(msg)
		return nil, fmt.Errorf(msg)
	}

	// Succesfully Created the record, return the result.
	logger.Log.Infof("Created %s: %s\n", tenmod.TenantMetaStr, result.GetXId())
	return result, nil
}

// UpdateTenantMeta - Update TenantMeta scoped to a single Tenant.
func (tsh *TenantServiceHandler) UpdateTenantMeta(ctx context.Context, meta *pb.TenantMetadata) (*pb.TenantMetadata, error) {
	// Validate the request to ensure no invalid data is stored:
	if err := validateTenantMetaRequest(meta, true); err != nil {
		msg := fmt.Sprintf("Unable to validate request to store %s: %s", tenmod.TenantMetaStr, err.Error())
		logger.Log.Error(msg)
		return nil, fmt.Errorf(msg)
	}

	logger.Log.Infof("Updating %s: %s", tenmod.TenantMetaStr, meta)

	// Issue request to DAO Layer to Update the record
	result, err := tsh.tenantDB.UpdateTenantMeta(meta)
	if err != nil {
		msg := fmt.Sprintf("Unable to store %s: %s", tenmod.TenantMetaStr, err.Error())
		logger.Log.Error(msg)
		return nil, fmt.Errorf(msg)
	}

	// Succesfully Updated the record, return the result.
	logger.Log.Infof("Updated %s: %s\n", tenmod.TenantMetaStr, result.GetXId())
	return result, nil
}

// DeleteTenantMeta - Delete TenantMeta scoped to a single Tenant.
func (tsh *TenantServiceHandler) DeleteTenantMeta(ctx context.Context, tenantID *wr.StringValue) (*pb.TenantMetadata, error) {

	logger.Log.Infof("Deleting %s for Tenant %s", tenmod.TenantMetaStr, tenantID.GetValue())

	// Issue request to DAO Layer to delete the record
	result, err := tsh.tenantDB.DeleteTenantMeta(tenantID.GetValue())
	if err != nil {
		msg := fmt.Sprintf("Unable to delete %s: %s", tenmod.TenantMetaStr, err.Error())
		logger.Log.Error(msg)
		return nil, fmt.Errorf(msg)
	}

	// Succesfully deleted the record, return the result.
	logger.Log.Infof("Deleted %s: %s\n", tenmod.TenantMetaStr, result.GetXId())
	return result, nil
}

// GetTenantMeta - Retrieve a User scoped to a single Tenant.
func (tsh *TenantServiceHandler) GetTenantMeta(ctx context.Context, tenantID *wr.StringValue) (*pb.TenantMetadata, error) {

	// Issue request to DAO Layer to fetch the record
	result, err := tsh.tenantDB.GetTenantMeta(tenantID.GetValue())
	if err != nil {
		msg := fmt.Sprintf("Unable to fetch %s: %s", tenmod.TenantMetaStr, err.Error())
		logger.Log.Error(msg)
		return nil, fmt.Errorf(msg)
	}

	// Succesfully fetched the record, return the result.
	logger.Log.Infof("Retrieved %s: %s\n", tenmod.TenantMetaStr, result.GetXId())
	return result, nil
}

// GetAllTenantThresholdProfiles - retieve all Tenant Thresholds.
func (tsh *TenantServiceHandler) GetAllTenantThresholdProfiles(ctx context.Context, tenantID *wr.StringValue) (*pb.TenantThresholdProfileList, error) {
	logger.Log.Infof("Retrieving all %ss for Tenant: %s", tenmod.TenantThresholdProfileStr, tenantID.Value)

	// Issue request to DAO Layer to fetch the records
	result, err := tsh.tenantDB.GetAllTenantThresholdProfile(tenantID.Value)
	if err != nil {
		msg := fmt.Sprintf("Unable to fetch %ss: %s", tenmod.TenantThresholdProfileStr, err.Error())
		logger.Log.Error(msg)
		return nil, fmt.Errorf(msg)
	}

	// Succesfully fetched the records, return the result.
	logger.Log.Infof("Retrieved %d %ss:\n", len(result.GetData()), tenmod.TenantThresholdProfileStr)
	return result, nil
}

// GetActiveTenantIngestionProfile - retrieves the active Ingestion Profile for a single Tenant.
func (tsh *TenantServiceHandler) GetActiveTenantIngestionProfile(ctx context.Context, tenantID *wr.StringValue) (*pb.TenantIngestionProfile, error) {
	// Issue request to DAO Layer to fetch the record
	result, err := tsh.tenantDB.GetActiveTenantIngestionProfile(tenantID.GetValue())
	if err != nil {
		msg := fmt.Sprintf("Unable to fetch %s: %s", tenmod.TenantIngestionProfileStr, err.Error())
		logger.Log.Error(msg)
		return nil, fmt.Errorf(msg)
	}

	// Succesfully fetched the record, return the result.
	logger.Log.Infof("Retrieved %s: %s\n", tenmod.TenantIngestionProfileStr, result.GetXId())
	return result, nil
}

// BulkInsertMonitoredObjects - perform a bulk operation on a set of Monitored Objects.
func (tsh *TenantServiceHandler) BulkInsertMonitoredObjects(ctx context.Context, value *pb.TenantMonitoredObjectSet) (*pb.BulkOperationResponse, error) {
	// Validate the request:
	if value == nil {
		msg := fmt.Sprintf("Unable to Update %ss in bulk: %s", tenmod.TenantMonitoredObjectStr, "No Monitored Object data provided")
		logger.Log.Error(msg)
		return nil, fmt.Errorf(msg)
	}

	for _, mo := range value.MonitoredObjectSet {
		if err := validateMonitoredObjectRequest(mo, false); err != nil {
			msg := fmt.Sprintf("Unable to Update %ss in bulk: %s", tenmod.TenantMonitoredObjectStr, err.Error())
			logger.Log.Error(msg)
			return nil, fmt.Errorf(msg)
		}

		if value.TenantId != mo.Data.TenantId {
			msg := fmt.Sprintf("Unable to Update %ss in bulk: %s", tenmod.TenantMonitoredObjectStr, "All Monitored Objects must have Tenant ID "+value.TenantId)
			logger.Log.Error(msg)
			return nil, fmt.Errorf(msg)
		}
	}

	// Issue request to DAO Layer to insert the MOs
	result, err := tsh.tenantDB.BulkInsertMonitoredObjects(value)
	if err != nil {
		msg := fmt.Sprintf("Unable to update %ss in bulk: %s", tenmod.TenantMonitoredObjectStr, err.Error())
		logger.Log.Error(msg)
		return nil, fmt.Errorf(msg)
	}

	// Succesfully inserted the MOs.
	logger.Log.Infof("Successfully completed bulk request on %ss\n", tenmod.TenantMonitoredObjectStr)
	return result, nil
}
