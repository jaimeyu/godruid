package gathergrpc

import (
	"encoding/json"

	admmod "github.com/accedian/adh-gather/models/admin"
	"github.com/accedian/adh-gather/models/common"
	tenmod "github.com/accedian/adh-gather/models/tenant"
)

var (
	pbStateEnumToStringMap = map[float64]string{
		2: string(common.UserActive),
		1: string(common.UserInvited),
		4: string(common.UserPendingDelete),
		3: string(common.UserSuspended),
		0: string(common.UserUnknown),
	}
	userStateToPBUserStateMap = map[string]float64{
		string(common.UserActive):        2,
		string(common.UserInvited):       1,
		string(common.UserPendingDelete): 4,
		string(common.UserSuspended):     3,
		string(common.UserUnknown):       0,
	}
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

	// handle state mapping
	switch pbObject.(type) {
	case *TenantUser, *AdminUser, *TenantDescriptor:
		// Convert the state enum to a string value
		if flattenedObject["state"] != nil {
			value := flattenedObject["state"].(float64)
			flattenedObject["state"] = pbStateEnumToStringMap[value]
		}
	}

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

	// Handle the state conversion
	switch initialObj.(type) {
	case *tenmod.User, *admmod.User, *admmod.Tenant:
		// Convert the state enum to a string value
		if genericObject["state"] != nil {
			value := genericObject["state"].(string)
			genericObject["state"] = userStateToPBUserStateMap[value]
		}
	}

	resultBytes, err := json.Marshal(expanded)
	if err != nil {
		return err
	}

	// Unmarshal into the desired type
	return json.Unmarshal(resultBytes, pbObjectContainer)
}
