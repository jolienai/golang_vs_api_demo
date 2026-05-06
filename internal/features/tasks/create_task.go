package tasks

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"strings"

	platformhttp "github.com/jolienaiviegas/golang_vs_api_demo/internal/platform/http"
)

type CreateTaskRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
}

type CreateTaskInput struct {
	Title       string
	Description string
}

func handleCreateTask(store Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var request CreateTaskRequest
		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			platformhttp.BadRequest(w, "invalid json body")
			return
		}

		input := CreateTaskInput{
			Title:       strings.TrimSpace(request.Title),
			Description: request.Description,
		}
		if err := validateTitle(input.Title); err != nil {
			platformhttp.ValidationError(w, err.Error())
			return
		}
		if err := validateDescription(input.Description); err != nil {
			platformhttp.ValidationError(w, err.Error())
			return
		}

		task, err := store.CreateTask(r.Context(), input)
		if err != nil {
			slog.ErrorContext(r.Context(), "create task failed",
				"operation", "tasks.create",
				"error", err,
			)
			platformhttp.InternalServerError(w)
			return
		}

		platformhttp.WriteJSON(w, http.StatusCreated, task)
	}
}

func (s *PostgresStore) CreateTask(ctx context.Context, input CreateTaskInput) (Task, error) {
	const query = `
		INSERT INTO tasks (title, description)
		VALUES ($1, $2)
		RETURNING id::text, title, description, completed, created_at, updated_at`

	var task Task
	err := s.db.QueryRow(ctx, query, input.Title, input.Description).Scan(
		&task.ID,
		&task.Title,
		&task.Description,
		&task.Completed,
		&task.CreatedAt,
		&task.UpdatedAt,
	)
	if err != nil {
		return Task{}, fmt.Errorf("create task: %w", err)
	}

	return task, nil
}
