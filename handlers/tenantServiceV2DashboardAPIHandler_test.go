package handlers_test

import (
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
	DashboardUrl = "http://deployment.test.cool/api/v2/dashboards"

	DashboardTypeString = "dashboards"
)

func TestDashboardCRUDV2(t *testing.T) {

	createdTenant := handlers.HandleCreateTenantV2(handlers.SkylightAdminRoleOnly, adminDB, tenantDB, stubbedMetricBaselineDatastore)(admin_provisioning_service_v2.CreateTenantV2Params{Body: generateRandomTenantCreationRequest(), HTTPRequest: createHttpRequestWithParams("", handlers.UserRoleSkylight, tenantURL, "POST")})
	castedCreateTeant := createdTenant.(*admin_provisioning_service_v2.CreateTenantV2Created)
	assert.NotNil(t, castedCreateTeant)
	assert.NotEmpty(t, castedCreateTeant.Payload.Data.ID)
	assert.NotEmpty(t, castedCreateTeant.Payload.Data.Attributes.Name)
	assert.NotEmpty(t, castedCreateTeant.Payload.Data.Attributes.URLSubdomain)
	assert.Equal(t, string(common.UserActive), *castedCreateTeant.Payload.Data.Attributes.State)
	assert.True(t, *castedCreateTeant.Payload.Data.Attributes.CreatedTimestamp > 0)
	assert.True(t, *castedCreateTeant.Payload.Data.Attributes.LastModifiedTimestamp > 0)

	// Make sure there are no existing Dashboards
	existing := handlers.HandleGetAllDashboardsV2(handlers.AllRoles, tenantDB)(tenant_provisioning_service_v2.GetAllDashboardsV2Params{HTTPRequest: createHttpRequestWithParams(*castedCreateTeant.Payload.Data.ID, handlers.UserRoleSkylight, DashboardUrl, "GET")})
	castedResponse := existing.(*tenant_provisioning_service_v2.GetAllDashboardsV2NotFound)
	assert.NotNil(t, castedResponse)

	created := handlers.HandleCreateDashboardV2(handlers.AllRoles, tenantDB)(tenant_provisioning_service_v2.CreateDashboardV2Params{Body: generateRandomTenantDashboardCreationRequest(), HTTPRequest: createHttpRequestWithParams(*castedCreateTeant.Payload.Data.ID, handlers.UserRoleSkylight, DashboardUrl, "POST")})
	castedCreate := created.(*tenant_provisioning_service_v2.CreateDashboardV2Created)
	assert.NotNil(t, castedCreate)
	assert.NotEmpty(t, castedCreate.Payload.Data.ID)
	assert.NotEmpty(t, castedCreate.Payload.Data.Attributes.Name)
	assert.NotEmpty(t, castedCreate.Payload.Data.Relationships.ThresholdProfile)
	assert.NotEmpty(t, castedCreate.Payload.Data.Attributes.Category)
	assert.NotEmpty(t, castedCreate.Payload.Data.Relationships.Cards.Data)
	assert.NotEmpty(t, castedCreate.Payload.Data.Attributes.Datatype)
	assert.True(t, *castedCreate.Payload.Data.Attributes.CreatedTimestamp > 0)
	assert.True(t, *castedCreate.Payload.Data.Attributes.LastModifiedTimestamp > 0)

	// Make sure we can retrieve this record:
	fetched := handlers.HandleGetDashboardV2(handlers.AllRoles, tenantDB)(tenant_provisioning_service_v2.GetDashboardV2Params{DashboardID: *castedCreate.Payload.Data.ID, HTTPRequest: createHttpRequestWithParams(*castedCreateTeant.Payload.Data.ID, handlers.UserRoleSkylight, DashboardUrl, "GET")})
	castedFetch := fetched.(*tenant_provisioning_service_v2.GetDashboardV2OK)
	assert.NotNil(t, castedFetch)
	assert.Equal(t, castedCreate.Payload.Data, castedFetch.Payload.Data)

	// Also retrieve the record as part of an array
	fetchList := handlers.HandleGetAllDashboardsV2(handlers.AllRoles, tenantDB)(tenant_provisioning_service_v2.GetAllDashboardsV2Params{HTTPRequest: createHttpRequestWithParams(*castedCreateTeant.Payload.Data.ID, handlers.UserRoleSkylight, DashboardUrl, "GET")})
	castedFetchList := fetchList.(*tenant_provisioning_service_v2.GetAllDashboardsV2OK)
	assert.NotNil(t, castedFetchList)
	assert.Equal(t, 1, len(castedFetchList.Payload.Data))
	assert.Equal(t, castedCreate.Payload.Data, castedFetchList.Payload.Data[0])

	// Make an update to the Record
	newName := fake.CharactersN(16)
	updateRequestBody := generateDashboardUpdateRequest(*castedCreate.Payload.Data.ID, *castedCreate.Payload.Data.Attributes.Rev, &newName, nil, nil)
	updated := handlers.HandleUpdateDashboardV2(handlers.AllRoles, tenantDB)(tenant_provisioning_service_v2.UpdateDashboardV2Params{DashboardID: *castedCreate.Payload.Data.ID, Body: updateRequestBody, HTTPRequest: createHttpRequestWithParams(*castedCreateTeant.Payload.Data.ID, handlers.UserRoleSkylight, DashboardUrl, "PATCH")})
	castedUpdate := updated.(*tenant_provisioning_service_v2.UpdateDashboardV2OK)
	assert.NotNil(t, castedUpdate)
	assert.NotEqual(t, castedCreate.Payload.Data, castedUpdate.Payload.Data)
	assert.NotEqual(t, castedCreate.Payload.Data.Attributes.Rev, castedUpdate.Payload.Data.Attributes.Rev)
	assert.Equal(t, newName, *castedUpdate.Payload.Data.Attributes.Name)
	assert.Equal(t, castedCreate.Payload.Data.Relationships.ThresholdProfile, castedUpdate.Payload.Data.Relationships.ThresholdProfile)
	assert.ElementsMatch(t, castedCreate.Payload.Data.Relationships.Cards.Data, castedUpdate.Payload.Data.Relationships.Cards.Data)

	// Delete the record
	deleted := handlers.HandleDeleteDashboardV2(handlers.AllRoles, tenantDB)(tenant_provisioning_service_v2.DeleteDashboardV2Params{DashboardID: *castedCreate.Payload.Data.ID, HTTPRequest: createHttpRequestWithParams(*castedCreateTeant.Payload.Data.ID, handlers.UserRoleSkylight, DashboardUrl, "DELETE")})
	castedDelete := deleted.(*tenant_provisioning_service_v2.DeleteDashboardV2OK)
	assert.NotNil(t, castedDelete)
	assert.Equal(t, castedUpdate.Payload.Data, castedDelete.Payload.Data)

	// Make sure there are no existing Dashboards
	existing = handlers.HandleGetAllDashboardsV2(handlers.AllRoles, tenantDB)(tenant_provisioning_service_v2.GetAllDashboardsV2Params{HTTPRequest: createHttpRequestWithParams(*castedCreateTeant.Payload.Data.ID, handlers.UserRoleSkylight, DashboardUrl, "GET")})
	castedResponse = existing.(*tenant_provisioning_service_v2.GetAllDashboardsV2NotFound)
	assert.NotNil(t, castedResponse)
}

func TestDashboardNotFoundV2(t *testing.T) {
	notFoundID := fake.CharactersN(20)

	// Get Dashboard
	fetched := handlers.HandleGetDashboardV2(handlers.AllRoles, tenantDB)(tenant_provisioning_service_v2.GetDashboardV2Params{DashboardID: notFoundID, HTTPRequest: createHttpRequestWithParams("", handlers.UserRoleSkylight, DashboardUrl, "GET")})
	castedFetch := fetched.(*tenant_provisioning_service_v2.GetDashboardV2NotFound)
	assert.NotNil(t, castedFetch)

	// Delete Dashboard
	deleted := handlers.HandleDeleteDashboardV2(handlers.AllRoles, tenantDB)(tenant_provisioning_service_v2.DeleteDashboardV2Params{DashboardID: notFoundID, HTTPRequest: createHttpRequestWithParams("", handlers.UserRoleSkylight, DashboardUrl, "DELETE")})
	castedDelete := deleted.(*tenant_provisioning_service_v2.DeleteDashboardV2NotFound)
	assert.NotNil(t, castedDelete)

	// Patch Dashboard
	updateRequest := generateDashboardUpdateRequest(notFoundID, "reviosionstuff", nil, nil, nil)
	updated := handlers.HandleUpdateDashboardV2(handlers.AllRoles, tenantDB)(tenant_provisioning_service_v2.UpdateDashboardV2Params{DashboardID: notFoundID, Body: updateRequest, HTTPRequest: createHttpRequestWithParams("", handlers.UserRoleSkylight, DashboardUrl, "PATCH")})
	castedUpdate := updated.(*tenant_provisioning_service_v2.UpdateDashboardV2NotFound)
	assert.NotNil(t, castedUpdate)
}

func TestDashboardBadRequestV2(t *testing.T) {

	// CreateDashboard
	created := handlers.HandleCreateDashboardV2(handlers.AllRoles, tenantDB)(tenant_provisioning_service_v2.CreateDashboardV2Params{Body: nil, HTTPRequest: createHttpRequestWithParams("", handlers.UserRoleSkylight, DashboardUrl, "POST")})
	castedCreate := created.(*tenant_provisioning_service_v2.CreateDashboardV2BadRequest)
	assert.NotNil(t, castedCreate)

	// Update Dashboard
	updated := handlers.HandleUpdateDashboardV2(handlers.AllRoles, tenantDB)(tenant_provisioning_service_v2.UpdateDashboardV2Params{DashboardID: fake.CharactersN(20), Body: nil, HTTPRequest: createHttpRequestWithParams("", handlers.UserRoleSkylight, DashboardUrl, "PATCH")})
	castedUpdate := updated.(*tenant_provisioning_service_v2.UpdateDashboardV2BadRequest)
	assert.NotNil(t, castedUpdate)
}

func TestDashboardConflictV2(t *testing.T) {

	createdTenant := handlers.HandleCreateTenantV2(handlers.SkylightAdminRoleOnly, adminDB, tenantDB, stubbedMetricBaselineDatastore)(admin_provisioning_service_v2.CreateTenantV2Params{Body: generateRandomTenantCreationRequest(), HTTPRequest: createHttpRequestWithParams("", handlers.UserRoleSkylight, tenantURL, "POST")})
	castedCreateTeant := createdTenant.(*admin_provisioning_service_v2.CreateTenantV2Created)
	assert.NotNil(t, castedCreateTeant)
	assert.NotEmpty(t, castedCreateTeant.Payload.Data.ID)
	assert.NotEmpty(t, castedCreateTeant.Payload.Data.Attributes.Name)
	assert.NotEmpty(t, castedCreateTeant.Payload.Data.Attributes.URLSubdomain)
	assert.Equal(t, string(common.UserActive), *castedCreateTeant.Payload.Data.Attributes.State)
	assert.True(t, *castedCreateTeant.Payload.Data.Attributes.CreatedTimestamp > 0)
	assert.True(t, *castedCreateTeant.Payload.Data.Attributes.LastModifiedTimestamp > 0)

	// Make sure there are no existing Dashboard
	existing := handlers.HandleGetAllDashboardsV2(handlers.AllRoles, tenantDB)(tenant_provisioning_service_v2.GetAllDashboardsV2Params{HTTPRequest: createHttpRequestWithParams(*castedCreateTeant.Payload.Data.ID, handlers.UserRoleSkylight, DashboardUrl, "GET")})
	castedResponse := existing.(*tenant_provisioning_service_v2.GetAllDashboardsV2NotFound)
	assert.NotNil(t, castedResponse)

	createReqBody := generateRandomTenantDashboardCreationRequest()
	created := handlers.HandleCreateDashboardV2(handlers.AllRoles, tenantDB)(tenant_provisioning_service_v2.CreateDashboardV2Params{Body: createReqBody, HTTPRequest: createHttpRequestWithParams(*castedCreateTeant.Payload.Data.ID, handlers.UserRoleSkylight, DashboardUrl, "POST")})
	castedCreate := created.(*tenant_provisioning_service_v2.CreateDashboardV2Created)
	assert.NotNil(t, castedCreate)
	assert.NotEmpty(t, castedCreate.Payload.Data.ID)
	assert.NotEmpty(t, castedCreate.Payload.Data.Attributes.Name)
	assert.NotEmpty(t, castedCreate.Payload.Data.Attributes.Category)
	assert.NotEmpty(t, castedCreate.Payload.Data.Relationships.ThresholdProfile)
	assert.NotEmpty(t, castedCreate.Payload.Data.Relationships.Cards.Data)
	assert.NotEmpty(t, castedCreate.Payload.Data.Attributes.Datatype)
	assert.True(t, *castedCreate.Payload.Data.Attributes.CreatedTimestamp > 0)
	assert.True(t, *castedCreate.Payload.Data.Attributes.LastModifiedTimestamp > 0)

	// Try to create the record again - should succeed as we are not guarding against name collisiones
	createdConflict := handlers.HandleCreateDashboardV2(handlers.AllRoles, tenantDB)(tenant_provisioning_service_v2.CreateDashboardV2Params{Body: createReqBody, HTTPRequest: createHttpRequestWithParams(*castedCreateTeant.Payload.Data.ID, handlers.UserRoleSkylight, DashboardUrl, "POST")})
	castedCreateConflictButOK := createdConflict.(*tenant_provisioning_service_v2.CreateDashboardV2Created)
	assert.NotNil(t, castedCreateConflictButOK)

	// Try the update with a bad revision
	newName := fake.CharactersN(16)
	updateRequestBody := generateDashboardUpdateRequest(*castedCreate.Payload.Data.ID, *castedCreate.Payload.Data.Attributes.Rev+"pork", &newName, nil, nil)
	updated := handlers.HandleUpdateDashboardV2(handlers.AllRoles, tenantDB)(tenant_provisioning_service_v2.UpdateDashboardV2Params{DashboardID: *castedCreate.Payload.Data.ID, Body: updateRequestBody, HTTPRequest: createHttpRequestWithParams(*castedCreateTeant.Payload.Data.ID, handlers.UserRoleSkylight, DashboardUrl, "PATCH")})
	castedUpdate := updated.(*tenant_provisioning_service_v2.UpdateDashboardV2Conflict)
	assert.NotNil(t, castedUpdate)

	// Delete the tenant
	deleted := handlers.HandleDeleteDashboardV2(handlers.AllRoles, tenantDB)(tenant_provisioning_service_v2.DeleteDashboardV2Params{DashboardID: *castedCreate.Payload.Data.ID, HTTPRequest: createHttpRequestWithParams(*castedCreateTeant.Payload.Data.ID, handlers.UserRoleSkylight, DashboardUrl, "DELETE")})
	castedDelete := deleted.(*tenant_provisioning_service_v2.DeleteDashboardV2OK)
	assert.NotNil(t, castedDelete)
	assert.NotNil(t, castedDelete.Payload.Data)
}

func TestDashboardAPIsProtectedByAuthV2(t *testing.T) {
	fakeTenantID := fake.CharactersN(20)
	// Get All - All Users
	existing := handlers.HandleGetAllDashboardsV2(handlers.AllRoles, tenantDB)(tenant_provisioning_service_v2.GetAllDashboardsV2Params{HTTPRequest: createHttpRequestWithParams(fakeTenantID, handlers.UserRoleUnknown, DashboardUrl, "GET")})
	castedResponse := existing.(*tenant_provisioning_service_v2.GetAllDashboardsV2Forbidden)
	assert.NotNil(t, castedResponse)

	fakeID := fake.CharactersN(20)

	// Get - All Users
	fetched := handlers.HandleGetDashboardV2(handlers.AllRoles, tenantDB)(tenant_provisioning_service_v2.GetDashboardV2Params{DashboardID: fakeID, HTTPRequest: createHttpRequestWithParams(fakeTenantID, handlers.UserRoleUnknown, DashboardUrl, "GET")})
	castedFetch := fetched.(*tenant_provisioning_service_v2.GetDashboardV2Forbidden)
	assert.NotNil(t, castedFetch)

	// Create - SkylightAdmin and TenantAdmin Only
	created := handlers.HandleCreateDashboardV2(handlers.SkylightAndTenantAdminRoles, tenantDB)(tenant_provisioning_service_v2.CreateDashboardV2Params{Body: generateRandomTenantDashboardCreationRequest(), HTTPRequest: createHttpRequestWithParams(fakeTenantID, handlers.UserRoleUnknown, DashboardUrl, "POST")})
	castedCreate := created.(*tenant_provisioning_service_v2.CreateDashboardV2Forbidden)
	assert.NotNil(t, castedCreate)

	created = handlers.HandleCreateDashboardV2(handlers.SkylightAndTenantAdminRoles, tenantDB)(tenant_provisioning_service_v2.CreateDashboardV2Params{Body: generateRandomTenantDashboardCreationRequest(), HTTPRequest: createHttpRequestWithParams(fakeTenantID, handlers.UserRoleTenantUser, DashboardUrl, "POST")})
	castedCreate = created.(*tenant_provisioning_service_v2.CreateDashboardV2Forbidden)
	assert.NotNil(t, castedCreate)

	// Update - SkylightAdmin and TenantAdmin Only
	updated := handlers.HandleUpdateDashboardV2(handlers.SkylightAndTenantAdminRoles, tenantDB)(tenant_provisioning_service_v2.UpdateDashboardV2Params{DashboardID: fakeID, Body: nil, HTTPRequest: createHttpRequestWithParams(fakeTenantID, handlers.UserRoleUnknown, DashboardUrl, "PATCH")})
	castedUpdate := updated.(*tenant_provisioning_service_v2.UpdateDashboardV2Forbidden)
	assert.NotNil(t, castedUpdate)

	updated = handlers.HandleUpdateDashboardV2(handlers.SkylightAndTenantAdminRoles, tenantDB)(tenant_provisioning_service_v2.UpdateDashboardV2Params{DashboardID: fakeID, Body: nil, HTTPRequest: createHttpRequestWithParams(fakeTenantID, handlers.UserRoleTenantUser, DashboardUrl, "PATCH")})
	castedUpdate = updated.(*tenant_provisioning_service_v2.UpdateDashboardV2Forbidden)
	assert.NotNil(t, castedUpdate)

	// Delete - SkylightAdmin and TenantAdmin Only
	deleted := handlers.HandleDeleteDashboardV2(handlers.SkylightAndTenantAdminRoles, tenantDB)(tenant_provisioning_service_v2.DeleteDashboardV2Params{DashboardID: fakeID, HTTPRequest: createHttpRequestWithParams(fakeTenantID, handlers.UserRoleUnknown, DashboardUrl, "DELETE")})
	castedDelete := deleted.(*tenant_provisioning_service_v2.DeleteDashboardV2Forbidden)
	assert.NotNil(t, castedDelete)

	deleted = handlers.HandleDeleteDashboardV2(handlers.SkylightAndTenantAdminRoles, tenantDB)(tenant_provisioning_service_v2.DeleteDashboardV2Params{DashboardID: fakeID, HTTPRequest: createHttpRequestWithParams(fakeTenantID, handlers.UserRoleTenantUser, DashboardUrl, "DELETE")})
	castedDelete = deleted.(*tenant_provisioning_service_v2.DeleteDashboardV2Forbidden)
	assert.NotNil(t, castedDelete)
}

func generateRandomTenantDashboardCreationRequest() *swagmodels.DashboardCreateRequest {
	name := fake.CharactersN(8)
	return &swagmodels.DashboardCreateRequest{
		Data: &swagmodels.DashboardCreateRequestData{
			Type: &DashboardTypeString,
			Attributes: &swagmodels.DashboardCreateRequestDataAttributes{
				Category: fake.CharactersN(12),
				Name:     &name,
			},
			Relationships: &swagmodels.DashboardRelationships{
				Cards: &swagmodels.JSONAPIRelationship{
					Data: []*swagmodels.JSONAPIRelationshipData{
						&swagmodels.JSONAPIRelationshipData{
							Type: "cards",
							ID:   fake.CharactersN(12),
						},
						&swagmodels.JSONAPIRelationshipData{
							Type: "cards",
							ID:   fake.CharactersN(12),
						},
					},
				},
				ThresholdProfile: &swagmodels.JSONAPISingleRelationship{
					Data: &swagmodels.JSONAPIRelationshipData{
						Type: "thresholdProfiles",
						ID:   fake.CharactersN(12),
					},
				},
			},
		},
	}
}

func generateDashboardUpdateRequest(id string, rev string, name *string, thresh *string, category *string) *swagmodels.DashboardUpdateRequest {
	result := &swagmodels.DashboardUpdateRequest{
		Data: &swagmodels.DashboardUpdateRequestData{
			Type:       &DashboardTypeString,
			ID:         &id,
			Attributes: &swagmodels.DashboardUpdateRequestDataAttributes{Rev: &rev},
		},
	}

	if name != nil {
		result.Data.Attributes.Name = *name
	}
	if thresh != nil {
		result.Data.Relationships.ThresholdProfile = &swagmodels.JSONAPISingleRelationship{
			Data: &swagmodels.JSONAPIRelationshipData{
				Type: "thresholdProfiles",
				ID:   *thresh,
			},
		}
	}

	if category != nil {
		result.Data.Attributes.Category = *category
	}

	return result
}
