package druid

import (
	"fmt"

	pb "github.com/accedian/adh-gather/gathergrpc"
	"github.com/accedian/godruid"
)

// HistogramQuery - Count of metrics per bucket for given interval.
func HistogramQuery(dataSource string, metric string, granularity string, interval string, resolution int32, granularityBuckets int32) *godruid.QueryTimeseries {

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
			godruid.AggHistoFold("thresholdBuckets", metric+"P95Histo", resolution, granularityBuckets, "0", "Infinity"),
		},
		Intervals: []string{interval},
	}
}

// FilterHelper - helper function to select correct druid filter based on
// a given event and metric
func FilterHelper(metric string, e *pb.TenantEvent) *godruid.Filter {

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

// ThresholdCrossingQuery - Query that returns a count of events that crossed a thresholds for metric/thresholds
// defined by the supplied threshold profile..
func ThresholdCrossingQuery(dataSource string, metric string, granularity string, interval string, objectType string, direction string, thresholdProfile *pb.TenantThresholdProfile) *godruid.QueryTimeseries {

	var aggregations []godruid.Aggregation

	aggregations = append(aggregations, godruid.AggCount("total"))

	for i, t := range thresholdProfile.GetThresholds() {
		fmt.Println(i, " ", t.ObjectType)
		for j, m := range t.GetMetrics() {
			fmt.Println(i, " ", j, " ", t.ObjectType, "-", m.Id)
			for _, d := range m.GetData() {
				for _, e := range d.GetEvents() {
					name := t.ObjectType + "-" + m.Id + "-" + e.GetType()
					aggregation := godruid.AggFiltered(
						godruid.FilterAnd(
							FilterHelper(metric, e),
						),
						&godruid.Aggregation{
							Type: "count",
							Name: name,
						},
					)
					aggregations = append(aggregations, aggregation)
				}
			}
		}
	}

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
