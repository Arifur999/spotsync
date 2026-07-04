package middleware

import (
	"net/http"

	"github.com/Arifur999/spotsync/utils"

	"github.com/labstack/echo/v4"
)

// RequireRole must run after JWTAuth. It rejects the request with 403 if the
// authenticated user's role is not one of the allowed roles.
func RequireRole(allowedRoles ...string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			role, ok := GetRole(c)
			if !ok {
				return utils.Fail(c, http.StatusForbidden, "Forbidden", "role not found in request context")
			}

			for _, allowed := range allowedRoles {
				if role == allowed {
					return next(c)
				}
			}

			return utils.Fail(c, http.StatusForbidden, "Forbidden", "insufficient permissions for this action")
		}
	}
}
