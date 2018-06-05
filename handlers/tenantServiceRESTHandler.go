package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	db "github.com/accedian/adh-gather/datastore"
	"github.com/accedian/adh-gather/logger"
	"github.com/accedian/adh-gather/models"
	metmod "github.com/accedian/adh-gather/models/metrics"
	tenmod "github.com/accedian/adh-gather/models/tenant"
	mon "github.com/accedian/adh-gather/monitoring"
	"github.com/accedian/adh-gather/scheduler"
	"github.com/accedian/adh-gather/server"
	"github.com/gorilla/mux"
	"github.com/manyminds/api2go/jsonapi"
)

// TenantServiceRESTHandler - handler of logic for REST calls made to the Tenant Service.
type TenantServiceRESTHandler struct {
	TenantDB db.TenantServiceDatastore
	routes   []server.Route
	notifH   *ChangeNotificationHandler
}

// CreateTenantServiceRESTHandler - used to create a Tenant Service REST handler which provides
// logic to serve the Admin Service REST calls
func CreateTenantServiceRESTHandler() *TenantServiceRESTHandler {
	result := new(TenantServiceRESTHandler)

	// Setup the DB implementation based on configuration
	tdb, err := getTenantServiceDatastore()
	if err != nil {
		logger.Log.Fatalf("Unable to instantiate TenantServiceRESTHandler: %s", err.Error())
	}
	result.TenantDB = tdb
	result.notifH = getChangeNotificationHandler()

	result.routes = []server.Route{
		server.Route{
			Name:        "CreateTenantUser",
			Method:      "POST",
			Pattern:     apiV1Prefix + tenantsAPIPrefix + "users",
			HandlerFunc: BuildRouteHandlerWithRAC([]string{userRoleSkylight, userRoleTenantAdmin}, result.CreateTenantUser),
		},
		server.Route{
			Name:        "UpdateTenantUser",
			Method:      "PUT",
			Pattern:     apiV1Prefix + tenantsAPIPrefix + "users",
			HandlerFunc: BuildRouteHandlerWithRAC([]string{userRoleSkylight, userRoleTenantAdmin}, result.UpdateTenantUser),
		},
		server.Route{
			Name:        "PatchTenantUser",
			Method:      "PATCH",
			Pattern:     apiV1Prefix + tenantsAPIPrefix + "users",
			HandlerFunc: BuildRouteHandlerWithRAC([]string{userRoleSkylight, userRoleTenantAdmin}, result.PatchTenantUser),
		},
		server.Route{
			Name:        "GetTenantUser",
			Method:      "GET",
			Pattern:     apiV1Prefix + tenantsAPIPrefix + "users/{userID}",
			HandlerFunc: result.GetTenantUser,
		},
		server.Route{
			Name:        "DeleteTenantUser",
			Method:      "DELETE",
			Pattern:     apiV1Prefix + tenantsAPIPrefix + "users/{userID}",
			HandlerFunc: BuildRouteHandlerWithRAC([]string{userRoleSkylight, userRoleTenantAdmin}, result.DeleteTenantUser),
		},
		server.Route{
			Name:        "GetAllTenantUsers",
			Method:      "GET",
			Pattern:     apiV1Prefix + tenantsAPIPrefix + "user-list",
			HandlerFunc: result.GetAllTenantUsers,
		},
		server.Route{
			Name:        "CreateTenantDomain",
			Method:      "POST",
			Pattern:     apiV1Prefix + tenantsAPIPrefix + "domains",
			HandlerFunc: BuildRouteHandlerWithRAC([]string{userRoleSkylight, userRoleTenantAdmin}, result.CreateTenantDomain),
		},
		server.Route{
			Name:        "UpdateTenantDomain",
			Method:      "PUT",
			Pattern:     apiV1Prefix + tenantsAPIPrefix + "domains/{domainID}",
			HandlerFunc: BuildRouteHandlerWithRAC([]string{userRoleSkylight, userRoleTenantAdmin}, result.UpdateTenantDomain),
		},
		server.Route{
			Name:        "PatchTenantDomain",
			Method:      "PATCH",
			Pattern:     apiV1Prefix + tenantsAPIPrefix + "domains/{domainID}",
			HandlerFunc: BuildRouteHandlerWithRAC([]string{userRoleSkylight, userRoleTenantAdmin}, result.PatchTenantDomain),
		},
		server.Route{
			Name:        "GetTenantDomain",
			Method:      "GET",
			Pattern:     apiV1Prefix + tenantsAPIPrefix + "domains/{domainID}",
			HandlerFunc: result.GetTenantDomain,
		},
		server.Route{
			Name:        "DeleteTenantDomain",
			Method:      "DELETE",
			Pattern:     apiV1Prefix + tenantsAPIPrefix + "domains/{domainID}",
			HandlerFunc: BuildRouteHandlerWithRAC([]string{userRoleSkylight, userRoleTenantAdmin}, result.DeleteTenantDomain),
		},
		server.Route{
			Name:        "GetAllTenantDomains",
			Method:      "GET",
			Pattern:     apiV1Prefix + tenantsAPIPrefix + "domain-list",
			HandlerFunc: result.GetAllTenantDomains,
		},
		server.Route{
			Name:        "CreateTenantConnectorConfig",
			Method:      "POST",
			Pattern:     apiV1Prefix + tenantsAPIPrefix + "connector-configs",
			HandlerFunc: BuildRouteHandlerWithRAC([]string{userRoleSkylight, userRoleTenantAdmin}, result.CreateTenantConnectorConfig),
		},
		server.Route{
			Name:        "UpdateTenantConnectorConfig",
			Method:      "PUT",
			Pattern:     apiV1Prefix + tenantsAPIPrefix + "connector-configs",
			HandlerFunc: BuildRouteHandlerWithRAC([]string{userRoleSkylight, userRoleTenantAdmin}, result.UpdateTenantConnectorConfig),
		},
		server.Route{
			Name:        "GetTenantConnectorConfig",
			Method:      "GET",
			Pattern:     apiV1Prefix + tenantsAPIPrefix + "connector-configs/{connectorID}",
			HandlerFunc: result.GetTenantConnectorConfig,
		},
		server.Route{
			Name:        "DeleteTenantConnectorConfig",
			Method:      "DELETE",
			Pattern:     apiV1Prefix + tenantsAPIPrefix + "connector-configs/{connectorID}",
			HandlerFunc: BuildRouteHandlerWithRAC([]string{userRoleSkylight, userRoleTenantAdmin}, result.DeleteTenantConnectorConfig),
		},
		server.Route{
			Name:        "GetAllTenantConnectorConfigs",
			Method:      "GET",
			Pattern:     apiV1Prefix + tenantsAPIPrefix + "connector-config-list",
			HandlerFunc: result.GetAllTenantConnectorConfigs,
		},
		server.Route{
			Name:        "CreateTenantConnectorInstance",
			Method:      "POST",
			Pattern:     apiV1Prefix + tenantsAPIPrefix + "connector-instances",
			HandlerFunc: BuildRouteHandlerWithRAC([]string{userRoleSkylight, userRoleTenantAdmin}, result.CreateTenantConnectorInstance),
		},
		server.Route{
			Name:        "UpdateTenantConnectorInstance",
			Method:      "PUT",
			Pattern:     apiV1Prefix + tenantsAPIPrefix + "connector-instances",
			HandlerFunc: BuildRouteHandlerWithRAC([]string{userRoleSkylight, userRoleTenantAdmin}, result.UpdateTenantConnectorInstance),
		},
		server.Route{
			Name:        "GetTenantConnectorInstance",
			Method:      "GET",
			Pattern:     apiV1Prefix + tenantsAPIPrefix + "connector-instances/{connectorID}",
			HandlerFunc: result.GetTenantConnectorInstance,
		},
		server.Route{
			Name:        "DeleteTenantConnectorInstance",
			Method:      "DELETE",
			Pattern:     apiV1Prefix + tenantsAPIPrefix + "connector-instances/{connectorID}",
			HandlerFunc: BuildRouteHandlerWithRAC([]string{userRoleSkylight, userRoleTenantAdmin}, result.DeleteTenantConnectorInstance),
		},
		server.Route{
			Name:        "GetAllTenantConnectorInstances",
			Method:      "GET",
			Pattern:     apiV1Prefix + tenantsAPIPrefix + "connector-instance-list",
			HandlerFunc: result.GetAllTenantConnectorInstances,
		},
		server.Route{
			Name:        "CreateTenantIngestionProfile",
			Method:      "POST",
			Pattern:     apiV1Prefix + tenantsAPIPrefix + "ingestion-profiles",
			HandlerFunc: BuildRouteHandlerWithRAC([]string{userRoleSkylight, userRoleTenantAdmin}, result.CreateTenantIngestionProfile),
		},
		server.Route{
			Name:        "UpdateTenantIngestionProfile",
			Method:      "PUT",
			Pattern:     apiV1Prefix + tenantsAPIPrefix + "ingestion-profiles",
			HandlerFunc: BuildRouteHandlerWithRAC([]string{userRoleSkylight, userRoleTenantAdmin}, result.UpdateTenantIngestionProfile),
		},
		server.Route{
			Name:        "PatchTenantIngestionProfile",
			Method:      "PATCH",
			Pattern:     apiV1Prefix + tenantsAPIPrefix + "ingestion-profiles",
			HandlerFunc: BuildRouteHandlerWithRAC([]string{userRoleSkylight, userRoleTenantAdmin}, result.PatchTenantIngestionProfile),
		},
		server.Route{
			Name:        "GetTenantIngestionProfile",
			Method:      "GET",
			Pattern:     apiV1Prefix + tenantsAPIPrefix + "ingestion-profiles/{dataID}",
			HandlerFunc: result.GetTenantIngestionProfile,
		},
		server.Route{
			Name:        "DeleteTenantIngestionProfile",
			Method:      "DELETE",
			Pattern:     apiV1Prefix + tenantsAPIPrefix + "ingestion-profiles/{dataID}",
			HandlerFunc: BuildRouteHandlerWithRAC([]string{userRoleSkylight, userRoleTenantAdmin}, result.DeleteTenantIngestionProfile),
		},
		server.Route{
			Name:        "GetActiveTenantIngestionProfile",
			Method:      "GET",
			Pattern:     apiV1Prefix + tenantsAPIPrefix + "active-ingestion-profile",
			HandlerFunc: result.GetActiveTenantIngestionProfile,
		},
		server.Route{
			Name:        "CreateTenantThresholdProfile",
			Method:      "POST",
			Pattern:     apiV1Prefix + tenantsAPIPrefix + "threshold-profiles",
			HandlerFunc: BuildRouteHandlerWithRAC([]string{userRoleSkylight, userRoleTenantAdmin}, result.CreateTenantThresholdProfile),
		},
		server.Route{
			Name:        "UpdateTenantThresholdProfile",
			Method:      "PUT",
			Pattern:     apiV1Prefix + tenantsAPIPrefix + "threshold-profiles",
			HandlerFunc: BuildRouteHandlerWithRAC([]string{userRoleSkylight, userRoleTenantAdmin}, result.UpdateTenantThresholdProfile),
		},
		server.Route{
			Name:        "PatchTenantThresholdProfile",
			Method:      "PATCH",
			Pattern:     apiV1Prefix + tenantsAPIPrefix + "threshold-profiles",
			HandlerFunc: BuildRouteHandlerWithRAC([]string{userRoleSkylight, userRoleTenantAdmin}, result.PatchTenantThresholdProfile),
		},
		server.Route{
			Name:        "GetTenantThresholdProfile",
			Method:      "GET",
			Pattern:     apiV1Prefix + tenantsAPIPrefix + "threshold-profiles/{dataID}",
			HandlerFunc: result.GetTenantThresholdProfile,
		},
		server.Route{
			Name:        "DeleteTenantThresholdProfile",
			Method:      "DELETE",
			Pattern:     apiV1Prefix + tenantsAPIPrefix + "threshold-profiles/{dataID}",
			HandlerFunc: BuildRouteHandlerWithRAC([]string{userRoleSkylight, userRoleTenantAdmin}, result.DeleteTenantThresholdProfile),
		},
		server.Route{
			Name:        "GetAllTenantThresholdProfiles",
			Method:      "GET",
			Pattern:     apiV1Prefix + tenantsAPIPrefix + "threshold-profile-list",
			HandlerFunc: result.GetAllTenantThresholdProfiles,
		},
		server.Route{
			Name:        "CreateMonitoredObject",
			Method:      "POST",
			Pattern:     apiV1Prefix + tenantsAPIPrefix + "monitored-objects",
			HandlerFunc: BuildRouteHandlerWithRAC([]string{userRoleSkylight, userRoleTenantAdmin}, result.CreateMonitoredObject),
		},
		server.Route{
			Name:        "BulkInsertMonitoredObject",
			Method:      "POST",
			Pattern:     apiV1Prefix + tenantsAPIPrefix + "bulk/insert/monitored-objects",
			HandlerFunc: BuildRouteHandlerWithRAC([]string{userRoleSkylight, userRoleTenantAdmin}, result.BulkInsertMonitoredObject),
		},
		server.Route{
			Name:        "BulkUpdateMonitoredObject",
			Method:      "PUT",
			Pattern:     apiV1Prefix + tenantsAPIPrefix + "bulk/insert/monitored-objects",
			HandlerFunc: BuildRouteHandlerWithRAC([]string{userRoleSkylight, userRoleTenantAdmin}, result.BulkUpdateMonitoredObject),
		},
		server.Route{
			Name:        "UpdateMonitoredObject",
			Method:      "PUT",
			Pattern:     apiV1Prefix + tenantsAPIPrefix + "monitored-objects",
			HandlerFunc: BuildRouteHandlerWithRAC([]string{userRoleSkylight, userRoleTenantAdmin}, result.UpdateMonitoredObject),
		},
		server.Route{
			Name:        "PatchMonitoredObject",
			Method:      "PATCH",
			Pattern:     apiV1Prefix + tenantsAPIPrefix + "monitored-objects",
			HandlerFunc: BuildRouteHandlerWithRAC([]string{userRoleSkylight, userRoleTenantAdmin}, result.PatchMonitoredObject),
		},
		server.Route{
			Name:        "GetMonitoredObject",
			Method:      "GET",
			Pattern:     apiV1Prefix + tenantsAPIPrefix + "monitored-objects/{dataID}",
			HandlerFunc: result.GetMonitoredObject,
		},
		server.Route{
			Name:        "DeleteMonitoredObject",
			Method:      "DELETE",
			Pattern:     apiV1Prefix + tenantsAPIPrefix + "monitored-objects/{dataID}",
			HandlerFunc: BuildRouteHandlerWithRAC([]string{userRoleSkylight, userRoleTenantAdmin}, result.DeleteMonitoredObject),
		},
		server.Route{
			Name:        "GetAllMonitoredObjects",
			Method:      "GET",
			Pattern:     apiV1Prefix + tenantsAPIPrefix + "monitored-object-list",
			HandlerFunc: result.GetAllMonitoredObjects,
		},
		server.Route{
			Name:        "GetMonitoredObjectToDomainMap",
			Method:      "POST",
			Pattern:     apiV1Prefix + tenantsAPIPrefix + "monitored-object-domain-map",
			HandlerFunc: result.GetMonitoredObjectToDomainMap,
		},
		server.Route{
			Name:        "CreateTenantMeta",
			Method:      "POST",
			Pattern:     apiV1Prefix + tenantsAPIPrefix + "meta",
			HandlerFunc: BuildRouteHandlerWithRAC([]string{userRoleSkylight, userRoleTenantAdmin}, result.CreateTenantMeta),
		},
		server.Route{
			Name:        "PatchTenantMeta",
			Method:      "PATCH",
			Pattern:     apiV1Prefix + tenantsAPIPrefix + "meta",
			HandlerFunc: BuildRouteHandlerWithRAC([]string{userRoleSkylight, userRoleTenantAdmin}, result.PatchTenantMeta),
		},
		server.Route{
			Name:        "UpdateTenantMeta",
			Method:      "PUT",
			Pattern:     apiV1Prefix + tenantsAPIPrefix + "meta",
			HandlerFunc: BuildRouteHandlerWithRAC([]string{userRoleSkylight, userRoleTenantAdmin}, result.UpdateTenantMeta),
		},
		server.Route{
			Name:        "GetTenantMeta",
			Method:      "GET",
			Pattern:     apiV1Prefix + tenantsAPIPrefix + "meta",
			HandlerFunc: result.GetTenantMeta,
		},
		server.Route{
			Name:        "DeleteTenantMeta",
			Method:      "DELETE",
			Pattern:     apiV1Prefix + tenantsAPIPrefix + "meta",
			HandlerFunc: BuildRouteHandlerWithRAC([]string{userRoleSkylight, userRoleTenantAdmin}, result.DeleteTenantMeta),
		},
		server.Route{
			Name:        "CreateReportScheduleConfig",
			Method:      "POST",
			Pattern:     apiV1Prefix + "tenants/{tenantID}/report-schedule-configs",
			HandlerFunc: BuildRouteHandlerWithRAC([]string{userRoleSkylight, userRoleTenantAdmin}, result.CreateReportScheduleConfig),
		},
		server.Route{
			Name:        "UpdateReportScheduleConfig",
			Method:      "PUT",
			Pattern:     apiV1Prefix + "tenants/{tenantID}/report-schedule-configs",
			HandlerFunc: BuildRouteHandlerWithRAC([]string{userRoleSkylight, userRoleTenantAdmin}, result.UpdateReportScheduleConfig),
		},
		server.Route{
			Name:        "GetReportScheduleConfig",
			Method:      "GET",
			Pattern:     apiV1Prefix + "tenants/{tenantID}/report-schedule-configs/{configID}",
			HandlerFunc: BuildRouteHandlerWithRAC([]string{userRoleSkylight, userRoleTenantAdmin}, result.GetReportScheduleConfig),
		},
		server.Route{
			Name:        "DeleteReportScheduleConfig",
			Method:      "DELETE",
			Pattern:     apiV1Prefix + "tenants/{tenantID}/report-schedule-configs/{configID}",
			HandlerFunc: BuildRouteHandlerWithRAC([]string{userRoleSkylight, userRoleTenantAdmin}, result.DeleteReportScheduleConfig),
		},
		server.Route{
			Name:        "GetAllReportScheduleConfigs",
			Method:      "GET",
			Pattern:     apiV1Prefix + "tenants/{tenantID}/report-schedule-config-list",
			HandlerFunc: result.GetAllReportScheduleConfigs,
		},
		server.Route{
			Name:        "GetSLAReport",
			Method:      "GET",
			Pattern:     apiV1Prefix + "tenants/{tenantID}/reports/{reportID}",
			HandlerFunc: result.GetSLAReport,
		},
		server.Route{
			Name:        "GetAllSLAReports",
			Method:      "GET",
			Pattern:     apiV1Prefix + "tenants/{tenantID}/report-list",
			HandlerFunc: result.GetAllSLAReports,
		},
	}

	return result
}

// RegisterAPIHandlers - will bind any REST API routes defined in this service
// to the passed in request multiplexor.
func (tsh *TenantServiceRESTHandler) RegisterAPIHandlers(router *mux.Router) {
	for _, route := range tsh.routes {
		logger.Log.Debugf("Registering endpoint: %v", route)
		router.
			Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			Handler(route.HandlerFunc)
	}
}

// CreateTenantUser - creates a tenant user
func (tsh *TenantServiceRESTHandler) CreateTenantUser(w http.ResponseWriter, r *http.Request) {
	startTime := time.Now()

	// Unmarshal the request
	requestBytes, err := getRequestBytes(r)
	if err != nil {
		msg := generateErrorMessage(http.StatusBadRequest, err.Error())
		reportError(w, startTime, "400", mon.CreateAdminUserStr, msg, http.StatusBadRequest)
		return
	}

	data := tenmod.User{}
	err = jsonapi.Unmarshal(requestBytes, &data)
	if err != nil {
		msg := generateErrorMessage(http.StatusBadRequest, err.Error())
		reportError(w, startTime, "400", mon.CreateAdminUserStr, msg, http.StatusBadRequest)
		return
	}

	err = data.Validate(false)
	if err != nil {
		msg := generateErrorMessage(http.StatusBadRequest, err.Error())
		reportError(w, startTime, "400", mon.CreateAdminUserStr, msg, http.StatusBadRequest)
		return
	}

	logger.Log.Infof("Creating %s: %s", tenmod.TenantUserStr, models.AsJSONString(&data))

	// Issue request to DAO Layer
	result, err := tsh.TenantDB.CreateTenantUser(&data)
	if err != nil {
		msg := fmt.Sprintf("Unable to store %s: %s", tenmod.TenantUserStr, err.Error())
		reportError(w, startTime, "500", mon.CreateTenantUserStr, msg, http.StatusInternalServerError)
		return
	}

	sendSuccessResponse(result, w, startTime, mon.CreateTenantUserStr, tenmod.TenantUserStr, "Created")
}

// PatchTenantUser - updates a tenant user
func (tsh *TenantServiceRESTHandler) PatchTenantUser(w http.ResponseWriter, r *http.Request) {
	startTime := time.Now()
	opStr := mon.PatchTenantUserStr

	// Unmarshal the request
	requestBytes, err := getRequestBytes(r)
	if err != nil {
		msg := generateErrorMessage(http.StatusBadRequest, err.Error())
		reportError(w, startTime, "400", opStr, msg, http.StatusBadRequest)
		return
	}

	data := tenmod.User{}
	err = jsonapi.Unmarshal(requestBytes, &data)
	if err != nil {
		msg := generateErrorMessage(http.StatusBadRequest, err.Error())
		reportError(w, startTime, "400", opStr, msg, http.StatusBadRequest)
		return
	}

	err = data.Validate(true)
	if err != nil {
		msg := generateErrorMessage(http.StatusBadRequest, err.Error())
		reportError(w, startTime, "400", opStr, msg, http.StatusBadRequest)
		return
	}

	// Issue request to DAO Layer
	oldData, err := tsh.TenantDB.GetTenantUser(data.TenantID, data.ID)
	if err != nil {
		msg := fmt.Sprintf("Unable to retrieve %s: %s", tenmod.TenantUserStr, err.Error())
		reportError(w, startTime, "500", mon.GetTenantUserStr, msg, http.StatusInternalServerError)
		return
	}

	errMerge := models.MergeObjWithMap(oldData, requestBytes)
	if errMerge != nil {
		msg := generateErrorMessage(http.StatusBadRequest, err.Error())
		reportError(w, startTime, "400", opStr, msg, http.StatusBadRequest)
		return

	}
	// This only checks if the ID&REV is set.
	err = oldData.Validate(true)
	if err != nil {
		msg := generateErrorMessage(http.StatusBadRequest, err.Error())
		reportError(w, startTime, "400", opStr, msg, http.StatusBadRequest)
		return
	}

	logger.Log.Infof("Patching %s: %s", tenmod.TenantUserStr, oldData)

	// Issue request to DAO Layer
	result, err := tsh.TenantDB.UpdateTenantUser(oldData)
	if err != nil {
		msg := fmt.Sprintf("Unable to store %s: %s", tenmod.TenantUserStr, err.Error())
		reportError(w, startTime, "500", mon.UpdateTenantUserStr, msg, http.StatusInternalServerError)
		return
	}

	sendSuccessResponse(result, w, startTime, mon.UpdateTenantUserStr, tenmod.TenantUserStr, "Patched")
}

// UpdateTenantUser - updates a tenant user
func (tsh *TenantServiceRESTHandler) UpdateTenantUser(w http.ResponseWriter, r *http.Request) {
	startTime := time.Now()

	// Unmarshal the request
	requestBytes, err := getRequestBytes(r)
	if err != nil {
		msg := generateErrorMessage(http.StatusBadRequest, err.Error())
		reportError(w, startTime, "400", mon.CreateAdminUserStr, msg, http.StatusBadRequest)
		return
	}

	data := tenmod.User{}
	err = jsonapi.Unmarshal(requestBytes, &data)
	if err != nil {
		msg := generateErrorMessage(http.StatusBadRequest, err.Error())
		reportError(w, startTime, "400", mon.CreateAdminUserStr, msg, http.StatusBadRequest)
		return
	}

	err = data.Validate(true)
	if err != nil {
		msg := generateErrorMessage(http.StatusBadRequest, err.Error())
		reportError(w, startTime, "400", mon.CreateAdminUserStr, msg, http.StatusBadRequest)
		return
	}

	logger.Log.Infof("Updating %s: %s", tenmod.TenantUserStr, models.AsJSONString(&data))

	// Issue request to DAO Layer
	result, err := tsh.TenantDB.UpdateTenantUser(&data)
	if err != nil {
		msg := fmt.Sprintf("Unable to store %s: %s", tenmod.TenantUserStr, err.Error())
		reportError(w, startTime, "500", mon.PatchTenantStr, msg, http.StatusInternalServerError)
		return
	}

	sendSuccessResponse(result, w, startTime, mon.PatchTenantStr, tenmod.TenantUserStr, "Updated")
}

// GetTenantUser - fetches a tenant user
func (tsh *TenantServiceRESTHandler) GetTenantUser(w http.ResponseWriter, r *http.Request) {
	startTime := time.Now()

	// Get the IDs from the URL
	tenantID := getDBFieldFromRequest(r, 4)
	userID := getDBFieldFromRequest(r, 6)

	logger.Log.Infof("Fetching %s: %s", tenmod.TenantUserStr, userID)

	// Issue request to DAO Layer
	result, err := tsh.TenantDB.GetTenantUser(tenantID, userID)
	if err != nil {
		msg := fmt.Sprintf("Unable to retrieve %s: %s", tenmod.TenantUserStr, err.Error())
		reportError(w, startTime, "500", mon.GetTenantUserStr, msg, http.StatusInternalServerError)
		return
	}

	sendSuccessResponse(result, w, startTime, mon.GetTenantUserStr, tenmod.TenantUserStr, "Retrieved")
}

// DeleteTenantUser - deletes a tenant user
func (tsh *TenantServiceRESTHandler) DeleteTenantUser(w http.ResponseWriter, r *http.Request) {
	startTime := time.Now()

	// Get the IDs from the URL
	tenantID := getDBFieldFromRequest(r, 4)
	userID := getDBFieldFromRequest(r, 6)

	logger.Log.Infof("Deleting %s: %s", tenmod.TenantUserStr, userID)

	// Issue request to DAO Layer
	result, err := tsh.TenantDB.DeleteTenantUser(tenantID, userID)
	if err != nil {
		msg := fmt.Sprintf("Unable to retrieve %s: %s", tenmod.TenantUserStr, err.Error())
		reportError(w, startTime, "500", mon.DeleteTenantUserStr, msg, http.StatusInternalServerError)
		return
	}

	sendSuccessResponse(result, w, startTime, mon.DeleteTenantUserStr, tenmod.TenantUserStr, "Deleted")
}

// GetAllTenantUsers - fetches list of tenant users
func (tsh *TenantServiceRESTHandler) GetAllTenantUsers(w http.ResponseWriter, r *http.Request) {
	startTime := time.Now()

	tenantID := getDBFieldFromRequest(r, 4)

	logger.Log.Infof("Fetching %s list for Tenant %s", tenmod.TenantUserStr, tenantID)

	// Issue request to DAO Layer
	result, err := tsh.TenantDB.GetAllTenantUsers(tenantID)
	if err != nil {
		msg := fmt.Sprintf("Unable to retrieve %s list: %s", tenmod.TenantUserStr, err.Error())
		reportError(w, startTime, "500", mon.GetAllTenantUserStr, msg, http.StatusInternalServerError)
		return
	}

	sendSuccessResponse(result, w, startTime, mon.GetAllTenantUserStr, tenmod.TenantUserStr, "Retrieved list of")
}

// CreateTenantConnectorConfig - creates a tenant ConnectorConfig
func (tsh *TenantServiceRESTHandler) CreateTenantConnectorConfig(w http.ResponseWriter, r *http.Request) {
	startTime := time.Now()

	// Unmarshal the request
	requestBytes, err := getRequestBytes(r)

	if err != nil {
		msg := generateErrorMessage(http.StatusBadRequest, err.Error())
		reportError(w, startTime, "400", mon.CreateTenantConnectorConfigStr, msg, http.StatusBadRequest)
		return
	}

	data := tenmod.ConnectorConfig{}
	err = jsonapi.Unmarshal(requestBytes, &data)
	if err != nil {
		msg := generateErrorMessage(http.StatusBadRequest, err.Error())
		reportError(w, startTime, "400", mon.CreateTenantConnectorConfigStr, msg, http.StatusBadRequest)
		return
	}

	err = data.Validate(false)
	if err != nil {
		msg := generateErrorMessage(http.StatusBadRequest, err.Error())
		reportError(w, startTime, "400", mon.CreateTenantConnectorConfigStr, msg, http.StatusBadRequest)
		return
	}

	logger.Log.Infof("Creating %s: %s", tenmod.TenantConnectorConfigStr, models.AsJSONString(&data))

	// Issue request to DAO Layer
	result, err := tsh.TenantDB.CreateTenantConnectorConfig(&data)
	if err != nil {
		msg := fmt.Sprintf("Unable to store %s: %s", tenmod.TenantConnectorConfigStr, err.Error())
		reportError(w, startTime, "500", mon.CreateTenantConnectorConfigStr, msg, http.StatusInternalServerError)
		return
	}

	sendSuccessResponse(result, w, startTime, mon.CreateTenantConnectorConfigStr, tenmod.TenantConnectorConfigStr, "Created")
}

// UpdateTenantConnectorConfig - updates a tenant ConnectorConfig
func (tsh *TenantServiceRESTHandler) UpdateTenantConnectorConfig(w http.ResponseWriter, r *http.Request) {
	startTime := time.Now()

	// Unmarshal the request
	requestBytes, err := getRequestBytes(r)
	if err != nil {
		msg := generateErrorMessage(http.StatusBadRequest, err.Error())
		reportError(w, startTime, "400", mon.UpdateTenantConnectorConfigStr, msg, http.StatusBadRequest)
		return
	}

	data := tenmod.ConnectorConfig{}
	err = jsonapi.Unmarshal(requestBytes, &data)
	if err != nil {
		msg := generateErrorMessage(http.StatusBadRequest, err.Error())
		reportError(w, startTime, "400", mon.UpdateTenantConnectorConfigStr, msg, http.StatusBadRequest)
		return
	}

	err = data.Validate(true)
	if err != nil {
		msg := generateErrorMessage(http.StatusBadRequest, err.Error())
		reportError(w, startTime, "400", mon.UpdateTenantConnectorConfigStr, msg, http.StatusBadRequest)
		return
	}

	logger.Log.Infof("Updating %s: %s", tenmod.TenantConnectorConfigStr, models.AsJSONString(&data))

	// Issue request to DAO Layer
	result, err := tsh.TenantDB.UpdateTenantConnectorConfig(&data)
	if err != nil {
		msg := fmt.Sprintf("Unable to store %s: %s", tenmod.TenantConnectorConfigStr, err.Error())
		reportError(w, startTime, "500", mon.UpdateTenantConnectorConfigStr, msg, http.StatusInternalServerError)
		return
	}

	sendSuccessResponse(result, w, startTime, mon.UpdateTenantConnectorConfigStr, tenmod.TenantConnectorConfigStr, "Updated")
}

// GetTenantConnectorConfig - fetches a tenant ConnectorConfig
func (tsh *TenantServiceRESTHandler) GetTenantConnectorConfig(w http.ResponseWriter, r *http.Request) {
	startTime := time.Now()

	// Get the IDs from the URL
	tenantID := getDBFieldFromRequest(r, 4)
	userID := getDBFieldFromRequest(r, 6)

	logger.Log.Infof("Fetching %s: %s", tenmod.TenantConnectorConfigStr, userID)

	// Issue request to DAO Layer
	result, err := tsh.TenantDB.GetTenantConnectorConfig(tenantID, userID)
	if err != nil {
		msg := fmt.Sprintf("Unable to retrieve %s: %s", tenmod.TenantConnectorConfigStr, err.Error())
		reportError(w, startTime, "500", mon.GetTenantConnectorConfigStr, msg, http.StatusInternalServerError)
		return
	}

	sendSuccessResponse(result, w, startTime, mon.GetTenantConnectorConfigStr, tenmod.TenantConnectorConfigStr, "Retrieved")
}

// DeleteTenantConnectorConfig - deletes a tenant ConnectorConfig
func (tsh *TenantServiceRESTHandler) DeleteTenantConnectorConfig(w http.ResponseWriter, r *http.Request) {
	startTime := time.Now()

	// Get the IDs from the URL
	tenantID := getDBFieldFromRequest(r, 4)
	userID := getDBFieldFromRequest(r, 6)

	logger.Log.Infof("Deleting %s: %s", tenmod.TenantConnectorConfigStr, userID)

	// Issue request to DAO Layer
	result, err := tsh.TenantDB.DeleteTenantConnectorConfig(tenantID, userID)
	if err != nil {
		msg := fmt.Sprintf("Unable to retrieve %s: %s", tenmod.TenantConnectorConfigStr, err.Error())
		reportError(w, startTime, "500", mon.DeleteTenantConnectorConfigStr, msg, http.StatusInternalServerError)
		return
	}

	sendSuccessResponse(result, w, startTime, mon.DeleteTenantConnectorConfigStr, tenmod.TenantConnectorConfigStr, "Deleted")
}

// GetAllTenantConnectorConfigs - fetches list of tenant ConnectorConfigs
func (tsh *TenantServiceRESTHandler) GetAllTenantConnectorConfigs(w http.ResponseWriter, r *http.Request) {
	startTime := time.Now()

	tenantID := getDBFieldFromRequest(r, 4)
	zones := r.URL.Query()["zone"]
	zone := ""

	if len(zones) > 0 {
		zone = zones[0]
	}

	logger.Log.Infof("Fetching %s list for Tenant %s", tenmod.TenantConnectorConfigStr, tenantID)

	// Issue request to DAO Layer
	result, err := tsh.TenantDB.GetAllTenantConnectorConfigs(tenantID, zone)
	if err != nil {
		msg := fmt.Sprintf("Unable to retrieve %s list: %s", tenmod.TenantConnectorConfigStr, err.Error())
		reportError(w, startTime, "500", mon.GetAllTenantConnectorConfigStr, msg, http.StatusInternalServerError)
		return
	}

	sendSuccessResponse(result, w, startTime, mon.GetAllTenantConnectorConfigStr, tenmod.TenantConnectorConfigStr, "Retrieved list of")
}

// CreateTenantConnectorInstance - creates a tenant ConnectorInstance
func (tsh *TenantServiceRESTHandler) CreateTenantConnectorInstance(w http.ResponseWriter, r *http.Request) {
	startTime := time.Now()

	// Unmarshal the request
	requestBytes, err := getRequestBytes(r)

	if err != nil {
		msg := generateErrorMessage(http.StatusBadRequest, err.Error())
		reportError(w, startTime, "400", mon.CreateTenantConnectorInstanceStr, msg, http.StatusBadRequest)
		return
	}

	data := tenmod.ConnectorInstance{}
	err = jsonapi.Unmarshal(requestBytes, &data)
	if err != nil {
		msg := generateErrorMessage(http.StatusBadRequest, err.Error())
		reportError(w, startTime, "400", mon.CreateTenantConnectorInstanceStr, msg, http.StatusBadRequest)
		return
	}

	err = data.Validate(false)
	if err != nil {
		msg := generateErrorMessage(http.StatusBadRequest, err.Error())
		reportError(w, startTime, "400", mon.CreateTenantConnectorInstanceStr, msg, http.StatusBadRequest)
		return
	}

	logger.Log.Infof("Creating %s: %s", tenmod.TenantConnectorInstanceStr, models.AsJSONString(&data))

	// Issue request to DAO Layer
	result, err := tsh.TenantDB.CreateTenantConnectorInstance(&data)
	if err != nil {
		msg := fmt.Sprintf("Unable to store %s: %s", tenmod.TenantConnectorInstanceStr, err.Error())
		reportError(w, startTime, "500", mon.CreateTenantConnectorInstanceStr, msg, http.StatusInternalServerError)
		return
	}

	sendSuccessResponse(result, w, startTime, mon.CreateTenantConnectorInstanceStr, tenmod.TenantConnectorInstanceStr, "Created")
}

// UpdateTenantConnectorInstance - updates a tenant ConnectorInstance
func (tsh *TenantServiceRESTHandler) UpdateTenantConnectorInstance(w http.ResponseWriter, r *http.Request) {
	startTime := time.Now()

	// Unmarshal the request
	requestBytes, err := getRequestBytes(r)
	if err != nil {
		msg := generateErrorMessage(http.StatusBadRequest, err.Error())
		reportError(w, startTime, "400", mon.UpdateTenantConnectorInstanceStr, msg, http.StatusBadRequest)
		return
	}

	data := tenmod.ConnectorInstance{}
	err = jsonapi.Unmarshal(requestBytes, &data)
	if err != nil {
		msg := generateErrorMessage(http.StatusBadRequest, err.Error())
		reportError(w, startTime, "400", mon.UpdateTenantConnectorInstanceStr, msg, http.StatusBadRequest)
		return
	}

	err = data.Validate(true)
	if err != nil {
		msg := generateErrorMessage(http.StatusBadRequest, err.Error())
		reportError(w, startTime, "400", mon.UpdateTenantConnectorInstanceStr, msg, http.StatusBadRequest)
		return
	}

	logger.Log.Infof("Updating %s: %s", tenmod.TenantConnectorInstanceStr, models.AsJSONString(&data))

	// Issue request to DAO Layer
	result, err := tsh.TenantDB.UpdateTenantConnectorInstance(&data)
	if err != nil {
		msg := fmt.Sprintf("Unable to store %s: %s", tenmod.TenantConnectorInstanceStr, err.Error())
		reportError(w, startTime, "500", mon.UpdateTenantConnectorInstanceStr, msg, http.StatusInternalServerError)
		return
	}

	sendSuccessResponse(result, w, startTime, mon.UpdateTenantConnectorInstanceStr, tenmod.TenantConnectorInstanceStr, "Updated")
}

// GetTenantConnectorInstance - fetches a tenant ConnectorInstance
func (tsh *TenantServiceRESTHandler) GetTenantConnectorInstance(w http.ResponseWriter, r *http.Request) {
	startTime := time.Now()

	// Get the IDs from the URL
	tenantID := getDBFieldFromRequest(r, 4)
	userID := getDBFieldFromRequest(r, 6)

	logger.Log.Infof("Fetching %s: %s", tenmod.TenantConnectorInstanceStr, userID)

	// Issue request to DAO Layer
	result, err := tsh.TenantDB.GetTenantConnectorInstance(tenantID, userID)
	if err != nil {
		msg := fmt.Sprintf("Unable to retrieve %s: %s", tenmod.TenantConnectorInstanceStr, err.Error())
		reportError(w, startTime, "500", mon.GetTenantConnectorInstanceStr, msg, http.StatusInternalServerError)
		return
	}

	sendSuccessResponse(result, w, startTime, mon.GetTenantConnectorInstanceStr, tenmod.TenantConnectorInstanceStr, "Retrieved")
}

// DeleteTenantConnectorInstance - deletes a tenant ConnectorInstance
func (tsh *TenantServiceRESTHandler) DeleteTenantConnectorInstance(w http.ResponseWriter, r *http.Request) {
	startTime := time.Now()

	// Get the IDs from the URL
	tenantID := getDBFieldFromRequest(r, 4)
	userID := getDBFieldFromRequest(r, 6)

	logger.Log.Infof("Deleting %s: %s", tenmod.TenantConnectorInstanceStr, userID)

	// Issue request to DAO Layer
	result, err := tsh.TenantDB.DeleteTenantConnectorInstance(tenantID, userID)
	if err != nil {
		msg := fmt.Sprintf("Unable to retrieve %s: %s", tenmod.TenantConnectorInstanceStr, err.Error())
		reportError(w, startTime, "500", mon.DeleteTenantConnectorInstanceStr, msg, http.StatusInternalServerError)
		return
	}

	sendSuccessResponse(result, w, startTime, mon.DeleteTenantConnectorInstanceStr, tenmod.TenantConnectorInstanceStr, "Deleted")
}

// GetAllTenantConnectorInstances - fetches list of tenant ConnectorInstances
func (tsh *TenantServiceRESTHandler) GetAllTenantConnectorInstances(w http.ResponseWriter, r *http.Request) {
	startTime := time.Now()

	tenantID := getDBFieldFromRequest(r, 4)

	logger.Log.Infof("Fetching %s list for Tenant %s", tenmod.TenantConnectorInstanceStr, tenantID)

	// Issue request to DAO Layer
	result, err := tsh.TenantDB.GetAllTenantConnectorInstances(tenantID)
	if err != nil {
		msg := fmt.Sprintf("Unable to retrieve %s list: %s", tenmod.TenantConnectorInstanceStr, err.Error())
		reportError(w, startTime, "500", mon.GetAllTenantConnectorInstanceStr, msg, http.StatusInternalServerError)
		return
	}

	sendSuccessResponse(result, w, startTime, mon.GetAllTenantConnectorInstanceStr, tenmod.TenantConnectorInstanceStr, "Retrieved list of")
}

// CreateTenantDomain - creates a tenant domain
func (tsh *TenantServiceRESTHandler) CreateTenantDomain(w http.ResponseWriter, r *http.Request) {
	startTime := time.Now()

	// Unmarshal the request
	requestBytes, err := getRequestBytes(r)
	if err != nil {
		msg := generateErrorMessage(http.StatusBadRequest, err.Error())
		reportError(w, startTime, "400", mon.CreateAdminUserStr, msg, http.StatusBadRequest)
		return
	}

	data := tenmod.Domain{}
	err = jsonapi.Unmarshal(requestBytes, &data)
	if err != nil {
		msg := generateErrorMessage(http.StatusBadRequest, err.Error())
		reportError(w, startTime, "400", mon.CreateAdminUserStr, msg, http.StatusBadRequest)
		return
	}

	err = data.Validate(false)
	if err != nil {
		msg := generateErrorMessage(http.StatusBadRequest, err.Error())
		reportError(w, startTime, "400", mon.CreateAdminUserStr, msg, http.StatusBadRequest)
		return
	}

	logger.Log.Infof("Creating %s: %s", tenmod.TenantDomainStr, models.AsJSONString(&data))

	// Issue request to DAO Layer
	result, err := tsh.TenantDB.CreateTenantDomain(&data)
	if err != nil {
		msg := fmt.Sprintf("Unable to store %s: %s", tenmod.TenantDomainStr, err.Error())
		reportError(w, startTime, "500", mon.CreateTenantDomainStr, msg, http.StatusInternalServerError)
		return
	}

	NotifyDomainCreated(data.TenantID, &data)
	sendSuccessResponse(result, w, startTime, mon.CreateTenantDomainStr, tenmod.TenantDomainStr, "Created")
}

// UpdateTenantDomain - updates a tenant domain
func (tsh *TenantServiceRESTHandler) UpdateTenantDomain(w http.ResponseWriter, r *http.Request) {
	startTime := time.Now()

	// Unmarshal the request
	requestBytes, err := getRequestBytes(r)
	if err != nil {
		msg := generateErrorMessage(http.StatusBadRequest, err.Error())
		reportError(w, startTime, "400", mon.CreateAdminUserStr, msg, http.StatusBadRequest)
		return
	}

	data := tenmod.Domain{}
	err = jsonapi.Unmarshal(requestBytes, &data)
	if err != nil {
		msg := generateErrorMessage(http.StatusBadRequest, err.Error())
		reportError(w, startTime, "400", mon.CreateAdminUserStr, msg, http.StatusBadRequest)
		return
	}

	err = data.Validate(true)
	if err != nil {
		msg := generateErrorMessage(http.StatusBadRequest, err.Error())
		reportError(w, startTime, "400", mon.CreateAdminUserStr, msg, http.StatusBadRequest)
		return
	}

	logger.Log.Infof("Updating %s: %s", tenmod.TenantDomainStr, models.AsJSONString(&data))

	// Issue request to DAO Layer
	result, err := tsh.TenantDB.UpdateTenantDomain(&data)
	if err != nil {
		msg := fmt.Sprintf("Unable to store %s: %s", tenmod.TenantDomainStr, err.Error())
		reportError(w, startTime, "500", mon.UpdateTenantDomainStr, msg, http.StatusInternalServerError)
		return
	}

	NotifyDomainUpdated(data.TenantID, &data)
	sendSuccessResponse(result, w, startTime, mon.UpdateTenantDomainStr, tenmod.TenantDomainStr, "Updated")
}

/*PatchTenantDomain - Patches a tenant domain
 * Takes in a partial Tenant domain payload and determines which elements need to be updated.
 * Works by
 * - Unmarshals the payload into a tenant domain struct
 * - Validates the data
 * - Gets the latest document from the DB.
 * - For each element that is a non-zero length string,
 * - Overwrite the latest document element.
 * - Store the changes.
 * Note that limited validation is done on the payload itself. I'm pushing the onus
 * to the database to reject the document if there is an error in the ID/REV.
 */
func (tsh *TenantServiceRESTHandler) PatchTenantDomain(w http.ResponseWriter, r *http.Request) {
	startTime := time.Now()
	opStr := mon.PatchTenantStr

	// Unmarshal the request
	requestBytes, err := getRequestBytes(r)
	if err != nil {
		msg := generateErrorMessage(http.StatusBadRequest, err.Error())
		reportError(w, startTime, "400", opStr, msg, http.StatusBadRequest)
		return
	}

	// Model
	data := tenmod.Domain{}
	err = jsonapi.Unmarshal(requestBytes, &data)
	if err != nil {
		msg := generateErrorMessage(http.StatusBadRequest, err.Error())
		reportError(w, startTime, "400", opStr, msg, http.StatusBadRequest)
		return
	}

	// This only checks if the ID&REV is set.
	err = data.Validate(true)
	if err != nil {
		msg := generateErrorMessage(http.StatusBadRequest, err.Error())
		reportError(w, startTime, "400", opStr, msg, http.StatusBadRequest)
		return
	}

	// Get the IDs from the URL
	tenantID := data.TenantID
	domainID := data.ID

	logger.Log.Infof("Fetching %s: %s", tenmod.TenantDomainStr, domainID)

	// Issue request to DAO Layer
	oldDomain, err := tsh.TenantDB.GetTenantDomain(tenantID, domainID)
	if err != nil {
		msg := fmt.Sprintf("Unable to retrieve %s: %s", tenmod.TenantDomainStr, err.Error())
		reportError(w, startTime, "500", mon.GetTenantDomainStr, msg, http.StatusInternalServerError)
		return
	}

	errMerge := models.MergeObjWithMap(oldDomain, requestBytes)
	if errMerge != nil {
		msg := generateErrorMessage(http.StatusBadRequest, err.Error())
		reportError(w, startTime, "400", opStr, msg, http.StatusBadRequest)
		return

	}
	// This only checks if the ID&REV is set.
	err = oldDomain.Validate(true)
	if err != nil {
		msg := generateErrorMessage(http.StatusBadRequest, err.Error())
		reportError(w, startTime, "400", opStr, msg, http.StatusBadRequest)
		return
	}
	// Done checking for differences
	logger.Log.Infof("Patching %s: %s", opStr, oldDomain)

	// Issue request to DAO Layer
	result, err := tsh.TenantDB.UpdateTenantDomain(oldDomain)
	if err != nil {
		msg := fmt.Sprintf("Unable to store %s: %s", tenmod.TenantDomainStr, err.Error())
		reportError(w, startTime, "500", opStr, msg, http.StatusInternalServerError)
		return
	}

	NotifyDomainUpdated(oldDomain.TenantID, oldDomain)
	sendSuccessResponse(result, w, startTime, opStr, tenmod.TenantDomainStr, "Patched")
}

// GetTenantDomain - fetches a tenant domain
func (tsh *TenantServiceRESTHandler) GetTenantDomain(w http.ResponseWriter, r *http.Request) {
	startTime := time.Now()

	// Get the IDs from the URL
	tenantID := getDBFieldFromRequest(r, 4)
	userID := getDBFieldFromRequest(r, 6)

	logger.Log.Infof("Fetching %s: %s", tenmod.TenantDomainStr, userID)

	// Issue request to DAO Layer
	result, err := tsh.TenantDB.GetTenantDomain(tenantID, userID)
	if err != nil {
		msg := fmt.Sprintf("Unable to retrieve %s: %s", tenmod.TenantDomainStr, err.Error())
		reportError(w, startTime, "500", mon.GetTenantDomainStr, msg, http.StatusInternalServerError)
		return
	}

	sendSuccessResponse(result, w, startTime, mon.GetTenantDomainStr, tenmod.TenantDomainStr, "Retrieved")
}

// DeleteTenantDomain - deletes a tenant domain
func (tsh *TenantServiceRESTHandler) DeleteTenantDomain(w http.ResponseWriter, r *http.Request) {
	startTime := time.Now()

	// Get the IDs from the URL
	tenantID := getDBFieldFromRequest(r, 4)
	domainID := getDBFieldFromRequest(r, 6)

	// Integrity Check - Monitored Objects
	moByDomainReq := tenmod.MonitoredObjectCountByDomainRequest{
		TenantID:  tenantID,
		ByCount:   true,
		DomainSet: []string{domainID},
	}
	moByDomainResp, err := tsh.TenantDB.GetMonitoredObjectToDomainMap(&moByDomainReq)
	if err != nil {
		msg := fmt.Sprintf("Unable to perform integrity check for %s deletion: %s", tenmod.TenantDomainStr, err.Error())
		reportError(w, startTime, "500", mon.DeleteTenantDomainStr, msg, http.StatusInternalServerError)
		return
	}

	logger.Log.Infof("%s got %s", models.AsJSONString(moByDomainReq), models.AsJSONString(moByDomainResp))
	if moByDomainResp.DomainToMonitoredObjectCountMap != nil {
		if count, exists := moByDomainResp.DomainToMonitoredObjectCountMap[domainID]; exists && count > 0 {
			msg := fmt.Sprintf("%s deletion failed integrity check: in use by at least one %s", tenmod.TenantDomainStr, tenmod.TenantMonitoredObjectStr)
			reportError(w, startTime, "500", mon.DeleteTenantDomainStr, msg, http.StatusInternalServerError)
			return
		}
	}

	// Integrity Check - Dashboards
	dashboardUsesDomain, err := tsh.TenantDB.HasDashboardsWithDomain(tenantID, domainID)
	if err != nil {
		msg := fmt.Sprintf("Unable to perform integrity check for %s deletion: %s", tenmod.TenantDomainStr, err.Error())
		reportError(w, startTime, "500", mon.DeleteTenantDomainStr, msg, http.StatusInternalServerError)
		return
	}
	if dashboardUsesDomain {
		msg := fmt.Sprintf("%s deletion failed integrity check: in use by at least one %s", tenmod.TenantDomainStr, tenmod.TenantDashboardStr)
		reportError(w, startTime, "500", mon.DeleteTenantDomainStr, msg, http.StatusInternalServerError)
		return
	}

	configs, err := tsh.TenantDB.GetAllReportScheduleConfigs(tenantID)
	if err != nil {
		msg := fmt.Sprintf("Unable to perform integrity check for %s deletion: %s", tenmod.TenantDomainStr, err.Error())
		reportError(w, startTime, "500", mon.DeleteTenantDomainStr, msg, http.StatusInternalServerError)
		return
	}
	for _, rep := range configs {
		if len(rep.Domains) == 0 {
			continue
		}
		for _, dom := range rep.Domains {
			if dom == domainID {
				msg := fmt.Sprintf("%s deletion failed integrity check: in use by at least one %s", tenmod.TenantDomainStr, tenmod.TenantReportScheduleConfigStr)
				reportError(w, startTime, "500", mon.DeleteTenantDomainStr, msg, http.StatusInternalServerError)
				return
			}
		}
	}

	logger.Log.Infof("Deleting %s: %s", tenmod.TenantDomainStr, domainID)

	// Issue request to DAO Layer
	result, err := tsh.TenantDB.DeleteTenantDomain(tenantID, domainID)
	if err != nil {
		msg := fmt.Sprintf("Unable to retrieve %s: %s", tenmod.TenantDomainStr, err.Error())
		reportError(w, startTime, "500", mon.DeleteTenantDomainStr, msg, http.StatusInternalServerError)
		return
	}

	NotifyDomainDeleted(tenantID, result)
	sendSuccessResponse(result, w, startTime, mon.DeleteTenantDomainStr, tenmod.TenantDomainStr, "Deleted")
}

// GetAllTenantDomains - fetches list of tenant domains
func (tsh *TenantServiceRESTHandler) GetAllTenantDomains(w http.ResponseWriter, r *http.Request) {
	startTime := time.Now()

	tenantID := getDBFieldFromRequest(r, 4)

	logger.Log.Infof("Fetching %s list for Tenant %s", tenmod.TenantDomainStr, tenantID)

	// Issue request to DAO Layer
	result, err := tsh.TenantDB.GetAllTenantDomains(tenantID)
	if err != nil {
		msg := fmt.Sprintf("Unable to retrieve %s list: %s", tenmod.TenantDomainStr, err.Error())
		reportError(w, startTime, "500", mon.GetAllTenantDomainStr, msg, http.StatusInternalServerError)
		return
	}

	sendSuccessResponse(result, w, startTime, mon.GetAllTenantDomainStr, tenmod.TenantDomainStr, "Retrieved list of")
}

// CreateTenantIngestionProfile - creates a tenant ingestion profile
func (tsh *TenantServiceRESTHandler) CreateTenantIngestionProfile(w http.ResponseWriter, r *http.Request) {
	startTime := time.Now()

	// Unmarshal the request
	requestBytes, err := getRequestBytes(r)
	if err != nil {
		msg := generateErrorMessage(http.StatusBadRequest, err.Error())
		reportError(w, startTime, "400", mon.CreateAdminUserStr, msg, http.StatusBadRequest)
		return
	}

	data := tenmod.IngestionProfile{}
	err = jsonapi.Unmarshal(requestBytes, &data)
	if err != nil {
		msg := generateErrorMessage(http.StatusBadRequest, err.Error())
		reportError(w, startTime, "400", mon.CreateAdminUserStr, msg, http.StatusBadRequest)
		return
	}

	err = data.Validate(false)
	if err != nil {
		msg := generateErrorMessage(http.StatusBadRequest, err.Error())
		reportError(w, startTime, "400", mon.CreateAdminUserStr, msg, http.StatusBadRequest)
		return
	}

	logger.Log.Infof("Creating %s: %s", tenmod.TenantIngestionProfileStr, models.AsJSONString(&data))

	// Issue request to DAO Layer
	result, err := tsh.TenantDB.CreateTenantIngestionProfile(&data)
	if err != nil {
		msg := fmt.Sprintf("Unable to store %s: %s", tenmod.TenantIngestionProfileStr, err.Error())
		reportError(w, startTime, "500", mon.CreateIngPrfStr, msg, http.StatusInternalServerError)
		return
	}

	sendSuccessResponse(result, w, startTime, mon.CreateIngPrfStr, tenmod.TenantIngestionProfileStr, "Created")
}

//PatchTenantIngestionProfile - updates a tenant ingestion profile
func (tsh *TenantServiceRESTHandler) PatchTenantIngestionProfile(w http.ResponseWriter, r *http.Request) {
	startTime := time.Now()
	opStr := mon.PatchIngPrfStr

	// Unmarshal the request
	requestBytes, err := getRequestBytes(r)
	if err != nil {
		msg := generateErrorMessage(http.StatusBadRequest, err.Error())
		reportError(w, startTime, "400", opStr, msg, http.StatusBadRequest)
		return
	}

	data := tenmod.IngestionProfile{}
	err = jsonapi.Unmarshal(requestBytes, &data)
	if err != nil {
		msg := generateErrorMessage(http.StatusBadRequest, err.Error())
		reportError(w, startTime, "400", opStr, msg, http.StatusBadRequest)
		return
	}

	err = data.Validate(true)
	if err != nil {
		msg := generateErrorMessage(http.StatusBadRequest, err.Error())
		reportError(w, startTime, "400", opStr, msg, http.StatusBadRequest)
		return
	}

	origData, err2 := tsh.TenantDB.GetTenantIngestionProfile(data.TenantID, data.ID)
	if err2 != nil {
		msg := generateErrorMessage(http.StatusBadRequest, err.Error())
		reportError(w, startTime, "400", opStr, msg, http.StatusBadRequest)
		return
	}

	errMerge := models.MergeObjWithMap(origData, requestBytes)
	if errMerge != nil {
		msg := generateErrorMessage(http.StatusBadRequest, err.Error())
		reportError(w, startTime, "400", opStr, msg, http.StatusBadRequest)
		return
	}

	logger.Log.Infof("Patching%s: %s", tenmod.TenantIngestionProfileStr, origData)

	// Issue request to DAO Layer
	result, err := tsh.TenantDB.UpdateTenantIngestionProfile(origData)
	if err != nil {
		msg := fmt.Sprintf("Unable to store %s: %s", tenmod.TenantIngestionProfileStr, err.Error())
		reportError(w, startTime, "500", opStr, msg, http.StatusInternalServerError)
		return
	}

	sendSuccessResponse(result, w, startTime, opStr, tenmod.TenantIngestionProfileStr, "Patched")
}

// UpdateTenantIngestionProfile - updates a tenant ingestion profile
func (tsh *TenantServiceRESTHandler) UpdateTenantIngestionProfile(w http.ResponseWriter, r *http.Request) {
	startTime := time.Now()

	// Unmarshal the request
	requestBytes, err := getRequestBytes(r)
	if err != nil {
		msg := generateErrorMessage(http.StatusBadRequest, err.Error())
		reportError(w, startTime, "400", mon.CreateAdminUserStr, msg, http.StatusBadRequest)
		return
	}

	data := tenmod.IngestionProfile{}
	err = jsonapi.Unmarshal(requestBytes, &data)
	if err != nil {
		msg := generateErrorMessage(http.StatusBadRequest, err.Error())
		reportError(w, startTime, "400", mon.CreateAdminUserStr, msg, http.StatusBadRequest)
		return
	}

	err = data.Validate(true)
	if err != nil {
		msg := generateErrorMessage(http.StatusBadRequest, err.Error())
		reportError(w, startTime, "400", mon.CreateAdminUserStr, msg, http.StatusBadRequest)
		return
	}

	logger.Log.Infof("Updating %s: %s", tenmod.TenantIngestionProfileStr, models.AsJSONString(&data))

	// Issue request to DAO Layer
	result, err := tsh.TenantDB.UpdateTenantIngestionProfile(&data)
	if err != nil {
		msg := fmt.Sprintf("Unable to store %s: %s", tenmod.TenantIngestionProfileStr, err.Error())
		reportError(w, startTime, "500", mon.UpdateIngPrfStr, msg, http.StatusInternalServerError)
		return
	}

	sendSuccessResponse(result, w, startTime, mon.UpdateIngPrfStr, tenmod.TenantIngestionProfileStr, "Updated")
}

// GetTenantIngestionProfile - fetches a tenant ingestion profile
func (tsh *TenantServiceRESTHandler) GetTenantIngestionProfile(w http.ResponseWriter, r *http.Request) {
	startTime := time.Now()

	// Get the IDs from the URL
	tenantID := getDBFieldFromRequest(r, 4)
	dataID := getDBFieldFromRequest(r, 6)

	logger.Log.Infof("Fetching %s: %s", tenmod.TenantIngestionProfileStr, dataID)

	// Issue request to DAO Layer
	result, err := tsh.TenantDB.GetTenantIngestionProfile(tenantID, dataID)
	if err != nil {
		msg := fmt.Sprintf("Unable to retrieve %s: %s", tenmod.TenantIngestionProfileStr, err.Error())
		reportError(w, startTime, "500", mon.GetIngPrfStr, msg, http.StatusInternalServerError)
		return
	}

	sendSuccessResponse(result, w, startTime, mon.GetIngPrfStr, tenmod.TenantIngestionProfileStr, "Retrieved")
}

// DeleteTenantIngestionProfile - deletes a tenant ingestion profile
func (tsh *TenantServiceRESTHandler) DeleteTenantIngestionProfile(w http.ResponseWriter, r *http.Request) {
	startTime := time.Now()

	// Get the IDs from the URL
	tenantID := getDBFieldFromRequest(r, 4)
	dataID := getDBFieldFromRequest(r, 6)

	logger.Log.Infof("Deleting %s: %s", tenmod.TenantIngestionProfileStr, dataID)

	// Issue request to DAO Layer
	result, err := tsh.TenantDB.DeleteTenantIngestionProfile(tenantID, dataID)
	if err != nil {
		msg := fmt.Sprintf("Unable to retrieve %s: %s", tenmod.TenantIngestionProfileStr, err.Error())
		reportError(w, startTime, "500", mon.DeleteIngPrfStr, msg, http.StatusInternalServerError)
		return
	}

	sendSuccessResponse(result, w, startTime, mon.DeleteIngPrfStr, tenmod.TenantIngestionProfileStr, "Deleted")
}

// GetActiveTenantIngestionProfile - fetches the active tenant ingestion profile
func (tsh *TenantServiceRESTHandler) GetActiveTenantIngestionProfile(w http.ResponseWriter, r *http.Request) {
	startTime := time.Now()

	// Get the IDs from the URL
	tenantID := getDBFieldFromRequest(r, 4)

	logger.Log.Infof("Fetching active %s for Tenant %s", tenmod.TenantIngestionProfileStr, tenantID)

	// Issue request to DAO Layer
	result, err := tsh.TenantDB.GetActiveTenantIngestionProfile(tenantID)
	if err != nil {
		msg := fmt.Sprintf("Unable to retrieve %s: %s", tenmod.TenantIngestionProfileStr, err.Error())
		reportError(w, startTime, "500", mon.GetActiveIngPrfStr, msg, http.StatusInternalServerError)
		return
	}

	sendSuccessResponse(result, w, startTime, mon.GetActiveIngPrfStr, tenmod.TenantIngestionProfileStr, "Retrieved")
}

// CreateTenantThresholdProfile - creates a tenant threshold profile
func (tsh *TenantServiceRESTHandler) CreateTenantThresholdProfile(w http.ResponseWriter, r *http.Request) {
	startTime := time.Now()

	// Unmarshal the request
	requestBytes, err := getRequestBytes(r)
	if err != nil {
		msg := generateErrorMessage(http.StatusBadRequest, err.Error())
		reportError(w, startTime, "400", mon.CreateAdminUserStr, msg, http.StatusBadRequest)
		return
	}

	data := tenmod.ThresholdProfile{}
	err = jsonapi.Unmarshal(requestBytes, &data)
	if err != nil {
		msg := generateErrorMessage(http.StatusBadRequest, err.Error())
		reportError(w, startTime, "400", mon.CreateAdminUserStr, msg, http.StatusBadRequest)
		return
	}

	err = data.Validate(false)
	if err != nil {
		msg := generateErrorMessage(http.StatusBadRequest, err.Error())
		reportError(w, startTime, "400", mon.CreateAdminUserStr, msg, http.StatusBadRequest)
		return
	}

	logger.Log.Infof("Creating %s: %s", tenmod.TenantThresholdProfileStr, models.AsJSONString(&data))

	// Issue request to DAO Layer
	result, err := tsh.TenantDB.CreateTenantThresholdProfile(&data)
	if err != nil {
		msg := fmt.Sprintf("Unable to store %s: %s", tenmod.TenantThresholdProfileStr, err.Error())
		reportError(w, startTime, "500", mon.CreateThrPrfStr, msg, http.StatusInternalServerError)
		return
	}

	sendSuccessResponse(result, w, startTime, mon.CreateThrPrfStr, tenmod.TenantThresholdProfileStr, "Created")
}

// UpdateTenantThresholdProfile - updates a tenant threshold profile
func (tsh *TenantServiceRESTHandler) PatchTenantThresholdProfile(w http.ResponseWriter, r *http.Request) {
	startTime := time.Now()
	opStr := mon.UpdateThrPrfStr

	// Unmarshal the request
	requestBytes, err := getRequestBytes(r)
	if err != nil {
		msg := generateErrorMessage(http.StatusBadRequest, err.Error())
		reportError(w, startTime, "400", opStr, msg, http.StatusBadRequest)
		return
	}

	data := tenmod.ThresholdProfile{}
	err = jsonapi.Unmarshal(requestBytes, &data)
	if err != nil {
		msg := generateErrorMessage(http.StatusBadRequest, err.Error())
		reportError(w, startTime, "400", opStr, msg, http.StatusBadRequest)
		return
	}

	err = data.Validate(true)
	if err != nil {
		msg := generateErrorMessage(http.StatusBadRequest, err.Error())
		reportError(w, startTime, "400", opStr, msg, http.StatusBadRequest)
		return
	}

	origData, err2 := tsh.TenantDB.GetTenantThresholdProfile(data.TenantID, data.ID)
	if err2 != nil {
		msg := generateErrorMessage(http.StatusBadRequest, err.Error())
		reportError(w, startTime, "400", opStr, msg, http.StatusBadRequest)
		return
	}

	errMerge := models.MergeObjWithMap(origData, requestBytes)
	if errMerge != nil {
		msg := generateErrorMessage(http.StatusBadRequest, err.Error())
		reportError(w, startTime, "400", opStr, msg, http.StatusBadRequest)
		return
	}

	logger.Log.Infof("Updating %s: %s", tenmod.TenantThresholdProfileStr, origData)

	// Issue request to DAO Layer
	result, err := tsh.TenantDB.UpdateTenantThresholdProfile(origData)
	if err != nil {
		msg := fmt.Sprintf("Unable to store %s: %s", tenmod.TenantThresholdProfileStr, err.Error())
		reportError(w, startTime, "500", opStr, msg, http.StatusInternalServerError)
		return
	}

	sendSuccessResponse(result, w, startTime, opStr, tenmod.TenantThresholdProfileStr, "Patched")
}

// UpdateTenantThresholdProfile - updates a tenant threshold profile
func (tsh *TenantServiceRESTHandler) UpdateTenantThresholdProfile(w http.ResponseWriter, r *http.Request) {
	startTime := time.Now()

	// Unmarshal the request
	requestBytes, err := getRequestBytes(r)
	if err != nil {
		msg := generateErrorMessage(http.StatusBadRequest, err.Error())
		reportError(w, startTime, "400", mon.CreateAdminUserStr, msg, http.StatusBadRequest)
		return
	}

	data := tenmod.ThresholdProfile{}
	err = jsonapi.Unmarshal(requestBytes, &data)
	if err != nil {
		msg := generateErrorMessage(http.StatusBadRequest, err.Error())
		reportError(w, startTime, "400", mon.CreateAdminUserStr, msg, http.StatusBadRequest)
		return
	}

	err = data.Validate(true)
	if err != nil {
		msg := generateErrorMessage(http.StatusBadRequest, err.Error())
		reportError(w, startTime, "400", mon.CreateAdminUserStr, msg, http.StatusBadRequest)
		return
	}

	logger.Log.Infof("Updating %s: %s", tenmod.TenantThresholdProfileStr, models.AsJSONString(&data))

	// Issue request to DAO Layer
	result, err := tsh.TenantDB.UpdateTenantThresholdProfile(&data)
	if err != nil {
		msg := fmt.Sprintf("Unable to store %s: %s", tenmod.TenantThresholdProfileStr, err.Error())
		reportError(w, startTime, "500", mon.UpdateThrPrfStr, msg, http.StatusInternalServerError)
		return
	}

	sendSuccessResponse(result, w, startTime, mon.UpdateThrPrfStr, tenmod.TenantThresholdProfileStr, "Updated")
}

// GetTenantThresholdProfile - fetches a tenant threshold profile
func (tsh *TenantServiceRESTHandler) GetTenantThresholdProfile(w http.ResponseWriter, r *http.Request) {
	startTime := time.Now()

	// Get the IDs from the URL
	tenantID := getDBFieldFromRequest(r, 4)
	dataID := getDBFieldFromRequest(r, 6)

	logger.Log.Infof("Fetching %s: %s", tenmod.TenantThresholdProfileStr, dataID)

	// Issue request to DAO Layer
	result, err := tsh.TenantDB.GetTenantThresholdProfile(tenantID, dataID)
	if err != nil {
		msg := fmt.Sprintf("Unable to retrieve %s: %s", tenmod.TenantThresholdProfileStr, err.Error())
		reportError(w, startTime, "500", mon.GetThrPrfStr, msg, http.StatusInternalServerError)
		return
	}

	sendSuccessResponse(result, w, startTime, mon.GetThrPrfStr, tenmod.TenantThresholdProfileStr, "Retrieved")
}

// DeleteTenantThresholdProfile - deletes a tenant threshold profile
func (tsh *TenantServiceRESTHandler) DeleteTenantThresholdProfile(w http.ResponseWriter, r *http.Request) {
	startTime := time.Now()

	// Get the IDs from the URL
	tenantID := getDBFieldFromRequest(r, 4)
	dataID := getDBFieldFromRequest(r, 6)

	logger.Log.Infof("Deleting %s: %s", tenmod.TenantThresholdProfileStr, dataID)

	// Integrity Check - SLA Reports
	configs, err := tsh.TenantDB.GetAllReportScheduleConfigs(tenantID)
	if err != nil {
		msg := fmt.Sprintf("Unable to perform integrity check for %s deletion: %s", tenmod.TenantThresholdProfileStr, err.Error())
		reportError(w, startTime, "500", mon.DeleteThrPrfStr, msg, http.StatusInternalServerError)
		return
	}
	for _, rep := range configs {
		if rep.ThresholdProfile == dataID {
			msg := fmt.Sprintf("%s deletion failed integrity check: in use by at least one %s", tenmod.TenantThresholdProfileStr, tenmod.TenantReportScheduleConfigStr)
			reportError(w, startTime, "500", mon.DeleteThrPrfStr, msg, http.StatusInternalServerError)
			return
		}
	}

	// Issue request to DAO Layer
	result, err := tsh.TenantDB.DeleteTenantThresholdProfile(tenantID, dataID)
	if err != nil {
		msg := fmt.Sprintf("Unable to retrieve %s: %s", tenmod.TenantThresholdProfileStr, err.Error())
		reportError(w, startTime, "500", mon.DeleteThrPrfStr, msg, http.StatusInternalServerError)
		return
	}

	sendSuccessResponse(result, w, startTime, mon.DeleteThrPrfStr, tenmod.TenantThresholdProfileStr, "Deleted")
}

// GetAllTenantThresholdProfiles - fetches list of tenant threshold profiles
func (tsh *TenantServiceRESTHandler) GetAllTenantThresholdProfiles(w http.ResponseWriter, r *http.Request) {
	startTime := time.Now()

	tenantID := getDBFieldFromRequest(r, 4)

	logger.Log.Infof("Fetching %s list for Tenant %s", tenmod.TenantThresholdProfileStr, tenantID)

	// Issue request to DAO Layer
	result, err := tsh.TenantDB.GetAllTenantThresholdProfile(tenantID)
	if err != nil {
		msg := fmt.Sprintf("Unable to retrieve %s list: %s", tenmod.TenantThresholdProfileStr, err.Error())
		reportError(w, startTime, "500", mon.GetAllThrPrfStr, msg, http.StatusInternalServerError)
		return
	}

	sendSuccessResponse(result, w, startTime, mon.GetAllThrPrfStr, tenmod.TenantThresholdProfileStr, "Retrieved list of")
}

// CreateMonitoredObject - creates a tenant monitored object
func (tsh *TenantServiceRESTHandler) CreateMonitoredObject(w http.ResponseWriter, r *http.Request) {
	startTime := time.Now()

	// Unmarshal the request
	requestBytes, err := getRequestBytes(r)
	if err != nil {
		msg := generateErrorMessage(http.StatusBadRequest, err.Error())
		reportError(w, startTime, "400", mon.CreateAdminUserStr, msg, http.StatusBadRequest)
		return
	}

	data := tenmod.MonitoredObject{}
	err = jsonapi.Unmarshal(requestBytes, &data)
	if err != nil {
		msg := generateErrorMessage(http.StatusBadRequest, err.Error())
		reportError(w, startTime, "400", mon.CreateAdminUserStr, msg, http.StatusBadRequest)
		return
	}

	err = data.Validate(false)
	if err != nil {
		msg := generateErrorMessage(http.StatusBadRequest, err.Error())
		reportError(w, startTime, "400", mon.CreateAdminUserStr, msg, http.StatusBadRequest)
		return
	}

	logger.Log.Infof("Creating %s: %s", tenmod.TenantMonitoredObjectStr, models.AsJSONString(&data))

	// Issue request to DAO Layer
	result, err := tsh.TenantDB.CreateMonitoredObject(&data)
	if err != nil {
		msg := fmt.Sprintf("Unable to store %s: %s", tenmod.TenantMonitoredObjectStr, err.Error())
		reportError(w, startTime, "500", mon.CreateMonObjStr, msg, http.StatusInternalServerError)
		return
	}

	NotifyMonitoredObjectCreated(data.TenantID, &data)
	sendSuccessResponse(result, w, startTime, mon.CreateMonObjStr, tenmod.TenantMonitoredObjectStr, "Created")
}

//PatchUpdateMonitoredObject - Patch a tenant monitored object
func (tsh *TenantServiceRESTHandler) PatchMonitoredObject(w http.ResponseWriter, r *http.Request) {
	startTime := time.Now()
	opStr := mon.PatchMonObjStr

	// Unmarshal the request
	requestBytes, err := getRequestBytes(r)
	if err != nil {
		msg := generateErrorMessage(http.StatusBadRequest, err.Error())
		reportError(w, startTime, "400", opStr, msg, http.StatusBadRequest)
		return
	}

	data := tenmod.MonitoredObject{}
	err = jsonapi.Unmarshal(requestBytes, &data)
	if err != nil {
		msg := generateErrorMessage(http.StatusBadRequest, err.Error())
		reportError(w, startTime, "400", opStr, msg, http.StatusBadRequest)
		return
	}

	err = data.Validate(true)
	if err != nil {
		msg := generateErrorMessage(http.StatusBadRequest, err.Error())
		reportError(w, startTime, "400", opStr, msg, http.StatusBadRequest)
		return
	}

	// Issue request to DAO Layer
	oldData, err := tsh.TenantDB.GetMonitoredObject(data.TenantID, data.ID)
	if err != nil {
		msg := fmt.Sprintf("Unable to retrieve %s: %s", tenmod.TenantMonitoredObjectStr, err.Error())
		reportError(w, startTime, "500", opStr, msg, http.StatusInternalServerError)
		return
	}

	errMerge := models.MergeObjWithMap(oldData, requestBytes)
	if errMerge != nil {
		msg := generateErrorMessage(http.StatusBadRequest, err.Error())
		reportError(w, startTime, "400", opStr, msg, http.StatusBadRequest)
		return

	}
	// This only checks if the ID&REV is set.
	err = oldData.Validate(true)
	if err != nil {
		msg := generateErrorMessage(http.StatusBadRequest, err.Error())
		reportError(w, startTime, "400", opStr, msg, http.StatusBadRequest)
		return
	}

	logger.Log.Infof("Patching %s: %s", tenmod.TenantMonitoredObjectStr, oldData)

	// Issue request to DAO Layer
	result, err := tsh.TenantDB.UpdateMonitoredObject(oldData)
	if err != nil {
		msg := fmt.Sprintf("Unable to store %s: %s", tenmod.TenantMonitoredObjectStr, err.Error())
		reportError(w, startTime, "500", opStr, msg, http.StatusInternalServerError)
		return
	}

	NotifyMonitoredObjectUpdated(data.TenantID, oldData)
	sendSuccessResponse(result, w, startTime, opStr, tenmod.TenantMonitoredObjectStr, "Patched")
}

// UpdateMonitoredObject - updates a tenant monitored object
func (tsh *TenantServiceRESTHandler) UpdateMonitoredObject(w http.ResponseWriter, r *http.Request) {
	startTime := time.Now()

	// Unmarshal the request
	requestBytes, err := getRequestBytes(r)
	if err != nil {
		msg := generateErrorMessage(http.StatusBadRequest, err.Error())
		reportError(w, startTime, "400", mon.CreateAdminUserStr, msg, http.StatusBadRequest)
		return
	}

	data := tenmod.MonitoredObject{}
	err = jsonapi.Unmarshal(requestBytes, &data)
	if err != nil {
		msg := generateErrorMessage(http.StatusBadRequest, err.Error())
		reportError(w, startTime, "400", mon.CreateAdminUserStr, msg, http.StatusBadRequest)
		return
	}

	err = data.Validate(true)
	if err != nil {
		msg := generateErrorMessage(http.StatusBadRequest, err.Error())
		reportError(w, startTime, "400", mon.CreateAdminUserStr, msg, http.StatusBadRequest)
		return
	}

	logger.Log.Infof("Updating %s: %s", tenmod.TenantMonitoredObjectStr, models.AsJSONString(&data))

	// Issue request to DAO Layer
	result, err := tsh.TenantDB.UpdateMonitoredObject(&data)
	if err != nil {
		msg := fmt.Sprintf("Unable to store %s: %s", tenmod.TenantMonitoredObjectStr, err.Error())
		reportError(w, startTime, "500", mon.UpdateMonObjStr, msg, http.StatusInternalServerError)
		return
	}

	NotifyMonitoredObjectUpdated(data.TenantID, &data)
	sendSuccessResponse(result, w, startTime, mon.UpdateMonObjStr, tenmod.TenantMonitoredObjectStr, "Updated")
}

// GetMonitoredObject - fetches a tenant monitored object
func (tsh *TenantServiceRESTHandler) GetMonitoredObject(w http.ResponseWriter, r *http.Request) {
	startTime := time.Now()

	// Get the IDs from the URL
	tenantID := getDBFieldFromRequest(r, 4)
	dataID := getDBFieldFromRequest(r, 6)

	logger.Log.Infof("Fetching %s: %s", tenmod.TenantMonitoredObjectStr, dataID)

	// Issue request to DAO Layer
	result, err := tsh.TenantDB.GetMonitoredObject(tenantID, dataID)
	if err != nil {
		msg := fmt.Sprintf("Unable to retrieve %s: %s", tenmod.TenantMonitoredObjectStr, err.Error())
		reportError(w, startTime, "500", mon.GetMonObjStr, msg, http.StatusInternalServerError)
		return
	}

	sendSuccessResponse(result, w, startTime, mon.GetMonObjStr, tenmod.TenantMonitoredObjectStr, "Retrieved")
}

// DeleteMonitoredObject - deletes a tenant monitored object
func (tsh *TenantServiceRESTHandler) DeleteMonitoredObject(w http.ResponseWriter, r *http.Request) {
	startTime := time.Now()

	// Get the IDs from the URL
	tenantID := getDBFieldFromRequest(r, 4)
	dataID := getDBFieldFromRequest(r, 6)

	logger.Log.Infof("Deleting %s: %s", tenmod.TenantMonitoredObjectStr, dataID)

	// Issue request to DAO Layer
	result, err := tsh.TenantDB.DeleteMonitoredObject(tenantID, dataID)
	if err != nil {
		msg := fmt.Sprintf("Unable to retrieve %s: %s", tenmod.TenantMonitoredObjectStr, err.Error())
		reportError(w, startTime, "500", mon.DeleteMonObjStr, msg, http.StatusInternalServerError)
		return
	}

	sendSuccessResponse(result, w, startTime, mon.DeleteMonObjStr, tenmod.TenantMonitoredObjectStr, "Deleted")
}

// GetAllMonitoredObjects - fetches list of tenant monitored objects
func (tsh *TenantServiceRESTHandler) GetAllMonitoredObjects(w http.ResponseWriter, r *http.Request) {
	startTime := time.Now()

	tenantID := getDBFieldFromRequest(r, 4)

	logger.Log.Infof("Fetching %s list for Tenant %s", tenmod.TenantMonitoredObjectStr, tenantID)

	// Issue request to DAO Layer
	result, err := tsh.TenantDB.GetAllMonitoredObjects(tenantID)
	if err != nil {
		msg := fmt.Sprintf("Unable to retrieve %s list: %s", tenmod.TenantMonitoredObjectStr, err.Error())
		reportError(w, startTime, "500", mon.GetAllMonObjStr, msg, http.StatusInternalServerError)
		return
	}

	sendSuccessResponse(result, w, startTime, mon.GetAllMonObjStr, tenmod.TenantMonitoredObjectStr, "Retrieved list of")
}

func (tsh *TenantServiceRESTHandler) GetMonitoredObjectToDomainMap(w http.ResponseWriter, r *http.Request) {

	startTime := time.Now()

	tenantID := getDBFieldFromRequest(r, 4)

	// Unmarshal the request
	data := tenmod.MonitoredObjectCountByDomainRequest{}
	data.TenantID = tenantID
	err := unmarshalRequest(r, &data, true)
	if err != nil {
		msg := generateErrorMessage(http.StatusBadRequest, err.Error())
		reportError(w, startTime, "400", mon.UpdateMonObjStr, msg, http.StatusBadRequest)
		return
	}

	logger.Log.Infof("Fetching %s for Tenant %s", tenmod.MonitoredObjectToDomainMapStr, data.TenantID)

	// Issue request to DAO Layer
	result, err := tsh.TenantDB.GetMonitoredObjectToDomainMap(&data)
	if err != nil {
		msg := fmt.Sprintf("Unable to retrieve %s : %s", tenmod.MonitoredObjectToDomainMapStr, err.Error())
		reportError(w, startTime, "500", mon.GetMonObjToDomMapStr, msg, http.StatusInternalServerError)
		return
	}

	res, err := json.Marshal(result)
	if err != nil {
		msg := fmt.Sprintf("Unable to marshal response for  %s : %s", tenmod.MonitoredObjectToDomainMapStr, err.Error())
		reportError(w, startTime, "500", mon.GetMonObjToDomMapStr, msg, http.StatusInternalServerError)
		return
	}

	logger.Log.Infof("Successfully retrieved %s for Tenant %s", tenmod.MonitoredObjectToDomainMapStr, data.TenantID)
	trackAPIMetrics(startTime, "200", mon.GetMonObjToDomMapStr)
	fmt.Fprintf(w, string(res))
}

// CreateTenantMeta - creates a tenant metadata
func (tsh *TenantServiceRESTHandler) CreateTenantMeta(w http.ResponseWriter, r *http.Request) {
	startTime := time.Now()

	// Unmarshal the request
	requestBytes, err := getRequestBytes(r)
	if err != nil {
		msg := generateErrorMessage(http.StatusBadRequest, err.Error())
		reportError(w, startTime, "400", mon.CreateAdminUserStr, msg, http.StatusBadRequest)
		return
	}

	data := tenmod.Metadata{}
	err = jsonapi.Unmarshal(requestBytes, &data)
	if err != nil {
		msg := generateErrorMessage(http.StatusBadRequest, err.Error())
		reportError(w, startTime, "400", mon.CreateAdminUserStr, msg, http.StatusBadRequest)
		return
	}

	err = data.Validate(false)
	if err != nil {
		msg := generateErrorMessage(http.StatusBadRequest, err.Error())
		reportError(w, startTime, "400", mon.CreateAdminUserStr, msg, http.StatusBadRequest)
		return
	}

	logger.Log.Infof("Creating %s: %s", tenmod.TenantMetaStr, models.AsJSONString(&data))

	// Issue request to DAO Layer
	result, err := tsh.TenantDB.CreateTenantMeta(&data)
	if err != nil {
		msg := fmt.Sprintf("Unable to store %s: %s", tenmod.TenantMetaStr, err.Error())
		reportError(w, startTime, "500", mon.CreateTenantMetaStr, msg, http.StatusInternalServerError)
		return
	}

	sendSuccessResponse(result, w, startTime, mon.CreateTenantMetaStr, tenmod.TenantMetaStr, "Created")
}

//PatchTenantMeta - Patch a tenant metadata
func (tsh *TenantServiceRESTHandler) PatchTenantMeta(w http.ResponseWriter, r *http.Request) {
	startTime := time.Now()
	opStr := mon.PatchTenantMetaStr

	// Unmarshal the request
	requestBytes, err := getRequestBytes(r)
	if err != nil {
		msg := generateErrorMessage(http.StatusBadRequest, err.Error())
		reportError(w, startTime, "400", opStr, msg, http.StatusBadRequest)
		return
	}

	data := tenmod.Metadata{}
	err = jsonapi.Unmarshal(requestBytes, &data)
	if err != nil {
		msg := generateErrorMessage(http.StatusBadRequest, err.Error())
		reportError(w, startTime, "400", opStr, msg, http.StatusBadRequest)
		return
	}

	err = data.Validate(true)
	if err != nil {
		msg := generateErrorMessage(http.StatusBadRequest, err.Error())
		reportError(w, startTime, "400", opStr, msg, http.StatusBadRequest)
		return
	}

	// Issue request to DAO Layer
	oldData, err := tsh.TenantDB.GetTenantMeta(data.TenantID)
	if err != nil {
		msg := fmt.Sprintf("Unable to retrieve %s: %s", tenmod.TenantMetaStr, err.Error())
		reportError(w, startTime, "500", opStr, msg, http.StatusInternalServerError)
		return
	}

	errMerge := models.MergeObjWithMap(oldData, requestBytes)
	if errMerge != nil {
		msg := generateErrorMessage(http.StatusBadRequest, err.Error())
		reportError(w, startTime, "400", opStr, msg, http.StatusBadRequest)
		return
	}

	err = data.Validate(true)
	if err != nil {
		msg := generateErrorMessage(http.StatusBadRequest, err.Error())
		reportError(w, startTime, "400", opStr, msg, http.StatusBadRequest)
		return
	}

	logger.Log.Infof("Patching %s: %s", tenmod.TenantMetaStr, oldData)

	// Issue request to DAO Layer
	result, err := tsh.TenantDB.UpdateTenantMeta(oldData)
	if err != nil {
		msg := fmt.Sprintf("Unable to store %s: %s", tenmod.TenantMetaStr, err.Error())
		reportError(w, startTime, "500", opStr, msg, http.StatusInternalServerError)
		return
	}

	sendSuccessResponse(result, w, startTime, opStr, tenmod.TenantMetaStr, "Patched")
}

// UpdateTenantMeta - updates a tenant metadata
func (tsh *TenantServiceRESTHandler) UpdateTenantMeta(w http.ResponseWriter, r *http.Request) {
	startTime := time.Now()

	// Unmarshal the request
	requestBytes, err := getRequestBytes(r)
	if err != nil {
		msg := generateErrorMessage(http.StatusBadRequest, err.Error())
		reportError(w, startTime, "400", mon.CreateAdminUserStr, msg, http.StatusBadRequest)
		return
	}

	data := tenmod.Metadata{}
	err = jsonapi.Unmarshal(requestBytes, &data)
	if err != nil {
		msg := generateErrorMessage(http.StatusBadRequest, err.Error())
		reportError(w, startTime, "400", mon.CreateAdminUserStr, msg, http.StatusBadRequest)
		return
	}

	err = data.Validate(true)
	if err != nil {
		msg := generateErrorMessage(http.StatusBadRequest, err.Error())
		reportError(w, startTime, "400", mon.CreateAdminUserStr, msg, http.StatusBadRequest)
		return
	}

	logger.Log.Infof("Updating %s: %s", tenmod.TenantMetaStr, models.AsJSONString(&data))

	// Issue request to DAO Layer
	result, err := tsh.TenantDB.UpdateTenantMeta(&data)
	if err != nil {
		msg := fmt.Sprintf("Unable to store %s: %s", tenmod.TenantMetaStr, err.Error())
		reportError(w, startTime, "500", mon.UpdateTenantMetaStr, msg, http.StatusInternalServerError)
		return
	}

	sendSuccessResponse(result, w, startTime, mon.UpdateTenantMetaStr, tenmod.TenantMetaStr, "Updated")
}

// GetTenantMeta - fetches a tenant metadata
func (tsh *TenantServiceRESTHandler) GetTenantMeta(w http.ResponseWriter, r *http.Request) {
	startTime := time.Now()

	// Get the IDs from the URL
	tenantID := getDBFieldFromRequest(r, 4)

	logger.Log.Infof("Fetching %s: %s", tenmod.TenantMetaStr)

	// Issue request to DAO Layer
	result, err := tsh.TenantDB.GetTenantMeta(tenantID)
	if err != nil {
		msg := fmt.Sprintf("Unable to retrieve %s: %s", tenmod.TenantMetaStr, err.Error())
		reportError(w, startTime, "500", mon.GetTenantMetaStr, msg, http.StatusInternalServerError)
		return
	}

	sendSuccessResponse(result, w, startTime, mon.GetTenantMetaStr, tenmod.TenantMetaStr, "Retrieved")
}

// DeleteTenantMeta - deletes a tenant metadata
func (tsh *TenantServiceRESTHandler) DeleteTenantMeta(w http.ResponseWriter, r *http.Request) {
	startTime := time.Now()

	// Get the IDs from the URL
	tenantID := getDBFieldFromRequest(r, 4)

	logger.Log.Infof("Deleting %s: %s", tenmod.TenantMetaStr)

	// Issue request to DAO Layer
	result, err := tsh.TenantDB.DeleteTenantMeta(tenantID)
	if err != nil {
		msg := fmt.Sprintf("Unable to retrieve %s: %s", tenmod.TenantMetaStr, err.Error())
		reportError(w, startTime, "500", mon.DeleteTenantMetaStr, msg, http.StatusInternalServerError)
		return
	}

	sendSuccessResponse(result, w, startTime, mon.DeleteTenantMetaStr, tenmod.TenantMetaStr, "Deleted")
}

// BulkInsertMonitoredObject - creates 1 or many monitored objects in one request
func (tsh *TenantServiceRESTHandler) BulkInsertMonitoredObject(w http.ResponseWriter, r *http.Request) {
	startTime := time.Now()

	tenantID := getDBFieldFromRequest(r, 4)

	// Unmarshal the request
	data := []*tenmod.MonitoredObject{}
	err := unmarshalData(r, &data)
	if err != nil {
		msg := generateErrorMessage(http.StatusBadRequest, err.Error())
		reportError(w, startTime, "400", mon.BulkInsertMonObjStr, msg, http.StatusBadRequest)
		return
	}

	// Validate the request data
	for _, obj := range data {
		if err = obj.Validate(false); err != nil {

		}

		if obj.TenantID != tenantID {
			msg := fmt.Sprintf("Unable to Update %ss in bulk: %s", tenmod.TenantMonitoredObjectStr, "All Monitored Objects must have Tenant ID "+tenantID)
			reportError(w, startTime, "400", mon.BulkInsertMonObjStr, msg, http.StatusBadRequest)
			return
		}
	}

	logger.Log.Infof("Bulk inserting %ss: %s", tenmod.TenantMonitoredObjectStr, models.AsJSONString(&data))

	// Issue request to DAO Layer
	result, err := tsh.TenantDB.BulkInsertMonitoredObjects(tenantID, data)
	if err != nil {
		msg := fmt.Sprintf("Unable to store %s: %s", tenmod.TenantMonitoredObjectStr, err.Error())
		reportError(w, startTime, "500", mon.BulkUpdateMonObjStr, msg, http.StatusInternalServerError)
		return
	}

	NotifyMonitoredObjectCreated(tenantID, data...)
	response := map[string]interface{}{}
	response["results"] = result

	res, err := json.Marshal(response)
	if err != nil {
		msg := fmt.Sprintf("Unable to marshal response for  %s : %s", tenmod.TenantMonitoredObjectStr, err.Error())
		reportError(w, startTime, "500", mon.BulkUpdateMonObjStr, msg, http.StatusInternalServerError)
		return
	}

	logger.Log.Infof("Completed bulk insert of %ss for Tenant %s", tenmod.TenantMonitoredObjectStr, tenantID)
	trackAPIMetrics(startTime, "200", mon.BulkUpdateMonObjStr)
	fmt.Fprintf(w, string(res))
}

// BulkUpdateMonitoredObject - updates 1 or many monitored objects in one request
func (tsh *TenantServiceRESTHandler) BulkUpdateMonitoredObject(w http.ResponseWriter, r *http.Request) {
	startTime := time.Now()

	tenantID := getDBFieldFromRequest(r, 4)

	// Unmarshal the request
	data := []*tenmod.MonitoredObject{}
	err := unmarshalData(r, &data)
	if err != nil {
		msg := generateErrorMessage(http.StatusBadRequest, err.Error())
		reportError(w, startTime, "400", mon.BulkUpdateMonObjStr, msg, http.StatusBadRequest)
		return
	}

	// Validate the request data
	for _, obj := range data {
		if err = obj.Validate(true); err != nil || obj.TenantID != tenantID {
			msg := fmt.Sprintf("Unable to Update %ss in bulk: %s", tenmod.TenantMonitoredObjectStr, "All Monitored Objects must have ID and revision")
			reportError(w, startTime, "400", mon.BulkUpdateMonObjStr, msg, http.StatusBadRequest)
			return
		}
		if obj.TenantID != tenantID {
			msg := fmt.Sprintf("Unable to Update %ss in bulk: %s", tenmod.TenantMonitoredObjectStr, "All Monitored Objects must have Tenant ID "+tenantID)
			reportError(w, startTime, "400", mon.BulkUpdateMonObjStr, msg, http.StatusBadRequest)
			return
		}

	}

	logger.Log.Infof("Bulk updating %ss: %s", tenmod.TenantMonitoredObjectStr, models.AsJSONString(&data))

	// Issue request to DAO Layer
	result, err := tsh.TenantDB.BulkUpdateMonitoredObjects(tenantID, data)
	if err != nil {
		msg := fmt.Sprintf("Unable to store %s: %s", tenmod.TenantMonitoredObjectStr, err.Error())
		reportError(w, startTime, "500", mon.BulkUpdateMonObjStr, msg, http.StatusInternalServerError)
		return
	}

	NotifyMonitoredObjectUpdated(tenantID, data...)
	response := map[string]interface{}{}
	response["results"] = result

	res, err := json.Marshal(response)
	if err != nil {
		msg := fmt.Sprintf("Unable to marshal response for  %s : %s", tenmod.TenantMonitoredObjectStr, err.Error())
		reportError(w, startTime, "500", mon.BulkUpdateMonObjStr, msg, http.StatusInternalServerError)
		return
	}

	logger.Log.Infof("Completed bulk insert of %ss for Tenant %s", tenmod.TenantMonitoredObjectStr, tenantID)
	trackAPIMetrics(startTime, "200", mon.BulkUpdateMonObjStr)
	fmt.Fprintf(w, string(res))
}

// Creates a report scheduling configuration that will ultimately prompt the scheduler to execute a report at the specified time with the specified parameters
// Refer to the model defined in scheduler.go to check what payload is acceptable to the REST endpoint
// Params:
//		w - the writer responsible for marshalling the response to the incoming http request
//		r - the initiating http request
func (tsh *TenantServiceRESTHandler) CreateReportScheduleConfig(w http.ResponseWriter, r *http.Request) {
	startTime := time.Now()

	// Ensure that we can unmarshal the provided report schedule payload into the model object
	requestBytes, err := getRequestBytes(r)
	if err != nil {
		msg := generateErrorMessage(http.StatusBadRequest, err.Error())
		reportError(w, startTime, "400", mon.CreateReportScheduleConfigStr, msg, http.StatusBadRequest)
		return
	}
	data := metmod.ReportScheduleConfig{}
	err = jsonapi.Unmarshal(requestBytes, &data)
	if err != nil {
		msg := generateErrorMessage(http.StatusBadRequest, err.Error())
		reportError(w, startTime, "400", mon.CreateReportScheduleConfigStr, msg, http.StatusBadRequest)
		return
	}
	// Seconds is not really used except for testing.
	data.Second = "0"

	// Ensure that the passed in data adheres to the model requirements
	err = data.Validate(false)
	if err != nil {
		msg := generateErrorMessage(http.StatusBadRequest, err.Error())
		reportError(w, startTime, "400", mon.CreateReportScheduleConfigStr, msg, http.StatusBadRequest)
		return
	}

	// Attempt to create a config entry in the datastore for the scheduler to pick up
	logger.Log.Infof("Creating %s: %s", metmod.ReportScheduleConfigStr, models.AsJSONString(&data))
	result, err := tsh.TenantDB.CreateReportScheduleConfig(&data)
	if err != nil {
		msg := fmt.Sprintf("Unable to store %s: %s", metmod.ReportScheduleConfigStr, err.Error())
		reportError(w, startTime, "500", mon.CreateReportScheduleConfigStr, msg, http.StatusInternalServerError)
		return
	}

	// Tell the scheduler to go and update based on the updated database
	err = scheduler.RebuildCronJobs()
	if err != nil {
		msg := fmt.Sprintf("Unable to start scheduled job%s: %s", metmod.ReportScheduleConfigStr, err.Error())
		reportError(w, startTime, "500", mon.CreateReportScheduleConfigStr, msg, http.StatusInternalServerError)
		return
	}
	sendSuccessResponse(result, w, startTime, mon.CreateReportScheduleConfigStr, metmod.ReportScheduleConfigStr, "Created")
}

// Updates a report scheduling configuration that will ultimately prompt the scheduler to execute a report at the specified time with the specified parameters
// Refer to the model defined in scheduler.go to check what payload is acceptable to the REST endpoint
// Params:
//		w - the writer responsible for marshalling the response to the incoming http request
//		r - the initiating http request
func (tsh *TenantServiceRESTHandler) UpdateReportScheduleConfig(w http.ResponseWriter, r *http.Request) {
	startTime := time.Now()

	// Ensure that we can unmarshal the provided report schedule payload into the model object
	requestBytes, err := getRequestBytes(r)
	if err != nil {
		msg := generateErrorMessage(http.StatusBadRequest, err.Error())
		reportError(w, startTime, "400", mon.UpdateReportScheduleConfigStr, msg, http.StatusBadRequest)
		return
	}
	data := metmod.ReportScheduleConfig{}
	err = jsonapi.Unmarshal(requestBytes, &data)
	if err != nil {
		msg := generateErrorMessage(http.StatusBadRequest, err.Error())
		reportError(w, startTime, "400", mon.UpdateReportScheduleConfigStr, msg, http.StatusBadRequest)
		return
	}

	// Seconds is not really used except for testing.
	data.Second = "0"

	// Ensure that the passed in data adheres to the model requirements
	err = data.Validate(true)
	if err != nil {
		msg := generateErrorMessage(http.StatusBadRequest, err.Error())
		reportError(w, startTime, "400", mon.UpdateReportScheduleConfigStr, msg, http.StatusBadRequest)
		return
	}

	// Attempt to update a config entry in the datastore for the scheduler to pick up
	logger.Log.Infof("Updating %s: %s", metmod.ReportScheduleConfigStr, models.AsJSONString(&data))
	result, err := tsh.TenantDB.UpdateReportScheduleConfig(&data)
	if err != nil {
		msg := fmt.Sprintf("Unable to store %s: %s", metmod.ReportScheduleConfigStr, err.Error())
		reportError(w, startTime, "500", mon.UpdateReportScheduleConfigStr, msg, http.StatusInternalServerError)
		return
	}

	// Tell the scheduler to go and update based on the updated database
	err = scheduler.RebuildCronJobs()
	if err != nil {
		msg := fmt.Sprintf("Unable to start scheduler %s: %s", metmod.ReportScheduleConfigStr, err.Error())
		reportError(w, startTime, "500", mon.DeleteReportScheduleConfigStr, msg, http.StatusInternalServerError)
		return
	}

	sendSuccessResponse(result, w, startTime, mon.UpdateReportScheduleConfigStr, metmod.ReportScheduleConfigStr, "Updated")
}

// Fetch a report scheduling configuration from the datastore for a particular tenant
// Params:
//		w - the writer responsible for marshalling the response to the incoming http request
//		r - the initiating http request
func (tsh *TenantServiceRESTHandler) GetReportScheduleConfig(w http.ResponseWriter, r *http.Request) {
	startTime := time.Now()

	// Retrieve the tenant ID and config ID from the URL
	tenantID := getDBFieldFromRequest(r, 4)
	configID := getDBFieldFromRequest(r, 6)

	// Attempt to fetch the config entry from the datastore
	logger.Log.Infof("Fetching %s: %s", metmod.ReportScheduleConfigStr, configID)
	result, err := tsh.TenantDB.GetReportScheduleConfig(tenantID, configID)
	if err != nil {
		msg := fmt.Sprintf("Unable to retrieve %s: %s", metmod.ReportScheduleConfigStr, err.Error())
		reportError(w, startTime, "500", mon.GetReportScheduleConfigStr, msg, http.StatusInternalServerError)
		return
	}

	sendSuccessResponse(result, w, startTime, mon.GetReportScheduleConfigStr, metmod.ReportScheduleConfigStr, "Retrieved")
}

// Fetch all report scheduling configurations from the datastore for a particular tenant
// Params:
//		w - the writer responsible for marshalling the response to the incoming http request
//		r - the initiating http request
func (tsh *TenantServiceRESTHandler) GetAllReportScheduleConfigs(w http.ResponseWriter, r *http.Request) {
	startTime := time.Now()

	// Retrieve the tenant ID and config ID from the URL
	tenantID := getDBFieldFromRequest(r, 4)

	// Attempt to fetch all the config entries from the datastore
	logger.Log.Infof("Fetching %s list for Tenant %s", metmod.ReportScheduleConfigStr, tenantID)
	result, err := tsh.TenantDB.GetAllReportScheduleConfigs(tenantID)
	if err != nil {
		msg := fmt.Sprintf("Unable to retrieve %s list: %s", metmod.ReportScheduleConfigStr, err.Error())
		reportError(w, startTime, "500", mon.GetAllReportScheduleConfigStr, msg, http.StatusInternalServerError)
		return
	}

	sendSuccessResponse(result, w, startTime, mon.GetAllReportScheduleConfigStr, metmod.ReportScheduleConfigStr, "Retrieved list of")
}

// Delete a report scheduling configuration from the datastore for a particular tenant
// Params:
//		w - the writer responsible for marshalling the response to the incoming http request
//		r - the initiating http request
func (tsh *TenantServiceRESTHandler) DeleteReportScheduleConfig(w http.ResponseWriter, r *http.Request) {
	startTime := time.Now()

	// Retrieve the tenant ID and config ID from the URL
	tenantID := getDBFieldFromRequest(r, 4)
	configID := getDBFieldFromRequest(r, 6)

	// Attempt to delete the specified configuration entry for the tenant
	logger.Log.Infof("Deleting %s: %s", metmod.ReportScheduleConfigStr, configID)
	result, err := tsh.TenantDB.DeleteReportScheduleConfig(tenantID, configID)
	if err != nil {
		msg := fmt.Sprintf("Unable to retrieve %s: %s", metmod.ReportScheduleConfigStr, err.Error())
		reportError(w, startTime, "500", mon.DeleteReportScheduleConfigStr, msg, http.StatusInternalServerError)
		return
	}

	// Tell the scheduler to go and update based on the updated database
	err = scheduler.RebuildCronJobs()
	if err != nil {
		msg := fmt.Sprintf("Unable to start scheduler %s: %s", metmod.ReportScheduleConfigStr, err.Error())
		reportError(w, startTime, "500", mon.DeleteReportScheduleConfigStr, msg, http.StatusInternalServerError)
		return
	}

	sendSuccessResponse(result, w, startTime, mon.DeleteReportScheduleConfigStr, metmod.ReportScheduleConfigStr, "Deleted")
}

// Fetch an SLA report from the datastore for a particular tenant
// Params:
//		w - the writer responsible for marshalling the response to the incoming http request
//		r - the initiating http request
func (tsh *TenantServiceRESTHandler) GetSLAReport(w http.ResponseWriter, r *http.Request) {
	startTime := time.Now()

	// Retrieve the tenant ID and config ID from the URL
	tenantID := getDBFieldFromRequest(r, 4)
	configID := getDBFieldFromRequest(r, 6)

	// Attempt to fetch the config entry from the datastore
	logger.Log.Infof("Fetching %s: %s", metmod.ReportStr, configID)
	result, err := tsh.TenantDB.GetSLAReport(tenantID, configID)
	if err != nil {
		msg := fmt.Sprintf("Unable to retrieve %s: %s", metmod.ReportStr, err.Error())
		reportError(w, startTime, "500", mon.GetSLAReportStr, msg, http.StatusInternalServerError)
		return
	}

	sendSuccessResponse(result, w, startTime, mon.GetSLAReportStr, metmod.ReportStr, "Retrieved")
}

// Fetch all SLA reports from the datastore for a particular tenant
// Params:
//		w - the writer responsible for marshalling the response to the incoming http request
//		r - the initiating http request
func (tsh *TenantServiceRESTHandler) GetAllSLAReports(w http.ResponseWriter, r *http.Request) {
	startTime := time.Now()

	// Retrieve the tenant ID and config ID from the URL
	tenantID := getDBFieldFromRequest(r, 4)

	// Attempt to fetch all the config entries from the datastore
	logger.Log.Infof("Fetching %s list for Tenant %s", metmod.ReportStr, tenantID)
	result, err := tsh.TenantDB.GetAllSLAReports(tenantID)
	if err != nil {
		msg := fmt.Sprintf("Unable to retrieve %s list: %s", metmod.ReportStr, err.Error())
		reportError(w, startTime, "500", mon.GetAllSLAReportStr, msg, http.StatusInternalServerError)
		return
	}

	sendSuccessResponse(result, w, startTime, mon.GetAllSLAReportStr, metmod.ReportStr, "Retrieved list of")
}
