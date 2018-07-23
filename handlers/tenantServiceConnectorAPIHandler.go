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

// HandleCreateTenantConnectorConfig - creates a connector configuration for a tenant
func HandleCreateTenantConnectorConfig(allowedRoles []string, tenantDB datastore.TenantServiceDatastore) func(params tenant_provisioning_service.CreateTenantConnectorConfigParams) middleware.Responder {
	return func(params tenant_provisioning_service.CreateTenantConnectorConfigParams) middleware.Responder {
		startTime := time.Now()
		incrementAPICounters(mon.APIRecieved, mon.TenantAPIRecieved)
		logger.Log.Infof("Creating %s %s for Tenant %s", tenmod.TenantConnectorConfigStr, params.Body.Data.Attributes.Name, params.TenantID)

		if !isRequestAuthorized(params.HTTPRequest, allowedRoles) {
			return tenant_provisioning_service.NewCreateTenantConnectorConfigForbidden().WithPayload(reportAPIError(fmt.Sprintf("Create %s operation not authorized for role: %s", tenmod.TenantConnectorConfigStr, params.HTTPRequest.Header.Get(XFwdUserRoles)), startTime, http.StatusForbidden, mon.CreateTenantConnectorConfigStr, mon.APICompleted, mon.TenantAPICompleted))
		}

		// Unmarshal the request
		requestBytes, err := json.Marshal(params.Body)
		if err != nil {
			return tenant_provisioning_service.NewCreateTenantConnectorConfigBadRequest().WithPayload(reportAPIError(generateErrorMessage(http.StatusBadRequest, err.Error()), startTime, http.StatusBadRequest, mon.CreateTenantConnectorConfigStr, mon.APICompleted, mon.TenantAPICompleted))
		}

		data := tenmod.ConnectorConfig{}
		err = jsonapi.Unmarshal(requestBytes, &data)
		if err != nil {
			return tenant_provisioning_service.NewCreateTenantConnectorConfigBadRequest().WithPayload(reportAPIError(generateErrorMessage(http.StatusBadRequest, err.Error()), startTime, http.StatusBadRequest, mon.CreateTenantConnectorConfigStr, mon.APICompleted, mon.TenantAPICompleted))
		}

		err = data.Validate(false)
		if err != nil {
			return tenant_provisioning_service.NewCreateTenantConnectorConfigBadRequest().WithPayload(reportAPIError(generateErrorMessage(http.StatusBadRequest, err.Error()), startTime, http.StatusBadRequest, mon.CreateTenantConnectorConfigStr, mon.APICompleted, mon.TenantAPICompleted))
		}

		// Issue request to DAO Layer
		result, err := tenantDB.CreateTenantConnectorConfig(&data)
		if err != nil {
			return tenant_provisioning_service.NewCreateTenantConnectorConfigInternalServerError().WithPayload(reportAPIError(fmt.Sprintf("Unable to store %s: %s", tenmod.TenantConnectorConfigStr, err.Error()), startTime, http.StatusInternalServerError, mon.CreateTenantConnectorConfigStr, mon.APICompleted, mon.TenantAPICompleted))
		}

		converted := swagmodels.JSONAPITenantConnectorConfig{}
		err = convertToJsonapiObject(result, &converted)
		if err != nil {
			return tenant_provisioning_service.NewCreateTenantConnectorConfigInternalServerError().WithPayload(reportAPIError(fmt.Sprintf("Unable to convert %s data to jsonapi return format: %s", tenmod.TenantConnectorConfigStr, err.Error()), startTime, http.StatusInternalServerError, mon.CreateTenantConnectorConfigStr, mon.APICompleted, mon.TenantAPICompleted))
		}

		reportAPICompletionState(startTime, http.StatusOK, mon.CreateTenantConnectorConfigStr, mon.APICompleted, mon.TenantAPICompleted)
		logger.Log.Infof("Created %s %s", tenmod.TenantConnectorConfigStr, models.AsJSONString(converted))
		return tenant_provisioning_service.NewCreateTenantConnectorConfigOK().WithPayload(&converted)
	}
}

// HandleUpdateTenantConnectorConfig - updates a connector configuration for a tenant
func HandleUpdateTenantConnectorConfig(allowedRoles []string, tenantDB datastore.TenantServiceDatastore) func(params tenant_provisioning_service.UpdateTenantConnectorConfigParams) middleware.Responder {
	return func(params tenant_provisioning_service.UpdateTenantConnectorConfigParams) middleware.Responder {
		startTime := time.Now()
		incrementAPICounters(mon.APIRecieved, mon.TenantAPIRecieved)
		logger.Log.Infof("Updating %s %s for Tenant %s", tenmod.TenantConnectorConfigStr, params.Body.Data.Attributes.Name, params.TenantID)

		if !isRequestAuthorized(params.HTTPRequest, allowedRoles) {
			return tenant_provisioning_service.NewUpdateTenantConnectorConfigForbidden().WithPayload(reportAPIError(fmt.Sprintf("Update %s operation not authorized for role: %s", tenmod.TenantConnectorConfigStr, params.HTTPRequest.Header.Get(XFwdUserRoles)), startTime, http.StatusForbidden, mon.UpdateTenantConnectorConfigStr, mon.APICompleted, mon.TenantAPICompleted))
		}

		// Unmarshal the request
		requestBytes, err := json.Marshal(params.Body)
		if err != nil {
			return tenant_provisioning_service.NewUpdateTenantConnectorConfigBadRequest().WithPayload(reportAPIError(generateErrorMessage(http.StatusBadRequest, err.Error()), startTime, http.StatusBadRequest, mon.UpdateTenantConnectorConfigStr, mon.APICompleted, mon.TenantAPICompleted))
		}

		data := tenmod.ConnectorConfig{}
		err = jsonapi.Unmarshal(requestBytes, &data)
		if err != nil {
			return tenant_provisioning_service.NewUpdateTenantConnectorConfigBadRequest().WithPayload(reportAPIError(generateErrorMessage(http.StatusBadRequest, err.Error()), startTime, http.StatusBadRequest, mon.UpdateTenantConnectorConfigStr, mon.APICompleted, mon.TenantAPICompleted))
		}

		err = data.Validate(true)
		if err != nil {
			return tenant_provisioning_service.NewUpdateTenantConnectorConfigBadRequest().WithPayload(reportAPIError(generateErrorMessage(http.StatusBadRequest, err.Error()), startTime, http.StatusBadRequest, mon.UpdateTenantConnectorConfigStr, mon.APICompleted, mon.TenantAPICompleted))
		}

		// Issue request to DAO Layer
		result, err := tenantDB.UpdateTenantConnectorConfig(&data)
		if err != nil {
			return tenant_provisioning_service.NewUpdateTenantConnectorConfigInternalServerError().WithPayload(reportAPIError(fmt.Sprintf("Unable to store %s: %s", tenmod.TenantConnectorConfigStr, err.Error()), startTime, http.StatusInternalServerError, mon.UpdateTenantConnectorConfigStr, mon.APICompleted, mon.TenantAPICompleted))
		}

		converted := swagmodels.JSONAPITenantConnectorConfig{}
		err = convertToJsonapiObject(result, &converted)
		if err != nil {
			return tenant_provisioning_service.NewUpdateTenantConnectorConfigInternalServerError().WithPayload(reportAPIError(fmt.Sprintf("Unable to convert %s data to jsonapi return format: %s", tenmod.TenantConnectorConfigStr, err.Error()), startTime, http.StatusInternalServerError, mon.UpdateTenantConnectorConfigStr, mon.APICompleted, mon.TenantAPICompleted))
		}

		reportAPICompletionState(startTime, http.StatusOK, mon.UpdateTenantConnectorConfigStr, mon.APICompleted, mon.TenantAPICompleted)
		logger.Log.Infof("Updated %s %s", tenmod.TenantConnectorConfigStr, models.AsJSONString(converted))
		return tenant_provisioning_service.NewUpdateTenantConnectorConfigOK().WithPayload(&converted)
	}
}

// HandleGetTenantConnectorConfig - fetch a connector configuration for a tenant
func HandleGetTenantConnectorConfig(allowedRoles []string, tenantDB datastore.TenantServiceDatastore) func(params tenant_provisioning_service.GetTenantConnectorConfigParams) middleware.Responder {
	return func(params tenant_provisioning_service.GetTenantConnectorConfigParams) middleware.Responder {
		startTime := time.Now()
		incrementAPICounters(mon.APIRecieved, mon.TenantAPIRecieved)
		logger.Log.Infof("Fetching %s %s for Tenant %s", tenmod.TenantConnectorConfigStr, params.ConnectorID, params.TenantID)

		if !isRequestAuthorized(params.HTTPRequest, allowedRoles) {
			return tenant_provisioning_service.NewGetTenantConnectorConfigForbidden().WithPayload(reportAPIError(fmt.Sprintf("Get %s operation not authorized for role: %s", tenmod.TenantConnectorConfigStr, params.HTTPRequest.Header.Get(XFwdUserRoles)), startTime, http.StatusForbidden, mon.GetTenantConnectorConfigStr, mon.APICompleted, mon.TenantAPICompleted))
		}

		// Issue request to DAO Layer
		result, err := tenantDB.GetTenantConnectorConfig(params.TenantID, params.ConnectorID)
		if err != nil {
			return tenant_provisioning_service.NewGetTenantConnectorConfigInternalServerError().WithPayload(reportAPIError(fmt.Sprintf("Unable to retrieve %s: %s", tenmod.TenantConnectorConfigStr, err.Error()), startTime, http.StatusInternalServerError, mon.GetTenantConnectorConfigStr, mon.APICompleted, mon.TenantAPICompleted))
		}

		converted := swagmodels.JSONAPITenantConnectorConfig{}
		err = convertToJsonapiObject(result, &converted)
		if err != nil {
			return tenant_provisioning_service.NewGetTenantConnectorConfigInternalServerError().WithPayload(reportAPIError(fmt.Sprintf("Unable to convert %s data to jsonapi return format: %s", tenmod.TenantConnectorConfigStr, err.Error()), startTime, http.StatusInternalServerError, mon.GetTenantConnectorConfigStr, mon.APICompleted, mon.TenantAPICompleted))
		}

		reportAPICompletionState(startTime, http.StatusOK, mon.GetTenantConnectorConfigStr, mon.APICompleted, mon.TenantAPICompleted)
		logger.Log.Infof("Retrieved %s %s", tenmod.TenantConnectorConfigStr, models.AsJSONString(converted))
		return tenant_provisioning_service.NewGetTenantConnectorConfigOK().WithPayload(&converted)
	}
}

// HandleDeleteTenantConnectorConfig - delete a connector configuration for a tenant
func HandleDeleteTenantConnectorConfig(allowedRoles []string, tenantDB datastore.TenantServiceDatastore) func(params tenant_provisioning_service.DeleteTenantConnectorConfigParams) middleware.Responder {
	return func(params tenant_provisioning_service.DeleteTenantConnectorConfigParams) middleware.Responder {
		startTime := time.Now()
		incrementAPICounters(mon.APIRecieved, mon.TenantAPIRecieved)
		logger.Log.Infof("Deleting %s %s for Tenant %s", tenmod.TenantConnectorConfigStr, params.ConnectorID, params.TenantID)

		if !isRequestAuthorized(params.HTTPRequest, allowedRoles) {
			return tenant_provisioning_service.NewGetTenantConnectorConfigForbidden().WithPayload(reportAPIError(fmt.Sprintf("Delete %s operation not authorized for role: %s", tenmod.TenantConnectorConfigStr, params.HTTPRequest.Header.Get(XFwdUserRoles)), startTime, http.StatusForbidden, mon.DeleteTenantConnectorConfigStr, mon.APICompleted, mon.TenantAPICompleted))
		}

		// Issue request to DAO Layer
		result, err := tenantDB.DeleteTenantConnectorConfig(params.TenantID, params.ConnectorID)
		if err != nil {
			return tenant_provisioning_service.NewGetTenantConnectorConfigInternalServerError().WithPayload(reportAPIError(fmt.Sprintf("Unable to delete %s: %s", tenmod.TenantConnectorConfigStr, err.Error()), startTime, http.StatusInternalServerError, mon.DeleteTenantConnectorConfigStr, mon.APICompleted, mon.TenantAPICompleted))
		}

		converted := swagmodels.JSONAPITenantConnectorConfig{}
		err = convertToJsonapiObject(result, &converted)
		if err != nil {
			return tenant_provisioning_service.NewGetTenantConnectorConfigInternalServerError().WithPayload(reportAPIError(fmt.Sprintf("Unable to convert %s data to jsonapi return format: %s", tenmod.TenantConnectorConfigStr, err.Error()), startTime, http.StatusInternalServerError, mon.DeleteTenantConnectorConfigStr, mon.APICompleted, mon.TenantAPICompleted))
		}

		reportAPICompletionState(startTime, http.StatusOK, mon.DeleteTenantConnectorConfigStr, mon.APICompleted, mon.TenantAPICompleted)
		logger.Log.Infof("Deleted %s %s", tenmod.TenantConnectorConfigStr, models.AsJSONString(converted))
		return tenant_provisioning_service.NewDeleteTenantConnectorConfigOK().WithPayload(&converted)
	}
}

// HandleGetAllTenantConnectorConfigs - fetch all connector configurations for a tenant
func HandleGetAllTenantConnectorConfigs(allowedRoles []string, tenantDB datastore.TenantServiceDatastore) func(params tenant_provisioning_service.GetAllTenantConnectorConfigsParams) middleware.Responder {
	return func(params tenant_provisioning_service.GetAllTenantConnectorConfigsParams) middleware.Responder {
		startTime := time.Now()
		incrementAPICounters(mon.APIRecieved, mon.TenantAPIRecieved)
		logger.Log.Infof("Fetching all %ss for Tenant %s", tenmod.TenantConnectorConfigStr, params.TenantID)

		if !isRequestAuthorized(params.HTTPRequest, allowedRoles) {
			return tenant_provisioning_service.NewGetAllTenantConnectorConfigsForbidden().WithPayload(reportAPIError(fmt.Sprintf("Get all %ss operation not authorized for role: %s", tenmod.TenantConnectorConfigStr, params.HTTPRequest.Header.Get(XFwdUserRoles)), startTime, http.StatusForbidden, mon.GetAllTenantConnectorConfigStr, mon.APICompleted, mon.TenantAPICompleted))
		}

		// Issue request to DAO Layer
		result, err := tenantDB.GetAllTenantConnectorConfigs(params.TenantID, *params.Zone)
		if err != nil {
			return tenant_provisioning_service.NewGetAllTenantConnectorConfigsInternalServerError().WithPayload(reportAPIError(fmt.Sprintf("Unable to retrieve %s list: %s", tenmod.TenantConnectorConfigStr, err.Error()), startTime, http.StatusInternalServerError, mon.GetAllTenantConnectorConfigStr, mon.APICompleted, mon.TenantAPICompleted))
		}

		converted := swagmodels.JSONAPITenantConnectorConfigList{}
		err = convertToJsonapiObject(result, &converted)
		if err != nil {
			return tenant_provisioning_service.NewGetAllTenantConnectorConfigsInternalServerError().WithPayload(reportAPIError(fmt.Sprintf("Unable to convert %s list data to jsonapi return format: %s", tenmod.TenantConnectorConfigStr, err.Error()), startTime, http.StatusInternalServerError, mon.GetAllTenantConnectorConfigStr, mon.APICompleted, mon.TenantAPICompleted))
		}

		reportAPICompletionState(startTime, http.StatusOK, mon.GetAllTenantConnectorConfigStr, mon.APICompleted, mon.TenantAPICompleted)
		logger.Log.Infof("Retrieved %d %ss", len(result), tenmod.TenantConnectorConfigStr)
		return tenant_provisioning_service.NewGetAllTenantConnectorConfigsOK().WithPayload(&converted)
	}
}

// HandleCreateTenantConnectorInstance - creates a connector instance for a tenant
func HandleCreateTenantConnectorInstance(allowedRoles []string, tenantDB datastore.TenantServiceDatastore) func(params tenant_provisioning_service.CreateTenantConnectorInstanceParams) middleware.Responder {
	return func(params tenant_provisioning_service.CreateTenantConnectorInstanceParams) middleware.Responder {
		startTime := time.Now()
		incrementAPICounters(mon.APIRecieved, mon.TenantAPIRecieved)
		logger.Log.Infof("Creating %s %s for Tenant %s", tenmod.TenantConnectorInstanceStr, params.Body.Data.Attributes.Hostname, params.TenantID)

		if !isRequestAuthorized(params.HTTPRequest, allowedRoles) {
			return tenant_provisioning_service.NewCreateTenantConnectorInstanceForbidden().WithPayload(reportAPIError(fmt.Sprintf("Create %s operation not authorized for role: %s", tenmod.TenantConnectorInstanceStr, params.HTTPRequest.Header.Get(XFwdUserRoles)), startTime, http.StatusForbidden, mon.CreateTenantConnectorInstanceStr, mon.APICompleted, mon.TenantAPICompleted))
		}

		// Unmarshal the request
		requestBytes, err := json.Marshal(params.Body)
		if err != nil {
			return tenant_provisioning_service.NewCreateTenantConnectorInstanceBadRequest().WithPayload(reportAPIError(generateErrorMessage(http.StatusBadRequest, err.Error()), startTime, http.StatusBadRequest, mon.CreateTenantConnectorInstanceStr, mon.APICompleted, mon.TenantAPICompleted))
		}

		data := tenmod.ConnectorInstance{}
		err = jsonapi.Unmarshal(requestBytes, &data)
		if err != nil {
			return tenant_provisioning_service.NewCreateTenantConnectorInstanceBadRequest().WithPayload(reportAPIError(generateErrorMessage(http.StatusBadRequest, err.Error()), startTime, http.StatusBadRequest, mon.CreateTenantConnectorInstanceStr, mon.APICompleted, mon.TenantAPICompleted))
		}

		err = data.Validate(false)
		if err != nil {
			return tenant_provisioning_service.NewCreateTenantConnectorInstanceBadRequest().WithPayload(reportAPIError(generateErrorMessage(http.StatusBadRequest, err.Error()), startTime, http.StatusBadRequest, mon.CreateTenantConnectorInstanceStr, mon.APICompleted, mon.TenantAPICompleted))
		}

		// Issue request to DAO Layer
		result, err := tenantDB.CreateTenantConnectorInstance(&data)
		if err != nil {
			return tenant_provisioning_service.NewCreateTenantConnectorInstanceInternalServerError().WithPayload(reportAPIError(fmt.Sprintf("Unable to store %s: %s", tenmod.TenantConnectorInstanceStr, err.Error()), startTime, http.StatusInternalServerError, mon.CreateTenantConnectorInstanceStr, mon.APICompleted, mon.TenantAPICompleted))
		}

		converted := swagmodels.JSONAPITenantConnectorInstance{}
		err = convertToJsonapiObject(result, &converted)
		if err != nil {
			return tenant_provisioning_service.NewCreateTenantConnectorInstanceInternalServerError().WithPayload(reportAPIError(fmt.Sprintf("Unable to convert %s data to jsonapi return format: %s", tenmod.TenantConnectorInstanceStr, err.Error()), startTime, http.StatusInternalServerError, mon.CreateTenantConnectorInstanceStr, mon.APICompleted, mon.TenantAPICompleted))
		}

		reportAPICompletionState(startTime, http.StatusOK, mon.CreateTenantConnectorInstanceStr, mon.APICompleted, mon.TenantAPICompleted)
		logger.Log.Infof("Created %s %s", tenmod.TenantConnectorInstanceStr, models.AsJSONString(converted))
		return tenant_provisioning_service.NewCreateTenantConnectorInstanceOK().WithPayload(&converted)
	}
}

// HandleUpdateTenantConnectorInstance - updates a connector instance for a tenant
func HandleUpdateTenantConnectorInstance(allowedRoles []string, tenantDB datastore.TenantServiceDatastore) func(params tenant_provisioning_service.UpdateTenantConnectorInstanceParams) middleware.Responder {
	return func(params tenant_provisioning_service.UpdateTenantConnectorInstanceParams) middleware.Responder {
		startTime := time.Now()
		incrementAPICounters(mon.APIRecieved, mon.TenantAPIRecieved)
		logger.Log.Infof("Updating %s %s for Tenant", tenmod.TenantConnectorInstanceStr, params.Body.Data.Attributes.Hostname, params.TenantID)

		if !isRequestAuthorized(params.HTTPRequest, allowedRoles) {
			return tenant_provisioning_service.NewUpdateTenantConnectorInstanceForbidden().WithPayload(reportAPIError(fmt.Sprintf("Update %s operation not authorized for role: %s", tenmod.TenantConnectorInstanceStr, params.HTTPRequest.Header.Get(XFwdUserRoles)), startTime, http.StatusForbidden, mon.UpdateTenantConnectorInstanceStr, mon.APICompleted, mon.TenantAPICompleted))
		}

		// Unmarshal the request
		requestBytes, err := json.Marshal(params.Body)
		if err != nil {
			return tenant_provisioning_service.NewUpdateTenantConnectorInstanceBadRequest().WithPayload(reportAPIError(generateErrorMessage(http.StatusBadRequest, err.Error()), startTime, http.StatusBadRequest, mon.UpdateTenantConnectorInstanceStr, mon.APICompleted, mon.TenantAPICompleted))
		}

		data := tenmod.ConnectorInstance{}
		err = jsonapi.Unmarshal(requestBytes, &data)
		if err != nil {
			return tenant_provisioning_service.NewUpdateTenantConnectorInstanceBadRequest().WithPayload(reportAPIError(generateErrorMessage(http.StatusBadRequest, err.Error()), startTime, http.StatusBadRequest, mon.UpdateTenantConnectorInstanceStr, mon.APICompleted, mon.TenantAPICompleted))
		}

		err = data.Validate(true)
		if err != nil {
			return tenant_provisioning_service.NewUpdateTenantConnectorInstanceBadRequest().WithPayload(reportAPIError(generateErrorMessage(http.StatusBadRequest, err.Error()), startTime, http.StatusBadRequest, mon.UpdateTenantConnectorInstanceStr, mon.APICompleted, mon.TenantAPICompleted))
		}

		// Issue request to DAO Layer
		result, err := tenantDB.UpdateTenantConnectorInstance(&data)
		if err != nil {
			return tenant_provisioning_service.NewUpdateTenantConnectorInstanceInternalServerError().WithPayload(reportAPIError(fmt.Sprintf("Unable to store %s: %s", tenmod.TenantConnectorInstanceStr, err.Error()), startTime, http.StatusInternalServerError, mon.UpdateTenantConnectorInstanceStr, mon.APICompleted, mon.TenantAPICompleted))
		}

		converted := swagmodels.JSONAPITenantConnectorInstance{}
		err = convertToJsonapiObject(result, &converted)
		if err != nil {
			return tenant_provisioning_service.NewUpdateTenantConnectorInstanceInternalServerError().WithPayload(reportAPIError(fmt.Sprintf("Unable to convert %s data to jsonapi return format: %s", tenmod.TenantConnectorInstanceStr, err.Error()), startTime, http.StatusInternalServerError, mon.UpdateTenantConnectorInstanceStr, mon.APICompleted, mon.TenantAPICompleted))
		}

		reportAPICompletionState(startTime, http.StatusOK, mon.UpdateTenantConnectorInstanceStr, mon.APICompleted, mon.TenantAPICompleted)
		logger.Log.Infof("Updated %s %s", tenmod.TenantConnectorInstanceStr, models.AsJSONString(converted))
		return tenant_provisioning_service.NewUpdateTenantConnectorInstanceOK().WithPayload(&converted)
	}
}

// HandleGetTenantConnectorInstance - fetch a connector instance for a tenant
func HandleGetTenantConnectorInstance(allowedRoles []string, tenantDB datastore.TenantServiceDatastore) func(params tenant_provisioning_service.GetTenantConnectorInstanceParams) middleware.Responder {
	return func(params tenant_provisioning_service.GetTenantConnectorInstanceParams) middleware.Responder {
		startTime := time.Now()
		incrementAPICounters(mon.APIRecieved, mon.TenantAPIRecieved)
		logger.Log.Infof("Fetching %s %s for Tenant %s", tenmod.TenantConnectorInstanceStr, params.ConnectorInstanceID, params.TenantID)

		if !isRequestAuthorized(params.HTTPRequest, allowedRoles) {
			return tenant_provisioning_service.NewGetTenantConnectorInstanceForbidden().WithPayload(reportAPIError(fmt.Sprintf("Get %s operation not authorized for role: %s", tenmod.TenantConnectorInstanceStr, params.HTTPRequest.Header.Get(XFwdUserRoles)), startTime, http.StatusForbidden, mon.GetTenantConnectorInstanceStr, mon.APICompleted, mon.TenantAPICompleted))
		}

		// Issue request to DAO Layer
		result, err := tenantDB.GetTenantConnectorInstance(params.TenantID, params.ConnectorInstanceID)
		if err != nil {
			return tenant_provisioning_service.NewGetTenantConnectorInstanceInternalServerError().WithPayload(reportAPIError(fmt.Sprintf("Unable to retrieve %s: %s", tenmod.TenantConnectorInstanceStr, err.Error()), startTime, http.StatusInternalServerError, mon.GetTenantConnectorInstanceStr, mon.APICompleted, mon.TenantAPICompleted))
		}

		converted := swagmodels.JSONAPITenantConnectorInstance{}
		err = convertToJsonapiObject(result, &converted)
		if err != nil {
			return tenant_provisioning_service.NewGetTenantConnectorInstanceInternalServerError().WithPayload(reportAPIError(fmt.Sprintf("Unable to convert %s data to jsonapi return format: %s", tenmod.TenantConnectorInstanceStr, err.Error()), startTime, http.StatusInternalServerError, mon.GetTenantConnectorInstanceStr, mon.APICompleted, mon.TenantAPICompleted))
		}

		reportAPICompletionState(startTime, http.StatusOK, mon.GetTenantConnectorInstanceStr, mon.APICompleted, mon.TenantAPICompleted)
		logger.Log.Infof("Retrieved %s %s", tenmod.TenantConnectorInstanceStr, models.AsJSONString(converted))
		return tenant_provisioning_service.NewGetTenantConnectorInstanceOK().WithPayload(&converted)
	}
}

// HandleDeleteTenantConnectorInstance - delete a connector instance for a tenant
func HandleDeleteTenantConnectorInstance(allowedRoles []string, tenantDB datastore.TenantServiceDatastore) func(params tenant_provisioning_service.DeleteTenantConnectorInstanceParams) middleware.Responder {
	return func(params tenant_provisioning_service.DeleteTenantConnectorInstanceParams) middleware.Responder {
		startTime := time.Now()
		incrementAPICounters(mon.APIRecieved, mon.TenantAPIRecieved)
		logger.Log.Infof("Deleting %s %s for Tenant %s", tenmod.TenantConnectorInstanceStr, params.ConnectorInstanceID, params.TenantID)

		if !isRequestAuthorized(params.HTTPRequest, allowedRoles) {
			return tenant_provisioning_service.NewGetTenantConnectorInstanceForbidden().WithPayload(reportAPIError(fmt.Sprintf("Delete %s operation not authorized for role: %s", tenmod.TenantConnectorInstanceStr, params.HTTPRequest.Header.Get(XFwdUserRoles)), startTime, http.StatusForbidden, mon.DeleteTenantConnectorInstanceStr, mon.APICompleted, mon.TenantAPICompleted))
		}

		// Issue request to DAO Layer
		result, err := tenantDB.DeleteTenantConnectorInstance(params.TenantID, params.ConnectorInstanceID)
		if err != nil {
			return tenant_provisioning_service.NewGetTenantConnectorInstanceInternalServerError().WithPayload(reportAPIError(fmt.Sprintf("Unable to delete %s: %s", tenmod.TenantConnectorInstanceStr, err.Error()), startTime, http.StatusInternalServerError, mon.DeleteTenantConnectorInstanceStr, mon.APICompleted, mon.TenantAPICompleted))
		}

		converted := swagmodels.JSONAPITenantConnectorInstance{}
		err = convertToJsonapiObject(result, &converted)
		if err != nil {
			return tenant_provisioning_service.NewGetTenantConnectorInstanceInternalServerError().WithPayload(reportAPIError(fmt.Sprintf("Unable to convert %s data to jsonapi return format: %s", tenmod.TenantConnectorInstanceStr, err.Error()), startTime, http.StatusInternalServerError, mon.DeleteTenantConnectorInstanceStr, mon.APICompleted, mon.TenantAPICompleted))
		}

		reportAPICompletionState(startTime, http.StatusOK, mon.DeleteTenantConnectorInstanceStr, mon.APICompleted, mon.TenantAPICompleted)
		logger.Log.Infof("Deleted %s %s", tenmod.TenantConnectorInstanceStr, models.AsJSONString(converted))
		return tenant_provisioning_service.NewDeleteTenantConnectorInstanceOK().WithPayload(&converted)
	}
}

// HandleGetAllTenantConnectorInstances - fetch all connector instances for a tenant
func HandleGetAllTenantConnectorInstances(allowedRoles []string, tenantDB datastore.TenantServiceDatastore) func(params tenant_provisioning_service.GetAllTenantConnectorInstancesParams) middleware.Responder {
	return func(params tenant_provisioning_service.GetAllTenantConnectorInstancesParams) middleware.Responder {
		startTime := time.Now()
		incrementAPICounters(mon.APIRecieved, mon.TenantAPIRecieved)
		logger.Log.Infof("Fetching all %ss for Tenant %s", tenmod.TenantConnectorInstanceStr, params.TenantID)

		if !isRequestAuthorized(params.HTTPRequest, allowedRoles) {
			return tenant_provisioning_service.NewGetAllTenantConnectorInstancesForbidden().WithPayload(reportAPIError(fmt.Sprintf("Get all %ss operation not authorized for role: %s", tenmod.TenantConnectorInstanceStr, params.HTTPRequest.Header.Get(XFwdUserRoles)), startTime, http.StatusForbidden, mon.GetAllTenantConnectorInstanceStr, mon.APICompleted, mon.TenantAPICompleted))
		}

		// Issue request to DAO Layer
		result, err := tenantDB.GetAllTenantConnectorInstances(params.TenantID)
		if err != nil {
			return tenant_provisioning_service.NewGetAllTenantConnectorInstancesInternalServerError().WithPayload(reportAPIError(fmt.Sprintf("Unable to retrieve %s list: %s", tenmod.TenantConnectorInstanceStr, err.Error()), startTime, http.StatusInternalServerError, mon.GetAllTenantConnectorInstanceStr, mon.APICompleted, mon.TenantAPICompleted))
		}

		converted := swagmodels.JSONAPITenantConnectorInstanceList{}
		err = convertToJsonapiObject(result, &converted)
		if err != nil {
			return tenant_provisioning_service.NewGetAllTenantConnectorInstancesInternalServerError().WithPayload(reportAPIError(fmt.Sprintf("Unable to convert %s list data to jsonapi return format: %s", tenmod.TenantConnectorInstanceStr, err.Error()), startTime, http.StatusInternalServerError, mon.GetAllTenantConnectorInstanceStr, mon.APICompleted, mon.TenantAPICompleted))
		}

		reportAPICompletionState(startTime, http.StatusOK, mon.GetAllTenantConnectorInstanceStr, mon.APICompleted, mon.TenantAPICompleted)
		logger.Log.Infof("Retrieved %d %ss", len(result), tenmod.TenantConnectorInstanceStr)
		return tenant_provisioning_service.NewGetAllTenantConnectorInstancesOK().WithPayload(&converted)
	}
}
