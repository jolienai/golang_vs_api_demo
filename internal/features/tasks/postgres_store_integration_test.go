package tasks_test

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"testing"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jolienaiviegas/golang_vs_api_demo/internal/app"
	"github.com/jolienaiviegas/golang_vs_api_demo/internal/features/tasks"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

func TestTasksAPIIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()

	container, err := postgres.RunContainer(ctx,
		testcontainers.WithImage("docker.io/postgres:16-alpine"),
		postgres.WithDatabase("todo"),
		postgres.WithUsername("todo"),
		postgres.WithPassword("todo"),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(30*time.Second),
		),
	)
	if err != nil {
		t.Fatalf("start postgres container: %v", err)
	}
	t.Cleanup(func() {
		terminateCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		if err := container.Terminate(terminateCtx); err != nil {
			t.Fatalf("terminate postgres container: %v", err)
		}
	})

	databaseURL, err := container.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		t.Fatalf("get postgres connection string: %v", err)
	}

	db, err := pgxpool.New(ctx, databaseURL)
	if err != nil {
		t.Fatalf("open database pool: %v", err)
	}
	t.Cleanup(db.Close)

	migrationsDir := filepath.Join("..", "..", "..", "migrations")
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	if err := app.RunMigrations(ctx, db, migrationsDir, logger); err != nil {
		t.Fatalf("run migrations: %v", err)
	}

	router := chi.NewRouter()
	tasks.RegisterRoutes(router, tasks.NewPostgresStore(db))

	createBody := bytes.NewBufferString(`{"title":"Integration test task","description":"Created through HTTP"}`)
	createRequest := httptest.NewRequest(http.MethodPost, "/tasks", createBody)
	createRequest.Header.Set("Content-Type", "application/json")
	createResponse := httptest.NewRecorder()

	router.ServeHTTP(createResponse, createRequest)

	if createResponse.Code != http.StatusCreated {
		t.Fatalf("expected POST /tasks status %d, got %d: %s", http.StatusCreated, createResponse.Code, createResponse.Body.String())
	}

	var created tasks.Task
	if err := json.NewDecoder(createResponse.Body).Decode(&created); err != nil {
		t.Fatalf("decode POST /tasks response: %v", err)
	}
	if created.ID == "" {
		t.Fatal("expected created task id")
	}
	if created.Title != "Integration test task" {
		t.Fatalf("expected created title %q, got %q", "Integration test task", created.Title)
	}

	getRequest := httptest.NewRequest(http.MethodGet, "/tasks/"+created.ID, nil)
	getResponse := httptest.NewRecorder()

	router.ServeHTTP(getResponse, getRequest)

	if getResponse.Code != http.StatusOK {
		t.Fatalf("expected GET /tasks/{id} status %d, got %d: %s", http.StatusOK, getResponse.Code, getResponse.Body.String())
	}

	var got tasks.Task
	if err := json.NewDecoder(getResponse.Body).Decode(&got); err != nil {
		t.Fatalf("decode GET /tasks/{id} response: %v", err)
	}
	if got.ID != created.ID {
		t.Fatalf("expected fetched id %q, got %q", created.ID, got.ID)
	}
	if got.Title != created.Title {
		t.Fatalf("expected fetched title %q, got %q", created.Title, got.Title)
	}
	if got.Description != created.Description {
		t.Fatalf("expected fetched description %q, got %q", created.Description, got.Description)
	}
}
