package druid_test

import (
	"testing"

	"github.com/accedian/adh-gather/datastore/druid"
	pb "github.com/accedian/adh-gather/gathergrpc"
	"github.com/accedian/adh-gather/logger"
	"github.com/accedian/godruid"
	"github.com/stretchr/testify/assert"
)

func TestFilterHelper(t *testing.T) {
	filter, err := druid.FilterHelper("foo", lowerEvent)
	if err != nil {
		logger.Log.Error("Filter helper error: ", err)
	}

	assert.Equal(t, &godruid.Filter{
		Type:        "bound",
		Dimension:   "foo",
		Ordering:    "numeric",
		Lower:       10000,
		LowerStrict: true,
	}, filter)

	filter, err = druid.FilterHelper("foo", upperEvent)
	if err != nil {
		logger.Log.Error("Filter helper error: ", err)
	}

	assert.Equal(t, &godruid.Filter{
		Type:        "bound",
		Dimension:   "foo",
		Ordering:    "numeric",
		Upper:       10000,
		UpperStrict: true,
	}, filter)

	filter, err = druid.FilterHelper("foo", bothEvent)
	if err != nil {
		logger.Log.Error("Filter helper error: ", err)
	}

	assert.Equal(t, &godruid.Filter{
		Type:        "bound",
		Dimension:   "foo",
		Ordering:    "numeric",
		Lower:       10000,
		LowerStrict: true,
		Upper:       10000,
		UpperStrict: true,
	}, filter)
}

func TestThresholdCrossingQuery(t *testing.T) {
	// PEYO TODO FIX THIS TEST
	// q, err := druid.ThresholdCrossingQuery("master", "druidTableName", "delayP95", "PT1H", "1900-11-02/2100-01-01", "TWAMP", "0", tp)

	// if err != nil {
	// 	logger.Log.Error("ThresholdCrossing query error: ", err)
	// }

	// qJson, _ := json.Marshal(q)
	// expectedJson, _ := json.Marshal(testThresholdCrossing1)

	// check to see if the number of bytes is equal since filters can be out of order
	// assert.Equal(t, len(expectedJson), len(qJson))
}

func TestBuildMetaFilter(t *testing.T) {

	testMetaMap := make(map[string][]string)

	colours := []string{"blue", "red"}
	cities := []string{"Ottawa", "Montreal"}

	testMetaMap["colour"] = colours
	testMetaMap["cities"] = cities

	allEntries := append(colours, cities...)
	comboFilter := druid.BuildMetaFilter("test", testMetaMap)

	if comboFilter.Type != "and" {
		t.Errorf("Incorrect filter built. Expecting '%s' but got '%s'", "and", comboFilter.Type)
	}

	andFilters := comboFilter.Fields

	if len(andFilters) != len(testMetaMap) {
		t.Errorf("Incorrect number of AND filters. Expecting '%d' but got '%d'", len(testMetaMap), len(andFilters))
	}

	for _, andFilter := range andFilters {
		if andFilter.Type != "or" {
			t.Errorf("Incorrect sub-filter built. Expecting '%s' but got '%s'", "or", andFilter.Type)
		}
		for _, orFilter := range andFilter.Fields {
			if orFilter.Type != "selector" {
				t.Errorf("Incorrect sub-filter built. Expecting '%s' but got '%s'", "selector", orFilter.Type)
			}

			found := false

			for i, testEntry := range allEntries {
				if testEntry == orFilter.Value {
					found = true
					allEntries = append(allEntries[:i], allEntries[i+1:]...)
					break
				}
			}
			if !found {
				t.Errorf("Extra filter with value %s was created but should not have been.", orFilter.Value)
			}
		}
	}

	if len(allEntries) > 0 {
		t.Errorf("The following meta values do not have OR filters against them: %v", allEntries)
	}
}

var metric1 = "delayP95"

var tp = &pb.TenantThresholdProfileData{
	Thresholds: &pb.TenantThresholdProfileData_VendorMap{
		VendorMap: map[string]*pb.TenantThresholdProfileData_MonitoredObjectTypeMap{
			"accedian": &pb.TenantThresholdProfileData_MonitoredObjectTypeMap{
				MonitoredObjectTypeMap: map[string]*pb.TenantThresholdProfileData_MetricMap{
					"TWAMP": &pb.TenantThresholdProfileData_MetricMap{
						MetricMap: map[string]*pb.TenantThresholdProfileData_DirectionMap{
							"delayP95": &pb.TenantThresholdProfileData_DirectionMap{
								DirectionMap: map[string]*pb.TenantThresholdProfileData_EventMap{
									"0": &pb.TenantThresholdProfileData_EventMap{
										EventMap: map[string]*pb.TenantThresholdProfileData_EventAttrMap{
											"minor": &pb.TenantThresholdProfileData_EventAttrMap{
												map[string]string{
													"upperLimit": "30000",
													"unit":       "ms",
												},
											},
											"major": &pb.TenantThresholdProfileData_EventAttrMap{
												map[string]string{
													"lowerLimit":  "30000",
													"lowerStrict": "true",
													"upperLimit":  "50000",
													"upperStrict": "false",
													"unit":        "ms",
												},
											},
											"critical": &pb.TenantThresholdProfileData_EventAttrMap{
												map[string]string{
													"lowerLimit":  "50000",
													"lowerStrict": "true",
													"unit":        "ms",
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},
		},
	},
}

var lowerEvent = &pb.TenantThresholdProfileData_EventAttrMap{
	EventAttrMap: map[string]string{
		"lowerLimit":  "10000",
		"lowerStrict": "true",
	},
}

var upperEvent = &pb.TenantThresholdProfileData_EventAttrMap{
	EventAttrMap: map[string]string{
		"upperLimit":  "10000",
		"upperStrict": "true",
	},
}

var bothEvent = &pb.TenantThresholdProfileData_EventAttrMap{
	EventAttrMap: map[string]string{
		"upperLimit":  "10000",
		"upperStrict": "true",
		"lowerLimit":  "10000",
		"lowerStrict": "true",
	},
}

var testThresholdCrossing1 = &godruid.QueryTimeseries{
	DataSource:  "druidTableName",
	Granularity: godruid.GranPeriod("PT1H", druid.TimeZoneUTC, ""),
	Context:     map[string]interface{}{"timeout": 60000},
	Aggregations: []godruid.Aggregation{
		godruid.AggCount("total"),
		godruid.AggFiltered(
			godruid.FilterAnd(
				godruid.FilterUpperBound(metric1, "numeric", 30000, false),
				godruid.FilterSelector("objectType", "TWAMP"),
				godruid.FilterSelector("tenantId", "master"),
				godruid.FilterSelector("direction", "0"),
			),
			&godruid.Aggregation{
				Type: "count",
				Name: "TWAMP.delayP95.minor.0",
			},
		),
		godruid.AggFiltered(
			godruid.FilterAnd(
				godruid.FilterLowerUpperBound(metric1, "numeric", 30000, true, 50000, false),
				godruid.FilterSelector("objectType", "TWAMP"),
				godruid.FilterSelector("tenantId", "master"),
				godruid.FilterSelector("direction", "0"),
			),
			&godruid.Aggregation{
				Type: "count",
				Name: "TWAMP.delayP95.major.0",
			},
		),
		godruid.AggFiltered(
			godruid.FilterAnd(
				godruid.FilterLowerBound(metric1, "numeric", 50000, true),
				godruid.FilterSelector("objectType", "TWAMP"),
				godruid.FilterSelector("tenantId", "master"),
				godruid.FilterSelector("direction", "0"),
			),
			&godruid.Aggregation{
				Type: "count",
				Name: "TWAMP.delayP95.critical.0",
			},
		),
	},
	Intervals: []string{"1900-11-02/2100-01-01"},
}
