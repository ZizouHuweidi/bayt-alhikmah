package echox

import (
	"net/http"
	"strconv"

	"github.com/gofrs/uuid/v5"
	"github.com/labstack/echo/v5"
)

type ErrorResponse struct {
	Error string `json:"error"`
}

func BindAndValidate(c *echo.Context, dst any) error {
	if err := c.Bind(dst); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid request body")
	}
	if err := c.Validate(dst); err != nil {
		return err
	}
	return nil
}

func Pagination(c *echo.Context) (limit, offset int) {
	limit = 100
	if value := c.QueryParam("limit"); value != "" {
		parsed, err := strconv.Atoi(value)
		if err == nil && parsed > 0 {
			limit = parsed
		}
	}
	if limit > 100 {
		limit = 100
	}

	if value := c.QueryParam("offset"); value != "" {
		parsed, err := strconv.Atoi(value)
		if err == nil && parsed >= 0 {
			offset = parsed
		}
	}
	return limit, offset
}

func ParamUUID(c *echo.Context, name, label string) (uuid.UUID, error) {
	id, err := uuid.FromString(c.Param(name))
	if err != nil {
		return uuid.Nil, echo.NewHTTPError(http.StatusBadRequest, "invalid "+label)
	}
	return id, nil
}

func OptionalUUID(value *string, label string) (*uuid.UUID, error) {
	if value == nil {
		return nil, nil
	}
	id, err := uuid.FromString(*value)
	if err != nil {
		return nil, echo.NewHTTPError(http.StatusBadRequest, "invalid "+label)
	}
	return &id, nil
}
