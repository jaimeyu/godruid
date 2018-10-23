package gather

import (
	"strings"
	"time"

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
			ip        string
			port      int
			batchsize int
		}
		grpc struct {
			ip   string
			port int
		}
		monitoring struct {
			port int
		}
		cors struct {
			allowedorigins []string
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
		testdatadb struct {
			impl DBImpl
		}
		maxConcurrentPouchAPICalls  uint64
		maxConcurrentProvAPICalls   uint64
		maxConcurrentMetricAPICalls uint64
		debug                       bool
		coltmef                     struct {
			enabled          bool
			server           string
			appID            string
			secret           string
			statusRetryCount int
			checkpoint1      float64
			checkpoint2      float64
			checkpoint3      float64
		}
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
	v.SetDefault(CK_server_datastore_ip.String(), "http://localhost")
	v.SetDefault(CK_server_datastore_port.String(), 5984)
	v.SetDefault(CK_server_datastore_batchsize.String(), 1000)
	v.SetDefault(CK_server_rest_ip.String(), "0.0.0.0")
	v.SetDefault(CK_server_rest_port.String(), 10001)
	v.SetDefault(CK_server_monitoring_port.String(), 9191)
	v.SetDefault(CK_server_profile_port.String(), 6060)
	v.SetDefault(CK_server_changenotif_refreshFreqSeconds.String(), 300)
	v.SetDefault(CK_server_grpc_ip.String(), "0.0.0.0")
	v.SetDefault(CK_server_grpc_port.String(), 10002)
	v.SetDefault(CK_args_admindb_name.String(), "adh-admin")
	v.SetDefault(CK_args_admindb_impl.String(), 1)
	v.SetDefault(CK_args_tenantdb_impl.String(), 1)
	v.SetDefault(CK_args_pouchplugindb_impl.String(), 1)
	v.SetDefault(CK_args_testdatadb_impl.String(), 1)
	v.SetDefault(CK_args_debug.String(), false)
	v.SetDefault(CK_args_debug.String(), false)
	v.SetDefault(CK_args_debug.String(), false)
	v.SetDefault(CK_args_debug.String(), false)
	v.SetDefault(CK_args_authorizationAAA.String(), false)
	v.SetDefault(CK_connector_maxSecondsWithoutHeartbeat.String(), 120)
	v.SetDefault(CK_args_maxConcurrentMetricAPICalls.String(), 500)
	v.SetDefault(CK_args_maxConcurrentProvAPICalls.String(), 1000)
	v.SetDefault(CK_args_maxConcurrentPouchAPICalls.String(), 1000)
	v.SetDefault(CK_kafka_broker.String(), "kafka:9092")
	v.SetDefault(CK_druid_timeoutsms_histogram.String(), 45000)
	v.SetDefault(CK_druid_timeoutsms_slareports.String(), 45000)
	v.SetDefault(CK_druid_timeoutsms_thresholdcrossing.String(), 45000)
	v.SetDefault(CK_druid_timeoutsms_aggregatedmetrics.String(), 45000)
	v.SetDefault(CK_druid_timeoutsms_rawmetrics.String(), 45000)
	v.SetDefault(CK_druid_timeoutsms_filteredrawmetrics.String(), 45000)
	v.SetDefault(CK_args_coltmef_enabled.String(), false)
	v.SetDefault(CK_args_coltmef_server.String(), "https://demo.ondemand.colt.net/api/performance/recommendation")
	v.SetDefault(CK_args_coltmef_appid.String(), "datahubtest")
	v.SetDefault(CK_args_coltmef_secret.String(), "test123")
	v.SetDefault(CK_args_coltmef_statusretrycount.String(), 10)
	v.SetDefault(CK_args_coltmef_checkpoint1.String(), (5 * time.Minute).Seconds())
	v.SetDefault(CK_args_coltmef_checkpoint2.String(), (10 * time.Minute).Seconds())
	v.SetDefault(CK_args_coltmef_checkpoint3.String(), (15 * time.Minute).Seconds())
}
