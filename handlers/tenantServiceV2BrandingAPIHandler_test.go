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
	BrandingUrl = "http://deployment.test.cool/api/v2/brandings"

	BrandingTypeString = "brandings"
)

func TestBrandingCRUDV2(t *testing.T) {

	createdTenant := handlers.HandleCreateTenantV2(handlers.SkylightAdminRoleOnly, adminDB, tenantDB)(admin_provisioning_service_v2.CreateTenantV2Params{Body: generateRandomTenantCreationRequest(), HTTPRequest: createHttpRequestWithParams("", handlers.UserRoleSkylight, tenantURL, "POST")})
	castedCreateTeant := createdTenant.(*admin_provisioning_service_v2.CreateTenantV2Created)
	assert.NotNil(t, castedCreateTeant)
	assert.NotEmpty(t, castedCreateTeant.Payload.Data.ID)
	assert.NotEmpty(t, castedCreateTeant.Payload.Data.Attributes.Name)
	assert.NotEmpty(t, castedCreateTeant.Payload.Data.Attributes.URLSubdomain)
	assert.Equal(t, string(common.UserActive), *castedCreateTeant.Payload.Data.Attributes.State)
	assert.True(t, *castedCreateTeant.Payload.Data.Attributes.CreatedTimestamp > 0)
	assert.True(t, *castedCreateTeant.Payload.Data.Attributes.LastModifiedTimestamp > 0)

	// Make sure there are no existing Brandings
	existing := handlers.HandleGetAllBrandingsV2(handlers.AllRoles, tenantDB)(tenant_provisioning_service_v2.GetAllBrandingsV2Params{HTTPRequest: createHttpRequestWithParams(*castedCreateTeant.Payload.Data.ID, handlers.UserRoleSkylight, BrandingUrl, "GET")})
	castedResponse := existing.(*tenant_provisioning_service_v2.GetAllBrandingsV2NotFound)
	assert.NotNil(t, castedResponse)

	created := handlers.HandleCreateBrandingV2(handlers.AllRoles, tenantDB)(tenant_provisioning_service_v2.CreateBrandingV2Params{Body: generateRandomTenantBrandingCreationRequest(), HTTPRequest: createHttpRequestWithParams(*castedCreateTeant.Payload.Data.ID, handlers.UserRoleSkylight, BrandingUrl, "POST")})
	castedCreate := created.(*tenant_provisioning_service_v2.CreateBrandingV2Created)
	assert.NotNil(t, castedCreate)
	assert.NotEmpty(t, castedCreate.Payload.Data.ID)
	assert.NotEmpty(t, castedCreate.Payload.Data.Attributes.Color)
	assert.NotEmpty(t, castedCreate.Payload.Data.Attributes.Logo)
	assert.NotEmpty(t, castedCreate.Payload.Data.Attributes.Datatype)
	assert.True(t, *castedCreate.Payload.Data.Attributes.CreatedTimestamp > 0)
	assert.True(t, *castedCreate.Payload.Data.Attributes.LastModifiedTimestamp > 0)

	// Make sure we can retrieve this record:
	fetched := handlers.HandleGetBrandingV2(handlers.AllRoles, tenantDB)(tenant_provisioning_service_v2.GetBrandingV2Params{BrandingID: *castedCreate.Payload.Data.ID, HTTPRequest: createHttpRequestWithParams(*castedCreateTeant.Payload.Data.ID, handlers.UserRoleSkylight, BrandingUrl, "GET")})
	castedFetch := fetched.(*tenant_provisioning_service_v2.GetBrandingV2OK)
	assert.NotNil(t, castedFetch)
	assert.Equal(t, castedCreate.Payload.Data, castedFetch.Payload.Data)

	// Also retrieve the record as part of an array
	fetchList := handlers.HandleGetAllBrandingsV2(handlers.AllRoles, tenantDB)(tenant_provisioning_service_v2.GetAllBrandingsV2Params{HTTPRequest: createHttpRequestWithParams(*castedCreateTeant.Payload.Data.ID, handlers.UserRoleSkylight, BrandingUrl, "GET")})
	castedFetchList := fetchList.(*tenant_provisioning_service_v2.GetAllBrandingsV2OK)
	assert.NotNil(t, castedFetchList)
	assert.Equal(t, 1, len(castedFetchList.Payload.Data))
	assert.Equal(t, castedCreate.Payload.Data, castedFetchList.Payload.Data[0])

	// Make an update to the Record
	newName := fake.CharactersN(16)
	updateRequestBody := generateBrandingUpdateRequest(*castedCreate.Payload.Data.ID, *castedCreate.Payload.Data.Attributes.Rev, &newName, nil)
	updated := handlers.HandleUpdateBrandingV2(handlers.AllRoles, tenantDB)(tenant_provisioning_service_v2.UpdateBrandingV2Params{BrandingID: *castedCreate.Payload.Data.ID, Body: updateRequestBody, HTTPRequest: createHttpRequestWithParams(*castedCreateTeant.Payload.Data.ID, handlers.UserRoleSkylight, BrandingUrl, "PATCH")})
	castedUpdate := updated.(*tenant_provisioning_service_v2.UpdateBrandingV2OK)
	assert.NotNil(t, castedUpdate)
	assert.NotEqual(t, castedCreate.Payload.Data, castedUpdate.Payload.Data)
	assert.NotEqual(t, castedCreate.Payload.Data.Attributes.Rev, castedUpdate.Payload.Data.Attributes.Rev)
	assert.Equal(t, newName, *castedUpdate.Payload.Data.Attributes.Color)
	assert.Equal(t, castedCreate.Payload.Data.Attributes.Logo, castedUpdate.Payload.Data.Attributes.Logo)

	// Delete the record
	deleted := handlers.HandleDeleteBrandingV2(handlers.AllRoles, tenantDB)(tenant_provisioning_service_v2.DeleteBrandingV2Params{BrandingID: *castedCreate.Payload.Data.ID, HTTPRequest: createHttpRequestWithParams(*castedCreateTeant.Payload.Data.ID, handlers.UserRoleSkylight, BrandingUrl, "DELETE")})
	castedDelete := deleted.(*tenant_provisioning_service_v2.DeleteBrandingV2OK)
	assert.NotNil(t, castedDelete)
	assert.Equal(t, castedUpdate.Payload.Data, castedDelete.Payload.Data)

	// Make sure there are no existing Brandings
	existing = handlers.HandleGetAllBrandingsV2(handlers.AllRoles, tenantDB)(tenant_provisioning_service_v2.GetAllBrandingsV2Params{HTTPRequest: createHttpRequestWithParams(*castedCreateTeant.Payload.Data.ID, handlers.UserRoleSkylight, BrandingUrl, "GET")})
	castedResponse = existing.(*tenant_provisioning_service_v2.GetAllBrandingsV2NotFound)
	assert.NotNil(t, castedResponse)
}

func TestBrandingNotFoundV2(t *testing.T) {
	notFoundID := fake.CharactersN(20)

	// Get Branding
	fetched := handlers.HandleGetBrandingV2(handlers.AllRoles, tenantDB)(tenant_provisioning_service_v2.GetBrandingV2Params{BrandingID: notFoundID, HTTPRequest: createHttpRequestWithParams("", handlers.UserRoleSkylight, BrandingUrl, "GET")})
	castedFetch := fetched.(*tenant_provisioning_service_v2.GetBrandingV2NotFound)
	assert.NotNil(t, castedFetch)

	// Delete Branding
	deleted := handlers.HandleDeleteBrandingV2(handlers.AllRoles, tenantDB)(tenant_provisioning_service_v2.DeleteBrandingV2Params{BrandingID: notFoundID, HTTPRequest: createHttpRequestWithParams("", handlers.UserRoleSkylight, BrandingUrl, "DELETE")})
	castedDelete := deleted.(*tenant_provisioning_service_v2.DeleteBrandingV2NotFound)
	assert.NotNil(t, castedDelete)

	// Patch Branding
	updateRequest := generateBrandingUpdateRequest(notFoundID, "reviosionstuff", nil, nil)
	updated := handlers.HandleUpdateBrandingV2(handlers.AllRoles, tenantDB)(tenant_provisioning_service_v2.UpdateBrandingV2Params{BrandingID: notFoundID, Body: updateRequest, HTTPRequest: createHttpRequestWithParams("", handlers.UserRoleSkylight, BrandingUrl, "PATCH")})
	castedUpdate := updated.(*tenant_provisioning_service_v2.UpdateBrandingV2NotFound)
	assert.NotNil(t, castedUpdate)
}

func TestBrandingBadRequestV2(t *testing.T) {

	// CreateBranding
	created := handlers.HandleCreateBrandingV2(handlers.AllRoles, tenantDB)(tenant_provisioning_service_v2.CreateBrandingV2Params{Body: nil, HTTPRequest: createHttpRequestWithParams("", handlers.UserRoleSkylight, BrandingUrl, "POST")})
	castedCreate := created.(*tenant_provisioning_service_v2.CreateBrandingV2BadRequest)
	assert.NotNil(t, castedCreate)

	// Update Branding
	updated := handlers.HandleUpdateBrandingV2(handlers.AllRoles, tenantDB)(tenant_provisioning_service_v2.UpdateBrandingV2Params{BrandingID: fake.CharactersN(20), Body: nil, HTTPRequest: createHttpRequestWithParams("", handlers.UserRoleSkylight, BrandingUrl, "PATCH")})
	castedUpdate := updated.(*tenant_provisioning_service_v2.UpdateBrandingV2BadRequest)
	assert.NotNil(t, castedUpdate)
}

func TestBrandingConflictV2(t *testing.T) {

	createdTenant := handlers.HandleCreateTenantV2(handlers.SkylightAdminRoleOnly, adminDB, tenantDB)(admin_provisioning_service_v2.CreateTenantV2Params{Body: generateRandomTenantCreationRequest(), HTTPRequest: createHttpRequestWithParams("", handlers.UserRoleSkylight, tenantURL, "POST")})
	castedCreateTeant := createdTenant.(*admin_provisioning_service_v2.CreateTenantV2Created)
	assert.NotNil(t, castedCreateTeant)
	assert.NotEmpty(t, castedCreateTeant.Payload.Data.ID)
	assert.NotEmpty(t, castedCreateTeant.Payload.Data.Attributes.Name)
	assert.NotEmpty(t, castedCreateTeant.Payload.Data.Attributes.URLSubdomain)
	assert.Equal(t, string(common.UserActive), *castedCreateTeant.Payload.Data.Attributes.State)
	assert.True(t, *castedCreateTeant.Payload.Data.Attributes.CreatedTimestamp > 0)
	assert.True(t, *castedCreateTeant.Payload.Data.Attributes.LastModifiedTimestamp > 0)

	// Make sure there are no existing Branding
	existing := handlers.HandleGetAllBrandingsV2(handlers.AllRoles, tenantDB)(tenant_provisioning_service_v2.GetAllBrandingsV2Params{HTTPRequest: createHttpRequestWithParams(*castedCreateTeant.Payload.Data.ID, handlers.UserRoleSkylight, BrandingUrl, "GET")})
	castedResponse := existing.(*tenant_provisioning_service_v2.GetAllBrandingsV2NotFound)
	assert.NotNil(t, castedResponse)

	createReqBody := generateRandomTenantBrandingCreationRequest()
	created := handlers.HandleCreateBrandingV2(handlers.AllRoles, tenantDB)(tenant_provisioning_service_v2.CreateBrandingV2Params{Body: createReqBody, HTTPRequest: createHttpRequestWithParams(*castedCreateTeant.Payload.Data.ID, handlers.UserRoleSkylight, BrandingUrl, "POST")})
	castedCreate := created.(*tenant_provisioning_service_v2.CreateBrandingV2Created)
	assert.NotNil(t, castedCreate)
	assert.NotEmpty(t, castedCreate.Payload.Data.ID)
	assert.NotEmpty(t, castedCreate.Payload.Data.Attributes.Color)
	assert.NotEmpty(t, castedCreate.Payload.Data.Attributes.Logo)
	assert.NotEmpty(t, castedCreate.Payload.Data.Attributes.Datatype)
	assert.True(t, *castedCreate.Payload.Data.Attributes.CreatedTimestamp > 0)
	assert.True(t, *castedCreate.Payload.Data.Attributes.LastModifiedTimestamp > 0)

	// Try to create the record again - should succeed as we are not guarding against name collisiones
	createdConflict := handlers.HandleCreateBrandingV2(handlers.AllRoles, tenantDB)(tenant_provisioning_service_v2.CreateBrandingV2Params{Body: createReqBody, HTTPRequest: createHttpRequestWithParams(*castedCreateTeant.Payload.Data.ID, handlers.UserRoleSkylight, BrandingUrl, "POST")})
	castedCreateConflictButOK := createdConflict.(*tenant_provisioning_service_v2.CreateBrandingV2Created)
	assert.NotNil(t, castedCreateConflictButOK)

	// Try the update with a bad revision
	newName := fake.CharactersN(16)
	updateRequestBody := generateBrandingUpdateRequest(*castedCreate.Payload.Data.ID, *castedCreate.Payload.Data.Attributes.Rev+"pork", &newName, nil)
	updated := handlers.HandleUpdateBrandingV2(handlers.AllRoles, tenantDB)(tenant_provisioning_service_v2.UpdateBrandingV2Params{BrandingID: *castedCreate.Payload.Data.ID, Body: updateRequestBody, HTTPRequest: createHttpRequestWithParams(*castedCreateTeant.Payload.Data.ID, handlers.UserRoleSkylight, BrandingUrl, "PATCH")})
	castedUpdate := updated.(*tenant_provisioning_service_v2.UpdateBrandingV2Conflict)
	assert.NotNil(t, castedUpdate)

	// Delete the tenant
	deleted := handlers.HandleDeleteBrandingV2(handlers.AllRoles, tenantDB)(tenant_provisioning_service_v2.DeleteBrandingV2Params{BrandingID: *castedCreate.Payload.Data.ID, HTTPRequest: createHttpRequestWithParams(*castedCreateTeant.Payload.Data.ID, handlers.UserRoleSkylight, BrandingUrl, "DELETE")})
	castedDelete := deleted.(*tenant_provisioning_service_v2.DeleteBrandingV2OK)
	assert.NotNil(t, castedDelete)
	assert.NotNil(t, castedDelete.Payload.Data)
}

func TestBrandingAPIsProtectedByAuthV2(t *testing.T) {
	fakeTenantID := fake.CharactersN(20)
	// Get All - All Users
	existing := handlers.HandleGetAllBrandingsV2(handlers.AllRoles, tenantDB)(tenant_provisioning_service_v2.GetAllBrandingsV2Params{HTTPRequest: createHttpRequestWithParams(fakeTenantID, handlers.UserRoleUnknown, BrandingUrl, "GET")})
	castedResponse := existing.(*tenant_provisioning_service_v2.GetAllBrandingsV2Forbidden)
	assert.NotNil(t, castedResponse)

	fakeID := fake.CharactersN(20)

	// Get - All Users
	fetched := handlers.HandleGetBrandingV2(handlers.AllRoles, tenantDB)(tenant_provisioning_service_v2.GetBrandingV2Params{BrandingID: fakeID, HTTPRequest: createHttpRequestWithParams(fakeTenantID, handlers.UserRoleUnknown, BrandingUrl, "GET")})
	castedFetch := fetched.(*tenant_provisioning_service_v2.GetBrandingV2Forbidden)
	assert.NotNil(t, castedFetch)

	// Create - SkylightAdmin and TenantAdmin Only
	created := handlers.HandleCreateBrandingV2(handlers.SkylightAndTenantAdminRoles, tenantDB)(tenant_provisioning_service_v2.CreateBrandingV2Params{Body: generateRandomTenantBrandingCreationRequest(), HTTPRequest: createHttpRequestWithParams(fakeTenantID, handlers.UserRoleUnknown, BrandingUrl, "POST")})
	castedCreate := created.(*tenant_provisioning_service_v2.CreateBrandingV2Forbidden)
	assert.NotNil(t, castedCreate)

	created = handlers.HandleCreateBrandingV2(handlers.SkylightAndTenantAdminRoles, tenantDB)(tenant_provisioning_service_v2.CreateBrandingV2Params{Body: generateRandomTenantBrandingCreationRequest(), HTTPRequest: createHttpRequestWithParams(fakeTenantID, handlers.UserRoleTenantUser, BrandingUrl, "POST")})
	castedCreate = created.(*tenant_provisioning_service_v2.CreateBrandingV2Forbidden)
	assert.NotNil(t, castedCreate)

	// Update - SkylightAdmin and TenantAdmin Only
	updated := handlers.HandleUpdateBrandingV2(handlers.SkylightAndTenantAdminRoles, tenantDB)(tenant_provisioning_service_v2.UpdateBrandingV2Params{BrandingID: fakeID, Body: nil, HTTPRequest: createHttpRequestWithParams(fakeTenantID, handlers.UserRoleUnknown, BrandingUrl, "PATCH")})
	castedUpdate := updated.(*tenant_provisioning_service_v2.UpdateBrandingV2Forbidden)
	assert.NotNil(t, castedUpdate)

	updated = handlers.HandleUpdateBrandingV2(handlers.SkylightAndTenantAdminRoles, tenantDB)(tenant_provisioning_service_v2.UpdateBrandingV2Params{BrandingID: fakeID, Body: nil, HTTPRequest: createHttpRequestWithParams(fakeTenantID, handlers.UserRoleTenantUser, BrandingUrl, "PATCH")})
	castedUpdate = updated.(*tenant_provisioning_service_v2.UpdateBrandingV2Forbidden)
	assert.NotNil(t, castedUpdate)

	// Delete - SkylightAdmin and TenantAdmin Only
	deleted := handlers.HandleDeleteBrandingV2(handlers.SkylightAndTenantAdminRoles, tenantDB)(tenant_provisioning_service_v2.DeleteBrandingV2Params{BrandingID: fakeID, HTTPRequest: createHttpRequestWithParams(fakeTenantID, handlers.UserRoleUnknown, BrandingUrl, "DELETE")})
	castedDelete := deleted.(*tenant_provisioning_service_v2.DeleteBrandingV2Forbidden)
	assert.NotNil(t, castedDelete)

	deleted = handlers.HandleDeleteBrandingV2(handlers.SkylightAndTenantAdminRoles, tenantDB)(tenant_provisioning_service_v2.DeleteBrandingV2Params{BrandingID: fakeID, HTTPRequest: createHttpRequestWithParams(fakeTenantID, handlers.UserRoleTenantUser, BrandingUrl, "DELETE")})
	castedDelete = deleted.(*tenant_provisioning_service_v2.DeleteBrandingV2Forbidden)
	assert.NotNil(t, castedDelete)
}

func generateRandomTenantBrandingCreationRequest() *swagmodels.BrandingCreateRequest {
	return &swagmodels.BrandingCreateRequest{
		Data: &swagmodels.BrandingCreateRequestData{
			Type: &BrandingTypeString,
			Attributes: &swagmodels.BrandingCreateRequestDataAttributes{
				Color: fake.CharactersN(6),
				Logo: &swagmodels.BrandingLogo{
					File: &swagmodels.BrandingLogoFile{
						ContentType: fake.CharactersN(12),
						Data:        fake.CharactersN(100),
					},
				},
			},
		},
	}
}

func generateBrandingUpdateRequest(id string, rev string, color *string, imageData *string) *swagmodels.BrandingUpdateRequest {
	result := &swagmodels.BrandingUpdateRequest{
		Data: &swagmodels.BrandingUpdateRequestData{
			Type:       &BrandingTypeString,
			ID:         &id,
			Attributes: &swagmodels.BrandingUpdateRequestDataAttributes{Rev: &rev},
		},
	}

	if color != nil {
		result.Data.Attributes.Color = *color
	}

	if imageData != nil {
		result.Data.Attributes.Logo.File.Data = *imageData
	}

	return result
}
