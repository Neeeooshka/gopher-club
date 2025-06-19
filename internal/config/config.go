package config

type Config struct {
	ServerAddress  string `env:"RUN_ADDRESS"`
	AccrualAddress string `env:"ACCRUAL_SYSTEM_ADDRESS"`
	DB             string `env:"DATABASE_URI"`
}

func NewConfig() Config {
	return Config{}
}
