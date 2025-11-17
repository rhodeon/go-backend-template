package internal

import (
	"time"

	"github.com/caarlos0/env/v11"
)

const configPrefix = "MIGRATIONS_"

//go:generate envdoc -output "../config.md" -types "Config" -env-prefix "MIGRATIONS_" -template "../../../templates/configdoc.md.go.tmpl" -title "Migrations"
//go:generate envdoc -output "../config.env" -types "Config" -env-prefix "MIGRATIONS_" -template "../../../templates/configdoc.dotenv.go.tmpl" -title "MIGRATIONS"
type Config struct {
	// Environment specifies the current running environment of the database migrations.
	Environment string `env:"ENVIRONMENT" envDefault:"development"`

	// DebugMode enables/disables detailed debugging output.
	DebugMode bool `env:"DEBUG_MODE" envDefault:"false"`

	Database DatabaseConfig
}

func ParseConfig() *Config {
	cfg := env.Must(env.ParseAsWithOptions[Config](env.Options{
		Prefix: configPrefix,
	}))
	return &cfg
}

type DatabaseConfig struct {
	// Host address of the database to connect to.
	Host string `env:"DB_ADDR" envDefault:"localhost"`

	// Port of the database to connect to.
	Port string `env:"DB_PORT" envDefault:"5432"`

	// User for the database authentication.
	User string `env:"DB_USER,notEmpty"`

	// Pass (password) for the database authentication.
	Pass string `env:"DB_PASS,notEmpty"`

	// Name of the database to connect to.
	Name string `env:"DB_NAME,notEmpty"`

	// SslMode of the database connection.
	SslMode string `env:"DB_SSL_MODE" envDefault:"disable"`

	// The connection defaults are those used by https://autostrada.dev and should be modified according to real usage if needed.

	// MaxConns is the maximum connections that can be created by the database connection pool.
	MaxConns int32 `env:"DB_MAX_CONNECTIONS" envDefault:"25"`

	// MaxConnLifetime is the duration since creation after which a connection will be automatically closed.
	MaxConnLifetime time.Duration `env:"DB_MAX_CONNECTION_LIFETIME" envDefault:"2h"`

	// MaxConnIdleTime is the duration after which an idle connection will be automatically closed.
	MaxConnIdleTime time.Duration `env:"DB_MAX_CONNECTION_IDLE_TIME" envDefault:"5m"`
}
