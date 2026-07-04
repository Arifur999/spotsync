package main

import (
	"log"
	"net/http"

	"github.com/Arifur999/spotsync/config"
	"github.com/Arifur999/spotsync/handler"
	"github.com/Arifur999/spotsync/models"
	"github.com/Arifur999/spotsync/repository"
	"github.com/Arifur999/spotsync/router"
	"github.com/Arifur999/spotsync/service"
	"github.com/Arifur999/spotsync/utils"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type customValidator struct {
	validator *validator.Validate
}

func (cv *customValidator) Validate(i any) error {
	return cv.validator.Struct(i)
}

func main() {
	cfg := config.LoadEnv()
	db := config.ConnectDatabase(cfg)

	if err := db.AutoMigrate(&models.User{}, &models.ParkingZone{}, &models.Reservation{}); err != nil {
		log.Fatalf("failed to run migrations: %v", err)
	}

	e := echo.New()
	e.Validator = &customValidator{validator: validator.New()}
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())

	e.GET("/health", func(c echo.Context) error {
		return utils.Success(c, http.StatusOK, "SpotSync API is running", nil)
	})

	// dependency injection: repository -> service -> handler
	userRepo := repository.NewUserRepository(db)
	authService := service.NewAuthService(userRepo, cfg.JWTSecret)
	authHandler := handler.NewAuthHandler(authService)

	zoneRepo := repository.NewZoneRepository(db)
	zoneService := service.NewZoneService(zoneRepo)
	zoneHandler := handler.NewZoneHandler(zoneService)

	router.SetupRoutes(e, authHandler, zoneHandler, cfg.JWTSecret)

	port := cfg.Port
	if port == "" {
		port = "8080"
	}

	log.Fatal(e.Start(":" + port))
}
