package handlers

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/accedian/adh-gather/models/metrics"

	"github.com/accedian/adh-gather/datastore"
	"github.com/accedian/adh-gather/logger"
	"github.com/accedian/adh-gather/models"
	admmod "github.com/accedian/adh-gather/models/admin"
	metmod "github.com/accedian/adh-gather/models/metrics"
	tenmod "github.com/accedian/adh-gather/models/tenant"
	mon "github.com/accedian/adh-gather/monitoring"
	"github.com/accedian/adh-gather/restapi/operations/tenant_provisioning_service_v2"
	"github.com/accedian/adh-gather/swagmodels"
	"github.com/go-openapi/runtime/middleware"
)

// HandleGetSLAReportV2 - retrieve an SLA Report
func HandleGetSLAReportV2(allowedRoles []string, tenantDB datastore.TenantServiceDatastore) func(params tenant_provisioning_service_v2.GetSLAReportV2Params) middleware.Responder {
	return func(params tenant_provisioning_service_v2.GetSLAReportV2Params) middleware.Responder {

		// Do the work
		startTime, responseCode, response, err := doGetSLAReportV2(allowedRoles, tenantDB, params)

		// Success Response
		if responseCode == http.StatusOK {
			return tenant_provisioning_service_v2.NewGetSLAReportV2OK().WithPayload(response)
		}

		// Error Responses
		errorMessage := reportAPIError(err.Error(), startTime, responseCode, mon.GetSLAReportStr, mon.APICompleted, mon.TenantAPICompleted)
		switch responseCode {
		case http.StatusForbidden:
			return tenant_provisioning_service_v2.NewGetSLAReportV2Forbidden().WithPayload(errorMessage)
		case http.StatusNotFound:
			return tenant_provisioning_service_v2.NewGetSLAReportV2NotFound().WithPayload(errorMessage)
		default:
			return tenant_provisioning_service_v2.NewGetSLAReportV2InternalServerError().WithPayload(errorMessage)
		}

	}
}

// HandleGetAllSLAReportsV2 - retrieve all SLA Reports for a Tenant
func HandleGetAllSLAReportsV2(allowedRoles []string, tenantDB datastore.TenantServiceDatastore) func(params tenant_provisioning_service_v2.GetAllSLAReportsV2Params) middleware.Responder {
	return func(params tenant_provisioning_service_v2.GetAllSLAReportsV2Params) middleware.Responder {

		// Do the work
		startTime, responseCode, response, err := doGetAllSLAReportsV2(allowedRoles, tenantDB, params)

		// Success Response
		if responseCode == http.StatusOK {
			return tenant_provisioning_service_v2.NewGetAllSLAReportsV2OK().WithPayload(response)
		}

		// Error Responses
		errorMessage := reportAPIError(err.Error(), startTime, responseCode, mon.GetAllSLAReportStr, mon.APICompleted, mon.TenantAPICompleted)
		switch responseCode {
		case http.StatusForbidden:
			return tenant_provisioning_service_v2.NewGetAllSLAReportsV2Forbidden().WithPayload(errorMessage)
		case http.StatusNotFound:
			return tenant_provisioning_service_v2.NewGetAllSLAReportsV2NotFound().WithPayload(errorMessage)
		default:
			return tenant_provisioning_service_v2.NewGetAllSLAReportsV2InternalServerError().WithPayload(errorMessage)
		}

	}
}

func doGetSLAReportV2(allowedRoles []string, tenantDB datastore.TenantServiceDatastore, params tenant_provisioning_service_v2.GetSLAReportV2Params) (time.Time, int, *swagmodels.GathergrpcJSONAPIObject, error) {
	tenantID := params.HTTPRequest.Header.Get(XFwdTenantId)
	isAuthorized, startTime := authorizeRequest(fmt.Sprintf("Fetching %s %s for %s %s", tenmod.TenantSLAReportStr, params.ReportID, admmod.TenantStr, tenantID), params.HTTPRequest, allowedRoles, mon.APIRecieved, mon.AdminAPIRecieved)

	if !isAuthorized {
		return startTime, http.StatusForbidden, nil, fmt.Errorf("Fetch %s operation not authorized for role: %s", tenmod.TenantReportScheduleConfigStr, params.HTTPRequest.Header.Get(XFwdUserRoles))
	}

	// Issue request to DAO Layer
	result, err := tenantDB.GetSLAReport(tenantID, params.ReportID)
	if err != nil {
		if strings.Contains(err.Error(), datastore.NotFoundStr) {
			return startTime, http.StatusNotFound, nil, err
		}
		return startTime, http.StatusInternalServerError, nil, fmt.Errorf("Unable to retrieve %s: %s", metmod.ReportStr, err.Error())
	}

	// Stick the result in an array....this is a hack for now due to improper modelling of V1 objects.
	// TODO: remove this array hack once the metrics re-work is done
	resultArray := []*metrics.SLAReport{result}

	converted := swagmodels.GathergrpcJSONAPIObject{}
	err = convertToJsonapiObject(resultArray, &converted)
	if err != nil {
		return startTime, http.StatusInternalServerError, nil, fmt.Errorf("Unable to convert %s data to jsonapi return format: %s", metmod.ReportStr, err.Error())
	}

	reportAPICompletionState(startTime, http.StatusOK, mon.GetSLAReportStr, mon.APICompleted, mon.TenantAPICompleted)
	logger.Log.Infof("Retrieved %s %s", metmod.ReportStr, models.AsJSONString(converted))
	return startTime, http.StatusOK, &converted, nil
}

func doGetAllSLAReportsV2(allowedRoles []string, tenantDB datastore.TenantServiceDatastore, params tenant_provisioning_service_v2.GetAllSLAReportsV2Params) (time.Time, int, *swagmodels.GathergrpcJSONAPIObjectList, error) {
	tenantID := params.HTTPRequest.Header.Get(XFwdTenantId)
	isAuthorized, startTime := authorizeRequest(fmt.Sprintf("Fetching %s list fot %s %s", tenmod.TenantSLAReportStr, admmod.TenantStr, tenantID), params.HTTPRequest, allowedRoles, mon.APIRecieved, mon.AdminAPIRecieved)

	if !isAuthorized {
		return startTime, http.StatusForbidden, nil, fmt.Errorf("Fetch %s operation not authorized for role: %s", tenmod.TenantReportScheduleConfigStr, params.HTTPRequest.Header.Get(XFwdUserRoles))
	}

	// Issue request to DAO Layer
	result, err := tenantDB.GetAllSLAReports(tenantID)
	if err != nil {
		return startTime, http.StatusInternalServerError, nil, fmt.Errorf("Unable to retrieve %s list: %s", metmod.ReportStr, err.Error())
	}

	if len(result) == 0 {
		return startTime, http.StatusNotFound, nil, fmt.Errorf(datastore.NotFoundStr)
	}

	converted := swagmodels.GathergrpcJSONAPIObjectList{}
	err = convertToJsonapiObject(result, &converted)
	if err != nil {
		return startTime, http.StatusInternalServerError, nil, fmt.Errorf("Unable to convert %s list data to jsonapi return format: %s", metmod.ReportStr, err.Error())
	}

	reportAPICompletionState(startTime, http.StatusOK, mon.GetAllSLAReportStr, mon.APICompleted, mon.TenantAPICompleted)
	logger.Log.Infof("Retrieved %d %ss", len(result), metmod.ReportStr)
	return startTime, http.StatusOK, &converted, nil
}
