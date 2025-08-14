package internal

import (
	"time"

	"github.com/caarlos0/env/v11"
	_ "github.com/joho/godotenv/autoload"
)

const configPrefix = "API_"

type Config struct {
	Environment string `env:"ENVIRONMENT" envDefault:"Development"`
	DebugMode   bool   `env:"DEBUG_MODE" envDefault:"false"`
	Database    DatabaseConfig
	Cache       CacheConfig
	Server      ServerConfig
	Auth        AuthConfig
	Smtp        SmtpConfig
}

func ParseConfig() *Config {
	cfg := env.Must(env.ParseAsWithOptions[Config](env.Options{
		Prefix: configPrefix,
	}))
	return &cfg
}

type ServerConfig struct {
	HttpPort int    `env:"HTTP_PORT,required"`
	BaseUrl  string `env:"BASE_URL"`

	// The below timeout defaults are those used by https://autostrada.dev and should be modified according to real usage if needed.
	IdleTimeout     time.Duration `env:"IDLE_TIMEOUT" envDefault:"1m"`
	ReadTimeout     time.Duration `env:"READ_TIMEOUT" envDefault:"5s"`
	WriteTimeout    time.Duration `env:"WRITE_TIMEOUT" envDefault:"15s"`
	ShutdownTimeout time.Duration `env:"SHUTDOWN_TIMEOUT" envDefault:"30s"`

	// RequestTimeout should be lower than WriteTimeout as no response will be returned if the request exceeds the write timeout.
	RequestTimeout time.Duration `env:"REQUEST_TIMEOUT" envDefault:"10s"`
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

type CacheConfig struct {
	Host     string `env:"CACHE_HOST" envDefault:"localhost"`
	Port     int    `env:"CACHE_PORT" envDefault:"6379"`
	Password string `env:"CACHE_PASSWORD" envDefault:""`
	Database int    `env:"CACHE_DATABASE" envDefault:"0"`
}

type AuthConfig struct {
	JwtIssuer               string        `env:"AUTH_JWT_ISSUER" envDefault:"go-backend-template"`
	JwtAccessTokenSecret    string        `env:"AUTH_JWT_ACCESS_TOKEN_SECRET"`
	JwtRefreshTokenSecret   string        `env:"AUTH_JWT_REFRESH_TOKEN_SECRET"`
	JwtAccessTokenDuration  time.Duration `env:"AUTH_JWT_ACCESS_TOKEN_DURATION" envDefault:"1h"`
	JwtRefreshTokenDuration time.Duration `env:"AUTH_JWT_REFRESH_TOKEN_DURATION" envDefault:"12h"`
	OtpDuration             time.Duration `env:"AUTH_OTP_DURATION" envDefault:"30s"`
}

type SmtpConfig struct {
	Host        string `env:"SMTP_HOST" envDefault:"localhost"`
	Port        int    `env:"SMTP_PORT" envDefault:"587"`
	User        string `env:"SMTP_USER,required"`
	Password    string `env:"SMTP_PASSWORD,required"`
	Sender      string `env:"SMTP_SENDER,required"`
	OtpDuration time.Duration
}
