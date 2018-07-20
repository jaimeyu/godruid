package handlers_test

import (
	"os"
	"testing"

	"github.com/accedian/adh-gather/swagmodels"

	"github.com/accedian/adh-gather/handlers"
	"github.com/accedian/adh-gather/logger"
	"github.com/accedian/adh-gather/restapi/operations/tenant_provisioning_service_v2"
	"github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {

	err := setupTestDatastore()
	if err != nil {
		logger.Log.Fatalf("Unable to setup datastore for Data Cleaning Profile tests: %s", err.Error())
	}

	code := m.Run()

	err = destroyTestDatastore()
	if err != nil {
		logger.Log.Errorf("Unable to remove test datastore for Data Cleaning Profile tests: %s", err.Error())
	}

	os.Exit(code)
}

func TestGetDataCleaningProfile(t *testing.T) {

	tenantDescriptor1 := getRandomTenantDescriptor()

	createdTenant, err := adminDB.CreateTenant(tenantDescriptor1)
	assert.Nil(t, err)
	assert.NotNil(t, createdTenant)
	assert.Equal(t, tenantDescriptor1.Name, createdTenant.Name)
	assert.Equal(t, tenantDescriptor1.URLSubdomain, createdTenant.URLSubdomain)
	assert.True(t, len(createdTenant.ID) > 0)

	// Make sure there are no existing Data Cleaning Profiles
	existing := handlers.HandleGetDataCleaningProfilesV2(handlers.SkylightAndTenantAdminRoles, tenantDB)(tenant_provisioning_service_v2.GetDataCleaningProfilesParams{HTTPRequest: createHttpRequest(createdTenant.ID, handlers.UserRoleSkylight)})
	castedResponse := existing.(*tenant_provisioning_service_v2.GetDataCleaningProfilesNotFound)
	assert.NotNil(t, castedResponse)

	// Create a record
	created := handlers.HandleCreateDataCleaningProfileV2(handlers.SkylightAndTenantAdminRoles, tenantDB)(tenant_provisioning_service_v2.CreateDataCleaningProfileParams{
		Body: &swagmodels.DataCleaningProfileCreateRequest{
			Data: createRandomDataCleaningProfile(),
		},
		HTTPRequest: createHttpRequest(createdTenant.ID, handlers.UserRoleSkylight)})
	castedCreate := created.(*tenant_provisioning_service_v2.CreateDataCleaningProfileCreated)
	assert.NotNil(t, castedCreate)
	assert.True(t, len(*castedCreate.Payload.Data.ID) != 0)

	// Make sure we can retrieve the record:
	fetched := handlers.HandleGetDataCleaningProfileV2(handlers.SkylightAndTenantAdminRoles, tenantDB)(tenant_provisioning_service_v2.GetDataCleaningProfileParams{
		ProfileID:   *castedCreate.Payload.Data.ID,
		HTTPRequest: createHttpRequest(createdTenant.ID, handlers.UserRoleSkylight)})
	castedFetch := fetched.(*tenant_provisioning_service_v2.GetDataCleaningProfileOK)
	assert.NotNil(t, castedFetch)
	assert.Equal(t, castedCreate.Payload.Data, castedFetch.Payload.Data)

	adminDB.DeleteTenant(createdTenant.ID)

}
