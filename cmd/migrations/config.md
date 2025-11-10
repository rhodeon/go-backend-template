# Environment Variables

## Config

 - `MIGRATIONS_ENVIRONMENT` (default: `Development`) - Environment specifies the current running environment of the database migrations.
 - `MIGRATIONS_DEBUG_MODE` (default: `false`) - DebugMode enables/disables detailed debugging output.
 - `MIGRATIONS_DB_ADDR` (default: `localhost`) - Host address of the database to connect to.
 - `MIGRATIONS_DB_PORT` (default: `5432`) - Port of the database to connect to.
 - `MIGRATIONS_DB_USER` (**required**) - User for the database authentication.
 - `MIGRATIONS_DB_PASS` (**required**) - Pass (password) for the database authentication.
 - `MIGRATIONS_DB_NAME` (**required**) - Name of the database to connect to.
 - `MIGRATIONS_DB_SSL_MODE` (default: `disable`) - SslMode of the database connection.
 - `MIGRATIONS_DB_MAX_CONNECTIONS` (default: `25`) - MaxConns is the maximum connections that can be created by the database connection pool.
 - `MIGRATIONS_DB_MAX_CONNECTION_LIFETIME` (default: `2h`) - MaxConnLifetime is the duration since creation after which a connection will be automatically closed.
 - `MIGRATIONS_DB_MAX_CONNECTION_IDLE_TIME` (default: `5m`) - MaxConnIdleTime is the duration after which an idle connection will be automatically closed.

