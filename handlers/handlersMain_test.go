package handlers_test

import (
	"flag"
	"os"
	"testing"

	"github.com/accedian/adh-gather/logger"
)

var (
	metricsIntegrationTests = flag.Bool("metrics", false, "Run metrics api integration tests")
)

func TestMain(m *testing.M) {

	flag.Parse()

	err := setupTestDatastore()
	if err != nil {
		logger.Log.Fatalf("Unable to setup datastore for Data Cleaning Profile tests: %s", err.Error())
	}

	code := m.Run()

	err = destroyTestDatastore()
	if err != nil {
		logger.Log.Errorf("Unable to remove test datastore for Data Cleaning Profile tests: %s", err.Error())
	}

	os.Exit(code)
}
