package druid

import (
	"fmt"

	"github.com/accedian/adh-gather/config"
	"github.com/accedian/adh-gather/gather"
	"github.com/evoleads/godruid"
)

type DruidDatastoreClient struct {
	server string
	cfg    config.Provider
}

func NewDruidDatasctoreClient() *DruidDatastoreClient {
	cfg := gather.GetConfig()

	return &DruidDatastoreClient{
		cfg:    cfg,
		server: "https://broker.proto.npav.accedian.net/druid/v2/?pretty=",
	}
}

func (dc *DruidDatastoreClient) GetNumberOfThesholdViolations(metric string, threshold string) (string, error) {

	query := &godruid.QueryTimeseries{
		QueryType:   "timeseries",
		DataSource:  "NPAVPM",
		Granularity: godruid.GranHour,
		Aggregations: []godruid.Aggregation{
			godruid.AggLongMin("delayP95Min", "delayP95"),
			godruid.AggLongSum("delayP95Sum", "delayP95"),
			godruid.AggLongMax("delayP95Max", "delayP95"),
			godruid.AggCount("rowCount"),
		},
		PostAggregations: []godruid.PostAggregation{
			godruid.PostAggregation{
				Type: "arithmetic",
				Name: "medianPoint",
				Fn:   "/",
				Fields: []godruid.PostAggregation{
					godruid.PostAggregation{
						Type:      "fieldAccess",
						FieldName: "rowCount",
					},
					godruid.PostAggregation{
						Type:      "constant",
						FieldName: "2",
					},
				},
			},
			godruid.PostAggregation{
				Type: "arithmetic",
				Name: "delayP95Mean",
				Fn:   "/",
				Fields: []godruid.PostAggregation{
					godruid.PostAggregation{
						Type:      "fieldAccess",
						FieldName: "delayp95Sum",
					},
					godruid.PostAggregation{
						Type:      "fieldAccess",
						FieldName: "rowcount",
					},
				},
			},
		},
		Intervals: []string{"2017-11-02/2100-01-01"},
	}
	client := godruid.Client{
		Url:   dc.server,
		Debug: true,
	}

	err := client.Query(query)

	if err != nil {
		fmt.Println("ERROR----", err)
	}

	fmt.Println("requst", client.LastRequest)

	fmt.Println("response", client.LastResponse)

	fmt.Printf("query.QueryResult:\n%v", query.QueryResult)

	return "nil", nil
}

func (dc *DruidDatastoreClient) GetStats(metric string) (string, error) {

	return "nil", nil
}
