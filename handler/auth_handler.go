package handler

import (
	"errors"
	"net/http"

	"github.com/Arifur999/spotsync/dto"
	"github.com/Arifur999/spotsync/repository"
	"github.com/Arifur999/spotsync/service"
	"github.com/Arifur999/spotsync/utils"

	"github.com/labstack/echo/v4"
)

type AuthHandler struct {
	authService service.AuthService
}

func NewAuthHandler(authService service.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

func (h *AuthHandler) Register(c echo.Context) error {
	var req dto.RegisterRequest
	if err := c.Bind(&req); err != nil {
		return utils.Fail(c, http.StatusBadRequest, "Invalid request payload", err.Error())
	}

	if err := c.Validate(&req); err != nil {
		return utils.Fail(c, http.StatusBadRequest, "Validation failed", err.Error())
	}

	user, err := h.authService.Register(req)
	if err != nil {
		if errors.Is(err, repository.ErrEmailAlreadyExists) {
			return utils.Fail(c, http.StatusBadRequest, "Registration failed", err.Error())
		}
		return utils.InternalError(c, err)
	}

	return utils.Success(c, http.StatusCreated, "User registered successfully", user)
}

func (h *AuthHandler) Login(c echo.Context) error {
	var req dto.LoginRequest
	if err := c.Bind(&req); err != nil {
		return utils.Fail(c, http.StatusBadRequest, "Invalid request payload", err.Error())
	}

	if err := c.Validate(&req); err != nil {
		return utils.Fail(c, http.StatusBadRequest, "Validation failed", err.Error())
	}

	loginResp, err := h.authService.Login(req)
	if err != nil {
		if errors.Is(err, service.ErrInvalidCredentials) {
			return utils.Fail(c, http.StatusUnauthorized, "Login failed", err.Error())
		}
		return utils.InternalError(c, err)
	}

	return utils.Success(c, http.StatusOK, "Login successful", loginResp)
}
