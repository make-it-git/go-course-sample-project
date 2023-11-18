package order

import (
	"context"
	"strconv"
	"time"

	"github.com/prometheus/client_golang/prometheus"

	"driver-service/internal/db/repository"
	"driver-service/internal/logger"
)

var orderCounter *prometheus.CounterVec

func init() {
	orderCounter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "order_create_count",
			Help: "No of orders created",
		},
		[]string{"assignAttempts"},
	)
	prometheus.MustRegister(orderCounter)
}

type OrderService struct {
	repo         OrderRepository
	driverSearch DriverSearch
	log          logger.Log
}

func NewOrderService(repo OrderRepository, driverSearch DriverSearch, log logger.Log) OrderService {
	return OrderService{repo: repo, driverSearch: driverSearch, log: log}
}

func (s OrderService) Create(ctx context.Context, orderCreate OrderCreate) (*OrderModel, error) {
	pickupLocation := repository.Location{
		Latitude:  orderCreate.PickupLocation.Latitude,
		Longitude: orderCreate.PickupLocation.Longitude,
	}
	dropoffLocation := repository.Location{
		Latitude:  orderCreate.DropoffLocation.Latitude,
		Longitude: orderCreate.DropoffLocation.Longitude,
	}
	order := repository.OrderModel{
		ID:              orderCreate.ID,
		PickupLocation:  pickupLocation,
		DropoffLocation: dropoffLocation,
		UserID:          orderCreate.UserID,
	}

	driverID, err := s.driverSearch.FindDriver(
		ctx,
		orderCreate.PickupLocation.Latitude,
		orderCreate.PickupLocation.Longitude,
		orderCreate.DropoffLocation.Latitude,
		orderCreate.DropoffLocation.Longitude,
	)
	if err == nil {
		order.DriverID = &driverID
		orderCounter.With(prometheus.Labels{"assignAttempts": "0"}).Inc()
	}

	err = s.repo.Create(ctx, &order)
	if err != nil {
		return nil, err
	}

	if order.DriverID == nil {
		go func() {
			attempts := 0
			for attempts < 10 {
				attempts++
				orderCounter.With(prometheus.Labels{"assignAttempts": strconv.Itoa(attempts)}).Inc()
				time.Sleep(time.Second * 1)
				driverID, err := s.driverSearch.FindDriver(
					ctx,
					orderCreate.PickupLocation.Latitude,
					orderCreate.PickupLocation.Longitude,
					orderCreate.DropoffLocation.Latitude,
					orderCreate.DropoffLocation.Longitude,
				)
				if err != nil {
					s.log.WithError(err, "find driver")
					continue
				}
				ok, err := s.repo.AssignDriver(context.Background(), orderCreate.ID, driverID)
				if err != nil {
					s.log.WithError(err, "assign driver", "orderID", orderCreate.ID)
					continue
				}
				if !ok {
					s.log.Warning("already assigned driver", "orderID", orderCreate.ID)
					return
				}
				s.log.Info("successfully assigned driver", "orderID", orderCreate.ID)
			}
		}()
	}

	return &OrderModel{
		ID:        orderCreate.ID,
		CreatedAt: orderCreate.CreatedAt,
		PickupLocation: Location{
			Latitude:  orderCreate.PickupLocation.Latitude,
			Longitude: orderCreate.PickupLocation.Longitude,
		},
		DropoffLocation: Location{
			Latitude:  orderCreate.DropoffLocation.Latitude,
			Longitude: orderCreate.DropoffLocation.Longitude,
		},
		DriverID: driverID,
	}, nil
}

func (s OrderService) UpdateLocation(ctx context.Context, id string, l Location) error {
	return s.repo.UpdateCurrentLocation(ctx, id, repository.Location{
		Latitude:  l.Latitude,
		Longitude: l.Longitude,
	})
}
