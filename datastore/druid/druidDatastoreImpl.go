package druid

import (
	"fmt"

	"github.com/accedian/adh-gather/config"
	"github.com/accedian/adh-gather/gather"
	"github.com/accedian/godruid"
)

type DruidDatastoreClient struct {
	server  string
	cfg     config.Provider
	dClient godruid.Client
}

func (dc *DruidDatastoreClient) executeQuery(query godruid.Query) (string, error) {
	client := dc.dClient

	err := client.Query(query)

	fmt.Println(string(query.GetRawJSON()))

	if err != nil {
		fmt.Println("ERROR----", err)
	}

	return "nil", nil
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

func (dc *DruidDatastoreClient) GetStats(metric string, threshold string) (string, error) {

	query := StatsQuery("NPAVKPI", "delayP95", "", "2017-11-02/2100-01-01")
	dc.executeQuery(query)
	return "nil", nil
}

func (dc *DruidDatastoreClient) GetThresholdCrossing(metric string) (string, error) {
	query := ThresholdCrossingQuery("NPAVKPI", "delayP95", "", "2017-11-02/2100-01-01")
	dc.executeQuery(query)

	return "nil", nil
}
