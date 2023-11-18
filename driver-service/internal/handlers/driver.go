package handlers

import (
	"context"
	"time"

	driver_order "driver-service/internal/generated/proto/driver.order"
	"driver-service/internal/services/order"
)

type Handler struct {
	driver_order.UnimplementedOrderServer
	orderService OrderService
}

func NewHandler(orderService OrderService) *Handler {
	return &Handler{
		orderService: orderService,
	}
}

func (h Handler) StartOrder(ctx context.Context, req *driver_order.StartOrderRequest) (*driver_order.StartOrderResponse, error) {
	if time.Now().Unix()%2 == 0 {
		time.Sleep(time.Millisecond * 500)
	}
	driverOrder, err := h.orderService.Create(ctx, order.OrderCreate{
		ID:        req.Id,
		CreatedAt: req.CreatedAt.AsTime(),
		PickupLocation: order.Location{
			Latitude:  req.PointA.Latitude,
			Longitude: req.PointA.Longitude,
		},
		DropoffLocation: order.Location{
			Latitude:  req.PointB.Latitude,
			Longitude: req.PointB.Longitude,
		},
		UserID: req.UserID,
	})
	if err != nil {
		return nil, err
	}
	return &driver_order.StartOrderResponse{DriverId: &driverOrder.DriverID}, nil
}
