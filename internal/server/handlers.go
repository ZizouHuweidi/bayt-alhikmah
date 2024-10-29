package server

import (
	"context"
	"net/http"

	"github.com/labstack/echo/v4"
)

func (s *Server) HelloWorldHandler(c echo.Context) error {
	resp := map[string]string{
		"message": "Hello World",
	}

	return c.JSON(http.StatusOK, resp)
}

func (s *Server) healthHandler(c echo.Context) error {
	return c.JSON(http.StatusOK, s.db.Health(context.Background()))
}

func (s *Server) GetAllBooks(c echo.Context) error {
	resp, err := s.db.Queries().GetBooks(context.Background())
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, resp)
}
