# API

| Name | Description | Default | Attributes |
|------|-------------|---------|------------|
| `API_ENVIRONMENT` | Environment specifies the current running environment of the API. | `development` |  |
| `API_DEBUG_MODE` | DebugMode enables/disables detailed debugging output. | `false` |  |
| `API_SERVER_HTTP_PORT` | HttpPort defines the port number on which the HTTP server will listen for incoming connections. |  | `REQUIRED` |
| `API_SERVER_BASE_URL` | BaseUrl specifies the base URL used for constructing server-related endpoints. |  |  |
| `API_SERVER_IDLE_TIMEOUT` | IdleTimeout is the duration the server will wait for the next request before closing idle connections when keep-alives are enabled. | `1m` |  |
| `API_SERVER_READ_TIMEOUT` | ReadTimeout specifies the maximum duration for reading an entire request, including the body. | `5s` |  |
| `API_SERVER_WRITE_TIMEOUT` | WriteTimeout defines the maximum duration for writing a response before timing out. | `15s` |  |
| `API_SERVER_REQUEST_TIMEOUT` | RequestTimeout specifies the maximum duration for handlers to run.<br>It should be lower than WriteTimeout as no response will be returned if a handler's runtime exceeds the write timeout. | `10s` |  |
| `API_SERVER_SHUTDOWN_TIMEOUT` | ShutdownTimeout specifies the duration the server will wait to wrap up active connections and background operations for a graceful shutdown. | `30s` |  |
| `API_DB_ADDR` | Host address of the database to connect to. | `localhost` |  |
| `API_DB_PORT` | Port of the database to connect to. | `5432` |  |
| `API_DB_USER` | User for the database authentication. |  | `REQUIRED` |
| `API_DB_PASS` | Pass (password) for the database authentication. |  | `REQUIRED` |
| `API_DB_NAME` | Name of the database to connect to. |  | `REQUIRED` |
| `API_DB_SSL_MODE` | SslMode of the database connection. | `disable` |  |
| `API_DB_MAX_CONNECTIONS` | MaxConns is the maximum connections that can be created by the database connection pool. | `25` |  |
| `API_DB_MAX_CONNECTION_LIFETIME` | MaxConnLifetime is the duration since creation after which a connection will be automatically closed. | `2h` |  |
| `API_DB_MAX_CONNECTION_IDLE_TIME` | MaxConnIdleTime is the duration after which an idle connection will be automatically closed. | `5m` |  |
| `API_CACHE_HOST` | Host specifies the address of the cache server. | `localhost` |  |
| `API_CACHE_PORT` | Port defines the port number on which the cache server will listen. | `6379` |  |
| `API_CACHE_PASSWORD` | Password of the cache server. |  |  |
| `API_CACHE_DATABASE` | Database specifies the cache database number to connect to. | `0` |  |
| `API_AUTH_JWT_ISSUER` | JwtIssuer specifies the issuer of the JWT, determining the entity responsible for generating the token. | `go-backend-template` |  |
| `API_AUTH_JWT_ACCESS_TOKEN_SECRET` | JwtAccessTokenSecret is the secret key used to sign and validate JWT access tokens. |  |  |
| `API_AUTH_JWT_REFRESH_TOKEN_SECRET` | JwtRefreshTokenSecret is the secret key used to sign and validate JWT refresh tokens. |  |  |
| `API_AUTH_JWT_ACCESS_TOKEN_DURATION` | JwtAccessTokenDuration defines the lifespan of JWT access tokens before they expire. | `1h` |  |
| `API_AUTH_JWT_REFRESH_TOKEN_DURATION` | JwtRefreshTokenDuration specifies the duration for which a JWT refresh token remains valid before expiration. | `12h` |  |
| `API_AUTH_OTP_DURATION` | OtpDuration defines the time duration for which an OTP remains valid before expiring. | `30s` |  |
| `API_SMTP_HOST` | Host specifies the SMTP server address. | `localhost` |  |
| `API_SMTP_PORT` | Port specifies the port number for the SMTP server. | `587` |  |
| `API_SMTP_USER` | User specifies the username required for SMTP authentication. |  |  |
| `API_SMTP_PASSWORD` | Password specifies the password required for SMTP authentication. |  |  |
| `API_SMTP_SENDER` | Sender defines the email address used as the sender in SMTP communications. |  |  |
| `API_OTEL_SERVICE_NAME` | ServiceName defines the name of the service used for observability and telemetry. | `go-backend-template` |  |
| `API_OTEL_OTLP_GRPC_HOST` | OtlpGrpcHost defines the host address for the OTLP gRPC exporter. | `localhost` |  |
| `API_OTEL_OTLP_GRPC_PORT` | OtlpGrpcPort specifies the port number for the OTLP gRPC exporter. | `4317` |  |
| `API_OTEL_OTLP_SECURE_CONNECTION` | OtlpSecureConnection determines if a secure (TLS) connection should be used for OTLP communication. | `false` |  |
