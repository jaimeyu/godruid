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

// DruidDatastoreClient - struct responsible for handling
// database operations for druid
type DruidDatastoreClient struct {
	server     string
	cfg        config.Provider
	dClient    godruid.Client
	AuthToken  string
	numRetries int
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
	client := godruid.Client{
		Url:   server,
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
	query := HistogramQuery(request.GetTenant(), table, request.Metric, request.Granularity, request.Interval, request.Resolution, request.GranularityBuckets)

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
	rr := &pb.JSONAPIObject{
		Data: []*pb.Data{
			&pb.Data{
				Id:         uuid.NewV4().String(),
				Type:       "event-distribution",
				Attributes: data,
			},
		},
	}
	return rr, nil
}

// GetThresholdCrossing - Executes a 'threshold crossing' query against druid. Wraps the
// result in a JSON API wrapper.
// peyo TODO: probably don't need to wrap JSON API here...should maybe do it elsewhere
func (dc *DruidDatastoreClient) GetThresholdCrossing(request *pb.ThresholdCrossingRequest, thresholdProfile *pb.TenantThresholdProfileResponse) (*pb.JSONAPIObject, error) {

	table := dc.cfg.GetString(gather.CK_druid_table.String())

	query := ThresholdCrossingQuery(request.GetTenant(), table, request.Metric, request.Granularity, request.Interval, request.ObjectType, request.Direction, thresholdProfile.Data)

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
	rr := &pb.JSONAPIObject{
		Data: []*pb.Data{
			&pb.Data{
				Id:         uuid.NewV4().String(),
				Type:       "threshold-crossing-report",
				Attributes: data,
			},
		},
	}

	return rr, nil
}
