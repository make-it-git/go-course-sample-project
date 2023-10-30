package order

import (
	"time"
)

type OrderCreate struct {
	PickupLocation  Location
	DropoffLocation Location
	UserID          int
	IdempotencyKey  string
}

type OrderModel struct {
	ID              string
	CreatedAt       time.Time
	CompletedAt     *time.Time
	PickupLocation  Location
	DropoffLocation Location
	TotalPrice      int
}

type Location struct {
	Latitude  float32
	Longitude float32
}
