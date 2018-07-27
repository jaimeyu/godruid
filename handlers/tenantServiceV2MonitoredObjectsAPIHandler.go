package handlers

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/accedian/adh-gather/datastore"
	"github.com/accedian/adh-gather/logger"
	tenmod "github.com/accedian/adh-gather/models/tenant"
	mon "github.com/accedian/adh-gather/monitoring"
	"github.com/accedian/adh-gather/restapi/operations/tenant_provisioning_service_v2"
	"github.com/accedian/adh-gather/swagmodels"
	"github.com/go-openapi/runtime/middleware"
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

func doGetAllMonitoredObjectsV2(allowedRoles []string, tenantDB datastore.TenantServiceDatastore, params tenant_provisioning_service_v2.GetAllMonitoredObjectsV2Params) (time.Time, int, *swagmodels.MonitoredObjectList, error) {
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

	converted := swagmodels.MonitoredObjectList{}
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
