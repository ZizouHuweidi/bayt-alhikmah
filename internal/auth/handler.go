package auth

import (
	"errors"
	"log/slog"
	"net"
	"net/http"
	"sync"
	"time"

	"github.com/gofrs/uuid/v5"
	"github.com/labstack/echo/v5"
	"github.com/zizouhuweidi/maktaba/internal/echox"
)

const refreshCookieName = "bh_refresh_token"

type Handler struct {
	service      *Service
	cookieSecure bool
	logger       *slog.Logger
	limiter      *rateLimiter
}

type registerRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Username string `json:"username" validate:"required,min=3,max=32"`
	Password string `json:"password" validate:"required,min=12"`
}

type loginRequest struct {
	Login    string `json:"login" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type refreshRequest struct {
	RefreshToken string `json:"refresh_token"`
}

type authResponse struct {
	User   *User      `json:"user"`
	Tokens AuthTokens `json:"tokens"`
}

func NewHandler(service *Service, cookieSecure bool, logger *slog.Logger) *Handler {
	return &Handler{service: service, cookieSecure: cookieSecure, logger: logger, limiter: newRateLimiter(10, time.Minute)}
}

func (h *Handler) RegisterRoutes(e *echo.Echo) {
	e.POST("/auth/register", h.Register)
	e.POST("/auth/login", h.Login)
	e.POST("/auth/refresh", h.Refresh)
	e.POST("/auth/logout", h.Logout)
}

func (h *Handler) RegisterProtectedRoutes(g *echo.Group) {
	g.GET("/me", h.Me)
}

func (h *Handler) Register(c *echo.Context) error {
	if err := h.allowAuthAttempt(c, "register"); err != nil {
		return err
	}

	var req registerRequest
	if err := echox.BindAndValidate(c, &req); err != nil {
		return err
	}

	user, tokens, err := h.service.Register(c.Request().Context(), req.Email, req.Username, req.Password)
	if errors.Is(err, ErrInvalidSignup) {
		return echo.NewHTTPError(http.StatusBadRequest, "email, username, and password are required; password must be at least 12 characters")
	}
	if err != nil {
		h.logger.Error("registration failed", "error", err)
		return echo.NewHTTPError(http.StatusConflict, "email or username is unavailable")
	}

	h.setRefreshCookie(c, tokens.RefreshToken)
	return c.JSON(http.StatusCreated, authResponse{User: user, Tokens: tokens})
}

func (h *Handler) Login(c *echo.Context) error {
	if err := h.allowAuthAttempt(c, "login"); err != nil {
		return err
	}

	var req loginRequest
	if err := echox.BindAndValidate(c, &req); err != nil {
		return err
	}

	user, tokens, err := h.service.Login(c.Request().Context(), req.Login, req.Password)
	if errors.Is(err, ErrInvalidCredentials) {
		return echo.NewHTTPError(http.StatusUnauthorized, "invalid credentials")
	}
	if err != nil {
		h.logger.Error("login failed", "error", err)
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to login")
	}

	h.setRefreshCookie(c, tokens.RefreshToken)
	return c.JSON(http.StatusOK, authResponse{User: user, Tokens: tokens})
}

func (h *Handler) Refresh(c *echo.Context) error {
	if err := h.allowAuthAttempt(c, "refresh"); err != nil {
		return err
	}

	var token string

	// Accept refresh_token from request body (mobile native apps)
	var req refreshRequest
	if err := c.Bind(&req); err == nil && req.RefreshToken != "" {
		token = req.RefreshToken
	} else {
		// Fall back to cookie (web SPA — automatic browser behavior)
		cookie, err := c.Cookie(refreshCookieName)
		if err != nil || cookie.Value == "" {
			return echo.NewHTTPError(http.StatusUnauthorized, "refresh token required")
		}
		token = cookie.Value
	}

	result, err := h.service.Refresh(c.Request().Context(), token)
	if errors.Is(err, ErrInvalidRefresh) {
		h.clearRefreshCookie(c)
		return echo.NewHTTPError(http.StatusUnauthorized, "invalid refresh token")
	}
	if err != nil {
		h.logger.Error("refresh failed", "error", err)
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to refresh token")
	}

	h.setRefreshCookie(c, result.RefreshToken)
	return c.JSON(http.StatusOK, result)
}

func (h *Handler) Logout(c *echo.Context) error {
	var token string

	var req refreshRequest
	if err := c.Bind(&req); err == nil && req.RefreshToken != "" {
		token = req.RefreshToken
	} else {
		cookie, err := c.Cookie(refreshCookieName)
		if err == nil && cookie.Value != "" {
			token = cookie.Value
		}
	}

	if token != "" {
		if err := h.service.Logout(c.Request().Context(), token); err != nil {
			h.logger.Error("logout failed", "error", err)
		}
	}

	h.clearRefreshCookie(c)
	return c.NoContent(http.StatusNoContent)
}

func (h *Handler) Me(c *echo.Context) error {
	userID, ok := UserID(c)
	if !ok {
		return echo.NewHTTPError(http.StatusUnauthorized, "authentication required")
	}

	user, err := h.service.GetUser(c.Request().Context(), userID)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, "user not found")
	}
	return c.JSON(http.StatusOK, user)
}

func (h *Handler) Middleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c *echo.Context) error {
		const prefix = "Bearer "
		header := c.Request().Header.Get(echo.HeaderAuthorization)
		if len(header) <= len(prefix) || header[:len(prefix)] != prefix {
			return echo.NewHTTPError(http.StatusUnauthorized, "authentication required")
		}
		token := header[len(prefix):]

		claims, err := h.service.VerifyAccessToken(token)
		if err != nil {
			return echo.NewHTTPError(http.StatusUnauthorized, "invalid access token")
		}
		userID, err := uuid.FromString(claims.Subject)
		if err != nil {
			return echo.NewHTTPError(http.StatusUnauthorized, "invalid access token subject")
		}

		SetUserID(c, userID)
		return next(c)
	}
}

func (h *Handler) setRefreshCookie(c *echo.Context, token string) {
	c.SetCookie(&http.Cookie{
		Name:     refreshCookieName,
		Value:    token,
		Path:     "/auth/refresh",
		HttpOnly: true,
		Secure:   h.cookieSecure,
		SameSite: http.SameSiteStrictMode,
		MaxAge:   int((30 * 24 * time.Hour).Seconds()),
	})
}

func (h *Handler) clearRefreshCookie(c *echo.Context) {
	c.SetCookie(&http.Cookie{Name: refreshCookieName, Path: "/auth/refresh", MaxAge: -1, HttpOnly: true, Secure: h.cookieSecure, SameSite: http.SameSiteStrictMode})
}

func (h *Handler) allowAuthAttempt(c *echo.Context, action string) error {
	key := clientIP(c) + ":" + action
	if h.limiter.allow(key) {
		return nil
	}
	return echo.NewHTTPError(http.StatusTooManyRequests, "too many attempts")
}

func clientIP(c *echo.Context) string {
	r := c.Request()
	if forwardedFor := r.Header.Get("X-Forwarded-For"); forwardedFor != "" {
		if host, _, err := net.SplitHostPort(forwardedFor); err == nil {
			return host
		}
		return forwardedFor
	}
	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return r.RemoteAddr
	}
	return host
}

type rateLimiter struct {
	mu       sync.Mutex
	limit    int
	window   time.Duration
	attempts map[string]rateLimitEntry
}

type rateLimitEntry struct {
	count     int
	resetTime time.Time
}

func newRateLimiter(limit int, window time.Duration) *rateLimiter {
	return &rateLimiter{limit: limit, window: window, attempts: make(map[string]rateLimitEntry)}
}

func (l *rateLimiter) allow(key string) bool {
	l.mu.Lock()
	defer l.mu.Unlock()

	now := time.Now()
	entry := l.attempts[key]
	if now.After(entry.resetTime) {
		l.attempts[key] = rateLimitEntry{count: 1, resetTime: now.Add(l.window)}
		return true
	}
	if entry.count >= l.limit {
		return false
	}
	entry.count++
	l.attempts[key] = entry
	return true
}
