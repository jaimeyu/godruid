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
	CardUrl = "http://deployment.test.cool/api/v2/cards"

	CardTypeString = "cards"
)

func TestCardCRUDV2(t *testing.T) {

	createdTenant := handlers.HandleCreateTenantV2(handlers.SkylightAdminRoleOnly, adminDB, tenantDB, stubbedMetricBaselineDatastore)(admin_provisioning_service_v2.CreateTenantV2Params{Body: generateRandomTenantCreationRequest(), HTTPRequest: createHttpRequestWithParams("", handlers.UserRoleSkylight, tenantURL, "POST")})
	castedCreateTeant := createdTenant.(*admin_provisioning_service_v2.CreateTenantV2Created)
	assert.NotNil(t, castedCreateTeant)
	assert.NotEmpty(t, castedCreateTeant.Payload.Data.ID)
	assert.NotEmpty(t, castedCreateTeant.Payload.Data.Attributes.Name)
	assert.NotEmpty(t, castedCreateTeant.Payload.Data.Attributes.URLSubdomain)
	assert.Equal(t, string(common.UserActive), *castedCreateTeant.Payload.Data.Attributes.State)
	assert.True(t, *castedCreateTeant.Payload.Data.Attributes.CreatedTimestamp > 0)
	assert.True(t, *castedCreateTeant.Payload.Data.Attributes.LastModifiedTimestamp > 0)

	// Make sure there are no existing Cards
	existing := handlers.HandleGetAllCardsV2(handlers.AllRoles, tenantDB)(tenant_provisioning_service_v2.GetAllCardsV2Params{HTTPRequest: createHttpRequestWithParams(*castedCreateTeant.Payload.Data.ID, handlers.UserRoleSkylight, CardUrl, "GET")})
	castedResponse := existing.(*tenant_provisioning_service_v2.GetAllCardsV2NotFound)
	assert.NotNil(t, castedResponse)

	created := handlers.HandleCreateCardV2(handlers.AllRoles, tenantDB)(tenant_provisioning_service_v2.CreateCardV2Params{Body: generateRandomTenantCardCreationRequest(), HTTPRequest: createHttpRequestWithParams(*castedCreateTeant.Payload.Data.ID, handlers.UserRoleSkylight, CardUrl, "POST")})
	castedCreate := created.(*tenant_provisioning_service_v2.CreateCardV2Created)
	assert.NotNil(t, castedCreate)
	assert.NotEmpty(t, castedCreate.Payload.Data.ID)
	assert.NotEmpty(t, castedCreate.Payload.Data.Attributes.Name)
	assert.NotEmpty(t, castedCreate.Payload.Data.Attributes.Description)
	assert.NotEmpty(t, castedCreate.Payload.Data.Attributes.Datatype)
	assert.True(t, *castedCreate.Payload.Data.Attributes.CreatedTimestamp > 0)
	assert.True(t, *castedCreate.Payload.Data.Attributes.LastModifiedTimestamp > 0)

	// Make sure we can retrieve this record:
	fetched := handlers.HandleGetCardV2(handlers.AllRoles, tenantDB)(tenant_provisioning_service_v2.GetCardV2Params{CardID: *castedCreate.Payload.Data.ID, HTTPRequest: createHttpRequestWithParams(*castedCreateTeant.Payload.Data.ID, handlers.UserRoleSkylight, CardUrl, "GET")})
	castedFetch := fetched.(*tenant_provisioning_service_v2.GetCardV2OK)
	assert.NotNil(t, castedFetch)
	assert.Equal(t, castedCreate.Payload.Data, castedFetch.Payload.Data)

	// Also retrieve the record as part of an array
	fetchList := handlers.HandleGetAllCardsV2(handlers.AllRoles, tenantDB)(tenant_provisioning_service_v2.GetAllCardsV2Params{HTTPRequest: createHttpRequestWithParams(*castedCreateTeant.Payload.Data.ID, handlers.UserRoleSkylight, CardUrl, "GET")})
	castedFetchList := fetchList.(*tenant_provisioning_service_v2.GetAllCardsV2OK)
	assert.NotNil(t, castedFetchList)
	assert.Equal(t, 1, len(castedFetchList.Payload.Data))
	assert.Equal(t, castedCreate.Payload.Data, castedFetchList.Payload.Data[0])

	// Make an update to the Record
	newName := fake.CharactersN(16)
	updateRequestBody := generateCardUpdateRequest(*castedCreate.Payload.Data.ID, *castedCreate.Payload.Data.Attributes.Rev, &newName, nil)
	updated := handlers.HandleUpdateCardV2(handlers.AllRoles, tenantDB)(tenant_provisioning_service_v2.UpdateCardV2Params{CardID: *castedCreate.Payload.Data.ID, Body: updateRequestBody, HTTPRequest: createHttpRequestWithParams(*castedCreateTeant.Payload.Data.ID, handlers.UserRoleSkylight, CardUrl, "PATCH")})
	castedUpdate := updated.(*tenant_provisioning_service_v2.UpdateCardV2OK)
	assert.NotNil(t, castedUpdate)
	assert.NotEqual(t, castedCreate.Payload.Data, castedUpdate.Payload.Data)
	assert.NotEqual(t, castedCreate.Payload.Data.Attributes.Rev, castedUpdate.Payload.Data.Attributes.Rev)
	assert.Equal(t, newName, *castedUpdate.Payload.Data.Attributes.Name)
	assert.Equal(t, castedCreate.Payload.Data.Attributes.Description, castedUpdate.Payload.Data.Attributes.Description)

	// Delete the record
	deleted := handlers.HandleDeleteCardV2(handlers.AllRoles, tenantDB)(tenant_provisioning_service_v2.DeleteCardV2Params{CardID: *castedCreate.Payload.Data.ID, HTTPRequest: createHttpRequestWithParams(*castedCreateTeant.Payload.Data.ID, handlers.UserRoleSkylight, CardUrl, "DELETE")})
	castedDelete := deleted.(*tenant_provisioning_service_v2.DeleteCardV2OK)
	assert.NotNil(t, castedDelete)
	assert.Equal(t, castedUpdate.Payload.Data, castedDelete.Payload.Data)

	// Make sure there are no existing Cards
	existing = handlers.HandleGetAllCardsV2(handlers.AllRoles, tenantDB)(tenant_provisioning_service_v2.GetAllCardsV2Params{HTTPRequest: createHttpRequestWithParams(*castedCreateTeant.Payload.Data.ID, handlers.UserRoleSkylight, CardUrl, "GET")})
	castedResponse = existing.(*tenant_provisioning_service_v2.GetAllCardsV2NotFound)
	assert.NotNil(t, castedResponse)
}

func TestCardNotFoundV2(t *testing.T) {
	notFoundID := fake.CharactersN(20)

	// Get Card
	fetched := handlers.HandleGetCardV2(handlers.AllRoles, tenantDB)(tenant_provisioning_service_v2.GetCardV2Params{CardID: notFoundID, HTTPRequest: createHttpRequestWithParams("", handlers.UserRoleSkylight, CardUrl, "GET")})
	castedFetch := fetched.(*tenant_provisioning_service_v2.GetCardV2NotFound)
	assert.NotNil(t, castedFetch)

	// Delete Card
	deleted := handlers.HandleDeleteCardV2(handlers.AllRoles, tenantDB)(tenant_provisioning_service_v2.DeleteCardV2Params{CardID: notFoundID, HTTPRequest: createHttpRequestWithParams("", handlers.UserRoleSkylight, CardUrl, "DELETE")})
	castedDelete := deleted.(*tenant_provisioning_service_v2.DeleteCardV2NotFound)
	assert.NotNil(t, castedDelete)

	// Patch Card
	updateRequest := generateCardUpdateRequest(notFoundID, "reviosionstuff", nil, nil)
	updated := handlers.HandleUpdateCardV2(handlers.AllRoles, tenantDB)(tenant_provisioning_service_v2.UpdateCardV2Params{CardID: notFoundID, Body: updateRequest, HTTPRequest: createHttpRequestWithParams("", handlers.UserRoleSkylight, CardUrl, "PATCH")})
	castedUpdate := updated.(*tenant_provisioning_service_v2.UpdateCardV2NotFound)
	assert.NotNil(t, castedUpdate)
}

func TestCardBadRequestV2(t *testing.T) {

	// CreateCard
	created := handlers.HandleCreateCardV2(handlers.AllRoles, tenantDB)(tenant_provisioning_service_v2.CreateCardV2Params{Body: nil, HTTPRequest: createHttpRequestWithParams("", handlers.UserRoleSkylight, CardUrl, "POST")})
	castedCreate := created.(*tenant_provisioning_service_v2.CreateCardV2BadRequest)
	assert.NotNil(t, castedCreate)

	// Update Card
	updated := handlers.HandleUpdateCardV2(handlers.AllRoles, tenantDB)(tenant_provisioning_service_v2.UpdateCardV2Params{CardID: fake.CharactersN(20), Body: nil, HTTPRequest: createHttpRequestWithParams("", handlers.UserRoleSkylight, CardUrl, "PATCH")})
	castedUpdate := updated.(*tenant_provisioning_service_v2.UpdateCardV2BadRequest)
	assert.NotNil(t, castedUpdate)
}

func TestCardConflictV2(t *testing.T) {

	createdTenant := handlers.HandleCreateTenantV2(handlers.SkylightAdminRoleOnly, adminDB, tenantDB, stubbedMetricBaselineDatastore)(admin_provisioning_service_v2.CreateTenantV2Params{Body: generateRandomTenantCreationRequest(), HTTPRequest: createHttpRequestWithParams("", handlers.UserRoleSkylight, tenantURL, "POST")})
	castedCreateTeant := createdTenant.(*admin_provisioning_service_v2.CreateTenantV2Created)
	assert.NotNil(t, castedCreateTeant)
	assert.NotEmpty(t, castedCreateTeant.Payload.Data.ID)
	assert.NotEmpty(t, castedCreateTeant.Payload.Data.Attributes.Name)
	assert.NotEmpty(t, castedCreateTeant.Payload.Data.Attributes.URLSubdomain)
	assert.Equal(t, string(common.UserActive), *castedCreateTeant.Payload.Data.Attributes.State)
	assert.True(t, *castedCreateTeant.Payload.Data.Attributes.CreatedTimestamp > 0)
	assert.True(t, *castedCreateTeant.Payload.Data.Attributes.LastModifiedTimestamp > 0)

	// Make sure there are no existing Card
	existing := handlers.HandleGetAllCardsV2(handlers.AllRoles, tenantDB)(tenant_provisioning_service_v2.GetAllCardsV2Params{HTTPRequest: createHttpRequestWithParams(*castedCreateTeant.Payload.Data.ID, handlers.UserRoleSkylight, CardUrl, "GET")})
	castedResponse := existing.(*tenant_provisioning_service_v2.GetAllCardsV2NotFound)
	assert.NotNil(t, castedResponse)

	createReqBody := generateRandomTenantCardCreationRequest()
	created := handlers.HandleCreateCardV2(handlers.AllRoles, tenantDB)(tenant_provisioning_service_v2.CreateCardV2Params{Body: createReqBody, HTTPRequest: createHttpRequestWithParams(*castedCreateTeant.Payload.Data.ID, handlers.UserRoleSkylight, CardUrl, "POST")})
	castedCreate := created.(*tenant_provisioning_service_v2.CreateCardV2Created)
	assert.NotNil(t, castedCreate)
	assert.NotEmpty(t, castedCreate.Payload.Data.ID)
	assert.NotEmpty(t, castedCreate.Payload.Data.Attributes.Name)
	assert.NotEmpty(t, castedCreate.Payload.Data.Attributes.Description)
	assert.NotEmpty(t, castedCreate.Payload.Data.Attributes.Datatype)
	assert.True(t, *castedCreate.Payload.Data.Attributes.CreatedTimestamp > 0)
	assert.True(t, *castedCreate.Payload.Data.Attributes.LastModifiedTimestamp > 0)

	// Try to create the record again - should succeed as we are not guarding against name collisiones
	createdConflict := handlers.HandleCreateCardV2(handlers.AllRoles, tenantDB)(tenant_provisioning_service_v2.CreateCardV2Params{Body: createReqBody, HTTPRequest: createHttpRequestWithParams(*castedCreateTeant.Payload.Data.ID, handlers.UserRoleSkylight, CardUrl, "POST")})
	castedCreateConflictButOK := createdConflict.(*tenant_provisioning_service_v2.CreateCardV2Created)
	assert.NotNil(t, castedCreateConflictButOK)

	// Try the update with a bad revision
	newName := fake.CharactersN(16)
	updateRequestBody := generateCardUpdateRequest(*castedCreate.Payload.Data.ID, *castedCreate.Payload.Data.Attributes.Rev+"pork", &newName, nil)
	updated := handlers.HandleUpdateCardV2(handlers.AllRoles, tenantDB)(tenant_provisioning_service_v2.UpdateCardV2Params{CardID: *castedCreate.Payload.Data.ID, Body: updateRequestBody, HTTPRequest: createHttpRequestWithParams(*castedCreateTeant.Payload.Data.ID, handlers.UserRoleSkylight, CardUrl, "PATCH")})
	castedUpdate := updated.(*tenant_provisioning_service_v2.UpdateCardV2Conflict)
	assert.NotNil(t, castedUpdate)

	// Delete the tenant
	deleted := handlers.HandleDeleteCardV2(handlers.AllRoles, tenantDB)(tenant_provisioning_service_v2.DeleteCardV2Params{CardID: *castedCreate.Payload.Data.ID, HTTPRequest: createHttpRequestWithParams(*castedCreateTeant.Payload.Data.ID, handlers.UserRoleSkylight, CardUrl, "DELETE")})
	castedDelete := deleted.(*tenant_provisioning_service_v2.DeleteCardV2OK)
	assert.NotNil(t, castedDelete)
	assert.NotNil(t, castedDelete.Payload.Data)
}

func TestCardAPIsProtectedByAuthV2(t *testing.T) {
	fakeTenantID := fake.CharactersN(20)
	// Get All - All Users
	existing := handlers.HandleGetAllCardsV2(handlers.AllRoles, tenantDB)(tenant_provisioning_service_v2.GetAllCardsV2Params{HTTPRequest: createHttpRequestWithParams(fakeTenantID, handlers.UserRoleUnknown, CardUrl, "GET")})
	castedResponse := existing.(*tenant_provisioning_service_v2.GetAllCardsV2Forbidden)
	assert.NotNil(t, castedResponse)

	fakeID := fake.CharactersN(20)

	// Get - All Users
	fetched := handlers.HandleGetCardV2(handlers.AllRoles, tenantDB)(tenant_provisioning_service_v2.GetCardV2Params{CardID: fakeID, HTTPRequest: createHttpRequestWithParams(fakeTenantID, handlers.UserRoleUnknown, CardUrl, "GET")})
	castedFetch := fetched.(*tenant_provisioning_service_v2.GetCardV2Forbidden)
	assert.NotNil(t, castedFetch)

	// Create - SkylightAdmin and TenantAdmin Only
	created := handlers.HandleCreateCardV2(handlers.SkylightAndTenantAdminRoles, tenantDB)(tenant_provisioning_service_v2.CreateCardV2Params{Body: generateRandomTenantCardCreationRequest(), HTTPRequest: createHttpRequestWithParams(fakeTenantID, handlers.UserRoleUnknown, CardUrl, "POST")})
	castedCreate := created.(*tenant_provisioning_service_v2.CreateCardV2Forbidden)
	assert.NotNil(t, castedCreate)

	created = handlers.HandleCreateCardV2(handlers.SkylightAndTenantAdminRoles, tenantDB)(tenant_provisioning_service_v2.CreateCardV2Params{Body: generateRandomTenantCardCreationRequest(), HTTPRequest: createHttpRequestWithParams(fakeTenantID, handlers.UserRoleTenantUser, CardUrl, "POST")})
	castedCreate = created.(*tenant_provisioning_service_v2.CreateCardV2Forbidden)
	assert.NotNil(t, castedCreate)

	// Update - SkylightAdmin and TenantAdmin Only
	updated := handlers.HandleUpdateCardV2(handlers.SkylightAndTenantAdminRoles, tenantDB)(tenant_provisioning_service_v2.UpdateCardV2Params{CardID: fakeID, Body: nil, HTTPRequest: createHttpRequestWithParams(fakeTenantID, handlers.UserRoleUnknown, CardUrl, "PATCH")})
	castedUpdate := updated.(*tenant_provisioning_service_v2.UpdateCardV2Forbidden)
	assert.NotNil(t, castedUpdate)

	updated = handlers.HandleUpdateCardV2(handlers.SkylightAndTenantAdminRoles, tenantDB)(tenant_provisioning_service_v2.UpdateCardV2Params{CardID: fakeID, Body: nil, HTTPRequest: createHttpRequestWithParams(fakeTenantID, handlers.UserRoleTenantUser, CardUrl, "PATCH")})
	castedUpdate = updated.(*tenant_provisioning_service_v2.UpdateCardV2Forbidden)
	assert.NotNil(t, castedUpdate)

	// Delete - SkylightAdmin and TenantAdmin Only
	deleted := handlers.HandleDeleteCardV2(handlers.SkylightAndTenantAdminRoles, tenantDB)(tenant_provisioning_service_v2.DeleteCardV2Params{CardID: fakeID, HTTPRequest: createHttpRequestWithParams(fakeTenantID, handlers.UserRoleUnknown, CardUrl, "DELETE")})
	castedDelete := deleted.(*tenant_provisioning_service_v2.DeleteCardV2Forbidden)
	assert.NotNil(t, castedDelete)

	deleted = handlers.HandleDeleteCardV2(handlers.SkylightAndTenantAdminRoles, tenantDB)(tenant_provisioning_service_v2.DeleteCardV2Params{CardID: fakeID, HTTPRequest: createHttpRequestWithParams(fakeTenantID, handlers.UserRoleTenantUser, CardUrl, "DELETE")})
	castedDelete = deleted.(*tenant_provisioning_service_v2.DeleteCardV2Forbidden)
	assert.NotNil(t, castedDelete)
}

func generateRandomTenantCardCreationRequest() *swagmodels.CardCreateRequest {
	name := fake.CharactersN(8)
	return &swagmodels.CardCreateRequest{
		Data: &swagmodels.CardCreateRequestData{
			Type: &CardTypeString,
			Attributes: &swagmodels.CardCreateRequestDataAttributes{
				Description: fake.CharactersN(12),
				Name:        &name,
				Visualization: &swagmodels.CardVisualization{
					Category: fake.CharactersN(12),
					Label:    fake.CharactersN(13),
				},
			},
		},
	}
}

func generateCardUpdateRequest(id string, rev string, name *string, description *string) *swagmodels.CardUpdateRequest {
	result := &swagmodels.CardUpdateRequest{
		Data: &swagmodels.CardUpdateRequestData{
			Type:       &CardTypeString,
			ID:         &id,
			Attributes: &swagmodels.CardUpdateRequestDataAttributes{Rev: &rev},
		},
	}

	if name != nil {
		result.Data.Attributes.Name = *name
	}

	if description != nil {
		result.Data.Attributes.Description = *description
	}

	return result
}
