# go-backend-template

## Stack

- [Huma web framework](https://huma.rocks/): This has the functionality of generating an OpenAPI documentation based on
  the code and automatically validates and serialises requests/response, letting the focus be more on the features of
  the platform (besides the initial setup, at least).
- [Sqlc](https://sqlc.dev/) for database operations. It generates Go structs and methods based on provided plain SQL
  queries.
- [PostgreSQL](https://www.postgresql.org/) because that's what I know best.

## Getting Started

- Install the [Mise](https://mise.jdx.dev) automation tool to run commands in the `mise.toml` file.
- Mise conveniently groups related commands together and is the preferred way of running different actions.
- Run `mise init` to register the `git-hooks` folder in your git config and automatically generate a `.env` file from `example.env`.
- Populate the required environment variables in the generated `.env` file.
- Mise also automatically installs necessary development tools on its first run.
- Run `mise migrations -- up` to set up the database.
- Run `mise api` to start the API server.
- An automatically generated OpenAPI spec can be viewed at the `/docs` path of the API.
- More actions can be found by running `mise tasks`.

## Project Structure

There are three primary layers with each depending on the next:

- The [API interface](./cmd/api): concerned with web-specific operations.
- The [domain (or service) level](./services): houses the business logic.
- The [repository level](./repositories): the lowest level covering details beyond the scope of the domain.

There's currently an implementation of a small API based on the blessed OpenAPI PetStore specification, which is meant to show the patterns for building more fleshed-out systems.

## Configuration
All required development configuration values are set in a single top-level `.env` file and documented under their respective main packages. `example.env` holds an exhaustive representation of what the `.env` file can contain.
- [API configuration documentation](./cmd/api/config.md)
- [Migrations configuration documentation](./cmd/migrations/config.md)

## Todo

- More package-level documentation.
- Tests for more endpoints.