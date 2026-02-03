package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/labstack/echo/v4"
	echomiddleware "github.com/labstack/echo/v4/middleware"
	"github.com/zizouhuweidi/maktaba/internal/auth"
	"github.com/zizouhuweidi/maktaba/internal/config"
	"github.com/zizouhuweidi/maktaba/internal/db"
	appmiddleware "github.com/zizouhuweidi/maktaba/internal/middleware"
	"github.com/zizouhuweidi/maktaba/internal/notes"
	"github.com/zizouhuweidi/maktaba/internal/sources"
	"github.com/zizouhuweidi/maktaba/pkg"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}))

	logger.Info("starting maktaba service")

	cfg := config.Load()

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

	// Initialize Ory Kratos client
	oryClient := auth.NewClient(cfg.Ory.PublicURL, cfg.Ory.AdminURL)
	logger.Info("initialized Ory Kratos client", "public_url", cfg.Ory.PublicURL)

	sourceRepo := sources.NewPostgresRepository(database)
	noteRepo := notes.NewPostgresRepository(database)

	sourceSvc := sources.NewService(sourceRepo, logger)
	noteSvc := notes.NewService(noteRepo, logger)

	sourceHndlr := sources.NewHandler(sourceSvc, logger)
	noteHndlr := notes.NewHandler(noteSvc, logger)

	e := echo.New()

	e.Use(echomiddleware.Recover())
	e.Use(echomiddleware.RequestID())
	e.Use(echomiddleware.CORS())
	e.Use(echomiddleware.Gzip())
	e.Use(echomiddleware.TimeoutWithConfig(echomiddleware.TimeoutConfig{
		Timeout: 60 * time.Second,
	}))

	e.Use(appmiddleware.LoggingMiddleware(logger))
	e.Use(appmiddleware.PrometheusMiddleware)

	e.GET("/health", pkg.HealthCheckHandler)
	e.GET("/metrics", pkg.MetricsHandler)

	// Public routes (no auth required)
	sourceHndlr.RegisterPublicRoutes(e)
	noteHndlr.RegisterPublicRoutes(e)

	// Protected routes (auth required)
	protected := e.Group("/api")
	protected.Use(oryClient.Middleware())
	sourceHndlr.RegisterProtectedRoutes(protected)
	noteHndlr.RegisterProtectedRoutes(protected)

	s := http.Server{
		Addr:         ":" + cfg.Server.Port,
		Handler:      e,
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
	}

	go func() {
		logger.Info("server starting", "port", cfg.Server.Port)
		if err := e.StartServer(&s); err != nil && err != http.ErrServerClosed {
			logger.Error("server failed", "error", err)
			os.Exit(1)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("server shutting down")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := e.Shutdown(ctx); err != nil {
		logger.Error("server forced to shutdown", "error", err)
		os.Exit(1)
	}

	logger.Info("server exited")
}
