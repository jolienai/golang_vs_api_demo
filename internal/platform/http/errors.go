package http

import (
	"encoding/json"
	"net/http"
)

type ProblemDetail struct {
	Type   string `json:"type"`
	Title  string `json:"title"`
	Status int    `json:"status"`
	Detail string `json:"detail"`
	Code   string `json:"code"`
}

type problemDefinition struct {
	Type  string
	Title string
	Code  string
}

var problemDefinitions = map[string]problemDefinition{
	"bad_request": {
		Type:  "https://example.com/problems/bad-request",
		Title: "Bad Request",
		Code:  "bad_request",
	},
	"validation_error": {
		Type:  "https://example.com/problems/validation-error",
		Title: "Validation Error",
		Code:  "validation_error",
	},
	"not_found": {
		Type:  "https://example.com/problems/not-found",
		Title: "Not Found",
		Code:  "not_found",
	},
	"internal_server_error": {
		Type:  "https://example.com/problems/internal-server-error",
		Title: "Internal Server Error",
		Code:  "internal_server_error",
	},
}

func WriteError(w http.ResponseWriter, status int, code, detail string) {
	definition, ok := problemDefinitions[code]
	if !ok {
		definition = problemDefinition{
			Type:  "https://example.com/problems/error",
			Title: http.StatusText(status),
			Code:  code,
		}
	}

	w.Header().Set("Content-Type", "application/problem+json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(ProblemDetail{
		Type:   definition.Type,
		Title:  definition.Title,
		Status: status,
		Detail: detail,
		Code:   definition.Code,
	})
}

func BadRequest(w http.ResponseWriter, message string) {
	WriteError(w, http.StatusBadRequest, "bad_request", message)
}

func ValidationError(w http.ResponseWriter, message string) {
	WriteError(w, http.StatusBadRequest, "validation_error", message)
}

func NotFound(w http.ResponseWriter, message string) {
	WriteError(w, http.StatusNotFound, "not_found", message)
}

func InternalServerError(w http.ResponseWriter) {
	WriteError(w, http.StatusInternalServerError, "internal_server_error", "internal server error")
}
