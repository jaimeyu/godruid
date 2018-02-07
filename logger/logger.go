package logger

import (
	"encoding/json"
	"os"

	"github.com/getlantern/deepcopy"

	logging "github.com/op/go-logging"
	pb "github.com/accedian/adh-gather/gathergrpc"
)

const (
	loggingModule = "asmImporter"

	LogRedactStr = "XXXXXXXX"
)

// Log is the project logger
var Log = logging.MustGetLogger(loggingModule)
var format = logging.MustStringFormatter(
	`%{color}%{time:15:04:05.000} %{shortfunc} ▶ %{level:.4s} %{id:03x}%{color:reset} %{message}`,
)

func init() {
	backend2 := logging.NewLogBackend(os.Stderr, "", 0)
	backend2Formatter := logging.NewBackendFormatter(backend2, format)
	backend2Leveled := logging.AddModuleLevel(backend2)
	backend2Leveled.SetLevel(logging.INFO, "")
	logging.SetBackend(backend2Formatter)

}

// SetDebugLevel enabled debug logging
func SetDebugLevel(t bool) {
	if t {
		logging.SetLevel(logging.DEBUG, loggingModule)
	} else {
		logging.SetLevel(logging.INFO, loggingModule)
	}
}

// IsDebugEnabled returns true if the current log level is set to debug
func IsDebugEnabled() bool {
	return logging.GetLevel(loggingModule) == logging.DEBUG
}

// AsJSONString - returns the object as a json string. If there is sensitive material in the object, 
// this method can be augmented to hide those details.
func AsJSONString(obj interface{}) string {
	switch obj.(type) {
	case *pb.AdminUser:
		user := obj.(*pb.AdminUser)
		userCopy := pb.AdminUser{}
		deepcopy.Copy(&userCopy, user)
		userCopy.Data.Password = LogRedactStr
		res, err := json.MarshalIndent(userCopy, "", "  ")
		if err != nil {
			return ""
		}
		return string(res)
	case *pb.TenantUser:
		user := obj.(*pb.TenantUser)
		userCopy := pb.TenantUser{}
		deepcopy.Copy(&userCopy, user)
		userCopy.Data.Password = LogRedactStr
		res, err := json.MarshalIndent(userCopy, "", "  ")
		if err != nil {
			return ""
		}
		return string(res)
	default:
		res, err := json.MarshalIndent(obj, "", "  ")
		if err != nil {
			return ""
		}
		return string(res)
	}
}