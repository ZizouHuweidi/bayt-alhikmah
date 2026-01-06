package sources

import (
	"net/http"
	"strconv"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"log/slog"
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
	Title       string   `json:"title"`
	Subtitle    *string  `json:"subtitle,omitempty"`
	Type        string   `json:"type"`
	Description *string  `json:"description,omitempty"`
	AuthorID    *string  `json:"author_id,omitempty"`
	Publisher   *string  `json:"publisher,omitempty"`
	ISBN        *string  `json:"isbn,omitempty"`
	DOI         *string  `json:"doi,omitempty"`
	URL         *string  `json:"url,omitempty"`
	ExternalID  *string  `json:"external_id,omitempty"`
	Tags        []string `json:"tags,omitempty"`
	PublishedAt *string  `json:"published_at,omitempty"`
}

type UpdateRequest struct {
	Title       *string  `json:"title,omitempty"`
	Subtitle    *string  `json:"subtitle,omitempty"`
	Type        *string  `json:"type,omitempty"`
	Description *string  `json:"description,omitempty"`
	AuthorID    *string  `json:"author_id,omitempty"`
	Publisher   *string  `json:"publisher,omitempty"`
	ISBN        *string  `json:"isbn,omitempty"`
	DOI         *string  `json:"doi,omitempty"`
	URL         *string  `json:"url,omitempty"`
	ExternalID  *string  `json:"external_id,omitempty"`
	Tags        []string `json:"tags,omitempty"`
	PublishedAt *string  `json:"published_at,omitempty"`
}

func (h *Handler) RegisterRoutes(e *echo.Echo) {
	g := e.Group("/sources")
	g.POST("", h.Create)
	g.GET("", h.List)
	g.GET("/search", h.Search)
	g.GET("/:id", h.GetByID)
	g.PUT("/:id", h.Update)
	g.DELETE("/:id", h.Delete)
}

func (h *Handler) Create(c echo.Context) error {
	var req CreateRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request body"})
	}

	var authorID *uuid.UUID
	if req.AuthorID != nil {
		id, err := uuid.Parse(*req.AuthorID)
		if err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid author_id"})
		}
		authorID = &id
	}

	source, err := h.service.Create(c.Request().Context(), CreateSourceParams{
		Title:       req.Title,
		Subtitle:    req.Subtitle,
		Type:        SourceType(req.Type),
		Description: req.Description,
		AuthorID:    authorID,
		Publisher:   req.Publisher,
		ISBN:        req.ISBN,
		DOI:         req.DOI,
		URL:         req.URL,
		ExternalID:  req.ExternalID,
		Tags:        req.Tags,
	})

	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to create source"})
	}

	return c.JSON(http.StatusCreated, source)
}

func (h *Handler) GetByID(c echo.Context) error {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid source ID"})
	}

	source, err := h.service.GetByID(c.Request().Context(), id)
	if err == ErrSourceNotFound {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "source not found"})
	}
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to get source"})
	}

	return c.JSON(http.StatusOK, source)
}

func (h *Handler) List(c echo.Context) error {
	limit, offset := h.parsePagination(c)
	sourceType := c.QueryParam("type")

	var sources []*Source
	var err error

	if sourceType != "" {
		sources, err = h.service.ListByType(c.Request().Context(), SourceType(sourceType), limit, offset)
	} else {
		sources, err = h.service.List(c.Request().Context(), limit, offset)
	}

	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to list sources"})
	}

	return c.JSON(http.StatusOK, sources)
}

func (h *Handler) Search(c echo.Context) error {
	query := c.QueryParam("q")
	if query == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "search query required"})
	}

	limit, offset := h.parsePagination(c)
	sources, err := h.service.Search(c.Request().Context(), query, limit, offset)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to search sources"})
	}

	return c.JSON(http.StatusOK, sources)
}

func (h *Handler) Update(c echo.Context) error {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid source ID"})
	}

	var req UpdateRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request body"})
	}

	var authorID *uuid.UUID
	if req.AuthorID != nil {
		aid, err := uuid.Parse(*req.AuthorID)
		if err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid author_id"})
		}
		authorID = &aid
	}

	var sourceType *SourceType
	if req.Type != nil {
		st := SourceType(*req.Type)
		sourceType = &st
	}

	source, err := h.service.Update(c.Request().Context(), id, UpdateSourceParams{
		Title:       req.Title,
		Subtitle:    req.Subtitle,
		Type:        sourceType,
		Description: req.Description,
		AuthorID:    authorID,
		Publisher:   req.Publisher,
		ISBN:        req.ISBN,
		DOI:         req.DOI,
		URL:         req.URL,
		ExternalID:  req.ExternalID,
		Tags:        req.Tags,
	})

	if err == ErrSourceNotFound {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "source not found"})
	}
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to update source"})
	}

	return c.JSON(http.StatusOK, source)
}

func (h *Handler) Delete(c echo.Context) error {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid source ID"})
	}

	if err := h.service.Delete(c.Request().Context(), id); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to delete source"})
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
