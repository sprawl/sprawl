package interfaces

// Config is an interface to viper
type Config interface {
	AddString(key string)
	AddBoolean(key string)
	AddStringE(key string) error
	AddBooleanE(key string) error
	ReadConfig(configPath string)
	GetDatabasePath() string
	GetExternalIP() string
	GetLogLevel() string
	GetLogFormat() string
	GetP2PPort() string
	GetRPCPort() string
	GetWebsocketPort() string
	GetWebsocketEnable() bool
	GetInMemoryDatabaseSetting() bool
	GetNATPortMapSetting() bool
	GetRelaySetting() bool
	GetAutoRelaySetting() bool
	GetDebugSetting() bool
	GetStackTraceSetting() bool
	GetIPFSPeerSetting() bool
}
