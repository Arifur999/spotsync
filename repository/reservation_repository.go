package repository

import (
	"errors"

	"github.com/Arifur999/spotsync/models"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

var (
	ErrZoneFull            = errors.New("parking zone is full")
	ErrReservationNotFound = errors.New("reservation not found")
)

type ReservationRepository interface {
	CreateWithLock(userID, zoneID uint, licensePlate string) (*models.Reservation, error)
	GetByUserID(userID uint) ([]models.Reservation, error)
	GetAll() ([]models.Reservation, error)
	GetByID(id uint) (*models.Reservation, error)
	Cancel(reservation *models.Reservation) error
}

type reservationRepository struct {
	db *gorm.DB
}

func NewReservationRepository(db *gorm.DB) ReservationRepository {
	return &reservationRepository{db: db}
}

// CreateWithLock is the concurrency-critical path: it locks the zone row
// FOR UPDATE inside a transaction so that two simultaneous requests for the
// last spot in a zone are serialized instead of both reading a stale
// available count and over-booking the zone.
func (r *reservationRepository) CreateWithLock(userID, zoneID uint, licensePlate string) (*models.Reservation, error) {
	var reservation models.Reservation

	err := r.db.Transaction(func(tx *gorm.DB) error {
		var zone models.ParkingZone

		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).First(&zone, zoneID).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return ErrZoneNotFound
			}
			return err
		}

		var activeCount int64
		if err := tx.Model(&models.Reservation{}).
			Where("zone_id = ? AND status = ?", zoneID, models.ReservationActive).
			Count(&activeCount).Error; err != nil {
			return err
		}

		if int(activeCount) >= zone.TotalCapacity {
			return ErrZoneFull
		}

		reservation = models.Reservation{
			UserID:       userID,
			ZoneID:       zoneID,
			LicensePlate: licensePlate,
			Status:       models.ReservationActive,
		}

		return tx.Create(&reservation).Error
	})

	if err != nil {
		return nil, err
	}

	return &reservation, nil
}

func (r *reservationRepository) GetByUserID(userID uint) ([]models.Reservation, error) {
	var reservations []models.Reservation

	err := r.db.Preload("Zone").
		Where("user_id = ?", userID).
		Order("id desc").
		Find(&reservations).Error
	if err != nil {
		return nil, err
	}

	return reservations, nil
}

func (r *reservationRepository) GetAll() ([]models.Reservation, error) {
	var reservations []models.Reservation

	err := r.db.Preload("User").Preload("Zone").
		Order("id desc").
		Find(&reservations).Error
	if err != nil {
		return nil, err
	}

	return reservations, nil
}

func (r *reservationRepository) GetByID(id uint) (*models.Reservation, error) {
	var reservation models.Reservation

	err := r.db.First(&reservation, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrReservationNotFound
		}
		return nil, err
	}

	return &reservation, nil
}

func (r *reservationRepository) Cancel(reservation *models.Reservation) error {
	reservation.Status = models.ReservationCancelled
	return r.db.Save(reservation).Error
}
