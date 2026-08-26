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
	o.log.Debug("Creating order | ID=%s | Customer=%s | Amount=%.2f | Status=%s", 
		order.ID, order.CustomerID, order.Amount, order.Status)
	
	query := "INSERT INTO ORDERS (id, customer_id, status, amount, created_at) VALUES ($1, $2, $3, $4, $5)"
	result, err := o.db.ExecContext(ctx, query, order.ID, order.CustomerID, order.Status, order.Amount, order.CreatedAt)
	
	if err != nil {
		o.log.Error("Failed to create order | ID=%s | Error: %v", order.ID, err)
		return err
	}
	
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		o.log.Error("Failed to get rows affected | ID=%s | Error: %v", order.ID, err)
		return err
	}
	
	o.log.Info("✅ Order created successfully | ID=%s | Rows affected: %d", order.ID, rowsAffected)
	return nil
}

func (o *orderRepo) GetByID(ctx context.Context, id uuid.UUID) (*domain.Order, error) {
	o.log.Debug("Fetching order from database | ID=%s", id)
	
	query := "SELECT * FROM ORDERS WHERE id = $1"
	var order domain.Order
	err := o.db.GetContext(ctx, &order, query, id)
	
	if err != nil {
		o.log.Error("Failed to fetch order | ID=%s | Error: %v", id, err)
		return nil, err
	}
	
	o.log.Debug("✅ Order fetched successfully | ID=%s | Customer=%s | Amount=%.2f", id, order.CustomerID, order.Amount)
	return &order, nil
}

func (o *orderRepo) UpdateStatus(ctx context.Context, id uuid.UUID, status string) error {
	o.log.Debug("Updating order status | ID=%s | New Status=%s", id, status)
	
	query := "UPDATE ORDERS SET status = $2 WHERE id = $1"
	result, err := o.db.ExecContext(ctx, query, id, status)
	
	if err != nil {
		o.log.Error("Failed to update order status | ID=%s | Error: %v", id, err)
		return err
	}
	
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		o.log.Error("Failed to get rows affected | ID=%s | Error: %v", id, err)
		return err
	}
	
	if rowsAffected == 0 {
		o.log.Error("Order not found for update | ID=%s", id)
		return domain.ErrOrderNotFound
	}
	
	o.log.Info("✅ Order status updated | ID=%s | New Status=%s", id, status)
	return nil
}
