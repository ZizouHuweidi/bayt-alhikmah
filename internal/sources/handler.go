package sources

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/gofrs/uuid/v5"
	"github.com/zizouhuweidi/maktaba/internal/httpx"
)

type Handler struct {
	service *Service
	logger  *slog.Logger
}

func NewHandler(service *Service, logger *slog.Logger) *Handler {
	return &Handler{service: service, logger: logger}
}

type CreateRequest struct {
	Title       string   `json:"title"`
	Subtitle    *string  `json:"subtitle,omitempty"`
	Type        string   `json:"type"`
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
	Title        string             `json:"title"`
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

func (h *Handler) RegisterPublicRoutes(r httpx.Router) {
	r.Get("/sources", h.List)
	r.Get("/sources/search", h.Search)
	r.Get("/sources/books/:id", h.GetBookByID)
	r.Get("/sources/:id", h.GetByID)
}

func (h *Handler) RegisterProtectedRoutes(r httpx.Router) {
	r.Post("/sources", h.Create)
	r.Post("/sources/books", h.CreateBook)
	r.Put("/sources/:id", h.Update)
	r.Delete("/sources/:id", h.Delete)
}

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	var req CreateRequest
	if err := httpx.ReadJSON(r, &req); err != nil {
		httpx.WriteError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	source, err := h.service.Create(r.Context(), CreateSourceParams{
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
		httpx.WriteError(w, http.StatusBadRequest, "invalid source")
		return
	}
	if err != nil {
		httpx.WriteError(w, http.StatusInternalServerError, "failed to create source")
		return
	}

	httpx.WriteJSON(w, http.StatusCreated, source)
}

func (h *Handler) CreateBook(w http.ResponseWriter, r *http.Request) {
	var req CreateBookRequest
	if err := httpx.ReadJSON(r, &req); err != nil {
		httpx.WriteError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	book, err := h.service.CreateBook(r.Context(), CreateBookParams{
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
		httpx.WriteError(w, http.StatusBadRequest, "invalid book")
		return
	}
	if err != nil {
		httpx.WriteError(w, http.StatusInternalServerError, "failed to create book")
		return
	}

	httpx.WriteJSON(w, http.StatusCreated, book)
}

func (h *Handler) GetByID(w http.ResponseWriter, r *http.Request) {
	id, ok := parsePathUUID(w, r, "id", "source ID")
	if !ok {
		return
	}

	source, err := h.service.GetByID(r.Context(), id)
	if errors.Is(err, ErrSourceNotFound) {
		httpx.WriteError(w, http.StatusNotFound, "source not found")
		return
	}
	if errors.Is(err, ErrInvalidSource) {
		httpx.WriteError(w, http.StatusBadRequest, "invalid source")
		return
	}
	if err != nil {
		httpx.WriteError(w, http.StatusInternalServerError, "failed to get source")
		return
	}

	httpx.WriteJSON(w, http.StatusOK, source)
}

func (h *Handler) GetBookByID(w http.ResponseWriter, r *http.Request) {
	id, ok := parsePathUUID(w, r, "id", "source ID")
	if !ok {
		return
	}

	book, err := h.service.GetBookByID(r.Context(), id)
	if errors.Is(err, ErrSourceNotFound) {
		httpx.WriteError(w, http.StatusNotFound, "book not found")
		return
	}
	if err != nil {
		httpx.WriteError(w, http.StatusInternalServerError, "failed to get book")
		return
	}

	httpx.WriteJSON(w, http.StatusOK, book)
}

func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	limit, offset := httpx.Pagination(r)
	sourceType := r.URL.Query().Get("type")

	var sources []*Source
	var err error
	if sourceType != "" {
		sources, err = h.service.ListByType(r.Context(), SourceType(sourceType), limit, offset)
	} else {
		sources, err = h.service.List(r.Context(), limit, offset)
	}
	if err != nil {
		httpx.WriteError(w, http.StatusInternalServerError, "failed to list sources")
		return
	}

	httpx.WriteJSON(w, http.StatusOK, sources)
}

func (h *Handler) Search(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("q")
	if query == "" {
		httpx.WriteError(w, http.StatusBadRequest, "search query required")
		return
	}

	limit, offset := httpx.Pagination(r)
	sources, err := h.service.Search(r.Context(), query, limit, offset)
	if err != nil {
		httpx.WriteError(w, http.StatusInternalServerError, "failed to search sources")
		return
	}

	httpx.WriteJSON(w, http.StatusOK, sources)
}

func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	id, ok := parsePathUUID(w, r, "id", "source ID")
	if !ok {
		return
	}

	var req UpdateRequest
	if err := httpx.ReadJSON(r, &req); err != nil {
		httpx.WriteError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	var sourceType *SourceType
	if req.Type != nil {
		st := SourceType(*req.Type)
		sourceType = &st
	}

	source, err := h.service.Update(r.Context(), id, UpdateSourceParams{
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
		httpx.WriteError(w, http.StatusNotFound, "source not found")
		return
	}
	if err != nil {
		httpx.WriteError(w, http.StatusInternalServerError, "failed to update source")
		return
	}

	httpx.WriteJSON(w, http.StatusOK, source)
}

func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	id, ok := parsePathUUID(w, r, "id", "source ID")
	if !ok {
		return
	}
	if err := h.service.Delete(r.Context(), id); err != nil {
		httpx.WriteError(w, http.StatusInternalServerError, "failed to delete source")
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func parsePathUUID(w http.ResponseWriter, r *http.Request, key, label string) (uuid.UUID, bool) {
	id, err := uuid.FromString(r.PathValue(key))
	if err != nil {
		httpx.WriteError(w, http.StatusBadRequest, "invalid "+label)
		return uuid.Nil, false
	}
	return id, true
}

func parseOptionalUUID(w http.ResponseWriter, value *string, label string) (*uuid.UUID, bool) {
	if value == nil {
		return nil, true
	}
	id, err := uuid.FromString(*value)
	if err != nil {
		httpx.WriteError(w, http.StatusBadRequest, "invalid "+label)
		return nil, false
	}
	return &id, true
}
