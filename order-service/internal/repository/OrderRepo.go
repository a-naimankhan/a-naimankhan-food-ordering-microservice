package repository

import (
	"context"
	"order-service/internal/domain"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type orderRepo struct {
	db *sqlx.DB
}

func NewOrderRepo(db *sqlx.DB) domain.OrderRepository {
	return &orderRepo{db}
}

func (o *orderRepo) Create(ctx context.Context, order *domain.Order) error {
	query := "INSERT INTO ORDERS (id, customer_id, status, amount, created_at) VALUES ($1, $2, $3, $4, $5)"
	_, err := o.db.ExecContext(ctx, query, order.ID, order.CustomerID, order.Status, order.Amount, order.CreatedAt)
	return err
}

func (o *orderRepo) GetByID(ctx context.Context, id uuid.UUID) (*domain.Order, error) {
	query := "SELECT * FROM ORDERS WHERE id = $1"
	var order domain.Order
	err := o.db.GetContext(ctx, &order, query, id)
	if err != nil {
		return nil, err
	}
	return &order, nil
}

func (o *orderRepo) UpdateStatus(ctx context.Context, id uuid.UUID, status string) error {
	query := "UPDATE ORDERS SET status = $2 WHERE id = $1"
	_, err := o.db.ExecContext(ctx, query, id, status)
	return err
}
