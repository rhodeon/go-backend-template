set dotenv-filename := ".env"

# Starts the API service
api:
    @cd ./cmd/api && go run . || true  # `ignore_error: true` translates to `|| true`

# Runs migrations. Run `just migrations` to see all commands.
migrations cli_args="":
    cd ./cmd/migrations && go run . {{cli_args}}

# Formats the codebase and vendors dependencies.
tidy:
    @echo "formatting codebase..."
    go fmt ./...
    go mod tidy -v
    go mod vendor

# Runs lint checks.
lint:
    golangci-lint run ./...

# Runs checks including linters and vetting sqlc queries.
vet:
    go mod verify
    go vet ./...
    golangci-lint run ./...
    ./sqlc.sh
    sqlc vet

# Runs all tests in codebase
test:
    go test ./...

# Regenerates sqlc queries.
sqlc:
    ./sqlc.sh
    sqlc vet
    sqlc generate


## Installs Go-based tools needed for running other tasks.
install-tools:
    go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest
    go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest