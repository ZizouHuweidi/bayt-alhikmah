package httpx

import (
	"net/http"

	"github.com/labstack/echo/v5"
)

type Router interface {
	Get(path string, handler http.HandlerFunc)
	Post(path string, handler http.HandlerFunc)
	Put(path string, handler http.HandlerFunc)
	Delete(path string, handler http.HandlerFunc)
}

type echoRouter interface {
	GET(path string, handler echo.HandlerFunc, middleware ...echo.MiddlewareFunc) echo.RouteInfo
	POST(path string, handler echo.HandlerFunc, middleware ...echo.MiddlewareFunc) echo.RouteInfo
	PUT(path string, handler echo.HandlerFunc, middleware ...echo.MiddlewareFunc) echo.RouteInfo
	DELETE(path string, handler echo.HandlerFunc, middleware ...echo.MiddlewareFunc) echo.RouteInfo
}

type EchoRouter struct {
	router echoRouter
}

func NewEchoRouter(router echoRouter) EchoRouter {
	return EchoRouter{router: router}
}

func (r EchoRouter) Get(path string, handler http.HandlerFunc) {
	r.router.GET(path, EchoHandler(handler))
}

func (r EchoRouter) Post(path string, handler http.HandlerFunc) {
	r.router.POST(path, EchoHandler(handler))
}

func (r EchoRouter) Put(path string, handler http.HandlerFunc) {
	r.router.PUT(path, EchoHandler(handler))
}

func (r EchoRouter) Delete(path string, handler http.HandlerFunc) {
	r.router.DELETE(path, EchoHandler(handler))
}

func EchoHandler(handler http.HandlerFunc) echo.HandlerFunc {
	return func(c *echo.Context) error {
		request := c.Request()
		for _, name := range c.RouteInfo().Parameters {
			request.SetPathValue(name, c.Param(name))
		}
		handler(c.Response(), request)
		return nil
	}
}
