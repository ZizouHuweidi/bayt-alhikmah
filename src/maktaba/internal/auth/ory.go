package auth

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
	kratos "github.com/ory/kratos-client-go"
)

const (
	SessionContextKey  = "session"
	IdentityContextKey = "identity"
)

// Client wraps the Ory Kratos client
type Client struct {
	publicClient *kratos.APIClient
	adminClient  *kratos.APIClient
}

// NewClient creates a new Ory Kratos client
func NewClient(publicURL, adminURL string) *Client {
	publicConfig := kratos.NewConfiguration()
	publicConfig.Servers = kratos.ServerConfigurations{
		{URL: publicURL},
	}

	adminConfig := kratos.NewConfiguration()
	adminConfig.Servers = kratos.ServerConfigurations{
		{URL: adminURL},
	}

	return &Client{
		publicClient: kratos.NewAPIClient(publicConfig),
		adminClient:  kratos.NewAPIClient(adminConfig),
	}
}

// ToSession validates the session token and returns the session
func (c *Client) ToSession(ctx context.Context, token string) (*kratos.Session, error) {
	req := c.publicClient.FrontendAPI.ToSession(ctx)
	if token != "" {
		req = req.XSessionToken(token)
	}

	session, _, err := req.Execute()
	if err != nil {
		return nil, fmt.Errorf("failed to validate session: %w", err)
	}

	return session, nil
}

// GetSessionFromRequest extracts and validates the session from the request
func (c *Client) GetSessionFromRequest(ctx context.Context, r *http.Request) (*kratos.Session, error) {
	// Check for session token in header
	token := r.Header.Get("X-Session-Token")
	if token != "" {
		return c.ToSession(ctx, token)
	}

	// Check for session cookie
	cookie, err := r.Cookie("ory_session")
	if err == nil && cookie.Value != "" {
		return c.ToSession(ctx, "")
	}

	// Check for Authorization header with Bearer token
	authHeader := r.Header.Get("Authorization")
	if strings.HasPrefix(authHeader, "Bearer ") {
		token = strings.TrimPrefix(authHeader, "Bearer ")
		return c.ToSession(ctx, token)
	}

	return nil, fmt.Errorf("no session token or cookie found")
}

// Middleware creates an Echo middleware that validates sessions
func (c *Client) Middleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(echoCtx echo.Context) error {
			ctx := echoCtx.Request().Context()

			session, err := c.GetSessionFromRequest(ctx, echoCtx.Request())
			if err != nil {
				return echoCtx.JSON(http.StatusUnauthorized, map[string]string{
					"error": "unauthorized",
				})
			}

			// Store session and identity in context
			echoCtx.Set(SessionContextKey, session)
			if session.Identity != nil {
				echoCtx.Set(IdentityContextKey, session.Identity)
			}

			return next(echoCtx)
		}
	}
}

// OptionalMiddleware creates middleware that validates sessions but allows anonymous access
func (c *Client) OptionalMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(echoCtx echo.Context) error {
			ctx := echoCtx.Request().Context()

			session, err := c.GetSessionFromRequest(ctx, echoCtx.Request())
			if err == nil && session != nil {
				echoCtx.Set(SessionContextKey, session)
				if session.Identity != nil {
					echoCtx.Set(IdentityContextKey, session.Identity)
				}
			}

			return next(echoCtx)
		}
	}
}

// GetSession retrieves the session from the Echo context
func GetSession(c echo.Context) *kratos.Session {
	session, ok := c.Get(SessionContextKey).(*kratos.Session)
	if !ok {
		return nil
	}
	return session
}

// GetIdentity retrieves the identity from the Echo context
func GetIdentity(c echo.Context) *kratos.Identity {
	identity, ok := c.Get(IdentityContextKey).(*kratos.Identity)
	if !ok {
		return nil
	}
	return identity
}

// GetUserID retrieves the user ID from the Echo context
func GetUserID(c echo.Context) string {
	identity := GetIdentity(c)
	if identity == nil {
		return ""
	}
	return identity.Id
}
