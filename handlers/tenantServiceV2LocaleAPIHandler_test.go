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
	LocaleUrl = "http://deployment.test.cool/api/v2/locales"

	LocaleTypeString = "locales"
)

func TestLocaleCRUDV2(t *testing.T) {

	createdTenant := handlers.HandleCreateTenantV2(handlers.SkylightAdminRoleOnly, adminDB, tenantDB)(admin_provisioning_service_v2.CreateTenantV2Params{Body: generateRandomTenantCreationRequest(), HTTPRequest: createHttpRequestWithParams("", handlers.UserRoleSkylight, tenantURL, "POST")})
	castedCreateTeant := createdTenant.(*admin_provisioning_service_v2.CreateTenantV2Created)
	assert.NotNil(t, castedCreateTeant)
	assert.NotEmpty(t, castedCreateTeant.Payload.Data.ID)
	assert.NotEmpty(t, castedCreateTeant.Payload.Data.Attributes.Name)
	assert.NotEmpty(t, castedCreateTeant.Payload.Data.Attributes.URLSubdomain)
	assert.Equal(t, string(common.UserActive), *castedCreateTeant.Payload.Data.Attributes.State)
	assert.True(t, *castedCreateTeant.Payload.Data.Attributes.CreatedTimestamp > 0)
	assert.True(t, *castedCreateTeant.Payload.Data.Attributes.LastModifiedTimestamp > 0)

	// Make sure there are no existing Locales
	existing := handlers.HandleGetAllLocalesV2(handlers.AllRoles, tenantDB)(tenant_provisioning_service_v2.GetAllLocalesV2Params{HTTPRequest: createHttpRequestWithParams(*castedCreateTeant.Payload.Data.ID, handlers.UserRoleSkylight, LocaleUrl, "GET")})
	castedResponse := existing.(*tenant_provisioning_service_v2.GetAllLocalesV2NotFound)
	assert.NotNil(t, castedResponse)

	created := handlers.HandleCreateLocaleV2(handlers.AllRoles, tenantDB)(tenant_provisioning_service_v2.CreateLocaleV2Params{Body: generateRandomTenantLocaleCreationRequest(), HTTPRequest: createHttpRequestWithParams(*castedCreateTeant.Payload.Data.ID, handlers.UserRoleSkylight, LocaleUrl, "POST")})
	castedCreate := created.(*tenant_provisioning_service_v2.CreateLocaleV2Created)
	assert.NotNil(t, castedCreate)
	assert.NotEmpty(t, castedCreate.Payload.Data.ID)
	assert.NotEmpty(t, castedCreate.Payload.Data.Attributes.Moment)
	assert.NotEmpty(t, castedCreate.Payload.Data.Attributes.Intl)
	assert.NotEmpty(t, castedCreate.Payload.Data.Attributes.Timezone)
	assert.NotEmpty(t, castedCreate.Payload.Data.Attributes.Datatype)
	assert.True(t, *castedCreate.Payload.Data.Attributes.CreatedTimestamp > 0)
	assert.True(t, *castedCreate.Payload.Data.Attributes.LastModifiedTimestamp > 0)

	// Make sure we can retrieve this record:
	fetched := handlers.HandleGetLocaleV2(handlers.AllRoles, tenantDB)(tenant_provisioning_service_v2.GetLocaleV2Params{LocaleID: *castedCreate.Payload.Data.ID, HTTPRequest: createHttpRequestWithParams(*castedCreateTeant.Payload.Data.ID, handlers.UserRoleSkylight, LocaleUrl, "GET")})
	castedFetch := fetched.(*tenant_provisioning_service_v2.GetLocaleV2OK)
	assert.NotNil(t, castedFetch)
	assert.Equal(t, castedCreate.Payload.Data, castedFetch.Payload.Data)

	// Also retrieve the record as part of an array
	fetchList := handlers.HandleGetAllLocalesV2(handlers.AllRoles, tenantDB)(tenant_provisioning_service_v2.GetAllLocalesV2Params{HTTPRequest: createHttpRequestWithParams(*castedCreateTeant.Payload.Data.ID, handlers.UserRoleSkylight, LocaleUrl, "GET")})
	castedFetchList := fetchList.(*tenant_provisioning_service_v2.GetAllLocalesV2OK)
	assert.NotNil(t, castedFetchList)
	assert.Equal(t, 1, len(castedFetchList.Payload.Data))
	assert.Equal(t, castedCreate.Payload.Data, castedFetchList.Payload.Data[0])

	// Make an update to the Record
	newName := fake.CharactersN(16)
	updateRequestBody := generateLocaleUpdateRequest(*castedCreate.Payload.Data.ID, *castedCreate.Payload.Data.Attributes.Rev, &newName, nil)
	updated := handlers.HandleUpdateLocaleV2(handlers.AllRoles, tenantDB)(tenant_provisioning_service_v2.UpdateLocaleV2Params{LocaleID: *castedCreate.Payload.Data.ID, Body: updateRequestBody, HTTPRequest: createHttpRequestWithParams(*castedCreateTeant.Payload.Data.ID, handlers.UserRoleSkylight, LocaleUrl, "PATCH")})
	castedUpdate := updated.(*tenant_provisioning_service_v2.UpdateLocaleV2OK)
	assert.NotNil(t, castedUpdate)
	assert.NotEqual(t, castedCreate.Payload.Data, castedUpdate.Payload.Data)
	assert.NotEqual(t, castedCreate.Payload.Data.Attributes.Rev, castedUpdate.Payload.Data.Attributes.Rev)
	assert.Equal(t, newName, *castedUpdate.Payload.Data.Attributes.Intl)
	assert.Equal(t, castedCreate.Payload.Data.Attributes.Moment, castedUpdate.Payload.Data.Attributes.Moment)

	// Delete the record
	deleted := handlers.HandleDeleteLocaleV2(handlers.AllRoles, tenantDB)(tenant_provisioning_service_v2.DeleteLocaleV2Params{LocaleID: *castedCreate.Payload.Data.ID, HTTPRequest: createHttpRequestWithParams(*castedCreateTeant.Payload.Data.ID, handlers.UserRoleSkylight, LocaleUrl, "DELETE")})
	castedDelete := deleted.(*tenant_provisioning_service_v2.DeleteLocaleV2OK)
	assert.NotNil(t, castedDelete)
	assert.Equal(t, castedUpdate.Payload.Data, castedDelete.Payload.Data)

	// Make sure there are no existing Locales
	existing = handlers.HandleGetAllLocalesV2(handlers.AllRoles, tenantDB)(tenant_provisioning_service_v2.GetAllLocalesV2Params{HTTPRequest: createHttpRequestWithParams(*castedCreateTeant.Payload.Data.ID, handlers.UserRoleSkylight, LocaleUrl, "GET")})
	castedResponse = existing.(*tenant_provisioning_service_v2.GetAllLocalesV2NotFound)
	assert.NotNil(t, castedResponse)
}

func TestLocaleNotFoundV2(t *testing.T) {
	notFoundID := fake.CharactersN(20)

	// Get Locale
	fetched := handlers.HandleGetLocaleV2(handlers.AllRoles, tenantDB)(tenant_provisioning_service_v2.GetLocaleV2Params{LocaleID: notFoundID, HTTPRequest: createHttpRequestWithParams("", handlers.UserRoleSkylight, LocaleUrl, "GET")})
	castedFetch := fetched.(*tenant_provisioning_service_v2.GetLocaleV2NotFound)
	assert.NotNil(t, castedFetch)

	// Delete Locale
	deleted := handlers.HandleDeleteLocaleV2(handlers.AllRoles, tenantDB)(tenant_provisioning_service_v2.DeleteLocaleV2Params{LocaleID: notFoundID, HTTPRequest: createHttpRequestWithParams("", handlers.UserRoleSkylight, LocaleUrl, "DELETE")})
	castedDelete := deleted.(*tenant_provisioning_service_v2.DeleteLocaleV2NotFound)
	assert.NotNil(t, castedDelete)

	// Patch Locale
	updateRequest := generateLocaleUpdateRequest(notFoundID, "reviosionstuff", nil, nil)
	updated := handlers.HandleUpdateLocaleV2(handlers.AllRoles, tenantDB)(tenant_provisioning_service_v2.UpdateLocaleV2Params{LocaleID: notFoundID, Body: updateRequest, HTTPRequest: createHttpRequestWithParams("", handlers.UserRoleSkylight, LocaleUrl, "PATCH")})
	castedUpdate := updated.(*tenant_provisioning_service_v2.UpdateLocaleV2NotFound)
	assert.NotNil(t, castedUpdate)
}

func TestLocaleBadRequestV2(t *testing.T) {

	// CreateLocale
	created := handlers.HandleCreateLocaleV2(handlers.AllRoles, tenantDB)(tenant_provisioning_service_v2.CreateLocaleV2Params{Body: nil, HTTPRequest: createHttpRequestWithParams("", handlers.UserRoleSkylight, LocaleUrl, "POST")})
	castedCreate := created.(*tenant_provisioning_service_v2.CreateLocaleV2BadRequest)
	assert.NotNil(t, castedCreate)

	// Update Locale
	updated := handlers.HandleUpdateLocaleV2(handlers.AllRoles, tenantDB)(tenant_provisioning_service_v2.UpdateLocaleV2Params{LocaleID: fake.CharactersN(20), Body: nil, HTTPRequest: createHttpRequestWithParams("", handlers.UserRoleSkylight, LocaleUrl, "PATCH")})
	castedUpdate := updated.(*tenant_provisioning_service_v2.UpdateLocaleV2BadRequest)
	assert.NotNil(t, castedUpdate)
}

func TestLocaleConflictV2(t *testing.T) {

	createdTenant := handlers.HandleCreateTenantV2(handlers.SkylightAdminRoleOnly, adminDB, tenantDB)(admin_provisioning_service_v2.CreateTenantV2Params{Body: generateRandomTenantCreationRequest(), HTTPRequest: createHttpRequestWithParams("", handlers.UserRoleSkylight, tenantURL, "POST")})
	castedCreateTeant := createdTenant.(*admin_provisioning_service_v2.CreateTenantV2Created)
	assert.NotNil(t, castedCreateTeant)
	assert.NotEmpty(t, castedCreateTeant.Payload.Data.ID)
	assert.NotEmpty(t, castedCreateTeant.Payload.Data.Attributes.Name)
	assert.NotEmpty(t, castedCreateTeant.Payload.Data.Attributes.URLSubdomain)
	assert.Equal(t, string(common.UserActive), *castedCreateTeant.Payload.Data.Attributes.State)
	assert.True(t, *castedCreateTeant.Payload.Data.Attributes.CreatedTimestamp > 0)
	assert.True(t, *castedCreateTeant.Payload.Data.Attributes.LastModifiedTimestamp > 0)

	// Make sure there are no existing Locale
	existing := handlers.HandleGetAllLocalesV2(handlers.AllRoles, tenantDB)(tenant_provisioning_service_v2.GetAllLocalesV2Params{HTTPRequest: createHttpRequestWithParams(*castedCreateTeant.Payload.Data.ID, handlers.UserRoleSkylight, LocaleUrl, "GET")})
	castedResponse := existing.(*tenant_provisioning_service_v2.GetAllLocalesV2NotFound)
	assert.NotNil(t, castedResponse)

	createReqBody := generateRandomTenantLocaleCreationRequest()
	created := handlers.HandleCreateLocaleV2(handlers.AllRoles, tenantDB)(tenant_provisioning_service_v2.CreateLocaleV2Params{Body: createReqBody, HTTPRequest: createHttpRequestWithParams(*castedCreateTeant.Payload.Data.ID, handlers.UserRoleSkylight, LocaleUrl, "POST")})
	castedCreate := created.(*tenant_provisioning_service_v2.CreateLocaleV2Created)
	assert.NotNil(t, castedCreate)
	assert.NotEmpty(t, castedCreate.Payload.Data.ID)
	assert.NotEmpty(t, castedCreate.Payload.Data.Attributes.Moment)
	assert.NotEmpty(t, castedCreate.Payload.Data.Attributes.Intl)
	assert.NotEmpty(t, castedCreate.Payload.Data.Attributes.Timezone)
	assert.NotEmpty(t, castedCreate.Payload.Data.Attributes.Datatype)
	assert.True(t, *castedCreate.Payload.Data.Attributes.CreatedTimestamp > 0)
	assert.True(t, *castedCreate.Payload.Data.Attributes.LastModifiedTimestamp > 0)

	// Try to create the record again - should succeed as we are not guarding against name collisiones
	createdConflict := handlers.HandleCreateLocaleV2(handlers.AllRoles, tenantDB)(tenant_provisioning_service_v2.CreateLocaleV2Params{Body: createReqBody, HTTPRequest: createHttpRequestWithParams(*castedCreateTeant.Payload.Data.ID, handlers.UserRoleSkylight, LocaleUrl, "POST")})
	castedCreateConflictButOK := createdConflict.(*tenant_provisioning_service_v2.CreateLocaleV2Created)
	assert.NotNil(t, castedCreateConflictButOK)

	// Try the update with a bad revision
	newName := fake.CharactersN(16)
	updateRequestBody := generateLocaleUpdateRequest(*castedCreate.Payload.Data.ID, *castedCreate.Payload.Data.Attributes.Rev+"pork", &newName, nil)
	updated := handlers.HandleUpdateLocaleV2(handlers.AllRoles, tenantDB)(tenant_provisioning_service_v2.UpdateLocaleV2Params{LocaleID: *castedCreate.Payload.Data.ID, Body: updateRequestBody, HTTPRequest: createHttpRequestWithParams(*castedCreateTeant.Payload.Data.ID, handlers.UserRoleSkylight, LocaleUrl, "PATCH")})
	castedUpdate := updated.(*tenant_provisioning_service_v2.UpdateLocaleV2Conflict)
	assert.NotNil(t, castedUpdate)

	// Delete the tenant
	deleted := handlers.HandleDeleteLocaleV2(handlers.AllRoles, tenantDB)(tenant_provisioning_service_v2.DeleteLocaleV2Params{LocaleID: *castedCreate.Payload.Data.ID, HTTPRequest: createHttpRequestWithParams(*castedCreateTeant.Payload.Data.ID, handlers.UserRoleSkylight, LocaleUrl, "DELETE")})
	castedDelete := deleted.(*tenant_provisioning_service_v2.DeleteLocaleV2OK)
	assert.NotNil(t, castedDelete)
	assert.NotNil(t, castedDelete.Payload.Data)
}

func TestLocaleAPIsProtectedByAuthV2(t *testing.T) {
	fakeTenantID := fake.CharactersN(20)
	// Get All - All Users
	existing := handlers.HandleGetAllLocalesV2(handlers.AllRoles, tenantDB)(tenant_provisioning_service_v2.GetAllLocalesV2Params{HTTPRequest: createHttpRequestWithParams(fakeTenantID, handlers.UserRoleUnknown, LocaleUrl, "GET")})
	castedResponse := existing.(*tenant_provisioning_service_v2.GetAllLocalesV2Forbidden)
	assert.NotNil(t, castedResponse)

	fakeID := fake.CharactersN(20)

	// Get - All Users
	fetched := handlers.HandleGetLocaleV2(handlers.AllRoles, tenantDB)(tenant_provisioning_service_v2.GetLocaleV2Params{LocaleID: fakeID, HTTPRequest: createHttpRequestWithParams(fakeTenantID, handlers.UserRoleUnknown, LocaleUrl, "GET")})
	castedFetch := fetched.(*tenant_provisioning_service_v2.GetLocaleV2Forbidden)
	assert.NotNil(t, castedFetch)

	// Create - SkylightAdmin and TenantAdmin Only
	created := handlers.HandleCreateLocaleV2(handlers.SkylightAndTenantAdminRoles, tenantDB)(tenant_provisioning_service_v2.CreateLocaleV2Params{Body: generateRandomTenantLocaleCreationRequest(), HTTPRequest: createHttpRequestWithParams(fakeTenantID, handlers.UserRoleUnknown, LocaleUrl, "POST")})
	castedCreate := created.(*tenant_provisioning_service_v2.CreateLocaleV2Forbidden)
	assert.NotNil(t, castedCreate)

	created = handlers.HandleCreateLocaleV2(handlers.SkylightAndTenantAdminRoles, tenantDB)(tenant_provisioning_service_v2.CreateLocaleV2Params{Body: generateRandomTenantLocaleCreationRequest(), HTTPRequest: createHttpRequestWithParams(fakeTenantID, handlers.UserRoleTenantUser, LocaleUrl, "POST")})
	castedCreate = created.(*tenant_provisioning_service_v2.CreateLocaleV2Forbidden)
	assert.NotNil(t, castedCreate)

	// Update - SkylightAdmin and TenantAdmin Only
	updated := handlers.HandleUpdateLocaleV2(handlers.SkylightAndTenantAdminRoles, tenantDB)(tenant_provisioning_service_v2.UpdateLocaleV2Params{LocaleID: fakeID, Body: nil, HTTPRequest: createHttpRequestWithParams(fakeTenantID, handlers.UserRoleUnknown, LocaleUrl, "PATCH")})
	castedUpdate := updated.(*tenant_provisioning_service_v2.UpdateLocaleV2Forbidden)
	assert.NotNil(t, castedUpdate)

	updated = handlers.HandleUpdateLocaleV2(handlers.SkylightAndTenantAdminRoles, tenantDB)(tenant_provisioning_service_v2.UpdateLocaleV2Params{LocaleID: fakeID, Body: nil, HTTPRequest: createHttpRequestWithParams(fakeTenantID, handlers.UserRoleTenantUser, LocaleUrl, "PATCH")})
	castedUpdate = updated.(*tenant_provisioning_service_v2.UpdateLocaleV2Forbidden)
	assert.NotNil(t, castedUpdate)

	// Delete - SkylightAdmin and TenantAdmin Only
	deleted := handlers.HandleDeleteLocaleV2(handlers.SkylightAndTenantAdminRoles, tenantDB)(tenant_provisioning_service_v2.DeleteLocaleV2Params{LocaleID: fakeID, HTTPRequest: createHttpRequestWithParams(fakeTenantID, handlers.UserRoleUnknown, LocaleUrl, "DELETE")})
	castedDelete := deleted.(*tenant_provisioning_service_v2.DeleteLocaleV2Forbidden)
	assert.NotNil(t, castedDelete)

	deleted = handlers.HandleDeleteLocaleV2(handlers.SkylightAndTenantAdminRoles, tenantDB)(tenant_provisioning_service_v2.DeleteLocaleV2Params{LocaleID: fakeID, HTTPRequest: createHttpRequestWithParams(fakeTenantID, handlers.UserRoleTenantUser, LocaleUrl, "DELETE")})
	castedDelete = deleted.(*tenant_provisioning_service_v2.DeleteLocaleV2Forbidden)
	assert.NotNil(t, castedDelete)
}

func generateRandomTenantLocaleCreationRequest() *swagmodels.LocaleCreateRequest {
	intl := fake.CharactersN(6)
	moment := fake.CharactersN(2)
	tz := fake.CharactersN(15)
	return &swagmodels.LocaleCreateRequest{
		Data: &swagmodels.LocaleCreateRequestData{
			Type: &LocaleTypeString,
			Attributes: &swagmodels.LocaleCreateRequestDataAttributes{
				Intl:     &intl,
				Moment:   &moment,
				Timezone: &tz,
			},
		},
	}
}

func generateLocaleUpdateRequest(id string, rev string, intl *string, moment *string) *swagmodels.LocaleUpdateRequest {
	result := &swagmodels.LocaleUpdateRequest{
		Data: &swagmodels.LocaleUpdateRequestData{
			Type:       &LocaleTypeString,
			ID:         &id,
			Attributes: &swagmodels.LocaleUpdateRequestDataAttributes{Rev: &rev},
		},
	}

	if intl != nil {
		result.Data.Attributes.Intl = *intl
	}

	if moment != nil {
		result.Data.Attributes.Moment = *moment
	}

	return result
}
