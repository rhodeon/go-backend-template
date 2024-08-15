package internal

import (
	"github.com/caarlos0/env/v11"
	_ "github.com/joho/godotenv/autoload"
)

type Config struct {
	Environment string `env:"ENVIRONMENT" envDefault:"Development"`
	DebugMode   bool   `env:"DEBUG_MODE" envDefault:"false"`
	HttpPort    int    `env:"HTTP_PORT,required"`
	BaseUrl     string `env:"BASE_URL"`
	Database    DatabaseConfig
}

func ParseConfig() *Config {
	cfg := env.Must(env.ParseAsWithOptions[Config](env.Options{
		Prefix: "API_",
	}))
	return &cfg
}

type DatabaseConfig struct {
	Addr string `env:"DB_ADDR" envDefault:"localhost:5432"`
	User string `env:"DB_USER,required"`
	Pass string `env:"DB_PASS,required"`
	Name string `env:"DB_NAME,required"`
}
