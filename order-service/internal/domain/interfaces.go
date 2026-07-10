package domain

import (
	"context"

	"github.com/google/uuid"
)

type EventPublisher interface {
	Publish(ctx context.Context, topic string, payload interface{}) error
}

type PaymentClient interface {
	ProcessPayment(ctx context.Context, orderID uuid.UUID, amount float64) error
}
