package order

import (
	"context"

	"driver-service/internal/db/repository"
)

type OrderRepository interface {
	Create(ctx context.Context, order *repository.OrderModel) error
	GetByID(ctx context.Context, id string) (*repository.OrderModel, error)
	AssignDriver(ctx context.Context, id string, driverID int64) (bool, error)
	UpdateCurrentLocation(ctx context.Context, id string, l repository.Location) error
}

type DriverSearch interface {
	FindDriver(ctx context.Context, lat1 float32, lng1 float32, lat2 float32, lng2 float32) (int64, error)
}
