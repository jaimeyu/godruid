package druid

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/accedian/adh-gather/config"
	"github.com/accedian/adh-gather/gather"
	"github.com/accedian/adh-gather/logger"
	"github.com/accedian/adh-gather/models"
	"github.com/accedian/adh-gather/models/metrics"
	"github.com/accedian/adh-gather/models/tenant"
	"github.com/accedian/godruid"

	db "github.com/accedian/adh-gather/datastore"
	pb "github.com/accedian/adh-gather/gathergrpc"

	"github.com/satori/go.uuid"
)

const (
	ThresholdCrossingReport = "threshold-crossing-report"
	EventDistribution       = "event-distribution"
	RawMetrics              = "raw-metrics"
	SLAReport               = "sla-report"
	TopNForMetric           = "top-n"
)

// DruidDatastoreClient - struct responsible for handling
// database operations for druid
type DruidDatastoreClient struct {
	server            string
	cfg               config.Provider
	dClient           godruid.Client
	AuthToken         string
	numRetries        int
	coordinatorServer string
	coordinatorPort   string
}

type ThresholdCrossingByMonitoredObjectResponse struct {
	Version   string
	Timestamp string
	Event     map[string]interface{}
}

type TopNThresholdCrossingByMonitoredObjectResponse struct {
	Timestamp string
	Result    []map[string]interface{}
}

type RawMetricsResponse struct {
	Timestamp string                 `json:"timestamp"`
	Result    map[string]interface{} `json:"result"`
}

type AggMetricsResponse struct {
	Timestamp string
	Result    map[string]interface{}
}

func makeHttpClient() *http.Client {
	// By default, use 60 second timeout unless specified otherwise
	// by the caller
	clientTimeout := 60 * time.Second

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	httpClient := &http.Client{
		Timeout:   clientTimeout,
		Transport: tr,
	}

	return httpClient
}

func (dc *DruidDatastoreClient) executeQuery(query godruid.Query) ([]byte, error) {

	client := dc.dClient

	err := client.Query(query, dc.AuthToken)

	if err != nil {
		if strings.Contains(err.Error(), "401") {
			logger.Log.Info("Auth token expired, refreshing token")
			dc.AuthToken = GetAuthCode(dc.cfg)
			err_retry := client.Query(query, dc.AuthToken)
			if err_retry != nil {
				logger.Log.Errorf("Druid Query RETRY failed due to: %s", err)
				return nil, err_retry
			}
			return query.GetRawJSON(), nil
		}
		logger.Log.Errorf("Druid Query failed due to: %s", err)
		return nil, err
	}

	return query.GetRawJSON(), nil
}

// NewDruidDatasctoreClient - Constructor for DruidDatastoreClient object
// initializes the godruid client, and retrieves auth token
// peyo TODO: the auth functionality here needs to be changed, this is only valid for dev
func NewDruidDatasctoreClient() *DruidDatastoreClient {
	cfg := gather.GetConfig()
	server := cfg.GetString(gather.CK_druid_broker_server.String())
	port := cfg.GetString(gather.CK_druid_broker_port.String())
	client := godruid.Client{
		Url:        server + ":" + port,
		Debug:      true,
		HttpClient: makeHttpClient(),
	}

	return &DruidDatastoreClient{
		cfg:               cfg,
		server:            server,
		dClient:           client,
		AuthToken:         GetAuthCode(cfg),
		coordinatorServer: cfg.GetString(gather.CK_druid_coordinator_server.String()),
		coordinatorPort:   cfg.GetString(gather.CK_druid_coordinator_port.String()),
	}
}

// peyo TODO: implement this query
func (dc *DruidDatastoreClient) GetHistogram(request *pb.HistogramRequest) (map[string]interface{}, error) {

	logger.Log.Debugf("Calling GetHistogram for request: %v", models.AsJSONString(request))
	table := dc.cfg.GetString(gather.CK_druid_broker_table.String())

	// peyo TODO we should have a better way to handle default query params
	timeout := request.GetTimeout()
	if timeout == 0 {
		timeout = 5000
	}

	query, err := HistogramQuery(request.GetTenant(), table, request.Metric, request.Granularity, request.Direction, request.Interval, request.Resolution, request.GranularityBuckets, request.GetVendor(), timeout)

	if err != nil {
		return nil, err
	}

	logger.Log.Debugf("Querying Druid for %s with query: %v", db.HistogramStr, models.AsJSONString(query))
	response, err := dc.executeQuery(query)

	if err != nil {
		return nil, err
	}

	histogram := []*pb.Histogram{}

	err = json.Unmarshal(response, &histogram)
	if err != nil {
		return nil, err
	}

	logger.Log.Debugf("Response from druid for %s: %v", db.HistogramStr, models.AsJSONString(histogram))

	resp := &pb.HistogramResponse{
		Data: histogram,
	}

	// peyo TODO: need to figure out where to get this ID and Type from.
	uuid := uuid.NewV4()
	data := make([]*pb.HistogramResponse, 0)
	data = append(data, resp)
	rr := map[string]interface{}{
		"data": map[string]interface{}{
			"id":         uuid.String(),
			"type":       ThresholdCrossingReport,
			"attributes": data,
		},
	}

	return rr, nil
}

// GetThresholdCrossing - Executes a 'threshold crossing' query against druid. Wraps the
// result in a JSON API wrapper.
// peyo TODO: probably don't need to wrap JSON API here...should maybe do it elsewhere
func (dc *DruidDatastoreClient) GetThresholdCrossing(request *pb.ThresholdCrossingRequest, thresholdProfile *pb.TenantThresholdProfile) (map[string]interface{}, error) {

	logger.Log.Debugf("Calling GetThresholdCrossing for request: %v", models.AsJSONString(request))
	table := dc.cfg.GetString(gather.CK_druid_broker_table.String())

	// peyo TODO we should have a better way to handle default query params
	timeout := request.GetTimeout()
	if timeout == 0 {
		timeout = 5000
	}

	query, err := ThresholdCrossingQuery(request.GetTenant(), table, request.Domain, request.Metric, request.Granularity, request.Interval, request.ObjectType, request.Direction, thresholdProfile.Data, request.GetVendor(), timeout)
	if err != nil {
		return nil, err
	}

	logger.Log.Debugf("Querying Druid for %s with query: %v", db.ThresholdCrossingStr, models.AsJSONString(query))
	response, err := dc.executeQuery(query)

	if err != nil {
		return nil, err
	}

	thresholdCrossing := []*pb.ThresholdCrossing{}
	err = json.Unmarshal(response, &thresholdCrossing)
	if err != nil {
		return nil, err
	}
	logger.Log.Debugf("Response from druid for %s: %v", db.ThresholdCrossingStr, models.AsJSONString(thresholdCrossing))

	formattedJSON, err := reformatThresholdCrossingResponse(thresholdCrossing)
	if err != nil {
		return nil, err
	}

	// peyo TODO: need to figure out where to get this ID and Type from.
	uuid := uuid.NewV4()
	data := []map[string]interface{}{}
	data = append(data, map[string]interface{}{
		"id":         uuid.String(),
		"type":       ThresholdCrossingReport,
		"attributes": formattedJSON,
	})
	rr := map[string]interface{}{
		"data": data,
	}

	return rr, nil
}

// GetThresholdCrossingByMonitoredObject - Executes a GroupBy 'threshold crossing' query against druid. Wraps the
// result in a JSON API wrapper.
// peyo TODO: probably don't need to wrap JSON API here...should maybe do it elsewhere
func (dc *DruidDatastoreClient) GetThresholdCrossingByMonitoredObject(request *pb.ThresholdCrossingRequest, thresholdProfile *pb.TenantThresholdProfile) (map[string]interface{}, error) {

	logger.Log.Debugf("Calling GetThresholdCrossingByMonitoredObject for request: %v", models.AsJSONString(request))
	table := dc.cfg.GetString(gather.CK_druid_broker_table.String())

	// peyo TODO we should have a better way to handle default query params
	timeout := request.GetTimeout()
	if timeout == 0 {
		timeout = 5000
	}

	query, err := ThresholdCrossingByMonitoredObjectQuery(request.GetTenant(), table, request.Domain, request.Metric, request.Granularity, request.Interval, request.ObjectType, request.Direction, thresholdProfile.Data, request.GetVendor(), timeout)
	if err != nil {
		return nil, err
	}

	logger.Log.Debugf("Querying Druid for %s with query: %v", db.ThresholdCrossingByMonitoredObjectStr, models.AsJSONString(query))
	response, err := dc.executeQuery(query)
	if err != nil {
		return nil, err
	}

	thresholdCrossing := make([]ThresholdCrossingByMonitoredObjectResponse, 0)
	err = json.Unmarshal(response, &thresholdCrossing)
	if err != nil {
		return nil, err
	}

	logger.Log.Debugf("Response from druid for %s: %v", db.ThresholdCrossingByMonitoredObjectStr, models.AsJSONString(thresholdCrossing))

	formattedJSON, err := reformatThresholdCrossingByMonitoredObjectResponse(thresholdCrossing)
	if err != nil {
		return nil, err
	}

	// peyo TODO: need to figure out where to get this ID and Type from.
	uuid := uuid.NewV4()
	data := []map[string]interface{}{}
	data = append(data, map[string]interface{}{
		"id":         uuid.String(),
		"type":       ThresholdCrossingReport,
		"attributes": formattedJSON,
	})
	rr := map[string]interface{}{
		"data": data,
	}

	return rr, nil
}

// GetTopNFor
func (dc *DruidDatastoreClient) GetTopNForMetric(request *metrics.TopNForMetric) (map[string]interface{}, error) {

	logger.Log.Debugf("Calling GetTopNFor for request: %v", models.AsJSONString(request))

	query, err := GetTopNForMetric(dc.cfg.GetString(gather.CK_druid_broker_table.String()), request)
	if err != nil {
		return nil, err
	}

	logger.Log.Errorf("Querying Druid for %s with query: %+v", db.TopNForMetricString, models.AsJSONString(query))
	response, err := dc.executeQuery(query)
	if err != nil {
		return nil, err
	}

	construct := fmt.Sprintf("{\"results\":%s}", string(response))

	responseMap := make(map[string]interface{})
	if err = json.Unmarshal([]byte(construct), &responseMap); err != nil {
		logger.Log.Errorf("Could not Unmarshal, %s", err.Error())
		return nil, err
	}

	logger.Log.Debugf("Response from druid for query %s ->  %+v", db.TopNForMetricString, models.AsJSONString(responseMap))

	// peyo TODO: need to figure out where to get this ID and Type from.
	uuid := uuid.NewV4()
	data := []map[string]interface{}{}
	data = append(data, map[string]interface{}{
		"id":         uuid.String(),
		"type":       TopNForMetric,
		"attributes": responseMap,
	})
	rr := map[string]interface{}{
		"data": data,
	}
	logger.Log.Debugf("Response to caller %+v", rr)

	return rr, nil
}

// GetThresholdCrossingByMonitoredObjectTopN - Executes a TopN 'threshold crossing' query against druid. Wraps the
// result in a JSON API wrapper.
// peyo TODO: probably don't need to wrap JSON API here...should maybe do it elsewhere
func (dc *DruidDatastoreClient) GetThresholdCrossingByMonitoredObjectTopN(request *metrics.ThresholdCrossingTopNRequest, thresholdProfile *pb.TenantThresholdProfile) (map[string]interface{}, error) {

	logger.Log.Debugf("Calling GetThresholdCrossingByMonitoredObject for request: %v", models.AsJSONString(request))
	table := dc.cfg.GetString(gather.CK_druid_broker_table.String())

	query, err := ThresholdCrossingByMonitoredObjectTopNQuery(request.TenantID, table, request.Domain, request.Metric, request.Granularity, request.Interval, request.ObjectType, request.Direction, thresholdProfile.Data, request.Vendor, request.Timeout, request.NumResults)

	if err != nil {
		return nil, err
	}

	logger.Log.Debugf("Querying Druid for %s with query: %v", db.TopNThresholdCrossingByMonitoredObjectStr, models.AsJSONString(query))
	response, err := dc.executeQuery(query)
	if err != nil {
		return nil, err
	}

	thresholdCrossing := make([]TopNThresholdCrossingByMonitoredObjectResponse, 0)
	err = json.Unmarshal(response, &thresholdCrossing)
	if err != nil {
		return nil, err
	}

	logger.Log.Debugf("Response from druid for %s: %v", db.TopNThresholdCrossingByMonitoredObjectStr, models.AsJSONString(thresholdCrossing))

	// peyo TODO: need to figure out where to get this ID and Type from.
	uuid := uuid.NewV4()
	data := []map[string]interface{}{}
	data = append(data, map[string]interface{}{
		"id":         uuid.String(),
		"type":       ThresholdCrossingReport,
		"attributes": thresholdCrossing,
	})
	rr := map[string]interface{}{
		"data": data,
	}

	return rr, nil
}

func (dc *DruidDatastoreClient) GetAggregatedMetrics(request *metrics.AggregateMetricsAPIRequest) (map[string]interface{}, error) {
	logger.Log.Debugf("Calling GetAggregatedMetrics for request: %v", models.AsJSONString(request))
	table := dc.cfg.GetString(gather.CK_druid_broker_table.String())

	query, pp, err := AggMetricsQuery(request.TenantID, table, request.Interval, request.DomainIDs, request.Aggregation, request.Metrics, request.Timeout, request.Granularity)
	if err != nil {
		return nil, err
	}

	logger.Log.Debugf("Querying Druid for %s with query: %v", db.AggMetricsStr, models.AsJSONString(query))
	druidResponse, err := dc.executeQuery(query)
	if err != nil {
		return nil, err
	}

	response := make([]AggMetricsResponse, 0)
	err = json.Unmarshal(druidResponse, &response)
	if err != nil {
		return nil, err
	}

	logger.Log.Debugf("Response from druid for %s: %v", db.AggMetricsStr, models.AsJSONString(response))

	response = (*pp).Apply(response)

	rr := map[string]interface{}{
		"results": response,
	}

	logger.Log.Debugf("Processed response from druid for %s: %v", db.AggMetricsStr, models.AsJSONString(rr))

	return rr, nil
}

type Debug struct {
	Data map[string]interface{} `json:"data"`
}

func (dc *DruidDatastoreClient) GetSLAReport(request *metrics.SLAReportRequest, thresholdProfile *pb.TenantThresholdProfile) (*metrics.SLAReport, error) {
	logger.Log.Debugf("Calling GetSLAReport for request: %v", models.AsJSONString(request))
	table := dc.cfg.GetString(gather.CK_druid_broker_table.String())
	var query godruid.Query

	timeout := request.Timeout
	if timeout == 0 {
		timeout = 5000
	}

	query, err := SLAViolationsQuery(request.TenantID, table, request.Domain, Granularity_All, request.Interval, thresholdProfile.Data, timeout)

	if err != nil {
		return nil, err
	}

	logger.Log.Debugf("Querying Druid for %s with query: %v", db.SLAReportStr, models.AsJSONString(query))
	response, err := dc.executeQuery(query)
	if err != nil {
		return nil, err
	}

	reportSummary, err := reformatReportSummary(response)
	if err != nil {
		return nil, err
	}
	logger.Log.Debugf("Result: %v", db.SLAReportStr, models.AsJSONString(reportSummary))

	query, err = SLAViolationsQuery(request.TenantID, table, request.Domain, request.Granularity, request.Interval, thresholdProfile.Data, timeout)
	if err != nil {
		return nil, err
	}

	logger.Log.Debugf("Querying Druid for %s with query: %v", db.SLAReportStr, models.AsJSONString(query))
	response, err = dc.executeQuery(query)
	if err != nil {
		return nil, err
	}

	slaTimeSeries, err := reformatSLATimeSeries(response)
	if err != nil {
		return nil, err
	}

	var hourOfDayBucketMap map[string]interface{}
	var dayOfWeekBucketMap map[string]interface{}

	for vk, v := range thresholdProfile.Data.GetThresholds().GetVendorMap() {
		for tk, t := range v.GetMonitoredObjectTypeMap() {
			for mk, m := range t.GetMetricMap() {
				for dk, d := range m.GetDirectionMap() {
					for ek, e := range d.GetEventMap() {
						if ek != "sla" {
							continue
						}
						query, err = SLATimeBucketQuery(request.TenantID, table, request.Domain, DayOfWeek, request.Timezone, vk, tk, mk, dk, "sla", e, Granularity_All, request.Interval, timeout)
						if err != nil {
							return nil, err
						}

						logger.Log.Debugf("Querying Druid for %s with query: %v", db.SLAReportStr, models.AsJSONString(query))
						response, err = dc.executeQuery(query)
						if err != nil {
							return nil, err
						}

						dayOfWeekBucketMap, err = reformatSLABucketResponse(response, dayOfWeekBucketMap)
						if err != nil {
							return nil, err
						}

						query, err = SLATimeBucketQuery(request.TenantID, table, request.Domain, HourOfDay, request.Timezone, vk, tk, mk, dk, "sla", e, Granularity_All, request.Interval, timeout)
						if err != nil {
							return nil, err
						}

						logger.Log.Debugf("Querying Druid for %s with query: %v", db.SLAReportStr, models.AsJSONString(query))
						response, err = dc.executeQuery(query)
						if err != nil {
							return nil, err
						}

						hourOfDayBucketMap, err = reformatSLABucketResponse(response, hourOfDayBucketMap)
						if err != nil {
							return nil, err
						}
					}
				}
			}
		}
	}

	reportID := uuid.NewV4().String()

	slaReport := metrics.SLAReport{
		ID:                   reportID,
		ReportCompletionTime: time.Now().UTC().Format(time.RFC3339),
		TenantID:             request.TenantID,
		ReportTimeRange:      request.Interval,
		ReportSummary:        *reportSummary,
		TimeSeriesResult:     slaTimeSeries,
		ByHourOfDayResult:    hourOfDayBucketMap,
		ByDayOfWeekResult:    dayOfWeekBucketMap,
		ReportScheduleConfig: request.SlaScheduleConfig,
	}

	/*
		data := []map[string]interface{}{}
		data = append(data, map[string]interface{}{
			"id":         reportID,
			"type":       SLAReport,
			"attributes": slaReport,
		})

		rr := map[string]interface{}{
			"data": data,
		}
	*/
	return &slaReport, nil
}

func (dc *DruidDatastoreClient) GetRawMetrics(request *pb.RawMetricsRequest) (map[string]interface{}, error) {

	logger.Log.Debugf("Calling GetRawMetrics for request: %v", models.AsJSONString(request))

	table := dc.cfg.GetString(gather.CK_druid_broker_table.String())

	// peyo TODO we should have a better way to handle default query params
	timeout := request.GetTimeout()
	if timeout == 0 {
		timeout = 30000
	}

	granularity := request.GetGranularity()
	if granularity == "" {
		granularity = "PT1M"
	}

	query, err := RawMetricsQuery(request.GetTenant(), table, request.Metric, request.GetInterval(), request.GetObjectType(), request.GetDirection(), request.GetMonitoredObjectId(), timeout, granularity)

	if err != nil {
		return nil, err
	}

	logger.Log.Debugf("Querying Druid for %s with query: '' %s ''", db.RawMetricStr, models.AsJSONString(query))
	response, err := dc.executeQuery(query)

	if err != nil {
		return nil, err
	}

	resp := make([]RawMetricsResponse, 0)

	err = json.Unmarshal(response, &resp)
	if err != nil {
		return nil, err
	}

	logger.Log.Debugf("Response from druid for %s: %v", db.RawMetricStr, models.AsJSONString(resp))

	formattedJSON := map[string]interface{}{}
	if len(resp) != 0 {
		formattedJSON, err = reformatRawMetricsResponse(resp)
	}

	if err != nil {
		return nil, err
	}

	uuid := uuid.NewV4()
	data := []map[string]interface{}{}
	data = append(data, map[string]interface{}{
		"id":         uuid.String(),
		"type":       RawMetrics,
		"attributes": formattedJSON,
	})
	rr := map[string]interface{}{
		"data": data,
	}

	return rr, nil
}

type lookup struct {
	Version                string    `json:"version"`
	LookupExtractorFactory mapLookup `json:"lookupExtractorFactory"`
	active                 bool
}

type mapLookup struct {
	LookupType string            `json:"type"`
	Data       map[string]string `json:"map"`
}

func (dc *DruidDatastoreClient) UpdateMonitoredObjectMetadata(tenantID string, monitoredObjects []*tenant.MonitoredObject, domains []*tenant.Domain, reset bool) error {
	version := time.Now().Format(time.RFC3339)
	lookupEndpoint := dc.coordinatorServer + ":" + dc.coordinatorPort + "/druid/coordinator/v1/lookups/config"

	// Create 1 lookup per domain. Lookups don't support multiple values so the solution is to create
	// 1 lookup per domain and each lookup has a map where key is monitoredObjectId that belongs in that domain.
	// Every domain should have a map even if it has no monitored objects.
	lookups := make(map[string]*lookup)
	for _, domain := range domains {
		lookupName := buildLookupName("dom", tenantID, domain.ID)
		domLookup := lookup{
			Version: version,
			LookupExtractorFactory: mapLookup{
				LookupType: "map",
				Data:       map[string]string{},
			},
		}
		lookups[lookupName] = &domLookup
	}

	// Fetch existing lookup names and delete any existing lookups on the server that are nolonger valid.
	// Use the lookup map created in the previous step to identify valid domains.
	lookupNames := []string{}
	url := lookupEndpoint + "/__default"
	result, err := sendRequest("GET", dc.dClient.HttpClient, url, dc.AuthToken, nil)
	if err != nil {
		if strings.Contains(err.Error(), "No lookups found") {
			logger.Log.Infof("No lookups found.  Need to initialize lookups before any are created")
			result, err = sendRequest("POST", dc.dClient.HttpClient, lookupEndpoint, dc.AuthToken, []byte("{}"))
			if err != nil {
				logger.Log.Errorf("Failed to initialize druid lookups", err.Error())
				return err
			}
			logger.Log.Infof("Lookups successfully initialized")
		} else {
			logger.Log.Errorf("Failed to fetch lookups", err.Error())
			return err
		}
	} else {
		err = json.Unmarshal(result, &lookupNames)
	}

	// Only delete orphaned domain lookups for this tenant
	lookupPrefix := buildLookupNamePrefix("dom", tenantID)
	for _, lookupName := range lookupNames {
		if !strings.HasPrefix(lookupName, lookupPrefix) {
			continue
		}

		if lookup, ok := lookups[lookupName]; !ok {
			url := lookupEndpoint + "/__default/" + lookupName
			logger.Log.Debugf("Deleting lookup %s, url is %s", lookupName, url)
			if _, err := sendRequest("DELETE", dc.dClient.HttpClient, url, dc.AuthToken, nil); err != nil {
				logger.Log.Errorf("Failed to delete lookup %s", lookupName, err.Error())
			}
		} else {
			lookup.active = true
		}
	}

	// Now fill in the contents of each lookup by traversing the monitoredObject-to-domain associations.
	for _, mo := range monitoredObjects {
		if len(mo.DomainSet) < 1 {
			continue
		}
		for _, domain := range mo.DomainSet {
			lookupName := buildLookupName("dom", tenantID, domain)
			domLookup, ok := lookups[lookupName]
			if ok {
				domLookup.LookupExtractorFactory.Data[mo.MonitoredObjectID] = domain
			}
		}
	}

	// Domain lookups are assigned to the __default tier
	b, err := json.Marshal(map[string]map[string]*lookup{"__default": lookups})
	if err != nil {
		logger.Log.Error("Failed to marshal lookupRequest", err.Error())
		return err
	}

	//logger.Log.Debugf("Sending lookup request %s", string(b))
	_, err = sendRequest("POST", dc.dClient.HttpClient, lookupEndpoint, dc.AuthToken, b)
	if err != nil {
		logger.Log.Errorf("Failed to update lookup", err.Error())
		return err
	}
	updateLookupCache(lookups)
	return nil
}
