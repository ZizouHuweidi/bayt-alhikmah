package server

import (
	"log/slog"
	"net/http"
	"runtime/debug"
	"slices"
	"time"

	"github.com/zizouhuweidi/maktaba/internal/auth"
	"github.com/zizouhuweidi/maktaba/internal/collections"
	"github.com/zizouhuweidi/maktaba/internal/config"
	"github.com/zizouhuweidi/maktaba/internal/db"
	"github.com/zizouhuweidi/maktaba/internal/health"
	"github.com/zizouhuweidi/maktaba/internal/library"
	"github.com/zizouhuweidi/maktaba/internal/notes"
	"github.com/zizouhuweidi/maktaba/internal/profiles"
	"github.com/zizouhuweidi/maktaba/internal/reviews"
	"github.com/zizouhuweidi/maktaba/internal/sources"
)

func New(cfg *config.Config, database *db.DB, logger *slog.Logger) (*http.Server, error) {
	tokenManager, err := auth.NewTokenManager(cfg.Auth.Issuer, cfg.Auth.Audience, cfg.Auth.Ed25519PrivateKey, cfg.Auth.AccessTokenLifetime)
	if err != nil {
		return nil, err
	}

	authRepo := auth.NewPostgresRepository(database)
	collectionRepo := collections.NewPostgresRepository(database)
	libraryRepo := library.NewPostgresRepository(database)
	sourceRepo := sources.NewPostgresRepository(database)
	noteRepo := notes.NewPostgresRepository(database)
	profileRepo := profiles.NewPostgresRepository(database)
	reviewRepo := reviews.NewPostgresRepository(database)

	authSvc := auth.NewService(authRepo, tokenManager, cfg.Auth.RefreshTokenLifetime, logger)
	collectionSvc := collections.NewService(collectionRepo, logger)
	librarySvc := library.NewService(libraryRepo, logger)
	sourceSvc := sources.NewService(sourceRepo, logger)
	noteSvc := notes.NewService(noteRepo, logger)
	profileSvc := profiles.NewService(profileRepo, logger)
	reviewSvc := reviews.NewService(reviewRepo, logger)

	authHndlr := auth.NewHandler(authSvc, cfg.Auth.CookieSecure, logger)
	collectionHndlr := collections.NewHandler(collectionSvc, logger)
	libraryHndlr := library.NewHandler(librarySvc, logger)
	sourceHndlr := sources.NewHandler(sourceSvc, logger)
	noteHndlr := notes.NewHandler(noteSvc, logger)
	profileHndlr := profiles.NewHandler(profileSvc, logger)
	reviewHndlr := reviews.NewHandler(reviewSvc, logger)

	mux := http.NewServeMux()
	mux.HandleFunc("GET /health", health.Handler)
	authHndlr.RegisterRoutes(mux)
	authHndlr.RegisterProtectedRoutes(mux)
	collectionHndlr.RegisterPublicRoutes(mux)
	libraryHndlr.RegisterPublicRoutes(mux)
	sourceHndlr.RegisterPublicRoutes(mux)
	noteHndlr.RegisterPublicRoutes(mux)
	profileHndlr.RegisterPublicRoutes(mux)
	reviewHndlr.RegisterPublicRoutes(mux)
	collectionHndlr.RegisterProtectedRoutes(mux, authHndlr.Middleware)
	libraryHndlr.RegisterProtectedRoutes(mux, authHndlr.Middleware)
	sourceHndlr.RegisterProtectedRoutes(mux, authHndlr.Middleware)
	noteHndlr.RegisterProtectedRoutes(mux, authHndlr.Middleware)
	profileHndlr.RegisterProtectedRoutes(mux, authHndlr.Middleware)
	reviewHndlr.RegisterProtectedRoutes(mux, authHndlr.Middleware)

	handler := recoverMiddleware(logger)(corsMiddleware(cfg.Server.CORSAllowedOrigins)(loggingMiddleware(logger)(mux)))

	return &http.Server{
		Addr:         ":" + cfg.Server.Port,
		Handler:      handler,
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
	}, nil
}

type statusRecorder struct {
	http.ResponseWriter
	status int
	bytes  int
}

func (r *statusRecorder) WriteHeader(status int) {
	r.status = status
	r.ResponseWriter.WriteHeader(status)
}

func (r *statusRecorder) Write(data []byte) (int, error) {
	if r.status == 0 {
		r.status = http.StatusOK
	}
	n, err := r.ResponseWriter.Write(data)
	r.bytes += n
	return n, err
}

func loggingMiddleware(logger *slog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			recorder := &statusRecorder{ResponseWriter: w}
			next.ServeHTTP(recorder, r)
			status := recorder.status
			if status == 0 {
				status = http.StatusOK
			}
			logger.Info("request completed", "method", r.Method, "path", r.URL.Path, "status", status, "bytes", recorder.bytes, "duration", time.Since(start))
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

func corsMiddleware(allowedOrigins []string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			origin := r.Header.Get("Origin")
			if origin != "" && slices.Contains(allowedOrigins, origin) {
				w.Header().Set("Access-Control-Allow-Origin", origin)
				w.Header().Set("Vary", "Origin")
				w.Header().Set("Access-Control-Allow-Credentials", "true")
				w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
				w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
			}
			if r.Method == http.MethodOptions {
				w.WriteHeader(http.StatusNoContent)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}
