package handler

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/Arifur999/spotsync/dto"
	"github.com/Arifur999/spotsync/repository"
	"github.com/Arifur999/spotsync/service"
	"github.com/Arifur999/spotsync/utils"

	"github.com/labstack/echo/v4"
)

type ZoneHandler struct {
	zoneService service.ZoneService
}

func NewZoneHandler(zoneService service.ZoneService) *ZoneHandler {
	return &ZoneHandler{zoneService: zoneService}
}

func (h *ZoneHandler) CreateZone(c echo.Context) error {
	var req dto.CreateZoneRequest
	if err := c.Bind(&req); err != nil {
		return utils.Fail(c, http.StatusBadRequest, "Invalid request payload", err.Error())
	}

	if err := c.Validate(&req); err != nil {
		return utils.Fail(c, http.StatusBadRequest, "Validation failed", err.Error())
	}

	zone, err := h.zoneService.CreateZone(req)
	if err != nil {
		return utils.InternalError(c, err)
	}

	return utils.Success(c, http.StatusCreated, "Parking zone created successfully", zone)
}

func (h *ZoneHandler) GetAllZones(c echo.Context) error {
	zones, err := h.zoneService.GetAllZones()
	if err != nil {
		return utils.InternalError(c, err)
	}

	return utils.Success(c, http.StatusOK, "Parking zones retrieved successfully", zones)
}

func (h *ZoneHandler) GetZoneByID(c echo.Context) error {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		return utils.Fail(c, http.StatusBadRequest, "Invalid zone id", err.Error())
	}

	zone, err := h.zoneService.GetZoneByID(uint(id))
	if err != nil {
		if errors.Is(err, repository.ErrZoneNotFound) {
			return utils.Fail(c, http.StatusNotFound, "Parking zone not found", err.Error())
		}
		return utils.InternalError(c, err)
	}

	return utils.Success(c, http.StatusOK, "Parking zone retrieved successfully", zone)
}
