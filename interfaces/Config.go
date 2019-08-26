package interfaces

type Config interface {
	ReadConfig(configPath string)
	Get(variable string) interface{}
	GetString(variable string) string
	GetUint(variable string) uint
}
