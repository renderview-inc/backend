package config

type ClickHouseConfig struct {
	Table string `mapstructure:"table"`
}

type LogConfig struct {
	Logger struct {
		Level      string           `mapstructure:"level"`
		ClickHouse ClickHouseConfig `mapstructure:"clickhouse"`
	} `mapstructure:"logger"`
}
