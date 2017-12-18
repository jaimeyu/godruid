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
func FilterHelper(metric string, e *pb.TenantThresholdProfile_EventAttrMap) *godruid.Filter {

	// if e.UpperBound != 0 && e.LowerBound != 0 {
	// 	return godruid.FilterLowerUpperBound(metric, godruid.NUMERIC, e.LowerBound, e.LowerStrict, e.UpperBound, e.UpperStrict)
	// }

	// if e.UpperBound != 0 {
	// 	return godruid.FilterUpperBound(metric, godruid.NUMERIC, e.UpperBound, e.UpperStrict)
	// }

	// if e.LowerBound != 0 {
	// 	return godruid.FilterLowerBound(metric, godruid.NUMERIC, e.LowerBound, e.LowerStrict)
	// }

	return godruid.FilterLowerBound(metric, godruid.NUMERIC, 1000, true)
	// return nil
}

// ThresholdCrossingQuery - Query that returns a count of events that crossed a thresholds for metric/thresholds
// defined by the supplied threshold profile..
func ThresholdCrossingQuery(tenant string, dataSource string, metric string, granularity string, interval string, objectType string, direction string, thresholdProfile *pb.TenantThresholdProfile) *godruid.QueryTimeseries {

	var aggregations []godruid.Aggregation
	metrics := strings.Split(metric, ",")
	objectTypes := strings.Split(objectType, ",")
	directions := strings.Split(direction, ",")

	aggregations = append(aggregations, godruid.AggCount("total"))

	// peyo TODO don't hardcode vendor
	for tk, t := range thresholdProfile.GetThresholds().GetVendorMap()["accedian"].GetMonitoredObjectTypeMap() {
		// if no objectTypes have been provided, use all of them, otherwise
		// only include the provided ones
		if contains(objectTypes, tk) || objectType == "" {
			for mk, m := range t.GetMetricMap() {
				// if no metrics have been provided, use all of them, otherwise
				// only include the provided ones
				if contains(metrics, mk) || metric == "" {
					for dk, d := range m.GetDirectionMap() {
						if contains(directions, dk) || direction == "" {
							for ek, e := range d.GetEventMap() {
								name := tk + "." + mk + "." + ek + "." + dk
								aggregation := godruid.AggFiltered(
									godruid.FilterAnd(
										FilterHelper(mk, e),
										godruid.FilterSelector("obejctType", tk),
										godruid.FilterSelector("tenantId", tenant),
										godruid.FilterSelector("direction", dk),
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
