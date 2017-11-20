package gather

import (
	"errors"
	"io/ioutil"

	"github.com/ghodss/yaml"

	"github.com/accedian/adh-gather/logger"
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
	ServerConfig struct {
		REST struct {
			BindIP   string
			BindPort int
		}
		Datastore struct {
			BindIP   string
			BindPort int
		}
		GRPC struct {
			BindIP   string
			BindPort int
		}
		StartupArgs struct {
			AdminDB struct {
				Name string
				Impl DBImpl
			}
			TenantDB      DBImpl
			PouchPluginDB DBImpl
			Debug         bool
		}
	}
}

// Stores the active configuration for the running instance.
var activeConfig *Config

func LoadConfig(cfgPath string) *Config {

	bData, err := ioutil.ReadFile(cfgPath)
	if err != nil {
		logger.Log.Panicf("Failed to open configuaration file '%s': %s", cfgPath, err.Error())

	}

	var cfg Config

	if err := yaml.Unmarshal(bData, &cfg); err != nil {
		logger.Log.Panicf("Failed to parse configuration file '%s': %s", cfgPath, err.Error())
	}

	activeConfig = &cfg

	return activeConfig
}

// GetActiveConfig - returns the active configuration for the running process,
// or an error if the configuration has not been loaded.
func GetActiveConfig() (*Config, error) {
	if activeConfig == nil {
		return nil, errors.New("Please run LoadConfig before trying to access Config Parameters")
	}

	return activeConfig, nil
}
