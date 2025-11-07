#!/bin/bash
set -euo pipefail

psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" --dbname "$POSTGRES_DB" \
  -v api_user="$API_DB_USER" \
  -v api_pass="$API_DB_PASS" \
  -v migration_user="$MIGRATIONS_DB_USER" \
  -v migration_pass="$MIGRATIONS_DB_PASS" <<-EOSQL
    CREATE DATABASE :db_name;
    CREATE USER :api_user WITH PASSWORD :api_pass;
    CREATE USER :migration_user WITH PASSWORD :migration_pass;
    GRANT ALL PRIVILEGES ON DATABASE :db_name TO :migration_user;
EOSQL
