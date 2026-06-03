package server

import (
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v5"
)

type requestValidator struct {
	validator *validator.Validate
}

func newRequestValidator() requestValidator {
	return requestValidator{validator: validator.New()}
}

func (v requestValidator) Validate(value any) error {
	if err := v.validator.Struct(value); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid request body")
	}
	return nil
}
