package models

import "time"

const (
	ZoneTypeGeneral    = "general"
	ZoneTypeEVCharging = "ev_charging"
	ZoneTypeCovered    = "covered"
)

type ParkingZone struct {
	ID            uint      `json:"id" gorm:"primaryKey"`
	Name          string    `json:"name" gorm:"type:varchar(150);not null"`
	Type          string    `json:"type" gorm:"type:varchar(20);not null;check:type IN ('general','ev_charging','covered')"`
	TotalCapacity int       `json:"total_capacity" gorm:"not null"`
	PricePerHour  float64   `json:"price_per_hour" gorm:"type:decimal(10,2);not null"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}
