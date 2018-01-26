package druid

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/accedian/adh-gather/config"
	"github.com/accedian/adh-gather/gather"
	"github.com/accedian/adh-gather/logger"
	"github.com/accedian/godruid"
	"github.com/golang/protobuf/jsonpb"
	"github.com/golang/protobuf/ptypes"

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
		Url:   server + ":" + port,
		Debug: true,
	}

	return &DruidDatastoreClient{
		cfg:       cfg,
		server:    server,
		dClient:   client,
		AuthToken: GetAuthCode(cfg),
	}
}

// peyo TODO: implement this query
func (dc *DruidDatastoreClient) GetHistogram(request *pb.HistogramRequest) (*pb.JSONAPIObject, error) {
	table := dc.cfg.GetString(gather.CK_druid_table.String())
	query, err := HistogramQuery(request.GetTenant(), table, request.Metric, request.Granularity, request.Direction, request.Interval, request.Resolution, request.GranularityBuckets)

	if err != nil {
		return nil, err
	}

	response, err := dc.executeQuery(query)

	if err != nil {
		return nil, err
	}

	histogram := []*pb.Histogram{}

	json.Unmarshal(response, &histogram)

	resp := &pb.HistogramResponse{
		Data: histogram,
	}

	data, err := ptypes.MarshalAny(resp)

	if err != nil {
		return nil, err
	}

	// peyo TODO: need to figure out where to get this ID and Type from.
	uuid := uuid.NewV4()
	rr := &pb.JSONAPIObject{
		Data: []*pb.Data{
			&pb.Data{
				Id:         uuid.String(),
				Type:       EventDistribution,
				Attributes: data,
			},
		},
	}
	return rr, nil
}

// GetThresholdCrossing - Executes a 'threshold crossing' query against druid. Wraps the
// result in a JSON API wrapper.
// peyo TODO: probably don't need to wrap JSON API here...should maybe do it elsewhere
func (dc *DruidDatastoreClient) GetThresholdCrossing(request *pb.ThresholdCrossingRequest, thresholdProfile *pb.TenantThresholdProfile) (*pb.JSONAPIObject, error) {

	table := dc.cfg.GetString(gather.CK_druid_table.String())

	query, err := ThresholdCrossingQuery(request.GetTenant(), table, request.Metric, request.Granularity, request.Interval, request.ObjectType, request.Direction, thresholdProfile.Data)

	if err != nil {
		return nil, err
	}

	response, err := dc.executeQuery(query)

	if err != nil {
		return nil, err
	}

	thresholdCrossing := []*pb.ThresholdCrossing{}

	json.Unmarshal(response, &thresholdCrossing)

	formattedJSON, err := reformatThresholdCrossingResponse(thresholdCrossing)

	if err != nil {
		return nil, err
	}

	resp := new(pb.ThresholdCrossingResponse)

	err = jsonpb.Unmarshal(bytes.NewReader(formattedJSON), resp)

	if err != nil {
		return nil, fmt.Errorf("Unable to unmarshal formatted JSON into ThresholdCrossingResponse. Err: %s", err)
	}

	data, err := ptypes.MarshalAny(resp)

	if err != nil {
		return nil, err
	}

	// peyo TODO: need to figure out where to get this ID and Type from.
	uuid := uuid.NewV4()
	rr := &pb.JSONAPIObject{
		Data: []*pb.Data{
			&pb.Data{
				Id:         uuid.String(),
				Type:       ThresholdCrossingReport,
				Attributes: data,
			},
		},
	}

	return rr, nil
}

// GetThresholdCrossing - Executes a 'threshold crossing' query against druid. Wraps the
// result in a JSON API wrapper.
// peyo TODO: probably don't need to wrap JSON API here...should maybe do it elsewhere
func (dc *DruidDatastoreClient) GetThresholdCrossingByMonitoredObject(request *pb.ThresholdCrossingRequest, thresholdProfile *pb.TenantThresholdProfile) (*pb.JSONAPIObject, error) {

	table := dc.cfg.GetString(gather.CK_druid_table.String())

	query, err := ThresholdCrossingByMonitoredObjectQuery(request.GetTenant(), table, request.Metric, request.Granularity, request.Interval, request.ObjectType, request.Direction, thresholdProfile.Data)

	if err != nil {
		return nil, err
	}

	response, err := dc.executeQuery(query)

	if err != nil {
		return nil, err
	}

	thresholdCrossing := make([]ThresholdCrossingByMonitoredObjectResponse, 0)

	err = json.Unmarshal(response, &thresholdCrossing)

	formattedJSON, err := reformatThresholdCrossingByMonitoredObjectResponse(thresholdCrossing)

	if err != nil {
		return nil, err
	}

	resp := new(pb.ThresholdCrossingByMonitoredObjectResponse)

	err = jsonpb.Unmarshal(bytes.NewReader(formattedJSON), resp)

	if err != nil {
		return nil, fmt.Errorf("Unable to unmarshal formatted JSON into ThresholdCrossingResponse. Err: %s", err)
	}

	data, err := ptypes.MarshalAny(resp)

	if err != nil {
		return nil, err
	}

	// peyo TODO: need to figure out where to get this ID and Type from.
	uuid := uuid.NewV4()
	rr := &pb.JSONAPIObject{
		Data: []*pb.Data{
			&pb.Data{
				Id:         uuid.String(),
				Type:       ThresholdCrossingReport,
				Attributes: data,
			},
		},
	}

	return rr, nil
}

func (dc *DruidDatastoreClient) GetRawMetrics(request *pb.RawMetricsRequest) (*pb.JSONAPIObject, error) {

	table := dc.cfg.GetString(gather.CK_druid_table.String())

	query, err := RawMetricsQuery(request.GetTenant(), table, request.GetMetric(), request.GetInterval(), request.GetObjectType(), request.GetDirection(), request.GetMonitoredObjectId())

	if err != nil {
		return nil, err
	}

	response, err := dc.executeQuery(query)

	//	fmt.Println(string(response))

	if err != nil {
		return nil, err
	}

	resp := make([]RawMetricsResponse, 0)

	json.Unmarshal(response, &resp)

	formattedJSON, err := reformatRawMetricsResponse(resp)

	if err != nil {
		return nil, err
	}

	rawMetricResp := new(pb.RawMetricsResponse)

	err = jsonpb.Unmarshal(bytes.NewReader(formattedJSON), rawMetricResp)

	if err != nil {
		return nil, fmt.Errorf("Unable to unmarshal formatted JSON into RawMetricsResponse. Err: %s", err)
	}

	data, err := ptypes.MarshalAny(rawMetricResp)

	if err != nil {
		return nil, err
	}

	uuid := uuid.NewV4()
	rr := &pb.JSONAPIObject{
		Data: []*pb.Data{
			&pb.Data{
				Id:         uuid.String(),
				Type:       RawMetrics,
				Attributes: data,
			},
		},
	}

	return rr, nil
}
