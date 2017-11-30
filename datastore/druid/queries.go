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

func filterHelper(metric string, e *pb.Event) *godruid.Filter {

	if e.UpperBound != 0 && e.LowerBound != 0 {
		return godruid.FilterLowerUpperBound(metric, "numeric", e.LowerBound, e.LowerStrict, e.UpperBound, e.UpperStrict)

	}

	if e.UpperBound != 0 {
		return godruid.FilterUpperBound(metric, "numeric", e.UpperBound, e.UpperStrict)
	}

	if e.LowerBound != 0 {
		return godruid.FilterLowerBound(metric, "numeric", e.LowerBound, e.LowerStrict)
	}
	return nil
}

func ThresholdCrossingQuery(dataSource string, metric string, granularity string, interval string, objectType string, direction string, events []*pb.Event) *godruid.QueryTimeseries {

	aggregations := make([]godruid.Aggregation, len(events)+1)
	postAggregations := make([]godruid.PostAggregation, len(events)*2)
	for i, e := range events {

		name := e.Severity + "Threshold"

		aggregations[i+1] = godruid.AggFiltered(
			godruid.FilterAnd(
				filterHelper(metric, e),
			),
			&godruid.Aggregation{
				Type: "count",
				Name: name,
			},
		)

		postAggregations[i] = godruid.PostAggArithmetic(e.Severity+"Ratio", "/", []godruid.PostAggregation{
			godruid.PostAggFieldAccessor(name),
			godruid.PostAggFieldAccessor("total"),
		})

		postAggregations[i+len(events)] = godruid.PostAggArithmetic(e.Severity+"Percent", "*", []godruid.PostAggregation{
			godruid.PostAggFieldAccessor(e.Severity + "Ratio"),
			godruid.PostAggConstant("", "100"),
		})

	}

	aggregations[0] = godruid.AggCount("total")

	return &godruid.QueryTimeseries{
		QueryType:        "timeseries",
		DataSource:       dataSource,
		Granularity:      godruid.GranHour,
		Context:          map[string]interface{}{"timeout": 60000},
		Aggregations:     aggregations,
		PostAggregations: postAggregations,
		Intervals:        []string{interval},
	}
}
