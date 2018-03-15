package gathergrpc

import (
	"encoding/json"
)

func ConvertFromPBObject(pbObject interface{}, dataContainer interface{}) error {
	// Convert the pb object to generic object
	objBytes, err := json.Marshal(pbObject)
	if err != nil {
		return err
	}

	genericObject := map[string]interface{}{}
	err = json.Unmarshal(objBytes, &genericObject)
	if err != nil {
		return err
	}

	// Flatten the object
	flattenedObject := genericObject["data"].(map[string]interface{})
	flattenedObject["_id"] = genericObject["_id"]
	flattenedObject["_rev"] = genericObject["_rev"]

	resultBytes, err := json.Marshal(flattenedObject)
	if err != nil {
		return err
	}

	// Unmarshal into the desired type
	return json.Unmarshal(resultBytes, dataContainer)
}

func ConvertToPBObject(initialObj interface{}, pbObjectContainer interface{}) error {
	// Convert the pb object to generic object
	objBytes, err := json.Marshal(initialObj)
	if err != nil {
		return err
	}

	genericObject := map[string]interface{}{}
	err = json.Unmarshal(objBytes, &genericObject)
	if err != nil {
		return err
	}

	// Expand the object
	expanded := map[string]interface{}{}
	expanded["data"] = genericObject
	expanded["_id"] = genericObject["_id"]
	expanded["_rev"] = genericObject["_rev"]

	resultBytes, err := json.Marshal(expanded)
	if err != nil {
		return err
	}

	// Unmarshal into the desired type
	return json.Unmarshal(resultBytes, pbObjectContainer)
}
