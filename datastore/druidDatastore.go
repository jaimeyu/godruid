package datastore

import (
	pb "github.com/accedian/adh-gather/gathergrpc"
)

const (
	// ThresholdCrossingStr - common name of the ThresholdCrossingStr data type for use in logs.
	ThresholdCrossingStr = "Threshold Crossing"

	// ThresholdCrossingByMonitoredObjectStr - common name for use in logs.
	ThresholdCrossingByMonitoredObjectStr = "Threshold Crossing by Monitored Object"

	// HistogramStr - common name for use in logs.
	HistogramStr = "Histogram"

	// RawMetricString - common name for use in logs.
	RawMetricStr = "Raw Metric"
)

type DruidDatastore interface {

	// Returns the the number of times a given metric crossed the
	// minor,major,critical thresholds of a given threshold object
	GetThresholdCrossing(request *pb.ThresholdCrossingRequest, thresholdProfile *pb.TenantThresholdProfile) (*pb.JSONAPIObject, error)

	// Returns the the number of times a given metric crossed the
	// minor,major,critical thresholds of a given threshold object
	GetThresholdCrossingByMonitoredObject(request *pb.ThresholdCrossingRequest, thresholdProfile *pb.TenantThresholdProfile) (*pb.JSONAPIObject, error)

	// Returns the min,max,avg,median for a given metric
	GetHistogram(request *pb.HistogramRequest) (*pb.JSONAPIObject, error)

	// Returns raw metrics from druid
	GetRawMetrics(request *pb.RawMetricsRequest) (*pb.JSONAPIObject, error)
}
