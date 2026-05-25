package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"runtime/debug"
	"syscall"
	"time"

	"github.com/zizouhuweidi/maktaba/internal/auth"
	"github.com/zizouhuweidi/maktaba/internal/config"
	"github.com/zizouhuweidi/maktaba/internal/db"
	"github.com/zizouhuweidi/maktaba/internal/notes"
	"github.com/zizouhuweidi/maktaba/internal/sources"
	"github.com/zizouhuweidi/maktaba/pkg"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	logger.Info("starting maktaba service")

	cfg := config.Load()
	if cfg.Auth.Ed25519PrivateKey == "" {
		logger.Warn("AUTH_ED25519_PRIVATE_KEY is not set; generated JWT keys will be ephemeral")
	}

	database, err := db.NewDB(cfg.Database.URL, cfg.Database.MaxOpenConns, cfg.Database.MaxIdleConns, cfg.Database.ConnMaxLifetime)
	if err != nil {
		logger.Error("failed to connect to database", "error", err)
		os.Exit(1)
	}
	defer database.Close()
	logger.Info("connected to database")

	tokenManager, err := auth.NewTokenManager(cfg.Auth.Issuer, cfg.Auth.Audience, cfg.Auth.Ed25519PrivateKey, cfg.Auth.AccessTokenLifetime)
	if err != nil {
		logger.Error("failed to initialize token manager", "error", err)
		os.Exit(1)
	}

	authRepo := auth.NewPostgresRepository(database)
	sourceRepo := sources.NewPostgresRepository(database)
	noteRepo := notes.NewPostgresRepository(database)

	authSvc := auth.NewService(authRepo, tokenManager, cfg.Auth.RefreshTokenLifetime, logger)
	sourceSvc := sources.NewService(sourceRepo, logger)
	noteSvc := notes.NewService(noteRepo, logger)

	authHndlr := auth.NewHandler(authSvc, cfg.Auth.CookieSecure, logger)
	sourceHndlr := sources.NewHandler(sourceSvc, logger)
	noteHndlr := notes.NewHandler(noteSvc, logger)

	mux := http.NewServeMux()
	mux.HandleFunc("GET /health", pkg.HealthCheckHandler)
	authHndlr.RegisterRoutes(mux)
	authHndlr.RegisterProtectedRoutes(mux)
	sourceHndlr.RegisterPublicRoutes(mux)
	noteHndlr.RegisterPublicRoutes(mux)
	sourceHndlr.RegisterProtectedRoutes(mux, authHndlr.Middleware)
	noteHndlr.RegisterProtectedRoutes(mux, authHndlr.Middleware)

	handler := recoverMiddleware(logger)(corsMiddleware(loggingMiddleware(logger)(mux)))
	server := http.Server{
		Addr:         ":" + cfg.Server.Port,
		Handler:      handler,
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
	}

	go func() {
		logger.Info("server starting", "port", cfg.Server.Port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
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
	if err := server.Shutdown(ctx); err != nil {
		logger.Error("server forced to shutdown", "error", err)
		os.Exit(1)
	}
	logger.Info("server exited")
}

func loggingMiddleware(logger *slog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			next.ServeHTTP(w, r)
			logger.Info("request completed", "method", r.Method, "path", r.URL.Path, "duration", time.Since(start))
		})
	}
}

func recoverMiddleware(logger *slog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if recovered := recover(); recovered != nil {
					logger.Error("panic recovered", "panic", recovered, "stack", string(debug.Stack()))
					http.Error(w, "internal server error", http.StatusInternalServerError)
				}
			}()
			next.ServeHTTP(w, r)
		})
	}
}

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		next.ServeHTTP(w, r)
	})
}
