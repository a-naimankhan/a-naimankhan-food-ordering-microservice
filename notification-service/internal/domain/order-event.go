package domain

import (
	"time"

	"github.com/google/uuid"
)

type OrderEvent struct {
	OrderID         uuid.UUID `json:"id"`
	CustomerID      uuid.UUID `json:"customer_id"`
	Status          string    `json:"status"`
	Amount          float64   `json:"amount"`
	CreatedAt       time.Time `json:"created_at"`
	CancelledReason *string   `json:"cancelled_reason,omitempty"`
}
