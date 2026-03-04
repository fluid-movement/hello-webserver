package main

import (
	"context"
	"embed"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/a-h/templ"
	"github.com/fluid-movement/hello-webserver/assets"
	"github.com/fluid-movement/hello-webserver/internal/db"
	"github.com/fluid-movement/hello-webserver/templates"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/pgx/v5"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

//go:embed migrations/*.sql
var migrationsFS embed.FS

type App struct {
	queries *db.Queries
}

func main() {
	godotenv.Load()

	ctx := context.Background()
	pool, err := pgxpool.New(ctx, os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatalf("unable to connect to database: %v", err)
	}
	defer pool.Close()

	if err := runMigrations(os.Getenv("DATABASE_URL")); err != nil {
		log.Fatalf("migration failed: %v", err)
	}

	app := &App{queries: db.New(pool)}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	mux := http.NewServeMux()
	setupAssetsRoutes(mux)
	mux.HandleFunc("GET /", app.handleIndex)
	mux.HandleFunc("POST /todos", app.handleCreate)
	mux.HandleFunc("PATCH /todos/{id}/toggle", app.handleToggle)
	mux.HandleFunc("DELETE /todos/{id}", app.handleDelete)

	log.Printf("listening on :%s", port)
	if err := http.ListenAndServe(":"+port, mux); err != nil {
		log.Fatalf("server error: %v", err)
	}
}

func (a *App) handleIndex(w http.ResponseWriter, r *http.Request) {
	todos, err := a.queries.ListTodos(r.Context())
	if err != nil {
		http.Error(w, "failed to list todos", http.StatusInternalServerError)
		return
	}
	templ.Handler(templates.Index(todos)).ServeHTTP(w, r)
}

func (a *App) handleCreate(w http.ResponseWriter, r *http.Request) {
	title := r.FormValue("title")
	if title == "" {
		http.Error(w, "title is required", http.StatusBadRequest)
		return
	}
	if len(title) > 500 {
		http.Error(w, "title too long", http.StatusBadRequest)
		return
	}
	todo, err := a.queries.CreateTodo(r.Context(), title)
	if err != nil {
		http.Error(w, "failed to create todo", http.StatusInternalServerError)
		return
	}
	templ.Handler(templates.TodoItem(todo)).ServeHTTP(w, r)
}

func (a *App) handleToggle(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}
	todo, err := a.queries.ToggleTodo(r.Context(), id)
	if err != nil {
		http.Error(w, "failed to toggle todo", http.StatusInternalServerError)
		return
	}
	templ.Handler(templates.TodoItem(todo)).ServeHTTP(w, r)
}

func setupAssetsRoutes(mux *http.ServeMux) {
	isDevelopment := os.Getenv("GO_ENV") != "production"

	assetHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if isDevelopment {
			w.Header().Set("Cache-Control", "no-store")
		}
		var fs http.Handler
		if isDevelopment {
			fs = http.FileServer(http.Dir("./assets"))
		} else {
			fs = http.FileServer(http.FS(assets.Assets))
		}
		fs.ServeHTTP(w, r)
	})

	mux.Handle("GET /assets/", http.StripPrefix("/assets/", assetHandler))
}

func runMigrations(databaseURL string) error {
	src, err := iofs.New(migrationsFS, "migrations")
	if err != nil {
		return err
	}
	migrateURL := strings.Replace(databaseURL, "postgres://", "pgx5://", 1)
	m, err := migrate.NewWithSourceInstance("iofs", src, migrateURL)
	if err != nil {
		return err
	}
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return err
	}
	return nil
}

func (a *App) handleDelete(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}
	if err := a.queries.DeleteTodo(r.Context(), id); err != nil {
		http.Error(w, "failed to delete todo", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}
