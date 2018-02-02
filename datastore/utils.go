package datastore

import (
	"fmt"
	"strings"
	"time"

	pb "github.com/accedian/adh-gather/gathergrpc"
	uuid "github.com/satori/go.uuid"
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
		return PrependToDataID(trimAndLowercase(cast.GetId()), dataType)
	default:
		uuid := uuid.NewV4()
		return PrependToDataID(uuid.String(), dataType)
	}
}

func trimAndLowercase(s string) string {
	return strings.TrimSpace(strings.ToLower(s))
}

// GetDataIDFromFullID - returns just the hash portion of an ID from the full ID.
func GetDataIDFromFullID(fullID string) string {
	parts := strings.Split(fullID, "_")

	if len(parts) == 0 {
		return ""
	}

	return parts[len(parts)-1]
}

// PrependToDataID - generates a full ID from the dataID and the dataType
func PrependToDataID(dataID string, dataType string) string {
	return fmt.Sprintf("%s%s%s", dataType, PouchDBIdBridgeStr, dataID)
}

// MakeTimestamp - get a timestamp from epoch in milliseconds
func MakeTimestamp() int64 {
    return time.Now().UnixNano() / int64(time.Millisecond)
}