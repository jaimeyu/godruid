package handlers_test

import (
	"os"
	"testing"

	"github.com/accedian/adh-gather/logger"
)

func TestMain(m *testing.M) {

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
