package datastore

import (
	"fmt"
	"strings"
	"time"

	pb "github.com/accedian/adh-gather/gathergrpc"
	tenmod "github.com/accedian/adh-gather/models/tenant"
	uuid "github.com/satori/go.uuid"
)

var (
	NotFoundStr = "status 404 - not found"
)

const (

	// PouchDBIdBridgeStr - value required for pouchDB to properly identify a an item. Used in the
	// following way to generate an ID:
	//      id = <dataType> + PouchDBIdBridgeStr + <generatedUUID>
	PouchDBIdBridgeStr = "_2_"
)

// GenerateID - generates an ID for an object based on the type of the
// provided object. Returns 2 versions of the ID, the first is the full ID used
// to refer to the object in Couch, the second is the hash portion of the Couch ID
// which is used to refer to the data object.
func GenerateID(obj interface{}, dataType string) string {

	switch obj.(type) {
	case *pb.MonitoredObjectData:
		cast := obj.(*pb.MonitoredObjectData)
		return PrependToDataID(strings.TrimSpace(cast.GetId()), dataType)
	case *tenmod.MonitoredObject:
		cast := obj.(*tenmod.MonitoredObject)
		return PrependToDataID(strings.TrimSpace(cast.MonitoredObjectID), dataType)
	default:
		uuid := uuid.NewV4()
		return PrependToDataID(uuid.String(), dataType)
	}
}

func trimAndLowercase(s string) string {
	return strings.TrimSpace(strings.ToLower(s))
}

// GetDataIDFromFullID - returns just the hash portion of an ID from the full ID.
// Note that for monitoredObjects, a hash is not stored in the ID field.
// The datatype is appended by a name passed in through ingestion and may have
// underscores in its name. There was a bug in the previous logic that incorrectly
// stripped all substrings prepended with _.
// The fix here is to find the index of the _FIRST TWO_ underscores in the fullID string.
// Then only providing the slice AFTER the second underscore.
// The behaviour here is based on the fact we prepend the unique ID with the
// document's data type, eg: 'monitoredObject_2_'. See #GenerateID() on the actual behavior.
func GetDataIDFromFullID(fullID string) string {
	const OFFSET int = 2
	const UNDERSCORE string = "_"
	if len(fullID) == 0 {
		return ""
	}
	loc1 := strings.Index(fullID, UNDERSCORE)
	if loc1 < 0 {
		return fullID
	}
	loc2 := strings.Index(fullID[loc1+1:], UNDERSCORE)
	if loc2 < 0 {
		return fullID
	}

	stripIdx := loc1 + loc2 + OFFSET
	stripped := fullID[stripIdx:]
	return stripped
}

// PrependToDataID - generates a full ID from the dataID and the dataType
func PrependToDataID(dataID string, dataType string) string {
	return fmt.Sprintf("%s%s%s", dataType, PouchDBIdBridgeStr, dataID)
}

// MakeTimestamp - get a timestamp from epoch in milliseconds
func MakeTimestamp() int64 {
	return time.Now().UnixNano() / int64(time.Millisecond)
}
