package druid_test

import (
	"testing"

	"github.com/accedian/adh-gather/datastore/druid"
	pb "github.com/accedian/adh-gather/gathergrpc"
	"github.com/accedian/godruid"
	"github.com/stretchr/testify/assert"
)

func TestFilterHelper(t *testing.T) {
	assert.Equal(t, druid.FilterHelper("foo", lowerEvent),
		&godruid.Filter{
			Type:        "bound",
			Dimension:   "foo",
			Ordering:    "numeric",
			Lower:       10000,
			LowerStrict: true,
		})

	assert.Equal(t, druid.FilterHelper("foo", upperEvent),
		&godruid.Filter{
			Type:        "bound",
			Dimension:   "foo",
			Ordering:    "numeric",
			Upper:       10000,
			UpperStrict: true,
		})

	assert.Equal(t, druid.FilterHelper("foo", bothEvent),
		&godruid.Filter{
			Type:        "bound",
			Dimension:   "foo",
			Ordering:    "numeric",
			Lower:       10000,
			LowerStrict: true,
			Upper:       10000,
			UpperStrict: true,
		})
}

func TestThresholdCrossingQuery(t *testing.T) {
	q := druid.ThresholdCrossingQuery("master", "druidTableName", "delayP95", "PT1H", "1900-11-02/2100-01-01", "TWAMP", "0", tp)

	assert.Equal(t, *q, *testThresholdCrossing1)
}

var metric1 = "delayP95"

var tp = &pb.TenantThresholdProfile{
	Thresholds: []*pb.TenantThreshold{
		&pb.TenantThreshold{
			ObjectType: "TWAMP",
			Metrics: []*pb.TenantMetric{
				&pb.TenantMetric{
					Id: "delayP95",
					Data: []*pb.TenantMetricData{
						&pb.TenantMetricData{
							Direction: "0",
							Events: []*pb.TenantEvent{
								&pb.TenantEvent{
									UpperBound:  30000,
									UpperStrict: true,
									Unit:        "percent",
									Type:        "minor",
								},
								&pb.TenantEvent{
									UpperBound:  50000,
									LowerBound:  30000,
									UpperStrict: true,
									LowerStrict: false,
									Unit:        "percent",
									Type:        "major",
								},
								&pb.TenantEvent{
									LowerBound:  50000,
									LowerStrict: false,
									Unit:        "percent",
									Type:        "critical",
								},
							},
						},
					},
				},
			},
		},
	},
}

var lowerEvent = &pb.TenantEvent{
	LowerBound:  10000,
	LowerStrict: true,
}

var upperEvent = &pb.TenantEvent{
	UpperBound:  10000,
	UpperStrict: true,
}

var bothEvent = &pb.TenantEvent{
	LowerBound:  10000,
	LowerStrict: true,
	UpperBound:  10000,
	UpperStrict: true,
}

var testThresholdCrossing1 = &godruid.QueryTimeseries{
	DataSource:  "druidTableName",
	Granularity: godruid.GranPeriod("PT1H", druid.TimeZoneUTC, ""),
	Context:     map[string]interface{}{"timeout": 60000},
	Aggregations: []godruid.Aggregation{
		godruid.AggCount("total"),
		godruid.AggFiltered(
			godruid.FilterAnd(
				godruid.FilterUpperBound(metric1, "numeric", 30000, true),
				godruid.FilterSelector("sessionType", "TWAMP"),
				godruid.FilterSelector("tenantId", "master"),
				godruid.FilterSelector("direction", "0"),
			),
			&godruid.Aggregation{
				Type: "count",
				Name: "TWAMP.delayP95.minor",
			},
		),
		godruid.AggFiltered(
			godruid.FilterAnd(
				godruid.FilterLowerUpperBound(metric1, "numeric", 30000, false, 50000, true),
				godruid.FilterSelector("sessionType", "TWAMP"),
				godruid.FilterSelector("tenantId", "master"),
				godruid.FilterSelector("direction", "0"),
			),
			&godruid.Aggregation{
				Type: "count",
				Name: "TWAMP.delayP95.major",
			},
		),
		godruid.AggFiltered(
			godruid.FilterAnd(
				godruid.FilterLowerBound(metric1, "numeric", 50000, false),
				godruid.FilterSelector("sessionType", "TWAMP"),
				godruid.FilterSelector("tenantId", "master"),
				godruid.FilterSelector("direction", "0"),
			),
			&godruid.Aggregation{
				Type: "count",
				Name: "TWAMP.delayP95.critical",
			},
		),
	},
	Intervals: []string{"1900-11-02/2100-01-01"},
}
