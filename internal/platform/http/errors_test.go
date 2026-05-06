package http

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestWriteErrorWritesProblemDetail(t *testing.T) {
	t.Parallel()

	response := httptest.NewRecorder()

	ValidationError(response, "title is required")

	if response.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d", http.StatusBadRequest, response.Code)
	}
	if contentType := response.Header().Get("Content-Type"); !strings.HasPrefix(contentType, "application/problem+json") {
		t.Fatalf("expected problem detail content type, got %q", contentType)
	}

	var problem ProblemDetail
	if err := json.NewDecoder(response.Body).Decode(&problem); err != nil {
		t.Fatalf("decode problem detail: %v", err)
	}

	if problem.Type != "https://example.com/problems/validation-error" {
		t.Fatalf("expected validation problem type, got %q", problem.Type)
	}
	if problem.Title != "Validation Error" {
		t.Fatalf("expected title %q, got %q", "Validation Error", problem.Title)
	}
	if problem.Status != http.StatusBadRequest {
		t.Fatalf("expected body status %d, got %d", http.StatusBadRequest, problem.Status)
	}
	if problem.Detail != "title is required" {
		t.Fatalf("expected detail %q, got %q", "title is required", problem.Detail)
	}
	if problem.Code != "validation_error" {
		t.Fatalf("expected code %q, got %q", "validation_error", problem.Code)
	}
}

func TestNotFoundWritesProblemDetail(t *testing.T) {
	t.Parallel()

	response := httptest.NewRecorder()

	NotFound(response, "task not found")

	if response.Code != http.StatusNotFound {
		t.Fatalf("expected status %d, got %d", http.StatusNotFound, response.Code)
	}

	var problem ProblemDetail
	if err := json.NewDecoder(response.Body).Decode(&problem); err != nil {
		t.Fatalf("decode problem detail: %v", err)
	}

	if problem.Type != "https://example.com/problems/not-found" {
		t.Fatalf("expected not found problem type, got %q", problem.Type)
	}
	if problem.Title != "Not Found" {
		t.Fatalf("expected title %q, got %q", "Not Found", problem.Title)
	}
	if problem.Status != http.StatusNotFound {
		t.Fatalf("expected body status %d, got %d", http.StatusNotFound, problem.Status)
	}
	if problem.Detail != "task not found" {
		t.Fatalf("expected detail %q, got %q", "task not found", problem.Detail)
	}
	if problem.Code != "not_found" {
		t.Fatalf("expected code %q, got %q", "not_found", problem.Code)
	}
}
