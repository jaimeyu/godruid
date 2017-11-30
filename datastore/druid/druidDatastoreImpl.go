package druid

import (
	"encoding/json"
	"fmt"

	"github.com/accedian/adh-gather/config"
	"github.com/accedian/adh-gather/gather"
	"github.com/accedian/godruid"
	"github.com/golang/protobuf/ptypes"

	pb "github.com/accedian/adh-gather/gathergrpc"
)

var bearerId = "eyJhbGciOiJSUzI1NiIsImtpZCI6IjYxNGQwZWQ5M2QzOWZiZjFiYzE4NDc5M2RhMDgwMWQ0MGY0MGI4MjIifQ.eyJhenAiOiI2Mjc2MzQ4Nzc3NjEtanNndDU2YjA4ODc1YWhkZG43MmRtaXBmcnA4NDhvdTQuYXBwcy5nb29nbGV1c2VyY29udGVudC5jb20iLCJhdWQiOiI2Mjc2MzQ4Nzc3NjEtZnVpdWhtbDI5Y2U3OTg1dWE1cmNqbTJzM2Fkazc5N3YuYXBwcy5nb29nbGV1c2VyY29udGVudC5jb20iLCJzdWIiOiIxMTQyOTMwODk5Mjc0MDUzNjM0NDUiLCJoZCI6ImFjY2VkaWFuLmNvbSIsImVtYWlsIjoicHR6b2xvdkBhY2NlZGlhbi5jb20iLCJlbWFpbF92ZXJpZmllZCI6dHJ1ZSwiYXRfaGFzaCI6ImJ5cXhEUjBBSFAyc0EwZW5YWTRxalEiLCJpc3MiOiJodHRwczovL2FjY291bnRzLmdvb2dsZS5jb20iLCJpYXQiOjE1MTIwNTkxOTYsImV4cCI6MTUxMjA2Mjc5Nn0.JiAIz2W25p9dkOVbhc5b01wEOp23QLErBskUMLFarb888l7ga1VHg37fJ_KiDDN9qVnwXrtN3LURy83ap73PCH-1y2nqr_hK10sMsLoKUNXtKvV5JRaVdWGIOL_i_HiaLzK6pwjfuMw4dECCuVdGc9EHB8Evz3oJh9Y1ELcr1hYLa2n4jkN9-QVjxH6lM5ATUTOkJ0fMJ5yggesRJ0bN1iYCMc_ibrwa05MfS8MqzVd-UhTMG8NDv5OiPu2WQMCdygSVXt4Z47mrvVNYe3vEK-RYsHXFZTuiJPpgORmA4Gp1Lc9e2hi02YZd5lERgZf18gr1y_rcck3romjgOWnoEA"

type DruidDatastoreClient struct {
	server  string
	cfg     config.Provider
	dClient godruid.Client
}

func (dc *DruidDatastoreClient) executeQuery(query godruid.Query) ([]byte, error) {
	client := dc.dClient

	err := client.Query(query)

	if err != nil {
		fmt.Println("ERROR---->", err)
	}

	return query.GetRawJSON(), nil
}

func NewDruidDatasctoreClient() *DruidDatastoreClient {
	cfg := gather.GetConfig()
	server := cfg.GetString(gather.CK_druid_server.String())
	client := godruid.Client{
		Url:       server,
		Debug:     true,
		AuthToken: bearerId,
	}

	return &DruidDatastoreClient{
		cfg:     cfg,
		server:  server,
		dClient: client,
	}
}

func (dc *DruidDatastoreClient) GetStats(metric string) (string, error) {
	table := dc.cfg.GetString(gather.CK_druid_table.String())
	query := StatsQuery(table, metric, "", "2017-11-02/2100-01-01")
	dc.executeQuery(query)
	return "nil", nil
}

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
									Bound:    "upper",
									Unit:     "percent",
									Severity: "minor",
									Value:    500,
								},
								&pb.Event{
									Bound:    "upper",
									Unit:     "percent",
									Severity: "major",
									Value:    1000,
								},
								&pb.Event{
									Bound:    "upper",
									Unit:     "percent",
									Severity: "critical",
									Value:    1500,
								},
							},
						},
					},
				},
			},
		},
	}
}

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
