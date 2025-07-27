package main

import (
	"context"
	"flag"
	"fmt"
	"os"

	"github.com/rhodeon/go-backend-template/cmd/migrations/internal"
	"github.com/rhodeon/go-backend-template/internal/database"
	"github.com/rhodeon/go-backend-template/internal/log"

	"github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
)

const usageText = `This program runs command on the db. Supported commands are:
  - up                   Migrate the DB to the most recent version available
  - up-by-one            Migrate the DB up by 1
  - up-to VERSION        Migrate the DB to a specific VERSION
  - down                 Roll back the version by 1
  - down-to VERSION      Roll back to a specific VERSION
  - redo                 Re-run the latest migration
  - reset                Roll back all migrations
  - status               Dump the migration status for the current DB
  - version              Print the current version of the database
  - create NAME [sql|go] Creates new migration file with the current timestamp
  - fix                  Apply sequential ordering to migrations
  - validate             Check migration files without running them

Usage:
  go run *.go <command> [args]`

func main() {
	flag.Usage = usage
	flag.Parse()

	// The provided flags are parsed for options to Goose.
	flagArgs := flag.Args()
	if len(flagArgs) == 0 {
		usage()
	}

	command := flagArgs[0]
	gooseArgs := []string{}

	if len(flagArgs) > 1 {
		gooseArgs = append(gooseArgs, flagArgs[1:]...)
	}

	// By default, Goose creates migrates in .go files if no format is specified.
	// This ensures it creates a .sql file instead.
	if command == "create" {
		if len(gooseArgs) == 1 {
			gooseArgs = []string{gooseArgs[0], "sql"}
		}
	}

	cfg := internal.ParseConfig()
	logger := log.NewLogger(cfg.DebugMode)

	dbConfig := database.Config(cfg.Database)
	dbPool, err := database.Connect(&dbConfig, logger, cfg.DebugMode)
	if err != nil {
		panic(err)
	}

	db := stdlib.OpenDBFromPool(dbPool)
	if err := goose.RunWithOptionsContext(context.Background(), command, db, "./schema", gooseArgs, goose.WithAllowMissing()); err != nil {
		panic(err)
	}
}

func usage() {
	fmt.Println(usageText)
	flag.PrintDefaults()
	os.Exit(0)
}
