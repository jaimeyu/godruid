package handlers_test

import (
	"encoding/json"
	"testing"

	"github.com/accedian/adh-gather/swagmodels"

	"github.com/accedian/adh-gather/handlers"
	"github.com/accedian/adh-gather/restapi/operations/tenant_provisioning_service_v2"
	"github.com/stretchr/testify/assert"
)

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
	castedResponse := existing.(*tenant_provisioning_service_v2.GetDataCleaningProfilesOK)
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
	castedResponse := existing.(*tenant_provisioning_service_v2.GetDataCleaningProfilesOK)
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

func TestCreateDataCleaningProfile(t *testing.T) {

	tenantDescriptor1 := getRandomTenantDescriptor()

	createdTenant, err := adminDB.CreateTenant(tenantDescriptor1)
	assert.Nil(t, err)
	assert.NotNil(t, createdTenant)
	assert.Equal(t, tenantDescriptor1.Name, createdTenant.Name)
	assert.Equal(t, tenantDescriptor1.URLSubdomain, createdTenant.URLSubdomain)
	assert.True(t, len(createdTenant.ID) > 0)

	// Make sure there are no existing Data Cleaning Profiles
	existing := handlers.HandleGetDataCleaningProfilesV2(handlers.SkylightAndTenantAdminRoles, tenantDB)(tenant_provisioning_service_v2.GetDataCleaningProfilesParams{HTTPRequest: createHttpRequest(createdTenant.ID, handlers.UserRoleSkylight)})
	castedResponse := existing.(*tenant_provisioning_service_v2.GetDataCleaningProfilesOK)
	assert.NotNil(t, castedResponse)

	// Make sure Tenant User is not able to make this call
	createdFrbdn := handlers.HandleCreateDataCleaningProfileV2(handlers.SkylightAndTenantAdminRoles, tenantDB)(tenant_provisioning_service_v2.CreateDataCleaningProfileParams{
		Body: &swagmodels.DataCleaningProfileCreateRequest{
			Data: createRandomDataCleaningProfileCreateRequest(),
		},
		HTTPRequest: createHttpRequest(createdTenant.ID, handlers.UserRoleTenantUser)})
	castedCreateFrbdn := createdFrbdn.(*tenant_provisioning_service_v2.CreateDataCleaningProfileForbidden)
	assert.NotNil(t, castedCreateFrbdn)

	// Try the create with bad data:
	createdBad := handlers.HandleCreateDataCleaningProfileV2(handlers.SkylightAndTenantAdminRoles, tenantDB)(tenant_provisioning_service_v2.CreateDataCleaningProfileParams{
		Body: &swagmodels.DataCleaningProfileCreateRequest{
			Data: nil,
		},
		HTTPRequest: createHttpRequest(createdTenant.ID, handlers.UserRoleSkylight)})
	castedCreateBad := createdBad.(*tenant_provisioning_service_v2.CreateDataCleaningProfileBadRequest)
	assert.NotNil(t, castedCreateBad)

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

	// Create another - faile with conflict
	createdConflict := handlers.HandleCreateDataCleaningProfileV2(handlers.SkylightAndTenantAdminRoles, tenantDB)(tenant_provisioning_service_v2.CreateDataCleaningProfileParams{
		Body: &swagmodels.DataCleaningProfileCreateRequest{
			Data: createRandomDataCleaningProfileCreateRequest(),
		},
		HTTPRequest: createHttpRequest(createdTenant.ID, handlers.UserRoleSkylight)})
	castedCreateConflict := createdConflict.(*tenant_provisioning_service_v2.CreateDataCleaningProfileConflict)
	assert.NotNil(t, castedCreateConflict)

	adminDB.DeleteTenant(createdTenant.ID)
}

func TestDeleteDataCleaningProfile(t *testing.T) {

	tenantDescriptor1 := getRandomTenantDescriptor()

	createdTenant, err := adminDB.CreateTenant(tenantDescriptor1)
	assert.Nil(t, err)
	assert.NotNil(t, createdTenant)
	assert.Equal(t, tenantDescriptor1.Name, createdTenant.Name)
	assert.Equal(t, tenantDescriptor1.URLSubdomain, createdTenant.URLSubdomain)
	assert.True(t, len(createdTenant.ID) > 0)

	// Make sure there are no existing Data Cleaning Profiles
	existing := handlers.HandleGetDataCleaningProfilesV2(handlers.SkylightAndTenantAdminRoles, tenantDB)(tenant_provisioning_service_v2.GetDataCleaningProfilesParams{HTTPRequest: createHttpRequest(createdTenant.ID, handlers.UserRoleSkylight)})
	castedResponse := existing.(*tenant_provisioning_service_v2.GetDataCleaningProfilesOK)
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

	// Test deleting a record that does not exist
	badID := "notgood nope"
	deleteDNE := handlers.HandleDeleteDataCleaningProfileV2(handlers.SkylightAndTenantAdminRoles, tenantDB)(tenant_provisioning_service_v2.DeleteDataCleaningProfileParams{
		ProfileID:   badID,
		HTTPRequest: createHttpRequest(createdTenant.ID, handlers.UserRoleSkylight)})
	castedDeleteDNE := deleteDNE.(*tenant_provisioning_service_v2.DeleteDataCleaningProfileNotFound)
	assert.NotNil(t, castedDeleteDNE)

	// Make sure Tenant User is not able to make this call
	deleteFrbdn := handlers.HandleDeleteDataCleaningProfileV2(handlers.SkylightAndTenantAdminRoles, tenantDB)(tenant_provisioning_service_v2.DeleteDataCleaningProfileParams{
		ProfileID:   *castedCreate.Payload.Data.ID,
		HTTPRequest: createHttpRequest(createdTenant.ID, handlers.UserRoleTenantUser)})
	castedDeleteFrbdn := deleteFrbdn.(*tenant_provisioning_service_v2.DeleteDataCleaningProfileForbidden)
	assert.NotNil(t, castedDeleteFrbdn)

	// Delete successfully
	deleted := handlers.HandleDeleteDataCleaningProfileV2(handlers.SkylightAndTenantAdminRoles, tenantDB)(tenant_provisioning_service_v2.DeleteDataCleaningProfileParams{
		ProfileID:   *castedCreate.Payload.Data.ID,
		HTTPRequest: createHttpRequest(createdTenant.ID, handlers.UserRoleSkylight)})
	castedDelete := deleted.(*tenant_provisioning_service_v2.DeleteDataCleaningProfileOK)
	assert.NotNil(t, castedDelete)
	assert.Equal(t, castedFetch.Payload.Data, castedDelete.Payload.Data)

	// Make sure the record no longer exists
	existing = handlers.HandleGetDataCleaningProfilesV2(handlers.SkylightAndTenantAdminRoles, tenantDB)(tenant_provisioning_service_v2.GetDataCleaningProfilesParams{HTTPRequest: createHttpRequest(createdTenant.ID, handlers.UserRoleSkylight)})
	castedResponse = existing.(*tenant_provisioning_service_v2.GetDataCleaningProfilesOK)
	assert.NotNil(t, castedResponse)

	adminDB.DeleteTenant(createdTenant.ID)
}

func TestGetDataCleaningProfiles(t *testing.T) {

	tenantDescriptor1 := getRandomTenantDescriptor()

	createdTenant, err := adminDB.CreateTenant(tenantDescriptor1)
	assert.Nil(t, err)
	assert.NotNil(t, createdTenant)
	assert.Equal(t, tenantDescriptor1.Name, createdTenant.Name)
	assert.Equal(t, tenantDescriptor1.URLSubdomain, createdTenant.URLSubdomain)
	assert.True(t, len(createdTenant.ID) > 0)

	// Make sure there are no existing Data Cleaning Profiles
	existing := handlers.HandleGetDataCleaningProfilesV2(handlers.SkylightAndTenantAdminRoles, tenantDB)(tenant_provisioning_service_v2.GetDataCleaningProfilesParams{HTTPRequest: createHttpRequest(createdTenant.ID, handlers.UserRoleSkylight)})
	castedResponse := existing.(*tenant_provisioning_service_v2.GetDataCleaningProfilesOK)
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
	fetched := handlers.HandleGetDataCleaningProfilesV2(handlers.SkylightAndTenantAdminRoles, tenantDB)(tenant_provisioning_service_v2.GetDataCleaningProfilesParams{
		HTTPRequest: createHttpRequest(createdTenant.ID, handlers.UserRoleSkylight)})
	castedFetch := fetched.(*tenant_provisioning_service_v2.GetDataCleaningProfilesOK)
	assert.NotNil(t, castedFetch)
	assert.Equal(t, 1, len(castedFetch.Payload.Data))

	// Make sure TenantAdmin can't amke this call
	fetchedFrbdn := handlers.HandleGetDataCleaningProfilesV2(handlers.SkylightAndTenantAdminRoles, tenantDB)(tenant_provisioning_service_v2.GetDataCleaningProfilesParams{
		HTTPRequest: createHttpRequest(createdTenant.ID, handlers.UserRoleTenantUser)})
	castedFetchFrbdn := fetchedFrbdn.(*tenant_provisioning_service_v2.GetDataCleaningProfilesForbidden)
	assert.NotNil(t, castedFetchFrbdn)

	adminDB.DeleteTenant(createdTenant.ID)
}
