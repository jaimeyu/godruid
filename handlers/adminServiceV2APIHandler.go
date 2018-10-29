package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/accedian/adh-gather/models"
	"github.com/accedian/adh-gather/swagmodels"
	"github.com/manyminds/api2go/jsonapi"

	"github.com/accedian/adh-gather/datastore"
	"github.com/accedian/adh-gather/logger"
	admmod "github.com/accedian/adh-gather/models/admin"
	tenmod "github.com/accedian/adh-gather/models/tenant"
	mon "github.com/accedian/adh-gather/monitoring"
	"github.com/accedian/adh-gather/restapi/operations/admin_provisioning_service_v2"
	"github.com/go-openapi/runtime/middleware"
)

// HandleCreateTenantV2 - create a new tenant
func HandleCreateTenantV2(allowedRoles []string, adminDB datastore.AdminServiceDatastore, tenantDB datastore.TenantServiceDatastore) func(params admin_provisioning_service_v2.CreateTenantV2Params) middleware.Responder {
	return func(params admin_provisioning_service_v2.CreateTenantV2Params) middleware.Responder {

		// Do the work
		startTime, responseCode, response, err := doCreateTenantV2(allowedRoles, adminDB, tenantDB, params)

		// Success Response
		if responseCode == http.StatusCreated {
			return admin_provisioning_service_v2.NewCreateTenantV2Created().WithPayload(response)
		}

		// Error Responses
		errorMessage := reportAPIError(err.Error(), startTime, responseCode, mon.CreateTenantStr, mon.APICompleted, mon.AdminAPICompleted)
		switch responseCode {
		case http.StatusForbidden:
			return admin_provisioning_service_v2.NewCreateTenantV2Forbidden().WithPayload(errorMessage)
		case http.StatusBadRequest:
			return admin_provisioning_service_v2.NewCreateTenantV2BadRequest().WithPayload(errorMessage)
		case http.StatusConflict:
			return admin_provisioning_service_v2.NewCreateTenantV2Conflict().WithPayload(errorMessage)
		default:
			return admin_provisioning_service_v2.NewCreateTenantV2InternalServerError().WithPayload(errorMessage)
		}
	}
}

// HandleGetTenantV2 - retrieve a tenant by the tenant ID.
func HandleGetTenantV2(allowedRoles []string, adminDB datastore.AdminServiceDatastore) func(params admin_provisioning_service_v2.GetTenantV2Params) middleware.Responder {
	return func(params admin_provisioning_service_v2.GetTenantV2Params) middleware.Responder {

		// Do the work
		startTime, responseCode, response, err := doGetTenantV2(allowedRoles, adminDB, params)

		// Success Response
		if responseCode == http.StatusOK {
			return admin_provisioning_service_v2.NewGetTenantV2OK().WithPayload(response)
		}

		// Error Responses
		errorMessage := reportAPIError(err.Error(), startTime, responseCode, mon.GetTenantStr, mon.APICompleted, mon.AdminAPICompleted)
		switch responseCode {
		case http.StatusForbidden:
			return admin_provisioning_service_v2.NewGetTenantV2Forbidden().WithPayload(errorMessage)
		case http.StatusNotFound:
			return admin_provisioning_service_v2.NewGetTenantV2NotFound().WithPayload(errorMessage)
		default:
			return admin_provisioning_service_v2.NewGetTenantV2InternalServerError().WithPayload(errorMessage)
		}

	}
}

// HandlePatchTenantV2 - update a Tenant record
func HandlePatchTenantV2(allowedRoles []string, adminDB datastore.AdminServiceDatastore) func(params admin_provisioning_service_v2.PatchTenantV2Params) middleware.Responder {
	return func(params admin_provisioning_service_v2.PatchTenantV2Params) middleware.Responder {

		// Do the work
		startTime, responseCode, response, err := doUpdateTenantV2(allowedRoles, adminDB, params)

		// Success Response
		if responseCode == http.StatusOK {
			return admin_provisioning_service_v2.NewPatchTenantV2OK().WithPayload(response)
		}

		// Error Responses
		errorMessage := reportAPIError(err.Error(), startTime, responseCode, mon.PatchTenantStr, mon.APICompleted, mon.AdminAPICompleted)
		switch responseCode {
		case http.StatusForbidden:
			return admin_provisioning_service_v2.NewPatchTenantV2Forbidden().WithPayload(errorMessage)
		case http.StatusBadRequest:
			return admin_provisioning_service_v2.NewPatchTenantV2BadRequest().WithPayload(errorMessage)
		case http.StatusNotFound:
			return admin_provisioning_service_v2.NewPatchTenantV2NotFound().WithPayload(errorMessage)
		case http.StatusConflict:
			return admin_provisioning_service_v2.NewPatchTenantV2Conflict().WithPayload(errorMessage)
		default:
			return admin_provisioning_service_v2.NewPatchTenantV2InternalServerError().WithPayload(errorMessage)
		}
	}
}

// HandleDeleteTenantV2 - delete a tenant by the tenant ID.
func HandleDeleteTenantV2(allowedRoles []string, adminDB datastore.AdminServiceDatastore) func(params admin_provisioning_service_v2.DeleteTenantV2Params) middleware.Responder {
	return func(params admin_provisioning_service_v2.DeleteTenantV2Params) middleware.Responder {

		// Do the work
		startTime, responseCode, response, err := doDeleteTenantV2(allowedRoles, adminDB, params)

		// Success Response
		if responseCode == http.StatusOK {
			return admin_provisioning_service_v2.NewDeleteTenantV2OK().WithPayload(response)
		}

		// Error Responses
		errorMessage := reportAPIError(err.Error(), startTime, responseCode, mon.DeleteTenantStr, mon.APICompleted, mon.AdminAPICompleted)
		switch responseCode {
		case http.StatusForbidden:
			return admin_provisioning_service_v2.NewDeleteTenantV2Forbidden().WithPayload(errorMessage)
		case http.StatusNotFound:
			return admin_provisioning_service_v2.NewDeleteTenantV2NotFound().WithPayload(errorMessage)
		default:
			return admin_provisioning_service_v2.NewDeleteTenantV2InternalServerError().WithPayload(errorMessage)
		}

	}
}

// HandleGetAllTenantsV2 - retrieve all tenants
func HandleGetAllTenantsV2(allowedRoles []string, adminDB datastore.AdminServiceDatastore) func(params admin_provisioning_service_v2.GetAllTenantsV2Params) middleware.Responder {
	return func(params admin_provisioning_service_v2.GetAllTenantsV2Params) middleware.Responder {

		// Do the work
		startTime, responseCode, response, err := doGetAllTenantsV2(allowedRoles, adminDB, params)

		// Success Response
		if responseCode == http.StatusOK {
			return admin_provisioning_service_v2.NewGetAllTenantsV2OK().WithPayload(response)
		}

		// Error Responses
		errorMessage := reportAPIError(err.Error(), startTime, responseCode, mon.GetAllTenantStr, mon.APICompleted, mon.AdminAPICompleted)
		switch responseCode {
		case http.StatusForbidden:
			return admin_provisioning_service_v2.NewGetAllTenantsV2Forbidden().WithPayload(errorMessage)
		case http.StatusNotFound:
			return admin_provisioning_service_v2.NewGetAllTenantsV2NotFound().WithPayload(errorMessage)
		default:
			return admin_provisioning_service_v2.NewGetAllTenantsV2InternalServerError().WithPayload(errorMessage)
		}

	}
}

// HandleGetTenantIDByAliasV2 - returns the tenant id as a string
func HandleGetTenantIDByAliasV2(adminDB datastore.AdminServiceDatastore) func(params admin_provisioning_service_v2.GetTenantIDByAliasV2Params) middleware.Responder {
	return func(params admin_provisioning_service_v2.GetTenantIDByAliasV2Params) middleware.Responder {

		// Do the work
		startTime, responseCode, response, err := doGetTenantIDByAliasV2(adminDB, params)

		// Success Response
		if responseCode == http.StatusOK {
			return admin_provisioning_service_v2.NewGetTenantIDByAliasV2OK().WithPayload(response)
		}

		// Error Responses
		errorMessage := reportAPIError(err.Error(), startTime, responseCode, mon.GetTenantIDByAliasStr, mon.APICompleted, mon.AdminAPICompleted)
		switch responseCode {
		case http.StatusNotFound:
			return admin_provisioning_service_v2.NewGetTenantIDByAliasV2NotFound().WithPayload(errorMessage)
		default:
			return admin_provisioning_service_v2.NewGetTenantIDByAliasV2InternalServerError().WithPayload(errorMessage)
		}

	}
}

// HandleGetTenantSummaryByAliasV2 - returns the tenant summary for an alias
func HandleGetTenantSummaryByAliasV2(adminDB datastore.AdminServiceDatastore) func(params admin_provisioning_service_v2.GetTenantSummaryByAliasV2Params) middleware.Responder {
	return func(params admin_provisioning_service_v2.GetTenantSummaryByAliasV2Params) middleware.Responder {

		// Do the work
		startTime, responseCode, response, err := doGetTenantSummaryByAliasV2(adminDB, params)

		// Success Response
		if responseCode == http.StatusOK {
			return admin_provisioning_service_v2.NewGetTenantSummaryByAliasV2OK().WithPayload(response)
		}

		// Error Responses
		errorMessage := reportAPIError(err.Error(), startTime, responseCode, mon.GetTenantSummaryByAliasStr, mon.APICompleted, mon.AdminAPICompleted)
		switch responseCode {
		case http.StatusNotFound:
			return admin_provisioning_service_v2.NewGetTenantSummaryByAliasV2NotFound().WithPayload(errorMessage)
		default:
			return admin_provisioning_service_v2.NewGetTenantSummaryByAliasV2InternalServerError().WithPayload(errorMessage)
		}
	}
}

// HandleGetIngestionDictionaryV2 - retrieve an ingestion dictionary
func HandleGetIngestionDictionaryV2(allowedRoles []string, adminDB datastore.AdminServiceDatastore) func(params admin_provisioning_service_v2.GetIngestionDictionaryV2Params) middleware.Responder {
	return func(params admin_provisioning_service_v2.GetIngestionDictionaryV2Params) middleware.Responder {

		// Do the work
		startTime, responseCode, response, err := doGetIngestionDictionaryV2(allowedRoles, adminDB, params)

		// Success Response
		if responseCode == http.StatusOK {
			return admin_provisioning_service_v2.NewGetIngestionDictionaryV2OK().WithPayload(response)
		}

		// Error Responses
		errorMessage := reportAPIError(err.Error(), startTime, responseCode, mon.GetIngDictStr, mon.APICompleted, mon.AdminAPICompleted)
		switch responseCode {
		case http.StatusForbidden:
			return admin_provisioning_service_v2.NewGetIngestionDictionaryV2Forbidden().WithPayload(errorMessage)
		default:
			return admin_provisioning_service_v2.NewGetIngestionDictionaryV2InternalServerError().WithPayload(errorMessage)
		}
	}
}

// HandleGetValidTypesV2 - retrieve an ingestion dictionary
func HandleGetValidTypesV2(allowedRoles []string, adminDB datastore.AdminServiceDatastore) func(params admin_provisioning_service_v2.GetValidTypesV2Params) middleware.Responder {
	return func(params admin_provisioning_service_v2.GetValidTypesV2Params) middleware.Responder {

		// Do the work
		startTime, responseCode, response, err := doGetValidTypesV2(allowedRoles, adminDB, params)

		// Success Response
		if responseCode == http.StatusOK {
			return admin_provisioning_service_v2.NewGetValidTypesV2OK().WithPayload(response)
		}

		// Error Responses
		errorMessage := reportAPIError(err.Error(), startTime, responseCode, mon.GetValidTypesStr, mon.APICompleted, mon.AdminAPICompleted)
		switch responseCode {
		case http.StatusForbidden:
			return admin_provisioning_service_v2.NewGetValidTypesV2Forbidden().WithPayload(errorMessage)
		default:
			return admin_provisioning_service_v2.NewGetValidTypesV2InternalServerError().WithPayload(errorMessage)
		}
	}
}

func doCreateTenantV2(allowedRoles []string, adminDB datastore.AdminServiceDatastore, tenantDB datastore.TenantServiceDatastore, params admin_provisioning_service_v2.CreateTenantV2Params) (time.Time, int, *swagmodels.TenantResponse, error) {
	isAuthorized, startTime := authorizeRequest(fmt.Sprintf("Creating %s: %s", admmod.TenantStr, models.AsJSONString(params.Body)), params.HTTPRequest, allowedRoles, mon.APIRecieved, mon.AdminAPIRecieved)

	if !isAuthorized {
		return startTime, http.StatusForbidden, nil, fmt.Errorf("Create %s operation not authorized for role: %s", admmod.TenantStr, params.HTTPRequest.Header.Get(XFwdUserRoles))
	}

	// Unmarshal the request
	requestBytes, err := json.Marshal(params.Body)
	if err != nil {
		return startTime, http.StatusBadRequest, nil, err
	}

	data := admmod.Tenant{}
	err = jsonapi.Unmarshal(requestBytes, &data)
	if err != nil {
		return startTime, http.StatusBadRequest, nil, err
	}

	// Check if a tenant already exists with this name.
	existingTenantByName, _ := adminDB.GetTenantIDByAlias(strings.ToLower(data.Name))
	if len(existingTenantByName) != 0 {
		return startTime, http.StatusConflict, nil, fmt.Errorf("Unable to create Tenant %s. A Tenant with this name already exists", data.Name)
	}

	// Issue request to DAO Layer to Create Tenant
	result, err := adminDB.CreateTenant(&data)
	if err != nil {
		return startTime, http.StatusInternalServerError, nil, fmt.Errorf("Unable to store %s: %s", admmod.TenantStr, err.Error())
	}

	// Create a default Ingestion Profile for the Tenant.
	idForTenant := result.ID
	ingPrfData := createDefaultTenantIngPrf(idForTenant)
	_, err = tenantDB.CreateTenantIngestionProfile(ingPrfData)
	if err != nil {
		return startTime, http.StatusInternalServerError, nil, fmt.Errorf("Unable to create default Ingestion Profile %s", err.Error())
	}

	// Create a default Metadata Config for the Tenant.
	_, err = tenantDB.CreateTenantMetadataConfig(&tenmod.MetadataConfig{TenantID: idForTenant})
	if err != nil {
		return startTime, http.StatusInternalServerError, nil, fmt.Errorf("Unable to create default Metadata Config %s", err.Error())
	}

	// Create a default Threshold Profile for the Tenant
	threshPrfData := createDefaultTenantThresholdPrf(idForTenant)
	threshProfileResponse, err := tenantDB.CreateTenantThresholdProfile(threshPrfData)
	if err != nil {
		return startTime, http.StatusInternalServerError, nil, fmt.Errorf("Unable to create default Threshold Profile %s", err.Error())
	}

	// Create a default Data Cleaning Profile for the Tenant
	dcp := &tenmod.DataCleaningProfile{
		TenantID: idForTenant,
		Rules:    []*tenmod.DataCleaningRule{},
	}
	_, err = tenantDB.CreateTenantDataCleaningProfile(dcp)
	if err != nil {
		return startTime, http.StatusInternalServerError, nil, fmt.Errorf("Unable to create Tenant Data Cleaning Profile %s", err.Error())
	}

	// Create the tenant metadata
	// For the IDs used as references inside other objects, need to strip off the 'thresholdProfile_2_'
	// as this is just relational pouch adaption:
	meta := createDefaultTenantMeta(idForTenant, threshProfileResponse.ID, result.Name)
	_, err = tenantDB.CreateTenantMeta(meta)
	if err != nil {
		return startTime, http.StatusInternalServerError, nil, fmt.Errorf("Unable to create Tenant metadata %s", err.Error())
	}

	converted := swagmodels.TenantResponse{}
	err = convertToJsonapiObject(result, &converted)
	if err != nil {
		return startTime, http.StatusInternalServerError, nil, fmt.Errorf("Unable to convert %s data to jsonapi return format: %s", admmod.TenantStr, err.Error())
	}

	reportAPICompletionState(startTime, http.StatusOK, mon.CreateTenantStr, mon.APICompleted, mon.AdminAPICompleted)
	logger.Log.Infof("Created %s %s", admmod.TenantStr, models.AsJSONString(converted))
	return startTime, http.StatusCreated, &converted, nil
}

func doGetTenantV2(allowedRoles []string, adminDB datastore.AdminServiceDatastore, params admin_provisioning_service_v2.GetTenantV2Params) (time.Time, int, *swagmodels.TenantResponse, error) {
	isAuthorized, startTime := authorizeRequest(fmt.Sprintf("Fetching %s: %s", admmod.TenantStr, params.TenantID), params.HTTPRequest, allowedRoles, mon.APIRecieved, mon.AdminAPIRecieved)

	if !isAuthorized {
		return startTime, http.StatusForbidden, nil, fmt.Errorf("Fetch %s operation not authorized for role: %s", admmod.TenantStr, params.HTTPRequest.Header.Get(XFwdUserRoles))
	}

	// Issue request to DAO Layer
	result, err := adminDB.GetTenantDescriptor(params.TenantID)
	if err != nil {
		if strings.Contains(err.Error(), datastore.NotFoundStr) {
			return startTime, http.StatusNotFound, nil, err
		}

		return startTime, http.StatusInternalServerError, nil, fmt.Errorf("Unable to retrieve %s: %s", admmod.TenantStr, err.Error())
	}

	converted := swagmodels.TenantResponse{}
	err = convertToJsonapiObject(result, &converted)
	if err != nil {
		return startTime, http.StatusInternalServerError, nil, fmt.Errorf("Unable to convert %s data to jsonapi return format: %s", admmod.TenantStr, err.Error())
	}

	reportAPICompletionState(startTime, http.StatusOK, mon.GetTenantStr, mon.APICompleted, mon.AdminAPICompleted)
	logger.Log.Infof("Retrieved %s %s", admmod.TenantStr, models.AsJSONString(converted))
	return startTime, http.StatusOK, &converted, nil
}

func doUpdateTenantV2(allowedRoles []string, adminDB datastore.AdminServiceDatastore, params admin_provisioning_service_v2.PatchTenantV2Params) (time.Time, int, *swagmodels.TenantResponse, error) {
	isAuthorized, startTime := authorizeRequest(fmt.Sprintf("Updating %s: %s", admmod.TenantStr, params.TenantID), params.HTTPRequest, allowedRoles, mon.APIRecieved, mon.AdminAPIRecieved)

	if !isAuthorized {
		return startTime, http.StatusForbidden, nil, fmt.Errorf("Update %s operation not authorized for role: %s", admmod.TenantStr, params.HTTPRequest.Header.Get(XFwdUserRoles))
	}

	// Retrieve tne request bytes from the payload in order to convert it to a map
	patchRequestBytes, err := json.Marshal(params.Body)
	if err != nil || params.Body == nil {
		return startTime, http.StatusBadRequest, nil, fmt.Errorf("Invalid request body: %s", models.AsJSONString(params.Body))
	}

	// Attempt to retrieve the tenant that we are trying to patch from the DB in order to do a merge
	fetchedTenant, err := adminDB.GetTenantDescriptor(params.TenantID)
	if err != nil {
		if strings.Contains(err.Error(), datastore.NotFoundStr) {
			return startTime, http.StatusNotFound, nil, err
		}

		return startTime, http.StatusInternalServerError, nil, fmt.Errorf("Unable to fetch %s: %s", mon.PatchTenantStr, err.Error())
	}

	// Merge the attributes passed in with the patch request to the tenant fetched from the datastore
	patchedTenant := &admmod.Tenant{}
	if err := models.MergeObjWithMap(patchedTenant, fetchedTenant, patchRequestBytes); err != nil {
		return startTime, http.StatusInternalServerError, nil, fmt.Errorf("Unable to patch tenant with id %s: %s", params.TenantID, err.Error())
	}

	// Finally update the tenant in the datastore with the merged map and fetched tenant
	result, err := adminDB.UpdateTenantDescriptor(patchedTenant)
	if err != nil {
		if strings.Contains(err.Error(), datastore.ConflictErrorStr) {
			return startTime, http.StatusConflict, nil, err
		}

		return startTime, http.StatusInternalServerError, nil, err
	}

	converted := swagmodels.TenantResponse{}
	err = convertToJsonapiObject(result, &converted)
	if err != nil {
		return startTime, http.StatusInternalServerError, nil, fmt.Errorf("Unable to convert %s data to jsonapi return format: %s", admmod.TenantStr, err.Error())
	}

	reportAPICompletionState(startTime, http.StatusOK, mon.PatchTenantStr, mon.APICompleted, mon.AdminAPICompleted)
	logger.Log.Infof("Updated %s %s", admmod.TenantStr, models.AsJSONString(converted))
	return startTime, http.StatusOK, &converted, nil
}

func doDeleteTenantV2(allowedRoles []string, adminDB datastore.AdminServiceDatastore, params admin_provisioning_service_v2.DeleteTenantV2Params) (time.Time, int, *swagmodels.TenantResponse, error) {
	isAuthorized, startTime := authorizeRequest(fmt.Sprintf("Deleting %s: %s", admmod.TenantStr, params.TenantID), params.HTTPRequest, allowedRoles, mon.APIRecieved, mon.AdminAPIRecieved)

	if !isAuthorized {
		return startTime, http.StatusForbidden, nil, fmt.Errorf("Delete %s operation not authorized for role: %s", admmod.TenantStr, params.HTTPRequest.Header.Get(XFwdUserRoles))
	}

	// Issue request to DAO Layer
	result, err := adminDB.DeleteTenant(params.TenantID)
	if err != nil {
		if strings.Contains(err.Error(), datastore.NotFoundStr) {
			return startTime, http.StatusNotFound, nil, err
		}

		return startTime, http.StatusInternalServerError, nil, fmt.Errorf("Unable to delete %s: %s", admmod.TenantStr, err.Error())
	}

	converted := swagmodels.TenantResponse{}
	err = convertToJsonapiObject(result, &converted)
	if err != nil {
		return startTime, http.StatusInternalServerError, nil, fmt.Errorf("Unable to convert %s data to jsonapi return format: %s", admmod.TenantStr, err.Error())
	}

	reportAPICompletionState(startTime, http.StatusOK, mon.DeleteTenantStr, mon.APICompleted, mon.AdminAPICompleted)
	logger.Log.Infof("Deleted %s %s", admmod.TenantStr, models.AsJSONString(converted))
	return startTime, http.StatusOK, &converted, nil
}

func doGetAllTenantsV2(allowedRoles []string, adminDB datastore.AdminServiceDatastore, params admin_provisioning_service_v2.GetAllTenantsV2Params) (time.Time, int, *swagmodels.TenantListResponse, error) {
	isAuthorized, startTime := authorizeRequest(fmt.Sprintf("Fetching %s list", admmod.TenantStr), params.HTTPRequest, allowedRoles, mon.APIRecieved, mon.AdminAPIRecieved)

	if !isAuthorized {
		return startTime, http.StatusForbidden, nil, fmt.Errorf("Fetch %s operation not authorized for role: %s", admmod.TenantStr, params.HTTPRequest.Header.Get(XFwdUserRoles))
	}

	// Issue request to DAO Layer
	result, err := adminDB.GetAllTenantDescriptors()
	if err != nil {
		if strings.Contains(err.Error(), datastore.NotFoundStr) {
			return startTime, http.StatusNotFound, nil, err
		}

		return startTime, http.StatusInternalServerError, nil, fmt.Errorf("Unable to retrieve %s list: %s", admmod.TenantStr, err.Error())
	}

	if len(result) == 0 {
		return startTime, http.StatusNotFound, nil, fmt.Errorf(datastore.NotFoundStr)
	}

	converted := swagmodels.TenantListResponse{}
	err = convertToJsonapiObject(result, &converted)
	if err != nil {
		return startTime, http.StatusInternalServerError, nil, fmt.Errorf("Unable to convert %s data to jsonapi return format: %s", admmod.TenantStr, err.Error())
	}

	reportAPICompletionState(startTime, http.StatusOK, mon.GetAllTenantStr, mon.APICompleted, mon.AdminAPICompleted)
	logger.Log.Infof("Retrieved %d %ss", len(result), admmod.TenantStr)
	return startTime, http.StatusOK, &converted, nil
}

func doGetTenantIDByAliasV2(adminDB datastore.AdminServiceDatastore, params admin_provisioning_service_v2.GetTenantIDByAliasV2Params) (time.Time, int, string, error) {
	startTime := time.Now()
	incrementAPICounters(mon.APIRecieved, mon.AdminAPIRecieved)
	logger.Log.Infof("Fetching id for %s %s", admmod.TenantStr, params.Value)

	// Issue request to DAO Layer
	result, err := adminDB.GetTenantIDByAlias(params.Value)
	if err != nil {
		if strings.Contains(err.Error(), datastore.NotFoundStr) {
			return startTime, http.StatusNotFound, "", err
		}

		return startTime, http.StatusInternalServerError, "", fmt.Errorf("Unable to retrieve %s ID for %s: %s", admmod.TenantStr, params.Value, err.Error())
	}

	if result == "" {
		return startTime, http.StatusNotFound, "", fmt.Errorf(datastore.NotFoundStr)
	}

	reportAPICompletionState(startTime, http.StatusOK, mon.GetTenantIDByAliasStr, mon.APICompleted, mon.AdminAPICompleted)
	logger.Log.Infof("Found ID %s for %s %s", result, admmod.TenantStr, params.Value)
	return startTime, http.StatusOK, result, nil
}

func doGetTenantSummaryByAliasV2(adminDB datastore.AdminServiceDatastore, params admin_provisioning_service_v2.GetTenantSummaryByAliasV2Params) (time.Time, int, *swagmodels.TenantSummaryResponse, error) {
	startTime := time.Now()
	incrementAPICounters(mon.APIRecieved, mon.AdminAPIRecieved)
	logger.Log.Infof("Fetching summary for %s %s", admmod.TenantStr, params.Value)

	// Issue request to DAO Layer
	result, err := adminDB.GetTenantIDByAlias(params.Value)
	if err != nil {
		if strings.Contains(err.Error(), datastore.NotFoundStr) {
			return startTime, http.StatusNotFound, nil, err
		}

		return startTime, http.StatusInternalServerError, nil, fmt.Errorf("Unable to retrieve %s summary for %s: %s", admmod.TenantStr, params.Value, err.Error())
	}

	if result == "" {
		return startTime, http.StatusNotFound, nil, fmt.Errorf(datastore.NotFoundStr)
	}

	summary := swagmodels.TenantSummaryResponse{
		Data: &swagmodels.TenantSummaryResponseData{
			ID:   result,
			Type: "tenantSummaries",
			Attributes: &swagmodels.TenantSummaryResponseDataAttributes{
				Alias: params.Value,
				ID:    result,
			},
		},
	}

	reportAPICompletionState(startTime, http.StatusOK, mon.GetTenantSummaryByAliasStr, mon.APICompleted, mon.AdminAPICompleted)
	logger.Log.Infof("Successfully retrieved ID %s for alias %s", result, params.Value)
	return startTime, http.StatusOK, &summary, nil
}

func doGetIngestionDictionaryV2(allowedRoles []string, adminDB datastore.AdminServiceDatastore, params admin_provisioning_service_v2.GetIngestionDictionaryV2Params) (time.Time, int, *swagmodels.IngestionDictionaryListResponse, error) {
	isAuthorized, startTime := authorizeRequest(fmt.Sprintf("Fetching %s", admmod.IngestionDictionaryStr), params.HTTPRequest, allowedRoles, mon.APIRecieved, mon.AdminAPIRecieved)

	if !isAuthorized {
		return startTime, http.StatusForbidden, nil, fmt.Errorf("Fetch %s operation not authorized for role: %s", admmod.IngestionDictionaryStr, params.HTTPRequest.Header.Get(XFwdUserRoles))
	}

	// Issue request to DAO Layer
	result := admmod.GetIngestionDictionaryFromFile()

	resultArray := []*admmod.IngestionDictionary{result}

	converted := swagmodels.IngestionDictionaryListResponse{}
	if err := convertToJsonapiObject(resultArray, &converted); err != nil {
		return startTime, http.StatusInternalServerError, nil, fmt.Errorf("Unable to convert %s data to jsonapi return format: %s", admmod.IngestionDictionaryStr, err.Error())
	}

	reportAPICompletionState(startTime, http.StatusOK, mon.GetIngDictStr, mon.APICompleted, mon.AdminAPICompleted)
	logger.Log.Infof("Retrieved %s %s", admmod.IngestionDictionaryStr, models.AsJSONString(converted))
	return startTime, http.StatusOK, &converted, nil

}

func doGetValidTypesV2(allowedRoles []string, adminDB datastore.AdminServiceDatastore, params admin_provisioning_service_v2.GetValidTypesV2Params) (time.Time, int, *swagmodels.ValidTypesListResponse, error) {
	isAuthorized, startTime := authorizeRequest(fmt.Sprintf("Fetching %s", admmod.ValidTypesStr), params.HTTPRequest, allowedRoles, mon.APIRecieved, mon.AdminAPIRecieved)

	if !isAuthorized {
		return startTime, http.StatusForbidden, nil, fmt.Errorf("Fetch %s operation not authorized for role: %s", admmod.ValidTypesStr, params.HTTPRequest.Header.Get(XFwdUserRoles))
	}

	// Issue request to DAO Layer
	result := admmod.GetValidTypes()
	resultArray := []*admmod.ValidTypes{result}

	converted := swagmodels.ValidTypesListResponse{}
	err := convertToJsonapiObject(resultArray, &converted)
	if err != nil {
		return startTime, http.StatusInternalServerError, nil, fmt.Errorf("Unable to convert %s data to jsonapi return format: %s", admmod.ValidTypesStr, err.Error())
	}

	reportAPICompletionState(startTime, http.StatusOK, mon.GetValidTypesStr, mon.APICompleted, mon.AdminAPICompleted)
	logger.Log.Infof("Retrieved %s %s", admmod.ValidTypesStr, models.AsJSONString(converted))
	return startTime, http.StatusOK, &converted, nil
}
