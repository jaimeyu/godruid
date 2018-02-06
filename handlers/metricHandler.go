package handlers

import (
	"context"
	"fmt"

	db "github.com/accedian/adh-gather/datastore"
	"github.com/accedian/adh-gather/datastore/druid"
	"github.com/accedian/adh-gather/logger"
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

	logger.Log.Infof("Retrieving %s for: %v", db.ThresholdCrossingStr, thresholdCrossingReq) 
	result, err := msh.druidDB.GetThresholdCrossing(thresholdCrossingReq, thresholdProfile)

	if err != nil {
		return nil, fmt.Errorf("Unable to retrieve Threshold Crossing. %s:", err.Error())
	}

	logger.Log.Infof("Completed %s fetch for: %v", db.ThresholdCrossingStr, thresholdCrossingReq) 

	return result, nil
}

// GetThresholdCrossingByMonitoredObject
func (msh *MetricServiceHandler) GetThresholdCrossingByMonitoredObject(ctx context.Context, thresholdCrossingReq *pb.ThresholdCrossingRequest, thresholdProfile *pb.TenantThresholdProfile) (*pb.JSONAPIObject, error) {

	logger.Log.Infof("Retrieving %s for: %v", db.ThresholdCrossingByMonitoredObjectStr, thresholdCrossingReq) 
	result, err := msh.druidDB.GetThresholdCrossingByMonitoredObject(thresholdCrossingReq, thresholdProfile)

	if err != nil {
		return nil, fmt.Errorf("Unable to retrieve Threshold Crossing. %s:", err.Error())
	}

	logger.Log.Infof("Completed %s fetch for: %v", db.ThresholdCrossingByMonitoredObjectStr, thresholdCrossingReq) 

	return result, nil
}

// GetThresholdHistogram
func (msh *MetricServiceHandler) GetHistogram(ctx context.Context, histogramReq *pb.HistogramRequest) (*pb.JSONAPIObject, error) {

	logger.Log.Infof("Retrieving %s for: %v", db.HistogramStr, histogramReq) 
	result, err := msh.druidDB.GetHistogram(histogramReq)

	if err != nil {
		return nil, fmt.Errorf("Unable to retrieve Histogram. %s:", err.Error())
	}

	logger.Log.Infof("Completed %s fetch for: %v", db.HistogramStr, histogramReq) 

	return result, nil
}

// GetRawMetrics
func (msh *MetricServiceHandler) GetRawMetrics(ctx context.Context, rawMetricReq *pb.RawMetricsRequest) (*pb.JSONAPIObject, error) {

	logger.Log.Infof("Retrieving %s for: %v", db.RawMetricStr, rawMetricReq) 
	result, err := msh.druidDB.GetRawMetrics(rawMetricReq)

	if err != nil {
		return nil, fmt.Errorf("Unable to retrieve RawMetrics. %s:", err.Error())
	}

	logger.Log.Infof("Completed %s fetch for: %v", db.RawMetricStr, rawMetricReq) 

	return result, nil
}
