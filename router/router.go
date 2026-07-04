package router

import (
	"github.com/Arifur999/spotsync/handler"
	"github.com/Arifur999/spotsync/middleware"
	"github.com/Arifur999/spotsync/models"

	"github.com/labstack/echo/v4"
)

func SetupRoutes(
	e *echo.Echo,
	authHandler *handler.AuthHandler,
	zoneHandler *handler.ZoneHandler,
	reservationHandler *handler.ReservationHandler,
	jwtSecret string,
) {
	api := e.Group("/api/v1")

	auth := api.Group("/auth")
	auth.POST("/register", authHandler.Register)
	auth.POST("/login", authHandler.Login)

	zones := api.Group("/zones")
	zones.GET("", zoneHandler.GetAllZones)
	zones.GET("/:id", zoneHandler.GetZoneByID)
	zones.POST("", zoneHandler.CreateZone, middleware.JWTAuth(jwtSecret), middleware.RequireRole(models.RoleAdmin))

	reservations := api.Group("/reservations", middleware.JWTAuth(jwtSecret))
	reservations.POST("", reservationHandler.CreateReservation)
	reservations.GET("/my-reservations", reservationHandler.GetMyReservations)
	reservations.DELETE("/:id", reservationHandler.CancelReservation)
	reservations.GET("", reservationHandler.GetAllReservations, middleware.RequireRole(models.RoleAdmin))
}
