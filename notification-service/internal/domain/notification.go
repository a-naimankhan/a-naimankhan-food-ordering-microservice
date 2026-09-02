package domain

import (
	"time"

	"github.com/google/uuid"
)

type Notification struct {
	ID         uuid.UUID
	OrderID    uuid.UUID
	CustomerID uuid.UUID
	Type       string
	Message    string
	CreatedAt  time.Time
}
