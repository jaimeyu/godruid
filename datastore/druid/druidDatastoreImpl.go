package druid

import (
	"encoding/json"
	"fmt"

	"github.com/accedian/adh-gather/config"
	"github.com/accedian/adh-gather/gather"
	"github.com/accedian/adh-gather/logger"
	"github.com/accedian/godruid"
	"github.com/golang/protobuf/ptypes"

	pb "github.com/accedian/adh-gather/gathergrpc"
)

type DruidDatastoreClient struct {
	server  string
	cfg     config.Provider
	dClient godruid.Client
}

func (dc *DruidDatastoreClient) executeQuery(query godruid.Query) ([]byte, error) {
	client := dc.dClient

	err := client.Query(query)

	//	fmt.Println(string(query.GetRawJSON()))

	if err != nil {
		fmt.Println("ERROR----", err)
	}

	return query.GetRawJSON(), nil
}

func NewDruidDatasctoreClient() *DruidDatastoreClient {
	cfg := gather.GetConfig()
	server := "https://broker.proto.npav.accedian.net"
	client := godruid.Client{
		Url:   server,
		Debug: true,
	}

	return &DruidDatastoreClient{
		cfg:     cfg,
		server:  server,
		dClient: client,
	}
}

func (dc *DruidDatastoreClient) GetStats(metric string) (string, error) {

	query := StatsQuery("NPAVKPI2", "delayP95", "", "2017-11-02/2100-01-01")
	dc.executeQuery(query)
	return "nil", nil
}

const googleApis = "type.googleapis.com/"

func (dc *DruidDatastoreClient) GetThresholdCrossing(metric string, threshold string) (*pb.JSONAPIObject, error) {

	query := ThresholdCrossingQuery("NPAVKPI2", "delayP95", "", "2017-11-02/2100-01-01")
	response, _ := dc.executeQuery(query)

	tt := []*pb.ThresholdCrossing{}

	json.Unmarshal(response, &tt)

	fmt.Println("RESPONSE, -->", tt)

	resp := &pb.ThresholdCrossingResponse{
		Data: tt,
	}

	data, err := ptypes.MarshalAny(resp)

	if err != nil {
		logger.Log.Errorf("Unable to marshal data. Err: %s", err)
	}

	rr := &pb.JSONAPIObject{
		Data: data,
	}

	return rr, nil
}
