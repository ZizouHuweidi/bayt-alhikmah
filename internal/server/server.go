package server

import (
	"fmt"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
)

type UserHandler interface {
	RegisterRoutes(g *echo.Group)
}

type Config struct {
	Port        int
	UserHandler UserHandler
	// BookHandler BookHandler
}

type Server struct {
	port        int
	userHandler UserHandler
	// bookHandler BookHandler ...
}

func NewServer(cfg Config) *http.Server {
	server := &Server{
		port:        cfg.Port,
		userHandler: cfg.UserHandler,
	}

	e := echo.New()

	server.RegisterRoutes(e)

	httpServer := &http.Server{
		Addr:         fmt.Sprintf(":%d", server.port),
		Handler:      e,
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	return httpServer
}
