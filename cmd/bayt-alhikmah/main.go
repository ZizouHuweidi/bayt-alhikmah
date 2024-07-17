package main

import (
	"github.com/labstack/echo/v4/middleware"
	"net/http"

	"github.com/labstack/echo/v4"
)

func greet(c echo.Context) error {
	return c.String(http.StatusOK, "hello")
}

func main() {
	e := echo.New()
	e.Use(middleware.Logger())
	e.GET("/", greet)

	e.Logger.Fatal(e.Start(":1323"))
}
