package handlers

import (
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
	data := tenmod.User{}
	err := unmarshalRequest(r, &data, false)
	if err != nil {
		msg := generateErrorMessage(http.StatusBadRequest, err.Error())
		reportError(w, startTime, "400", mon.CreateTenantUserStr, msg, http.StatusBadRequest)
		return
	}

	logger.Log.Infof("Creating %s: %s", tenmod.TenantUserStr, models.AsJSONString(&data))

	// Issue request to DAO Layer
	result, err := tsh.tenantDB.CreateTenantUser(&data)
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
	data := tenmod.User{}
	err := unmarshalRequest(r, &data, true)
	if err != nil {
		msg := generateErrorMessage(http.StatusBadRequest, err.Error())
		reportError(w, startTime, "400", mon.UpdateTenantUserStr, msg, http.StatusBadRequest)
		return
	}

	logger.Log.Infof("Updating %s: %s", tenmod.TenantUserStr, models.AsJSONString(&data))

	// Issue request to DAO Layer
	result, err := tsh.tenantDB.UpdateTenantUser(&data)
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
