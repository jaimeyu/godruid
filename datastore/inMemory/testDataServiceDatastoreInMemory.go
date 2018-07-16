package inMemory

import (
	"errors"
)

// TestDataServiceDatastoreInMemory - test data InMemory impl.
type TestDataServiceDatastoreInMemory struct {
}

// CreateTestDataServiceDAO - returns an in-memory implementation of the Admin Service
// datastore.
func CreateTestDataServiceDAO() (*TestDataServiceDatastoreInMemory, error) {
	res := new(TestDataServiceDatastoreInMemory)

	return res, nil
}

// GetAllDocsByDatatype - InMemory implementation of GetAllDocsByDatatype
func (testDB *TestDataServiceDatastoreInMemory) GetAllDocsByDatatype(dbName string, datatype string) ([]map[string]interface{}, error) {
	// Stub to implement
	return nil, errors.New("GetAllDocsByDatatype() not implemented for InMemory DB")
}

// InsertTenantViews - InMemory implementation of InsertTenantViews
func (testDB *TestDataServiceDatastoreInMemory) InsertTenantViews(dbName string) string {
	// Stub to implement
	return "InsertTenantViews() not implemented for InMemory DB"
}
