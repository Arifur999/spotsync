package middleware

import (
	"net/http"
	"strings"

	"github.com/Arifur999/spotsync/utils"

	"github.com/labstack/echo/v4"
)

const (
	ContextUserIDKey = "user_id"
	ContextRoleKey   = "role"
)

// JWTAuth verifies the Bearer token on protected routes and injects the
// requester's id and role into the Echo context for downstream handlers.
func JWTAuth(secret string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			authHeader := c.Request().Header.Get("Authorization")
			if authHeader == "" {
				return utils.Fail(c, http.StatusUnauthorized, "Unauthorized", "missing authorization header")
			}

			parts := strings.SplitN(authHeader, " ", 2)
			if len(parts) != 2 || parts[0] != "Bearer" {
				return utils.Fail(c, http.StatusUnauthorized, "Unauthorized", "invalid authorization header format")
			}

			claims, err := utils.ValidateToken(parts[1], secret)
			if err != nil {
				return utils.Fail(c, http.StatusUnauthorized, "Unauthorized", err.Error())
			}

			c.Set(ContextUserIDKey, claims.UserID)
			c.Set(ContextRoleKey, claims.Role)

			return next(c)
		}
	}
}

// GetUserID reads the authenticated user's id set by JWTAuth.
func GetUserID(c echo.Context) (uint, bool) {
	id, ok := c.Get(ContextUserIDKey).(uint)
	return id, ok
}

// GetRole reads the authenticated user's role set by JWTAuth.
func GetRole(c echo.Context) (string, bool) {
	role, ok := c.Get(ContextRoleKey).(string)
	return role, ok
}
