package druid

import (
	"fmt"
	"strconv"
	"strings"

	pb "github.com/accedian/adh-gather/gathergrpc"
	"github.com/accedian/adh-gather/logger"
	"github.com/accedian/adh-gather/models/metrics"
	"github.com/accedian/godruid"
)

type TimeBucket int

const (
	TimeZoneUTC                = "UTC"
	Granularity_All            = "all"
	HourOfDay       TimeBucket = 0
	DayOfWeek       TimeBucket = 1
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

func HistogramCustomQuery(tenant string, domains []string, dataSource string, interval string, granularity string, timeout int32, metricsOne []map[string]interface{}) (*godruid.QueryTimeseries, error) {

	var aggregations []godruid.Aggregation

	for _, met := range metricsOne {

		metName := met["name"].(string)
		metVendor := met["vendor"].(string)
		metDirection := met["direction"].(string)
		metObjectType := met["objectType"].(string)

		for _, bucket := range met["buckets"].([]interface{}) {
			bucketMap := bucket.(map[string]interface{})
			metUpper := bucketMap["upper"].(float64)
			metLower := bucketMap["lower"].(float64)
			metIndex := bucketMap["index"]

			name := fmt.Sprintf("%s.%s.%s.%s.%s", metVendor, metObjectType, metName, metDirection, metIndex)

			filter := godruid.FilterLowerUpperBound(metName, godruid.NUMERIC, float32(metLower), false, float32(metUpper), true)

			aggregation := godruid.AggFiltered(
				godruid.FilterAnd(
					filter,
					godruid.FilterSelector("objectType", metObjectType),
					godruid.FilterSelector("direction", metDirection),
				),
				&godruid.Aggregation{
					Type: "count",
					Name: name,
				},
			)
			aggregations = append(aggregations, aggregation)
		}
	}

	return &godruid.QueryTimeseries{
		DataSource:  dataSource,
		Granularity: toGranularity(granularity),
		Context:     map[string]interface{}{"timeout": timeout, "skipEmptyBuckets": true},
		Filter: godruid.FilterAnd(
			godruid.FilterSelector("tenantId", strings.ToLower(tenant)),
			//buildDomainFilter(tenant, domains), //TODO PUT THIS BACK IN!!!!
		),
		Aggregations: aggregations,
		Intervals:    []string{interval}}, nil
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
			buildDomainFilter(tenant, domains),
		),
		Aggregations: aggregations,
		Intervals:    []string{interval}}, nil
}

func SLAViolationsQuery(tenant string, dataSource string, domains []string, granularity string, interval string, thresholdProfile *pb.TenantThresholdProfileData, timeout int32) (*godruid.QueryTimeseries, error) {
	var aggregations []godruid.Aggregation
	var postAggregations []godruid.PostAggregation
	var violationCountAggs []string
	var totalDurationAggs []string
	var violationDurationAggs []string
	var objectDirectionFilters []*godruid.Filter

	type objectTypeDirectionFilters struct {
		BaseFilter       *godruid.Filter
		ThresholdFilters []*godruid.Filter
	}
	for vk, v := range thresholdProfile.GetThresholds().GetVendorMap() {
		for tk, t := range v.GetMonitoredObjectTypeMap() {
			perDirectionFilters := make(map[string]*objectTypeDirectionFilters)

			for mk, m := range t.GetMetricMap() {
				for dk, d := range m.GetDirectionMap() {
					for ek, e := range d.GetEventMap() {
						if ek != "sla" {
							continue
						}

						objectTypeAndDirectionFilter := godruid.FilterAnd(
							godruid.FilterSelector("objectType", tk),
							godruid.FilterSelector("direction", dk),
						)

						thresholdFilter, err := FilterHelper(mk, e)
						if err != nil {
							return nil, err
						}

						dirFilters, ok := perDirectionFilters[vk+"."+tk+"."+dk]
						if !ok {
							perDirectionFilters[vk+"."+tk+"."+dk] = &objectTypeDirectionFilters{objectTypeAndDirectionFilter, []*godruid.Filter{thresholdFilter}}
						} else {
							dirFilters.ThresholdFilters = append(dirFilters.ThresholdFilters, thresholdFilter)
						}

						aggNamePrefix := vk + "." + tk + "." + mk + "." + ek + "." + dk
						// Count violations for this metric
						aggregations = append(aggregations, godruid.AggFiltered(
							godruid.FilterAnd(
								thresholdFilter,
								objectTypeAndDirectionFilter,
							),
							&godruid.Aggregation{
								Type: "count",
								Name: aggNamePrefix + ".violationCount",
							},
						))
						violationCountAggs = append(violationCountAggs, aggNamePrefix+".violationCount")

						// Sum the total duration this metric was measured.
						aggregations = append(aggregations, godruid.AggFiltered(
							objectTypeAndDirectionFilter,
							&godruid.Aggregation{
								Type:      "longSum",
								FieldName: "duration",
								Name:      aggNamePrefix + ".totalDuration",
							},
						))

						// Sum the duration while this metric was in violation.
						aggregations = append(aggregations, godruid.AggFiltered(
							godruid.FilterAnd(
								thresholdFilter,
								objectTypeAndDirectionFilter,
							),
							&godruid.Aggregation{
								Type:      "longSum",
								FieldName: "duration",
								Name:      aggNamePrefix + ".violationDuration",
							},
						))

					}
				}
			}

			if len(perDirectionFilters) > 0 {

				// Sum the duration per vendor/objecttype/direction
				for k, v := range perDirectionFilters {
					objectDirectionFilters = append(objectDirectionFilters, v.BaseFilter)

					aggregations = append(aggregations, godruid.AggFiltered(
						v.BaseFilter,
						&godruid.Aggregation{
							Type:      "longSum",
							FieldName: "duration",
							Name:      k + ".totalDuration",
						},
					))
					totalDurationAggs = append(totalDurationAggs, k+".totalDuration")

					if len(v.ThresholdFilters) > 0 {
						// Sum the violation duration per vendor/objecttype/direction
						aggregations = append(aggregations, godruid.AggFiltered(
							godruid.FilterAnd(
								v.BaseFilter,
								godruid.FilterOr(v.ThresholdFilters...),
							),
							&godruid.Aggregation{
								Type:      "longSum",
								FieldName: "duration",
								Name:      k + ".violationDuration",
							},
						))
						violationDurationAggs = append(violationDurationAggs, k+".violationDuration")
					}

				}

			}
		}
	}

	// Count the monitored objects
	if len(objectDirectionFilters) > 0 {

		aggregations = append(aggregations, godruid.AggFiltered(
			godruid.FilterOr(
				objectDirectionFilters...,
			),
			&godruid.Aggregation{
				Type:       "cardinality",
				Name:       "objectCount",
				FieldNames: []string{"monitoredObjectId"},
				ByRow:      false,
				Round:      true,
			},
		))
	}

	if len(violationCountAggs) > 0 {
		// Sum the violation count per metric to get an overall total.
		postAggregations = append(postAggregations, godruid.PostAggArithmetic(
			"totalViolationCount",
			"+",
			buildPostAggregationFields(violationCountAggs)))
	}
	if len(violationDurationAggs) > 0 {
		// Sum the violation duration per metric to get an overal violation duration.
		postAggregations = append(postAggregations, godruid.PostAggArithmetic(
			"totalViolationDuration",
			"+",
			buildPostAggregationFields(violationDurationAggs)))
	}
	if len(totalDurationAggs) > 0 {
		// Sum the total duration per metric to get an overall total duration.
		postAggregations = append(postAggregations, godruid.PostAggArithmetic(
			"totalDuration",
			"+",
			buildPostAggregationFields(totalDurationAggs)))
	}

	return &godruid.QueryTimeseries{
		QueryType:   godruid.TIMESERIES,
		DataSource:  dataSource,
		Granularity: toGranularity(granularity),
		Context:     map[string]interface{}{"timeout": timeout, "skipEmptyBuckets": true},
		Filter: godruid.FilterAnd(
			godruid.FilterSelector("tenantId", strings.ToLower(tenant)),
			buildDomainFilter(tenant, domains),
		),
		Aggregations:     aggregations,
		PostAggregations: postAggregations,
		Intervals:        []string{interval}}, nil
}

func SLATimeBucketQuery(tenant string, dataSource string, domains []string, timeBucket TimeBucket, timeZone string, vendor, objectType, metric, direction, event string, eventAttr *pb.TenantThresholdProfileData_EventAttrMap, granularity string, interval string, timeout int32) (*godruid.QueryTopN, error) {
	var aggregations []godruid.Aggregation
	var dimension godruid.DimSpec
	threshold := 0
	if timeBucket == DayOfWeek {
		threshold = 7
		dimension = godruid.TimeExtractionDimensionSpec{
			Type:       "extraction",
			Dimension:  "__time",
			OutputName: "dayOfWeek",
			ExtractionFunction: godruid.TimeExtractionFn{
				Type:     "timeFormat",
				Format:   "e",
				TimeZone: timeZone,
				Locale:   "en",
			},
		}
	} else if timeBucket == HourOfDay {
		threshold = 24
		dimension = godruid.TimeExtractionDimensionSpec{
			Type:       "extraction",
			Dimension:  "__time",
			OutputName: "hourOfDay",
			ExtractionFunction: godruid.TimeExtractionFn{
				Type:     "timeFormat",
				Format:   "HH",
				TimeZone: timeZone,
				Locale:   "en",
			},
		}
	} else {
		return nil, fmt.Errorf("Invalid value for 'timeBucket' : %v", timeBucket)
	}

	thresholdFilter, err := FilterHelper(metric, eventAttr)
	if err != nil {
		return nil, err
	}

	// Count violations for this metric
	countName := vendor + "." + objectType + "." + metric + "." + event + "." + direction + ".violationCount"
	aggregations = append(aggregations, godruid.Aggregation{
		Type: "count",
		Name: countName,
	})

	return &godruid.QueryTopN{
		DataSource:  dataSource,
		Granularity: toGranularity(granularity),
		Context:     map[string]interface{}{"timeout": timeout},
		Filter: godruid.FilterAnd(
			godruid.FilterSelector("tenantId", strings.ToLower(tenant)),
			buildDomainFilter(tenant, domains),
			thresholdFilter,
			godruid.FilterSelector("objectType", objectType),
			godruid.FilterSelector("direction", direction),
		),
		Metric:       countName,
		Dimension:    dimension,
		Threshold:    threshold,
		Aggregations: aggregations,
		Intervals:    []string{interval},
	}, nil

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
			buildDomainFilter(tenant, domains),
		),
		Intervals: []string{interval},
		Dimensions: []godruid.DimSpec{
			godruid.Dimension{
				Dimension:  "monitoredObjectId",
				OutputName: "monitoredObjectId",
			},
		}}, nil
}

// ThresholdCrossingByMonitoredObjectQuery - Query that returns a count of events that crossed a thresholds for metric/thresholds
// defined by the supplied threshold profile. Groups results my monitored object ID.
func ThresholdCrossingByMonitoredObjectTopNQuery(tenant string, dataSource string, domains []string, metric string, granularity string, interval string, objectType string, direction string, thresholdProfile *pb.TenantThresholdProfileData, vendor string, timeout int32, numResults int32) (*godruid.QueryTopN, error) {

	var aggregations []godruid.Aggregation
	var postAggregations []godruid.PostAggregation

	var eventWeights = make(map[string]float32)
	eventWeights["minor"] = 0.0001
	eventWeights["major"] = 0.001
	eventWeights["critical"] = 1

	aggregations = append(aggregations, godruid.AggCount("total"))

	vendorMap := thresholdProfile.GetThresholds().GetVendorMap()
	events := vendorMap[vendor].GetMonitoredObjectTypeMap()[objectType].GetMetricMap()[metric].GetDirectionMap()[direction].GetEventMap()

	for ek, e := range events {
		name := ek
		filter, err := FilterHelper(metric, e)
		if err != nil {
			return nil, err
		}
		aggregation := godruid.AggFiltered(
			godruid.FilterAnd(
				filter,
			),
			&godruid.Aggregation{
				Type: "count",
				Name: name,
			},
		)

		postAggregation := godruid.PostAggArithmetic("", "*", []godruid.PostAggregation{
			godruid.PostAggConstant("", eventWeights[ek]),
			godruid.PostAggFieldAccessor(ek),
		})

		postAggregations = append(postAggregations, postAggregation)
		aggregations = append(aggregations, aggregation)
	}

	var scoredPostAggregation []godruid.PostAggregation
	if len(postAggregations) == 0 {
		scoredPostAggregation = []godruid.PostAggregation{}
	} else {
		scoredPostAggregation = []godruid.PostAggregation{
			godruid.PostAggArithmetic("scored", "+", postAggregations),
		}
	}

	return &godruid.QueryTopN{
		DataSource:   dataSource,
		Granularity:  toGranularity(granularity),
		Context:      map[string]interface{}{"timeout": timeout},
		Aggregations: aggregations,
		Filter: godruid.FilterAnd(
			godruid.FilterSelector("tenantId", strings.ToLower(tenant)),
			buildDomainFilter(tenant, domains),
			godruid.FilterSelector("objectType", objectType),
			godruid.FilterSelector("direction", direction),
		),
		PostAggregations: scoredPostAggregation,
		Intervals:        []string{interval},
		Metric:           "scored",
		Threshold:        int(numResults),
		Dimension:        "monitoredObjectId",
	}, nil
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

//AggMetricsQuery  - Query that returns a aggregated metric values
func AggMetricsQuery(tenant string, dataSource string, interval string, domains []string, aggregationFunc metrics.AggregationSpec, metrics []metrics.MetricIdentifier, timeout int32, granularity string) (*godruid.QueryTimeseries, *PostProcessor, error) {

	var aggregations []godruid.Aggregation
	var pp PostProcessor
	postAggs := []godruid.PostAggregation{}

	keyToDrop := []string{}
	countKeys := map[string][]string{}

	for _, metric := range metrics {
		countName := metric.Name + "Count"
		keyToDrop = append(keyToDrop, countName)
		countKeys[countName] = []string{metric.Name}
		aggregations = append(aggregations, buildMetricAggregation("count", &metric, countName))
		if aggregationFunc.Name == "max" {
			aggregations = append(aggregations, buildMetricAggregation("doubleMax", &metric))

		} else if aggregationFunc.Name == "min" {
			aggregations = append(aggregations, buildMetricAggregation("doubleMin", &metric))

		} else if aggregationFunc.Name == "avg" {

			aggregations = append(aggregations, buildMetricAggregation("doubleSum", &metric, metric.Name+"Sum"))

			keyToDrop = append(keyToDrop, metric.Name+"Sum")

			postAgg := godruid.PostAggArithmetic(
				metric.Name,
				"/",
				[]godruid.PostAggregation{godruid.PostAggFieldAccessor(metric.Name + "Sum"), godruid.PostAggFieldAccessor(metric.Name + "Count")},
			)
			postAggs = append(postAggs, postAgg)
		} else {
			return nil, nil, fmt.Errorf("Invalid value for 'aggregation' : %v", aggregationFunc)
		}

	}

	// Drop the intermediate sum and count aggregations after the response returns from druid.
	// There doesn't seem to be an option in the druid query to do this.
	pp = DropKeysPostprocessor{
		keysToDrop: keyToDrop,
		countKeys:  countKeys,
	}

	return &godruid.QueryTimeseries{
		DataSource:   dataSource,
		Granularity:  toGranularity(granularity),
		Context:      map[string]interface{}{"timeout": timeout, "skipEmptyBuckets": true},
		Aggregations: aggregations,
		Filter: godruid.FilterAnd(
			godruid.FilterSelector("tenantId", strings.ToLower(tenant)),
			buildDomainFilter(tenant, domains),
		),
		Intervals:        []string{interval},
		PostAggregations: postAggs,
	}, &pp, nil
}

func buildMetricAggregation(aggType string, metric *metrics.MetricIdentifier, name ...string) godruid.Aggregation {
	var aggName string
	if len(name) == 0 {
		aggName = metric.Name
	} else {
		aggName = name[0]
	}
	return godruid.AggFiltered(
		godruid.FilterAnd(
			godruid.FilterSelector("objectType", metric.ObjectType),
			godruid.FilterSelector("direction", metric.Direction),
		),
		&godruid.Aggregation{
			Type:      aggType,
			Name:      aggName,
			FieldName: metric.Name,
		})

}

func buildDomainFilter(tenantID string, domains []string) *godruid.Filter {
	if len(domains) < 1 {
		return nil
	}

	filters := make([]*godruid.Filter, len(domains))
	atLeastOneDomainFilter := false
	for i, domID := range domains {
		var ef godruid.ExtractionFn

		lookupName, exists := getLookupName("dom", tenantID, domID)
		if !exists {
			logger.Log.Warningf("No lookup (%s) found for domain ID %s. It will be excluded from the domain filter", lookupName, domID)
			continue
		}

		atLeastOneDomainFilter = true
		ef = godruid.RegisteredLookupExtractionFn{
			Type:   "registeredLookup",
			Lookup: lookupName,
		}

		filters[i] = &godruid.Filter{
			Type:         "selector",
			Dimension:    "monitoredObjectId",
			Value:        domID,
			ExtractionFn: &ef,
		}
	}

	if !atLeastOneDomainFilter {
		// This is a hack to get around no domain lookups being ready yet.  Basically want to
		// create an 'always false filter".
		logger.Log.Debugf("No domains found in cached using false filter")
		return godruid.FilterAnd(
			godruid.FilterSelector("tenantId", strings.ToLower(tenantID)),
			godruid.FilterNot(godruid.FilterSelector("tenantId", strings.ToLower(tenantID))),
		)
	}

	return godruid.FilterOr(filters...)
}

func buildOrFilter(dimensionName string, values []string) *godruid.Filter {
	if len(values) < 1 {
		return nil
	}
	filters := make([]*godruid.Filter, len(values))
	for i, value := range values {
		filters[i] = godruid.FilterSelector(dimensionName, value)
	}
	return godruid.FilterOr(filters...)

}

func buildPostAggregationFields(fieldNames []string) []godruid.PostAggregation {

	fields := make([]godruid.PostAggregation, len(fieldNames))
	for i, name := range fieldNames {
		fields[i] = godruid.PostAggFieldAccessor(name)
	}

	if len(fields) == 1 {
		// If we are only summing 1 value, satisfy the post aggregator by adding 0.
		fields = append(fields, godruid.PostAggConstant("const", 0))
	}
	return fields
}

func toGranularity(granularityStr string) godruid.Granlarity {
	if granularityStr == Granularity_All {
		return godruid.GranAll
	}
	return godruid.GranPeriod(granularityStr, TimeZoneUTC, "")
}
