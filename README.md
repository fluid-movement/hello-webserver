# Todo App

A simple todo app in Go.

## Stack

- **[Go](https://go.dev/)** — HTTP server (`net/http`)
- **[templ](https://templ.guide/)** — HTML templates compiled to Go
- **[templui](https://templui.io/)** — UI component library (Button, Input, Checkbox, Card)
- **[HTMX](https://htmx.org/)** — frontend interactivity, no JS written
- **[Tailwind CSS v4](https://tailwindcss.com/)** — styling (Catppuccin Frappé theme)
- **[sqlc](https://sqlc.dev/)** — type-safe SQL (generated Go code)
- **[pgx](https://github.com/jackc/pgx)** — PostgreSQL driver
- **[golang-migrate](https://github.com/golang-migrate/migrate)** — DB migrations, run on startup

## Development

**Prerequisites:** Go, PostgreSQL, [templ CLI](https://templ.guide/quick-start/installation), [Task](https://taskfile.dev/)

```sh
cp .env.example .env  # set DATABASE_URL
task dev              # starts server with live reload on :8080
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
