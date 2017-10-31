package logger

import (
	"os"

	logging "github.com/op/go-logging"
)

const (
	loggingModule = "asmImporter"
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
