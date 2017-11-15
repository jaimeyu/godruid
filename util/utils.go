package util

import (
	"encoding/json"
	"fmt"

	"github.com/accedian/adh-gather/logger"
)

func convertGenericDataToObject(genericData map[string]interface{}, dataContainer interface{}, dataTypeStr string) error {
	genericDataInBytes, err := convertGenericObjectToBytes(genericData)
	if err != nil {
		return err
	}

	err = json.Unmarshal(genericDataInBytes, &dataContainer)
	if err != nil {
		return fmt.Errorf("Error converting generic data to %s type: %s", dataTypeStr, err.Error())
	}

	logger.Log.Debugf("Converted generic data to %s: %v\n", dataTypeStr, dataContainer)

	return nil
}

func convertGenericObjectToBytes(genericObject map[string]interface{}) ([]byte, error) {
	genericUserInBytes, err := json.Marshal(genericObject)
	if err != nil {
		return nil, fmt.Errorf("Error converting generic data to bytes: %s", err.Error())
	}

	return genericUserInBytes, nil
}
