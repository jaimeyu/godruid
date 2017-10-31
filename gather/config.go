package gather

import (
	"io/ioutil"

	"github.com/ghodss/yaml"

	"github.com/accedian/adh-gather/logger"
)

// Config represents the global adh-fedex configuration parameters, as loaded from the config file
type Config struct {
}

func LoadConfig(cfgPath string) *Config {

	bData, err := ioutil.ReadFile(cfgPath)
	if err != nil {
		logger.Log.Panicf("Failed to open configuaration file '%s': %s", cfgPath, err.Error())

	}

	var cfg Config

	if err := yaml.Unmarshal(bData, &cfg); err != nil {
		logger.Log.Panicf("Failed to parse configuration file '%s': %s", cfgPath, err.Error())
	}

	return &cfg
}
