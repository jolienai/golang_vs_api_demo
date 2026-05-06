package tasks

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"

	platformhttp "github.com/jolienaiviegas/golang_vs_api_demo/internal/platform/http"
)

type UpdateTaskRequest struct {
	Title       *string `json:"title"`
	Description *string `json:"description"`
	Completed   *bool   `json:"completed"`
}

type UpdateTaskInput struct {
	Title       *string
	Description *string
	Completed   *bool
}

func handleUpdateTask(store Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		if _, err := uuid.Parse(id); err != nil {
			platformhttp.NotFound(w, "task not found")
			return
		}

		var request UpdateTaskRequest
		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			platformhttp.BadRequest(w, "invalid json body")
			return
		}

		input := UpdateTaskInput{
			Title:       request.Title,
			Description: request.Description,
			Completed:   request.Completed,
		}
		if input.Title != nil {
			title := strings.TrimSpace(*input.Title)
			input.Title = &title
			if err := validateTitle(title); err != nil {
				platformhttp.ValidationError(w, err.Error())
				return
			}
		}
		if input.Description != nil {
			if err := validateDescription(*input.Description); err != nil {
				platformhttp.ValidationError(w, err.Error())
				return
			}
		}

		task, err := store.UpdateTask(r.Context(), id, input)
		if err != nil {
			if errors.Is(err, ErrTaskNotFound) {
				platformhttp.NotFound(w, "task not found")
				return
			}
			slog.ErrorContext(r.Context(), "update task failed",
				"operation", "tasks.update",
				"task_id", id,
				"error", err,
			)
			platformhttp.InternalServerError(w)
			return
		}

		platformhttp.WriteJSON(w, http.StatusOK, task)
	}
}

func (s *PostgresStore) UpdateTask(ctx context.Context, id string, input UpdateTaskInput) (Task, error) {
	const query = `
		UPDATE tasks
		SET
			title = COALESCE($2, title),
			description = COALESCE($3, description),
			completed = COALESCE($4, completed),
			updated_at = now()
		WHERE id = $1
		RETURNING id::text, title, description, completed, created_at, updated_at`

	var task Task
	err := s.db.QueryRow(ctx, query, id, input.Title, input.Description, input.Completed).Scan(
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
		return Task{}, fmt.Errorf("update task: %w", err)
	}

	return task, nil
}
