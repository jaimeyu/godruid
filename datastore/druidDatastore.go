package datastore

type DruidDatastore interface {

	// Returns the the number of times a given metric crossed the
	// minor,major,critical thresholds of a given threshold object
	GetNumberOfThesholdViolations(metric string, threshold string) (string, error)

	// Returns the min,max,avg,median for a given metric
	GetStats(metric string) (string, error)
}
