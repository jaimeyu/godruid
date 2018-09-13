package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/accedian/adh-gather/models"
	"github.com/accedian/adh-gather/scheduler"
	"github.com/accedian/adh-gather/swagmodels"
	"github.com/manyminds/api2go/jsonapi"

	"github.com/accedian/adh-gather/datastore"
	"github.com/accedian/adh-gather/logger"
	admmod "github.com/accedian/adh-gather/models/admin"
	metmod "github.com/accedian/adh-gather/models/metrics"
	tenmod "github.com/accedian/adh-gather/models/tenant"
	mon "github.com/accedian/adh-gather/monitoring"
	"github.com/accedian/adh-gather/restapi/operations/tenant_provisioning_service_v2"
	"github.com/go-openapi/runtime/middleware"
)

// HandleCreateReportScheduleConfigV2 - create a new ReportSchedule Configuration
func HandleCreateReportScheduleConfigV2(allowedRoles []string, tenantDB datastore.TenantServiceDatastore) func(params tenant_provisioning_service_v2.CreateReportScheduleConfigV2Params) middleware.Responder {
	return func(params tenant_provisioning_service_v2.CreateReportScheduleConfigV2Params) middleware.Responder {

		// Do the work
		startTime, responseCode, response, err := doCreateReportScheduleConfigV2(allowedRoles, tenantDB, params)

		// Success Response
		if responseCode == http.StatusCreated {
			return tenant_provisioning_service_v2.NewCreateReportScheduleConfigV2Created().WithPayload(response)
		}

		// Error Responses
		errorMessage := reportAPIError(err.Error(), startTime, responseCode, mon.CreateReportScheduleConfigStr, mon.APICompleted, mon.TenantAPICompleted)
		switch responseCode {
		case http.StatusForbidden:
			return tenant_provisioning_service_v2.NewCreateReportScheduleConfigV2Forbidden().WithPayload(errorMessage)
		case http.StatusBadRequest:
			return tenant_provisioning_service_v2.NewCreateReportScheduleConfigV2BadRequest().WithPayload(errorMessage)
		case http.StatusConflict:
			return tenant_provisioning_service_v2.NewCreateReportScheduleConfigV2Conflict().WithPayload(errorMessage)
		default:
			return tenant_provisioning_service_v2.NewCreateReportScheduleConfigV2InternalServerError().WithPayload(errorMessage)
		}
	}
}

// HandleGetReportScheduleConfigV2 - retrieve a ReportSchedule Config by the config ID.
func HandleGetReportScheduleConfigV2(allowedRoles []string, tenantDB datastore.TenantServiceDatastore) func(params tenant_provisioning_service_v2.GetReportScheduleConfigV2Params) middleware.Responder {
	return func(params tenant_provisioning_service_v2.GetReportScheduleConfigV2Params) middleware.Responder {

		// Do the work
		startTime, responseCode, response, err := doGetReportScheduleConfigV2(allowedRoles, tenantDB, params)

		// Success Response
		if responseCode == http.StatusOK {
			return tenant_provisioning_service_v2.NewGetReportScheduleConfigV2OK().WithPayload(response)
		}

		// Error Responses
		errorMessage := reportAPIError(err.Error(), startTime, responseCode, mon.GetReportScheduleConfigStr, mon.APICompleted, mon.TenantAPICompleted)
		switch responseCode {
		case http.StatusForbidden:
			return tenant_provisioning_service_v2.NewGetReportScheduleConfigV2Forbidden().WithPayload(errorMessage)
		case http.StatusNotFound:
			return tenant_provisioning_service_v2.NewGetReportScheduleConfigV2NotFound().WithPayload(errorMessage)
		default:
			return tenant_provisioning_service_v2.NewGetReportScheduleConfigV2InternalServerError().WithPayload(errorMessage)
		}

	}
}

// HandleUpdateReportScheduleConfigV2 - update a ReportScheduleConfig record
func HandleUpdateReportScheduleConfigV2(allowedRoles []string, tenantDB datastore.TenantServiceDatastore) func(params tenant_provisioning_service_v2.UpdateReportScheduleConfigV2Params) middleware.Responder {
	return func(params tenant_provisioning_service_v2.UpdateReportScheduleConfigV2Params) middleware.Responder {

		// Do the work
		startTime, responseCode, response, err := doUpdateReportScheduleConfigV2(allowedRoles, tenantDB, params)

		// Success Response
		if responseCode == http.StatusOK {
			return tenant_provisioning_service_v2.NewUpdateReportScheduleConfigV2OK().WithPayload(response)
		}

		// Error Responses
		errorMessage := reportAPIError(err.Error(), startTime, responseCode, mon.UpdateReportScheduleConfigStr, mon.APICompleted, mon.TenantAPICompleted)
		switch responseCode {
		case http.StatusForbidden:
			return tenant_provisioning_service_v2.NewUpdateReportScheduleConfigV2Forbidden().WithPayload(errorMessage)
		case http.StatusBadRequest:
			return tenant_provisioning_service_v2.NewUpdateReportScheduleConfigV2BadRequest().WithPayload(errorMessage)
		case http.StatusNotFound:
			return tenant_provisioning_service_v2.NewUpdateReportScheduleConfigV2NotFound().WithPayload(errorMessage)
		case http.StatusConflict:
			return tenant_provisioning_service_v2.NewUpdateReportScheduleConfigV2Conflict().WithPayload(errorMessage)
		default:
			return tenant_provisioning_service_v2.NewUpdateReportScheduleConfigV2InternalServerError().WithPayload(errorMessage)
		}
	}
}

// HandleDeleteReportScheduleConfigV2 - delete a tenant by the tenant ID.
func HandleDeleteReportScheduleConfigV2(allowedRoles []string, tenantDB datastore.TenantServiceDatastore) func(params tenant_provisioning_service_v2.DeleteReportScheduleConfigV2Params) middleware.Responder {
	return func(params tenant_provisioning_service_v2.DeleteReportScheduleConfigV2Params) middleware.Responder {

		// Do the work
		startTime, responseCode, response, err := doDeleteReportScheduleConfigV2(allowedRoles, tenantDB, params)

		// Success Response
		if responseCode == http.StatusOK {
			return tenant_provisioning_service_v2.NewDeleteReportScheduleConfigV2OK().WithPayload(response)
		}

		// Error Responses
		errorMessage := reportAPIError(err.Error(), startTime, responseCode, mon.DeleteReportScheduleConfigStr, mon.APICompleted, mon.TenantAPICompleted)
		switch responseCode {
		case http.StatusForbidden:
			return tenant_provisioning_service_v2.NewDeleteReportScheduleConfigV2Forbidden().WithPayload(errorMessage)
		case http.StatusNotFound:
			return tenant_provisioning_service_v2.NewDeleteReportScheduleConfigV2NotFound().WithPayload(errorMessage)
		default:
			return tenant_provisioning_service_v2.NewDeleteReportScheduleConfigV2InternalServerError().WithPayload(errorMessage)
		}

	}
}

// HandleGetAllReportScheduleConfigsV2 - retrieve all tenants
func HandleGetAllReportScheduleConfigsV2(allowedRoles []string, tenantDB datastore.TenantServiceDatastore) func(params tenant_provisioning_service_v2.GetAllReportScheduleConfigsV2Params) middleware.Responder {
	return func(params tenant_provisioning_service_v2.GetAllReportScheduleConfigsV2Params) middleware.Responder {

		// Do the work
		startTime, responseCode, response, err := doGetAllReportScheduleConfigsV2(allowedRoles, tenantDB, params)

		// Success Response
		if responseCode == http.StatusOK {
			return tenant_provisioning_service_v2.NewGetAllReportScheduleConfigsV2OK().WithPayload(response)
		}

		// Error Responses
		errorMessage := reportAPIError(err.Error(), startTime, responseCode, mon.GetAllReportScheduleConfigStr, mon.APICompleted, mon.TenantAPICompleted)
		switch responseCode {
		case http.StatusForbidden:
			return tenant_provisioning_service_v2.NewGetAllReportScheduleConfigsV2Forbidden().WithPayload(errorMessage)
		case http.StatusNotFound:
			return tenant_provisioning_service_v2.NewGetAllReportScheduleConfigsV2NotFound().WithPayload(errorMessage)
		default:
			return tenant_provisioning_service_v2.NewGetAllReportScheduleConfigsV2InternalServerError().WithPayload(errorMessage)
		}

	}
}

func doCreateReportScheduleConfigV2(allowedRoles []string, tenantDB datastore.TenantServiceDatastore, params tenant_provisioning_service_v2.CreateReportScheduleConfigV2Params) (time.Time, int, *swagmodels.ReportScheduleConfigResponse, error) {
	tenantID := params.HTTPRequest.Header.Get(XFwdTenantId)
	isAuthorized, startTime := authorizeRequest(fmt.Sprintf("Creating %s for tenant %s", tenmod.TenantReportScheduleConfigStr, tenantID), params.HTTPRequest, allowedRoles, mon.APIRecieved, mon.TenantAPIRecieved)

	if !isAuthorized {
		return startTime, http.StatusForbidden, nil, fmt.Errorf("Create %s operation not authorized for role: %s", tenmod.TenantReportScheduleConfigStr, params.HTTPRequest.Header.Get(XFwdUserRoles))
	}

	// Unmarshal the request
	requestBytes, err := json.Marshal(params.Body)
	if err != nil {
		return startTime, http.StatusBadRequest, nil, err
	}

	data := metmod.ReportScheduleConfig{}
	if err = jsonapi.Unmarshal(requestBytes, &data); err != nil {
		return startTime, http.StatusBadRequest, nil, err
	}

	data.TenantID = tenantID

	// Ensure that the passed in data adheres to the model requirements
	if err = data.Validate(false); err != nil {
		return startTime, http.StatusBadRequest, nil, err
	}

	// Issue request to DAO Layer to Create Record
	result, err := tenantDB.CreateReportScheduleConfig(&data)
	if err != nil {
		if strings.Contains(err.Error(), datastore.ConflictErrorStr) {
			return startTime, http.StatusConflict, nil, err
		}
		return startTime, http.StatusInternalServerError, nil, fmt.Errorf("Unable to store %s: %s", tenmod.TenantReportScheduleConfigStr, err.Error())
	}

	// Tell the scheduler to go and update based on the updated database
	if err = scheduler.RebuildCronJobs(); err != nil {
		return startTime, http.StatusInternalServerError, nil, fmt.Errorf("Unable to start scheduled job%s: %s", metmod.ReportScheduleConfigStr, err.Error())
	}

	converted := swagmodels.ReportScheduleConfigResponse{}
	if err = convertToJsonapiObject(result, &converted); err != nil {
		return startTime, http.StatusInternalServerError, nil, fmt.Errorf("Unable to convert %s data to jsonapi return format: %s", tenmod.TenantReportScheduleConfigStr, err.Error())
	}

	reportAPICompletionState(startTime, http.StatusOK, mon.CreateReportScheduleConfigStr, mon.APICompleted, mon.TenantAPICompleted)
	logger.Log.Infof("Created %s %s", tenmod.TenantReportScheduleConfigStr, models.AsJSONString(converted))
	return startTime, http.StatusCreated, &converted, nil
}

func doGetReportScheduleConfigV2(allowedRoles []string, tenantDB datastore.TenantServiceDatastore, params tenant_provisioning_service_v2.GetReportScheduleConfigV2Params) (time.Time, int, *swagmodels.ReportScheduleConfigResponse, error) {
	tenantID := params.HTTPRequest.Header.Get(XFwdTenantId)
	isAuthorized, startTime := authorizeRequest(fmt.Sprintf("Fetching %s %s for %s %s", tenmod.TenantReportScheduleConfigStr, params.ConfigID, admmod.TenantStr, tenantID), params.HTTPRequest, allowedRoles, mon.APIRecieved, mon.AdminAPIRecieved)

	if !isAuthorized {
		return startTime, http.StatusForbidden, nil, fmt.Errorf("Fetch %s operation not authorized for role: %s", tenmod.TenantReportScheduleConfigStr, params.HTTPRequest.Header.Get(XFwdUserRoles))
	}

	// Issue request to DAO Layer
	result, err := tenantDB.GetReportScheduleConfig(tenantID, params.ConfigID)
	if err != nil {
		if strings.Contains(err.Error(), datastore.NotFoundStr) {
			return startTime, http.StatusNotFound, nil, err
		}

		return startTime, http.StatusInternalServerError, nil, fmt.Errorf("Unable to retrieve %s: %s", tenmod.TenantReportScheduleConfigStr, err.Error())
	}

	converted := swagmodels.ReportScheduleConfigResponse{}
	err = convertToJsonapiObject(result, &converted)
	if err != nil {
		return startTime, http.StatusInternalServerError, nil, fmt.Errorf("Unable to convert %s data to jsonapi return format: %s", tenmod.TenantReportScheduleConfigStr, err.Error())
	}

	reportAPICompletionState(startTime, http.StatusOK, mon.GetReportScheduleConfigStr, mon.APICompleted, mon.TenantAPICompleted)
	logger.Log.Infof("Retrieved %s %s", tenmod.TenantReportScheduleConfigStr, models.AsJSONString(converted))
	return startTime, http.StatusOK, &converted, nil
}

func doUpdateReportScheduleConfigV2(allowedRoles []string, tenantDB datastore.TenantServiceDatastore, params tenant_provisioning_service_v2.UpdateReportScheduleConfigV2Params) (time.Time, int, *swagmodels.ReportScheduleConfigResponse, error) {
	tenantID := params.HTTPRequest.Header.Get(XFwdTenantId)
	isAuthorized, startTime := authorizeRequest(fmt.Sprintf("Updating %s %s for %s %s", tenmod.TenantReportScheduleConfigStr, params.ConfigID, admmod.TenantStr, tenantID), params.HTTPRequest, allowedRoles, mon.APIRecieved, mon.AdminAPIRecieved)

	if !isAuthorized {
		return startTime, http.StatusForbidden, nil, fmt.Errorf("Update %s operation not authorized for role: %s", tenmod.TenantReportScheduleConfigStr, params.HTTPRequest.Header.Get(XFwdUserRoles))
	}

	// Retrieve tne request bytes from the payload in order to convert it to a map
	patchRequestBytes, err := json.Marshal(params.Body)
	if err != nil || params.Body == nil {
		return startTime, http.StatusBadRequest, nil, fmt.Errorf("Invalid request body: %s", models.AsJSONString(params.Body))
	}

	// Attempt to retrieve the record that we are trying to patch from the DB in order to do a merge
	fetched, err := tenantDB.GetReportScheduleConfig(tenantID, params.ConfigID)
	if err != nil {
		if strings.Contains(err.Error(), datastore.NotFoundStr) {
			return startTime, http.StatusNotFound, nil, err
		}

		return startTime, http.StatusInternalServerError, nil, fmt.Errorf("Unable to fetch %s: %s", tenmod.TenantReportScheduleConfigStr, err.Error())
	}

	// Merge the attributes passed in with the patch request to the record fetched from the datastore
	var patched *metmod.ReportScheduleConfig
	if err := models.MergeObjWithMap(fetched, patchRequestBytes); err != nil {
		return startTime, http.StatusInternalServerError, nil, fmt.Errorf("Unable to patch %s with id %s: %s", tenmod.TenantReportScheduleConfigStr, params.ConfigID, err.Error())
	}
	patched = fetched
	patched.TenantID = tenantID

	// Before updating, make sure to handle any relationship data:
	if params.Body.Data.Relationships != nil {
		tp := params.Body.Data.Relationships.ThresholdProfile
		if tp != nil {
			patched.ThresholdProfile = params.Body.Data.Relationships.ThresholdProfile.Data.ID
		}
	}

	// Ensure that the passed in data adheres to the model requirements
	err = patched.Validate(true)
	if err != nil {
		return startTime, http.StatusBadRequest, nil, err
	}

	// Finally update the record in the datastore with the merged map and fetched tenant
	result, err := tenantDB.UpdateReportScheduleConfig(patched)
	if err != nil {
		if strings.Contains(err.Error(), datastore.ConflictErrorStr) {
			return startTime, http.StatusConflict, nil, err
		}

		return startTime, http.StatusInternalServerError, nil, err
	}

	// Tell the scheduler to go and update based on the updated database
	if err = scheduler.RebuildCronJobs(); err != nil {
		return startTime, http.StatusInternalServerError, nil, fmt.Errorf("Unable to start scheduled job%s: %s", metmod.ReportScheduleConfigStr, err.Error())
	}

	converted := swagmodels.ReportScheduleConfigResponse{}
	err = convertToJsonapiObject(result, &converted)
	if err != nil {
		return startTime, http.StatusInternalServerError, nil, fmt.Errorf("Unable to convert %s data to jsonapi return format: %s", tenmod.TenantReportScheduleConfigStr, err.Error())
	}

	reportAPICompletionState(startTime, http.StatusOK, mon.UpdateReportScheduleConfigStr, mon.APICompleted, mon.TenantAPICompleted)
	logger.Log.Infof("Updated %s %s", tenmod.TenantReportScheduleConfigStr, models.AsJSONString(converted))
	return startTime, http.StatusOK, &converted, nil
}

func doDeleteReportScheduleConfigV2(allowedRoles []string, tenantDB datastore.TenantServiceDatastore, params tenant_provisioning_service_v2.DeleteReportScheduleConfigV2Params) (time.Time, int, *swagmodels.ReportScheduleConfigResponse, error) {
	tenantID := params.HTTPRequest.Header.Get(XFwdTenantId)
	isAuthorized, startTime := authorizeRequest(fmt.Sprintf("Deleting %s %s for %s %s", tenmod.TenantReportScheduleConfigStr, params.ConfigID, admmod.TenantStr, tenantID), params.HTTPRequest, allowedRoles, mon.APIRecieved, mon.AdminAPIRecieved)

	if !isAuthorized {
		return startTime, http.StatusForbidden, nil, fmt.Errorf("Delete %s operation not authorized for role: %s", tenmod.TenantReportScheduleConfigStr, params.HTTPRequest.Header.Get(XFwdUserRoles))
	}

	// Issue request to DAO Layer
	result, err := tenantDB.DeleteReportScheduleConfig(tenantID, params.ConfigID)
	if err != nil {
		if strings.Contains(err.Error(), datastore.NotFoundStr) {
			return startTime, http.StatusNotFound, nil, err
		}

		return startTime, http.StatusInternalServerError, nil, fmt.Errorf("Unable to delete %s: %s", tenmod.TenantReportScheduleConfigStr, err.Error())
	}

	// Tell the scheduler to go and update based on the updated database
	if err = scheduler.RebuildCronJobs(); err != nil {
		return startTime, http.StatusInternalServerError, nil, fmt.Errorf("Unable to start scheduled job%s: %s", metmod.ReportScheduleConfigStr, err.Error())
	}

	converted := swagmodels.ReportScheduleConfigResponse{}
	err = convertToJsonapiObject(result, &converted)
	if err != nil {
		return startTime, http.StatusInternalServerError, nil, fmt.Errorf("Unable to convert %s data to jsonapi return format: %s", tenmod.TenantReportScheduleConfigStr, err.Error())
	}

	reportAPICompletionState(startTime, http.StatusOK, mon.DeleteReportScheduleConfigStr, mon.APICompleted, mon.TenantAPICompleted)
	logger.Log.Infof("Deleted %s %s", tenmod.TenantReportScheduleConfigStr, models.AsJSONString(converted))
	return startTime, http.StatusOK, &converted, nil
}

func doGetAllReportScheduleConfigsV2(allowedRoles []string, tenantDB datastore.TenantServiceDatastore, params tenant_provisioning_service_v2.GetAllReportScheduleConfigsV2Params) (time.Time, int, *swagmodels.ReportScheduleConfigListResponse, error) {
	tenantID := params.HTTPRequest.Header.Get(XFwdTenantId)
	isAuthorized, startTime := authorizeRequest(fmt.Sprintf("Fetching %s list for %s %s", tenmod.TenantReportScheduleConfigStr, admmod.TenantStr, tenantID), params.HTTPRequest, allowedRoles, mon.APIRecieved, mon.AdminAPIRecieved)

	if !isAuthorized {
		return startTime, http.StatusForbidden, nil, fmt.Errorf("Fetch %s operation not authorized for role: %s", tenmod.TenantReportScheduleConfigStr, params.HTTPRequest.Header.Get(XFwdUserRoles))
	}

	// Issue request to DAO Layer
	result, err := tenantDB.GetAllReportScheduleConfigs(tenantID)
	if err != nil {
		if strings.Contains(err.Error(), datastore.NotFoundStr) {
			return startTime, http.StatusNotFound, nil, err
		}

		return startTime, http.StatusInternalServerError, nil, fmt.Errorf("Unable to retrieve %s list: %s", tenmod.TenantReportScheduleConfigStr, err.Error())
	}

	if len(result) == 0 {
		return startTime, http.StatusNotFound, nil, fmt.Errorf(datastore.NotFoundStr)
	}

	converted := swagmodels.ReportScheduleConfigListResponse{}
	err = convertToJsonapiObject(result, &converted)
	if err != nil {
		return startTime, http.StatusInternalServerError, nil, fmt.Errorf("Unable to convert %s data to jsonapi return format: %s", tenmod.TenantReportScheduleConfigStr, err.Error())
	}

	reportAPICompletionState(startTime, http.StatusOK, mon.GetAllReportScheduleConfigStr, mon.APICompleted, mon.TenantAPICompleted)
	logger.Log.Infof("Retrieved %d %ss", len(result), tenmod.TenantReportScheduleConfigStr)
	return startTime, http.StatusOK, &converted, nil
}
