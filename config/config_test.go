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
const dbInMemoryVar string = "database.inMemory"
const rpcPortVar string = "rpc.port"
const websocketPortVar string = "websocket.port"
const p2pDebugVar string = "p2p.debug"
const errorsEnableStackTraceVar string = "errors.enableStackTrace"
const websocketEnableVar string = "websocket.enable"
const defaultDBPath string = "/var/lib/sprawl/data"
const defaultAPIPort uint = 1337
const defaultWebsocketPort uint = 3000
const testDBPath string = "/var/lib/sprawl/test"
const dbPathEnvVar string = "SPRAWL_DATABASE_PATH"
const useInMemoryEnvVar string = "SPRAWL_DATABASE_INMEMORY"
const rpcPortEnvVar string = "SPRAWL_RPC_PORT"
const websocketPortEnvVar string = "SPRAWL_WEBSOCKET_PORT"
const p2pDebugEnvVar string = "SPRAWL_P2P_DEBUG"
const errorsEnableStackTraceEnvVar string = "SPRAWL_ERRORS_ENABLESTACKTRACE"
const websocketEnableEnvVar string = "SPRAWL_WEBSOCKET_ENABLE"
const envTestDBPath string = "/var/lib/sprawl/justforthistest"
const envTestAPIPort uint = 9001
const envTestWebsocketPort uint = 8000
const envTestP2PDebug string = "true"
const envTestErrorsEnableStackTrace string = "true"
const envTestUseInMemory string = "true"
const envTestWebsocketEnable string = "true"

var logger *zap.Logger
var log *zap.SugaredLogger
var config interfaces.Config
var databasePath string
var useInMemory bool

var websocketPort uint
var rpcPort uint
var p2pDebug bool
var errorsEnableStackTrace bool
var websocketEnable bool

func init() {
	logger = zap.NewNop()
	log = logger.Sugar()
	config = &Config{Logger: log}
}

func resetEnv() {
	os.Unsetenv(dbPathEnvVar)
	os.Unsetenv(rpcPortEnvVar)
	os.Unsetenv(websocketPortEnvVar)
	os.Unsetenv(p2pDebugEnvVar)
	os.Unsetenv(errorsEnableStackTraceEnvVar)
	os.Unsetenv(useInMemoryEnvVar)
	os.Unsetenv(websocketEnableEnvVar)
}

func TestPanics(t *testing.T) {
	resetEnv()
	// Tests for panics when not initialized with a config file
	assert.Panics(t, func() { databasePath = config.GetString(dbPathVar) }, "Config.GetString should panic when no config file or environment variables are provided")
	assert.Panics(t, func() { rpcPort = config.GetUint(rpcPortVar) }, "Config.GetUint should panic when no config file or environment variables are provided")
	assert.Panics(t, func() { websocketPort = config.GetUint(websocketPortVar) }, "Config.GetUint should panic when no config file or environment variables are provided")
	assert.Panics(t, func() { config.Get(dbPathVar) }, "Config.Get should panic when no config file or environment variables are provided")
	assert.Panics(t, func() { config.Get(dbPathVar) }, "Config.Get should panic when no config file or environment variables are provided")
	assert.Equal(t, databasePath, "")
	assert.Equal(t, rpcPort, uint(0))
	assert.Equal(t, websocketPort, uint(0))
}

func TestDefaults(t *testing.T) {
	resetEnv()
	// Tests for defaults
	config.ReadConfig(defaultConfigPath)
	databasePath = config.GetString(dbPathVar)
	rpcPort = config.GetUint(rpcPortVar)
	websocketPort = config.GetUint(websocketPortVar)
	p2pDebug = config.GetBool(p2pDebugVar)
	errorsEnableStackTrace = config.GetBool(errorsEnableStackTraceVar)
	useInMemory = config.GetBool(dbInMemoryVar)
	websocketEnable = config.GetBool(websocketEnableVar)
	assert.Equal(t, databasePath, defaultDBPath)
	assert.Equal(t, rpcPort, defaultAPIPort)
	assert.Equal(t, websocketPort, defaultWebsocketPort)
	assert.False(t, p2pDebug)
	assert.False(t, errorsEnableStackTrace)
	assert.False(t, useInMemory)
	assert.False(t, websocketEnable)
}

func TestTestVariables(t *testing.T) {
	resetEnv()
	config.ReadConfig(testConfigPath)
	databasePath = config.GetString(dbPathVar)
	rpcPort = config.GetUint(rpcPortVar)
	websocketPort = config.GetUint(websocketPortVar)
	p2pDebug = config.GetBool(p2pDebugVar)
	errorsEnableStackTrace = config.GetBool(errorsEnableStackTraceVar)
	useInMemory = config.GetBool(dbInMemoryVar)
	websocketEnable = config.GetBool(websocketEnableVar)
	assert.Equal(t, databasePath, testDBPath)
	assert.Equal(t, rpcPort, defaultAPIPort)
	assert.Equal(t, websocketPort, defaultWebsocketPort)
	assert.False(t, p2pDebug)
	assert.False(t, errorsEnableStackTrace)
	assert.True(t, useInMemory)
	assert.True(t, websocketEnable)

}

// TestEnvironment tests that environment variables overwrite any other configuration
func TestEnvironment(t *testing.T) {
	os.Setenv(dbPathEnvVar, envTestDBPath)
	os.Setenv(rpcPortEnvVar, strconv.FormatUint(uint64(envTestAPIPort), 10))
	os.Setenv(websocketPortEnvVar, strconv.FormatUint(uint64(envTestWebsocketPort), 10))
	os.Setenv(p2pDebugEnvVar, string(envTestP2PDebug))
	os.Setenv(errorsEnableStackTraceEnvVar, string(envTestErrorsEnableStackTrace))
	os.Setenv(useInMemoryEnvVar, string(envTestUseInMemory))
	os.Setenv(websocketEnableEnvVar, string(envTestWebsocketEnable))

	config.ReadConfig("")
	databasePath = config.GetString(dbPathVar)
	rpcPort = config.GetUint(rpcPortVar)
	websocketPort = config.GetUint(websocketPortVar)
	p2pDebug = config.GetBool(p2pDebugVar)
	errorsEnableStackTrace = config.GetBool(errorsEnableStackTraceVar)
	useInMemory = config.GetBool(dbInMemoryVar)
	websocketEnable = config.GetBool(websocketEnableVar)

	assert.Equal(t, databasePath, envTestDBPath)
	assert.Equal(t, rpcPort, envTestAPIPort)
	assert.Equal(t, websocketPort, envTestWebsocketPort)
	assert.True(t, p2pDebug)
	assert.True(t, errorsEnableStackTrace)
	assert.True(t, useInMemory)
	assert.True(t, websocketEnable)

	resetEnv()
}
