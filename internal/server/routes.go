package server

import (
	"net/http"

	_ "bayt-alhikmah/docs" // Import docs
	customMiddleware "bayt-alhikmah/internal/middleware"

	"github.com/labstack/echo/v4"
	echoSwagger "github.com/swaggo/echo-swagger"
)

func (s *Server) RegisterRoutes() {
	e := s.Echo

	e.GET("/health", s.healthHandler)

	api := e.Group("/api/v1")

	// Swagger
	e.GET("/swagger/*", echoSwagger.WrapHandler)

	// Auth Routes
	s.AuthHandler.RegisterRoutes(api)

	// Protected Routes
	protected := api.Group("")
	protected.Use(customMiddleware.Auth(s.Config.JWTSecret))
	protected.Use(customMiddleware.CasbinMiddleware(s.RBACService))
	s.UserHandler.RegisterRoutes(protected)
}

func (s *Server) healthHandler(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]interface{}{
		"status": "up",
		"db":     s.DB.Health(),
		"redis":  s.Redis.Health(),
	})
}
