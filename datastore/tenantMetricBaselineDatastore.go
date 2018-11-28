package datastore

import (
	"fmt"

	"github.com/accedian/adh-gather/models/common"
	tenmod "github.com/accedian/adh-gather/models/tenant"
	"github.com/accedian/adh-gather/swagmodels"
)

type TenantMetricBaselineDatastore interface {
	CreateMetricBaseline(baseline *tenmod.MetricBaseline) (*tenmod.MetricBaseline, error)
	UpdateMetricBaseline(baseline *tenmod.MetricBaseline) (*tenmod.MetricBaseline, error)
	UpdateMetricBaselineForHourOfWeek(tenantID string, monObjID string, hourOfWeek int32, baselineData *tenmod.MetricBaselineData) (*tenmod.MetricBaseline, error)
	UpdateMetricBaselineForHourOfWeekWithCollection(tenantID string, monObjID string, hourOfWeek int32, baselineDataCollection []*tenmod.MetricBaselineData) (*tenmod.MetricBaseline, error)
	GetMetricBaseline(tenantID string, dataID string) (*tenmod.MetricBaseline, error)
	GetMetricBaselineForMonitoredObject(tenantID string, monObjID string) ([]*tenmod.MetricBaseline, error)
	DeleteMetricBaselineForMonitoredObject(tenantID string, monObjID string, reset bool) error
	GetMetricBaselineForMonitoredObjectForHourOfWeek(tenantID string, monObjID string, hourOfWeek int32) ([]*tenmod.MetricBaselineData, error)
	DeleteMetricBaseline(tenantID string, dataID string) (*tenmod.MetricBaseline, error)
	// GetMetricBaselinesForMOsIn - note that this function will return results that are not stored in the DB as new "empty" items so that they can be populated
	// in a subsequent bulk PUT call.
	GetMetricBaselinesFor(tenantID string, moIDToHourOfWeekMap map[string][]int32, addNotFoundValuesInResponse bool) ([]*tenmod.MetricBaseline, error)
	BulkUpdateMetricBaselinesFromList(tenantID string, baselineUpdateList []*tenmod.MetricBaseline) ([]*common.BulkOperationResult, error)
	BulkUpdateMetricBaselines(tenantID string, entries []*swagmodels.MetricBaselineBulkUpdateRequestDataAttributesItems0) ([]*common.BulkOperationResult, error)

	CreateMetricBaselineDB(tenantID string) error
	DeleteMetricBaselineDB(tenantID string) error
}

// StubbedTenantMetricBaselineDatastore - Used to keep tests from failing when not testing multiple service interactions
type StubbedTenantMetricBaselineDatastore struct{}

func (stub *StubbedTenantMetricBaselineDatastore) CreateMetricBaseline(baseline *tenmod.MetricBaseline) (*tenmod.MetricBaseline, error) {
	return nil, fmt.Errorf("STUBBED OUT")
}
func (stub *StubbedTenantMetricBaselineDatastore) UpdateMetricBaseline(baseline *tenmod.MetricBaseline) (*tenmod.MetricBaseline, error) {
	return nil, fmt.Errorf("STUBBED OUT")
}
func (stub *StubbedTenantMetricBaselineDatastore) UpdateMetricBaselineForHourOfWeek(tenantID string, monObjID string, hourOfWeek int32, baselineData *tenmod.MetricBaselineData) (*tenmod.MetricBaseline, error) {
	return nil, fmt.Errorf("STUBBED OUT")
}
func (stub *StubbedTenantMetricBaselineDatastore) UpdateMetricBaselineForHourOfWeekWithCollection(tenantID string, monObjID string, hourOfWeek int32, baselineDataCollection []*tenmod.MetricBaselineData) (*tenmod.MetricBaseline, error) {
	return nil, fmt.Errorf("STUBBED OUT")
}
func (stub *StubbedTenantMetricBaselineDatastore) GetMetricBaseline(tenantID string, dataID string) (*tenmod.MetricBaseline, error) {
	return nil, fmt.Errorf("STUBBED OUT")
}
func (stub *StubbedTenantMetricBaselineDatastore) GetMetricBaselineForMonitoredObject(tenantID string, monObjID string) ([]*tenmod.MetricBaseline, error) {
	return nil, fmt.Errorf("STUBBED OUT")
}
func (stub *StubbedTenantMetricBaselineDatastore) DeleteMetricBaselineForMonitoredObject(tenantID string, monObjID string, reset bool) error {
	return fmt.Errorf("STUBBED OUT")
}
func (stub *StubbedTenantMetricBaselineDatastore) GetMetricBaselineForMonitoredObjectForHourOfWeek(tenantID string, monObjID string, hourOfWeek int32) ([]*tenmod.MetricBaselineData, error) {
	return nil, fmt.Errorf("STUBBED OUT")
}
func (stub *StubbedTenantMetricBaselineDatastore) DeleteMetricBaseline(tenantID string, dataID string) (*tenmod.MetricBaseline, error) {
	return nil, fmt.Errorf("STUBBED OUT")
}

func (stub *StubbedTenantMetricBaselineDatastore) GetMetricBaselinesFor(tenantID string, moIDToHourOfWeekMap map[string][]int32, addNotFoundValuesInResponse bool) ([]*tenmod.MetricBaseline, error) {
	return nil, fmt.Errorf("STUBBED OUT")
}
func (stub *StubbedTenantMetricBaselineDatastore) BulkUpdateMetricBaselinesFromList(tenantID string, baselineUpdateList []*tenmod.MetricBaseline) ([]*common.BulkOperationResult, error) {
	return nil, fmt.Errorf("STUBBED OUT")
}
func (stub *StubbedTenantMetricBaselineDatastore) BulkUpdateMetricBaselines(tenantID string, entries []*swagmodels.MetricBaselineBulkUpdateRequestDataAttributesItems0) ([]*common.BulkOperationResult, error) {
	return nil, fmt.Errorf("STUBBED OUT")
}

func (stub *StubbedTenantMetricBaselineDatastore) CreateMetricBaselineDB(tenantID string) error {
	return fmt.Errorf("STUBBED OUT")
}
func (stub *StubbedTenantMetricBaselineDatastore) DeleteMetricBaselineDB(tenantID string) error {
	return fmt.Errorf("STUBBED OUT")
}
