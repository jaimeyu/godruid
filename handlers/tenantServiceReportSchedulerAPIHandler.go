package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/accedian/adh-gather/datastore"
	"github.com/accedian/adh-gather/logger"
	"github.com/accedian/adh-gather/models"
	metmod "github.com/accedian/adh-gather/models/metrics"
	mon "github.com/accedian/adh-gather/monitoring"
	"github.com/accedian/adh-gather/restapi/operations/tenant_provisioning_service"
	"github.com/accedian/adh-gather/scheduler"
	"github.com/accedian/adh-gather/swagmodels"
	"github.com/go-openapi/runtime/middleware"
	"github.com/manyminds/api2go/jsonapi"
)

// HandleCreateReportScheduleConfig - creates a report schedule config for a tenant
func HandleCreateReportScheduleConfig(allowedRoles []string, tenantDB datastore.TenantServiceDatastore) func(params tenant_provisioning_service.CreateReportScheduleConfigParams) middleware.Responder {
	return func(params tenant_provisioning_service.CreateReportScheduleConfigParams) middleware.Responder {
		startTime := time.Now()
		incrementAPICounters(mon.APIRecieved, mon.TenantAPIRecieved)
		logger.Log.Infof("Creating %s %s for Tenant %s", metmod.ReportScheduleConfigStr, params.Body.Data.Attributes.Name, params.TenantID)

		if !isRequestAuthorized(params.HTTPRequest, allowedRoles) {
			return tenant_provisioning_service.NewCreateReportScheduleConfigForbidden().WithPayload(reportAPIError(fmt.Sprintf("Update %s operation not authorized for role: %s", metmod.ReportScheduleConfigStr, params.HTTPRequest.Header.Get(XFwdUserRoles)), startTime, http.StatusForbidden, mon.CreateReportScheduleConfigStr, mon.APICompleted, mon.TenantAPICompleted))
		}

		// Ensure that we can unmarshal the provided report schedule payload into the model object
		requestBytes, err := json.Marshal(params.Body)
		if err != nil {
			return tenant_provisioning_service.NewCreateReportScheduleConfigBadRequest().WithPayload(reportAPIError(generateErrorMessage(http.StatusBadRequest, err.Error()), startTime, http.StatusBadRequest, mon.CreateReportScheduleConfigStr, mon.APICompleted, mon.TenantAPICompleted))
		}
		data := metmod.ReportScheduleConfig{}
		err = jsonapi.Unmarshal(requestBytes, &data)
		if err != nil {
			return tenant_provisioning_service.NewCreateReportScheduleConfigBadRequest().WithPayload(reportAPIError(generateErrorMessage(http.StatusBadRequest, err.Error()), startTime, http.StatusBadRequest, mon.CreateReportScheduleConfigStr, mon.APICompleted, mon.TenantAPICompleted))
		}
		// Seconds is not really used except for testing.
		data.Second = "0"

		// Ensure that the passed in data adheres to the model requirements
		err = data.Validate(false)
		if err != nil {
			return tenant_provisioning_service.NewCreateReportScheduleConfigBadRequest().WithPayload(reportAPIError(generateErrorMessage(http.StatusBadRequest, err.Error()), startTime, http.StatusBadRequest, mon.CreateReportScheduleConfigStr, mon.APICompleted, mon.TenantAPICompleted))
		}

		// Attempt to create a config entry in the datastore for the scheduler to pick up
		result, err := tenantDB.CreateReportScheduleConfig(&data)
		if err != nil {
			return tenant_provisioning_service.NewCreateReportScheduleConfigInternalServerError().WithPayload(reportAPIError(fmt.Sprintf("Unable to store %s: %s", metmod.ReportScheduleConfigStr, err.Error()), startTime, http.StatusInternalServerError, mon.CreateReportScheduleConfigStr, mon.APICompleted, mon.TenantAPICompleted))
		}

		// Tell the scheduler to go and update based on the updated database
		err = scheduler.RebuildCronJobs()
		if err != nil {
			return tenant_provisioning_service.NewCreateReportScheduleConfigInternalServerError().WithPayload(reportAPIError(fmt.Sprintf("Unable to start scheduled job%s: %s", metmod.ReportScheduleConfigStr, err.Error()), startTime, http.StatusInternalServerError, mon.CreateReportScheduleConfigStr, mon.APICompleted, mon.TenantAPICompleted))
		}

		converted := swagmodels.JSONAPITenantReportScheduleConfig{}
		err = convertToJsonapiObject(result, &converted)
		if err != nil {
			return tenant_provisioning_service.NewCreateReportScheduleConfigInternalServerError().WithPayload(reportAPIError(fmt.Sprintf("Unable to convert %s data to jsonapi return format: %s", metmod.ReportScheduleConfigStr, err.Error()), startTime, http.StatusInternalServerError, mon.CreateReportScheduleConfigStr, mon.APICompleted, mon.TenantAPICompleted))
		}

		reportAPICompletionState(startTime, http.StatusOK, mon.CreateReportScheduleConfigStr, mon.APICompleted, mon.TenantAPICompleted)
		logger.Log.Infof("Created %s %s", metmod.ReportScheduleConfigStr, models.AsJSONString(converted))
		return tenant_provisioning_service.NewCreateReportScheduleConfigOK().WithPayload(&converted)
	}
}

// HandleUpdateReportScheduleConfig - updates a report schedule config for a tenant
func HandleUpdateReportScheduleConfig(allowedRoles []string, tenantDB datastore.TenantServiceDatastore) func(params tenant_provisioning_service.UpdateReportScheduleConfigParams) middleware.Responder {
	return func(params tenant_provisioning_service.UpdateReportScheduleConfigParams) middleware.Responder {
		startTime := time.Now()
		incrementAPICounters(mon.APIRecieved, mon.TenantAPIRecieved)
		logger.Log.Infof("Updating %s %s for Tenant %s", metmod.ReportScheduleConfigStr, params.Body.Data.Attributes.Name, params.TenantID)

		if !isRequestAuthorized(params.HTTPRequest, allowedRoles) {
			return tenant_provisioning_service.NewUpdateReportScheduleConfigForbidden().WithPayload(reportAPIError(fmt.Sprintf("Update %s operation not authorized for role: %s", metmod.ReportScheduleConfigStr, params.HTTPRequest.Header.Get(XFwdUserRoles)), startTime, http.StatusForbidden, mon.UpdateReportScheduleConfigStr, mon.APICompleted, mon.TenantAPICompleted))
		}

		// Ensure that we can unmarshal the provided report schedule payload into the model object
		requestBytes, err := json.Marshal(params.Body)
		if err != nil {
			return tenant_provisioning_service.NewUpdateReportScheduleConfigBadRequest().WithPayload(reportAPIError(generateErrorMessage(http.StatusBadRequest, err.Error()), startTime, http.StatusBadRequest, mon.UpdateReportScheduleConfigStr, mon.APICompleted, mon.TenantAPICompleted))
		}
		data := metmod.ReportScheduleConfig{}
		err = jsonapi.Unmarshal(requestBytes, &data)
		if err != nil {
			return tenant_provisioning_service.NewUpdateReportScheduleConfigBadRequest().WithPayload(reportAPIError(generateErrorMessage(http.StatusBadRequest, err.Error()), startTime, http.StatusBadRequest, mon.UpdateReportScheduleConfigStr, mon.APICompleted, mon.TenantAPICompleted))
		}
		// Seconds is not really used except for testing.
		data.Second = "0"

		// Ensure that the passed in data adheres to the model requirements
		err = data.Validate(true)
		if err != nil {
			return tenant_provisioning_service.NewUpdateReportScheduleConfigBadRequest().WithPayload(reportAPIError(generateErrorMessage(http.StatusBadRequest, err.Error()), startTime, http.StatusBadRequest, mon.UpdateReportScheduleConfigStr, mon.APICompleted, mon.TenantAPICompleted))
		}

		// Attempt to create a config entry in the datastore for the scheduler to pick up
		result, err := tenantDB.UpdateReportScheduleConfig(&data)
		if err != nil {
			return tenant_provisioning_service.NewUpdateReportScheduleConfigInternalServerError().WithPayload(reportAPIError(fmt.Sprintf("Unable to store %s: %s", metmod.ReportScheduleConfigStr, err.Error()), startTime, http.StatusInternalServerError, mon.UpdateReportScheduleConfigStr, mon.APICompleted, mon.TenantAPICompleted))
		}

		// Tell the scheduler to go and update based on the updated database
		err = scheduler.RebuildCronJobs()
		if err != nil {
			return tenant_provisioning_service.NewUpdateReportScheduleConfigInternalServerError().WithPayload(reportAPIError(fmt.Sprintf("Unable to start scheduled job%s: %s", metmod.ReportScheduleConfigStr, err.Error()), startTime, http.StatusInternalServerError, mon.UpdateReportScheduleConfigStr, mon.APICompleted, mon.TenantAPICompleted))
		}

		converted := swagmodels.JSONAPITenantReportScheduleConfig{}
		err = convertToJsonapiObject(result, &converted)
		if err != nil {
			return tenant_provisioning_service.NewUpdateReportScheduleConfigInternalServerError().WithPayload(reportAPIError(fmt.Sprintf("Unable to convert %s data to jsonapi return format: %s", metmod.ReportScheduleConfigStr, err.Error()), startTime, http.StatusInternalServerError, mon.UpdateReportScheduleConfigStr, mon.APICompleted, mon.TenantAPICompleted))
		}

		reportAPICompletionState(startTime, http.StatusOK, mon.UpdateReportScheduleConfigStr, mon.APICompleted, mon.TenantAPICompleted)
		logger.Log.Infof("Updated %s %s", metmod.ReportScheduleConfigStr, models.AsJSONString(converted))
		return tenant_provisioning_service.NewUpdateReportScheduleConfigOK().WithPayload(&converted)
	}
}

// HandleGetReportScheduleConfig - fetch a report schedule config for a tenant
func HandleGetReportScheduleConfig(allowedRoles []string, tenantDB datastore.TenantServiceDatastore) func(params tenant_provisioning_service.GetReportScheduleConfigParams) middleware.Responder {
	return func(params tenant_provisioning_service.GetReportScheduleConfigParams) middleware.Responder {
		startTime := time.Now()
		incrementAPICounters(mon.APIRecieved, mon.TenantAPIRecieved)
		logger.Log.Infof("Fetching %s %s for Tenant %s", metmod.ReportScheduleConfigStr, params.ConfigID, params.TenantID)

		if !isRequestAuthorized(params.HTTPRequest, allowedRoles) {
			return tenant_provisioning_service.NewGetReportScheduleConfigForbidden().WithPayload(reportAPIError(fmt.Sprintf("Get %s operation not authorized for role: %s", metmod.ReportScheduleConfigStr, params.HTTPRequest.Header.Get(XFwdUserRoles)), startTime, http.StatusForbidden, mon.GetReportScheduleConfigStr, mon.APICompleted, mon.TenantAPICompleted))
		}

		// Issue request to DAO Layer
		result, err := tenantDB.GetReportScheduleConfig(params.TenantID, params.ConfigID)
		if err != nil {
			return tenant_provisioning_service.NewGetReportScheduleConfigInternalServerError().WithPayload(reportAPIError(fmt.Sprintf("Unable to retrieve %s: %s", metmod.ReportScheduleConfigStr, err.Error()), startTime, http.StatusInternalServerError, mon.GetReportScheduleConfigStr, mon.APICompleted, mon.TenantAPICompleted))
		}

		converted := swagmodels.JSONAPITenantReportScheduleConfig{}
		err = convertToJsonapiObject(result, &converted)
		if err != nil {
			return tenant_provisioning_service.NewGetReportScheduleConfigInternalServerError().WithPayload(reportAPIError(fmt.Sprintf("Unable to convert %s data to jsonapi return format: %s", metmod.ReportScheduleConfigStr, err.Error()), startTime, http.StatusInternalServerError, mon.GetReportScheduleConfigStr, mon.APICompleted, mon.TenantAPICompleted))
		}

		reportAPICompletionState(startTime, http.StatusOK, mon.GetReportScheduleConfigStr, mon.APICompleted, mon.TenantAPICompleted)
		logger.Log.Infof("Retrieved %s %s", metmod.ReportScheduleConfigStr, models.AsJSONString(converted))
		return tenant_provisioning_service.NewGetReportScheduleConfigOK().WithPayload(&converted)
	}
}

// HandleGetAllReportScheduleConfigs - fetch all domains for a tenant
func HandleGetAllReportScheduleConfigs(allowedRoles []string, tenantDB datastore.TenantServiceDatastore) func(params tenant_provisioning_service.GetAllReportScheduleConfigParams) middleware.Responder {
	return func(params tenant_provisioning_service.GetAllReportScheduleConfigParams) middleware.Responder {
		startTime := time.Now()
		incrementAPICounters(mon.APIRecieved, mon.TenantAPIRecieved)
		logger.Log.Infof("Fetching all %ss for Tenant %s", metmod.ReportScheduleConfigStr, params.TenantID)

		if !isRequestAuthorized(params.HTTPRequest, allowedRoles) {
			return tenant_provisioning_service.NewGetAllReportScheduleConfigForbidden().WithPayload(reportAPIError(fmt.Sprintf("Get all %ss operation not authorized for role: %s", metmod.ReportScheduleConfigStr, params.HTTPRequest.Header.Get(XFwdUserRoles)), startTime, http.StatusForbidden, mon.GetAllReportScheduleConfigStr, mon.APICompleted, mon.TenantAPICompleted))
		}

		// Issue request to DAO Layer
		result, err := tenantDB.GetAllReportScheduleConfigs(params.TenantID)
		if err != nil {
			return tenant_provisioning_service.NewGetAllReportScheduleConfigInternalServerError().WithPayload(reportAPIError(fmt.Sprintf("Unable to retrieve %s list: %s", metmod.ReportScheduleConfigStr, err.Error()), startTime, http.StatusInternalServerError, mon.GetAllReportScheduleConfigStr, mon.APICompleted, mon.TenantAPICompleted))
		}

		converted := swagmodels.JSONAPITenantReportScheduleConfigList{}
		err = convertToJsonapiObject(result, &converted)
		if err != nil {
			return tenant_provisioning_service.NewGetAllReportScheduleConfigInternalServerError().WithPayload(reportAPIError(fmt.Sprintf("Unable to convert %s list data to jsonapi return format: %s", metmod.ReportScheduleConfigStr, err.Error()), startTime, http.StatusInternalServerError, mon.GetAllReportScheduleConfigStr, mon.APICompleted, mon.TenantAPICompleted))
		}

		reportAPICompletionState(startTime, http.StatusOK, mon.GetAllReportScheduleConfigStr, mon.APICompleted, mon.TenantAPICompleted)
		logger.Log.Infof("Retrieved %d %ss", len(result), metmod.ReportScheduleConfigStr)
		return tenant_provisioning_service.NewGetAllReportScheduleConfigOK().WithPayload(&converted)
	}
}

// HandleDeleteReportScheduleConfig - delete a domain for a tenant
func HandleDeleteReportScheduleConfig(allowedRoles []string, tenantDB datastore.TenantServiceDatastore) func(params tenant_provisioning_service.DeleteReportScheduleConfigParams) middleware.Responder {
	return func(params tenant_provisioning_service.DeleteReportScheduleConfigParams) middleware.Responder {
		startTime := time.Now()
		incrementAPICounters(mon.APIRecieved, mon.TenantAPIRecieved)
		logger.Log.Infof("Deleting %s %s for Tenant %s", metmod.ReportScheduleConfigStr, params.ConfigID, params.TenantID)

		if !isRequestAuthorized(params.HTTPRequest, allowedRoles) {
			return tenant_provisioning_service.NewDeleteReportScheduleConfigForbidden().WithPayload(reportAPIError(fmt.Sprintf("Get %s operation not authorized for role: %s", metmod.ReportScheduleConfigStr, params.HTTPRequest.Header.Get(XFwdUserRoles)), startTime, http.StatusForbidden, mon.DeleteReportScheduleConfigStr, mon.APICompleted, mon.TenantAPICompleted))
		}

		// Issue request to DAO Layer
		result, err := tenantDB.DeleteReportScheduleConfig(params.TenantID, params.ConfigID)
		if err != nil {
			return tenant_provisioning_service.NewDeleteReportScheduleConfigInternalServerError().WithPayload(reportAPIError(fmt.Sprintf("Unable to retrieve %s: %s", metmod.ReportScheduleConfigStr, err.Error()), startTime, http.StatusInternalServerError, mon.DeleteReportScheduleConfigStr, mon.APICompleted, mon.TenantAPICompleted))
		}

		// Tell the scheduler to go and update based on the updated database
		err = scheduler.RebuildCronJobs()
		if err != nil {
			return tenant_provisioning_service.NewDeleteReportScheduleConfigInternalServerError().WithPayload(reportAPIError(fmt.Sprintf("Unable to start scheduled job%s: %s", metmod.ReportScheduleConfigStr, err.Error()), startTime, http.StatusInternalServerError, mon.DeleteReportScheduleConfigStr, mon.APICompleted, mon.TenantAPICompleted))
		}

		converted := swagmodels.JSONAPITenantReportScheduleConfig{}
		err = convertToJsonapiObject(result, &converted)
		if err != nil {
			return tenant_provisioning_service.NewDeleteReportScheduleConfigInternalServerError().WithPayload(reportAPIError(fmt.Sprintf("Unable to convert %s data to jsonapi return format: %s", metmod.ReportScheduleConfigStr, err.Error()), startTime, http.StatusInternalServerError, mon.DeleteReportScheduleConfigStr, mon.APICompleted, mon.TenantAPICompleted))
		}

		reportAPICompletionState(startTime, http.StatusOK, mon.DeleteReportScheduleConfigStr, mon.APICompleted, mon.TenantAPICompleted)
		logger.Log.Infof("Deleted %s %s", metmod.ReportScheduleConfigStr, models.AsJSONString(converted))
		return tenant_provisioning_service.NewDeleteReportScheduleConfigOK().WithPayload(&converted)
	}
}
