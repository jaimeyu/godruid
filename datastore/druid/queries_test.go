package druid_test

import (
	"encoding/json"
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
	q, err := druid.ThresholdCrossingQuery("master", "druidTableName", "delayP95", "PT1H", "1900-11-02/2100-01-01", "TWAMP", "0", tp)

	if err != nil {
		logger.Log.Error("ThresholdCrossing query error: ", err)
	}

	qJson, _ := json.Marshal(q)
	expectedJson, _ := json.Marshal(testThresholdCrossing1)

	// check to see if the number of bytes is equal since filters can be out of order
	assert.Equal(t, len(expectedJson), len(qJson))
}

var metric1 = "delayP95"

var tp = &pb.TenantThresholdProfile{
	Thresholds: &pb.TenantThresholdProfile_VendorMap{
		VendorMap: map[string]*pb.TenantThresholdProfile_MonitoredObjectTypeMap{
			"accedian": &pb.TenantThresholdProfile_MonitoredObjectTypeMap{
				MonitoredObjectTypeMap: map[string]*pb.TenantThresholdProfile_MetricMap{
					"TWAMP": &pb.TenantThresholdProfile_MetricMap{
						MetricMap: map[string]*pb.TenantThresholdProfile_DirectionMap{
							"delayP95": &pb.TenantThresholdProfile_DirectionMap{
								DirectionMap: map[string]*pb.TenantThresholdProfile_EventMap{
									"0": &pb.TenantThresholdProfile_EventMap{
										EventMap: map[string]*pb.TenantThresholdProfile_EventAttrMap{
											"minor": &pb.TenantThresholdProfile_EventAttrMap{
												map[string]string{
													"upperLimit": "30000",
													"unit":       "ms",
												},
											},
											"major": &pb.TenantThresholdProfile_EventAttrMap{
												map[string]string{
													"lowerLimit":  "30000",
													"lowerStrict": "true",
													"upperLimit":  "50000",
													"upperStrict": "false",
													"unit":        "ms",
												},
											},
											"critical": &pb.TenantThresholdProfile_EventAttrMap{
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

var lowerEvent = &pb.TenantThresholdProfile_EventAttrMap{
	EventAttrMap: map[string]string{
		"lowerLimit":  "10000",
		"lowerStrict": "true",
	},
}

var upperEvent = &pb.TenantThresholdProfile_EventAttrMap{
	EventAttrMap: map[string]string{
		"upperLimit":  "10000",
		"upperStrict": "true",
	},
}

var bothEvent = &pb.TenantThresholdProfile_EventAttrMap{
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
