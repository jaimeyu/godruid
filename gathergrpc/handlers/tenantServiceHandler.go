package handlers

import (
	"context"
	"errors"

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
		logger.Log.Fatalf("Unable to instantiate TenantServiceHandler: %v", err)
	}
	result.tenantDB = db

	return result
}

func getTenantServiceDatastore() (db.TenantServiceDatastore, error) {
	cfg, err := gather.GetActiveConfig()
	if err != nil {
		logger.Log.Errorf("Falied to instantiate TenantServiceHandler: %v", err)
		return nil, err
	}

	dbType := cfg.ServerConfig.StartupArgs.TenantDB
	if dbType == gather.COUCH {
		logger.Log.Debug("TenantService DB is using CouchDB Implementation")
		return couchDB.CreateTenantServiceDAO(), nil
	} else if dbType == gather.MEM {
		logger.Log.Debug("TenantService DB is using InMemory Implementation")
		return inMemory.CreateTenantServiceDAO(), nil
	}

	return nil, errors.New("No DB implementation provided for Admin Service. Check configuration")
}

// CreateTenantUser - creates a user scoped to a single Tenant.
func (tsh *TenantServiceHandler) CreateTenantUser(ctx context.Context, tenantUserReq *pb.TenantUserRequest) (*pb.TenantUserResponse, error) {
	// Validate the request to ensure no invalid data is stored:
	if err := validateTenantUserRequest(tenantUserReq); err != nil {
		return nil, err
	}

	logger.Log.Infof("Creating %s: %s", db.TenantUserStr, tenantUserReq.GetUser())

	// Issue request to DAO Layer to Create the Tenant User
	result, err := tsh.tenantDB.CreateTenantUser(tenantUserReq)
	if err != nil {
		logger.Log.Errorf("Unable to store %s: %v\n", db.TenantUserStr, err)
		return nil, err
	}

	// Succesfully Created the User, return the result.
	logger.Log.Infof("Created %s: %v\n", db.TenantUserStr, result)
	return &pb.TenantUserResponse{TenantId: tenantUserReq.GetTenantId(), User: result}, nil
}

// UpdateTenantUser - updates a user scoped to a single Tenant.
func (tsh *TenantServiceHandler) UpdateTenantUser(ctx context.Context, tenantUserReq *pb.TenantUserRequest) (*pb.TenantUserResponse, error) {
	// Validate the request to ensure no invalid data is stored:
	if err := validateTenantUserRequest(tenantUserReq); err != nil {
		return nil, err
	}

	logger.Log.Infof("Updating %s: %s", db.TenantUserStr, tenantUserReq.GetUser())

	// Issue request to DAO Layer to Update the Tenant User
	result, err := tsh.tenantDB.UpdateTenantUser(tenantUserReq)
	if err != nil {
		logger.Log.Errorf("Unable to store %s: %v\n", db.TenantUserStr, err)
		return nil, err
	}

	// Succesfully Updated the User, return the result.
	logger.Log.Infof("Updated %s: %v\n", db.TenantUserStr, result)
	return &pb.TenantUserResponse{TenantId: tenantUserReq.GetTenantId(), User: result}, nil
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
		logger.Log.Errorf("Unable to delete %s: %v\n", db.TenantUserStr, err)
		return nil, err
	}

	// Succesfully Deleted the User, return the result.
	logger.Log.Infof("Deleted %s: %v\n", db.TenantUserStr, result)
	return &pb.TenantUserResponse{TenantId: tenantUserIDReq.GetTenantId(), User: result}, nil
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
		logger.Log.Errorf("Unable to retrieve %s: %v\n", db.TenantUserStr, err)
		return nil, err
	}

	// Succesfully fetched the User, return the result.
	logger.Log.Infof("Retrieved %s: %v\n", db.TenantUserStr, result)
	return &pb.TenantUserResponse{TenantId: tenantUserIDReq.GetTenantId(), User: result}, nil
}

// GetAllTenantUsers - retrieves all users scoped to a single Tenant.
func (tsh *TenantServiceHandler) GetAllTenantUsers(ctx context.Context, tenantID *wr.StringValue) (*pb.TenantUserListResponse, error) {
	// Validate the request to ensure this operatin is valid:

	logger.Log.Infof("Retrieving all %ss for Tenant: %s", db.TenantUserStr, tenantID.Value)

	// Issue request to DAO Layer to fetch the Tenant Users
	result, err := tsh.tenantDB.GetAllTenantUsers(tenantID.Value)
	if err != nil {
		logger.Log.Errorf("Unable to retrieve %ss: %v\n", db.TenantUserStr, err)
		return nil, err
	}

	// Succesfully fetched the Users, return the result.
	logger.Log.Infof("Retrieved %d %ss:\n", db.TenantUserStr, len(result.List))
	return &pb.TenantUserListResponse{TenantId: tenantID.Value, List: result}, nil
}

// CreateTenantDomain - creates a Domain scoped to a single Tenant.
func (tsh *TenantServiceHandler) CreateTenantDomain(ctx context.Context, tenantDomainRequest *pb.TenantDomainRequest) (*pb.TenantDomainResponse, error) {
	// Validate the request to ensure no invalid data is stored:
	if err := validateTenantDomainRequest(tenantDomainRequest); err != nil {
		return nil, err
	}

	logger.Log.Infof("Creating %s: %s", db.TenantDomainStr, tenantDomainRequest.GetDomain())

	// Issue request to DAO Layer to Create the Tenant Domain
	result, err := tsh.tenantDB.CreateTenantDomain(tenantDomainRequest)
	if err != nil {
		logger.Log.Errorf("Unable to store %s: %v\n", db.TenantDomainStr, err)
		return nil, err
	}

	// Succesfully Created the Domain, return the result.
	logger.Log.Infof("Created %s: %v\n", db.TenantDomainStr, result)
	return &pb.TenantDomainResponse{TenantId: tenantDomainRequest.GetTenantId(), Domain: result}, nil
}

// UpdateTenantDomain - updates a Domain scoped to a single Tenant.
func (tsh *TenantServiceHandler) UpdateTenantDomain(ctx context.Context, tenantDomainRequest *pb.TenantDomainRequest) (*pb.TenantDomainResponse, error) {
	// Validate the request to ensure no invalid data is stored:
	if err := validateTenantDomainRequest(tenantDomainRequest); err != nil {
		return nil, err
	}

	logger.Log.Infof("Updating %s: %s", db.TenantDomainStr, tenantDomainRequest.GetDomain())

	// Issue request to DAO Layer to Update the Tenant Domain
	result, err := tsh.tenantDB.UpdateTenantDomain(tenantDomainRequest)
	if err != nil {
		logger.Log.Errorf("Unable to store %s: %v\n", db.TenantDomainStr, err)
		return nil, err
	}

	// Succesfully Updated the Domain, return the result.
	logger.Log.Infof("Updated %s: %v\n", db.TenantDomainStr, result)
	return &pb.TenantDomainResponse{TenantId: tenantDomainRequest.GetTenantId(), Domain: result}, nil
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
		logger.Log.Errorf("Unable to delete %s: %v\n", db.TenantDomainStr, err)
		return nil, err
	}

	// Succesfully Deleted the Domain, return the result.
	logger.Log.Infof("Deleted %s: %v\n", db.TenantDomainStr, result)
	return &pb.TenantDomainResponse{TenantId: tenantDomainIDRequest.GetTenantId(), Domain: result}, nil
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
		logger.Log.Errorf("Unable to retrieve %s: %v\n", db.TenantDomainStr, err)
		return nil, err
	}

	// Succesfully fetched the Domain, return the result.
	logger.Log.Infof("Retrieved %s: %v\n", db.TenantDomainStr, result)
	return &pb.TenantDomainResponse{TenantId: tenantDomainIDRequest.GetTenantId(), Domain: result}, nil
}

// GetAllTenantDomains - retrieves all Domains scoped to a single Tenant.
func (tsh *TenantServiceHandler) GetAllTenantDomains(ctx context.Context, tenantID *wr.StringValue) (*pb.TenantDomainListResponse, error) {
	// Validate the request to ensure this operation is valid:

	logger.Log.Infof("Retrieving all %ss for Tenant: %s", db.TenantDomainStr, tenantID.Value)

	// Issue request to DAO Layer to fetch the Tenant Domains
	result, err := tsh.tenantDB.GetAllTenantDomains(tenantID.Value)
	if err != nil {
		logger.Log.Errorf("Unable to retrieve %ss: %v\n", db.TenantDomainStr, err)
		return nil, err
	}

	// Succesfully fetched the Domains, return the result.
	logger.Log.Infof("Retrieved %d %ss:\n", db.TenantDomainStr, len(result.List))
	return &pb.TenantDomainListResponse{TenantId: tenantID.Value, List: result}, nil
}

// CreateTenantIngestionProfile - creates an Ingestion Profile scoped to a specific Tenant.
func (tsh *TenantServiceHandler) CreateTenantIngestionProfile(ctx context.Context, tenantIngPrfReq *pb.TenantIngestionProfileRequest) (*pb.TenantIngestionProfileResponse, error) {
	// Validate the request to ensure no invalid data is stored:
	if err := validateTenantIngPrfRequest(tenantIngPrfReq); err != nil {
		return nil, err
	}

	logger.Log.Infof("Creating %s: %s", db.TenantIngestionProfileStr, tenantIngPrfReq.GetIngestionProfile())

	// Issue request to DAO Layer to Create the Tenant Ingestion Profile
	result, err := tsh.tenantDB.CreateTenantIngestionProfile(tenantIngPrfReq)
	if err != nil {
		logger.Log.Errorf("Unable to store %s: %v\n", db.TenantIngestionProfileStr, err)
		return nil, err
	}

	// Succesfully Created the Ingestion Profile, return the result.
	logger.Log.Infof("Created %s: %v\n", db.TenantIngestionProfileStr, result)
	return &pb.TenantIngestionProfileResponse{TenantId: tenantIngPrfReq.GetTenantId(), IngestionProfile: result}, nil
}

// UpdateTenantIngestionProfile - updates an Ingestion Profile scoped to a specific Tenant.
func (tsh *TenantServiceHandler) UpdateTenantIngestionProfile(ctx context.Context, tenantIngPrfReq *pb.TenantIngestionProfileRequest) (*pb.TenantIngestionProfileResponse, error) {
	// Validate the request to ensure no invalid data is stored:
	if err := validateTenantIngPrfRequest(tenantIngPrfReq); err != nil {
		return nil, err
	}

	logger.Log.Infof("Updating %s: %s", db.TenantIngestionProfileStr, tenantIngPrfReq.GetIngestionProfile())

	// Issue request to DAO Layer to Update the Tenant Ingestion Profile
	result, err := tsh.tenantDB.UpdateTenantIngestionProfile(tenantIngPrfReq)
	if err != nil {
		logger.Log.Errorf("Unable to store %s: %v\n", db.TenantIngestionProfileStr, err)
		return nil, err
	}

	// Succesfully Updated the Ingestion Profile, return the result.
	logger.Log.Infof("Updated %s: %v\n", db.TenantIngestionProfileStr, result)
	return &pb.TenantIngestionProfileResponse{TenantId: tenantIngPrfReq.GetTenantId(), IngestionProfile: result}, nil
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
		logger.Log.Errorf("Unable to retrieve %s: %v\n", db.TenantIngestionProfileStr, err)
		return nil, err
	}

	// Succesfully fetched the Ingestion Profile, return the result.
	logger.Log.Infof("Retrieved %s: %v\n", db.TenantIngestionProfileStr, result)
	return &pb.TenantIngestionProfileResponse{TenantId: tenantIngPrfIDReq.GetTenantId(), IngestionProfile: result}, nil
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
		logger.Log.Errorf("Unable to delete %s: %v\n", db.TenantIngestionProfileStr, err)
		return nil, err
	}

	// Succesfully deleted the Ingestion Profile, return the result.
	logger.Log.Infof("Deleted %s: %v\n", db.TenantIngestionProfileStr, result)
	return &pb.TenantIngestionProfileResponse{TenantId: tenantIngPrfIDReq.GetTenantId(), IngestionProfile: result}, nil
}

func validateTenantUserRequest(request *pb.TenantUserRequest) error {
	if request == nil || request.GetUser() == nil {
		return errors.New("Invalid TenantUserRequest: no Tenant User data provided")
	}

	if len(request.GetTenantId()) == 0 {
		return errors.New("Invalid TenantUserRequest: no Tenant Id provided")
	}

	return nil
}

func validateTenantUserIDRequest(request *pb.TenantUserIdRequest) error {
	if request == nil || len(request.GetUserId()) == 0 {
		return errors.New("Invalid TenantUserIdRequest: no Tenant User ID data provided")
	}

	if len(request.GetTenantId()) == 0 {
		return errors.New("Invalid TenantUserIdRequest: no Tenant Id provided")
	}

	return nil
}

func validateTenantDomainRequest(request *pb.TenantDomainRequest) error {
	if request == nil || request.GetDomain() == nil {
		return errors.New("Invalid TenantDomainRequest: no Tenant Domain data provided")
	}

	if len(request.GetTenantId()) == 0 {
		return errors.New("Invalid TenantDomainRequest: no Tenant Id provided")
	}

	return nil
}

func validateTenantDomainIDRequest(request *pb.TenantDomainIdRequest) error {
	if request == nil || len(request.GetDomainId()) == 0 {
		return errors.New("Invalid TenantDomainIdRequest: no Tenant Domain ID data provided")
	}

	if len(request.GetTenantId()) == 0 {
		return errors.New("Invalid TenantDomainIdRequest: no Tenant Id provided")
	}

	return nil
}

func validateTenantIngPrfRequest(request *pb.TenantIngestionProfileRequest) error {
	if request == nil || request.GetIngestionProfile() == nil {
		return errors.New("Invalid TenantIngestionProfileRequest: no Tenant Ingestion Profile data provided")
	}

	if len(request.GetTenantId()) == 0 {
		return errors.New("Invalid TenantIngestionProfileRequest: no Tenant Id provided")
	}

	return nil
}

func validateTenantIngPrfIDRequest(request *pb.TenantIngestionProfileIdRequest) error {
	if request == nil || len(request.GetIngestionProfileId()) == 0 {
		return errors.New("Invalid TenantIngestionProfileIdRequest: no Ingestion Profile ID data provided")
	}

	if len(request.GetTenantId()) == 0 {
		return errors.New("Invalid TenantIngestionProfileIdRequest: no Tenant Id provided")
	}

	return nil
}
