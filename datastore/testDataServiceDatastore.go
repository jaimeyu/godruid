package datastore

// TestDataType - data type descriptors for objects stored by the datastore
type TestDataType string

const (
	// DomainSlaReportType - datatype string used to identify an Admin User in the datastore record
	DomainSlaReportType TestDataType = "domainSlaReportInstance"
)

const (
	// DomainSlaReportStr - common name of the Domain SLA Report data type for use in logs.
	DomainSlaReportStr = "Domain SLA Report"
)

// TestDataServiceDatastore - interface which provides the functionality
// of the TestData Service Datastore.
type TestDataServiceDatastore interface {
	GetAllDocsByDatatype(dbName string, datatype string) ([]map[string]interface{}, error)
}
