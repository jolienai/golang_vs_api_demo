package tasks

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"

	platformhttp "github.com/jolienaiviegas/golang_vs_api_demo/internal/platform/http"
)

var ErrTaskNotFound = errors.New("task not found")

func handleGetTask(store Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		if _, err := uuid.Parse(id); err != nil {
			platformhttp.NotFound(w, "task not found")
			return
		}

		task, err := store.GetTask(r.Context(), id)
		if err != nil {
			if errors.Is(err, ErrTaskNotFound) {
				platformhttp.NotFound(w, "task not found")
				return
			}
			slog.ErrorContext(r.Context(), "get task failed",
				"operation", "tasks.get",
				"task_id", id,
				"error", err,
			)
			platformhttp.InternalServerError(w)
			return
		}

		platformhttp.WriteJSON(w, http.StatusOK, task)
	}
}

func (s *PostgresStore) GetTask(ctx context.Context, id string) (Task, error) {
	const query = `
		SELECT id::text, title, description, completed, created_at, updated_at
		FROM tasks
		WHERE id = $1`

	var task Task
	err := s.db.QueryRow(ctx, query, id).Scan(
		&task.ID,
		&task.Title,
		&task.Description,
		&task.Completed,
		&task.CreatedAt,
		&task.UpdatedAt,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return Task{}, ErrTaskNotFound
	}
	if err != nil {
		return Task{}, fmt.Errorf("get task: %w", err)
	}

	return task, nil
}
