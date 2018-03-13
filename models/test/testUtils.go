package models

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/accedian/adh-gather/logger"
	"github.com/manyminds/api2go/jsonapi"
	"github.com/stretchr/testify/assert"
)

// AttributeKeyPair - pair used to keep track of how a key for a property in a real object
// is mapped to a key in the attributes section of the same object serialized in jsonapi format.
type AttributeKeyPair struct {
	OriginalKey string
	MappedKey   string
}

// RunSerializationTest - a basic Marshal/Unmarshal test for any generic object which is
// structiured for jsonapi format. Paramaters are:
// t - the testing interface context
// original - the populated object you want to test serialization with.
// unmarshalled - an empty object you wish to use for unmarshalling the object back into.
// id - the id of the original object
// attrNames - list of string names of attributes to compare between the original and generated objects
func RunSerializationTest(t *testing.T, original interface{}, unmarshalled interface{}, id string, attrNames []string) {
	out, err := jsonapi.Marshal(original)
	assert.Nil(t, err)

	var jsonData map[string]interface{}
	err = json.Unmarshal(out, &jsonData)
	assert.Nil(t, err)

	logger.Log.Debugf("Original object after conversion to generic: %v", jsonData)

	// Validate the JSON version:
	validateObject(t, original, jsonData["data"].(map[string]interface{}), id, attrNames)

	// Now unmarshal directly into the object:
	err = jsonapi.Unmarshal(out, unmarshalled)
	assert.Nil(t, err)
	logger.Log.Debugf("Unmarshalled object from jsonapi content: %v", unmarshalled)
	assert.Equal(t, original, unmarshalled)
}

func validateObject(t *testing.T, original interface{}, generated map[string]interface{}, id string, attrNames []string) {
	assert.NotNil(t, original)
	assert.NotNil(t, generated)

	originalBytes, err := json.Marshal(original)
	assert.Nil(t, err)
	assert.NotNil(t, originalBytes)
	genericOriginal := map[string]interface{}{}
	err = json.Unmarshal(originalBytes, &genericOriginal)
	assert.Nil(t, err)

	assert.Equal(t, id, generated["id"], "The known id is not the same as the id in the generated object")
	assert.Equal(t, jsonapi.Pluralize(genericOriginal["datatype"].(string)), generated["type"], "The type of the generated object does not match the type in the generated object")

	attributes, ok := generated["attributes"].(map[string]interface{})
	assert.True(t, ok)
	assert.NotNil(t, attributes)

	for _, key := range attrNames {
		assert.NotNil(t, genericOriginal[key], fmt.Sprintf("Expect orignal to have property %s", key))
		assert.NotNil(t, attributes[key], fmt.Sprintf("Expect generated object to have property %s", key))
		assert.Equal(t, genericOriginal[key], attributes[key], fmt.Sprintf("Expected original object and generated object to have the same value for %s", key))
	}
}
