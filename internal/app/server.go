package app

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/jolienaiviegas/golang_vs_api_demo/internal/features/tasks"
	platformhttp "github.com/jolienaiviegas/golang_vs_api_demo/internal/platform/http"
)

func NewServer(cfg Config, logger *slog.Logger, db *pgxpool.Pool) *http.Server {
	router := chi.NewRouter()

	router.Use(platformhttp.RequestID)
	router.Use(platformhttp.Recoverer(logger))
	router.Use(platformhttp.RequestLogger(logger))

	router.Get("/healthz", func(w http.ResponseWriter, r *http.Request) {
		platformhttp.WriteJSON(w, http.StatusOK, map[string]string{"status": "ok"})
	})

	tasks.RegisterRoutes(router, tasks.NewPostgresStore(db))

	return &http.Server{
		Addr:              cfg.HTTPAddr,
		Handler:           router,
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      10 * time.Second,
		IdleTimeout:       60 * time.Second,
	}
}
