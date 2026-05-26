package library

import (
	"errors"
	"log/slog"
	"net/http"
	"time"

	"github.com/gofrs/uuid/v5"
	"github.com/zizouhuweidi/maktaba/internal/auth"
	"github.com/zizouhuweidi/maktaba/internal/httpx"
)

type Handler struct {
	service *Service
	logger  *slog.Logger
}

type createRequest struct {
	SourceID      string        `json:"source_id"`
	Status        string        `json:"status"`
	ProgressValue *int          `json:"progress_value,omitempty"`
	ProgressUnit  *ProgressUnit `json:"progress_unit,omitempty"`
	Visibility    string        `json:"visibility,omitempty"`
	StartedAt     *time.Time    `json:"started_at,omitempty"`
	CompletedAt   *time.Time    `json:"completed_at,omitempty"`
}

type updateRequest struct {
	Status        *string       `json:"status,omitempty"`
	ProgressValue *int          `json:"progress_value,omitempty"`
	ProgressUnit  *ProgressUnit `json:"progress_unit,omitempty"`
	Visibility    *string       `json:"visibility,omitempty"`
	StartedAt     *time.Time    `json:"started_at,omitempty"`
	CompletedAt   *time.Time    `json:"completed_at,omitempty"`
}

func NewHandler(service *Service, logger *slog.Logger) *Handler {
	return &Handler{service: service, logger: logger}
}

func (h *Handler) RegisterPublicRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /users/{user_id}/library", h.ListPublicByUser)
}

func (h *Handler) RegisterProtectedRoutes(mux *http.ServeMux, middleware func(http.Handler) http.Handler) {
	mux.Handle("POST /api/library/items", middleware(http.HandlerFunc(h.Create)))
	mux.Handle("GET /api/library/items", middleware(http.HandlerFunc(h.ListMine)))
	mux.Handle("GET /api/library/items/{id}", middleware(http.HandlerFunc(h.GetMine)))
	mux.Handle("PUT /api/library/items/{id}", middleware(http.HandlerFunc(h.Update)))
	mux.Handle("DELETE /api/library/items/{id}", middleware(http.HandlerFunc(h.Delete)))
}

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	userID, ok := auth.UserIDFromContext(r.Context())
	if !ok {
		httpx.WriteError(w, http.StatusUnauthorized, "authentication required")
		return
	}

	var req createRequest
	if err := httpx.ReadJSON(r, &req); err != nil {
		httpx.WriteError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	sourceID, err := uuid.FromString(req.SourceID)
	if err != nil {
		httpx.WriteError(w, http.StatusBadRequest, "invalid source_id")
		return
	}
	visibility := Visibility(req.Visibility)
	if visibility == "" {
		visibility = VisibilityPrivate
	}

	item, err := h.service.Create(r.Context(), CreateItemParams{
		UserID:        userID,
		SourceID:      sourceID,
		Status:        Status(req.Status),
		ProgressValue: req.ProgressValue,
		ProgressUnit:  req.ProgressUnit,
		Visibility:    visibility,
		StartedAt:     req.StartedAt,
		CompletedAt:   req.CompletedAt,
	})
	if errors.Is(err, ErrInvalidItem) {
		httpx.WriteError(w, http.StatusBadRequest, "invalid library item")
		return
	}
	if err != nil {
		h.logger.Error("failed to create library item", "error", err)
		httpx.WriteError(w, http.StatusInternalServerError, "failed to create library item")
		return
	}

	httpx.WriteJSON(w, http.StatusCreated, item)
}

func (h *Handler) ListMine(w http.ResponseWriter, r *http.Request) {
	userID, ok := auth.UserIDFromContext(r.Context())
	if !ok {
		httpx.WriteError(w, http.StatusUnauthorized, "authentication required")
		return
	}
	limit, offset := httpx.Pagination(r)
	items, err := h.service.ListByUser(r.Context(), userID, limit, offset)
	if err != nil {
		httpx.WriteError(w, http.StatusInternalServerError, "failed to list library items")
		return
	}
	httpx.WriteJSON(w, http.StatusOK, items)
}

func (h *Handler) GetMine(w http.ResponseWriter, r *http.Request) {
	userID, ok := auth.UserIDFromContext(r.Context())
	if !ok {
		httpx.WriteError(w, http.StatusUnauthorized, "authentication required")
		return
	}
	item, ok := h.getOwnedItem(w, r, userID)
	if !ok {
		return
	}
	httpx.WriteJSON(w, http.StatusOK, item)
}

func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	userID, ok := auth.UserIDFromContext(r.Context())
	if !ok {
		httpx.WriteError(w, http.StatusUnauthorized, "authentication required")
		return
	}
	item, ok := h.getOwnedItem(w, r, userID)
	if !ok {
		return
	}

	var req updateRequest
	if err := httpx.ReadJSON(r, &req); err != nil {
		httpx.WriteError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	var status *Status
	if req.Status != nil {
		parsed := Status(*req.Status)
		status = &parsed
	}
	var visibility *Visibility
	if req.Visibility != nil {
		parsed := Visibility(*req.Visibility)
		visibility = &parsed
	}

	updated, err := h.service.Update(r.Context(), item.ID, UpdateItemParams{
		Status:        status,
		ProgressValue: req.ProgressValue,
		ProgressUnit:  req.ProgressUnit,
		Visibility:    visibility,
		StartedAt:     req.StartedAt,
		CompletedAt:   req.CompletedAt,
	})
	if errors.Is(err, ErrInvalidItem) {
		httpx.WriteError(w, http.StatusBadRequest, "invalid library item")
		return
	}
	if err != nil {
		httpx.WriteError(w, http.StatusInternalServerError, "failed to update library item")
		return
	}
	httpx.WriteJSON(w, http.StatusOK, updated)
}

func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	userID, ok := auth.UserIDFromContext(r.Context())
	if !ok {
		httpx.WriteError(w, http.StatusUnauthorized, "authentication required")
		return
	}
	item, ok := h.getOwnedItem(w, r, userID)
	if !ok {
		return
	}
	if err := h.service.Delete(r.Context(), item.ID); err != nil {
		httpx.WriteError(w, http.StatusInternalServerError, "failed to delete library item")
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) ListPublicByUser(w http.ResponseWriter, r *http.Request) {
	userID, err := uuid.FromString(r.PathValue("user_id"))
	if err != nil {
		httpx.WriteError(w, http.StatusBadRequest, "invalid user ID")
		return
	}
	limit, offset := httpx.Pagination(r)
	items, err := h.service.ListPublicByUser(r.Context(), userID, limit, offset)
	if err != nil {
		httpx.WriteError(w, http.StatusInternalServerError, "failed to list public library items")
		return
	}
	httpx.WriteJSON(w, http.StatusOK, items)
}

func (h *Handler) getOwnedItem(w http.ResponseWriter, r *http.Request, userID uuid.UUID) (*Item, bool) {
	id, err := uuid.FromString(r.PathValue("id"))
	if err != nil {
		httpx.WriteError(w, http.StatusBadRequest, "invalid library item ID")
		return nil, false
	}
	item, err := h.service.GetByID(r.Context(), id)
	if errors.Is(err, ErrItemNotFound) {
		httpx.WriteError(w, http.StatusNotFound, "library item not found")
		return nil, false
	}
	if err != nil {
		httpx.WriteError(w, http.StatusInternalServerError, "failed to get library item")
		return nil, false
	}
	if item.UserID != userID {
		httpx.WriteError(w, http.StatusForbidden, "cannot access another user's library item")
		return nil, false
	}
	return item, true
}
