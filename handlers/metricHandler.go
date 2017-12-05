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
func (msh *MetricServiceHandler) GetThresholdCrossing(ctx context.Context, thresholdCrossingReq *pb.ThresholdCrossingRequest, ingestionProfile *pb.TenantIngestionProfileResponse) (*pb.JSONAPIObject, error) {

	result, err := msh.druidDB.GetThresholdCrossing(thresholdCrossingReq, ingestionProfile)

	if err != nil {
		return nil, fmt.Errorf("Unable to retrieve Threshold Crossing. %s:", err.Error())
	}

	return result, nil
}

// GetThresholdCrossing
func (msh *MetricServiceHandler) GetHistogram(ctx context.Context, histogramReq *pb.HistogramRequest) (*pb.JSONAPIObject, error) {

	result, err := msh.druidDB.GetHistogram(histogramReq)

	if err != nil {
		return nil, fmt.Errorf("Unable to retrieve Histogram. %s:", err.Error())
	}

	return result, nil
}
