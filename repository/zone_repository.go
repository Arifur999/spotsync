package repository

import (
	"errors"

	"github.com/Arifur999/spotsync/models"

	"gorm.io/gorm"
)

var ErrZoneNotFound = errors.New("parking zone not found")

// zoneAvailabilitySelect computes available_spots dynamically as
// total_capacity minus the count of that zone's active reservations.
const zoneAvailabilitySelect = `parking_zones.*, ` +
	`(parking_zones.total_capacity - COALESCE((` +
	`SELECT COUNT(*) FROM reservations ` +
	`WHERE reservations.zone_id = parking_zones.id AND reservations.status = 'active'` +
	`), 0)) AS available_spots`

// ZoneWithAvailability wraps a ParkingZone with its dynamically computed
// available spot count, used by the list/get-by-id queries.
type ZoneWithAvailability struct {
	models.ParkingZone
	AvailableSpots int `json:"-"`
}

type ZoneRepository interface {
	Create(zone *models.ParkingZone) error
	GetAll() ([]ZoneWithAvailability, error)
	GetByID(id uint) (*ZoneWithAvailability, error)
}

type zoneRepository struct {
	db *gorm.DB
}

func NewZoneRepository(db *gorm.DB) ZoneRepository {
	return &zoneRepository{db: db}
}

func (r *zoneRepository) Create(zone *models.ParkingZone) error {
	return r.db.Create(zone).Error
}

func (r *zoneRepository) GetAll() ([]ZoneWithAvailability, error) {
	var zones []ZoneWithAvailability

	err := r.db.Model(&models.ParkingZone{}).
		Select(zoneAvailabilitySelect).
		Order("parking_zones.id").
		Scan(&zones).Error
	if err != nil {
		return nil, err
	}

	return zones, nil
}

func (r *zoneRepository) GetByID(id uint) (*ZoneWithAvailability, error) {
	var zone ZoneWithAvailability

	err := r.db.Model(&models.ParkingZone{}).
		Select(zoneAvailabilitySelect).
		Where("parking_zones.id = ?", id).
		Scan(&zone).Error
	if err != nil {
		return nil, err
	}

	if zone.ID == 0 {
		return nil, ErrZoneNotFound
	}

	return &zone, nil
}
