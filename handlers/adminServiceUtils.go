package handlers

import (
	"encoding/json"
	"io/ioutil"

	"github.com/accedian/adh-gather/gather"
	"github.com/accedian/adh-gather/logger"
	admmod "github.com/accedian/adh-gather/models/admin"
	tenmod "github.com/accedian/adh-gather/models/tenant"
)

var (
	// ValidMonitoredObjectTypes - known Monitored Object types in the system
	ValidMonitoredObjectTypes = map[string]tenmod.MonitoredObjectType{
		"pe": tenmod.TwampPE,
		"sf": tenmod.TwampSF,
		"sl": tenmod.TwampSL,
		string(tenmod.TwampPE): tenmod.TwampPE,
		string(tenmod.TwampSF): tenmod.TwampSF,
		string(tenmod.TwampSL): tenmod.TwampSL}

	// ValidMonitoredObjectDeviceTypes - known Monitored Object Device types in the system.
	ValidMonitoredObjectDeviceTypes = map[string]tenmod.MonitoredObjectDeviceType{
		string(tenmod.AccedianNID):  tenmod.AccedianNID,
		string(tenmod.AccedianVNID): tenmod.AccedianVNID}

	// DefaultValidTypes - default values for the valid types supported by datahub
	DefaultValidTypes = &admmod.ValidTypes{}

	// DefaultIngestionDictionary - default values for the Ingestion Dictionary
	DefaultIngestionDictionary = &admmod.IngestionDictionary{}
)

// getIngestionDictionaryFromFile - retrieves the contents of the IngestionDictionary.
func getIngestionDictionaryFromFile() (*admmod.IngestionDictionary, error) {
	if DefaultIngestionDictionary == nil || DefaultIngestionDictionary.ID == "" {
		cfg := gather.GetConfig()
		ingDictFilePath := cfg.GetString("ingDict")
		defaultDictionaryBytes, err := ioutil.ReadFile(ingDictFilePath)
		if err != nil {
			logger.Log.Fatalf("Unable to read Default Ingestion Dictionary from file: %s", err.Error())
		}

		DefaultIngestionDictionary = &admmod.IngestionDictionary{}
		if err = json.Unmarshal(defaultDictionaryBytes, DefaultIngestionDictionary); err != nil {
			logger.Log.Fatalf("Unable to construct Default Ingestion Dictionary from file: %s", err.Error())
		}

		DefaultIngestionDictionary.ID = "1"
		DefaultIngestionDictionary.REV = "1"
		DefaultIngestionDictionary.Datatype = "ingestionDictionary"
	}

	return DefaultIngestionDictionary, nil
}

// getValidTypes - retrieves the known valid types
func getValidTypes() *admmod.ValidTypes {
	if DefaultValidTypes == nil || DefaultValidTypes.ID == "" {
		validMonObjTypes := make(map[string]string, 0)
		validMonObjDevTypes := make(map[string]string, 0)

		for key, val := range ValidMonitoredObjectTypes {
			validMonObjTypes[key] = string(val)
		}
		for key, val := range ValidMonitoredObjectDeviceTypes {
			validMonObjDevTypes[key] = string(val)
		}

		DefaultValidTypes = &admmod.ValidTypes{
			ID:       "1",
			REV:      "1",
			Datatype: "validTypes",

			MonitoredObjectTypes:       validMonObjTypes,
			MonitoredObjectDeviceTypes: validMonObjDevTypes}
	}

	return DefaultValidTypes
}
