package collections

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
	Name        string   `json:"name"`
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

func (h *Handler) RegisterPublicRoutes(r httpx.Router) {
	r.Get("/collections", h.List)
	r.Get("/collections/:id", h.GetByID)
}

func (h *Handler) RegisterProtectedRoutes(r httpx.Router) {
	r.Post("/collections", h.Create)
	r.Get("/collections", h.ListOwn)
	r.Put("/collections/:id", h.Update)
	r.Delete("/collections/:id", h.Delete)
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
	sourceIDs, ok := parseUUIDs(w, req.SourceIDs, "source_ids")
	if !ok {
		return
	}

	collection, err := h.service.Create(r.Context(), CreateCollectionParams{
		UserID:      userID,
		Name:        req.Name,
		Description: req.Description,
		IsPublic:    req.IsPublic,
		SourceIDs:   sourceIDs,
	})
	if errors.Is(err, ErrInvalidCollection) {
		httpx.WriteError(w, http.StatusBadRequest, "invalid collection")
		return
	}
	if errors.Is(err, ErrSourceNotFound) {
		httpx.WriteError(w, http.StatusNotFound, "source not found")
		return
	}
	if err != nil {
		httpx.WriteError(w, http.StatusInternalServerError, "failed to create collection")
		return
	}

	httpx.WriteJSON(w, http.StatusCreated, collection)
}

func (h *Handler) GetByID(w http.ResponseWriter, r *http.Request) {
	id, ok := parsePathUUID(w, r, "id", "collection ID")
	if !ok {
		return
	}

	collection, err := h.service.GetByID(r.Context(), id)
	if errors.Is(err, ErrCollectionNotFound) {
		httpx.WriteError(w, http.StatusNotFound, "collection not found")
		return
	}
	if err != nil {
		httpx.WriteError(w, http.StatusInternalServerError, "failed to get collection")
		return
	}
	if !collection.IsPublic {
		httpx.WriteError(w, http.StatusForbidden, "collection is private")
		return
	}

	httpx.WriteJSON(w, http.StatusOK, collection)
}

func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	limit, offset := httpx.Pagination(r)
	userIDStr := r.URL.Query().Get("user_id")
	if userIDStr == "" {
		httpx.WriteError(w, http.StatusBadRequest, "user_id required")
		return
	}
	userID, err := uuid.FromString(userIDStr)
	if err != nil {
		httpx.WriteError(w, http.StatusBadRequest, "invalid user_id")
		return
	}

	collections, err := h.service.ListPublicByUser(r.Context(), userID, limit, offset)
	if err != nil {
		httpx.WriteError(w, http.StatusInternalServerError, "failed to list collections")
		return
	}

	httpx.WriteJSON(w, http.StatusOK, collections)
}

func (h *Handler) ListOwn(w http.ResponseWriter, r *http.Request) {
	userID, ok := auth.UserIDFromContext(r.Context())
	if !ok {
		httpx.WriteError(w, http.StatusUnauthorized, "authentication required")
		return
	}
	limit, offset := httpx.Pagination(r)

	collections, err := h.service.ListByUser(r.Context(), userID, limit, offset)
	if err != nil {
		httpx.WriteError(w, http.StatusInternalServerError, "failed to list collections")
		return
	}

	httpx.WriteJSON(w, http.StatusOK, collections)
}

func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	userID, ok := auth.UserIDFromContext(r.Context())
	if !ok {
		httpx.WriteError(w, http.StatusUnauthorized, "authentication required")
		return
	}
	id, ok := parsePathUUID(w, r, "id", "collection ID")
	if !ok {
		return
	}

	existing, err := h.service.GetByID(r.Context(), id)
	if errors.Is(err, ErrCollectionNotFound) {
		httpx.WriteError(w, http.StatusNotFound, "collection not found")
		return
	}
	if err != nil {
		httpx.WriteError(w, http.StatusInternalServerError, "failed to get collection")
		return
	}
	if existing.UserID != userID {
		httpx.WriteError(w, http.StatusForbidden, "cannot update another user's collection")
		return
	}

	var req UpdateRequest
	if err := httpx.ReadJSON(r, &req); err != nil {
		httpx.WriteError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	sourceIDs, ok := parseUUIDs(w, req.SourceIDs, "source_ids")
	if !ok {
		return
	}

	collection, err := h.service.Update(r.Context(), id, UpdateCollectionParams{
		Name:        req.Name,
		Description: req.Description,
		IsPublic:    req.IsPublic,
		SourceIDs:   sourceIDs,
	})
	if errors.Is(err, ErrInvalidCollection) {
		httpx.WriteError(w, http.StatusBadRequest, "invalid collection")
		return
	}
	if errors.Is(err, ErrSourceNotFound) {
		httpx.WriteError(w, http.StatusNotFound, "source not found")
		return
	}
	if err != nil {
		httpx.WriteError(w, http.StatusInternalServerError, "failed to update collection")
		return
	}

	httpx.WriteJSON(w, http.StatusOK, collection)
}

func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	userID, ok := auth.UserIDFromContext(r.Context())
	if !ok {
		httpx.WriteError(w, http.StatusUnauthorized, "authentication required")
		return
	}
	id, ok := parsePathUUID(w, r, "id", "collection ID")
	if !ok {
		return
	}

	existing, err := h.service.GetByID(r.Context(), id)
	if errors.Is(err, ErrCollectionNotFound) {
		httpx.WriteError(w, http.StatusNotFound, "collection not found")
		return
	}
	if err != nil {
		httpx.WriteError(w, http.StatusInternalServerError, "failed to get collection")
		return
	}
	if existing.UserID != userID {
		httpx.WriteError(w, http.StatusForbidden, "cannot delete another user's collection")
		return
	}

	if err := h.service.Delete(r.Context(), id); err != nil {
		httpx.WriteError(w, http.StatusInternalServerError, "failed to delete collection")
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

func parseUUIDs(w http.ResponseWriter, values []string, label string) ([]uuid.UUID, bool) {
	if values == nil {
		return nil, true
	}
	ids := make([]uuid.UUID, 0, len(values))
	for _, value := range values {
		id, err := uuid.FromString(value)
		if err != nil {
			httpx.WriteError(w, http.StatusBadRequest, "invalid "+label)
			return nil, false
		}
		ids = append(ids, id)
	}
	return ids, true
}
