package models

import "time"

const (
	ReservationActive    = "active"
	ReservationCompleted = "completed"
	ReservationCancelled = "cancelled"
)

type Reservation struct {
	ID           uint      `json:"id" gorm:"primaryKey"`
	UserID       uint      `json:"user_id" gorm:"not null;index"`
	ZoneID       uint      `json:"zone_id" gorm:"not null;index"`
	LicensePlate string    `json:"license_plate" gorm:"type:varchar(15);not null"`
	Status       string    `json:"status" gorm:"type:varchar(20);not null;default:active;check:status IN ('active','completed','cancelled')"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`

	User User        `json:"user,omitempty" gorm:"foreignKey:UserID"`
	Zone ParkingZone `json:"zone,omitempty" gorm:"foreignKey:ZoneID"`
}
