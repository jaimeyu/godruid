package druid

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"net/http"
	"sort"
	"strings"
	"time"

	"github.com/accedian/adh-gather/config"
	"github.com/accedian/adh-gather/gather"
	"github.com/accedian/adh-gather/logger"
	"github.com/accedian/adh-gather/models"
	"github.com/accedian/adh-gather/models/metrics"
	"github.com/accedian/adh-gather/models/tenant"
	"github.com/accedian/adh-gather/swagmodels"
	"github.com/accedian/godruid"

	db "github.com/accedian/adh-gather/datastore"
	pb "github.com/accedian/adh-gather/gathergrpc"
	mon "github.com/accedian/adh-gather/monitoring"

	"github.com/satori/go.uuid"
)

const (
	ThresholdCrossingReport = "threshold-crossing-report"
	EventDistribution       = "event-distribution"
	RawMetrics              = "raw-metrics"
	SLAReport               = "sla-report"
	TopNForMetric           = "top-n"

	errorCode   = "500"
	successCode = "200"
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

type BaseDruidResponse struct {
	Timestamp string                 `json:"timestamp"`
	Result    map[string]interface{} `json:"result"`
}

func makeHttpClient() *http.Client {
	// By default, use 60 second timeout unless specified otherwise
	// by the caller
	clientTimeout := 600 * time.Second

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
	dc.AuthToken = GetAuthCode(dc.cfg)

	err := client.Query(query, dc.AuthToken)

	if err != nil {
		if strings.Contains(err.Error(), "405") || strings.Contains(err.Error(), "401") {
			logger.Log.Info("Auth token expired, refreshing token. error:%s", err.Error())
			dc.AuthToken = GetAuthCode(dc.cfg)
			err_retry := client.Query(query, dc.AuthToken)
			if err_retry != nil {
				logger.Log.Errorf("Druid Query RETRY failed due to: %s", err)
				return nil, err_retry
			}
			return query.GetRawJSON(), nil
		}
		logger.Log.Errorf("Druid Query failed due to: %s for response: %s", err.Error(), models.AsJSONString(query.GetRawJSON()))
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

	var path string

	if port == "" {
		path = server
	} else {
		path = server + ":" + port
	}

	client := godruid.Client{
		Url:        path,
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

// Retrieves a histogram for specified metrics based on custom defined buckets
func (dc *DruidDatastoreClient) GetHistogram(request *metrics.Histogram, metaMOs []string) ([]metrics.TimeseriesEntryResponse, *db.QueryKeySpec, error) {
	methodStartTime := time.Now()
	if logger.IsDebugEnabled() {
		logger.Log.Debugf("Calling GetHistogram for request: %v", models.AsJSONString(request))
	}
	table := dc.cfg.GetString(gather.CK_druid_broker_table.String())

	timeout := request.Timeout
	if timeout == 0 {
		timeout = int32(gather.GetConfig().GetInt(gather.CK_druid_timeoutsms_histogram.String()))
	}

	// Build out the actual druid query to send
	query, querySpec, err := HistogramQuery(request.TenantID, metaMOs, table, request.Interval, request.Granularity, timeout, request.IgnoreCleaning, request.MetricBucketRequests)

	if err != nil {
		mon.TrackDruidTimeMetricInSeconds(mon.DruidAPIMethodDurationType, methodStartTime, errorCode, mon.GetHistogramObjStr)
		return nil, nil, err
	}

	// Execute the druid query
	queryStartTime := time.Now()
	if logger.IsDebugEnabled() {
		logger.Log.Debugf("Query Key Spec map contains the following entries: %v", querySpec)
		logger.Log.Debugf("Querying Druid for %s with query: %v", db.HistogramStr, models.AsJSONString(query))
	}

	druidResponse, err := dc.executeQuery(query)

	if err != nil {
		mon.TrackDruidTimeMetricInSeconds(mon.DruidQueryDurationType, queryStartTime, errorCode, mon.GetHistogramObjStr)
		mon.TrackDruidTimeMetricInSeconds(mon.DruidAPIMethodDurationType, methodStartTime, errorCode, mon.GetHistogramObjStr)
		return nil, nil, err
	}
	mon.TrackDruidTimeMetricInSeconds(mon.DruidQueryDurationType, queryStartTime, successCode, mon.GetHistogramObjStr)

	// Reformat the druid response from a flat structure to a json api structure
	if logger.IsDebugEnabled() {
		logger.Log.Debugf("Response from druid for %s: %v", db.HistogramStr, string(druidResponse))
	}

	response := make([]metrics.TimeseriesEntryResponse, 0)
	err = json.Unmarshal(druidResponse, &response)
	if err != nil {
		mon.TrackDruidTimeMetricInSeconds(mon.DruidAPIMethodDurationType, methodStartTime, errorCode, mon.QueryAggregatedMetricsStr)
		return nil, nil, err
	}

	mon.TrackDruidTimeMetricInSeconds(mon.DruidAPIMethodDurationType, methodStartTime, successCode, mon.GetHistogramObjStr)
	return response, querySpec, nil
}

// DEPRECATED - Remove once v1 queries have been moved
// Retrieves a histogram for specified metrics based on custom defined buckets
func (dc *DruidDatastoreClient) GetHistogramV1(request *metrics.HistogramV1, metaMOs []string) (map[string]interface{}, error) {
	methodStartTime := time.Now()
	if logger.IsDebugEnabled() {
		logger.Log.Debugf("Calling GetHistogram for request: %v", models.AsJSONString(request))
	}
	table := dc.cfg.GetString(gather.CK_druid_broker_table.String())

	timeout := request.Timeout
	if timeout == 0 {
		timeout = int32(gather.GetConfig().GetInt(gather.CK_druid_timeoutsms_histogram.String()))
	}

	// Split out the request into a set of request metrics keyed off of the metric vendor, objectType, name, and direction
	metrics := make([]map[string]interface{}, len(request.MetricBucketRequests))
	for i, mb := range request.MetricBucketRequests {
		metricsMap, err := models.ConvertObj2Map(mb)
		if err != nil {
			return nil, err
		}
		v, _ := metricsMap["name"]
		metricsMap["metric"] = v
		delete(metricsMap, "name")

		metrics[i] = metricsMap
	}

	// Build out the actual druid query to send
	query, err := HistogramQueryV1(request.TenantID, metaMOs, table, request.Interval, request.Granularity, timeout, metrics)

	if err != nil {
		mon.TrackDruidTimeMetricInSeconds(mon.DruidAPIMethodDurationType, methodStartTime, errorCode, mon.GetHistogramObjStr)
		return nil, err
	}

	// Execute the druid query
	queryStartTime := time.Now()
	if logger.IsDebugEnabled() {
		logger.Log.Debugf("Querying Druid for %s with query: %v", db.HistogramStr, models.AsJSONString(query))
	}
	response, err := dc.executeQuery(query)

	if err != nil {
		mon.TrackDruidTimeMetricInSeconds(mon.DruidQueryDurationType, queryStartTime, errorCode, mon.GetHistogramObjStr)
		mon.TrackDruidTimeMetricInSeconds(mon.DruidAPIMethodDurationType, methodStartTime, errorCode, mon.GetHistogramObjStr)
		return nil, err
	}
	mon.TrackDruidTimeMetricInSeconds(mon.DruidQueryDurationType, queryStartTime, successCode, mon.GetHistogramObjStr)

	// Reformat the druid response from a flat structure to a json api structure
	if logger.IsDebugEnabled() {
		logger.Log.Debugf("Response from druid for %s: %v", db.HistogramStr, string(response))
	}
	rr, err := reformatHistogramResponseV1(request.TenantID, request.Meta, request.Interval, string(response))

	if err != nil {
		mon.TrackDruidTimeMetricInSeconds(mon.DruidAPIMethodDurationType, methodStartTime, errorCode, mon.GetHistogramObjStr)
		return nil, err
	}

	mon.TrackDruidTimeMetricInSeconds(mon.DruidAPIMethodDurationType, methodStartTime, successCode, mon.GetHistogramObjStr)
	return rr, nil
}

// New version of threshold-crossing
func (dc *DruidDatastoreClient) QueryThresholdCrossing(request *metrics.ThresholdCrossing, thresholdProfile *tenant.ThresholdProfile, metaMOs []string) (map[string]interface{}, error) {
	methodStartTime := time.Now()
	if logger.IsDebugEnabled() {
		logger.Log.Debugf("Calling QueryThresholdCrossing for request: %v", models.AsJSONString(request))
	}
	table := dc.cfg.GetString(gather.CK_druid_broker_table.String())

	timeout := request.Timeout
	if timeout == 0 {
		timeout = int32(gather.GetConfig().GetInt(gather.CK_druid_timeoutsms_thresholdcrossing.String()))
	}

	query, err := ThresholdViolationsQuery(request.TenantID, table, metaMOs, request.Granularity, request.Interval, request.IgnoreCleaning, request.Metrics, thresholdProfile, timeout)

	if err != nil {
		mon.TrackDruidTimeMetricInSeconds(mon.DruidAPIMethodDurationType, methodStartTime, errorCode, mon.GetThrCrossStr)
		return nil, err
	}
	queryStartTime := time.Now()
	if logger.IsDebugEnabled() {
		logger.Log.Debugf("Querying Druid for %s with query: %v", db.ThresholdCrossingStr, models.AsJSONString(query))
	}
	druidResponse, err := dc.executeQuery(query)

	response := make([]BaseDruidResponse, 0)
	err = json.Unmarshal(druidResponse, &response)
	if err != nil {
		mon.TrackDruidTimeMetricInSeconds(mon.DruidQueryDurationType, queryStartTime, errorCode, mon.GetThrCrossStr)
		mon.TrackDruidTimeMetricInSeconds(mon.DruidAPIMethodDurationType, methodStartTime, errorCode, mon.GetThrCrossStr)
		return nil, err
	}
	mon.TrackDruidTimeMetricInSeconds(mon.DruidQueryDurationType, queryStartTime, successCode, mon.GetThrCrossStr)

	if logger.IsDebugEnabled() {
		logger.Log.Debugf("Response from druid for %s: %v", db.ThresholdCrossingStr, models.AsJSONString(response))
	}

	reformatted, err := reformatThresholdCrossingTimeSeries(druidResponse)
	if err != nil {
		mon.TrackDruidTimeMetricInSeconds(mon.DruidAPIMethodDurationType, methodStartTime, errorCode, mon.GetThrCrossStr)
		return nil, err
	}

	rr := map[string]interface{}{
		"results": reformatted,
	}

	if logger.IsDebugEnabled() {
		logger.Log.Debugf("Processed response from druid for %s: %v", db.ThresholdCrossingStr, models.AsJSONString(rr))
	}
	mon.TrackDruidTimeMetricInSeconds(mon.DruidAPIMethodDurationType, methodStartTime, successCode, mon.GetThrCrossStr)

	return rr, nil
}

// New version of threshold-crossing
func (dc *DruidDatastoreClient) QueryThresholdCrossingV1(request *metrics.ThresholdCrossingV1, thresholdProfile *pb.TenantThresholdProfile, metaMOs []string) (map[string]interface{}, error) {
	methodStartTime := time.Now()
	if logger.IsDebugEnabled() {
		logger.Log.Debugf("Calling QueryThresholdCrossing for request: %v", models.AsJSONString(request))
	}
	table := dc.cfg.GetString(gather.CK_druid_broker_table.String())

	timeout := request.Timeout
	if timeout == 0 {
		timeout = int32(gather.GetConfig().GetInt(gather.CK_druid_timeoutsms_thresholdcrossing.String()))
	}

	query, err := ThresholdViolationsQueryV1(request.TenantID, table, metaMOs, request.Granularity, request.Interval, request.Metrics, thresholdProfile.Data, timeout)

	if err != nil {
		mon.TrackDruidTimeMetricInSeconds(mon.DruidAPIMethodDurationType, methodStartTime, errorCode, mon.GetThrCrossStr)
		return nil, err
	}
	queryStartTime := time.Now()
	if logger.IsDebugEnabled() {
		logger.Log.Debugf("Querying Druid for %s with query: %v", db.ThresholdCrossingStr, models.AsJSONString(query))
	}
	druidResponse, err := dc.executeQuery(query)

	response := make([]BaseDruidResponse, 0)
	err = json.Unmarshal(druidResponse, &response)
	if err != nil {
		mon.TrackDruidTimeMetricInSeconds(mon.DruidQueryDurationType, queryStartTime, errorCode, mon.GetThrCrossStr)
		mon.TrackDruidTimeMetricInSeconds(mon.DruidAPIMethodDurationType, methodStartTime, errorCode, mon.GetThrCrossStr)
		return nil, err
	}
	mon.TrackDruidTimeMetricInSeconds(mon.DruidQueryDurationType, queryStartTime, successCode, mon.GetThrCrossStr)

	if logger.IsDebugEnabled() {
		logger.Log.Debugf("Response from druid for %s: %v", db.ThresholdCrossingStr, models.AsJSONString(response))
	}

	reformatted, err := reformatThresholdCrossingTimeSeries(druidResponse)
	if err != nil {
		mon.TrackDruidTimeMetricInSeconds(mon.DruidAPIMethodDurationType, methodStartTime, errorCode, mon.GetThrCrossStr)
		return nil, err
	}

	rr := map[string]interface{}{
		"results": reformatted,
	}

	if logger.IsDebugEnabled() {
		logger.Log.Debugf("Processed response from druid for %s: %v", db.ThresholdCrossingStr, models.AsJSONString(rr))
	}
	mon.TrackDruidTimeMetricInSeconds(mon.DruidAPIMethodDurationType, methodStartTime, successCode, mon.GetThrCrossStr)

	return rr, nil
}

// GetTopNFor - Executes a TopN on a given metric, based on its min/max/avg.
func (dc *DruidDatastoreClient) GetTopNForMetric(request *metrics.TopNForMetric, metaMOs []string) ([]metrics.TopNEntryResponse, error) {

	methodStartTime := time.Now()

	if logger.IsDebugEnabled() {
		logger.Log.Debugf("Calling GetTopNFor for request: %v", models.AsJSONString(request))
	}

	timeout := request.Timeout
	if timeout == 0 {
		timeout = int32(gather.GetConfig().GetInt(gather.CK_druid_timeoutsms_topn.String()))
	}

	query, err := GetTopNForMetric(dc.cfg.GetString(gather.CK_druid_broker_table.String()), request, timeout, request.IgnoreCleaning, metaMOs)
	if err != nil {
		mon.TrackDruidTimeMetricInSeconds(mon.DruidAPIMethodDurationType, methodStartTime, errorCode, mon.GetTopNReqStr)
		return nil, fmt.Errorf("Failed to generate a druid query while processing request: %s: '%s'", models.AsJSONString(request), err.Error())
	}

	if logger.IsDebugEnabled() {
		logger.Log.Debugf("Querying Druid for %s with query: %s", db.TopNForMetricStr, models.AsJSONString(query))
	}

	queryStartTime := time.Now()
	response, err := dc.executeQuery(query)
	if err != nil {
		mon.TrackDruidTimeMetricInSeconds(mon.DruidQueryDurationType, queryStartTime, errorCode, mon.GetTopNReqStr)
		mon.TrackDruidTimeMetricInSeconds(mon.DruidAPIMethodDurationType, methodStartTime, errorCode, mon.GetTopNReqStr)
		return nil, fmt.Errorf("Failed to get TopN result from druid for request %s: %s", models.AsJSONString(query), err.Error())
	}
	mon.TrackDruidTimeMetricInSeconds(mon.DruidQueryDurationType, queryStartTime, successCode, mon.GetTopNReqStr)

	topN := make([]map[string]interface{}, 0)
	err = json.Unmarshal(response, &topN)
	if err != nil {
		mon.TrackDruidTimeMetricInSeconds(mon.DruidAPIMethodDurationType, methodStartTime, errorCode, mon.GetTopNReqStr)
		return nil, err
	}

	if logger.IsDebugEnabled() {
		logger.Log.Debugf("Response from druid for query %s ->  %+v", db.TopNForMetricStr, models.AsJSONString(topN))
	}

	topNResults := make([]metrics.TopNEntryResponse, 0)

	if len(topN) != 0 {
		topNResponseHead, ok := topN[0]["result"] // There should be only a single entry with the granularity all
		if !ok {
			return topNResults, nil
		}

		topNResponse, ok := topNResponseHead.([]interface{}) // There should be only a single entry with the granularity all
		if !ok {
			logger.Log.Errorf("Could not cast topN response")
			mon.TrackDruidTimeMetricInSeconds(mon.DruidAPIMethodDurationType, methodStartTime, errorCode, mon.GetTopNReqStr)
			return nil, err
		}
		topNResults = make([]metrics.TopNEntryResponse, 0)

		for _, r := range topNResponse {
			rawMap := r.(map[string]interface{})
			if rawMap["monitoredObjectId"] != nil {
				toAdd := metrics.TopNEntryResponse{MonitoredObjectId: rawMap["monitoredObjectId"].(string)}
				delete(rawMap, "monitoredObjectId")
				toAdd.Result = rawMap
				topNResults = append(topNResults, toAdd)
			}
		}

	}

	mon.TrackDruidTimeMetricInSeconds(mon.DruidAPIMethodDurationType, methodStartTime, successCode, mon.GetTopNReqStr)
	return topNResults, nil
}

// DEPRECATED - Remove once v1 queries have been moved
// GetTopNFor - Executes a TopN on a given metric, based on its min/max/avg.
func (dc *DruidDatastoreClient) GetTopNForMetricV1(request *metrics.TopNForMetricV1, metaMOs []string) (map[string]interface{}, error) {
	stat := "druid_topn_get"
	methodStartTime := time.Now()

	if logger.IsDebugEnabled() {
		logger.Log.Debugf("Calling GetTopNFor for request: %v", models.AsJSONString(request))
	}

	query, err := GetTopNForMetricV1(dc.cfg.GetString(gather.CK_druid_broker_table.String()), request, metaMOs)
	if err != nil {
		mon.TrackDruidTimeMetricInSeconds(mon.DruidAPIMethodDurationType, methodStartTime, errorCode, stat)
		return nil, fmt.Errorf("Failed to generate a druid query while processing request: %s: '%s'", models.AsJSONString(request), err.Error())
	}

	if logger.IsDebugEnabled() {
		logger.Log.Debugf("Querying Druid for %s with query: %s", db.TopNForMetricStr, models.AsJSONString(query))
	}

	queryStartTime := time.Now()
	response, err := dc.executeQuery(query)
	if err != nil {
		mon.TrackDruidTimeMetricInSeconds(mon.DruidQueryDurationType, queryStartTime, errorCode, stat)
		mon.TrackDruidTimeMetricInSeconds(mon.DruidAPIMethodDurationType, methodStartTime, errorCode, stat)
		return nil, fmt.Errorf("Failed to get TopN result from druid for request %s: %s", models.AsJSONString(query), err.Error())
	}
	mon.TrackDruidTimeMetricInSeconds(mon.DruidQueryDurationType, queryStartTime, successCode, "QUERY_"+stat)

	construct := fmt.Sprintf("{\"results\":%s}", string(response))

	responseMap := make(map[string]interface{})
	if err = json.Unmarshal([]byte(construct), &responseMap); err != nil {
		mon.TrackDruidTimeMetricInSeconds(mon.DruidAPIMethodDurationType, methodStartTime, errorCode, stat)
		return nil, fmt.Errorf("Unable to unmarshal response from druid for request %s: %s", models.AsJSONString(request), err.Error())
	}

	if logger.IsDebugEnabled() {
		logger.Log.Debugf("Response from druid for query %s ->  %+v", db.TopNForMetricStr, models.AsJSONString(responseMap))
	}

	mon.TrackDruidTimeMetricInSeconds(mon.DruidAPIMethodDurationType, methodStartTime, successCode, "METHOD_"+stat)
	return responseMap, nil
}

func (dc *DruidDatastoreClient) GetThresholdCrossingByMonitoredObjectTopN(request *metrics.ThresholdCrossingTopN, thresholdProfile *tenant.ThresholdProfile, metaMOs []string) ([]metrics.TopNEntryResponse, error) {
	methodStartTime := time.Now()
	if logger.IsDebugEnabled() {
		logger.Log.Debugf("Calling GetThresholdCrossingByMonitoredObject for request: %v", models.AsJSONString(request))
	}
	table := dc.cfg.GetString(gather.CK_druid_broker_table.String())

	timeout := request.Timeout
	if timeout == 0 {
		timeout = int32(gather.GetConfig().GetInt(gather.CK_druid_timeoutsms_thresholdcrossingtopn.String()))
	}

	query, err := ThresholdCrossingByMonitoredObjectTopNQuery(request.TenantID, table, metaMOs, request.Metric, request.Granularity, request.Interval, request.IgnoreCleaning, thresholdProfile, timeout, request.NumResults)

	if err != nil {
		mon.TrackDruidTimeMetricInSeconds(mon.DruidAPIMethodDurationType, methodStartTime, errorCode, mon.GetThrCrossByMonObjTopNStr)
		return nil, err
	}

	queryStartTime := time.Now()
	if logger.IsDebugEnabled() {
		logger.Log.Debugf("Querying Druid for %s with query: %v", db.TopNThresholdCrossingByMonitoredObjectStr, models.AsJSONString(query))
	}
	response, err := dc.executeQuery(query)
	if err != nil {
		mon.TrackDruidTimeMetricInSeconds(mon.DruidQueryDurationType, queryStartTime, errorCode, mon.GetThrCrossByMonObjTopNStr)
		mon.TrackDruidTimeMetricInSeconds(mon.DruidAPIMethodDurationType, methodStartTime, errorCode, mon.GetThrCrossByMonObjTopNStr)
		return nil, err
	}
	mon.TrackDruidTimeMetricInSeconds(mon.DruidQueryDurationType, queryStartTime, errorCode, mon.GetThrCrossByMonObjTopNStr)

	topN := make([]map[string]interface{}, 0)
	err = json.Unmarshal(response, &topN)
	if err != nil {
		mon.TrackDruidTimeMetricInSeconds(mon.DruidAPIMethodDurationType, methodStartTime, errorCode, mon.GetThrCrossByMonObjTopNStr)
		return nil, err
	}

	if logger.IsDebugEnabled() {
		logger.Log.Debugf("Response from druid for %s: %v", db.TopNThresholdCrossingByMonitoredObjectStr, models.AsJSONString(topN))
	}

	topNResults := make([]metrics.TopNEntryResponse, 0)

	if len(topN) != 0 {
		topNResponseHead, ok := topN[0]["result"] // There should be only a single entry with the granularity all
		if !ok {
			return topNResults, nil
		}

		topNResponse, ok := topNResponseHead.([]interface{}) // There should be only a single entry with the granularity all
		if !ok {
			logger.Log.Errorf("Could not cast topN response")
			mon.TrackDruidTimeMetricInSeconds(mon.DruidAPIMethodDurationType, methodStartTime, errorCode, mon.GetTopNReqStr)
			return nil, err
		}
		topNResults = make([]metrics.TopNEntryResponse, 0)

		for _, r := range topNResponse {
			rawMap := r.(map[string]interface{})
			if score := rawMap["scored"]; rawMap["monitoredObjectId"] != nil && score != nil && score.(float64) != 0 {
				toAdd := metrics.TopNEntryResponse{MonitoredObjectId: rawMap["monitoredObjectId"].(string)}
				delete(rawMap, "monitoredObjectId")
				toAdd.Result = rawMap
				topNResults = append(topNResults, toAdd)
			}
		}

	}

	mon.TrackDruidTimeMetricInSeconds(mon.DruidAPIMethodDurationType, methodStartTime, successCode, mon.GetThrCrossByMonObjTopNStr)
	return topNResults, nil
}

// DEPRECATED - Remove once v1 queries have been moved
// GetThresholdCrossingByMonitoredObjectTopN - Executes a TopN 'threshold crossing' query against druid. Wraps the
// result in a JSON API wrapper.
// peyo TODO: probably don't need to wrap JSON API here...should maybe do it elsewhere
func (dc *DruidDatastoreClient) GetThresholdCrossingByMonitoredObjectTopNV1(request *metrics.ThresholdCrossingTopNV1, thresholdProfile *pb.TenantThresholdProfile, metaMOs []string) (map[string]interface{}, error) {
	methodStartTime := time.Now()
	if logger.IsDebugEnabled() {
		logger.Log.Debugf("Calling GetThresholdCrossingByMonitoredObject for request: %v", models.AsJSONString(request))
	}
	table := dc.cfg.GetString(gather.CK_druid_broker_table.String())

	direction := fmt.Sprintf("%v", request.Metric.Direction)

	query, err := ThresholdCrossingByMonitoredObjectTopNQueryV1(request.TenantID, table, metaMOs, request.Metric.Name, request.Granularity, request.Interval, request.Metric.ObjectType, direction, thresholdProfile.Data, request.Metric.Vendor, request.Timeout, request.NumResults)

	if err != nil {
		mon.TrackDruidTimeMetricInSeconds(mon.DruidAPIMethodDurationType, methodStartTime, errorCode, mon.GetThrCrossByMonObjTopNStr)
		return nil, err
	}

	queryStartTime := time.Now()
	if logger.IsDebugEnabled() {
		logger.Log.Debugf("Querying Druid for %s with query: %v", db.TopNThresholdCrossingByMonitoredObjectStr, models.AsJSONString(query))
	}
	response, err := dc.executeQuery(query)
	if err != nil {
		mon.TrackDruidTimeMetricInSeconds(mon.DruidQueryDurationType, queryStartTime, errorCode, mon.GetThrCrossByMonObjTopNStr)
		mon.TrackDruidTimeMetricInSeconds(mon.DruidAPIMethodDurationType, methodStartTime, errorCode, mon.GetThrCrossByMonObjTopNStr)
		return nil, err
	}
	mon.TrackDruidTimeMetricInSeconds(mon.DruidQueryDurationType, queryStartTime, errorCode, mon.GetThrCrossByMonObjTopNStr)

	thresholdCrossing := make([]TopNThresholdCrossingByMonitoredObjectResponse, 0)
	err = json.Unmarshal(response, &thresholdCrossing)
	if err != nil {
		mon.TrackDruidTimeMetricInSeconds(mon.DruidAPIMethodDurationType, methodStartTime, errorCode, mon.GetThrCrossByMonObjTopNStr)
		return nil, err
	}

	if logger.IsDebugEnabled() {
		logger.Log.Debugf("Response from druid for %s: %v", db.TopNThresholdCrossingByMonitoredObjectStr, models.AsJSONString(thresholdCrossing))
	}

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

	mon.TrackDruidTimeMetricInSeconds(mon.DruidAPIMethodDurationType, methodStartTime, successCode, mon.GetThrCrossByMonObjTopNStr)
	return rr, nil
}

func (dc *DruidDatastoreClient) GetAggregatedMetrics(request *metrics.AggregateMetrics, metaMOs []string) ([]metrics.TimeseriesEntryResponse, *db.QueryKeySpec, error) {
	methodStartTime := time.Now()
	if logger.IsDebugEnabled() {
		logger.Log.Debugf("Calling GetAggregatedMetrics for request: %v", models.AsJSONString(request))
	}
	table := dc.cfg.GetString(gather.CK_druid_broker_table.String())

	timeout := request.Timeout
	if timeout == 0 {
		timeout = int32(gather.GetConfig().GetInt(gather.CK_druid_timeoutsms_aggregatedmetrics.String()))
	}

	aggregateOnMeta := metaMOs != nil
	var monitoredObjectIds []string

	// If we're provided metadata then we want to aggregate on the set of monitored objects that pass that metadata filter
	// Otherwise we use the explicit list of monitored objects that are passed in and each individual monitored object is processed individually
	if aggregateOnMeta {
		monitoredObjectIds = metaMOs
	} else {
		monitoredObjectIds = request.MonitoredObjects
	}

	query, pp, queryKeySpec, err := AggMetricsQuery(request.TenantID, table, request.Interval, monitoredObjectIds, aggregateOnMeta, request.Aggregation, request.Metrics, request.IgnoreCleaning, timeout, request.Granularity)
	if err != nil {
		mon.TrackDruidTimeMetricInSeconds(mon.DruidAPIMethodDurationType, methodStartTime, errorCode, mon.QueryAggregatedMetricsStr)
		return nil, nil, err
	}

	queryStartTime := time.Now()
	if logger.IsDebugEnabled() {
		logger.Log.Debugf("Query Key Spec map contains the following entries: %v", queryKeySpec)
		logger.Log.Debugf("Querying Druid for %s with query: %v", db.AggMetricsStr, models.AsJSONString(query))
	}
	druidResponse, err := dc.executeQuery(query)
	if err != nil {
		mon.TrackDruidTimeMetricInSeconds(mon.DruidQueryDurationType, queryStartTime, errorCode, mon.QueryAggregatedMetricsStr)
		mon.TrackDruidTimeMetricInSeconds(mon.DruidAPIMethodDurationType, methodStartTime, errorCode, mon.QueryAggregatedMetricsStr)
		return nil, nil, err
	}
	mon.TrackDruidTimeMetricInSeconds(mon.DruidQueryDurationType, queryStartTime, successCode, mon.QueryAggregatedMetricsStr)

	response := make([]metrics.TimeseriesEntryResponse, 0)
	err = json.Unmarshal(druidResponse, &response)
	if err != nil {
		mon.TrackDruidTimeMetricInSeconds(mon.DruidAPIMethodDurationType, methodStartTime, errorCode, mon.QueryAggregatedMetricsStr)
		return nil, nil, err
	}

	if logger.IsDebugEnabled() {
		logger.Log.Debugf("Response from druid for %s: %v", db.AggMetricsStr, models.AsJSONString(response))
	}

	response = (*pp).Apply(response)

	if logger.IsDebugEnabled() {
		logger.Log.Debugf("Processed response from druid for %s: %v", db.AggMetricsStr, response)
	}

	mon.TrackDruidTimeMetricInSeconds(mon.DruidAPIMethodDurationType, methodStartTime, successCode, mon.QueryAggregatedMetricsStr)
	return response, queryKeySpec, nil
}

// DEPRECATED - Remove once v1 queries have been moved
func (dc *DruidDatastoreClient) GetAggregatedMetricsV1(request *metrics.AggregateMetricsV1, metaMOs []string) (map[string]interface{}, error) {
	methodStartTime := time.Now()
	if logger.IsDebugEnabled() {
		logger.Log.Debugf("Calling GetAggregatedMetrics for request: %v", models.AsJSONString(request))
	}
	table := dc.cfg.GetString(gather.CK_druid_broker_table.String())

	timeout := request.Timeout
	if timeout == 0 {
		timeout = int32(gather.GetConfig().GetInt(gather.CK_druid_timeoutsms_aggregatedmetrics.String()))
	}

	query, pp, err := AggMetricsQueryV1(request.TenantID, table, request.Interval, metaMOs, request.Aggregation, request.Metrics, timeout, request.Granularity)
	if err != nil {
		mon.TrackDruidTimeMetricInSeconds(mon.DruidAPIMethodDurationType, methodStartTime, errorCode, mon.QueryAggregatedMetricsStr)
		return nil, err
	}

	queryStartTime := time.Now()
	if logger.IsDebugEnabled() {
		logger.Log.Debugf("Querying Druid for %s with query: %v", db.AggMetricsStr, models.AsJSONString(query))
	}
	druidResponse, err := dc.executeQuery(query)
	if err != nil {
		mon.TrackDruidTimeMetricInSeconds(mon.DruidQueryDurationType, queryStartTime, errorCode, mon.QueryAggregatedMetricsStr)
		mon.TrackDruidTimeMetricInSeconds(mon.DruidAPIMethodDurationType, methodStartTime, errorCode, mon.QueryAggregatedMetricsStr)
		return nil, err
	}
	mon.TrackDruidTimeMetricInSeconds(mon.DruidQueryDurationType, queryStartTime, successCode, mon.QueryAggregatedMetricsStr)

	response := make([]metrics.TimeseriesEntryResponseV1, 0)
	err = json.Unmarshal(druidResponse, &response)
	if err != nil {
		mon.TrackDruidTimeMetricInSeconds(mon.DruidAPIMethodDurationType, methodStartTime, errorCode, mon.QueryAggregatedMetricsStr)
		return nil, err
	}

	if logger.IsDebugEnabled() {
		logger.Log.Debugf("Response from druid for %s: %v", db.AggMetricsStr, models.AsJSONString(response))
	}

	response = (*pp).Apply(response)

	rr := map[string]interface{}{
		"results": response,
	}

	if logger.IsDebugEnabled() {
		logger.Log.Debugf("Processed response from druid for %s: %v", db.AggMetricsStr, models.AsJSONString(rr))
	}

	mon.TrackDruidTimeMetricInSeconds(mon.DruidAPIMethodDurationType, methodStartTime, successCode, mon.QueryAggregatedMetricsStr)
	return rr, nil
}

type Debug struct {
	Data map[string]interface{} `json:"data"`
}

func (dc *DruidDatastoreClient) GetSLAReportV1(request *metrics.SLAReportRequest, thresholdProfile *tenant.ThresholdProfile, metaMOs []string) (*metrics.SLAReport, error) {
	methodStartTime := time.Now()
	if logger.IsDebugEnabled() {
		logger.Log.Debugf("Calling GetSLAReport for request: %v", models.AsJSONString(request))
	}
	table := dc.cfg.GetString(gather.CK_druid_broker_table.String())
	var query godruid.Query

	timeout := request.Timeout
	if timeout == 0 {
		timeout = int32(gather.GetConfig().GetInt(gather.CK_druid_timeoutsms_slareports.String()))
	}

	query, _, err := SLAViolationsQuery(request.TenantID, table, metaMOs, GranularityAll, request.Interval, false, thresholdProfile, timeout)

	if err != nil {
		mon.TrackDruidTimeMetricInSeconds(mon.DruidAPIMethodDurationType, methodStartTime, errorCode, mon.GetSLAReportStr)
		return nil, err
	}

	queryStartTime := time.Now()
	if logger.IsDebugEnabled() {
		logger.Log.Debugf("Querying Druid for %s with query: %v", db.SLAReportStr, models.AsJSONString(query))
	}
	response, err := dc.executeQuery(query)
	if err != nil {
		mon.TrackDruidTimeMetricInSeconds(mon.DruidQueryDurationType, queryStartTime, errorCode, mon.SLAViolationsQueryStr)
		mon.TrackDruidTimeMetricInSeconds(mon.DruidAPIMethodDurationType, methodStartTime, errorCode, mon.GetSLAReportStr)
		return nil, err
	}
	mon.TrackDruidTimeMetricInSeconds(mon.DruidQueryDurationType, queryStartTime, successCode, mon.SLAViolationsQueryStr)

	reportSummary, err := reformatReportSummary(response)
	if err != nil {
		mon.TrackDruidTimeMetricInSeconds(mon.DruidQueryDurationType, queryStartTime, successCode, mon.GetSLAReportStr)
		return nil, err
	}
	if logger.IsDebugEnabled() {
		logger.Log.Debugf("Result: %v", db.SLAReportStr, models.AsJSONString(reportSummary))
	}

	query, _, err = SLAViolationsQuery(request.TenantID, table, metaMOs, request.Granularity, request.Interval, false, thresholdProfile, timeout)
	if err != nil {
		mon.TrackDruidTimeMetricInSeconds(mon.DruidQueryDurationType, queryStartTime, successCode, mon.GetSLAReportStr)
		return nil, err
	}

	queryStartTime = time.Now()
	if logger.IsDebugEnabled() {
		logger.Log.Debugf("Querying Druid for %s with query: %v", db.SLAReportStr, models.AsJSONString(query))
	}
	response, err = dc.executeQuery(query)
	if err != nil {
		mon.TrackDruidTimeMetricInSeconds(mon.DruidQueryDurationType, queryStartTime, errorCode, mon.SLAViolationsQueryStr)
		mon.TrackDruidTimeMetricInSeconds(mon.DruidAPIMethodDurationType, methodStartTime, errorCode, mon.GetSLAReportStr)
		return nil, err
	}
	mon.TrackDruidTimeMetricInSeconds(mon.DruidQueryDurationType, queryStartTime, successCode, mon.SLAViolationsQueryStr)

	slaTimeSeries, err := reformatSLATimeSeries(response)
	if err != nil {
		mon.TrackDruidTimeMetricInSeconds(mon.DruidAPIMethodDurationType, methodStartTime, errorCode, mon.GetSLAReportStr)
		return nil, err
	}

	var hourOfDayBucketMap map[string]interface{}
	var dayOfWeekBucketMap map[string]interface{}

	for vk, v := range thresholdProfile.Thresholds.VendorMap {
		for tk, t := range v.MonitoredObjectTypeMap {
			if tk != "twamp-sf" {
				continue
			}

			for mk, m := range t.MetricMap {
				for dk, d := range m.DirectionMap {
					for ek, e := range d.EventMap {
						if ek != "sla" {
							continue
						}
						query, _, err = SLATimeBucketQuery(request.TenantID, table, metaMOs, DayOfWeek, false, request.Timezone, vk, tk, mk, dk, "sla", e, GranularityAll, request.Interval, timeout)
						if err != nil {
							mon.TrackDruidTimeMetricInSeconds(mon.DruidAPIMethodDurationType, methodStartTime, errorCode, mon.GetSLAReportStr)
							return nil, err
						}

						queryStartTime = time.Now()
						if logger.IsDebugEnabled() {
							logger.Log.Debugf("Querying Druid for %s with query: %v", db.SLAReportStr, models.AsJSONString(query))
						}
						response, err = dc.executeQuery(query)
						if err != nil {
							mon.TrackDruidTimeMetricInSeconds(mon.DruidQueryDurationType, queryStartTime, errorCode, mon.SLATimeBucketQueryStr)
							mon.TrackDruidTimeMetricInSeconds(mon.DruidAPIMethodDurationType, methodStartTime, errorCode, mon.GetSLAReportStr)
							return nil, err
						}
						mon.TrackDruidTimeMetricInSeconds(mon.DruidQueryDurationType, queryStartTime, successCode, mon.SLATimeBucketQueryStr)

						dayOfWeekBucketMap, err = reformatSLABucketResponse(response, dayOfWeekBucketMap)
						if err != nil {
							mon.TrackDruidTimeMetricInSeconds(mon.DruidAPIMethodDurationType, methodStartTime, errorCode, mon.GetSLAReportStr)
							return nil, err
						}

						query, _, err = SLATimeBucketQuery(request.TenantID, table, metaMOs, HourOfDay, false, request.Timezone, vk, tk, mk, dk, "sla", e, GranularityAll, request.Interval, timeout)
						if err != nil {
							mon.TrackDruidTimeMetricInSeconds(mon.DruidAPIMethodDurationType, methodStartTime, errorCode, mon.GetSLAReportStr)
							return nil, err
						}

						queryStartTime = time.Now()
						if logger.IsDebugEnabled() {
							logger.Log.Debugf("Querying Druid for %s with query: %v", db.SLAReportStr, models.AsJSONString(query))
						}
						response, err = dc.executeQuery(query)
						if err != nil {
							mon.TrackDruidTimeMetricInSeconds(mon.DruidQueryDurationType, queryStartTime, errorCode, mon.SLATimeBucketQueryStr)
							mon.TrackDruidTimeMetricInSeconds(mon.DruidAPIMethodDurationType, methodStartTime, errorCode, mon.GetSLAReportStr)
							return nil, err
						}
						mon.TrackDruidTimeMetricInSeconds(mon.DruidQueryDurationType, queryStartTime, successCode, mon.SLATimeBucketQueryStr)

						hourOfDayBucketMap, err = reformatSLABucketResponse(response, hourOfDayBucketMap)
						if err != nil {
							mon.TrackDruidTimeMetricInSeconds(mon.DruidAPIMethodDurationType, methodStartTime, errorCode, mon.GetSLAReportStr)
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
		ReportScheduleConfig: request.SLAScheduleConfig,
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
	mon.TrackDruidTimeMetricInSeconds(mon.DruidAPIMethodDurationType, methodStartTime, successCode, mon.GetSLAReportStr)
	return &slaReport, nil
}

func (dc *DruidDatastoreClient) GetSLAViolationsQueryAllGranularity(request *metrics.SLAReportRequest, thresholdProfile *tenant.ThresholdProfile, metaMOs []string) ([]byte, metrics.DruidViolationsMap, error) {
	methodStartTime := time.Now()
	if logger.IsDebugEnabled() {
		logger.Log.Debugf("Calling GetSLAReport for request: %v", models.AsJSONString(request))
	}
	table := dc.cfg.GetString(gather.CK_druid_broker_table.String())
	var query godruid.Query

	timeout := request.Timeout
	if timeout == 0 {
		timeout = int32(gather.GetConfig().GetInt(gather.CK_druid_timeoutsms_slareports.String()))
	}

	query, respSchema, err := SLAViolationsQuery(request.TenantID, table, metaMOs, GranularityAll, request.Interval, false, thresholdProfile, timeout)

	if err != nil {
		mon.TrackDruidTimeMetricInSeconds(mon.DruidAPIMethodDurationType, methodStartTime, errorCode, mon.GetSLAReportStr)
		return nil, nil, err
	}

	queryStartTime := time.Now()
	if logger.IsDebugEnabled() {
		logger.Log.Debugf("Querying Druid for %s with query: %v", db.SLAReportStr, models.AsJSONString(query))
	}
	response, err := dc.executeQuery(query)
	if err != nil {
		mon.TrackDruidTimeMetricInSeconds(mon.DruidQueryDurationType, queryStartTime, errorCode, mon.SLAViolationsQueryStr)
		mon.TrackDruidTimeMetricInSeconds(mon.DruidAPIMethodDurationType, methodStartTime, errorCode, mon.GetSLAReportStr)
		return nil, nil, err
	}
	mon.TrackDruidTimeMetricInSeconds(mon.DruidQueryDurationType, queryStartTime, successCode, mon.SLAViolationsQueryStr)

	if logger.IsDebugEnabled() {
		logger.Log.Debugf("Result: %v", db.SLAReportStr, models.AsJSONString(response))
	}

	return response, respSchema, nil
}

func (dc *DruidDatastoreClient) GetSLAViolationsQueryWithGranularity(request *metrics.SLAReportRequest, thresholdProfile *tenant.ThresholdProfile, metaMOs []string) ([]byte, metrics.DruidViolationsMap, error) {
	methodStartTime := time.Now()

	if logger.IsDebugEnabled() {
		logger.Log.Debugf("Calling GetSLAReport for request: %v", models.AsJSONString(request))
	}
	table := dc.cfg.GetString(gather.CK_druid_broker_table.String())
	var query godruid.Query

	timeout := request.Timeout
	if timeout == 0 {
		timeout = int32(gather.GetConfig().GetInt(gather.CK_druid_timeoutsms_slareports.String()))
	}
	query, schema, err := SLAViolationsQuery(request.TenantID, table, metaMOs, request.Granularity, request.Interval, request.IgnoreCleaning, thresholdProfile, timeout)
	if err != nil {
		mon.TrackDruidTimeMetricInSeconds(mon.DruidQueryDurationType, methodStartTime, successCode, mon.GetSLAReportStr)
		return nil, nil, err
	}

	queryStartTime := time.Now()
	if logger.IsDebugEnabled() {
		logger.Log.Debugf("Querying Druid for %s with query: %v", db.SLAReportStr, models.AsJSONString(query))
	}
	response, err := dc.executeQuery(query)
	if err != nil {
		mon.TrackDruidTimeMetricInSeconds(mon.DruidQueryDurationType, queryStartTime, errorCode, mon.SLAViolationsQueryStr)
		mon.TrackDruidTimeMetricInSeconds(mon.DruidAPIMethodDurationType, methodStartTime, errorCode, mon.GetSLAReportStr)
		return nil, nil, err
	}
	mon.TrackDruidTimeMetricInSeconds(mon.DruidQueryDurationType, queryStartTime, successCode, mon.SLAViolationsQueryStr)

	return response, schema, nil
}

func (dc *DruidDatastoreClient) GetTopNTimeByBuckets(request *metrics.SLAReportRequest, extractFn int, vendor, objType, metric, direction, event string, eventAttr *tenant.ThrPrfEventAttrMap, metaMOs []string) ([]byte, metrics.DruidViolationsMap, error) {
	methodStartTime := time.Now()
	if logger.IsDebugEnabled() {
		logger.Log.Debugf("Calling GetSLAReport for request: %v", models.AsJSONString(request))
	}
	table := dc.cfg.GetString(gather.CK_druid_broker_table.String())

	timeout := request.Timeout
	if timeout == 0 {
		timeout = int32(gather.GetConfig().GetInt(gather.CK_druid_timeoutsms_slareports.String()))
	}
	query, schema, err := SLATimeBucketQuery(request.TenantID, table, metaMOs /*DayOfWeek*/, extractFn, false, request.Timezone, vendor, objType, metric, direction /*event*/, "sla", eventAttr, GranularityAll, request.Interval, timeout)
	if err != nil {
		mon.TrackDruidTimeMetricInSeconds(mon.DruidAPIMethodDurationType, methodStartTime, errorCode, mon.GetSLAReportStr)
		return nil, nil, err
	}

	if logger.IsDebugEnabled() {
		logger.Log.Debugf("Querying Druid for %s with query: %v", db.SLAReportStr, models.AsJSONString(schema))
	}

	queryStartTime := time.Now()

	response, err := dc.executeQuery(query)

	if err != nil {
		mon.TrackDruidTimeMetricInSeconds(mon.DruidQueryDurationType, queryStartTime, errorCode, mon.SLATimeBucketQueryStr)
		mon.TrackDruidTimeMetricInSeconds(mon.DruidAPIMethodDurationType, methodStartTime, errorCode, mon.GetSLAReportStr)
		return nil, nil, err
	}
	mon.TrackDruidTimeMetricInSeconds(mon.DruidQueryDurationType, queryStartTime, successCode, mon.SLATimeBucketQueryStr)

	return response, schema, nil
}

func (dc *DruidDatastoreClient) GetSLATimeSeries(request *metrics.SLAReportRequest, thresholdProfile *tenant.ThresholdProfile, metaMOs []string) (map[string]interface{}, map[string]interface{}, error) {
	methodStartTime := time.Now()
	if logger.IsDebugEnabled() {
		logger.Log.Debugf("Calling GetSLAReport for request: %v", models.AsJSONString(request))
	}
	table := dc.cfg.GetString(gather.CK_druid_broker_table.String())

	timeout := request.Timeout
	if timeout == 0 {
		timeout = int32(gather.GetConfig().GetInt(gather.CK_druid_timeoutsms_slareports.String()))
	}

	var hourOfDayBucketMap map[string]interface{}
	var dayOfWeekBucketMap map[string]interface{}

	for vk, v := range thresholdProfile.Thresholds.VendorMap {
		for tk, t := range v.MonitoredObjectTypeMap {

			for mk, m := range t.MetricMap {
				for dk, d := range m.DirectionMap {
					for ek, e := range d.EventMap {
						if ek != "sla" {
							continue
						}

						// split into its own function
						query, schema, err := SLATimeBucketQuery(request.TenantID, table, metaMOs, DayOfWeek, false, request.Timezone, vk, tk, mk, dk, "sla", e, GranularityAll, request.Interval, timeout)
						if err != nil {
							mon.TrackDruidTimeMetricInSeconds(mon.DruidAPIMethodDurationType, methodStartTime, errorCode, mon.GetSLAReportStr)
							return nil, nil, err
						}
						logger.Log.Debugf("Querying Druid for %s with query: %v", db.SLAReportStr, models.AsJSONString(schema))

						queryStartTime := time.Now()

						response, err := dc.executeQuery(query)
						if err != nil {
							mon.TrackDruidTimeMetricInSeconds(mon.DruidQueryDurationType, queryStartTime, errorCode, mon.SLATimeBucketQueryStr)
							mon.TrackDruidTimeMetricInSeconds(mon.DruidAPIMethodDurationType, methodStartTime, errorCode, mon.GetSLAReportStr)
							return nil, nil, err
						}
						mon.TrackDruidTimeMetricInSeconds(mon.DruidQueryDurationType, queryStartTime, successCode, mon.SLATimeBucketQueryStr)

						dayOfWeekBucketMap, err = reformatSLABucketResponse(response, dayOfWeekBucketMap)

						if err != nil {
							mon.TrackDruidTimeMetricInSeconds(mon.DruidAPIMethodDurationType, methodStartTime, errorCode, mon.GetSLAReportStr)
							return nil, nil, err
						}

						query, schema, err = SLATimeBucketQuery(request.TenantID, table, metaMOs, HourOfDay, false, request.Timezone, vk, tk, mk, dk, "sla", e, GranularityAll, request.Interval, timeout)
						if err != nil {
							mon.TrackDruidTimeMetricInSeconds(mon.DruidAPIMethodDurationType, methodStartTime, errorCode, mon.GetSLAReportStr)
							return nil, nil, err
						}

						queryStartTime = time.Now()
						if logger.IsDebugEnabled() {
							logger.Log.Debugf("Querying Druid for %s with query: %v", db.SLAReportStr, models.AsJSONString(query))
						}
						response, err = dc.executeQuery(query)
						if err != nil {
							mon.TrackDruidTimeMetricInSeconds(mon.DruidQueryDurationType, queryStartTime, errorCode, mon.SLATimeBucketQueryStr)
							mon.TrackDruidTimeMetricInSeconds(mon.DruidAPIMethodDurationType, methodStartTime, errorCode, mon.GetSLAReportStr)
							return nil, nil, err
						}
						mon.TrackDruidTimeMetricInSeconds(mon.DruidQueryDurationType, queryStartTime, successCode, mon.SLATimeBucketQueryStr)

						hourOfDayBucketMap, err = reformatSLABucketResponse(response, hourOfDayBucketMap)
						if err != nil {
							mon.TrackDruidTimeMetricInSeconds(mon.DruidAPIMethodDurationType, methodStartTime, errorCode, mon.GetSLAReportStr)
							return nil, nil, err
						}
					}
				}
			}
		}
	}

	return hourOfDayBucketMap, dayOfWeekBucketMap, nil
}

func (dc *DruidDatastoreClient) GetSLATimeSeriesV1(request *metrics.SLAReportRequest, thresholdProfile *tenant.ThresholdProfile, metaMOs []string) (map[string]interface{}, map[string]interface{}, error) {
	methodStartTime := time.Now()
	var err error
	if logger.IsDebugEnabled() {
		logger.Log.Debugf("Calling GetSLAReport for request: %v", models.AsJSONString(request))
	}
	table := dc.cfg.GetString(gather.CK_druid_broker_table.String())
	var query godruid.Query

	timeout := request.Timeout
	if timeout == 0 {
		timeout = int32(gather.GetConfig().GetInt(gather.CK_druid_timeoutsms_slareports.String()))
	}

	var hourOfDayBucketMap map[string]interface{}
	var dayOfWeekBucketMap map[string]interface{}

	for vk, v := range thresholdProfile.Thresholds.VendorMap {
		for tk, t := range v.MonitoredObjectTypeMap {

			for mk, m := range t.MetricMap {
				for dk, d := range m.DirectionMap {
					for ek, e := range d.EventMap {
						if ek != "sla" {
							continue
						}

						// split into its own function
						query, _, err = SLATimeBucketQuery(request.TenantID, table, metaMOs, DayOfWeek, false, request.Timezone, vk, tk, mk, dk, "sla", e, GranularityAll, request.Interval, timeout)
						if err != nil {
							mon.TrackDruidTimeMetricInSeconds(mon.DruidAPIMethodDurationType, methodStartTime, errorCode, mon.GetSLAReportStr)
							return nil, nil, err
						}

						queryStartTime := time.Now()
						if logger.IsDebugEnabled() {
							logger.Log.Debugf("Querying Druid for %s with query: %v", db.SLAReportStr, models.AsJSONString(query))
						}
						response, err := dc.executeQuery(query)
						logger.Log.Debugf("dayOfWeekBucketMap: %s", models.AsJSONString(string(response)))
						if err != nil {
							mon.TrackDruidTimeMetricInSeconds(mon.DruidQueryDurationType, queryStartTime, errorCode, mon.SLATimeBucketQueryStr)
							mon.TrackDruidTimeMetricInSeconds(mon.DruidAPIMethodDurationType, methodStartTime, errorCode, mon.GetSLAReportStr)
							return nil, nil, err
						}
						mon.TrackDruidTimeMetricInSeconds(mon.DruidQueryDurationType, queryStartTime, successCode, mon.SLATimeBucketQueryStr)

						dayOfWeekBucketMap, err = reformatSLABucketResponse(response, dayOfWeekBucketMap)

						logger.Log.Debugf(" rendered dayOfWeekBucketMap: %s", dayOfWeekBucketMap)
						if err != nil {
							mon.TrackDruidTimeMetricInSeconds(mon.DruidAPIMethodDurationType, methodStartTime, errorCode, mon.GetSLAReportStr)
							return nil, nil, err
						}

						query, _, err = SLATimeBucketQuery(request.TenantID, table, metaMOs, HourOfDay, false, request.Timezone, vk, tk, mk, dk, "sla", e, GranularityAll, request.Interval, timeout)
						if err != nil {
							mon.TrackDruidTimeMetricInSeconds(mon.DruidAPIMethodDurationType, methodStartTime, errorCode, mon.GetSLAReportStr)
							return nil, nil, err
						}

						queryStartTime = time.Now()
						if logger.IsDebugEnabled() {
							logger.Log.Debugf("Querying Druid for %s with query: %v", db.SLAReportStr, models.AsJSONString(query))
						}
						response, err = dc.executeQuery(query)
						if err != nil {
							mon.TrackDruidTimeMetricInSeconds(mon.DruidQueryDurationType, queryStartTime, errorCode, mon.SLATimeBucketQueryStr)
							mon.TrackDruidTimeMetricInSeconds(mon.DruidAPIMethodDurationType, methodStartTime, errorCode, mon.GetSLAReportStr)
							return nil, nil, err
						}
						mon.TrackDruidTimeMetricInSeconds(mon.DruidQueryDurationType, queryStartTime, successCode, mon.SLATimeBucketQueryStr)

						hourOfDayBucketMap, err = reformatSLABucketResponse(response, hourOfDayBucketMap)
						if err != nil {
							mon.TrackDruidTimeMetricInSeconds(mon.DruidAPIMethodDurationType, methodStartTime, errorCode, mon.GetSLAReportStr)
							return nil, nil, err
						}
					}
				}
			}
		}
	}

	return hourOfDayBucketMap, dayOfWeekBucketMap, nil
}

// GetSLAReportV2 - Get SLA report but returns in a v2 compatible format
func (dc *DruidDatastoreClient) GetSLAReportV2(request *metrics.SLAReportRequest, thresholdProfile *tenant.ThresholdProfile, metaMOs []string) (*metrics.SLAReport, error) {
	methodStartTime := time.Now()
	if logger.IsDebugEnabled() {
		logger.Log.Debugf("Calling GetSLAReport for request: %v", models.AsJSONString(request))
	}
	table := dc.cfg.GetString(gather.CK_druid_broker_table.String())
	var query godruid.Query

	timeout := request.Timeout
	if timeout == 0 {
		timeout = int32(gather.GetConfig().GetInt(gather.CK_druid_timeoutsms_slareports.String()))
	}

	query, _, err := SLAViolationsQuery(request.TenantID, table, metaMOs, GranularityAll, request.Interval, request.IgnoreCleaning, thresholdProfile, timeout)

	if err != nil {
		mon.TrackDruidTimeMetricInSeconds(mon.DruidAPIMethodDurationType, methodStartTime, errorCode, mon.GetSLAReportStr)
		return nil, err
	}

	queryStartTime := time.Now()
	if logger.IsDebugEnabled() {
		logger.Log.Debugf("Querying Druid for %s with query: %v", db.SLAReportStr, models.AsJSONString(query))
	}
	response, err := dc.executeQuery(query)
	if err != nil {
		mon.TrackDruidTimeMetricInSeconds(mon.DruidQueryDurationType, queryStartTime, errorCode, mon.SLAViolationsQueryStr)
		mon.TrackDruidTimeMetricInSeconds(mon.DruidAPIMethodDurationType, methodStartTime, errorCode, mon.GetSLAReportStr)
		return nil, err
	}
	mon.TrackDruidTimeMetricInSeconds(mon.DruidQueryDurationType, queryStartTime, successCode, mon.SLAViolationsQueryStr)

	reportSummary, err := reformatReportSummary(response)
	if err != nil {
		mon.TrackDruidTimeMetricInSeconds(mon.DruidQueryDurationType, queryStartTime, successCode, mon.GetSLAReportStr)
		return nil, err
	}
	if logger.IsDebugEnabled() {
		logger.Log.Debugf("Result: %v", db.SLAReportStr, models.AsJSONString(reportSummary))
	}

	query, _, err = SLAViolationsQuery(request.TenantID, table, metaMOs, request.Granularity, request.Interval, request.IgnoreCleaning, thresholdProfile, timeout)
	if err != nil {
		mon.TrackDruidTimeMetricInSeconds(mon.DruidQueryDurationType, queryStartTime, successCode, mon.GetSLAReportStr)
		return nil, err
	}

	queryStartTime = time.Now()
	if logger.IsDebugEnabled() {
		logger.Log.Debugf("Querying Druid for %s with query: %v", db.SLAReportStr, models.AsJSONString(query))
	}
	response, err = dc.executeQuery(query)
	if err != nil {
		mon.TrackDruidTimeMetricInSeconds(mon.DruidQueryDurationType, queryStartTime, errorCode, mon.SLAViolationsQueryStr)
		mon.TrackDruidTimeMetricInSeconds(mon.DruidAPIMethodDurationType, methodStartTime, errorCode, mon.GetSLAReportStr)
		return nil, err
	}
	mon.TrackDruidTimeMetricInSeconds(mon.DruidQueryDurationType, queryStartTime, successCode, mon.SLAViolationsQueryStr)

	slaTimeSeries, err := reformatSLATimeSeries(response)
	if err != nil {
		mon.TrackDruidTimeMetricInSeconds(mon.DruidAPIMethodDurationType, methodStartTime, errorCode, mon.GetSLAReportStr)
		return nil, err
	}

	var hourOfDayBucketMap map[string]interface{}
	var dayOfWeekBucketMap map[string]interface{}

	for vk, v := range thresholdProfile.Thresholds.VendorMap {
		for tk, t := range v.MonitoredObjectTypeMap {
			if tk != "twamp-sf" {
				continue
			}

			for mk, m := range t.MetricMap {
				for dk, d := range m.DirectionMap {
					for ek, e := range d.EventMap {
						if ek != "sla" {
							continue
						}
						query, _, err = SLATimeBucketQuery(request.TenantID, table, metaMOs, DayOfWeek, request.IgnoreCleaning, request.Timezone, vk, tk, mk, dk, "sla", e, GranularityAll, request.Interval, timeout)
						if err != nil {
							mon.TrackDruidTimeMetricInSeconds(mon.DruidAPIMethodDurationType, methodStartTime, errorCode, mon.GetSLAReportStr)
							return nil, err
						}

						queryStartTime = time.Now()
						if logger.IsDebugEnabled() {
							logger.Log.Debugf("Querying Druid for %s with query: %v", db.SLAReportStr, models.AsJSONString(query))
						}
						response, err = dc.executeQuery(query)
						if err != nil {
							mon.TrackDruidTimeMetricInSeconds(mon.DruidQueryDurationType, queryStartTime, errorCode, mon.SLATimeBucketQueryStr)
							mon.TrackDruidTimeMetricInSeconds(mon.DruidAPIMethodDurationType, methodStartTime, errorCode, mon.GetSLAReportStr)
							return nil, err
						}
						mon.TrackDruidTimeMetricInSeconds(mon.DruidQueryDurationType, queryStartTime, successCode, mon.SLATimeBucketQueryStr)

						dayOfWeekBucketMap, err = reformatSLABucketResponse(response, dayOfWeekBucketMap)
						if err != nil {
							mon.TrackDruidTimeMetricInSeconds(mon.DruidAPIMethodDurationType, methodStartTime, errorCode, mon.GetSLAReportStr)
							return nil, err
						}

						query, _, err = SLATimeBucketQuery(request.TenantID, table, metaMOs, HourOfDay, request.IgnoreCleaning, request.Timezone, vk, tk, mk, dk, "sla", e, GranularityAll, request.Interval, timeout)
						if err != nil {
							mon.TrackDruidTimeMetricInSeconds(mon.DruidAPIMethodDurationType, methodStartTime, errorCode, mon.GetSLAReportStr)
							return nil, err
						}

						queryStartTime = time.Now()
						if logger.IsDebugEnabled() {
							logger.Log.Debugf("Querying Druid for %s with query: %v", db.SLAReportStr, models.AsJSONString(query))
						}
						response, err = dc.executeQuery(query)
						if err != nil {
							mon.TrackDruidTimeMetricInSeconds(mon.DruidQueryDurationType, queryStartTime, errorCode, mon.SLATimeBucketQueryStr)
							mon.TrackDruidTimeMetricInSeconds(mon.DruidAPIMethodDurationType, methodStartTime, errorCode, mon.GetSLAReportStr)
							return nil, err
						}
						mon.TrackDruidTimeMetricInSeconds(mon.DruidQueryDurationType, queryStartTime, successCode, mon.SLATimeBucketQueryStr)

						hourOfDayBucketMap, err = reformatSLABucketResponse(response, hourOfDayBucketMap)
						if err != nil {
							mon.TrackDruidTimeMetricInSeconds(mon.DruidAPIMethodDurationType, methodStartTime, errorCode, mon.GetSLAReportStr)
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
		ReportScheduleConfig: request.SLAScheduleConfig,
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
	mon.TrackDruidTimeMetricInSeconds(mon.DruidAPIMethodDurationType, methodStartTime, successCode, mon.GetSLAReportStr)
	return &slaReport, nil
}

// DEPRECATED - Remove once v1 queries have been moved
func (dc *DruidDatastoreClient) GetRawMetricsV1(request *pb.RawMetricsRequest) (map[string]interface{}, error) {
	methodStartTime := time.Now()
	if logger.IsDebugEnabled() {
		logger.Log.Debugf("Calling GetRawMetrics for request: %v", models.AsJSONString(request))
	}

	table := dc.cfg.GetString(gather.CK_druid_broker_table.String())

	// peyo TODO we should have a better way to handle default query params
	timeout := request.GetTimeout()
	if timeout == 0 {
		timeout = int32(gather.GetConfig().GetInt(gather.CK_druid_timeoutsms_rawmetrics.String()))
	}

	granularity := request.GetGranularity()
	if granularity == "" {
		granularity = "PT1M"
	}

	query, err := RawMetricsQueryV1(request.GetTenant(), table, request.Metric, request.GetInterval(), request.GetObjectType(), request.GetDirection(), request.GetMonitoredObjectId(), timeout, granularity, request.GetCleanOnly())

	if err != nil {
		mon.TrackDruidTimeMetricInSeconds(mon.DruidAPIMethodDurationType, methodStartTime, errorCode, mon.GetRawMetricStr)
		return nil, err
	}

	queryStartTime := time.Now()
	if logger.IsDebugEnabled() {
		logger.Log.Debugf("Querying Druid for %s with query: '' %s ''", db.RawMetricStr, models.AsJSONString(query))
	}
	response, err := dc.executeQuery(query)

	if err != nil {
		mon.TrackDruidTimeMetricInSeconds(mon.DruidQueryDurationType, queryStartTime, errorCode, mon.GetRawMetricStr)
		mon.TrackDruidTimeMetricInSeconds(mon.DruidAPIMethodDurationType, methodStartTime, errorCode, mon.GetRawMetricStr)
		return nil, err
	}
	mon.TrackDruidTimeMetricInSeconds(mon.DruidQueryDurationType, queryStartTime, successCode, mon.GetRawMetricStr)

	resp := make([]map[string]interface{}, 0)

	err = json.Unmarshal(response, &resp)
	if err != nil {
		mon.TrackDruidTimeMetricInSeconds(mon.DruidAPIMethodDurationType, methodStartTime, errorCode, mon.GetRawMetricStr)
		return nil, err
	}

	if logger.IsDebugEnabled() {
		logger.Log.Debugf("Response from druid for %s: %v", db.RawMetricStr, models.AsJSONString(resp))
	}

	mon.TrackDruidTimeMetricInSeconds(mon.DruidAPIMethodDurationType, methodStartTime, successCode, mon.GetRawMetricStr)
	return map[string]interface{}{"results": resp}, nil
}

// DEPRECATED - Remove once COLT is no longer dependent on this API
func (dc *DruidDatastoreClient) GetFilteredRawMetrics(request *metrics.RawMetrics, metaMOs []string) (map[string]interface{}, error) {
	methodStartTime := time.Now()
	if logger.IsDebugEnabled() {
		logger.Log.Debugf("Calling GetFilteredRawMetrics for request: %v", models.AsJSONString(request))
	}

	table := dc.cfg.GetString(gather.CK_druid_broker_table.String())

	timeout := request.Timeout
	if timeout == 0 {
		timeout = int32(gather.GetConfig().GetInt(gather.CK_druid_timeoutsms_rawmetrics.String()))
	}

	granularity := request.Granularity
	if granularity == "" {
		granularity = "PT1M"
	}

	query, err := RawMetricsQueryV1(request.TenantID, table, request.Metrics, request.Interval, request.ObjectType, request.Directions, metaMOs, timeout, granularity, false)

	if err != nil {
		mon.TrackDruidTimeMetricInSeconds(mon.DruidAPIMethodDurationType, methodStartTime, errorCode, mon.GetRawMetricStr)
		return nil, err
	}

	queryStartTime := time.Now()
	if logger.IsDebugEnabled() {
		logger.Log.Debugf("Querying Druid for %s with query: '' %s ''", db.RawMetricStr, models.AsJSONString(query))
	}
	response, err := dc.executeQuery(query)

	if err != nil {
		mon.TrackDruidTimeMetricInSeconds(mon.DruidQueryDurationType, queryStartTime, errorCode, mon.GetRawMetricStr)
		mon.TrackDruidTimeMetricInSeconds(mon.DruidAPIMethodDurationType, methodStartTime, errorCode, mon.GetRawMetricStr)
		return nil, err
	}
	mon.TrackDruidTimeMetricInSeconds(mon.DruidQueryDurationType, queryStartTime, successCode, mon.GetRawMetricStr)

	resp := make([]metrics.TimeseriesEntryResponse, 0)

	err = json.Unmarshal(response, &resp)
	if err != nil {
		mon.TrackDruidTimeMetricInSeconds(mon.DruidAPIMethodDurationType, methodStartTime, errorCode, mon.GetRawMetricStr)
		return nil, err
	}

	if logger.IsDebugEnabled() {
		logger.Log.Debugf("Response from druid for %s: %v", db.RawMetricStr, models.AsJSONString(resp))
	}

	formattedJSON := map[string]interface{}{}
	if len(resp) != 0 {
		formattedJSON, err = reformatRawMetricsResponse(resp)
	}

	if err != nil {
		mon.TrackDruidTimeMetricInSeconds(mon.DruidAPIMethodDurationType, methodStartTime, errorCode, mon.GetRawMetricStr)
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

	mon.TrackDruidTimeMetricInSeconds(mon.DruidAPIMethodDurationType, methodStartTime, successCode, mon.GetRawMetricStr)
	return rr, nil
}

type ActiveRuleDetails struct {
	startTime int64
	endTime   int64
	isError   bool
	testFreq  int64
}

func (dc *DruidDatastoreClient) GetDataCleaningHistory(tenantID string, monitoredObjectID string, interval string) ([]*swagmodels.DataCleaningTransition, error) {
	methodStartTime := time.Now()
	if logger.IsDebugEnabled() {
		logger.Log.Debugf("Calling GetRaGetDataCleaningHistorywMetrics for tenant %s monitoredObject %s interval %s", tenantID, monitoredObjectID, interval)
	}
	table := dc.cfg.GetString(gather.CK_druid_broker_table.String())

	query := &godruid.QueryScan{
		DataSource:   table,
		Intervals:    interval,
		ResultFormat: "list",
		Filter: godruid.FilterAnd(
			godruid.FilterSelector("tenantId", strings.ToLower(tenantID)),
			godruid.FilterSelector("monitoredObjectId", monitoredObjectID),
			godruid.FilterSelector("cleanStatus", "-1"),
		),
		Columns: []string{"__time", "direction", "errorCode", "failedRules", "duration"},
	}

	queryStartTime := time.Now()
	if logger.IsDebugEnabled() {
		logger.Log.Debugf("Querying Druid for %s with query: '' %s ''", db.DataCleaningStr, models.AsJSONString(query))
	}

	_, err := dc.executeQuery(query)

	if err != nil {
		mon.TrackDruidTimeMetricInSeconds(mon.DruidQueryDurationType, queryStartTime, errorCode, mon.GetDataCleaningHistoryStr)
		mon.TrackDruidTimeMetricInSeconds(mon.DruidAPIMethodDurationType, methodStartTime, errorCode, mon.GetDataCleaningHistoryStr)
		return nil, err
	}
	mon.TrackDruidTimeMetricInSeconds(mon.DruidQueryDurationType, queryStartTime, successCode, mon.GetDataCleaningHistoryStr)

	scanBlobs := query.QueryResult

	if logger.IsDebugEnabled() {
		logger.Log.Debugf("got result %v", scanBlobs)
	}

	// The scan blobs have provided all the rows that were tagged for cleaning and the reason (failed rules or error code).
	// To reduce this to a series of raise/clear events they are processed in chronological order.  When a rule/errorCode
	// is first seen, this is a 'raise'.  When a rule/errorCode is seen immediately after the last time it was seen (i.e. continuous)
	// the timestamp for the 'clear' is increased.  When a rule/erroCode is seen again after some time gap, the 'clear'
	// for the previous occurance of the rule/errorCode is created and the 'raise' for the new occurance is created.
	// Basically we are detecting start/end events of when a rule/erroCode was active.

	// Make sure segments are processed chronologically.
	sort.Slice(scanBlobs, func(i, j int) bool { return scanBlobs[i].SegmentID < scanBlobs[j].SegmentID })
	activeRules := map[string]map[string]*ActiveRuleDetails{}
	history := map[int64]*swagmodels.DataCleaningTransition{}
	for _, b := range scanBlobs {

		// Make sure events are processed chronologially.
		sort.Slice(b.Events, func(i, j int) bool { return b.Events[i]["__time"].(float64) < b.Events[j]["__time"].(float64) })

		for _, e := range b.Events {
			logger.Log.Debugf("processing event %v", e)

			// Coerse the properties into the appropriate data types.
			timestamp := int64(e["__time"].(float64))
			direction := e["direction"].(string)
			duration := int64(e["duration"].(float64))
			errCode := int32(e["errorCode"].(float64))

			// Extract the failed rule(s) string. This will be the key for tracking active rules.
			failedRules := parseFailedRules(e["failedRules"])

			if len(failedRules) == 0 && errCode == 0 {
				// Really shouldn't happen.
				continue
			}

			for _, failedRule := range failedRules {

				if _, ok := activeRules[direction]; !ok {
					// Add an entry for this direction
					activeRules[direction] = map[string]*ActiveRuleDetails{}
				}

				// Check if the failed rule is actually an indication of an invalid test (errorCode != 0)
				isErrCodeFailure := false
				if failedRule == "msg:errorCode!=0" {
					if errCode == 0 {
						continue
					}
					isErrCodeFailure = true
					failedRule = fmt.Sprintf("%d", errCode)
				}

				if activeRule, ok := activeRules[direction][failedRule]; ok {
					// We've seen this rule before. It's either a continuation or a new occurance
					logger.Log.Debugf("entry found for for direction %v, ruleKey %v, value %v", direction, failedRule, *activeRule)
					//TODO not 100% sure using duration is good enough in case there are small gaps of missing reports.
					//Maybe we need to use the duration in the rule.
					if timestamp-duration <= activeRule.endTime {
						// It's a continuation. Extend the end time so we can keep track of when the clear happens.
						activeRule.endTime = timestamp + duration
					} else {
						// It's a whole new trigger period for this rule.
						// Add the clear event to history.
						updateHistory(history, activeRule.endTime, direction, false, failedRule, isErrCodeFailure)

						// Reset the active rule for the new raise
						activeRule.startTime = timestamp
						activeRule.endTime = timestamp + duration

						// Add the raise event to history
						updateHistory(history, activeRule.startTime, direction, true, failedRule, isErrCodeFailure)
					}

				} else {
					// First time seeing this rule. Start an active rule and add the raise event to history.
					activeRules[direction][failedRule] = &ActiveRuleDetails{
						startTime: timestamp,
						endTime:   timestamp + duration,
						isError:   isErrCodeFailure,
						testFreq:  duration,
					}
					updateHistory(history, timestamp, direction, true, failedRule, isErrCodeFailure)
				}
			}
		}
	}

	// All the active rules now need a clear event for closure.
	for direction, ruleSet := range activeRules {
		for rule, event := range ruleSet {
			updateHistory(history, event.endTime, direction, false, rule, event.isError)
		}
	}

	// Extract array and sort by time.
	res := make([]*swagmodels.DataCleaningTransition, len(history))
	i := 0
	for _, v := range history {
		res[i] = v
		i++
	}
	sort.Slice(res, func(i, j int) bool { return res[i].Timestamp < res[j].Timestamp })

	return res, nil
}

func updateHistory(history map[int64]*swagmodels.DataCleaningTransition, timestamp int64, dir string, raised bool, rule string, isErrCode bool) {
	var transition, ok = history[timestamp]
	if !ok {
		transition = &swagmodels.DataCleaningTransition{Timestamp: timestamp}
		history[timestamp] = transition
	}

	if isErrCode {
		errCodeTrans := &swagmodels.DataCleaningTransitionError{ErrorCode: rule, Direction: dir}
		if raised {
			transition.ErrorsRaised = append(transition.ErrorsRaised, errCodeTrans)
		} else {
			transition.ErrorsCleared = append(transition.ErrorsCleared, errCodeTrans)
		}
	} else {
		parsedRule := swagmodels.DataCleaningRule{}
		err := json.Unmarshal([]byte(rule), &parsedRule)
		if err != nil {
			logger.Log.Warningf("Failed to parse cleaning rule %s", rule)
		}
		ruleTrans := &swagmodels.DataCleaningTransitionRule{Rule: &parsedRule, Direction: dir}
		if raised {
			transition.RulesRaised = append(transition.RulesRaised, ruleTrans)
		} else {
			transition.RulesCleared = append(transition.RulesCleared, ruleTrans)
		}
	}

}

func parseFailedRules(input interface{}) []string {
	if input != nil {
		if singleRule, ok := input.(string); ok {
			return []string{singleRule}
		}

		result := []string{}
		for _, val := range input.([]interface{}) {
			result = append(result, val.(string))
		}
		return result
	}

	return []string{}
}

type lookup struct {
	Version                string    `json:"version"`
	LookupExtractorFactory mapLookup `json:"lookupExtractorFactory"`
	active                 bool
	count                  int
}

type mapLookup struct {
	LookupType string            `json:"type"`
	Data       map[string]string `json:"map"`
}

func buildLookup(datatype, tenant, key, val string, partition int) *lookup {

	version := time.Now().Format(time.RFC3339)
	domLookup := lookup{
		count:   0,
		Version: version,
		LookupExtractorFactory: mapLookup{
			LookupType: "map",
			Data:       map[string]string{},
		},
	}

	return &domLookup
}

func (dc *DruidDatastoreClient) buildNewDruidLookup(datatype, tenant, key, val, itemKey, itemVal string, partition int) *lookup {
	// Can't find lookup, so create it
	domLookup := buildLookup(datatype, tenant, key, val, partition)
	// Now add the first item for this lookup
	domLookup.LookupExtractorFactory.Data[itemKey] = itemVal
	domLookup.count = 1
	return domLookup
}

func (dc *DruidDatastoreClient) addToLookup(lookups map[string]*lookup, existingLookups DruidLookupStatus, datatype, tenant, key, val, itemKey, itemVal string, partition int) {

	lookupName := buildLookupName(datatype, tenant, key, val, partition)
	domLookup, ok := lookups[lookupName]
	if ok {
		// Ok, we have an existing look up for this key, check if there are too many values in this bucket and needs to spill over
		if domLookup.count >= 50000 {
			// Yes, there are too many items, spill over to the next bucket
			dc.addToLookup(lookups, existingLookups, datatype, tenant, key, val, itemKey, itemVal, partition+1)
		} else {
			// No, we can continue to use this bucket
			domLookup.LookupExtractorFactory.Data[itemKey] = lookupName
			domLookup.count = domLookup.count + 1
			existingLookups[lookupName] = true
		}
	} else {
		// First time encountering this item, let's create a lookup for it
		newLookup := dc.buildNewDruidLookup(datatype, tenant, key, val, itemKey, itemVal, partition)
		// Now append the lookups
		lookups[lookupName] = newLookup
		existingLookups[lookupName] = true
	}
}

/*
Order of operations
* Get all the monitored objects we want to work with
* Walk the monitored objects
	* Go through all the monitored object's metadata and add it to the lookup table
	* Add them to the lookups
	* if lookup is >= 50,000 rows
		* Spill over to the next lookup bucket
* For each lookups
		* Delete existing lookups (cleans up the lookups)
		* Send each Lookup  to druid

*/
func (dc *DruidDatastoreClient) updateMetadataLookup(lookupEndpoint string, tenantID string, datatype string, monitoredObjects []*tenant.MonitoredObject) (map[string]*lookup, error) {

	methodStartTime := time.Now()

	existingNames, err := dc.GetDruidLookupNames()
	if err != nil {
		mon.TrackDruidTimeMetricInSeconds(mon.DruidAPIMethodDurationType, methodStartTime, errorCode, mon.UpdateDruidMetaLookups)
		return nil, fmt.Errorf("Failed to update lookup %s", err.Error())
	}

	//logger.Log.Debugf("Lookups to add for %+v", lookups)
	lookups := make(map[string]*lookup, 0)

	// debugging

	// Now fill in the contents of each lookup by traversing the monitoredObject-to-domain associations.
	for _, mo := range monitoredObjects {
		dc.addToLookup(lookups, existingNames, datatype, tenantID, "allobjs", "test", mo.MonitoredObjectID, mo.MonitoredObjectID, 0)
		// Special exception case for domains

		for key, val := range mo.Meta {
			dc.addToLookup(lookups, existingNames, datatype, tenantID, key, val, mo.MonitoredObjectID, mo.MonitoredObjectID, 0)
		}
	}

	// Now delete all the orphaned lookups
	// We delete the active lookups only when we're about to update them.
	for name, _ := range existingNames {
		// Delete the look up first (don't worry about errors, best effort)
		//if status == false {
		dc.deleteItemToLookup(lookupEndpoint, name)
		//}
	}

	// Debugging only
	if logger.IsDebugEnabled() && false {
		for key, val := range lookups {
			logger.Log.Debugf("{%s,\t%+v }", key, val.LookupExtractorFactory.Data)
		}
	}

	// now post it
	// Domain lookups are assigned to the __default tier
	// The second argument is empty because lookupname is already part of the request

	logger.Log.Infof("Sending Lookup table to druid")
	waitForCompletion := make(chan string, 25)
	for key, val := range lookups {
		logger.Log.Debugf("Sending lookups to druid", key)

		val.active = true
		// Looks up are costly, let's see if we can parallalize the operations
		go func(look string, key string, val *lookup, waitForCompletion chan string) {
			// Update the lookups
			err := dc.addItemToLookup(lookupEndpoint, key, val)
			if err != nil {
				mon.TrackDruidTimeMetricInSeconds(mon.DruidAPIMethodDurationType, methodStartTime, errorCode, mon.UpdateDruidMetaLookups)
				logger.Log.Errorf("Failed to update lookup %s", err.Error())
			}
			waitForCompletion <- key
		}(lookupEndpoint, key, val, waitForCompletion)
	}
	for {
		select {
		case key := <-waitForCompletion:
			lookups[key].active = false
			// If lookups are still active, wait for the next response
			for _, lk := range lookups {
				if lk.active {
					break
				}
			}
			// No active lookups, return success
			mon.TrackDruidTimeMetricInSeconds(mon.DruidAPIMethodDurationType, methodStartTime, successCode, mon.UpdateDruidMetaLookups)
			return lookups, nil
		case <-time.After(5 * 60 * time.Second):
			mon.TrackDruidTimeMetricInSeconds(mon.DruidAPIMethodDurationType, methodStartTime, errorCode, mon.UpdateDruidMetaLookups)
			return nil, fmt.Errorf("Timed out trying to update lookup tables, lookups:%s", models.AsJSONString(lookups))
		}
	}

	// mon.TrackDruidTimeMetricInSeconds(mon.DruidAPIMethodDurationType, methodStartTime, successCode, mon.UpdateDruidMetaLookups)
	// return lookups, nil
}

// GetDruidLookupFor - Returns a list of lookup tables with partial matches
func (dc *DruidDatastoreClient) GetDruidLookupFor(keyval string) ([]string, error) {

	var matches []string
	list, err := dc.GetDruidLookupNames()
	if err != nil {
		return nil, fmt.Errorf("Could not get look up list from druid %s", err.Error())
	}

	for k := range list {
		if idx := strings.Index(k, keyval); idx != -1 {
			// Check that the next character is a * because that is how we delim the names
			if k[idx+1] == druidLookupSeparator[0] {
				matches = append(matches, k)
			}
		}
	}
	return matches, nil
}

// DruidLookupStatus - This is a map of all the active lookups tables in Druid
type DruidLookupStatus map[string]bool

// GetDruidLookupNames - Goes to Druid and grabs the list of monitored objects for a lookup
func (dc *DruidDatastoreClient) GetDruidLookupNames() (DruidLookupStatus, error) {

	methodStartTime := time.Now()

	lookupEndpoint := dc.generateDruidCoordinatorURI(druidLookUpConfig, druidLookUpTierName)

	logger.Log.Infof("Getting Druid look up: %s", lookupEndpoint)

	resp, err := sendRequest("GET", dc.dClient.HttpClient, lookupEndpoint, dc.AuthToken, nil)
	if err != nil {
		mon.TrackDruidTimeMetricInSeconds(mon.DruidAPIMethodDurationType, methodStartTime, errorCode, mon.GetDruidLookups)
		return nil, fmt.Errorf("Could not get look up list from druid %s", err.Error())
	}

	// Druid returns an array of strings which is a problem because go/encoding/json expects a starting brace
	// According to https://stackoverflow.com/questions/5034444/can-json-start-with druid is sending valid json
	construct := fmt.Sprintf("{\"results\":%s}", string(resp))

	type convert2json struct {
		Results []string `json:"results"`
	}
	var responseMap convert2json
	if err = json.Unmarshal([]byte(construct), &responseMap); err != nil {
		mon.TrackDruidTimeMetricInSeconds(mon.DruidAPIMethodDurationType, methodStartTime, errorCode, mon.GetDruidLookups)
		return nil, fmt.Errorf("Unable to unmarshal response from druid for request %s: %s", models.AsJSONString(resp), err.Error())
	}

	if logger.IsDebugEnabled() {
		logger.Log.Debugf("Response from druid for query %s ->  %s", mon.GetDruidLookups, models.AsJSONString(responseMap))
	}

	lookupMap := make(DruidLookupStatus)
	for _, val := range responseMap.Results {
		lookupMap[val] = false
	}

	mon.TrackDruidTimeMetricInSeconds(mon.DruidAPIMethodDurationType, methodStartTime, successCode, mon.GetDruidLookups)

	return lookupMap, nil
}

const (
	druidLookUpTierName  = "__default"
	druidLookUpConfig    = "lookups/config"
	druidCoordinatorPath = "/druid/coordinator/v1"
)

func (dc *DruidDatastoreClient) generateDruidCoordinatorURI(paths ...string) string {
	var lookupEndpoint string
	var path string

	for _, p := range paths {
		path = path + "/" + p
	}
	// @TODO: We should make this more explicit than -1. Maybe empty string is better? It's only used for debugging
	if dc.coordinatorPort == "-1" {
		//lookupEndpoint = dc.coordinatorServer + "/druid/coordinator/v1/lookups/config"
		lookupEndpoint = dc.coordinatorServer + druidCoordinatorPath + path
	} else {
		lookupEndpoint = dc.coordinatorServer + ":" + dc.coordinatorPort + druidCoordinatorPath + path
	}

	return lookupEndpoint
}
func (dc *DruidDatastoreClient) deleteItemToLookup(host string, lookupName string) error {
	startTime := time.Now()
	url := host + "/__default/" + lookupName
	if logger.IsDebugEnabled() {
		logger.Log.Debugf("Deleting lookup %s, url is %s", lookupName, url)
	}

	_, err := sendRequest("DELETE", dc.dClient.HttpClient, url, dc.AuthToken, nil)
	if err != nil {
		logger.Log.Errorf("Failed to Delete lookup %s because %s", lookupName, err.Error())

		mon.TrackDruidTimeMetricInSeconds(mon.DruidAPIMethodDurationType, startTime, errorCode, mon.DeleteDruidMetaLookups)
		return err
	}

	mon.TrackDruidTimeMetricInSeconds(mon.DruidAPIMethodDurationType, startTime, successCode, mon.DeleteDruidMetaLookups)
	return nil
}

func (dc *DruidDatastoreClient) addItemToLookup(host string, lookupName string, payload *lookup) error {
	startTime := time.Now()

	url := host //+ "/__default/" + lookupName
	// Domain lookups are assigned to the __default tier
	b, err := json.Marshal(map[string]map[string]*lookup{"__default": map[string]*lookup{lookupName: payload}})
	if err != nil {
		logger.Log.Error("Failed to marshal lookupRequest", err.Error())
		mon.TrackDruidTimeMetricInSeconds(mon.DruidAPIMethodDurationType, startTime, errorCode, mon.AddDruidMetaLookups)
		return err
	}

	// if logger.IsDebugEnabled() {
	// 	logger.Log.Debugf("Dumping url: %s", url)
	// 	logger.Log.Debugf("Sending lookup request %s, payload: %s", url, string(b))
	// }

	// Delete the look up first (don't worry about errors, best effort)
	//dc.deleteItemToLookup(host, lookupName)

	_, err = sendRequest("POST", dc.dClient.HttpClient, url, dc.AuthToken, b)
	if err != nil {
		logger.Log.Errorf("Failed to update lookup %s", err.Error())
		mon.TrackDruidTimeMetricInSeconds(mon.DruidAPIMethodDurationType, startTime, errorCode, mon.AddDruidMetaLookups)
		return err
	}
	mon.TrackDruidTimeMetricInSeconds(mon.DruidAPIMethodDurationType, startTime, successCode, mon.AddDruidMetaLookups)

	return nil
}

// Fetch existing lookup names and delete any existing lookups on the server that are nolonger valid.
// Use the lookup map created in the previous step to identify valid domains.
func (dc *DruidDatastoreClient) checkAndPostDefaultLookup(lookupEndpoint string) ([]string, error) {
	url := lookupEndpoint + "/__default"
	var lookupNames []string
	result, err := sendRequest("GET", dc.dClient.HttpClient, url, dc.AuthToken, nil)
	if err != nil {
		if strings.Contains(err.Error(), "No lookups found") {
			logger.Log.Infof("No lookups found.  Need to initialize lookups before any are created")
			result, err = sendRequest("POST", dc.dClient.HttpClient, lookupEndpoint, dc.AuthToken, []byte("{}"))
			if err != nil {
				logger.Log.Errorf("Failed to initialize druid lookups", err.Error())
				return nil, err
			}
			logger.Log.Infof("Lookups successfully initialized")
		} else {
			logger.Log.Errorf("Failed to fetch lookups", err.Error())
			return nil, err
		}
	} else {
		err = json.Unmarshal(result, &lookupNames)
		if err != nil {
			logger.Log.Errorf("Failed to unmarshal lookup:%s", result, err.Error())
			return nil, err
		}
	}

	return lookupNames, nil
}

// AddMonitoredObjectToLookup - Adds a monitored object to the druid look ups
func (dc *DruidDatastoreClient) AddMonitoredObjectToLookup(tenantID string, monitoredObjects []*tenant.MonitoredObject, datatype string) error {
	methodStartTime := time.Now()
	// version := time.Now().Format(time.RFC3339)
	lookupEndpoint := dc.generateDruidCoordinatorURI(druidLookUpConfig)
	logger.Log.Info("Druid Lookup URI:%s", lookupEndpoint)

	var lookups map[string]*lookup
	// Fetch existing lookup names and delete any existing lookups on the server that are nolonger valid.
	// Use the lookup map created in the previous step to identify valid domains.
	lookupNames, err := dc.checkAndPostDefaultLookup(lookupEndpoint)
	if err != nil {
		mon.TrackDruidTimeMetricInSeconds(mon.DruidAPIMethodDurationType, methodStartTime, errorCode, mon.UpdateDruidLookups)

		return err
	}

	if logger.IsDebugEnabled() {
		logger.Log.Infof("Dumping lookupNames: %+v", lookupNames)
	}

	if datatype == "meta" {
		lookups, err = dc.updateMetadataLookup(lookupEndpoint, tenantID, datatype, monitoredObjects)
		if err != nil {
			mon.TrackDruidTimeMetricInSeconds(mon.DruidAPIMethodDurationType, methodStartTime, errorCode, mon.UpdateDruidLookups)

			return err
		}
	}
	updateLookupCache(lookups)
	mon.TrackDruidTimeMetricInSeconds(mon.DruidAPIMethodDurationType, methodStartTime, successCode, mon.UpdateDruidLookups)

	return nil
}
