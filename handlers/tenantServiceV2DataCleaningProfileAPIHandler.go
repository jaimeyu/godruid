package handlers

import (
	"fmt"
	"net/http"

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
		tenantID := params.HTTPRequest.Header.Get(xFwdTenantId)
		isAuthorized, startTime := authorizeRequest(fmt.Sprintf("Fetching %s: %s", tenmod.TenantDataCleaningProfileStr, tenantID), params.HTTPRequest, allowedRoles, mon.APIRecieved, mon.TenantAPIRecieved)

		if !isAuthorized {
			return tenant_provisioning_service_v2.NewGetDataCleaningProfileForbidden().WithPayload(reportAPIError(fmt.Sprintf("Get %s operation not authorized for role: %s", tenmod.TenantDataCleaningProfileStr, params.HTTPRequest.Header.Get(xFwdUserRoles)), startTime, http.StatusForbidden, mon.GetDataCleaningProfileStr, mon.APICompleted, mon.TenantAPICompleted))
		}

		// Issue request to DAO Layer
		result, err := tenantDB.GetTenantDataCleaningProfile(tenantID, params.ProfileID)
		if err != nil {
			if checkForNotFound(err.Error()) {
				return tenant_provisioning_service_v2.NewGetDataCleaningProfileNotFound().WithPayload(reportAPIError(generateErrorMessage(http.StatusNotFound, err.Error()), startTime, http.StatusBadRequest, mon.GetDataCleaningProfileStr, mon.APICompleted, mon.TenantAPICompleted))
			}
			return tenant_provisioning_service_v2.NewGetDataCleaningProfileInternalServerError().WithPayload(reportAPIError(fmt.Sprintf("Unable to retrieve %s: %s", tenmod.TenantDataCleaningProfileStr, err.Error()), startTime, http.StatusInternalServerError, mon.GetDataCleaningProfileStr, mon.APICompleted, mon.TenantAPICompleted))
		}

		converted := swagmodels.DataCleaningProfileResponse{}
		err = convertToJsonapiObject(result, &converted)
		if err != nil {
			return tenant_provisioning_service_v2.NewGetDataCleaningProfileInternalServerError().WithPayload(reportAPIError(fmt.Sprintf("Unable to convert %s data to jsonapi return format: %s", tenmod.TenantDataCleaningProfileStr, err.Error()), startTime, http.StatusInternalServerError, mon.GetDataCleaningProfileStr, mon.APICompleted, mon.TenantAPICompleted))
		}

		reportAPICompletionState(startTime, http.StatusOK, mon.GetDataCleaningProfileStr, mon.APICompleted, mon.TenantAPICompleted)
		logger.Log.Infof("Retrieved %s %s", tenmod.TenantDataCleaningProfileStr, models.AsJSONString(converted))
		return tenant_provisioning_service_v2.NewGetDataCleaningProfileOK().WithPayload(&converted)
	}
}

// HandleGetDataCleaningProfilesV2 - retrieve all Data Cleaning Profile for a Tenant
func HandleGetDataCleaningProfilesV2(allowedRoles []string, tenantDB datastore.TenantServiceDatastore) func(params tenant_provisioning_service_v2.GetDataCleaningProfilesParams) middleware.Responder {
	return func(params tenant_provisioning_service_v2.GetDataCleaningProfilesParams) middleware.Responder {
		tenantID := params.HTTPRequest.Header.Get(xFwdTenantId)
		isAuthorized, startTime := authorizeRequest(fmt.Sprintf("Fetching all %s: %s", tenmod.TenantDataCleaningProfileStr, tenantID), params.HTTPRequest, allowedRoles, mon.APIRecieved, mon.TenantAPIRecieved)

		if !isAuthorized {
			return tenant_provisioning_service_v2.NewGetDataCleaningProfilesForbidden().WithPayload(reportAPIError(fmt.Sprintf("Get all %s operation not authorized for role: %s", tenmod.TenantDataCleaningProfileStr, params.HTTPRequest.Header.Get(xFwdUserRoles)), startTime, http.StatusForbidden, mon.GetAllDataCleaningProfileStr, mon.APICompleted, mon.TenantAPICompleted))
		}

		// Issue request to DAO Layer
		result, err := tenantDB.GetAllTenantDataCleaningProfiles(tenantID)
		if err != nil {
			if checkForNotFound(err.Error()) {
				return tenant_provisioning_service_v2.NewGetDataCleaningProfilesNotFound().WithPayload(reportAPIError(generateErrorMessage(http.StatusNotFound, err.Error()), startTime, http.StatusBadRequest, mon.GetAllDataCleaningProfileStr, mon.APICompleted, mon.TenantAPICompleted))
			}
			return tenant_provisioning_service_v2.NewGetDataCleaningProfilesInternalServerError().WithPayload(reportAPIError(fmt.Sprintf("Unable to retrieve %s: %s", tenmod.TenantDataCleaningProfileStr, err.Error()), startTime, http.StatusInternalServerError, mon.GetAllDataCleaningProfileStr, mon.APICompleted, mon.TenantAPICompleted))
		}

		converted := swagmodels.DataCleaningProfileListResponse{}
		err = convertToJsonapiObject(result, &converted)
		if err != nil {
			return tenant_provisioning_service_v2.NewGetDataCleaningProfilesInternalServerError().WithPayload(reportAPIError(fmt.Sprintf("Unable to convert %s data to jsonapi return format: %s", tenmod.TenantDataCleaningProfileStr, err.Error()), startTime, http.StatusInternalServerError, mon.GetAllDataCleaningProfileStr, mon.APICompleted, mon.TenantAPICompleted))
		}

		reportAPICompletionState(startTime, http.StatusOK, mon.GetAllDataCleaningProfileStr, mon.APICompleted, mon.TenantAPICompleted)
		logger.Log.Infof("Retrieved %d %ss", len(converted.Data), tenmod.TenantDataCleaningProfileStr)
		return tenant_provisioning_service_v2.NewGetDataCleaningProfilesOK().WithPayload(&converted)
	}
}

// HandleDeleteDataCleaningProfileV2 - retrieve the Data Cleaning Profile for a Tenant
func HandleDeleteDataCleaningProfileV2(allowedRoles []string, tenantDB datastore.TenantServiceDatastore) func(params tenant_provisioning_service_v2.DeleteDataCleaningProfileParams) middleware.Responder {
	return func(params tenant_provisioning_service_v2.DeleteDataCleaningProfileParams) middleware.Responder {
		tenantID := params.HTTPRequest.Header.Get(xFwdTenantId)
		isAuthorized, startTime := authorizeRequest(fmt.Sprintf("Deleting %s: %s", tenmod.TenantDataCleaningProfileStr, tenantID), params.HTTPRequest, allowedRoles, mon.APIRecieved, mon.TenantAPIRecieved)

		if !isAuthorized {
			return tenant_provisioning_service_v2.NewDeleteDataCleaningProfileForbidden().WithPayload(reportAPIError(fmt.Sprintf("Delete %s operation not authorized for role: %s", tenmod.TenantDataCleaningProfileStr, params.HTTPRequest.Header.Get(xFwdUserRoles)), startTime, http.StatusForbidden, mon.DeleteDataCleaningProfileStr, mon.APICompleted, mon.TenantAPICompleted))
		}

		// Issue request to DAO Layer
		result, err := tenantDB.DeleteTenantDataCleaningProfile(tenantID, params.ProfileID)
		if err != nil {
			if checkForNotFound(err.Error()) {
				return tenant_provisioning_service_v2.NewDeleteDataCleaningProfileNotFound().WithPayload(reportAPIError(generateErrorMessage(http.StatusNotFound, err.Error()), startTime, http.StatusBadRequest, mon.DeleteDataCleaningProfileStr, mon.APICompleted, mon.TenantAPICompleted))
			}
			return tenant_provisioning_service_v2.NewDeleteDataCleaningProfileInternalServerError().WithPayload(reportAPIError(fmt.Sprintf("Unable to delete %s: %s", tenmod.TenantDataCleaningProfileStr, err.Error()), startTime, http.StatusInternalServerError, mon.DeleteDataCleaningProfileStr, mon.APICompleted, mon.TenantAPICompleted))
		}

		converted := swagmodels.DataCleaningProfileResponse{}
		err = convertToJsonapiObject(result, &converted)
		if err != nil {
			return tenant_provisioning_service_v2.NewDeleteDataCleaningProfileInternalServerError().WithPayload(reportAPIError(fmt.Sprintf("Unable to convert %s data to jsonapi return format: %s", tenmod.TenantDataCleaningProfileStr, err.Error()), startTime, http.StatusInternalServerError, mon.DeleteDataCleaningProfileStr, mon.APICompleted, mon.TenantAPICompleted))
		}

		reportAPICompletionState(startTime, http.StatusOK, mon.DeleteDataCleaningProfileStr, mon.APICompleted, mon.TenantAPICompleted)
		logger.Log.Infof("Deleted %s %s", tenmod.TenantDataCleaningProfileStr, models.AsJSONString(converted))
		return tenant_provisioning_service_v2.NewDeleteDataCleaningProfileOK().WithPayload(&converted)
	}
}

// HandleUpdateDataCleaningProfileV2 - update the Data Cleaning Profile for a Tenant
func HandleUpdateDataCleaningProfileV2(allowedRoles []string, tenantDB datastore.TenantServiceDatastore) func(params tenant_provisioning_service_v2.UpdateDataCleaningProfileParams) middleware.Responder {
	return func(params tenant_provisioning_service_v2.UpdateDataCleaningProfileParams) middleware.Responder {
		tenantID := params.HTTPRequest.Header.Get(xFwdTenantId)
		isAuthorized, startTime := authorizeRequest(fmt.Sprintf("Updating %s: %s", tenmod.TenantDataCleaningProfileStr, tenantID), params.HTTPRequest, allowedRoles, mon.APIRecieved, mon.TenantAPIRecieved)

		if !isAuthorized {
			return tenant_provisioning_service_v2.NewUpdateDataCleaningProfileForbidden().WithPayload(reportAPIError(fmt.Sprintf("Update %s operation not authorized for role: %s", tenmod.TenantDataCleaningProfileStr, params.HTTPRequest.Header.Get(xFwdUserRoles)), startTime, http.StatusForbidden, mon.UpdateDataCleaningProfileStr, mon.APICompleted, mon.TenantAPICompleted))
		}

		// Fetch the existing record
		existing, err := tenantDB.GetTenantDataCleaningProfile(tenantID, params.ProfileID)
		if err != nil {
			return tenant_provisioning_service_v2.NewUpdateDataCleaningProfileConflict().WithPayload(reportAPIError(fmt.Sprintf("Unable to update %s: %s", tenmod.TenantDataCleaningProfileStr, err.Error()), startTime, http.StatusConflict, mon.UpdateDataCleaningProfileStr, mon.APICompleted, mon.TenantAPICompleted))
		}

		// Convert the request to a db model type:
		data := tenmod.DataCleaningProfile{}
		err = convertRequestBodyToDBModel(params.Body, &data)
		if err != nil {
			return tenant_provisioning_service_v2.NewUpdateDataCleaningProfileBadRequest().WithPayload(reportAPIError(generateErrorMessage(http.StatusBadRequest, err.Error()), startTime, http.StatusBadRequest, mon.UpdateDataCleaningProfileStr, mon.APICompleted, mon.TenantAPICompleted))
		}

		existing.Rules = data.Rules

		// Issue request to DAO Layer
		result, err := tenantDB.UpdateTenantDataCleaningProfile(existing)
		if err != nil {
			if checkForNotFound(err.Error()) {
				return tenant_provisioning_service_v2.NewUpdateDataCleaningProfileNotFound().WithPayload(reportAPIError(generateErrorMessage(http.StatusNotFound, err.Error()), startTime, http.StatusBadRequest, mon.UpdateDataCleaningProfileStr, mon.APICompleted, mon.TenantAPICompleted))
			}
			return tenant_provisioning_service_v2.NewUpdateDataCleaningProfileInternalServerError().WithPayload(reportAPIError(fmt.Sprintf("Unable to update %s: %s", tenmod.TenantDataCleaningProfileStr, err.Error()), startTime, http.StatusInternalServerError, mon.UpdateDataCleaningProfileStr, mon.APICompleted, mon.TenantAPICompleted))
		}

		converted := swagmodels.DataCleaningProfileResponse{}
		err = convertToJsonapiObject(result, &converted)
		if err != nil {
			return tenant_provisioning_service_v2.NewUpdateDataCleaningProfileInternalServerError().WithPayload(reportAPIError(fmt.Sprintf("Unable to convert %s data to jsonapi return format: %s", tenmod.TenantDataCleaningProfileStr, err.Error()), startTime, http.StatusInternalServerError, mon.UpdateDataCleaningProfileStr, mon.APICompleted, mon.TenantAPICompleted))
		}

		reportAPICompletionState(startTime, http.StatusOK, mon.UpdateDataCleaningProfileStr, mon.APICompleted, mon.TenantAPICompleted)
		logger.Log.Infof("Updated %s %s", tenmod.TenantDataCleaningProfileStr, models.AsJSONString(converted))
		return tenant_provisioning_service_v2.NewUpdateDataCleaningProfileOK().WithPayload(&converted)
	}
}

// HandleCreateDataCleaningProfileV2 - update the Data Cleaning Profile for a Tenant
func HandleCreateDataCleaningProfileV2(allowedRoles []string, tenantDB datastore.TenantServiceDatastore) func(params tenant_provisioning_service_v2.CreateDataCleaningProfileParams) middleware.Responder {
	return func(params tenant_provisioning_service_v2.CreateDataCleaningProfileParams) middleware.Responder {
		tenantID := params.HTTPRequest.Header.Get(xFwdTenantId)
		isAuthorized, startTime := authorizeRequest(fmt.Sprintf("Creating %s: %s", tenmod.TenantDataCleaningProfileStr, tenantID), params.HTTPRequest, allowedRoles, mon.APIRecieved, mon.TenantAPIRecieved)

		if !isAuthorized {
			return tenant_provisioning_service_v2.NewCreateDataCleaningProfileForbidden().WithPayload(reportAPIError(fmt.Sprintf("Update %s operation not authorized for role: %s", tenmod.TenantDataCleaningProfileStr, params.HTTPRequest.Header.Get(xFwdUserRoles)), startTime, http.StatusForbidden, mon.CreateDataCleaningProfileStr, mon.APICompleted, mon.TenantAPICompleted))
		}

		// Convert the request to a db model type:
		data := tenmod.DataCleaningProfile{}
		err := convertRequestBodyToDBModel(params.Body, &data)
		if err != nil {
			return tenant_provisioning_service_v2.NewCreateDataCleaningProfileBadRequest().WithPayload(reportAPIError(generateErrorMessage(http.StatusBadRequest, err.Error()), startTime, http.StatusBadRequest, mon.CreateDataCleaningProfileStr, mon.APICompleted, mon.TenantAPICompleted))
		}

		data.TenantID = tenantID

		// Issue request to DAO Layer
		result, err := tenantDB.CreateTenantDataCleaningProfile(&data)
		if err != nil {
			return tenant_provisioning_service_v2.NewCreateDataCleaningProfileInternalServerError().WithPayload(reportAPIError(fmt.Sprintf("Unable to create %s: %s", tenmod.TenantDataCleaningProfileStr, err.Error()), startTime, http.StatusInternalServerError, mon.CreateDataCleaningProfileStr, mon.APICompleted, mon.TenantAPICompleted))
		}

		converted := swagmodels.DataCleaningProfileResponse{}
		err = convertToJsonapiObject(result, &converted)
		if err != nil {
			return tenant_provisioning_service_v2.NewCreateDataCleaningProfileInternalServerError().WithPayload(reportAPIError(fmt.Sprintf("Unable to convert %s data to jsonapi return format: %s", tenmod.TenantDataCleaningProfileStr, err.Error()), startTime, http.StatusInternalServerError, mon.CreateDataCleaningProfileStr, mon.APICompleted, mon.TenantAPICompleted))
		}

		reportAPICompletionState(startTime, http.StatusOK, mon.CreateDataCleaningProfileStr, mon.APICompleted, mon.TenantAPICompleted)
		logger.Log.Infof("Created %s %s", tenmod.TenantDataCleaningProfileStr, models.AsJSONString(converted))
		return tenant_provisioning_service_v2.NewCreateDataCleaningProfileCreated().WithPayload(&converted)
	}
}
