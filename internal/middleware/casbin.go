package middleware

import (
	"yalla-go/internal/json"
	"yalla-go/internal/jwt"
	"yalla-go/internal/rbac"

	"github.com/labstack/echo/v4"
)

func CasbinMiddleware(rbacService *rbac.Service) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			user := c.Get("user")
			if user == nil {
				return json.Unauthorized(c, "Missing user context")
			}

			claims, ok := user.(*jwt.Claims)
			if !ok {
				return json.Unauthorized(c, "Invalid user context")
			}

			// Subject: Username (simplifies policy management)
			// Object: Request Path
			// Action: Request Method
			sub := claims.Username // Assuming Username is in claims, if not we might need to fetch it or use ID
			obj := c.Path()
			act := c.Request().Method

			allowed, err := rbacService.Enforce(sub, obj, act)
			if err != nil {
				return json.InternalServerError(c, err)
			}

			if !allowed {
				return json.Forbidden(c, "Access denied")
			}

			return next(c)
		}
	}
}
