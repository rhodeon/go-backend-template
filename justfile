set dotenv-filename := ".env"

dev_tools_dir := justfile_directory() + "/.dev-tools"

# The dev tools directory is added to $PATH to make the installed binaries usable by recipes.
export PATH := dev_tools_dir + ":" + env_var("PATH")

# Default recipe to display all available recipes and their information.
default:
    @just --list

# This should be the first non-default recipe run after cloning the repository. It makes setups needed for the lifetime of the project.
init:
    @# Since the .git folder doesn't get tracked as path of the repository, hooks are stored in a separate directory which is tracked instead.
    @# This command makes git use the tracked directly for hooks instead of the default.
    @echo "configuring git hooks directory..."
    @git config core.hooksPath .git-hooks

    @# For a fresh clone, all development tools need to be installed.
    @just install-dev-tools

# Installs tools needed for running other tasks to aid development.
install-dev-tools:
    @mkdir -p {{ dev_tools_dir }}
    @echo "installing development tools into {{ dev_tools_dir }}..."
    GOBIN={{ dev_tools_dir }} go install github.com/sqlc-dev/sqlc/cmd/sqlc@v1.27.0
    GOBIN={{ dev_tools_dir }} go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.62.2

# Starts the API service
api:
    @cd ./cmd/api && go run .

# Runs migrations. Run `just migrations` to see all commands.
migrations *args:
    cd ./cmd/migrations && go run . {{ args }}

# Runs all tests in codebase
test:
    go test ./...

# Formats the codebase uniformly and vendors dependencies.
tidy:
    @echo "formatting codebase..."
    @golangci-lint run --fix --enable-only gofmt,gofumpt

    @echo "tidying dependencies..."
    go mod tidy -v
    go mod vendor

# Runs lint checks.
lint:
    golangci-lint run ./...

# Runs checks including linters and vetting sqlc queries.
vet:
    @echo "vetting Go code..."
    @go mod verify
    @golangci-lint run ./...

    @echo "vetting SQL code..."
    @./sqlc.sh
    @sqlc vet

# Regenerates sqlc queries.
sqlc:
    ./sqlc.sh
    sqlc vet
    sqlc generate
