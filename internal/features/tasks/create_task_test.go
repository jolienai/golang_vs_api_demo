package tasks

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/go-chi/chi/v5"
)

type fakeStore struct {
	createTask func(ctx context.Context, input CreateTaskInput) (Task, error)
}

func (s fakeStore) CreateTask(ctx context.Context, input CreateTaskInput) (Task, error) {
	return s.createTask(ctx, input)
}

func (s fakeStore) GetTask(ctx context.Context, id string) (Task, error) {
	return Task{}, errors.New("not implemented")
}

func (s fakeStore) ListTasks(ctx context.Context, filter ListTasksFilter) ([]Task, error) {
	return nil, errors.New("not implemented")
}

func (s fakeStore) UpdateTask(ctx context.Context, id string, input UpdateTaskInput) (Task, error) {
	return Task{}, errors.New("not implemented")
}

func (s fakeStore) DeleteTask(ctx context.Context, id string) error {
	return errors.New("not implemented")
}

func TestCreateTaskHandler(t *testing.T) {
	t.Parallel()

	now := time.Date(2026, 5, 6, 12, 0, 0, 0, time.UTC)
	store := fakeStore{
		createTask: func(ctx context.Context, input CreateTaskInput) (Task, error) {
			if input.Title != "Buy milk" {
				t.Fatalf("expected trimmed title, got %q", input.Title)
			}
			if input.Description != "Optional description" {
				t.Fatalf("expected description, got %q", input.Description)
			}
			return Task{
				ID:          "e24dfbbe-aab8-4db0-9966-e8828c01f472",
				Title:       input.Title,
				Description: input.Description,
				Completed:   false,
				CreatedAt:   now,
				UpdatedAt:   now,
			}, nil
		},
	}

	router := chi.NewRouter()
	RegisterRoutes(router, store)

	body := bytes.NewBufferString(`{"title":" Buy milk ","description":"Optional description"}`)
	request := httptest.NewRequest(http.MethodPost, "/tasks", body)
	response := httptest.NewRecorder()

	router.ServeHTTP(response, request)

	if response.Code != http.StatusCreated {
		t.Fatalf("expected status %d, got %d: %s", http.StatusCreated, response.Code, response.Body.String())
	}

	var task Task
	if err := json.NewDecoder(response.Body).Decode(&task); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if task.ID == "" {
		t.Fatal("expected task id")
	}
	if task.Title != "Buy milk" {
		t.Fatalf("expected title, got %q", task.Title)
	}
}
