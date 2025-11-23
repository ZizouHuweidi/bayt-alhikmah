package server

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"time"

	"yalla-go/internal/auth"
	"yalla-go/internal/config"
	"yalla-go/internal/database"
	customMiddleware "yalla-go/internal/middleware"
	"yalla-go/internal/rbac"
	"yalla-go/internal/redis"
	"yalla-go/internal/user"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"go.opentelemetry.io/contrib/instrumentation/github.com/labstack/echo/otelecho"
)

type Server struct {
	Echo        *echo.Echo
	Config      *config.Config
	DB          database.Service
	Redis       *redis.Client
	AuthHandler *auth.Handler
	UserHandler *user.Handler
	RBACService *rbac.Service
}

func NewServer(
	cfg *config.Config,
	db database.Service,
	redis *redis.Client,
	authHandler *auth.Handler,
	userHandler *user.Handler,
) *Server {
	// Init RBAC
	rbacService, err := rbac.NewService("config/rbac_model.conf", "config/rbac_policy.csv")
	if err != nil {
		// In production, we might want to fail hard or log error
		// For starter, let's panic or log fatal
		panic(err)
	}

	e := echo.New()
	e.HideBanner = true

	// Logger
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	e.Use(customMiddleware.SlogLogger(logger))

	// OTel
	e.Use(otelecho.Middleware("go-backend-template"))

	// Recover
	e.Use(middleware.Recover())

	// CORS
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete, http.MethodOptions},
		AllowHeaders:     []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderAuthorization},
		AllowCredentials: true,
	}))

	// Rate Limit (Global - 100 req/min)
	e.Use(customMiddleware.RateLimit(redis, 100, 1*time.Minute))

	s := &Server{
		Echo:        e,
		Config:      cfg,
		DB:          db,
		Redis:       redis,
		AuthHandler: authHandler,
		UserHandler: userHandler,
		RBACService: rbacService,
	}

	s.RegisterRoutes()

	return s
}

func (s *Server) Start() error {
	return s.Echo.Start(fmt.Sprintf(":%d", s.Config.Port))
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.Echo.Shutdown(ctx)
}
