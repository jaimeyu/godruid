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

// HandleCreateTenantIngestionProfile - creates an ingestion profile for a tenant
func HandleCreateTenantIngestionProfile(allowedRoles []string, tenantDB datastore.TenantServiceDatastore) func(params tenant_provisioning_service.CreateTenantIngestionProfileParams) middleware.Responder {
	return func(params tenant_provisioning_service.CreateTenantIngestionProfileParams) middleware.Responder {
		startTime := time.Now()
		incrementAPICounters(mon.APIRecieved, mon.TenantAPIRecieved)
		logger.Log.Infof("Creating %s for Tenant %s", tenmod.TenantIngestionProfileStr, params.TenantID)

		if !isRequestAuthorized(params.HTTPRequest, allowedRoles) {
			return tenant_provisioning_service.NewCreateTenantIngestionProfileForbidden().WithPayload(reportAPIError(fmt.Sprintf("Create %s operation not authorized for role: %s", tenmod.TenantIngestionProfileStr, params.HTTPRequest.Header.Get(xFwdUserRoles)), startTime, http.StatusForbidden, mon.CreateIngPrfStr, mon.APICompleted, mon.TenantAPICompleted))
		}

		// Unmarshal the request
		requestBytes, err := json.Marshal(params.Body)
		if err != nil {
			return tenant_provisioning_service.NewCreateTenantIngestionProfileBadRequest().WithPayload(reportAPIError(generateErrorMessage(http.StatusBadRequest, err.Error()), startTime, http.StatusBadRequest, mon.CreateIngPrfStr, mon.APICompleted, mon.TenantAPICompleted))
		}

		data := tenmod.IngestionProfile{}
		err = jsonapi.Unmarshal(requestBytes, &data)
		if err != nil {
			return tenant_provisioning_service.NewCreateTenantIngestionProfileBadRequest().WithPayload(reportAPIError(generateErrorMessage(http.StatusBadRequest, err.Error()), startTime, http.StatusBadRequest, mon.CreateIngPrfStr, mon.APICompleted, mon.TenantAPICompleted))
		}

		err = data.Validate(false)
		if err != nil {
			return tenant_provisioning_service.NewCreateTenantIngestionProfileBadRequest().WithPayload(reportAPIError(generateErrorMessage(http.StatusBadRequest, err.Error()), startTime, http.StatusBadRequest, mon.CreateIngPrfStr, mon.APICompleted, mon.TenantAPICompleted))
		}

		// Issue request to DAO Layer
		result, err := tenantDB.CreateTenantIngestionProfile(&data)
		if err != nil {
			return tenant_provisioning_service.NewCreateTenantIngestionProfileInternalServerError().WithPayload(reportAPIError(fmt.Sprintf("Unable to store %s: %s", tenmod.TenantIngestionProfileStr, err.Error()), startTime, http.StatusInternalServerError, mon.CreateIngPrfStr, mon.APICompleted, mon.TenantAPICompleted))
		}

		converted := swagmodels.JSONAPITenantIngestionProfile{}
		err = convertToJsonapiObject(result, &converted)
		if err != nil {
			return tenant_provisioning_service.NewCreateTenantIngestionProfileInternalServerError().WithPayload(reportAPIError(fmt.Sprintf("Unable to convert %s data to jsonapi return format: %s", tenmod.TenantIngestionProfileStr, err.Error()), startTime, http.StatusInternalServerError, mon.CreateIngPrfStr, mon.APICompleted, mon.TenantAPICompleted))
		}

		reportAPICompletionState(startTime, http.StatusOK, mon.CreateIngPrfStr, mon.APICompleted, mon.TenantAPICompleted)
		logger.Log.Infof("Created %s %s", tenmod.TenantIngestionProfileStr, models.AsJSONString(converted))
		return tenant_provisioning_service.NewCreateTenantIngestionProfileOK().WithPayload(&converted)
	}
}

// HandleUpdateTenantIngestionProfile - updates a n ingestion profile for a tenant
func HandleUpdateTenantIngestionProfile(allowedRoles []string, tenantDB datastore.TenantServiceDatastore) func(params tenant_provisioning_service.UpdateTenantIngestionProfileParams) middleware.Responder {
	return func(params tenant_provisioning_service.UpdateTenantIngestionProfileParams) middleware.Responder {
		startTime := time.Now()
		incrementAPICounters(mon.APIRecieved, mon.TenantAPIRecieved)
		logger.Log.Infof("Updating %s for Tenant %s", tenmod.TenantIngestionProfileStr, params.TenantID)

		if !isRequestAuthorized(params.HTTPRequest, allowedRoles) {
			return tenant_provisioning_service.NewUpdateTenantIngestionProfileForbidden().WithPayload(reportAPIError(fmt.Sprintf("Update %s operation not authorized for role: %s", tenmod.TenantIngestionProfileStr, params.HTTPRequest.Header.Get(xFwdUserRoles)), startTime, http.StatusForbidden, mon.UpdateIngPrfStr, mon.APICompleted, mon.TenantAPICompleted))
		}

		// Unmarshal the request
		requestBytes, err := json.Marshal(params.Body)
		if err != nil {
			return tenant_provisioning_service.NewUpdateTenantIngestionProfileBadRequest().WithPayload(reportAPIError(generateErrorMessage(http.StatusBadRequest, err.Error()), startTime, http.StatusBadRequest, mon.UpdateIngPrfStr, mon.APICompleted, mon.TenantAPICompleted))
		}

		data := tenmod.IngestionProfile{}
		err = jsonapi.Unmarshal(requestBytes, &data)
		if err != nil {
			return tenant_provisioning_service.NewUpdateTenantIngestionProfileBadRequest().WithPayload(reportAPIError(generateErrorMessage(http.StatusBadRequest, err.Error()), startTime, http.StatusBadRequest, mon.UpdateIngPrfStr, mon.APICompleted, mon.TenantAPICompleted))
		}

		err = data.Validate(true)
		if err != nil {
			return tenant_provisioning_service.NewUpdateTenantIngestionProfileBadRequest().WithPayload(reportAPIError(generateErrorMessage(http.StatusBadRequest, err.Error()), startTime, http.StatusBadRequest, mon.UpdateIngPrfStr, mon.APICompleted, mon.TenantAPICompleted))
		}

		// Issue request to DAO Layer
		result, err := tenantDB.UpdateTenantIngestionProfile(&data)
		if err != nil {
			return tenant_provisioning_service.NewUpdateTenantIngestionProfileInternalServerError().WithPayload(reportAPIError(fmt.Sprintf("Unable to store %s: %s", tenmod.TenantIngestionProfileStr, err.Error()), startTime, http.StatusInternalServerError, mon.UpdateIngPrfStr, mon.APICompleted, mon.TenantAPICompleted))
		}

		converted := swagmodels.JSONAPITenantIngestionProfile{}
		err = convertToJsonapiObject(result, &converted)
		if err != nil {
			return tenant_provisioning_service.NewUpdateTenantIngestionProfileInternalServerError().WithPayload(reportAPIError(fmt.Sprintf("Unable to convert %s data to jsonapi return format: %s", tenmod.TenantIngestionProfileStr, err.Error()), startTime, http.StatusInternalServerError, mon.UpdateIngPrfStr, mon.APICompleted, mon.TenantAPICompleted))
		}

		reportAPICompletionState(startTime, http.StatusOK, mon.UpdateIngPrfStr, mon.APICompleted, mon.TenantAPICompleted)
		logger.Log.Infof("Updated %s %s", tenmod.TenantIngestionProfileStr, models.AsJSONString(converted))
		return tenant_provisioning_service.NewUpdateTenantIngestionProfileOK().WithPayload(&converted)
	}
}

// HandleGetTenantIngestionProfile - fetch an ingestion profile for a tenant
func HandleGetTenantIngestionProfile(allowedRoles []string, tenantDB datastore.TenantServiceDatastore) func(params tenant_provisioning_service.GetTenantIngestionProfileParams) middleware.Responder {
	return func(params tenant_provisioning_service.GetTenantIngestionProfileParams) middleware.Responder {
		startTime := time.Now()
		incrementAPICounters(mon.APIRecieved, mon.TenantAPIRecieved)
		logger.Log.Infof("Fetching %s %s for Tenant %s", tenmod.TenantIngestionProfileStr, params.IngestionProfileID, params.TenantID)

		if !isRequestAuthorized(params.HTTPRequest, allowedRoles) {
			return tenant_provisioning_service.NewGetTenantIngestionProfileForbidden().WithPayload(reportAPIError(fmt.Sprintf("Get %s operation not authorized for role: %s", tenmod.TenantIngestionProfileStr, params.HTTPRequest.Header.Get(xFwdUserRoles)), startTime, http.StatusForbidden, mon.GetIngPrfStr, mon.APICompleted, mon.TenantAPICompleted))
		}

		// Issue request to DAO Layer
		result, err := tenantDB.GetTenantIngestionProfile(params.TenantID, params.IngestionProfileID)
		if err != nil {
			return tenant_provisioning_service.NewGetTenantIngestionProfileInternalServerError().WithPayload(reportAPIError(fmt.Sprintf("Unable to retrieve %s: %s", tenmod.TenantIngestionProfileStr, err.Error()), startTime, http.StatusInternalServerError, mon.GetIngPrfStr, mon.APICompleted, mon.TenantAPICompleted))
		}

		converted := swagmodels.JSONAPITenantIngestionProfile{}
		err = convertToJsonapiObject(result, &converted)
		if err != nil {
			return tenant_provisioning_service.NewGetTenantIngestionProfileInternalServerError().WithPayload(reportAPIError(fmt.Sprintf("Unable to convert %s data to jsonapi return format: %s", tenmod.TenantIngestionProfileStr, err.Error()), startTime, http.StatusInternalServerError, mon.GetIngPrfStr, mon.APICompleted, mon.TenantAPICompleted))
		}

		reportAPICompletionState(startTime, http.StatusOK, mon.GetIngPrfStr, mon.APICompleted, mon.TenantAPICompleted)
		logger.Log.Infof("Retrieved %s %s", tenmod.TenantIngestionProfileStr, models.AsJSONString(converted))
		return tenant_provisioning_service.NewGetTenantIngestionProfileOK().WithPayload(&converted)
	}
}

// HandleGetActiveTenantIngestionProfile - fetch an ingestion profile for a tenant
func HandleGetActiveTenantIngestionProfile(allowedRoles []string, tenantDB datastore.TenantServiceDatastore) func(params tenant_provisioning_service.GetActiveTenantIngestionProfileParams) middleware.Responder {
	return func(params tenant_provisioning_service.GetActiveTenantIngestionProfileParams) middleware.Responder {
		startTime := time.Now()
		incrementAPICounters(mon.APIRecieved, mon.TenantAPIRecieved)
		logger.Log.Infof("Fetching active %s for Tenant %s", tenmod.TenantIngestionProfileStr, params.TenantID)

		if !isRequestAuthorized(params.HTTPRequest, allowedRoles) {
			return tenant_provisioning_service.NewGetActiveTenantIngestionProfileForbidden().WithPayload(reportAPIError(fmt.Sprintf("Get %s operation not authorized for role: %s", tenmod.TenantIngestionProfileStr, params.HTTPRequest.Header.Get(xFwdUserRoles)), startTime, http.StatusForbidden, mon.GetActiveIngPrfStr, mon.APICompleted, mon.TenantAPICompleted))
		}

		// Issue request to DAO Layer
		result, err := tenantDB.GetActiveTenantIngestionProfile(params.TenantID)
		if err != nil {
			return tenant_provisioning_service.NewGetActiveTenantIngestionProfileInternalServerError().WithPayload(reportAPIError(fmt.Sprintf("Unable to retrieve %s: %s", tenmod.TenantIngestionProfileStr, err.Error()), startTime, http.StatusInternalServerError, mon.GetActiveIngPrfStr, mon.APICompleted, mon.TenantAPICompleted))
		}

		converted := swagmodels.JSONAPITenantIngestionProfile{}
		err = convertToJsonapiObject(result, &converted)
		if err != nil {
			return tenant_provisioning_service.NewGetActiveTenantIngestionProfileInternalServerError().WithPayload(reportAPIError(fmt.Sprintf("Unable to convert %s data to jsonapi return format: %s", tenmod.TenantIngestionProfileStr, err.Error()), startTime, http.StatusInternalServerError, mon.GetActiveIngPrfStr, mon.APICompleted, mon.TenantAPICompleted))
		}

		reportAPICompletionState(startTime, http.StatusOK, mon.GetActiveIngPrfStr, mon.APICompleted, mon.TenantAPICompleted)
		logger.Log.Infof("Retrieved %s %s", tenmod.TenantIngestionProfileStr, models.AsJSONString(converted))
		return tenant_provisioning_service.NewGetActiveTenantIngestionProfileOK().WithPayload(&converted)
	}
}

// HandleDeleteTenantIngestionProfile - delete an ingestion profile for a tenant
func HandleDeleteTenantIngestionProfile(allowedRoles []string, tenantDB datastore.TenantServiceDatastore) func(params tenant_provisioning_service.DeleteTenantIngestionProfileParams) middleware.Responder {
	return func(params tenant_provisioning_service.DeleteTenantIngestionProfileParams) middleware.Responder {
		startTime := time.Now()
		incrementAPICounters(mon.APIRecieved, mon.TenantAPIRecieved)
		logger.Log.Infof("Deleting %s %s for Tenant %s", tenmod.TenantIngestionProfileStr, params.IngestionProfileID, params.TenantID)

		if !isRequestAuthorized(params.HTTPRequest, allowedRoles) {
			return tenant_provisioning_service.NewDeleteTenantIngestionProfileForbidden().WithPayload(reportAPIError(fmt.Sprintf("Delete %s operation not authorized for role: %s", tenmod.TenantIngestionProfileStr, params.HTTPRequest.Header.Get(xFwdUserRoles)), startTime, http.StatusForbidden, mon.DeleteIngPrfStr, mon.APICompleted, mon.TenantAPICompleted))
		}

		// Issue request to DAO Layer
		result, err := tenantDB.DeleteTenantConnectorConfig(params.TenantID, params.IngestionProfileID)
		if err != nil {
			return tenant_provisioning_service.NewDeleteTenantIngestionProfileInternalServerError().WithPayload(reportAPIError(fmt.Sprintf("Unable to delete %s: %s", tenmod.TenantIngestionProfileStr, err.Error()), startTime, http.StatusInternalServerError, mon.DeleteIngPrfStr, mon.APICompleted, mon.TenantAPICompleted))
		}

		converted := swagmodels.JSONAPITenantIngestionProfile{}
		err = convertToJsonapiObject(result, &converted)
		if err != nil {
			return tenant_provisioning_service.NewDeleteTenantIngestionProfileInternalServerError().WithPayload(reportAPIError(fmt.Sprintf("Unable to convert %s data to jsonapi return format: %s", tenmod.TenantIngestionProfileStr, err.Error()), startTime, http.StatusInternalServerError, mon.DeleteIngPrfStr, mon.APICompleted, mon.TenantAPICompleted))
		}

		reportAPICompletionState(startTime, http.StatusOK, mon.DeleteIngPrfStr, mon.APICompleted, mon.TenantAPICompleted)
		logger.Log.Infof("Deleted %s %s", tenmod.TenantIngestionProfileStr, models.AsJSONString(converted))
		return tenant_provisioning_service.NewDeleteTenantIngestionProfileOK().WithPayload(&converted)
	}
}

// HandlePatchTenantIngestionProfile - patch an ingestion profile for a tenant
func HandlePatchTenantIngestionProfile(allowedRoles []string, tenantDB datastore.TenantServiceDatastore) func(params tenant_provisioning_service.PatchTenantIngestionProfileParams) middleware.Responder {
	return func(params tenant_provisioning_service.PatchTenantIngestionProfileParams) middleware.Responder {
		startTime := time.Now()
		incrementAPICounters(mon.APIRecieved, mon.TenantAPIRecieved)
		logger.Log.Infof("Patching %s %s for Tenant %s", tenmod.TenantIngestionProfileStr, params.Body.Data.ID, params.TenantID)

		if !isRequestAuthorized(params.HTTPRequest, allowedRoles) {
			return tenant_provisioning_service.NewPatchTenantIngestionProfileForbidden().WithPayload(reportAPIError(fmt.Sprintf("Patch %s operation not authorized for role: %s", tenmod.TenantIngestionProfileStr, params.HTTPRequest.Header.Get(xFwdUserRoles)), startTime, http.StatusForbidden, mon.PatchIngPrfStr, mon.APICompleted, mon.TenantAPICompleted))
		}

		// Unmarshal the request
		requestBytes, err := json.Marshal(params.Body)
		if err != nil {
			return tenant_provisioning_service.NewPatchTenantIngestionProfileBadRequest().WithPayload(reportAPIError(generateErrorMessage(http.StatusBadRequest, err.Error()), startTime, http.StatusBadRequest, mon.PatchIngPrfStr, mon.APICompleted, mon.TenantAPICompleted))
		}

		data := tenmod.IngestionProfile{}
		err = jsonapi.Unmarshal(requestBytes, &data)
		if err != nil {
			return tenant_provisioning_service.NewPatchTenantIngestionProfileBadRequest().WithPayload(reportAPIError(generateErrorMessage(http.StatusBadRequest, err.Error()), startTime, http.StatusBadRequest, mon.PatchIngPrfStr, mon.APICompleted, mon.TenantAPICompleted))
		}

		err = data.Validate(true)
		if err != nil {
			return tenant_provisioning_service.NewPatchTenantIngestionProfileBadRequest().WithPayload(reportAPIError(generateErrorMessage(http.StatusBadRequest, err.Error()), startTime, http.StatusBadRequest, mon.PatchIngPrfStr, mon.APICompleted, mon.TenantAPICompleted))
		}

		origData, err2 := tenantDB.GetTenantIngestionProfile(data.TenantID, data.ID)
		if err2 != nil {
			return tenant_provisioning_service.NewPatchTenantIngestionProfileBadRequest().WithPayload(reportAPIError(generateErrorMessage(http.StatusBadRequest, err.Error()), startTime, http.StatusBadRequest, mon.PatchIngPrfStr, mon.APICompleted, mon.TenantAPICompleted))
		}

		errMerge := models.MergeObjWithMap(origData, requestBytes)
		if errMerge != nil {
			return tenant_provisioning_service.NewPatchTenantIngestionProfileBadRequest().WithPayload(reportAPIError(generateErrorMessage(http.StatusBadRequest, err.Error()), startTime, http.StatusBadRequest, mon.PatchIngPrfStr, mon.APICompleted, mon.TenantAPICompleted))
		}

		// Issue request to DAO Layer
		result, err := tenantDB.UpdateTenantIngestionProfile(origData)
		if err != nil {
			return tenant_provisioning_service.NewPatchTenantIngestionProfileInternalServerError().WithPayload(reportAPIError(fmt.Sprintf("Unable to store %s: %s", tenmod.TenantIngestionProfileStr, err.Error()), startTime, http.StatusInternalServerError, mon.PatchIngPrfStr, mon.APICompleted, mon.TenantAPICompleted))
		}

		converted := swagmodels.JSONAPITenantIngestionProfile{}
		err = convertToJsonapiObject(result, &converted)
		if err != nil {
			return tenant_provisioning_service.NewPatchTenantIngestionProfileInternalServerError().WithPayload(reportAPIError(fmt.Sprintf("Unable to convert %s data to jsonapi return format: %s", tenmod.TenantIngestionProfileStr, err.Error()), startTime, http.StatusInternalServerError, mon.PatchIngPrfStr, mon.APICompleted, mon.TenantAPICompleted))
		}

		reportAPICompletionState(startTime, http.StatusOK, mon.PatchIngPrfStr, mon.APICompleted, mon.TenantAPICompleted)
		logger.Log.Infof("Patched %s %s", tenmod.TenantIngestionProfileStr, models.AsJSONString(converted))
		return tenant_provisioning_service.NewPatchTenantIngestionProfileOK().WithPayload(&converted)
	}
}
