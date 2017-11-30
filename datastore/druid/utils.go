package druid

import (
	"fmt"

	pb "github.com/accedian/adh-gather/gathergrpc"
)

func getThreshold(thresholdProfile *pb.ThresholdProfile, objectType string) (*pb.Threshold, error) {

	switch objectType {
	case "twamp":
		return thresholdProfile.Twamp, nil
	default:
		return nil, fmt.Errorf("No threshold profile available for object type: %s", objectType)
	}

}

func getMetric(threshold *pb.Threshold, metricName string, objectType string) (*pb.Metric, error) {

	for _, m := range threshold.Metrics {
		if m.Id == metricName {
			return m, nil
		}

	}
	return nil, fmt.Errorf("No threshold information available for object type: %s, and metric: %s", objectType, metricName)
}

func getEvents(metric *pb.Metric, direction string, objectType string) ([]*pb.Event, error) {
	for _, md := range metric.Data {
		if md.Direction == direction {
			return md.Events, nil
		}
	}
	return nil, fmt.Errorf("No events information available for object type: %s, metric: %s, and direction: %s", objectType, metric.Id, direction)
}
