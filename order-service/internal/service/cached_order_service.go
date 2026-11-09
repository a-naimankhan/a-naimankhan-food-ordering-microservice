package service

import (
	"context"
	"errors"

	"order-service/internal/domain"
	"order-service/internal/logger"

	"github.com/google/uuid"
)

// cachedOrderService wraps domain.OrderService with a cache-aside layer for GetOrder.
// Writes that change an order invalidate its cache entry. Cache errors never fail
// a request: on any cache problem the call falls through to the wrapped service.
type cachedOrderService struct {
	next  domain.OrderService
	cache domain.OrderCache
	log   *logger.Logger
}

func NewCachedOrderService(next domain.OrderService, cache domain.OrderCache) domain.OrderService {
	return &cachedOrderService{
		next:  next,
		cache: cache,
		log:   logger.GetLogger(),
	}
}

func (s *cachedOrderService) GetOrder(ctx context.Context, id uuid.UUID) (*domain.Order, error) {
	order, err := s.cache.Get(ctx, id)
	switch {
	case err == nil:
		s.log.Debug("Order cache hit", "id", id)
		return order, nil
	case errors.Is(err, domain.ErrCacheMiss):
		s.log.Debug("Order cache miss", "id", id)
	default:
		s.log.Warn("Order cache read failed, falling back to database", "id", id, "error", err.Error())
	}

	order, err = s.next.GetOrder(ctx, id)
	if err != nil {
		return nil, err
	}

	if err := s.cache.Set(ctx, order); err != nil {
		s.log.Warn("Failed to cache order", "id", id, "error", err.Error())
	}
	return order, nil
}

func (s *cachedOrderService) CreateOrder(ctx context.Context, order *domain.Order) (*domain.Order, error) {
	return s.next.CreateOrder(ctx, order)
}

func (s *cachedOrderService) UpdateOrderStatus(ctx context.Context, id uuid.UUID, status string) error {
	if err := s.next.UpdateOrderStatus(ctx, id, status); err != nil {
		return err
	}
	s.invalidate(ctx, id)
	return nil
}

func (s *cachedOrderService) CancelOrder(ctx context.Context, id uuid.UUID, reason string) error {
	if err := s.next.CancelOrder(ctx, id, reason); err != nil {
		return err
	}
	s.invalidate(ctx, id)
	return nil
}

func (s *cachedOrderService) invalidate(ctx context.Context, id uuid.UUID) {
	if err := s.cache.Delete(ctx, id); err != nil {
		s.log.Warn("Failed to invalidate order cache", "id", id, "error", err.Error())
	}
}
