package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type Order struct {
	ID         uuid.UUID `json:"id" db:"id" `
	CustomerID uuid.UUID `json:"customer_id" db:"customer_id"`
	Status     string    `json:"status" db:"status"`
	Amount     float64   `json:"amount" db:"amount"`
	CreatedAt  time.Time `json:"created_at" db:"created_at"`
}

type OrderRepository interface {
	Create(ctx context.Context, order *Order) error
	GetByID(ctx context.Context, id uuid.UUID) (*Order, error)
	UpdateStatus(ctx context.Context, id uuid.UUID, status string) error
}

type OrderService interface {
	PlaceOrder(ctx context.Context, order *Order) error
	GetOrder(ctx context.Context, id uuid.UUID) (*Order, error)
}
