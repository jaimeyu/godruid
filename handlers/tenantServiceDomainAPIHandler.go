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

// HandleCreateTenantDomain - creates a domain for a tenant
func HandleCreateTenantDomain(allowedRoles []string, tenantDB datastore.TenantServiceDatastore) func(params tenant_provisioning_service.CreateTenantDomainParams) middleware.Responder {
	return func(params tenant_provisioning_service.CreateTenantDomainParams) middleware.Responder {
		startTime := time.Now()
		incrementAPICounters(mon.APIRecieved, mon.TenantAPIRecieved)
		logger.Log.Infof("Creating %s %s for Tenant %s", tenmod.TenantDomainStr, params.Body.Data.Attributes.Name, params.TenantID)

		if !isRequestAuthorized(params.HTTPRequest, allowedRoles) {
			return tenant_provisioning_service.NewCreateTenantDomainForbidden().WithPayload(reportAPIError(fmt.Sprintf("Create %s operation not authorized for role: %s", tenmod.TenantDomainStr, params.HTTPRequest.Header.Get(XFwdUserRoles)), startTime, http.StatusForbidden, mon.CreateTenantDomainStr, mon.APICompleted, mon.TenantAPICompleted))
		}

		// Unmarshal the request
		requestBytes, err := json.Marshal(params.Body)
		if err != nil {
			return tenant_provisioning_service.NewCreateTenantDomainBadRequest().WithPayload(reportAPIError(generateErrorMessage(http.StatusBadRequest, err.Error()), startTime, http.StatusBadRequest, mon.CreateTenantDomainStr, mon.APICompleted, mon.TenantAPICompleted))
		}

		data := tenmod.Domain{}
		err = jsonapi.Unmarshal(requestBytes, &data)
		if err != nil {
			return tenant_provisioning_service.NewCreateTenantDomainBadRequest().WithPayload(reportAPIError(generateErrorMessage(http.StatusBadRequest, err.Error()), startTime, http.StatusBadRequest, mon.CreateTenantDomainStr, mon.APICompleted, mon.TenantAPICompleted))
		}

		err = data.Validate(false)
		if err != nil {
			return tenant_provisioning_service.NewCreateTenantDomainBadRequest().WithPayload(reportAPIError(generateErrorMessage(http.StatusBadRequest, err.Error()), startTime, http.StatusBadRequest, mon.CreateTenantDomainStr, mon.APICompleted, mon.TenantAPICompleted))
		}

		// Issue request to DAO Layer
		result, err := tenantDB.CreateTenantDomain(&data)
		if err != nil {
			return tenant_provisioning_service.NewCreateTenantDomainInternalServerError().WithPayload(reportAPIError(fmt.Sprintf("Unable to store %s: %s", tenmod.TenantDomainStr, err.Error()), startTime, http.StatusInternalServerError, mon.CreateTenantDomainStr, mon.APICompleted, mon.TenantAPICompleted))
		}

		if changeNotificationEnabled {
			NotifyDomainCreated(data.TenantID, &data)
		}

		converted := swagmodels.JSONAPITenantDomain{}
		err = convertToJsonapiObject(result, &converted)
		if err != nil {
			return tenant_provisioning_service.NewCreateTenantDomainInternalServerError().WithPayload(reportAPIError(fmt.Sprintf("Unable to convert %s data to jsonapi return format: %s", tenmod.TenantDomainStr, err.Error()), startTime, http.StatusInternalServerError, mon.CreateTenantDomainStr, mon.APICompleted, mon.TenantAPICompleted))
		}

		reportAPICompletionState(startTime, http.StatusOK, mon.CreateTenantDomainStr, mon.APICompleted, mon.TenantAPICompleted)
		logger.Log.Infof("Created %s %s", tenmod.TenantDomainStr, models.AsJSONString(converted))
		return tenant_provisioning_service.NewCreateTenantDomainOK().WithPayload(&converted)
	}
}

// HandleUpdateTenantDomain - updates a domain for a tenant
func HandleUpdateTenantDomain(allowedRoles []string, tenantDB datastore.TenantServiceDatastore) func(params tenant_provisioning_service.UpdateTenantDomainParams) middleware.Responder {
	return func(params tenant_provisioning_service.UpdateTenantDomainParams) middleware.Responder {
		startTime := time.Now()
		incrementAPICounters(mon.APIRecieved, mon.TenantAPIRecieved)
		logger.Log.Infof("Updating %s %s for Tenant", tenmod.TenantDomainStr, params.Body.Data.Attributes.Name, params.TenantID)

		if !isRequestAuthorized(params.HTTPRequest, allowedRoles) {
			return tenant_provisioning_service.NewUpdateTenantDomainForbidden().WithPayload(reportAPIError(fmt.Sprintf("Update %s operation not authorized for role: %s", tenmod.TenantDomainStr, params.HTTPRequest.Header.Get(XFwdUserRoles)), startTime, http.StatusForbidden, mon.UpdateTenantDomainStr, mon.APICompleted, mon.TenantAPICompleted))
		}

		// Unmarshal the request
		requestBytes, err := json.Marshal(params.Body)
		if err != nil {
			return tenant_provisioning_service.NewUpdateTenantDomainBadRequest().WithPayload(reportAPIError(generateErrorMessage(http.StatusBadRequest, err.Error()), startTime, http.StatusBadRequest, mon.UpdateTenantDomainStr, mon.APICompleted, mon.TenantAPICompleted))
		}

		data := tenmod.Domain{}
		err = jsonapi.Unmarshal(requestBytes, &data)
		if err != nil {
			return tenant_provisioning_service.NewUpdateTenantDomainBadRequest().WithPayload(reportAPIError(generateErrorMessage(http.StatusBadRequest, err.Error()), startTime, http.StatusBadRequest, mon.UpdateTenantDomainStr, mon.APICompleted, mon.TenantAPICompleted))
		}

		err = data.Validate(true)
		if err != nil {
			return tenant_provisioning_service.NewUpdateTenantDomainBadRequest().WithPayload(reportAPIError(generateErrorMessage(http.StatusBadRequest, err.Error()), startTime, http.StatusBadRequest, mon.UpdateTenantDomainStr, mon.APICompleted, mon.TenantAPICompleted))
		}

		// Issue request to DAO Layer
		result, err := tenantDB.UpdateTenantDomain(&data)
		if err != nil {
			return tenant_provisioning_service.NewUpdateTenantDomainInternalServerError().WithPayload(reportAPIError(fmt.Sprintf("Unable to store %s: %s", tenmod.TenantDomainStr, err.Error()), startTime, http.StatusInternalServerError, mon.UpdateTenantDomainStr, mon.APICompleted, mon.TenantAPICompleted))
		}

		if changeNotificationEnabled {
			NotifyDomainUpdated(data.TenantID, &data)
		}

		converted := swagmodels.JSONAPITenantDomain{}
		err = convertToJsonapiObject(result, &converted)
		if err != nil {
			return tenant_provisioning_service.NewUpdateTenantDomainInternalServerError().WithPayload(reportAPIError(fmt.Sprintf("Unable to convert %s data to jsonapi return format: %s", tenmod.TenantDomainStr, err.Error()), startTime, http.StatusInternalServerError, mon.UpdateTenantDomainStr, mon.APICompleted, mon.TenantAPICompleted))
		}

		reportAPICompletionState(startTime, http.StatusOK, mon.UpdateTenantDomainStr, mon.APICompleted, mon.TenantAPICompleted)
		logger.Log.Infof("Updated %s %s", tenmod.TenantDomainStr, models.AsJSONString(converted))
		return tenant_provisioning_service.NewUpdateTenantDomainOK().WithPayload(&converted)
	}
}

// HandlePatchTenantDomain - patches a domain for a tenant
func HandlePatchTenantDomain(allowedRoles []string, tenantDB datastore.TenantServiceDatastore) func(params tenant_provisioning_service.PatchTenantDomainParams) middleware.Responder {
	return func(params tenant_provisioning_service.PatchTenantDomainParams) middleware.Responder {
		startTime := time.Now()
		incrementAPICounters(mon.APIRecieved, mon.TenantAPIRecieved)
		logger.Log.Infof("Patching %s %s for Tenant", tenmod.TenantDomainStr, params.Body.Data.ID, params.TenantID)

		if !isRequestAuthorized(params.HTTPRequest, allowedRoles) {
			return tenant_provisioning_service.NewPatchTenantDomainForbidden().WithPayload(reportAPIError(fmt.Sprintf("Patch %s operation not authorized for role: %s", tenmod.TenantDomainStr, params.HTTPRequest.Header.Get(XFwdUserRoles)), startTime, http.StatusForbidden, mon.PatchTenantDomainStr, mon.APICompleted, mon.TenantAPICompleted))
		}

		// Unmarshal the request
		requestBytes, err := json.Marshal(params.Body)
		if err != nil {
			return tenant_provisioning_service.NewPatchTenantDomainBadRequest().WithPayload(reportAPIError(generateErrorMessage(http.StatusBadRequest, err.Error()), startTime, http.StatusBadRequest, mon.PatchTenantDomainStr, mon.APICompleted, mon.TenantAPICompleted))
		}

		// Model
		data := tenmod.Domain{}
		err = jsonapi.Unmarshal(requestBytes, &data)
		if err != nil {
			return tenant_provisioning_service.NewPatchTenantDomainBadRequest().WithPayload(reportAPIError(generateErrorMessage(http.StatusBadRequest, err.Error()), startTime, http.StatusBadRequest, mon.PatchTenantDomainStr, mon.APICompleted, mon.TenantAPICompleted))
		}

		// This only checks if the ID&REV is set.
		err = data.Validate(true)
		if err != nil {
			return tenant_provisioning_service.NewPatchTenantDomainBadRequest().WithPayload(reportAPIError(generateErrorMessage(http.StatusBadRequest, err.Error()), startTime, http.StatusBadRequest, mon.PatchTenantDomainStr, mon.APICompleted, mon.TenantAPICompleted))
		}

		// Issue request to DAO Layer
		oldDomain, err := tenantDB.GetTenantDomain(data.TenantID, data.ID)
		if err != nil {
			return tenant_provisioning_service.NewPatchTenantDomainInternalServerError().WithPayload(reportAPIError(fmt.Sprintf("Unable to retrieve %s: %s", tenmod.TenantDomainStr, err.Error()), startTime, http.StatusInternalServerError, mon.PatchTenantDomainStr, mon.APICompleted, mon.TenantAPICompleted))
		}

		errMerge := models.MergeObjWithMap(oldDomain, requestBytes)
		if errMerge != nil {
			return tenant_provisioning_service.NewPatchTenantDomainBadRequest().WithPayload(reportAPIError(generateErrorMessage(http.StatusBadRequest, err.Error()), startTime, http.StatusBadRequest, mon.PatchTenantDomainStr, mon.APICompleted, mon.TenantAPICompleted))
		}

		// This only checks if the ID&REV is set.
		err = oldDomain.Validate(true)
		if err != nil {
			return tenant_provisioning_service.NewPatchTenantDomainBadRequest().WithPayload(reportAPIError(generateErrorMessage(http.StatusBadRequest, err.Error()), startTime, http.StatusBadRequest, mon.PatchTenantDomainStr, mon.APICompleted, mon.TenantAPICompleted))
		}

		// Issue request to DAO Layer
		result, err := tenantDB.UpdateTenantDomain(oldDomain)
		if err != nil {
			return tenant_provisioning_service.NewPatchTenantDomainInternalServerError().WithPayload(reportAPIError(fmt.Sprintf("Unable to store %s: %s", tenmod.TenantDomainStr, err.Error()), startTime, http.StatusInternalServerError, mon.PatchTenantDomainStr, mon.APICompleted, mon.TenantAPICompleted))
		}

		if changeNotificationEnabled {
			NotifyDomainUpdated(oldDomain.TenantID, oldDomain)
		}

		converted := swagmodels.JSONAPITenantDomain{}
		err = convertToJsonapiObject(result, &converted)
		if err != nil {
			return tenant_provisioning_service.NewPatchTenantDomainInternalServerError().WithPayload(reportAPIError(fmt.Sprintf("Unable to convert %s data to jsonapi return format: %s", tenmod.TenantDomainStr, err.Error()), startTime, http.StatusInternalServerError, mon.PatchTenantDomainStr, mon.APICompleted, mon.TenantAPICompleted))
		}

		reportAPICompletionState(startTime, http.StatusOK, mon.PatchTenantDomainStr, mon.APICompleted, mon.TenantAPICompleted)
		logger.Log.Infof("Updated %s %s", tenmod.TenantDomainStr, models.AsJSONString(converted))
		return tenant_provisioning_service.NewPatchTenantDomainOK().WithPayload(&converted)
	}
}

// HandleGetTenantDomain - fetch a domain for a tenant
func HandleGetTenantDomain(allowedRoles []string, tenantDB datastore.TenantServiceDatastore) func(params tenant_provisioning_service.GetTenantDomainParams) middleware.Responder {
	return func(params tenant_provisioning_service.GetTenantDomainParams) middleware.Responder {
		startTime := time.Now()
		incrementAPICounters(mon.APIRecieved, mon.TenantAPIRecieved)
		logger.Log.Infof("Fetching %s %s for Tenant %s", tenmod.TenantDomainStr, params.DomainID, params.TenantID)

		if !isRequestAuthorized(params.HTTPRequest, allowedRoles) {
			return tenant_provisioning_service.NewGetTenantDomainForbidden().WithPayload(reportAPIError(fmt.Sprintf("Get %s operation not authorized for role: %s", tenmod.TenantDomainStr, params.HTTPRequest.Header.Get(XFwdUserRoles)), startTime, http.StatusForbidden, mon.GetTenantDomainStr, mon.APICompleted, mon.TenantAPICompleted))
		}

		// Issue request to DAO Layer
		result, err := tenantDB.GetTenantDomain(params.TenantID, params.DomainID)
		if err != nil {
			return tenant_provisioning_service.NewGetTenantDomainInternalServerError().WithPayload(reportAPIError(fmt.Sprintf("Unable to retrieve %s: %s", tenmod.TenantDomainStr, err.Error()), startTime, http.StatusInternalServerError, mon.GetTenantDomainStr, mon.APICompleted, mon.TenantAPICompleted))
		}

		converted := swagmodels.JSONAPITenantDomain{}
		err = convertToJsonapiObject(result, &converted)
		if err != nil {
			return tenant_provisioning_service.NewGetTenantDomainInternalServerError().WithPayload(reportAPIError(fmt.Sprintf("Unable to convert %s data to jsonapi return format: %s", tenmod.TenantDomainStr, err.Error()), startTime, http.StatusInternalServerError, mon.GetTenantDomainStr, mon.APICompleted, mon.TenantAPICompleted))
		}

		reportAPICompletionState(startTime, http.StatusOK, mon.GetTenantDomainStr, mon.APICompleted, mon.TenantAPICompleted)
		logger.Log.Infof("Retrieved %s %s", tenmod.TenantDomainStr, models.AsJSONString(converted))
		return tenant_provisioning_service.NewGetTenantDomainOK().WithPayload(&converted)
	}
}

// HandleGetAllTenantDomains - fetch all domains for a tenant
func HandleGetAllTenantDomains(allowedRoles []string, tenantDB datastore.TenantServiceDatastore) func(params tenant_provisioning_service.GetAllTenantDomainsParams) middleware.Responder {
	return func(params tenant_provisioning_service.GetAllTenantDomainsParams) middleware.Responder {
		startTime := time.Now()
		incrementAPICounters(mon.APIRecieved, mon.TenantAPIRecieved)
		logger.Log.Infof("Fetching all %ss for Tenant %s", tenmod.TenantDomainStr, params.TenantID)

		if !isRequestAuthorized(params.HTTPRequest, allowedRoles) {
			return tenant_provisioning_service.NewGetAllTenantDomainsForbidden().WithPayload(reportAPIError(fmt.Sprintf("Get all %ss operation not authorized for role: %s", tenmod.TenantDomainStr, params.HTTPRequest.Header.Get(XFwdUserRoles)), startTime, http.StatusForbidden, mon.GetAllTenantDomainStr, mon.APICompleted, mon.TenantAPICompleted))
		}

		// Issue request to DAO Layer
		result, err := tenantDB.GetAllTenantDomains(params.TenantID)
		if err != nil {
			return tenant_provisioning_service.NewGetAllTenantDomainsInternalServerError().WithPayload(reportAPIError(fmt.Sprintf("Unable to retrieve %s list: %s", tenmod.TenantDomainStr, err.Error()), startTime, http.StatusInternalServerError, mon.GetAllTenantDomainStr, mon.APICompleted, mon.TenantAPICompleted))
		}

		converted := swagmodels.JSONAPITenantDomainList{}
		err = convertToJsonapiObject(result, &converted)
		if err != nil {
			return tenant_provisioning_service.NewGetAllTenantDomainsInternalServerError().WithPayload(reportAPIError(fmt.Sprintf("Unable to convert %s list data to jsonapi return format: %s", tenmod.TenantDomainStr, err.Error()), startTime, http.StatusInternalServerError, mon.GetAllTenantDomainStr, mon.APICompleted, mon.TenantAPICompleted))
		}

		reportAPICompletionState(startTime, http.StatusOK, mon.GetAllTenantDomainStr, mon.APICompleted, mon.TenantAPICompleted)
		logger.Log.Infof("Retrieved %d %ss", len(result), tenmod.TenantDomainStr)
		return tenant_provisioning_service.NewGetAllTenantDomainsOK().WithPayload(&converted)
	}
}

// HandleDeleteTenantDomain - delete a domain for a tenant
func HandleDeleteTenantDomain(allowedRoles []string, tenantDB datastore.TenantServiceDatastore) func(params tenant_provisioning_service.DeleteTenantDomainParams) middleware.Responder {
	return func(params tenant_provisioning_service.DeleteTenantDomainParams) middleware.Responder {
		startTime := time.Now()
		incrementAPICounters(mon.APIRecieved, mon.TenantAPIRecieved)
		logger.Log.Infof("Deleting %s %s for Tenant %s", tenmod.TenantDomainStr, params.DomainID, params.TenantID)

		if !isRequestAuthorized(params.HTTPRequest, allowedRoles) {
			return tenant_provisioning_service.NewDeleteTenantDomainForbidden().WithPayload(reportAPIError(fmt.Sprintf("Delete %s operation not authorized for role: %s", tenmod.TenantDomainStr, params.HTTPRequest.Header.Get(XFwdUserRoles)), startTime, http.StatusForbidden, mon.DeleteTenantDomainStr, mon.APICompleted, mon.TenantAPICompleted))
		}

		// Integrity Check - Monitored Objects
		moByDomainReq := tenmod.MonitoredObjectCountByDomainRequest{
			TenantID:  params.TenantID,
			ByCount:   true,
			DomainSet: []string{params.DomainID},
		}
		moByDomainResp, err := tenantDB.GetMonitoredObjectToDomainMap(&moByDomainReq)
		if err != nil {
			return tenant_provisioning_service.NewDeleteTenantDomainInternalServerError().WithPayload(reportAPIError(fmt.Sprintf("Unable to perform integrity check for %s deletion: %s", tenmod.TenantDomainStr, err.Error()), startTime, http.StatusInternalServerError, mon.DeleteTenantDomainStr, mon.APICompleted, mon.TenantAPICompleted))
		}

		logger.Log.Infof("%s got %s", models.AsJSONString(moByDomainReq), models.AsJSONString(moByDomainResp))
		if moByDomainResp.DomainToMonitoredObjectCountMap != nil {
			if count, exists := moByDomainResp.DomainToMonitoredObjectCountMap[params.DomainID]; exists && count > 0 {
				return tenant_provisioning_service.NewDeleteTenantDomainInternalServerError().WithPayload(reportAPIError(fmt.Sprintf("%s deletion failed integrity check: in use by at least one %s", tenmod.TenantDomainStr, tenmod.TenantMonitoredObjectStr), startTime, http.StatusInternalServerError, mon.DeleteTenantDomainStr, mon.APICompleted, mon.TenantAPICompleted))
			}
		}

		// Integrity Check - Dashboards
		dashboardUsesDomain, err := tenantDB.HasDashboardsWithDomain(params.TenantID, params.DomainID)
		if err != nil {
			return tenant_provisioning_service.NewDeleteTenantDomainInternalServerError().WithPayload(reportAPIError(fmt.Sprintf("Unable to perform integrity check for %s deletion: %s", tenmod.TenantDomainStr, err.Error()), startTime, http.StatusInternalServerError, mon.DeleteTenantDomainStr, mon.APICompleted, mon.TenantAPICompleted))
		}
		if dashboardUsesDomain {
			return tenant_provisioning_service.NewDeleteTenantDomainInternalServerError().WithPayload(reportAPIError(fmt.Sprintf("%s deletion failed integrity check: in use by at least one %s", tenmod.TenantDomainStr, tenmod.TenantDashboardStr), startTime, http.StatusInternalServerError, mon.DeleteTenantDomainStr, mon.APICompleted, mon.TenantAPICompleted))
		}

		configs, err := tenantDB.GetAllReportScheduleConfigs(params.TenantID)
		if err != nil {
			return tenant_provisioning_service.NewDeleteTenantDomainInternalServerError().WithPayload(reportAPIError(fmt.Sprintf("Unable to perform integrity check for %s deletion: %s", tenmod.TenantDomainStr, err.Error()), startTime, http.StatusInternalServerError, mon.DeleteTenantDomainStr, mon.APICompleted, mon.TenantAPICompleted))
		}
		for _, rep := range configs {
			if len(rep.Domains) == 0 {
				continue
			}
			for _, dom := range rep.Domains {
				if dom == params.DomainID {
					return tenant_provisioning_service.NewDeleteTenantDomainInternalServerError().WithPayload(reportAPIError(fmt.Sprintf("%s deletion failed integrity check: in use by at least one %s", tenmod.TenantDomainStr, tenmod.TenantReportScheduleConfigStr), startTime, http.StatusInternalServerError, mon.DeleteTenantDomainStr, mon.APICompleted, mon.TenantAPICompleted))
				}
			}
		}

		// Issue request to DAO Layer
		result, err := tenantDB.DeleteTenantDomain(params.TenantID, params.DomainID)
		if err != nil {
			return tenant_provisioning_service.NewDeleteTenantDomainInternalServerError().WithPayload(reportAPIError(fmt.Sprintf("Unable to retrieve %s: %s", tenmod.TenantDomainStr, err.Error()), startTime, http.StatusInternalServerError, mon.DeleteTenantDomainStr, mon.APICompleted, mon.TenantAPICompleted))
		}

		if changeNotificationEnabled {
			NotifyDomainDeleted(params.TenantID, result)
		}

		converted := swagmodels.JSONAPITenantDomain{}
		err = convertToJsonapiObject(result, &converted)
		if err != nil {
			return tenant_provisioning_service.NewDeleteTenantDomainInternalServerError().WithPayload(reportAPIError(fmt.Sprintf("Unable to convert %s list data to jsonapi return format: %s", tenmod.TenantDomainStr, err.Error()), startTime, http.StatusInternalServerError, mon.DeleteTenantDomainStr, mon.APICompleted, mon.TenantAPICompleted))
		}

		reportAPICompletionState(startTime, http.StatusOK, mon.DeleteTenantDomainStr, mon.APICompleted, mon.TenantAPICompleted)
		logger.Log.Infof("Retrieved %s %s", tenmod.TenantDomainStr, models.AsJSONString(result))
		return tenant_provisioning_service.NewDeleteTenantDomainOK().WithPayload(&converted)
	}
}
