package driver_sender

import (
	"context"
	"time"

	"google.golang.org/protobuf/types/known/timestamppb"

	driver_order "rider-service/internal/generated/proto/driver.order"
)

type DriverSenderService struct {
	orderClient driver_order.OrderClient
}

func NewDriverSenderService(orderClient driver_order.OrderClient) DriverSenderService {
	return DriverSenderService{
		orderClient: orderClient,
	}
}

type Location struct {
	Latitude  float32
	Longitude float32
}

type Order struct {
	ID              string
	CreatedAt       time.Time
	PickupLocation  Location
	DropoffLocation Location
	UserID          int64
}

func (d DriverSenderService) SendToDriver(ctx context.Context, order Order) error {
	_, err := d.orderClient.StartOrder(ctx, &driver_order.StartOrderRequest{
		Id:        order.ID,
		CreatedAt: timestamppb.New(order.CreatedAt),
		PointA: &driver_order.Location{
			Latitude:  order.PickupLocation.Latitude,
			Longitude: order.PickupLocation.Longitude,
		},
		PointB: &driver_order.Location{
			Latitude:  order.DropoffLocation.Latitude,
			Longitude: order.DropoffLocation.Longitude,
		},
		UserID: order.UserID,
	})
	return err
}
