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

// HandleCreateTenantMetadata - creates an ingestion profile for a tenant
func HandleCreateTenantMetadata(allowedRoles []string, tenantDB datastore.TenantServiceDatastore) func(params tenant_provisioning_service.CreateTenantMetadataParams) middleware.Responder {
	return func(params tenant_provisioning_service.CreateTenantMetadataParams) middleware.Responder {
		startTime := time.Now()
		incrementAPICounters(mon.APIRecieved, mon.TenantAPIRecieved)
		logger.Log.Infof("Creating %s for Tenant %s", tenmod.TenantMetaStr, params.TenantID)

		if !isRequestAuthorized(params.HTTPRequest, allowedRoles) {
			return tenant_provisioning_service.NewCreateTenantMetadataForbidden().WithPayload(reportAPIError(fmt.Sprintf("Create %s operation not authorized for role: %s", tenmod.TenantMetaStr, params.HTTPRequest.Header.Get(XFwdUserRoles)), startTime, http.StatusForbidden, mon.CreateTenantMetaStr, mon.APICompleted, mon.TenantAPICompleted))
		}

		// Unmarshal the request
		requestBytes, err := json.Marshal(params.Body)
		if err != nil {
			return tenant_provisioning_service.NewCreateTenantMetadataBadRequest().WithPayload(reportAPIError(generateErrorMessage(http.StatusBadRequest, err.Error()), startTime, http.StatusBadRequest, mon.CreateTenantMetaStr, mon.APICompleted, mon.TenantAPICompleted))
		}

		data := tenmod.Metadata{}
		err = jsonapi.Unmarshal(requestBytes, &data)
		if err != nil {
			return tenant_provisioning_service.NewCreateTenantMetadataBadRequest().WithPayload(reportAPIError(generateErrorMessage(http.StatusBadRequest, err.Error()), startTime, http.StatusBadRequest, mon.CreateTenantMetaStr, mon.APICompleted, mon.TenantAPICompleted))
		}

		err = data.Validate(false)
		if err != nil {
			return tenant_provisioning_service.NewCreateTenantMetadataBadRequest().WithPayload(reportAPIError(generateErrorMessage(http.StatusBadRequest, err.Error()), startTime, http.StatusBadRequest, mon.CreateTenantMetaStr, mon.APICompleted, mon.TenantAPICompleted))
		}

		// Issue request to DAO Layer
		result, err := tenantDB.CreateTenantMeta(&data)
		if err != nil {
			return tenant_provisioning_service.NewCreateTenantMetadataInternalServerError().WithPayload(reportAPIError(fmt.Sprintf("Unable to store %s: %s", tenmod.TenantMetaStr, err.Error()), startTime, http.StatusInternalServerError, mon.CreateTenantMetaStr, mon.APICompleted, mon.TenantAPICompleted))
		}

		converted := swagmodels.JSONAPITenantMetadata{}
		err = convertToJsonapiObject(result, &converted)
		if err != nil {
			return tenant_provisioning_service.NewCreateTenantMetadataInternalServerError().WithPayload(reportAPIError(fmt.Sprintf("Unable to convert %s data to jsonapi return format: %s", tenmod.TenantMetaStr, err.Error()), startTime, http.StatusInternalServerError, mon.CreateTenantMetaStr, mon.APICompleted, mon.TenantAPICompleted))
		}

		reportAPICompletionState(startTime, http.StatusOK, mon.CreateTenantMetaStr, mon.APICompleted, mon.TenantAPICompleted)
		logger.Log.Infof("Created %s %s", tenmod.TenantMetaStr, models.AsJSONString(converted))
		return tenant_provisioning_service.NewCreateTenantMetadataOK().WithPayload(&converted)
	}
}

// HandleUpdateTenantMetadata - updates a n ingestion profile for a tenant
func HandleUpdateTenantMetadata(allowedRoles []string, tenantDB datastore.TenantServiceDatastore) func(params tenant_provisioning_service.UpdateTenantMetadataParams) middleware.Responder {
	return func(params tenant_provisioning_service.UpdateTenantMetadataParams) middleware.Responder {
		startTime := time.Now()
		incrementAPICounters(mon.APIRecieved, mon.TenantAPIRecieved)
		logger.Log.Infof("Updating %s for Tenant %s", tenmod.TenantMetaStr, params.TenantID)

		if !isRequestAuthorized(params.HTTPRequest, allowedRoles) {
			return tenant_provisioning_service.NewUpdateTenantMetadataForbidden().WithPayload(reportAPIError(fmt.Sprintf("Update %s operation not authorized for role: %s", tenmod.TenantMetaStr, params.HTTPRequest.Header.Get(XFwdUserRoles)), startTime, http.StatusForbidden, mon.UpdateTenantMetaStr, mon.APICompleted, mon.TenantAPICompleted))
		}

		// Unmarshal the request
		requestBytes, err := json.Marshal(params.Body)
		if err != nil {
			return tenant_provisioning_service.NewUpdateTenantMetadataBadRequest().WithPayload(reportAPIError(generateErrorMessage(http.StatusBadRequest, err.Error()), startTime, http.StatusBadRequest, mon.UpdateTenantMetaStr, mon.APICompleted, mon.TenantAPICompleted))
		}

		data := tenmod.Metadata{}
		err = jsonapi.Unmarshal(requestBytes, &data)
		if err != nil {
			return tenant_provisioning_service.NewUpdateTenantMetadataBadRequest().WithPayload(reportAPIError(generateErrorMessage(http.StatusBadRequest, err.Error()), startTime, http.StatusBadRequest, mon.UpdateTenantMetaStr, mon.APICompleted, mon.TenantAPICompleted))
		}

		err = data.Validate(true)
		if err != nil {
			return tenant_provisioning_service.NewUpdateTenantMetadataBadRequest().WithPayload(reportAPIError(generateErrorMessage(http.StatusBadRequest, err.Error()), startTime, http.StatusBadRequest, mon.UpdateTenantMetaStr, mon.APICompleted, mon.TenantAPICompleted))
		}

		// Issue request to DAO Layer
		result, err := tenantDB.UpdateTenantMeta(&data)
		if err != nil {
			return tenant_provisioning_service.NewUpdateTenantMetadataInternalServerError().WithPayload(reportAPIError(fmt.Sprintf("Unable to store %s: %s", tenmod.TenantMetaStr, err.Error()), startTime, http.StatusInternalServerError, mon.UpdateTenantMetaStr, mon.APICompleted, mon.TenantAPICompleted))
		}

		converted := swagmodels.JSONAPITenantMetadata{}
		err = convertToJsonapiObject(result, &converted)
		if err != nil {
			return tenant_provisioning_service.NewUpdateTenantMetadataInternalServerError().WithPayload(reportAPIError(fmt.Sprintf("Unable to convert %s data to jsonapi return format: %s", tenmod.TenantMetaStr, err.Error()), startTime, http.StatusInternalServerError, mon.UpdateTenantMetaStr, mon.APICompleted, mon.TenantAPICompleted))
		}

		reportAPICompletionState(startTime, http.StatusOK, mon.UpdateTenantMetaStr, mon.APICompleted, mon.TenantAPICompleted)
		logger.Log.Infof("Updated %s %s", tenmod.TenantMetaStr, models.AsJSONString(converted))
		return tenant_provisioning_service.NewUpdateTenantMetadataOK().WithPayload(&converted)
	}
}

// HandleGetTenantMetadata - fetch an ingestion profile for a tenant
func HandleGetTenantMetadata(allowedRoles []string, tenantDB datastore.TenantServiceDatastore) func(params tenant_provisioning_service.GetTenantMetadataParams) middleware.Responder {
	return func(params tenant_provisioning_service.GetTenantMetadataParams) middleware.Responder {
		startTime := time.Now()
		incrementAPICounters(mon.APIRecieved, mon.TenantAPIRecieved)
		logger.Log.Infof("Fetching %s for Tenant %s", tenmod.TenantMetaStr, params.TenantID)

		if !isRequestAuthorized(params.HTTPRequest, allowedRoles) {
			return tenant_provisioning_service.NewGetTenantMetadataForbidden().WithPayload(reportAPIError(fmt.Sprintf("Get %s operation not authorized for role: %s", tenmod.TenantMetaStr, params.HTTPRequest.Header.Get(XFwdUserRoles)), startTime, http.StatusForbidden, mon.GetTenantMetaStr, mon.APICompleted, mon.TenantAPICompleted))
		}

		// Issue request to DAO Layer
		result, err := tenantDB.GetTenantMeta(params.TenantID)
		if err != nil {
			return tenant_provisioning_service.NewGetTenantMetadataInternalServerError().WithPayload(reportAPIError(fmt.Sprintf("Unable to retrieve %s: %s", tenmod.TenantMetaStr, err.Error()), startTime, http.StatusInternalServerError, mon.GetTenantMetaStr, mon.APICompleted, mon.TenantAPICompleted))
		}

		converted := swagmodels.JSONAPITenantMetadata{}
		err = convertToJsonapiObject(result, &converted)
		if err != nil {
			return tenant_provisioning_service.NewGetTenantMetadataInternalServerError().WithPayload(reportAPIError(fmt.Sprintf("Unable to convert %s data to jsonapi return format: %s", tenmod.TenantMetaStr, err.Error()), startTime, http.StatusInternalServerError, mon.GetTenantMetaStr, mon.APICompleted, mon.TenantAPICompleted))
		}

		reportAPICompletionState(startTime, http.StatusOK, mon.GetTenantMetaStr, mon.APICompleted, mon.TenantAPICompleted)
		logger.Log.Infof("Retrieved %s %s", tenmod.TenantMetaStr, models.AsJSONString(converted))
		return tenant_provisioning_service.NewGetTenantMetadataOK().WithPayload(&converted)
	}
}

// HandleDeleteTenantMetadata - delete an ingestion profile for a tenant
func HandleDeleteTenantMetadata(allowedRoles []string, tenantDB datastore.TenantServiceDatastore) func(params tenant_provisioning_service.DeleteTenantMetadataParams) middleware.Responder {
	return func(params tenant_provisioning_service.DeleteTenantMetadataParams) middleware.Responder {
		startTime := time.Now()
		incrementAPICounters(mon.APIRecieved, mon.TenantAPIRecieved)
		logger.Log.Infof("Deleting %s for Tenant %s", tenmod.TenantMetaStr, params.TenantID)

		if !isRequestAuthorized(params.HTTPRequest, allowedRoles) {
			return tenant_provisioning_service.NewDeleteTenantMetadataForbidden().WithPayload(reportAPIError(fmt.Sprintf("Delete %s operation not authorized for role: %s", tenmod.TenantMetaStr, params.HTTPRequest.Header.Get(XFwdUserRoles)), startTime, http.StatusForbidden, mon.DeleteTenantMetaStr, mon.APICompleted, mon.TenantAPICompleted))
		}

		// Issue request to DAO Layer
		result, err := tenantDB.DeleteTenantMeta(params.TenantID)
		if err != nil {
			return tenant_provisioning_service.NewDeleteTenantMetadataInternalServerError().WithPayload(reportAPIError(fmt.Sprintf("Unable to delete %s: %s", tenmod.TenantMetaStr, err.Error()), startTime, http.StatusInternalServerError, mon.DeleteTenantMetaStr, mon.APICompleted, mon.TenantAPICompleted))
		}

		converted := swagmodels.JSONAPITenantMetadata{}
		err = convertToJsonapiObject(result, &converted)
		if err != nil {
			return tenant_provisioning_service.NewDeleteTenantMetadataInternalServerError().WithPayload(reportAPIError(fmt.Sprintf("Unable to convert %s data to jsonapi return format: %s", tenmod.TenantMetaStr, err.Error()), startTime, http.StatusInternalServerError, mon.DeleteTenantMetaStr, mon.APICompleted, mon.TenantAPICompleted))
		}

		reportAPICompletionState(startTime, http.StatusOK, mon.DeleteTenantMetaStr, mon.APICompleted, mon.TenantAPICompleted)
		logger.Log.Infof("Deleted %s %s", tenmod.TenantMetaStr, models.AsJSONString(converted))
		return tenant_provisioning_service.NewDeleteTenantMetadataOK().WithPayload(&converted)
	}
}

// HandlePatchTenantMetadata - patch an ingestion profile for a tenant
func HandlePatchTenantMetadata(allowedRoles []string, tenantDB datastore.TenantServiceDatastore) func(params tenant_provisioning_service.PatchTenantMetadataParams) middleware.Responder {
	return func(params tenant_provisioning_service.PatchTenantMetadataParams) middleware.Responder {
		startTime := time.Now()
		incrementAPICounters(mon.APIRecieved, mon.TenantAPIRecieved)
		logger.Log.Infof("Patching %s %s for Tenant %s", tenmod.TenantMetaStr, params.Body.Data.ID, params.TenantID)

		if !isRequestAuthorized(params.HTTPRequest, allowedRoles) {
			return tenant_provisioning_service.NewPatchTenantMetadataForbidden().WithPayload(reportAPIError(fmt.Sprintf("Patch %s operation not authorized for role: %s", tenmod.TenantMetaStr, params.HTTPRequest.Header.Get(XFwdUserRoles)), startTime, http.StatusForbidden, mon.PatchTenantMetaStr, mon.APICompleted, mon.TenantAPICompleted))
		}

		// Unmarshal the request
		requestBytes, err := json.Marshal(params.Body)
		if err != nil {
			return tenant_provisioning_service.NewPatchTenantMetadataBadRequest().WithPayload(reportAPIError(generateErrorMessage(http.StatusBadRequest, err.Error()), startTime, http.StatusBadRequest, mon.PatchTenantMetaStr, mon.APICompleted, mon.TenantAPICompleted))
		}

		data := tenmod.Metadata{}
		err = jsonapi.Unmarshal(requestBytes, &data)
		if err != nil {
			return tenant_provisioning_service.NewPatchTenantMetadataBadRequest().WithPayload(reportAPIError(generateErrorMessage(http.StatusBadRequest, err.Error()), startTime, http.StatusBadRequest, mon.PatchTenantMetaStr, mon.APICompleted, mon.TenantAPICompleted))
		}

		err = data.Validate(true)
		if err != nil {
			return tenant_provisioning_service.NewPatchTenantMetadataBadRequest().WithPayload(reportAPIError(generateErrorMessage(http.StatusBadRequest, err.Error()), startTime, http.StatusBadRequest, mon.PatchTenantMetaStr, mon.APICompleted, mon.TenantAPICompleted))
		}

		origData, err2 := tenantDB.GetTenantMeta(data.TenantID)
		if err2 != nil {
			return tenant_provisioning_service.NewPatchTenantMetadataBadRequest().WithPayload(reportAPIError(generateErrorMessage(http.StatusBadRequest, err.Error()), startTime, http.StatusBadRequest, mon.PatchTenantMetaStr, mon.APICompleted, mon.TenantAPICompleted))
		}

		errMerge := models.MergeObjWithMap(origData, requestBytes)
		if errMerge != nil {
			return tenant_provisioning_service.NewPatchTenantMetadataBadRequest().WithPayload(reportAPIError(generateErrorMessage(http.StatusBadRequest, err.Error()), startTime, http.StatusBadRequest, mon.PatchTenantMetaStr, mon.APICompleted, mon.TenantAPICompleted))
		}

		// Issue request to DAO Layer
		result, err := tenantDB.UpdateTenantMeta(origData)
		if err != nil {
			return tenant_provisioning_service.NewPatchTenantMetadataInternalServerError().WithPayload(reportAPIError(fmt.Sprintf("Unable to store %s: %s", tenmod.TenantMetaStr, err.Error()), startTime, http.StatusInternalServerError, mon.PatchTenantMetaStr, mon.APICompleted, mon.TenantAPICompleted))
		}

		converted := swagmodels.JSONAPITenantMetadata{}
		err = convertToJsonapiObject(result, &converted)
		if err != nil {
			return tenant_provisioning_service.NewPatchTenantMetadataInternalServerError().WithPayload(reportAPIError(fmt.Sprintf("Unable to convert %s data to jsonapi return format: %s", tenmod.TenantMetaStr, err.Error()), startTime, http.StatusInternalServerError, mon.PatchTenantMetaStr, mon.APICompleted, mon.TenantAPICompleted))
		}

		reportAPICompletionState(startTime, http.StatusOK, mon.PatchTenantMetaStr, mon.APICompleted, mon.TenantAPICompleted)
		logger.Log.Infof("Patched %s %s", tenmod.TenantMetaStr, models.AsJSONString(converted))
		return tenant_provisioning_service.NewPatchTenantMetadataOK().WithPayload(&converted)
	}
}
