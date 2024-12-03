set dotenv-filename := ".env"

dev_tools_dir := "dev-tools"

# Default recipe to display all available recipes and their information.
default:
  @just --list

# Installs tools needed for running other tasks to aid development.
install-dev-tools:
    @mkdir -p {{ dev_tools_dir }}
    GOBIN=$(pwd)/dev-tools go install github.com/sqlc-dev/sqlc/cmd/sqlc@v1.27.0
    GOBIN=$(pwd)/dev-tools go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.62.2

# Adds the relatively installed dev tools to $PATH to make them usable by other recipes. Should be depended on by any recipe which uses a dev tool.
set-dev-tools:
    @export PATH="{{ dev_tools_dir }}:$PATH"

# Starts the API service
api:
    @cd ./cmd/api && go run .

# Runs migrations. Run `just migrations` to see all commands.
migrations *args:
    cd ./cmd/migrations && go run . {{args}}

# Formats the codebase and vendors dependencies.
tidy:
    @echo "formatting codebase..."
    go fmt ./...
    go mod tidy -v
    go mod vendor

# Runs lint checks.
lint: set-dev-tools
    golangci-lint run ./...

# Runs checks including linters and vetting sqlc queries.
vet: set-dev-tools
    go mod verify
    go vet ./...
    golangci-lint run ./...
    ./sqlc.sh
    sqlc vet

# Runs all tests in codebase
test:
    go test ./...

# Regenerates sqlc queries.
sqlc: set-dev-tools
    ./sqlc.sh
    sqlc vet
    sqlc generate
