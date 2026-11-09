package service

import (
	"context"
	"errors"
	"testing"

	"order-service/internal/domain"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type fakeCache struct {
	items     map[uuid.UUID]*domain.Order
	getErr    error
	setErr    error
	deleteErr error
	deleted   []uuid.UUID
}

func newFakeCache() *fakeCache {
	return &fakeCache{items: map[uuid.UUID]*domain.Order{}}
}

func (c *fakeCache) Get(_ context.Context, id uuid.UUID) (*domain.Order, error) {
	if c.getErr != nil {
		return nil, c.getErr
	}
	o, ok := c.items[id]
	if !ok {
		return nil, domain.ErrCacheMiss
	}
	return o, nil
}

func (c *fakeCache) Set(_ context.Context, o *domain.Order) error {
	if c.setErr != nil {
		return c.setErr
	}
	c.items[o.ID] = o
	return nil
}

func (c *fakeCache) Delete(_ context.Context, id uuid.UUID) error {
	c.deleted = append(c.deleted, id)
	if c.deleteErr != nil {
		return c.deleteErr
	}
	delete(c.items, id)
	return nil
}

type fakeOrderService struct {
	getCalls  int
	getFn     func(id uuid.UUID) (*domain.Order, error)
	updateErr error
	cancelErr error
}

func (f *fakeOrderService) CreateOrder(_ context.Context, o *domain.Order) (*domain.Order, error) {
	return o, nil
}

func (f *fakeOrderService) GetOrder(_ context.Context, id uuid.UUID) (*domain.Order, error) {
	f.getCalls++
	return f.getFn(id)
}

func (f *fakeOrderService) UpdateOrderStatus(context.Context, uuid.UUID, string) error {
	return f.updateErr
}

func (f *fakeOrderService) CancelOrder(context.Context, uuid.UUID, string) error {
	return f.cancelErr
}

func orderGetter(order *domain.Order) func(uuid.UUID) (*domain.Order, error) {
	return func(uuid.UUID) (*domain.Order, error) { return order, nil }
}

func TestCachedGetOrder_MissThenHit(t *testing.T) {
	order := &domain.Order{ID: uuid.New(), Status: domain.StatusPending}
	inner := &fakeOrderService{getFn: orderGetter(order)}
	cache := newFakeCache()
	svc := NewCachedOrderService(inner, cache)

	got, err := svc.GetOrder(context.Background(), order.ID)
	require.NoError(t, err)
	assert.Equal(t, order, got)
	assert.Equal(t, 1, inner.getCalls)
	assert.Contains(t, cache.items, order.ID)

	got, err = svc.GetOrder(context.Background(), order.ID)
	require.NoError(t, err)
	assert.Equal(t, order, got)
	assert.Equal(t, 1, inner.getCalls, "second call must be served from cache")
}

func TestCachedGetOrder_NotFoundIsNotCached(t *testing.T) {
	inner := &fakeOrderService{getFn: func(uuid.UUID) (*domain.Order, error) {
		return nil, domain.ErrOrderNotFound
	}}
	cache := newFakeCache()
	svc := NewCachedOrderService(inner, cache)

	_, err := svc.GetOrder(context.Background(), uuid.New())
	assert.ErrorIs(t, err, domain.ErrOrderNotFound)
	assert.Empty(t, cache.items)
}

func TestCachedGetOrder_CacheErrorsFallBackToService(t *testing.T) {
	order := &domain.Order{ID: uuid.New()}
	inner := &fakeOrderService{getFn: orderGetter(order)}
	cache := newFakeCache()
	cache.getErr = errors.New("redis down")
	cache.setErr = errors.New("redis down")
	svc := NewCachedOrderService(inner, cache)

	got, err := svc.GetOrder(context.Background(), order.ID)
	require.NoError(t, err)
	assert.Equal(t, order, got)
	assert.Equal(t, 1, inner.getCalls)
}

func TestCachedUpdateOrderStatus_Invalidates(t *testing.T) {
	id := uuid.New()
	cache := newFakeCache()
	cache.items[id] = &domain.Order{ID: id, Status: domain.StatusPending}
	svc := NewCachedOrderService(&fakeOrderService{}, cache)

	require.NoError(t, svc.UpdateOrderStatus(context.Background(), id, domain.StatusAccepted))
	assert.NotContains(t, cache.items, id)
}

func TestCachedUpdateOrderStatus_ErrorKeepsCache(t *testing.T) {
	id := uuid.New()
	cache := newFakeCache()
	cache.items[id] = &domain.Order{ID: id}
	svc := NewCachedOrderService(&fakeOrderService{updateErr: domain.ErrInvalidTransaction}, cache)

	err := svc.UpdateOrderStatus(context.Background(), id, domain.StatusDelivered)
	assert.ErrorIs(t, err, domain.ErrInvalidTransaction)
	assert.Empty(t, cache.deleted)
}

func TestCachedCancelOrder_Invalidates(t *testing.T) {
	id := uuid.New()
	cache := newFakeCache()
	cache.items[id] = &domain.Order{ID: id, Status: domain.StatusPending}
	svc := NewCachedOrderService(&fakeOrderService{}, cache)

	require.NoError(t, svc.CancelOrder(context.Background(), id, "changed my mind"))
	assert.NotContains(t, cache.items, id)
}

func TestCachedCancelOrder_InvalidateErrorIsIgnored(t *testing.T) {
	cache := newFakeCache()
	cache.deleteErr = errors.New("redis down")
	svc := NewCachedOrderService(&fakeOrderService{}, cache)

	assert.NoError(t, svc.CancelOrder(context.Background(), uuid.New(), "reason"))
}
