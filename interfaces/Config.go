package interfaces

// Config is an interface to viper
type Config interface {
	ReadConfig(configPath string)
	Get(variable string) interface{}
	GetString(variable string) string
	GetUint(variable string) uint
}
