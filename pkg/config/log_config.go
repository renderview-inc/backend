package config

type LogConfig struct {
	Logger struct {
		Level      string `mapstructure:"level"`
		ClickHouse struct {
			DSN   string `mapstructure:"dsn"`
			Table string `mapstructure:"table"`
		} `mapstructure:"clickhouse"`
	} `mapstructure:"logger"`
}
