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

// HandleCreateTenantMonitoredObject - creates a MonitoredObject for a tenant
func HandleCreateTenantMonitoredObject(allowedRoles []string, tenantDB datastore.TenantServiceDatastore) func(params tenant_provisioning_service.CreateTenantMonitoredObjectParams) middleware.Responder {
	return func(params tenant_provisioning_service.CreateTenantMonitoredObjectParams) middleware.Responder {
		startTime := time.Now()
		incrementAPICounters(mon.APIRecieved, mon.TenantAPIRecieved)
		logger.Log.Infof("Creating %s %s for Tenant %s", tenmod.TenantMonitoredObjectStr, params.Body.Data.Attributes.ObjectID, params.TenantID)

		if !isRequestAuthorized(params.HTTPRequest, allowedRoles) {
			return tenant_provisioning_service.NewCreateTenantMonitoredObjectForbidden().WithPayload(reportAPIError(fmt.Sprintf("Create %s operation not authorized for role: %s", tenmod.TenantMonitoredObjectStr, params.HTTPRequest.Header.Get(XFwdUserRoles)), startTime, http.StatusForbidden, mon.CreateMonObjStr, mon.APICompleted, mon.TenantAPICompleted))
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
		logger.Log.Infof("Updating %s %s for Tenant %s", tenmod.TenantMonitoredObjectStr, params.Body.Data.Attributes.ObjectID, params.TenantID)

		if !isRequestAuthorized(params.HTTPRequest, allowedRoles) {
			return tenant_provisioning_service.NewUpdateTenantMonitoredObjectForbidden().WithPayload(reportAPIError(fmt.Sprintf("Update %s operation not authorized for role: %s", tenmod.TenantMonitoredObjectStr, params.HTTPRequest.Header.Get(XFwdUserRoles)), startTime, http.StatusForbidden, mon.UpdateMonObjStr, mon.APICompleted, mon.TenantAPICompleted))
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
		logger.Log.Infof("Patching %s %s for Tenant %s", tenmod.TenantMonitoredObjectStr, params.Body.Data.ID, params.TenantID)

		if !isRequestAuthorized(params.HTTPRequest, allowedRoles) {
			return tenant_provisioning_service.NewPatchTenantMonitoredObjectForbidden().WithPayload(reportAPIError(fmt.Sprintf("Patch %s operation not authorized for role: %s", tenmod.TenantMonitoredObjectStr, params.HTTPRequest.Header.Get(XFwdUserRoles)), startTime, http.StatusForbidden, mon.PatchMonObjStr, mon.APICompleted, mon.TenantAPICompleted))
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
			return tenant_provisioning_service.NewGetTenantMonitoredObjectForbidden().WithPayload(reportAPIError(fmt.Sprintf("Get %s operation not authorized for role: %s", tenmod.TenantMonitoredObjectStr, params.HTTPRequest.Header.Get(XFwdUserRoles)), startTime, http.StatusForbidden, mon.GetMonObjStr, mon.APICompleted, mon.TenantAPICompleted))
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
			return tenant_provisioning_service.NewGetAllTenantMonitoredObjectsForbidden().WithPayload(reportAPIError(fmt.Sprintf("Get all %ss operation not authorized for role: %s", tenmod.TenantMonitoredObjectStr, params.HTTPRequest.Header.Get(XFwdUserRoles)), startTime, http.StatusForbidden, mon.GetAllMonObjStr, mon.APICompleted, mon.TenantAPICompleted))
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
			return tenant_provisioning_service.NewDeleteTenantMonitoredObjectForbidden().WithPayload(reportAPIError(fmt.Sprintf("Delete %s operation not authorized for role: %s", tenmod.TenantMonitoredObjectStr, params.HTTPRequest.Header.Get(XFwdUserRoles)), startTime, http.StatusForbidden, mon.DeleteMonObjStr, mon.APICompleted, mon.TenantAPICompleted))
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
		logger.Log.Infof("Retrieved %s %s", tenmod.TenantMonitoredObjectStr, models.AsJSONString(result))
		return tenant_provisioning_service.NewDeleteTenantMonitoredObjectOK().WithPayload(&converted)
	}
}

// HandleGetDomainToMonitoredObjectMap - retrieves a mapping of which domains are associated with which monitored objects
func HandleGetDomainToMonitoredObjectMap(allowedRoles []string, tenantDB datastore.TenantServiceDatastore) func(params tenant_provisioning_service.GetDomainToMonitoredObjectMapParams) middleware.Responder {
	return func(params tenant_provisioning_service.GetDomainToMonitoredObjectMapParams) middleware.Responder {
		startTime := time.Now()
		incrementAPICounters(mon.APIRecieved, mon.TenantAPIRecieved)
		logger.Log.Infof("Fetching %s for Tenant %s", tenmod.MonitoredObjectToDomainMapStr, params.TenantID)

		if !isRequestAuthorized(params.HTTPRequest, allowedRoles) {
			return tenant_provisioning_service.NewGetDomainToMonitoredObjectMapForbidden().WithPayload(reportAPIError(fmt.Sprintf("Get %s operation not authorized for role: %s", tenmod.MonitoredObjectToDomainMapStr, params.HTTPRequest.Header.Get(XFwdUserRoles)), startTime, http.StatusForbidden, mon.GetMonObjToDomMapStr, mon.APICompleted, mon.TenantAPICompleted))
		}

		// Convert the request
		data := tenmod.MonitoredObjectCountByDomainRequest{
			TenantID:  params.TenantID,
			DomainSet: params.Body.DomainSet,
			ByCount:   params.Body.ByCount,
		}

		// Issue request to DAO Layer
		result, err := tenantDB.GetMonitoredObjectToDomainMap(&data)
		if err != nil {
			return tenant_provisioning_service.NewGetDomainToMonitoredObjectMapInternalServerError().WithPayload(reportAPIError(fmt.Sprintf("Unable to retrieve %s : %s", tenmod.MonitoredObjectToDomainMapStr, err.Error()), startTime, http.StatusInternalServerError, mon.GetMonObjToDomMapStr, mon.APICompleted, mon.TenantAPICompleted))
		}

		converted := swagmodels.MonitoredObjectCountByDomainResponse{
			DomainToMonitoredObjectCountMap: result.DomainToMonitoredObjectCountMap,
			DomainToMonitoredObjectSetMap:   result.DomainToMonitoredObjectSetMap,
		}

		reportAPICompletionState(startTime, http.StatusOK, mon.DeleteMonObjStr, mon.APICompleted, mon.TenantAPICompleted)
		logger.Log.Infof("Retrieved %s %s", tenmod.TenantMonitoredObjectStr, models.AsJSONString(result))
		return tenant_provisioning_service.NewGetDomainToMonitoredObjectMapOK().WithPayload(&converted)
	}
}

// HandleBulkInsertMonitoredObjects - inserts monitored objects in bulk for a tenant
func HandleBulkInsertMonitoredObjects(allowedRoles []string, tenantDB datastore.TenantServiceDatastore) func(params tenant_provisioning_service.BulkInsertMonitoredObjectParams) middleware.Responder {
	return func(params tenant_provisioning_service.BulkInsertMonitoredObjectParams) middleware.Responder {
		startTime := time.Now()
		incrementAPICounters(mon.APIRecieved, mon.TenantAPIRecieved)
		logger.Log.Infof("Creating %s in bulk for Tenant %s", tenmod.TenantMonitoredObjectStr, params.TenantID)

		if !isRequestAuthorized(params.HTTPRequest, allowedRoles) {
			return tenant_provisioning_service.NewBulkInsertMonitoredObjectForbidden().WithPayload(reportAPIError(fmt.Sprintf("Bulk insert %s operation not authorized for role: %s", tenmod.TenantMonitoredObjectStr, params.HTTPRequest.Header.Get(XFwdUserRoles)), startTime, http.StatusForbidden, mon.BulkInsertMonObjStr, mon.APICompleted, mon.TenantAPICompleted))
		}

		logger.Log.Infof("Recieved: %s", models.AsJSONString(params.Body))
		requestBytes, err := json.Marshal(params.Body)
		if err != nil {
			return tenant_provisioning_service.NewBulkInsertMonitoredObjectBadRequest().WithPayload(reportAPIError(generateErrorMessage(http.StatusBadRequest, err.Error()), startTime, http.StatusBadRequest, mon.BulkInsertMonObjStr, mon.APICompleted, mon.TenantAPICompleted))
		}

		// Unmarshal the request
		data := []*tenmod.MonitoredObject{}
		err = json.Unmarshal(requestBytes, &data)
		if err != nil {
			return tenant_provisioning_service.NewBulkInsertMonitoredObjectBadRequest().WithPayload(reportAPIError(generateErrorMessage(http.StatusBadRequest, err.Error()), startTime, http.StatusBadRequest, mon.BulkInsertMonObjStr, mon.APICompleted, mon.TenantAPICompleted))
		}

		if len(data) == 0 {
			return tenant_provisioning_service.NewBulkInsertMonitoredObjectBadRequest().WithPayload(reportAPIError(generateErrorMessage(http.StatusBadRequest, "No Monitored Objects in provided in the request"), startTime, http.StatusBadRequest, mon.BulkInsertMonObjStr, mon.APICompleted, mon.TenantAPICompleted))
		}

		// Validate the request data
		for _, obj := range data {
			if err = obj.Validate(false); err != nil {

			}

			if obj.TenantID != params.TenantID {
				return tenant_provisioning_service.NewBulkInsertMonitoredObjectBadRequest().WithPayload(reportAPIError(generateErrorMessage(http.StatusBadRequest, err.Error()), startTime, http.StatusBadRequest, mon.BulkInsertMonObjStr, mon.APICompleted, mon.TenantAPICompleted))
			}
		}

		// Issue request to DAO Layer
		result, err := tenantDB.BulkInsertMonitoredObjects(params.TenantID, data)
		if err != nil {
			return tenant_provisioning_service.NewBulkInsertMonitoredObjectInternalServerError().WithPayload(reportAPIError(fmt.Sprintf("Unable to store %s: %s", tenmod.TenantMonitoredObjectStr, err.Error()), startTime, http.StatusInternalServerError, mon.BulkInsertMonObjStr, mon.APICompleted, mon.TenantAPICompleted))
		}

		if changeNotificationEnabled {
			NotifyMonitoredObjectCreated(params.TenantID, data...)
		}

		res, err := json.Marshal(result)
		if err != nil {
			return tenant_provisioning_service.NewBulkInsertMonitoredObjectInternalServerError().WithPayload(reportAPIError(fmt.Sprintf("Unable to serialize bulk insert %s: %s", tenmod.TenantMonitoredObjectStr, err.Error()), startTime, http.StatusInternalServerError, mon.BulkInsertMonObjStr, mon.APICompleted, mon.TenantAPICompleted))
		}

		converted := swagmodels.BulkOperationResponse{}
		err = json.Unmarshal(res, &converted)
		if err != nil {
			return tenant_provisioning_service.NewBulkInsertMonitoredObjectInternalServerError().WithPayload(reportAPIError(fmt.Sprintf("Unable to format bulk insert %s: %s", tenmod.TenantMonitoredObjectStr, err.Error()), startTime, http.StatusInternalServerError, mon.BulkInsertMonObjStr, mon.APICompleted, mon.TenantAPICompleted))
		}

		reportAPICompletionState(startTime, http.StatusOK, mon.BulkInsertMonObjStr, mon.APICompleted, mon.TenantAPICompleted)
		logger.Log.Infof("Bulk insertion of %ss complete", tenmod.TenantMonitoredObjectStr)
		return tenant_provisioning_service.NewBulkInsertMonitoredObjectOK().WithPayload(converted)
	}
}

// HandleBulkUpdateMonitoredObjects - inserts monitored objects in bulk for a tenant
func HandleBulkUpdateMonitoredObjects(allowedRoles []string, tenantDB datastore.TenantServiceDatastore) func(params tenant_provisioning_service.BulkUpdateMonitoredObjectParams) middleware.Responder {
	return func(params tenant_provisioning_service.BulkUpdateMonitoredObjectParams) middleware.Responder {
		startTime := time.Now()
		incrementAPICounters(mon.APIRecieved, mon.TenantAPIRecieved)
		logger.Log.Infof("Updating %s in bulk for Tenant %s", tenmod.TenantMonitoredObjectStr, params.TenantID)

		if !isRequestAuthorized(params.HTTPRequest, allowedRoles) {
			return tenant_provisioning_service.NewBulkUpdateMonitoredObjectForbidden().WithPayload(reportAPIError(fmt.Sprintf("Bulk update %s operation not authorized for role: %s", tenmod.TenantMonitoredObjectStr, params.HTTPRequest.Header.Get(XFwdUserRoles)), startTime, http.StatusForbidden, mon.BulkUpdateMonObjStr, mon.APICompleted, mon.TenantAPICompleted))
		}

		requestBytes, err := json.Marshal(params.Body)
		if err != nil {
			return tenant_provisioning_service.NewBulkUpdateMonitoredObjectBadRequest().WithPayload(reportAPIError(generateErrorMessage(http.StatusBadRequest, err.Error()), startTime, http.StatusBadRequest, mon.BulkUpdateMonObjStr, mon.APICompleted, mon.TenantAPICompleted))
		}

		// Unmarshal the request
		data := []*tenmod.MonitoredObject{}
		err = json.Unmarshal(requestBytes, &data)
		if err != nil {
			return tenant_provisioning_service.NewBulkUpdateMonitoredObjectBadRequest().WithPayload(reportAPIError(generateErrorMessage(http.StatusBadRequest, err.Error()), startTime, http.StatusBadRequest, mon.BulkUpdateMonObjStr, mon.APICompleted, mon.TenantAPICompleted))
		}

		// Validate the request data
		for _, obj := range data {
			if err = obj.Validate(true); err != nil || obj.TenantID != params.TenantID {
				return tenant_provisioning_service.NewBulkUpdateMonitoredObjectBadRequest().WithPayload(reportAPIError(fmt.Sprintf("Unable to Update %ss in bulk: %s", tenmod.TenantMonitoredObjectStr, err.Error()), startTime, http.StatusBadRequest, mon.BulkUpdateMonObjStr, mon.APICompleted, mon.TenantAPICompleted))
			}
			if obj.TenantID != params.TenantID {
				return tenant_provisioning_service.NewBulkUpdateMonitoredObjectBadRequest().WithPayload(reportAPIError(fmt.Sprintf("Unable to Update %ss in bulk: %s", tenmod.TenantMonitoredObjectStr, "All Monitored Objects must have Tenant ID "+params.TenantID), startTime, http.StatusBadRequest, mon.BulkUpdateMonObjStr, mon.APICompleted, mon.TenantAPICompleted))
			}

		}

		// Issue request to DAO Layer
		result, err := tenantDB.BulkUpdateMonitoredObjects(params.TenantID, data)
		if err != nil {
			return tenant_provisioning_service.NewBulkUpdateMonitoredObjectInternalServerError().WithPayload(reportAPIError(fmt.Sprintf("Unable to store %s: %s", tenmod.TenantMonitoredObjectStr, err.Error()), startTime, http.StatusInternalServerError, mon.BulkUpdateMonObjStr, mon.APICompleted, mon.TenantAPICompleted))
		}

		if changeNotificationEnabled {
			NotifyMonitoredObjectUpdated(params.TenantID, data...)
		}

		res, err := json.Marshal(result)
		if err != nil {
			return tenant_provisioning_service.NewBulkUpdateMonitoredObjectInternalServerError().WithPayload(reportAPIError(fmt.Sprintf("Unable to serialize bulk update %s: %s", tenmod.TenantMonitoredObjectStr, err.Error()), startTime, http.StatusInternalServerError, mon.BulkUpdateMonObjStr, mon.APICompleted, mon.TenantAPICompleted))
		}

		converted := swagmodels.BulkOperationResponse{}
		err = json.Unmarshal(res, &converted)
		if err != nil {
			return tenant_provisioning_service.NewBulkUpdateMonitoredObjectInternalServerError().WithPayload(reportAPIError(fmt.Sprintf("Unable to format bulk update %s: %s", tenmod.TenantMonitoredObjectStr, err.Error()), startTime, http.StatusInternalServerError, mon.BulkUpdateMonObjStr, mon.APICompleted, mon.TenantAPICompleted))
		}

		reportAPICompletionState(startTime, http.StatusOK, mon.BulkUpdateMonObjStr, mon.APICompleted, mon.TenantAPICompleted)
		logger.Log.Infof("Bulk insertion of %ss complete", tenmod.TenantMonitoredObjectStr)
		return tenant_provisioning_service.NewBulkUpdateMonitoredObjectOK().WithPayload(converted)
	}
}

func HandleBulkUpsertMonitoredObjectsMeta(allowedRoles []string, tenantDB datastore.TenantServiceDatastore) func(params tenant_provisioning_service.BulkUpsertMonitoredObjectMetaParams) middleware.Responder {
	return func(params tenant_provisioning_service.BulkUpsertMonitoredObjectMetaParams) middleware.Responder {

		startTime := time.Now()
		incrementAPICounters(mon.APIRecieved, mon.TenantAPIRecieved)
		logger.Log.Infof("Updating %s meta data in bulk for Tenant %s", tenmod.TenantMonitoredObjectStr, params.TenantID)

		tenantID := params.TenantID

		if !isRequestAuthorized(params.HTTPRequest, allowedRoles) {
			return tenant_provisioning_service.NewBulkUpsertMonitoredObjectMetaForbidden().WithPayload(reportAPIError(fmt.Sprintf("Bulk upsert %s meta operation not authorized for role: %s", tenmod.TenantMonitoredObjectStr, params.HTTPRequest.Header.Get(xFwdUserRoles)), startTime, http.StatusForbidden, mon.BulkUpsertMonObjMetaStr, mon.APICompleted, mon.TenantAPICompleted))
		}

		requestBytes, err := json.Marshal(params.Body)
		if err != nil {
			return tenant_provisioning_service.NewBulkUpsertMonitoredObjectMetaBadRequest().WithPayload(reportAPIError(generateErrorMessage(http.StatusBadRequest, err.Error()), startTime, http.StatusBadRequest, mon.BulkUpsertMonObjMetaStr, mon.APICompleted, mon.TenantAPICompleted))
		}

		// Unmarshal the request
		data := tenmod.MonitoredObjectBulkMetadata{}
		err = json.Unmarshal(requestBytes, &data)
		if err != nil {
			return tenant_provisioning_service.NewBulkUpsertMonitoredObjectMetaBadRequest().WithPayload(reportAPIError(generateErrorMessage(http.StatusBadRequest, err.Error()), startTime, http.StatusBadRequest, mon.BulkUpsertMonObjMetaStr, mon.APICompleted, mon.TenantAPICompleted))
		}

		for _, item := range data.Items {
			// Issue request to DAO Layer
			existingMonitoredObject, err := tenantDB.GetMonitoredObjectByObjectName(item.MetadataKey, tenantID)
			if err != nil {
				msg := fmt.Sprintf("Unable to retrieve %s: %s", tenmod.TenantMonitoredObjectStr, err.Error())
				return tenant_provisioning_service.NewBulkUpsertMonitoredObjectMetaInternalServerError().WithPayload(reportAPIError(generateErrorMessage(http.StatusNotFound, msg), startTime, http.StatusBadRequest, mon.BulkUpsertMonObjMetaStr, mon.APICompleted, mon.TenantAPICompleted))
			}

			logger.Log.Infof("Patching metadata for %s with name %s", tenmod.TenantMonitoredObjectStr, existingMonitoredObject.ObjectName)

			existingMonitoredObject.Meta = item.Metadata
			// Hack to emulate an external request. If this is not done, then the monitored object prefix will be added again causing a 409 conflict
			existingMonitoredObject.ID = existingMonitoredObject.MonitoredObjectID

			// Issue request to DAO Layer
			_, err = tenantDB.UpdateMonitoredObject(existingMonitoredObject)
			if err != nil {
				msg := fmt.Sprintf("Unable to store %s: %s", tenmod.TenantMonitoredObjectStr, err.Error())
				return tenant_provisioning_service.NewBulkUpsertMonitoredObjectMetaInternalServerError().WithPayload(reportAPIError(generateErrorMessage(http.StatusInternalServerError, msg), startTime, http.StatusBadRequest, mon.BulkUpsertMonObjMetaStr, mon.APICompleted, mon.TenantAPICompleted))
			}

			err = tenantDB.MonitoredObjectKeysUpdate(tenantID, existingMonitoredObject)
			if err != nil {
				msg := fmt.Sprintf("Unable to update monitored object keys %s: %s -> %s", tenmod.TenantMonitoredObjectStr, err.Error(), models.AsJSONString(existingMonitoredObject))
				return tenant_provisioning_service.NewBulkUpsertMonitoredObjectMetaInternalServerError().WithPayload(reportAPIError(generateErrorMessage(http.StatusInternalServerError, msg), startTime, http.StatusBadRequest, mon.BulkUpsertMonObjMetaStr, mon.APICompleted, mon.TenantAPICompleted))
			}

			logger.Log.Debugf("Sending notification of update to monitored object %s", existingMonitoredObject.ObjectName)
			NotifyMonitoredObjectUpdated(existingMonitoredObject.TenantID, existingMonitoredObject)

		}

		reportAPICompletionState(startTime, http.StatusOK, mon.BulkUpsertMonObjMetaStr, mon.APICompleted, mon.TenantAPICompleted)
		logger.Log.Infof("Bulk insertion of %ss meta data complete", tenmod.TenantMonitoredObjectStr)
		return tenant_provisioning_service.NewBulkUpsertMonitoredObjectMetaOK()
	}
}
