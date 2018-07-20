package handlers

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/accedian/adh-gather/datastore"
	"github.com/accedian/adh-gather/logger"
	"github.com/accedian/adh-gather/models"
	tenmod "github.com/accedian/adh-gather/models/tenant"
	mon "github.com/accedian/adh-gather/monitoring"
	"github.com/accedian/adh-gather/restapi/operations/tenant_provisioning_service_v2"
	"github.com/accedian/adh-gather/swagmodels"
	"github.com/go-openapi/runtime/middleware"
)

// HandleGetDataCleaningProfileV2 - retrieve the Data Cleaning Profile for a Tenant
func HandleGetDataCleaningProfileV2(allowedRoles []string, tenantDB datastore.TenantServiceDatastore) func(params tenant_provisioning_service_v2.GetDataCleaningProfileParams) middleware.Responder {
	return func(params tenant_provisioning_service_v2.GetDataCleaningProfileParams) middleware.Responder {
		// Do the work
		startTime, responseCode, response, err := doGetDataCleaningProfileV2(allowedRoles, tenantDB, params)

		// Success Response
		if responseCode == http.StatusOK {
			return tenant_provisioning_service_v2.NewGetDataCleaningProfileOK().WithPayload(response)
		}

		// Error Responses
		errorMessage := reportAPIError(err.Error(), startTime, responseCode, mon.GetDataCleaningProfileStr, mon.APICompleted, mon.TenantAPICompleted)
		switch responseCode {
		case http.StatusForbidden:
			return tenant_provisioning_service_v2.NewGetDataCleaningProfileForbidden().WithPayload(errorMessage)
		case http.StatusNotFound:
			return tenant_provisioning_service_v2.NewGetDataCleaningProfileNotFound().WithPayload(errorMessage)
		default:
			return tenant_provisioning_service_v2.NewGetDataCleaningProfileInternalServerError().WithPayload(errorMessage)

		}
	}
}

// HandleGetDataCleaningProfilesV2 - retrieve all Data Cleaning Profile for a Tenant
func HandleGetDataCleaningProfilesV2(allowedRoles []string, tenantDB datastore.TenantServiceDatastore) func(params tenant_provisioning_service_v2.GetDataCleaningProfilesParams) middleware.Responder {
	return func(params tenant_provisioning_service_v2.GetDataCleaningProfilesParams) middleware.Responder {
		// Do the work
		startTime, responseCode, response, err := doGetDataCleaningProfilesV2(allowedRoles, tenantDB, params)

		// Success Response
		if responseCode == http.StatusOK {
			return tenant_provisioning_service_v2.NewGetDataCleaningProfilesOK().WithPayload(response)
		}

		// Error Responses
		errorMessage := reportAPIError(err.Error(), startTime, responseCode, mon.GetAllDataCleaningProfileStr, mon.APICompleted, mon.TenantAPICompleted)
		switch responseCode {
		case http.StatusForbidden:
			return tenant_provisioning_service_v2.NewGetDataCleaningProfilesForbidden().WithPayload(errorMessage)
		case http.StatusNotFound:
			return tenant_provisioning_service_v2.NewGetDataCleaningProfilesNotFound().WithPayload(errorMessage)
		default:
			return tenant_provisioning_service_v2.NewGetDataCleaningProfilesInternalServerError().WithPayload(errorMessage)

		}
	}
}

// HandleDeleteDataCleaningProfileV2 - retrieve the Data Cleaning Profile for a Tenant
func HandleDeleteDataCleaningProfileV2(allowedRoles []string, tenantDB datastore.TenantServiceDatastore) func(params tenant_provisioning_service_v2.DeleteDataCleaningProfileParams) middleware.Responder {
	return func(params tenant_provisioning_service_v2.DeleteDataCleaningProfileParams) middleware.Responder {
		// Do the work
		startTime, responseCode, response, err := doDeleteDataCleaningProfileV2(allowedRoles, tenantDB, params)

		// Success Response
		if responseCode == http.StatusOK {
			return tenant_provisioning_service_v2.NewDeleteDataCleaningProfileOK().WithPayload(response)
		}

		// Error Responses
		errorMessage := reportAPIError(err.Error(), startTime, responseCode, mon.DeleteDataCleaningProfileStr, mon.APICompleted, mon.TenantAPICompleted)
		switch responseCode {
		case http.StatusForbidden:
			return tenant_provisioning_service_v2.NewDeleteDataCleaningProfileForbidden().WithPayload(errorMessage)
		case http.StatusNotFound:
			return tenant_provisioning_service_v2.NewDeleteDataCleaningProfileNotFound().WithPayload(errorMessage)
		default:
			return tenant_provisioning_service_v2.NewDeleteDataCleaningProfileInternalServerError().WithPayload(errorMessage)

		}
	}
}

// HandleUpdateDataCleaningProfileV2 - update the Data Cleaning Profile for a Tenant
func HandleUpdateDataCleaningProfileV2(allowedRoles []string, tenantDB datastore.TenantServiceDatastore) func(params tenant_provisioning_service_v2.UpdateDataCleaningProfileParams) middleware.Responder {
	return func(params tenant_provisioning_service_v2.UpdateDataCleaningProfileParams) middleware.Responder {
		// Do the work
		startTime, responseCode, response, err := doUpdateDataCleaningProfileV2(allowedRoles, tenantDB, params)

		// Success Response
		if responseCode == http.StatusOK {
			return tenant_provisioning_service_v2.NewUpdateDataCleaningProfileOK().WithPayload(response)
		}

		// Error Responses
		errorMessage := reportAPIError(err.Error(), startTime, responseCode, mon.UpdateDataCleaningProfileStr, mon.APICompleted, mon.TenantAPICompleted)
		switch responseCode {
		case http.StatusForbidden:
			return tenant_provisioning_service_v2.NewUpdateDataCleaningProfileForbidden().WithPayload(errorMessage)
		case http.StatusNotFound:
			return tenant_provisioning_service_v2.NewUpdateDataCleaningProfileNotFound().WithPayload(errorMessage)
		case http.StatusBadRequest:
			return tenant_provisioning_service_v2.NewUpdateDataCleaningProfileBadRequest().WithPayload(errorMessage)
		case http.StatusConflict:
			return tenant_provisioning_service_v2.NewUpdateDataCleaningProfileConflict().WithPayload(errorMessage)
		default:
			return tenant_provisioning_service_v2.NewUpdateDataCleaningProfileInternalServerError().WithPayload(errorMessage)

		}
	}
}

// HandleCreateDataCleaningProfileV2 - update the Data Cleaning Profile for a Tenant
func HandleCreateDataCleaningProfileV2(allowedRoles []string, tenantDB datastore.TenantServiceDatastore) func(params tenant_provisioning_service_v2.CreateDataCleaningProfileParams) middleware.Responder {
	return func(params tenant_provisioning_service_v2.CreateDataCleaningProfileParams) middleware.Responder {
		// Do the work
		startTime, responseCode, response, err := doCreateDataCleaningProfileV2(allowedRoles, tenantDB, params)

		// Success Response
		if responseCode == http.StatusCreated {
			return tenant_provisioning_service_v2.NewCreateDataCleaningProfileCreated().WithPayload(response)
		}

		// Error Responses
		errorMessage := reportAPIError(err.Error(), startTime, responseCode, mon.CreateDataCleaningProfileStr, mon.APICompleted, mon.TenantAPICompleted)
		switch responseCode {
		case http.StatusForbidden:
			return tenant_provisioning_service_v2.NewCreateDataCleaningProfileForbidden().WithPayload(errorMessage)
		case http.StatusBadRequest:
			return tenant_provisioning_service_v2.NewCreateDataCleaningProfileBadRequest().WithPayload(errorMessage)
		case http.StatusConflict:
			return tenant_provisioning_service_v2.NewCreateDataCleaningProfileConflict().WithPayload(errorMessage)
		default:
			return tenant_provisioning_service_v2.NewCreateDataCleaningProfileInternalServerError().WithPayload(errorMessage)

		}
	}
}

func doGetDataCleaningProfileV2(allowedRoles []string, tenantDB datastore.TenantServiceDatastore, params tenant_provisioning_service_v2.GetDataCleaningProfileParams) (time.Time, int, *swagmodels.DataCleaningProfileResponse, error) {
	tenantID := params.HTTPRequest.Header.Get(XFwdTenantId)
	isAuthorized, startTime := authorizeRequest(fmt.Sprintf("Fetching %s: %s", tenmod.TenantDataCleaningProfileStr, tenantID), params.HTTPRequest, allowedRoles, mon.APIRecieved, mon.TenantAPIRecieved)

	if !isAuthorized {
		err := fmt.Errorf("Get %s operation not authorized for role: %s", tenmod.TenantDataCleaningProfileStr, params.HTTPRequest.Header.Get(XFwdUserRoles))
		return startTime, http.StatusForbidden, nil, err
	}

	// Issue request to DAO Layer
	result, err := tenantDB.GetTenantDataCleaningProfile(tenantID, params.ProfileID)
	if err != nil {
		if checkForNotFound(err.Error()) {
			return startTime, http.StatusNotFound, nil, err
		}

		errResp := fmt.Errorf("Unable to retrieve %s: %s", tenmod.TenantDataCleaningProfileStr, err.Error())
		return startTime, http.StatusInternalServerError, nil, errResp
	}

	converted := swagmodels.DataCleaningProfileResponse{}
	err = convertToJsonapiObject(result, &converted)
	if err != nil {
		errResp := fmt.Errorf("Unable to convert %s data to jsonapi return format: %s", tenmod.TenantDataCleaningProfileStr, err.Error())
		return startTime, http.StatusInternalServerError, nil, errResp
	}

	reportAPICompletionState(startTime, http.StatusOK, mon.GetDataCleaningProfileStr, mon.APICompleted, mon.TenantAPICompleted)
	logger.Log.Infof("Retrieved %s %s", tenmod.TenantDataCleaningProfileStr, models.AsJSONString(converted))
	return startTime, http.StatusOK, &converted, nil
}

func doGetDataCleaningProfilesV2(allowedRoles []string, tenantDB datastore.TenantServiceDatastore, params tenant_provisioning_service_v2.GetDataCleaningProfilesParams) (time.Time, int, *swagmodels.DataCleaningProfileListResponse, error) {
	tenantID := params.HTTPRequest.Header.Get(XFwdTenantId)
	isAuthorized, startTime := authorizeRequest(fmt.Sprintf("Fetching all %s: %s", tenmod.TenantDataCleaningProfileStr, tenantID), params.HTTPRequest, allowedRoles, mon.APIRecieved, mon.TenantAPIRecieved)

	if !isAuthorized {
		err := fmt.Errorf("Get all %s operation not authorized for role: %s", tenmod.TenantDataCleaningProfileStr, params.HTTPRequest.Header.Get(XFwdUserRoles))
		return startTime, http.StatusForbidden, nil, err
	}

	// Issue request to DAO Layer
	result, err := tenantDB.GetAllTenantDataCleaningProfiles(tenantID)
	if err != nil {
		if checkForNotFound(err.Error()) {
			return startTime, http.StatusNotFound, nil, err
		}

		errResp := fmt.Errorf("Unable to retrieve %s: %s", tenmod.TenantDataCleaningProfileStr, err.Error())
		return startTime, http.StatusInternalServerError, nil, errResp
	}

	converted := swagmodels.DataCleaningProfileListResponse{}
	err = convertToJsonapiObject(result, &converted)
	if err != nil {
		errResp := fmt.Errorf("Unable to convert %s data to jsonapi return format: %s", tenmod.TenantDataCleaningProfileStr, err.Error())
		return startTime, http.StatusInternalServerError, nil, errResp
	}

	reportAPICompletionState(startTime, http.StatusOK, mon.GetAllDataCleaningProfileStr, mon.APICompleted, mon.TenantAPICompleted)
	logger.Log.Infof("Retrieved %d %ss", len(converted.Data), tenmod.TenantDataCleaningProfileStr)
	return startTime, http.StatusOK, &converted, nil
}

func doDeleteDataCleaningProfileV2(allowedRoles []string, tenantDB datastore.TenantServiceDatastore, params tenant_provisioning_service_v2.DeleteDataCleaningProfileParams) (time.Time, int, *swagmodels.DataCleaningProfileResponse, error) {
	tenantID := params.HTTPRequest.Header.Get(XFwdTenantId)
	isAuthorized, startTime := authorizeRequest(fmt.Sprintf("Deleting %s: %s", tenmod.TenantDataCleaningProfileStr, tenantID), params.HTTPRequest, allowedRoles, mon.APIRecieved, mon.TenantAPIRecieved)

	if !isAuthorized {
		err := fmt.Errorf("Delete %s operation not authorized for role: %s", tenmod.TenantDataCleaningProfileStr, params.HTTPRequest.Header.Get(XFwdUserRoles))
		return startTime, http.StatusForbidden, nil, err
	}

	// Issue request to DAO Layer
	result, err := tenantDB.DeleteTenantDataCleaningProfile(tenantID, params.ProfileID)
	if err != nil {
		if checkForNotFound(err.Error()) {
			return startTime, http.StatusNotFound, nil, err
		}

		errResp := fmt.Errorf("Unable to delete %s: %s", tenmod.TenantDataCleaningProfileStr, err.Error())
		return startTime, http.StatusInternalServerError, nil, errResp
	}

	converted := swagmodels.DataCleaningProfileResponse{}
	err = convertToJsonapiObject(result, &converted)
	if err != nil {
		errResp := fmt.Errorf("Unable to convert %s data to jsonapi return format: %s", tenmod.TenantDataCleaningProfileStr, err.Error())
		return startTime, http.StatusInternalServerError, nil, errResp
	}

	reportAPICompletionState(startTime, http.StatusOK, mon.DeleteDataCleaningProfileStr, mon.APICompleted, mon.TenantAPICompleted)
	logger.Log.Infof("Deleted %s %s", tenmod.TenantDataCleaningProfileStr, models.AsJSONString(converted))
	return startTime, http.StatusOK, &converted, nil
}

func doUpdateDataCleaningProfileV2(allowedRoles []string, tenantDB datastore.TenantServiceDatastore, params tenant_provisioning_service_v2.UpdateDataCleaningProfileParams) (time.Time, int, *swagmodels.DataCleaningProfileResponse, error) {
	tenantID := params.HTTPRequest.Header.Get(XFwdTenantId)
	isAuthorized, startTime := authorizeRequest(fmt.Sprintf("Updating %s: %s", tenmod.TenantDataCleaningProfileStr, tenantID), params.HTTPRequest, allowedRoles, mon.APIRecieved, mon.TenantAPIRecieved)

	if !isAuthorized {
		err := fmt.Errorf("Update %s operation not authorized for role: %s", tenmod.TenantDataCleaningProfileStr, params.HTTPRequest.Header.Get(XFwdUserRoles))
		return startTime, http.StatusForbidden, nil, err
	}

	// Fetch the existing record
	existing, err := tenantDB.GetTenantDataCleaningProfile(tenantID, params.ProfileID)
	if err != nil {
		errResp := fmt.Errorf("Unable to update %s: %s", tenmod.TenantDataCleaningProfileStr, err.Error())
		return startTime, http.StatusConflict, nil, errResp
	}

	// Convert the request to a db model type:
	data := tenmod.DataCleaningProfile{}
	err = convertRequestBodyToDBModel(params.Body, &data)
	if err != nil {
		return startTime, http.StatusBadRequest, nil, err
	}

	existing.Rules = data.Rules

	// Issue request to DAO Layer
	result, err := tenantDB.UpdateTenantDataCleaningProfile(existing)
	if err != nil {
		if checkForNotFound(err.Error()) {
			return startTime, http.StatusNotFound, nil, err
		}
		errResp := fmt.Errorf("Unable to update %s: %s", tenmod.TenantDataCleaningProfileStr, err.Error())
		return startTime, http.StatusInternalServerError, nil, errResp
	}

	converted := swagmodels.DataCleaningProfileResponse{}
	err = convertToJsonapiObject(result, &converted)
	if err != nil {
		errResp := fmt.Errorf("Unable to convert %s data to jsonapi return format: %s", tenmod.TenantDataCleaningProfileStr, err.Error())
		return startTime, http.StatusInternalServerError, nil, errResp
	}

	reportAPICompletionState(startTime, http.StatusOK, mon.UpdateDataCleaningProfileStr, mon.APICompleted, mon.TenantAPICompleted)
	logger.Log.Infof("Updated %s %s", tenmod.TenantDataCleaningProfileStr, models.AsJSONString(converted))
	return startTime, http.StatusOK, &converted, nil
}

func doCreateDataCleaningProfileV2(allowedRoles []string, tenantDB datastore.TenantServiceDatastore, params tenant_provisioning_service_v2.CreateDataCleaningProfileParams) (time.Time, int, *swagmodels.DataCleaningProfileResponse, error) {
	tenantID := params.HTTPRequest.Header.Get(XFwdTenantId)
	isAuthorized, startTime := authorizeRequest(fmt.Sprintf("Creating %s: %s", tenmod.TenantDataCleaningProfileStr, tenantID), params.HTTPRequest, allowedRoles, mon.APIRecieved, mon.TenantAPIRecieved)

	if !isAuthorized {
		return startTime, http.StatusForbidden, nil, fmt.Errorf("Update %s operation not authorized for role: %s", tenmod.TenantDataCleaningProfileStr, params.HTTPRequest.Header.Get(XFwdUserRoles))
	}

	// Convert the request to a db model type:
	data := tenmod.DataCleaningProfile{}
	err := convertRequestBodyToDBModel(params.Body, &data)
	if err != nil {
		return startTime, http.StatusBadRequest, nil, err
	}

	data.TenantID = tenantID

	// Issue request to DAO Layer
	result, err := tenantDB.CreateTenantDataCleaningProfile(&data)
	if err != nil {
		if strings.Contains(err.Error(), string(conflict)) {
			return startTime, http.StatusConflict, nil, err
		}
		errResp := fmt.Errorf("Unable to create %s: %s", tenmod.TenantDataCleaningProfileStr, err.Error())
		return startTime, http.StatusInternalServerError, nil, errResp
	}

	converted := swagmodels.DataCleaningProfileResponse{}
	err = convertToJsonapiObject(result, &converted)
	if err != nil {
		errResp := fmt.Errorf("Unable to convert %s data to jsonapi return format: %s", tenmod.TenantDataCleaningProfileStr, err.Error())
		return startTime, http.StatusInternalServerError, nil, errResp
	}

	reportAPICompletionState(startTime, http.StatusOK, mon.CreateDataCleaningProfileStr, mon.APICompleted, mon.TenantAPICompleted)
	logger.Log.Infof("Created %s %s", tenmod.TenantDataCleaningProfileStr, models.AsJSONString(converted))
	return startTime, http.StatusCreated, &converted, nil
}
