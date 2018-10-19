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
	MetadataConfigUrl = "http://deployment.test.cool/api/v2/MetadataConfigs"

	MetadataConfigTypeString = "metadataConfigs"
)

func TestMetadataConfigCRUDV2(t *testing.T) {

	createdTenant := handlers.HandleCreateTenantV2(handlers.SkylightAdminRoleOnly, adminDB, tenantDB)(admin_provisioning_service_v2.CreateTenantV2Params{Body: generateRandomTenantCreationRequest(), HTTPRequest: createHttpRequestWithParams("", handlers.UserRoleSkylight, tenantURL, "POST")})
	castedCreateTeant := createdTenant.(*admin_provisioning_service_v2.CreateTenantV2Created)
	assert.NotNil(t, castedCreateTeant)
	assert.NotEmpty(t, castedCreateTeant.Payload.Data.ID)
	assert.NotEmpty(t, castedCreateTeant.Payload.Data.Attributes.Name)
	assert.NotEmpty(t, castedCreateTeant.Payload.Data.Attributes.URLSubdomain)
	assert.Equal(t, string(common.UserActive), *castedCreateTeant.Payload.Data.Attributes.State)
	assert.True(t, *castedCreateTeant.Payload.Data.Attributes.CreatedTimestamp > 0)
	assert.True(t, *castedCreateTeant.Payload.Data.Attributes.LastModifiedTimestamp > 0)

	// Make sure there record is created when the Tenant is
	existing := handlers.HandleGetAllMetadataConfigsV2(handlers.AllRoles, tenantDB)(tenant_provisioning_service_v2.GetAllMetadataConfigsV2Params{HTTPRequest: createHttpRequestWithParams(*castedCreateTeant.Payload.Data.ID, handlers.UserRoleSkylight, MetadataConfigUrl, "GET")})
	castedResponse := existing.(*tenant_provisioning_service_v2.GetAllMetadataConfigsV2OK)
	assert.Equal(t, 1, len(castedResponse.Payload.Data))

	existingRecord := castedResponse.Payload.Data[0]

	// Should only be able to create 1 profile
	created := handlers.HandleCreateMetadataConfigV2(handlers.AllRoles, tenantDB)(tenant_provisioning_service_v2.CreateMetadataConfigV2Params{Body: generateTenantMetadataConfigCreationRequest(), HTTPRequest: createHttpRequestWithParams(*castedCreateTeant.Payload.Data.ID, handlers.UserRoleSkylight, MetadataConfigUrl, "POST")})
	castedCreate := created.(*tenant_provisioning_service_v2.CreateMetadataConfigV2Conflict)
	assert.NotNil(t, castedCreate)

	// Make sure we can retrieve this record:
	fetched := handlers.HandleGetMetadataConfigV2(handlers.AllRoles, tenantDB)(tenant_provisioning_service_v2.GetMetadataConfigV2Params{MetadataConfigID: *existingRecord.ID, HTTPRequest: createHttpRequestWithParams(*castedCreateTeant.Payload.Data.ID, handlers.UserRoleSkylight, MetadataConfigUrl, "GET")})
	castedFetch := fetched.(*tenant_provisioning_service_v2.GetMetadataConfigV2OK)
	assert.NotNil(t, castedFetch)
	assert.Equal(t, *existingRecord.Attributes.StartPoint, *castedFetch.Payload.Data.Attributes.StartPoint)
	assert.Equal(t, *existingRecord.Attributes.EndPoint, *castedFetch.Payload.Data.Attributes.EndPoint)
	assert.Empty(t, castedFetch.Payload.Data.Attributes.MidPoints)

	// Also retrieve the record as part of an array
	fetchList := handlers.HandleGetAllMetadataConfigsV2(handlers.AllRoles, tenantDB)(tenant_provisioning_service_v2.GetAllMetadataConfigsV2Params{HTTPRequest: createHttpRequestWithParams(*castedCreateTeant.Payload.Data.ID, handlers.UserRoleSkylight, MetadataConfigUrl, "GET")})
	castedFetchList := fetchList.(*tenant_provisioning_service_v2.GetAllMetadataConfigsV2OK)
	assert.NotNil(t, castedFetchList)
	assert.Equal(t, 1, len(castedFetchList.Payload.Data))
	assert.Equal(t, *existingRecord.Attributes.StartPoint, *castedFetchList.Payload.Data[0].Attributes.StartPoint)
	assert.Equal(t, *existingRecord.Attributes.EndPoint, *castedFetchList.Payload.Data[0].Attributes.EndPoint)
	assert.Empty(t, castedFetchList.Payload.Data[0].Attributes.MidPoints)

	// Make an update to the Record
	start := "makenew"
	end := "point"
	updateRequestBody := generateMetadataConfigUpdateRequest(*existingRecord.ID, *existingRecord.Attributes.Rev, &start, &end)
	updated := handlers.HandleUpdateMetadataConfigV2(handlers.AllRoles, tenantDB)(tenant_provisioning_service_v2.UpdateMetadataConfigV2Params{MetadataConfigID: *existingRecord.ID, Body: updateRequestBody, HTTPRequest: createHttpRequestWithParams(*castedCreateTeant.Payload.Data.ID, handlers.UserRoleSkylight, MetadataConfigUrl, "PATCH")})
	castedUpdate := updated.(*tenant_provisioning_service_v2.UpdateMetadataConfigV2OK)
	assert.NotNil(t, castedUpdate)
	assert.NotEqual(t, existingRecord, castedUpdate.Payload.Data)
	assert.NotEqual(t, existingRecord.Attributes.Rev, castedUpdate.Payload.Data.Attributes.Rev)
	assert.NotEqual(t, *existingRecord.Attributes.StartPoint, *castedUpdate.Payload.Data.Attributes.StartPoint)
	assert.NotEqual(t, *existingRecord.Attributes.EndPoint, *castedUpdate.Payload.Data.Attributes.EndPoint)
	assert.Equal(t, start, *castedUpdate.Payload.Data.Attributes.StartPoint)
	assert.Equal(t, end, *castedUpdate.Payload.Data.Attributes.EndPoint)

	// Delete the record
	deleted := handlers.HandleDeleteMetadataConfigV2(handlers.AllRoles, tenantDB)(tenant_provisioning_service_v2.DeleteMetadataConfigV2Params{MetadataConfigID: *existingRecord.ID, HTTPRequest: createHttpRequestWithParams(*castedCreateTeant.Payload.Data.ID, handlers.UserRoleSkylight, MetadataConfigUrl, "DELETE")})
	castedDelete := deleted.(*tenant_provisioning_service_v2.DeleteMetadataConfigV2OK)
	assert.NotNil(t, castedDelete)
	assert.Equal(t, castedUpdate.Payload.Data, castedDelete.Payload.Data)

	// Make sure there are no existing MetadataConfigs
	existingDNE := handlers.HandleGetAllMetadataConfigsV2(handlers.AllRoles, tenantDB)(tenant_provisioning_service_v2.GetAllMetadataConfigsV2Params{HTTPRequest: createHttpRequestWithParams(*castedCreateTeant.Payload.Data.ID, handlers.UserRoleSkylight, MetadataConfigUrl, "GET")})
	castedResponseDNE := existingDNE.(*tenant_provisioning_service_v2.GetAllMetadataConfigsV2NotFound)
	assert.NotNil(t, castedResponseDNE)
}

func TestMetadataConfigNotFoundV2(t *testing.T) {
	notFoundID := fake.CharactersN(20)

	// Get MetadataConfig
	fetched := handlers.HandleGetMetadataConfigV2(handlers.AllRoles, tenantDB)(tenant_provisioning_service_v2.GetMetadataConfigV2Params{MetadataConfigID: notFoundID, HTTPRequest: createHttpRequestWithParams("", handlers.UserRoleSkylight, MetadataConfigUrl, "GET")})
	castedFetch := fetched.(*tenant_provisioning_service_v2.GetMetadataConfigV2NotFound)
	assert.NotNil(t, castedFetch)

	// Delete MetadataConfig
	deleted := handlers.HandleDeleteMetadataConfigV2(handlers.AllRoles, tenantDB)(tenant_provisioning_service_v2.DeleteMetadataConfigV2Params{MetadataConfigID: notFoundID, HTTPRequest: createHttpRequestWithParams("", handlers.UserRoleSkylight, MetadataConfigUrl, "DELETE")})
	castedDelete := deleted.(*tenant_provisioning_service_v2.DeleteMetadataConfigV2NotFound)
	assert.NotNil(t, castedDelete)

	// Patch MetadataConfig
	updateRequest := generateMetadataConfigUpdateRequest(notFoundID, "reviosionstuff", nil, nil)
	updated := handlers.HandleUpdateMetadataConfigV2(handlers.AllRoles, tenantDB)(tenant_provisioning_service_v2.UpdateMetadataConfigV2Params{MetadataConfigID: notFoundID, Body: updateRequest, HTTPRequest: createHttpRequestWithParams("", handlers.UserRoleSkylight, MetadataConfigUrl, "PATCH")})
	castedUpdate := updated.(*tenant_provisioning_service_v2.UpdateMetadataConfigV2NotFound)
	assert.NotNil(t, castedUpdate)
}

func TestMetadataConfigBadRequestV2(t *testing.T) {

	// CreateMetadataConfig
	created := handlers.HandleCreateMetadataConfigV2(handlers.AllRoles, tenantDB)(tenant_provisioning_service_v2.CreateMetadataConfigV2Params{Body: nil, HTTPRequest: createHttpRequestWithParams("", handlers.UserRoleSkylight, MetadataConfigUrl, "POST")})
	castedCreate := created.(*tenant_provisioning_service_v2.CreateMetadataConfigV2BadRequest)
	assert.NotNil(t, castedCreate)

	// Update MetadataConfig
	updated := handlers.HandleUpdateMetadataConfigV2(handlers.AllRoles, tenantDB)(tenant_provisioning_service_v2.UpdateMetadataConfigV2Params{MetadataConfigID: fake.CharactersN(20), Body: nil, HTTPRequest: createHttpRequestWithParams("", handlers.UserRoleSkylight, MetadataConfigUrl, "PATCH")})
	castedUpdate := updated.(*tenant_provisioning_service_v2.UpdateMetadataConfigV2BadRequest)
	assert.NotNil(t, castedUpdate)
}

func TestMetadataConfigConflictV2(t *testing.T) {

	createdTenant := handlers.HandleCreateTenantV2(handlers.SkylightAdminRoleOnly, adminDB, tenantDB)(admin_provisioning_service_v2.CreateTenantV2Params{Body: generateRandomTenantCreationRequest(), HTTPRequest: createHttpRequestWithParams("", handlers.UserRoleSkylight, tenantURL, "POST")})
	castedCreateTeant := createdTenant.(*admin_provisioning_service_v2.CreateTenantV2Created)
	assert.NotNil(t, castedCreateTeant)
	assert.NotEmpty(t, castedCreateTeant.Payload.Data.ID)
	assert.NotEmpty(t, castedCreateTeant.Payload.Data.Attributes.Name)
	assert.NotEmpty(t, castedCreateTeant.Payload.Data.Attributes.URLSubdomain)
	assert.Equal(t, string(common.UserActive), *castedCreateTeant.Payload.Data.Attributes.State)
	assert.True(t, *castedCreateTeant.Payload.Data.Attributes.CreatedTimestamp > 0)
	assert.True(t, *castedCreateTeant.Payload.Data.Attributes.LastModifiedTimestamp > 0)

	// Make sure the record is created when the Tenant is
	existing := handlers.HandleGetAllMetadataConfigsV2(handlers.AllRoles, tenantDB)(tenant_provisioning_service_v2.GetAllMetadataConfigsV2Params{HTTPRequest: createHttpRequestWithParams(*castedCreateTeant.Payload.Data.ID, handlers.UserRoleSkylight, MetadataConfigUrl, "GET")})
	castedResponse := existing.(*tenant_provisioning_service_v2.GetAllMetadataConfigsV2OK)
	assert.Equal(t, 1, len(castedResponse.Payload.Data))

	existingRecord := castedResponse.Payload.Data[0]

	// Try the update with a bad revision
	updateRequestBody := generateMetadataConfigUpdateRequest(*existingRecord.ID, *existingRecord.Attributes.Rev+"pork", nil, nil)
	updated := handlers.HandleUpdateMetadataConfigV2(handlers.AllRoles, tenantDB)(tenant_provisioning_service_v2.UpdateMetadataConfigV2Params{MetadataConfigID: *existingRecord.ID, Body: updateRequestBody, HTTPRequest: createHttpRequestWithParams(*castedCreateTeant.Payload.Data.ID, handlers.UserRoleSkylight, MetadataConfigUrl, "PATCH")})
	castedUpdate := updated.(*tenant_provisioning_service_v2.UpdateMetadataConfigV2Conflict)
	assert.NotNil(t, castedUpdate)

	// Delete the tenant
	deleted := handlers.HandleDeleteMetadataConfigV2(handlers.AllRoles, tenantDB)(tenant_provisioning_service_v2.DeleteMetadataConfigV2Params{MetadataConfigID: *existingRecord.ID, HTTPRequest: createHttpRequestWithParams(*castedCreateTeant.Payload.Data.ID, handlers.UserRoleSkylight, MetadataConfigUrl, "DELETE")})
	castedDelete := deleted.(*tenant_provisioning_service_v2.DeleteMetadataConfigV2OK)
	assert.NotNil(t, castedDelete)
	assert.NotNil(t, castedDelete.Payload.Data)
}

func TestMetadataConfigAPIsProtectedByAuthV2(t *testing.T) {
	fakeTenantID := fake.CharactersN(20)
	// Get All - All Users
	existing := handlers.HandleGetAllMetadataConfigsV2(handlers.AllRoles, tenantDB)(tenant_provisioning_service_v2.GetAllMetadataConfigsV2Params{HTTPRequest: createHttpRequestWithParams(fakeTenantID, handlers.UserRoleUnknown, MetadataConfigUrl, "GET")})
	castedResponse := existing.(*tenant_provisioning_service_v2.GetAllMetadataConfigsV2Forbidden)
	assert.NotNil(t, castedResponse)

	fakeID := fake.CharactersN(20)

	// Get - All Users
	fetched := handlers.HandleGetMetadataConfigV2(handlers.AllRoles, tenantDB)(tenant_provisioning_service_v2.GetMetadataConfigV2Params{MetadataConfigID: fakeID, HTTPRequest: createHttpRequestWithParams(fakeTenantID, handlers.UserRoleUnknown, MetadataConfigUrl, "GET")})
	castedFetch := fetched.(*tenant_provisioning_service_v2.GetMetadataConfigV2Forbidden)
	assert.NotNil(t, castedFetch)

	// Create - SkylightAdmin and TenantAdmin Only
	created := handlers.HandleCreateMetadataConfigV2(handlers.SkylightAdminRoleOnly, tenantDB)(tenant_provisioning_service_v2.CreateMetadataConfigV2Params{Body: generateTenantMetadataConfigCreationRequest(), HTTPRequest: createHttpRequestWithParams(fakeTenantID, handlers.UserRoleTenantAdmin, MetadataConfigUrl, "POST")})
	castedCreate := created.(*tenant_provisioning_service_v2.CreateMetadataConfigV2Forbidden)
	assert.NotNil(t, castedCreate)

	created = handlers.HandleCreateMetadataConfigV2(handlers.SkylightAdminRoleOnly, tenantDB)(tenant_provisioning_service_v2.CreateMetadataConfigV2Params{Body: generateTenantMetadataConfigCreationRequest(), HTTPRequest: createHttpRequestWithParams(fakeTenantID, handlers.UserRoleTenantUser, MetadataConfigUrl, "POST")})
	castedCreate = created.(*tenant_provisioning_service_v2.CreateMetadataConfigV2Forbidden)
	assert.NotNil(t, castedCreate)

	// Update - SkylightAdmin and TenantAdmin Only
	updated := handlers.HandleUpdateMetadataConfigV2(handlers.SkylightAndTenantAdminRoles, tenantDB)(tenant_provisioning_service_v2.UpdateMetadataConfigV2Params{MetadataConfigID: fakeID, Body: nil, HTTPRequest: createHttpRequestWithParams(fakeTenantID, handlers.UserRoleUnknown, MetadataConfigUrl, "PATCH")})
	castedUpdate := updated.(*tenant_provisioning_service_v2.UpdateMetadataConfigV2Forbidden)
	assert.NotNil(t, castedUpdate)

	updated = handlers.HandleUpdateMetadataConfigV2(handlers.SkylightAndTenantAdminRoles, tenantDB)(tenant_provisioning_service_v2.UpdateMetadataConfigV2Params{MetadataConfigID: fakeID, Body: nil, HTTPRequest: createHttpRequestWithParams(fakeTenantID, handlers.UserRoleTenantUser, MetadataConfigUrl, "PATCH")})
	castedUpdate = updated.(*tenant_provisioning_service_v2.UpdateMetadataConfigV2Forbidden)
	assert.NotNil(t, castedUpdate)

	// Delete - SkylightAdmin and TenantAdmin Only
	deleted := handlers.HandleDeleteMetadataConfigV2(handlers.SkylightAdminRoleOnly, tenantDB)(tenant_provisioning_service_v2.DeleteMetadataConfigV2Params{MetadataConfigID: fakeID, HTTPRequest: createHttpRequestWithParams(fakeTenantID, handlers.UserRoleTenantAdmin, MetadataConfigUrl, "DELETE")})
	castedDelete := deleted.(*tenant_provisioning_service_v2.DeleteMetadataConfigV2Forbidden)
	assert.NotNil(t, castedDelete)

	deleted = handlers.HandleDeleteMetadataConfigV2(handlers.SkylightAdminRoleOnly, tenantDB)(tenant_provisioning_service_v2.DeleteMetadataConfigV2Params{MetadataConfigID: fakeID, HTTPRequest: createHttpRequestWithParams(fakeTenantID, handlers.UserRoleTenantUser, MetadataConfigUrl, "DELETE")})
	castedDelete = deleted.(*tenant_provisioning_service_v2.DeleteMetadataConfigV2Forbidden)
	assert.NotNil(t, castedDelete)
}

func generateTenantMetadataConfigCreationRequest() *swagmodels.MetadataConfigCreateRequest {
	start := fake.CharactersN(12)
	end := fake.CharactersN(12)
	mid := []string{fake.CharactersN(12), fake.CharactersN(6), fake.CharactersN(9)}
	result := swagmodels.MetadataConfigCreateRequest{
		Data: &swagmodels.MetadataConfigCreateRequestData{
			Type: &MetadataConfigTypeString,
			Attributes: &swagmodels.MetadataConfigCreateRequestDataAttributes{
				StartPoint: &start,
				EndPoint:   &end,
				MidPoints:  mid,
			},
		},
	}

	return &result
}

func generateMetadataConfigUpdateRequest(id string, rev string, start *string, end *string) *swagmodels.MetadataConfigUpdateRequest {
	result := swagmodels.MetadataConfigUpdateRequest{
		Data: &swagmodels.MetadataConfigUpdateRequestData{
			Type:       &MetadataConfigTypeString,
			Attributes: &swagmodels.MetadataConfigUpdateRequestDataAttributes{},
		},
	}

	result.Data.ID = &id
	result.Data.Attributes.Rev = &rev

	if start != nil {
		result.Data.Attributes.StartPoint = *start
	}

	if end != nil {
		result.Data.Attributes.EndPoint = *end
	}
	return &result
}
