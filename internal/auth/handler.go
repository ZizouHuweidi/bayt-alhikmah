package auth

import (
	"errors"
	"log/slog"
	"net/http"
	"time"

	"github.com/gofrs/uuid/v5"
	"github.com/zizouhuweidi/maktaba/internal/httpx"
)

const refreshCookieName = "bh_refresh_token"

type Handler struct {
	service      *Service
	cookieSecure bool
	logger       *slog.Logger
}

type registerRequest struct {
	Email    string `json:"email"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type loginRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type authResponse struct {
	User   *User      `json:"user"`
	Tokens AuthTokens `json:"tokens"`
}

func NewHandler(service *Service, cookieSecure bool, logger *slog.Logger) *Handler {
	return &Handler{service: service, cookieSecure: cookieSecure, logger: logger}
}

func (h *Handler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /auth/register", h.Register)
	mux.HandleFunc("POST /auth/login", h.Login)
	mux.HandleFunc("POST /auth/refresh", h.Refresh)
}

func (h *Handler) RegisterProtectedRoutes(mux *http.ServeMux) {
	mux.Handle("GET /api/me", h.Middleware(http.HandlerFunc(h.Me)))
}

func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
	var req registerRequest
	if err := httpx.ReadJSON(r, &req); err != nil {
		httpx.WriteError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	user, tokens, err := h.service.Register(r.Context(), req.Email, req.Username, req.Password)
	if errors.Is(err, ErrInvalidSignup) {
		httpx.WriteError(w, http.StatusBadRequest, "email, username, and password are required; password must be at least 12 characters")
		return
	}
	if err != nil {
		h.logger.Error("registration failed", "error", err)
		httpx.WriteError(w, http.StatusConflict, "email or username is unavailable")
		return
	}

	h.setRefreshCookie(w, tokens.RefreshToken)
	tokens.RefreshToken = ""
	httpx.WriteJSON(w, http.StatusCreated, authResponse{User: user, Tokens: tokens})
}

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	var req loginRequest
	if err := httpx.ReadJSON(r, &req); err != nil {
		httpx.WriteError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	user, tokens, err := h.service.Login(r.Context(), req.Login, req.Password)
	if errors.Is(err, ErrInvalidCredentials) {
		httpx.WriteError(w, http.StatusUnauthorized, "invalid credentials")
		return
	}
	if err != nil {
		h.logger.Error("login failed", "error", err)
		httpx.WriteError(w, http.StatusInternalServerError, "failed to login")
		return
	}

	h.setRefreshCookie(w, tokens.RefreshToken)
	tokens.RefreshToken = ""
	httpx.WriteJSON(w, http.StatusOK, authResponse{User: user, Tokens: tokens})
}

func (h *Handler) Refresh(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie(refreshCookieName)
	if err != nil || cookie.Value == "" {
		httpx.WriteError(w, http.StatusUnauthorized, "refresh token required")
		return
	}

	tokens, err := h.service.Refresh(r.Context(), cookie.Value)
	if errors.Is(err, ErrInvalidRefresh) {
		h.clearRefreshCookie(w)
		httpx.WriteError(w, http.StatusUnauthorized, "invalid refresh token")
		return
	}
	if err != nil {
		h.logger.Error("refresh failed", "error", err)
		httpx.WriteError(w, http.StatusInternalServerError, "failed to refresh token")
		return
	}

	h.setRefreshCookie(w, tokens.RefreshToken)
	tokens.RefreshToken = ""
	httpx.WriteJSON(w, http.StatusOK, tokens)
}

func (h *Handler) Me(w http.ResponseWriter, r *http.Request) {
	userID, ok := UserIDFromContext(r.Context())
	if !ok {
		httpx.WriteError(w, http.StatusUnauthorized, "authentication required")
		return
	}

	user, err := h.service.GetUser(r.Context(), userID)
	if err != nil {
		httpx.WriteError(w, http.StatusNotFound, "user not found")
		return
	}
	httpx.WriteJSON(w, http.StatusOK, user)
}

func (h *Handler) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token, err := httpx.BearerToken(r)
		if err != nil {
			httpx.WriteError(w, http.StatusUnauthorized, "authentication required")
			return
		}

		claims, err := h.service.VerifyAccessToken(token)
		if err != nil {
			httpx.WriteError(w, http.StatusUnauthorized, "invalid access token")
			return
		}
		userID, err := uuid.FromString(claims.Subject)
		if err != nil {
			httpx.WriteError(w, http.StatusUnauthorized, "invalid access token subject")
			return
		}

		next.ServeHTTP(w, r.WithContext(ContextWithUserID(r.Context(), userID)))
	})
}

func (h *Handler) setRefreshCookie(w http.ResponseWriter, token string) {
	http.SetCookie(w, &http.Cookie{
		Name:     refreshCookieName,
		Value:    token,
		Path:     "/auth/refresh",
		HttpOnly: true,
		Secure:   h.cookieSecure,
		SameSite: http.SameSiteLaxMode,
		MaxAge:   int((30 * 24 * time.Hour).Seconds()),
	})
}

func (h *Handler) clearRefreshCookie(w http.ResponseWriter) {
	http.SetCookie(w, &http.Cookie{Name: refreshCookieName, Path: "/auth/refresh", MaxAge: -1, HttpOnly: true, Secure: h.cookieSecure, SameSite: http.SameSiteLaxMode})
}
