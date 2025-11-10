# go-backend-template

This repository is an attempt to formalise a framework (for myself and potentially others) for writing Web APIs in Go
based on patterns I've observed and implemented over the years across different codebases. I expect it to be in flux
and evolve over time with changes in my opinions and experience.

The case study used here is based on the "blessed"
[OpenAPI PetStore specification](https://learn.openapis.org/examples/v3.0/petstore.html) with changes to make it more
functional.

## Stack

- [Huma](https://huma.rocks/) as the web framework. Huma generates OpenAPI documentation from the code and automatically
  validates and (de)serialises requests/responses, letting the focus be more on the features of the service (after the
  initial setup, at least).
- [Sqlc](https://sqlc.dev/) for database operations. Sqlc generates Go models and methods from pure (or very close to
  pure) SQL queries.
- [PostgreSQL](https://www.postgresql.org/) as the database because that's what I know and use best.

## Requirements

Only two programs need to be explicitly installed: [Docker](https://www.docker.com/) and [Mise](https://mise.jdx.dev/).
They both serve the complementary purpose of providing and managing other tools and services required for development.

Docker is used to handle long-running services like the database and cache servers, while Mise is used for one-off CLI
tools which are run occasionally like linters and code generators.

- Docker is already ubiquitous, so I won't go further into its usage here.
- Mise is used for running tasks and managing environment variables in addition to the tool management mentioned above.
- At the time of writing, the project was tested on an environment running Docker Engine version `26.1.1` and Mise
  version `2025.11.2`.
- This project has been tested only on Linux but should work on macOS out of the box as well. Scripts make use of
  UNIX commands which aren't preinstalled on Windows, so one of the various options for running a UNIX shell will be
  needed there.

## Getting Started

- Ensure both Docker and Mise are installed before proceeding.
- Run `mise run init` to register the `git-hooks` folder in your git config and automatically generate a `.env` file
  from `example.env`. Mise also automatically installs necessary development tools for the project on its first run.
- Populate the required environment variables in the generated `.env` file.
- Run `mise run migrations -- up` to set up the database.
- Run `mise run api` to start the API server.
- An automatically generated OpenAPI spec can be viewed at the `/docs` path of the running API.
- More actions can be found by running `mise tasks`.

## Project Structure

There are three primary layers with each depending on the next:

- The [API layer](./cmd/api) is concerned with web-specific operations.
- The [service layer](./services) houses the core business logic of the overall platform.
- The [repository level](./repositories) is the lowest level covering details beyond the scope of the domain like the
  operations of external services (the database, cache, email provider and so on).

## Configuration

All development configuration values are set as environment variables in a single top-level `.env` file and documented
under their respective main packages. [example.env](example.env) holds an exhaustive (and automatically synced)
representation of what the `.env` file can contain.

- [API config documentation](./cmd/api/config.md)
- [Migrations config documentation](./cmd/migrations/config.md)

## Todo

- Add more package-level documentation.
- Add tests for more endpoints.