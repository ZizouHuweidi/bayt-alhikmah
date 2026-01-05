package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	chiMiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/zizouhuweidi/maktaba/internal/config"
	"github.com/zizouhuweidi/maktaba/internal/db"
	"github.com/zizouhuweidi/maktaba/internal/middleware"
	"github.com/zizouhuweidi/maktaba/internal/notes"
	"github.com/zizouhuweidi/maktaba/internal/pkg"
	"github.com/zizouhuweidi/maktaba/internal/sources"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}))

	logger.Info("starting maktaba service")

	// Load configuration
	cfg := config.Load()

	// Initialize database
	database, err := db.NewDB(
		cfg.Database.URL,
		cfg.Database.MaxOpenConns,
		cfg.Database.MaxIdleConns,
		cfg.Database.ConnMaxLifetime,
	)
	if err != nil {
		logger.Error("failed to connect to database", "error", err)
		os.Exit(1)
	}
	defer database.Close()

	logger.Info("connected to database")

	// Initialize repositories
	sourceRepo := sources.NewPostgresRepository(database)
	noteRepo := notes.NewPostgresRepository(database)

	// Initialize services
	sourceSvc := sources.NewService(sourceRepo, logger)
	noteSvc := notes.NewService(noteRepo, logger)

	// Initialize handlers
	sourceHndlr := sources.NewHandler(sourceSvc, logger)
	noteHndlr := notes.NewHandler(noteSvc, logger)

	// Setup router
	r := chi.NewRouter()

	// Standard middleware
	r.Use(chiMiddleware.Recoverer)
	r.Use(chiMiddleware.RequestID)
	r.Use(chiMiddleware.RealIP)
	r.Use(chiMiddleware.Timeout(60 * time.Second))
	r.Use(chiMiddleware.SetHeader("Content-Type", "application/json"))

	// Custom middleware
	r.Use(middleware.LoggingMiddleware(logger))
	r.Use(middleware.PrometheusMiddleware)

	// Health and metrics
	r.Get("/health", pkg.HealthCheckHandler)
	r.Get("/metrics", pkg.MetricsHandler)

	// Register domain routes
	sourceHndlr.RegisterRoutes(r)
	noteHndlr.RegisterRoutes(r)

	// Server configuration
	srv := &http.Server{
		Addr:         ":" + cfg.Server.Port,
		Handler:      r,
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
	}

	// Start server
	go func() {
		logger.Info("server starting", "port", cfg.Server.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Error("server failed", "error", err)
			os.Exit(1)
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("server shutting down")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		logger.Error("server forced to shutdown", "error", err)
		os.Exit(1)
	}

	logger.Info("server exited")
}
