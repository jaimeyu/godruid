package tenant

import (
	"testing"
	"time"

	ds "github.com/accedian/adh-gather/datastore"
	admmod "github.com/accedian/adh-gather/models/admin"
	"github.com/accedian/adh-gather/models/common"
	tenmod "github.com/accedian/adh-gather/models/tenant"
	"github.com/getlantern/deepcopy"
	"github.com/stretchr/testify/assert"
)

func (runner *TenantServiceDatastoreTestRunner) RunTenantMetadataConfigCRUD(t *testing.T) {
	const COMPANY1 = "MetadataConfigCompany"
	const SUBDOMAIN1 = "subdom1"

	// Create a tenant
	data := admmod.Tenant{
		Name:         COMPANY1,
		URLSubdomain: SUBDOMAIN1,
		State:        string(common.UserActive)}
	tenantDescriptor, err := runner.adminDB.CreateTenant(&data)
	assert.Nil(t, err)
	assert.NotNil(t, tenantDescriptor)
	assert.Equal(t, COMPANY1, tenantDescriptor.Name)

	TENANT := ds.GetDataIDFromFullID(tenantDescriptor.ID)

	// Validate that there are currently no records
	rec, err := runner.tenantDB.GetActiveTenantMetadataConfig(TENANT)
	assert.NotNil(t, err)
	assert.Nil(t, rec)

	// Try to fetch a record even though none exist:
	fail, err := runner.tenantDB.GetTenantMetadataConfig(TENANT, "someID")
	assert.NotNil(t, err)
	assert.Nil(t, fail)

	// Try to Update a record that does not exist:
	fail, err = runner.tenantDB.UpdateTenantMetadataConfig(&tenmod.MetadataConfig{})
	assert.NotNil(t, err)
	assert.Nil(t, fail)

	// Create a record
	metaCfg := tenmod.MetadataConfig{
		Datatype:   string(tenmod.TenantMetadataConfigType),
		TenantID:   TENANT,
		StartPoint: "Start",
		EndPoint:   "End",
		MidPoints:  []string{},
	}
	created, err := runner.tenantDB.CreateTenantMetadataConfig(&metaCfg)
	assert.Nil(t, err)
	assert.NotNil(t, created)
	assert.NotEmpty(t, created.ID)
	assert.NotEmpty(t, created.REV)
	assert.Equal(t, string(tenmod.TenantMetadataConfigType), created.Datatype)
	assert.Equal(t, metaCfg.StartPoint, created.StartPoint, "Start point not the same")
	assert.Equal(t, metaCfg.EndPoint, created.EndPoint, "End point point not the same")
	assert.ElementsMatch(t, metaCfg.MidPoints, created.MidPoints, "Mid points not the same")
	assert.Equal(t, created.TenantID, TENANT, "Tenant ID not the same")
	assert.True(t, created.CreatedTimestamp > 0, "CreatedTimestamp was not set")
	assert.True(t, created.LastModifiedTimestamp > 0, "LastmodifiedTimestamp was not set")

	// Get a record
	fetched, err := runner.tenantDB.GetTenantMetadataConfig(TENANT, created.ID)
	assert.Nil(t, err)
	assert.ElementsMatch(t, created.MidPoints, fetched.MidPoints, "The retrieved record should have the same mid pointssame as the created record")
	assert.Equal(t, created.StartPoint, fetched.StartPoint, "Start point not the same")
	assert.Equal(t, created.EndPoint, fetched.EndPoint, "End point point not the same")

	time.Sleep(time.Millisecond * 2)

	// Update a record
	updateRecord := tenmod.MetadataConfig{}
	deepcopy.Copy(&updateRecord, fetched)
	updateRecord.StartPoint = "nerd"
	updated, err := runner.tenantDB.UpdateTenantMetadataConfig(&updateRecord)
	assert.Nil(t, err)
	assert.NotNil(t, updated)
	assert.Equal(t, updated.ID, fetched.ID)
	assert.NotEqual(t, updated.REV, fetched.REV)
	assert.Equal(t, string(tenmod.TenantMetadataConfigType), updated.Datatype)
	assert.Equal(t, updated.TenantID, TENANT, "Tenant ID not the same")
	assert.NotEqual(t, fetched.StartPoint, updated.StartPoint, "Start point was not updated")
	assert.ElementsMatch(t, fetched.MidPoints, updated.MidPoints, "The retrieved record should have the same mid pointssame as the created record")
	assert.Equal(t, fetched.EndPoint, updated.EndPoint, "End point point not the same")
	assert.Equal(t, updated.CreatedTimestamp, fetched.CreatedTimestamp, "CreatedTimestamp should not be updated")
	assert.True(t, updated.LastModifiedTimestamp > fetched.LastModifiedTimestamp, "LastmodifiedTimestamp was not updated")

	// Add a second record - should fail.
	tenantMetaCfg2 := tenmod.MetadataConfig{
		Datatype:   string(tenmod.TenantMetadataConfigType),
		TenantID:   TENANT,
		StartPoint: "who",
		EndPoint:   "Cares",
	}
	created2, err := runner.tenantDB.CreateTenantMetadataConfig(&tenantMetaCfg2)
	assert.NotNil(t, err)
	assert.Nil(t, created2)

	// Get active records
	active, err := runner.tenantDB.GetActiveTenantMetadataConfig(TENANT)
	assert.Nil(t, err)
	assert.NotNil(t, active)
	assert.Equal(t, updated.StartPoint, active.StartPoint, "Start point not the same")
	assert.Equal(t, updated.EndPoint, active.EndPoint, "End point point not the same")
	assert.ElementsMatch(t, updated.MidPoints, active.MidPoints, "Mid points not the same")

	// Delete a record.
	deleted, err := runner.tenantDB.DeleteTenantMetadataConfig(TENANT, active.ID)
	assert.Nil(t, err)
	assert.NotNil(t, deleted)
	assert.NotEmpty(t, deleted.ID)
	assert.NotEmpty(t, deleted.REV)
	assert.Equal(t, updated.StartPoint, deleted.StartPoint, "Start point not the same")
	assert.Equal(t, updated.EndPoint, deleted.EndPoint, "End point point not the same")
	assert.ElementsMatch(t, updated.MidPoints, deleted.MidPoints, "Mid points not the same")

	// Delete a record that oes not exist
	deleteDNE, err := runner.tenantDB.DeleteTenantDomain(TENANT, deleted.ID)
	assert.NotNil(t, err)
	assert.Nil(t, deleteDNE)
}
