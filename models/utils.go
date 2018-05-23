package models

import (
	"encoding/json"

	pb "github.com/accedian/adh-gather/gathergrpc"
	admmod "github.com/accedian/adh-gather/models/admin"
	metmod "github.com/accedian/adh-gather/models/metrics"
	tenmod "github.com/accedian/adh-gather/models/tenant"
	"github.com/getlantern/deepcopy"
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
	case *metmod.ReportScheduleConfig:
		user := obj.(*metmod.ReportScheduleConfig)
		userCopy := metmod.ReportScheduleConfig{}
		deepcopy.Copy(&userCopy, user)
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
