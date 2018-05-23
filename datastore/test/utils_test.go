package test

import (
	"github.com/accedian/adh-gather/logger"
	//"log"
	ds "github.com/accedian/adh-gather/datastore"
	"testing"
)

func TestStripPrefixFromFullID(t *testing.T) {
	type qa struct {
		q string
		a string
	}
	mt := []qa{
		{q: "adminUser_2_someID", a: "someID"},
		{q: "hello_world", a: ""},
		{q: "hello", a: ""},
		{q: "monitoredObject_2_monitoredObject_id", a: "monitoredObject_id"},
		{q: "tenant_2_123456789", a: "123456789"},
	}

	t.Log("mt:", mt)

	for idx, e := range mt {
		output := ds.GetDataIDFromFullID(e.q)
		logger.Log.Debugf("[%d]Input:%s\tExpected:%s\tOutput:%s", idx, e.q, e.a, output)
		if output != e.a {
			logger.Log.Fatalf("Test failed, expected %s but got %s", e.a, output)
			t.Fail()
		}
	}
}
