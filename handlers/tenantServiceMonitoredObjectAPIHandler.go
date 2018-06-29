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
	"github.com/accedian/skylight-aaa/utils"
	"github.com/go-openapi/runtime/middleware"
	"github.com/manyminds/api2go/jsonapi"
)

// HandleCreateTenantMonitoredObject - creates a MonitoredObject for a tenant
func HandleCreateTenantMonitoredObject(allowedRoles []string, tenantDB datastore.TenantServiceDatastore) func(params tenant_provisioning_service.CreateTenantMonitoredObjectParams) middleware.Responder {
	return func(params tenant_provisioning_service.CreateTenantMonitoredObjectParams) middleware.Responder {
		startTime := time.Now()
		incrementAPICounters(mon.APIRecieved, mon.TenantAPIRecieved)
		logger.Log.Infof("Creating %s %s for Tenant %s", tenmod.TenantMonitoredObjectStr, params.Body.Data.Attributes.ObjectID, params.TenantID)

		if !isRequestAuthorized(params.HTTPRequest, allowedRoles) {
			return tenant_provisioning_service.NewCreateTenantMonitoredObjectForbidden().WithPayload(reportAPIError(fmt.Sprintf("Create %s operation not authorized for role: %s", tenmod.TenantMonitoredObjectStr, params.HTTPRequest.Header.Get(xFwdUserRoles)), startTime, http.StatusForbidden, mon.CreateMonObjStr, mon.APICompleted, mon.TenantAPICompleted))
		}

		// Unmarshal the request
		requestBytes, err := json.Marshal(params.Body)
		if err != nil {
			return tenant_provisioning_service.NewCreateTenantMonitoredObjectBadRequest().WithPayload(reportAPIError(generateErrorMessage(http.StatusBadRequest, err.Error()), startTime, http.StatusBadRequest, mon.CreateMonObjStr, mon.APICompleted, mon.TenantAPICompleted))
		}

		data := tenmod.MonitoredObject{}
		err = jsonapi.Unmarshal(requestBytes, &data)
		if err != nil {
			return tenant_provisioning_service.NewCreateTenantMonitoredObjectBadRequest().WithPayload(reportAPIError(generateErrorMessage(http.StatusBadRequest, err.Error()), startTime, http.StatusBadRequest, mon.CreateMonObjStr, mon.APICompleted, mon.TenantAPICompleted))
		}

		err = data.Validate(false)
		if err != nil {
			return tenant_provisioning_service.NewCreateTenantMonitoredObjectBadRequest().WithPayload(reportAPIError(generateErrorMessage(http.StatusBadRequest, err.Error()), startTime, http.StatusBadRequest, mon.CreateMonObjStr, mon.APICompleted, mon.TenantAPICompleted))
		}

		// Issue request to DAO Layer
		result, err := tenantDB.CreateMonitoredObject(&data)
		if err != nil {
			return tenant_provisioning_service.NewCreateTenantMonitoredObjectInternalServerError().WithPayload(reportAPIError(fmt.Sprintf("Unable to store %s: %s", tenmod.TenantMonitoredObjectStr, err.Error()), startTime, http.StatusInternalServerError, mon.CreateMonObjStr, mon.APICompleted, mon.TenantAPICompleted))
		}

		if changeNotificationEnabled {
			NotifyMonitoredObjectCreated(data.TenantID, &data)
		}

		converted := swagmodels.JSONAPITenantMonitoredObject{}
		err = convertToJsonapiObject(result, &converted)
		if err != nil {
			return tenant_provisioning_service.NewCreateTenantMonitoredObjectInternalServerError().WithPayload(reportAPIError(fmt.Sprintf("Unable to convert %s data to jsonapi return format: %s", tenmod.TenantMonitoredObjectStr, err.Error()), startTime, http.StatusInternalServerError, mon.CreateMonObjStr, mon.APICompleted, mon.TenantAPICompleted))
		}

		reportAPICompletionState(startTime, http.StatusOK, mon.CreateMonObjStr, mon.APICompleted, mon.TenantAPICompleted)
		logger.Log.Infof("Created %s %s", tenmod.TenantMonitoredObjectStr, models.AsJSONString(converted))
		return tenant_provisioning_service.NewCreateTenantMonitoredObjectOK().WithPayload(&converted)
	}
}

// HandleUpdateTenantMonitoredObject - updates a MonitoredObject for a tenant
func HandleUpdateTenantMonitoredObject(allowedRoles []string, tenantDB datastore.TenantServiceDatastore) func(params tenant_provisioning_service.UpdateTenantMonitoredObjectParams) middleware.Responder {
	return func(params tenant_provisioning_service.UpdateTenantMonitoredObjectParams) middleware.Responder {
		startTime := time.Now()
		incrementAPICounters(mon.APIRecieved, mon.TenantAPIRecieved)
		logger.Log.Infof("Updating %s %s for Tenant", tenmod.TenantMonitoredObjectStr, params.Body.Data.Attributes.ObjectID, params.TenantID)

		if !isRequestAuthorized(params.HTTPRequest, allowedRoles) {
			return tenant_provisioning_service.NewUpdateTenantMonitoredObjectForbidden().WithPayload(reportAPIError(fmt.Sprintf("Update %s operation not authorized for role: %s", tenmod.TenantMonitoredObjectStr, params.HTTPRequest.Header.Get(xFwdUserRoles)), startTime, http.StatusForbidden, mon.UpdateMonObjStr, mon.APICompleted, mon.TenantAPICompleted))
		}

		// Unmarshal the request
		requestBytes, err := json.Marshal(params.Body)
		if err != nil {
			return tenant_provisioning_service.NewUpdateTenantMonitoredObjectBadRequest().WithPayload(reportAPIError(generateErrorMessage(http.StatusBadRequest, err.Error()), startTime, http.StatusBadRequest, mon.UpdateMonObjStr, mon.APICompleted, mon.TenantAPICompleted))
		}

		data := tenmod.MonitoredObject{}
		err = jsonapi.Unmarshal(requestBytes, &data)
		if err != nil {
			return tenant_provisioning_service.NewUpdateTenantMonitoredObjectBadRequest().WithPayload(reportAPIError(generateErrorMessage(http.StatusBadRequest, err.Error()), startTime, http.StatusBadRequest, mon.UpdateMonObjStr, mon.APICompleted, mon.TenantAPICompleted))
		}

		err = data.Validate(true)
		if err != nil {
			return tenant_provisioning_service.NewUpdateTenantMonitoredObjectBadRequest().WithPayload(reportAPIError(generateErrorMessage(http.StatusBadRequest, err.Error()), startTime, http.StatusBadRequest, mon.UpdateMonObjStr, mon.APICompleted, mon.TenantAPICompleted))
		}

		// Issue request to DAO Layer
		result, err := tenantDB.UpdateMonitoredObject(&data)
		if err != nil {
			return tenant_provisioning_service.NewUpdateTenantMonitoredObjectInternalServerError().WithPayload(reportAPIError(fmt.Sprintf("Unable to store %s: %s", tenmod.TenantMonitoredObjectStr, err.Error()), startTime, http.StatusInternalServerError, mon.UpdateMonObjStr, mon.APICompleted, mon.TenantAPICompleted))
		}

		if changeNotificationEnabled {
			NotifyMonitoredObjectUpdated(data.TenantID, &data)
		}

		converted := swagmodels.JSONAPITenantMonitoredObject{}
		err = convertToJsonapiObject(result, &converted)
		if err != nil {
			return tenant_provisioning_service.NewUpdateTenantMonitoredObjectInternalServerError().WithPayload(reportAPIError(fmt.Sprintf("Unable to convert %s data to jsonapi return format: %s", tenmod.TenantMonitoredObjectStr, err.Error()), startTime, http.StatusInternalServerError, mon.UpdateMonObjStr, mon.APICompleted, mon.TenantAPICompleted))
		}

		reportAPICompletionState(startTime, http.StatusOK, mon.UpdateMonObjStr, mon.APICompleted, mon.TenantAPICompleted)
		logger.Log.Infof("Updated %s %s", tenmod.TenantMonitoredObjectStr, models.AsJSONString(converted))
		return tenant_provisioning_service.NewUpdateTenantMonitoredObjectOK().WithPayload(&converted)
	}
}

// HandlePatchTenantMonitoredObject - patches a MonitoredObject for a tenant
func HandlePatchTenantMonitoredObject(allowedRoles []string, tenantDB datastore.TenantServiceDatastore) func(params tenant_provisioning_service.PatchTenantMonitoredObjectParams) middleware.Responder {
	return func(params tenant_provisioning_service.PatchTenantMonitoredObjectParams) middleware.Responder {
		startTime := time.Now()
		incrementAPICounters(mon.APIRecieved, mon.TenantAPIRecieved)
		logger.Log.Infof("Patching %s %s for Tenant", tenmod.TenantMonitoredObjectStr, params.Body.Data.ID, params.TenantID)

		if !isRequestAuthorized(params.HTTPRequest, allowedRoles) {
			return tenant_provisioning_service.NewPatchTenantMonitoredObjectForbidden().WithPayload(reportAPIError(fmt.Sprintf("Patch %s operation not authorized for role: %s", tenmod.TenantMonitoredObjectStr, params.HTTPRequest.Header.Get(xFwdUserRoles)), startTime, http.StatusForbidden, mon.PatchMonObjStr, mon.APICompleted, mon.TenantAPICompleted))
		}

		// Unmarshal the request
		requestBytes, err := json.Marshal(params.Body)
		if err != nil {
			return tenant_provisioning_service.NewPatchTenantMonitoredObjectBadRequest().WithPayload(reportAPIError(generateErrorMessage(http.StatusBadRequest, err.Error()), startTime, http.StatusBadRequest, mon.PatchMonObjStr, mon.APICompleted, mon.TenantAPICompleted))
		}

		// Model
		data := tenmod.MonitoredObject{}
		err = jsonapi.Unmarshal(requestBytes, &data)
		if err != nil {
			return tenant_provisioning_service.NewPatchTenantMonitoredObjectBadRequest().WithPayload(reportAPIError(generateErrorMessage(http.StatusBadRequest, err.Error()), startTime, http.StatusBadRequest, mon.PatchMonObjStr, mon.APICompleted, mon.TenantAPICompleted))
		}

		// This only checks if the ID&REV is set.
		err = data.Validate(true)
		if err != nil {
			return tenant_provisioning_service.NewPatchTenantMonitoredObjectBadRequest().WithPayload(reportAPIError(generateErrorMessage(http.StatusBadRequest, err.Error()), startTime, http.StatusBadRequest, mon.PatchMonObjStr, mon.APICompleted, mon.TenantAPICompleted))
		}

		// Issue request to DAO Layer
		oldMonitoredObject, err := tenantDB.GetMonitoredObject(data.TenantID, data.ID)
		if err != nil {
			return tenant_provisioning_service.NewPatchTenantMonitoredObjectInternalServerError().WithPayload(reportAPIError(fmt.Sprintf("Unable to retrieve %s: %s", tenmod.TenantMonitoredObjectStr, err.Error()), startTime, http.StatusInternalServerError, mon.PatchMonObjStr, mon.APICompleted, mon.TenantAPICompleted))
		}

		errMerge := models.MergeObjWithMap(oldMonitoredObject, requestBytes)
		if errMerge != nil {
			return tenant_provisioning_service.NewPatchTenantMonitoredObjectBadRequest().WithPayload(reportAPIError(generateErrorMessage(http.StatusBadRequest, err.Error()), startTime, http.StatusBadRequest, mon.PatchMonObjStr, mon.APICompleted, mon.TenantAPICompleted))
		}

		// This only checks if the ID&REV is set.
		err = oldMonitoredObject.Validate(true)
		if err != nil {
			return tenant_provisioning_service.NewPatchTenantMonitoredObjectBadRequest().WithPayload(reportAPIError(generateErrorMessage(http.StatusBadRequest, err.Error()), startTime, http.StatusBadRequest, mon.PatchMonObjStr, mon.APICompleted, mon.TenantAPICompleted))
		}

		// Issue request to DAO Layer
		result, err := tenantDB.UpdateMonitoredObject(oldMonitoredObject)
		if err != nil {
			return tenant_provisioning_service.NewPatchTenantMonitoredObjectInternalServerError().WithPayload(reportAPIError(fmt.Sprintf("Unable to store %s: %s", tenmod.TenantMonitoredObjectStr, err.Error()), startTime, http.StatusInternalServerError, mon.PatchMonObjStr, mon.APICompleted, mon.TenantAPICompleted))
		}

		if changeNotificationEnabled {
			NotifyMonitoredObjectUpdated(oldMonitoredObject.TenantID, oldMonitoredObject)
		}

		converted := swagmodels.JSONAPITenantMonitoredObject{}
		err = convertToJsonapiObject(result, &converted)
		if err != nil {
			return tenant_provisioning_service.NewPatchTenantMonitoredObjectInternalServerError().WithPayload(reportAPIError(fmt.Sprintf("Unable to convert %s data to jsonapi return format: %s", tenmod.TenantMonitoredObjectStr, err.Error()), startTime, http.StatusInternalServerError, mon.PatchMonObjStr, mon.APICompleted, mon.TenantAPICompleted))
		}

		reportAPICompletionState(startTime, http.StatusOK, mon.PatchMonObjStr, mon.APICompleted, mon.TenantAPICompleted)
		logger.Log.Infof("Updated %s %s", tenmod.TenantMonitoredObjectStr, models.AsJSONString(converted))
		return tenant_provisioning_service.NewPatchTenantMonitoredObjectOK().WithPayload(&converted)
	}
}

// HandleGetTenantMonitoredObject - fetch a MonitoredObject for a tenant
func HandleGetTenantMonitoredObject(allowedRoles []string, tenantDB datastore.TenantServiceDatastore) func(params tenant_provisioning_service.GetTenantMonitoredObjectParams) middleware.Responder {
	return func(params tenant_provisioning_service.GetTenantMonitoredObjectParams) middleware.Responder {
		startTime := time.Now()
		incrementAPICounters(mon.APIRecieved, mon.TenantAPIRecieved)
		logger.Log.Infof("Fetching %s %s for Tenant %s", tenmod.TenantMonitoredObjectStr, params.MonObjID, params.TenantID)

		if !isRequestAuthorized(params.HTTPRequest, allowedRoles) {
			return tenant_provisioning_service.NewGetTenantMonitoredObjectForbidden().WithPayload(reportAPIError(fmt.Sprintf("Get %s operation not authorized for role: %s", tenmod.TenantMonitoredObjectStr, params.HTTPRequest.Header.Get(xFwdUserRoles)), startTime, http.StatusForbidden, mon.GetMonObjStr, mon.APICompleted, mon.TenantAPICompleted))
		}

		// Issue request to DAO Layer
		result, err := tenantDB.GetMonitoredObject(params.TenantID, params.MonObjID)
		if err != nil {
			return tenant_provisioning_service.NewGetTenantMonitoredObjectInternalServerError().WithPayload(reportAPIError(fmt.Sprintf("Unable to retrieve %s: %s", tenmod.TenantMonitoredObjectStr, err.Error()), startTime, http.StatusInternalServerError, mon.GetMonObjStr, mon.APICompleted, mon.TenantAPICompleted))
		}

		converted := swagmodels.JSONAPITenantMonitoredObject{}
		err = convertToJsonapiObject(result, &converted)
		if err != nil {
			return tenant_provisioning_service.NewGetTenantMonitoredObjectInternalServerError().WithPayload(reportAPIError(fmt.Sprintf("Unable to convert %s data to jsonapi return format: %s", tenmod.TenantMonitoredObjectStr, err.Error()), startTime, http.StatusInternalServerError, mon.GetMonObjStr, mon.APICompleted, mon.TenantAPICompleted))
		}

		reportAPICompletionState(startTime, http.StatusOK, mon.GetMonObjStr, mon.APICompleted, mon.TenantAPICompleted)
		logger.Log.Infof("Retrieved %s %s", tenmod.TenantMonitoredObjectStr, models.AsJSONString(converted))
		return tenant_provisioning_service.NewGetTenantMonitoredObjectOK().WithPayload(&converted)
	}
}

// HandleGetAllTenantMonitoredObjects - fetch all MonitoredObjects for a tenant
func HandleGetAllTenantMonitoredObjects(allowedRoles []string, tenantDB datastore.TenantServiceDatastore) func(params tenant_provisioning_service.GetAllTenantMonitoredObjectsParams) middleware.Responder {
	return func(params tenant_provisioning_service.GetAllTenantMonitoredObjectsParams) middleware.Responder {
		startTime := time.Now()
		incrementAPICounters(mon.APIRecieved, mon.TenantAPIRecieved)
		logger.Log.Infof("Fetching all %ss for Tenant %s", tenmod.TenantMonitoredObjectStr, params.TenantID)

		if !isRequestAuthorized(params.HTTPRequest, allowedRoles) {
			return tenant_provisioning_service.NewGetAllTenantMonitoredObjectsForbidden().WithPayload(reportAPIError(fmt.Sprintf("Get all %ss operation not authorized for role: %s", tenmod.TenantMonitoredObjectStr, params.HTTPRequest.Header.Get(xFwdUserRoles)), startTime, http.StatusForbidden, mon.GetAllMonObjStr, mon.APICompleted, mon.TenantAPICompleted))
		}

		// Issue request to DAO Layer
		result, err := tenantDB.GetAllMonitoredObjects(params.TenantID)
		if err != nil {
			return tenant_provisioning_service.NewGetAllTenantMonitoredObjectsInternalServerError().WithPayload(reportAPIError(fmt.Sprintf("Unable to retrieve %s list: %s", tenmod.TenantMonitoredObjectStr, err.Error()), startTime, http.StatusInternalServerError, mon.GetAllMonObjStr, mon.APICompleted, mon.TenantAPICompleted))
		}

		converted := swagmodels.JSONAPITenantMonitoredObjectList{}
		err = convertToJsonapiObject(result, &converted)
		if err != nil {
			return tenant_provisioning_service.NewGetAllTenantMonitoredObjectsInternalServerError().WithPayload(reportAPIError(fmt.Sprintf("Unable to convert %s list data to jsonapi return format: %s", tenmod.TenantMonitoredObjectStr, err.Error()), startTime, http.StatusInternalServerError, mon.GetAllMonObjStr, mon.APICompleted, mon.TenantAPICompleted))
		}

		reportAPICompletionState(startTime, http.StatusOK, mon.GetAllMonObjStr, mon.APICompleted, mon.TenantAPICompleted)
		logger.Log.Infof("Retrieved %d %ss", len(result), tenmod.TenantMonitoredObjectStr)
		return tenant_provisioning_service.NewGetAllTenantMonitoredObjectsOK().WithPayload(&converted)
	}
}

// HandleDeleteTenantMonitoredObject - delete a MonitoredObject for a tenant
func HandleDeleteTenantMonitoredObject(allowedRoles []string, tenantDB datastore.TenantServiceDatastore) func(params tenant_provisioning_service.DeleteTenantMonitoredObjectParams) middleware.Responder {
	return func(params tenant_provisioning_service.DeleteTenantMonitoredObjectParams) middleware.Responder {
		startTime := time.Now()
		incrementAPICounters(mon.APIRecieved, mon.TenantAPIRecieved)
		logger.Log.Infof("Deleting %s %s for Tenant %s", tenmod.TenantMonitoredObjectStr, params.MonObjID, params.TenantID)

		if !isRequestAuthorized(params.HTTPRequest, allowedRoles) {
			return tenant_provisioning_service.NewDeleteTenantMonitoredObjectForbidden().WithPayload(reportAPIError(fmt.Sprintf("Delete %s operation not authorized for role: %s", tenmod.TenantMonitoredObjectStr, params.HTTPRequest.Header.Get(xFwdUserRoles)), startTime, http.StatusForbidden, mon.DeleteMonObjStr, mon.APICompleted, mon.TenantAPICompleted))
		}

		// Issue request to DAO Layer
		result, err := tenantDB.DeleteMonitoredObject(params.TenantID, params.MonObjID)
		if err != nil {
			return tenant_provisioning_service.NewDeleteTenantMonitoredObjectInternalServerError().WithPayload(reportAPIError(fmt.Sprintf("Unable to retrieve %s: %s", tenmod.TenantMonitoredObjectStr, err.Error()), startTime, http.StatusInternalServerError, mon.DeleteMonObjStr, mon.APICompleted, mon.TenantAPICompleted))
		}

		if changeNotificationEnabled {
			NotifyMonitoredObjectDeleted(params.TenantID, result)
		}

		converted := swagmodels.JSONAPITenantMonitoredObject{}
		err = convertToJsonapiObject(result, &converted)
		if err != nil {
			return tenant_provisioning_service.NewDeleteTenantMonitoredObjectInternalServerError().WithPayload(reportAPIError(fmt.Sprintf("Unable to convert %s list data to jsonapi return format: %s", tenmod.TenantMonitoredObjectStr, err.Error()), startTime, http.StatusInternalServerError, mon.DeleteMonObjStr, mon.APICompleted, mon.TenantAPICompleted))
		}

		reportAPICompletionState(startTime, http.StatusOK, mon.DeleteMonObjStr, mon.APICompleted, mon.TenantAPICompleted)
		logger.Log.Infof("Retrieved %s %s", tenmod.TenantMonitoredObjectStr, utils.AsJSONString(result))
		return tenant_provisioning_service.NewDeleteTenantMonitoredObjectOK().WithPayload(&converted)
	}
}
