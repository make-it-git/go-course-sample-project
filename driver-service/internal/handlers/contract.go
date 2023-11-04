package handlers

import (
	"context"

	"driver-service/internal/services/order"
)

type OrderService interface {
	Create(ctx context.Context, orderCreate order.OrderCreate) (*order.OrderModel, error)
	UpdateLocation(ctx context.Context, id string, l order.Location) error
}
