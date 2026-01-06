package notes

import (
	"log/slog"
	"net/http"
	"strconv"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type Handler struct {
	service *Service
	logger  *slog.Logger
}

func NewHandler(service *Service, logger *slog.Logger) *Handler {
	return &Handler{
		service: service,
		logger:  logger,
	}
}

type CreateRequest struct {
	UserID      string   `json:"user_id"`
	SourceID    *string  `json:"source_id,omitempty"`
	Content     string   `json:"content"`
	ContentType string   `json:"content_type"`
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

func (h *Handler) RegisterRoutes(e *echo.Echo) {
	g := e.Group("/notes")
	g.POST("", h.Create)
	g.GET("", h.List)
	g.GET("/:id", h.GetByID)
	g.PUT("/:id", h.Update)
	g.DELETE("/:id", h.Delete)
}

func (h *Handler) Create(c echo.Context) error {
	var req CreateRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request body"})
	}

	userID, err := uuid.Parse(req.UserID)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid user_id"})
	}

	var sourceID *uuid.UUID
	if req.SourceID != nil {
		id, err := uuid.Parse(*req.SourceID)
		if err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid source_id"})
		}
		sourceID = &id
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
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to create note"})
	}

	return c.JSON(http.StatusCreated, note)
}

func (h *Handler) GetByID(c echo.Context) error {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid note ID"})
	}

	note, err := h.service.GetByID(c.Request().Context(), id)
	if err == ErrNoteNotFound {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "note not found"})
	}
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to get note"})
	}

	return c.JSON(http.StatusOK, note)
}

func (h *Handler) List(c echo.Context) error {
	limit, offset := h.parsePagination(c)
	userIDStr := c.QueryParam("user_id")
	sourceIDStr := c.QueryParam("source_id")
	publicOnly := c.QueryParam("public") == "true"

	var notes []*Note
	var err error

	if userIDStr != "" {
		userID, err := uuid.Parse(userIDStr)
		if err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid user_id"})
		}
		notes, err = h.service.ListByUser(c.Request().Context(), userID, limit, offset)
		if err != nil {
			return err
		}

	} else if sourceIDStr != "" {
		sourceID, err := uuid.Parse(sourceIDStr)
		if err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid source_id"})
		}
		notes, err = h.service.ListBySource(c.Request().Context(), sourceID, limit, offset)
		if err != nil {
			return err
		}
	} else if publicOnly {
		notes, err = h.service.ListPublic(c.Request().Context(), limit, offset)
	} else {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "user_id, source_id, or public=true required"})
	}

	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to list notes"})
	}

	return c.JSON(http.StatusOK, notes)
}

func (h *Handler) Update(c echo.Context) error {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid note ID"})
	}

	var req UpdateRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request body"})
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

	if err == ErrNoteNotFound {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "note not found"})
	}
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to update note"})
	}

	return c.JSON(http.StatusOK, note)
}

func (h *Handler) Delete(c echo.Context) error {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid note ID"})
	}

	if err := h.service.Delete(c.Request().Context(), id); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to delete note"})
	}

	return c.NoContent(http.StatusNoContent)
}

func (h *Handler) parsePagination(c echo.Context) (limit, offset int) {
	limit = 100
	offset = 0

	if l := c.QueryParam("limit"); l != "" {
		if parsed, err := strconv.Atoi(l); err == nil && parsed > 0 {
			limit = parsed
		}
	}

	if o := c.QueryParam("offset"); o != "" {
		if parsed, err := strconv.Atoi(o); err == nil && parsed >= 0 {
			offset = parsed
		}
	}

	return limit, offset
}
