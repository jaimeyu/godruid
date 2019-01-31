package handlers_test

import (
	"encoding/json"
	"testing"

	"github.com/icrowley/fake"
	"github.com/stretchr/testify/assert"

	"github.com/accedian/adh-gather/handlers"
	"github.com/accedian/adh-gather/models/common"
	"github.com/accedian/adh-gather/restapi/operations/admin_provisioning_service_v2"
	"github.com/accedian/adh-gather/restapi/operations/tenant_provisioning_service_v2"
	"github.com/accedian/adh-gather/swagmodels"
)

var (
	ThresholdProfileUrl = "http://deployment.test.cool/api/v2/threshold-profiles"

	ThresholdProfileTypeString = "thresholdProfiles"
)

func TestThresholdProfileCRUDV2(t *testing.T) {

	createdTenant := handlers.HandleCreateTenantV2(handlers.SkylightAdminRoleOnly, adminDB, tenantDB, stubbedMetricBaselineDatastore)(admin_provisioning_service_v2.CreateTenantV2Params{Body: generateRandomTenantCreationRequest(), HTTPRequest: createHttpRequestWithParams("", handlers.UserRoleSkylight, tenantURL, "POST")})
	castedCreateTeant := createdTenant.(*admin_provisioning_service_v2.CreateTenantV2Created)
	assert.NotNil(t, castedCreateTeant)
	assert.NotEmpty(t, castedCreateTeant.Payload.Data.ID)
	assert.NotEmpty(t, castedCreateTeant.Payload.Data.Attributes.Name)
	assert.NotEmpty(t, castedCreateTeant.Payload.Data.Attributes.URLSubdomain)
	assert.Equal(t, string(common.UserActive), *castedCreateTeant.Payload.Data.Attributes.State)
	assert.True(t, *castedCreateTeant.Payload.Data.Attributes.CreatedTimestamp > 0)
	assert.True(t, *castedCreateTeant.Payload.Data.Attributes.LastModifiedTimestamp > 0)

	// Make sure a record is created when the tenant is created
	existing := handlers.HandleGetAllThresholdProfilesV2(handlers.AllRoles, tenantDB)(tenant_provisioning_service_v2.GetAllThresholdProfilesV2Params{HTTPRequest: createHttpRequestWithParams(*castedCreateTeant.Payload.Data.ID, handlers.UserRoleSkylight, ThresholdProfileUrl, "GET")})
	castedResponse := existing.(*tenant_provisioning_service_v2.GetAllThresholdProfilesV2OK)
	assert.NotNil(t, castedResponse)
	assert.Equal(t, 1, len(castedResponse.Payload.Data))
	assert.Equal(t, "Default", *castedResponse.Payload.Data[0].Attributes.Name)

	created := handlers.HandleCreateThresholdProfileV2(handlers.AllRoles, tenantDB)(tenant_provisioning_service_v2.CreateThresholdProfileV2Params{Body: generateRandomTenantThresholdProfileCreationRequest(), HTTPRequest: createHttpRequestWithParams(*castedCreateTeant.Payload.Data.ID, handlers.UserRoleSkylight, ThresholdProfileUrl, "POST")})
	castedCreate := created.(*tenant_provisioning_service_v2.CreateThresholdProfileV2Created)
	assert.NotNil(t, castedCreate)
	assert.NotEmpty(t, castedCreate.Payload.Data.ID)
	assert.NotEmpty(t, castedCreate.Payload.Data.Attributes.Name)
	assert.NotEmpty(t, castedCreate.Payload.Data.Attributes.Thresholds)
	assert.NotEmpty(t, castedCreate.Payload.Data.Attributes.Datatype)
	assert.True(t, *castedCreate.Payload.Data.Attributes.CreatedTimestamp > 0)
	assert.True(t, *castedCreate.Payload.Data.Attributes.LastModifiedTimestamp > 0)

	// Make sure we can retrieve this record:
	fetched := handlers.HandleGetThresholdProfileV2(handlers.AllRoles, tenantDB)(tenant_provisioning_service_v2.GetThresholdProfileV2Params{ThrPrfID: *castedCreate.Payload.Data.ID, HTTPRequest: createHttpRequestWithParams(*castedCreateTeant.Payload.Data.ID, handlers.UserRoleSkylight, ThresholdProfileUrl, "GET")})
	castedFetch := fetched.(*tenant_provisioning_service_v2.GetThresholdProfileV2OK)
	assert.NotNil(t, castedFetch)
	assert.Equal(t, castedCreate.Payload.Data.Attributes.Thresholds, castedFetch.Payload.Data.Attributes.Thresholds)
	assert.NotEmpty(t, castedFetch.Payload.Data.Attributes.ThresholdList)

	// Make sure there are now multiple records
	fetchList := handlers.HandleGetAllThresholdProfilesV2(handlers.AllRoles, tenantDB)(tenant_provisioning_service_v2.GetAllThresholdProfilesV2Params{HTTPRequest: createHttpRequestWithParams(*castedCreateTeant.Payload.Data.ID, handlers.UserRoleSkylight, ThresholdProfileUrl, "GET")})
	castedFetchList := fetchList.(*tenant_provisioning_service_v2.GetAllThresholdProfilesV2OK)
	assert.NotNil(t, castedFetchList)
	assert.Equal(t, 2, len(castedFetchList.Payload.Data))

	// Make an update to the Record
	newName := fake.CharactersN(16)
	updateRequestBody := generateThresholdProfileUpdateRequest(*castedCreate.Payload.Data.ID, *castedCreate.Payload.Data.Attributes.Rev, &newName)
	updated := handlers.HandleUpdateThresholdProfileV2(handlers.AllRoles, tenantDB)(tenant_provisioning_service_v2.UpdateThresholdProfileV2Params{ThrPrfID: *castedCreate.Payload.Data.ID, Body: updateRequestBody, HTTPRequest: createHttpRequestWithParams(*castedCreateTeant.Payload.Data.ID, handlers.UserRoleSkylight, ThresholdProfileUrl, "PATCH")})
	castedUpdate := updated.(*tenant_provisioning_service_v2.UpdateThresholdProfileV2OK)
	assert.NotNil(t, castedUpdate)
	assert.NotEqual(t, castedCreate.Payload.Data, castedUpdate.Payload.Data)
	assert.NotEqual(t, castedCreate.Payload.Data.Attributes.Rev, castedUpdate.Payload.Data.Attributes.Rev)
	assert.Equal(t, newName, *castedUpdate.Payload.Data.Attributes.Name)
	assert.NotEmpty(t, castedUpdate.Payload.Data.Attributes.Thresholds)

	// Delete the record
	deleted := handlers.HandleDeleteThresholdProfileV2(handlers.AllRoles, tenantDB)(tenant_provisioning_service_v2.DeleteThresholdProfileV2Params{ThrPrfID: *castedCreate.Payload.Data.ID, HTTPRequest: createHttpRequestWithParams(*castedCreateTeant.Payload.Data.ID, handlers.UserRoleSkylight, ThresholdProfileUrl, "DELETE")})
	castedDelete := deleted.(*tenant_provisioning_service_v2.DeleteThresholdProfileV2OK)
	assert.NotNil(t, castedDelete)
	assert.Equal(t, castedUpdate.Payload.Data, castedDelete.Payload.Data)

	// Make sure there is only 1 record left
	existing = handlers.HandleGetAllThresholdProfilesV2(handlers.AllRoles, tenantDB)(tenant_provisioning_service_v2.GetAllThresholdProfilesV2Params{HTTPRequest: createHttpRequestWithParams(*castedCreateTeant.Payload.Data.ID, handlers.UserRoleSkylight, ThresholdProfileUrl, "GET")})
	castedResponse = existing.(*tenant_provisioning_service_v2.GetAllThresholdProfilesV2OK)
	assert.NotNil(t, castedResponse)
	assert.Equal(t, 1, len(castedResponse.Payload.Data))
}

func TestThresholdProfileNotFoundV2(t *testing.T) {
	notFoundID := fake.CharactersN(20)

	// Get All ThresholdProfile
	fetchedAll := handlers.HandleGetAllThresholdProfilesV2(handlers.AllRoles, tenantDB)(tenant_provisioning_service_v2.GetAllThresholdProfilesV2Params{HTTPRequest: createHttpRequestWithParams("", handlers.UserRoleSkylight, ThresholdProfileUrl, "GET")})
	castedFetchAll := fetchedAll.(*tenant_provisioning_service_v2.GetAllThresholdProfilesV2NotFound)
	assert.NotNil(t, castedFetchAll)

	// Get ThresholdProfile
	fetched := handlers.HandleGetThresholdProfileV2(handlers.AllRoles, tenantDB)(tenant_provisioning_service_v2.GetThresholdProfileV2Params{ThrPrfID: notFoundID, HTTPRequest: createHttpRequestWithParams("", handlers.UserRoleSkylight, ThresholdProfileUrl, "GET")})
	castedFetch := fetched.(*tenant_provisioning_service_v2.GetThresholdProfileV2NotFound)
	assert.NotNil(t, castedFetch)

	// Delete ThresholdProfile
	deleted := handlers.HandleDeleteThresholdProfileV2(handlers.AllRoles, tenantDB)(tenant_provisioning_service_v2.DeleteThresholdProfileV2Params{ThrPrfID: notFoundID, HTTPRequest: createHttpRequestWithParams("", handlers.UserRoleSkylight, ThresholdProfileUrl, "DELETE")})
	castedDelete := deleted.(*tenant_provisioning_service_v2.DeleteThresholdProfileV2NotFound)
	assert.NotNil(t, castedDelete)

	// Patch ThresholdProfile
	updateRequest := generateThresholdProfileUpdateRequest(notFoundID, "reviosionstuff", nil)
	updated := handlers.HandleUpdateThresholdProfileV2(handlers.AllRoles, tenantDB)(tenant_provisioning_service_v2.UpdateThresholdProfileV2Params{ThrPrfID: notFoundID, Body: updateRequest, HTTPRequest: createHttpRequestWithParams("", handlers.UserRoleSkylight, ThresholdProfileUrl, "PATCH")})
	castedUpdate := updated.(*tenant_provisioning_service_v2.UpdateThresholdProfileV2NotFound)
	assert.NotNil(t, castedUpdate)
}

func TestThresholdProfileBadRequestV2(t *testing.T) {

	// CreateThresholdProfile
	created := handlers.HandleCreateThresholdProfileV2(handlers.AllRoles, tenantDB)(tenant_provisioning_service_v2.CreateThresholdProfileV2Params{Body: nil, HTTPRequest: createHttpRequestWithParams("", handlers.UserRoleSkylight, ThresholdProfileUrl, "POST")})
	castedCreate := created.(*tenant_provisioning_service_v2.CreateThresholdProfileV2BadRequest)
	assert.NotNil(t, castedCreate)

	// Update ThresholdProfile
	updated := handlers.HandleUpdateThresholdProfileV2(handlers.AllRoles, tenantDB)(tenant_provisioning_service_v2.UpdateThresholdProfileV2Params{ThrPrfID: fake.CharactersN(20), Body: nil, HTTPRequest: createHttpRequestWithParams("", handlers.UserRoleSkylight, ThresholdProfileUrl, "PATCH")})
	castedUpdate := updated.(*tenant_provisioning_service_v2.UpdateThresholdProfileV2BadRequest)
	assert.NotNil(t, castedUpdate)
}

func TestThresholdProfileConflictV2(t *testing.T) {

	createdTenant := handlers.HandleCreateTenantV2(handlers.SkylightAdminRoleOnly, adminDB, tenantDB, stubbedMetricBaselineDatastore)(admin_provisioning_service_v2.CreateTenantV2Params{Body: generateRandomTenantCreationRequest(), HTTPRequest: createHttpRequestWithParams("", handlers.UserRoleSkylight, tenantURL, "POST")})
	castedCreateTeant := createdTenant.(*admin_provisioning_service_v2.CreateTenantV2Created)
	assert.NotNil(t, castedCreateTeant)
	assert.NotEmpty(t, castedCreateTeant.Payload.Data.ID)
	assert.NotEmpty(t, castedCreateTeant.Payload.Data.Attributes.Name)
	assert.NotEmpty(t, castedCreateTeant.Payload.Data.Attributes.URLSubdomain)
	assert.Equal(t, string(common.UserActive), *castedCreateTeant.Payload.Data.Attributes.State)
	assert.True(t, *castedCreateTeant.Payload.Data.Attributes.CreatedTimestamp > 0)
	assert.True(t, *castedCreateTeant.Payload.Data.Attributes.LastModifiedTimestamp > 0)

	// Make sure a record is created when the tenant is created
	existing := handlers.HandleGetAllThresholdProfilesV2(handlers.AllRoles, tenantDB)(tenant_provisioning_service_v2.GetAllThresholdProfilesV2Params{HTTPRequest: createHttpRequestWithParams(*castedCreateTeant.Payload.Data.ID, handlers.UserRoleSkylight, ThresholdProfileUrl, "GET")})
	castedResponse := existing.(*tenant_provisioning_service_v2.GetAllThresholdProfilesV2OK)
	assert.NotNil(t, castedResponse)
	assert.Equal(t, 1, len(castedResponse.Payload.Data))

	name := fake.CharactersN(15)
	createReqBody := generateNamedTenantThresholdProfileCreationRequest(name)
	created := handlers.HandleCreateThresholdProfileV2(handlers.AllRoles, tenantDB)(tenant_provisioning_service_v2.CreateThresholdProfileV2Params{Body: createReqBody, HTTPRequest: createHttpRequestWithParams(*castedCreateTeant.Payload.Data.ID, handlers.UserRoleSkylight, ThresholdProfileUrl, "POST")})
	castedCreate := created.(*tenant_provisioning_service_v2.CreateThresholdProfileV2Created)
	assert.NotNil(t, castedCreate)
	assert.NotEmpty(t, castedCreate.Payload.Data.ID)
	assert.NotEmpty(t, castedCreate.Payload.Data.Attributes.Name)
	assert.NotEmpty(t, castedCreate.Payload.Data.Attributes.Thresholds)
	assert.NotEmpty(t, castedCreate.Payload.Data.Attributes.Datatype)
	assert.True(t, *castedCreate.Payload.Data.Attributes.CreatedTimestamp > 0)
	assert.True(t, *castedCreate.Payload.Data.Attributes.LastModifiedTimestamp > 0)

	// Try to create the record again - should succeed as we are not guarding against name collisiones
	createdConflict := handlers.HandleCreateThresholdProfileV2(handlers.AllRoles, tenantDB)(tenant_provisioning_service_v2.CreateThresholdProfileV2Params{Body: createReqBody, HTTPRequest: createHttpRequestWithParams(*castedCreateTeant.Payload.Data.ID, handlers.UserRoleSkylight, ThresholdProfileUrl, "POST")})
	castedCreateConflictButOK := createdConflict.(*tenant_provisioning_service_v2.CreateThresholdProfileV2Created)
	assert.NotNil(t, castedCreateConflictButOK)

	// Try the update with a bad revision
	newName := fake.CharactersN(16)
	updateRequestBody := generateThresholdProfileUpdateRequest(*castedCreate.Payload.Data.ID, *castedCreate.Payload.Data.Attributes.Rev+"pork", &newName)
	updated := handlers.HandleUpdateThresholdProfileV2(handlers.AllRoles, tenantDB)(tenant_provisioning_service_v2.UpdateThresholdProfileV2Params{ThrPrfID: *castedCreate.Payload.Data.ID, Body: updateRequestBody, HTTPRequest: createHttpRequestWithParams(*castedCreateTeant.Payload.Data.ID, handlers.UserRoleSkylight, ThresholdProfileUrl, "PATCH")})
	castedUpdate := updated.(*tenant_provisioning_service_v2.UpdateThresholdProfileV2Conflict)
	assert.NotNil(t, castedUpdate)

	// Delete the profile
	deleted := handlers.HandleDeleteThresholdProfileV2(handlers.AllRoles, tenantDB)(tenant_provisioning_service_v2.DeleteThresholdProfileV2Params{ThrPrfID: *castedCreate.Payload.Data.ID, HTTPRequest: createHttpRequestWithParams(*castedCreateTeant.Payload.Data.ID, handlers.UserRoleSkylight, ThresholdProfileUrl, "DELETE")})
	castedDelete := deleted.(*tenant_provisioning_service_v2.DeleteThresholdProfileV2OK)
	assert.NotNil(t, castedDelete)
	assert.NotNil(t, castedDelete.Payload.Data)
}

func TestThresholdProfileAPIsProtectedByAuthV2(t *testing.T) {
	fakeTenantID := fake.CharactersN(20)
	// Get All - All Users
	existing := handlers.HandleGetAllThresholdProfilesV2(handlers.AllRoles, tenantDB)(tenant_provisioning_service_v2.GetAllThresholdProfilesV2Params{HTTPRequest: createHttpRequestWithParams(fakeTenantID, handlers.UserRoleUnknown, ThresholdProfileUrl, "GET")})
	castedResponse := existing.(*tenant_provisioning_service_v2.GetAllThresholdProfilesV2Forbidden)
	assert.NotNil(t, castedResponse)

	fakeID := fake.CharactersN(20)

	// Get - All Users
	fetched := handlers.HandleGetThresholdProfileV2(handlers.AllRoles, tenantDB)(tenant_provisioning_service_v2.GetThresholdProfileV2Params{ThrPrfID: fakeID, HTTPRequest: createHttpRequestWithParams(fakeTenantID, handlers.UserRoleUnknown, ThresholdProfileUrl, "GET")})
	castedFetch := fetched.(*tenant_provisioning_service_v2.GetThresholdProfileV2Forbidden)
	assert.NotNil(t, castedFetch)

	// Create - SkylightAdmin and TenantAdmin Only
	created := handlers.HandleCreateThresholdProfileV2(handlers.SkylightAndTenantAdminRoles, tenantDB)(tenant_provisioning_service_v2.CreateThresholdProfileV2Params{Body: generateRandomTenantThresholdProfileCreationRequest(), HTTPRequest: createHttpRequestWithParams(fakeTenantID, handlers.UserRoleUnknown, ThresholdProfileUrl, "POST")})
	castedCreate := created.(*tenant_provisioning_service_v2.CreateThresholdProfileV2Forbidden)
	assert.NotNil(t, castedCreate)

	created = handlers.HandleCreateThresholdProfileV2(handlers.SkylightAndTenantAdminRoles, tenantDB)(tenant_provisioning_service_v2.CreateThresholdProfileV2Params{Body: generateRandomTenantThresholdProfileCreationRequest(), HTTPRequest: createHttpRequestWithParams(fakeTenantID, handlers.UserRoleTenantUser, ThresholdProfileUrl, "POST")})
	castedCreate = created.(*tenant_provisioning_service_v2.CreateThresholdProfileV2Forbidden)
	assert.NotNil(t, castedCreate)

	created = handlers.HandleCreateThresholdProfileV2(handlers.SkylightAndTenantAdminRoles, tenantDB)(tenant_provisioning_service_v2.CreateThresholdProfileV2Params{Body: generateRandomTenantThresholdProfileCreationRequest(), HTTPRequest: createHttpRequestWithParams(fakeTenantID, handlers.UserRoleTenantContributor, ThresholdProfileUrl, "POST")})
	castedCreate = created.(*tenant_provisioning_service_v2.CreateThresholdProfileV2Forbidden)
	assert.NotNil(t, castedCreate)

	// Update - SkylightAdmin and TenantAdmin Only
	updated := handlers.HandleUpdateThresholdProfileV2(handlers.SkylightAndTenantAdminRoles, tenantDB)(tenant_provisioning_service_v2.UpdateThresholdProfileV2Params{ThrPrfID: fakeID, Body: nil, HTTPRequest: createHttpRequestWithParams(fakeTenantID, handlers.UserRoleUnknown, ThresholdProfileUrl, "PATCH")})
	castedUpdate := updated.(*tenant_provisioning_service_v2.UpdateThresholdProfileV2Forbidden)
	assert.NotNil(t, castedUpdate)

	updated = handlers.HandleUpdateThresholdProfileV2(handlers.SkylightAndTenantAdminRoles, tenantDB)(tenant_provisioning_service_v2.UpdateThresholdProfileV2Params{ThrPrfID: fakeID, Body: nil, HTTPRequest: createHttpRequestWithParams(fakeTenantID, handlers.UserRoleTenantUser, ThresholdProfileUrl, "PATCH")})
	castedUpdate = updated.(*tenant_provisioning_service_v2.UpdateThresholdProfileV2Forbidden)
	assert.NotNil(t, castedUpdate)

	updated = handlers.HandleUpdateThresholdProfileV2(handlers.SkylightAndTenantAdminRoles, tenantDB)(tenant_provisioning_service_v2.UpdateThresholdProfileV2Params{ThrPrfID: fakeID, Body: nil, HTTPRequest: createHttpRequestWithParams(fakeTenantID, handlers.UserRoleTenantContributor, ThresholdProfileUrl, "PATCH")})
	castedUpdate = updated.(*tenant_provisioning_service_v2.UpdateThresholdProfileV2Forbidden)
	assert.NotNil(t, castedUpdate)

	// Delete - SkylightAdmin and TenantAdmin Only
	deleted := handlers.HandleDeleteThresholdProfileV2(handlers.SkylightAndTenantAdminRoles, tenantDB)(tenant_provisioning_service_v2.DeleteThresholdProfileV2Params{ThrPrfID: fakeID, HTTPRequest: createHttpRequestWithParams(fakeTenantID, handlers.UserRoleUnknown, ThresholdProfileUrl, "DELETE")})
	castedDelete := deleted.(*tenant_provisioning_service_v2.DeleteThresholdProfileV2Forbidden)
	assert.NotNil(t, castedDelete)

	deleted = handlers.HandleDeleteThresholdProfileV2(handlers.SkylightAndTenantAdminRoles, tenantDB)(tenant_provisioning_service_v2.DeleteThresholdProfileV2Params{ThrPrfID: fakeID, HTTPRequest: createHttpRequestWithParams(fakeTenantID, handlers.UserRoleTenantUser, ThresholdProfileUrl, "DELETE")})
	castedDelete = deleted.(*tenant_provisioning_service_v2.DeleteThresholdProfileV2Forbidden)
	assert.NotNil(t, castedDelete)

	deleted = handlers.HandleDeleteThresholdProfileV2(handlers.SkylightAndTenantAdminRoles, tenantDB)(tenant_provisioning_service_v2.DeleteThresholdProfileV2Params{ThrPrfID: fakeID, HTTPRequest: createHttpRequestWithParams(fakeTenantID, handlers.UserRoleTenantContributor, ThresholdProfileUrl, "DELETE")})
	castedDelete = deleted.(*tenant_provisioning_service_v2.DeleteThresholdProfileV2Forbidden)
	assert.NotNil(t, castedDelete)
}

func TestThresholdProfileDeleteIntegrityCheckV2(t *testing.T) {

	createdTenant := handlers.HandleCreateTenantV2(handlers.SkylightAdminRoleOnly, adminDB, tenantDB, stubbedMetricBaselineDatastore)(admin_provisioning_service_v2.CreateTenantV2Params{Body: generateRandomTenantCreationRequest(), HTTPRequest: createHttpRequestWithParams("", handlers.UserRoleSkylight, tenantURL, "POST")})
	castedCreateTeant := createdTenant.(*admin_provisioning_service_v2.CreateTenantV2Created)
	assert.NotNil(t, castedCreateTeant)
	assert.NotEmpty(t, castedCreateTeant.Payload.Data.ID)
	assert.NotEmpty(t, castedCreateTeant.Payload.Data.Attributes.Name)
	assert.NotEmpty(t, castedCreateTeant.Payload.Data.Attributes.URLSubdomain)
	assert.Equal(t, string(common.UserActive), *castedCreateTeant.Payload.Data.Attributes.State)
	assert.True(t, *castedCreateTeant.Payload.Data.Attributes.CreatedTimestamp > 0)
	assert.True(t, *castedCreateTeant.Payload.Data.Attributes.LastModifiedTimestamp > 0)

	// Make sure a record is created when the tenant is created
	existing := handlers.HandleGetAllThresholdProfilesV2(handlers.AllRoles, tenantDB)(tenant_provisioning_service_v2.GetAllThresholdProfilesV2Params{HTTPRequest: createHttpRequestWithParams(*castedCreateTeant.Payload.Data.ID, handlers.UserRoleSkylight, ThresholdProfileUrl, "GET")})
	castedResponse := existing.(*tenant_provisioning_service_v2.GetAllThresholdProfilesV2OK)
	assert.NotNil(t, castedResponse)
	assert.Equal(t, 1, len(castedResponse.Payload.Data))
	assert.Equal(t, "Default", *castedResponse.Payload.Data[0].Attributes.Name)

	created := handlers.HandleCreateThresholdProfileV2(handlers.AllRoles, tenantDB)(tenant_provisioning_service_v2.CreateThresholdProfileV2Params{Body: generateRandomTenantThresholdProfileCreationRequest(), HTTPRequest: createHttpRequestWithParams(*castedCreateTeant.Payload.Data.ID, handlers.UserRoleSkylight, ThresholdProfileUrl, "POST")})
	castedCreate := created.(*tenant_provisioning_service_v2.CreateThresholdProfileV2Created)
	assert.NotNil(t, castedCreate)
	assert.NotEmpty(t, castedCreate.Payload.Data.ID)
	assert.NotEmpty(t, castedCreate.Payload.Data.Attributes.Name)
	assert.NotEmpty(t, castedCreate.Payload.Data.Attributes.Thresholds)
	assert.NotEmpty(t, castedCreate.Payload.Data.Attributes.Datatype)
	assert.True(t, *castedCreate.Payload.Data.Attributes.CreatedTimestamp > 0)
	assert.True(t, *castedCreate.Payload.Data.Attributes.LastModifiedTimestamp > 0)

	// Create a Dashboard that uses the Profile
	dashReq := generateRandomTenantDashboardCreationRequest()
	dashReq.Data.Relationships.ThresholdProfile.Data.ID = *castedCreate.Payload.Data.ID
	createdDash := handlers.HandleCreateDashboardV2(handlers.AllRoles, tenantDB)(tenant_provisioning_service_v2.CreateDashboardV2Params{Body: dashReq, HTTPRequest: createHttpRequestWithParams(*castedCreateTeant.Payload.Data.ID, handlers.UserRoleSkylight, DashboardUrl, "POST")})
	castedCreateDash := createdDash.(*tenant_provisioning_service_v2.CreateDashboardV2Created)
	assert.NotNil(t, castedCreateDash)
	assert.NotEmpty(t, castedCreateDash.Payload.Data.ID)
	assert.NotEmpty(t, castedCreateDash.Payload.Data.Attributes.Name)
	assert.NotEmpty(t, castedCreateDash.Payload.Data.Relationships.ThresholdProfile)
	assert.NotEmpty(t, castedCreateDash.Payload.Data.Attributes.Category)
	assert.NotEmpty(t, castedCreateDash.Payload.Data.Relationships.Cards.Data)
	assert.NotEmpty(t, castedCreateDash.Payload.Data.Attributes.Datatype)
	assert.True(t, *castedCreateDash.Payload.Data.Attributes.CreatedTimestamp > 0)
	assert.True(t, *castedCreateDash.Payload.Data.Attributes.LastModifiedTimestamp > 0)

	// Create a SLA Report Configuration that uses the profile
	createdSLA := handlers.HandleCreateReportScheduleConfigV2(handlers.AllRoles, tenantDB)(tenant_provisioning_service_v2.CreateReportScheduleConfigV2Params{Body: generateRandomTenantReportScheduleConfigCreationRequest(*castedCreate.Payload.Data.ID), HTTPRequest: createHttpRequestWithParams(*castedCreateTeant.Payload.Data.ID, handlers.UserRoleSkylight, ReportScheduleConfigUrl, "POST")})
	castedCreateSLA := createdSLA.(*tenant_provisioning_service_v2.CreateReportScheduleConfigV2Created)
	assert.NotNil(t, castedCreateSLA)
	assert.NotEmpty(t, castedCreateSLA.Payload.Data.ID)
	assert.NotEmpty(t, castedCreateSLA.Payload.Data.Attributes.Name)
	assert.NotEmpty(t, castedCreateSLA.Payload.Data.Relationships)
	assert.NotEmpty(t, castedCreateSLA.Payload.Data.Attributes.Hour)
	assert.NotEmpty(t, castedCreateSLA.Payload.Data.Attributes.Minute)
	assert.NotEmpty(t, castedCreateSLA.Payload.Data.Attributes.DayMonth)
	assert.NotEmpty(t, castedCreateSLA.Payload.Data.Attributes.DayWeek)
	assert.NotEmpty(t, castedCreateSLA.Payload.Data.Attributes.Datatype)
	assert.True(t, *castedCreateSLA.Payload.Data.Attributes.CreatedTimestamp > 0)
	assert.True(t, *castedCreateSLA.Payload.Data.Attributes.LastModifiedTimestamp > 0)

	// Delete - fails due to integrity check
	deleted := handlers.HandleDeleteThresholdProfileV2(handlers.AllRoles, tenantDB)(tenant_provisioning_service_v2.DeleteThresholdProfileV2Params{ThrPrfID: *castedCreate.Payload.Data.ID, HTTPRequest: createHttpRequestWithParams(*castedCreateTeant.Payload.Data.ID, handlers.UserRoleSkylight, ThresholdProfileUrl, "DELETE")})
	castedDeleteError := deleted.(*tenant_provisioning_service_v2.DeleteThresholdProfileV2InternalServerError)
	assert.NotNil(t, castedDeleteError)

	// Delete the Dashboard
	deletedDash := handlers.HandleDeleteDashboardV2(handlers.AllRoles, tenantDB)(tenant_provisioning_service_v2.DeleteDashboardV2Params{DashboardID: *castedCreateDash.Payload.Data.ID, HTTPRequest: createHttpRequestWithParams(*castedCreateTeant.Payload.Data.ID, handlers.UserRoleSkylight, DashboardUrl, "DELETE")})
	castedDeleteDash := deletedDash.(*tenant_provisioning_service_v2.DeleteDashboardV2OK)
	assert.NotNil(t, castedDeleteDash)
	assert.Equal(t, castedCreateDash.Payload.Data, castedDeleteDash.Payload.Data)

	// Delete - still fails due to SLA Report Configuration
	deleted = handlers.HandleDeleteThresholdProfileV2(handlers.AllRoles, tenantDB)(tenant_provisioning_service_v2.DeleteThresholdProfileV2Params{ThrPrfID: *castedCreate.Payload.Data.ID, HTTPRequest: createHttpRequestWithParams(*castedCreateTeant.Payload.Data.ID, handlers.UserRoleSkylight, ThresholdProfileUrl, "DELETE")})
	castedDeleteError = deleted.(*tenant_provisioning_service_v2.DeleteThresholdProfileV2InternalServerError)
	assert.NotNil(t, castedDeleteError)

	// Delete the SLA Report Configuration
	deletedSLA := handlers.HandleDeleteReportScheduleConfigV2(handlers.AllRoles, tenantDB)(tenant_provisioning_service_v2.DeleteReportScheduleConfigV2Params{ConfigID: *castedCreateSLA.Payload.Data.ID, HTTPRequest: createHttpRequestWithParams(*castedCreateTeant.Payload.Data.ID, handlers.UserRoleSkylight, ReportScheduleConfigUrl, "DELETE")})
	castedDeleteSLA := deletedSLA.(*tenant_provisioning_service_v2.DeleteReportScheduleConfigV2OK)
	assert.NotNil(t, castedDeleteSLA)
	assert.Equal(t, castedCreateSLA.Payload.Data, castedDeleteSLA.Payload.Data)

	// Delete - success
	deleted = handlers.HandleDeleteThresholdProfileV2(handlers.AllRoles, tenantDB)(tenant_provisioning_service_v2.DeleteThresholdProfileV2Params{ThrPrfID: *castedCreate.Payload.Data.ID, HTTPRequest: createHttpRequestWithParams(*castedCreateTeant.Payload.Data.ID, handlers.UserRoleSkylight, ThresholdProfileUrl, "DELETE")})
	castedDeleteOK := deleted.(*tenant_provisioning_service_v2.DeleteThresholdProfileV2OK)
	assert.NotNil(t, castedDeleteOK)
	assert.Equal(t, castedCreate.Payload.Data, castedDeleteOK.Payload.Data)
}

func generateRandomTenantThresholdProfileCreationRequest() *swagmodels.ThresholdProfileCreateRequest {
	name := fake.CharactersN(12)
	thresholds := generateThresholdsObject()

	return &swagmodels.ThresholdProfileCreateRequest{
		Data: &swagmodels.ThresholdProfileCreateRequestData{
			Type: &ThresholdProfileTypeString,
			Attributes: &swagmodels.ThresholdProfileCreateRequestDataAttributes{
				Name:       &name,
				Thresholds: thresholds,
			},
		},
	}
}

func generateNamedTenantThresholdProfileCreationRequest(name string) *swagmodels.ThresholdProfileCreateRequest {
	thresholds := generateThresholdsObject()

	return &swagmodels.ThresholdProfileCreateRequest{
		Data: &swagmodels.ThresholdProfileCreateRequestData{
			Type: &ThresholdProfileTypeString,
			Attributes: &swagmodels.ThresholdProfileCreateRequestDataAttributes{
				Name:       &name,
				Thresholds: thresholds,
			},
		},
	}
}

func generateThresholdProfileUpdateRequest(id string, rev string, name *string) *swagmodels.ThresholdProfileUpdateRequest {
	result := &swagmodels.ThresholdProfileUpdateRequest{
		Data: &swagmodels.ThresholdProfileUpdateRequestData{
			Type:       &ThresholdProfileTypeString,
			ID:         &id,
			Attributes: &swagmodels.ThresholdProfileUpdateRequestDataAttributes{Rev: &rev},
		},
	}

	if name != nil {
		result.Data.Attributes.Name = *name
	}

	result.Data.Attributes.Thresholds = generateThresholdsObject()

	return result
}

var (
	thresholdBytes = []byte(`{
		"vendorMap": {
			"accedian-flowmeter": {
				"monitoredObjectTypeMap": {
					"flowmeter": {
						"metricMap": {
							"throughputAvg": {
								"directionMap": {
									"0": {
										"eventMap": {
											"critical": {
												"eventAttrMap": {
													"lowerLimit": "25000000",
													"lowerStrict": "true",
													"unit": "bps"
												}
											},
											"major": {
												"eventAttrMap": {
													"lowerLimit": "20000000",
													"lowerStrict": "true",
													"unit": "bps",
													"upperLimit": "25000000",
													"upperStrict": "false"
												}
											},
											"minor": {
												"eventAttrMap": {
													"lowerLimit": "18000000",
													"lowerStrict": "true",
													"unit": "bps",
													"upperLimit": "20000000"
												}
											}
										}
									}
								}
							}
						}
					}
				}
			},
			"accedian-twamp": {
				"monitoredObjectTypeMap": {
					"twamp-pe": {
						"metricMap": {
							"delayP95": {
								"directionMap": {
									"0": {
										"eventMap": {
											"critical": {
												"eventAttrMap": {
													"lowerLimit": "100000",
													"lowerStrict": "true",
													"unit": "ms"
												}
											},
											"major": {
												"eventAttrMap": {
													"lowerLimit": "95000",
													"lowerStrict": "true",
													"unit": "ms",
													"upperLimit": "100000",
													"upperStrict": "false"
												}
											},
											"minor": {
												"eventAttrMap": {
													"lowerLimit": "92500",
													"lowerStrict": "true",
													"unit": "ms",
													"upperLimit": "95000"
												}
											}
										}
									}
								}
							},
							"jitterP95": {
								"directionMap": {
									"0": {
										"eventMap": {
											"critical": {
												"eventAttrMap": {
													"lowerLimit": "30000",
													"lowerStrict": "true",
													"unit": "ms"
												}
											},
											"major": {
												"eventAttrMap": {
													"lowerLimit": "20000",
													"lowerStrict": "true",
													"unit": "ms",
													"upperLimit": "30000",
													"upperStrict": "false"
												}
											},
											"minor": {
												"eventAttrMap": {
													"lowerLimit": "15000",
													"lowerStrict": "true",
													"unit": "ms",
													"upperLimit": "20000"
												}
											}
										}
									}
								}
							},
							"packetsLostPct": {
								"directionMap": {
									"0": {
										"eventMap": {
											"critical": {
												"eventAttrMap": {
													"lowerLimit": "0.8",
													"lowerStrict": "true",
													"unit": "pct"
												}
											},
											"major": {
												"eventAttrMap": {
													"lowerLimit": "0.3",
													"lowerStrict": "true",
													"unit": "pct",
													"upperLimit": "0.8",
													"upperStrict": "false"
												}
											},
											"minor": {
												"eventAttrMap": {
													"lowerLimit": "0.1",
													"lowerStrict": "true",
													"unit": "pct",
													"upperLimit": "0.3"
												}
											}
										}
									}
								}
							}
						}
					}
				}
			}
		}
	}`)
)

func generateThresholdsObject() *swagmodels.ThresholdsObject {
	result := swagmodels.ThresholdsObject{}

	json.Unmarshal(thresholdBytes, &result)

	return &result
}
