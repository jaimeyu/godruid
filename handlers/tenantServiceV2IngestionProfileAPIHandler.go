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
	"github.com/accedian/adh-gather/restapi/operations/tenant_provisioning_service_v2"
	"github.com/go-openapi/runtime/middleware"
)

// HandleCreateIngestionProfileV2 - create a new Ingestion Profile
func HandleCreateIngestionProfileV2(allowedRoles []string, tenantDB datastore.TenantServiceDatastore) func(params tenant_provisioning_service_v2.CreateIngestionProfileV2Params) middleware.Responder {
	return func(params tenant_provisioning_service_v2.CreateIngestionProfileV2Params) middleware.Responder {

		// Do the work
		startTime, responseCode, response, err := doCreateIngestionProfileV2(allowedRoles, tenantDB, params)

		// Success Response
		if responseCode == http.StatusCreated {
			return tenant_provisioning_service_v2.NewCreateIngestionProfileV2Created().WithPayload(response)
		}

		// Error Responses
		errorMessage := reportAPIError(err.Error(), startTime, responseCode, mon.CreateIngPrfStr, mon.APICompleted, mon.TenantAPICompleted)
		switch responseCode {
		case http.StatusForbidden:
			return tenant_provisioning_service_v2.NewCreateIngestionProfileV2Forbidden().WithPayload(errorMessage)
		case http.StatusBadRequest:
			return tenant_provisioning_service_v2.NewCreateIngestionProfileV2BadRequest().WithPayload(errorMessage)
		case http.StatusConflict:
			return tenant_provisioning_service_v2.NewCreateIngestionProfileV2Conflict().WithPayload(errorMessage)
		default:
			return tenant_provisioning_service_v2.NewCreateIngestionProfileV2InternalServerError().WithPayload(errorMessage)
		}
	}
}

// HandleGetIngestionProfileV2 - retrieve a Ingestion Profile by the config ID.
func HandleGetIngestionProfileV2(allowedRoles []string, tenantDB datastore.TenantServiceDatastore) func(params tenant_provisioning_service_v2.GetIngestionProfileV2Params) middleware.Responder {
	return func(params tenant_provisioning_service_v2.GetIngestionProfileV2Params) middleware.Responder {

		// Do the work
		startTime, responseCode, response, err := doGetIngestionProfileV2(allowedRoles, tenantDB, params)

		// Success Response
		if responseCode == http.StatusOK {
			return tenant_provisioning_service_v2.NewGetIngestionProfileV2OK().WithPayload(response)
		}

		// Error Responses
		errorMessage := reportAPIError(err.Error(), startTime, responseCode, mon.GetIngPrfStr, mon.APICompleted, mon.TenantAPICompleted)
		switch responseCode {
		case http.StatusForbidden:
			return tenant_provisioning_service_v2.NewGetIngestionProfileV2Forbidden().WithPayload(errorMessage)
		case http.StatusNotFound:
			return tenant_provisioning_service_v2.NewGetIngestionProfileV2NotFound().WithPayload(errorMessage)
		default:
			return tenant_provisioning_service_v2.NewGetIngestionProfileV2InternalServerError().WithPayload(errorMessage)
		}

	}
}

// HandleUpdateIngestionProfileV2 - update a IngestionProfile record
func HandleUpdateIngestionProfileV2(allowedRoles []string, tenantDB datastore.TenantServiceDatastore) func(params tenant_provisioning_service_v2.UpdateIngestionProfileV2Params) middleware.Responder {
	return func(params tenant_provisioning_service_v2.UpdateIngestionProfileV2Params) middleware.Responder {

		// Do the work
		startTime, responseCode, response, err := doUpdateIngestionProfileV2(allowedRoles, tenantDB, params)

		// Success Response
		if responseCode == http.StatusOK {
			return tenant_provisioning_service_v2.NewUpdateIngestionProfileV2OK().WithPayload(response)
		}

		// Error Responses
		errorMessage := reportAPIError(err.Error(), startTime, responseCode, mon.UpdateIngPrfStr, mon.APICompleted, mon.TenantAPICompleted)
		switch responseCode {
		case http.StatusForbidden:
			return tenant_provisioning_service_v2.NewUpdateIngestionProfileV2Forbidden().WithPayload(errorMessage)
		case http.StatusBadRequest:
			return tenant_provisioning_service_v2.NewUpdateIngestionProfileV2BadRequest().WithPayload(errorMessage)
		case http.StatusNotFound:
			return tenant_provisioning_service_v2.NewUpdateIngestionProfileV2NotFound().WithPayload(errorMessage)
		case http.StatusConflict:
			return tenant_provisioning_service_v2.NewUpdateIngestionProfileV2Conflict().WithPayload(errorMessage)
		default:
			return tenant_provisioning_service_v2.NewUpdateIngestionProfileV2InternalServerError().WithPayload(errorMessage)
		}
	}
}

// HandleDeleteIngestionProfileV2 - delete a Ingestion Profile by ID
func HandleDeleteIngestionProfileV2(allowedRoles []string, tenantDB datastore.TenantServiceDatastore) func(params tenant_provisioning_service_v2.DeleteIngestionProfileV2Params) middleware.Responder {
	return func(params tenant_provisioning_service_v2.DeleteIngestionProfileV2Params) middleware.Responder {

		// Do the work
		startTime, responseCode, response, err := doDeleteIngestionProfileV2(allowedRoles, tenantDB, params)

		// Success Response
		if responseCode == http.StatusOK {
			return tenant_provisioning_service_v2.NewDeleteIngestionProfileV2OK().WithPayload(response)
		}

		// Error Responses
		errorMessage := reportAPIError(err.Error(), startTime, responseCode, mon.DeleteIngPrfStr, mon.APICompleted, mon.TenantAPICompleted)
		switch responseCode {
		case http.StatusForbidden:
			return tenant_provisioning_service_v2.NewDeleteIngestionProfileV2Forbidden().WithPayload(errorMessage)
		case http.StatusNotFound:
			return tenant_provisioning_service_v2.NewDeleteIngestionProfileV2NotFound().WithPayload(errorMessage)
		default:
			return tenant_provisioning_service_v2.NewDeleteIngestionProfileV2InternalServerError().WithPayload(errorMessage)
		}

	}
}

// HandleGetAllIngestionProfilesV2 - retrieve all Ingestion Profiles
func HandleGetAllIngestionProfilesV2(allowedRoles []string, tenantDB datastore.TenantServiceDatastore) func(params tenant_provisioning_service_v2.GetAllIngestionProfilesV2Params) middleware.Responder {
	return func(params tenant_provisioning_service_v2.GetAllIngestionProfilesV2Params) middleware.Responder {

		// Do the work
		startTime, responseCode, response, err := doGetAllIngestionProfilesV2(allowedRoles, tenantDB, params)

		// Success Response
		if responseCode == http.StatusOK {
			return tenant_provisioning_service_v2.NewGetAllIngestionProfilesV2OK().WithPayload(response)
		}

		// Error Responses
		errorMessage := reportAPIError(err.Error(), startTime, responseCode, mon.GetAllIngPrfStr, mon.APICompleted, mon.TenantAPICompleted)
		switch responseCode {
		case http.StatusForbidden:
			return tenant_provisioning_service_v2.NewGetAllIngestionProfilesV2Forbidden().WithPayload(errorMessage)
		case http.StatusNotFound:
			return tenant_provisioning_service_v2.NewGetAllIngestionProfilesV2NotFound().WithPayload(errorMessage)
		default:
			return tenant_provisioning_service_v2.NewGetAllIngestionProfilesV2InternalServerError().WithPayload(errorMessage)
		}

	}
}

func doCreateIngestionProfileV2(allowedRoles []string, tenantDB datastore.TenantServiceDatastore, params tenant_provisioning_service_v2.CreateIngestionProfileV2Params) (time.Time, int, *swagmodels.IngestionProfileResponse, error) {
	tenantID := params.HTTPRequest.Header.Get(XFwdTenantId)
	isAuthorized, startTime := authorizeRequest(fmt.Sprintf("Creating %s for tenant %s", tenmod.TenantIngestionProfileStr, tenantID), params.HTTPRequest, allowedRoles, mon.APIRecieved, mon.TenantAPIRecieved)

	if !isAuthorized {
		return startTime, http.StatusForbidden, nil, fmt.Errorf("Create %s operation not authorized for role: %s", tenmod.TenantIngestionProfileStr, params.HTTPRequest.Header.Get(XFwdUserRoles))
	}

	// Unmarshal the request
	requestBytes, err := json.Marshal(params.Body)
	if err != nil {
		return startTime, http.StatusBadRequest, nil, err
	}

	data := tenmod.IngestionProfile{}
	err = jsonapi.Unmarshal(requestBytes, &data)
	if err != nil {
		return startTime, http.StatusBadRequest, nil, err
	}

	data.TenantID = tenantID

	// Issue request to DAO Layer to Create Record
	result, err := tenantDB.CreateTenantIngestionProfile(&data)
	if err != nil {
		if strings.Contains(err.Error(), datastore.ConflictStr) {
			return startTime, http.StatusConflict, nil, err
		}
		return startTime, http.StatusInternalServerError, nil, fmt.Errorf("Unable to store %s: %s", tenmod.TenantIngestionProfileStr, err.Error())
	}

	converted := swagmodels.IngestionProfileResponse{}
	err = convertToJsonapiObject(result, &converted)
	if err != nil {
		return startTime, http.StatusInternalServerError, nil, fmt.Errorf("Unable to convert %s data to jsonapi return format: %s", tenmod.TenantIngestionProfileStr, err.Error())
	}

	reportAPICompletionState(startTime, http.StatusOK, mon.CreateIngPrfStr, mon.APICompleted, mon.TenantAPICompleted)
	logger.Log.Infof("Created %s %s", tenmod.TenantIngestionProfileStr, models.AsJSONString(converted))
	return startTime, http.StatusCreated, &converted, nil
}

func doGetIngestionProfileV2(allowedRoles []string, tenantDB datastore.TenantServiceDatastore, params tenant_provisioning_service_v2.GetIngestionProfileV2Params) (time.Time, int, *swagmodels.IngestionProfileResponse, error) {
	tenantID := params.HTTPRequest.Header.Get(XFwdTenantId)
	isAuthorized, startTime := authorizeRequest(fmt.Sprintf("Fetching %s %s for %s %s", tenmod.TenantIngestionProfileStr, params.IngestionProfileID, admmod.TenantStr, tenantID), params.HTTPRequest, allowedRoles, mon.APIRecieved, mon.AdminAPIRecieved)

	if !isAuthorized {
		return startTime, http.StatusForbidden, nil, fmt.Errorf("Fetch %s operation not authorized for role: %s", tenmod.TenantIngestionProfileStr, params.HTTPRequest.Header.Get(XFwdUserRoles))
	}

	// Issue request to DAO Layer
	result, err := tenantDB.GetTenantIngestionProfile(tenantID, params.IngestionProfileID)
	if err != nil {
		if strings.Contains(err.Error(), datastore.NotFoundStr) {
			return startTime, http.StatusNotFound, nil, err
		}

		return startTime, http.StatusInternalServerError, nil, fmt.Errorf("Unable to retrieve %s: %s", tenmod.TenantIngestionProfileStr, err.Error())
	}

	converted := swagmodels.IngestionProfileResponse{}
	err = convertToJsonapiObject(result, &converted)
	if err != nil {
		return startTime, http.StatusInternalServerError, nil, fmt.Errorf("Unable to convert %s data to jsonapi return format: %s", tenmod.TenantIngestionProfileStr, err.Error())
	}

	reportAPICompletionState(startTime, http.StatusOK, mon.GetIngPrfStr, mon.APICompleted, mon.TenantAPICompleted)
	logger.Log.Infof("Retrieved %s %s", tenmod.TenantIngestionProfileStr, models.AsJSONString(converted))
	return startTime, http.StatusOK, &converted, nil
}

func doUpdateIngestionProfileV2(allowedRoles []string, tenantDB datastore.TenantServiceDatastore, params tenant_provisioning_service_v2.UpdateIngestionProfileV2Params) (time.Time, int, *swagmodels.IngestionProfileResponse, error) {
	tenantID := params.HTTPRequest.Header.Get(XFwdTenantId)
	isAuthorized, startTime := authorizeRequest(fmt.Sprintf("Updating %s %s for %s %s", tenmod.TenantIngestionProfileStr, params.IngestionProfileID, admmod.TenantStr, tenantID), params.HTTPRequest, allowedRoles, mon.APIRecieved, mon.AdminAPIRecieved)

	if !isAuthorized {
		return startTime, http.StatusForbidden, nil, fmt.Errorf("Update %s operation not authorized for role: %s", tenmod.TenantIngestionProfileStr, params.HTTPRequest.Header.Get(XFwdUserRoles))
	}

	// Retrieve tne request bytes from the payload in order to convert it to a map
	patchRequestBytes, err := json.Marshal(params.Body)
	if err != nil || params.Body == nil {
		return startTime, http.StatusBadRequest, nil, fmt.Errorf("Invalid request body: %s", models.AsJSONString(params.Body))
	}

	// Attempt to retrieve the record that we are trying to patch from the DB in order to do a merge
	fetched, err := tenantDB.GetTenantIngestionProfile(tenantID, params.IngestionProfileID)
	if err != nil {
		if strings.Contains(err.Error(), datastore.NotFoundStr) {
			return startTime, http.StatusNotFound, nil, err
		}

		return startTime, http.StatusInternalServerError, nil, fmt.Errorf("Unable to fetch %s: %s", tenmod.TenantIngestionProfileStr, err.Error())
	}

	// Merge the attributes passed in with the patch request to the record fetched from the datastore
	var patched *tenmod.IngestionProfile
	if err := models.MergeObjWithMap(fetched, patchRequestBytes); err != nil {
		return startTime, http.StatusInternalServerError, nil, fmt.Errorf("Unable to patch %s with id %s: %s", tenmod.TenantIngestionProfileStr, params.IngestionProfileID, err.Error())
	}
	patched = fetched
	patched.TenantID = tenantID

	// Finally update the record in the datastore with the merged map and fetched tenant
	result, err := tenantDB.UpdateTenantIngestionProfile(patched)
	if err != nil {
		if strings.Contains(err.Error(), datastore.ConflictErrorStr) {
			return startTime, http.StatusConflict, nil, err
		}

		return startTime, http.StatusInternalServerError, nil, err
	}

	converted := swagmodels.IngestionProfileResponse{}
	err = convertToJsonapiObject(result, &converted)
	if err != nil {
		return startTime, http.StatusInternalServerError, nil, fmt.Errorf("Unable to convert %s data to jsonapi return format: %s", tenmod.TenantIngestionProfileStr, err.Error())
	}

	reportAPICompletionState(startTime, http.StatusOK, mon.UpdateIngPrfStr, mon.APICompleted, mon.TenantAPICompleted)
	logger.Log.Infof("Updated %s %s", tenmod.TenantIngestionProfileStr, models.AsJSONString(converted))
	return startTime, http.StatusOK, &converted, nil
}

func doDeleteIngestionProfileV2(allowedRoles []string, tenantDB datastore.TenantServiceDatastore, params tenant_provisioning_service_v2.DeleteIngestionProfileV2Params) (time.Time, int, *swagmodels.IngestionProfileResponse, error) {
	tenantID := params.HTTPRequest.Header.Get(XFwdTenantId)
	isAuthorized, startTime := authorizeRequest(fmt.Sprintf("Deleting %s %s for %s %s", tenmod.TenantIngestionProfileStr, params.IngestionProfileID, admmod.TenantStr, tenantID), params.HTTPRequest, allowedRoles, mon.APIRecieved, mon.AdminAPIRecieved)

	if !isAuthorized {
		return startTime, http.StatusForbidden, nil, fmt.Errorf("Delete %s operation not authorized for role: %s", tenmod.TenantIngestionProfileStr, params.HTTPRequest.Header.Get(XFwdUserRoles))
	}

	// Issue request to DAO Layer
	result, err := tenantDB.DeleteTenantIngestionProfile(tenantID, params.IngestionProfileID)
	if err != nil {
		if strings.Contains(err.Error(), datastore.NotFoundStr) {
			return startTime, http.StatusNotFound, nil, err
		}

		return startTime, http.StatusInternalServerError, nil, fmt.Errorf("Unable to delete %s: %s", tenmod.TenantIngestionProfileStr, err.Error())
	}

	converted := swagmodels.IngestionProfileResponse{}
	err = convertToJsonapiObject(result, &converted)
	if err != nil {
		return startTime, http.StatusInternalServerError, nil, fmt.Errorf("Unable to convert %s data to jsonapi return format: %s", tenmod.TenantIngestionProfileStr, err.Error())
	}

	reportAPICompletionState(startTime, http.StatusOK, mon.DeleteIngPrfStr, mon.APICompleted, mon.TenantAPICompleted)
	logger.Log.Infof("Deleted %s %s", tenmod.TenantIngestionProfileStr, models.AsJSONString(converted))
	return startTime, http.StatusOK, &converted, nil
}

func doGetAllIngestionProfilesV2(allowedRoles []string, tenantDB datastore.TenantServiceDatastore, params tenant_provisioning_service_v2.GetAllIngestionProfilesV2Params) (time.Time, int, *swagmodels.IngestionProfileListResponse, error) {
	tenantID := params.HTTPRequest.Header.Get(XFwdTenantId)
	isAuthorized, startTime := authorizeRequest(fmt.Sprintf("Fetching %s list for %s %s", tenmod.TenantIngestionProfileStr, admmod.TenantStr, tenantID), params.HTTPRequest, allowedRoles, mon.APIRecieved, mon.AdminAPIRecieved)

	if !isAuthorized {
		return startTime, http.StatusForbidden, nil, fmt.Errorf("Fetch %s operation not authorized for role: %s", tenmod.TenantIngestionProfileStr, params.HTTPRequest.Header.Get(XFwdUserRoles))
	}

	// Issue request to DAO Layer
	result, err := tenantDB.GetActiveTenantIngestionProfile(tenantID)
	if err != nil {
		if strings.Contains(err.Error(), datastore.NotFoundStr) {
			return startTime, http.StatusNotFound, nil, err
		}

		return startTime, http.StatusInternalServerError, nil, fmt.Errorf("Unable to retrieve %s list: %s", tenmod.TenantIngestionProfileStr, err.Error())
	}

	// Add to an array for proper result
	resultArray := []*tenmod.IngestionProfile{result}

	converted := swagmodels.IngestionProfileListResponse{}
	err = convertToJsonapiObject(resultArray, &converted)
	if err != nil {
		return startTime, http.StatusInternalServerError, nil, fmt.Errorf("Unable to convert %s data to jsonapi return format: %s", tenmod.TenantIngestionProfileStr, err.Error())
	}

	reportAPICompletionState(startTime, http.StatusOK, mon.GetAllTenantConnectorInstanceStr, mon.APICompleted, mon.TenantAPICompleted)
	logger.Log.Infof("Retrieved %d %ss", len(resultArray), tenmod.TenantIngestionProfileStr)
	return startTime, http.StatusOK, &converted, nil
}
