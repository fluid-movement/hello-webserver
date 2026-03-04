# Todo App

A simple todo app in Go.

## Stack

- **[Go](https://go.dev/)** — HTTP server (`net/http`)
- **[templ](https://templ.guide/)** — HTML templates (compiled to Go)
- **[HTMX](https://htmx.org/)** — frontend interactivity, no JS written
- **[sqlc](https://sqlc.dev/)** — type-safe SQL queries (generated Go code)
- **[pgx](https://github.com/jackc/pgx)** — PostgreSQL driver
- **[golang-migrate](https://github.com/golang-migrate/migrate)** — DB migrations, run automatically on startup
- **[air](https://github.com/air-verse/air)** — live reload for development

## Development

**Prerequisites:** Go, PostgreSQL, [air](https://github.com/air-verse/air), [templ CLI](https://templ.guide/quick-start/installation)

```sh
cp .env.example .env   # set DATABASE_URL
air                    # starts server with live reload on :8080
```

Migrations run automatically on startup. To add a migration:

```sh
migrations/000002_name.up.sql
migrations/000002_name.down.sql
```

To regenerate DB queries after editing `db/query.sql`:

```sh
sqlc generate
```

To regenerate templates after editing `*.templ` files:

```sh
templ generate
```
