package internal

import (
	"time"

	"github.com/caarlos0/env/v11"
)

const configPrefix = "API_"

//go:generate envdoc -output "../config.md" -types "Config" -env-prefix "API_" -template "../../../templates/configdoc.md.go.tmpl" -title "API"
//go:generate envdoc -output "../config.env" -types "Config" -env-prefix "API_" -template "../../../templates/configdoc.dotenv.go.tmpl" -title "API"
type Config struct {
	// Environment specifies the current running environment of the API.
	Environment string `env:"ENVIRONMENT" envDefault:"development"`

	// DebugMode enables/disables detailed debugging output.
	DebugMode bool `env:"DEBUG_MODE" envDefault:"false"`

	Server   ServerConfig
	Database DatabaseConfig
	Cache    CacheConfig
	Auth     AuthConfig
	Smtp     SmtpConfig
	Otel     OtelConfig
}

func ParseConfig() *Config {
	cfg := env.Must(env.ParseAsWithOptions[Config](env.Options{
		Prefix: configPrefix,
	}))
	return &cfg
}

type ServerConfig struct {
	// HttpPort defines the port number on which the HTTP server will listen for incoming connections.
	HttpPort int `env:"SERVER_HTTP_PORT,notEmpty"`

	// BaseUrl specifies the base URL used for constructing server-related endpoints.
	BaseUrl string `env:"SERVER_BASE_URL"`

	// The timeout defaults are those used by https://autostrada.dev and should be modified according to real usage if needed.

	// IdleTimeout is the duration the server will wait for the next request before closing idle connections when keep-alives are enabled.
	IdleTimeout time.Duration `env:"SERVER_IDLE_TIMEOUT" envDefault:"1m"`

	// ReadTimeout specifies the maximum duration for reading an entire request, including the body.
	ReadTimeout time.Duration `env:"SERVER_READ_TIMEOUT" envDefault:"5s"`

	// WriteTimeout defines the maximum duration for writing a response before timing out.
	WriteTimeout time.Duration `env:"SERVER_WRITE_TIMEOUT" envDefault:"15s"`

	// RequestTimeout specifies the maximum duration for handlers to run.
	// It should be lower than WriteTimeout as no response will be returned if a handler's runtime exceeds the write timeout.
	RequestTimeout time.Duration `env:"SERVER_REQUEST_TIMEOUT" envDefault:"10s"`

	// ShutdownTimeout specifies the duration the server will wait to wrap up active connections and background operations for a graceful shutdown.
	ShutdownTimeout time.Duration `env:"SERVER_SHUTDOWN_TIMEOUT" envDefault:"30s"`
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

type CacheConfig struct {
	// Host specifies the address of the cache server.
	Host string `env:"CACHE_HOST" envDefault:"localhost"`

	// Port defines the port number on which the cache server will listen.
	Port int `env:"CACHE_PORT" envDefault:"6379"`

	// Password of the cache server.
	Password string `env:"CACHE_PASSWORD" envDefault:""`

	// Database specifies the cache database number to connect to.
	Database int `env:"CACHE_DATABASE" envDefault:"0"`
}

type AuthConfig struct {
	// JwtIssuer specifies the issuer of the JWT, determining the entity responsible for generating the token.
	JwtIssuer string `env:"AUTH_JWT_ISSUER" envDefault:"go-backend-template"`

	// JwtAccessTokenSecret is the secret key used to sign and validate JWT access tokens.
	JwtAccessTokenSecret string `env:"AUTH_JWT_ACCESS_TOKEN_SECRET"`

	// JwtRefreshTokenSecret is the secret key used to sign and validate JWT refresh tokens.
	JwtRefreshTokenSecret string `env:"AUTH_JWT_REFRESH_TOKEN_SECRET"`

	// JwtAccessTokenDuration defines the lifespan of JWT access tokens before they expire.
	JwtAccessTokenDuration time.Duration `env:"AUTH_JWT_ACCESS_TOKEN_DURATION" envDefault:"1h"`

	// JwtRefreshTokenDuration specifies the duration for which a JWT refresh token remains valid before expiration.
	JwtRefreshTokenDuration time.Duration `env:"AUTH_JWT_REFRESH_TOKEN_DURATION" envDefault:"12h"`

	// OtpDuration defines the time duration for which an OTP remains valid before expiring.
	OtpDuration time.Duration `env:"AUTH_OTP_DURATION" envDefault:"30s"`
}

type SmtpConfig struct {
	// Host specifies the SMTP server address.
	Host string `env:"SMTP_HOST" envDefault:"localhost"`

	// Port specifies the port number for the SMTP server.
	Port int `env:"SMTP_PORT" envDefault:"587"`

	// User specifies the username required for SMTP authentication.
	User string `env:"SMTP_USER"`

	// Password specifies the password required for SMTP authentication.
	Password string `env:"SMTP_PASSWORD"`

	// Sender defines the email address used as the sender in SMTP communications.
	Sender string `env:"SMTP_SENDER"`
}

type OtelConfig struct {
	// ServiceName defines the name of the service used for observability and telemetry.
	ServiceName string `env:"OTEL_SERVICE_NAME" envDefault:"go-backend-template"`

	// OtlpGrpcHost defines the host address for the OTLP gRPC exporter.
	OtlpGrpcHost string `env:"OTEL_OTLP_GRPC_HOST" envDefault:"localhost"`

	// OtlpGrpcPort specifies the port number for the OTLP gRPC exporter.
	OtlpGrpcPort int `env:"OTEL_OTLP_GRPC_PORT" envDefault:"4317"`

	// OtlpSecureConnection determines if a secure (TLS) connection should be used for OTLP communication.
	OtlpSecureConnection bool `env:"OTEL_OTLP_SECURE_CONNECTION" envDefault:"false"`
}
