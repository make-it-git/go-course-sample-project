package repository

import (
	"time"
)

type OrderModel struct {
	ID                 string     `pg:"id"`
	CompletedAt        *time.Time `pg:"completed_at"`
	PickupLocation     Location   `pg:"pickup_location"`
	DropoffLocation    Location   `pg:"dropoff_location"`
	LastActiveLocation Location   `pg:"last_active_location"`
	UserID             int64      `pg:"user_id"`
	DriverID           *int64     `pg:"user_id"`
}

type Location struct {
	Latitude  float32 `json:"latitude"`
	Longitude float32 `json:"longitude"`
}
