package config

import (
	"os"
	"testing"

	"github.com/eqlabs/sprawl/interfaces"
	"github.com/stretchr/testify/assert"
)

var config interfaces.Config = &Config{}
var databasePath string
var apiPort uint

func resetEnv() {
	os.Unsetenv("SPRAWL_DATABASE_PATH")
	os.Unsetenv("SPRAWL_API_PORT")
}

func TestPanics(t *testing.T) {
	// Tests for panics when not initialized with a config file
	assert.Panics(t, func() { databasePath = config.GetString("database.path") }, "Config.GetString should panic when no config file or environment variables are provided")
	assert.Panics(t, func() { apiPort = config.GetUint("api.port") }, "Config.GetUint should panic when no config file or environment variables are provided")
	assert.Panics(t, func() { config.Get("database.path") }, "Config.Get should panic when no config file or environment variables are provided")
	assert.Equal(t, databasePath, "")
	assert.Equal(t, apiPort, uint(0))
}

func TestDefaults(t *testing.T) {
	// Tests for defaults
	config.ReadConfig("default")
	databasePath = config.GetString("database.path")
	apiPort = config.GetUint("api.port")
	assert.Equal(t, databasePath, "/var/lib/sprawl/data")
	assert.Equal(t, apiPort, uint(1337))
}

func TestTestVariables(t *testing.T) {
	config.ReadConfig("test")
	databasePath = config.GetString("database.path")
	apiPort = config.GetUint("api.port")
	assert.Equal(t, databasePath, "/var/lib/sprawl/test")
	assert.Equal(t, apiPort, uint(1337))
}

// TestEnvironment tests that environment variables overwrite any other configuration
func TestEnvironment(t *testing.T) {
	os.Setenv("SPRAWL_DATABASE_PATH", "/var/lib/sprawl/justforthistest")
	os.Setenv("SPRAWL_API_PORT", "9001")
	config.ReadConfig("asd")
	databasePath = config.GetString("database.path")
	apiPort = config.GetUint("api.port")
	assert.Equal(t, databasePath, "/var/lib/sprawl/justforthistest")
	assert.Equal(t, apiPort, uint(9001))
	resetEnv()
}
