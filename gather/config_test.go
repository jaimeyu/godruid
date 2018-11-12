package gather_test

import (
	"os"
	"testing"

	"github.com/accedian/adh-gather/config"
	"github.com/accedian/adh-gather/gather"
	"github.com/magiconair/properties/assert"
	"github.com/spf13/viper"
)

func cleanUp() {
	os.Clearenv()
	viper.Reset()
}

func TestLoadingConfig(t *testing.T) {
	defer cleanUp()

	cfg := gather.LoadConfig("../config/adh-gather-debug.yml", viper.New())

	assert.Equal(t, cfg.GetString(gather.CK_server_rest_ip.String()), "0.0.0.0")
	assert.Equal(t, cfg.GetString(gather.CK_server_datastore_ip.String()), "http://localhost")
	assert.Equal(t, cfg.GetInt(gather.CK_server_datastore_batchsize.String()), 1000)
	assert.Equal(t, cfg.GetBool(gather.CK_args_debug.String()), true)
	assert.Equal(t, cfg.GetString(gather.CK_args_admindb_name.String()), "adh-admin")
	assert.Equal(t, cfg.GetInt(gather.CK_server_monitoring_port.String()), 9191)
	assert.Equal(t, cfg.GetBool(gather.CK_args_coltmef_enabled.String()), true)
	assert.Equal(t, cfg.GetInt(gather.CK_args_coltmef_statusretrycount.String()), 10)
	assert.Equal(t, cfg.GetFloat64(gather.CK_args_coltmef_checkpoint1.String()), float64(300))
	assert.Equal(t, cfg.GetFloat64(gather.CK_args_coltmef_checkpoint3.String()), float64(900))
	assert.Equal(t, cfg.GetInt(gather.CK_args_metricbaselines_numworkers.String()), 10)
}

func TestLoadingDefaults(t *testing.T) {
	defer cleanUp()

	var cfg config.Provider
	v := viper.New()
	gather.LoadDefaults(v)
	cfg = v

	assert.Equal(t, cfg.GetString(gather.CK_server_rest_ip.String()), "0.0.0.0")
	assert.Equal(t, cfg.GetString(gather.CK_server_datastore_ip.String()), "http://localhost")
	assert.Equal(t, cfg.GetBool(gather.CK_args_debug.String()), false)
	assert.Equal(t, cfg.GetString(gather.CK_args_admindb_name.String()), "adh-admin")
	assert.Equal(t, cfg.GetInt(gather.CK_server_monitoring_port.String()), 9191)
	assert.Equal(t, cfg.GetInt(gather.CK_args_maxConcurrentMetricAPICalls.String()), 500)
	assert.Equal(t, cfg.GetInt(gather.CK_args_maxConcurrentPouchAPICalls.String()), 1000)
	assert.Equal(t, cfg.GetInt(gather.CK_args_maxConcurrentProvAPICalls.String()), 1000)
	assert.Equal(t, cfg.GetInt(gather.CK_server_datastore_batchsize.String()), 1000)
	assert.Equal(t, cfg.GetBool(gather.CK_args_coltmef_enabled.String()), false)
	assert.Equal(t, cfg.GetInt(gather.CK_args_coltmef_statusretrycount.String()), 10)
	assert.Equal(t, cfg.GetFloat64(gather.CK_args_coltmef_checkpoint1.String()), float64(300))
	assert.Equal(t, cfg.GetFloat64(gather.CK_args_coltmef_checkpoint3.String()), float64(900))
	assert.Equal(t, cfg.GetInt(gather.CK_args_metricbaselines_numworkers.String()), 10)
}

func TestWithEnvironmentVariables(t *testing.T) {
	defer cleanUp()

	cfg := gather.LoadConfig("../config/adh-gather.yml", viper.New())

	os.Setenv("SERVER_REST_IP", "1.1.1.1")
	assert.Equal(t, cfg.GetString(gather.CK_server_rest_ip.String()), "1.1.1.1")

	os.Setenv("SERVER_DATASTORE_BATCHSIZE", "2000")
	assert.Equal(t, cfg.GetInt(gather.CK_server_datastore_batchsize.String()), 2000)

	os.Setenv("ARGS_DEBUG", "true")
	assert.Equal(t, cfg.GetBool(gather.CK_args_debug.String()), true)

	os.Setenv("ARGS_COLTMEF_ENABLED", "true")
	assert.Equal(t, cfg.GetBool(gather.CK_args_coltmef_enabled.String()), true)

	os.Setenv("ARGS_ADMINDB_NAME", "testname")
	assert.Equal(t, cfg.GetString(gather.CK_args_admindb_name.String()), "testname")

	os.Setenv("SERVER_CORS_ALLOWEDORIGINS", "http://stuff.com http://other.com")
	assert.Equal(t, cfg.GetStringSlice(gather.CK_server_cors_allowedorigins.String()), []string{"http://stuff.com", "http://other.com"})

	os.Setenv("ARGS_COLTMEF_STATUSRETRYCOUNT", "2000")
	assert.Equal(t, cfg.GetInt(gather.CK_args_coltmef_statusretrycount.String()), 2000)

	os.Setenv("ARGS_COLTMEF_CHECKPOINT1", "1000")
	assert.Equal(t, cfg.GetFloat64(gather.CK_args_coltmef_checkpoint1.String()), float64(1000))

	os.Setenv("ARGS_METRICBASELINES_NUMWORKERS", "2000")
	assert.Equal(t, cfg.GetInt(gather.CK_args_metricbaselines_numworkers.String()), 2000)
}
