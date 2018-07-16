package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	db "github.com/accedian/adh-gather/datastore"
	"github.com/accedian/adh-gather/logger"
	"github.com/accedian/adh-gather/models"
	admmod "github.com/accedian/adh-gather/models/admin"
	mon "github.com/accedian/adh-gather/monitoring"
	"github.com/accedian/adh-gather/server"
	"github.com/gorilla/mux"
	"github.com/manyminds/api2go/jsonapi"
)

// AdminServiceRESTHandler - handler of logic for REST calls made to the Admin Service.
type AdminServiceRESTHandler struct {
	adminDB  db.AdminServiceDatastore
	tenantDB db.TenantServiceDatastore
	routes   []server.Route
}

// CreateAdminServiceRESTHandler - used to create a Admin Service REST handler which provides
// logic to serve the Admin Service REST calls
func CreateAdminServiceRESTHandler() *AdminServiceRESTHandler {
	result := new(AdminServiceRESTHandler)

	// Setup the DB implementation based on configuration
	db, err := GetAdminServiceDatastore()
	if err != nil {
		logger.Log.Fatalf("Unable to instantiate AdminServiceRESTHandler: %s", err.Error())
	}
	result.adminDB = db

	tdb, err := GetTenantServiceDatastore()
	if err != nil {
		logger.Log.Fatalf("Unable to instantiate AdminServiceRESTHandler: %s", err.Error())
	}
	result.tenantDB = tdb

	result.routes = []server.Route{
		server.Route{
			Name:        "CreateAdminUser",
			Method:      "POST",
			Pattern:     apiV1Prefix + "admin",
			HandlerFunc: BuildRouteHandlerWithRAC([]string{userRoleSkylight}, result.CreateAdminUser),
		},
		server.Route{
			Name:        "UpdateAdminUser",
			Method:      "PUT",
			Pattern:     apiV1Prefix + "admin",
			HandlerFunc: BuildRouteHandlerWithRAC([]string{userRoleSkylight}, result.UpdateAdminUser),
		},
		server.Route{
			Name:        "GetAdminUser",
			Method:      "GET",
			Pattern:     apiV1Prefix + "admin/{userID}",
			HandlerFunc: BuildRouteHandlerWithRAC([]string{userRoleSkylight}, result.GetAdminUser),
		},
		server.Route{
			Name:        "DeleteAdminUser",
			Method:      "DELETE",
			Pattern:     apiV1Prefix + "admin/{userID}",
			HandlerFunc: BuildRouteHandlerWithRAC([]string{userRoleSkylight}, result.DeleteAdminUser),
		},
		server.Route{
			Name:        "GetAllAdminUsers",
			Method:      "GET",
			Pattern:     apiV1Prefix + "admin-user-list",
			HandlerFunc: BuildRouteHandlerWithRAC([]string{userRoleSkylight}, result.GetAllAdminUsers),
		},
		server.Route{
			Name:        "CreateTenant",
			Method:      "POST",
			Pattern:     apiV1Prefix + "tenants",
			HandlerFunc: BuildRouteHandlerWithRAC([]string{userRoleSkylight}, result.CreateTenant),
		},
		server.Route{
			Name:        "UpdateTenant",
			Method:      "PUT",
			Pattern:     apiV1Prefix + "tenants",
			HandlerFunc: BuildRouteHandlerWithRAC([]string{userRoleSkylight}, result.UpdateTenant),
		},
		server.Route{
			Name:        "PatchTenant",
			Method:      "PATCH",
			Pattern:     apiV1Prefix + "tenants/{tenantID}",
			HandlerFunc: BuildRouteHandlerWithRAC([]string{userRoleSkylight}, result.PatchTenant),
		},
		server.Route{
			Name:        "GetTenant",
			Method:      "GET",
			Pattern:     apiV1Prefix + "tenants/{tenantID}",
			HandlerFunc: BuildRouteHandlerWithRAC([]string{userRoleSkylight}, result.GetTenant),
		},
		server.Route{
			Name:        "DeleteTenant",
			Method:      "DELETE",
			Pattern:     apiV1Prefix + "tenants/{tenantID}",
			HandlerFunc: BuildRouteHandlerWithRAC([]string{userRoleSkylight}, result.DeleteTenant),
		},
		server.Route{
			Name:        "GetAllTenants",
			Method:      "GET",
			Pattern:     apiV1Prefix + "tenant-list",
			HandlerFunc: BuildRouteHandlerWithRAC([]string{userRoleSkylight}, result.GetAllTenants),
		},
		// TODO: Make a "system" role for internal callers to be able to use to bypass auth checks. For
		// now, just commenting out auth for these calls.
		server.Route{
			Name:    "GetTenantIDByAlias",
			Method:  "GET",
			Pattern: apiV1Prefix + "tenant-by-alias/{value}",
			// HandlerFunc: BuildRouteHandlerWithRAC([]string{userRoleSkylight}, result.GetTenantIDByAlias),
			HandlerFunc: result.GetTenantIDByAlias,
		},
		server.Route{
			Name:    "GetTenantSummaryByAlias",
			Method:  "GET",
			Pattern: apiV1Prefix + "tenant-summary-by-alias/{value}",
			// HandlerFunc: BuildRouteHandlerWithRAC([]string{userRoleSkylight}, result.GetTenantSummaryByAlias),
			HandlerFunc: result.GetTenantSummaryByAlias,
		},
		server.Route{
			Name:        "CreateIngestionDictionary",
			Method:      "POST",
			Pattern:     apiV1Prefix + "ingestion-dictionaries",
			HandlerFunc: BuildRouteHandlerWithRAC([]string{userRoleSkylight}, result.CreateIngestionDictionary),
		},
		server.Route{
			Name:        "UpdateIngestionDictionary",
			Method:      "PUT",
			Pattern:     apiV1Prefix + "ingestion-dictionaries",
			HandlerFunc: BuildRouteHandlerWithRAC([]string{userRoleSkylight}, result.UpdateIngestionDictionary),
		},
		server.Route{
			Name:        "GetIngestionDictionary",
			Method:      "GET",
			Pattern:     apiV1Prefix + "ingestion-dictionaries",
			HandlerFunc: BuildRouteHandlerWithRAC([]string{userRoleSkylight}, result.GetIngestionDictionary),
		},
		server.Route{
			Name:        "DeleteIngestionDictionary",
			Method:      "DELETE",
			Pattern:     apiV1Prefix + "ingestion-dictionaries",
			HandlerFunc: BuildRouteHandlerWithRAC([]string{userRoleSkylight}, result.DeleteIngestionDictionary),
		},
		server.Route{
			Name:        "CreateValidTypes",
			Method:      "POST",
			Pattern:     apiV1Prefix + "valid-types",
			HandlerFunc: BuildRouteHandlerWithRAC([]string{userRoleSkylight}, result.CreateValidTypes),
		},
		server.Route{
			Name:        "UpdateValidTypes",
			Method:      "PUT",
			Pattern:     apiV1Prefix + "valid-types",
			HandlerFunc: BuildRouteHandlerWithRAC([]string{userRoleSkylight}, result.UpdateValidTypes),
		},
		server.Route{
			Name:        "GetValidTypes",
			Method:      "GET",
			Pattern:     apiV1Prefix + "valid-types",
			HandlerFunc: BuildRouteHandlerWithRAC([]string{userRoleSkylight}, result.GetValidTypes),
		},
		server.Route{
			Name:        "DeleteValidTypes",
			Method:      "DELETE",
			Pattern:     apiV1Prefix + "valid-types",
			HandlerFunc: BuildRouteHandlerWithRAC([]string{userRoleSkylight}, result.DeleteValidTypes),
		},
		server.Route{
			Name:        "GetSpecificValidTypes",
			Method:      "GET",
			Pattern:     apiV1Prefix + "specific-valid-types",
			HandlerFunc: BuildRouteHandlerWithRAC([]string{userRoleSkylight}, result.GetSpecificValidTypes),
		},
	}

	return result
}

// RegisterAPIHandlers - will bind any REST API routes defined in this service
// to the passed in request multiplexor.
func (ash *AdminServiceRESTHandler) RegisterAPIHandlers(router *mux.Router) {
	for _, route := range ash.routes {
		logger.Log.Debugf("Registering endpoint: %v", route)
		router.
			Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			Handler(route.HandlerFunc)
	}
}

// AddAdminViews - Add admin views to Admin DB.
func (ash *AdminServiceRESTHandler) AddAdminViews() error {
	logger.Log.Info("Adding Views to Admin DB")

	// Issue request to DAO Layer to Get the requested Tenant ID
	return ash.adminDB.AddAdminViews()
}

// CreateAdminUser - creates an admin user
func (ash *AdminServiceRESTHandler) CreateAdminUser(w http.ResponseWriter, r *http.Request) {
	startTime := time.Now()

	// Unmarshal the request
	requestBytes, err := getRequestBytes(r)
	if err != nil {
		msg := generateErrorMessage(http.StatusBadRequest, err.Error())
		reportError(w, startTime, "400", mon.CreateAdminUserStr, msg, http.StatusBadRequest)
		return
	}

	data := admmod.User{}
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

	logger.Log.Infof("Creating %s: %s", admmod.AdminUserStr, models.AsJSONString(&data))

	// Issue request to DAO Layer to Create the Admin User
	result, err := ash.adminDB.CreateAdminUser(&data)
	if err != nil {
		msg := fmt.Sprintf("Unable to store %s: %s", admmod.AdminUserStr, err.Error())
		reportError(w, startTime, "500", mon.CreateAdminUserStr, msg, http.StatusInternalServerError)
		return
	}

	sendSuccessResponse(result, w, startTime, mon.CreateAdminUserStr, admmod.AdminUserStr, "Created")
}

// UpdateAdminUser - updates an admin user
func (ash *AdminServiceRESTHandler) UpdateAdminUser(w http.ResponseWriter, r *http.Request) {
	startTime := time.Now()

	// Unmarshal the request
	requestBytes, err := getRequestBytes(r)
	if err != nil {
		msg := generateErrorMessage(http.StatusBadRequest, err.Error())
		reportError(w, startTime, "400", mon.CreateAdminUserStr, msg, http.StatusBadRequest)
		return
	}

	data := admmod.User{}
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

	logger.Log.Infof("Updating %s: %s", admmod.AdminUserStr, models.AsJSONString(&data))

	// Issue request to DAO Layer
	result, err := ash.adminDB.UpdateAdminUser(&data)
	if err != nil {
		msg := fmt.Sprintf("Unable to store %s: %s", admmod.AdminUserStr, err.Error())
		reportError(w, startTime, "500", mon.UpdateAdminUserStr, msg, http.StatusInternalServerError)
		return
	}

	sendSuccessResponse(result, w, startTime, mon.UpdateAdminUserStr, admmod.AdminUserStr, "Updated")
}

// GetAdminUser - fetches an admin user
func (ash *AdminServiceRESTHandler) GetAdminUser(w http.ResponseWriter, r *http.Request) {
	startTime := time.Now()

	// Get the ID from the URL
	userID := getDBFieldFromRequest(r, 4)

	logger.Log.Infof("Fetching %s: %s", admmod.AdminUserStr, userID)

	// Issue request to DAO Layer
	result, err := ash.adminDB.GetAdminUser(userID)
	if err != nil {
		msg := fmt.Sprintf("Unable to retrieve %s: %s", admmod.AdminUserStr, err.Error())
		reportError(w, startTime, "500", mon.GetAdminUserStr, msg, http.StatusInternalServerError)
		return
	}

	sendSuccessResponse(result, w, startTime, mon.GetAdminUserStr, admmod.AdminUserStr, "Retrieved")
}

// DeleteAdminUser - deletes an admin user
func (ash *AdminServiceRESTHandler) DeleteAdminUser(w http.ResponseWriter, r *http.Request) {
	startTime := time.Now()

	// Get the ID from the URL
	userID := getDBFieldFromRequest(r, 4)

	logger.Log.Infof("Deleting %s: %s", admmod.AdminUserStr, userID)

	// Issue request to DAO Layer
	result, err := ash.adminDB.DeleteAdminUser(userID)
	if err != nil {
		msg := fmt.Sprintf("Unable to retrieve %s: %s", admmod.AdminUserStr, err.Error())
		reportError(w, startTime, "500", mon.DeleteAdminUserStr, msg, http.StatusInternalServerError)
		return
	}

	sendSuccessResponse(result, w, startTime, mon.DeleteAdminUserStr, admmod.AdminUserStr, "Deleted")
}

// GetAllAdminUsers - fetches list of admin users
func (ash *AdminServiceRESTHandler) GetAllAdminUsers(w http.ResponseWriter, r *http.Request) {
	startTime := time.Now()

	logger.Log.Infof("Fetching %s list", admmod.AdminUserStr)

	// Issue request to DAO Layer
	result, err := ash.adminDB.GetAllAdminUsers()
	if err != nil {
		msg := fmt.Sprintf("Unable to retrieve %s list: %s", admmod.AdminUserStr, err.Error())
		reportError(w, startTime, "500", mon.GetAllAdminUserStr, msg, http.StatusInternalServerError)
		return
	}

	sendSuccessResponse(result, w, startTime, mon.GetAllAdminUserStr, admmod.AdminUserStr, "Retrieved list of")
}

// CreateTenant - creates a tenant including the default values for a user, ingestion profile,
// and threshood profile
func (ash *AdminServiceRESTHandler) CreateTenant(w http.ResponseWriter, r *http.Request) {
	startTime := time.Now()

	// Unmarshal the request
	requestBytes, err := getRequestBytes(r)
	if err != nil {
		msg := generateErrorMessage(http.StatusBadRequest, err.Error())
		reportError(w, startTime, "400", mon.CreateTenantStr, msg, http.StatusBadRequest)
		return
	}

	data := admmod.Tenant{}
	err = jsonapi.Unmarshal(requestBytes, &data)
	if err != nil {
		msg := generateErrorMessage(http.StatusBadRequest, err.Error())
		reportError(w, startTime, "400", mon.CreateTenantStr, msg, http.StatusBadRequest)
		return
	}

	err = data.Validate(false)
	if err != nil {
		msg := generateErrorMessage(http.StatusBadRequest, err.Error())
		reportError(w, startTime, "400", mon.CreateTenantStr, msg, http.StatusBadRequest)
		return
	}

	// Check if a tenant already exists with this name.
	existingTenantByName, _ := ash.adminDB.GetTenantIDByAlias(strings.ToLower(data.Name))
	if len(existingTenantByName) != 0 {
		msg := fmt.Sprintf("Unable to create Tenant %s. A Tenant with this name already exists", data.Name)
		reportError(w, startTime, "409", mon.CreateTenantStr, msg, http.StatusConflict)
		return
	}

	logger.Log.Infof("Creating %s: %s", admmod.TenantStr, models.AsJSONString(&data))

	// Issue request to DAO Layer to Create the Admin User
	result, err := ash.adminDB.CreateTenant(&data)
	if err != nil {
		msg := fmt.Sprintf("Unable to store %s: %s", admmod.TenantStr, err.Error())
		reportError(w, startTime, "500", mon.CreateTenantStr, msg, http.StatusInternalServerError)
		return
	}

	// Create a default Ingestion Profile for the Tenant.
	idForTenant := result.ID
	ingPrfData := createDefaultTenantIngPrf(idForTenant)
	logger.Log.Debugf("Sending to DAO: %s", models.AsJSONString(ingPrfData))
	prf, err := ash.tenantDB.CreateTenantIngestionProfile(ingPrfData)
	if err != nil {
		msg := fmt.Sprintf("Unable to create default Ingestion Profile %s", err.Error())
		reportError(w, startTime, "500", mon.CreateTenantStr, msg, http.StatusInternalServerError)
		return
	}
	logger.Log.Debugf("Got back from DAO: %s", models.AsJSONString(prf))

	// Create a default Threshold Profile for the Tenant
	threshPrfData := createDefaultTenantThresholdPrf(idForTenant)
	threshProfileResponse, err := ash.tenantDB.CreateTenantThresholdProfile(threshPrfData)
	if err != nil {
		msg := fmt.Sprintf("Unable to create default Threshold Profile %s", err.Error())
		reportError(w, startTime, "500", mon.CreateTenantStr, msg, http.StatusInternalServerError)
		return
	}

	// Create the tenant metadata
	// For the IDs used as references inside other objects, need to strip off the 'thresholdProfile_2_'
	// as this is just relational pouch adaption:
	meta := createDefaultTenantMeta(idForTenant, threshProfileResponse.ID, result.Name)
	_, err = ash.tenantDB.CreateTenantMeta(meta)
	if err != nil {
		msg := fmt.Sprintf("Unable to create Tenant metadata %s", err.Error())
		reportError(w, startTime, "500", mon.CreateTenantStr, msg, http.StatusInternalServerError)
		return
	}

	sendSuccessResponse(result, w, startTime, mon.CreateTenantStr, admmod.TenantStr, "Created")
}

// UpdateTenant - updates a Tenant
func (ash *AdminServiceRESTHandler) UpdateTenant(w http.ResponseWriter, r *http.Request) {
	startTime := time.Now()

	// Unmarshal the request
	requestBytes, err := getRequestBytes(r)
	if err != nil {
		msg := generateErrorMessage(http.StatusBadRequest, err.Error())
		reportError(w, startTime, "400", mon.UpdateTenantStr, msg, http.StatusBadRequest)
		return
	}

	data := admmod.Tenant{}
	err = jsonapi.Unmarshal(requestBytes, &data)
	if err != nil {
		msg := generateErrorMessage(http.StatusBadRequest, err.Error())
		reportError(w, startTime, "400", mon.UpdateTenantStr, msg, http.StatusBadRequest)
		return
	}

	err = data.Validate(true)
	if err != nil {
		msg := generateErrorMessage(http.StatusBadRequest, err.Error())
		reportError(w, startTime, "400", mon.UpdateTenantStr, msg, http.StatusBadRequest)
		return
	}

	logger.Log.Infof("Updating %s: %s", admmod.TenantStr, models.AsJSONString(&data))

	// Issue request to DAO Layer
	result, err := ash.adminDB.UpdateTenantDescriptor(&data)
	if err != nil {
		msg := fmt.Sprintf("Unable to store %s: %s", admmod.TenantStr, err.Error())
		reportError(w, startTime, "500", mon.UpdateTenantStr, msg, http.StatusInternalServerError)
		return
	}

	sendSuccessResponse(result, w, startTime, mon.UpdateTenantStr, admmod.TenantStr, "Updated")
}

//PatchTenant - Logic to update a tenant based on a potential subset of Tenant attribute value pairs from an HTTP request.
// Reports HTTP error responses for the following conditions:
// 	400 - if the request could not be parsed or if it failed Tenant validation logic
//  500 - if the tenant could not be retrieved from the datastore, if the merge failed, or if the merged data could not be pushed into the datastore
// Params:
//  w - the HTTP response to populate based on the success/failure of the request
//  r - the initiating HTTP patch request to run the business logic on
// Returns:
//  void
func (ash *AdminServiceRESTHandler) PatchTenant(w http.ResponseWriter, r *http.Request) {
	startTime := time.Now()

	// Get the ID from the URL in order to fetch the tenant from the DB
	tenantID := getDBFieldFromRequest(r, 4)

	// Retrieve tne request bytes from the payload in order to convert it to a map
	patchRequestBytes, err := getRequestBytes(r)
	if err != nil {
		msg := generateErrorMessage(http.StatusBadRequest, err.Error())
		reportError(w, startTime, "400", mon.PatchTenantStr, msg, http.StatusBadRequest)
		return
	}

	// Attempt to retrieve the tenant that we are trying to patch from the DB in order to do a merge
	fetchedTenant, err := ash.adminDB.GetTenantDescriptor(tenantID)
	if err != nil {
		//TODO we should try to return a 404 if the tenant is indeed not found. Unfortunately the response code from the db is buried in an error string
		msg := fmt.Sprintf("Unable to retrieve %s: %s", mon.PatchTenantStr, err.Error())
		reportError(w, startTime, "500", mon.PatchTenantStr, msg, http.StatusInternalServerError)
		return
	}

	// Merge the attributes passed in with the patch request to the tenant fetched from the datastore
	var patchedTenant *admmod.Tenant
	if err := models.MergeObjWithMap(fetchedTenant, patchRequestBytes); err != nil {
		msg := fmt.Sprintf("Unable to patch tenant with id %s: %s", tenantID, err.Error())
		reportError(w, startTime, "500", mon.PatchTenantStr, msg, http.StatusInternalServerError)
		return
	}
	patchedTenant = fetchedTenant

	// Ensure that the tenant is properly constructed following the merge prior to updating the record in the datastore
	err = patchedTenant.Validate(true)
	if err != nil {
		msg := generateErrorMessage(http.StatusBadRequest, err.Error())
		reportError(w, startTime, "400", mon.PatchTenantStr, msg, http.StatusBadRequest)
		return
	}

	// Finally update the tenant in the datastore with the merged map and fetched tenant
	result, err := ash.adminDB.UpdateTenantDescriptor(patchedTenant)
	if err != nil {
		msg := generateErrorMessage(http.StatusBadRequest, err.Error())
		reportError(w, startTime, "500", mon.PatchTenantStr, msg, http.StatusBadRequest)
		return
	}

	sendSuccessResponse(result, w, startTime, mon.PatchTenantStr, admmod.TenantStr, "Patched")
}

// GetTenant - fetches a Tenant
func (ash *AdminServiceRESTHandler) GetTenant(w http.ResponseWriter, r *http.Request) {
	startTime := time.Now()

	// Get the ID from the URL
	tenantID := getDBFieldFromRequest(r, 4)

	logger.Log.Infof("Fetching %s: %s", admmod.TenantStr, tenantID)

	// Issue request to DAO Layer
	result, err := ash.adminDB.GetTenantDescriptor(tenantID)
	if err != nil {
		msg := fmt.Sprintf("Unable to retrieve %s: %s", admmod.TenantStr, err.Error())
		reportError(w, startTime, "500", mon.GetTenantStr, msg, http.StatusInternalServerError)
		return
	}

	sendSuccessResponse(result, w, startTime, mon.GetTenantStr, admmod.TenantStr, "Retrieved")
}

// DeleteTenant - deletes a Tenant
func (ash *AdminServiceRESTHandler) DeleteTenant(w http.ResponseWriter, r *http.Request) {
	startTime := time.Now()

	// Get the ID from the URL
	tenantID := getDBFieldFromRequest(r, 4)

	logger.Log.Infof("Deleting %s: %s", admmod.TenantStr, tenantID)

	// Issue request to DAO Layer
	result, err := ash.adminDB.DeleteTenant(tenantID)
	if err != nil {
		msg := fmt.Sprintf("Unable to retrieve %s: %s", admmod.TenantStr, err.Error())
		reportError(w, startTime, "500", mon.DeleteTenantStr, msg, http.StatusInternalServerError)
		return
	}

	sendSuccessResponse(result, w, startTime, mon.DeleteTenantStr, admmod.TenantStr, "Deleted")
}

// GetAllTenants - fetches list of Tenants
func (ash *AdminServiceRESTHandler) GetAllTenants(w http.ResponseWriter, r *http.Request) {
	startTime := time.Now()

	logger.Log.Infof("Fetching %s list", admmod.TenantStr)

	// Issue request to DAO Layer
	result, err := ash.adminDB.GetAllTenantDescriptors()
	if err != nil {
		msg := fmt.Sprintf("Unable to retrieve %s list: %s", admmod.TenantStr, err.Error())
		reportError(w, startTime, "500", mon.GetAllTenantStr, msg, http.StatusInternalServerError)
		return
	}

	sendSuccessResponse(result, w, startTime, mon.GetAllTenantStr, admmod.TenantStr, "Retrieved list of")
}

// GetTenantIDByAlias - fetches a Tenant ID by its known alias.
func (ash *AdminServiceRESTHandler) GetTenantIDByAlias(w http.ResponseWriter, r *http.Request) {
	startTime := time.Now()

	// Get the ID from the URL
	tenantName := getDBFieldFromRequest(r, 4)

	logger.Log.Infof("Fetching ID for Tenant %s", tenantName)

	// Issue request to DAO Layer
	result, err := ash.adminDB.GetTenantIDByAlias(tenantName)
	if err != nil {
		msg := fmt.Sprintf("Unable to retrieve %s ID: %s", admmod.TenantStr, err.Error())
		reportError(w, startTime, "500", mon.GetTenantIDByAliasStr, msg, http.StatusInternalServerError)
		return
	}

	logger.Log.Infof("Successfully retrieved ID %s for alias %s", result, tenantName)
	trackAPIMetrics(startTime, "200", mon.GetTenantIDByAliasStr)
	fmt.Fprintf(w, result)
}

// GetTenantSummaryByAlias - fetches a Tenant summary by its known alias.
func (ash *AdminServiceRESTHandler) GetTenantSummaryByAlias(w http.ResponseWriter, r *http.Request) {
	startTime := time.Now()

	// Get the ID from the URL
	tenantName := getDBFieldFromRequest(r, 4)

	logger.Log.Infof("Fetching ID for Tenant %s", tenantName)

	// Issue request to DAO Layer
	result, err := ash.adminDB.GetTenantIDByAlias(tenantName)
	if err != nil {
		msg := fmt.Sprintf("Unable to retrieve %s summary for %s: %s", admmod.TenantStr, tenantName, err.Error())
		reportError(w, startTime, "500", mon.GetTenantSummaryByAliasStr, msg, http.StatusInternalServerError)
		return
	}

	summary := admmod.TenantSummary{Alias: tenantName, ID: result}
	response, err := json.Marshal(summary)
	if err != nil {
		msg := fmt.Sprintf("Unable to marshal response summary response for %s %s: %s", admmod.TenantStr, tenantName, err.Error())
		reportError(w, startTime, "500", mon.GetTenantSummaryByAliasStr, msg, http.StatusInternalServerError)
		return
	}

	logger.Log.Infof("Successfully retrieved ID %s for alias %s", result, tenantName)
	trackAPIMetrics(startTime, "200", mon.GetTenantSummaryByAliasStr)
	fmt.Fprintf(w, string(response))
}

// CreateIngestionDictionaryInternal - provides access to operation to create an ingestion dictionary without
// the need of a REST call.
func (ash *AdminServiceRESTHandler) CreateIngestionDictionaryInternal(dict *admmod.IngestionDictionary) (*admmod.IngestionDictionary, error) {
	return ash.adminDB.CreateIngestionDictionary(dict)
}

// CreateIngestionDictionary - creates an Ingestion Dictionary
func (ash *AdminServiceRESTHandler) CreateIngestionDictionary(w http.ResponseWriter, r *http.Request) {
	startTime := time.Now()

	// Unmarshal the request
	requestBytes, err := getRequestBytes(r)
	if err != nil {
		msg := generateErrorMessage(http.StatusBadRequest, err.Error())
		reportError(w, startTime, "400", mon.CreateAdminUserStr, msg, http.StatusBadRequest)
		return
	}

	data := admmod.IngestionDictionary{}
	logger.Log.Debugf("byte data: %s", string(requestBytes))
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

	logger.Log.Infof("Creating %s: %s", admmod.IngestionDictionaryStr, models.AsJSONString(&data))

	// Issue request to DAO Layer to Create the record
	result, err := ash.CreateIngestionDictionaryInternal(&data)
	if err != nil {
		msg := fmt.Sprintf("Unable to store %s: %s", admmod.IngestionDictionaryStr, err.Error())
		reportError(w, startTime, "500", mon.CreateIngDictStr, msg, http.StatusInternalServerError)
		return
	}

	sendSuccessResponse(result, w, startTime, mon.CreateIngDictStr, admmod.IngestionDictionaryStr, "Created")
}

// UpdateIngestionDictionaryInternal - provides access to operation to update an ingestion dictionary without
// the need of a REST call.
func (ash *AdminServiceRESTHandler) UpdateIngestionDictionaryInternal(dict *admmod.IngestionDictionary) (*admmod.IngestionDictionary, error) {
	return ash.adminDB.UpdateIngestionDictionary(dict)
}

// UpdateIngestionDictionary - updates an Ingestion Dictionary
func (ash *AdminServiceRESTHandler) UpdateIngestionDictionary(w http.ResponseWriter, r *http.Request) {
	startTime := time.Now()

	// Unmarshal the request
	requestBytes, err := getRequestBytes(r)
	if err != nil {
		msg := generateErrorMessage(http.StatusBadRequest, err.Error())
		reportError(w, startTime, "400", mon.CreateAdminUserStr, msg, http.StatusBadRequest)
		return
	}

	data := admmod.IngestionDictionary{}
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

	logger.Log.Infof("Updating %s: %s", admmod.IngestionDictionaryStr, models.AsJSONString(&data))

	// Issue request to DAO Layer
	result, err := ash.UpdateIngestionDictionaryInternal(&data)
	if err != nil {
		msg := fmt.Sprintf("Unable to store %s: %s", admmod.IngestionDictionaryStr, err.Error())
		reportError(w, startTime, "500", mon.UpdateIngDictStr, msg, http.StatusInternalServerError)
		return
	}

	sendSuccessResponse(result, w, startTime, mon.UpdateIngDictStr, admmod.IngestionDictionaryStr, "Updated")
}

// GetIngestionDictionaryInternal - provides access to operation to fetch an ingestion dictionary without
// the need of a REST call.
func (ash *AdminServiceRESTHandler) GetIngestionDictionaryInternal() (*admmod.IngestionDictionary, error) {
	return ash.adminDB.GetIngestionDictionary()
}

// GetIngestionDictionary - fetches an Ingestion Dictionary
func (ash *AdminServiceRESTHandler) GetIngestionDictionary(w http.ResponseWriter, r *http.Request) {
	startTime := time.Now()

	logger.Log.Infof("Fetching %s", admmod.IngestionDictionaryStr)

	// Issue request to DAO Layer
	result, err := ash.GetIngestionDictionaryInternal()
	if err != nil {
		msg := fmt.Sprintf("Unable to retrieve %s: %s", admmod.IngestionDictionaryStr, err.Error())
		reportError(w, startTime, "500", mon.GetIngDictStr, msg, http.StatusInternalServerError)
		return
	}

	sendSuccessResponse(result, w, startTime, mon.GetIngDictStr, admmod.IngestionDictionaryStr, "Retrieved")
}

// DeleteIngestionDictionary - deletes an Ingestion Dictionary
func (ash *AdminServiceRESTHandler) DeleteIngestionDictionary(w http.ResponseWriter, r *http.Request) {
	startTime := time.Now()

	logger.Log.Infof("Deleting %s", admmod.IngestionDictionaryStr)

	// Issue request to DAO Layer
	result, err := ash.adminDB.DeleteIngestionDictionary()
	if err != nil {
		msg := fmt.Sprintf("Unable to retrieve %s: %s", admmod.IngestionDictionaryStr, err.Error())
		reportError(w, startTime, "500", mon.DeleteIngDictStr, msg, http.StatusInternalServerError)
		return
	}

	sendSuccessResponse(result, w, startTime, mon.DeleteIngDictStr, admmod.IngestionDictionaryStr, "Deleted")
}

// CreateValidTypesInternal - provides access to operation to create a valid types object without
// the need of a REST call.
func (ash *AdminServiceRESTHandler) CreateValidTypesInternal(vt *admmod.ValidTypes) (*admmod.ValidTypes, error) {
	return ash.adminDB.CreateValidTypes(vt)
}

// CreateValidTypes - creates a Valid Types object
func (ash *AdminServiceRESTHandler) CreateValidTypes(w http.ResponseWriter, r *http.Request) {
	startTime := time.Now()

	// Unmarshal the request
	requestBytes, err := getRequestBytes(r)
	if err != nil {
		msg := generateErrorMessage(http.StatusBadRequest, err.Error())
		reportError(w, startTime, "400", mon.CreateAdminUserStr, msg, http.StatusBadRequest)
		return
	}

	data := admmod.ValidTypes{}
	logger.Log.Debugf("byte data: %s", string(requestBytes))
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

	logger.Log.Infof("Creating %s: %s", admmod.ValidTypesStr, models.AsJSONString(&data))

	// Issue request to DAO Layer to Create the record
	result, err := ash.CreateValidTypesInternal(&data)
	if err != nil {
		msg := fmt.Sprintf("Unable to store %s: %s", admmod.ValidTypesStr, err.Error())
		reportError(w, startTime, "500", mon.CreateValidTypesStr, msg, http.StatusInternalServerError)
		return
	}

	sendSuccessResponse(result, w, startTime, mon.CreateValidTypesStr, admmod.ValidTypesStr, "Created")
}

// UpdateValidTypesInternal - provides access to operation to update a valid types object without
// the need of a REST call.
func (ash *AdminServiceRESTHandler) UpdateValidTypesInternal(vt *admmod.ValidTypes) (*admmod.ValidTypes, error) {
	return ash.adminDB.UpdateValidTypes(vt)
}

// UpdateValidTypes - updates a Valid Types object
func (ash *AdminServiceRESTHandler) UpdateValidTypes(w http.ResponseWriter, r *http.Request) {
	startTime := time.Now()

	// Unmarshal the request
	requestBytes, err := getRequestBytes(r)
	if err != nil {
		msg := generateErrorMessage(http.StatusBadRequest, err.Error())
		reportError(w, startTime, "400", mon.CreateAdminUserStr, msg, http.StatusBadRequest)
		return
	}

	data := []admmod.ValidTypes{}
	logger.Log.Debugf("byte data: %s", string(requestBytes))
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

	logger.Log.Infof("Updating %s: %s", admmod.ValidTypesStr, models.AsJSONString(&data[0]))

	// Issue request to DAO Layer
	result, err := ash.UpdateValidTypesInternal(&data[0])
	if err != nil {
		msg := fmt.Sprintf("Unable to store %s: %s", admmod.ValidTypesStr, err.Error())
		reportError(w, startTime, "500", mon.UpdateValidTypesStr, msg, http.StatusInternalServerError)
		return
	}

	sendSuccessResponse(result, w, startTime, mon.UpdateValidTypesStr, admmod.ValidTypesStr, "Updated")
}

// GetValidTypesInternal - provides access to operation to fetch a valid types object without
// the need of a REST call.
func (ash *AdminServiceRESTHandler) GetValidTypesInternal() (*admmod.ValidTypes, error) {
	return ash.adminDB.GetValidTypes()
}

// GetValidTypes - fetches a Valid Types Object
func (ash *AdminServiceRESTHandler) GetValidTypes(w http.ResponseWriter, r *http.Request) {
	startTime := time.Now()

	logger.Log.Infof("Fetching %s", admmod.ValidTypesStr)

	// Issue request to DAO Layer
	result, err := ash.GetValidTypesInternal()
	if err != nil {
		msg := fmt.Sprintf("Unable to retrieve %s: %s", admmod.ValidTypesStr, err.Error())
		reportError(w, startTime, "500", mon.GetValidTypesStr, msg, http.StatusInternalServerError)
		return
	}

	sendSuccessResponse(result, w, startTime, mon.GetValidTypesStr, admmod.ValidTypesStr, "Retrieved")
}

// GetSpecificValidTypes - fetches specific portions of a Valid Types object
func (ash *AdminServiceRESTHandler) GetSpecificValidTypes(w http.ResponseWriter, r *http.Request) {
	startTime := time.Now()

	data := admmod.ValidTypesRequest{}
	queryParams := r.URL.Query()

	queryVal := queryParams.Get("monitoredObjectTypes")
	if queryVal != "" {
		boolVal, err := strconv.ParseBool(queryVal)
		if err != nil {
			msg := generateErrorMessage(http.StatusBadRequest, err.Error())
			reportError(w, startTime, "400", mon.GetSpecificValidTypesStr, msg, http.StatusBadRequest)
			return
		}
		data.MonitoredObjectTypes = boolVal
	}
	queryVal = queryParams.Get("monitoredObjectDeviceTypes")
	if queryVal != "" {
		boolVal, err := strconv.ParseBool(queryVal)
		if err != nil {
			msg := generateErrorMessage(http.StatusBadRequest, err.Error())
			reportError(w, startTime, "400", mon.GetSpecificValidTypesStr, msg, http.StatusBadRequest)
			return
		}
		data.MonitoredObjectDeviceTypes = boolVal
	}

	logger.Log.Infof("Fetching specific %s for filter %v", admmod.ValidTypesStr, models.AsJSONString(data))

	// Issue request to DAO Layer
	result, err := ash.adminDB.GetSpecificValidTypes(&data)
	if err != nil {
		msg := fmt.Sprintf("Unable to retrieve %s: %s", admmod.ValidTypesStr, err.Error())
		reportError(w, startTime, "500", mon.GetSpecificValidTypesStr, msg, http.StatusInternalServerError)
		return
	}

	sendSuccessResponse(result, w, startTime, mon.GetSpecificValidTypesStr, admmod.ValidTypesStr, "Retrieved")
}

// DeleteValidTypes - deletes a Valid Types object
func (ash *AdminServiceRESTHandler) DeleteValidTypes(w http.ResponseWriter, r *http.Request) {
	startTime := time.Now()

	logger.Log.Infof("Deleting %s", admmod.ValidTypesStr)

	// Issue request to DAO Layer
	result, err := ash.adminDB.DeleteValidTypes()
	if err != nil {
		msg := fmt.Sprintf("Unable to retrieve %s: %s", admmod.ValidTypesStr, err.Error())
		reportError(w, startTime, "500", mon.DeleteValidTypesStr, msg, http.StatusInternalServerError)
		return
	}

	sendSuccessResponse(result, w, startTime, mon.DeleteValidTypesStr, admmod.ValidTypesStr, "Deleted")
}
