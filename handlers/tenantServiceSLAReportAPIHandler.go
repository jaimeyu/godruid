package handlers

import (
	"fmt"
	"net/http"
	"time"

	"github.com/accedian/adh-gather/datastore"
	"github.com/accedian/adh-gather/logger"
	"github.com/accedian/adh-gather/models"
	metmod "github.com/accedian/adh-gather/models/metrics"
	mon "github.com/accedian/adh-gather/monitoring"
	"github.com/accedian/adh-gather/restapi/operations/tenant_provisioning_service"
	"github.com/accedian/adh-gather/swagmodels"
	"github.com/go-openapi/runtime/middleware"
)

// HandleGetSLAReport - fetch a SLA report for a tenant
func HandleGetSLAReport(allowedRoles []string, tenantDB datastore.TenantServiceDatastore) func(params tenant_provisioning_service.GetSLAReportParams) middleware.Responder {
	return func(params tenant_provisioning_service.GetSLAReportParams) middleware.Responder {
		startTime := time.Now()
		incrementAPICounters(mon.APIRecieved, mon.TenantAPIRecieved)
		logger.Log.Infof("Fetching %s %s for Tenant %s", metmod.ReportStr, params.ReportID, params.TenantID)

		if !isRequestAuthorized(params.HTTPRequest, allowedRoles) {
			return tenant_provisioning_service.NewGetSLAReportForbidden().WithPayload(reportAPIError(fmt.Sprintf("Get %s operation not authorized for role: %s", metmod.ReportStr, params.HTTPRequest.Header.Get(xFwdUserRoles)), startTime, http.StatusForbidden, mon.GetSLAReportStr, mon.APICompleted, mon.TenantAPICompleted))
		}

		// Issue request to DAO Layer
		result, err := tenantDB.GetSLAReport(params.TenantID, params.ReportID)
		if err != nil {
			return tenant_provisioning_service.NewGetSLAReportInternalServerError().WithPayload(reportAPIError(fmt.Sprintf("Unable to retrieve %s: %s", metmod.ReportStr, err.Error()), startTime, http.StatusInternalServerError, mon.GetSLAReportStr, mon.APICompleted, mon.TenantAPICompleted))
		}

		converted := swagmodels.GathergrpcJSONAPIObject{}
		err = convertToJsonapiObject(result, &converted)
		if err != nil {
			return tenant_provisioning_service.NewGetSLAReportInternalServerError().WithPayload(reportAPIError(fmt.Sprintf("Unable to convert %s data to jsonapi return format: %s", metmod.ReportStr, err.Error()), startTime, http.StatusInternalServerError, mon.GetSLAReportStr, mon.APICompleted, mon.TenantAPICompleted))
		}

		reportAPICompletionState(startTime, http.StatusOK, mon.GetSLAReportStr, mon.APICompleted, mon.TenantAPICompleted)
		logger.Log.Infof("Retrieved %s %s", metmod.ReportStr, models.AsJSONString(converted))
		return tenant_provisioning_service.NewGetSLAReportOK().WithPayload(&converted)
	}
}

// HandleGetAllSLAReports - fetch all SLA reports for a tenant
func HandleGetAllSLAReports(allowedRoles []string, tenantDB datastore.TenantServiceDatastore) func(params tenant_provisioning_service.GetAllSLAReportsParams) middleware.Responder {
	return func(params tenant_provisioning_service.GetAllSLAReportsParams) middleware.Responder {
		startTime := time.Now()
		incrementAPICounters(mon.APIRecieved, mon.TenantAPIRecieved)
		logger.Log.Infof("Fetching all %ss for Tenant %s", metmod.ReportStr, params.TenantID)

		if !isRequestAuthorized(params.HTTPRequest, allowedRoles) {
			return tenant_provisioning_service.NewGetAllSLAReportsForbidden().WithPayload(reportAPIError(fmt.Sprintf("Get all %ss operation not authorized for role: %s", metmod.ReportStr, params.HTTPRequest.Header.Get(xFwdUserRoles)), startTime, http.StatusForbidden, mon.GetAllSLAReportStr, mon.APICompleted, mon.TenantAPICompleted))
		}

		// Issue request to DAO Layer
		result, err := tenantDB.GetAllSLAReports(params.TenantID)
		if err != nil {
			return tenant_provisioning_service.NewGetAllSLAReportsInternalServerError().WithPayload(reportAPIError(fmt.Sprintf("Unable to retrieve %s list: %s", metmod.ReportStr, err.Error()), startTime, http.StatusInternalServerError, mon.GetAllSLAReportStr, mon.APICompleted, mon.TenantAPICompleted))
		}

		converted := swagmodels.GathergrpcJSONAPIObjectList{}
		err = convertToJsonapiObject(result, &converted)
		if err != nil {
			return tenant_provisioning_service.NewGetAllSLAReportsInternalServerError().WithPayload(reportAPIError(fmt.Sprintf("Unable to convert %s list data to jsonapi return format: %s", metmod.ReportStr, err.Error()), startTime, http.StatusInternalServerError, mon.GetAllSLAReportStr, mon.APICompleted, mon.TenantAPICompleted))
		}

		reportAPICompletionState(startTime, http.StatusOK, mon.GetAllSLAReportStr, mon.APICompleted, mon.TenantAPICompleted)
		logger.Log.Infof("Retrieved %d %ss", len(result), metmod.ReportStr)
		return tenant_provisioning_service.NewGetAllSLAReportsOK().WithPayload(&converted)
	}
}
