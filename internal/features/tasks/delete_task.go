package tasks

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"

	platformhttp "github.com/jolienaiviegas/golang_vs_api_demo/internal/platform/http"
)

func handleDeleteTask(store Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		if _, err := uuid.Parse(id); err != nil {
			platformhttp.NotFound(w, "task not found")
			return
		}

		if err := store.DeleteTask(r.Context(), id); err != nil {
			if errors.Is(err, ErrTaskNotFound) {
				platformhttp.NotFound(w, "task not found")
				return
			}
			slog.ErrorContext(r.Context(), "delete task failed",
				"operation", "tasks.delete",
				"task_id", id,
				"error", err,
			)
			platformhttp.InternalServerError(w)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}

func (s *PostgresStore) DeleteTask(ctx context.Context, id string) error {
	const query = `DELETE FROM tasks WHERE id = $1`

	result, err := s.db.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("delete task: %w", err)
	}
	if result.RowsAffected() == 0 {
		return ErrTaskNotFound
	}

	return nil
}
