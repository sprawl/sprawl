package config

import (
	"fmt"
	"strings"

	"github.com/spf13/viper"
	"github.com/sprawl/sprawl/errors"
)

const dbPathVar string = "database.path"
const rpcPortVar string = "rpc.port"
const p2pExternalIPVar string = "p2p.externalIP"
const p2pPortVar string = "p2p.port"
const p2pDebugVar string = "p2p.debug"
const p2pRelayVar string = "p2p.enableRelay"
const p2pAutoRelayVar string = "p2p.enableAutoRelay"
const p2pNATPortMapVar string = "p2p.enableNATPortMap"
const ipfsPeerVar string = "p2p.useIPFSPeers"
const errorsEnableStackTraceVar string = "errors.enableStackTrace"
const logLevelVar string = "log.level"
const logFormatVar string = "log.format"

// Config has an initialized version of spf13/viper
type Config struct {
	v *viper.Viper
}

// ReadConfig opens the configuration file and initializes viper
func (c *Config) ReadConfig(configPath string) {
	// Init viper
	c.v = viper.New()

	// Define where viper tries to get config information
	envPrefix := "sprawl"

	// Set environment variable prefix, automatically transformed to uppercase
	c.v.SetEnvPrefix(envPrefix)

	// Set replacer to env variables, replacing dots with underscores
	c.v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	// Automatically try to fetch all configs from env
	c.v.AutomaticEnv()

	// Initialize viper with Sprawl-specific options
	c.v.SetConfigName("config")

	// Use toml format for config files
	c.v.SetConfigType("toml")

	// Allow build to disable config file directories
	if configPath != "" {
		// Check for overriding config files
		c.v.AddConfigPath(".")
		// Check for user submitted config path
		c.v.AddConfigPath(configPath)
	}

	// Read config file
	if err := c.v.ReadInConfig(); !errors.IsEmpty(err) {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			fmt.Println("Config file not found, using ENV")
		} else {
			fmt.Println("Config file invalid!")
		}
	} else {
		fmt.Println("Config successfully loaded.")
	}
}

// GetDatabasePath defines the host directory for the database
func (c *Config) GetDatabasePath() string {
	return c.v.GetString(dbPathVar)
}

// GetExternalIP defines the listened external IP for P2P
func (c *Config) GetExternalIP() string {
	return c.v.GetString(p2pExternalIPVar)
}

// GetP2PPort defines the listened P2P port
func (c *Config) GetP2PPort() string {
	return c.v.GetString(p2pPortVar)
}

// GetRPCPort defines the port the gRPC is running at
func (c *Config) GetRPCPort() string {
	return c.v.GetString(rpcPortVar)
}

// GetNATPortMapSetting defines whether to use NAT port mapping or not
func (c *Config) GetNATPortMapSetting() bool {
	return c.v.GetBool(p2pNATPortMapVar)
}

// GetRelaySetting defines whether to run the node in relay mode or not
func (c *Config) GetRelaySetting() bool {
	return c.v.GetBool(p2pRelayVar)
}

// GetAutoRelaySetting defines whether to run the node in autorelay mode or not
func (c *Config) GetAutoRelaySetting() bool {
	return c.v.GetBool(p2pAutoRelayVar)
}

// GetDebugSetting defines whether to run the debug pinger or not
func (c *Config) GetDebugSetting() bool {
	return c.v.GetBool(p2pDebugVar)
}

// GetStackTraceSetting defines whether to run the debug pinger or not
func (c *Config) GetStackTraceSetting() bool {
	return c.v.GetBool(errorsEnableStackTraceVar)
}

// GetIPFSPeerSetting defines if we use IPFS bootstrap peers for discovery or just our own
func (c *Config) GetIPFSPeerSetting() bool {
	return c.v.GetBool(ipfsPeerVar)
}

// GetLogLevel gets configured log level for uber/zap
func (c *Config) GetLogLevel() string {
	return c.v.GetString(logLevelVar)
}

// GetLogFormat gets configured log format for uber/zap
func (c *Config) GetLogFormat() string {
	return c.v.GetString(logFormatVar)
}
