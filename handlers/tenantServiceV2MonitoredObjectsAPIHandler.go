package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/accedian/adh-gather/datastore"
	"github.com/accedian/adh-gather/logger"
	"github.com/accedian/adh-gather/models"
	admmod "github.com/accedian/adh-gather/models/admin"
	"github.com/accedian/adh-gather/models/common"
	tenmod "github.com/accedian/adh-gather/models/tenant"
	mon "github.com/accedian/adh-gather/monitoring"
	"github.com/accedian/adh-gather/restapi/operations/tenant_provisioning_service_v2"
	"github.com/accedian/adh-gather/swagmodels"
	"github.com/go-openapi/runtime/middleware"
	"github.com/manyminds/api2go/jsonapi"
)

// HandleGetAllMonitoredObjectsV2 - retrieve all MonitoredObjects for a Tenant
func HandleGetAllMonitoredObjectsV2(allowedRoles []string, tenantDB datastore.TenantServiceDatastore) func(params tenant_provisioning_service_v2.GetAllMonitoredObjectsV2Params) middleware.Responder {
	return func(params tenant_provisioning_service_v2.GetAllMonitoredObjectsV2Params) middleware.Responder {
		// Do the work
		startTime, responseCode, response, err := doGetAllMonitoredObjectsV2(allowedRoles, tenantDB, params)

		// Success Response
		if responseCode == http.StatusOK {
			return tenant_provisioning_service_v2.NewGetAllMonitoredObjectsV2OK().WithPayload(response)
		}

		// Error Responses
		errorMessage := reportAPIError(err.Error(), startTime, responseCode, mon.GetAllMonObjStr, mon.APICompleted, mon.TenantAPICompleted)
		switch responseCode {
		case http.StatusForbidden:
			return tenant_provisioning_service_v2.NewGetAllMonitoredObjectsV2Forbidden().WithPayload(errorMessage)
		case http.StatusNotFound:
			return tenant_provisioning_service_v2.NewGetAllMonitoredObjectsV2NotFound().WithPayload(errorMessage)
		default:
			return tenant_provisioning_service_v2.NewGetAllMonitoredObjectsV2InternalServerError().WithPayload(errorMessage)

		}
	}
}

// HandleCreateMonitoredObjectV2 - create a new Monitored Object
func HandleCreateMonitoredObjectV2(allowedRoles []string, tenantDB datastore.TenantServiceDatastore) func(params tenant_provisioning_service_v2.CreateMonitoredObjectV2Params) middleware.Responder {
	return func(params tenant_provisioning_service_v2.CreateMonitoredObjectV2Params) middleware.Responder {

		// Do the work
		startTime, responseCode, response, err := doCreateMonitoredObjectV2(allowedRoles, tenantDB, params)

		// Success Response
		if responseCode == http.StatusCreated {
			return tenant_provisioning_service_v2.NewCreateMonitoredObjectV2Created().WithPayload(response)
		}

		// Error Responses
		errorMessage := reportAPIError(err.Error(), startTime, responseCode, mon.CreateMonObjStr, mon.APICompleted, mon.TenantAPICompleted)
		switch responseCode {
		case http.StatusForbidden:
			return tenant_provisioning_service_v2.NewCreateMonitoredObjectV2Forbidden().WithPayload(errorMessage)
		case http.StatusBadRequest:
			return tenant_provisioning_service_v2.NewCreateMonitoredObjectV2BadRequest().WithPayload(errorMessage)
		case http.StatusConflict:
			return tenant_provisioning_service_v2.NewCreateMonitoredObjectV2Conflict().WithPayload(errorMessage)
		default:
			return tenant_provisioning_service_v2.NewCreateMonitoredObjectV2InternalServerError().WithPayload(errorMessage)
		}
	}
}

// HandleGetMonitoredObjectV2 - retrieve a Threshold Profile by the ID.
func HandleGetMonitoredObjectV2(allowedRoles []string, tenantDB datastore.TenantServiceDatastore) func(params tenant_provisioning_service_v2.GetMonitoredObjectV2Params) middleware.Responder {
	return func(params tenant_provisioning_service_v2.GetMonitoredObjectV2Params) middleware.Responder {

		// Do the work
		startTime, responseCode, response, err := doGetMonitoredObjectV2(allowedRoles, tenantDB, params)

		// Success Response
		if responseCode == http.StatusOK {
			return tenant_provisioning_service_v2.NewGetMonitoredObjectV2OK().WithPayload(response)
		}

		// Error Responses
		errorMessage := reportAPIError(err.Error(), startTime, responseCode, mon.GetMonObjStr, mon.APICompleted, mon.TenantAPICompleted)
		switch responseCode {
		case http.StatusForbidden:
			return tenant_provisioning_service_v2.NewGetMonitoredObjectV2Forbidden().WithPayload(errorMessage)
		case http.StatusNotFound:
			return tenant_provisioning_service_v2.NewGetMonitoredObjectV2NotFound().WithPayload(errorMessage)
		default:
			return tenant_provisioning_service_v2.NewGetMonitoredObjectV2InternalServerError().WithPayload(errorMessage)
		}

	}
}

// HandleUpdateMonitoredObjectV2 - update a MonitoredObject record
func HandleUpdateMonitoredObjectV2(allowedRoles []string, tenantDB datastore.TenantServiceDatastore) func(params tenant_provisioning_service_v2.UpdateMonitoredObjectV2Params) middleware.Responder {
	return func(params tenant_provisioning_service_v2.UpdateMonitoredObjectV2Params) middleware.Responder {

		// Do the work
		startTime, responseCode, response, err := doUpdateMonitoredObjectV2(allowedRoles, tenantDB, params)

		// Success Response
		if responseCode == http.StatusOK {
			return tenant_provisioning_service_v2.NewUpdateMonitoredObjectV2OK().WithPayload(response)
		}

		// Error Responses
		errorMessage := reportAPIError(err.Error(), startTime, responseCode, mon.UpdateMonObjStr, mon.APICompleted, mon.TenantAPICompleted)
		switch responseCode {
		case http.StatusForbidden:
			return tenant_provisioning_service_v2.NewUpdateMonitoredObjectV2Forbidden().WithPayload(errorMessage)
		case http.StatusBadRequest:
			return tenant_provisioning_service_v2.NewUpdateMonitoredObjectV2BadRequest().WithPayload(errorMessage)
		case http.StatusNotFound:
			return tenant_provisioning_service_v2.NewUpdateMonitoredObjectV2NotFound().WithPayload(errorMessage)
		case http.StatusConflict:
			return tenant_provisioning_service_v2.NewUpdateMonitoredObjectV2Conflict().WithPayload(errorMessage)
		default:
			return tenant_provisioning_service_v2.NewUpdateMonitoredObjectV2InternalServerError().WithPayload(errorMessage)
		}
	}
}

// HandleDeleteMonitoredObjectV2 - delete a Threshold Profile by the ID.
func HandleDeleteMonitoredObjectV2(allowedRoles []string, tenantDB datastore.TenantServiceDatastore) func(params tenant_provisioning_service_v2.DeleteMonitoredObjectV2Params) middleware.Responder {
	return func(params tenant_provisioning_service_v2.DeleteMonitoredObjectV2Params) middleware.Responder {

		// Do the work
		startTime, responseCode, response, err := doDeleteMonitoredObjectV2(allowedRoles, tenantDB, params)

		// Success Response
		if responseCode == http.StatusOK {
			return tenant_provisioning_service_v2.NewDeleteMonitoredObjectV2OK().WithPayload(response)
		}

		// Error Responses
		errorMessage := reportAPIError(err.Error(), startTime, responseCode, mon.DeleteMonObjStr, mon.APICompleted, mon.TenantAPICompleted)
		switch responseCode {
		case http.StatusForbidden:
			return tenant_provisioning_service_v2.NewDeleteMonitoredObjectV2Forbidden().WithPayload(errorMessage)
		case http.StatusNotFound:
			return tenant_provisioning_service_v2.NewDeleteMonitoredObjectV2NotFound().WithPayload(errorMessage)
		default:
			return tenant_provisioning_service_v2.NewDeleteMonitoredObjectV2InternalServerError().WithPayload(errorMessage)
		}

	}
}

// HandleBulkCreateMonitoredObjectsV2 - insert more than 1 Monitored Object in a single request
func HandleBulkCreateMonitoredObjectsV2(allowedRoles []string, tenantDB datastore.TenantServiceDatastore) func(params tenant_provisioning_service_v2.BulkInsertMonitoredObjectsV2Params) middleware.Responder {
	return func(params tenant_provisioning_service_v2.BulkInsertMonitoredObjectsV2Params) middleware.Responder {
		// Do the work
		startTime, responseCode, response, err := doBulkInsertMonitoredObjectsV2(allowedRoles, tenantDB, params)

		// Success Response
		if responseCode == http.StatusOK {
			return tenant_provisioning_service_v2.NewBulkInsertMonitoredObjectsV2OK().WithPayload(response)
		}

		// Error Responses
		errorMessage := reportAPIError(err.Error(), startTime, responseCode, mon.BulkInsertMonObjStr, mon.APICompleted, mon.TenantAPICompleted)
		switch responseCode {
		case http.StatusForbidden:
			return tenant_provisioning_service_v2.NewBulkInsertMonitoredObjectsV2Forbidden().WithPayload(errorMessage)
		case http.StatusBadRequest:
			return tenant_provisioning_service_v2.NewBulkInsertMonitoredObjectsV2BadRequest().WithPayload(errorMessage)
		default:
			return tenant_provisioning_service_v2.NewBulkInsertMonitoredObjectsV2InternalServerError().WithPayload(errorMessage)

		}
	}
}

// HandleBulkUpdateMonitoredObjectsV2 - update more than 1 Monitored Object in a single request
func HandleBulkUpdateMonitoredObjectsV2(allowedRoles []string, tenantDB datastore.TenantServiceDatastore) func(params tenant_provisioning_service_v2.BulkUpdateMonitoredObjectsV2Params) middleware.Responder {
	return func(params tenant_provisioning_service_v2.BulkUpdateMonitoredObjectsV2Params) middleware.Responder {
		// Do the work
		startTime, responseCode, response, err := doBulkUpdateMonitoredObjectsV2(allowedRoles, tenantDB, params)

		// Success Response
		if responseCode == http.StatusOK {
			return tenant_provisioning_service_v2.NewBulkUpdateMonitoredObjectsV2OK().WithPayload(response)
		}

		// Error Responses
		errorMessage := reportAPIError(err.Error(), startTime, responseCode, mon.BulkUpdateMonObjStr, mon.APICompleted, mon.TenantAPICompleted)
		switch responseCode {
		case http.StatusForbidden:
			return tenant_provisioning_service_v2.NewBulkUpdateMonitoredObjectsV2Forbidden().WithPayload(errorMessage)
		case http.StatusBadRequest:
			return tenant_provisioning_service_v2.NewBulkUpdateMonitoredObjectsV2BadRequest().WithPayload(errorMessage)
		default:
			return tenant_provisioning_service_v2.NewBulkUpdateMonitoredObjectsV2InternalServerError().WithPayload(errorMessage)

		}
	}
}

func doCreateMonitoredObjectV2(allowedRoles []string, tenantDB datastore.TenantServiceDatastore, params tenant_provisioning_service_v2.CreateMonitoredObjectV2Params) (time.Time, int, *swagmodels.MonitoredObjectResponse, error) {
	tenantID := params.HTTPRequest.Header.Get(XFwdTenantId)
	isAuthorized, startTime := authorizeRequest(fmt.Sprintf("Creating %s for tenant %s", tenmod.TenantMonitoredObjectStr, tenantID), params.HTTPRequest, allowedRoles, mon.APIRecieved, mon.TenantAPIRecieved)

	if !isAuthorized {
		return startTime, http.StatusForbidden, nil, fmt.Errorf("Create %s operation not authorized for role: %s", tenmod.TenantMonitoredObjectStr, params.HTTPRequest.Header.Get(XFwdUserRoles))
	}

	// Unmarshal the request
	requestBytes, err := json.Marshal(params.Body)
	if err != nil {
		return startTime, http.StatusBadRequest, nil, err
	}

	data := tenmod.MonitoredObject{}
	err = jsonapi.Unmarshal(requestBytes, &data)
	if err != nil {
		return startTime, http.StatusBadRequest, nil, err
	}

	data.TenantID = tenantID

	// Issue request to DAO Layer to Create Record
	result, err := tenantDB.CreateMonitoredObject(&data)
	if err != nil {
		if strings.Contains(err.Error(), datastore.ConflictErrorStr) {
			return startTime, http.StatusConflict, nil, err
		}
		return startTime, http.StatusInternalServerError, nil, fmt.Errorf("Unable to store %s: %s", tenmod.TenantMonitoredObjectStr, err.Error())
	}

	if changeNotificationEnabled {
		NotifyMonitoredObjectUpdated(tenantID, &data)
	}

	converted := swagmodels.MonitoredObjectResponse{}
	err = convertToJsonapiObject(result, &converted)
	if err != nil {
		return startTime, http.StatusInternalServerError, nil, fmt.Errorf("Unable to convert %s data to jsonapi return format: %s", tenmod.TenantMonitoredObjectStr, err.Error())
	}

	reportAPICompletionState(startTime, http.StatusOK, mon.CreateThrPrfStr, mon.APICompleted, mon.TenantAPICompleted)
	logger.Log.Infof("Created %s %s", tenmod.TenantMonitoredObjectStr, models.AsJSONString(converted))
	return startTime, http.StatusCreated, &converted, nil
}

func doGetMonitoredObjectV2(allowedRoles []string, tenantDB datastore.TenantServiceDatastore, params tenant_provisioning_service_v2.GetMonitoredObjectV2Params) (time.Time, int, *swagmodels.MonitoredObjectResponse, error) {
	tenantID := params.HTTPRequest.Header.Get(XFwdTenantId)
	isAuthorized, startTime := authorizeRequest(fmt.Sprintf("Fetching %s %s for %s %s", tenmod.TenantMonitoredObjectStr, params.MonObjID, admmod.TenantStr, tenantID), params.HTTPRequest, allowedRoles, mon.APIRecieved, mon.AdminAPIRecieved)

	if !isAuthorized {
		return startTime, http.StatusForbidden, nil, fmt.Errorf("Fetch %s operation not authorized for role: %s", tenmod.TenantMonitoredObjectStr, params.HTTPRequest.Header.Get(XFwdUserRoles))
	}

	// Issue request to DAO Layer
	result, err := tenantDB.GetMonitoredObject(tenantID, params.MonObjID)
	if err != nil {
		if strings.Contains(err.Error(), datastore.NotFoundStr) {
			return startTime, http.StatusNotFound, nil, err
		}

		return startTime, http.StatusInternalServerError, nil, fmt.Errorf("Unable to retrieve %s: %s", tenmod.TenantMonitoredObjectStr, err.Error())
	}

	converted := swagmodels.MonitoredObjectResponse{}
	err = convertToJsonapiObject(result, &converted)
	if err != nil {
		return startTime, http.StatusInternalServerError, nil, fmt.Errorf("Unable to convert %s data to jsonapi return format: %s", tenmod.TenantMonitoredObjectStr, err.Error())
	}

	reportAPICompletionState(startTime, http.StatusOK, mon.GetThrPrfStr, mon.APICompleted, mon.TenantAPICompleted)
	logger.Log.Infof("Retrieved %s %s", tenmod.TenantMonitoredObjectStr, models.AsJSONString(converted))
	return startTime, http.StatusOK, &converted, nil
}

func doUpdateMonitoredObjectV2(allowedRoles []string, tenantDB datastore.TenantServiceDatastore, params tenant_provisioning_service_v2.UpdateMonitoredObjectV2Params) (time.Time, int, *swagmodels.MonitoredObjectResponse, error) {
	tenantID := params.HTTPRequest.Header.Get(XFwdTenantId)
	isAuthorized, startTime := authorizeRequest(fmt.Sprintf("Updating %s %s for %s %s", tenmod.TenantMonitoredObjectStr, params.MonObjID, admmod.TenantStr, tenantID), params.HTTPRequest, allowedRoles, mon.APIRecieved, mon.AdminAPIRecieved)

	if !isAuthorized {
		return startTime, http.StatusForbidden, nil, fmt.Errorf("Update %s operation not authorized for role: %s", tenmod.TenantMonitoredObjectStr, params.HTTPRequest.Header.Get(XFwdUserRoles))
	}

	// Retrieve tne request bytes from the payload in order to convert it to a map
	patchRequestBytes, err := json.Marshal(params.Body)
	if err != nil || params.Body == nil {
		return startTime, http.StatusBadRequest, nil, fmt.Errorf("Invalid request body: %s", models.AsJSONString(params.Body))
	}

	// Attempt to retrieve the record that we are trying to patch from the DB in order to do a merge
	fetched, err := tenantDB.GetMonitoredObject(tenantID, params.MonObjID)
	if err != nil {
		if strings.Contains(err.Error(), datastore.NotFoundStr) {
			return startTime, http.StatusNotFound, nil, err
		}

		return startTime, http.StatusInternalServerError, nil, fmt.Errorf("Unable to fetch %s: %s", tenmod.TenantMonitoredObjectStr, err.Error())
	}

	// Merge the attributes passed in with the patch request to the record fetched from the datastore
	var patched *tenmod.MonitoredObject
	if err := models.MergeObjWithMap(fetched, patchRequestBytes); err != nil {
		return startTime, http.StatusInternalServerError, nil, fmt.Errorf("Unable to patch %s with id %s: %s", tenmod.TenantMonitoredObjectStr, params.MonObjID, err.Error())
	}
	patched = fetched

	// Finally update the record in the datastore with the merged map and fetched tenant
	result, err := tenantDB.UpdateMonitoredObject(patched)
	if err != nil {
		if strings.Contains(err.Error(), datastore.ConflictErrorStr) {
			return startTime, http.StatusConflict, nil, err
		}

		return startTime, http.StatusInternalServerError, nil, err
	}

	if changeNotificationEnabled {
		NotifyMonitoredObjectUpdated(tenantID, fetched)
	}

	converted := swagmodels.MonitoredObjectResponse{}
	err = convertToJsonapiObject(result, &converted)
	if err != nil {
		return startTime, http.StatusInternalServerError, nil, fmt.Errorf("Unable to convert %s data to jsonapi return format: %s", tenmod.TenantMonitoredObjectStr, err.Error())
	}

	reportAPICompletionState(startTime, http.StatusOK, mon.UpdateThrPrfStr, mon.APICompleted, mon.TenantAPICompleted)
	logger.Log.Infof("Updated %s %s", tenmod.TenantMonitoredObjectStr, models.AsJSONString(converted))
	return startTime, http.StatusOK, &converted, nil
}

func doDeleteMonitoredObjectV2(allowedRoles []string, tenantDB datastore.TenantServiceDatastore, params tenant_provisioning_service_v2.DeleteMonitoredObjectV2Params) (time.Time, int, *swagmodels.MonitoredObjectResponse, error) {
	tenantID := params.HTTPRequest.Header.Get(XFwdTenantId)
	isAuthorized, startTime := authorizeRequest(fmt.Sprintf("Deleting %s %s for %s %s", tenmod.TenantMonitoredObjectStr, params.MonObjID, admmod.TenantStr, tenantID), params.HTTPRequest, allowedRoles, mon.APIRecieved, mon.AdminAPIRecieved)

	if !isAuthorized {
		return startTime, http.StatusForbidden, nil, fmt.Errorf("Delete %s operation not authorized for role: %s", tenmod.TenantMonitoredObjectStr, params.HTTPRequest.Header.Get(XFwdUserRoles))
	}

	// Issue request to DAO Layer
	result, err := tenantDB.DeleteMonitoredObject(tenantID, params.MonObjID)
	if err != nil {
		if strings.Contains(err.Error(), datastore.NotFoundStr) {
			return startTime, http.StatusNotFound, nil, err
		}

		return startTime, http.StatusInternalServerError, nil, fmt.Errorf("Unable to delete %s: %s", tenmod.TenantMonitoredObjectStr, err.Error())
	}

	if changeNotificationEnabled {
		NotifyMonitoredObjectDeleted(tenantID, result)
	}

	converted := swagmodels.MonitoredObjectResponse{}
	err = convertToJsonapiObject(result, &converted)
	if err != nil {
		return startTime, http.StatusInternalServerError, nil, fmt.Errorf("Unable to convert %s data to jsonapi return format: %s", tenmod.TenantMonitoredObjectStr, err.Error())
	}

	reportAPICompletionState(startTime, http.StatusOK, mon.DeleteThrPrfStr, mon.APICompleted, mon.TenantAPICompleted)
	logger.Log.Infof("Deleted %s %s", tenmod.TenantMonitoredObjectStr, models.AsJSONString(converted))
	return startTime, http.StatusOK, &converted, nil
}

func doGetAllMonitoredObjectsV2(allowedRoles []string, tenantDB datastore.TenantServiceDatastore, params tenant_provisioning_service_v2.GetAllMonitoredObjectsV2Params) (time.Time, int, *swagmodels.MonitoredObjectListResponse, error) {
	tenantID := params.HTTPRequest.Header.Get(XFwdTenantId)
	isAuthorized, startTime := authorizeRequest(fmt.Sprintf("Fetching all %s: %s", tenmod.TenantMonitoredObjectStr, tenantID), params.HTTPRequest, allowedRoles, mon.APIRecieved, mon.TenantAPIRecieved)

	if !isAuthorized {
		err := fmt.Errorf("Get all %s operation not authorized for role: %s", tenmod.TenantMonitoredObjectStr, params.HTTPRequest.Header.Get(XFwdUserRoles))
		return startTime, http.StatusForbidden, nil, err
	}

	// Issue request to DAO Layer
	startKey := ""
	limit := int64(0)
	if params.StartKey != nil {
		startKey = *params.StartKey
	}
	if params.Limit != nil {
		limit = *params.Limit
	}

	result, paginationOffsets, err := tenantDB.GetAllMonitoredObjectsByPage(tenantID, startKey, limit)
	if err != nil {
		if checkForNotFound(err.Error()) {
			return startTime, http.StatusNotFound, nil, err
		}

		errResp := fmt.Errorf("Unable to retrieve %s: %s", tenmod.TenantMonitoredObjectStr, err.Error())
		return startTime, http.StatusInternalServerError, nil, errResp
	}

	// Make sure the IDs are properly trimmed:
	for _, mo := range result {
		mo.ID = datastore.GetDataIDFromFullID(mo.ID)
	}

	converted := swagmodels.MonitoredObjectListResponse{}
	err = convertToJsonapiObject(result, &converted)
	if err != nil {
		errResp := fmt.Errorf("Unable to convert %s data to jsonapi return format: %s", tenmod.TenantMonitoredObjectStr, err.Error())
		return startTime, http.StatusInternalServerError, nil, errResp
	}

	// Add in the links section of the payload
	converted.Links = generateLinks(strings.Join([]string{params.HTTPRequest.URL.Scheme, params.HTTPRequest.URL.Host, params.HTTPRequest.URL.Path}, ""), paginationOffsets, limit)

	reportAPICompletionState(startTime, http.StatusOK, mon.GetAllMonObjStr, mon.APICompleted, mon.TenantAPICompleted)
	logger.Log.Infof("Retrieved %d %ss", len(converted.Data), tenmod.TenantMonitoredObjectStr)
	return startTime, http.StatusOK, &converted, nil
}

func doBulkInsertMonitoredObjectsV2(allowedRoles []string, tenantDB datastore.TenantServiceDatastore, params tenant_provisioning_service_v2.BulkInsertMonitoredObjectsV2Params) (time.Time, int, *swagmodels.BulkOperationResponseV2, error) {
	tenantID := params.HTTPRequest.Header.Get(XFwdTenantId)
	isAuthorized, startTime := authorizeRequest(fmt.Sprintf("Attempting Bulk Insert of %s for Tenant %s", tenmod.TenantMonitoredObjectStr, tenantID), params.HTTPRequest, allowedRoles, mon.APIRecieved, mon.TenantAPIRecieved)

	if !isAuthorized {
		err := fmt.Errorf("Bulk Insert %s operation not authorized for role: %s", tenmod.TenantMonitoredObjectStr, params.HTTPRequest.Header.Get(XFwdUserRoles))
		return startTime, http.StatusForbidden, nil, err
	}

	// Unmarshal the request
	data := []*tenmod.MonitoredObject{}
	for _, val := range params.Body.Data {
		addItem, err := convertMOBulkItemToDBMonitoredObject(val.Attributes, tenantID)
		if err != nil {
			return startTime, http.StatusBadRequest, nil, err
		}
		data = append(data, addItem)
	}

	if len(data) == 0 {
		return startTime, http.StatusBadRequest, nil, fmt.Errorf("No Monitored Objects in provided in the request")
	}

	// Issue request to DAO Layer
	result, err := tenantDB.BulkInsertMonitoredObjects(tenantID, data)
	if err != nil {
		return startTime, http.StatusInternalServerError, nil, fmt.Errorf("Unable to store %s: %s", tenmod.TenantMonitoredObjectStr, err.Error())
	}

	if changeNotificationEnabled {
		NotifyMonitoredObjectCreated(tenantID, data...)
	}

	converted, err := convertToBulkMOResponse(result)
	// err = json.Unmarshal(res, &converted)
	if err != nil {
		return startTime, http.StatusInternalServerError, nil, fmt.Errorf("Unable to format bulk insert %s: %s", tenmod.TenantMonitoredObjectStr, err.Error())
	}

	reportAPICompletionState(startTime, http.StatusOK, mon.BulkInsertMonObjStr, mon.APICompleted, mon.TenantAPICompleted)
	logger.Log.Infof("Bulk insertion of %ss complete", tenmod.TenantMonitoredObjectStr)
	return startTime, http.StatusOK, &converted, nil
}

func doBulkUpdateMonitoredObjectsV2(allowedRoles []string, tenantDB datastore.TenantServiceDatastore, params tenant_provisioning_service_v2.BulkUpdateMonitoredObjectsV2Params) (time.Time, int, *swagmodels.BulkOperationResponseV2, error) {
	tenantID := params.HTTPRequest.Header.Get(XFwdTenantId)
	isAuthorized, startTime := authorizeRequest(fmt.Sprintf("Attempting Bulk Update of %s for Tenant %s", tenmod.TenantMonitoredObjectStr, tenantID), params.HTTPRequest, allowedRoles, mon.APIRecieved, mon.TenantAPIRecieved)

	if !isAuthorized {
		err := fmt.Errorf("Bulk Update %s operation not authorized for role: %s", tenmod.TenantMonitoredObjectStr, params.HTTPRequest.Header.Get(XFwdUserRoles))
		return startTime, http.StatusForbidden, nil, err
	}

	// Unmarshal the request
	data := []*tenmod.MonitoredObject{}
	for _, val := range params.Body.Data {
		addItem, err := convertMOBulkItemToDBMonitoredObject(val.Attributes, tenantID)
		if err != nil {
			return startTime, http.StatusBadRequest, nil, err
		}
		data = append(data, addItem)
	}

	if len(data) == 0 {
		return startTime, http.StatusBadRequest, nil, fmt.Errorf("No Monitored Objects in provided in the request")
	}

	// Issue request to DAO Layer
	result, err := tenantDB.BulkUpdateMonitoredObjects(tenantID, data)
	if err != nil {
		return startTime, http.StatusInternalServerError, nil, fmt.Errorf("Unable to store %s: %s", tenmod.TenantMonitoredObjectStr, err.Error())
	}

	if changeNotificationEnabled {
		NotifyMonitoredObjectUpdated(tenantID, data...)
	}

	converted, err := convertToBulkMOResponse(result)
	// err = json.Unmarshal(res, &converted)
	if err != nil {
		return startTime, http.StatusInternalServerError, nil, fmt.Errorf("Unable to format bulk update %s: %s", tenmod.TenantMonitoredObjectStr, err.Error())
	}

	reportAPICompletionState(startTime, http.StatusOK, mon.BulkUpdateMonObjStr, mon.APICompleted, mon.TenantAPICompleted)
	logger.Log.Infof("Bulk insertion of %ss complete", tenmod.TenantMonitoredObjectStr)
	return startTime, http.StatusOK, &converted, nil
}

func convertMOBulkItemToDBMonitoredObject(item interface{}, tenantID string) (*tenmod.MonitoredObject, error) {
	if item == nil {
		return nil, nil
	}

	attrBytes, err := json.Marshal(item)
	if err != nil {
		return nil, fmt.Errorf("Unable to marshal Monitored Object %s", models.AsJSONString(item))
	}

	result := tenmod.MonitoredObject{}
	if err = json.Unmarshal(attrBytes, &result); err != nil {
		return nil, fmt.Errorf("Unable to unmarshal Monitored Object %s", models.AsJSONString(item))
	}

	result.TenantID = tenantID
	result.ID = result.MonitoredObjectID

	return &result, nil
}

func convertToBulkMOResponse(results []*common.BulkOperationResult) (swagmodels.BulkOperationResponseV2, error) {
	responseType := "bulkOperationResponses"

	response := swagmodels.BulkOperationResponseV2{}
	if results == nil || len(results) == 0 {
		return response, nil
	}

	items := []*swagmodels.BulkOperationResponseV2DataItems0{}
	for _, result := range results {
		attrBytes, err := json.Marshal(result)
		if err != nil {
			return response, err
		}

		addItem := swagmodels.BulkOperationResponseV2DataItems0Attributes{}
		if err = json.Unmarshal(attrBytes, &addItem); err != nil {
			return response, err
		}

		items = append(items, &swagmodels.BulkOperationResponseV2DataItems0{
			Attributes: &addItem,
			Type:       &responseType,
			ID:         &addItem.ID,
		})
	}

	response.Data = items
	return response, nil
}
