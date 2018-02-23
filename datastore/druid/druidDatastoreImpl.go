package druid

import (
	"crypto/tls"
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/accedian/adh-gather/config"
	"github.com/accedian/adh-gather/gather"
	"github.com/accedian/adh-gather/logger"
	"github.com/accedian/godruid"

	db "github.com/accedian/adh-gather/datastore"
	pb "github.com/accedian/adh-gather/gathergrpc"

	"github.com/satori/go.uuid"
)

const (
	ThresholdCrossingReport = "threshold-crossing-report"
	EventDistribution       = "event-distribution"
	RawMetrics              = "raw-metrics"
)

// DruidDatastoreClient - struct responsible for handling
// database operations for druid
type DruidDatastoreClient struct {
	server     string
	cfg        config.Provider
	dClient    godruid.Client
	AuthToken  string
	numRetries int
}

type ThresholdCrossingByMonitoredObjectResponse struct {
	Version   string
	Timestamp string
	Event     map[string]interface{}
}

type RawMetricsEvents struct {
	Event map[string]interface{}
}

type RawMetricsResult struct {
	Events []RawMetricsEvents
}

type RawMetricsResponse struct {
	Timestamp string
	Result    RawMetricsResult
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
			err := client.Query(query, dc.AuthToken)
			if err != nil {
				return nil, err
			}
			return query.GetRawJSON(), nil
		}
		return nil, err
	}

	return query.GetRawJSON(), nil
}

// NewDruidDatasctoreClient - Constructor for DruidDatastoreClient object
// initializes the godruid client, and retrieves auth token
// peyo TODO: the auth functionality here needs to be changed, this is only valid for dev
func NewDruidDatasctoreClient() *DruidDatastoreClient {
	cfg := gather.GetConfig()
	server := cfg.GetString(gather.CK_druid_server.String())
	port := cfg.GetString(gather.CK_druid_port.String())
	client := godruid.Client{
		Url:        server + ":" + port,
		Debug:      true,
		HttpClient: makeHttpClient(),
	}

	return &DruidDatastoreClient{
		cfg:       cfg,
		server:    server,
		dClient:   client,
		AuthToken: GetAuthCode(cfg),
	}
}

// peyo TODO: implement this query
func (dc *DruidDatastoreClient) GetHistogram(request *pb.HistogramRequest) (map[string]interface{}, error) {

	logger.Log.Debugf("Calling GetHistogram for request: %v", logger.AsJSONString(request))
	table := dc.cfg.GetString(gather.CK_druid_table.String())

	// peyo TODO we should have a better way to handle default query params
	timeout := request.GetTimeout()
	if timeout == 0 {
		timeout = 5000
	}

	query, err := HistogramQuery(request.GetTenant(), table, request.Metric, request.Granularity, request.Direction, request.Interval, request.Resolution, request.GranularityBuckets, request.GetVendor(), timeout)

	if err != nil {
		return nil, err
	}

	logger.Log.Debugf("Querying Druid for %s with query: %v", db.HistogramStr, logger.AsJSONString(query))
	response, err := dc.executeQuery(query)

	if err != nil {
		return nil, err
	}

	histogram := []*pb.Histogram{}

	err = json.Unmarshal(response, &histogram)
	if err != nil {
		return nil, err
	}

	logger.Log.Debugf("Response from druid for %s: %v", db.HistogramStr, logger.AsJSONString(histogram))

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

	logger.Log.Debugf("Calling GetThresholdCrossing for request: %v", logger.AsJSONString(request))
	table := dc.cfg.GetString(gather.CK_druid_table.String())

	// peyo TODO we should have a better way to handle default query params
	timeout := request.GetTimeout()
	if timeout == 0 {
		timeout = 5000
	}

	query, err := ThresholdCrossingQuery(request.GetTenant(), table, request.Metric, request.Granularity, request.Interval, request.ObjectType, request.Direction, thresholdProfile.Data, request.GetVendor(), timeout)

	if err != nil {
		return nil, err
	}

	logger.Log.Debugf("Querying Druid for %s with query: %v", db.ThresholdCrossingStr, logger.AsJSONString(query))
	response, err := dc.executeQuery(query)

	if err != nil {
		return nil, err
	}

	thresholdCrossing := []*pb.ThresholdCrossing{}
	err = json.Unmarshal(response, &thresholdCrossing)
	if err != nil {
		return nil, err
	}
	logger.Log.Debugf("Response from druid for %s: %v", db.ThresholdCrossingStr, logger.AsJSONString(thresholdCrossing))

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

// GetThresholdCrossing - Executes a 'threshold crossing' query against druid. Wraps the
// result in a JSON API wrapper.
// peyo TODO: probably don't need to wrap JSON API here...should maybe do it elsewhere
func (dc *DruidDatastoreClient) GetThresholdCrossingByMonitoredObject(request *pb.ThresholdCrossingRequest, thresholdProfile *pb.TenantThresholdProfile) (map[string]interface{}, error) {

	logger.Log.Debugf("Calling GetThresholdCrossingByMonitoredObject for request: %v", logger.AsJSONString(request))
	table := dc.cfg.GetString(gather.CK_druid_table.String())

	// peyo TODO we should have a better way to handle default query params
	timeout := request.GetTimeout()
	if timeout == 0 {
		timeout = 5000
	}

	query, err := ThresholdCrossingByMonitoredObjectQuery(request.GetTenant(), table, request.Metric, request.Granularity, request.Interval, request.ObjectType, request.Direction, thresholdProfile.Data, request.GetVendor(), timeout)
	if err != nil {
		return nil, err
	}

	logger.Log.Debugf("Querying Druid for %s with query: %v", db.ThresholdCrossingByMonitoredObjectStr, logger.AsJSONString(query))
	response, err := dc.executeQuery(query)
	if err != nil {
		return nil, err
	}

	thresholdCrossing := make([]ThresholdCrossingByMonitoredObjectResponse, 0)
	err = json.Unmarshal(response, &thresholdCrossing)
	if err != nil {
		return nil, err
	}

	logger.Log.Debugf("Response from druid for %s: %v", db.ThresholdCrossingByMonitoredObjectStr, logger.AsJSONString(thresholdCrossing))

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

func (dc *DruidDatastoreClient) GetRawMetrics(request *pb.RawMetricsRequest) (map[string]interface{}, error) {

	logger.Log.Debugf("Calling GetRawMetrics for request: %v", logger.AsJSONString(request))

	table := dc.cfg.GetString(gather.CK_druid_table.String())

	// peyo TODO we should have a better way to handle default query params
	timeout := request.GetTimeout()
	if timeout == 0 {
		timeout = 30000
	}

	query, err := RawMetricsQuery(request.GetTenant(), table, request.GetMetric(), request.GetInterval(), request.GetObjectType(), request.GetDirection(), request.GetMonitoredObjectId(), timeout)

	if err != nil {
		return nil, err
	}

	logger.Log.Debugf("Querying Druid for %s with query: %v", db.RawMetricStr, logger.AsJSONString(query))
	response, err := dc.executeQuery(query)

	//	fmt.Println(string(response))

	if err != nil {
		return nil, err
	}

	resp := make([]RawMetricsResponse, 0)

	err = json.Unmarshal(response, &resp)
	if err != nil {
		return nil, err
	}

	logger.Log.Debugf("Response from druid for %s: %v", db.RawMetricStr, logger.AsJSONString(resp))

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
