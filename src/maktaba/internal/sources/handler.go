package sources

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

// Handler handles HTTP requests for sources
type Handler struct {
	service *Service
	logger  *slog.Logger
}

// NewHandler creates a new sources handler
func NewHandler(service *Service, logger *slog.Logger) *Handler {
	return &Handler{
		service: service,
		logger:  logger,
	}
}

// RegisterRoutes registers the sources routes
func (h *Handler) RegisterRoutes(r chi.Router) {
	r.Route("/sources", func(r chi.Router) {
		r.Post("/", h.Create)
		r.Get("/", h.List)
		r.Get("/search", h.Search)
		r.Get("/{id}", h.GetByID)
		r.Put("/{id}", h.Update)
		r.Delete("/{id}", h.Delete)
	})
}

// CreateRequest represents the request body for creating a source
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

// UpdateRequest represents the request body for updating a source
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

// Create handles POST /sources
func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	var req CreateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.respondError(w, "invalid request body", http.StatusBadRequest)
		return
	}

	var authorID *uuid.UUID
	if req.AuthorID != nil {
		id, err := uuid.Parse(*req.AuthorID)
		if err != nil {
			h.respondError(w, "invalid author_id", http.StatusBadRequest)
			return
		}
		authorID = &id
	}

	source, err := h.service.Create(r.Context(), CreateSourceParams{
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
		h.respondError(w, "failed to create source", http.StatusInternalServerError)
		return
	}

	h.respondJSON(w, source, http.StatusCreated)
}

// GetByID handles GET /sources/{id}
func (h *Handler) GetByID(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		h.respondError(w, "invalid source ID", http.StatusBadRequest)
		return
	}

	source, err := h.service.GetByID(r.Context(), id)
	if err == ErrSourceNotFound {
		h.respondError(w, "source not found", http.StatusNotFound)
		return
	}
	if err != nil {
		h.respondError(w, "failed to get source", http.StatusInternalServerError)
		return
	}

	h.respondJSON(w, source, http.StatusOK)
}

// List handles GET /sources
func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	limit, offset := h.parsePagination(r)
	sourceType := r.URL.Query().Get("type")

	var sources []*Source
	var err error

	if sourceType != "" {
		sources, err = h.service.ListByType(r.Context(), SourceType(sourceType), limit, offset)
	} else {
		sources, err = h.service.List(r.Context(), limit, offset)
	}

	if err != nil {
		h.respondError(w, "failed to list sources", http.StatusInternalServerError)
		return
	}

	h.respondJSON(w, sources, http.StatusOK)
}

// Search handles GET /sources/search
func (h *Handler) Search(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("q")
	if query == "" {
		h.respondError(w, "search query required", http.StatusBadRequest)
		return
	}

	limit, offset := h.parsePagination(r)
	sources, err := h.service.Search(r.Context(), query, limit, offset)
	if err != nil {
		h.respondError(w, "failed to search sources", http.StatusInternalServerError)
		return
	}

	h.respondJSON(w, sources, http.StatusOK)
}

// Update handles PUT /sources/{id}
func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		h.respondError(w, "invalid source ID", http.StatusBadRequest)
		return
	}

	var req UpdateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.respondError(w, "invalid request body", http.StatusBadRequest)
		return
	}

	var authorID *uuid.UUID
	if req.AuthorID != nil {
		aid, err := uuid.Parse(*req.AuthorID)
		if err != nil {
			h.respondError(w, "invalid author_id", http.StatusBadRequest)
			return
		}
		authorID = &aid
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
		AuthorID:    authorID,
		Publisher:   req.Publisher,
		ISBN:        req.ISBN,
		DOI:         req.DOI,
		URL:         req.URL,
		ExternalID:  req.ExternalID,
		Tags:        req.Tags,
	})

	if err == ErrSourceNotFound {
		h.respondError(w, "source not found", http.StatusNotFound)
		return
	}
	if err != nil {
		h.respondError(w, "failed to update source", http.StatusInternalServerError)
		return
	}

	h.respondJSON(w, source, http.StatusOK)
}

// Delete handles DELETE /sources/{id}
func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		h.respondError(w, "invalid source ID", http.StatusBadRequest)
		return
	}

	if err := h.service.Delete(r.Context(), id); err != nil {
		h.respondError(w, "failed to delete source", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) parsePagination(r *http.Request) (limit, offset int) {
	limit = 100
	offset = 0

	if l := r.URL.Query().Get("limit"); l != "" {
		if parsed, err := strconv.Atoi(l); err == nil && parsed > 0 {
			limit = parsed
		}
	}

	if o := r.URL.Query().Get("offset"); o != "" {
		if parsed, err := strconv.Atoi(o); err == nil && parsed >= 0 {
			offset = parsed
		}
	}

	return limit, offset
}

func (h *Handler) respondJSON(w http.ResponseWriter, data interface{}, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func (h *Handler) respondError(w http.ResponseWriter, message string, status int) {
	h.respondJSON(w, map[string]string{"error": message}, status)
}
