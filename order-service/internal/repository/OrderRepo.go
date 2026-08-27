package repository

import (
	"context"
	"order-service/internal/domain"
	"order-service/internal/logger"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type orderRepo struct {
	db  *sqlx.DB
	log *logger.Logger
}

func NewOrderRepo(db *sqlx.DB) domain.OrderRepository {
	return &orderRepo{
		db:  db,
		log: logger.GetLogger(),
	}
}

func (o *orderRepo) Create(ctx context.Context, order *domain.Order) error {
	o.log.Debug("Creating order", "id", order.ID, "customer_id", order.CustomerID, "amount", order.Amount, "status", order.Status)
	
	query := "INSERT INTO ORDERS (id, customer_id, status, amount, created_at) VALUES ($1, $2, $3, $4, $5)"
	result, err := o.db.ExecContext(ctx, query, order.ID, order.CustomerID, order.Status, order.Amount, order.CreatedAt)
	
	if err != nil {
		o.log.Error("Failed to create order", "id", order.ID, "error", err.Error())
		return err
	}
	
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		o.log.Error("Failed to get rows affected", "id", order.ID, "error", err.Error())
		return err
	}
	
	o.log.Info("Order created successfully", "id", order.ID, "rows_affected", rowsAffected)
	return nil
}

func (o *orderRepo) GetByID(ctx context.Context, id uuid.UUID) (*domain.Order, error) {
	o.log.Debug("Fetching order from database", "id", id)
	
	query := "SELECT * FROM ORDERS WHERE id = $1"
	var order domain.Order
	err := o.db.GetContext(ctx, &order, query, id)
	
	if err != nil {
		o.log.Error("Failed to fetch order", "id", id, "error", err.Error())
		return nil, err
	}
	
	o.log.Debug("Order fetched successfully", "id", id, "customer_id", order.CustomerID, "amount", order.Amount)
	return &order, nil
}

func (o *orderRepo) UpdateStatus(ctx context.Context, id uuid.UUID, status string) error {
	o.log.Debug("Updating order status", "id", id, "new_status", status)
	
	query := "UPDATE ORDERS SET status = $2 WHERE id = $1"
	result, err := o.db.ExecContext(ctx, query, id, status)
	
	if err != nil {
		o.log.Error("Failed to update order status", "id", id, "error", err.Error())
		return err
	}
	
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		o.log.Error("Failed to get rows affected", "id", id, "error", err.Error())
		return err
	}
	
	if rowsAffected == 0 {
		o.log.Error("Order not found for update", "id", id)
		return domain.ErrOrderNotFound
	}
	
	o.log.Info("Order status updated", "id", id, "new_status", status)
	return nil
}
