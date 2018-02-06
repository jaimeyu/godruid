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
func ThresholdCrossingQuery(tenant string, dataSource string, metric string, granularity string, interval string, objectType string, direction string, thresholdProfile *pb.TenantThresholdProfileData, vendor string, timeout int32) (*godruid.QueryTimeseries, error) {

	var aggregations []godruid.Aggregation
	metrics := strings.Split(metric, ",")
	objectTypes := strings.Split(objectType, ",")
	directions := strings.Split(direction, ",")
	vendors := strings.Split(vendor, ",")

	aggregations = append(aggregations, godruid.AggCount("total"))

	for vk, v := range thresholdProfile.GetThresholds().GetVendorMap() {
		// if no vendors have been provided, use all of them, otherwise
		// only include the provided ones

		if contains(vendors, vk) || vendor == "" {
			for tk, t := range v.GetMonitoredObjectTypeMap() {
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
		DataSource:   dataSource,
		Granularity:  godruid.GranPeriod(granularity, TimeZoneUTC, ""),
		Context:      map[string]interface{}{"timeout": timeout, "skipEmptyBuckets": true},
		Filter:       godruid.FilterSelector("tenantId", strings.ToLower(tenant)),
		Aggregations: aggregations,
		Intervals:    []string{interval}}, nil
}

// ThresholdCrossingByMonitoredObjectQuery - Query that returns a count of events that crossed a thresholds for metric/thresholds
// defined by the supplied threshold profile. Groups results my monitored object ID.
func ThresholdCrossingByMonitoredObjectQuery(tenant string, dataSource string, metric string, granularity string, interval string, objectType string, direction string, thresholdProfile *pb.TenantThresholdProfileData, vendor string, timeout int32) (*godruid.QueryGroupBy, error) {

	var aggregations []godruid.Aggregation
	metrics := strings.Split(metric, ",")
	directions := strings.Split(direction, ",")
	vendors := strings.Split(vendor, ",")
	objectTypes := strings.Split(objectType, ",")

	aggregations = append(aggregations, godruid.AggCount("total"))

	for vk, v := range thresholdProfile.GetThresholds().GetVendorMap() {
		if contains(vendors, vk) || vendor == "" {
			for tk, t := range v.GetMonitoredObjectTypeMap() {
				if contains(objectTypes, tk) || objectType == "" {
					for mk, m := range t.GetMetricMap() {
						// if no metrics have been provided, use all of them, otherwise
						// only include the provided ones
						if contains(metrics, mk) || metric == "" {
							for dk, d := range m.GetDirectionMap() {
								if contains(directions, dk) || direction == "" {
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
		Filter:       godruid.FilterSelector("tenantId", strings.ToLower(tenant)),
		Intervals:    []string{interval},
		Dimensions: []godruid.DimSpec{
			godruid.Dimension{
				Dimension:  "monitoredObjectId",
				OutputName: "monitoredObjectId",
			},
		}}, nil
}

//RawMetricsQuery  - Query that returns a raw metric values
func RawMetricsQuery(tenant string, dataSource string, metric string, interval string, objectType string, direction string, monitoredObjectId string, timeout int32) (*godruid.QuerySelect, error) {

	metrics := strings.Split(metric, ",")

	return &godruid.QuerySelect{
		DataSource:  dataSource,
		Granularity: godruid.GranAll,
		Context:     map[string]interface{}{"timeout": timeout},
		Filter: godruid.FilterAnd(
			godruid.FilterSelector("tenantId", strings.ToLower(tenant)),
			godruid.FilterSelector("monitoredObjectId", monitoredObjectId),
			godruid.FilterSelector("objectType", objectType),
			godruid.FilterSelector("direction", direction),
		),
		Intervals:  []string{interval},
		Metrics:    metrics,
		PagingSpec: map[string]interface{}{"threshold": 2000}, // TODO peyo hardcoding threshold, will have to find out what value we want here
		Dimensions: []godruid.DimSpec{
			godruid.Dimension{
				Dimension:  "monitoredObjectId",
				OutputName: "monitoredObjectId",
			},
		},
	}, nil
}
