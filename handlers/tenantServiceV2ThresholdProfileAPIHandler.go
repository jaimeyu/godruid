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
	tenmod "github.com/accedian/adh-gather/models/tenant"
	mon "github.com/accedian/adh-gather/monitoring"
	"github.com/accedian/adh-gather/restapi/operations/tenant_provisioning_service_v2"
	"github.com/go-openapi/runtime/middleware"
)

// HandleCreateThresholdProfileV2 - create a new Threshold Profile
func HandleCreateThresholdProfileV2(allowedRoles []string, tenantDB datastore.TenantServiceDatastore) func(params tenant_provisioning_service_v2.CreateThresholdProfileV2Params) middleware.Responder {
	return func(params tenant_provisioning_service_v2.CreateThresholdProfileV2Params) middleware.Responder {

		// Do the work
		startTime, responseCode, response, err := doCreateThresholdProfileV2(allowedRoles, tenantDB, params)

		// Success Response
		if responseCode == http.StatusCreated {
			return tenant_provisioning_service_v2.NewCreateThresholdProfileV2Created().WithPayload(response)
		}

		// Error Responses
		errorMessage := reportAPIError(err.Error(), startTime, responseCode, mon.CreateThrPrfStr, mon.APICompleted, mon.TenantAPICompleted)
		switch responseCode {
		case http.StatusForbidden:
			return tenant_provisioning_service_v2.NewCreateThresholdProfileV2Forbidden().WithPayload(errorMessage)
		case http.StatusBadRequest:
			return tenant_provisioning_service_v2.NewCreateThresholdProfileV2BadRequest().WithPayload(errorMessage)
		case http.StatusConflict:
			return tenant_provisioning_service_v2.NewCreateThresholdProfileV2Conflict().WithPayload(errorMessage)
		default:
			return tenant_provisioning_service_v2.NewCreateThresholdProfileV2InternalServerError().WithPayload(errorMessage)
		}
	}
}

// HandleGetThresholdProfileV2 - retrieve a Threshold Profile by the ID.
func HandleGetThresholdProfileV2(allowedRoles []string, tenantDB datastore.TenantServiceDatastore) func(params tenant_provisioning_service_v2.GetThresholdProfileV2Params) middleware.Responder {
	return func(params tenant_provisioning_service_v2.GetThresholdProfileV2Params) middleware.Responder {

		// Do the work
		startTime, responseCode, response, err := doGetThresholdProfileV2(allowedRoles, tenantDB, params)

		// Success Response
		if responseCode == http.StatusOK {
			return tenant_provisioning_service_v2.NewGetThresholdProfileV2OK().WithPayload(response)
		}

		// Error Responses
		errorMessage := reportAPIError(err.Error(), startTime, responseCode, mon.GetThrPrfStr, mon.APICompleted, mon.TenantAPICompleted)
		switch responseCode {
		case http.StatusForbidden:
			return tenant_provisioning_service_v2.NewGetThresholdProfileV2Forbidden().WithPayload(errorMessage)
		case http.StatusNotFound:
			return tenant_provisioning_service_v2.NewGetThresholdProfileV2NotFound().WithPayload(errorMessage)
		default:
			return tenant_provisioning_service_v2.NewGetThresholdProfileV2InternalServerError().WithPayload(errorMessage)
		}

	}
}

// HandleUpdateThresholdProfileV2 - update a ThresholdProfile record
func HandleUpdateThresholdProfileV2(allowedRoles []string, tenantDB datastore.TenantServiceDatastore) func(params tenant_provisioning_service_v2.UpdateThresholdProfileV2Params) middleware.Responder {
	return func(params tenant_provisioning_service_v2.UpdateThresholdProfileV2Params) middleware.Responder {

		// Do the work
		startTime, responseCode, response, err := doUpdateThresholdProfileV2(allowedRoles, tenantDB, params)

		// Success Response
		if responseCode == http.StatusOK {
			return tenant_provisioning_service_v2.NewUpdateThresholdProfileV2OK().WithPayload(response)
		}

		// Error Responses
		errorMessage := reportAPIError(err.Error(), startTime, responseCode, mon.UpdateThrPrfStr, mon.APICompleted, mon.TenantAPICompleted)
		switch responseCode {
		case http.StatusForbidden:
			return tenant_provisioning_service_v2.NewUpdateThresholdProfileV2Forbidden().WithPayload(errorMessage)
		case http.StatusBadRequest:
			return tenant_provisioning_service_v2.NewUpdateThresholdProfileV2BadRequest().WithPayload(errorMessage)
		case http.StatusNotFound:
			return tenant_provisioning_service_v2.NewUpdateThresholdProfileV2NotFound().WithPayload(errorMessage)
		case http.StatusConflict:
			return tenant_provisioning_service_v2.NewUpdateThresholdProfileV2Conflict().WithPayload(errorMessage)
		default:
			return tenant_provisioning_service_v2.NewUpdateThresholdProfileV2InternalServerError().WithPayload(errorMessage)
		}
	}
}

// HandleDeleteThresholdProfileV2 - delete a Threshold Profile by the ID.
func HandleDeleteThresholdProfileV2(allowedRoles []string, tenantDB datastore.TenantServiceDatastore) func(params tenant_provisioning_service_v2.DeleteThresholdProfileV2Params) middleware.Responder {
	return func(params tenant_provisioning_service_v2.DeleteThresholdProfileV2Params) middleware.Responder {

		// Do the work
		startTime, responseCode, response, err := doDeleteThresholdProfileV2(allowedRoles, tenantDB, params)

		// Success Response
		if responseCode == http.StatusOK {
			return tenant_provisioning_service_v2.NewDeleteThresholdProfileV2OK().WithPayload(response)
		}

		// Error Responses
		errorMessage := reportAPIError(err.Error(), startTime, responseCode, mon.DeleteThrPrfStr, mon.APICompleted, mon.TenantAPICompleted)
		switch responseCode {
		case http.StatusForbidden:
			return tenant_provisioning_service_v2.NewDeleteThresholdProfileV2Forbidden().WithPayload(errorMessage)
		case http.StatusNotFound:
			return tenant_provisioning_service_v2.NewDeleteThresholdProfileV2NotFound().WithPayload(errorMessage)
		default:
			return tenant_provisioning_service_v2.NewDeleteThresholdProfileV2InternalServerError().WithPayload(errorMessage)
		}

	}
}

// HandleGetAllThresholdProfilesV2 - retrieve all Threshold Profiles
func HandleGetAllThresholdProfilesV2(allowedRoles []string, tenantDB datastore.TenantServiceDatastore) func(params tenant_provisioning_service_v2.GetAllThresholdProfilesV2Params) middleware.Responder {
	return func(params tenant_provisioning_service_v2.GetAllThresholdProfilesV2Params) middleware.Responder {

		// Do the work
		startTime, responseCode, response, err := doGetAllThresholdProfilesV2(allowedRoles, tenantDB, params)

		// Success Response
		if responseCode == http.StatusOK {
			return tenant_provisioning_service_v2.NewGetAllThresholdProfilesV2OK().WithPayload(response)
		}

		// Error Responses
		errorMessage := reportAPIError(err.Error(), startTime, responseCode, mon.GetAllThrPrfStr, mon.APICompleted, mon.TenantAPICompleted)
		switch responseCode {
		case http.StatusForbidden:
			return tenant_provisioning_service_v2.NewGetAllThresholdProfilesV2Forbidden().WithPayload(errorMessage)
		case http.StatusNotFound:
			return tenant_provisioning_service_v2.NewGetAllThresholdProfilesV2NotFound().WithPayload(errorMessage)
		default:
			return tenant_provisioning_service_v2.NewGetAllThresholdProfilesV2InternalServerError().WithPayload(errorMessage)
		}

	}
}

func doCreateThresholdProfileV2(allowedRoles []string, tenantDB datastore.TenantServiceDatastore, params tenant_provisioning_service_v2.CreateThresholdProfileV2Params) (time.Time, int, *swagmodels.ThresholdProfileResponse, error) {
	tenantID := params.HTTPRequest.Header.Get(XFwdTenantId)
	isAuthorized, startTime := authorizeRequest(fmt.Sprintf("Creating %s for tenant %s", tenmod.TenantThresholdProfileStr, tenantID), params.HTTPRequest, allowedRoles, mon.APIRecieved, mon.TenantAPIRecieved)

	if !isAuthorized {
		return startTime, http.StatusForbidden, nil, fmt.Errorf("Create %s operation not authorized for role: %s", tenmod.TenantThresholdProfileStr, params.HTTPRequest.Header.Get(XFwdUserRoles))
	}

	// Unmarshal the request
	requestBytes, err := json.Marshal(params.Body)
	if err != nil {
		return startTime, http.StatusBadRequest, nil, err
	}

	data := tenmod.ThresholdProfile{}
	err = jsonapi.Unmarshal(requestBytes, &data)
	if err != nil {
		return startTime, http.StatusBadRequest, nil, err
	}

	data.TenantID = tenantID

	// Issue request to DAO Layer to Create Record
	result, err := tenantDB.CreateTenantThresholdProfile(&data)
	if err != nil {
		if strings.Contains(err.Error(), datastore.ConflictErrorStr) {
			return startTime, http.StatusConflict, nil, err
		}
		return startTime, http.StatusInternalServerError, nil, fmt.Errorf("Unable to store %s: %s", tenmod.TenantThresholdProfileStr, err.Error())
	}

	converted := swagmodels.ThresholdProfileResponse{}
	err = convertToJsonapiObject(result, &converted)
	if err != nil {
		return startTime, http.StatusInternalServerError, nil, fmt.Errorf("Unable to convert %s data to jsonapi return format: %s", tenmod.TenantThresholdProfileStr, err.Error())
	}

	reportAPICompletionState(startTime, http.StatusOK, mon.CreateThrPrfStr, mon.APICompleted, mon.TenantAPICompleted)
	logger.Log.Infof("Created %s %s", tenmod.TenantThresholdProfileStr, models.AsJSONString(converted))
	return startTime, http.StatusCreated, &converted, nil
}

func doGetThresholdProfileV2(allowedRoles []string, tenantDB datastore.TenantServiceDatastore, params tenant_provisioning_service_v2.GetThresholdProfileV2Params) (time.Time, int, *swagmodels.ThresholdProfileResponse, error) {
	tenantID := params.HTTPRequest.Header.Get(XFwdTenantId)
	isAuthorized, startTime := authorizeRequest(fmt.Sprintf("Fetching %s %s for %s %s", tenmod.TenantThresholdProfileStr, params.ThrPrfID, admmod.TenantStr, tenantID), params.HTTPRequest, allowedRoles, mon.APIRecieved, mon.AdminAPIRecieved)

	if !isAuthorized {
		return startTime, http.StatusForbidden, nil, fmt.Errorf("Fetch %s operation not authorized for role: %s", tenmod.TenantThresholdProfileStr, params.HTTPRequest.Header.Get(XFwdUserRoles))
	}

	// Issue request to DAO Layer
	result, err := tenantDB.GetTenantThresholdProfile(tenantID, params.ThrPrfID)
	if err != nil {
		if strings.Contains(err.Error(), datastore.NotFoundStr) {
			return startTime, http.StatusNotFound, nil, err
		}

		return startTime, http.StatusInternalServerError, nil, fmt.Errorf("Unable to retrieve %s: %s", tenmod.TenantThresholdProfileStr, err.Error())
	}

	converted := swagmodels.ThresholdProfileResponse{}
	err = convertToJsonapiObject(result, &converted)
	if err != nil {
		return startTime, http.StatusInternalServerError, nil, fmt.Errorf("Unable to convert %s data to jsonapi return format: %s", tenmod.TenantThresholdProfileStr, err.Error())
	}

	reportAPICompletionState(startTime, http.StatusOK, mon.GetThrPrfStr, mon.APICompleted, mon.TenantAPICompleted)
	logger.Log.Infof("Retrieved %s %s", tenmod.TenantThresholdProfileStr, models.AsJSONString(converted))
	return startTime, http.StatusOK, &converted, nil
}

func doUpdateThresholdProfileV2(allowedRoles []string, tenantDB datastore.TenantServiceDatastore, params tenant_provisioning_service_v2.UpdateThresholdProfileV2Params) (time.Time, int, *swagmodels.ThresholdProfileResponse, error) {
	tenantID := params.HTTPRequest.Header.Get(XFwdTenantId)
	isAuthorized, startTime := authorizeRequest(fmt.Sprintf("Updating %s %s for %s %s", tenmod.TenantThresholdProfileStr, params.ThrPrfID, admmod.TenantStr, tenantID), params.HTTPRequest, allowedRoles, mon.APIRecieved, mon.AdminAPIRecieved)

	if !isAuthorized {
		return startTime, http.StatusForbidden, nil, fmt.Errorf("Update %s operation not authorized for role: %s", tenmod.TenantThresholdProfileStr, params.HTTPRequest.Header.Get(XFwdUserRoles))
	}

	// Retrieve tne request bytes from the payload in order to convert it to a map
	patchRequestBytes, err := json.Marshal(params.Body)
	if err != nil || params.Body == nil {
		return startTime, http.StatusBadRequest, nil, fmt.Errorf("Invalid request body: %s", models.AsJSONString(params.Body))
	}

	// Attempt to retrieve the record that we are trying to patch from the DB in order to do a merge
	fetched, err := tenantDB.GetTenantThresholdProfile(tenantID, params.ThrPrfID)
	if err != nil {
		if strings.Contains(err.Error(), datastore.NotFoundStr) {
			return startTime, http.StatusNotFound, nil, err
		}

		return startTime, http.StatusInternalServerError, nil, fmt.Errorf("Unable to fetch %s: %s", tenmod.TenantThresholdProfileStr, err.Error())
	}

	// Merge the attributes passed in with the patch request to the record fetched from the datastore
	patched := &tenmod.ThresholdProfile{}
	if err := models.MergeObjWithMap(patched, fetched, patchRequestBytes); err != nil {
		return startTime, http.StatusInternalServerError, nil, fmt.Errorf("Unable to patch %s with id %s: %s", tenmod.TenantThresholdProfileStr, params.ThrPrfID, err.Error())
	}
	patched.TenantID = tenantID

	// Finally update the record in the datastore with the merged map and fetched tenant
	result, err := tenantDB.UpdateTenantThresholdProfile(patched)
	if err != nil {
		if strings.Contains(err.Error(), datastore.ConflictErrorStr) {
			return startTime, http.StatusConflict, nil, err
		}

		return startTime, http.StatusInternalServerError, nil, err
	}

	converted := swagmodels.ThresholdProfileResponse{}
	err = convertToJsonapiObject(result, &converted)
	if err != nil {
		return startTime, http.StatusInternalServerError, nil, fmt.Errorf("Unable to convert %s data to jsonapi return format: %s", tenmod.TenantThresholdProfileStr, err.Error())
	}

	reportAPICompletionState(startTime, http.StatusOK, mon.UpdateThrPrfStr, mon.APICompleted, mon.TenantAPICompleted)
	logger.Log.Infof("Updated %s %s", tenmod.TenantThresholdProfileStr, models.AsJSONString(converted))
	return startTime, http.StatusOK, &converted, nil
}

func doDeleteThresholdProfileV2(allowedRoles []string, tenantDB datastore.TenantServiceDatastore, params tenant_provisioning_service_v2.DeleteThresholdProfileV2Params) (time.Time, int, *swagmodels.ThresholdProfileResponse, error) {
	tenantID := params.HTTPRequest.Header.Get(XFwdTenantId)
	isAuthorized, startTime := authorizeRequest(fmt.Sprintf("Deleting %s %s for %s %s", tenmod.TenantThresholdProfileStr, params.ThrPrfID, admmod.TenantStr, tenantID), params.HTTPRequest, allowedRoles, mon.APIRecieved, mon.AdminAPIRecieved)

	if !isAuthorized {
		return startTime, http.StatusForbidden, nil, fmt.Errorf("Delete %s operation not authorized for role: %s", tenmod.TenantThresholdProfileStr, params.HTTPRequest.Header.Get(XFwdUserRoles))
	}

	// Issue request to DAO Layer
	result, err := tenantDB.DeleteTenantThresholdProfile(tenantID, params.ThrPrfID)
	if err != nil {
		if strings.Contains(err.Error(), datastore.NotFoundStr) {
			return startTime, http.StatusNotFound, nil, err
		}

		return startTime, http.StatusInternalServerError, nil, fmt.Errorf("Unable to delete %s: %s", tenmod.TenantThresholdProfileStr, err.Error())
	}

	converted := swagmodels.ThresholdProfileResponse{}
	err = convertToJsonapiObject(result, &converted)
	if err != nil {
		return startTime, http.StatusInternalServerError, nil, fmt.Errorf("Unable to convert %s data to jsonapi return format: %s", tenmod.TenantThresholdProfileStr, err.Error())
	}

	reportAPICompletionState(startTime, http.StatusOK, mon.DeleteThrPrfStr, mon.APICompleted, mon.TenantAPICompleted)
	logger.Log.Infof("Deleted %s %s", tenmod.TenantThresholdProfileStr, models.AsJSONString(converted))
	return startTime, http.StatusOK, &converted, nil
}

func doGetAllThresholdProfilesV2(allowedRoles []string, tenantDB datastore.TenantServiceDatastore, params tenant_provisioning_service_v2.GetAllThresholdProfilesV2Params) (time.Time, int, *swagmodels.ThresholdProfileListResponse, error) {
	tenantID := params.HTTPRequest.Header.Get(XFwdTenantId)
	isAuthorized, startTime := authorizeRequest(fmt.Sprintf("Fetching %s list for %s %s", tenmod.TenantThresholdProfileStr, admmod.TenantStr, tenantID), params.HTTPRequest, allowedRoles, mon.APIRecieved, mon.AdminAPIRecieved)

	if !isAuthorized {
		return startTime, http.StatusForbidden, nil, fmt.Errorf("Fetch %s operation not authorized for role: %s", tenmod.TenantThresholdProfileStr, params.HTTPRequest.Header.Get(XFwdUserRoles))
	}

	// Issue request to DAO Layer
	result, err := tenantDB.GetAllTenantThresholdProfile(tenantID)
	if err != nil {
		if strings.Contains(err.Error(), datastore.NotFoundStr) {
			return startTime, http.StatusNotFound, nil, err
		}

		return startTime, http.StatusInternalServerError, nil, fmt.Errorf("Unable to retrieve %s list: %s", tenmod.TenantThresholdProfileStr, err.Error())
	}

	if len(result) == 0 {
		return startTime, http.StatusNotFound, nil, fmt.Errorf(datastore.NotFoundStr)
	}

	converted := swagmodels.ThresholdProfileListResponse{}
	err = convertToJsonapiObject(result, &converted)
	if err != nil {
		return startTime, http.StatusInternalServerError, nil, fmt.Errorf("Unable to convert %s data to jsonapi return format: %s", tenmod.TenantThresholdProfileStr, err.Error())
	}

	reportAPICompletionState(startTime, http.StatusOK, mon.GetAllThrPrfStr, mon.APICompleted, mon.TenantAPICompleted)
	logger.Log.Infof("Retrieved %d %ss", len(result), tenmod.TenantThresholdProfileStr)
	return startTime, http.StatusOK, &converted, nil
}
