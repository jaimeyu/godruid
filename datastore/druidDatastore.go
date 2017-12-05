package datastore

import (
	pb "github.com/accedian/adh-gather/gathergrpc"
)

type DruidDatastore interface {

	// Returns the the number of times a given metric crossed the
	// minor,major,critical thresholds of a given threshold object
	GetThresholdCrossing(request *pb.ThresholdCrossingRequest, ingestionProfile *pb.TenantIngestionProfileResponse) (*pb.JSONAPIObject, error)

	// Returns the min,max,avg,median for a given metric
	GetHistogram(request *pb.HistogramRequest) (*pb.JSONAPIObject, error)
}
