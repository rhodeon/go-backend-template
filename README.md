# go-backend-template
# Stack
- [Huma web framework](https://huma.rocks/): This has the functionality of generating an OpenAPI documentation based on the code and automatically validates and serialises requests/response, letting the focus be more on the features of the platform (besides the initial setup, at least).
- [sqlc](https://sqlc.dev/) for database operations. It generates Go structs and methods based on provided plain SQL queries.
- Postgres (should be self-explanatory).

## Getting Started
- Install the [just](https://just.systems/) automation tool to run commands in the `Justfile`.
- Run `just install-tools` to install the binaries needed for other phases.
- Set up a new `.env` file using the vars set in .env.example as a guide.
- Run `just migrations up` to set up the database.
- Run `just api` to start the server.
- An automatically generated OpenAPI spec can be viewed at the `/docs` path of the API.
- More actions can be found by running `just --list`.

## Project Structure
There are 3 primary layers:
- The API interface: concerned with web-specific operations.
- The domain (or service) level: houses the business logic.
- The repository level: the lowest level covering details beyond the scope of the domain.

Each depends on the next.

There's currently an implementation of a small API, meant to show the patterns for building a more fleshed out system.

## Todo
- Add tests.
- More package-level comments.
- Utilise more features of Huma (and OpenAPI) for a more substantial documentation.
- Integrate non-database services.