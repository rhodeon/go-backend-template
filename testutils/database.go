package testutils

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/pkg/errors"
	"github.com/pressly/goose/v3"
	"github.com/rhodeon/go-backend-template/internal/database"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

func setupDatabaseContainer(ctx context.Context) error {
	goose.SetLogger(goose.NopLogger())

	postgresContainer, err := postgres.Run(ctx,
		config.PostgresContainer,
		postgres.WithDatabase(config.Database.Name),
		postgres.WithUsername(config.Database.User),
		postgres.WithPassword(config.Database.Pass),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(5*time.Second)),
	)
	if err != nil {
		return errors.Wrap(err, "creating Postgres container instance")
	}

	mappedPort, err := postgresContainer.MappedPort(ctx, "5432")
	if err != nil {
		return errors.Wrap(err, "getting mapped Postgres container ports")
	}

	config.Database.Port = mappedPort.Port()

	if err = postgresContainer.Start(ctx); err != nil {
		return errors.Wrap(err, "starting Postgres container")
	}

	return nil
}

func ConnectDb(ctx context.Context) (*pgxpool.Pool, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	// A new random database is created for each test to prevent conflicts in operations.
	dbName := "test_" + strings.ReplaceAll(uuid.NewString(), "-", "")

	dbPool, err := database.Connect(config.Database, nil, false)
	if err != nil {
		return nil, errors.Wrap(err, "connecting to initial database")
	}

	if _, err := dbPool.Exec(ctx, fmt.Sprintf("CREATE DATABASE %s", dbName)); err != nil {
		return nil, errors.Wrap(err, "creating test database")
	}
	dbPool.Close()

	testDbConfig := *config.Database
	testDbConfig.Name = dbName

	dbPool, err = database.Connect(&testDbConfig, nil, false)
	if err != nil {
		return nil, errors.Wrap(err, "connecting to test database")
	}

	if err := migrateSchema(dbPool.Config().ConnConfig.ConnString()); err != nil {
		return nil, errors.Wrap(err, "migrating schema")
	}

	if err := seedData(ctx, dbPool); err != nil {
		return nil, errors.Wrap(err, "seeding data")
	}

	return dbPool, nil
}

func migrateSchema(connString string) error {
	db, err := sql.Open("pgx/v5", connString)
	if err != nil {
		return errors.Wrap(err, "opening database")
	}

	migrationsPath := filepath.Join(projectRootDir, "cmd", "migrations", "schema")
	if err := goose.Up(db, migrationsPath); err != nil {
		return errors.Wrap(err, "applying up migrations")
	}

	return nil
}

// seedData parses the individual seed files in the hub-api repo, each named after a table,
// and populates their tables.
func seedData(ctx context.Context, dbPool *pgxpool.Pool) error {
	seedsPath := filepath.Join(projectRootDir, "testdata", "database_seeds")
	seedFiles, err := os.ReadDir(seedsPath)
	if err != nil {
		return errors.Wrap(err, "reading database seeds directory")
	}

	dbTx, rollback, err := database.BeginTransaction(ctx, dbPool)
	if err != nil {
		return errors.Wrap(err, "unable to begin transaction to seed data")
	}
	defer rollback(ctx, dbTx)

	// Triggers are disabled to prevent foreign key constraints from raising an error when seeding data out of order.
	if _, err := dbTx.Exec(ctx, "SET session_replication_role = replica"); err != nil {
		return errors.Wrap(err, "disabling triggers")
	}

	for _, entry := range seedFiles {
		if !entry.IsDir() && strings.HasSuffix(entry.Name(), ".json") {
			tableName := strings.TrimSuffix(entry.Name(), ".json")

			data, err := os.ReadFile(path.Join(seedsPath, entry.Name()))
			if err != nil {
				return errors.Wrap(err, fmt.Sprintf("reading seed file for table %q", tableName))
			}

			insertionQuery := fmt.Sprintf("INSERT INTO %s OVERRIDING SYSTEM VALUE SELECT * FROM json_populate_recordset(NULL::%s, $1)", tableName, tableName)
			if _, err := dbTx.Exec(ctx, insertionQuery, data); err != nil {
				return errors.Wrap(err, fmt.Sprintf("inserting data for table %q", tableName))
			}

			// Although this isn't directly stated in the Postgres docs (as far as I can tell),
			// inserting an auto-incrementing row with an explicit ID rather than depending on the generated default prevents the next value from being updated by Postgres internally.
			// So this needs to be done manually in order to prevent unique key conflicts from any future operations that insert into the affected table.
			// More details: https://stackoverflow.com/a/9091979
			// This can be removed entirely if all the primary keys are non-incrementing values like UUIDs.
			type identityTable struct {
				Id int `json:"id"`
			}

			err = json.Unmarshal(data, &[]identityTable{})
			if err != nil {
				// If the deserialization fails, the "id" column of the table is not an integer and hence cannot auto-increment. This phase is skipped in such a case.
				continue
			}

			var sequenceName string
			err = dbTx.QueryRow(ctx, "SELECT pg_get_serial_sequence($1, 'id')", tableName).Scan(&sequenceName)
			if err != nil {
				return errors.Wrap(err, fmt.Sprintf("getting sequence name for table %q", tableName))
			}

			columnName := "id"
			// With this, the next insertion to the table will start from a valid ID.
			setValQuery := fmt.Sprintf("SELECT setval('%s', (SELECT MAX(%s) FROM %s))", sequenceName, columnName, tableName)
			if _, err := dbTx.Exec(ctx, setValQuery); err != nil {
				return errors.Wrap(err, fmt.Sprintf("updating starting ID value for table %q", tableName))
			}
		}
	}

	if _, err := dbTx.Exec(ctx, "SET session_replication_role = DEFAULT"); err != nil {
		return errors.Wrap(err, "re-enabling triggers")
	}

	_ = dbTx.Commit(ctx)
	return nil
}
