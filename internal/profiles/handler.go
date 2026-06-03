package profiles

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/labstack/echo/v5"
	"github.com/zizouhuweidi/maktaba/internal/auth"
	"github.com/zizouhuweidi/maktaba/internal/echox"
)

type Handler struct {
	service *Service
	logger  *slog.Logger
}

func NewHandler(service *Service, logger *slog.Logger) *Handler {
	return &Handler{service: service, logger: logger}
}

type UpdateRequest struct {
	DisplayName   *string `json:"display_name,omitempty"`
	Bio           *string `json:"bio,omitempty"`
	PublicProfile *bool   `json:"public_profile,omitempty"`
}

func (h *Handler) RegisterPublicRoutes(e *echo.Echo) {
	e.GET("/users/:username/profile", h.GetPublicByUsername)
}

func (h *Handler) RegisterProtectedRoutes(g *echo.Group) {
	g.GET("/profile", h.GetOwn)
	g.PUT("/profile", h.Update)
}

func (h *Handler) GetOwn(c *echo.Context) error {
	userID, ok := auth.UserID(c)
	if !ok {
		return echo.NewHTTPError(http.StatusUnauthorized, "authentication required")
	}

	profile, err := h.service.GetOwn(c.Request().Context(), userID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to get profile")
	}

	return c.JSON(http.StatusOK, profile)
}

func (h *Handler) GetPublicByUsername(c *echo.Context) error {
	profile, err := h.service.GetPublicByUsername(c.Request().Context(), c.Param("username"))
	if errors.Is(err, ErrInvalidProfile) {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid username")
	}
	if errors.Is(err, ErrProfileNotFound) {
		return echo.NewHTTPError(http.StatusNotFound, "profile not found")
	}
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to get profile")
	}

	return c.JSON(http.StatusOK, profile)
}

func (h *Handler) Update(c *echo.Context) error {
	userID, ok := auth.UserID(c)
	if !ok {
		return echo.NewHTTPError(http.StatusUnauthorized, "authentication required")
	}

	var req UpdateRequest
	if err := echox.BindAndValidate(c, &req); err != nil {
		return err
	}

	profile, err := h.service.Update(c.Request().Context(), UpdateProfileParams{
		UserID:        userID,
		DisplayName:   req.DisplayName,
		Bio:           req.Bio,
		PublicProfile: req.PublicProfile,
	})
	if errors.Is(err, ErrInvalidProfile) {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid profile")
	}
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to update profile")
	}

	return c.JSON(http.StatusOK, profile)
}
