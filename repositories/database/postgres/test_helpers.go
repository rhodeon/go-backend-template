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

	"github.com/go-errors/errors"
	"github.com/google/uuid"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
	"github.com/testcontainers/testcontainers-go"
	tcpostgres "github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

const templateDbName = "test_template_db"

// testConfig represents the configuration of the database set up in the container.
// The port is determined after the container has been created.
var testConfig = &database.Config{
	Host:     "localhost",
	User:     "test_user",
	Pass:     "test_pass",
	Name:     "test_initial_db",
	SslMode:  "disable",
	MaxConns: 1,
}

// SetupTestContainer establishes a Postgres instance in a container to be used for testing.
// During this, a template database is also set up. All tests which need a database connection create one by first cloning the template database.
// This is much faster than running the migrations for each test, especially when there's a lot of data involved.
func SetupTestContainer(ctx context.Context, image string, projectRootDir string) (*tcpostgres.PostgresContainer, error) {
	goose.SetLogger(goose.NopLogger())

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

	if err := setupTemplateDb(ctx, projectRootDir); err != nil {
		return nil, err
	}

	return postgresContainer, nil
}

// setupTemplateDb generates and populates the template database which all tests with a dependency of the database clone.
func setupTemplateDb(ctx context.Context, projectRootDir string) error {
	initialDb, closeInitialDb, err := database.Connect(ctx, testConfig, false)
	if err != nil {
		return errors.Errorf("connecting to initial database: %w", err)
	}

	if _, err := initialDb.Pool().Exec(ctx, fmt.Sprintf("CREATE DATABASE %s", templateDbName)); err != nil {
		return errors.Errorf("creating template database: %w", err)
	}

	closeInitialDb()

	templateDbConfig := *testConfig
	templateDbConfig.Name = templateDbName
	templateDb, closeTemplateDb, err := database.Connect(ctx, &templateDbConfig, false)
	if err != nil {
		return errors.Errorf("connecting to test database: %w", err)
	}

	if err := migrateSchema(ctx, templateDb.Pool().Config().ConnString(), projectRootDir); err != nil {
		return errors.Errorf("migrating schema: %w", err)
	}

	if err := seedData(ctx, templateDb, projectRootDir); err != nil {
		return errors.Errorf("seeding data: %w", err)
	}

	closeTemplateDb()

	return nil
}

// ConnectTestDb connects to the initial database, creates a new database using the seeded templated and returns a connection to said database.
func ConnectTestDb(ctx context.Context) (*database.Db, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	initialDb, closeInitialDb, err := database.Connect(ctx, testConfig, false)
	if err != nil {
		return nil, errors.Errorf("connecting to initial database: %w", err)
	}

	// A new random database is created for each test to prevent conflicts.
	dbName := "test_" + strings.ReplaceAll(uuid.NewString(), "-", "")

	if _, err := initialDb.Pool().Exec(ctx, fmt.Sprintf("CREATE DATABASE %s TEMPLATE %s", dbName, templateDbName)); err != nil {
		return nil, errors.Errorf("creating test database: %w", err)
	}

	closeInitialDb()

	// The database pool of the test will be automatically dropped when the test ends, so there's no need to add extra complexity by trying to close it.
	dbConfig := *testConfig
	dbConfig.Name = dbName
	db, _, err := database.Connect(ctx, &dbConfig, false)
	if err != nil {
		return nil, errors.Errorf("connecting to test database: %w", err)
	}

	return db, nil
}

func migrateSchema(ctx context.Context, connString string, projectRootDir string) error {
	db, err := sql.Open("pgx/v5", connString)
	if err != nil {
		return errors.Errorf("opening database: %w", err)
	}
	defer db.Close() //nolint:errcheck

	migrationsPath := filepath.Join(projectRootDir, "cmd", "migrations", "schema")
	if err := goose.UpContext(ctx, db, migrationsPath); err != nil {
		return errors.Errorf("applying up migrations: %w", err)
	}

	return nil
}

// seedData parses the individual seed files from the test data (each named after the corresponding table) and populates their tables.
func seedData(ctx context.Context, db *database.Db, projectRootDir string) error {
	seedsPath := filepath.Join(projectRootDir, "testdata", "database_seeds")
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

			// Although this isn't directly stated in the Postgres docs (as far as I can tell), inserting an auto-incrementing row
			// with an explicit ID rather than depending on the generated default prevents the next value from being updated by Postgres internally.
			// So this needs to be done manually to prevent unique key conflicts from any future operations that insert into the affected table.
			// More details: https://stackoverflow.com/a/9091979
			// This can be removed entirely if all the primary keys are non-incrementing values, like UUIDs.
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

	// Unfortunately, re-enabling doesn't retroactively enforce the integrity of the already-seeded data,
	// but it guards against any future changes during the test.
	if _, err := dbTx.Exec(ctx, "SET session_replication_role = DEFAULT"); err != nil {
		return errors.Errorf("re-enabling triggers: %w", err)
	}

	if err := commit(ctx); err != nil {
		return errors.Errorf("committing transaction: %w", err)
	}

	return nil
}
