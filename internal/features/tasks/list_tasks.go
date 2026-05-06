package tasks

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"

	platformhttp "github.com/jolienaiviegas/golang_vs_api_demo/internal/platform/http"
)

var (
	errInvalidCompleted = errors.New("completed must be true or false")
	errInvalidLimit     = errors.New("limit must be an integer")
	errInvalidOffset    = errors.New("offset must be an integer")
)

type ListTasksResponse struct {
	Tasks  []Task `json:"tasks"`
	Limit  int    `json:"limit"`
	Offset int    `json:"offset"`
}

type ListTasksFilter struct {
	Completed *bool
	Limit     int
	Offset    int
}

func handleListTasks(store Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		filter, err := parseListTasksFilter(r)
		if err != nil {
			platformhttp.ValidationError(w, err.Error())
			return
		}

		tasks, err := store.ListTasks(r.Context(), filter)
		if err != nil {
			slog.ErrorContext(r.Context(), "list tasks failed",
				"operation", "tasks.list",
				"completed", filter.Completed,
				"limit", filter.Limit,
				"offset", filter.Offset,
				"error", err,
			)
			platformhttp.InternalServerError(w)
			return
		}

		platformhttp.WriteJSON(w, http.StatusOK, ListTasksResponse{
			Tasks:  tasks,
			Limit:  filter.Limit,
			Offset: filter.Offset,
		})
	}
}

func parseListTasksFilter(r *http.Request) (ListTasksFilter, error) {
	query := r.URL.Query()
	filter := ListTasksFilter{}

	if completedValue := query.Get("completed"); completedValue != "" {
		completed, err := strconv.ParseBool(completedValue)
		if err != nil {
			return ListTasksFilter{}, errInvalidCompleted
		}
		filter.Completed = &completed
	}

	if limitValue := query.Get("limit"); limitValue != "" {
		limit, err := strconv.Atoi(limitValue)
		if err != nil {
			return ListTasksFilter{}, errInvalidLimit
		}
		filter.Limit = limit
	}

	if offsetValue := query.Get("offset"); offsetValue != "" {
		offset, err := strconv.Atoi(offsetValue)
		if err != nil {
			return ListTasksFilter{}, errInvalidOffset
		}
		filter.Offset = offset
	}

	filter.Limit, filter.Offset = normalizeLimitOffset(filter.Limit, filter.Offset)
	return filter, nil
}

func (s *PostgresStore) ListTasks(ctx context.Context, filter ListTasksFilter) ([]Task, error) {
	const query = `
		SELECT id::text, title, description, completed, created_at, updated_at
		FROM tasks
		WHERE ($1::boolean IS NULL OR completed = $1)
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3`

	rows, err := s.db.Query(ctx, query, filter.Completed, filter.Limit, filter.Offset)
	if err != nil {
		return nil, fmt.Errorf("list tasks query: %w", err)
	}
	defer rows.Close()

	tasks := make([]Task, 0)
	for rows.Next() {
		var task Task
		if err := rows.Scan(
			&task.ID,
			&task.Title,
			&task.Description,
			&task.Completed,
			&task.CreatedAt,
			&task.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("scan task: %w", err)
		}
		tasks = append(tasks, task)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate tasks: %w", err)
	}

	return tasks, nil
}
