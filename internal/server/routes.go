package server

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func (s *Server) RegisterRoutes(e *echo.Echo) {
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		// TODO: Change for prod
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete},
		AllowHeaders:     []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderAuthorization},
		AllowCredentials: true,
	}))

	e.GET("/health", s.healthHandler)

	apiV1 := e.Group("/api/v1")

	userGroup := apiV1.Group("/users")
	s.userHandler.RegisterRoutes(userGroup)

	// bookGroup := apiV1.Group("/books")
	// s.bookHandler.RegisterRoutes(bookGroup)
}

func (s *Server) healthHandler(c echo.Context) error {
	// TODO: Add database health check

	return c.JSON(http.StatusOK, map[string]string{"status": "ok"})
}
