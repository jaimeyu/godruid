package handlers_test

import (
	"encoding/json"
	"testing"

	"github.com/icrowley/fake"
	"github.com/stretchr/testify/assert"

	"github.com/accedian/adh-gather/handlers"
	"github.com/accedian/adh-gather/logger"
	"github.com/accedian/adh-gather/models/common"
	"github.com/accedian/adh-gather/restapi/operations/admin_provisioning_service_v2"
	"github.com/accedian/adh-gather/restapi/operations/tenant_provisioning_service_v2"
	"github.com/accedian/adh-gather/swagmodels"
)

var (
	IngestionProfileUrl = "http://deployment.test.cool/api/v2/ingestionProfiles"

	IngestionProfileTypeString = "ingestionProfiles"
)

func TestIngestionProfileCRUDV2(t *testing.T) {

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
	existing := handlers.HandleGetAllIngestionProfilesV2(handlers.AllRoles, tenantDB)(tenant_provisioning_service_v2.GetAllIngestionProfilesV2Params{HTTPRequest: createHttpRequestWithParams(*castedCreateTeant.Payload.Data.ID, handlers.UserRoleSkylight, IngestionProfileUrl, "GET")})
	castedResponse := existing.(*tenant_provisioning_service_v2.GetAllIngestionProfilesV2OK)
	assert.Equal(t, 1, len(castedResponse.Payload.Data))

	existingRecord := castedResponse.Payload.Data[0]

	// Should only be able to create 1 profile
	created := handlers.HandleCreateIngestionProfileV2(handlers.AllRoles, tenantDB)(tenant_provisioning_service_v2.CreateIngestionProfileV2Params{Body: generateTenantIngestionProfileCreationRequest(), HTTPRequest: createHttpRequestWithParams(*castedCreateTeant.Payload.Data.ID, handlers.UserRoleSkylight, IngestionProfileUrl, "POST")})
	castedCreate := created.(*tenant_provisioning_service_v2.CreateIngestionProfileV2Conflict)
	assert.NotNil(t, castedCreate)

	// Make sure we can retrieve this record:
	fetched := handlers.HandleGetIngestionProfileV2(handlers.AllRoles, tenantDB)(tenant_provisioning_service_v2.GetIngestionProfileV2Params{IngestionProfileID: *existingRecord.ID, HTTPRequest: createHttpRequestWithParams(*castedCreateTeant.Payload.Data.ID, handlers.UserRoleSkylight, IngestionProfileUrl, "GET")})
	castedFetch := fetched.(*tenant_provisioning_service_v2.GetIngestionProfileV2OK)
	assert.NotNil(t, castedFetch)
	assert.Equal(t, existingRecord.Attributes.Metrics, castedFetch.Payload.Data.Attributes.Metrics)
	assert.NotEmpty(t, castedFetch.Payload.Data.Attributes.MetricList)

	// Also retrieve the record as part of an array
	fetchList := handlers.HandleGetAllIngestionProfilesV2(handlers.AllRoles, tenantDB)(tenant_provisioning_service_v2.GetAllIngestionProfilesV2Params{HTTPRequest: createHttpRequestWithParams(*castedCreateTeant.Payload.Data.ID, handlers.UserRoleSkylight, IngestionProfileUrl, "GET")})
	castedFetchList := fetchList.(*tenant_provisioning_service_v2.GetAllIngestionProfilesV2OK)
	assert.NotNil(t, castedFetchList)
	assert.Equal(t, 1, len(castedFetchList.Payload.Data))
	assert.Equal(t, existingRecord.Attributes.Metrics, castedFetchList.Payload.Data[0].Attributes.Metrics)
	assert.NotEmpty(t, castedFetchList.Payload.Data[0].Attributes.MetricList)

	// Make an update to the Record
	updateRequestBody := generateIngestionProfileUpdateRequest(*existingRecord.ID, *existingRecord.Attributes.Rev)
	updated := handlers.HandleUpdateIngestionProfileV2(handlers.AllRoles, tenantDB)(tenant_provisioning_service_v2.UpdateIngestionProfileV2Params{IngestionProfileID: *existingRecord.ID, Body: updateRequestBody, HTTPRequest: createHttpRequestWithParams(*castedCreateTeant.Payload.Data.ID, handlers.UserRoleSkylight, IngestionProfileUrl, "PATCH")})
	castedUpdate := updated.(*tenant_provisioning_service_v2.UpdateIngestionProfileV2OK)
	assert.NotNil(t, castedUpdate)
	assert.NotEqual(t, existingRecord, castedUpdate.Payload.Data)
	assert.NotEqual(t, existingRecord.Attributes.Rev, castedUpdate.Payload.Data.Attributes.Rev)
	assert.NotEqual(t, existingRecord.Attributes.Metrics, castedUpdate.Payload.Data.Attributes.Metrics)

	// Delete the record
	deleted := handlers.HandleDeleteIngestionProfileV2(handlers.AllRoles, tenantDB)(tenant_provisioning_service_v2.DeleteIngestionProfileV2Params{IngestionProfileID: *existingRecord.ID, HTTPRequest: createHttpRequestWithParams(*castedCreateTeant.Payload.Data.ID, handlers.UserRoleSkylight, IngestionProfileUrl, "DELETE")})
	castedDelete := deleted.(*tenant_provisioning_service_v2.DeleteIngestionProfileV2OK)
	assert.NotNil(t, castedDelete)
	assert.Equal(t, castedUpdate.Payload.Data, castedDelete.Payload.Data)

	// Make sure there are no existing IngestionProfiles
	existingDNE := handlers.HandleGetAllIngestionProfilesV2(handlers.AllRoles, tenantDB)(tenant_provisioning_service_v2.GetAllIngestionProfilesV2Params{HTTPRequest: createHttpRequestWithParams(*castedCreateTeant.Payload.Data.ID, handlers.UserRoleSkylight, IngestionProfileUrl, "GET")})
	castedResponseDNE := existingDNE.(*tenant_provisioning_service_v2.GetAllIngestionProfilesV2NotFound)
	assert.NotNil(t, castedResponseDNE)
}

func TestIngestionProfileNotFoundV2(t *testing.T) {
	notFoundID := fake.CharactersN(20)

	// Get IngestionProfile
	fetched := handlers.HandleGetIngestionProfileV2(handlers.AllRoles, tenantDB)(tenant_provisioning_service_v2.GetIngestionProfileV2Params{IngestionProfileID: notFoundID, HTTPRequest: createHttpRequestWithParams("", handlers.UserRoleSkylight, IngestionProfileUrl, "GET")})
	castedFetch := fetched.(*tenant_provisioning_service_v2.GetIngestionProfileV2NotFound)
	assert.NotNil(t, castedFetch)

	// Delete IngestionProfile
	deleted := handlers.HandleDeleteIngestionProfileV2(handlers.AllRoles, tenantDB)(tenant_provisioning_service_v2.DeleteIngestionProfileV2Params{IngestionProfileID: notFoundID, HTTPRequest: createHttpRequestWithParams("", handlers.UserRoleSkylight, IngestionProfileUrl, "DELETE")})
	castedDelete := deleted.(*tenant_provisioning_service_v2.DeleteIngestionProfileV2NotFound)
	assert.NotNil(t, castedDelete)

	// Patch IngestionProfile
	updateRequest := generateIngestionProfileUpdateRequest(notFoundID, "reviosionstuff")
	updated := handlers.HandleUpdateIngestionProfileV2(handlers.AllRoles, tenantDB)(tenant_provisioning_service_v2.UpdateIngestionProfileV2Params{IngestionProfileID: notFoundID, Body: updateRequest, HTTPRequest: createHttpRequestWithParams("", handlers.UserRoleSkylight, IngestionProfileUrl, "PATCH")})
	castedUpdate := updated.(*tenant_provisioning_service_v2.UpdateIngestionProfileV2NotFound)
	assert.NotNil(t, castedUpdate)
}

func TestIngestionProfileBadRequestV2(t *testing.T) {

	// CreateIngestionProfile
	created := handlers.HandleCreateIngestionProfileV2(handlers.AllRoles, tenantDB)(tenant_provisioning_service_v2.CreateIngestionProfileV2Params{Body: nil, HTTPRequest: createHttpRequestWithParams("", handlers.UserRoleSkylight, IngestionProfileUrl, "POST")})
	castedCreate := created.(*tenant_provisioning_service_v2.CreateIngestionProfileV2BadRequest)
	assert.NotNil(t, castedCreate)

	// Update IngestionProfile
	updated := handlers.HandleUpdateIngestionProfileV2(handlers.AllRoles, tenantDB)(tenant_provisioning_service_v2.UpdateIngestionProfileV2Params{IngestionProfileID: fake.CharactersN(20), Body: nil, HTTPRequest: createHttpRequestWithParams("", handlers.UserRoleSkylight, IngestionProfileUrl, "PATCH")})
	castedUpdate := updated.(*tenant_provisioning_service_v2.UpdateIngestionProfileV2BadRequest)
	assert.NotNil(t, castedUpdate)
}

func TestIngestionProfileConflictV2(t *testing.T) {

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
	existing := handlers.HandleGetAllIngestionProfilesV2(handlers.AllRoles, tenantDB)(tenant_provisioning_service_v2.GetAllIngestionProfilesV2Params{HTTPRequest: createHttpRequestWithParams(*castedCreateTeant.Payload.Data.ID, handlers.UserRoleSkylight, IngestionProfileUrl, "GET")})
	castedResponse := existing.(*tenant_provisioning_service_v2.GetAllIngestionProfilesV2OK)
	assert.Equal(t, 1, len(castedResponse.Payload.Data))

	existingRecord := castedResponse.Payload.Data[0]

	// Try the update with a bad revision
	updateRequestBody := generateIngestionProfileUpdateRequest(*existingRecord.ID, *existingRecord.Attributes.Rev+"pork")
	updated := handlers.HandleUpdateIngestionProfileV2(handlers.AllRoles, tenantDB)(tenant_provisioning_service_v2.UpdateIngestionProfileV2Params{IngestionProfileID: *existingRecord.ID, Body: updateRequestBody, HTTPRequest: createHttpRequestWithParams(*castedCreateTeant.Payload.Data.ID, handlers.UserRoleSkylight, IngestionProfileUrl, "PATCH")})
	castedUpdate := updated.(*tenant_provisioning_service_v2.UpdateIngestionProfileV2Conflict)
	assert.NotNil(t, castedUpdate)

	// Delete the tenant
	deleted := handlers.HandleDeleteIngestionProfileV2(handlers.AllRoles, tenantDB)(tenant_provisioning_service_v2.DeleteIngestionProfileV2Params{IngestionProfileID: *existingRecord.ID, HTTPRequest: createHttpRequestWithParams(*castedCreateTeant.Payload.Data.ID, handlers.UserRoleSkylight, IngestionProfileUrl, "DELETE")})
	castedDelete := deleted.(*tenant_provisioning_service_v2.DeleteIngestionProfileV2OK)
	assert.NotNil(t, castedDelete)
	assert.NotNil(t, castedDelete.Payload.Data)
}

func TestIngestionProfileAPIsProtectedByAuthV2(t *testing.T) {
	fakeTenantID := fake.CharactersN(20)
	// Get All - All Users
	existing := handlers.HandleGetAllIngestionProfilesV2(handlers.AllRoles, tenantDB)(tenant_provisioning_service_v2.GetAllIngestionProfilesV2Params{HTTPRequest: createHttpRequestWithParams(fakeTenantID, handlers.UserRoleUnknown, IngestionProfileUrl, "GET")})
	castedResponse := existing.(*tenant_provisioning_service_v2.GetAllIngestionProfilesV2Forbidden)
	assert.NotNil(t, castedResponse)

	fakeID := fake.CharactersN(20)

	// Get - All Users
	fetched := handlers.HandleGetIngestionProfileV2(handlers.AllRoles, tenantDB)(tenant_provisioning_service_v2.GetIngestionProfileV2Params{IngestionProfileID: fakeID, HTTPRequest: createHttpRequestWithParams(fakeTenantID, handlers.UserRoleUnknown, IngestionProfileUrl, "GET")})
	castedFetch := fetched.(*tenant_provisioning_service_v2.GetIngestionProfileV2Forbidden)
	assert.NotNil(t, castedFetch)

	// Create - SkylightAdmin and TenantAdmin Only
	created := handlers.HandleCreateIngestionProfileV2(handlers.SkylightAdminRoleOnly, tenantDB)(tenant_provisioning_service_v2.CreateIngestionProfileV2Params{Body: generateTenantIngestionProfileCreationRequest(), HTTPRequest: createHttpRequestWithParams(fakeTenantID, handlers.UserRoleTenantAdmin, IngestionProfileUrl, "POST")})
	castedCreate := created.(*tenant_provisioning_service_v2.CreateIngestionProfileV2Forbidden)
	assert.NotNil(t, castedCreate)

	created = handlers.HandleCreateIngestionProfileV2(handlers.SkylightAdminRoleOnly, tenantDB)(tenant_provisioning_service_v2.CreateIngestionProfileV2Params{Body: generateTenantIngestionProfileCreationRequest(), HTTPRequest: createHttpRequestWithParams(fakeTenantID, handlers.UserRoleTenantUser, IngestionProfileUrl, "POST")})
	castedCreate = created.(*tenant_provisioning_service_v2.CreateIngestionProfileV2Forbidden)
	assert.NotNil(t, castedCreate)

	// Update - SkylightAdmin and TenantAdmin Only
	updated := handlers.HandleUpdateIngestionProfileV2(handlers.SkylightAndTenantAdminRoles, tenantDB)(tenant_provisioning_service_v2.UpdateIngestionProfileV2Params{IngestionProfileID: fakeID, Body: nil, HTTPRequest: createHttpRequestWithParams(fakeTenantID, handlers.UserRoleUnknown, IngestionProfileUrl, "PATCH")})
	castedUpdate := updated.(*tenant_provisioning_service_v2.UpdateIngestionProfileV2Forbidden)
	assert.NotNil(t, castedUpdate)

	updated = handlers.HandleUpdateIngestionProfileV2(handlers.SkylightAndTenantAdminRoles, tenantDB)(tenant_provisioning_service_v2.UpdateIngestionProfileV2Params{IngestionProfileID: fakeID, Body: nil, HTTPRequest: createHttpRequestWithParams(fakeTenantID, handlers.UserRoleTenantUser, IngestionProfileUrl, "PATCH")})
	castedUpdate = updated.(*tenant_provisioning_service_v2.UpdateIngestionProfileV2Forbidden)
	assert.NotNil(t, castedUpdate)

	// Delete - SkylightAdmin and TenantAdmin Only
	deleted := handlers.HandleDeleteIngestionProfileV2(handlers.SkylightAdminRoleOnly, tenantDB)(tenant_provisioning_service_v2.DeleteIngestionProfileV2Params{IngestionProfileID: fakeID, HTTPRequest: createHttpRequestWithParams(fakeTenantID, handlers.UserRoleTenantAdmin, IngestionProfileUrl, "DELETE")})
	castedDelete := deleted.(*tenant_provisioning_service_v2.DeleteIngestionProfileV2Forbidden)
	assert.NotNil(t, castedDelete)

	deleted = handlers.HandleDeleteIngestionProfileV2(handlers.SkylightAdminRoleOnly, tenantDB)(tenant_provisioning_service_v2.DeleteIngestionProfileV2Params{IngestionProfileID: fakeID, HTTPRequest: createHttpRequestWithParams(fakeTenantID, handlers.UserRoleTenantUser, IngestionProfileUrl, "DELETE")})
	castedDelete = deleted.(*tenant_provisioning_service_v2.DeleteIngestionProfileV2Forbidden)
	assert.NotNil(t, castedDelete)
}

func generateTenantIngestionProfileCreationRequest() *swagmodels.IngestionProfileCreateRequest {
	bytes := []byte(`{
		"data": {
			"attributes": {
				"metrics": {
					"vendorMap": {
						"accedian-flowmeter": {
							"monitoredObjectTypeMap": {
								"flowmeter": {
									"metricMap": {
										"bytesReceived": true,
										"packetsReceived": true,
										"throughputAvg": true,
										"throughputMax": true,
										"throughputMin": true
									}
								}
							}
						}
					}
				}
			},
			"type": "ingestionProfiles"
		}
	}`)

	result := swagmodels.IngestionProfileCreateRequest{}
	if err := json.Unmarshal(bytes, &result); err != nil {
		logger.Log.Error(err.Error())
	}

	return &result
}

func generateIngestionProfileUpdateRequest(id string, rev string) *swagmodels.IngestionProfileUpdateRequest {
	bytes := []byte(`{
		"data": {
			"attributes": {
				"_rev": "1-c9951e5eac4202d85d181bd4dd8453c2",
				"metrics": {
					"vendorMap": {
						"accedian-flowmeter": {
							"monitoredObjectTypeMap": {
								"flowmeter": {
									"metricMap": {
										"bytesReceived": false,
										"packetsReceived": false,
										"throughputAvg": true,
										"throughputMax": true,
										"throughputMin": true
									}
								}
							}
						}
					}
				}
			},
			"id": "0be6f37c-f1b6-4a79-98d3-eb1cb2cbf887",
			"type": "ingestionProfiles"
		}
	}`)

	result := swagmodels.IngestionProfileUpdateRequest{}
	if err := json.Unmarshal(bytes, &result); err != nil {
		logger.Log.Error(err.Error())
	}

	result.Data.ID = &id
	result.Data.Attributes.Rev = &rev
	return &result
}
