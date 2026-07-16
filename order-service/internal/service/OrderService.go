package service

import (
	"context"
	"log"
	"order-service/internal/domain"
	"time"

	"github.com/google/uuid"
)

type orderService struct {
	orderRepo      domain.OrderRepository
	eventPublisher domain.EventPublisher
}

func NewOrderService(orderRepo domain.OrderRepository, eventPublisher domain.EventPublisher) domain.OrderService {
	return &orderService{
		orderRepo:      orderRepo,
		eventPublisher: eventPublisher,
	}
}

func (s *orderService) GetOrder(ctx context.Context, id uuid.UUID) (*domain.Order, error) {
	order, err := s.orderRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if order == nil {
		return nil, domain.ErrOrderNotFound
	}
	return order, nil
}

func (s *orderService) CreateOrder(ctx context.Context, order *domain.Order) (*domain.Order, error) {
	if order == nil {
		return nil, domain.ErrInvalidCustomerID
	}

	if order.Amount <= 0 {
		return nil, domain.ErrInvalidAmount
	}

	if order.Status == "" {
		return nil, domain.ErrStatusEmpty
	}

	if !domain.IsValidStatus(order.Status) {
		return nil, domain.ErrInvalidStatus
	}

	if order.CustomerID == uuid.Nil {
		return nil, domain.ErrInvalidCustomerID
	}

	if order.ID == uuid.Nil {
		order.ID = uuid.New()
	}
	if order.CreatedAt.IsZero() {
		order.CreatedAt = time.Now()
	}

	if err := s.orderRepo.Create(ctx, order); err != nil {
		return nil, err
	}

	// Publish event asynchronously (non-blocking)
	if s.eventPublisher != nil {
		if err := s.eventPublisher.Publish(ctx, "order.created", order); err != nil {
			// Log the error but don't fail the request - event publishing is not critical path
			log.Printf("failed to publish order.created event: %v", err)
		}
	}

	return order, nil
}

func (s *orderService) UpdateOrderStatus(ctx context.Context, id uuid.UUID, status string) error {
	if status == "" {
		return domain.ErrStatusEmpty
	}
	if !domain.IsValidStatus(status) {
		return domain.ErrInvalidStatus
	}
	return s.orderRepo.UpdateStatus(ctx, id, status)
}
