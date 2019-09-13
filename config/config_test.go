package config

import (
	"os"
	"strconv"
	"testing"

	"github.com/eqlabs/sprawl/interfaces"
	"github.com/stretchr/testify/assert"
)

const defaultConfigPath string = "default"
const testConfigPath string = "test"
const dbPathVar string = "database.path"
const apiPortVar string = "api.port"
const p2pDebugVar string = "p2p.debug"
const defaultDBPath string = "/var/lib/sprawl/data"
const defaultAPIPort uint = 1337
const testDBPath string = "/var/lib/sprawl/test"
const dbPathEnvVar string = "SPRAWL_DATABASE_PATH"
const apiPortEnvVar string = "SPRAWL_API_PORT"
const p2pDebugEnvVar string = "SPRAWL_P2P_DEBUG"
const envTestDBPath string = "/var/lib/sprawl/justforthistest"
const envTestAPIPort uint = 9001
const envTestP2PDebug string = "true"

var config interfaces.Config = &Config{}
var databasePath string
var apiPort uint
var p2pDebug bool

func resetEnv() {
	os.Unsetenv(dbPathEnvVar)
	os.Unsetenv(apiPortEnvVar)
}

func TestPanics(t *testing.T) {
	resetEnv()
	// Tests for panics when not initialized with a config file
	assert.Panics(t, func() { databasePath = config.GetString(dbPathVar) }, "Config.GetString should panic when no config file or environment variables are provided")
	assert.Panics(t, func() { apiPort = config.GetUint(apiPortVar) }, "Config.GetUint should panic when no config file or environment variables are provided")
	assert.Panics(t, func() { config.Get(dbPathVar) }, "Config.Get should panic when no config file or environment variables are provided")
	assert.Equal(t, databasePath, "")
	assert.Equal(t, apiPort, uint(0))
}

func TestDefaults(t *testing.T) {
	resetEnv()
	// Tests for defaults
	config.ReadConfig(defaultConfigPath)
	databasePath = config.GetString(dbPathVar)
	apiPort = config.GetUint(apiPortVar)
	p2pDebug = config.GetBool(p2pDebugVar)
	assert.Equal(t, databasePath, defaultDBPath)
	assert.Equal(t, apiPort, defaultAPIPort)
	assert.False(t, p2pDebug)
}

func TestTestVariables(t *testing.T) {
	resetEnv()
	config.ReadConfig(testConfigPath)
	databasePath = config.GetString(dbPathVar)
	apiPort = config.GetUint(apiPortVar)
	p2pDebug = config.GetBool(p2pDebugVar)
	assert.Equal(t, databasePath, testDBPath)
	assert.Equal(t, apiPort, defaultAPIPort)
	assert.False(t, p2pDebug)
}

// TestEnvironment tests that environment variables overwrite any other configuration
func TestEnvironment(t *testing.T) {
	os.Setenv(dbPathEnvVar, envTestDBPath)
	os.Setenv(apiPortEnvVar, strconv.FormatUint(uint64(envTestAPIPort), 10))
	os.Setenv(p2pDebugEnvVar, string(envTestP2PDebug))

	config.ReadConfig("")
	databasePath = config.GetString(dbPathVar)
	apiPort = config.GetUint(apiPortVar)
	p2pDebug = config.GetBool(p2pDebugVar)

	assert.Equal(t, databasePath, envTestDBPath)
	assert.Equal(t, apiPort, envTestAPIPort)
	assert.True(t, envTestP2PDebug)

	resetEnv()
}
