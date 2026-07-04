package utils

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

type SuccessResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
}

type ErrorResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Errors  any    `json:"errors,omitempty"`
}

// Success writes the standard {success, message, data} envelope.
func Success(c echo.Context, status int, message string, data any) error {
	return c.JSON(status, SuccessResponse{
		Success: true,
		Message: message,
		Data:    data,
	})
}

// Fail writes the standard {success, message, errors} envelope.
func Fail(c echo.Context, status int, message string, errDetails any) error {
	return c.JSON(status, ErrorResponse{
		Success: false,
		Message: message,
		Errors:  errDetails,
	})
}

// InternalError logs the real (possibly raw GORM/DB) error server-side and
// returns a generic 500 response so internal details never reach the client.
func InternalError(c echo.Context, err error) error {
	c.Logger().Error(err)
	return Fail(c, http.StatusInternalServerError, "Something went wrong", nil)
}
