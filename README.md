# go-backend-template
## Stack
- [Huma web framework](https://huma.rocks/): This has the functionality of generating an OpenAPI documentation based on the code and automatically validates and serialises requests/response, letting the focus be more on the features of the platform (besides the initial setup, at least).
- [sqlc](https://sqlc.dev/) for database operations. It generates Go structs and methods based on provided plain SQL queries.
- Postgres because that's what I know best.

## Getting Started
- Install the [Mise](https://mise.jdx.dev) automation tool to run commands in the `mise.toml` file.
- Mise conveniently groups related commands together and is the preferred way of running different actions. 
- Run `mise init` to register the `git-hooks` folder in your git config.
- Mise also automatically installs needed development tools on its first run.
- Set up a new `.env` file using the vars set in .env.example as a guide.
- Run `mise migrations -- up` to set up the database.
- Run `mise api` to start the server.
- An automatically generated OpenAPI spec can be viewed at the `/docs` path of the API.
- More actions can be found by running `mise run`.

## Project Structure
There are 3 primary layers:
- The API interface: concerned with web-specific operations.
- The domain (or service) level: houses the business logic.
- The repository level: the lowest level covering details beyond the scope of the domain.

Each depends on the next.

There's currently an implementation of a small API, meant to show the patterns for building a more fleshed out system.

## Todo
- More package-level comments.
- Utilise more features of Huma (and OpenAPI) for more substantial documentation.
- Integrate non-database services.