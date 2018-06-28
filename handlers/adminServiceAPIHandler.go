package handlers

import (
	"fmt"
	"net/http"
	"time"

	"github.com/accedian/adh-gather/models"
	"github.com/accedian/adh-gather/swagmodels"

	"github.com/accedian/adh-gather/datastore"
	"github.com/accedian/adh-gather/logger"
	admmod "github.com/accedian/adh-gather/models/admin"
	mon "github.com/accedian/adh-gather/monitoring"
	"github.com/accedian/adh-gather/restapi/operations/admin_provisioning_service"
	"github.com/go-openapi/runtime/middleware"
)

// HandleGetTenant - retrieve a tenant by the tenant ID.
func HandleGetTenant(allowedRoles []string, adminDB datastore.AdminServiceDatastore) func(params admin_provisioning_service.GetTenantParams) middleware.Responder {
	return func(params admin_provisioning_service.GetTenantParams) middleware.Responder {
		startTime := time.Now()

		logger.Log.Infof("Fetching %s: %s", admmod.TenantStr, params.TenantID)

		if !isRequestAuthorized(params.HTTPRequest, allowedRoles) {
			return admin_provisioning_service.NewGetTenantForbidden().WithPayload(reportAPIError(fmt.Sprintf("Get %s operation not authorized for role: %s", admmod.TenantStr, params.HTTPRequest.Header.Get(xFwdUserRoles)), startTime, http.StatusForbidden, mon.GetTenantStr, mon.APICompleted, mon.AdminAPICompleted))
		}

		// Issue request to DAO Layer
		result, err := adminDB.GetTenantDescriptor(params.TenantID)
		if err != nil {
			return admin_provisioning_service.NewGetTenantInternalServerError().WithPayload(reportAPIError(fmt.Sprintf("Unable to retrieve %s: %s", admmod.TenantStr, err.Error()), startTime, http.StatusInternalServerError, mon.GetTenantStr, mon.APICompleted, mon.AdminAPICompleted))
		}

		converted := swagmodels.JSONAPITenant{}
		err = convertToJsonapiObject(result, &converted)
		if err != nil {
			return admin_provisioning_service.NewGetTenantInternalServerError().WithPayload(reportAPIError(fmt.Sprintf("Unable to convert %s data to jsonapi return format: %s", admmod.TenantStr, err.Error()), startTime, http.StatusInternalServerError, mon.GetTenantStr, mon.APICompleted, mon.AdminAPICompleted))
		}

		reportAPICompletionState(startTime, http.StatusOK, mon.GetTenantStr, mon.APICompleted, mon.AdminAPICompleted)
		logger.Log.Infof("Retrieved %s %s", admmod.TenantStr, models.AsJSONString(converted))
		return admin_provisioning_service.NewGetTenantOK().WithPayload(&converted)
	}
}
