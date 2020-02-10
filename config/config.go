package config

import (
	"fmt"
	"strings"

	"github.com/spf13/cast"
	"github.com/spf13/viper"
	"github.com/sprawl/sprawl/errors"
)

const dbPathVar string = "database.path"
const dbInMemoryVar string = "database.inMemory"
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
const websocketEnableVar string = "websocket.enable"
const websocketPortVar string = "websocket.port"

// Config has an initialized version of spf13/viper
type Config struct {
	v        *viper.Viper
	strings  map[string]string
	booleans map[string]bool
}

// ReadConfig opens the configuration file and initializes viper
func (c *Config) ReadConfig(configPath string) {
	// Init viper
	c.v = viper.New()

	// Init maps where the config will be stored
	c.strings = make(map[string]string)
	c.booleans = make(map[string]bool)

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

	c.AddString(dbPathVar)
	c.AddString(p2pExternalIPVar)
	c.AddString(p2pPortVar)
	c.AddString(rpcPortVar)
	c.AddString(websocketPortVar)
	c.AddString(logLevelVar)
	c.AddString(logFormatVar)
	c.AddBoolean(websocketEnableVar)
	c.AddBoolean(dbInMemoryVar)
	c.AddBoolean(p2pNATPortMapVar)
	c.AddBoolean(p2pRelayVar)
	c.AddBoolean(p2pAutoRelayVar)
	c.AddBoolean(p2pDebugVar)
	c.AddBoolean(errorsEnableStackTraceVar)
	c.AddBoolean(ipfsPeerVar)

}

// AddString to config and print a message, if default is used.
func (c *Config) AddString(key string) {
	err := c.AddStringE(key)
	if err != nil {
		fmt.Println(key + ": set to \"\"")
	}
}

// AddBoolean to config and print a message, if default is used.
func (c *Config) AddBoolean(key string) {
	err := c.AddBooleanE(key)
	if err != nil {
		fmt.Println(key + ": set to false")
	}
}

// AddStringE (default "") to config and return error
func (c *Config) AddStringE(key string) error {
	s, err := cast.ToStringE(c.v.Get(key))
	c.strings[key] = s
	return err
}

// AddBooleanE (default false) to config and return error
func (c *Config) AddBooleanE(key string) error {
	b, err := cast.ToBoolE(c.v.Get(key))
	c.booleans[key] = b
	return err
}

// GetDatabasePath defines the host directory for the database
func (c *Config) GetDatabasePath() string {
	return c.strings[dbPathVar]
}

// GetExternalIP defines the listened external IP for P2P
func (c *Config) GetExternalIP() string {
	return c.strings[p2pExternalIPVar]
}

// GetP2PPort defines the listened P2P port
func (c *Config) GetP2PPort() string {
	return c.strings[p2pPortVar]
}

// GetRPCPort defines the port the gRPC is running at
func (c *Config) GetRPCPort() string {
	return c.strings[rpcPortVar]
}

// GetWebsocketPort defines port for websocket connections. websocket.enable must be true or the port is not used.
func (c *Config) GetWebsocketPort() string {
	return c.strings[websocketPortVar]
}

// GetLogLevel gets configured log level for uber/zap
func (c *Config) GetLogLevel() string {
	return c.strings[logLevelVar]
}

// GetLogFormat gets configured log format for uber/zap
func (c *Config) GetLogFormat() string {
	return c.strings[logFormatVar]
}

// GetWebsocketEnable defines if websocket connections are allowed. Starts waiting http request using websocket.port
func (c *Config) GetWebsocketEnable() bool {
	return c.booleans[websocketEnableVar]
}

// GetInMemoryDatabaseSetting defines if RAM is used instead of LevelDB for storage
func (c *Config) GetInMemoryDatabaseSetting() bool {
	return c.booleans[dbInMemoryVar]
}

// GetNATPortMapSetting defines whether to use NAT port mapping or not
func (c *Config) GetNATPortMapSetting() bool {
	return c.booleans[p2pNATPortMapVar]
}

// GetRelaySetting defines whether to run the node in relay mode or not
func (c *Config) GetRelaySetting() bool {
	return c.booleans[p2pRelayVar]
}

// GetAutoRelaySetting defines whether to run the node in autorelay mode or not
func (c *Config) GetAutoRelaySetting() bool {
	return c.booleans[p2pAutoRelayVar]
}

// GetDebugSetting defines whether to run the debug pinger or not
func (c *Config) GetDebugSetting() bool {
	return c.booleans[p2pDebugVar]
}

// GetStackTraceSetting defines whether to run the debug pinger or not
func (c *Config) GetStackTraceSetting() bool {
	return c.booleans[errorsEnableStackTraceVar]
}

// GetIPFSPeerSetting defines if we use IPFS bootstrap peers for discovery or just our own
func (c *Config) GetIPFSPeerSetting() bool {
	return c.booleans[ipfsPeerVar]
}
