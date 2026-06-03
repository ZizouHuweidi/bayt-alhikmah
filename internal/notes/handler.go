package notes

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/gofrs/uuid/v5"
	"github.com/zizouhuweidi/maktaba/internal/auth"
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

func (h *Handler) RegisterPublicRoutes(r httpx.Router) {
	r.Get("/notes", h.List)
	r.Get("/notes/:id", h.GetByID)
}

func (h *Handler) RegisterProtectedRoutes(r httpx.Router) {
	r.Post("/notes", h.Create)
	r.Get("/notes", h.ListMine)
	r.Put("/notes/:id", h.Update)
	r.Delete("/notes/:id", h.Delete)
}

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	userID, ok := auth.UserIDFromContext(r.Context())
	if !ok {
		httpx.WriteError(w, http.StatusUnauthorized, "authentication required")
		return
	}

	var req CreateRequest
	if err := httpx.ReadJSON(r, &req); err != nil {
		httpx.WriteError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	sourceID, ok := parseOptionalUUID(w, req.SourceID, "source_id")
	if !ok {
		return
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
	if errors.Is(err, ErrInvalidNote) {
		httpx.WriteError(w, http.StatusBadRequest, "invalid note")
		return
	}
	if errors.Is(err, ErrSourceNotFound) {
		httpx.WriteError(w, http.StatusNotFound, "source not found")
		return
	}
	if err != nil {
		httpx.WriteError(w, http.StatusInternalServerError, "failed to create note")
		return
	}

	httpx.WriteJSON(w, http.StatusCreated, note)
}

func (h *Handler) GetByID(w http.ResponseWriter, r *http.Request) {
	id, ok := parsePathUUID(w, r, "id", "note ID")
	if !ok {
		return
	}

	note, err := h.service.GetByID(r.Context(), id)
	if errors.Is(err, ErrNoteNotFound) {
		httpx.WriteError(w, http.StatusNotFound, "note not found")
		return
	}
	if err != nil {
		httpx.WriteError(w, http.StatusInternalServerError, "failed to get note")
		return
	}
	if !note.IsPublic {
		httpx.WriteError(w, http.StatusForbidden, "note is private")
		return
	}

	httpx.WriteJSON(w, http.StatusOK, note)
}

func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	limit, offset := httpx.Pagination(r)
	userIDStr := r.URL.Query().Get("user_id")
	sourceIDStr := r.URL.Query().Get("source_id")
	publicOnly := r.URL.Query().Get("public") == "true"

	var result []*Note
	var err error
	if userIDStr != "" {
		userID, parseErr := uuid.FromString(userIDStr)
		if parseErr != nil {
			httpx.WriteError(w, http.StatusBadRequest, "invalid user_id")
			return
		}
		result, err = h.service.ListPublicByUser(r.Context(), userID, limit, offset)
	} else if sourceIDStr != "" {
		sourceID, parseErr := uuid.FromString(sourceIDStr)
		if parseErr != nil {
			httpx.WriteError(w, http.StatusBadRequest, "invalid source_id")
			return
		}
		result, err = h.service.ListPublicBySource(r.Context(), sourceID, limit, offset)
	} else if publicOnly {
		result, err = h.service.ListPublic(r.Context(), limit, offset)
	} else {
		httpx.WriteError(w, http.StatusBadRequest, "user_id, source_id, or public=true required")
		return
	}
	if err != nil {
		httpx.WriteError(w, http.StatusInternalServerError, "failed to list notes")
		return
	}

	httpx.WriteJSON(w, http.StatusOK, result)
}

func (h *Handler) ListMine(w http.ResponseWriter, r *http.Request) {
	userID, ok := auth.UserIDFromContext(r.Context())
	if !ok {
		httpx.WriteError(w, http.StatusUnauthorized, "authentication required")
		return
	}
	limit, offset := httpx.Pagination(r)

	notes, err := h.service.ListByUser(r.Context(), userID, limit, offset)
	if err != nil {
		httpx.WriteError(w, http.StatusInternalServerError, "failed to list notes")
		return
	}

	httpx.WriteJSON(w, http.StatusOK, notes)
}

func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	userID, ok := auth.UserIDFromContext(r.Context())
	if !ok {
		httpx.WriteError(w, http.StatusUnauthorized, "authentication required")
		return
	}
	id, ok := parsePathUUID(w, r, "id", "note ID")
	if !ok {
		return
	}

	existing, err := h.service.GetByID(r.Context(), id)
	if errors.Is(err, ErrNoteNotFound) {
		httpx.WriteError(w, http.StatusNotFound, "note not found")
		return
	}
	if err != nil {
		httpx.WriteError(w, http.StatusInternalServerError, "failed to get note")
		return
	}
	if existing.UserID != userID {
		httpx.WriteError(w, http.StatusForbidden, "cannot update another user's note")
		return
	}

	var req UpdateRequest
	if err := httpx.ReadJSON(r, &req); err != nil {
		httpx.WriteError(w, http.StatusBadRequest, "invalid request body")
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
	if err != nil {
		httpx.WriteError(w, http.StatusInternalServerError, "failed to update note")
		return
	}

	httpx.WriteJSON(w, http.StatusOK, note)
}

func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	userID, ok := auth.UserIDFromContext(r.Context())
	if !ok {
		httpx.WriteError(w, http.StatusUnauthorized, "authentication required")
		return
	}
	id, ok := parsePathUUID(w, r, "id", "note ID")
	if !ok {
		return
	}

	existing, err := h.service.GetByID(r.Context(), id)
	if errors.Is(err, ErrNoteNotFound) {
		httpx.WriteError(w, http.StatusNotFound, "note not found")
		return
	}
	if err != nil {
		httpx.WriteError(w, http.StatusInternalServerError, "failed to get note")
		return
	}
	if existing.UserID != userID {
		httpx.WriteError(w, http.StatusForbidden, "cannot delete another user's note")
		return
	}

	if err := h.service.Delete(r.Context(), id); err != nil {
		httpx.WriteError(w, http.StatusInternalServerError, "failed to delete note")
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
