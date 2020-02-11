package interfaces

// Config is an interface to viper
type Config interface {
	AddString(key string)
	AddBoolean(key string)
	AddUint(key string)
	AddStringE(key string) error
	AddBooleanE(key string) error
	AddUintE(key string) error
	ReadConfig(configPath string)
	GetDatabasePath() string
	GetExternalIP() string
	GetLogLevel() string
	GetLogFormat() string
	GetP2PPort() uint
	GetRPCPort() uint
	GetWebsocketPort() uint
	GetWebsocketEnable() bool
	GetInMemoryDatabaseSetting() bool
	GetNATPortMapSetting() bool
	GetRelaySetting() bool
	GetAutoRelaySetting() bool
	GetDebugSetting() bool
	GetStackTraceSetting() bool
	GetIPFSPeerSetting() bool
}
