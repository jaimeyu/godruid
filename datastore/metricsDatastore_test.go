package datastore_test

import (
	"testing"

	"github.com/accedian/adh-gather/datastore"
	"github.com/stretchr/testify/assert"
)

func TestKeySpecSorting(t *testing.T) {

	ordered := map[string]interface{}{"vendor": "Accedian", "monitoredObjectTypes": []string{"test1", "test2", "test3"}}
	unordered := map[string]interface{}{"vendor": "Accedian", "monitoredObjectTypes": []string{"test3", "test1", "test2"}}

	qks := datastore.QueryKeySpec{}

	// Ensure that internally we sort the entries for a key and put out the same id
	assert.Equal(t, qks.AddKeySpec(ordered), qks.AddKeySpec(unordered))
	for i := 0; i < len(ordered); i++ {
		assert.Equal(t, ordered["monitoredObjectTypes"].([]string)[i], unordered["monitoredObjectTypes"].([]string)[i])
	}
}
