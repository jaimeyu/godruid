package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	db "github.com/accedian/adh-gather/datastore"
	"github.com/accedian/adh-gather/logger"
	"github.com/accedian/adh-gather/models"
	tenmod "github.com/accedian/adh-gather/models/tenant"
	mon "github.com/accedian/adh-gather/monitoring"
	"github.com/accedian/adh-gather/server"
	"github.com/gorilla/mux"
	"github.com/manyminds/api2go/jsonapi"
)

// TenantServiceRESTHandler - handler of logic for REST calls made to the Tenant Service.
type TenantServiceRESTHandler struct {
	tenantDB db.TenantServiceDatastore
	routes   []server.Route
}

// CreateTenantServiceRESTHandler - used to create a Tenant Service REST handler which provides
// logic to serve the Admin Service REST calls
func CreateTenantServiceRESTHandler() *TenantServiceRESTHandler {
	result := new(TenantServiceRESTHandler)

	// Setup the DB implementation based on configuration
	tdb, err := getTenantServiceDatastore()
	if err != nil {
		logger.Log.Fatalf("Unable to instantiate AdminServiceRESTHandler: %s", err.Error())
	}
	result.tenantDB = tdb

	result.routes = []server.Route{
		server.Route{
			Name:        "CreateTenantUser",
			Method:      "POST",
			Pattern:     apiV1Prefix + tenantsAPIPrefix + "users",
			HandlerFunc: result.CreateTenantUser,
		},
		server.Route{
			Name:        "UpdateTenantUser",
			Method:      "PUT",
			Pattern:     apiV1Prefix + tenantsAPIPrefix + "users",
			HandlerFunc: result.UpdateTenantUser,
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
			HandlerFunc: result.DeleteTenantUser,
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
			HandlerFunc: result.CreateTenantDomain,
		},
		server.Route{
			Name:        "UpdateTenantDomain",
			Method:      "PUT",
			Pattern:     apiV1Prefix + tenantsAPIPrefix + "domains",
			HandlerFunc: result.UpdateTenantDomain,
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
			HandlerFunc: result.DeleteTenantDomain,
		},
		server.Route{
			Name:        "GetAllTenantDomains",
			Method:      "GET",
			Pattern:     apiV1Prefix + tenantsAPIPrefix + "domain-list",
			HandlerFunc: result.GetAllTenantDomains,
		},
		server.Route{
			Name:        "CreateTenantIngestionProfile",
			Method:      "POST",
			Pattern:     apiV1Prefix + tenantsAPIPrefix + "ingestion-profiles",
			HandlerFunc: result.CreateTenantIngestionProfile,
		},
		server.Route{
			Name:        "UpdateTenantIngestionProfile",
			Method:      "PUT",
			Pattern:     apiV1Prefix + tenantsAPIPrefix + "ingestion-profiles",
			HandlerFunc: result.UpdateTenantIngestionProfile,
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
			HandlerFunc: result.DeleteTenantIngestionProfile,
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
			HandlerFunc: result.CreateTenantThresholdProfile,
		},
		server.Route{
			Name:        "UpdateTenantThresholdProfile",
			Method:      "PUT",
			Pattern:     apiV1Prefix + tenantsAPIPrefix + "threshold-profiles",
			HandlerFunc: result.UpdateTenantThresholdProfile,
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
			HandlerFunc: result.DeleteTenantThresholdProfile,
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
			HandlerFunc: result.CreateMonitoredObject,
		},
		server.Route{
			Name:        "BulkInsertMonitoredObject",
			Method:      "POST",
			Pattern:     apiV1Prefix + tenantsAPIPrefix + "bulk/insert/monitored-objects",
			HandlerFunc: result.BulkInsertMonitoredObject,
		},
		server.Route{
			Name:        "UpdateMonitoredObject",
			Method:      "PUT",
			Pattern:     apiV1Prefix + tenantsAPIPrefix + "monitored-objects",
			HandlerFunc: result.UpdateMonitoredObject,
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
			HandlerFunc: result.DeleteMonitoredObject,
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
			HandlerFunc: result.CreateTenantMeta,
		},
		server.Route{
			Name:        "UpdateTenantMeta",
			Method:      "PUT",
			Pattern:     apiV1Prefix + tenantsAPIPrefix + "meta",
			HandlerFunc: result.UpdateTenantMeta,
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
			HandlerFunc: result.DeleteTenantMeta,
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

	data := []tenmod.User{}
	err = jsonapi.Unmarshal(requestBytes, &data)
	if err != nil {
		msg := generateErrorMessage(http.StatusBadRequest, err.Error())
		reportError(w, startTime, "400", mon.CreateAdminUserStr, msg, http.StatusBadRequest)
		return
	}

	err = data[0].Validate(false)
	if err != nil {
		msg := generateErrorMessage(http.StatusBadRequest, err.Error())
		reportError(w, startTime, "400", mon.CreateAdminUserStr, msg, http.StatusBadRequest)
		return
	}

	logger.Log.Infof("Creating %s: %s", tenmod.TenantUserStr, models.AsJSONString(&data))

	// Issue request to DAO Layer
	result, err := tsh.tenantDB.CreateTenantUser(&data[0])
	if err != nil {
		msg := fmt.Sprintf("Unable to store %s: %s", tenmod.TenantUserStr, err.Error())
		reportError(w, startTime, "500", mon.CreateTenantUserStr, msg, http.StatusInternalServerError)
		return
	}

	sendSuccessResponse(result, w, startTime, mon.CreateTenantUserStr, tenmod.TenantUserStr, "Created")
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

	data := []tenmod.User{}
	err = jsonapi.Unmarshal(requestBytes, &data)
	if err != nil {
		msg := generateErrorMessage(http.StatusBadRequest, err.Error())
		reportError(w, startTime, "400", mon.CreateAdminUserStr, msg, http.StatusBadRequest)
		return
	}

	err = data[0].Validate(true)
	if err != nil {
		msg := generateErrorMessage(http.StatusBadRequest, err.Error())
		reportError(w, startTime, "400", mon.CreateAdminUserStr, msg, http.StatusBadRequest)
		return
	}

	logger.Log.Infof("Updating %s: %s", tenmod.TenantUserStr, models.AsJSONString(&data))

	// Issue request to DAO Layer
	result, err := tsh.tenantDB.UpdateTenantUser(&data[0])
	if err != nil {
		msg := fmt.Sprintf("Unable to store %s: %s", tenmod.TenantUserStr, err.Error())
		reportError(w, startTime, "500", mon.UpdateTenantUserStr, msg, http.StatusInternalServerError)
		return
	}

	sendSuccessResponse(result, w, startTime, mon.UpdateTenantUserStr, tenmod.TenantUserStr, "Updated")
}

// GetTenantUser - fetches a tenant user
func (tsh *TenantServiceRESTHandler) GetTenantUser(w http.ResponseWriter, r *http.Request) {
	startTime := time.Now()

	// Get the IDs from the URL
	tenantID := getDBFieldFromRequest(r, 4)
	userID := getDBFieldFromRequest(r, 6)

	logger.Log.Infof("Fetching %s: %s", tenmod.TenantUserStr, userID)

	// Issue request to DAO Layer
	result, err := tsh.tenantDB.GetTenantUser(tenantID, userID)
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
	result, err := tsh.tenantDB.DeleteTenantUser(tenantID, userID)
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
	result, err := tsh.tenantDB.GetAllTenantUsers(tenantID)
	if err != nil {
		msg := fmt.Sprintf("Unable to retrieve %s list: %s", tenmod.TenantUserStr, err.Error())
		reportError(w, startTime, "500", mon.GetAllTenantUserStr, msg, http.StatusInternalServerError)
		return
	}

	sendSuccessResponse(result, w, startTime, mon.GetAllTenantUserStr, tenmod.TenantUserStr, "Retrieved list of")
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

	data := []tenmod.Domain{}
	err = jsonapi.Unmarshal(requestBytes, &data)
	if err != nil {
		msg := generateErrorMessage(http.StatusBadRequest, err.Error())
		reportError(w, startTime, "400", mon.CreateAdminUserStr, msg, http.StatusBadRequest)
		return
	}

	err = data[0].Validate(false)
	if err != nil {
		msg := generateErrorMessage(http.StatusBadRequest, err.Error())
		reportError(w, startTime, "400", mon.CreateAdminUserStr, msg, http.StatusBadRequest)
		return
	}

	logger.Log.Infof("Creating %s: %s", tenmod.TenantDomainStr, models.AsJSONString(&data))

	// Issue request to DAO Layer
	result, err := tsh.tenantDB.CreateTenantDomain(&data[0])
	if err != nil {
		msg := fmt.Sprintf("Unable to store %s: %s", tenmod.TenantDomainStr, err.Error())
		reportError(w, startTime, "500", mon.CreateTenantDomainStr, msg, http.StatusInternalServerError)
		return
	}

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

	data := []tenmod.Domain{}
	err = jsonapi.Unmarshal(requestBytes, &data)
	if err != nil {
		msg := generateErrorMessage(http.StatusBadRequest, err.Error())
		reportError(w, startTime, "400", mon.CreateAdminUserStr, msg, http.StatusBadRequest)
		return
	}

	err = data[0].Validate(true)
	if err != nil {
		msg := generateErrorMessage(http.StatusBadRequest, err.Error())
		reportError(w, startTime, "400", mon.CreateAdminUserStr, msg, http.StatusBadRequest)
		return
	}

	logger.Log.Infof("Updating %s: %s", tenmod.TenantDomainStr, models.AsJSONString(&data))

	// Issue request to DAO Layer
	result, err := tsh.tenantDB.UpdateTenantDomain(&data[0])
	if err != nil {
		msg := fmt.Sprintf("Unable to store %s: %s", tenmod.TenantDomainStr, err.Error())
		reportError(w, startTime, "500", mon.UpdateTenantDomainStr, msg, http.StatusInternalServerError)
		return
	}

	sendSuccessResponse(result, w, startTime, mon.UpdateTenantDomainStr, tenmod.TenantDomainStr, "Updated")
}

// GetTenantDomain - fetches a tenant domain
func (tsh *TenantServiceRESTHandler) GetTenantDomain(w http.ResponseWriter, r *http.Request) {
	startTime := time.Now()

	// Get the IDs from the URL
	tenantID := getDBFieldFromRequest(r, 4)
	userID := getDBFieldFromRequest(r, 6)

	logger.Log.Infof("Fetching %s: %s", tenmod.TenantDomainStr, userID)

	// Issue request to DAO Layer
	result, err := tsh.tenantDB.GetTenantDomain(tenantID, userID)
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
	userID := getDBFieldFromRequest(r, 6)

	logger.Log.Infof("Deleting %s: %s", tenmod.TenantDomainStr, userID)

	// Issue request to DAO Layer
	result, err := tsh.tenantDB.DeleteTenantDomain(tenantID, userID)
	if err != nil {
		msg := fmt.Sprintf("Unable to retrieve %s: %s", tenmod.TenantDomainStr, err.Error())
		reportError(w, startTime, "500", mon.DeleteTenantDomainStr, msg, http.StatusInternalServerError)
		return
	}

	sendSuccessResponse(result, w, startTime, mon.DeleteTenantDomainStr, tenmod.TenantDomainStr, "Deleted")
}

// GetAllTenantDomains - fetches list of tenant domains
func (tsh *TenantServiceRESTHandler) GetAllTenantDomains(w http.ResponseWriter, r *http.Request) {
	startTime := time.Now()

	tenantID := getDBFieldFromRequest(r, 4)

	logger.Log.Infof("Fetching %s list for Tenant %s", tenmod.TenantDomainStr, tenantID)

	// Issue request to DAO Layer
	result, err := tsh.tenantDB.GetAllTenantDomains(tenantID)
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
	data := tenmod.IngestionProfile{}
	err := unmarshalRequest(r, &data, false)
	if err != nil {
		msg := generateErrorMessage(http.StatusBadRequest, err.Error())
		reportError(w, startTime, "400", mon.CreateIngPrfStr, msg, http.StatusBadRequest)
		return
	}

	logger.Log.Infof("Creating %s: %s", tenmod.TenantIngestionProfileStr, models.AsJSONString(&data))

	// Issue request to DAO Layer
	result, err := tsh.tenantDB.CreateTenantIngestionProfile(&data)
	if err != nil {
		msg := fmt.Sprintf("Unable to store %s: %s", tenmod.TenantIngestionProfileStr, err.Error())
		reportError(w, startTime, "500", mon.CreateIngPrfStr, msg, http.StatusInternalServerError)
		return
	}

	sendSuccessResponse(result, w, startTime, mon.CreateIngPrfStr, tenmod.TenantIngestionProfileStr, "Created")
}

// UpdateTenantIngestionProfile - updates a tenant ingestion profile
func (tsh *TenantServiceRESTHandler) UpdateTenantIngestionProfile(w http.ResponseWriter, r *http.Request) {
	startTime := time.Now()

	// Unmarshal the request
	data := tenmod.IngestionProfile{}
	err := unmarshalRequest(r, &data, true)
	if err != nil {
		msg := generateErrorMessage(http.StatusBadRequest, err.Error())
		reportError(w, startTime, "400", mon.UpdateIngPrfStr, msg, http.StatusBadRequest)
		return
	}

	logger.Log.Infof("Updating %s: %s", tenmod.TenantIngestionProfileStr, models.AsJSONString(&data))

	// Issue request to DAO Layer
	result, err := tsh.tenantDB.UpdateTenantIngestionProfile(&data)
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
	result, err := tsh.tenantDB.GetTenantIngestionProfile(tenantID, dataID)
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
	result, err := tsh.tenantDB.DeleteTenantIngestionProfile(tenantID, dataID)
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
	result, err := tsh.tenantDB.GetActiveTenantIngestionProfile(tenantID)
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
	data := tenmod.ThresholdProfile{}
	err := unmarshalRequest(r, &data, false)
	if err != nil {
		msg := generateErrorMessage(http.StatusBadRequest, err.Error())
		reportError(w, startTime, "400", mon.CreateThrPrfStr, msg, http.StatusBadRequest)
		return
	}

	logger.Log.Infof("Creating %s: %s", tenmod.TenantThresholdProfileStr, models.AsJSONString(&data))

	// Issue request to DAO Layer
	result, err := tsh.tenantDB.CreateTenantThresholdProfile(&data)
	if err != nil {
		msg := fmt.Sprintf("Unable to store %s: %s", tenmod.TenantThresholdProfileStr, err.Error())
		reportError(w, startTime, "500", mon.CreateThrPrfStr, msg, http.StatusInternalServerError)
		return
	}

	sendSuccessResponse(result, w, startTime, mon.CreateThrPrfStr, tenmod.TenantThresholdProfileStr, "Created")
}

// UpdateTenantThresholdProfile - updates a tenant threshold profile
func (tsh *TenantServiceRESTHandler) UpdateTenantThresholdProfile(w http.ResponseWriter, r *http.Request) {
	startTime := time.Now()

	// Unmarshal the request
	data := tenmod.ThresholdProfile{}
	err := unmarshalRequest(r, &data, true)
	if err != nil {
		msg := generateErrorMessage(http.StatusBadRequest, err.Error())
		reportError(w, startTime, "400", mon.UpdateThrPrfStr, msg, http.StatusBadRequest)
		return
	}

	logger.Log.Infof("Updating %s: %s", tenmod.TenantThresholdProfileStr, models.AsJSONString(&data))

	// Issue request to DAO Layer
	result, err := tsh.tenantDB.UpdateTenantThresholdProfile(&data)
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
	result, err := tsh.tenantDB.GetTenantThresholdProfile(tenantID, dataID)
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

	// Issue request to DAO Layer
	result, err := tsh.tenantDB.DeleteTenantThresholdProfile(tenantID, dataID)
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
	result, err := tsh.tenantDB.GetAllTenantThresholdProfile(tenantID)
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

	data := []tenmod.MonitoredObject{}
	err = jsonapi.Unmarshal(requestBytes, &data)
	if err != nil {
		msg := generateErrorMessage(http.StatusBadRequest, err.Error())
		reportError(w, startTime, "400", mon.CreateAdminUserStr, msg, http.StatusBadRequest)
		return
	}

	err = data[0].Validate(false)
	if err != nil {
		msg := generateErrorMessage(http.StatusBadRequest, err.Error())
		reportError(w, startTime, "400", mon.CreateAdminUserStr, msg, http.StatusBadRequest)
		return
	}

	logger.Log.Infof("Creating %s: %s", tenmod.TenantMonitoredObjectStr, models.AsJSONString(&data))

	// Issue request to DAO Layer
	result, err := tsh.tenantDB.CreateMonitoredObject(&data[0])
	if err != nil {
		msg := fmt.Sprintf("Unable to store %s: %s", tenmod.TenantMonitoredObjectStr, err.Error())
		reportError(w, startTime, "500", mon.CreateMonObjStr, msg, http.StatusInternalServerError)
		return
	}

	sendSuccessResponse(result, w, startTime, mon.CreateMonObjStr, tenmod.TenantMonitoredObjectStr, "Created")
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

	data := []tenmod.MonitoredObject{}
	err = jsonapi.Unmarshal(requestBytes, &data)
	if err != nil {
		msg := generateErrorMessage(http.StatusBadRequest, err.Error())
		reportError(w, startTime, "400", mon.CreateAdminUserStr, msg, http.StatusBadRequest)
		return
	}

	err = data[0].Validate(true)
	if err != nil {
		msg := generateErrorMessage(http.StatusBadRequest, err.Error())
		reportError(w, startTime, "400", mon.CreateAdminUserStr, msg, http.StatusBadRequest)
		return
	}

	logger.Log.Infof("Updating %s: %s", tenmod.TenantMonitoredObjectStr, models.AsJSONString(&data))

	// Issue request to DAO Layer
	result, err := tsh.tenantDB.UpdateMonitoredObject(&data[0])
	if err != nil {
		msg := fmt.Sprintf("Unable to store %s: %s", tenmod.TenantMonitoredObjectStr, err.Error())
		reportError(w, startTime, "500", mon.UpdateMonObjStr, msg, http.StatusInternalServerError)
		return
	}

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
	result, err := tsh.tenantDB.GetMonitoredObject(tenantID, dataID)
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
	result, err := tsh.tenantDB.DeleteMonitoredObject(tenantID, dataID)
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
	result, err := tsh.tenantDB.GetAllMonitoredObjects(tenantID)
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
	result, err := tsh.tenantDB.GetMonitoredObjectToDomainMap(&data)
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

	logger.Log.Infof("Suuceesfully retrieved %s for Tenant %s", tenmod.MonitoredObjectToDomainMapStr, data.TenantID)
	trackAPIMetrics(startTime, "200", mon.GetMonObjToDomMapStr)
	fmt.Fprintf(w, string(res))
}

// CreateTenantMeta - creates a tenant metadata
func (tsh *TenantServiceRESTHandler) CreateTenantMeta(w http.ResponseWriter, r *http.Request) {
	startTime := time.Now()

	// Unmarshal the request
	data := tenmod.Metadata{}
	err := unmarshalRequest(r, &data, false)
	if err != nil {
		msg := generateErrorMessage(http.StatusBadRequest, err.Error())
		reportError(w, startTime, "400", mon.CreateTenantMetaStr, msg, http.StatusBadRequest)
		return
	}

	logger.Log.Infof("Creating %s: %s", tenmod.TenantMetaStr, models.AsJSONString(&data))

	// Issue request to DAO Layer
	result, err := tsh.tenantDB.CreateTenantMeta(&data)
	if err != nil {
		msg := fmt.Sprintf("Unable to store %s: %s", tenmod.TenantMetaStr, err.Error())
		reportError(w, startTime, "500", mon.CreateTenantMetaStr, msg, http.StatusInternalServerError)
		return
	}

	sendSuccessResponse(result, w, startTime, mon.CreateTenantMetaStr, tenmod.TenantMetaStr, "Created")
}

// UpdateTenantMeta - updates a tenant metadata
func (tsh *TenantServiceRESTHandler) UpdateTenantMeta(w http.ResponseWriter, r *http.Request) {
	startTime := time.Now()

	// Unmarshal the request
	data := tenmod.Metadata{}
	err := unmarshalRequest(r, &data, true)
	if err != nil {
		msg := generateErrorMessage(http.StatusBadRequest, err.Error())
		reportError(w, startTime, "400", mon.UpdateTenantMetaStr, msg, http.StatusBadRequest)
		return
	}

	logger.Log.Infof("Updating %s: %s", tenmod.TenantMetaStr, models.AsJSONString(&data))

	// Issue request to DAO Layer
	result, err := tsh.tenantDB.UpdateTenantMeta(&data)
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
	result, err := tsh.tenantDB.GetTenantMeta(tenantID)
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
	result, err := tsh.tenantDB.DeleteTenantMeta(tenantID)
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
		reportError(w, startTime, "400", mon.BulkUpdateMonObjStr, msg, http.StatusBadRequest)
		return
	}

	// Validate the request data
	for _, obj := range data {
		if err = obj.Validate(false); err != nil {

		}

		if obj.TenantID != tenantID {
			msg := fmt.Sprintf("Unable to Update %ss in bulk: %s", tenmod.TenantMonitoredObjectStr, "All Monitored Objects must have Tenant ID "+tenantID)
			reportError(w, startTime, "400", mon.BulkUpdateMonObjStr, msg, http.StatusBadRequest)
			return
		}
	}

	logger.Log.Infof("Bulk instering %ss: %s", tenmod.TenantMonitoredObjectStr, models.AsJSONString(&data))

	// Issue request to DAO Layer
	result, err := tsh.tenantDB.BulkInsertMonitoredObjects(tenantID, data)
	if err != nil {
		msg := fmt.Sprintf("Unable to store %s: %s", tenmod.TenantMonitoredObjectStr, err.Error())
		reportError(w, startTime, "500", mon.BulkUpdateMonObjStr, msg, http.StatusInternalServerError)
		return
	}

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
