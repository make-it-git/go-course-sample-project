package repository

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type OrderRepository interface {
	Create(ctx context.Context, order *OrderModel) error
	GetByID(ctx context.Context, id string) (*OrderModel, error)
	AssignDriver(ctx context.Context, id string, driverID int64) (bool, error)
	UpdateCurrentLocation(ctx context.Context, id string, l Location) error
}

func NewOrderRepository(conn *pgxpool.Pool) OrderRepository {
	return &OrderRepositoryImpl{conn: conn}
}

type OrderRepositoryImpl struct {
	conn *pgxpool.Pool
}

func (r *OrderRepositoryImpl) Create(ctx context.Context, order *OrderModel) error {
	query := `INSERT INTO orders 
    	(id, pickup_location, dropoff_location, user_id, driver_id) VALUES ($1, $2, $3, $4, $5)
    	ON CONFLICT (id) DO NOTHING`
	args := make([]interface{}, 5)
	args[0] = order.ID
	args[1] = order.PickupLocation
	args[2] = order.DropoffLocation
	args[3] = order.UserID
	if order.DriverID == nil {
		args[4] = nil
	} else {
		args[4] = *order.DriverID
	}
	_, err := r.conn.Exec(ctx, query, args...)
	return err
}

func (r *OrderRepositoryImpl) GetByID(ctx context.Context, id string) (*OrderModel, error) {
	var order OrderModel
	sql := `SELECT id, pickup_location, dropoff_location, last_active_location, completed_at, user_id, driver_id FROM orders WHERE id = $1`
	err := r.conn.QueryRow(ctx, sql, id).Scan(&order)
	if err != nil {
		return nil, err
	}

	return &order, nil
}

func (r *OrderRepositoryImpl) UpdateCurrentLocation(ctx context.Context, id string, l Location) error {
	query := `UPDATE orders SET last_active_location = $1 WHERE id = $2`
	_, err := r.conn.Exec(ctx, query, l, id)
	return err
}

func (r *OrderRepositoryImpl) AssignDriver(ctx context.Context, id string, driverID int64) (bool, error) {
	query := `UPDATE orders SET driver_id = $1 WHERE id = $2 AND driver_id IS NULL`
	cmd, err := r.conn.Exec(ctx, query, driverID, id)
	if err != nil {
		return false, err
	}
	return cmd.RowsAffected() > 0, nil
}
