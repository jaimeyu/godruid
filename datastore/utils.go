package datastore

import (
	"fmt"

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
// provided object.
func GenerateID(obj interface{}, dataType string) string {
	switch obj.(type) {
	case *pb.MonitoredObject:
		cast := obj.(*pb.MonitoredObject)
		return fmt.Sprintf("%s%s%s", dataType, PouchDBIdBridgeStr, cast.GetId())
	default:
		return fmt.Sprintf("%s%s%s", dataType, PouchDBIdBridgeStr, uuid.NewV4().String())
	}
}
