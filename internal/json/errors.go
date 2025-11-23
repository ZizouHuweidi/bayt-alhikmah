package json

import (
	"net/http"

	"yalla-go/internal/response"

	"github.com/labstack/echo/v4"
)

func BadRequest(c echo.Context, err error) error {
	return response.ErrorJSON(c, http.StatusBadRequest, "BAD_REQUEST", err.Error(), nil)
}

func InternalServerError(c echo.Context, err error) error {
	return response.ErrorJSON(c, http.StatusInternalServerError, "INTERNAL_SERVER_ERROR", "Something went wrong", nil)
}

func Unauthorized(c echo.Context, message string) error {
	return response.ErrorJSON(c, http.StatusUnauthorized, "UNAUTHORIZED", message, nil)
}

func NotFound(c echo.Context, message string) error {
	return response.ErrorJSON(c, http.StatusNotFound, "NOT_FOUND", message, nil)
}

func Forbidden(c echo.Context, message string) error {
	return response.ErrorJSON(c, http.StatusForbidden, "FORBIDDEN", message, nil)
}
