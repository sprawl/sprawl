package config

import (
	"os"
	"strconv"
	"testing"

	"github.com/sprawl/sprawl/interfaces"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

const defaultConfigPath string = "default"
const testConfigPath string = "test"
const dbPathVar string = "database.path"
const rpcPortVar string = "rpc.port"
const p2pDebugVar string = "p2p.debug"
const errorsEnableStackTraceVar string = "errors.enableStackTrace"
const defaultDBPath string = "/var/lib/sprawl/data"
const defaultAPIPort uint = 1337
const testDBPath string = "/var/lib/sprawl/test"
const dbPathEnvVar string = "SPRAWL_DATABASE_PATH"
const rpcPortEnvVar string = "SPRAWL_RPC_PORT"
const p2pDebugEnvVar string = "SPRAWL_P2P_DEBUG"
const errorsEnableStackTraceEnvVar string = "SPRAWL_ERRORS_ENABLESTACKTRACE"
const envTestDBPath string = "/var/lib/sprawl/justforthistest"
const envTestAPIPort uint = 9001
const envTestP2PDebug string = "true"
const envTestErrorsEnableStackTrace string = "true"

var logger *zap.Logger
var log *zap.SugaredLogger
var config interfaces.Config
var databasePath string
var rpcPort uint
var p2pDebug bool
var errorsEnableStackTrace bool

func init() {
	logger, _ = zap.NewProduction()
	log = logger.Sugar()
	config = &Config{Logger: log}
}

func resetEnv() {
	os.Unsetenv(dbPathEnvVar)
	os.Unsetenv(rpcPortEnvVar)
	os.Unsetenv(p2pDebugEnvVar)
	os.Unsetenv(errorsEnableStackTraceEnvVar)
}

func TestPanics(t *testing.T) {
	resetEnv()
	// Tests for panics when not initialized with a config file
	assert.Panics(t, func() { databasePath = config.GetString(dbPathVar) }, "Config.GetString should panic when no config file or environment variables are provided")
	assert.Panics(t, func() { rpcPort = config.GetUint(rpcPortVar) }, "Config.GetUint should panic when no config file or environment variables are provided")
	assert.Panics(t, func() { config.Get(dbPathVar) }, "Config.Get should panic when no config file or environment variables are provided")
	assert.Equal(t, databasePath, "")
	assert.Equal(t, rpcPort, uint(0))
}

func TestDefaults(t *testing.T) {
	resetEnv()
	// Tests for defaults
	config.ReadConfig(defaultConfigPath)
	databasePath = config.GetString(dbPathVar)
	rpcPort = config.GetUint(rpcPortVar)
	p2pDebug = config.GetBool(p2pDebugVar)
	errorsEnableStackTrace = config.GetBool(errorsEnableStackTraceEnvVar)
	assert.Equal(t, databasePath, defaultDBPath)
	assert.Equal(t, rpcPort, defaultAPIPort)
	assert.False(t, p2pDebug)
	assert.False(t, errorsEnableStackTrace)
}

func TestTestVariables(t *testing.T) {
	resetEnv()
	config.ReadConfig(testConfigPath)
	databasePath = config.GetString(dbPathVar)
	rpcPort = config.GetUint(rpcPortVar)
	p2pDebug = config.GetBool(p2pDebugVar)
	errorsEnableStackTrace = config.GetBool(errorsEnableStackTraceVar)
	assert.Equal(t, databasePath, testDBPath)
	assert.Equal(t, rpcPort, defaultAPIPort)
	assert.False(t, p2pDebug)
	assert.False(t, errorsEnableStackTrace)
}

// TestEnvironment tests that environment variables overwrite any other configuration
func TestEnvironment(t *testing.T) {
	os.Setenv(dbPathEnvVar, envTestDBPath)
	os.Setenv(rpcPortEnvVar, strconv.FormatUint(uint64(envTestAPIPort), 10))
	os.Setenv(p2pDebugEnvVar, string(envTestP2PDebug))
	os.Setenv(errorsEnableStackTraceEnvVar, string(envTestErrorsEnableStackTrace))

	config.ReadConfig("")
	databasePath = config.GetString(dbPathVar)
	rpcPort = config.GetUint(rpcPortVar)
	p2pDebug = config.GetBool(p2pDebugVar)
	errorsEnableStackTrace = config.GetBool(errorsEnableStackTraceVar)

	assert.Equal(t, databasePath, envTestDBPath)
	assert.Equal(t, rpcPort, envTestAPIPort)
	assert.True(t, p2pDebug)
	assert.True(t, errorsEnableStackTrace)

	resetEnv()
}
