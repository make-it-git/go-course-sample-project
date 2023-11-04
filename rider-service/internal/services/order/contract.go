package order

import (
	"context"

	"rider-service/internal/db/repository"
	"rider-service/internal/services/driver_sender"
)

type OrderRepository interface {
	List(ctx context.Context, userID int) ([]repository.OrderModel, error)
	CreateAndGetID(ctx context.Context, order *repository.OrderModel) (string, error)
}

type RidePriceEstimator interface {
	Estimate(ctx context.Context, lat1 float32, lng1 float32, lat2 float32, lng2 float32) (int, error)
}

type DriverSenderService interface {
	SendToDriver(ctx context.Context, order driver_sender.Order) error
}
