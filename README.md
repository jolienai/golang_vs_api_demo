# Golang Vertical Slice API Demo

Small REST API demo for managing tasks/todos with Go, PostgreSQL, pgx, chi, and structured logging.

The project is organized by feature slice. Task creation, listing, reading, updating, and deletion each keep their HTTP DTOs, handler logic, and SQL close together under `internal/features/tasks`. Shared application concerns such as config, database startup, middleware, and JSON responses live outside the slice.

## Stack

- Go 1.22+
- PostgreSQL
- `pgx/v5` with `pgxpool`
- `chi` router
- `log/slog`
- Docker Compose for local PostgreSQL
- Plain SQL migrations

## Project Structure

```text
.
├── cmd/api/main.go
├── internal/app
├── internal/platform/http
├── internal/features/tasks
├── migrations/001_create_tasks.sql
├── docker-compose.yml
├── .env.example
├── go.mod
└── README.md
```

## Prerequisites

- Go 1.22 or newer
- Docker and Docker Compose

## Setup

```bash
cp .env.example .env
docker compose up -d
```

This starts both PostgreSQL and the API. The API runs migrations on startup and is exposed at `http://localhost:8080`.

Run the API:

```bash
export DATABASE_URL="postgres://todo:todo@localhost:5432/todo?sslmode=disable"
go run ./cmd/api
```

Or rebuild and run it through Docker Compose:

```bash
docker compose up --build
```

Optional environment variables:

```bash
export HTTP_ADDR=":8080"
export LOG_LEVEL="info"
export APP_ENV="local"
export MIGRATIONS_DIR="migrations"
```

## Endpoints

Health check:

```bash
curl http://localhost:8080/healthz
```

Create a task:

```bash
curl -i -X POST http://localhost:8080/tasks \
  -H "Content-Type: application/json" \
  -d '{"title":"Buy milk","description":"Optional description"}'
```

List tasks:

```bash
curl http://localhost:8080/tasks
curl "http://localhost:8080/tasks?completed=false&limit=25&offset=0"
```

Get a task:

```bash
curl http://localhost:8080/tasks/e24dfbbe-aab8-4db0-9966-e8828c01f472
```

Patch a task:

```bash
curl -i -X PATCH http://localhost:8080/tasks/e24dfbbe-aab8-4db0-9966-e8828c01f472 \
  -H "Content-Type: application/json" \
  -d '{"completed":true}'
```

Delete a task:

```bash
curl -i -X DELETE http://localhost:8080/tasks/e24dfbbe-aab8-4db0-9966-e8828c01f472
```

Errors are returned consistently:

```json
{
  "type": "https://example.com/problems/validation-error",
  "title": "Validation Error",
  "status": 400,
  "detail": "title is required",
  "code": "validation_error"
}
```

## Vertical Slice Notes

To add another feature, create a sibling package under `internal/features`, for example `internal/features/projects`. Put that feature's routes, request/response DTOs, validation, handlers, and SQL in the package. Register the feature from `internal/app/server.go`.

This keeps business workflows readable without introducing generic handler, service, or repository layers before the demo needs them.

## Tests

```bash
go test ./...
```

The test suite includes a Testcontainers integration test for the PostgreSQL task store. Use short mode to run only unit tests:

```bash
go test -short ./...
```
