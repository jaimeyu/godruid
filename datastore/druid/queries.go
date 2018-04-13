package druid

import (
	"fmt"
	"strconv"
	"strings"

	pb "github.com/accedian/adh-gather/gathergrpc"
	"github.com/accedian/godruid"
)

const (
	TimeZoneUTC = "UTC"
)

// HistogramQuery - Count of metrics per bucket for given interval.
func HistogramQuery(tenant string, dataSource string, metric string, granularity string, direction string, interval string, resolution int32, granularityBuckets int32, vendor string, timeout int32) (*godruid.QueryTimeseries, error) {

	//peyo TODO need to figure out a better way than just appending Histo
	aggHist := godruid.AggHistoFold("thresholdBuckets", metric+"Histo", resolution, granularityBuckets, "0", "Infinity")

	return &godruid.QueryTimeseries{
		DataSource:  dataSource,
		Granularity: godruid.GranPeriod(granularity, TimeZoneUTC, ""),
		Context:     map[string]interface{}{"timeout": timeout},
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
	}, nil
}

// FilterHelper - helper function to select correct druid filter based on
// a given event and metric
func FilterHelper(metric string, e *pb.TenantThresholdProfileData_EventAttrMap) (*godruid.Filter, error) {

	event := e.GetEventAttrMap()

	upperStrict, err := strconv.ParseBool(event["upperStrict"])
	if err != nil && event["upperStrict"] != "" {
		return nil, fmt.Errorf("Invalid value for 'upperStrict' : %v. Must be a boolean", upperStrict)
	}

	lowerStrict, err := strconv.ParseBool(event["lowerStrict"])
	if err != nil && event["lowerStrict"] != "" {
		return nil, fmt.Errorf("Invalid value for 'lowerStrict' : %v. Must be a boolean", lowerStrict)
	}

	lowerLimit, err := strconv.ParseFloat(event["lowerLimit"], 32)
	if err != nil && event["lowerLimit"] != "" {
		return nil, fmt.Errorf("Invalid value for 'lowerLimit' : %v. Must be a number", lowerLimit)
	}

	upperLimit, err := strconv.ParseFloat(event["upperLimit"], 32)
	if err != nil && event["upperLimit"] != "" {
		return nil, fmt.Errorf("Invalid value for 'upperLimit' : %v. Must be a number", upperLimit)
	}

	if upperLimit != 0 && lowerLimit != 0 {
		return godruid.FilterLowerUpperBound(metric, godruid.NUMERIC, float32(lowerLimit), lowerStrict, float32(upperLimit), upperStrict), nil
	}

	if upperLimit != 0 {
		return godruid.FilterUpperBound(metric, godruid.NUMERIC, float32(upperLimit), upperStrict), nil
	}

	if lowerLimit != 0 {
		return godruid.FilterLowerBound(metric, godruid.NUMERIC, float32(lowerLimit), lowerStrict), nil
	}

	return nil, fmt.Errorf("Unable to consume threshold profile for: %v", metric)
}

// ThresholdCrossingQuery - Query that returns a count of events that crossed a thresholds for metric/thresholds
// defined by the supplied threshold profile..
func ThresholdCrossingQuery(tenant string, dataSource string, domains []string, metrics []string, granularity string, interval string, objectTypes []string, directions []string, thresholdProfile *pb.TenantThresholdProfileData, vendors []string, timeout int32) (*godruid.QueryTimeseries, error) {

	var aggregations []godruid.Aggregation

	aggregations = append(aggregations, godruid.AggCount("total"))

	for vk, v := range thresholdProfile.GetThresholds().GetVendorMap() {
		// if no vendors have been provided, use all of them, otherwise
		// only include the provided ones

		if vendors == nil || contains(vendors, vk) {
			for tk, t := range v.GetMonitoredObjectTypeMap() {
				// if no objectTypes have been provided, use all of them, otherwise
				// only include the provided ones
				if objectTypes == nil || contains(objectTypes, tk) {
					for mk, m := range t.GetMetricMap() {
						// if no metrics have been provided, use all of them, otherwise
						// only include the provided ones
						if metrics == nil || contains(metrics, mk) {
							for dk, d := range m.GetDirectionMap() {
								if directions == nil || contains(directions, dk) {
									for ek, e := range d.GetEventMap() {
										name := vk + "." + tk + "." + mk + "." + ek + "." + dk
										filter, err := FilterHelper(mk, e)
										if err != nil {
											return nil, err
										}
										aggregation := godruid.AggFiltered(
											godruid.FilterAnd(
												filter,
												godruid.FilterSelector("objectType", tk),
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
		}

	}

	return &godruid.QueryTimeseries{
		DataSource:  dataSource,
		Granularity: godruid.GranPeriod(granularity, TimeZoneUTC, ""),
		Context:     map[string]interface{}{"timeout": timeout, "skipEmptyBuckets": true},
		Filter: godruid.FilterAnd(
			godruid.FilterSelector("tenantId", strings.ToLower(tenant)),
			buildDomainFilter(domains),
		),
		Aggregations: aggregations,
		Intervals:    []string{interval}}, nil
}

// ThresholdCrossingByMonitoredObjectQuery - Query that returns a count of events that crossed a thresholds for metric/thresholds
// defined by the supplied threshold profile. Groups results my monitored object ID.
func ThresholdCrossingByMonitoredObjectQuery(tenant string, dataSource string, domains []string, metrics []string, granularity string, interval string, objectTypes []string, directions []string, thresholdProfile *pb.TenantThresholdProfileData, vendors []string, timeout int32) (*godruid.QueryGroupBy, error) {

	var aggregations []godruid.Aggregation

	aggregations = append(aggregations, godruid.AggCount("total"))

	for vk, v := range thresholdProfile.GetThresholds().GetVendorMap() {
		if vendors == nil || contains(vendors, vk) {
			for tk, t := range v.GetMonitoredObjectTypeMap() {
				if objectTypes == nil || contains(objectTypes, tk) {
					for mk, m := range t.GetMetricMap() {
						// if no metrics have been provided, use all of them, otherwise
						// only include the provided ones
						if metrics == nil || contains(metrics, mk) {
							for dk, d := range m.GetDirectionMap() {
								if directions == nil || contains(directions, dk) {
									for ek, e := range d.GetEventMap() {
										name := vk + "." + tk + "." + mk + "." + ek + "." + dk
										filter, err := FilterHelper(mk, e)
										if err != nil {
											return nil, err
										}
										aggregation := godruid.AggFiltered(
											godruid.FilterAnd(
												filter,
												godruid.FilterSelector("objectType", tk),
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
		}
	}

	return &godruid.QueryGroupBy{
		DataSource:   dataSource,
		Granularity:  godruid.GranPeriod(granularity, TimeZoneUTC, ""),
		Context:      map[string]interface{}{"timeout": timeout},
		Aggregations: aggregations,
		Filter: godruid.FilterAnd(
			godruid.FilterSelector("tenantId", strings.ToLower(tenant)),
			buildDomainFilter(domains),
		),
		Intervals: []string{interval},
		Dimensions: []godruid.DimSpec{
			godruid.Dimension{
				Dimension:  "monitoredObjectId",
				OutputName: "monitoredObjectId",
			},
		}}, nil
}

//RawMetricsQuery  - Query that returns a raw metric values
func RawMetricsQuery(tenant string, dataSource string, metrics []string, interval string, objectType string, direction string, monitoredObjects []string, timeout int32, granularity string) (*godruid.QueryTimeseries, error) {

	var aggregations []godruid.Aggregation

	for _, monObj := range monitoredObjects {
		for _, metric := range metrics {
			aggregationMax := godruid.AggFiltered(
				godruid.FilterSelector("monitoredObjectId", monObj),
				&godruid.Aggregation{
					Type:      "doubleMax",
					Name:      monObj + "." + metric,
					FieldName: metric,
				},
			)
			aggregations = append(aggregations, aggregationMax)
		}
	}

	return &godruid.QueryTimeseries{
		DataSource:   dataSource,
		Granularity:  godruid.GranPeriod(granularity, TimeZoneUTC, ""),
		Context:      map[string]interface{}{"timeout": timeout, "skipEmptyBuckets": true},
		Aggregations: aggregations,
		Filter: godruid.FilterAnd(
			godruid.FilterSelector("tenantId", strings.ToLower(tenant)),
			godruid.FilterSelector("objectType", objectType),
			godruid.FilterSelector("direction", direction),
		),
		Intervals: []string{interval},
	}, nil
}

func buildDomainFilter(domains []string) *godruid.Filter {
	if len(domains) < 1 {
		return nil
	}
	filters := make([]*godruid.Filter, len(domains))
	for i, domain := range domains {
		filters[i] = godruid.FilterSelector("domains", domain)
	}
	return godruid.FilterOr(filters...)

}
