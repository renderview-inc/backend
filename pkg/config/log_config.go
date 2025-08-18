package config

type LogConfig struct {
	Logger struct {
		Level string `mapstructure:"level"`
	} `mapstructure:"logger"`
}
