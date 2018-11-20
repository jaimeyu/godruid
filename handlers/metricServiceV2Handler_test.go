package handlers_test

import (
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/accedian/adh-gather/datastore/druid"
	"github.com/accedian/adh-gather/handlers"
	metmod "github.com/accedian/adh-gather/models/metrics"
	tenmod "github.com/accedian/adh-gather/models/tenant"
	"github.com/accedian/adh-gather/restapi/operations/metrics_service_v2"
	"github.com/accedian/adh-gather/swagmodels"
	"github.com/stretchr/testify/assert"
)

var allTwampMonitoredObjectTypes []string = []string{"twamp-sf", "twamp-sl", "twamp-pe"}
var allDirections []string = []string{"0", "1", "2"}

// Query Defaults
const (
	defaultInterval    = "2018-01-01/2025-01-01"
	defaultVendor      = "accedian-twamp"
	defaultGranularity = "PT1H"
)

// Request Defaults
var (
	typeThresholdCrossings = "thresholdCrossings"
)

// Threshold Crossing Tests
func TestThresholdCrossingRegularCrossing(t *testing.T) {

	if !*metricsIntegrationTests {
		return
	}

	testTenant := "testthresholdcrossing"
	metric := "delayMax"

	testProfile := constructThresholdProfile([]tenmod.ThresholdProfileThreshold{
		tenmod.ThresholdProfileThreshold{
			Vendor:              defaultVendor,
			MonitoredObjectType: "twamp-sf",
			Direction:           "0",
			Metric:              metric,
			Events: []map[string]string{
				map[string]string{
					"eventName":  "critical",
					"upperLimit": "300",
					"lowerLimit": "200",
				},
				map[string]string{
					"eventName":  "major",
					"upperLimit": "200",
					"lowerLimit": "100",
				},
				map[string]string{
					"eventName":  "minor",
					"upperLimit": "100",
					"lowerLimit": "0",
				},
			},
		},
	})

	metricDatastoreClient := druid.NewDruidDatasctoreClient()
	tenantDatastore := TenantServiceDatastoreStub{thresholdProfile: testProfile}

	req := constructThresholdCrossingRequest(testTenant, defaultInterval, defaultGranularity, []metmod.MetricIdentifierFilter{
		metmod.MetricIdentifierFilter{
			Vendor:     defaultVendor,
			ObjectType: allTwampMonitoredObjectTypes,
			Metric:     metric,
			Direction:  []string{"0"}},
	})

	expected := swagmodels.ThresholdCrossingReportMetric{
		Critical: []*swagmodels.ThresholdCrossingViolations{
			&swagmodels.ThresholdCrossingViolations{
				ViolationCount:    iAddress(1),
				ViolationDuration: iAddress(50),
			},
		},
		Major: []*swagmodels.ThresholdCrossingViolations{
			&swagmodels.ThresholdCrossingViolations{
				ViolationCount:    iAddress(2),
				ViolationDuration: iAddress(100),
			},
		},
		Minor: []*swagmodels.ThresholdCrossingViolations{
			&swagmodels.ThresholdCrossingViolations{
				ViolationCount:    iAddress(1),
				ViolationDuration: iAddress(50),
			},
		},
	}

	response := handlers.HandleGetThresholdCrossingV2(handlers.SkylightAdminRoleOnly, metricDatastoreClient, tenantDatastore)(req)
	assert.IsType(t, &metrics_service_v2.QueryThresholdCrossingV2OK{}, response)
	metricReport := response.(*metrics_service_v2.QueryThresholdCrossingV2OK)
	assert.NotNil(t, metricReport)
	metricResponses := metricReport.Payload.Data.Attributes.Result.Metric

	valid, err := validateMetricThresholdCrossingEntries(expected, *metricResponses[0])
	assert.Truef(t, valid, "Expected threshold crossing response issue: %v", err)
}

func TestThresholdCrossingBaselineStaticCrossing(t *testing.T) {

	if !*metricsIntegrationTests {
		return
	}

	testTenant := "testthresholdcrossing"
	metric := "delayMax"

	testProfile := constructThresholdProfile([]tenmod.ThresholdProfileThreshold{
		tenmod.ThresholdProfileThreshold{
			Vendor:              defaultVendor,
			MonitoredObjectType: "twamp-sf",
			Direction:           "0",
			Metric:              metric,
			Events: []map[string]string{
				map[string]string{
					"eventName":  "critical",
					"eventType":  "baseline_static",
					"lowerLimit": "100",
				},
				map[string]string{
					"eventName":   "major",
					"eventType":   "baseline_static",
					"upperLimit":  "100",
					"upperStrict": "true",
					"lowerLimit":  "50",
				},
				map[string]string{
					"eventName":   "minor",
					"eventType":   "baseline_static",
					"upperLimit":  "50",
					"upperStrict": "true",
					"lowerLimit":  "20",
				},
			},
		},
	})

	metricDatastoreClient := druid.NewDruidDatasctoreClient()
	tenantDatastore := TenantServiceDatastoreStub{thresholdProfile: testProfile}

	req := constructThresholdCrossingRequest(testTenant, defaultInterval, defaultGranularity, []metmod.MetricIdentifierFilter{
		metmod.MetricIdentifierFilter{
			Vendor:     defaultVendor,
			ObjectType: allTwampMonitoredObjectTypes,
			Metric:     metric,
			Direction:  []string{"0"}},
	})

	expected := swagmodels.ThresholdCrossingReportMetric{
		Critical: []*swagmodels.ThresholdCrossingViolations{
			&swagmodels.ThresholdCrossingViolations{
				ViolationCount:    iAddress(1),
				ViolationDuration: iAddress(50),
			},
		},
		Major: []*swagmodels.ThresholdCrossingViolations{
			&swagmodels.ThresholdCrossingViolations{
				ViolationCount:    iAddress(1),
				ViolationDuration: iAddress(50),
			},
		},
		Minor: []*swagmodels.ThresholdCrossingViolations{
			&swagmodels.ThresholdCrossingViolations{
				ViolationCount:    iAddress(1),
				ViolationDuration: iAddress(50),
			},
		},
	}

	response := handlers.HandleGetThresholdCrossingV2(handlers.SkylightAdminRoleOnly, metricDatastoreClient, tenantDatastore)(req)
	assert.IsType(t, &metrics_service_v2.QueryThresholdCrossingV2OK{}, response)
	metricReport := response.(*metrics_service_v2.QueryThresholdCrossingV2OK)
	assert.NotNil(t, metricReport)
	metricResponses := metricReport.Payload.Data.Attributes.Result.Metric

	valid, err := validateMetricThresholdCrossingEntries(expected, *metricResponses[0])
	assert.Truef(t, valid, "Expected threshold crossing response issue: %v", err)
}

func TestThresholdCrossingBaselinePercentageCrossing(t *testing.T) {

	if !*metricsIntegrationTests {
		return
	}

	testTenant := "testthresholdcrossing"
	metric := "delayMax"

	testProfile := constructThresholdProfile([]tenmod.ThresholdProfileThreshold{
		tenmod.ThresholdProfileThreshold{
			Vendor:              defaultVendor,
			MonitoredObjectType: "twamp-sf",
			Direction:           "0",
			Metric:              metric,
			Events: []map[string]string{
				map[string]string{
					"eventName":  "critical",
					"eventType":  "baseline_percentage",
					"lowerLimit": "1",
				},
				map[string]string{
					"eventName":   "major",
					"eventType":   "baseline_percentage",
					"upperLimit":  "1",
					"upperStrict": "true",
					"lowerLimit":  "0.5",
				},
				map[string]string{
					"eventName":   "minor",
					"eventType":   "baseline_percentage",
					"upperLimit":  "0.5",
					"upperStrict": "true",
					"lowerLimit":  "0.19",
				},
			},
		},
	})

	metricDatastoreClient := druid.NewDruidDatasctoreClient()
	tenantDatastore := TenantServiceDatastoreStub{thresholdProfile: testProfile}

	req := constructThresholdCrossingRequest(testTenant, defaultInterval, defaultGranularity, []metmod.MetricIdentifierFilter{
		metmod.MetricIdentifierFilter{
			Vendor:     defaultVendor,
			ObjectType: allTwampMonitoredObjectTypes,
			Metric:     metric,
			Direction:  []string{"0"}},
	})

	expected := swagmodels.ThresholdCrossingReportMetric{
		Critical: []*swagmodels.ThresholdCrossingViolations{
			&swagmodels.ThresholdCrossingViolations{
				ViolationCount:    iAddress(1),
				ViolationDuration: iAddress(50),
			},
		},
		Major: []*swagmodels.ThresholdCrossingViolations{
			&swagmodels.ThresholdCrossingViolations{
				ViolationCount:    iAddress(1),
				ViolationDuration: iAddress(50),
			},
		},
		Minor: []*swagmodels.ThresholdCrossingViolations{
			&swagmodels.ThresholdCrossingViolations{
				ViolationCount:    iAddress(1),
				ViolationDuration: iAddress(50),
			},
		},
	}

	response := handlers.HandleGetThresholdCrossingV2(handlers.SkylightAdminRoleOnly, metricDatastoreClient, tenantDatastore)(req)
	assert.IsType(t, &metrics_service_v2.QueryThresholdCrossingV2OK{}, response)
	metricReport := response.(*metrics_service_v2.QueryThresholdCrossingV2OK)
	assert.NotNil(t, metricReport)
	metricResponses := metricReport.Payload.Data.Attributes.Result.Metric

	valid, err := validateMetricThresholdCrossingEntries(expected, *metricResponses[0])
	assert.Truef(t, valid, "Expected threshold crossing response issue: %v", err)
}

func constructThresholdCrossingRequest(tenant string, interval string, granularity string, metricFilters []metmod.MetricIdentifierFilter) metrics_service_v2.QueryThresholdCrossingV2Params {

	metFilters := make([]*swagmodels.MetricIdentifierFilter, 0)

	for _, filter := range metricFilters {
		addFilter := swagmodels.MetricIdentifierFilter{
			Vendor:     &filter.Vendor,
			ObjectType: filter.ObjectType,
			Metric:     &filter.Metric,
			Direction:  filter.Direction,
		}
		metFilters = append(metFilters, &addFilter)
	}

	tcQuery := metrics_service_v2.QueryThresholdCrossingV2Params{
		HTTPRequest: createHttpRequestWithParams(tenant, handlers.UserRoleSkylight, "", "POST"),
		Body: &swagmodels.JSONAPIThresholdCrossingRequest{
			Data: &swagmodels.JSONAPIThresholdCrossingRequestData{
				Attributes: &swagmodels.ThresholdCrossingConfig{
					Interval:    &interval,
					Granularity: granularity,
					Metrics:     metFilters},
				Type: &typeThresholdCrossings,
			},
		},
	}

	return tcQuery
}

func constructThresholdProfile(thresholds []tenmod.ThresholdProfileThreshold) *tenmod.ThresholdProfile {
	// TODO change this once we move to flattened threshold profile
	tp := tenmod.ThresholdProfile{
		Thresholds: &tenmod.ThrPrfVendorMap{
			map[string]*tenmod.ThrPrfMetric{},
		},
	}

	for _, threshold := range thresholds {
		vMap := tp.Thresholds.VendorMap[threshold.Vendor]
		if vMap == nil {
			vMap = &tenmod.ThrPrfMetric{
				MonitoredObjectTypeMap: map[string]*tenmod.ThrPrfMetricMap{},
			}
		}
		oMap := vMap.MonitoredObjectTypeMap[threshold.MonitoredObjectType]
		if oMap == nil {
			oMap = &tenmod.ThrPrfMetricMap{
				MetricMap: map[string]*tenmod.ThrPrfDirectionMap{},
			}
		}
		mMap := oMap.MetricMap[threshold.Metric]
		if mMap == nil {
			mMap = &tenmod.ThrPrfDirectionMap{
				DirectionMap: map[string]*tenmod.ThrPrfEventMap{},
			}
		}
		dMap := mMap.DirectionMap[threshold.Direction]
		if dMap == nil {
			dMap = &tenmod.ThrPrfEventMap{
				EventMap: map[string]*tenmod.ThrPrfEventAttrMap{},
			}
		}
		for _, event := range threshold.Events {
			eventName := event["eventName"]
			eMap := dMap.EventMap[eventName]
			if eMap == nil {
				eMap = &tenmod.ThrPrfEventAttrMap{}
			}
			eMap.EventAttrMap = event

			dMap.EventMap[eventName] = eMap
		}

		mMap.DirectionMap[threshold.Direction] = dMap
		oMap.MetricMap[threshold.Metric] = mMap
		vMap.MonitoredObjectTypeMap[threshold.MonitoredObjectType] = oMap
		tp.Thresholds.VendorMap[threshold.Vendor] = vMap
	}

	return &tp
}

func validateMetricThresholdCrossingEntries(expected swagmodels.ThresholdCrossingReportMetric, actual swagmodels.ThresholdCrossingReportMetric) (bool, error) {
	if valid, err := validateThresholdCrossingViolationsEntry(expected.Critical, actual.Critical, false); !valid {
		return false, err
	}
	if valid, err := validateThresholdCrossingViolationsEntry(expected.Major, actual.Major, false); !valid {
		return false, err
	}
	if valid, err := validateThresholdCrossingViolationsEntry(expected.Minor, actual.Minor, false); !valid {
		return false, err
	}
	if valid, err := validateThresholdCrossingViolationsEntry(expected.Warning, actual.Warning, false); !valid {
		return false, err
	}
	if valid, err := validateThresholdCrossingViolationsEntry(expected.SLA, actual.SLA, false); !valid {
		return false, err
	}

	return true, nil
}

func validateThresholdCrossingViolationsEntry(expected []*swagmodels.ThresholdCrossingViolations, actual []*swagmodels.ThresholdCrossingViolations, validateTimestamps bool) (bool, error) {

	if len(expected) != len(actual) {
		return false, fmt.Errorf("Actual number of threshold crossing violations entries %d did not match the expected number of %d", len(actual), len(expected))
	}

	if valid, err := validateOrdered(actual); !valid || err != nil {
		return false, err
	}

	for i, expectedEntry := range expected {
		actualEntry := actual[i]
		if validateTimestamps {
			if expectedEntry.Timestamp != actualEntry.Timestamp {
				return false, fmt.Errorf("Expected timestamp %v did not match up with actual timestamp %v", expectedEntry.Timestamp, actualEntry.Timestamp)
			}
		}
		if *expectedEntry.ViolationCount != *actualEntry.ViolationCount {
			return false, fmt.Errorf("Expected violation count %v did not match up with actual violation count %v for timestamp %v", *expectedEntry.ViolationCount, *actualEntry.ViolationCount, *actualEntry.Timestamp)
		}
		if *expectedEntry.ViolationDuration != *actualEntry.ViolationDuration {
			return false, fmt.Errorf("Expected violation duration %v did not match up with actual violation duration %v for timestamp %v", *expectedEntry.ViolationDuration, *actualEntry.ViolationDuration, *actualEntry.Timestamp)
		}
	}

	return true, nil
}

// Validates that all of the metrics that were requested in the filter have a response associated with them.
// The argument should be a list of response entries that contain a metric identifier at the root level
func validateMetricResultEntries(expectedMetricEntries []metmod.MetricIdentifierFilter, metricList interface{}) (bool, error) {
	metricListContainer, err := convertToListOfMaps(metricList)
	if err != nil {
		return false, err
	}

	if len(metricListContainer) != len(expectedMetricEntries) {
		return false, fmt.Errorf("Mismatch between expected number of metric entries %d and actual metric entries %d", len(expectedMetricEntries), len(metricListContainer))
	}

	for _, expectedMetric := range expectedMetricEntries {
		if indexOfMetric(expectedMetric, metricListContainer) == -1 {
			return false, fmt.Errorf("Could not find metric result that matches metric identifier %v", expectedMetric)
		}
	}

	return true, nil
}

// Validates that metric timeseries data is in chronological order.
// The argument should be a list of objects that can be broken down into an array of string->string maps that contain a "timestamp" key
func validateOrdered(timestampedKeyList interface{}) (bool, error) {

	timestampedListContainer, err := convertToListOfMaps(timestampedKeyList)
	if err != nil {
		return false, err
	}

	var t1 *time.Time
	for _, entry := range timestampedListContainer {
		t2, err := time.Parse(time.RFC3339, entry["timestamp"].(string))
		if err != nil {
			return false, nil
		}

		if t1 != nil && !t1.Before(t2) {
			return false, fmt.Errorf("Timeseries entry with timestamp %v is not before %v", t1, t2)
		}

		t1 = &t2
	}

	return true, nil
}

func convertToListOfMaps(container interface{}) ([]map[string]interface{}, error) {
	if container == nil {
		return nil, fmt.Errorf("Provided list of timestamp keyed entries is nil")
	}

	itemBytes, err := json.Marshal(container)
	if err != nil {
		return nil, err
	}

	listContainer := make([]map[string]interface{}, 0)
	err = json.Unmarshal(itemBytes, &listContainer)
	if err != nil {
		return nil, err
	}

	return listContainer, nil
}

func indexOfMetric(mid metmod.MetricIdentifierFilter, metricEntries []map[string]interface{}) int {

	for i, metricEntry := range metricEntries {
		if testVendor, ok := metricEntry["vendor"]; !ok || testVendor != mid.Vendor {
			continue
		}
		if testMetric, ok := metricEntry["metric"]; !ok || testMetric != mid.Metric {
			continue
		}
		if testDirections, ok := metricEntry["direction"]; !ok || !isArrayContentEqual(testDirections.([]string), mid.Direction) {
			continue
		}
		if testObjectTypes, ok := metricEntry["objectType"]; !ok || !isArrayContentEqual(testObjectTypes.([]string), mid.ObjectType) {
			continue
		}
		return i
	}

	return -1
}

func checkArrayComparisonPresence(a1 []string, a2 []string) bool {
	if (a1 == nil) != (a2 == nil) {
		return false
	}
	if len(a1) != len(a2) {
		return false
	}
	return true
}

func isArrayContentEqual(a1 []string, a2 []string) bool {

	if ok := checkArrayComparisonPresence(a1, a2); !ok {
		return false
	}

	for i := range a1 {
		if a1[i] != a2[i] {
			return false
		}
	}
	return true
}

func iAddress(n int64) *int64 {
	return &n
}

// func TestAggregate(t *testing.T) {
// 	metricDatastoreClient := druid.NewDruidDatasctoreClient()
// 	tenantDatastore, err := handlers.GetTenantServiceDatastore()
// 	if err != nil {
// 		fmt.Errorf("Error while attempting to connect to configure datastore: %v", err)
// 		return
// 	}

// 	aggregation := "avg"
// 	interval := "2018-01-01/2025-01-01"
// 	vendor := "accedian-twamp"
// 	metric := "delayMin"
// 	rType := "aggregateMetrics"

// 	req := metrics_service_v2.QueryAggregateMetricsV2Params{
// 		HTTPRequest: createHttpRequestWithParams("", handlers.UserRoleSkylight, "", "POST"),
// 		Body: &swagmodels.JSONAPIAggregationRequest{
// 			Data: &swagmodels.JSONAPIAggregationRequestData{
// 				Attributes: &swagmodels.AggregationConfig{
// 					Aggregation: &aggregation,
// 					Granularity: "PT1H",
// 					Interval:    &interval,
// 					Metrics: []*swagmodels.MetricIdentifierFilter{
// 						&swagmodels.MetricIdentifierFilter{
// 							Vendor:     &vendor,
// 							ObjectType: []string{"twamp-sl", "twamp-sf"},
// 							Metric:     &metric,
// 							Direction:  []string{"0"},
// 						},
// 					},
// 					Timeout: 30000,
// 				},
// 				Type: &rType,
// 			},
// 		}}
// 	req.HTTPRequest.Header.Add(XFwdTenantId, testTenant)
// 	doGetAggregateMetricsV2(SkylightAdminRoleOnly, metricDatastoreClient, tenantDatastore, req)
// }

type TenantServiceDatastoreStub struct {
	filteredMonitoredObjectList []string
	moids                       []string
	thresholdProfile            *tenmod.ThresholdProfile
}

func (tsd TenantServiceDatastoreStub) GetFilteredMonitoredObjectList(tenantId string, meta map[string][]string) ([]string, error) {
	return tsd.filteredMonitoredObjectList, nil
}
func (tsd TenantServiceDatastoreStub) GetAllMonitoredObjectsIDs(tenantID string) ([]string, error) {
	return tsd.moids, nil
}
func (tsd TenantServiceDatastoreStub) GetTenantThresholdProfile(tenantID string, dataID string) (*tenmod.ThresholdProfile, error) {
	return tsd.thresholdProfile, nil
}
