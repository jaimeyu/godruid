package handlers_test

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/accedian/adh-gather/swagmodels"

	"github.com/accedian/adh-gather/handlers"
	"github.com/accedian/adh-gather/logger"
	"github.com/accedian/adh-gather/restapi/operations/tenant_provisioning_service_v2"
	"github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {

	err := setupTestDatastore()
	if err != nil {
		logger.Log.Fatalf("Unable to setup datastore for Data Cleaning Profile tests: %s", err.Error())
	}

	code := m.Run()

	err = destroyTestDatastore()
	if err != nil {
		logger.Log.Errorf("Unable to remove test datastore for Data Cleaning Profile tests: %s", err.Error())
	}

	os.Exit(code)
}

func TestGetDataCleaningProfile(t *testing.T) {

	tenantDescriptor1 := getRandomTenantDescriptor()

	createdTenant, err := adminDB.CreateTenant(tenantDescriptor1)
	assert.Nil(t, err)
	assert.NotNil(t, createdTenant)
	assert.Equal(t, tenantDescriptor1.Name, createdTenant.Name)
	assert.Equal(t, tenantDescriptor1.URLSubdomain, createdTenant.URLSubdomain)
	assert.True(t, len(createdTenant.ID) > 0)

	// Make sure there are no existing Data Cleaning Profiles
	existing := handlers.HandleGetDataCleaningProfilesV2(handlers.SkylightAndTenantAdminRoles, tenantDB)(tenant_provisioning_service_v2.GetDataCleaningProfilesParams{HTTPRequest: createHttpRequest(createdTenant.ID, handlers.UserRoleSkylight)})
	castedResponse := existing.(*tenant_provisioning_service_v2.GetDataCleaningProfilesNotFound)
	assert.NotNil(t, castedResponse)

	// Create a record
	created := handlers.HandleCreateDataCleaningProfileV2(handlers.SkylightAndTenantAdminRoles, tenantDB)(tenant_provisioning_service_v2.CreateDataCleaningProfileParams{
		Body: &swagmodels.DataCleaningProfileCreateRequest{
			Data: createRandomDataCleaningProfileCreateRequest(),
		},
		HTTPRequest: createHttpRequest(createdTenant.ID, handlers.UserRoleSkylight)})
	castedCreate := created.(*tenant_provisioning_service_v2.CreateDataCleaningProfileCreated)
	assert.NotNil(t, castedCreate)
	assert.True(t, len(*castedCreate.Payload.Data.ID) != 0)

	// Make sure we can retrieve the record:
	fetched := handlers.HandleGetDataCleaningProfileV2(handlers.SkylightAndTenantAdminRoles, tenantDB)(tenant_provisioning_service_v2.GetDataCleaningProfileParams{
		ProfileID:   *castedCreate.Payload.Data.ID,
		HTTPRequest: createHttpRequest(createdTenant.ID, handlers.UserRoleSkylight)})
	castedFetch := fetched.(*tenant_provisioning_service_v2.GetDataCleaningProfileOK)
	assert.NotNil(t, castedFetch)
	assert.Equal(t, castedCreate.Payload.Data, castedFetch.Payload.Data)

	// Test fetching a record that does not exist
	fakeID := "porkypork"
	fetchedDNE := handlers.HandleGetDataCleaningProfileV2(handlers.SkylightAndTenantAdminRoles, tenantDB)(tenant_provisioning_service_v2.GetDataCleaningProfileParams{
		ProfileID:   fakeID,
		HTTPRequest: createHttpRequest(createdTenant.ID, handlers.UserRoleSkylight)})
	castedFetchDNE := fetchedDNE.(*tenant_provisioning_service_v2.GetDataCleaningProfileNotFound)
	assert.NotNil(t, castedFetchDNE)

	// Make sure Tenant User is not able to make this call
	fetchedFrbdn := handlers.HandleGetDataCleaningProfileV2(handlers.SkylightAndTenantAdminRoles, tenantDB)(tenant_provisioning_service_v2.GetDataCleaningProfileParams{
		ProfileID:   fakeID,
		HTTPRequest: createHttpRequest(createdTenant.ID, handlers.UserRoleTenantUser)})
	castedFetchFrbdn := fetchedFrbdn.(*tenant_provisioning_service_v2.GetDataCleaningProfileForbidden)
	assert.NotNil(t, castedFetchFrbdn)

	adminDB.DeleteTenant(createdTenant.ID)
}

func TestUpdateDataCleaningProfile(t *testing.T) {

	tenantDescriptor1 := getRandomTenantDescriptor()

	createdTenant, err := adminDB.CreateTenant(tenantDescriptor1)
	assert.Nil(t, err)
	assert.NotNil(t, createdTenant)
	assert.Equal(t, tenantDescriptor1.Name, createdTenant.Name)
	assert.Equal(t, tenantDescriptor1.URLSubdomain, createdTenant.URLSubdomain)
	assert.True(t, len(createdTenant.ID) > 0)

	// Make sure there are no existing Data Cleaning Profiles
	existing := handlers.HandleGetDataCleaningProfilesV2(handlers.SkylightAndTenantAdminRoles, tenantDB)(tenant_provisioning_service_v2.GetDataCleaningProfilesParams{HTTPRequest: createHttpRequest(createdTenant.ID, handlers.UserRoleSkylight)})
	castedResponse := existing.(*tenant_provisioning_service_v2.GetDataCleaningProfilesNotFound)
	assert.NotNil(t, castedResponse)

	// Create a record
	created := handlers.HandleCreateDataCleaningProfileV2(handlers.SkylightAndTenantAdminRoles, tenantDB)(tenant_provisioning_service_v2.CreateDataCleaningProfileParams{
		Body: &swagmodels.DataCleaningProfileCreateRequest{
			Data: createRandomDataCleaningProfileCreateRequest(),
		},
		HTTPRequest: createHttpRequest(createdTenant.ID, handlers.UserRoleSkylight)})
	castedCreate := created.(*tenant_provisioning_service_v2.CreateDataCleaningProfileCreated)
	assert.NotNil(t, castedCreate)
	assert.True(t, len(*castedCreate.Payload.Data.ID) != 0)

	// Make sure we can retrieve the record:
	fetched := handlers.HandleGetDataCleaningProfileV2(handlers.SkylightAndTenantAdminRoles, tenantDB)(tenant_provisioning_service_v2.GetDataCleaningProfileParams{
		ProfileID:   *castedCreate.Payload.Data.ID,
		HTTPRequest: createHttpRequest(createdTenant.ID, handlers.UserRoleSkylight)})
	castedFetch := fetched.(*tenant_provisioning_service_v2.GetDataCleaningProfileOK)
	assert.NotNil(t, castedFetch)
	assert.Equal(t, castedCreate.Payload.Data, castedFetch.Payload.Data)

	// Make an update
	updateRequest := swagmodels.DataCleaningProfileUpdateRequest{}
	fetchedInBytes, err := json.Marshal(castedFetch.Payload)
	assert.Nil(t, err)
	err = json.Unmarshal(fetchedInBytes, &updateRequest)
	assert.Nil(t, err)

	knownString := "hopeIgetThisBack"
	updateRequest.Data.Attributes.Rules[0].MetricVendor = &knownString

	updated := handlers.HandleUpdateDataCleaningProfileV2(handlers.SkylightAndTenantAdminRoles, tenantDB)(tenant_provisioning_service_v2.UpdateDataCleaningProfileParams{
		ProfileID:   *updateRequest.Data.ID,
		Body:        &updateRequest,
		HTTPRequest: createHttpRequest(createdTenant.ID, handlers.UserRoleSkylight)})
	castedUpdate := updated.(*tenant_provisioning_service_v2.UpdateDataCleaningProfileOK)
	assert.NotNil(t, castedUpdate)
	assert.Equal(t, knownString, *castedUpdate.Payload.Data.Attributes.Rules[0].MetricVendor)

	// Try the update without bad revision
	badRev := "whoKnows"
	updateRequest.Data.Attributes.Rev = &badRev
	updatedConflict := handlers.HandleUpdateDataCleaningProfileV2(handlers.SkylightAndTenantAdminRoles, tenantDB)(tenant_provisioning_service_v2.UpdateDataCleaningProfileParams{
		ProfileID:   *updateRequest.Data.ID,
		Body:        &updateRequest,
		HTTPRequest: createHttpRequest(createdTenant.ID, handlers.UserRoleSkylight)})
	castedConflict := updatedConflict.(*tenant_provisioning_service_v2.UpdateDataCleaningProfileConflict)
	assert.NotNil(t, castedConflict)

	// Test updating a record for a tenantthat does not exist
	updateDNE := handlers.HandleUpdateDataCleaningProfileV2(handlers.SkylightAndTenantAdminRoles, tenantDB)(tenant_provisioning_service_v2.UpdateDataCleaningProfileParams{
		Body:        &updateRequest,
		HTTPRequest: createHttpRequest("I am  not real", handlers.UserRoleSkylight)})
	castedUpdateDNE := updateDNE.(*tenant_provisioning_service_v2.UpdateDataCleaningProfileNotFound)
	assert.NotNil(t, castedUpdateDNE)

	// Make sure Tenant User is not able to make this call
	updateFrbdn := handlers.HandleUpdateDataCleaningProfileV2(handlers.SkylightAndTenantAdminRoles, tenantDB)(tenant_provisioning_service_v2.UpdateDataCleaningProfileParams{
		Body:        &updateRequest,
		HTTPRequest: createHttpRequest(createdTenant.ID, handlers.UserRoleTenantUser)})
	castedUpdateFrbdn := updateFrbdn.(*tenant_provisioning_service_v2.UpdateDataCleaningProfileForbidden)
	assert.NotNil(t, castedUpdateFrbdn)

	// Try the update with bad data:
	updateBad := handlers.HandleUpdateDataCleaningProfileV2(handlers.SkylightAndTenantAdminRoles, tenantDB)(tenant_provisioning_service_v2.UpdateDataCleaningProfileParams{
		ProfileID:   *updateRequest.Data.ID,
		Body:        nil,
		HTTPRequest: createHttpRequest(createdTenant.ID, handlers.UserRoleTenantAdmin)})
	castedBad := updateBad.(*tenant_provisioning_service_v2.UpdateDataCleaningProfileBadRequest)
	assert.NotNil(t, castedBad)

	adminDB.DeleteTenant(createdTenant.ID)
}
