package models

import (
	"encoding/json"
	"errors"

	pb "github.com/accedian/adh-gather/gathergrpc"
	admmod "github.com/accedian/adh-gather/models/admin"
	tenmod "github.com/accedian/adh-gather/models/tenant"
	"github.com/getlantern/deepcopy"
	"github.com/peterbourgon/mergemap"
)

const (
	LogRedactStr = "XXXXXXXX"
)

// AsJSONString - returns the object as a json string. If there is sensitive material in the object,
// this method can be augmented to hide those details.
func AsJSONString(obj interface{}) string {
	switch obj.(type) {
	case *pb.AdminUser:
		user := obj.(*pb.AdminUser)
		userCopy := pb.AdminUser{}
		deepcopy.Copy(&userCopy, user)
		userCopy.Data.Password = LogRedactStr
		res, err := json.Marshal(userCopy)
		if err != nil {
			return ""
		}
		return string(res)
	case *admmod.User:
		user := obj.(*admmod.User)
		userCopy := admmod.User{}
		deepcopy.Copy(&userCopy, user)
		userCopy.Password = LogRedactStr
		res, err := json.Marshal(userCopy)
		if err != nil {
			return ""
		}
		return string(res)
	case *tenmod.User:
		user := obj.(*tenmod.User)
		userCopy := tenmod.User{}
		deepcopy.Copy(&userCopy, user)
		userCopy.Password = LogRedactStr
		res, err := json.Marshal(userCopy)
		if err != nil {
			return ""
		}
		return string(res)
	case *pb.TenantUser:
		user := obj.(*pb.TenantUser)
		userCopy := pb.TenantUser{}
		deepcopy.Copy(&userCopy, user)
		userCopy.Data.Password = LogRedactStr
		res, err := json.Marshal(userCopy)
		if err != nil {
			return ""
		}
		return string(res)
	default:
		res, err := json.Marshal(obj)
		if err != nil {
			return ""
		}
		return string(res)
	}
}

/* MergeMaps - Merges two maps
 * It takes two maps and the dst map is the one getting modified.
 * If the dst map is empty, then the function should just copy the src to dst.
 * Otherwise, you can pre-populate the dst map with some data and the function will
 * update the keys with the new values.
 * Note dst is modified on error, potentially could be avoid by copying dst into a new container.
 */
func MergeMaps(dst map[string]interface{}, src map[string]interface{}) {
	mergemap.Merge(dst, src)
	return
}

/*MergeObjWithMap
 * Merges a JSON API request into an object
 * Trying to be as generic as possible to accept different kinds of structs
 * Works by transforming data several times to do the merge.
 *
 * Transforms the orig interface into a map.
 * Transforms the JSON []byte into a map.
 * Calls mergeMap which takes in two maps and does a merge.
 * Gets the merged map as an output.
 * Transforms the map back into the originating interface.
 *
 * It works because I'm trying to keep the input interface's struct meta data.
 * If we only deal with maps, we may not know about how to transform the maps
 * back into the struct.
 */
func MergeObjWithMap(orig interface{}, reqJson []byte) error {
	const KEY_DATA = "data"
	const KEY_ATTR = "attributes"
	requestMap := make(map[string]interface{})

	// Convert the request JSON into a map
	errMap := json.Unmarshal(reqJson, &requestMap)
	if errMap != nil {
		return errMap
	}

	omap, errconv := ConvertObj2Map(orig)
	if errconv != nil {
		return errconv
	}

	// Assumes the request map is based on JSON-API
	// Get the element with data
	if requestMap[KEY_DATA] == nil {
		return errors.New("Cannot merge as 'data' key is missing.")
	}
	mapData := requestMap[KEY_DATA].(map[string]interface{})
	// Get the element for attributes
	if mapData[KEY_ATTR] == nil {
		return errors.New("Cannot merge as 'attributes' key is missing.")
	}
	req := mapData[KEY_ATTR].(map[string]interface{})

	// Merge the request map into the original data map
	MergeMaps(omap, req)

	// Marshall the map into a JSON
	//fmt.Println("merged!")
	jstr, errT2 := json.Marshal(omap)
	if errT2 != nil {
		return errT2
	}

	// Unmarshal the data into an known struct.
	//fmt.Println("unshalling json to obj ", string(jstr), orig, reflect.TypeOf(orig))
	errUnMarsh := json.Unmarshal(jstr, orig)
	if errUnMarsh != nil {
		return errUnMarsh
	}
	//fmt.Printf("New merged object: %+v\n", orig)

	return nil
}

/*ConvertObj2Map --  Converts an object into a generic map */
func ConvertObj2Map(item interface{}) (map[string]interface{}, error) {
	orig, umerr := json.Marshal(item)
	if umerr != nil {
		return nil, umerr
	}

	var omap map[string]interface{}
	err := json.Unmarshal(orig, &omap)
	if err != nil {
		return nil, err
	}
	return omap, nil

}
