package config

type Config struct {
	HTTPPort string   `mapstructure:"HTTP_PORT"`
	BaseURL  string   `mapstructure:"BASE_URL"`
	LogLevel string   `mapstructure:"LOG_LEVEL"`
	DB       DBConfig `mapstructure:"DB"`
}

type DBConfig struct {
	DSN string `mapstructure:"DSN"`
}
