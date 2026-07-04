package service

import (
	"errors"

	"github.com/Arifur999/spotsync/dto"
	"github.com/Arifur999/spotsync/models"
	"github.com/Arifur999/spotsync/repository"
)

var (
	ErrForbiddenReservationAccess  = errors.New("you do not own this reservation")
	ErrReservationAlreadyCancelled = errors.New("reservation is already cancelled")
)

type ReservationService interface {
	CreateReservation(userID uint, req dto.CreateReservationRequest) (*dto.ReservationResponse, error)
	GetMyReservations(userID uint) ([]dto.MyReservationResponse, error)
	CancelReservation(userID uint, role string, reservationID uint) error
	GetAllReservations() ([]dto.AdminReservationResponse, error)
}

type reservationService struct {
	reservationRepo repository.ReservationRepository
}

func NewReservationService(reservationRepo repository.ReservationRepository) ReservationService {
	return &reservationService{reservationRepo: reservationRepo}
}

func (s *reservationService) CreateReservation(userID uint, req dto.CreateReservationRequest) (*dto.ReservationResponse, error) {
	reservation, err := s.reservationRepo.CreateWithLock(userID, req.ZoneID, req.LicensePlate)
	if err != nil {
		return nil, err
	}

	return &dto.ReservationResponse{
		ID:           reservation.ID,
		UserID:       reservation.UserID,
		ZoneID:       reservation.ZoneID,
		LicensePlate: reservation.LicensePlate,
		Status:       reservation.Status,
		CreatedAt:    reservation.CreatedAt,
		UpdatedAt:    reservation.UpdatedAt,
	}, nil
}

func (s *reservationService) GetMyReservations(userID uint) ([]dto.MyReservationResponse, error) {
	reservations, err := s.reservationRepo.GetByUserID(userID)
	if err != nil {
		return nil, err
	}

	responses := make([]dto.MyReservationResponse, len(reservations))
	for i, res := range reservations {
		responses[i] = dto.MyReservationResponse{
			ID:           res.ID,
			LicensePlate: res.LicensePlate,
			Status:       res.Status,
			Zone: dto.ZoneSummary{
				ID:   res.Zone.ID,
				Name: res.Zone.Name,
				Type: res.Zone.Type,
			},
			CreatedAt: res.CreatedAt,
		}
	}

	return responses, nil
}

func (s *reservationService) CancelReservation(userID uint, role string, reservationID uint) error {
	reservation, err := s.reservationRepo.GetByID(reservationID)
	if err != nil {
		return err
	}

	if reservation.UserID != userID && role != models.RoleAdmin {
		return ErrForbiddenReservationAccess
	}

	if reservation.Status == models.ReservationCancelled {
		return ErrReservationAlreadyCancelled
	}

	return s.reservationRepo.Cancel(reservation)
}

func (s *reservationService) GetAllReservations() ([]dto.AdminReservationResponse, error) {
	reservations, err := s.reservationRepo.GetAll()
	if err != nil {
		return nil, err
	}

	responses := make([]dto.AdminReservationResponse, len(reservations))
	for i, res := range reservations {
		responses[i] = dto.AdminReservationResponse{
			ID:           res.ID,
			LicensePlate: res.LicensePlate,
			Status:       res.Status,
			User: dto.ReservationUserSummary{
				ID:    res.User.ID,
				Name:  res.User.Name,
				Email: res.User.Email,
			},
			Zone: dto.ZoneSummary{
				ID:   res.Zone.ID,
				Name: res.Zone.Name,
				Type: res.Zone.Type,
			},
			CreatedAt: res.CreatedAt,
			UpdatedAt: res.UpdatedAt,
		}
	}

	return responses, nil
}
