package config

type LogConfig struct {
	Logger struct {
		Level    string `mapstructure:"level"`
		LogsDir  string `mapstructure:"logs-dir"`
		LogsFile string `mapstructure:"logs-file"`
	} `mapstructure:"logger"`
}
