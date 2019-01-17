package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/accedian/adh-gather/gather"
	"github.com/accedian/adh-gather/models"
	"github.com/accedian/adh-gather/swagmodels"
	"github.com/manyminds/api2go/jsonapi"

	"github.com/accedian/adh-gather/datastore"
	"github.com/accedian/adh-gather/logger"
	admmod "github.com/accedian/adh-gather/models/admin"
	tenmod "github.com/accedian/adh-gather/models/tenant"
	mon "github.com/accedian/adh-gather/monitoring"
	"github.com/accedian/adh-gather/restapi/operations/tenant_provisioning_service_v2"
	"github.com/go-openapi/runtime/middleware"
)

// HandleCreateMetricBaselineV2 - create a new MetricBaseline record
func HandleCreateMetricBaselineV2(allowedRoles []string, tenantDB datastore.TenantMetricBaselineDatastore) func(params tenant_provisioning_service_v2.CreateMetricBaselineV2Params) middleware.Responder {
	return func(params tenant_provisioning_service_v2.CreateMetricBaselineV2Params) middleware.Responder {

		// Do the work
		startTime, responseCode, response, err := doCreateMetricBaselineV2(allowedRoles, tenantDB, params)

		// Success Response
		if responseCode == http.StatusCreated {
			return tenant_provisioning_service_v2.NewCreateMetricBaselineV2Created().WithPayload(response)
		}

		// Error Responses
		errorMessage := reportAPIError(err.Error(), startTime, responseCode, mon.CreateTenantMetricBaselineStr, mon.APICompleted, mon.TenantAPICompleted)
		switch responseCode {
		case http.StatusForbidden:
			return tenant_provisioning_service_v2.NewCreateMetricBaselineV2Forbidden().WithPayload(errorMessage)
		case http.StatusBadRequest:
			return tenant_provisioning_service_v2.NewCreateMetricBaselineV2BadRequest().WithPayload(errorMessage)
		case http.StatusConflict:
			return tenant_provisioning_service_v2.NewCreateMetricBaselineV2Conflict().WithPayload(errorMessage)
		default:
			return tenant_provisioning_service_v2.NewCreateMetricBaselineV2InternalServerError().WithPayload(errorMessage)
		}
	}
}

// HandleGetMetricBaselineV2 - retrieve a MetricBaseline record by the ID
func HandleGetMetricBaselineV2(allowedRoles []string, tenantDB datastore.TenantMetricBaselineDatastore) func(params tenant_provisioning_service_v2.GetMetricBaselineV2Params) middleware.Responder {
	return func(params tenant_provisioning_service_v2.GetMetricBaselineV2Params) middleware.Responder {

		// Do the work
		startTime, responseCode, response, err := doGetMetricBaselineV2(allowedRoles, tenantDB, params)

		// Success Response
		if responseCode == http.StatusOK {
			return tenant_provisioning_service_v2.NewGetMetricBaselineV2OK().WithPayload(response)
		}

		// Error Responses
		errorMessage := reportAPIError(err.Error(), startTime, responseCode, mon.GetTenantMetricBaselineStr, mon.APICompleted, mon.TenantAPICompleted)
		switch responseCode {
		case http.StatusForbidden:
			return tenant_provisioning_service_v2.NewGetMetricBaselineV2Forbidden().WithPayload(errorMessage)
		case http.StatusNotFound:
			return tenant_provisioning_service_v2.NewGetMetricBaselineV2NotFound().WithPayload(errorMessage)
		default:
			return tenant_provisioning_service_v2.NewGetMetricBaselineV2InternalServerError().WithPayload(errorMessage)
		}

	}
}

// HandleUpdateMetricBaselineV2 - update a MetricBaseline record
func HandleUpdateMetricBaselineV2(allowedRoles []string, tenantDB datastore.TenantMetricBaselineDatastore) func(params tenant_provisioning_service_v2.UpdateMetricBaselineV2Params) middleware.Responder {
	return func(params tenant_provisioning_service_v2.UpdateMetricBaselineV2Params) middleware.Responder {

		// Do the work
		startTime, responseCode, response, err := doUpdateMetricBaselineV2(allowedRoles, tenantDB, params)

		// Success Response
		if responseCode == http.StatusOK {
			return tenant_provisioning_service_v2.NewUpdateMetricBaselineV2OK().WithPayload(response)
		}

		// Error Responses
		errorMessage := reportAPIError(err.Error(), startTime, responseCode, mon.UpdateTenantMetricBaselineStr, mon.APICompleted, mon.TenantAPICompleted)
		switch responseCode {
		case http.StatusForbidden:
			return tenant_provisioning_service_v2.NewUpdateMetricBaselineV2Forbidden().WithPayload(errorMessage)
		case http.StatusBadRequest:
			return tenant_provisioning_service_v2.NewUpdateMetricBaselineV2BadRequest().WithPayload(errorMessage)
		case http.StatusNotFound:
			return tenant_provisioning_service_v2.NewUpdateMetricBaselineV2NotFound().WithPayload(errorMessage)
		case http.StatusConflict:
			return tenant_provisioning_service_v2.NewUpdateMetricBaselineV2Conflict().WithPayload(errorMessage)
		default:
			return tenant_provisioning_service_v2.NewUpdateMetricBaselineV2InternalServerError().WithPayload(errorMessage)
		}
	}
}

// HandleDeleteMetricBaselineV2 - delete a MetricBaseline record by ID.
func HandleDeleteMetricBaselineV2(allowedRoles []string, tenantDB datastore.TenantMetricBaselineDatastore) func(params tenant_provisioning_service_v2.DeleteMetricBaselineV2Params) middleware.Responder {
	return func(params tenant_provisioning_service_v2.DeleteMetricBaselineV2Params) middleware.Responder {

		// Do the work
		startTime, responseCode, response, err := doDeleteMetricBaselineV2(allowedRoles, tenantDB, params)

		// Success Response
		if responseCode == http.StatusOK {
			return tenant_provisioning_service_v2.NewDeleteMetricBaselineV2OK().WithPayload(response)
		}

		// Error Responses
		errorMessage := reportAPIError(err.Error(), startTime, responseCode, mon.DeleteTenantMetricBaselineStr, mon.APICompleted, mon.TenantAPICompleted)
		switch responseCode {
		case http.StatusForbidden:
			return tenant_provisioning_service_v2.NewDeleteMetricBaselineV2Forbidden().WithPayload(errorMessage)
		case http.StatusNotFound:
			return tenant_provisioning_service_v2.NewDeleteMetricBaselineV2NotFound().WithPayload(errorMessage)
		default:
			return tenant_provisioning_service_v2.NewDeleteMetricBaselineV2InternalServerError().WithPayload(errorMessage)
		}

	}
}

func HandleGetMetricBaselineByMonitoredObjectIdForHourOfWeekV2(allowedRoles []string, tenantDB datastore.TenantMetricBaselineDatastore, fetchLimiter *MetricBaselineFetchLimiter) func(params tenant_provisioning_service_v2.GetMetricBaselineByMonitoredObjectIDForHourOfWeekV2Params) middleware.Responder {
	return func(params tenant_provisioning_service_v2.GetMetricBaselineByMonitoredObjectIDForHourOfWeekV2Params) middleware.Responder {

		// Do the work
		startTime, responseCode, response, err := doGetMetricBaselineByMonitoredObjectIDForHourOfWeekV2(allowedRoles, tenantDB, params, fetchLimiter)

		// Success Response
		if responseCode == http.StatusOK {
			return tenant_provisioning_service_v2.NewGetMetricBaselineByMonitoredObjectIDForHourOfWeekV2OK().WithPayload(response)
		}

		// Error Responses
		errorMessage := reportAPIError(err.Error(), startTime, responseCode, mon.GetMetricBaselineByMonitoredObjectIdForHourOfWeekStr, mon.APICompleted, mon.TenantAPICompleted)
		switch responseCode {
		case http.StatusForbidden:
			return tenant_provisioning_service_v2.NewGetMetricBaselineByMonitoredObjectIDForHourOfWeekV2Forbidden().WithPayload(errorMessage)
		case http.StatusNotFound:
			return tenant_provisioning_service_v2.NewGetMetricBaselineByMonitoredObjectIDForHourOfWeekV2NotFound().WithPayload(errorMessage)
		default:
			return tenant_provisioning_service_v2.NewGetMetricBaselineByMonitoredObjectIDForHourOfWeekV2InternalServerError().WithPayload(errorMessage)
		}

	}
}

func HandleUpdateMetricBaselineForHourOfWeekV2(allowedRoles []string, tenantDB datastore.TenantMetricBaselineDatastore) func(params tenant_provisioning_service_v2.UpdateMetricBaselineForHourOfWeekV2Params) middleware.Responder {
	return func(params tenant_provisioning_service_v2.UpdateMetricBaselineForHourOfWeekV2Params) middleware.Responder {
		// Do the work
		startTime, responseCode, response, err := doUpdateMetricBaselineForHourOfWeekV2(allowedRoles, tenantDB, params)

		// Success Response
		if responseCode == http.StatusOK {
			return tenant_provisioning_service_v2.NewUpdateMetricBaselineForHourOfWeekV2OK().WithPayload(response)
		}

		// Error Responses
		errorMessage := reportAPIError(err.Error(), startTime, responseCode, mon.UpdateMetricBaselineForHourOfWeekV2Str, mon.APICompleted, mon.TenantAPICompleted)
		switch responseCode {
		case http.StatusForbidden:
			return tenant_provisioning_service_v2.NewUpdateMetricBaselineForHourOfWeekV2Forbidden().WithPayload(errorMessage)
		case http.StatusNotFound:
			return tenant_provisioning_service_v2.NewUpdateMetricBaselineForHourOfWeekV2NotFound().WithPayload(errorMessage)
		default:
			return tenant_provisioning_service_v2.NewUpdateMetricBaselineForHourOfWeekV2InternalServerError().WithPayload(errorMessage)
		}
	}
}

// HandleBulkUpdateMetricBaselinesV2 - update multiple MetricBaseline records for Multiple Monitored Objects
func HandleBulkUpdateMetricBaselinesV2(allowedRoles []string, tenantDB datastore.TenantMetricBaselineDatastore, bulkUpdateManager *MetricBaselineBulkUpdateManager) func(params tenant_provisioning_service_v2.BulkUpdateMetricBaselinesV2Params) middleware.Responder {
	return func(params tenant_provisioning_service_v2.BulkUpdateMetricBaselinesV2Params) middleware.Responder {

		// Do the work
		startTime, responseCode, response, err := doBulkUpdateMetricBaselineV2(allowedRoles, tenantDB, params, bulkUpdateManager)

		// Success Response
		if responseCode == http.StatusOK {
			return tenant_provisioning_service_v2.NewBulkUpdateMetricBaselinesV2OK().WithPayload(response)
		}

		// Error Responses
		errorMessage := reportAPIError(err.Error(), startTime, responseCode, mon.BulkUpdateTenantMetricBaselineStr, mon.APICompleted, mon.TenantAPICompleted)
		switch responseCode {
		case http.StatusForbidden:
			return tenant_provisioning_service_v2.NewBulkUpdateMetricBaselinesV2Forbidden().WithPayload(errorMessage)
		case http.StatusBadRequest:
			return tenant_provisioning_service_v2.NewBulkUpdateMetricBaselinesV2BadRequest().WithPayload(errorMessage)
		default:
			return tenant_provisioning_service_v2.NewBulkUpdateMetricBaselinesV2InternalServerError().WithPayload(errorMessage)
		}
	}
}

func HandleGetMetricBaselineByMonitoredObjectIDV2(allowedRoles []string, tenantDB datastore.TenantMetricBaselineDatastore) func(params tenant_provisioning_service_v2.GetMetricBaselineByMonitoredObjectIDV2Params) middleware.Responder {
	return func(params tenant_provisioning_service_v2.GetMetricBaselineByMonitoredObjectIDV2Params) middleware.Responder {

		// Do the work
		startTime, responseCode, response, err := doGetMetricBaselineByMonitoredObjectIDV2(allowedRoles, tenantDB, params)

		// Success Response
		if responseCode == http.StatusOK {
			return tenant_provisioning_service_v2.NewGetMetricBaselineByMonitoredObjectIDV2OK().WithPayload(response)
		}

		// Error Responses
		errorMessage := reportAPIError(err.Error(), startTime, responseCode, mon.GetMetricBaselineByMonitoredObjectIdStr, mon.APICompleted, mon.TenantAPICompleted)
		switch responseCode {
		case http.StatusForbidden:
			return tenant_provisioning_service_v2.NewGetMetricBaselineByMonitoredObjectIDV2Forbidden().WithPayload(errorMessage)
		case http.StatusNotFound:
			return tenant_provisioning_service_v2.NewGetMetricBaselineByMonitoredObjectIDV2NotFound().WithPayload(errorMessage)
		default:
			return tenant_provisioning_service_v2.NewGetMetricBaselineByMonitoredObjectIDV2InternalServerError().WithPayload(errorMessage)
		}

	}
}

func HandleDeleteMetricBaselineByMonitoredObjectIDV2(allowedRoles []string, tenantDB datastore.TenantMetricBaselineDatastore) func(params tenant_provisioning_service_v2.DeleteMetricBaselineByMonitoredObjectIDV2Params) middleware.Responder {
	return func(params tenant_provisioning_service_v2.DeleteMetricBaselineByMonitoredObjectIDV2Params) middleware.Responder {

		// Do the work
		startTime, responseCode, response, err := doDeleteMetricBaselineByMonitoredObjectIDV2(allowedRoles, tenantDB, params)

		// Success Response
		if responseCode == http.StatusOK {
			return tenant_provisioning_service_v2.NewDeleteMetricBaselineByMonitoredObjectIDV2OK().WithPayload(response)
		}

		// Error Responses
		errorMessage := reportAPIError(err.Error(), startTime, responseCode, mon.GetMetricBaselineByMonitoredObjectIdStr, mon.APICompleted, mon.TenantAPICompleted)
		switch responseCode {
		case http.StatusForbidden:
			return tenant_provisioning_service_v2.NewDeleteMetricBaselineByMonitoredObjectIDV2Forbidden().WithPayload(errorMessage)
		case http.StatusNotFound:
			return tenant_provisioning_service_v2.NewDeleteMetricBaselineByMonitoredObjectIDV2NotFound().WithPayload(errorMessage)
		default:
			return tenant_provisioning_service_v2.NewDeleteMetricBaselineByMonitoredObjectIDV2InternalServerError().WithPayload(errorMessage)
		}

	}
}

func doCreateMetricBaselineV2(allowedRoles []string, tenantDB datastore.TenantMetricBaselineDatastore, params tenant_provisioning_service_v2.CreateMetricBaselineV2Params) (time.Time, int, *swagmodels.MetricBaselineResponse, error) {
	tenantID := params.HTTPRequest.Header.Get(XFwdTenantId)
	isAuthorized, startTime := authorizeRequest(fmt.Sprintf("Creating %s for tenant %s", tenmod.TenantMetricBaselineStr, tenantID), params.HTTPRequest, allowedRoles, mon.APIRecieved, mon.TenantAPIRecieved)

	if !isAuthorized {
		return startTime, http.StatusForbidden, nil, fmt.Errorf("Create %s operation not authorized for role: %s", tenmod.TenantMetricBaselineStr, params.HTTPRequest.Header.Get(XFwdUserRoles))
	}

	// Unmarshal the request
	requestBytes, err := json.Marshal(params.Body)
	if err != nil {
		return startTime, http.StatusBadRequest, nil, err
	}

	data := tenmod.MetricBaseline{}
	err = jsonapi.Unmarshal(requestBytes, &data)
	if err != nil {
		return startTime, http.StatusBadRequest, nil, err
	}

	data.TenantID = tenantID

	// Issue request to DAO Layer to Create Record
	result, err := tenantDB.CreateMetricBaseline(&data)
	if err != nil {
		if strings.Contains(err.Error(), datastore.ConflictErrorStr) {
			return startTime, http.StatusConflict, nil, err
		}
		return startTime, http.StatusInternalServerError, nil, fmt.Errorf("Unable to store %s: %s", tenmod.TenantMetricBaselineStr, err.Error())
	}

	converted := swagmodels.MetricBaselineResponse{}
	err = convertToJsonapiObject(result, &converted)
	if err != nil {
		return startTime, http.StatusInternalServerError, nil, fmt.Errorf("Unable to convert %s data to jsonapi return format: %s", tenmod.TenantMetricBaselineStr, err.Error())
	}

	reportAPICompletionState(startTime, http.StatusOK, mon.CreateTenantMetricBaselineStr, mon.APICompleted, mon.TenantAPICompleted)
	logger.Log.Infof("Created %s %s", tenmod.TenantMetricBaselineStr, *converted.Data.ID)
	return startTime, http.StatusCreated, &converted, nil
}

func doGetMetricBaselineV2(allowedRoles []string, tenantDB datastore.TenantMetricBaselineDatastore, params tenant_provisioning_service_v2.GetMetricBaselineV2Params) (time.Time, int, *swagmodels.MetricBaselineResponse, error) {
	tenantID := params.HTTPRequest.Header.Get(XFwdTenantId)
	isAuthorized, startTime := authorizeRequest(fmt.Sprintf("Fetching %s %s for %s %s", tenmod.TenantMetricBaselineStr, params.MetricBaselineID, admmod.TenantStr, tenantID), params.HTTPRequest, allowedRoles, mon.APIRecieved, mon.AdminAPIRecieved)

	if !isAuthorized {
		return startTime, http.StatusForbidden, nil, fmt.Errorf("Fetch %s operation not authorized for role: %s", tenmod.TenantMetricBaselineStr, params.HTTPRequest.Header.Get(XFwdUserRoles))
	}

	// Issue request to DAO Layer
	result, err := tenantDB.GetMetricBaseline(tenantID, params.MetricBaselineID)
	if err != nil {
		if strings.Contains(err.Error(), datastore.NotFoundStr) {
			return startTime, http.StatusNotFound, nil, err
		}

		return startTime, http.StatusInternalServerError, nil, fmt.Errorf("Unable to retrieve %s: %s", tenmod.TenantMetricBaselineStr, err.Error())
	}

	converted := swagmodels.MetricBaselineResponse{}
	err = convertToJsonapiObject(result, &converted)
	if err != nil {
		return startTime, http.StatusInternalServerError, nil, fmt.Errorf("Unable to convert %s data to jsonapi return format: %s", tenmod.TenantMetricBaselineStr, err.Error())
	}

	reportAPICompletionState(startTime, http.StatusOK, mon.GetTenantMetricBaselineStr, mon.APICompleted, mon.TenantAPICompleted)
	logger.Log.Infof("Retrieved %s %s", tenmod.TenantMetricBaselineStr, params.MetricBaselineID)
	return startTime, http.StatusOK, &converted, nil
}

func doUpdateMetricBaselineV2(allowedRoles []string, tenantDB datastore.TenantMetricBaselineDatastore, params tenant_provisioning_service_v2.UpdateMetricBaselineV2Params) (time.Time, int, *swagmodels.MetricBaselineResponse, error) {
	tenantID := params.HTTPRequest.Header.Get(XFwdTenantId)
	isAuthorized, startTime := authorizeRequest(fmt.Sprintf("Updating %s %s for %s %s", tenmod.TenantMetricBaselineStr, params.MetricBaselineID, admmod.TenantStr, tenantID), params.HTTPRequest, allowedRoles, mon.APIRecieved, mon.AdminAPIRecieved)

	if !isAuthorized {
		return startTime, http.StatusForbidden, nil, fmt.Errorf("Update %s operation not authorized for role: %s", tenmod.TenantMetricBaselineStr, params.HTTPRequest.Header.Get(XFwdUserRoles))
	}

	// Retrieve tne request bytes from the payload in order to convert it to a map
	patchRequestBytes, err := json.Marshal(params.Body)
	if err != nil || params.Body == nil {
		return startTime, http.StatusBadRequest, nil, fmt.Errorf("Invalid request body: %s", models.AsJSONString(params.Body))
	}

	// Attempt to retrieve the record that we are trying to patch from the DB in order to do a merge
	fetched, err := tenantDB.GetMetricBaseline(tenantID, params.MetricBaselineID)
	if err != nil {
		if strings.Contains(err.Error(), datastore.NotFoundStr) {
			return startTime, http.StatusNotFound, nil, err
		}

		return startTime, http.StatusInternalServerError, nil, fmt.Errorf("Unable to fetch %s: %s", tenmod.TenantMetricBaselineStr, err.Error())
	}

	// Merge the attributes passed in with the patch request to the record fetched from the datastore
	patched := &tenmod.MetricBaseline{}
	if err := models.MergeObjWithMap(patched, fetched, patchRequestBytes); err != nil {
		return startTime, http.StatusInternalServerError, nil, fmt.Errorf("Unable to patch %s with id %s: %s", tenmod.TenantMetricBaselineStr, params.MetricBaselineID, err.Error())
	}
	patched.TenantID = tenantID

	// Finally update the record in the datastore with the merged map and fetched tenant
	result, err := tenantDB.UpdateMetricBaseline(patched)
	if err != nil {
		if strings.Contains(err.Error(), datastore.ConflictErrorStr) {
			return startTime, http.StatusConflict, nil, err
		}

		return startTime, http.StatusInternalServerError, nil, err
	}

	converted := swagmodels.MetricBaselineResponse{}
	err = convertToJsonapiObject(result, &converted)
	if err != nil {
		return startTime, http.StatusInternalServerError, nil, fmt.Errorf("Unable to convert %s data to jsonapi return format: %s", tenmod.TenantMetricBaselineStr, err.Error())
	}

	reportAPICompletionState(startTime, http.StatusOK, mon.UpdateTenantMetricBaselineStr, mon.APICompleted, mon.TenantAPICompleted)
	logger.Log.Infof("Updated %s %s", tenmod.TenantMetricBaselineStr, params.MetricBaselineID)
	return startTime, http.StatusOK, &converted, nil
}

func doDeleteMetricBaselineV2(allowedRoles []string, tenantDB datastore.TenantMetricBaselineDatastore, params tenant_provisioning_service_v2.DeleteMetricBaselineV2Params) (time.Time, int, *swagmodels.MetricBaselineResponse, error) {
	tenantID := params.HTTPRequest.Header.Get(XFwdTenantId)
	isAuthorized, startTime := authorizeRequest(fmt.Sprintf("Deleting %s %s for %s %s", tenmod.TenantMetricBaselineStr, params.MetricBaselineID, admmod.TenantStr, tenantID), params.HTTPRequest, allowedRoles, mon.APIRecieved, mon.AdminAPIRecieved)

	if !isAuthorized {
		return startTime, http.StatusForbidden, nil, fmt.Errorf("Delete %s operation not authorized for role: %s", tenmod.TenantMetricBaselineStr, params.HTTPRequest.Header.Get(XFwdUserRoles))
	}

	// Issue request to DAO Layer
	result, err := tenantDB.DeleteMetricBaseline(tenantID, params.MetricBaselineID)
	if err != nil {
		if strings.Contains(err.Error(), datastore.NotFoundStr) {
			return startTime, http.StatusNotFound, nil, err
		}

		return startTime, http.StatusInternalServerError, nil, fmt.Errorf("Unable to delete %s: %s", tenmod.TenantMetricBaselineStr, err.Error())
	}

	converted := swagmodels.MetricBaselineResponse{}
	err = convertToJsonapiObject(result, &converted)
	if err != nil {
		return startTime, http.StatusInternalServerError, nil, fmt.Errorf("Unable to convert %s data to jsonapi return format: %s", tenmod.TenantMetricBaselineStr, err.Error())
	}

	reportAPICompletionState(startTime, http.StatusOK, mon.DeleteTenantMetricBaselineStr, mon.APICompleted, mon.TenantAPICompleted)
	logger.Log.Infof("Deleted %s %s", tenmod.TenantMetricBaselineStr, params.MetricBaselineID)
	return startTime, http.StatusOK, &converted, nil
}

func doGetMetricBaselineByMonitoredObjectIDForHourOfWeekV2(allowedRoles []string, tenantDB datastore.TenantMetricBaselineDatastore, params tenant_provisioning_service_v2.GetMetricBaselineByMonitoredObjectIDForHourOfWeekV2Params, fetchLimiter *MetricBaselineFetchLimiter) (time.Time, int, *swagmodels.MetricBaselineResponse, error) {
	tenantID := params.HTTPRequest.Header.Get(XFwdTenantId)
	isAuthorized, startTime := authorizeRequest(fmt.Sprintf("Fetching %s for %s %s for %s %s and hourOfWeek %d", tenmod.TenantMetricBaselineStr, admmod.TenantStr, tenantID, tenmod.TenantMonitoredObjectStr, params.MonitoredObjectID, params.HourOfWeek), params.HTTPRequest, allowedRoles, mon.APIRecieved, mon.AdminAPIRecieved)

	if !isAuthorized {
		return startTime, http.StatusForbidden, nil, fmt.Errorf("Fetch %s operation not authorized for role: %s", tenmod.TenantMetricBaselineStr, params.HTTPRequest.Header.Get(XFwdUserRoles))
	}

	// Issue request to DAO Layer
	dataID := fmt.Sprintf("%s_%d", params.MonitoredObjectID, params.HourOfWeek)
	// result, err := tenantDB.GetMetricBaseline(tenantID, dataID)
	// if err != nil {
	// 	if strings.Contains(err.Error(), datastore.NotFoundStr) {
	// 		return startTime, http.StatusNotFound, nil, err
	// 	}

	// 	return startTime, http.StatusInternalServerError, nil, fmt.Errorf("Unable to retrieve %s: %s", tenmod.TenantMetricBaselineStr, err.Error())
	// }
	daoResult := make(chan limitedFetchResult, 1)
	job := limitedFetchJob{
		dataID:    dataID,
		startTime: startTime,
		tenantID:  tenantID,
		result:    daoResult,
	}
	fetchLimiter.jobs <- job

	jobResp := <-daoResult
	if jobResp.err != nil {
		return startTime, jobResp.resultCode, nil, fmt.Errorf("Unable to fetch %s for %s %s and hourOfWeek %d: %s", tenmod.TenantMetricBaselineStr, tenmod.TenantMonitoredObjectStr, params.MonitoredObjectID, params.HourOfWeek, jobResp.err.Error())
	}

	close(daoResult)

	converted := swagmodels.MetricBaselineResponse{}
	err := convertToJsonapiObject(jobResp.result, &converted)
	if err != nil {
		return startTime, http.StatusInternalServerError, nil, fmt.Errorf("Unable to convert %s data to jsonapi return format: %s", tenmod.TenantMetricBaselineStr, err.Error())
	}

	reportAPICompletionState(startTime, http.StatusOK, mon.GetMetricBaselineByMonitoredObjectIdForHourOfWeekStr, mon.APICompleted, mon.TenantAPICompleted)
	logger.Log.Infof("Retrieved %s for Monitored Object %s", tenmod.TenantMetricBaselineStr, params.MonitoredObjectID)
	return startTime, http.StatusOK, &converted, nil
}

func doUpdateMetricBaselineForHourOfWeekV2(allowedRoles []string, tenantDB datastore.TenantMetricBaselineDatastore, params tenant_provisioning_service_v2.UpdateMetricBaselineForHourOfWeekV2Params) (time.Time, int, *swagmodels.MetricBaselineResponse, error) {
	tenantID := params.HTTPRequest.Header.Get(XFwdTenantId)
	isAuthorized, startTime := authorizeRequest(fmt.Sprintf("Updateing %s for %s %s for %s %s", tenmod.TenantMetricBaselineStr, admmod.TenantStr, tenantID, tenmod.TenantMonitoredObjectStr, params.MonitoredObjectID), params.HTTPRequest, allowedRoles, mon.APIRecieved, mon.AdminAPIRecieved)

	if !isAuthorized {
		return startTime, http.StatusForbidden, nil, fmt.Errorf("Fetch %s operation not authorized for role: %s", tenmod.TenantMetricBaselineStr, params.HTTPRequest.Header.Get(XFwdUserRoles))
	}

	// Retrieve tne request bytes from the payload in order to convert it to a map
	patchRequestBytes, err := json.Marshal(params.Body.Data.Attributes.Baselines)
	if err != nil || params.Body == nil {
		return startTime, http.StatusBadRequest, nil, fmt.Errorf("Invalid request body: %s", models.AsJSONString(params.Body))
	}

	patchObject := []*tenmod.MetricBaselineData{}
	err = json.Unmarshal(patchRequestBytes, &patchObject)
	if err != nil {
		return startTime, http.StatusBadRequest, nil, fmt.Errorf("Invalid request body: %s", models.AsJSONString(params.Body))
	}

	// Issue request to DAO Layer
	result, err := tenantDB.UpdateMetricBaselineForHourOfWeekWithCollection(tenantID, params.MonitoredObjectID, params.HourOfWeek, patchObject)
	if err != nil {
		if strings.Contains(err.Error(), datastore.NotFoundStr) {
			return startTime, http.StatusNotFound, nil, err
		}

		return startTime, http.StatusInternalServerError, nil, fmt.Errorf("Unable to retrieve %s: %s", tenmod.TenantMetricBaselineStr, err.Error())
	}

	converted := swagmodels.MetricBaselineResponse{}
	err = convertToJsonapiObject(result, &converted)
	if err != nil {
		return startTime, http.StatusInternalServerError, nil, fmt.Errorf("Unable to convert %s data to jsonapi return format: %s", tenmod.TenantMetricBaselineStr, err.Error())
	}

	reportAPICompletionState(startTime, http.StatusOK, mon.UpdateMetricBaselineForHourOfWeekV2Str, mon.APICompleted, mon.TenantAPICompleted)
	logger.Log.Infof("Updated %s for Monitored Object %s", tenmod.TenantMetricBaselineStr, params.MonitoredObjectID)
	return startTime, http.StatusOK, &converted, nil
}

func doBulkUpdateMetricBaselineV2(allowedRoles []string, tenantDB datastore.TenantMetricBaselineDatastore, params tenant_provisioning_service_v2.BulkUpdateMetricBaselinesV2Params, bulkUpdateManager *MetricBaselineBulkUpdateManager) (time.Time, int, string, error) {
	startTime := time.Now()

	tenantID := params.HTTPRequest.Header.Get(XFwdTenantId)
	isAuthorized, startTime := authorizeRequest(fmt.Sprintf("Bulk Updating %ss for %s %s", tenmod.TenantMetricBaselineStr, admmod.TenantStr, tenantID), params.HTTPRequest, allowedRoles, mon.APIRecieved, mon.AdminAPIRecieved)

	if !isAuthorized {
		return startTime, http.StatusForbidden, "", fmt.Errorf("Bulk Update %s operation not authorized for role: %s", tenmod.TenantMetricBaselineStr, params.HTTPRequest.Header.Get(XFwdUserRoles))
	}

	go func() {
		job := bulkUpdateJob{
			params:    params,
			startTime: startTime,
			tenantID:  tenantID,
		}
		bulkUpdateManager.jobs <- job
	}()

	reportAPICompletionState(startTime, http.StatusOK, mon.BulkUpdateTenantMetricBaselineStr, mon.APICompleted, mon.TenantAPICompleted)
	return startTime, http.StatusOK, "Bulk update request accepted", nil
}

func doGetMetricBaselineByMonitoredObjectIDV2(allowedRoles []string, tenantDB datastore.TenantMetricBaselineDatastore, params tenant_provisioning_service_v2.GetMetricBaselineByMonitoredObjectIDV2Params) (time.Time, int, *swagmodels.MetricBaselineListResponse, error) {
	tenantID := params.HTTPRequest.Header.Get(XFwdTenantId)
	isAuthorized, startTime := authorizeRequest(fmt.Sprintf("Fetching all %ss for %s %s for %s %s", tenmod.TenantMetricBaselineStr, admmod.TenantStr, tenantID, tenmod.TenantMonitoredObjectStr, params.MonitoredObjectID), params.HTTPRequest, allowedRoles, mon.APIRecieved, mon.AdminAPIRecieved)

	if !isAuthorized {
		return startTime, http.StatusForbidden, nil, fmt.Errorf("Fetch all %s operation not authorized for role: %s", tenmod.TenantMetricBaselineStr, params.HTTPRequest.Header.Get(XFwdUserRoles))
	}

	// Issue request to DAO Layer
	result, err := tenantDB.GetMetricBaselineForMonitoredObject(tenantID, params.MonitoredObjectID)
	if err != nil {
		if strings.Contains(err.Error(), datastore.NotFoundStr) {
			return startTime, http.StatusNotFound, nil, err
		}

		return startTime, http.StatusInternalServerError, nil, fmt.Errorf("Unable to retrieve %ss for %s %s: %s", tenmod.TenantMetricBaselineStr, tenmod.TenantMonitoredObjectStr, params.MonitoredObjectID, err.Error())
	}

	converted := swagmodels.MetricBaselineListResponse{}
	err = convertToJsonapiObject(result, &converted)
	if err != nil {
		return startTime, http.StatusInternalServerError, nil, fmt.Errorf("Unable to convert %s data to jsonapi return format: %s", tenmod.TenantMetricBaselineStr, err.Error())
	}

	reportAPICompletionState(startTime, http.StatusOK, mon.GetMetricBaselineByMonitoredObjectIdStr, mon.APICompleted, mon.TenantAPICompleted)
	logger.Log.Infof("Retrieved %d %ss for %s %s", len(result), tenmod.TenantMetricBaselineStr, tenmod.TenantMonitoredObjectStr, params.MonitoredObjectID)
	return startTime, http.StatusOK, &converted, nil
}

func doDeleteMetricBaselineByMonitoredObjectIDV2(allowedRoles []string, tenantDB datastore.TenantMetricBaselineDatastore, params tenant_provisioning_service_v2.DeleteMetricBaselineByMonitoredObjectIDV2Params) (time.Time, int, string, error) {
	tenantID := params.HTTPRequest.Header.Get(XFwdTenantId)
	operationType := "Delete"
	if params.Reset {
		operationType = "Reset"
	}
	isAuthorized, startTime := authorizeRequest(fmt.Sprintf("%s all %ss for %s %s for %s %s", operationType, tenmod.TenantMetricBaselineStr, admmod.TenantStr, tenantID, tenmod.TenantMonitoredObjectStr, params.MonitoredObjectID), params.HTTPRequest, allowedRoles, mon.APIRecieved, mon.AdminAPIRecieved)

	if !isAuthorized {
		return startTime, http.StatusForbidden, "", fmt.Errorf("%s all %s operation not authorized for role: %s", operationType, tenmod.TenantMetricBaselineStr, params.HTTPRequest.Header.Get(XFwdUserRoles))
	}

	// Issue request to DAO Layer
	err := tenantDB.DeleteMetricBaselineForMonitoredObject(tenantID, params.MonitoredObjectID, params.Reset)
	if err != nil {
		if strings.Contains(err.Error(), datastore.NotFoundStr) {
			return startTime, http.StatusNotFound, "", err
		}

		return startTime, http.StatusInternalServerError, "", fmt.Errorf("Unable to %s %ss for %s %s: %s", strings.ToLower(operationType), tenmod.TenantMetricBaselineStr, tenmod.TenantMonitoredObjectStr, params.MonitoredObjectID, err.Error())
	}

	reportAPICompletionState(startTime, http.StatusOK, mon.DeleteMetricBaselineByMonitoredObjectIdStr, mon.APICompleted, mon.TenantAPICompleted)
	successMsg := fmt.Sprintf("%s successful for %ss for %s %s", operationType, tenmod.TenantMetricBaselineStr, tenmod.TenantMonitoredObjectStr, params.MonitoredObjectID)
	logger.Log.Infof(successMsg)
	return startTime, http.StatusOK, successMsg, nil
}

type MetricBaselineBulkUpdateManager struct {
	jobs     chan bulkUpdateJob
	results  chan bulkUpdateResult
	tenantDB datastore.TenantMetricBaselineDatastore
}

type bulkUpdateJob struct {
	startTime time.Time
	tenantID  string
	params    tenant_provisioning_service_v2.BulkUpdateMetricBaselinesV2Params
}

type bulkUpdateResult struct {
	startTime  time.Time
	resultCode int
	// result     *swagmodels.MetricBaselineBulkUpdateResponse
	err error
}

func CreateMetricBaselineBulkUpdateManager(db datastore.TenantMetricBaselineDatastore) *MetricBaselineBulkUpdateManager {
	cfg := gather.GetConfig()

	numJobs := cfg.GetInt(gather.CK_args_metricbaselines_maxnumjobs.String())
	numWorkers := cfg.GetInt(gather.CK_args_metricbaselines_numworkers.String())

	result := MetricBaselineBulkUpdateManager{
		jobs:     make(chan bulkUpdateJob, numJobs),
		results:  make(chan bulkUpdateResult, numJobs),
		tenantDB: db,
	}

	for w := 1; w <= numWorkers; w++ {
		go result.bulkMetricBaselineUpdateWorker(w, result.jobs, result.results)
	}

	return &result
}

func (manager *MetricBaselineBulkUpdateManager) bulkMetricBaselineUpdateWorker(id int, jobs <-chan bulkUpdateJob, results chan<- bulkUpdateResult) {
	for j := range jobs {
		// startTime, resultCode, result, err := manager.performBulkUpdate(j.startTime, j.tenantID, j.params)
		_, _, err := manager.performBulkUpdate(j.startTime, j.tenantID, j.params)
		if err != nil {
			logger.Log.Error(err.Error())
		}
		// results <- bulkUpdateResult{
		// 	result:     result,
		// 	err:        err,
		// 	resultCode: resultCode,
		// 	startTime:  startTime,
		// }
	}
}

func (manager *MetricBaselineBulkUpdateManager) performBulkUpdate(startTime time.Time, tenantID string, params tenant_provisioning_service_v2.BulkUpdateMetricBaselinesV2Params) (time.Time, int, error) {
	_, err := manager.tenantDB.BulkUpdateMetricBaselines(tenantID, params.Body.Data.Attributes)
	if err != nil {
		return startTime, http.StatusInternalServerError, fmt.Errorf("Unable to bulk update %s: %s", tenmod.TenantMetricBaselineStr, err.Error())
	}

	// reportAPICompletionState(startTime, http.StatusOK, mon.BulkUpdateTenantMetricBaselineStr, mon.APICompleted, mon.TenantAPICompleted)
	logger.Log.Infof("Bulk Update of %ss for Tenant %s complete", tenmod.TenantMetricBaselineStr, tenantID)
	return startTime, http.StatusOK, nil
}

type MetricBaselineFetchLimiter struct {
	jobs     chan limitedFetchJob
	tenantDB datastore.TenantMetricBaselineDatastore
}

type limitedFetchJob struct {
	startTime time.Time
	tenantID  string
	dataID    string
	result    chan limitedFetchResult
}

type limitedFetchResult struct {
	startTime  time.Time
	resultCode int
	result     *tenmod.MetricBaseline
	err        error
}

func CreateLimitedFetchManager(db datastore.TenantMetricBaselineDatastore) *MetricBaselineFetchLimiter {
	cfg := gather.GetConfig()

	numJobs := cfg.GetInt(gather.CK_args_metricbaselines_maxnumjobs.String())
	numWorkers := cfg.GetInt(gather.CK_args_metricbaselines_numworkers.String())

	result := MetricBaselineFetchLimiter{
		jobs:     make(chan limitedFetchJob, numJobs),
		tenantDB: db,
	}

	for w := 1; w <= numWorkers; w++ {
		go result.limitedFetchWorker(w, result.jobs)
	}

	return &result
}

func (manager *MetricBaselineFetchLimiter) limitedFetchWorker(id int, jobs <-chan limitedFetchJob) {
	for j := range jobs {
		// // startTime, resultCode, result, err := manager.performBulkUpdate(j.startTime, j.tenantID, j.params)
		// fetchRes, code, err := manager.performLimitedFetch(j.startTime, j.tenantID, j.dataID)
		// if err != nil {
		// 	logger.Log.Error(err.Error())
		// }

		fetchRes, err := manager.tenantDB.GetMetricBaseline(j.tenantID, j.dataID)
		if err != nil {
			if strings.Contains(err.Error(), datastore.NotFoundStr) {
				res := limitedFetchResult{
					startTime:  j.startTime,
					resultCode: http.StatusNotFound,
					result:     nil,
					err:        err,
				}

				j.result <- res
				continue
			}

			res := limitedFetchResult{
				startTime:  j.startTime,
				resultCode: http.StatusInternalServerError,
				result:     nil,
				err:        fmt.Errorf("Unable to retrieve %s: %s", tenmod.TenantMetricBaselineStr, err.Error()),
			}

			j.result <- res
			continue
		}

		res := limitedFetchResult{
			startTime:  j.startTime,
			resultCode: http.StatusOK,
			result:     fetchRes,
			err:        err,
		}

		j.result <- res
	}
}
