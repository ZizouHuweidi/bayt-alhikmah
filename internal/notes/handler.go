package notes

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
	SourceID    *string  `json:"source_id,omitempty"`
	Content     string   `json:"content" validate:"required"`
	ContentType string   `json:"content_type" validate:"required"`
	IsPublic    bool     `json:"is_public"`
	Annotations []string `json:"annotations,omitempty"`
	Tags        []string `json:"tags,omitempty"`
}

type UpdateRequest struct {
	Content     *string  `json:"content,omitempty"`
	ContentType *string  `json:"content_type,omitempty"`
	IsPublic    *bool    `json:"is_public,omitempty"`
	Annotations []string `json:"annotations,omitempty"`
	Tags        []string `json:"tags,omitempty"`
}

func (h *Handler) RegisterPublicRoutes(e *echo.Echo) {
	e.GET("/notes", h.List)
	e.GET("/notes/:id", h.GetByID)
}

func (h *Handler) RegisterProtectedRoutes(g *echo.Group) {
	g.POST("/notes", h.Create)
	g.GET("/notes", h.ListMine)
	g.PUT("/notes/:id", h.Update)
	g.DELETE("/notes/:id", h.Delete)
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

	sourceID, err := echox.OptionalUUID(req.SourceID, "source_id")
	if err != nil {
		return err
	}

	note, err := h.service.Create(c.Request().Context(), CreateNoteParams{
		UserID:      userID,
		SourceID:    sourceID,
		Content:     req.Content,
		ContentType: ContentType(req.ContentType),
		IsPublic:    req.IsPublic,
		Annotations: req.Annotations,
		Tags:        req.Tags,
	})
	if errors.Is(err, ErrInvalidNote) {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid note")
	}
	if errors.Is(err, ErrSourceNotFound) {
		return echo.NewHTTPError(http.StatusNotFound, "source not found")
	}
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to create note")
	}

	return c.JSON(http.StatusCreated, note)
}

func (h *Handler) GetByID(c *echo.Context) error {
	id, err := echox.ParamUUID(c, "id", "note ID")
	if err != nil {
		return err
	}

	note, err := h.service.GetByID(c.Request().Context(), id)
	if errors.Is(err, ErrNoteNotFound) {
		return echo.NewHTTPError(http.StatusNotFound, "note not found")
	}
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to get note")
	}
	if !note.IsPublic {
		return echo.NewHTTPError(http.StatusForbidden, "note is private")
	}

	return c.JSON(http.StatusOK, note)
}

func (h *Handler) List(c *echo.Context) error {
	limit, offset := echox.Pagination(c)
	userIDStr := c.QueryParam("user_id")
	sourceIDStr := c.QueryParam("source_id")
	publicOnly := c.QueryParam("public") == "true"

	var result []*Note
	var err error
	if userIDStr != "" {
		userID, parseErr := uuid.FromString(userIDStr)
		if parseErr != nil {
			return echo.NewHTTPError(http.StatusBadRequest, "invalid user_id")
		}
		result, err = h.service.ListPublicByUser(c.Request().Context(), userID, limit, offset)
	} else if sourceIDStr != "" {
		sourceID, parseErr := uuid.FromString(sourceIDStr)
		if parseErr != nil {
			return echo.NewHTTPError(http.StatusBadRequest, "invalid source_id")
		}
		result, err = h.service.ListPublicBySource(c.Request().Context(), sourceID, limit, offset)
	} else if publicOnly {
		result, err = h.service.ListPublic(c.Request().Context(), limit, offset)
	} else {
		return echo.NewHTTPError(http.StatusBadRequest, "user_id, source_id, or public=true required")
	}
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to list notes")
	}

	return c.JSON(http.StatusOK, result)
}

func (h *Handler) ListMine(c *echo.Context) error {
	userID, ok := auth.UserID(c)
	if !ok {
		return echo.NewHTTPError(http.StatusUnauthorized, "authentication required")
	}
	limit, offset := echox.Pagination(c)

	notes, err := h.service.ListByUser(c.Request().Context(), userID, limit, offset)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to list notes")
	}

	return c.JSON(http.StatusOK, notes)
}

func (h *Handler) Update(c *echo.Context) error {
	userID, ok := auth.UserID(c)
	if !ok {
		return echo.NewHTTPError(http.StatusUnauthorized, "authentication required")
	}
	id, err := echox.ParamUUID(c, "id", "note ID")
	if err != nil {
		return err
	}

	existing, err := h.service.GetByID(c.Request().Context(), id)
	if errors.Is(err, ErrNoteNotFound) {
		return echo.NewHTTPError(http.StatusNotFound, "note not found")
	}
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to get note")
	}
	if existing.UserID != userID {
		return echo.NewHTTPError(http.StatusForbidden, "cannot update another user's note")
	}

	var req UpdateRequest
	if err := echox.BindAndValidate(c, &req); err != nil {
		return err
	}

	var contentType *ContentType
	if req.ContentType != nil {
		ct := ContentType(*req.ContentType)
		contentType = &ct
	}

	note, err := h.service.Update(c.Request().Context(), id, UpdateNoteParams{
		Content:     req.Content,
		ContentType: contentType,
		IsPublic:    req.IsPublic,
		Annotations: req.Annotations,
		Tags:        req.Tags,
	})
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to update note")
	}

	return c.JSON(http.StatusOK, note)
}

func (h *Handler) Delete(c *echo.Context) error {
	userID, ok := auth.UserID(c)
	if !ok {
		return echo.NewHTTPError(http.StatusUnauthorized, "authentication required")
	}
	id, err := echox.ParamUUID(c, "id", "note ID")
	if err != nil {
		return err
	}

	existing, err := h.service.GetByID(c.Request().Context(), id)
	if errors.Is(err, ErrNoteNotFound) {
		return echo.NewHTTPError(http.StatusNotFound, "note not found")
	}
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to get note")
	}
	if existing.UserID != userID {
		return echo.NewHTTPError(http.StatusForbidden, "cannot delete another user's note")
	}

	if err := h.service.Delete(c.Request().Context(), id); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to delete note")
	}
	return c.NoContent(http.StatusNoContent)
}
