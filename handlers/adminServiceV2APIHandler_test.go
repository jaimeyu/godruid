package handlers_test

import (
	"testing"

	"github.com/icrowley/fake"

	"github.com/accedian/adh-gather/models/common"
	"github.com/accedian/adh-gather/swagmodels"

	"github.com/accedian/adh-gather/handlers"
	"github.com/accedian/adh-gather/restapi/operations/admin_provisioning_service_v2"
	"github.com/stretchr/testify/assert"
)

var (
	tenantURL               = "http://deployment.test.cool/api/v2/tenants"
	tenantIDByAliasURL      = "http://deployment.test.cool/api/v2/tenant-id-by-alias"
	tenantSummaryByAliasURL = "http://deployment.test.cool/api/v2/tenant-summary-by-alias"
	ingestionDictionaryURL  = "http://deployment.test.cool/api/v2/ingestion-dictionaries"
	validTypesURL           = "http://deployment.test.cool/api/v2/valid-types"

	tenantsTypeStr = "tenants"
)

func TestTeantCRUDV2(t *testing.T) {

	// Make sure there are no existing Tenants
	existing := handlers.HandleGetAllTenantsV2(handlers.SkylightAdminRoleOnly, adminDB)(admin_provisioning_service_v2.GetAllTenantsV2Params{HTTPRequest: createHttpRequestWithParams("", handlers.UserRoleSkylight, tenantURL, "GET")})
	castedResponse := existing.(*admin_provisioning_service_v2.GetAllTenantsV2NotFound)
	assert.NotNil(t, castedResponse)

	created := handlers.HandleCreateTenantV2(handlers.SkylightAdminRoleOnly, adminDB, tenantDB)(admin_provisioning_service_v2.CreateTenantV2Params{Body: generateRandomTenantCreationRequest(), HTTPRequest: createHttpRequestWithParams("", handlers.UserRoleSkylight, tenantURL, "POST")})
	castedCreate := created.(*admin_provisioning_service_v2.CreateTenantV2Created)
	assert.NotNil(t, castedCreate)
	assert.NotEmpty(t, castedCreate.Payload.Data.ID)
	assert.NotEmpty(t, castedCreate.Payload.Data.Attributes.Name)
	assert.NotEmpty(t, castedCreate.Payload.Data.Attributes.URLSubdomain)
	assert.Equal(t, string(common.UserActive), *castedCreate.Payload.Data.Attributes.State)
	assert.True(t, *castedCreate.Payload.Data.Attributes.CreatedTimestamp > 0)
	assert.True(t, *castedCreate.Payload.Data.Attributes.LastModifiedTimestamp > 0)

	// Make sure we can retrieve this record:
	fetched := handlers.HandleGetTenantV2(handlers.SkylightAdminRoleOnly, adminDB)(admin_provisioning_service_v2.GetTenantV2Params{TenantID: *castedCreate.Payload.Data.ID, HTTPRequest: createHttpRequestWithParams("", handlers.UserRoleSkylight, tenantURL, "GET")})
	castedFetch := fetched.(*admin_provisioning_service_v2.GetTenantV2OK)
	assert.NotNil(t, castedFetch)
	assert.Equal(t, castedCreate.Payload.Data, castedFetch.Payload.Data)

	// Also retrieve the record as part of an array
	fetchList := handlers.HandleGetAllTenantsV2(handlers.SkylightAdminRoleOnly, adminDB)(admin_provisioning_service_v2.GetAllTenantsV2Params{HTTPRequest: createHttpRequestWithParams("", handlers.UserRoleSkylight, tenantURL, "GET")})
	castedFetchList := fetchList.(*admin_provisioning_service_v2.GetAllTenantsV2OK)
	assert.NotNil(t, castedFetchList)
	assert.Equal(t, 1, len(castedFetchList.Payload.Data))
	assert.Equal(t, castedCreate.Payload.Data, castedFetchList.Payload.Data[0])

	// Make an update to the Record
	newName := fake.CharactersN(16)
	updateRequestBody := generateTenantUpdateRequest(*castedCreate.Payload.Data.ID, *castedCreate.Payload.Data.Attributes.Rev, &newName, nil, nil)
	updated := handlers.HandlePatchTenantV2(handlers.SkylightAdminRoleOnly, adminDB)(admin_provisioning_service_v2.PatchTenantV2Params{TenantID: *castedCreate.Payload.Data.ID, Body: updateRequestBody, HTTPRequest: createHttpRequestWithParams("", handlers.UserRoleSkylight, tenantURL, "PATCH")})
	castedUpdate := updated.(*admin_provisioning_service_v2.PatchTenantV2OK)
	assert.NotNil(t, castedUpdate)
	assert.NotEqual(t, castedCreate.Payload.Data, castedUpdate.Payload.Data)
	assert.NotEqual(t, castedCreate.Payload.Data.Attributes.Rev, castedUpdate.Payload.Data.Attributes.Rev)
	assert.Equal(t, newName, *castedUpdate.Payload.Data.Attributes.Name)
	assert.Equal(t, castedCreate.Payload.Data.Attributes.State, castedUpdate.Payload.Data.Attributes.State)
	assert.Equal(t, castedCreate.Payload.Data.Attributes.URLSubdomain, castedUpdate.Payload.Data.Attributes.URLSubdomain)

	// Delete the record
	deleted := handlers.HandleDeleteTenantV2(handlers.SkylightAdminRoleOnly, adminDB)(admin_provisioning_service_v2.DeleteTenantV2Params{TenantID: *castedCreate.Payload.Data.ID, HTTPRequest: createHttpRequestWithParams("", handlers.UserRoleSkylight, tenantURL, "DELETE")})
	castedDelete := deleted.(*admin_provisioning_service_v2.DeleteTenantV2OK)
	assert.NotNil(t, castedDelete)
	assert.Equal(t, castedUpdate.Payload.Data, castedDelete.Payload.Data)

	// Make sure there are no existing Tenants
	existing = handlers.HandleGetAllTenantsV2(handlers.SkylightAdminRoleOnly, adminDB)(admin_provisioning_service_v2.GetAllTenantsV2Params{HTTPRequest: createHttpRequestWithParams("", handlers.UserRoleSkylight, tenantURL, "GET")})
	castedResponse = existing.(*admin_provisioning_service_v2.GetAllTenantsV2NotFound)
	assert.NotNil(t, castedResponse)
}

func TestTenantNotFoundV2(t *testing.T) {
	notFoundID := fake.CharactersN(20)

	// Get Tenant
	fetched := handlers.HandleGetTenantV2(handlers.SkylightAdminRoleOnly, adminDB)(admin_provisioning_service_v2.GetTenantV2Params{TenantID: notFoundID, HTTPRequest: createHttpRequestWithParams("", handlers.UserRoleSkylight, tenantURL, "GET")})
	castedFetch := fetched.(*admin_provisioning_service_v2.GetTenantV2NotFound)
	assert.NotNil(t, castedFetch)

	// Delete Tenant
	deleted := handlers.HandleDeleteTenantV2(handlers.SkylightAdminRoleOnly, adminDB)(admin_provisioning_service_v2.DeleteTenantV2Params{TenantID: notFoundID, HTTPRequest: createHttpRequestWithParams("", handlers.UserRoleSkylight, tenantURL, "DELETE")})
	castedDelete := deleted.(*admin_provisioning_service_v2.DeleteTenantV2NotFound)
	assert.NotNil(t, castedDelete)

	// Patch tenant
	updateRequest := generateTenantUpdateRequest(notFoundID, "reviosionstuff", nil, nil, nil)
	updated := handlers.HandlePatchTenantV2(handlers.SkylightAdminRoleOnly, adminDB)(admin_provisioning_service_v2.PatchTenantV2Params{TenantID: notFoundID, Body: updateRequest, HTTPRequest: createHttpRequestWithParams("", handlers.UserRoleSkylight, tenantURL, "PATCH")})
	castedUpdate := updated.(*admin_provisioning_service_v2.PatchTenantV2NotFound)
	assert.NotNil(t, castedUpdate)
}

func TestTenantBadRequestV2(t *testing.T) {

	// CreateTenant
	created := handlers.HandleCreateTenantV2(handlers.SkylightAdminRoleOnly, adminDB, tenantDB)(admin_provisioning_service_v2.CreateTenantV2Params{Body: nil, HTTPRequest: createHttpRequestWithParams("", handlers.UserRoleSkylight, tenantURL, "POST")})
	castedCreate := created.(*admin_provisioning_service_v2.CreateTenantV2BadRequest)
	assert.NotNil(t, castedCreate)

	// Update Tenant
	updated := handlers.HandlePatchTenantV2(handlers.SkylightAdminRoleOnly, adminDB)(admin_provisioning_service_v2.PatchTenantV2Params{TenantID: fake.CharactersN(20), Body: nil, HTTPRequest: createHttpRequestWithParams("", handlers.UserRoleSkylight, tenantURL, "PATCH")})
	castedUpdate := updated.(*admin_provisioning_service_v2.PatchTenantV2BadRequest)
	assert.NotNil(t, castedUpdate)
}

func TestTenantConflictV2(t *testing.T) {

	// Make sure there are no existing Tenants
	existing := handlers.HandleGetAllTenantsV2(handlers.SkylightAdminRoleOnly, adminDB)(admin_provisioning_service_v2.GetAllTenantsV2Params{HTTPRequest: createHttpRequestWithParams("", handlers.UserRoleSkylight, tenantURL, "GET")})
	castedResponse := existing.(*admin_provisioning_service_v2.GetAllTenantsV2NotFound)
	assert.NotNil(t, castedResponse)

	createReqBody := generateRandomTenantCreationRequest()
	created := handlers.HandleCreateTenantV2(handlers.SkylightAdminRoleOnly, adminDB, tenantDB)(admin_provisioning_service_v2.CreateTenantV2Params{Body: createReqBody, HTTPRequest: createHttpRequestWithParams("", handlers.UserRoleSkylight, tenantURL, "POST")})
	castedCreate := created.(*admin_provisioning_service_v2.CreateTenantV2Created)
	assert.NotNil(t, castedCreate)
	assert.NotEmpty(t, castedCreate.Payload.Data.ID)
	assert.NotEmpty(t, castedCreate.Payload.Data.Attributes.Name)
	assert.NotEmpty(t, castedCreate.Payload.Data.Attributes.URLSubdomain)
	assert.Equal(t, string(common.UserActive), *castedCreate.Payload.Data.Attributes.State)
	assert.True(t, *castedCreate.Payload.Data.Attributes.CreatedTimestamp > 0)
	assert.True(t, *castedCreate.Payload.Data.Attributes.LastModifiedTimestamp > 0)

	// Try to create the record again
	createdConflict := handlers.HandleCreateTenantV2(handlers.SkylightAdminRoleOnly, adminDB, tenantDB)(admin_provisioning_service_v2.CreateTenantV2Params{Body: createReqBody, HTTPRequest: createHttpRequestWithParams("", handlers.UserRoleSkylight, tenantURL, "POST")})
	castedCreateConflict := createdConflict.(*admin_provisioning_service_v2.CreateTenantV2Conflict)
	assert.NotNil(t, castedCreateConflict)

	// Try the update with a bad revision
	newName := fake.CharactersN(16)
	updateRequestBody := generateTenantUpdateRequest(*castedCreate.Payload.Data.ID, *castedCreate.Payload.Data.Attributes.Rev+"pork", &newName, nil, nil)
	updated := handlers.HandlePatchTenantV2(handlers.SkylightAdminRoleOnly, adminDB)(admin_provisioning_service_v2.PatchTenantV2Params{TenantID: *castedCreate.Payload.Data.ID, Body: updateRequestBody, HTTPRequest: createHttpRequestWithParams("", handlers.UserRoleSkylight, tenantURL, "PATCH")})
	castedUpdate := updated.(*admin_provisioning_service_v2.PatchTenantV2Conflict)
	assert.NotNil(t, castedUpdate)

	// Delete the tenant
	deleted := handlers.HandleDeleteTenantV2(handlers.SkylightAdminRoleOnly, adminDB)(admin_provisioning_service_v2.DeleteTenantV2Params{TenantID: *castedCreate.Payload.Data.ID, HTTPRequest: createHttpRequestWithParams("", handlers.UserRoleSkylight, tenantURL, "DELETE")})
	castedDelete := deleted.(*admin_provisioning_service_v2.DeleteTenantV2OK)
	assert.NotNil(t, castedDelete)
	assert.NotNil(t, castedDelete.Payload.Data)
}

func TestAdminServiceAPIsProtectedByAuthV2(t *testing.T) {
	// Get All Tenants - Skylight Admin only
	existing := handlers.HandleGetAllTenantsV2(handlers.SkylightAdminRoleOnly, adminDB)(admin_provisioning_service_v2.GetAllTenantsV2Params{HTTPRequest: createHttpRequestWithParams("", handlers.UserRoleTenantAdmin, tenantURL, "GET")})
	castedResponse := existing.(*admin_provisioning_service_v2.GetAllTenantsV2Forbidden)
	assert.NotNil(t, castedResponse)

	existing = handlers.HandleGetAllTenantsV2(handlers.SkylightAdminRoleOnly, adminDB)(admin_provisioning_service_v2.GetAllTenantsV2Params{HTTPRequest: createHttpRequestWithParams("", handlers.UserRoleTenantUser, tenantURL, "GET")})
	castedResponse = existing.(*admin_provisioning_service_v2.GetAllTenantsV2Forbidden)
	assert.NotNil(t, castedResponse)

	fakeID := fake.CharactersN(20)

	// Get Tenant - Skylight Admin only
	fetched := handlers.HandleGetTenantV2(handlers.SkylightAdminRoleOnly, adminDB)(admin_provisioning_service_v2.GetTenantV2Params{TenantID: fakeID, HTTPRequest: createHttpRequestWithParams("", handlers.UserRoleTenantAdmin, tenantURL, "GET")})
	castedFetch := fetched.(*admin_provisioning_service_v2.GetTenantV2Forbidden)
	assert.NotNil(t, castedFetch)

	fetched = handlers.HandleGetTenantV2(handlers.SkylightAdminRoleOnly, adminDB)(admin_provisioning_service_v2.GetTenantV2Params{TenantID: fakeID, HTTPRequest: createHttpRequestWithParams("", handlers.UserRoleTenantUser, tenantURL, "GET")})
	castedFetch = fetched.(*admin_provisioning_service_v2.GetTenantV2Forbidden)
	assert.NotNil(t, castedFetch)

	// Create Tenant - Skylight Admin Only
	created := handlers.HandleCreateTenantV2(handlers.SkylightAdminRoleOnly, adminDB, tenantDB)(admin_provisioning_service_v2.CreateTenantV2Params{Body: generateRandomTenantCreationRequest(), HTTPRequest: createHttpRequestWithParams("", handlers.UserRoleTenantAdmin, tenantURL, "POST")})
	castedCreate := created.(*admin_provisioning_service_v2.CreateTenantV2Forbidden)
	assert.NotNil(t, castedCreate)

	created = handlers.HandleCreateTenantV2(handlers.SkylightAdminRoleOnly, adminDB, tenantDB)(admin_provisioning_service_v2.CreateTenantV2Params{Body: generateRandomTenantCreationRequest(), HTTPRequest: createHttpRequestWithParams("", handlers.UserRoleTenantUser, tenantURL, "POST")})
	castedCreate = created.(*admin_provisioning_service_v2.CreateTenantV2Forbidden)
	assert.NotNil(t, castedCreate)

	// Update Tenant - Skylight Admin only
	updated := handlers.HandlePatchTenantV2(handlers.SkylightAdminRoleOnly, adminDB)(admin_provisioning_service_v2.PatchTenantV2Params{TenantID: fakeID, Body: nil, HTTPRequest: createHttpRequestWithParams("", handlers.UserRoleTenantAdmin, tenantURL, "PATCH")})
	castedUpdate := updated.(*admin_provisioning_service_v2.PatchTenantV2Forbidden)
	assert.NotNil(t, castedUpdate)

	updated = handlers.HandlePatchTenantV2(handlers.SkylightAdminRoleOnly, adminDB)(admin_provisioning_service_v2.PatchTenantV2Params{TenantID: fakeID, Body: nil, HTTPRequest: createHttpRequestWithParams("", handlers.UserRoleTenantUser, tenantURL, "PATCH")})
	castedUpdate = updated.(*admin_provisioning_service_v2.PatchTenantV2Forbidden)
	assert.NotNil(t, castedUpdate)

	// Delete record - Skylight Admin only
	deleted := handlers.HandleDeleteTenantV2(handlers.SkylightAdminRoleOnly, adminDB)(admin_provisioning_service_v2.DeleteTenantV2Params{TenantID: fakeID, HTTPRequest: createHttpRequestWithParams("", handlers.UserRoleTenantAdmin, tenantURL, "DELETE")})
	castedDelete := deleted.(*admin_provisioning_service_v2.DeleteTenantV2Forbidden)
	assert.NotNil(t, castedDelete)

	deleted = handlers.HandleDeleteTenantV2(handlers.SkylightAdminRoleOnly, adminDB)(admin_provisioning_service_v2.DeleteTenantV2Params{TenantID: fakeID, HTTPRequest: createHttpRequestWithParams("", handlers.UserRoleTenantUser, tenantURL, "DELETE")})
	castedDelete = deleted.(*admin_provisioning_service_v2.DeleteTenantV2Forbidden)
	assert.NotNil(t, castedDelete)

	// Get ID and Summary by Alias both unprotected
	idFromAlais := handlers.HandleGetTenantIDByAliasV2(adminDB)(admin_provisioning_service_v2.GetTenantIDByAliasV2Params{Value: fake.CharactersN(6), HTTPRequest: createHttpRequestWithParams("", handlers.UserRoleUnknown, tenantIDByAliasURL, "GET")})
	castedIdFromAlias := idFromAlais.(*admin_provisioning_service_v2.GetTenantIDByAliasV2NotFound)
	assert.NotNil(t, castedIdFromAlias)

	summaryFromAlais := handlers.HandleGetTenantSummaryByAliasV2(adminDB)(admin_provisioning_service_v2.GetTenantSummaryByAliasV2Params{Value: fake.CharactersN(6), HTTPRequest: createHttpRequestWithParams("", handlers.UserRoleUnknown, tenantSummaryByAliasURL, "GET")})
	castedSummaryFromAlias := summaryFromAlais.(*admin_provisioning_service_v2.GetTenantSummaryByAliasV2NotFound)
	assert.NotNil(t, castedSummaryFromAlias)

	// Get IngestionDictionary and GetValidTypers should be available to all roles
	fetchedIngDict := handlers.HandleGetIngestionDictionaryV2(handlers.AllRoles, adminDB)(admin_provisioning_service_v2.GetIngestionDictionaryV2Params{HTTPRequest: createHttpRequestWithParams("", handlers.UserRoleUnknown, ingestionDictionaryURL, "GET")})
	castedFetchIngDict := fetchedIngDict.(*admin_provisioning_service_v2.GetIngestionDictionaryV2Forbidden)
	assert.NotNil(t, castedFetchIngDict)

	fetchedVT := handlers.HandleGetValidTypesV2(handlers.AllRoles, adminDB)(admin_provisioning_service_v2.GetValidTypesV2Params{HTTPRequest: createHttpRequestWithParams("", handlers.UserRoleUnknown, validTypesURL, "GET")})
	castedFetchVT := fetchedVT.(*admin_provisioning_service_v2.GetValidTypesV2Forbidden)
	assert.NotNil(t, castedFetchVT)
}

func TestGetValidTypesV2(t *testing.T) {
	fetched := handlers.HandleGetValidTypesV2(handlers.AllRoles, adminDB)(admin_provisioning_service_v2.GetValidTypesV2Params{HTTPRequest: createHttpRequestWithParams("", handlers.UserRoleTenantUser, validTypesURL, "GET")})
	castedFetchVT := fetched.(*admin_provisioning_service_v2.GetValidTypesV2OK)
	assert.NotEmpty(t, castedFetchVT.Payload.Data[0])
	assert.NotEmpty(t, castedFetchVT.Payload.Data[0].Attributes.MonitoredObjectDeviceTypes)
	assert.NotEmpty(t, castedFetchVT.Payload.Data[0].Attributes.MonitoredObjectTypes)
}

func TestGetByAliasV2(t *testing.T) {
	// Make sure there are no existing Tenants
	existing := handlers.HandleGetAllTenantsV2(handlers.SkylightAdminRoleOnly, adminDB)(admin_provisioning_service_v2.GetAllTenantsV2Params{HTTPRequest: createHttpRequestWithParams("", handlers.UserRoleSkylight, tenantURL, "GET")})
	castedResponse := existing.(*admin_provisioning_service_v2.GetAllTenantsV2NotFound)
	assert.NotNil(t, castedResponse)

	createReqBody := generateRandomTenantCreationRequest()
	created := handlers.HandleCreateTenantV2(handlers.SkylightAdminRoleOnly, adminDB, tenantDB)(admin_provisioning_service_v2.CreateTenantV2Params{Body: createReqBody, HTTPRequest: createHttpRequestWithParams("", handlers.UserRoleSkylight, tenantURL, "POST")})
	castedCreate := created.(*admin_provisioning_service_v2.CreateTenantV2Created)
	assert.NotNil(t, castedCreate)
	assert.NotEmpty(t, castedCreate.Payload.Data.ID)
	assert.NotEmpty(t, castedCreate.Payload.Data.Attributes.Name)
	assert.NotEmpty(t, castedCreate.Payload.Data.Attributes.URLSubdomain)
	assert.Equal(t, string(common.UserActive), *castedCreate.Payload.Data.Attributes.State)
	assert.True(t, *castedCreate.Payload.Data.Attributes.CreatedTimestamp > 0)
	assert.True(t, *castedCreate.Payload.Data.Attributes.LastModifiedTimestamp > 0)

	idByAlias := handlers.HandleGetTenantIDByAliasV2(adminDB)(admin_provisioning_service_v2.GetTenantIDByAliasV2Params{Value: *castedCreate.Payload.Data.Attributes.Name, HTTPRequest: createHttpRequestWithParams("", handlers.UserRoleTenantUser, tenantIDByAliasURL, "GET")})
	castedIdByAlias := idByAlias.(*admin_provisioning_service_v2.GetTenantIDByAliasV2OK)
	assert.Equal(t, *castedCreate.Payload.Data.ID, castedIdByAlias.Payload)

	summaryByAlias := handlers.HandleGetTenantSummaryByAliasV2(adminDB)(admin_provisioning_service_v2.GetTenantSummaryByAliasV2Params{Value: *castedCreate.Payload.Data.Attributes.Name, HTTPRequest: createHttpRequestWithParams("", handlers.UserRoleTenantUser, tenantSummaryByAliasURL, "GET")})
	castedSummaryByAlias := summaryByAlias.(*admin_provisioning_service_v2.GetTenantSummaryByAliasV2OK)
	assert.Equal(t, *castedCreate.Payload.Data.ID, castedSummaryByAlias.Payload.Data.Attributes.ID)
}

func TestGetIngestionDictionaryV2(t *testing.T) {
	ingDict := handlers.HandleGetIngestionDictionaryV2(handlers.AllRoles, adminDB)(admin_provisioning_service_v2.GetIngestionDictionaryV2Params{HTTPRequest: createHttpRequestWithParams("", handlers.UserRoleTenantUser, ingestionDictionaryURL, "GET")})
	castedIngDict := ingDict.(*admin_provisioning_service_v2.GetIngestionDictionaryV2OK)
	assert.NotEmpty(t, castedIngDict.Payload.Data)
	assert.NotEmpty(t, castedIngDict.Payload.Data[0].Attributes.Metrics)
}

func generateRandomTenantCreationRequest() *swagmodels.TenantCreationRequest {
	name := fake.CharactersN(12)
	domain := fake.DomainName()
	state := string(common.UserActive)
	return &swagmodels.TenantCreationRequest{
		Data: &swagmodels.TenantCreationObject{
			Type: &tenantsTypeStr,
			Attributes: &swagmodels.TenantCreationObjectAttributes{
				Name:         &name,
				URLSubdomain: &domain,
				State:        &state,
			},
		},
	}
}

func generateTenantUpdateRequest(id string, rev string, name *string, domain *string, state *string) *swagmodels.TenantUpdateRequest {
	result := &swagmodels.TenantUpdateRequest{
		Data: &swagmodels.TenantUpdateObject{
			Type:       &tenantsTypeStr,
			ID:         &id,
			Attributes: &swagmodels.TenantUpdateObjectAttributes{Rev: &rev},
		},
	}

	if name != nil {
		result.Data.Attributes.Name = *name
	}
	if domain != nil {
		result.Data.Attributes.URLSubdomain = *domain
	}
	if state != nil {
		result.Data.Attributes.State = *state
	}
	return result
}
