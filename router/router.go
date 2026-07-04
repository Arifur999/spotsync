package router

import (
	"github.com/Arifur999/spotsync/handler"

	"github.com/labstack/echo/v4"
)

func SetupRoutes(e *echo.Echo, authHandler *handler.AuthHandler) {
	api := e.Group("/api/v1")

	auth := api.Group("/auth")
	auth.POST("/register", authHandler.Register)
	auth.POST("/login", authHandler.Login)
}
