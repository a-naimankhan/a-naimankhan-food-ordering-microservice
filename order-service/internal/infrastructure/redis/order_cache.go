package redis

import (
	"context"
	"errors"
	"time"

	"order-service/internal/domain"

	"github.com/goccy/go-json"
	"github.com/google/uuid"
	goredis "github.com/redis/go-redis/v9"
)

const keyPrefix = "order:"

// OrderCache is a Redis-backed implementation of domain.OrderCache.
type OrderCache struct {
	client *goredis.Client
	ttl    time.Duration
}

func NewOrderCache(client *goredis.Client, ttl time.Duration) *OrderCache {
	return &OrderCache{client: client, ttl: ttl}
}

func key(id uuid.UUID) string {
	return keyPrefix + id.String()
}

func (c *OrderCache) Get(ctx context.Context, id uuid.UUID) (*domain.Order, error) {
	data, err := c.client.Get(ctx, key(id)).Bytes()
	if errors.Is(err, goredis.Nil) {
		return nil, domain.ErrCacheMiss
	}
	if err != nil {
		return nil, err
	}

	var order domain.Order
	if err := json.Unmarshal(data, &order); err != nil {
		return nil, err
	}
	return &order, nil
}

func (c *OrderCache) Set(ctx context.Context, order *domain.Order) error {
	data, err := json.Marshal(order)
	if err != nil {
		return err
	}
	return c.client.Set(ctx, key(order.ID), data, c.ttl).Err()
}

func (c *OrderCache) Delete(ctx context.Context, id uuid.UUID) error {
	return c.client.Del(ctx, key(id)).Err()
}
