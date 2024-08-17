set dotenv-load

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

# Runs checks including linters and vetting sqlc queries.
vet:
    go mod verify
    go vet ./...
    golangci-lint run ./...
    ./sqlc.sh
    sqlc vet

# Regenerates sqlc queries.
sqlc:
    ./sqlc.sh
    sqlc vet
    sqlc generate
