package collections

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/gofrs/uuid/v5"
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

type CreateRequest struct {
	Name        string   `json:"name" validate:"required"`
	Description *string  `json:"description,omitempty"`
	IsPublic    bool     `json:"is_public"`
	SourceIDs   []string `json:"source_ids,omitempty"`
}

type UpdateRequest struct {
	Name        *string  `json:"name,omitempty"`
	Description *string  `json:"description,omitempty"`
	IsPublic    *bool    `json:"is_public,omitempty"`
	SourceIDs   []string `json:"source_ids,omitempty"`
}

func (h *Handler) RegisterPublicRoutes(e *echo.Echo) {
	e.GET("/collections", h.List)
	e.GET("/collections/:id", h.GetByID)
}

func (h *Handler) RegisterProtectedRoutes(g *echo.Group) {
	g.POST("/collections", h.Create)
	g.GET("/collections", h.ListOwn)
	g.PUT("/collections/:id", h.Update)
	g.DELETE("/collections/:id", h.Delete)
}

func (h *Handler) Create(c *echo.Context) error {
	userID, ok := auth.UserID(c)
	if !ok {
		return echo.NewHTTPError(http.StatusUnauthorized, "authentication required")
	}

	var req CreateRequest
	if err := echox.BindAndValidate(c, &req); err != nil {
		return err
	}
	sourceIDs, err := parseUUIDs(req.SourceIDs, "source_ids")
	if err != nil {
		return err
	}

	collection, err := h.service.Create(c.Request().Context(), CreateCollectionParams{UserID: userID, Name: req.Name, Description: req.Description, IsPublic: req.IsPublic, SourceIDs: sourceIDs})
	if errors.Is(err, ErrInvalidCollection) {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid collection")
	}
	if errors.Is(err, ErrSourceNotFound) {
		return echo.NewHTTPError(http.StatusNotFound, "source not found")
	}
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to create collection")
	}

	return c.JSON(http.StatusCreated, collection)
}

func (h *Handler) GetByID(c *echo.Context) error {
	id, err := echox.ParamUUID(c, "id", "collection ID")
	if err != nil {
		return err
	}

	collection, err := h.service.GetByID(c.Request().Context(), id)
	if errors.Is(err, ErrCollectionNotFound) {
		return echo.NewHTTPError(http.StatusNotFound, "collection not found")
	}
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to get collection")
	}
	if !collection.IsPublic {
		return echo.NewHTTPError(http.StatusForbidden, "collection is private")
	}

	return c.JSON(http.StatusOK, collection)
}

func (h *Handler) List(c *echo.Context) error {
	limit, offset := echox.Pagination(c)
	userIDStr := c.QueryParam("user_id")
	if userIDStr == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "user_id required")
	}
	userID, err := uuid.FromString(userIDStr)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid user_id")
	}

	collections, err := h.service.ListPublicByUser(c.Request().Context(), userID, limit, offset)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to list collections")
	}

	return c.JSON(http.StatusOK, collections)
}

func (h *Handler) ListOwn(c *echo.Context) error {
	userID, ok := auth.UserID(c)
	if !ok {
		return echo.NewHTTPError(http.StatusUnauthorized, "authentication required")
	}
	limit, offset := echox.Pagination(c)

	collections, err := h.service.ListByUser(c.Request().Context(), userID, limit, offset)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to list collections")
	}

	return c.JSON(http.StatusOK, collections)
}

func (h *Handler) Update(c *echo.Context) error {
	userID, ok := auth.UserID(c)
	if !ok {
		return echo.NewHTTPError(http.StatusUnauthorized, "authentication required")
	}
	id, err := echox.ParamUUID(c, "id", "collection ID")
	if err != nil {
		return err
	}

	existing, err := h.service.GetByID(c.Request().Context(), id)
	if errors.Is(err, ErrCollectionNotFound) {
		return echo.NewHTTPError(http.StatusNotFound, "collection not found")
	}
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to get collection")
	}
	if existing.UserID != userID {
		return echo.NewHTTPError(http.StatusForbidden, "cannot update another user's collection")
	}

	var req UpdateRequest
	if err := echox.BindAndValidate(c, &req); err != nil {
		return err
	}
	sourceIDs, err := parseUUIDs(req.SourceIDs, "source_ids")
	if err != nil {
		return err
	}

	collection, err := h.service.Update(c.Request().Context(), id, UpdateCollectionParams{Name: req.Name, Description: req.Description, IsPublic: req.IsPublic, SourceIDs: sourceIDs})
	if errors.Is(err, ErrInvalidCollection) {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid collection")
	}
	if errors.Is(err, ErrSourceNotFound) {
		return echo.NewHTTPError(http.StatusNotFound, "source not found")
	}
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to update collection")
	}

	return c.JSON(http.StatusOK, collection)
}

func (h *Handler) Delete(c *echo.Context) error {
	userID, ok := auth.UserID(c)
	if !ok {
		return echo.NewHTTPError(http.StatusUnauthorized, "authentication required")
	}
	id, err := echox.ParamUUID(c, "id", "collection ID")
	if err != nil {
		return err
	}

	existing, err := h.service.GetByID(c.Request().Context(), id)
	if errors.Is(err, ErrCollectionNotFound) {
		return echo.NewHTTPError(http.StatusNotFound, "collection not found")
	}
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to get collection")
	}
	if existing.UserID != userID {
		return echo.NewHTTPError(http.StatusForbidden, "cannot delete another user's collection")
	}

	if err := h.service.Delete(c.Request().Context(), id); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to delete collection")
	}
	return c.NoContent(http.StatusNoContent)
}

func parseUUIDs(values []string, label string) ([]uuid.UUID, error) {
	if values == nil {
		return nil, nil
	}
	ids := make([]uuid.UUID, 0, len(values))
	for _, value := range values {
		id, err := uuid.FromString(value)
		if err != nil {
			return nil, echo.NewHTTPError(http.StatusBadRequest, "invalid "+label)
		}
		ids = append(ids, id)
	}
	return ids, nil
}
