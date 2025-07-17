package server

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"github.com/zizouhuweidi/bayt-alhikmah/internal/database"
	"github.com/zizouhuweidi/bayt-alhikmah/internal/handlers"
	"github.com/zizouhuweidi/bayt-alhikmah/internal/repository"
)

func (s *Server) RegisterRoutes() http.Handler {
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:     []string{"https://*", "http://*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},
		AllowHeaders:     []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	bookRepo := repository.NewBookRepository(database.New().DB())
	bookHandler := handlers.NewBookHandler(bookRepo)

	// Setup routes.
	//
	e.GET("/books", bookHandler.ListBooks)
	e.GET("/health", s.healthHandler)
	e.GET("/", s.HelloWorldHandler)
	return e
}

func (s *Server) HelloWorldHandler(c echo.Context) error {
	resp := map[string]string{
		"message": "hello, world!",
	}

	return c.JSON(http.StatusOK, resp)
}

func (s *Server) healthHandler(c echo.Context) error {
	return c.JSON(http.StatusOK, s.db.Health())
}
