# Environment Variables

## Config

 - `API_ENVIRONMENT` (default: `development`) - Environment specifies the current running environment of the API.
 - `API_DEBUG_MODE` (default: `false`) - DebugMode enables/disables detailed debugging output.
 - `API_DB_ADDR` (default: `localhost`) - Host address of the database to connect to.
 - `API_DB_PORT` (default: `5432`) - Port of the database to connect to.
 - `API_DB_USER` (**required**) - User for the database authentication.
 - `API_DB_PASS` (**required**) - Pass (password) for the database authentication.
 - `API_DB_NAME` (**required**) - Name of the database to connect to.
 - `API_DB_SSL_MODE` (default: `disable`) - SslMode of the database connection.
 - `API_DB_MAX_CONNECTIONS` (default: `25`) - MaxConns is the maximum connections that can be created by the database connection pool.
 - `API_DB_MAX_CONNECTION_LIFETIME` (default: `2h`) - MaxConnLifetime is the duration since creation after which a connection will be automatically closed.
 - `API_DB_MAX_CONNECTION_IDLE_TIME` (default: `5m`) - MaxConnIdleTime is the duration after which an idle connection will be automatically closed.
 - `API_CACHE_HOST` (default: `localhost`) - Host specifies the address of the cache server.
 - `API_CACHE_PORT` (default: `6379`) - Port defines the port number on which the cache server will listen.
 - `API_CACHE_PASSWORD` - Password of the cache server.
 - `API_CACHE_DATABASE` (default: `0`) - Database specifies the cache database number to connect to.
 - `API_SERVER_HTTP_PORT` (**required**) - HttpPort defines the port number on which the HTTP server will listen for incoming connections.
 - `API_SERVER_BASE_URL` - BaseUrl specifies the base URL used for constructing server-related endpoints.
 - `API_SERVER_IDLE_TIMEOUT` (default: `1m`) - IdleTimeout is the duration the server will wait for the next request before closing idle connections when keep-alives are enabled.
 - `API_SERVER_READ_TIMEOUT` (default: `5s`) - ReadTimeout specifies the maximum duration for reading an entire request, including the body.
 - `API_SERVER_WRITE_TIMEOUT` (default: `15s`) - WriteTimeout defines the maximum duration for writing a response before timing out.
 - `API_SERVER_REQUEST_TIMEOUT` (default: `10s`) - RequestTimeout specifies the maximum duration for handlers to run.
It should be lower than WriteTimeout as no response will be returned if a handler's runtime exceeds the write timeout.
 - `API_SERVER_SHUTDOWN_TIMEOUT` (default: `30s`) - ShutdownTimeout specifies the duration the server will wait to wrap up active connections and background operations for a graceful shutdown.
 - `API_AUTH_JWT_ISSUER` (default: `go-backend-template`) - JwtIssuer specifies the issuer of the JWT, determining the entity responsible for generating the token.
 - `API_AUTH_JWT_ACCESS_TOKEN_SECRET` - JwtAccessTokenSecret is the secret key used to sign and validate JWT access tokens.
 - `API_AUTH_JWT_REFRESH_TOKEN_SECRET` - JwtRefreshTokenSecret is the secret key used to sign and validate JWT refresh tokens.
 - `API_AUTH_JWT_ACCESS_TOKEN_DURATION` (default: `1h`) - JwtAccessTokenDuration defines the lifespan of JWT access tokens before they expire.
 - `API_AUTH_JWT_REFRESH_TOKEN_DURATION` (default: `12h`) - JwtRefreshTokenDuration specifies the duration for which a JWT refresh token remains valid before expiration.
 - `API_AUTH_OTP_DURATION` (default: `30s`) - OtpDuration defines the time duration for which an OTP remains valid before expiring.
 - `API_SMTP_HOST` (default: `localhost`) - Host specifies the SMTP server address.
 - `API_SMTP_PORT` (default: `587`) - Port specifies the port number for the SMTP server.
 - `API_SMTP_USER` - User specifies the username required for SMTP authentication.
 - `API_SMTP_PASSWORD` - Password specifies the password required for SMTP authentication.
 - `API_SMTP_SENDER` - Sender defines the email address used as the sender in SMTP communications.
 - `API_OTEL_SERVICE_NAME` (default: `go-backend-template`) - ServiceName defines the name of the service used for observability and telemetry.
 - `API_OTEL_OTLP_GRPC_HOST` (default: `localhost`) - OtlpGrpcHost defines the host address for the OTLP gRPC exporter.
 - `API_OTEL_OTLP_GRPC_PORT` (default: `4317`) - OtlpGrpcPort specifies the port number for the OTLP gRPC exporter.
 - `API_OTEL_OTLP_SECURE_CONNECTION` (default: `false`) - OtlpSecureConnection determines if a secure (TLS) connection should be used for OTLP communication.

