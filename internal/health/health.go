package health

import (
	"context"
	"net/http"
	"time"

	"github.com/labstack/echo/v5"
	"github.com/zizouhuweidi/maktaba/internal/db"
)

type Response struct {
	Status string `json:"status"`
}

func Handler(c *echo.Context) error {
	return c.JSON(http.StatusOK, Response{Status: "ok"})
}

func ReadyHandler(database *db.DB) echo.HandlerFunc {
	return func(c *echo.Context) error {
		ctx, cancel := context.WithTimeout(c.Request().Context(), 2*time.Second)
		defer cancel()

		if err := database.Ping(ctx); err != nil {
			return c.JSON(http.StatusServiceUnavailable, Response{Status: "unavailable"})
		}

		return c.JSON(http.StatusOK, Response{Status: "ready"})
	}
}
