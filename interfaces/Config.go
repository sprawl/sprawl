package interfaces

// Config is an interface to viper
type Config interface {
	ReadConfig(configPath string)
	GetDatabasePath() string
	GetExternalIP() string
	GetLogLevel() string
	GetLogFormat() string
	GetP2PPort() string
	GetRPCPort() string
	GetInMemoryDatabaseSetting() bool
	GetNATPortMapSetting() bool
	GetRelaySetting() bool
	GetAutoRelaySetting() bool
	GetDebugSetting() bool
	GetStackTraceSetting() bool
	GetIPFSPeerSetting() bool
}
