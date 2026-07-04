package handler

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/Arifur999/spotsync/dto"
	"github.com/Arifur999/spotsync/middleware"
	"github.com/Arifur999/spotsync/repository"
	"github.com/Arifur999/spotsync/service"
	"github.com/Arifur999/spotsync/utils"

	"github.com/labstack/echo/v4"
)

type ReservationHandler struct {
	reservationService service.ReservationService
}

func NewReservationHandler(reservationService service.ReservationService) *ReservationHandler {
	return &ReservationHandler{reservationService: reservationService}
}

func (h *ReservationHandler) CreateReservation(c echo.Context) error {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		return utils.Fail(c, http.StatusUnauthorized, "Unauthorized", "missing user in request context")
	}

	var req dto.CreateReservationRequest
	if err := c.Bind(&req); err != nil {
		return utils.Fail(c, http.StatusBadRequest, "Invalid request payload", err.Error())
	}

	if err := c.Validate(&req); err != nil {
		return utils.Fail(c, http.StatusBadRequest, "Validation failed", err.Error())
	}

	reservation, err := h.reservationService.CreateReservation(userID, req)
	if err != nil {
		return reservationErrorResponse(c, err)
	}

	return utils.Success(c, http.StatusCreated, "Reservation confirmed successfully", reservation)
}

func (h *ReservationHandler) GetMyReservations(c echo.Context) error {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		return utils.Fail(c, http.StatusUnauthorized, "Unauthorized", "missing user in request context")
	}

	reservations, err := h.reservationService.GetMyReservations(userID)
	if err != nil {
		return utils.Fail(c, http.StatusInternalServerError, "Failed to retrieve reservations", err.Error())
	}

	return utils.Success(c, http.StatusOK, "My reservations retrieved successfully", reservations)
}

func (h *ReservationHandler) CancelReservation(c echo.Context) error {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		return utils.Fail(c, http.StatusUnauthorized, "Unauthorized", "missing user in request context")
	}
	role, _ := middleware.GetRole(c)

	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		return utils.Fail(c, http.StatusBadRequest, "Invalid reservation id", err.Error())
	}

	if err := h.reservationService.CancelReservation(userID, role, uint(id)); err != nil {
		return reservationErrorResponse(c, err)
	}

	return utils.Success(c, http.StatusOK, "Reservation cancelled successfully", nil)
}

func (h *ReservationHandler) GetAllReservations(c echo.Context) error {
	reservations, err := h.reservationService.GetAllReservations()
	if err != nil {
		return utils.Fail(c, http.StatusInternalServerError, "Failed to retrieve reservations", err.Error())
	}

	return utils.Success(c, http.StatusOK, "Reservations retrieved successfully", reservations)
}

func reservationErrorResponse(c echo.Context, err error) error {
	switch {
	case errors.Is(err, repository.ErrZoneNotFound):
		return utils.Fail(c, http.StatusNotFound, "Parking zone not found", err.Error())
	case errors.Is(err, repository.ErrZoneFull):
		return utils.Fail(c, http.StatusConflict, "Parking zone is full", err.Error())
	case errors.Is(err, repository.ErrReservationNotFound):
		return utils.Fail(c, http.StatusNotFound, "Reservation not found", err.Error())
	case errors.Is(err, service.ErrForbiddenReservationAccess):
		return utils.Fail(c, http.StatusForbidden, "Forbidden", err.Error())
	case errors.Is(err, service.ErrReservationAlreadyCancelled):
		return utils.Fail(c, http.StatusConflict, "Reservation already cancelled", err.Error())
	default:
		return utils.Fail(c, http.StatusInternalServerError, "Something went wrong", err.Error())
	}
}
