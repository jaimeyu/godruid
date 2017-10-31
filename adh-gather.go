package main

import (
	"flag"
	"fmt"

	"github.com/accedian/adh-gather/gather"
	"github.com/accedian/adh-gather/logger"
)

var (
	configFilePath string
	debug          bool
)

func init() {
	flag.StringVar(&configFilePath, "config", "config/adh-gather.yml", "Specify a configuration file to use")
	flag.BoolVar(&debug, "debug", false, "Enable debug mode (and logs)")
}

func main() {
	flag.Parse()

	if debug {
		logger.SetDebugLevel(true)
	} else {
		logger.SetDebugLevel(false)
	}

	logger.Log.Infof("Starting adh-gather broker with config '%s'", configFilePath)

	cfg := gather.LoadConfig(configFilePath)

	fmt.Printf("Your config is %+v \n", cfg)

	logger.Log.Infof("Stopping broker")
}
