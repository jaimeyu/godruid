package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
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

// HandleGetMonitoredObjectListV2 - retrieve all monitored object ids for a Tenant based on the specified search criteria
func HandleGetFilteredMonitoredObjectListV2(allowedRoles []string, tenantDB datastore.TenantServiceDatastore) func(params tenant_provisioning_service_v2.GetFilteredMonitoredObjectListV2Params) middleware.Responder {
	return func(params tenant_provisioning_service_v2.GetFilteredMonitoredObjectListV2Params) middleware.Responder {
		// Do the work
		startTime, responseCode, response, err := doGetFilteredMonitoredObjectListV2(allowedRoles, tenantDB, params)

		// Success Response
		if responseCode == http.StatusOK {
			return tenant_provisioning_service_v2.NewGetFilteredMonitoredObjectListV2OK().WithPayload(response)
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

func HandleBulkInsertMonitoredObjectsMetaV2(allowedRoles []string, tenantDB datastore.TenantServiceDatastore) func(params tenant_provisioning_service_v2.BulkInsertMonitoredObjectsMetaV2Params) middleware.Responder {
	return func(params tenant_provisioning_service_v2.BulkInsertMonitoredObjectsMetaV2Params) middleware.Responder {
		// Do the work
		startTime, responseCode, response, err := doBulkInsertMonitoredObjectsMetaV2(allowedRoles, tenantDB, params)

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

	err = data.Validate(false)
	if err != nil {
		return startTime, http.StatusBadRequest, nil, fmt.Errorf("Validation failed due to %s", err.Error())
	}

	// Issue request to DAO Layer to Create Record
	result, err := tenantDB.CreateMonitoredObject(&data)
	if err != nil {
		if strings.Contains(err.Error(), datastore.ConflictErrorStr) {
			return startTime, http.StatusConflict, nil, err
		}
		return startTime, http.StatusInternalServerError, nil, fmt.Errorf("Unable to store %s: %s", tenmod.TenantMonitoredObjectStr, err.Error())
	}

	// Build up the monitored object indices
	err = tenantDB.UpdateMonitoredObjectMetadataViews(tenantID, result.Meta)
	if err != nil {
		return startTime, http.StatusInternalServerError, nil, fmt.Errorf("Unable to update metadata views %s %s", tenmod.TenantMonitoredObjectStr, err.Error())
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
	patched := &tenmod.MonitoredObject{}
	logger.Log.Debugf("THE PATCH REQ is: %s", string(patchRequestBytes))
	if err := models.MergeObjWithMap(patched, fetched, patchRequestBytes); err != nil {
		return startTime, http.StatusInternalServerError, nil, fmt.Errorf("Unable to patch %s with id %s: %s", tenmod.TenantMonitoredObjectStr, params.MonObjID, err.Error())
	}
	patched.TenantID = tenantID

	err = patched.Validate(true)
	if err != nil {
		return startTime, http.StatusInternalServerError, nil, fmt.Errorf("Unable to patch %s with id %s: %s", tenmod.TenantMonitoredObjectStr, params.MonObjID, err.Error())
	}

	// Finally update the record in the datastore with the merged map and fetched tenant
	result, err := tenantDB.UpdateMonitoredObject(patched)
	if err != nil {
		if strings.Contains(err.Error(), datastore.ConflictErrorStr) {
			return startTime, http.StatusConflict, nil, err
		}

		return startTime, http.StatusInternalServerError, nil, err
	}

	// Build up the monitored object indices
	err = tenantDB.UpdateMonitoredObjectMetadataViews(tenantID, result.Meta)
	if err != nil {
		return startTime, http.StatusInternalServerError, nil, fmt.Errorf("Unable to update metadata views %s %s", tenmod.TenantMonitoredObjectStr, err.Error())
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

	//Clean up the views now that a monitored object is deleted.
	tenantDB.UpdateMonitoredObjectMetadataViews(tenantID, result.Meta)

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

func doGetFilteredMonitoredObjectListV2(allowedRoles []string, tenantDB datastore.TenantServiceDatastore, params tenant_provisioning_service_v2.GetFilteredMonitoredObjectListV2Params) (time.Time, int, *swagmodels.MonitoredObjectFilteredListResponse, error) {
	tenantID := params.HTTPRequest.Header.Get(XFwdTenantId)
	isAuthorized, startTime := authorizeRequest(fmt.Sprintf("Fetching all %s: %s", tenmod.TenantMonitoredObjectKeysStr, tenantID), params.HTTPRequest, allowedRoles, mon.APIRecieved, mon.TenantAPIRecieved)

	if !isAuthorized {
		err := fmt.Errorf("Get all %s operation not authorized for role: %s", tenmod.TenantMonitoredObjectKeysStr, params.HTTPRequest.Header.Get(XFwdUserRoles))
		return startTime, http.StatusForbidden, nil, err
	}

	if tenantID == "" {
		return startTime, http.StatusBadRequest, nil, fmt.Errorf("Invalid request missing tenant ID")
	}

	meta := make(map[string][]string)

	requestBytes, err := json.Marshal(params.Body.Meta)
	if err != nil || params.Body == nil {
		return startTime, http.StatusBadRequest, nil, fmt.Errorf("Invalid request body: %s", models.AsJSONString(params.Body))
	}
	// Convert the request JSON into a map
	err = json.Unmarshal(requestBytes, &meta)
	if err != nil {
		errResp := fmt.Errorf("Unable to unmarshal metadata for %s request for tenant %s: %s", tenmod.TenantMonitoredObjectKeysStr, tenantID, err.Error())
		return startTime, http.StatusInternalServerError, nil, errResp
	}

	logger.Log.Debugf("Retrieving %s for tenant %s with meta filters %v", tenmod.TenantMonitoredObjectStr, tenantID, meta)

	resourceIdentifierList, err := tenantDB.GetFilteredMonitoredObjectList(tenantID, meta)

	if err != nil {
		if checkForNotFound(err.Error()) {
			return startTime, http.StatusNotFound, nil, err
		}

		errResp := fmt.Errorf("Unable to retrieve %s: %s", tenmod.TenantMonitoredObjectKeysStr, err.Error())
		return startTime, http.StatusInternalServerError, nil, errResp
	}

	responseMap := wrapJsonAPIObject(map[string]interface{}{"resourceIdentifiers": resourceIdentifierList}, "1", "filteredResourceIdentifierList")

	responseBytes, err := json.Marshal(&responseMap)
	if err != nil {
		errResp := fmt.Errorf("Unable to convert %s data to json return format: %s", tenmod.TenantMonitoredObjectKeysStr, err.Error())
		return startTime, http.StatusInternalServerError, nil, errResp
	}

	converted := swagmodels.MonitoredObjectFilteredListResponse{}
	err = json.Unmarshal(responseBytes, &converted)
	if err != nil {
		errResp := fmt.Errorf("Unable to convert %s data to json return format: %s", tenmod.TenantMonitoredObjectKeysStr, err.Error())
		return startTime, http.StatusInternalServerError, nil, errResp
	}

	reportAPICompletionState(startTime, http.StatusOK, mon.GetAllMonObjStr, mon.APICompleted, mon.TenantAPICompleted)
	logger.Log.Infof("Retrieved %d %ss", len(converted.Data.Attributes.ResourceIdentifiers), tenmod.TenantMonitoredObjectKeysStr)
	return startTime, http.StatusOK, &converted, nil
}

func doBulkInsertMonitoredObjectsV2(allowedRoles []string, tenantDB datastore.TenantServiceDatastore, params tenant_provisioning_service_v2.BulkInsertMonitoredObjectsV2Params) (time.Time, int, *swagmodels.BulkOperationResponseV2, error) {
	tenantID := params.HTTPRequest.Header.Get(XFwdTenantId)
	isAuthorized, startTime := authorizeRequest(fmt.Sprintf("Attempting Bulk Insert of %s for Tenant %s", tenmod.TenantMonitoredObjectStr, tenantID), params.HTTPRequest, allowedRoles, mon.APIRecieved, mon.TenantAPIRecieved)

	if !isAuthorized {
		err := fmt.Errorf("Bulk Insert %s operation not authorized for role: %s", tenmod.TenantMonitoredObjectStr, params.HTTPRequest.Header.Get(XFwdUserRoles))
		return startTime, http.StatusForbidden, nil, err
	}

	metaKeys := make(map[string]string)

	// Unmarshal the request
	data := []*tenmod.MonitoredObject{}
	for _, val := range params.Body.Data {
		addItem, err := convertMOBulkItemToDBMonitoredObject(val.Attributes, tenantID)
		if err != nil {
			return startTime, http.StatusBadRequest, nil, err
		}

		err = addItem.Validate(false)
		if err != nil {
			return startTime, http.StatusBadRequest, nil, fmt.Errorf("Validation failure of %s because: %s", models.AsJSONString(addItem), err.Error())
		}

		// Track all distinct metadata items to be index processed after all are items are worked through
		for k := range addItem.Meta {
			metaKeys[k] = ""
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

	// Build up the monitored object indices
	err = tenantDB.UpdateMonitoredObjectMetadataViews(tenantID, metaKeys)
	if err != nil {
		return startTime, http.StatusInternalServerError, nil, fmt.Errorf("Unable to update metadata views %s %s", tenmod.TenantMonitoredObjectStr, err.Error())
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

	metaKeys := make(map[string]string)

	// Unmarshal the request
	data := []*tenmod.MonitoredObject{}
	for _, val := range params.Body.Data {
		addItem, err := convertMOBulkItemToDBMonitoredObject(val.Attributes, tenantID)
		if err != nil {
			return startTime, http.StatusBadRequest, nil, err
		}

		err = addItem.Validate(true)
		if err != nil {
			return startTime, http.StatusBadRequest, nil, fmt.Errorf("Validation failure of %s because: %s", models.AsJSONString(addItem), err.Error())
		}

		// Track all distinct metadata items to be index processed after all are items are worked through
		for k := range addItem.Meta {
			metaKeys[k] = ""
		}

		data = append(data, addItem)
	}

	if len(data) == 0 {
		return startTime, http.StatusBadRequest, nil, fmt.Errorf("No Monitored Objects provided in the request")
	}

	// Issue request to DAO Layer
	result, err := tenantDB.BulkUpdateMonitoredObjects(tenantID, data)
	if err != nil {
		return startTime, http.StatusInternalServerError, nil, fmt.Errorf("Unable to store %s: %s", tenmod.TenantMonitoredObjectStr, err.Error())
	}

	if changeNotificationEnabled {
		NotifyMonitoredObjectUpdated(tenantID, data...)
	}

	// Build up the monitored object indices
	err = tenantDB.UpdateMonitoredObjectMetadataViews(tenantID, metaKeys)
	if err != nil {
		return startTime, http.StatusInternalServerError, nil, fmt.Errorf("Unable to update metadata views %s %s", tenmod.TenantMonitoredObjectStr, err.Error())
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

func doBulkInsertMonitoredObjectsMetaV2(allowedRoles []string, tenantDB datastore.TenantServiceDatastore, params tenant_provisioning_service_v2.BulkInsertMonitoredObjectsMetaV2Params) (time.Time, int, *swagmodels.BulkOperationResponseV2, error) {
	tenantID := params.HTTPRequest.Header.Get(XFwdTenantId)
	isAuthorized, startTime := authorizeRequest(fmt.Sprintf("Attempting Bulk Insert of %s metadata for Tenant %s", tenmod.TenantMonitoredObjectStr, tenantID), params.HTTPRequest, allowedRoles, mon.APIRecieved, mon.TenantAPIRecieved)

	if !isAuthorized {
		err := fmt.Errorf("Bulk Insert %s metadata operation not authorized for role: %s", tenmod.TenantMonitoredObjectStr, params.HTTPRequest.Header.Get(XFwdUserRoles))
		return startTime, http.StatusForbidden, nil, err
	}

	incrementAPICounters(mon.APIRecieved, mon.TenantAPIRecieved)
	logger.Log.Infof("Inserting %s metadata in bulk for Tenant %s", tenmod.TenantMonitoredObjectStr, tenantID)

	requestBytes, err := json.Marshal(params.Body)
	if err != nil {
		return startTime, http.StatusBadRequest, nil, err
	}

	// Unmarshal the request
	data := tenmod.BulkMetadataEntries{}
	err = jsonapi.Unmarshal(requestBytes, &data)
	if err != nil {
		return startTime, http.StatusBadRequest, nil, err
	}

	response := make([]*common.BulkOperationResult, len(data.MetadataEntries))

	// Internal function responsible for managing error scenarios for individual result items
	itemError := func(position int, itemResponse *common.BulkOperationResult, reason int, itemErr string) {
		itemResponse.OK = false
		itemResponse.REASON = strconv.Itoa(reason)
		itemResponse.ERROR = itemErr

		logger.Log.Errorf(generateErrorMessage(reason, itemErr))

		response[position] = itemResponse
	}

	metaKeys := make(map[string]string)
	const idSep = "_"

	for i, item := range data.MetadataEntries {
		itemResponse := common.BulkOperationResult{
			ID: item.ObjectName,
		}
		// Issue request to DAO Layer
		existingMonitoredObjects, err := tenantDB.GetMonitoredObjectsByObjectName(item.ObjectName, tenantID)
		if err != nil {
			itemError(i, &itemResponse, http.StatusNotFound, fmt.Sprintf("Unable to retrieve %s %s", tenmod.TenantMonitoredObjectStr, err.Error()))
			continue
		}

		revlist := "" // Create a comma separated list of revisions since multiple monitored objects could be associated with this key
		var moErr error

		for i, existingMO := range existingMonitoredObjects {
			logger.Log.Infof("Patching metadata for %s with name %s and id %s", tenmod.TenantMonitoredObjectStr, existingMO.ObjectName, existingMO.ID)

			existingMO.Meta = item.Metadata

			splitID := datastore.GetDataIDFromFullID(existingMO.ID)

			existingMO.ID = splitID

			moErr = existingMO.Validate(true)
			if moErr != nil {
				itemError(i, &itemResponse, http.StatusInternalServerError, fmt.Sprintf("Data did not validate %s: %s", tenmod.TenantMonitoredObjectStr, err.Error()))
				break
			}
			// Issue request to DAO Layer
			updatedMonitoredObject, moErr := tenantDB.UpdateMonitoredObject(existingMO)
			if moErr != nil {
				itemError(i, &itemResponse, http.StatusInternalServerError, fmt.Sprintf("Unable to store %s: %s", tenmod.TenantMonitoredObjectStr, err.Error()))
				break
			}

			if i > 0 {
				revlist += ","
			}
			revlist += updatedMonitoredObject.REV

			// Track all distinct metadata items to be index processed after all are items are worked through
			for k := range item.Metadata {
				metaKeys[k] = ""
			}

			logger.Log.Debugf("Sending notification of update to monitored object %s", existingMO.ObjectName)
			NotifyMonitoredObjectUpdated(existingMO.TenantID, existingMO)
		}

		// If there was an error against a monitored object for a set of monitored objects that share the same object ID then continue through the loop.
		// The inner loop will have already reported the proble
		if moErr != nil {
			continue
		}

		itemResponse.OK = true
		itemResponse.REV = revlist
		response[i] = &itemResponse
	}

	// Build up the monitored object indices
	err = tenantDB.UpdateMonitoredObjectMetadataViews(tenantID, metaKeys)
	if err != nil {
		return startTime, http.StatusInternalServerError, nil, err
	}

	converted, err := convertToBulkMOResponse(response)
	if err != nil {
		return startTime, http.StatusInternalServerError, nil, err
	}

	reportAPICompletionState(startTime, http.StatusOK, mon.BulkUpsertMonObjMetaStr, mon.APICompleted, mon.TenantAPICompleted)
	logger.Log.Infof("Bulk insertion of %ss meta data complete", tenmod.TenantMonitoredObjectStr)

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
