package repository

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type OrderRepository interface {
	GetByID(ctx context.Context, id string) (*OrderModel, error)
	List(ctx context.Context, userID int) ([]OrderModel, error)
	CreateAndGetID(ctx context.Context, order *OrderModel) (string, error)
}

func NewOrderRepository(conn *pgxpool.Pool) OrderRepository {
	return &OrderRepositoryImpl{conn: conn}
}

type OrderRepositoryImpl struct {
	conn *pgxpool.Pool
}

func (r *OrderRepositoryImpl) GetByID(ctx context.Context, id string) (*OrderModel, error) {
	var order OrderModel
	sql := `SELECT id, created_at, completed_at, pickup_location, dropoff_location, total_price FROM orders WHERE id = $1`
	err := r.conn.QueryRow(ctx, sql, id).Scan(&order)
	if err != nil {
		return nil, err
	}

	return &OrderModel{
		ID:              order.ID,
		CreatedAt:       order.CreatedAt,
		CompletedAt:     order.CompletedAt,
		PickupLocation:  order.PickupLocation,
		DropoffLocation: order.DropoffLocation,
		TotalPrice:      order.TotalPrice,
	}, nil
}

func (r *OrderRepositoryImpl) List(ctx context.Context, userID int) ([]OrderModel, error) {
	sql := `SELECT id, created_at, completed_at, pickup_location, dropoff_location, total_price FROM orders WHERE user_id = $1`
	rows, err := r.conn.Query(ctx, sql, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var convertedOrders []OrderModel
	for rows.Next() {
		var order OrderModel
		err := rows.Scan(&order.ID, &order.CreatedAt, &order.CompletedAt, &order.PickupLocation, &order.DropoffLocation, &order.TotalPrice)
		if err != nil {
			return nil, err
		}
		convertedOrders = append(convertedOrders, order)
	}
	return convertedOrders, nil
}

func (r *OrderRepositoryImpl) CreateAndGetID(ctx context.Context, order *OrderModel) (string, error) {
	query := `INSERT INTO orders 
    	(id, created_at, completed_at, pickup_location, dropoff_location, total_price, user_id, idempotency_key) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
    	ON CONFLICT (idempotency_key) DO NOTHING RETURNING id`
	var id string
	err := r.conn.QueryRow(ctx, query, order.ID, order.CreatedAt, order.CompletedAt, order.PickupLocation, order.DropoffLocation, order.TotalPrice, order.UserID, order.IdempotencyKey).Scan(&id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			err = r.conn.QueryRow(ctx, "select id from orders WHERE idempotency_key = $1 AND user_id = $2", order.IdempotencyKey, order.UserID).Scan(&id)
			return id, err
		}
		return "", err
	}

	return id, nil
}
