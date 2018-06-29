package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/accedian/adh-gather/datastore"
	"github.com/accedian/adh-gather/logger"
	"github.com/accedian/adh-gather/models"
	tenmod "github.com/accedian/adh-gather/models/tenant"
	mon "github.com/accedian/adh-gather/monitoring"
	"github.com/accedian/adh-gather/restapi/operations/tenant_provisioning_service"
	"github.com/accedian/adh-gather/swagmodels"
	"github.com/go-openapi/runtime/middleware"
	"github.com/manyminds/api2go/jsonapi"
)

// HandleCreateTenantThresholdProfile - creates an Threshold profile for a tenant
func HandleCreateTenantThresholdProfile(allowedRoles []string, tenantDB datastore.TenantServiceDatastore) func(params tenant_provisioning_service.CreateTenantThresholdProfileParams) middleware.Responder {
	return func(params tenant_provisioning_service.CreateTenantThresholdProfileParams) middleware.Responder {
		startTime := time.Now()
		incrementAPICounters(mon.APIRecieved, mon.TenantAPIRecieved)
		logger.Log.Infof("Creating %s for Tenant %s", tenmod.TenantThresholdProfileStr, params.TenantID)

		if !isRequestAuthorized(params.HTTPRequest, allowedRoles) {
			return tenant_provisioning_service.NewCreateTenantThresholdProfileForbidden().WithPayload(reportAPIError(fmt.Sprintf("Create %s operation not authorized for role: %s", tenmod.TenantThresholdProfileStr, params.HTTPRequest.Header.Get(xFwdUserRoles)), startTime, http.StatusForbidden, mon.CreateThrPrfStr, mon.APICompleted, mon.TenantAPICompleted))
		}

		// Unmarshal the request
		requestBytes, err := json.Marshal(params.Body)
		if err != nil {
			return tenant_provisioning_service.NewCreateTenantThresholdProfileBadRequest().WithPayload(reportAPIError(generateErrorMessage(http.StatusBadRequest, err.Error()), startTime, http.StatusBadRequest, mon.CreateThrPrfStr, mon.APICompleted, mon.TenantAPICompleted))
		}

		data := tenmod.ThresholdProfile{}
		err = jsonapi.Unmarshal(requestBytes, &data)
		if err != nil {
			return tenant_provisioning_service.NewCreateTenantThresholdProfileBadRequest().WithPayload(reportAPIError(generateErrorMessage(http.StatusBadRequest, err.Error()), startTime, http.StatusBadRequest, mon.CreateThrPrfStr, mon.APICompleted, mon.TenantAPICompleted))
		}

		err = data.Validate(false)
		if err != nil {
			return tenant_provisioning_service.NewCreateTenantThresholdProfileBadRequest().WithPayload(reportAPIError(generateErrorMessage(http.StatusBadRequest, err.Error()), startTime, http.StatusBadRequest, mon.CreateThrPrfStr, mon.APICompleted, mon.TenantAPICompleted))
		}

		// Issue request to DAO Layer
		result, err := tenantDB.CreateTenantThresholdProfile(&data)
		if err != nil {
			return tenant_provisioning_service.NewCreateTenantThresholdProfileInternalServerError().WithPayload(reportAPIError(fmt.Sprintf("Unable to store %s: %s", tenmod.TenantThresholdProfileStr, err.Error()), startTime, http.StatusInternalServerError, mon.CreateThrPrfStr, mon.APICompleted, mon.TenantAPICompleted))
		}

		converted := swagmodels.JSONAPITenantThresholdProfile{}
		err = convertToJsonapiObject(result, &converted)
		if err != nil {
			return tenant_provisioning_service.NewCreateTenantThresholdProfileInternalServerError().WithPayload(reportAPIError(fmt.Sprintf("Unable to convert %s data to jsonapi return format: %s", tenmod.TenantThresholdProfileStr, err.Error()), startTime, http.StatusInternalServerError, mon.CreateThrPrfStr, mon.APICompleted, mon.TenantAPICompleted))
		}

		reportAPICompletionState(startTime, http.StatusOK, mon.CreateThrPrfStr, mon.APICompleted, mon.TenantAPICompleted)
		logger.Log.Infof("Created %s %s", tenmod.TenantThresholdProfileStr, models.AsJSONString(converted))
		return tenant_provisioning_service.NewCreateTenantThresholdProfileOK().WithPayload(&converted)
	}
}

// HandleUpdateTenantThresholdProfile - updates a n Threshold profile for a tenant
func HandleUpdateTenantThresholdProfile(allowedRoles []string, tenantDB datastore.TenantServiceDatastore) func(params tenant_provisioning_service.UpdateTenantThresholdProfileParams) middleware.Responder {
	return func(params tenant_provisioning_service.UpdateTenantThresholdProfileParams) middleware.Responder {
		startTime := time.Now()
		incrementAPICounters(mon.APIRecieved, mon.TenantAPIRecieved)
		logger.Log.Infof("Updating %s for Tenant %s", tenmod.TenantThresholdProfileStr, params.TenantID)

		if !isRequestAuthorized(params.HTTPRequest, allowedRoles) {
			return tenant_provisioning_service.NewUpdateTenantThresholdProfileForbidden().WithPayload(reportAPIError(fmt.Sprintf("Update %s operation not authorized for role: %s", tenmod.TenantThresholdProfileStr, params.HTTPRequest.Header.Get(xFwdUserRoles)), startTime, http.StatusForbidden, mon.UpdateThrPrfStr, mon.APICompleted, mon.TenantAPICompleted))
		}

		// Unmarshal the request
		requestBytes, err := json.Marshal(params.Body)
		if err != nil {
			return tenant_provisioning_service.NewUpdateTenantThresholdProfileBadRequest().WithPayload(reportAPIError(generateErrorMessage(http.StatusBadRequest, err.Error()), startTime, http.StatusBadRequest, mon.UpdateThrPrfStr, mon.APICompleted, mon.TenantAPICompleted))
		}

		data := tenmod.ThresholdProfile{}
		err = jsonapi.Unmarshal(requestBytes, &data)
		if err != nil {
			return tenant_provisioning_service.NewUpdateTenantThresholdProfileBadRequest().WithPayload(reportAPIError(generateErrorMessage(http.StatusBadRequest, err.Error()), startTime, http.StatusBadRequest, mon.UpdateThrPrfStr, mon.APICompleted, mon.TenantAPICompleted))
		}

		err = data.Validate(true)
		if err != nil {
			return tenant_provisioning_service.NewUpdateTenantThresholdProfileBadRequest().WithPayload(reportAPIError(generateErrorMessage(http.StatusBadRequest, err.Error()), startTime, http.StatusBadRequest, mon.UpdateThrPrfStr, mon.APICompleted, mon.TenantAPICompleted))
		}

		// Issue request to DAO Layer
		result, err := tenantDB.UpdateTenantThresholdProfile(&data)
		if err != nil {
			return tenant_provisioning_service.NewUpdateTenantThresholdProfileInternalServerError().WithPayload(reportAPIError(fmt.Sprintf("Unable to store %s: %s", tenmod.TenantThresholdProfileStr, err.Error()), startTime, http.StatusInternalServerError, mon.UpdateThrPrfStr, mon.APICompleted, mon.TenantAPICompleted))
		}

		converted := swagmodels.JSONAPITenantThresholdProfile{}
		err = convertToJsonapiObject(result, &converted)
		if err != nil {
			return tenant_provisioning_service.NewUpdateTenantThresholdProfileInternalServerError().WithPayload(reportAPIError(fmt.Sprintf("Unable to convert %s data to jsonapi return format: %s", tenmod.TenantThresholdProfileStr, err.Error()), startTime, http.StatusInternalServerError, mon.UpdateThrPrfStr, mon.APICompleted, mon.TenantAPICompleted))
		}

		reportAPICompletionState(startTime, http.StatusOK, mon.UpdateThrPrfStr, mon.APICompleted, mon.TenantAPICompleted)
		logger.Log.Infof("Updated %s %s", tenmod.TenantThresholdProfileStr, models.AsJSONString(converted))
		return tenant_provisioning_service.NewUpdateTenantThresholdProfileOK().WithPayload(&converted)
	}
}

// HandleGetTenantThresholdProfile - fetch an Threshold profile for a tenant
func HandleGetTenantThresholdProfile(allowedRoles []string, tenantDB datastore.TenantServiceDatastore) func(params tenant_provisioning_service.GetTenantThresholdProfileParams) middleware.Responder {
	return func(params tenant_provisioning_service.GetTenantThresholdProfileParams) middleware.Responder {
		startTime := time.Now()
		incrementAPICounters(mon.APIRecieved, mon.TenantAPIRecieved)
		logger.Log.Infof("Fetching %s %s for Tenant %s", tenmod.TenantThresholdProfileStr, params.ThrPrfID, params.TenantID)

		if !isRequestAuthorized(params.HTTPRequest, allowedRoles) {
			return tenant_provisioning_service.NewGetTenantThresholdProfileForbidden().WithPayload(reportAPIError(fmt.Sprintf("Get %s operation not authorized for role: %s", tenmod.TenantThresholdProfileStr, params.HTTPRequest.Header.Get(xFwdUserRoles)), startTime, http.StatusForbidden, mon.GetThrPrfStr, mon.APICompleted, mon.TenantAPICompleted))
		}

		// Issue request to DAO Layer
		result, err := tenantDB.GetTenantThresholdProfile(params.TenantID, params.ThrPrfID)
		if err != nil {
			return tenant_provisioning_service.NewGetTenantThresholdProfileInternalServerError().WithPayload(reportAPIError(fmt.Sprintf("Unable to retrieve %s: %s", tenmod.TenantThresholdProfileStr, err.Error()), startTime, http.StatusInternalServerError, mon.GetThrPrfStr, mon.APICompleted, mon.TenantAPICompleted))
		}

		converted := swagmodels.JSONAPITenantThresholdProfile{}
		err = convertToJsonapiObject(result, &converted)
		if err != nil {
			return tenant_provisioning_service.NewGetTenantThresholdProfileInternalServerError().WithPayload(reportAPIError(fmt.Sprintf("Unable to convert %s data to jsonapi return format: %s", tenmod.TenantThresholdProfileStr, err.Error()), startTime, http.StatusInternalServerError, mon.GetThrPrfStr, mon.APICompleted, mon.TenantAPICompleted))
		}

		reportAPICompletionState(startTime, http.StatusOK, mon.GetThrPrfStr, mon.APICompleted, mon.TenantAPICompleted)
		logger.Log.Infof("Retrieved %s %s", tenmod.TenantThresholdProfileStr, models.AsJSONString(converted))
		return tenant_provisioning_service.NewGetTenantThresholdProfileOK().WithPayload(&converted)
	}
}

// HandleGetAllTenantThresholdProfiles - fetch an Threshold profile for a tenant
func HandleGetAllTenantThresholdProfiles(allowedRoles []string, tenantDB datastore.TenantServiceDatastore) func(params tenant_provisioning_service.GetAllTenantThresholdProfilesParams) middleware.Responder {
	return func(params tenant_provisioning_service.GetAllTenantThresholdProfilesParams) middleware.Responder {
		startTime := time.Now()
		incrementAPICounters(mon.APIRecieved, mon.TenantAPIRecieved)
		logger.Log.Infof("Fetching all %ss for Tenant %s", tenmod.TenantThresholdProfileStr, params.TenantID)

		if !isRequestAuthorized(params.HTTPRequest, allowedRoles) {
			return tenant_provisioning_service.NewGetAllTenantThresholdProfilesForbidden().WithPayload(reportAPIError(fmt.Sprintf("Get all %ss operation not authorized for role: %s", tenmod.TenantThresholdProfileStr, params.HTTPRequest.Header.Get(xFwdUserRoles)), startTime, http.StatusForbidden, mon.GetAllThrPrfStr, mon.APICompleted, mon.TenantAPICompleted))
		}

		// Issue request to DAO Layer
		result, err := tenantDB.GetAllTenantThresholdProfile(params.TenantID)
		if err != nil {
			return tenant_provisioning_service.NewGetAllTenantThresholdProfilesInternalServerError().WithPayload(reportAPIError(fmt.Sprintf("Unable to retrieve %s list: %s", tenmod.TenantThresholdProfileStr, err.Error()), startTime, http.StatusInternalServerError, mon.GetAllThrPrfStr, mon.APICompleted, mon.TenantAPICompleted))
		}

		converted := swagmodels.JSONAPITenantThresholdProfileList{}
		err = convertToJsonapiObject(result, &converted)
		if err != nil {
			return tenant_provisioning_service.NewGetAllTenantThresholdProfilesInternalServerError().WithPayload(reportAPIError(fmt.Sprintf("Unable to convert %s list data to jsonapi return format: %s", tenmod.TenantThresholdProfileStr, err.Error()), startTime, http.StatusInternalServerError, mon.GetAllThrPrfStr, mon.APICompleted, mon.TenantAPICompleted))
		}

		reportAPICompletionState(startTime, http.StatusOK, mon.GetAllThrPrfStr, mon.APICompleted, mon.TenantAPICompleted)
		logger.Log.Infof("Retrieved %d %ss", len(result), tenmod.TenantThresholdProfileStr)
		return tenant_provisioning_service.NewGetAllTenantThresholdProfilesOK().WithPayload(&converted)
	}
}

// HandleDeleteTenantThresholdProfile - delete an Threshold profile for a tenant
func HandleDeleteTenantThresholdProfile(allowedRoles []string, tenantDB datastore.TenantServiceDatastore) func(params tenant_provisioning_service.DeleteTenantThresholdProfileParams) middleware.Responder {
	return func(params tenant_provisioning_service.DeleteTenantThresholdProfileParams) middleware.Responder {
		startTime := time.Now()
		incrementAPICounters(mon.APIRecieved, mon.TenantAPIRecieved)
		logger.Log.Infof("Deleting %s %s for Tenant %s", tenmod.TenantThresholdProfileStr, params.ThrPrfID, params.TenantID)

		if !isRequestAuthorized(params.HTTPRequest, allowedRoles) {
			return tenant_provisioning_service.NewDeleteTenantThresholdProfileForbidden().WithPayload(reportAPIError(fmt.Sprintf("Delete %s operation not authorized for role: %s", tenmod.TenantThresholdProfileStr, params.HTTPRequest.Header.Get(xFwdUserRoles)), startTime, http.StatusForbidden, mon.DeleteThrPrfStr, mon.APICompleted, mon.TenantAPICompleted))
		}

		// Issue request to DAO Layer
		result, err := tenantDB.DeleteTenantThresholdProfile(params.TenantID, params.ThrPrfID)
		if err != nil {
			return tenant_provisioning_service.NewDeleteTenantThresholdProfileInternalServerError().WithPayload(reportAPIError(fmt.Sprintf("Unable to delete %s: %s", tenmod.TenantThresholdProfileStr, err.Error()), startTime, http.StatusInternalServerError, mon.DeleteThrPrfStr, mon.APICompleted, mon.TenantAPICompleted))
		}

		converted := swagmodels.JSONAPITenantThresholdProfile{}
		err = convertToJsonapiObject(result, &converted)
		if err != nil {
			return tenant_provisioning_service.NewDeleteTenantThresholdProfileInternalServerError().WithPayload(reportAPIError(fmt.Sprintf("Unable to convert %s data to jsonapi return format: %s", tenmod.TenantThresholdProfileStr, err.Error()), startTime, http.StatusInternalServerError, mon.DeleteThrPrfStr, mon.APICompleted, mon.TenantAPICompleted))
		}

		reportAPICompletionState(startTime, http.StatusOK, mon.DeleteThrPrfStr, mon.APICompleted, mon.TenantAPICompleted)
		logger.Log.Infof("Deleted %s %s", tenmod.TenantThresholdProfileStr, models.AsJSONString(converted))
		return tenant_provisioning_service.NewDeleteTenantThresholdProfileOK().WithPayload(&converted)
	}
}

// HandlePatchTenantThresholdProfile - patch an Threshold profile for a tenant
func HandlePatchTenantThresholdProfile(allowedRoles []string, tenantDB datastore.TenantServiceDatastore) func(params tenant_provisioning_service.PatchTenantThresholdProfileParams) middleware.Responder {
	return func(params tenant_provisioning_service.PatchTenantThresholdProfileParams) middleware.Responder {
		startTime := time.Now()
		incrementAPICounters(mon.APIRecieved, mon.TenantAPIRecieved)
		logger.Log.Infof("Patching %s %s for Tenant %s", tenmod.TenantThresholdProfileStr, params.Body.Data.ID, params.TenantID)

		if !isRequestAuthorized(params.HTTPRequest, allowedRoles) {
			return tenant_provisioning_service.NewPatchTenantThresholdProfileForbidden().WithPayload(reportAPIError(fmt.Sprintf("Patch %s operation not authorized for role: %s", tenmod.TenantThresholdProfileStr, params.HTTPRequest.Header.Get(xFwdUserRoles)), startTime, http.StatusForbidden, mon.PatchThrPrfStr, mon.APICompleted, mon.TenantAPICompleted))
		}

		// Unmarshal the request
		requestBytes, err := json.Marshal(params.Body)
		if err != nil {
			return tenant_provisioning_service.NewPatchTenantThresholdProfileBadRequest().WithPayload(reportAPIError(generateErrorMessage(http.StatusBadRequest, err.Error()), startTime, http.StatusBadRequest, mon.PatchThrPrfStr, mon.APICompleted, mon.TenantAPICompleted))
		}

		data := tenmod.ThresholdProfile{}
		err = jsonapi.Unmarshal(requestBytes, &data)
		if err != nil {
			return tenant_provisioning_service.NewPatchTenantThresholdProfileBadRequest().WithPayload(reportAPIError(generateErrorMessage(http.StatusBadRequest, err.Error()), startTime, http.StatusBadRequest, mon.PatchThrPrfStr, mon.APICompleted, mon.TenantAPICompleted))
		}

		err = data.Validate(true)
		if err != nil {
			return tenant_provisioning_service.NewPatchTenantThresholdProfileBadRequest().WithPayload(reportAPIError(generateErrorMessage(http.StatusBadRequest, err.Error()), startTime, http.StatusBadRequest, mon.PatchThrPrfStr, mon.APICompleted, mon.TenantAPICompleted))
		}

		origData, err2 := tenantDB.GetTenantThresholdProfile(data.TenantID, data.ID)
		if err2 != nil {
			return tenant_provisioning_service.NewPatchTenantThresholdProfileBadRequest().WithPayload(reportAPIError(generateErrorMessage(http.StatusBadRequest, err.Error()), startTime, http.StatusBadRequest, mon.PatchThrPrfStr, mon.APICompleted, mon.TenantAPICompleted))
		}

		errMerge := models.MergeObjWithMap(origData, requestBytes)
		if errMerge != nil {
			return tenant_provisioning_service.NewPatchTenantThresholdProfileBadRequest().WithPayload(reportAPIError(generateErrorMessage(http.StatusBadRequest, err.Error()), startTime, http.StatusBadRequest, mon.PatchThrPrfStr, mon.APICompleted, mon.TenantAPICompleted))
		}

		// Issue request to DAO Layer
		result, err := tenantDB.UpdateTenantThresholdProfile(origData)
		if err != nil {
			return tenant_provisioning_service.NewPatchTenantThresholdProfileInternalServerError().WithPayload(reportAPIError(fmt.Sprintf("Unable to store %s: %s", tenmod.TenantThresholdProfileStr, err.Error()), startTime, http.StatusInternalServerError, mon.PatchThrPrfStr, mon.APICompleted, mon.TenantAPICompleted))
		}

		converted := swagmodels.JSONAPITenantThresholdProfile{}
		err = convertToJsonapiObject(result, &converted)
		if err != nil {
			return tenant_provisioning_service.NewPatchTenantThresholdProfileInternalServerError().WithPayload(reportAPIError(fmt.Sprintf("Unable to convert %s data to jsonapi return format: %s", tenmod.TenantThresholdProfileStr, err.Error()), startTime, http.StatusInternalServerError, mon.PatchThrPrfStr, mon.APICompleted, mon.TenantAPICompleted))
		}

		reportAPICompletionState(startTime, http.StatusOK, mon.PatchThrPrfStr, mon.APICompleted, mon.TenantAPICompleted)
		logger.Log.Infof("Patched %s %s", tenmod.TenantThresholdProfileStr, models.AsJSONString(converted))
		return tenant_provisioning_service.NewPatchTenantThresholdProfileOK().WithPayload(&converted)
	}
}
