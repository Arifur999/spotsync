package main

import (
	"log"
	"net/http"

	"github.com/Arifur999/spotsync/config"
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
	_ = db // wired into repositories in later steps

	e := echo.New()
	e.Validator = &customValidator{validator: validator.New()}
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())

	e.GET("/health", func(c echo.Context) error {
		return utils.Success(c, http.StatusOK, "SpotSync API is running", nil)
	})

	port := cfg.Port
	if port == "" {
		port = "8080"
	}

	log.Fatal(e.Start(":" + port))
}
