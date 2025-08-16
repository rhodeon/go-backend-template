package postgres

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

	"github.com/rhodeon/go-backend-template/internal/database"
	"github.com/rhodeon/go-backend-template/utils/testutils/temp"

	"github.com/go-errors/errors"
	"github.com/google/uuid"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
	"github.com/testcontainers/testcontainers-go"
	tcpostgres "github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

var testConfig *database.Config = &database.Config{}

// SetupTestContainer establishes a Postgres instance in a container to be used for testing.
func SetupTestContainer(ctx context.Context, image string) (*tcpostgres.PostgresContainer, error) {
	goose.SetLogger(goose.NopLogger())

	testConfig = &database.Config{
		Host:     "localhost",
		User:     "test_user",
		Pass:     "test_pass",
		Name:     "test_db",
		SslMode:  "disable",
		MaxConns: 1,
	}

	postgresContainer, err := tcpostgres.Run(ctx,
		image,
		tcpostgres.WithDatabase(testConfig.Name),
		tcpostgres.WithUsername(testConfig.User),
		tcpostgres.WithPassword(testConfig.Pass),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(5*time.Second)),
	)
	if err != nil {
		return nil, errors.Errorf("creating Postgres container instance: %w", err)
	}

	mappedPort, err := postgresContainer.MappedPort(ctx, "5432")
	if err != nil {
		return nil, errors.Errorf("getting mapped Postgres container ports: %w", err)
	}

	testConfig.Port = mappedPort.Port()

	if err = postgresContainer.Start(ctx); err != nil {
		return nil, errors.Errorf("starting Postgres container: %w", err)
	}

	return postgresContainer, nil
}

func ConnectDb(ctx context.Context) (*database.Db, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	// A new random database is created for each test to prevent conflicts in operations.
	dbName := "test_" + strings.ReplaceAll(uuid.NewString(), "-", "")

	initialDb, closeInitialDb, err := database.Connect(ctx, testConfig, false)
	if err != nil {
		return nil, errors.Errorf("connecting to initial database: %w", err)
	}

	if _, err := initialDb.Pool().Exec(ctx, fmt.Sprintf("CREATE DATABASE %s", dbName)); err != nil {
		return nil, errors.Errorf("creating test database: %w", err)
	}
	closeInitialDb()

	testConfig.Name = dbName

	// The database pool of the test will be short-lived and automatically dropped when the test ends,
	// so there's no need to add extra complexity by trying to close it.
	db, _, err := database.Connect(ctx, testConfig, false)
	if err != nil {
		return nil, errors.Errorf("connecting to test database: %w", err)
	}

	if err := migrateSchema(db.Pool().Config().ConnConfig.ConnString()); err != nil {
		return nil, errors.Errorf("migrating schema: %w", err)
	}

	if err := seedData(ctx, db); err != nil {
		return nil, errors.Errorf("seeding data: %w", err)
	}

	return db, nil
}

func migrateSchema(connString string) error {
	db, err := sql.Open("pgx/v5", connString)
	if err != nil {
		return errors.Errorf("opening database: %w", err)
	}

	migrationsPath := filepath.Join(temp.ProjectRootDir, "cmd", "migrations", "schema")
	if err := goose.Up(db, migrationsPath); err != nil {
		return errors.Errorf("applying up migrations: %w", err)
	}

	return nil
}

// seedData parses the individual seed files in the hub-api repo, each named after a table,
// and populates their tables.
func seedData(ctx context.Context, db *database.Db) error {
	seedsPath := filepath.Join(temp.ProjectRootDir, "testdata", "database_seeds")
	seedFiles, err := os.ReadDir(seedsPath)
	if err != nil {
		return errors.Errorf("reading database seeds directory: %w", err)
	}

	dbTx, commit, rollback, err := db.BeginTx(ctx)
	if err != nil {
		return errors.Errorf("beginning transaction to seed data: %w", err)
	}
	defer rollback(ctx)

	// Triggers are disabled to prevent foreign key constraints from raising an error when seeding data out of order.
	if _, err := dbTx.Exec(ctx, "SET session_replication_role = replica"); err != nil {
		return errors.Errorf("disabling triggers: %w", err)
	}

	for _, entry := range seedFiles {
		if !entry.IsDir() && strings.HasSuffix(entry.Name(), ".json") {
			tableName := strings.TrimSuffix(entry.Name(), ".json")

			data, err := os.ReadFile(path.Join(seedsPath, entry.Name()))
			if err != nil {
				return errors.Errorf("reading seed file for table %q: %w", tableName, err)
			}

			insertionQuery := fmt.Sprintf("INSERT INTO %s OVERRIDING SYSTEM VALUE SELECT * FROM json_populate_recordset(NULL::%s, $1)", tableName, tableName)
			if _, err := dbTx.Exec(ctx, insertionQuery, data); err != nil {
				return errors.Errorf("inserting data for table %q: %w", tableName, err)
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
				return errors.Errorf("getting sequence name for table %q: %w", tableName, err)
			}

			columnName := "id"
			// With this, the next insertion to the table will start from a valid ID.
			setValQuery := fmt.Sprintf("SELECT setval('%s', (SELECT MAX(%s) FROM %s))", sequenceName, columnName, tableName)
			if _, err := dbTx.Exec(ctx, setValQuery); err != nil {
				return errors.Errorf("updating starting ID value for table %q: %w", tableName, err)
			}
		}
	}

	if _, err := dbTx.Exec(ctx, "SET session_replication_role = DEFAULT"); err != nil {
		return errors.Errorf("re-enabling triggers: %w", err)
	}

	if err := commit(ctx); err != nil {
		return errors.Errorf("committing transaction: %w", err)
	}

	return nil
}
