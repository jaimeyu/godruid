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

	assert.Equal(t, cfg.GetString("server.rest.ip"), "0.0.0.0")
	assert.Equal(t, cfg.GetString("server.datastore.ip"), "http://localhost")
	assert.Equal(t, cfg.GetBool("args.debug"), true)
	assert.Equal(t, cfg.GetString("args.admindb.name"), "aod-ui")
}

func TestLoadingDefaults(t *testing.T) {
	defer cleanUp()

	var cfg config.Provider
	v := viper.New()
	gather.LoadDefaults(v)
	cfg = v

	assert.Equal(t, cfg.GetString("server.rest.ip"), "0.0.0.0")
	assert.Equal(t, cfg.GetString("server.datastore.ip"), "http://localhost")
	assert.Equal(t, cfg.GetBool("args.debug"), false)
	assert.Equal(t, cfg.GetString("args.admindb.name"), "adh-admin")
}

func TestWithEnvironmentVariables(t *testing.T) {
	defer cleanUp()

	cfg := gather.LoadConfig("../config/adh-gather.yml", viper.New())

	os.Setenv("SERVER_REST_IP", "1.1.1.1")
	assert.Equal(t, cfg.GetString("server.rest.ip"), "1.1.1.1")

	os.Setenv("ARGS_DEBUG", "true")
	assert.Equal(t, cfg.GetBool("args.debug"), true)

	os.Setenv("ARGS_ADMINDB_NAME", "testname")
	assert.Equal(t, cfg.GetString("args.admindb.name"), "testname")
}
