package druid

import (
	"encoding/json"
	"fmt"

	pb "github.com/accedian/adh-gather/gathergrpc"
	"github.com/accedian/godruid"
)

func getThreshold(thresholdProfile *pb.TenantThresholdProfile, objectType string) (*pb.TenantThreshold, error) {

	for _, tp := range thresholdProfile.Thresholds {
		if tp.ObjectType == objectType {
			return tp, nil
		}
	}

	return nil, fmt.Errorf("No threshold profile available for object type: %s", objectType)

}

func getMetric(threshold *pb.TenantThreshold, metricName string, objectType string) (*pb.TenantMetric, error) {

	for _, m := range threshold.Metrics {
		if m.Id == metricName {
			return m, nil
		}

	}
	return nil, fmt.Errorf("No threshold information available for object type: %s, and metric: %s", objectType, metricName)
}

func getEvents(metric *pb.TenantMetric, direction string, objectType string) ([]*pb.TenantEvent, error) {
	for _, md := range metric.Data {
		if md.Direction == direction {
			return md.Events, nil
		}
	}
	return nil, fmt.Errorf("No events information available for object type: %s, metric: %s, and direction: %s", objectType, metric.Id, direction)
}

func queryToString(query godruid.Query, debug bool) string {
	var reqJson []byte
	var err error

	if debug {
		reqJson, err = json.MarshalIndent(query, "", "  ")
	} else {
		reqJson, err = json.Marshal(query)
	}

	if err != nil {
		return ""
	}

	return string(reqJson)
}
