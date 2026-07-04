package service

import (
	"github.com/Arifur999/spotsync/dto"
	"github.com/Arifur999/spotsync/models"
	"github.com/Arifur999/spotsync/repository"
)

type ZoneService interface {
	CreateZone(req dto.CreateZoneRequest) (*dto.ZoneResponse, error)
	GetAllZones() ([]dto.ZoneAvailabilityResponse, error)
	GetZoneByID(id uint) (*dto.ZoneAvailabilityResponse, error)
}

type zoneService struct {
	zoneRepo repository.ZoneRepository
}

func NewZoneService(zoneRepo repository.ZoneRepository) ZoneService {
	return &zoneService{zoneRepo: zoneRepo}
}

func (s *zoneService) CreateZone(req dto.CreateZoneRequest) (*dto.ZoneResponse, error) {
	zone := models.ParkingZone{
		Name:          req.Name,
		Type:          req.Type,
		TotalCapacity: req.TotalCapacity,
		PricePerHour:  req.PricePerHour,
	}

	if err := s.zoneRepo.Create(&zone); err != nil {
		return nil, err
	}

	return &dto.ZoneResponse{
		ID:            zone.ID,
		Name:          zone.Name,
		Type:          zone.Type,
		TotalCapacity: zone.TotalCapacity,
		PricePerHour:  zone.PricePerHour,
		CreatedAt:     zone.CreatedAt,
		UpdatedAt:     zone.UpdatedAt,
	}, nil
}

func (s *zoneService) GetAllZones() ([]dto.ZoneAvailabilityResponse, error) {
	zones, err := s.zoneRepo.GetAll()
	if err != nil {
		return nil, err
	}

	responses := make([]dto.ZoneAvailabilityResponse, len(zones))
	for i, z := range zones {
		responses[i] = toZoneAvailabilityResponse(z)
	}

	return responses, nil
}

func (s *zoneService) GetZoneByID(id uint) (*dto.ZoneAvailabilityResponse, error) {
	zone, err := s.zoneRepo.GetByID(id)
	if err != nil {
		return nil, err
	}

	resp := toZoneAvailabilityResponse(*zone)
	return &resp, nil
}

func toZoneAvailabilityResponse(z repository.ZoneWithAvailability) dto.ZoneAvailabilityResponse {
	return dto.ZoneAvailabilityResponse{
		ID:             z.ID,
		Name:           z.Name,
		Type:           z.Type,
		TotalCapacity:  z.TotalCapacity,
		AvailableSpots: z.AvailableSpots,
		PricePerHour:   z.PricePerHour,
		CreatedAt:      z.CreatedAt,
	}
}
