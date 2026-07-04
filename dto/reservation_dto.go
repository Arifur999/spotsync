package dto

import "time"

type CreateReservationRequest struct {
	ZoneID       uint   `json:"zone_id" validate:"required"`
	LicensePlate string `json:"license_plate" validate:"required,max=15"`
}

// ReservationResponse mirrors the create-reservation response shape.
type ReservationResponse struct {
	ID           uint      `json:"id"`
	UserID       uint      `json:"user_id"`
	ZoneID       uint      `json:"zone_id"`
	LicensePlate string    `json:"license_plate"`
	Status       string    `json:"status"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// MyReservationResponse is used by GET /reservations/my-reservations,
// with the zone preloaded and nested instead of raw foreign keys.
type MyReservationResponse struct {
	ID           uint        `json:"id"`
	LicensePlate string      `json:"license_plate"`
	Status       string      `json:"status"`
	Zone         ZoneSummary `json:"zone"`
	CreatedAt    time.Time   `json:"created_at"`
}

type ReservationUserSummary struct {
	ID    uint   `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

// AdminReservationResponse is used by GET /reservations (admin), with both
// the user and zone preloaded and nested.
type AdminReservationResponse struct {
	ID           uint                   `json:"id"`
	LicensePlate string                 `json:"license_plate"`
	Status       string                 `json:"status"`
	User         ReservationUserSummary `json:"user"`
	Zone         ZoneSummary            `json:"zone"`
	CreatedAt    time.Time              `json:"created_at"`
	UpdatedAt    time.Time              `json:"updated_at"`
}
