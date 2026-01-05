package notes

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

// Handler handles HTTP requests for notes
type Handler struct {
	service *Service
	logger  *slog.Logger
}

// NewHandler creates a new notes handler
func NewHandler(service *Service, logger *slog.Logger) *Handler {
	return &Handler{
		service: service,
		logger:  logger,
	}
}

// RegisterRoutes registers the notes routes
func (h *Handler) RegisterRoutes(r chi.Router) {
	r.Route("/notes", func(r chi.Router) {
		r.Post("/", h.Create)
		r.Get("/", h.List)
		r.Get("/{id}", h.GetByID)
		r.Put("/{id}", h.Update)
		r.Delete("/{id}", h.Delete)
	})
}

// CreateRequest represents the request body for creating a note
type CreateRequest struct {
	UserID      string   `json:"user_id"`
	SourceID    *string  `json:"source_id,omitempty"`
	Content     string   `json:"content"`
	ContentType string   `json:"content_type"`
	IsPublic    bool     `json:"is_public"`
	Annotations []string `json:"annotations,omitempty"`
	Tags        []string `json:"tags,omitempty"`
}

// UpdateRequest represents the request body for updating a note
type UpdateRequest struct {
	Content     *string  `json:"content,omitempty"`
	ContentType *string  `json:"content_type,omitempty"`
	IsPublic    *bool    `json:"is_public,omitempty"`
	Annotations []string `json:"annotations,omitempty"`
	Tags        []string `json:"tags,omitempty"`
}

// Create handles POST /notes
func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	var req CreateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.respondError(w, "invalid request body", http.StatusBadRequest)
		return
	}

	userID, err := uuid.Parse(req.UserID)
	if err != nil {
		h.respondError(w, "invalid user_id", http.StatusBadRequest)
		return
	}

	var sourceID *uuid.UUID
	if req.SourceID != nil {
		id, err := uuid.Parse(*req.SourceID)
		if err != nil {
			h.respondError(w, "invalid source_id", http.StatusBadRequest)
			return
		}
		sourceID = &id
	}

	note, err := h.service.Create(r.Context(), CreateNoteParams{
		UserID:      userID,
		SourceID:    sourceID,
		Content:     req.Content,
		ContentType: ContentType(req.ContentType),
		IsPublic:    req.IsPublic,
		Annotations: req.Annotations,
		Tags:        req.Tags,
	})

	if err != nil {
		h.respondError(w, "failed to create note", http.StatusInternalServerError)
		return
	}

	h.respondJSON(w, note, http.StatusCreated)
}

// GetByID handles GET /notes/{id}
func (h *Handler) GetByID(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		h.respondError(w, "invalid note ID", http.StatusBadRequest)
		return
	}

	note, err := h.service.GetByID(r.Context(), id)
	if err == ErrNoteNotFound {
		h.respondError(w, "note not found", http.StatusNotFound)
		return
	}
	if err != nil {
		h.respondError(w, "failed to get note", http.StatusInternalServerError)
		return
	}

	h.respondJSON(w, note, http.StatusOK)
}

// List handles GET /notes
func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	limit, offset := h.parsePagination(r)
	userIDStr := r.URL.Query().Get("user_id")
	sourceIDStr := r.URL.Query().Get("source_id")
	publicOnly := r.URL.Query().Get("public") == "true"

	var notes []*Note
	var err error

	if userIDStr != "" {
		userID, err := uuid.Parse(userIDStr)
		if err != nil {
			h.respondError(w, "invalid user_id", http.StatusBadRequest)
			return
		}
		notes, err = h.service.ListByUser(r.Context(), userID, limit, offset)
	} else if sourceIDStr != "" {
		sourceID, err := uuid.Parse(sourceIDStr)
		if err != nil {
			h.respondError(w, "invalid source_id", http.StatusBadRequest)
			return
		}
		notes, err = h.service.ListBySource(r.Context(), sourceID, limit, offset)
	} else if publicOnly {
		notes, err = h.service.ListPublic(r.Context(), limit, offset)
	} else {
		h.respondError(w, "user_id, source_id, or public=true required", http.StatusBadRequest)
		return
	}

	if err != nil {
		h.respondError(w, "failed to list notes", http.StatusInternalServerError)
		return
	}

	h.respondJSON(w, notes, http.StatusOK)
}

// Update handles PUT /notes/{id}
func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		h.respondError(w, "invalid note ID", http.StatusBadRequest)
		return
	}

	var req UpdateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.respondError(w, "invalid request body", http.StatusBadRequest)
		return
	}

	var contentType *ContentType
	if req.ContentType != nil {
		ct := ContentType(*req.ContentType)
		contentType = &ct
	}

	note, err := h.service.Update(r.Context(), id, UpdateNoteParams{
		Content:     req.Content,
		ContentType: contentType,
		IsPublic:    req.IsPublic,
		Annotations: req.Annotations,
		Tags:        req.Tags,
	})

	if err == ErrNoteNotFound {
		h.respondError(w, "note not found", http.StatusNotFound)
		return
	}
	if err != nil {
		h.respondError(w, "failed to update note", http.StatusInternalServerError)
		return
	}

	h.respondJSON(w, note, http.StatusOK)
}

// Delete handles DELETE /notes/{id}
func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		h.respondError(w, "invalid note ID", http.StatusBadRequest)
		return
	}

	if err := h.service.Delete(r.Context(), id); err != nil {
		h.respondError(w, "failed to delete note", http.StatusInternalServerError)
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
