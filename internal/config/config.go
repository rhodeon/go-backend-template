package config

import (
	"github.com/caarlos0/env/v11"
	_ "github.com/joho/godotenv/autoload"
	"time"
)

type Config struct {
	Environment string `env:"ENVIRONMENT" envDefault:"Development"`
	DebugMode   bool   `env:"DEBUG_MODE" envDefault:"false"`
	HttpPort    int    `env:"HTTP_PORT,required"`
	BaseUrl     string `env:"BASE_URL"`
	Database    DatabaseConfig
}

func Parse(prefix string) *Config {
	cfg := env.Must(env.ParseAsWithOptions[Config](env.Options{
		Prefix: prefix + "_",
	}))
	return &cfg
}

type DatabaseConfig struct {
	Host    string `env:"DB_ADDR" envDefault:"localhost"`
	Port    string `env:"DB_ADDR" envDefault:"5432"`
	User    string `env:"DB_USER,required"`
	Pass    string `env:"DB_PASS,required"`
	Name    string `env:"DB_NAME,required"`
	SslMode string `env:"DB_SSL_MODE" envDefault:"disable"`

	// The connection defaults are those used by https://autostrada.dev and should be modified according to real usage if needed.
	MaxConns        int32         `env:"DB_MAX_CONNECTIONS" envDefault:"25"`
	MaxConnLifetime time.Duration `env:"DB_MAX_CONNECTION_LIFETIME" envDefault:"2h"`
	MaxConnIdleTime time.Duration `env:"DB_MAX_CONNECTION_IDLE_TIME" envDefault:"5m"`
}
