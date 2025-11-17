# Migrations

| Name | Description | Default | Attributes |
|------|-------------|---------|------------|
| `MIGRATIONS_ENVIRONMENT` | Environment specifies the current running environment of the database migrations. | `development` |  |
| `MIGRATIONS_DEBUG_MODE` | DebugMode enables/disables detailed debugging output. | `false` |  |
| `MIGRATIONS_DB_ADDR` | Host address of the database to connect to. | `localhost` |  |
| `MIGRATIONS_DB_PORT` | Port of the database to connect to. | `5432` |  |
| `MIGRATIONS_DB_USER` | User for the database authentication. |  | `REQUIRED` |
| `MIGRATIONS_DB_PASS` | Pass (password) for the database authentication. |  | `REQUIRED` |
| `MIGRATIONS_DB_NAME` | Name of the database to connect to. |  | `REQUIRED` |
| `MIGRATIONS_DB_SSL_MODE` | SslMode of the database connection. | `disable` |  |
| `MIGRATIONS_DB_MAX_CONNECTIONS` | MaxConns is the maximum connections that can be created by the database connection pool. | `25` |  |
| `MIGRATIONS_DB_MAX_CONNECTION_LIFETIME` | MaxConnLifetime is the duration since creation after which a connection will be automatically closed. | `2h` |  |
| `MIGRATIONS_DB_MAX_CONNECTION_IDLE_TIME` | MaxConnIdleTime is the duration after which an idle connection will be automatically closed. | `5m` |  |
