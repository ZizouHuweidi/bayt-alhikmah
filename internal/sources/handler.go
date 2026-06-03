package sources

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/labstack/echo/v5"
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
	Title       string   `json:"title" validate:"required"`
	Subtitle    *string  `json:"subtitle,omitempty"`
	Type        string   `json:"type" validate:"required"`
	Description *string  `json:"description,omitempty"`
	Publisher   *string  `json:"publisher,omitempty"`
	ISBN        *string  `json:"isbn,omitempty"`
	DOI         *string  `json:"doi,omitempty"`
	URL         *string  `json:"url,omitempty"`
	ExternalID  *string  `json:"external_id,omitempty"`
	Tags        []string `json:"tags,omitempty"`
}

type UpdateRequest struct {
	Title       *string  `json:"title,omitempty"`
	Subtitle    *string  `json:"subtitle,omitempty"`
	Type        *string  `json:"type,omitempty"`
	Description *string  `json:"description,omitempty"`
	Publisher   *string  `json:"publisher,omitempty"`
	ISBN        *string  `json:"isbn,omitempty"`
	DOI         *string  `json:"doi,omitempty"`
	URL         *string  `json:"url,omitempty"`
	ExternalID  *string  `json:"external_id,omitempty"`
	Tags        []string `json:"tags,omitempty"`
}

type CreateBookRequest struct {
	Title        string             `json:"title" validate:"required"`
	Subtitle     *string            `json:"subtitle,omitempty"`
	Description  *string            `json:"description,omitempty"`
	URL          *string            `json:"url,omitempty"`
	ExternalID   *string            `json:"external_id,omitempty"`
	Tags         []string           `json:"tags,omitempty"`
	ISBN10       *string            `json:"isbn_10,omitempty"`
	ISBN13       *string            `json:"isbn_13,omitempty"`
	Publisher    *string            `json:"publisher,omitempty"`
	PageCount    *int               `json:"page_count,omitempty"`
	Language     *string            `json:"language,omitempty"`
	CoverURL     *string            `json:"cover_url,omitempty"`
	Contributors []ContributorInput `json:"contributors,omitempty"`
}

func (h *Handler) RegisterPublicRoutes(e *echo.Echo) {
	e.GET("/sources", h.List)
	e.GET("/sources/search", h.Search)
	e.GET("/sources/books/:id", h.GetBookByID)
	e.GET("/sources/:id", h.GetByID)
}

func (h *Handler) RegisterProtectedRoutes(g *echo.Group) {
	g.POST("/sources", h.Create)
	g.POST("/sources/books", h.CreateBook)
	g.PUT("/sources/:id", h.Update)
	g.DELETE("/sources/:id", h.Delete)
}

func (h *Handler) Create(c *echo.Context) error {
	var req CreateRequest
	if err := echox.BindAndValidate(c, &req); err != nil {
		return err
	}

	source, err := h.service.Create(c.Request().Context(), CreateSourceParams{
		Title:       req.Title,
		Subtitle:    req.Subtitle,
		Type:        SourceType(req.Type),
		Description: req.Description,
		Publisher:   req.Publisher,
		ISBN:        req.ISBN,
		DOI:         req.DOI,
		URL:         req.URL,
		ExternalID:  req.ExternalID,
		Tags:        req.Tags,
	})
	if errors.Is(err, ErrInvalidSource) {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid source")
	}
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to create source")
	}

	return c.JSON(http.StatusCreated, source)
}

func (h *Handler) CreateBook(c *echo.Context) error {
	var req CreateBookRequest
	if err := echox.BindAndValidate(c, &req); err != nil {
		return err
	}

	book, err := h.service.CreateBook(c.Request().Context(), CreateBookParams{
		Title:        req.Title,
		Subtitle:     req.Subtitle,
		Description:  req.Description,
		URL:          req.URL,
		ExternalID:   req.ExternalID,
		Tags:         req.Tags,
		ISBN10:       req.ISBN10,
		ISBN13:       req.ISBN13,
		Publisher:    req.Publisher,
		PageCount:    req.PageCount,
		Language:     req.Language,
		CoverURL:     req.CoverURL,
		Contributors: req.Contributors,
	})
	if errors.Is(err, ErrInvalidSource) {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid book")
	}
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to create book")
	}

	return c.JSON(http.StatusCreated, book)
}

func (h *Handler) GetByID(c *echo.Context) error {
	id, err := echox.ParamUUID(c, "id", "source ID")
	if err != nil {
		return err
	}

	source, err := h.service.GetByID(c.Request().Context(), id)
	if errors.Is(err, ErrSourceNotFound) {
		return echo.NewHTTPError(http.StatusNotFound, "source not found")
	}
	if errors.Is(err, ErrInvalidSource) {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid source")
	}
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to get source")
	}

	return c.JSON(http.StatusOK, source)
}

func (h *Handler) GetBookByID(c *echo.Context) error {
	id, err := echox.ParamUUID(c, "id", "source ID")
	if err != nil {
		return err
	}

	book, err := h.service.GetBookByID(c.Request().Context(), id)
	if errors.Is(err, ErrSourceNotFound) {
		return echo.NewHTTPError(http.StatusNotFound, "book not found")
	}
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to get book")
	}

	return c.JSON(http.StatusOK, book)
}

func (h *Handler) List(c *echo.Context) error {
	limit, offset := echox.Pagination(c)
	sourceType := c.QueryParam("type")

	var sources []*Source
	var err error
	if sourceType != "" {
		sources, err = h.service.ListByType(c.Request().Context(), SourceType(sourceType), limit, offset)
	} else {
		sources, err = h.service.List(c.Request().Context(), limit, offset)
	}
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to list sources")
	}

	return c.JSON(http.StatusOK, sources)
}

func (h *Handler) Search(c *echo.Context) error {
	query := c.QueryParam("q")
	if query == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "search query required")
	}

	limit, offset := echox.Pagination(c)
	sources, err := h.service.Search(c.Request().Context(), query, limit, offset)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to search sources")
	}

	return c.JSON(http.StatusOK, sources)
}

func (h *Handler) Update(c *echo.Context) error {
	id, err := echox.ParamUUID(c, "id", "source ID")
	if err != nil {
		return err
	}

	var req UpdateRequest
	if err := echox.BindAndValidate(c, &req); err != nil {
		return err
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
		Publisher:   req.Publisher,
		ISBN:        req.ISBN,
		DOI:         req.DOI,
		URL:         req.URL,
		ExternalID:  req.ExternalID,
		Tags:        req.Tags,
	})
	if errors.Is(err, ErrSourceNotFound) {
		return echo.NewHTTPError(http.StatusNotFound, "source not found")
	}
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to update source")
	}

	return c.JSON(http.StatusOK, source)
}

func (h *Handler) Delete(c *echo.Context) error {
	id, err := echox.ParamUUID(c, "id", "source ID")
	if err != nil {
		return err
	}
	if err := h.service.Delete(c.Request().Context(), id); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to delete source")
	}
	return c.NoContent(http.StatusNoContent)
}
