package order

import (
	"time"
)

type OrderCreate struct {
	ID              string
	CreatedAt       time.Time
	PickupLocation  Location
	DropoffLocation Location
	UserID          int64
}

type OrderModel struct {
	ID              string
	CreatedAt       time.Time
	CompletedAt     *time.Time
	PickupLocation  Location
	DropoffLocation Location
	DriverID        int64
}

type Location struct {
	Latitude  float32
	Longitude float32
}
