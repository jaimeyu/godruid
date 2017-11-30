package druid

import (
	pb "github.com/accedian/adh-gather/gathergrpc"
	"github.com/accedian/godruid"
)

func StatsQuery(dataSource string, metric string, threshold string, interval string) *godruid.QueryTimeseries {
	return &godruid.QueryTimeseries{
		QueryType:   "timeseries",
		DataSource:  dataSource,
		Granularity: godruid.GranHour,
		Aggregations: []godruid.Aggregation{
			godruid.AggLongMin(metric+"Min", metric),
			godruid.AggLongSum(metric+"Sum", metric),
			godruid.AggLongMax(metric+"Max", metric),
			godruid.AggCount("rowCount"),
		},
		PostAggregations: []godruid.PostAggregation{
			godruid.PostAggArithmetic("medianPoint", "/", []godruid.PostAggregation{
				godruid.PostAggFieldAccessor("rowCount"),
				godruid.PostAggConstant("", "2"),
			}),

			godruid.PostAggArithmetic("delayP95Mean", "/", []godruid.PostAggregation{
				godruid.PostAggFieldAccessor("delayP95Sum"),
				godruid.PostAggFieldAccessor("rowCount"),
			}),
		},
		Intervals: []string{interval},
	}
}

func ThresholdCrossingQuery(dataSource string, metric string, granularity string, interval string, objectType string, direction string, events []*pb.Event) *godruid.QueryTimeseries {

	return &godruid.QueryTimeseries{
		QueryType:   "timeseries",
		DataSource:  dataSource,
		Granularity: godruid.GranHour,
		Context:     map[string]interface{}{"timeout": 60000},

		Aggregations: []godruid.Aggregation{
			godruid.AggCount("total"),
			godruid.AggFiltered(
				godruid.FilterAnd(
					godruid.FilterUpperBound(metric, "numeric", "30000", true),
				),
				&godruid.Aggregation{
					Type: "count",
					Name: "minorThreshold",
				},
			),
			godruid.AggFiltered(
				godruid.FilterAnd(
					godruid.FilterLowerUpperBound(metric, "numeric", "30000", false, "75000", true),
				),
				&godruid.Aggregation{
					Type: "count",
					Name: "majorThreshold",
				},
			),
			godruid.AggFiltered(
				godruid.FilterAnd(
					godruid.FilterLowerBound(metric, "numeric", "75000", false),
				),
				&godruid.Aggregation{
					Type: "count",
					Name: "criticalThreshold",
				},
			),
		},
		PostAggregations: []godruid.PostAggregation{
			godruid.PostAggArithmetic("minorRatio", "/", []godruid.PostAggregation{
				godruid.PostAggFieldAccessor("minorThreshold"),
				godruid.PostAggFieldAccessor("total"),
			}),

			godruid.PostAggArithmetic("majorRatio", "/", []godruid.PostAggregation{
				godruid.PostAggFieldAccessor("majorThreshold"),
				godruid.PostAggFieldAccessor("total"),
			}),

			godruid.PostAggArithmetic("criticalRatio", "/", []godruid.PostAggregation{
				godruid.PostAggFieldAccessor("criticalThreshold"),
				godruid.PostAggFieldAccessor("total"),
			}),

			godruid.PostAggArithmetic("minorPercent", "*", []godruid.PostAggregation{
				godruid.PostAggFieldAccessor("minorRatio"),
				godruid.PostAggConstant("", "100"),
			}),

			godruid.PostAggArithmetic("majorPercent", "*", []godruid.PostAggregation{
				godruid.PostAggFieldAccessor("majorRatio"),
				godruid.PostAggConstant("", "100"),
			}),

			godruid.PostAggArithmetic("criticalPercent", "*", []godruid.PostAggregation{
				godruid.PostAggFieldAccessor("criticalRatio"),
				godruid.PostAggConstant("", "100"),
			}),
		},
		Intervals: []string{interval},
	}
}
