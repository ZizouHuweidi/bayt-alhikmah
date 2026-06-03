package reviews

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
	SourceID string  `json:"source_id"`
	Rating   int     `json:"rating"`
	Content  *string `json:"content,omitempty"`
	IsPublic bool    `json:"is_public"`
}

type UpdateRequest struct {
	Rating   *int    `json:"rating,omitempty"`
	Content  *string `json:"content,omitempty"`
	IsPublic *bool   `json:"is_public,omitempty"`
}

func (h *Handler) RegisterPublicRoutes(r httpx.Router) {
	r.Get("/reviews", h.List)
	r.Get("/reviews/:id", h.GetByID)
}

func (h *Handler) RegisterProtectedRoutes(r httpx.Router) {
	r.Post("/reviews", h.Create)
	r.Get("/reviews", h.ListOwn)
	r.Put("/reviews/:id", h.Update)
	r.Delete("/reviews/:id", h.Delete)
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
	sourceID, err := uuid.FromString(req.SourceID)
	if err != nil {
		httpx.WriteError(w, http.StatusBadRequest, "invalid source_id")
		return
	}

	review, err := h.service.Create(r.Context(), CreateReviewParams{
		UserID:   userID,
		SourceID: sourceID,
		Rating:   req.Rating,
		Content:  req.Content,
		IsPublic: req.IsPublic,
	})
	if errors.Is(err, ErrInvalidReview) {
		httpx.WriteError(w, http.StatusBadRequest, "invalid review")
		return
	}
	if errors.Is(err, ErrReviewExists) {
		httpx.WriteError(w, http.StatusConflict, "review already exists for source")
		return
	}
	if errors.Is(err, ErrSourceNotFound) {
		httpx.WriteError(w, http.StatusNotFound, "source not found")
		return
	}
	if errors.Is(err, ErrReviewConflict) {
		httpx.WriteError(w, http.StatusConflict, "review conflict")
		return
	}
	if err != nil {
		httpx.WriteError(w, http.StatusInternalServerError, "failed to create review")
		return
	}

	httpx.WriteJSON(w, http.StatusCreated, review)
}

func (h *Handler) GetByID(w http.ResponseWriter, r *http.Request) {
	id, ok := parsePathUUID(w, r, "id", "review ID")
	if !ok {
		return
	}

	review, err := h.service.GetByID(r.Context(), id)
	if errors.Is(err, ErrReviewNotFound) {
		httpx.WriteError(w, http.StatusNotFound, "review not found")
		return
	}
	if err != nil {
		httpx.WriteError(w, http.StatusInternalServerError, "failed to get review")
		return
	}
	if !review.IsPublic {
		httpx.WriteError(w, http.StatusForbidden, "review is private")
		return
	}

	httpx.WriteJSON(w, http.StatusOK, review)
}

func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	limit, offset := httpx.Pagination(r)
	userIDStr := r.URL.Query().Get("user_id")
	sourceIDStr := r.URL.Query().Get("source_id")

	var result []*Review
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
	} else {
		httpx.WriteError(w, http.StatusBadRequest, "user_id or source_id required")
		return
	}
	if err != nil {
		httpx.WriteError(w, http.StatusInternalServerError, "failed to list reviews")
		return
	}

	httpx.WriteJSON(w, http.StatusOK, result)
}

func (h *Handler) ListOwn(w http.ResponseWriter, r *http.Request) {
	userID, ok := auth.UserIDFromContext(r.Context())
	if !ok {
		httpx.WriteError(w, http.StatusUnauthorized, "authentication required")
		return
	}
	limit, offset := httpx.Pagination(r)

	reviews, err := h.service.ListByUser(r.Context(), userID, limit, offset)
	if err != nil {
		httpx.WriteError(w, http.StatusInternalServerError, "failed to list reviews")
		return
	}

	httpx.WriteJSON(w, http.StatusOK, reviews)
}

func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	userID, ok := auth.UserIDFromContext(r.Context())
	if !ok {
		httpx.WriteError(w, http.StatusUnauthorized, "authentication required")
		return
	}
	id, ok := parsePathUUID(w, r, "id", "review ID")
	if !ok {
		return
	}

	existing, err := h.service.GetByID(r.Context(), id)
	if errors.Is(err, ErrReviewNotFound) {
		httpx.WriteError(w, http.StatusNotFound, "review not found")
		return
	}
	if err != nil {
		httpx.WriteError(w, http.StatusInternalServerError, "failed to get review")
		return
	}
	if existing.UserID != userID {
		httpx.WriteError(w, http.StatusForbidden, "cannot update another user's review")
		return
	}

	var req UpdateRequest
	if err := httpx.ReadJSON(r, &req); err != nil {
		httpx.WriteError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	review, err := h.service.Update(r.Context(), id, UpdateReviewParams{
		Rating:   req.Rating,
		Content:  req.Content,
		IsPublic: req.IsPublic,
	})
	if errors.Is(err, ErrInvalidReview) {
		httpx.WriteError(w, http.StatusBadRequest, "invalid review")
		return
	}
	if err != nil {
		httpx.WriteError(w, http.StatusInternalServerError, "failed to update review")
		return
	}

	httpx.WriteJSON(w, http.StatusOK, review)
}

func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	userID, ok := auth.UserIDFromContext(r.Context())
	if !ok {
		httpx.WriteError(w, http.StatusUnauthorized, "authentication required")
		return
	}
	id, ok := parsePathUUID(w, r, "id", "review ID")
	if !ok {
		return
	}

	existing, err := h.service.GetByID(r.Context(), id)
	if errors.Is(err, ErrReviewNotFound) {
		httpx.WriteError(w, http.StatusNotFound, "review not found")
		return
	}
	if err != nil {
		httpx.WriteError(w, http.StatusInternalServerError, "failed to get review")
		return
	}
	if existing.UserID != userID {
		httpx.WriteError(w, http.StatusForbidden, "cannot delete another user's review")
		return
	}

	if err := h.service.Delete(r.Context(), id); err != nil {
		httpx.WriteError(w, http.StatusInternalServerError, "failed to delete review")
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
