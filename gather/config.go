package gather

import (
	"strings"

	"github.com/accedian/adh-gather/config"
	"github.com/accedian/adh-gather/logger"
	"github.com/spf13/viper"
)

// DBImpl - type which describes a Database Implementation technology.
type DBImpl int

// This set of constants acts as an enumeration for Database Implementation types.
const (
	MEM DBImpl = iota
	COUCH
)

// Config represents the global adh-gather configuration parameters, as loaded from the config file
type Config struct {
	server struct {
		rest struct {
			ip   string
			port int
		}
		datastore struct {
			ip   string
			port int
		}
		grpc struct {
			ip   string
			port int
		}
	}
	args struct {
		admindb struct {
			name string
			impl DBImpl
		}
		tenantdb struct {
			impl DBImpl
		}
		pouchplugindb struct {
			impl DBImpl
		}
		debug bool
	}
}

// Stores the active configuration for the running instance.
var cfg config.Provider

// GetConfig - returns the current configuration.
func GetConfig() config.Provider {
	return cfg
}

// LoadConfig - implements configuration based on the provided file
func LoadConfig(cfgPath string, v *viper.Viper) config.Provider {
	replacer := strings.NewReplacer(".", "_")
	v.SetEnvKeyReplacer(replacer)

	v.AutomaticEnv()

	v.SetConfigFile(cfgPath)
	v.SetConfigType("yaml")
	err := v.ReadInConfig()
	if err != nil {
		logger.Log.Panicf("Failed to parse configuration file '%s': %s",
			cfgPath, err.Error())
	}

	LoadDefaults(v)

	cfg = v
	return cfg
}

// LoadDefaults - loads default values for the configuration
func LoadDefaults(v *viper.Viper) {
	v.SetDefault("server.datastore.ip", "http://localhost")
	v.SetDefault("server.datastore.port", 5984)
	v.SetDefault("server.rest.ip", "0.0.0.0")
	v.SetDefault("server.rest.port", 10001)
	v.SetDefault("server.grpc.ip", "0.0.0.0")
	v.SetDefault("server.grpc.port", 10002)
	v.SetDefault("args.admindb.name", "adh-admin")
	v.SetDefault("args.admindb.impl", 1)
	v.SetDefault("args.tenantdb.impl", 1)
	v.SetDefault("args.pouchplugindb.impl", 1)
	v.SetDefault("args.debug", false)
}
