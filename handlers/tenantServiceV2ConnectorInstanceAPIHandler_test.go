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
	connectorInstanceUrl = "http://deployment.test.cool/api/v2/connector-instances"

	connectorInstanceTypeString = "connectorInstances"
)

func TestConnectorInstanceCRUDV2(t *testing.T) {

	createdTenant := handlers.HandleCreateTenantV2(handlers.SkylightAdminRoleOnly, adminDB, tenantDB, stubbedMetricBaselineDatastore)(admin_provisioning_service_v2.CreateTenantV2Params{Body: generateRandomTenantCreationRequest(), HTTPRequest: createHttpRequestWithParams("", handlers.UserRoleSkylight, tenantURL, "POST")})
	castedCreateTeant := createdTenant.(*admin_provisioning_service_v2.CreateTenantV2Created)
	assert.NotNil(t, castedCreateTeant)
	assert.NotEmpty(t, castedCreateTeant.Payload.Data.ID)
	assert.NotEmpty(t, castedCreateTeant.Payload.Data.Attributes.Name)
	assert.NotEmpty(t, castedCreateTeant.Payload.Data.Attributes.URLSubdomain)
	assert.Equal(t, string(common.UserActive), *castedCreateTeant.Payload.Data.Attributes.State)
	assert.True(t, *castedCreateTeant.Payload.Data.Attributes.CreatedTimestamp > 0)
	assert.True(t, *castedCreateTeant.Payload.Data.Attributes.LastModifiedTimestamp > 0)

	// Make sure there are no existing ConnectorInstances
	existing := handlers.HandleGetAllConnectorInstancesV2(handlers.AllRoles, tenantDB)(tenant_provisioning_service_v2.GetAllConnectorInstancesV2Params{HTTPRequest: createHttpRequestWithParams(*castedCreateTeant.Payload.Data.ID, handlers.UserRoleSkylight, connectorInstanceUrl, "GET")})
	castedResponse := existing.(*tenant_provisioning_service_v2.GetAllConnectorInstancesV2NotFound)
	assert.NotNil(t, castedResponse)

	created := handlers.HandleCreateConnectorInstanceV2(handlers.AllRoles, tenantDB)(tenant_provisioning_service_v2.CreateConnectorInstanceV2Params{Body: generateRandomTenantConnectorInstanceCreationRequest(), HTTPRequest: createHttpRequestWithParams(*castedCreateTeant.Payload.Data.ID, handlers.UserRoleSkylight, connectorInstanceUrl, "POST")})
	castedCreate := created.(*tenant_provisioning_service_v2.CreateConnectorInstanceV2Created)
	assert.NotNil(t, castedCreate)
	assert.NotEmpty(t, castedCreate.Payload.Data.ID)
	assert.NotEmpty(t, castedCreate.Payload.Data.Attributes.Hostname)
	assert.NotEmpty(t, castedCreate.Payload.Data.Attributes.Status)
	assert.NotEmpty(t, castedCreate.Payload.Data.Attributes.Datatype)
	assert.True(t, *castedCreate.Payload.Data.Attributes.CreatedTimestamp > 0)
	assert.True(t, *castedCreate.Payload.Data.Attributes.LastModifiedTimestamp > 0)

	// Make sure we can retrieve this record:
	fetched := handlers.HandleGetConnectorInstanceV2(handlers.AllRoles, tenantDB)(tenant_provisioning_service_v2.GetConnectorInstanceV2Params{ConnectorInstanceID: *castedCreate.Payload.Data.ID, HTTPRequest: createHttpRequestWithParams(*castedCreateTeant.Payload.Data.ID, handlers.UserRoleSkylight, connectorInstanceUrl, "GET")})
	castedFetch := fetched.(*tenant_provisioning_service_v2.GetConnectorInstanceV2OK)
	assert.NotNil(t, castedFetch)
	assert.Equal(t, castedCreate.Payload.Data, castedFetch.Payload.Data)

	// Also retrieve the record as part of an array
	fetchList := handlers.HandleGetAllConnectorInstancesV2(handlers.AllRoles, tenantDB)(tenant_provisioning_service_v2.GetAllConnectorInstancesV2Params{HTTPRequest: createHttpRequestWithParams(*castedCreateTeant.Payload.Data.ID, handlers.UserRoleSkylight, connectorInstanceUrl, "GET")})
	castedFetchList := fetchList.(*tenant_provisioning_service_v2.GetAllConnectorInstancesV2OK)
	assert.NotNil(t, castedFetchList)
	assert.Equal(t, 1, len(castedFetchList.Payload.Data))
	assert.Equal(t, castedCreate.Payload.Data, castedFetchList.Payload.Data[0])

	// Make an update to the Record
	newName := fake.CharactersN(16)
	updateRequestBody := generateConnectorInstanceUpdateRequest(*castedCreate.Payload.Data.ID, *castedCreate.Payload.Data.Attributes.Rev, &newName, nil)
	updated := handlers.HandleUpdateConnectorInstanceV2(handlers.AllRoles, tenantDB)(tenant_provisioning_service_v2.UpdateConnectorInstanceV2Params{ConnectorInstanceID: *castedCreate.Payload.Data.ID, Body: updateRequestBody, HTTPRequest: createHttpRequestWithParams(*castedCreateTeant.Payload.Data.ID, handlers.UserRoleSkylight, connectorInstanceUrl, "PATCH")})
	castedUpdate := updated.(*tenant_provisioning_service_v2.UpdateConnectorInstanceV2OK)
	assert.NotNil(t, castedUpdate)
	assert.NotEqual(t, castedCreate.Payload.Data, castedUpdate.Payload.Data)
	assert.NotEqual(t, castedCreate.Payload.Data.Attributes.Rev, castedUpdate.Payload.Data.Attributes.Rev)
	assert.Equal(t, newName, *castedUpdate.Payload.Data.Attributes.Hostname)
	assert.Equal(t, castedCreate.Payload.Data.Attributes.Status, castedUpdate.Payload.Data.Attributes.Status)

	// Delete the record
	deleted := handlers.HandleDeleteConnectorInstanceV2(handlers.AllRoles, tenantDB)(tenant_provisioning_service_v2.DeleteConnectorInstanceV2Params{ConnectorInstanceID: *castedCreate.Payload.Data.ID, HTTPRequest: createHttpRequestWithParams(*castedCreateTeant.Payload.Data.ID, handlers.UserRoleSkylight, connectorInstanceUrl, "DELETE")})
	castedDelete := deleted.(*tenant_provisioning_service_v2.DeleteConnectorInstanceV2OK)
	assert.NotNil(t, castedDelete)
	assert.Equal(t, castedUpdate.Payload.Data, castedDelete.Payload.Data)

	// Make sure there are no existing ConnectorInstances
	existing = handlers.HandleGetAllConnectorInstancesV2(handlers.AllRoles, tenantDB)(tenant_provisioning_service_v2.GetAllConnectorInstancesV2Params{HTTPRequest: createHttpRequestWithParams(*castedCreateTeant.Payload.Data.ID, handlers.UserRoleSkylight, connectorInstanceUrl, "GET")})
	castedResponse = existing.(*tenant_provisioning_service_v2.GetAllConnectorInstancesV2NotFound)
	assert.NotNil(t, castedResponse)
}

func TestConnectorInstanceNotFoundV2(t *testing.T) {
	notFoundID := fake.CharactersN(20)

	// Get ConnectorInstance
	fetched := handlers.HandleGetConnectorInstanceV2(handlers.AllRoles, tenantDB)(tenant_provisioning_service_v2.GetConnectorInstanceV2Params{ConnectorInstanceID: notFoundID, HTTPRequest: createHttpRequestWithParams("", handlers.UserRoleSkylight, connectorInstanceUrl, "GET")})
	castedFetch := fetched.(*tenant_provisioning_service_v2.GetConnectorInstanceV2NotFound)
	assert.NotNil(t, castedFetch)

	// Delete ConnectorInstance
	deleted := handlers.HandleDeleteConnectorInstanceV2(handlers.AllRoles, tenantDB)(tenant_provisioning_service_v2.DeleteConnectorInstanceV2Params{ConnectorInstanceID: notFoundID, HTTPRequest: createHttpRequestWithParams("", handlers.UserRoleSkylight, connectorInstanceUrl, "DELETE")})
	castedDelete := deleted.(*tenant_provisioning_service_v2.DeleteConnectorInstanceV2NotFound)
	assert.NotNil(t, castedDelete)

	// Patch ConnectorInstance
	updateRequest := generateConnectorInstanceUpdateRequest(notFoundID, "reviosionstuff", nil, nil)
	updated := handlers.HandleUpdateConnectorInstanceV2(handlers.AllRoles, tenantDB)(tenant_provisioning_service_v2.UpdateConnectorInstanceV2Params{ConnectorInstanceID: notFoundID, Body: updateRequest, HTTPRequest: createHttpRequestWithParams("", handlers.UserRoleSkylight, connectorInstanceUrl, "PATCH")})
	castedUpdate := updated.(*tenant_provisioning_service_v2.UpdateConnectorInstanceV2NotFound)
	assert.NotNil(t, castedUpdate)
}

func TestConnectorInstanceBadRequestV2(t *testing.T) {

	// CreateConnectorInstance
	created := handlers.HandleCreateConnectorInstanceV2(handlers.AllRoles, tenantDB)(tenant_provisioning_service_v2.CreateConnectorInstanceV2Params{Body: nil, HTTPRequest: createHttpRequestWithParams("", handlers.UserRoleSkylight, connectorInstanceUrl, "POST")})
	castedCreate := created.(*tenant_provisioning_service_v2.CreateConnectorInstanceV2BadRequest)
	assert.NotNil(t, castedCreate)

	// Update ConnectorInstance
	updated := handlers.HandleUpdateConnectorInstanceV2(handlers.AllRoles, tenantDB)(tenant_provisioning_service_v2.UpdateConnectorInstanceV2Params{ConnectorInstanceID: fake.CharactersN(20), Body: nil, HTTPRequest: createHttpRequestWithParams("", handlers.UserRoleSkylight, connectorInstanceUrl, "PATCH")})
	castedUpdate := updated.(*tenant_provisioning_service_v2.UpdateConnectorInstanceV2BadRequest)
	assert.NotNil(t, castedUpdate)
}

func TestConnectorInstanceConflictV2(t *testing.T) {

	createdTenant := handlers.HandleCreateTenantV2(handlers.SkylightAdminRoleOnly, adminDB, tenantDB, stubbedMetricBaselineDatastore)(admin_provisioning_service_v2.CreateTenantV2Params{Body: generateRandomTenantCreationRequest(), HTTPRequest: createHttpRequestWithParams("", handlers.UserRoleSkylight, tenantURL, "POST")})
	castedCreateTeant := createdTenant.(*admin_provisioning_service_v2.CreateTenantV2Created)
	assert.NotNil(t, castedCreateTeant)
	assert.NotEmpty(t, castedCreateTeant.Payload.Data.ID)
	assert.NotEmpty(t, castedCreateTeant.Payload.Data.Attributes.Name)
	assert.NotEmpty(t, castedCreateTeant.Payload.Data.Attributes.URLSubdomain)
	assert.Equal(t, string(common.UserActive), *castedCreateTeant.Payload.Data.Attributes.State)
	assert.True(t, *castedCreateTeant.Payload.Data.Attributes.CreatedTimestamp > 0)
	assert.True(t, *castedCreateTeant.Payload.Data.Attributes.LastModifiedTimestamp > 0)

	// Make sure there are no existing ConnectorInstance
	existing := handlers.HandleGetAllConnectorInstancesV2(handlers.AllRoles, tenantDB)(tenant_provisioning_service_v2.GetAllConnectorInstancesV2Params{HTTPRequest: createHttpRequestWithParams(*castedCreateTeant.Payload.Data.ID, handlers.UserRoleSkylight, connectorInstanceUrl, "GET")})
	castedResponse := existing.(*tenant_provisioning_service_v2.GetAllConnectorInstancesV2NotFound)
	assert.NotNil(t, castedResponse)

	createReqBody := generateRandomTenantConnectorInstanceCreationRequest()
	created := handlers.HandleCreateConnectorInstanceV2(handlers.AllRoles, tenantDB)(tenant_provisioning_service_v2.CreateConnectorInstanceV2Params{Body: createReqBody, HTTPRequest: createHttpRequestWithParams(*castedCreateTeant.Payload.Data.ID, handlers.UserRoleSkylight, connectorInstanceUrl, "POST")})
	castedCreate := created.(*tenant_provisioning_service_v2.CreateConnectorInstanceV2Created)
	assert.NotNil(t, castedCreate)
	assert.NotEmpty(t, castedCreate.Payload.Data.ID)
	assert.NotEmpty(t, castedCreate.Payload.Data.Attributes.Hostname)
	assert.NotEmpty(t, castedCreate.Payload.Data.Attributes.Status)
	assert.NotEmpty(t, castedCreate.Payload.Data.Attributes.Datatype)
	assert.True(t, *castedCreate.Payload.Data.Attributes.CreatedTimestamp > 0)
	assert.True(t, *castedCreate.Payload.Data.Attributes.LastModifiedTimestamp > 0)

	// Try to create the record again - should succeed as we are not guarding against name collisiones
	createdConflict := handlers.HandleCreateConnectorInstanceV2(handlers.AllRoles, tenantDB)(tenant_provisioning_service_v2.CreateConnectorInstanceV2Params{Body: createReqBody, HTTPRequest: createHttpRequestWithParams(*castedCreateTeant.Payload.Data.ID, handlers.UserRoleSkylight, connectorInstanceUrl, "POST")})
	castedCreateConflictButOK := createdConflict.(*tenant_provisioning_service_v2.CreateConnectorInstanceV2Created)
	assert.NotNil(t, castedCreateConflictButOK)

	// Try the update with a bad revision
	newName := fake.CharactersN(16)
	updateRequestBody := generateConnectorInstanceUpdateRequest(*castedCreate.Payload.Data.ID, *castedCreate.Payload.Data.Attributes.Rev+"pork", &newName, nil)
	updated := handlers.HandleUpdateConnectorInstanceV2(handlers.AllRoles, tenantDB)(tenant_provisioning_service_v2.UpdateConnectorInstanceV2Params{ConnectorInstanceID: *castedCreate.Payload.Data.ID, Body: updateRequestBody, HTTPRequest: createHttpRequestWithParams(*castedCreateTeant.Payload.Data.ID, handlers.UserRoleSkylight, connectorInstanceUrl, "PATCH")})
	castedUpdate := updated.(*tenant_provisioning_service_v2.UpdateConnectorInstanceV2Conflict)
	assert.NotNil(t, castedUpdate)

	// Delete the tenant
	deleted := handlers.HandleDeleteConnectorInstanceV2(handlers.AllRoles, tenantDB)(tenant_provisioning_service_v2.DeleteConnectorInstanceV2Params{ConnectorInstanceID: *castedCreate.Payload.Data.ID, HTTPRequest: createHttpRequestWithParams(*castedCreateTeant.Payload.Data.ID, handlers.UserRoleSkylight, connectorInstanceUrl, "DELETE")})
	castedDelete := deleted.(*tenant_provisioning_service_v2.DeleteConnectorInstanceV2OK)
	assert.NotNil(t, castedDelete)
	assert.NotNil(t, castedDelete.Payload.Data)
}

func TestConnectorInstanceAPIsProtectedByAuthV2(t *testing.T) {
	fakeTenantID := fake.CharactersN(20)
	// Get All - All Users
	existing := handlers.HandleGetAllConnectorInstancesV2(handlers.AllRoles, tenantDB)(tenant_provisioning_service_v2.GetAllConnectorInstancesV2Params{HTTPRequest: createHttpRequestWithParams(fakeTenantID, handlers.UserRoleUnknown, connectorInstanceUrl, "GET")})
	castedResponse := existing.(*tenant_provisioning_service_v2.GetAllConnectorInstancesV2Forbidden)
	assert.NotNil(t, castedResponse)

	fakeID := fake.CharactersN(20)

	// Get - All Users
	fetched := handlers.HandleGetConnectorInstanceV2(handlers.AllRoles, tenantDB)(tenant_provisioning_service_v2.GetConnectorInstanceV2Params{ConnectorInstanceID: fakeID, HTTPRequest: createHttpRequestWithParams(fakeTenantID, handlers.UserRoleUnknown, connectorInstanceUrl, "GET")})
	castedFetch := fetched.(*tenant_provisioning_service_v2.GetConnectorInstanceV2Forbidden)
	assert.NotNil(t, castedFetch)

	// Create - SkylightAdmin and TenantAdmin Only
	created := handlers.HandleCreateConnectorInstanceV2(handlers.SkylightAndTenantAdminRoles, tenantDB)(tenant_provisioning_service_v2.CreateConnectorInstanceV2Params{Body: generateRandomTenantConnectorInstanceCreationRequest(), HTTPRequest: createHttpRequestWithParams(fakeTenantID, handlers.UserRoleUnknown, connectorInstanceUrl, "POST")})
	castedCreate := created.(*tenant_provisioning_service_v2.CreateConnectorInstanceV2Forbidden)
	assert.NotNil(t, castedCreate)

	created = handlers.HandleCreateConnectorInstanceV2(handlers.SkylightAndTenantAdminRoles, tenantDB)(tenant_provisioning_service_v2.CreateConnectorInstanceV2Params{Body: generateRandomTenantConnectorInstanceCreationRequest(), HTTPRequest: createHttpRequestWithParams(fakeTenantID, handlers.UserRoleTenantUser, connectorInstanceUrl, "POST")})
	castedCreate = created.(*tenant_provisioning_service_v2.CreateConnectorInstanceV2Forbidden)
	assert.NotNil(t, castedCreate)

	created = handlers.HandleCreateConnectorInstanceV2(handlers.SkylightAndTenantAdminRoles, tenantDB)(tenant_provisioning_service_v2.CreateConnectorInstanceV2Params{Body: generateRandomTenantConnectorInstanceCreationRequest(), HTTPRequest: createHttpRequestWithParams(fakeTenantID, handlers.UserRoleTenantContributor, connectorInstanceUrl, "POST")})
	castedCreate = created.(*tenant_provisioning_service_v2.CreateConnectorInstanceV2Forbidden)
	assert.NotNil(t, castedCreate)

	// Update - SkylightAdmin and TenantAdmin Only
	updated := handlers.HandleUpdateConnectorInstanceV2(handlers.SkylightAndTenantAdminRoles, tenantDB)(tenant_provisioning_service_v2.UpdateConnectorInstanceV2Params{ConnectorInstanceID: fakeID, Body: nil, HTTPRequest: createHttpRequestWithParams(fakeTenantID, handlers.UserRoleUnknown, connectorInstanceUrl, "PATCH")})
	castedUpdate := updated.(*tenant_provisioning_service_v2.UpdateConnectorInstanceV2Forbidden)
	assert.NotNil(t, castedUpdate)

	updated = handlers.HandleUpdateConnectorInstanceV2(handlers.SkylightAndTenantAdminRoles, tenantDB)(tenant_provisioning_service_v2.UpdateConnectorInstanceV2Params{ConnectorInstanceID: fakeID, Body: nil, HTTPRequest: createHttpRequestWithParams(fakeTenantID, handlers.UserRoleTenantUser, connectorInstanceUrl, "PATCH")})
	castedUpdate = updated.(*tenant_provisioning_service_v2.UpdateConnectorInstanceV2Forbidden)
	assert.NotNil(t, castedUpdate)

	updated = handlers.HandleUpdateConnectorInstanceV2(handlers.SkylightAndTenantAdminRoles, tenantDB)(tenant_provisioning_service_v2.UpdateConnectorInstanceV2Params{ConnectorInstanceID: fakeID, Body: nil, HTTPRequest: createHttpRequestWithParams(fakeTenantID, handlers.UserRoleTenantContributor, connectorInstanceUrl, "PATCH")})
	castedUpdate = updated.(*tenant_provisioning_service_v2.UpdateConnectorInstanceV2Forbidden)
	assert.NotNil(t, castedUpdate)

	// Delete - SkylightAdmin and TenantAdmin Only
	deleted := handlers.HandleDeleteConnectorInstanceV2(handlers.SkylightAndTenantAdminRoles, tenantDB)(tenant_provisioning_service_v2.DeleteConnectorInstanceV2Params{ConnectorInstanceID: fakeID, HTTPRequest: createHttpRequestWithParams(fakeTenantID, handlers.UserRoleUnknown, connectorInstanceUrl, "DELETE")})
	castedDelete := deleted.(*tenant_provisioning_service_v2.DeleteConnectorInstanceV2Forbidden)
	assert.NotNil(t, castedDelete)

	deleted = handlers.HandleDeleteConnectorInstanceV2(handlers.SkylightAndTenantAdminRoles, tenantDB)(tenant_provisioning_service_v2.DeleteConnectorInstanceV2Params{ConnectorInstanceID: fakeID, HTTPRequest: createHttpRequestWithParams(fakeTenantID, handlers.UserRoleTenantUser, connectorInstanceUrl, "DELETE")})
	castedDelete = deleted.(*tenant_provisioning_service_v2.DeleteConnectorInstanceV2Forbidden)
	assert.NotNil(t, castedDelete)

	deleted = handlers.HandleDeleteConnectorInstanceV2(handlers.SkylightAndTenantAdminRoles, tenantDB)(tenant_provisioning_service_v2.DeleteConnectorInstanceV2Params{ConnectorInstanceID: fakeID, HTTPRequest: createHttpRequestWithParams(fakeTenantID, handlers.UserRoleTenantContributor, connectorInstanceUrl, "DELETE")})
	castedDelete = deleted.(*tenant_provisioning_service_v2.DeleteConnectorInstanceV2Forbidden)
	assert.NotNil(t, castedDelete)
}

func generateRandomTenantConnectorInstanceCreationRequest() *swagmodels.ConnectorInstanceCreateRequest {
	name := fake.CharactersN(12)
	status := fake.CharactersN(12)

	return &swagmodels.ConnectorInstanceCreateRequest{
		Data: &swagmodels.ConnectorInstanceCreateRequestData{
			Type: &connectorInstanceTypeString,
			Attributes: &swagmodels.ConnectorInstanceCreateRequestDataAttributes{
				Hostname: &name,
				Status:   &status,
			},
		},
	}
}

func generateConnectorInstanceUpdateRequest(id string, rev string, name *string, status *string) *swagmodels.ConnectorInstanceUpdateRequest {
	result := &swagmodels.ConnectorInstanceUpdateRequest{
		Data: &swagmodels.ConnectorInstanceUpdateRequestData{
			Type:       &connectorInstanceTypeString,
			ID:         &id,
			Attributes: &swagmodels.ConnectorInstanceUpdateRequestDataAttributes{Rev: &rev},
		},
	}

	if name != nil {
		result.Data.Attributes.Hostname = *name
	}
	if status != nil {
		result.Data.Attributes.Status = *status
	}

	return result
}
