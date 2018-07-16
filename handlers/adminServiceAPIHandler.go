package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/accedian/adh-gather/models"
	"github.com/accedian/adh-gather/swagmodels"
	"github.com/manyminds/api2go/jsonapi"

	"github.com/accedian/adh-gather/datastore"
	"github.com/accedian/adh-gather/logger"
	admmod "github.com/accedian/adh-gather/models/admin"
	mon "github.com/accedian/adh-gather/monitoring"
	"github.com/accedian/adh-gather/restapi/operations/admin_provisioning_service"
	"github.com/go-openapi/runtime/middleware"
)

// HandleCreateTenant - create a new tenant
func HandleCreateTenant(allowedRoles []string, adminDB datastore.AdminServiceDatastore, tenantDB datastore.TenantServiceDatastore) func(params admin_provisioning_service.CreateTenantParams) middleware.Responder {
	return func(params admin_provisioning_service.CreateTenantParams) middleware.Responder {
		startTime := time.Now()
		incrementAPICounters(mon.APIRecieved, mon.AdminAPIRecieved)
		logger.Log.Infof("Creating %s: %s", admmod.TenantStr, params.Body.Data.Attributes.Name)

		if !isRequestAuthorized(params.HTTPRequest, allowedRoles) {
			return admin_provisioning_service.NewCreateTenantForbidden().WithPayload(reportAPIError(fmt.Sprintf("Create %s operation not authorized for role: %s", admmod.TenantStr, params.HTTPRequest.Header.Get(xFwdUserRoles)), startTime, http.StatusForbidden, mon.CreateTenantStr, mon.APICompleted, mon.AdminAPICompleted))
		}

		// Unmarshal the request
		requestBytes, err := json.Marshal(params.Body)
		if err != nil {
			return admin_provisioning_service.NewCreateTenantBadRequest().WithPayload(reportAPIError(generateErrorMessage(http.StatusBadRequest, err.Error()), startTime, http.StatusBadRequest, mon.CreateTenantStr, mon.APICompleted, mon.AdminAPICompleted))
		}

		data := admmod.Tenant{}
		err = jsonapi.Unmarshal(requestBytes, &data)
		if err != nil {
			return admin_provisioning_service.NewCreateTenantBadRequest().WithPayload(reportAPIError(generateErrorMessage(http.StatusBadRequest, err.Error()), startTime, http.StatusBadRequest, mon.CreateTenantStr, mon.APICompleted, mon.AdminAPICompleted))
		}

		err = data.Validate(false)
		if err != nil {
			return admin_provisioning_service.NewCreateTenantBadRequest().WithPayload(reportAPIError(generateErrorMessage(http.StatusBadRequest, err.Error()), startTime, http.StatusBadRequest, mon.CreateTenantStr, mon.APICompleted, mon.AdminAPICompleted))
		}

		// Check if a tenant already exists with this name.
		existingTenantByName, _ := adminDB.GetTenantIDByAlias(strings.ToLower(data.Name))
		if len(existingTenantByName) != 0 {
			return admin_provisioning_service.NewCreateTenantConflict().WithPayload(reportAPIError(fmt.Sprintf("Unable to create Tenant %s. A Tenant with this name already exists", data.Name), startTime, http.StatusConflict, mon.CreateTenantStr, mon.APICompleted, mon.AdminAPICompleted))
		}

		// Issue request to DAO Layer to Create the Admin User
		result, err := adminDB.CreateTenant(&data)
		if err != nil {
			return admin_provisioning_service.NewCreateTenantInternalServerError().WithPayload(reportAPIError(fmt.Sprintf("Unable to store %s: %s", admmod.TenantStr, err.Error()), startTime, http.StatusInternalServerError, mon.CreateTenantStr, mon.APICompleted, mon.AdminAPICompleted))
		}

		// Create a default Ingestion Profile for the Tenant.
		idForTenant := result.ID
		ingPrfData := createDefaultTenantIngPrf(idForTenant)
		_, err = tenantDB.CreateTenantIngestionProfile(ingPrfData)
		if err != nil {
			return admin_provisioning_service.NewCreateTenantInternalServerError().WithPayload(reportAPIError(fmt.Sprintf("Unable to create default Ingestion Profile %s", err.Error()), startTime, http.StatusInternalServerError, mon.CreateTenantStr, mon.APICompleted, mon.AdminAPICompleted))
		}

		// Create a default Threshold Profile for the Tenant
		threshPrfData := createDefaultTenantThresholdPrf(idForTenant)
		threshProfileResponse, err := tenantDB.CreateTenantThresholdProfile(threshPrfData)
		if err != nil {
			return admin_provisioning_service.NewCreateTenantInternalServerError().WithPayload(reportAPIError(fmt.Sprintf("Unable to create default Threshold Profile %s", err.Error()), startTime, http.StatusInternalServerError, mon.CreateTenantStr, mon.APICompleted, mon.AdminAPICompleted))
		}

		// Create the tenant metadata
		// For the IDs used as references inside other objects, need to strip off the 'thresholdProfile_2_'
		// as this is just relational pouch adaption:
		meta := createDefaultTenantMeta(idForTenant, threshProfileResponse.ID, result.Name)
		_, err = tenantDB.CreateTenantMeta(meta)
		if err != nil {
			return admin_provisioning_service.NewCreateTenantInternalServerError().WithPayload(reportAPIError(fmt.Sprintf("Unable to create Tenant metadata %s", err.Error()), startTime, http.StatusInternalServerError, mon.CreateTenantStr, mon.APICompleted, mon.AdminAPICompleted))
		}

		converted := swagmodels.JSONAPITenant{}
		err = convertToJsonapiObject(result, &converted)
		if err != nil {
			return admin_provisioning_service.NewCreateTenantInternalServerError().WithPayload(reportAPIError(fmt.Sprintf("Unable to convert %s data to jsonapi return format: %s", admmod.TenantStr, err.Error()), startTime, http.StatusInternalServerError, mon.CreateTenantStr, mon.APICompleted, mon.AdminAPICompleted))
		}

		reportAPICompletionState(startTime, http.StatusOK, mon.CreateTenantStr, mon.APICompleted, mon.AdminAPICompleted)
		logger.Log.Infof("Created %s %s", admmod.TenantStr, models.AsJSONString(converted))
		return admin_provisioning_service.NewCreateTenantOK().WithPayload(&converted)
	}
}

// HandleGetTenant - retrieve a tenant by the tenant ID.
func HandleGetTenant(allowedRoles []string, adminDB datastore.AdminServiceDatastore) func(params admin_provisioning_service.GetTenantParams) middleware.Responder {
	return func(params admin_provisioning_service.GetTenantParams) middleware.Responder {
		startTime := time.Now()
		incrementAPICounters(mon.APIRecieved, mon.AdminAPIRecieved)

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

func HandlePatchTenant(allowedRoles []string, adminDB datastore.AdminServiceDatastore) func(params admin_provisioning_service.PatchTenantParams) middleware.Responder {
	return func(params admin_provisioning_service.PatchTenantParams) middleware.Responder {
		startTime := time.Now()
		incrementAPICounters(mon.APIRecieved, mon.AdminAPIRecieved)
		logger.Log.Infof("Updating %s: %s", admmod.TenantStr, params.TenantID)

		// Retrieve tne request bytes from the payload in order to convert it to a map
		patchRequestBytes, err := json.Marshal(params.Body)
		if err != nil {
			return admin_provisioning_service.NewPatchTenantBadRequest().WithPayload(reportAPIError(generateErrorMessage(http.StatusBadRequest, err.Error()), startTime, http.StatusBadRequest, mon.PatchTenantStr, mon.APICompleted, mon.AdminAPICompleted))
		}

		// Attempt to retrieve the tenant that we are trying to patch from the DB in order to do a merge
		fetchedTenant, err := adminDB.GetTenantDescriptor(params.TenantID)
		if err != nil {
			//TODO we should try to return a 404 if the tenant is indeed not found. Unfortunately the response code from the db is buried in an error string
			return admin_provisioning_service.NewPatchTenantInternalServerError().WithPayload(reportAPIError(fmt.Sprintf("Unable to retrieve %s: %s", mon.PatchTenantStr, err.Error()), startTime, http.StatusInternalServerError, mon.PatchTenantStr, mon.APICompleted, mon.AdminAPICompleted))
		}

		// Merge the attributes passed in with the patch request to the tenant fetched from the datastore
		var patchedTenant *admmod.Tenant
		if err := models.MergeObjWithMap(fetchedTenant, patchRequestBytes); err != nil {
			return admin_provisioning_service.NewPatchTenantInternalServerError().WithPayload(reportAPIError(fmt.Sprintf("Unable to patch tenant with id %s: %s", params.TenantID, err.Error()), startTime, http.StatusInternalServerError, mon.PatchTenantStr, mon.APICompleted, mon.AdminAPICompleted))
		}
		patchedTenant = fetchedTenant

		// Ensure that the tenant is properly constructed following the merge prior to updating the record in the datastore
		err = patchedTenant.Validate(true)
		if err != nil {
			return admin_provisioning_service.NewPatchTenantBadRequest().WithPayload(reportAPIError(generateErrorMessage(http.StatusBadRequest, err.Error()), startTime, http.StatusBadRequest, mon.PatchTenantStr, mon.APICompleted, mon.AdminAPICompleted))
		}

		// Finally update the tenant in the datastore with the merged map and fetched tenant
		result, err := adminDB.UpdateTenantDescriptor(patchedTenant)
		if err != nil {
			return admin_provisioning_service.NewPatchTenantInternalServerError().WithPayload(reportAPIError(generateErrorMessage(http.StatusBadRequest, err.Error()), startTime, http.StatusInternalServerError, mon.PatchTenantStr, mon.APICompleted, mon.AdminAPICompleted))
		}

		converted := swagmodels.JSONAPITenant{}
		err = convertToJsonapiObject(result, &converted)
		if err != nil {
			return admin_provisioning_service.NewPatchTenantInternalServerError().WithPayload(reportAPIError(fmt.Sprintf("Unable to convert %s data to jsonapi return format: %s", admmod.TenantStr, err.Error()), startTime, http.StatusInternalServerError, mon.PatchTenantStr, mon.APICompleted, mon.AdminAPICompleted))
		}

		reportAPICompletionState(startTime, http.StatusOK, mon.PatchTenantStr, mon.APICompleted, mon.AdminAPICompleted)
		logger.Log.Infof("Updated %s %s", admmod.TenantStr, models.AsJSONString(converted))
		return admin_provisioning_service.NewPatchTenantOK().WithPayload(&converted)
	}
}

// HandleDeleteTenant - delete a tenant by the tenant ID.
func HandleDeleteTenant(allowedRoles []string, adminDB datastore.AdminServiceDatastore) func(params admin_provisioning_service.DeleteTenantParams) middleware.Responder {
	return func(params admin_provisioning_service.DeleteTenantParams) middleware.Responder {
		startTime := time.Now()
		incrementAPICounters(mon.APIRecieved, mon.AdminAPIRecieved)
		logger.Log.Infof("Deleting %s: %s", admmod.TenantStr, params.TenantID)

		if !isRequestAuthorized(params.HTTPRequest, allowedRoles) {
			return admin_provisioning_service.NewDeleteTenantForbidden().WithPayload(reportAPIError(fmt.Sprintf("Delete %s operation not authorized for role: %s", admmod.TenantStr, params.HTTPRequest.Header.Get(xFwdUserRoles)), startTime, http.StatusForbidden, mon.DeleteTenantStr, mon.APICompleted, mon.AdminAPICompleted))
		}

		// Issue request to DAO Layer
		result, err := adminDB.DeleteTenant(params.TenantID)
		if err != nil {
			return admin_provisioning_service.NewDeleteTenantInternalServerError().WithPayload(reportAPIError(fmt.Sprintf("Unable to delete %s: %s", admmod.TenantStr, err.Error()), startTime, http.StatusInternalServerError, mon.DeleteTenantStr, mon.APICompleted, mon.AdminAPICompleted))
		}

		converted := swagmodels.JSONAPITenant{}
		err = convertToJsonapiObject(result, &converted)
		if err != nil {
			return admin_provisioning_service.NewDeleteTenantInternalServerError().WithPayload(reportAPIError(fmt.Sprintf("Unable to convert %s data to jsonapi return format: %s", admmod.TenantStr, err.Error()), startTime, http.StatusInternalServerError, mon.DeleteTenantStr, mon.APICompleted, mon.AdminAPICompleted))
		}

		reportAPICompletionState(startTime, http.StatusOK, mon.DeleteTenantStr, mon.APICompleted, mon.AdminAPICompleted)
		logger.Log.Infof("Deleted %s %s", admmod.TenantStr, models.AsJSONString(converted))
		return admin_provisioning_service.NewGetTenantOK().WithPayload(&converted)
	}
}

// HandleGetAllTenants - retrieve all tenants
func HandleGetAllTenants(allowedRoles []string, adminDB datastore.AdminServiceDatastore) func(params admin_provisioning_service.GetAllTenantsParams) middleware.Responder {
	return func(params admin_provisioning_service.GetAllTenantsParams) middleware.Responder {
		startTime := time.Now()
		incrementAPICounters(mon.APIRecieved, mon.AdminAPIRecieved)
		logger.Log.Infof("Fetching %s list", admmod.TenantStr)

		if !isRequestAuthorized(params.HTTPRequest, allowedRoles) {
			return admin_provisioning_service.NewGetAllTenantsForbidden().WithPayload(reportAPIError(fmt.Sprintf("Get All %ss operation not authorized for role: %s", admmod.TenantStr, params.HTTPRequest.Header.Get(xFwdUserRoles)), startTime, http.StatusForbidden, mon.GetAllTenantStr, mon.APICompleted, mon.AdminAPICompleted))
		}

		// Issue request to DAO Layer
		result, err := adminDB.GetAllTenantDescriptors()
		if err != nil {
			return admin_provisioning_service.NewGetAllTenantsInternalServerError().WithPayload(reportAPIError(fmt.Sprintf("Unable to retrieve %s list: %s", admmod.TenantStr, err.Error()), startTime, http.StatusInternalServerError, mon.GetAllTenantStr, mon.APICompleted, mon.AdminAPICompleted))
		}

		converted := swagmodels.JSONAPITenantList{}
		err = convertToJsonapiObject(result, &converted)
		if err != nil {
			return admin_provisioning_service.NewGetAllTenantsInternalServerError().WithPayload(reportAPIError(fmt.Sprintf("Unable to convert %s data to jsonapi return format: %s", admmod.TenantStr, err.Error()), startTime, http.StatusInternalServerError, mon.GetAllTenantStr, mon.APICompleted, mon.AdminAPICompleted))
		}

		reportAPICompletionState(startTime, http.StatusOK, mon.GetAllTenantStr, mon.APICompleted, mon.AdminAPICompleted)
		logger.Log.Infof("Retrieved %d %ss", len(result), admmod.TenantStr)
		return admin_provisioning_service.NewGetAllTenantsOK().WithPayload(&converted)
	}
}

// HandleGetTenantIDByAlias - returns the tenant id as a string
func HandleGetTenantIDByAlias(adminDB datastore.AdminServiceDatastore) func(params admin_provisioning_service.GetTenantIDByAliasParams) middleware.Responder {
	return func(params admin_provisioning_service.GetTenantIDByAliasParams) middleware.Responder {
		startTime := time.Now()
		incrementAPICounters(mon.APIRecieved, mon.AdminAPIRecieved)
		logger.Log.Infof("Fetching ID for %s %s", admmod.TenantStr, params.Value)

		// Issue request to DAO Layer
		result, err := adminDB.GetTenantIDByAlias(params.Value)
		if err != nil {
			return admin_provisioning_service.NewGetTenantIDByAliasInternalServerError().WithPayload(reportAPIError(fmt.Sprintf("Unable to retrieve %s ID for %s: %s", admmod.TenantStr, params.Value, err.Error()), startTime, http.StatusInternalServerError, mon.GetTenantIDByAliasStr, mon.APICompleted, mon.AdminAPICompleted))
		}

		reportAPICompletionState(startTime, http.StatusOK, mon.GetTenantIDByAliasStr, mon.APICompleted, mon.AdminAPICompleted)
		logger.Log.Infof("Found ID %s for %s %s", result, admmod.TenantStr, params.Value)
		return admin_provisioning_service.NewGetTenantIDByAliasOK().WithPayload(result)
	}
}

// HandleGetTenantSummaryByAlias - returns the tenant summary for an alias
func HandleGetTenantSummaryByAlias(adminDB datastore.AdminServiceDatastore) func(params admin_provisioning_service.GetTenantSummaryByAliasParams) middleware.Responder {
	return func(params admin_provisioning_service.GetTenantSummaryByAliasParams) middleware.Responder {
		startTime := time.Now()
		incrementAPICounters(mon.APIRecieved, mon.AdminAPIRecieved)
		logger.Log.Infof("Fetching summary for %s %s", admmod.TenantStr, params.Value)

		// Issue request to DAO Layer
		result, err := adminDB.GetTenantIDByAlias(params.Value)
		if err != nil {
			return admin_provisioning_service.NewGetTenantSummaryByAliasInternalServerError().WithPayload(reportAPIError(fmt.Sprintf("Unable to retrieve %s summary for %s: %s", admmod.TenantStr, params.Value, err.Error()), startTime, http.StatusInternalServerError, mon.GetTenantSummaryByAliasStr, mon.APICompleted, mon.AdminAPICompleted))
		}

		summary := swagmodels.TenantSummary{Alias: params.Value, ID: result}

		reportAPICompletionState(startTime, http.StatusOK, mon.GetTenantSummaryByAliasStr, mon.APICompleted, mon.AdminAPICompleted)
		logger.Log.Infof("Successfully retrieved ID %s for alias %s", result, params.Value)
		return admin_provisioning_service.NewGetTenantSummaryByAliasOK().WithPayload(&summary)
	}
}

// HandleCreateIngestionDictionary - creates an ingestion dictionary
func HandleCreateIngestionDictionary(allowedRoles []string, adminDB datastore.AdminServiceDatastore) func(params admin_provisioning_service.CreateIngestionDictionaryParams) middleware.Responder {
	return func(params admin_provisioning_service.CreateIngestionDictionaryParams) middleware.Responder {
		startTime := time.Now()
		incrementAPICounters(mon.APIRecieved, mon.AdminAPIRecieved)
		logger.Log.Infof("Creating %s", admmod.IngestionDictionaryStr)

		if !isRequestAuthorized(params.HTTPRequest, allowedRoles) {
			return admin_provisioning_service.NewCreateIngestionDictionaryForbidden().WithPayload(reportAPIError(fmt.Sprintf("Create %s operation not authorized for role: %s", admmod.IngestionDictionaryStr, params.HTTPRequest.Header.Get(xFwdUserRoles)), startTime, http.StatusForbidden, mon.CreateIngDictStr, mon.APICompleted, mon.AdminAPICompleted))
		}

		// Unmarshal the request
		requestBytes, err := json.Marshal(params.Body)
		if err != nil {
			return admin_provisioning_service.NewCreateIngestionDictionaryBadRequest().WithPayload(reportAPIError(generateErrorMessage(http.StatusBadRequest, err.Error()), startTime, http.StatusForbidden, mon.CreateIngDictStr, mon.APICompleted, mon.AdminAPICompleted))
		}

		data := admmod.IngestionDictionary{}
		err = jsonapi.Unmarshal(requestBytes, &data)
		if err != nil {
			return admin_provisioning_service.NewCreateIngestionDictionaryBadRequest().WithPayload(reportAPIError(generateErrorMessage(http.StatusBadRequest, err.Error()), startTime, http.StatusBadRequest, mon.CreateIngDictStr, mon.APICompleted, mon.AdminAPICompleted))
		}

		err = data.Validate(false)
		if err != nil {
			return admin_provisioning_service.NewCreateIngestionDictionaryBadRequest().WithPayload(reportAPIError(generateErrorMessage(http.StatusBadRequest, err.Error()), startTime, http.StatusBadRequest, mon.CreateIngDictStr, mon.APICompleted, mon.AdminAPICompleted))
		}

		// Issue request to DAO Layer to Create the record
		result, err := adminDB.CreateIngestionDictionary(&data)
		if err != nil {
			return admin_provisioning_service.NewCreateIngestionDictionaryInternalServerError().WithPayload(reportAPIError(fmt.Sprintf("Unable to store %s: %s", admmod.IngestionDictionaryStr, err.Error()), startTime, http.StatusInternalServerError, mon.CreateIngDictStr, mon.APICompleted, mon.AdminAPICompleted))
		}

		converted := swagmodels.JSONAPIIngestionDictionary{}
		err = convertToJsonapiObject(result, &converted)
		if err != nil {
			return admin_provisioning_service.NewCreateIngestionDictionaryInternalServerError().WithPayload(reportAPIError(fmt.Sprintf("Unable to convert %s data to jsonapi return format: %s", admmod.IngestionDictionaryStr, err.Error()), startTime, http.StatusInternalServerError, mon.CreateIngDictStr, mon.APICompleted, mon.AdminAPICompleted))
		}

		reportAPICompletionState(startTime, http.StatusOK, mon.CreateIngDictStr, mon.APICompleted, mon.AdminAPICompleted)
		logger.Log.Infof("Created %s", admmod.IngestionDictionaryStr)
		return admin_provisioning_service.NewCreateIngestionDictionaryOK().WithPayload(&converted)
	}
}

// HandleUpdateIngestionDictionary - updates an ingestion dictionary
func HandleUpdateIngestionDictionary(allowedRoles []string, adminDB datastore.AdminServiceDatastore) func(params admin_provisioning_service.UpdateIngestionDictionaryParams) middleware.Responder {
	return func(params admin_provisioning_service.UpdateIngestionDictionaryParams) middleware.Responder {
		startTime := time.Now()
		incrementAPICounters(mon.APIRecieved, mon.AdminAPIRecieved)
		logger.Log.Infof("Updating %s", admmod.IngestionDictionaryStr)

		if !isRequestAuthorized(params.HTTPRequest, allowedRoles) {
			return admin_provisioning_service.NewUpdateIngestionDictionaryForbidden().WithPayload(reportAPIError(fmt.Sprintf("Update %s operation not authorized for role: %s", admmod.IngestionDictionaryStr, params.HTTPRequest.Header.Get(xFwdUserRoles)), startTime, http.StatusForbidden, mon.UpdateIngDictStr, mon.APICompleted, mon.AdminAPICompleted))
		}

		// Unmarshal the request
		requestBytes, err := json.Marshal(params.Body)
		if err != nil {
			return admin_provisioning_service.NewUpdateIngestionDictionaryBadRequest().WithPayload(reportAPIError(generateErrorMessage(http.StatusBadRequest, err.Error()), startTime, http.StatusForbidden, mon.UpdateIngDictStr, mon.APICompleted, mon.AdminAPICompleted))
		}

		data := admmod.IngestionDictionary{}
		err = jsonapi.Unmarshal(requestBytes, &data)
		if err != nil {
			return admin_provisioning_service.NewUpdateIngestionDictionaryBadRequest().WithPayload(reportAPIError(generateErrorMessage(http.StatusBadRequest, err.Error()), startTime, http.StatusBadRequest, mon.UpdateIngDictStr, mon.APICompleted, mon.AdminAPICompleted))
		}

		err = data.Validate(true)
		if err != nil {
			return admin_provisioning_service.NewUpdateIngestionDictionaryBadRequest().WithPayload(reportAPIError(generateErrorMessage(http.StatusBadRequest, err.Error()), startTime, http.StatusBadRequest, mon.UpdateIngDictStr, mon.APICompleted, mon.AdminAPICompleted))
		}

		// Issue request to DAO Layer to Create the record
		result, err := adminDB.UpdateIngestionDictionary(&data)
		if err != nil {
			return admin_provisioning_service.NewUpdateIngestionDictionaryInternalServerError().WithPayload(reportAPIError(fmt.Sprintf("Unable to store %s: %s", admmod.IngestionDictionaryStr, err.Error()), startTime, http.StatusInternalServerError, mon.UpdateIngDictStr, mon.APICompleted, mon.AdminAPICompleted))
		}

		converted := swagmodels.JSONAPIIngestionDictionary{}
		err = convertToJsonapiObject(result, &converted)
		if err != nil {
			return admin_provisioning_service.NewUpdateIngestionDictionaryInternalServerError().WithPayload(reportAPIError(fmt.Sprintf("Unable to convert %s data to jsonapi return format: %s", admmod.IngestionDictionaryStr, err.Error()), startTime, http.StatusInternalServerError, mon.UpdateIngDictStr, mon.APICompleted, mon.AdminAPICompleted))
		}

		reportAPICompletionState(startTime, http.StatusOK, mon.UpdateIngDictStr, mon.APICompleted, mon.AdminAPICompleted)
		logger.Log.Infof("Updated %s", admmod.IngestionDictionaryStr)
		return admin_provisioning_service.NewUpdateIngestionDictionaryOK().WithPayload(&converted)
	}
}

// HandleGetIngestionDictionary - retrieve an ingestion dictionary
func HandleGetIngestionDictionary(allowedRoles []string, adminDB datastore.AdminServiceDatastore) func(params admin_provisioning_service.GetIngestionDictionaryParams) middleware.Responder {
	return func(params admin_provisioning_service.GetIngestionDictionaryParams) middleware.Responder {
		startTime := time.Now()
		incrementAPICounters(mon.APIRecieved, mon.AdminAPIRecieved)
		logger.Log.Infof("Fetching %s", admmod.IngestionDictionaryStr)

		if !isRequestAuthorized(params.HTTPRequest, allowedRoles) {
			return admin_provisioning_service.NewGetIngestionDictionaryForbidden().WithPayload(reportAPIError(fmt.Sprintf("Get %s operation not authorized for role: %s", admmod.IngestionDictionaryStr, params.HTTPRequest.Header.Get(xFwdUserRoles)), startTime, http.StatusForbidden, mon.GetIngDictStr, mon.APICompleted, mon.AdminAPICompleted))
		}

		// Issue request to DAO Layer
		result, err := adminDB.GetIngestionDictionary()
		if err != nil {
			return admin_provisioning_service.NewGetIngestionDictionaryInternalServerError().WithPayload(reportAPIError(fmt.Sprintf("Unable to retrieve %s: %s", admmod.IngestionDictionaryStr, err.Error()), startTime, http.StatusInternalServerError, mon.GetIngDictStr, mon.APICompleted, mon.AdminAPICompleted))
		}

		converted := swagmodels.JSONAPIIngestionDictionary{}
		err = convertToJsonapiObject(result, &converted)
		if err != nil {
			return admin_provisioning_service.NewGetIngestionDictionaryInternalServerError().WithPayload(reportAPIError(fmt.Sprintf("Unable to convert %s data to jsonapi return format: %s", admmod.IngestionDictionaryStr, err.Error()), startTime, http.StatusInternalServerError, mon.GetIngDictStr, mon.APICompleted, mon.AdminAPICompleted))
		}

		reportAPICompletionState(startTime, http.StatusOK, mon.GetIngDictStr, mon.APICompleted, mon.AdminAPICompleted)
		logger.Log.Infof("Retrieved %s %s", admmod.IngestionDictionaryStr, models.AsJSONString(converted))
		return admin_provisioning_service.NewGetIngestionDictionaryOK().WithPayload(&converted)
	}
}

// HandleDeleteIngestionDictionary - delete an ingestion dictionary
func HandleDeleteIngestionDictionary(allowedRoles []string, adminDB datastore.AdminServiceDatastore) func(params admin_provisioning_service.DeleteIngestionDictionaryParams) middleware.Responder {
	return func(params admin_provisioning_service.DeleteIngestionDictionaryParams) middleware.Responder {
		startTime := time.Now()
		incrementAPICounters(mon.APIRecieved, mon.AdminAPIRecieved)
		logger.Log.Infof("Deleting %s", admmod.IngestionDictionaryStr)

		if !isRequestAuthorized(params.HTTPRequest, allowedRoles) {
			return admin_provisioning_service.NewDeleteIngestionDictionaryForbidden().WithPayload(reportAPIError(fmt.Sprintf("Get %s operation not authorized for role: %s", admmod.IngestionDictionaryStr, params.HTTPRequest.Header.Get(xFwdUserRoles)), startTime, http.StatusForbidden, mon.DeleteIngDictStr, mon.APICompleted, mon.AdminAPICompleted))
		}

		// Issue request to DAO Layer
		result, err := adminDB.DeleteIngestionDictionary()
		if err != nil {
			return admin_provisioning_service.NewDeleteIngestionDictionaryInternalServerError().WithPayload(reportAPIError(fmt.Sprintf("Unable to retrieve %s: %s", admmod.IngestionDictionaryStr, err.Error()), startTime, http.StatusInternalServerError, mon.DeleteIngDictStr, mon.APICompleted, mon.AdminAPICompleted))
		}

		converted := swagmodels.JSONAPIIngestionDictionary{}
		err = convertToJsonapiObject(result, &converted)
		if err != nil {
			return admin_provisioning_service.NewDeleteIngestionDictionaryInternalServerError().WithPayload(reportAPIError(fmt.Sprintf("Unable to convert %s data to jsonapi return format: %s", admmod.IngestionDictionaryStr, err.Error()), startTime, http.StatusInternalServerError, mon.DeleteIngDictStr, mon.APICompleted, mon.AdminAPICompleted))
		}

		reportAPICompletionState(startTime, http.StatusOK, mon.DeleteIngDictStr, mon.APICompleted, mon.AdminAPICompleted)
		logger.Log.Infof("Deleted %s %s", admmod.IngestionDictionaryStr, models.AsJSONString(converted))
		return admin_provisioning_service.NewDeleteIngestionDictionaryOK().WithPayload(&converted)
	}
}

// HandleCreateValidTypes - creates valid types object
func HandleCreateValidTypes(allowedRoles []string, adminDB datastore.AdminServiceDatastore) func(params admin_provisioning_service.CreateValidTypesParams) middleware.Responder {
	return func(params admin_provisioning_service.CreateValidTypesParams) middleware.Responder {
		startTime := time.Now()
		incrementAPICounters(mon.APIRecieved, mon.AdminAPIRecieved)
		logger.Log.Infof("Creating %s", admmod.ValidTypesStr)

		if !isRequestAuthorized(params.HTTPRequest, allowedRoles) {
			return admin_provisioning_service.NewCreateValidTypesForbidden().WithPayload(reportAPIError(fmt.Sprintf("Create %s operation not authorized for role: %s", admmod.ValidTypesStr, params.HTTPRequest.Header.Get(xFwdUserRoles)), startTime, http.StatusForbidden, mon.CreateValidTypesStr, mon.APICompleted, mon.AdminAPICompleted))
		}

		// Unmarshal the request
		requestBytes, err := json.Marshal(params.Body)
		if err != nil {
			return admin_provisioning_service.NewCreateValidTypesForbidden().WithPayload(reportAPIError(generateErrorMessage(http.StatusBadRequest, err.Error()), startTime, http.StatusForbidden, mon.CreateValidTypesStr, mon.APICompleted, mon.AdminAPICompleted))
		}

		data := admmod.ValidTypes{}
		err = jsonapi.Unmarshal(requestBytes, &data)
		if err != nil {
			return admin_provisioning_service.NewCreateValidTypesForbidden().WithPayload(reportAPIError(generateErrorMessage(http.StatusBadRequest, err.Error()), startTime, http.StatusBadRequest, mon.CreateValidTypesStr, mon.APICompleted, mon.AdminAPICompleted))
		}

		err = data.Validate(false)
		if err != nil {
			return admin_provisioning_service.NewCreateValidTypesForbidden().WithPayload(reportAPIError(generateErrorMessage(http.StatusBadRequest, err.Error()), startTime, http.StatusBadRequest, mon.CreateValidTypesStr, mon.APICompleted, mon.AdminAPICompleted))
		}

		// Issue request to DAO Layer to Create the record
		result, err := adminDB.CreateValidTypes(&data)
		if err != nil {
			return admin_provisioning_service.NewCreateValidTypesInternalServerError().WithPayload(reportAPIError(fmt.Sprintf("Unable to store %s: %s", admmod.ValidTypesStr, err.Error()), startTime, http.StatusInternalServerError, mon.CreateValidTypesStr, mon.APICompleted, mon.AdminAPICompleted))
		}

		converted := swagmodels.JSONAPIValidTypes{}
		err = convertToJsonapiObject(result, &converted)
		if err != nil {
			return admin_provisioning_service.NewCreateValidTypesInternalServerError().WithPayload(reportAPIError(fmt.Sprintf("Unable to convert %s data to jsonapi return format: %s", admmod.ValidTypesStr, err.Error()), startTime, http.StatusInternalServerError, mon.CreateValidTypesStr, mon.APICompleted, mon.AdminAPICompleted))
		}

		reportAPICompletionState(startTime, http.StatusOK, mon.CreateValidTypesStr, mon.APICompleted, mon.AdminAPICompleted)
		logger.Log.Infof("Created %s", admmod.ValidTypesStr)
		return admin_provisioning_service.NewCreateValidTypesOK().WithPayload(&converted)
	}
}

// HandleUpdateValidTypes - updates valid types object
func HandleUpdateValidTypes(allowedRoles []string, adminDB datastore.AdminServiceDatastore) func(params admin_provisioning_service.UpdateValidTypesParams) middleware.Responder {
	return func(params admin_provisioning_service.UpdateValidTypesParams) middleware.Responder {
		startTime := time.Now()
		incrementAPICounters(mon.APIRecieved, mon.AdminAPIRecieved)
		logger.Log.Infof("Updating %s", admmod.ValidTypesStr)

		if !isRequestAuthorized(params.HTTPRequest, allowedRoles) {
			return admin_provisioning_service.NewCreateValidTypesForbidden().WithPayload(reportAPIError(fmt.Sprintf("Update %s operation not authorized for role: %s", admmod.ValidTypesStr, params.HTTPRequest.Header.Get(xFwdUserRoles)), startTime, http.StatusForbidden, mon.UpdateValidTypesStr, mon.APICompleted, mon.AdminAPICompleted))
		}

		// Unmarshal the request
		requestBytes, err := json.Marshal(params.Body)
		if err != nil {
			return admin_provisioning_service.NewCreateValidTypesForbidden().WithPayload(reportAPIError(generateErrorMessage(http.StatusBadRequest, err.Error()), startTime, http.StatusForbidden, mon.UpdateValidTypesStr, mon.APICompleted, mon.AdminAPICompleted))
		}

		data := admmod.ValidTypes{}
		err = jsonapi.Unmarshal(requestBytes, &data)
		if err != nil {
			return admin_provisioning_service.NewCreateValidTypesForbidden().WithPayload(reportAPIError(generateErrorMessage(http.StatusBadRequest, err.Error()), startTime, http.StatusBadRequest, mon.UpdateValidTypesStr, mon.APICompleted, mon.AdminAPICompleted))
		}

		err = data.Validate(true)
		if err != nil {
			return admin_provisioning_service.NewCreateValidTypesForbidden().WithPayload(reportAPIError(generateErrorMessage(http.StatusBadRequest, err.Error()), startTime, http.StatusBadRequest, mon.UpdateValidTypesStr, mon.APICompleted, mon.AdminAPICompleted))
		}

		// Issue request to DAO Layer to Create the record
		result, err := adminDB.UpdateValidTypes(&data)
		if err != nil {
			return admin_provisioning_service.NewCreateValidTypesInternalServerError().WithPayload(reportAPIError(fmt.Sprintf("Unable to store %s: %s", admmod.ValidTypesStr, err.Error()), startTime, http.StatusInternalServerError, mon.UpdateValidTypesStr, mon.APICompleted, mon.AdminAPICompleted))
		}

		converted := swagmodels.JSONAPIValidTypes{}
		err = convertToJsonapiObject(result, &converted)
		if err != nil {
			return admin_provisioning_service.NewCreateValidTypesInternalServerError().WithPayload(reportAPIError(fmt.Sprintf("Unable to convert %s data to jsonapi return format: %s", admmod.ValidTypesStr, err.Error()), startTime, http.StatusInternalServerError, mon.UpdateValidTypesStr, mon.APICompleted, mon.AdminAPICompleted))
		}

		reportAPICompletionState(startTime, http.StatusOK, mon.UpdateValidTypesStr, mon.APICompleted, mon.AdminAPICompleted)
		logger.Log.Infof("Updated %s", admmod.ValidTypesStr)
		return admin_provisioning_service.NewUpdateValidTypesOK().WithPayload(&converted)
	}
}

// HandleGetValidTypes - retrieve an ingestion dictionary
func HandleGetValidTypes(allowedRoles []string, adminDB datastore.AdminServiceDatastore) func(params admin_provisioning_service.GetValidTypesParams) middleware.Responder {
	return func(params admin_provisioning_service.GetValidTypesParams) middleware.Responder {
		startTime := time.Now()
		incrementAPICounters(mon.APIRecieved, mon.AdminAPIRecieved)
		logger.Log.Infof("Fetching %s", admmod.ValidTypesStr)

		if !isRequestAuthorized(params.HTTPRequest, allowedRoles) {
			return admin_provisioning_service.NewGetIngestionDictionaryForbidden().WithPayload(reportAPIError(fmt.Sprintf("Get %s operation not authorized for role: %s", admmod.ValidTypesStr, params.HTTPRequest.Header.Get(xFwdUserRoles)), startTime, http.StatusForbidden, mon.GetValidTypesStr, mon.APICompleted, mon.AdminAPICompleted))
		}

		// Issue request to DAO Layer
		result, err := adminDB.GetValidTypes()
		if err != nil {
			return admin_provisioning_service.NewGetIngestionDictionaryInternalServerError().WithPayload(reportAPIError(fmt.Sprintf("Unable to retrieve %s: %s", admmod.ValidTypesStr, err.Error()), startTime, http.StatusInternalServerError, mon.GetValidTypesStr, mon.APICompleted, mon.AdminAPICompleted))
		}

		converted := swagmodels.JSONAPIValidTypes{}
		err = convertToJsonapiObject(result, &converted)
		if err != nil {
			return admin_provisioning_service.NewGetIngestionDictionaryInternalServerError().WithPayload(reportAPIError(fmt.Sprintf("Unable to convert %s data to jsonapi return format: %s", admmod.ValidTypesStr, err.Error()), startTime, http.StatusInternalServerError, mon.GetValidTypesStr, mon.APICompleted, mon.AdminAPICompleted))
		}

		reportAPICompletionState(startTime, http.StatusOK, mon.GetValidTypesStr, mon.APICompleted, mon.AdminAPICompleted)
		logger.Log.Infof("Retrieved %s %s", admmod.ValidTypesStr, models.AsJSONString(converted))
		return admin_provisioning_service.NewGetValidTypesOK().WithPayload(&converted)
	}
}

// HandleDeleteValidTypes - delete valid types object
func HandleDeleteValidTypes(allowedRoles []string, adminDB datastore.AdminServiceDatastore) func(params admin_provisioning_service.DeleteValidTypesParams) middleware.Responder {
	return func(params admin_provisioning_service.DeleteValidTypesParams) middleware.Responder {
		startTime := time.Now()
		incrementAPICounters(mon.APIRecieved, mon.AdminAPIRecieved)
		logger.Log.Infof("Deleting %s", admmod.ValidTypesStr)

		if !isRequestAuthorized(params.HTTPRequest, allowedRoles) {
			return admin_provisioning_service.NewDeleteIngestionDictionaryForbidden().WithPayload(reportAPIError(fmt.Sprintf("Get %s operation not authorized for role: %s", admmod.ValidTypesStr, params.HTTPRequest.Header.Get(xFwdUserRoles)), startTime, http.StatusForbidden, mon.DeleteValidTypesStr, mon.APICompleted, mon.AdminAPICompleted))
		}

		// Issue request to DAO Layer
		result, err := adminDB.DeleteValidTypes()
		if err != nil {
			return admin_provisioning_service.NewDeleteIngestionDictionaryInternalServerError().WithPayload(reportAPIError(fmt.Sprintf("Unable to retrieve %s: %s", admmod.ValidTypesStr, err.Error()), startTime, http.StatusInternalServerError, mon.DeleteValidTypesStr, mon.APICompleted, mon.AdminAPICompleted))
		}

		converted := swagmodels.JSONAPIValidTypes{}
		err = convertToJsonapiObject(result, &converted)
		if err != nil {
			return admin_provisioning_service.NewDeleteIngestionDictionaryInternalServerError().WithPayload(reportAPIError(fmt.Sprintf("Unable to convert %s data to jsonapi return format: %s", admmod.ValidTypesStr, err.Error()), startTime, http.StatusInternalServerError, mon.DeleteValidTypesStr, mon.APICompleted, mon.AdminAPICompleted))
		}

		reportAPICompletionState(startTime, http.StatusOK, mon.DeleteValidTypesStr, mon.APICompleted, mon.AdminAPICompleted)
		logger.Log.Infof("Deleted %s %s", admmod.ValidTypesStr, models.AsJSONString(converted))
		return admin_provisioning_service.NewDeleteValidTypesOK().WithPayload(&converted)
	}
}
