package internal

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pkg/errors"
	"log/slog"
	"strings"
	"time"
)

// ConnectToDb establishes a connection to the given Postgres database and returns a connection pool to be used for further access.
func ConnectToDb(cfg *Config, logger *slog.Logger) (*pgxpool.Pool, error) {
	dbCfg := cfg.Database

	dsn := fmt.Sprintf(
		"user=%s password=%s host=%s port=%s dbname=%s sslmode=%s",
		dbCfg.User, dbCfg.Pass, dbCfg.Host, dbCfg.Port, dbCfg.Name, dbCfg.SslMode,
	)

	pgxPoolCfg, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, errors.Wrap(err, "unable to parse pgx pool config")
	}
	pgxPoolCfg.ConnConfig.Tracer = newDbTracer(logger, cfg.DebugMode)
	pgxPoolCfg.MaxConns = dbCfg.MaxConns
	pgxPoolCfg.MaxConnLifetime = dbCfg.MaxConnLifetime
	pgxPoolCfg.MaxConnIdleTime = dbCfg.MaxConnIdleTime

	connPool, err := pgxpool.NewWithConfig(context.Background(), pgxPoolCfg)
	if err != nil {
		return nil, errors.Wrap(err, "unable to create db connection pool")
	}

	// The database is pinged to ensure the connection was established.
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := connPool.Ping(ctx); err != nil {
		return nil, errors.Wrap(err, "error pinging postgres")
	}

	return connPool, nil
}

// dbTracer implements the pgx.QueryTracer interface to provider debugging and tracing capabilities for queries.
type dbTracer struct {
	logger    *slog.Logger
	debugMode bool
}

func newDbTracer(logger *slog.Logger, debugMode bool) dbTracer {
	return dbTracer{
		logger:    logger,
		debugMode: debugMode,
	}
}

// TraceQueryStart logs each database query triggered when in debug mode.
// The debugMode is used rather than simply calling the logger.Debug in order to
// prevent wasting time on formatting the output when not in debug mode.
func (d dbTracer) TraceQueryStart(ctx context.Context, conn *pgx.Conn, data pgx.TraceQueryStartData) context.Context {
	if d.debugMode {
		// Prevent pollution of logs with transaction boundary commands.
		if strings.EqualFold(data.SQL, "begin") ||
			strings.EqualFold(data.SQL, "rollback") ||
			strings.EqualFold(data.SQL, "commit") ||
			strings.EqualFold(data.SQL, "end") {
			return ctx
		}
		var render = "\n"
		for i, arg := range data.Args {
			switch arg.(type) {
			case []byte:
				render += fmt.Sprintf("$%d:\t%s\n", i+1, arg)
			default:
				render += fmt.Sprintf("$%d:\t%v\n", i+1, arg)
			}
		}

		d.logger.Debug(fmt.Sprintf("Executing db query:\n%s\nargs:%s", data.SQL, render))
	}

	return ctx
}

func (d dbTracer) TraceQueryEnd(ctx context.Context, conn *pgx.Conn, data pgx.TraceQueryEndData) {}
