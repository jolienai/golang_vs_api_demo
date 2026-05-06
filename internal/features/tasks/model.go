package tasks

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Task struct {
	ID          string    `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Completed   bool      `json:"completed"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type Store interface {
	CreateTask(ctx context.Context, input CreateTaskInput) (Task, error)
	GetTask(ctx context.Context, id string) (Task, error)
	ListTasks(ctx context.Context, filter ListTasksFilter) ([]Task, error)
	UpdateTask(ctx context.Context, id string, input UpdateTaskInput) (Task, error)
	DeleteTask(ctx context.Context, id string) error
}

type PostgresStore struct {
	db *pgxpool.Pool
}

func NewPostgresStore(db *pgxpool.Pool) *PostgresStore {
	return &PostgresStore{db: db}
}
