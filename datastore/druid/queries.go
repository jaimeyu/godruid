package druid

import (
	pb "github.com/accedian/adh-gather/gathergrpc"
	"github.com/accedian/godruid"
)

// HistogramQuery -
func HistogramQuery(dataSource string, metric string, granularity string, interval string, resolution int32, numBuckets int32) *godruid.QueryTimeseries {

	return &godruid.QueryTimeseries{
		QueryType:  "timeseries",
		DataSource: dataSource,
		Granularity: godruid.GranPeriod{
			Type:     "period",
			Period:   granularity,
			TimeZone: "UTC",
		},
		Context: map[string]interface{}{"timeout": 60000},
		Aggregations: []godruid.Aggregation{
			godruid.AggHistoFold("thresholdBuckets", metric+"P95Histo", resolution, numBuckets),
		},
		Intervals: []string{interval},
	}
}

// FilterHelper - helper function to select correct druid filter based on
// a given event and metric
func FilterHelper(metric string, e *pb.Event) *godruid.Filter {

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

// ThresholdCrossingQuery - Query that returns a count of events that crossed a given
// threshold, for a given metric.
func ThresholdCrossingQuery(dataSource string, metric string, granularity string, interval string, objectType string, direction string, events []*pb.Event) *godruid.QueryTimeseries {

	aggregations := make([]godruid.Aggregation, len(events)+1)

	for i, e := range events {

		name := e.Type + "Threshold"

		aggregations[i+1] = godruid.AggFiltered(
			godruid.FilterAnd(
				FilterHelper(metric, e),
			),
			&godruid.Aggregation{
				Type: "count",
				Name: name,
			},
		)
	}

	aggregations[0] = godruid.AggCount("total")

	return &godruid.QueryTimeseries{
		QueryType:  "timeseries",
		DataSource: dataSource,
		Granularity: godruid.GranPeriod{
			Type:     "period",
			Period:   granularity,
			TimeZone: "UTC",
		},
		Context:      map[string]interface{}{"timeout": 60000},
		Aggregations: aggregations,
		Intervals:    []string{interval},
	}
}
