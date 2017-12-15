package druid

import (
	"strings"

	pb "github.com/accedian/adh-gather/gathergrpc"
	"github.com/accedian/godruid"
)

const (
	TimeZoneUTC = "UTC"
)

// HistogramQuery - Count of metrics per bucket for given interval.
func HistogramQuery(tenant string, dataSource string, metric string, granularity string, direction string, interval string, resolution int32, granularityBuckets int32) *godruid.QueryTimeseries {

	//peyo TODO need to figure out a better way than just appending Histo
	aggHist := godruid.AggHistoFold("thresholdBuckets", metric+"Histo", resolution, granularityBuckets, "0", "Infinity")

	return &godruid.QueryTimeseries{
		DataSource:  dataSource,
		Granularity: godruid.GranPeriod(granularity, TimeZoneUTC, ""),
		Context:     map[string]interface{}{"timeout": 60000},
		Aggregations: []godruid.Aggregation{
			godruid.AggFiltered(
				godruid.FilterAnd(
					godruid.FilterSelector("tenantId", tenant),
					godruid.FilterSelector("direction", direction),
				),
				&aggHist,
			),
		},
		Intervals: []string{interval},
	}
}

// FilterHelper - helper function to select correct druid filter based on
// a given event and metric
func FilterHelper(metric string, e *pb.TenantEvent) *godruid.Filter {

	if e.UpperBound != 0 && e.LowerBound != 0 {
		return godruid.FilterLowerUpperBound(metric, godruid.NUMERIC, e.LowerBound, e.LowerStrict, e.UpperBound, e.UpperStrict)
	}

	if e.UpperBound != 0 {
		return godruid.FilterUpperBound(metric, godruid.NUMERIC, e.UpperBound, e.UpperStrict)
	}

	if e.LowerBound != 0 {
		return godruid.FilterLowerBound(metric, godruid.NUMERIC, e.LowerBound, e.LowerStrict)
	}
	return nil
}

// ThresholdCrossingQuery - Query that returns a count of events that crossed a thresholds for metric/thresholds
// defined by the supplied threshold profile..
func ThresholdCrossingQuery(tenant string, dataSource string, metric string, granularity string, interval string, objectType string, direction string, thresholdProfile *pb.TenantThresholdProfile) *godruid.QueryTimeseries {

	var aggregations []godruid.Aggregation
	metrics := strings.Split(metric, ",")
	objectTypes := strings.Split(objectType, ",")
	directions := strings.Split(direction, ",")

	aggregations = append(aggregations, godruid.AggCount("total"))

	for _, t := range thresholdProfile.GetThresholds() {
		// if no objectTypes have been provided, use all of them, otherwise
		// only include the provided ones
		if contains(objectTypes, t.GetObjectType()) || objectType == "" {
			for _, m := range t.GetMetrics() {
				// if no metrics have been provided, use all of them, otherwise
				// only include the provided ones
				if contains(metrics, m.GetId()) || metric == "" {
					for _, d := range m.GetData() {
						if contains(directions, d.GetDirection()) || direction == "" {
							for _, e := range d.GetEvents() {
								name := t.GetObjectType() + "." + m.GetId() + "." + e.GetType() + "." + d.GetDirection()
								aggregation := godruid.AggFiltered(
									godruid.FilterAnd(
										FilterHelper(m.GetId(), e),
										godruid.FilterSelector("obejctType", t.GetObjectType()),
										godruid.FilterSelector("tenantId", tenant),
										godruid.FilterSelector("direction", d.GetDirection()),
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
			}
		}
	}

	return &godruid.QueryTimeseries{
		DataSource:   dataSource,
		Granularity:  godruid.GranPeriod(granularity, TimeZoneUTC, ""),
		Context:      map[string]interface{}{"timeout": 60000},
		Aggregations: aggregations,
		Intervals:    []string{interval},
	}
}
