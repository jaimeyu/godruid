package gathergrpc

import (
	"testing"

	tenmod "github.com/accedian/adh-gather/models/tenant"
	"github.com/stretchr/testify/assert"
)

func TestMonitoredObjectConversion(t *testing.T) {
	original := tenmod.MonitoredObject{
		ID:                    "theID",
		ActuatorName:          "Act1",
		ActuatorType:          "Good",
		MonitoredObjectID:     "BestID",
		Datatype:              "moType",
		ObjectName:            "Casper",
		CreatedTimestamp:      1234,
		LastModifiedTimestamp: 5678,
	}

	pbVersion := MonitoredObject{}
	err := ConvertToPBObject(original, &pbVersion)
	assert.Nil(t, err)
	assert.Equal(t, original.ID, pbVersion.XId)
	assert.Equal(t, original.REV, pbVersion.XRev)
	assert.NotNil(t, pbVersion.Data)
	assert.Equal(t, original.ActuatorName, pbVersion.Data.ActuatorName)
	assert.Equal(t, original.ActuatorType, pbVersion.Data.ActuatorType)
	assert.Equal(t, original.MonitoredObjectID, pbVersion.Data.Id)
	assert.Equal(t, original.Datatype, pbVersion.Data.Datatype)
	assert.Equal(t, original.ObjectName, pbVersion.Data.ObjectName)
	assert.Equal(t, original.ObjectType, pbVersion.Data.ObjectType)
	assert.Equal(t, original.CreatedTimestamp, pbVersion.Data.CreatedTimestamp)
	assert.Equal(t, original.LastModifiedTimestamp, pbVersion.Data.LastModifiedTimestamp)

	convertedBack := tenmod.MonitoredObject{}
	err = ConvertFromPBObject(pbVersion, &convertedBack)
	assert.Nil(t, err)
	assert.NotNil(t, convertedBack)
	assert.Equal(t, original, convertedBack)
}
