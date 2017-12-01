package druid

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/accedian/adh-gather/config"
	"github.com/accedian/adh-gather/gather"
	"github.com/accedian/adh-gather/logger"
	"github.com/accedian/godruid"
	"github.com/golang/protobuf/ptypes"

	pb "github.com/accedian/adh-gather/gathergrpc"
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
			dc.numRetries++
			if dc.numRetries > 3 {
				return nil, fmt.Errorf("Unable to refresh valid auth token. Please contact administrator")
			}
			return dc.executeQuery(query)
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
func (dc *DruidDatastoreClient) GetStats(metric string) (string, error) {
	table := dc.cfg.GetString(gather.CK_druid_table.String())
	query := StatsQuery(table, metric, "", "2017-11-02/2100-01-01")
	dc.executeQuery(query)
	return "nil", nil
}

// peyo TODO : this should query couchDB and get the actual threshold profile
// will do this once we have threshold profiles in couch
func lookupThresholdProfile() *pb.ThresholdProfile {

	return &pb.ThresholdProfile{
		Twamp: &pb.Threshold{
			Metrics: []*pb.Metric{
				&pb.Metric{
					Id: "delayP95",
					Data: []*pb.MetricData{
						&pb.MetricData{
							Direction: "az",
							Events: []*pb.Event{
								&pb.Event{
									UpperBound:  30000,
									UpperStrict: true,
									Unit:        "percent",
									Type:        "minor",
								},
								&pb.Event{
									UpperBound:  50000,
									LowerBound:  30000,
									UpperStrict: true,
									LowerStrict: false,
									Unit:        "percent",
									Type:        "major",
								},
								&pb.Event{
									LowerBound:  50000,
									LowerStrict: false,
									Unit:        "percent",
									Type:        "critical",
								},
							},
						},
					},
				},
			},
		},
	}
}

// GetThresholdCrossing - Executes a 'threshold crossing' query against druid. Wraps the
// result in a JSON API wrapper.
// peyo TODO: probably don't need to wrap JSON API here...should maybe do it elsewhere
func (dc *DruidDatastoreClient) GetThresholdCrossing(request *pb.ThresholdCrossingRequest) (*pb.JSONAPIObject, error) {
	table := dc.cfg.GetString(gather.CK_druid_table.String())

	thresholdProfile := lookupThresholdProfile()
	threshold, err := getThreshold(thresholdProfile, "twamp")
	if err != nil {
		return nil, err
	}
	metric, err := getMetric(threshold, request.Metric, "twamp")
	if err != nil {
		return nil, err
	}
	events, err := getEvents(metric, request.Direction, "twamp")
	if err != nil {
		return nil, err
	}

	query := ThresholdCrossingQuery(table, request.Metric, request.Granularity, request.Interval, request.ObjectType, request.Direction, events)

	response, err := dc.executeQuery(query)

	if err != nil {
		return nil, err
	}

	tt := []*pb.ThresholdCrossing{}

	json.Unmarshal(response, &tt)

	resp := &pb.ThresholdCrossingResponse{
		Data: tt,
	}

	data, err := ptypes.MarshalAny(resp)

	if err != nil {
		return nil, err
	}

	// peyo TODO: need to figure out where to get this ID and Type from.
	rr := &pb.JSONAPIObject{
		Data: []*pb.Data{
			&pb.Data{
				Id:         "some-uuid",
				Type:       "report",
				Attributes: data,
			},
		},
	}

	return rr, nil
}
