package handlers

import (
	"context"
	"fmt"

	db "github.com/accedian/adh-gather/datastore"
	"github.com/accedian/adh-gather/datastore/druid"
	pb "github.com/accedian/adh-gather/gathergrpc"
)

type MetricServiceHandler struct {
	druidDB db.DruidDatastore
}

func CreateMetricServiceHandler() *MetricServiceHandler {
	result := new(MetricServiceHandler)

	db := druid.NewDruidDatasctoreClient()

	result.druidDB = db

	return result
}

// GetThresholdCrossing
func (msh *MetricServiceHandler) GetThresholdCrossing(ctx context.Context, thresholdCrossingReq *pb.ThresholdCrossingRequest, thresholdProfile *pb.TenantThresholdProfile) (*pb.JSONAPIObject, error) {

	result, err := msh.druidDB.GetThresholdCrossing(thresholdCrossingReq, thresholdProfile)

	if err != nil {
		return nil, fmt.Errorf("Unable to retrieve Threshold Crossing. %s:", err.Error())
	}

	return result, nil
}

// GetThresholdCrossingByMonitoredObject
func (msh *MetricServiceHandler) GetThresholdCrossingByMonitoredObject(ctx context.Context, thresholdCrossingReq *pb.ThresholdCrossingRequest, thresholdProfile *pb.TenantThresholdProfile) (*pb.JSONAPIObject, error) {

	result, err := msh.druidDB.GetThresholdCrossingByMonitoredObject(thresholdCrossingReq, thresholdProfile)

	if err != nil {
		return nil, fmt.Errorf("Unable to retrieve Threshold Crossing. %s:", err.Error())
	}

	return result, nil
}

// GetThresholdHistogram
func (msh *MetricServiceHandler) GetHistogram(ctx context.Context, histogramReq *pb.HistogramRequest) (*pb.JSONAPIObject, error) {

	result, err := msh.druidDB.GetHistogram(histogramReq)

	if err != nil {
		return nil, fmt.Errorf("Unable to retrieve Histogram. %s:", err.Error())
	}

	return result, nil
}
