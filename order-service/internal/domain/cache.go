package domain

import (
	"context"
	"errors"

	"github.com/google/uuid"
)

// ErrCacheMiss is returned by OrderCache.Get when the order is not cached.
var ErrCacheMiss = errors.New("cache miss")

// OrderCache stores orders by ID (cache-aside for GET /orders/:id).
type OrderCache interface {
	Get(ctx context.Context, id uuid.UUID) (*Order, error)
	Set(ctx context.Context, order *Order) error
	Delete(ctx context.Context, id uuid.UUID) error
}
