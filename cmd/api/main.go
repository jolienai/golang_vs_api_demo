package main

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/jolienaiviegas/golang_vs_api_demo/internal/app"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	cfg, err := app.LoadConfig()
	if err != nil {
		slog.Error("load config", "error", err)
		os.Exit(1)
	}

	logger := app.NewLogger(cfg.LogLevel)
	slog.SetDefault(logger)

	db, err := app.OpenDB(ctx, cfg.DatabaseURL)
	if err != nil {
		logger.Error("open database", "error", err)
		os.Exit(1)
	}
	defer db.Close()

	migrationCtx, cancelMigrations := context.WithTimeout(ctx, 30*time.Second)
	defer cancelMigrations()

	if err := app.RunMigrations(migrationCtx, db, cfg.MigrationsDir, logger); err != nil {
		logger.Error("run migrations", "error", err)
		os.Exit(1)
	}

	server := app.NewServer(cfg, logger, db)

	go func() {
		logger.Info("starting http server", "addr", cfg.HTTPAddr, "env", cfg.AppEnv)
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Error("http server failed", "error", err)
			stop()
		}
	}()

	<-ctx.Done()

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	logger.Info("shutting down http server")
	if err := server.Shutdown(shutdownCtx); err != nil {
		logger.Error("http server shutdown failed", "error", err)
		os.Exit(1)
	}
}
