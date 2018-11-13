package tenant

import (
	"testing"

	ds "github.com/accedian/adh-gather/datastore"
	admmod "github.com/accedian/adh-gather/models/admin"
	"github.com/accedian/adh-gather/models/common"
	tenmod "github.com/accedian/adh-gather/models/tenant"
	"github.com/stretchr/testify/assert"
)

func (runner *TenantServiceDatastoreTestRunner) RunTenantMetricBaselineCRUD(t *testing.T) {
	const COMPANY1 = "MetricBaselineCompany"
	const SUBDOMAIN1 = "subdom1"
	const MONOBJ1 = "MONOBJ1"
	const MONOBJ2 = "MONOBJ2"
	const DELAYP95 = "delayP95"

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

	// Try to fetch a record even though none exist:
	fail, err := runner.tenantDB.GetMetricBaseline(TENANT, "someID")
	assert.NotNil(t, err)
	assert.Nil(t, fail)

	fail, err = runner.tenantDB.GetMetricBaseline(TENANT, "someID")
	assert.NotNil(t, err)
	assert.Nil(t, fail)

	failArray, err := runner.tenantDB.GetMetricBaselineForMonitoredObjectForHourOfWeek(TENANT, "someID", int32(1))
	assert.NotNil(t, err)
	assert.Nil(t, failArray)

	// Try to Update a record that does not exist:
	fail, err = runner.tenantDB.UpdateMetricBaseline(&tenmod.MetricBaseline{})
	assert.NotNil(t, err)
	assert.Nil(t, fail)

	// Try to update an hour in a record that does not exist, should creatre the record
	upsert, err := runner.tenantDB.UpdateMetricBaselineForHourOfWeek(TENANT, MONOBJ1, &tenmod.MetricBaselineData{Metric: DELAYP95, Direction: "0", HourOfWeek: 150})
	assert.Nil(t, err)
	assert.NotNil(t, upsert)
	assert.NotEmpty(t, upsert.ID)
	assert.NotEmpty(t, upsert.REV)
	assert.Equal(t, string(tenmod.TenantMetricBaselineType), upsert.Datatype)
	assert.Equal(t, DELAYP95, upsert.Baselines[0].Metric, "Metric not the same")
	assert.Equal(t, "0", upsert.Baselines[0].Direction, "Direction not the same")
	assert.Equal(t, int32(150), upsert.Baselines[0].HourOfWeek, "HourofWeek not the same")
	assert.Equal(t, TENANT, upsert.TenantID, "Tenant ID not the same")
	assert.True(t, upsert.CreatedTimestamp > 0, "CreatedTimestamp was not set")
	assert.True(t, upsert.LastModifiedTimestamp > 0, "LastmodifiedTimestamp was not set")

	// Try to create the same record, should fail
	createObj := tenmod.MetricBaseline{
		Datatype:          string(tenmod.TenantMetricBaselineType),
		TenantID:          TENANT,
		MonitoredObjectID: MONOBJ1,
	}
	failedCreate, err := runner.tenantDB.CreateMetricBaseline(&createObj)
	assert.NotNil(t, err)
	assert.Nil(t, failedCreate)

	// Create a record successfully
	metricBaseline1 := tenmod.MetricBaseline{
		Datatype:          string(tenmod.TenantMetricBaselineType),
		TenantID:          TENANT,
		MonitoredObjectID: MONOBJ2,
		Baselines:         []*tenmod.MetricBaselineData{},
	}
	created, err := runner.tenantDB.CreateMetricBaseline(&metricBaseline1)
	assert.Nil(t, err)
	assert.NotNil(t, created)
	assert.NotEmpty(t, created.ID)
	assert.NotEmpty(t, created.REV)
	assert.Equal(t, string(tenmod.TenantMetricBaselineType), created.Datatype)
	assert.Equal(t, TENANT, upsert.TenantID, "Tenant ID not the same")
	assert.Empty(t, created.Baselines, "Baseline array should be empty")
	assert.True(t, created.CreatedTimestamp > 0, "CreatedTimestamp was not set")
	assert.True(t, created.LastModifiedTimestamp > 0, "LastmodifiedTimestamp was not set")

	// Get a record
	fetched, err := runner.tenantDB.GetMetricBaseline(TENANT, created.ID)
	assert.Nil(t, err)
	assert.ElementsMatch(t, created.Baselines, fetched.Baselines, "The retrieved record should have the same baselines as the created record")
	assert.Equal(t, created.MonitoredObjectID, fetched.MonitoredObjectID, "Monitored object not the same")

	// Get a record by monitored object ID
	fetched, err = runner.tenantDB.GetMetricBaseline(TENANT, created.MonitoredObjectID)
	assert.Nil(t, err)
	assert.ElementsMatch(t, created.Baselines, fetched.Baselines, "The retrieved record should have the same baselines as the created record")
	assert.Equal(t, created.MonitoredObjectID, fetched.MonitoredObjectID, "Monitored object not the same")

	// Add new baseline data to existing record
	updated, err := runner.tenantDB.UpdateMetricBaselineForHourOfWeek(TENANT, upsert.MonitoredObjectID, &tenmod.MetricBaselineData{Metric: DELAYP95, Direction: "1", HourOfWeek: 250})
	assert.Nil(t, err)
	assert.NotNil(t, updated)
	assert.Equal(t, upsert.ID, updated.ID, "ID not the same")
	assert.Equal(t, TENANT, updated.TenantID, "Tenant ID not the same")
	assert.Equal(t, 2, len(updated.Baselines), "Baseline array should have 2 elements")

	// add another bit of data for hour 150
	_, err = runner.tenantDB.UpdateMetricBaselineForHourOfWeek(TENANT, upsert.MonitoredObjectID, &tenmod.MetricBaselineData{Metric: DELAYP95, Direction: "1", HourOfWeek: 150})
	assert.Nil(t, err)

	// Get baselines for an hour of the week - success
	baselineArray, err := runner.tenantDB.GetMetricBaselineForMonitoredObjectForHourOfWeek(TENANT, upsert.MonitoredObjectID, 150)
	assert.Nil(t, err)
	assert.NotNil(t, baselineArray)
	assert.Equal(t, 2, len(baselineArray))

	baselineArray, err = runner.tenantDB.GetMetricBaselineForMonitoredObjectForHourOfWeek(TENANT, upsert.MonitoredObjectID, 250)
	assert.Nil(t, err)
	assert.NotNil(t, baselineArray)
	assert.Equal(t, 1, len(baselineArray))

	// Update an entire record
	fetched, err = runner.tenantDB.GetMetricBaseline(TENANT, upsert.MonitoredObjectID)
	assert.Nil(t, err)
	assert.NotNil(t, fetched)
	assert.Equal(t, 3, len(fetched.Baselines))

	fetched.Baselines = []*tenmod.MetricBaselineData{}
	fetched.ID = ds.GetDataIDFromFullID(fetched.ID)
	replaced, err := runner.tenantDB.UpdateMetricBaseline(fetched)
	assert.Nil(t, err)
	assert.NotNil(t, replaced)
	assert.Equal(t, 0, len(replaced.Baselines))

	// Now create a monitored Object that is linked to a metric baseline so that it can be deleted and make sure the baseline is deleted as well
	tenantMonObj := tenmod.MonitoredObject{
		MonitoredObjectID: MONOBJ1,
		ObjectName:        "Glummstein",
		TenantID:          TENANT,
	}
	createdMO, err := runner.tenantDB.CreateMonitoredObject(&tenantMonObj)
	assert.Nil(t, err)
	assert.NotNil(t, createdMO)

	deletedMO, err := runner.tenantDB.DeleteMonitoredObject(TENANT, ds.GetDataIDFromFullID(createdMO.ID))
	assert.Nil(t, err)
	assert.NotNil(t, deletedMO)

	fetched, err = runner.tenantDB.GetMetricBaseline(TENANT, MONOBJ1)
	assert.NotNil(t, err)
	assert.Nil(t, fetched)

	// Now delete a metric baseline successfully by ID
	deleted, err := runner.tenantDB.DeleteMetricBaseline(TENANT, ds.GetDataIDFromFullID(created.ID))
	assert.Nil(t, err)
	assert.NotNil(t, deleted)

	fetched, err = runner.tenantDB.GetMetricBaseline(TENANT, created.MonitoredObjectID)
	assert.NotNil(t, err)
	assert.Nil(t, fetched)
}
