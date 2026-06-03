package library

import (
	"errors"
	"log/slog"
	"net/http"
	"time"

	"github.com/gofrs/uuid/v5"
	"github.com/labstack/echo/v5"
	"github.com/zizouhuweidi/maktaba/internal/auth"
	"github.com/zizouhuweidi/maktaba/internal/echox"
)

type Handler struct {
	service *Service
	logger  *slog.Logger
}

type createRequest struct {
	SourceID      string        `json:"source_id" validate:"required"`
	Status        string        `json:"status" validate:"required"`
	ProgressValue *int          `json:"progress_value,omitempty"`
	ProgressUnit  *ProgressUnit `json:"progress_unit,omitempty"`
	Visibility    string        `json:"visibility,omitempty"`
	StartedAt     *time.Time    `json:"started_at,omitempty"`
	CompletedAt   *time.Time    `json:"completed_at,omitempty"`
}

type updateRequest struct {
	Status        *string       `json:"status,omitempty"`
	ProgressValue *int          `json:"progress_value,omitempty"`
	ProgressUnit  *ProgressUnit `json:"progress_unit,omitempty"`
	Visibility    *string       `json:"visibility,omitempty"`
	StartedAt     *time.Time    `json:"started_at,omitempty"`
	CompletedAt   *time.Time    `json:"completed_at,omitempty"`
}

func NewHandler(service *Service, logger *slog.Logger) *Handler {
	return &Handler{service: service, logger: logger}
}

func (h *Handler) RegisterPublicRoutes(e *echo.Echo) {
	e.GET("/users/:user/library", h.ListPublicLibrary)
	e.GET("/users/:user/library/with-sources", h.ListPublicLibraryWithSources)
}

func (h *Handler) RegisterProtectedRoutes(g *echo.Group) {
	g.POST("/library/items", h.Create)
	g.GET("/library/items", h.ListMine)
	g.GET("/library/items/with-sources", h.ListMineWithSources)
	g.GET("/library/items/:id", h.GetMine)
	g.PUT("/library/items/:id", h.Update)
	g.DELETE("/library/items/:id", h.Delete)
}

func (h *Handler) Create(c *echo.Context) error {
	userID, ok := auth.UserID(c)
	if !ok {
		return echo.NewHTTPError(http.StatusUnauthorized, "authentication required")
	}

	var req createRequest
	if err := echox.BindAndValidate(c, &req); err != nil {
		return err
	}
	sourceID, err := uuid.FromString(req.SourceID)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid source_id")
	}
	visibility := Visibility(req.Visibility)
	if visibility == "" {
		visibility = VisibilityPrivate
	}

	item, err := h.service.Create(c.Request().Context(), CreateItemParams{
		UserID:        userID,
		SourceID:      sourceID,
		Status:        Status(req.Status),
		ProgressValue: req.ProgressValue,
		ProgressUnit:  req.ProgressUnit,
		Visibility:    visibility,
		StartedAt:     req.StartedAt,
		CompletedAt:   req.CompletedAt,
	})
	if errors.Is(err, ErrInvalidItem) {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid library item")
	}
	if errors.Is(err, ErrSourceNotFound) {
		return echo.NewHTTPError(http.StatusNotFound, "source not found")
	}
	if errors.Is(err, ErrItemExists) {
		return echo.NewHTTPError(http.StatusConflict, "source already exists in library")
	}
	if errors.Is(err, ErrLibraryConflict) {
		return echo.NewHTTPError(http.StatusConflict, "library item conflict")
	}
	if err != nil {
		h.logger.Error("failed to create library item", "error", err)
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to create library item")
	}

	return c.JSON(http.StatusCreated, item)
}

func (h *Handler) ListMine(c *echo.Context) error {
	userID, ok := auth.UserID(c)
	if !ok {
		return echo.NewHTTPError(http.StatusUnauthorized, "authentication required")
	}
	limit, offset := echox.Pagination(c)
	items, err := h.service.ListByUser(c.Request().Context(), userID, limit, offset)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to list library items")
	}
	return c.JSON(http.StatusOK, items)
}

func (h *Handler) ListMineWithSources(c *echo.Context) error {
	userID, ok := auth.UserID(c)
	if !ok {
		return echo.NewHTTPError(http.StatusUnauthorized, "authentication required")
	}
	limit, offset := echox.Pagination(c)
	items, err := h.service.ListByUserWithSources(c.Request().Context(), userID, limit, offset)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to list library items")
	}
	return c.JSON(http.StatusOK, items)
}

func (h *Handler) GetMine(c *echo.Context) error {
	userID, ok := auth.UserID(c)
	if !ok {
		return echo.NewHTTPError(http.StatusUnauthorized, "authentication required")
	}
	item, err := h.getOwnedItem(c, userID)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, item)
}

func (h *Handler) Update(c *echo.Context) error {
	userID, ok := auth.UserID(c)
	if !ok {
		return echo.NewHTTPError(http.StatusUnauthorized, "authentication required")
	}
	item, err := h.getOwnedItem(c, userID)
	if err != nil {
		return err
	}

	var req updateRequest
	if err := echox.BindAndValidate(c, &req); err != nil {
		return err
	}
	var status *Status
	if req.Status != nil {
		parsed := Status(*req.Status)
		status = &parsed
	}
	var visibility *Visibility
	if req.Visibility != nil {
		parsed := Visibility(*req.Visibility)
		visibility = &parsed
	}

	updated, err := h.service.Update(c.Request().Context(), item.ID, UpdateItemParams{
		Status:        status,
		ProgressValue: req.ProgressValue,
		ProgressUnit:  req.ProgressUnit,
		Visibility:    visibility,
		StartedAt:     req.StartedAt,
		CompletedAt:   req.CompletedAt,
	})
	if errors.Is(err, ErrInvalidItem) {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid library item")
	}
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to update library item")
	}
	return c.JSON(http.StatusOK, updated)
}

func (h *Handler) Delete(c *echo.Context) error {
	userID, ok := auth.UserID(c)
	if !ok {
		return echo.NewHTTPError(http.StatusUnauthorized, "authentication required")
	}
	item, err := h.getOwnedItem(c, userID)
	if err != nil {
		return err
	}
	if err := h.service.Delete(c.Request().Context(), item.ID); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to delete library item")
	}
	return c.NoContent(http.StatusNoContent)
}

func (h *Handler) ListPublicLibrary(c *echo.Context) error {
	limit, offset := echox.Pagination(c)
	user := c.Param("user")

	userID, err := uuid.FromString(user)
	if err == nil {
		items, err := h.service.ListPublicByUser(c.Request().Context(), userID, limit, offset)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, "failed to list public library items")
		}
		return c.JSON(http.StatusOK, items)
	}

	items, err := h.service.ListPublicByUsername(c.Request().Context(), user, limit, offset)
	if errors.Is(err, ErrInvalidUser) {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid username")
	}
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to list public library items")
	}
	return c.JSON(http.StatusOK, items)
}

func (h *Handler) ListPublicLibraryWithSources(c *echo.Context) error {
	limit, offset := echox.Pagination(c)
	items, err := h.service.ListPublicByUsernameWithSources(c.Request().Context(), c.Param("user"), limit, offset)
	if errors.Is(err, ErrInvalidUser) {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid username")
	}
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to list public library items")
	}
	return c.JSON(http.StatusOK, items)
}

func (h *Handler) getOwnedItem(c *echo.Context, userID uuid.UUID) (*Item, error) {
	id, err := echox.ParamUUID(c, "id", "library item ID")
	if err != nil {
		return nil, err
	}
	item, err := h.service.GetByID(c.Request().Context(), id)
	if errors.Is(err, ErrItemNotFound) {
		return nil, echo.NewHTTPError(http.StatusNotFound, "library item not found")
	}
	if err != nil {
		return nil, echo.NewHTTPError(http.StatusInternalServerError, "failed to get library item")
	}
	if item.UserID != userID {
		return nil, echo.NewHTTPError(http.StatusForbidden, "cannot access another user's library item")
	}
	return item, nil
}
