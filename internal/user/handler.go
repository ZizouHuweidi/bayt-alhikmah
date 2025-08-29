package user

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

type Handler struct {
	service Service
}

func NewHandler(svc Service) *Handler {
	return &Handler{service: svc}
}

func (h *Handler) RegisterRoutes(g *echo.Group) {
	g.POST("/register", h.handleRegister)
}

type registerRequest struct {
	Username string `json:"username" validate:"required"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
}

func (h *Handler) handleRegister(c echo.Context) error {
	var req registerRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request body"})
	}

	// TODO: add validation using go-playground/validator to check the request struct tags.

	user, err := h.service.Register(c.Request().Context(), req.Username, req.Email, req.Password)
	if err != nil {
		// This is where you could check for specific domain errors
		// and return different status codes (e.g., 409 Conflict for duplicate email).
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "could not register user"})
	}

	return c.JSON(http.StatusCreated, user)
}
