package redis

import (
	"context"
	"testing"
	"time"

	"order-service/internal/domain"

	"github.com/alicebob/miniredis/v2"
	"github.com/google/uuid"
	goredis "github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func newTestCache(t *testing.T, ttl time.Duration) (*OrderCache, *miniredis.Miniredis) {
	t.Helper()
	mr := miniredis.RunT(t)
	client := goredis.NewClient(&goredis.Options{Addr: mr.Addr()})
	t.Cleanup(func() { _ = client.Close() })
	return NewOrderCache(client, ttl), mr
}

func TestOrderCache_SetGetDelete(t *testing.T) {
	cache, mr := newTestCache(t, time.Minute)
	ctx := context.Background()

	order := &domain.Order{
		ID:         uuid.New(),
		CustomerID: uuid.New(),
		Status:     domain.StatusPending,
		Amount:     42.5,
		CreatedAt:  time.Now().UTC().Truncate(time.Second),
	}

	_, err := cache.Get(ctx, order.ID)
	assert.ErrorIs(t, err, domain.ErrCacheMiss)

	require.NoError(t, cache.Set(ctx, order))
	assert.Equal(t, time.Minute, mr.TTL("order:"+order.ID.String()))

	got, err := cache.Get(ctx, order.ID)
	require.NoError(t, err)
	assert.Equal(t, order.ID, got.ID)
	assert.Equal(t, order.CustomerID, got.CustomerID)
	assert.Equal(t, order.Status, got.Status)
	assert.Equal(t, order.Amount, got.Amount)
	assert.True(t, order.CreatedAt.Equal(got.CreatedAt))

	require.NoError(t, cache.Delete(ctx, order.ID))
	_, err = cache.Get(ctx, order.ID)
	assert.ErrorIs(t, err, domain.ErrCacheMiss)
}

func TestOrderCache_Expires(t *testing.T) {
	cache, mr := newTestCache(t, time.Minute)
	ctx := context.Background()
	order := &domain.Order{ID: uuid.New(), Status: domain.StatusPending}

	require.NoError(t, cache.Set(ctx, order))
	mr.FastForward(2 * time.Minute)

	_, err := cache.Get(ctx, order.ID)
	assert.ErrorIs(t, err, domain.ErrCacheMiss)
}

func TestOrderCache_RedisDown(t *testing.T) {
	cache, mr := newTestCache(t, time.Minute)
	mr.Close()

	_, err := cache.Get(context.Background(), uuid.New())
	assert.Error(t, err)
	assert.NotErrorIs(t, err, domain.ErrCacheMiss)
}
