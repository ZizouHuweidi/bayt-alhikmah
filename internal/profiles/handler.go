package profiles

import (
	"errors"
	"log/slog"
	"net/http"

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

type UpdateRequest struct {
	DisplayName   *string `json:"display_name,omitempty"`
	Bio           *string `json:"bio,omitempty"`
	PublicProfile *bool   `json:"public_profile,omitempty"`
}

func (h *Handler) RegisterPublicRoutes(r httpx.Router) {
	r.Get("/users/:username/profile", h.GetPublicByUsername)
}

func (h *Handler) RegisterProtectedRoutes(r httpx.Router) {
	r.Get("/profile", h.GetOwn)
	r.Put("/profile", h.Update)
}

func (h *Handler) GetOwn(w http.ResponseWriter, r *http.Request) {
	userID, ok := auth.UserIDFromContext(r.Context())
	if !ok {
		httpx.WriteError(w, http.StatusUnauthorized, "authentication required")
		return
	}

	profile, err := h.service.GetOwn(r.Context(), userID)
	if err != nil {
		httpx.WriteError(w, http.StatusInternalServerError, "failed to get profile")
		return
	}

	httpx.WriteJSON(w, http.StatusOK, profile)
}

func (h *Handler) GetPublicByUsername(w http.ResponseWriter, r *http.Request) {
	profile, err := h.service.GetPublicByUsername(r.Context(), r.PathValue("username"))
	if errors.Is(err, ErrInvalidProfile) {
		httpx.WriteError(w, http.StatusBadRequest, "invalid username")
		return
	}
	if errors.Is(err, ErrProfileNotFound) {
		httpx.WriteError(w, http.StatusNotFound, "profile not found")
		return
	}
	if err != nil {
		httpx.WriteError(w, http.StatusInternalServerError, "failed to get profile")
		return
	}

	httpx.WriteJSON(w, http.StatusOK, profile)
}

func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	userID, ok := auth.UserIDFromContext(r.Context())
	if !ok {
		httpx.WriteError(w, http.StatusUnauthorized, "authentication required")
		return
	}

	var req UpdateRequest
	if err := httpx.ReadJSON(r, &req); err != nil {
		httpx.WriteError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	profile, err := h.service.Update(r.Context(), UpdateProfileParams{
		UserID:        userID,
		DisplayName:   req.DisplayName,
		Bio:           req.Bio,
		PublicProfile: req.PublicProfile,
	})
	if errors.Is(err, ErrInvalidProfile) {
		httpx.WriteError(w, http.StatusBadRequest, "invalid profile")
		return
	}
	if err != nil {
		httpx.WriteError(w, http.StatusInternalServerError, "failed to update profile")
		return
	}

	httpx.WriteJSON(w, http.StatusOK, profile)
}
