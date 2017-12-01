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
	druid.ThresholdCrossingQuery("druidTableName", metric1, "1h", "1900-11-02/2100-01-01", "twamp", "az", tp.Twamp.Metrics[0].Data[0].Events)

}

var metric1 = "delayP95"

var tp = &pb.ThresholdProfile{
	Twamp: &pb.Threshold{
		Metrics: []*pb.Metric{
			&pb.Metric{
				Id: "delayP95",
				Data: []*pb.MetricData{
					&pb.MetricData{
						Direction: "az",
						Events: []*pb.Event{
							&pb.Event{
								UpperBound:  30000,
								UpperStrict: true,
								Unit:        "percent",
								Type:        "minor",
							},
							&pb.Event{
								UpperBound:  50000,
								LowerBound:  30000,
								UpperStrict: true,
								LowerStrict: false,
								Unit:        "percent",
								Type:        "major",
							},
							&pb.Event{
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
}

var lowerEvent = &pb.Event{
	LowerBound:  10000,
	LowerStrict: true,
}

var upperEvent = &pb.Event{
	UpperBound:  10000,
	UpperStrict: true,
}

var bothEvent = &pb.Event{
	LowerBound:  10000,
	LowerStrict: true,
	UpperBound:  10000,
	UpperStrict: true,
}

var testThresholdCrossing1 = &godruid.QueryTimeseries{
	QueryType:   "timeseries",
	DataSource:  "druidTableName",
	Granularity: godruid.GranHour,
	Context:     map[string]interface{}{"timeout": 60000},
	Aggregations: []godruid.Aggregation{
		godruid.AggCount("total"),
		godruid.AggFiltered(
			godruid.FilterAnd(
				godruid.FilterLowerBound(metric1, "numeric", 75000, false),
			),
			&godruid.Aggregation{
				Type: "count",
				Name: "minorThreshold",
			},
		),
		godruid.AggFiltered(
			godruid.FilterAnd(
				godruid.FilterUpperBound(metric1, "numeric", 30000, true),
			),
			&godruid.Aggregation{
				Type: "count",
				Name: "majorThreshold",
			},
		),
		godruid.AggFiltered(
			godruid.FilterAnd(
				godruid.FilterLowerUpperBound(metric1, "numeric", 30000, false, 75000, true),
			),
			&godruid.Aggregation{
				Type: "count",
				Name: "criticalThreshold",
			},
		),
	},
	Intervals: []string{"1900-11-02/2100-01-01"},
}
