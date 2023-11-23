package order

import (
	"context"

	"github.com/google/uuid"

	"rider-service/internal/db/repository"
	"rider-service/internal/now_time"
	otel2 "rider-service/internal/otel"
	"rider-service/internal/services/driver_sender"
)

type OrderService struct {
	repo                OrderRepository
	now                 now_time.NowType
	priceEstimator      RidePriceEstimator
	driverSenderService DriverSenderService
}

func NewOrderService(repo OrderRepository, priceEstimator RidePriceEstimator, now now_time.NowType, driverSenderService DriverSenderService) OrderService {
	return OrderService{repo: repo, priceEstimator: priceEstimator, now: now, driverSenderService: driverSenderService}
}

func (s OrderService) Create(ctx context.Context, orderCreate OrderCreate) (*OrderModel, error) {
	ctx, span := otel2.GetTracer().Start(ctx, "serviceCreateOrder")
	defer span.End()

	price, err := s.priceEstimator.Estimate(
		ctx,
		orderCreate.PickupLocation.Latitude,
		orderCreate.PickupLocation.Longitude,
		orderCreate.DropoffLocation.Latitude,
		orderCreate.DropoffLocation.Longitude,
	)
	if err != nil {
		return nil, err
	}

	pickupLocation := repository.Location{
		Latitude:  orderCreate.PickupLocation.Latitude,
		Longitude: orderCreate.PickupLocation.Longitude,
	}
	dropoffLocation := repository.Location{
		Latitude:  orderCreate.DropoffLocation.Latitude,
		Longitude: orderCreate.DropoffLocation.Longitude,
	}
	now := s.now()
	order := repository.OrderModel{
		CreatedAt:       now,
		PickupLocation:  pickupLocation,
		ID:              uuid.New().String(),
		DropoffLocation: dropoffLocation,
		TotalPrice:      price,
		IdempotencyKey:  orderCreate.IdempotencyKey,
		UserID:          orderCreate.UserID,
	}
	id, err := s.repo.CreateAndGetID(ctx, &order)
	if err != nil {
		return nil, err
	}

	err = s.driverSenderService.SendToDriver(ctx, driver_sender.Order{
		ID:        id,
		CreatedAt: now,
		PickupLocation: driver_sender.Location{
			Latitude:  orderCreate.PickupLocation.Latitude,
			Longitude: orderCreate.PickupLocation.Longitude,
		},
		DropoffLocation: driver_sender.Location{
			Latitude:  orderCreate.DropoffLocation.Latitude,
			Longitude: orderCreate.DropoffLocation.Longitude,
		},
		UserID: int64(orderCreate.UserID),
	})
	if err != nil {
		return nil, err
	}

	return &OrderModel{
		ID:        id,
		CreatedAt: order.CreatedAt,
		PickupLocation: Location{
			Latitude:  orderCreate.PickupLocation.Latitude,
			Longitude: orderCreate.PickupLocation.Longitude,
		},
		DropoffLocation: Location{
			Latitude:  orderCreate.DropoffLocation.Latitude,
			Longitude: orderCreate.DropoffLocation.Longitude,
		},
		TotalPrice: price,
	}, nil
}

func (s OrderService) List(ctx context.Context, userID int) ([]OrderModel, error) {
	list, err := s.repo.List(ctx, userID)
	if err != nil {
		return nil, err
	}

	result := make([]OrderModel, 0, len(list))
	for _, l := range list {
		result = append(result, OrderModel{
			ID:          l.ID,
			CreatedAt:   l.CreatedAt,
			CompletedAt: l.CompletedAt,
			PickupLocation: Location{
				Latitude:  l.PickupLocation.Latitude,
				Longitude: l.PickupLocation.Longitude,
			},
			DropoffLocation: Location{
				Latitude:  l.DropoffLocation.Latitude,
				Longitude: l.DropoffLocation.Longitude,
			},
			TotalPrice: l.TotalPrice,
		})
	}

	return result, nil
}
