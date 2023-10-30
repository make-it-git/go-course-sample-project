package repository

import (
	"time"
)

type OrderModel struct {
	ID              string     `pg:"id"`
	CreatedAt       time.Time  `pg:"created_at"`
	CompletedAt     *time.Time `pg:"completed_at"`
	PickupLocation  Location   `pg:"pickup_location"`
	DropoffLocation Location   `pg:"dropoff_location"`
	TotalPrice      int        `pg:"total_price"`
	UserID          int        `pg:"user_id"`
	IdempotencyKey  string     `pg:"idempotency_key"`
}

type Location struct {
	Latitude  float32 `json:"latitude"`
	Longitude float32 `json:"longitude"`
}
