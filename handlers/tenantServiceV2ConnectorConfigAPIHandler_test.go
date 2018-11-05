package handlers_test

import (
	"math/rand"
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
	connectorConfigUrl = "http://deployment.test.cool/api/v2/connector-configs"

	connectorConfigTypeString = "connectorConfigs"
)

func TestConnectorConfigCRUDV2(t *testing.T) {

	createdTenant := handlers.HandleCreateTenantV2(handlers.SkylightAdminRoleOnly, adminDB, tenantDB)(admin_provisioning_service_v2.CreateTenantV2Params{Body: generateRandomTenantCreationRequest(), HTTPRequest: createHttpRequestWithParams("", handlers.UserRoleSkylight, tenantURL, "POST")})
	castedCreateTeant := createdTenant.(*admin_provisioning_service_v2.CreateTenantV2Created)
	assert.NotNil(t, castedCreateTeant)
	assert.NotEmpty(t, castedCreateTeant.Payload.Data.ID)
	assert.NotEmpty(t, castedCreateTeant.Payload.Data.Attributes.Name)
	assert.NotEmpty(t, castedCreateTeant.Payload.Data.Attributes.URLSubdomain)
	assert.Equal(t, string(common.UserActive), *castedCreateTeant.Payload.Data.Attributes.State)
	assert.True(t, *castedCreateTeant.Payload.Data.Attributes.CreatedTimestamp > 0)
	assert.True(t, *castedCreateTeant.Payload.Data.Attributes.LastModifiedTimestamp > 0)

	// Make sure there are no existing ConnectorConfigs
	existing := handlers.HandleGetAllConnectorConfigsV2(handlers.AllRoles, tenantDB)(tenant_provisioning_service_v2.GetAllConnectorConfigsV2Params{HTTPRequest: createHttpRequestWithParams(*castedCreateTeant.Payload.Data.ID, handlers.UserRoleSkylight, connectorConfigUrl, "GET")})
	castedResponse := existing.(*tenant_provisioning_service_v2.GetAllConnectorConfigsV2NotFound)
	assert.NotNil(t, castedResponse)

	created := handlers.HandleCreateConnectorConfigV2(handlers.AllRoles, tenantDB)(tenant_provisioning_service_v2.CreateConnectorConfigV2Params{Body: generateRandomTenantConnectorConfigCreationRequest(), HTTPRequest: createHttpRequestWithParams(*castedCreateTeant.Payload.Data.ID, handlers.UserRoleSkylight, connectorConfigUrl, "POST")})
	castedCreate := created.(*tenant_provisioning_service_v2.CreateConnectorConfigV2Created)
	assert.NotNil(t, castedCreate)
	assert.NotEmpty(t, castedCreate.Payload.Data.ID)
	assert.NotEmpty(t, castedCreate.Payload.Data.Attributes.Name)
	assert.NotEmpty(t, castedCreate.Payload.Data.Attributes.URL)
	assert.NotEmpty(t, castedCreate.Payload.Data.Attributes.Username)
	assert.NotEmpty(t, castedCreate.Payload.Data.Attributes.Password)
	assert.NotEmpty(t, castedCreate.Payload.Data.Attributes.ExportGroup)
	assert.True(t, castedCreate.Payload.Data.Attributes.Port > 0)
	assert.True(t, *castedCreate.Payload.Data.Attributes.PollingFrequency > 0)
	assert.True(t, castedCreate.Payload.Data.Attributes.DatahubConnectionRetryFrequency > 0)
	assert.True(t, castedCreate.Payload.Data.Attributes.DatahubHeartbeatFrequency > 0)
	assert.NotEmpty(t, castedCreate.Payload.Data.Attributes.Datatype)
	assert.True(t, *castedCreate.Payload.Data.Attributes.CreatedTimestamp > 0)
	assert.True(t, *castedCreate.Payload.Data.Attributes.LastModifiedTimestamp > 0)

	// Make sure we can retrieve this record:
	fetched := handlers.HandleGetConnectorConfigV2(handlers.AllRoles, tenantDB)(tenant_provisioning_service_v2.GetConnectorConfigV2Params{ConnectorID: *castedCreate.Payload.Data.ID, HTTPRequest: createHttpRequestWithParams(*castedCreateTeant.Payload.Data.ID, handlers.UserRoleSkylight, connectorConfigUrl, "GET")})
	castedFetch := fetched.(*tenant_provisioning_service_v2.GetConnectorConfigV2OK)
	assert.NotNil(t, castedFetch)
	assert.Equal(t, castedCreate.Payload.Data, castedFetch.Payload.Data)

	// Also retrieve the record as part of an array
	fetchList := handlers.HandleGetAllConnectorConfigsV2(handlers.AllRoles, tenantDB)(tenant_provisioning_service_v2.GetAllConnectorConfigsV2Params{HTTPRequest: createHttpRequestWithParams(*castedCreateTeant.Payload.Data.ID, handlers.UserRoleSkylight, connectorConfigUrl, "GET")})
	castedFetchList := fetchList.(*tenant_provisioning_service_v2.GetAllConnectorConfigsV2OK)
	assert.NotNil(t, castedFetchList)
	assert.Equal(t, 1, len(castedFetchList.Payload.Data))
	assert.Equal(t, castedCreate.Payload.Data, castedFetchList.Payload.Data[0])

	// Make an update to the Record
	newName := fake.CharactersN(16)
	updateRequestBody := generateConnectorConfigUpdateRequest(*castedCreate.Payload.Data.ID, *castedCreate.Payload.Data.Attributes.Rev, &newName, nil, nil, nil, nil, nil, nil, nil, nil, nil)
	updated := handlers.HandleUpdateConnectorConfigV2(handlers.AllRoles, tenantDB)(tenant_provisioning_service_v2.UpdateConnectorConfigV2Params{ConnectorID: *castedCreate.Payload.Data.ID, Body: updateRequestBody, HTTPRequest: createHttpRequestWithParams(*castedCreateTeant.Payload.Data.ID, handlers.UserRoleSkylight, connectorConfigUrl, "PATCH")})
	castedUpdate := updated.(*tenant_provisioning_service_v2.UpdateConnectorConfigV2OK)
	assert.NotNil(t, castedUpdate)
	assert.NotEqual(t, castedCreate.Payload.Data, castedUpdate.Payload.Data)
	assert.NotEqual(t, castedCreate.Payload.Data.Attributes.Rev, castedUpdate.Payload.Data.Attributes.Rev)
	assert.Equal(t, newName, castedUpdate.Payload.Data.Attributes.Name)
	assert.Equal(t, castedCreate.Payload.Data.Attributes.URL, castedUpdate.Payload.Data.Attributes.URL)
	assert.Equal(t, castedCreate.Payload.Data.Attributes.Port, castedUpdate.Payload.Data.Attributes.Port)

	// Delete the record
	deleted := handlers.HandleDeleteConnectorConfigV2(handlers.AllRoles, tenantDB)(tenant_provisioning_service_v2.DeleteConnectorConfigV2Params{ConnectorID: *castedCreate.Payload.Data.ID, HTTPRequest: createHttpRequestWithParams(*castedCreateTeant.Payload.Data.ID, handlers.UserRoleSkylight, connectorConfigUrl, "DELETE")})
	castedDelete := deleted.(*tenant_provisioning_service_v2.DeleteConnectorConfigV2OK)
	assert.NotNil(t, castedDelete)
	assert.Equal(t, castedUpdate.Payload.Data, castedDelete.Payload.Data)

	// Make sure there are no existing ConnectorConfigs
	existing = handlers.HandleGetAllConnectorConfigsV2(handlers.AllRoles, tenantDB)(tenant_provisioning_service_v2.GetAllConnectorConfigsV2Params{HTTPRequest: createHttpRequestWithParams(*castedCreateTeant.Payload.Data.ID, handlers.UserRoleSkylight, connectorConfigUrl, "GET")})
	castedResponse = existing.(*tenant_provisioning_service_v2.GetAllConnectorConfigsV2NotFound)
	assert.NotNil(t, castedResponse)
}

func TestConnectorConfigNotFoundV2(t *testing.T) {
	notFoundID := fake.CharactersN(20)

	// Get ConnectorConfig
	fetched := handlers.HandleGetConnectorConfigV2(handlers.AllRoles, tenantDB)(tenant_provisioning_service_v2.GetConnectorConfigV2Params{ConnectorID: notFoundID, HTTPRequest: createHttpRequestWithParams("", handlers.UserRoleSkylight, connectorConfigUrl, "GET")})
	castedFetch := fetched.(*tenant_provisioning_service_v2.GetConnectorConfigV2NotFound)
	assert.NotNil(t, castedFetch)

	// Delete ConnectorConfig
	deleted := handlers.HandleDeleteConnectorConfigV2(handlers.AllRoles, tenantDB)(tenant_provisioning_service_v2.DeleteConnectorConfigV2Params{ConnectorID: notFoundID, HTTPRequest: createHttpRequestWithParams("", handlers.UserRoleSkylight, connectorConfigUrl, "DELETE")})
	castedDelete := deleted.(*tenant_provisioning_service_v2.DeleteConnectorConfigV2NotFound)
	assert.NotNil(t, castedDelete)

	// Patch ConnectorConfig
	updateRequest := generateConnectorConfigUpdateRequest(notFoundID, "reviosionstuff", nil, nil, nil, nil, nil, nil, nil, nil, nil, nil)
	updated := handlers.HandleUpdateConnectorConfigV2(handlers.AllRoles, tenantDB)(tenant_provisioning_service_v2.UpdateConnectorConfigV2Params{ConnectorID: notFoundID, Body: updateRequest, HTTPRequest: createHttpRequestWithParams("", handlers.UserRoleSkylight, connectorConfigUrl, "PATCH")})
	castedUpdate := updated.(*tenant_provisioning_service_v2.UpdateConnectorConfigV2NotFound)
	assert.NotNil(t, castedUpdate)
}

func TestConnectorConfigBadRequestV2(t *testing.T) {

	// CreateConnectorConfig
	created := handlers.HandleCreateConnectorConfigV2(handlers.AllRoles, tenantDB)(tenant_provisioning_service_v2.CreateConnectorConfigV2Params{Body: nil, HTTPRequest: createHttpRequestWithParams("", handlers.UserRoleSkylight, connectorConfigUrl, "POST")})
	castedCreate := created.(*tenant_provisioning_service_v2.CreateConnectorConfigV2BadRequest)
	assert.NotNil(t, castedCreate)

	// Update ConnectorConfig
	updated := handlers.HandleUpdateConnectorConfigV2(handlers.AllRoles, tenantDB)(tenant_provisioning_service_v2.UpdateConnectorConfigV2Params{ConnectorID: fake.CharactersN(20), Body: nil, HTTPRequest: createHttpRequestWithParams("", handlers.UserRoleSkylight, connectorConfigUrl, "PATCH")})
	castedUpdate := updated.(*tenant_provisioning_service_v2.UpdateConnectorConfigV2BadRequest)
	assert.NotNil(t, castedUpdate)
}

func TestConnectorConfigConflictV2(t *testing.T) {

	createdTenant := handlers.HandleCreateTenantV2(handlers.SkylightAdminRoleOnly, adminDB, tenantDB)(admin_provisioning_service_v2.CreateTenantV2Params{Body: generateRandomTenantCreationRequest(), HTTPRequest: createHttpRequestWithParams("", handlers.UserRoleSkylight, tenantURL, "POST")})
	castedCreateTeant := createdTenant.(*admin_provisioning_service_v2.CreateTenantV2Created)
	assert.NotNil(t, castedCreateTeant)
	assert.NotEmpty(t, castedCreateTeant.Payload.Data.ID)
	assert.NotEmpty(t, castedCreateTeant.Payload.Data.Attributes.Name)
	assert.NotEmpty(t, castedCreateTeant.Payload.Data.Attributes.URLSubdomain)
	assert.Equal(t, string(common.UserActive), *castedCreateTeant.Payload.Data.Attributes.State)
	assert.True(t, *castedCreateTeant.Payload.Data.Attributes.CreatedTimestamp > 0)
	assert.True(t, *castedCreateTeant.Payload.Data.Attributes.LastModifiedTimestamp > 0)

	// Make sure there are no existing ConnectorConfig
	existing := handlers.HandleGetAllConnectorConfigsV2(handlers.AllRoles, tenantDB)(tenant_provisioning_service_v2.GetAllConnectorConfigsV2Params{HTTPRequest: createHttpRequestWithParams(*castedCreateTeant.Payload.Data.ID, handlers.UserRoleSkylight, connectorConfigUrl, "GET")})
	castedResponse := existing.(*tenant_provisioning_service_v2.GetAllConnectorConfigsV2NotFound)
	assert.NotNil(t, castedResponse)

	createReqBody := generateRandomTenantConnectorConfigCreationRequest()
	created := handlers.HandleCreateConnectorConfigV2(handlers.AllRoles, tenantDB)(tenant_provisioning_service_v2.CreateConnectorConfigV2Params{Body: createReqBody, HTTPRequest: createHttpRequestWithParams(*castedCreateTeant.Payload.Data.ID, handlers.UserRoleSkylight, connectorConfigUrl, "POST")})
	castedCreate := created.(*tenant_provisioning_service_v2.CreateConnectorConfigV2Created)
	assert.NotNil(t, castedCreate)
	assert.NotEmpty(t, castedCreate.Payload.Data.ID)
	assert.NotEmpty(t, castedCreate.Payload.Data.Attributes.Name)
	assert.NotEmpty(t, castedCreate.Payload.Data.Attributes.URL)
	assert.NotEmpty(t, castedCreate.Payload.Data.Attributes.Username)
	assert.NotEmpty(t, castedCreate.Payload.Data.Attributes.Password)
	assert.NotEmpty(t, castedCreate.Payload.Data.Attributes.ExportGroup)
	assert.True(t, castedCreate.Payload.Data.Attributes.Port > 0)
	assert.True(t, *castedCreate.Payload.Data.Attributes.PollingFrequency > 0)
	assert.True(t, castedCreate.Payload.Data.Attributes.DatahubConnectionRetryFrequency > 0)
	assert.True(t, castedCreate.Payload.Data.Attributes.DatahubHeartbeatFrequency > 0)
	assert.NotEmpty(t, castedCreate.Payload.Data.Attributes.Datatype)
	assert.True(t, *castedCreate.Payload.Data.Attributes.CreatedTimestamp > 0)
	assert.True(t, *castedCreate.Payload.Data.Attributes.LastModifiedTimestamp > 0)

	// Try to create the record again - should succeed as we are not guarding against name collisiones
	createdConflict := handlers.HandleCreateConnectorConfigV2(handlers.AllRoles, tenantDB)(tenant_provisioning_service_v2.CreateConnectorConfigV2Params{Body: createReqBody, HTTPRequest: createHttpRequestWithParams(*castedCreateTeant.Payload.Data.ID, handlers.UserRoleSkylight, connectorConfigUrl, "POST")})
	castedCreateConflictButOK := createdConflict.(*tenant_provisioning_service_v2.CreateConnectorConfigV2Created)
	assert.NotNil(t, castedCreateConflictButOK)

	// Try the update with a bad revision
	newName := fake.CharactersN(16)
	updateRequestBody := generateConnectorConfigUpdateRequest(*castedCreate.Payload.Data.ID, *castedCreate.Payload.Data.Attributes.Rev+"pork", &newName, nil, nil, nil, nil, nil, nil, nil, nil, nil)
	updated := handlers.HandleUpdateConnectorConfigV2(handlers.AllRoles, tenantDB)(tenant_provisioning_service_v2.UpdateConnectorConfigV2Params{ConnectorID: *castedCreate.Payload.Data.ID, Body: updateRequestBody, HTTPRequest: createHttpRequestWithParams(*castedCreateTeant.Payload.Data.ID, handlers.UserRoleSkylight, connectorConfigUrl, "PATCH")})
	castedUpdate := updated.(*tenant_provisioning_service_v2.UpdateConnectorConfigV2Conflict)
	assert.NotNil(t, castedUpdate)

	// Delete the tenant
	deleted := handlers.HandleDeleteConnectorConfigV2(handlers.AllRoles, tenantDB)(tenant_provisioning_service_v2.DeleteConnectorConfigV2Params{ConnectorID: *castedCreate.Payload.Data.ID, HTTPRequest: createHttpRequestWithParams(*castedCreateTeant.Payload.Data.ID, handlers.UserRoleSkylight, connectorConfigUrl, "DELETE")})
	castedDelete := deleted.(*tenant_provisioning_service_v2.DeleteConnectorConfigV2OK)
	assert.NotNil(t, castedDelete)
	assert.NotNil(t, castedDelete.Payload.Data)
}

func TestConnectorConfigAPIsProtectedByAuthV2(t *testing.T) {
	fakeTenantID := fake.CharactersN(20)
	// Get All - All Users
	existing := handlers.HandleGetAllConnectorConfigsV2(handlers.AllRoles, tenantDB)(tenant_provisioning_service_v2.GetAllConnectorConfigsV2Params{HTTPRequest: createHttpRequestWithParams(fakeTenantID, handlers.UserRoleUnknown, connectorConfigUrl, "GET")})
	castedResponse := existing.(*tenant_provisioning_service_v2.GetAllConnectorConfigsV2Forbidden)
	assert.NotNil(t, castedResponse)

	fakeID := fake.CharactersN(20)

	// Get - All Users
	fetched := handlers.HandleGetConnectorConfigV2(handlers.AllRoles, tenantDB)(tenant_provisioning_service_v2.GetConnectorConfigV2Params{ConnectorID: fakeID, HTTPRequest: createHttpRequestWithParams(fakeTenantID, handlers.UserRoleUnknown, connectorConfigUrl, "GET")})
	castedFetch := fetched.(*tenant_provisioning_service_v2.GetConnectorConfigV2Forbidden)
	assert.NotNil(t, castedFetch)

	// Create - SkylightAdmin and TenantAdmin Only
	created := handlers.HandleCreateConnectorConfigV2(handlers.SkylightAndTenantAdminRoles, tenantDB)(tenant_provisioning_service_v2.CreateConnectorConfigV2Params{Body: generateRandomTenantConnectorConfigCreationRequest(), HTTPRequest: createHttpRequestWithParams(fakeTenantID, handlers.UserRoleUnknown, connectorConfigUrl, "POST")})
	castedCreate := created.(*tenant_provisioning_service_v2.CreateConnectorConfigV2Forbidden)
	assert.NotNil(t, castedCreate)

	created = handlers.HandleCreateConnectorConfigV2(handlers.SkylightAndTenantAdminRoles, tenantDB)(tenant_provisioning_service_v2.CreateConnectorConfigV2Params{Body: generateRandomTenantConnectorConfigCreationRequest(), HTTPRequest: createHttpRequestWithParams(fakeTenantID, handlers.UserRoleTenantUser, connectorConfigUrl, "POST")})
	castedCreate = created.(*tenant_provisioning_service_v2.CreateConnectorConfigV2Forbidden)
	assert.NotNil(t, castedCreate)

	// Update - SkylightAdmin and TenantAdmin Only
	updated := handlers.HandleUpdateConnectorConfigV2(handlers.SkylightAndTenantAdminRoles, tenantDB)(tenant_provisioning_service_v2.UpdateConnectorConfigV2Params{ConnectorID: fakeID, Body: nil, HTTPRequest: createHttpRequestWithParams(fakeTenantID, handlers.UserRoleUnknown, connectorConfigUrl, "PATCH")})
	castedUpdate := updated.(*tenant_provisioning_service_v2.UpdateConnectorConfigV2Forbidden)
	assert.NotNil(t, castedUpdate)

	updated = handlers.HandleUpdateConnectorConfigV2(handlers.SkylightAndTenantAdminRoles, tenantDB)(tenant_provisioning_service_v2.UpdateConnectorConfigV2Params{ConnectorID: fakeID, Body: nil, HTTPRequest: createHttpRequestWithParams(fakeTenantID, handlers.UserRoleTenantUser, connectorConfigUrl, "PATCH")})
	castedUpdate = updated.(*tenant_provisioning_service_v2.UpdateConnectorConfigV2Forbidden)
	assert.NotNil(t, castedUpdate)

	// Delete - SkylightAdmin and TenantAdmin Only
	deleted := handlers.HandleDeleteConnectorConfigV2(handlers.SkylightAndTenantAdminRoles, tenantDB)(tenant_provisioning_service_v2.DeleteConnectorConfigV2Params{ConnectorID: fakeID, HTTPRequest: createHttpRequestWithParams(fakeTenantID, handlers.UserRoleUnknown, connectorConfigUrl, "DELETE")})
	castedDelete := deleted.(*tenant_provisioning_service_v2.DeleteConnectorConfigV2Forbidden)
	assert.NotNil(t, castedDelete)

	deleted = handlers.HandleDeleteConnectorConfigV2(handlers.SkylightAndTenantAdminRoles, tenantDB)(tenant_provisioning_service_v2.DeleteConnectorConfigV2Params{ConnectorID: fakeID, HTTPRequest: createHttpRequestWithParams(fakeTenantID, handlers.UserRoleTenantUser, connectorConfigUrl, "DELETE")})
	castedDelete = deleted.(*tenant_provisioning_service_v2.DeleteConnectorConfigV2Forbidden)
	assert.NotNil(t, castedDelete)
}

func generateRandomTenantConnectorConfigCreationRequest() *swagmodels.ConnectorConfigCreateRequest {
	name := fake.CharactersN(12)
	username := fake.CharactersN(12)
	password := fake.CharactersN(12)
	exportGroup := fake.CharactersN(12)
	someType := fake.CharactersN(12)
	url := fake.DomainName()
	port := int64(rand.Intn(50000))
	pollingFrequency := int64(rand.Intn(50000))
	datahubHeartbeatFrequency := int64(rand.Intn(50000))
	datahubConnectionRetryFrequency := int64(rand.Intn(50000))

	return &swagmodels.ConnectorConfigCreateRequest{
		Data: &swagmodels.ConnectorConfigCreateRequestData{
			Type: &connectorConfigTypeString,
			Attributes: &swagmodels.ConnectorConfigCreateRequestDataAttributes{
				Name:                            name,
				URL:                             &url,
				Username:                        username,
				Password:                        password,
				ExportGroup:                     exportGroup,
				Type:                            &someType,
				Port:                            port,
				PollingFrequency:                &pollingFrequency,
				DatahubHeartbeatFrequency:       datahubHeartbeatFrequency,
				DatahubConnectionRetryFrequency: datahubConnectionRetryFrequency,
			},
		},
	}
}

func generateConnectorConfigUpdateRequest(id string, rev string, name *string, url *string, username *string, password *string, exportGroup *string, someType *string, port *int64, pollFreq *int64, hbFreq *int64, crFreq *int64) *swagmodels.ConnectorConfigUpdateRequest {
	result := &swagmodels.ConnectorConfigUpdateRequest{
		Data: &swagmodels.ConnectorConfigUpdateRequestData{
			Type:       &connectorConfigTypeString,
			ID:         &id,
			Attributes: &swagmodels.ConnectorConfigUpdateRequestDataAttributes{Rev: &rev},
		},
	}

	if name != nil {
		result.Data.Attributes.Name = *name
	}
	if url != nil {
		result.Data.Attributes.URL = *url
	}
	if username != nil {
		result.Data.Attributes.Username = *username
	}
	if password != nil {
		result.Data.Attributes.Password = *password
	}
	if exportGroup != nil {
		result.Data.Attributes.ExportGroup = *exportGroup
	}
	if someType != nil {
		result.Data.Attributes.Type = *someType
	}

	if port != nil {
		result.Data.Attributes.Port = *port
	}
	if pollFreq != nil {
		result.Data.Attributes.Port = *pollFreq
	}
	if hbFreq != nil {
		result.Data.Attributes.Port = *hbFreq
	}
	if crFreq != nil {
		result.Data.Attributes.Port = *crFreq
	}
	return result
}
