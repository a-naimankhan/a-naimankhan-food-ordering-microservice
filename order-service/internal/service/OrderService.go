package service

import (
	"context"
	"order-service/internal/domain"
	"order-service/internal/logger"
	"time"

	"github.com/google/uuid"
)

type orderService struct {
	orderRepo      domain.OrderRepository
	eventPublisher domain.EventPublisher
	log            *logger.Logger
}

func NewOrderService(orderRepo domain.OrderRepository, eventPublisher domain.EventPublisher) domain.OrderService {
	return &orderService{
		orderRepo:      orderRepo,
		eventPublisher: eventPublisher,
		log:            logger.GetLogger(),
	}
}

func (s *orderService) GetOrder(ctx context.Context, id uuid.UUID) (*domain.Order, error) {
	s.log.Debug("Service: Getting order | ID=%s", id)
	
	order, err := s.orderRepo.GetByID(ctx, id)
	if err != nil {
		s.log.Error("Service: Failed to get order | ID=%s | Error: %v", id, err)
		return nil, err
	}
	
	if order == nil {
		s.log.Error("Service: Order not found | ID=%s", id)
		return nil, domain.ErrOrderNotFound
	}
	
	s.log.Info("✅ Service: Order retrieved | ID=%s | Status=%s | Amount=%.2f", id, order.Status, order.Amount)
	return order, nil
}

func (s *orderService) CreateOrder(ctx context.Context, order *domain.Order) (*domain.Order, error) {
	s.log.Debug("Service: Creating order | CustomerID=%s | Amount=%.2f | Status=%s", 
		order.CustomerID, order.Amount, order.Status)
	
	if order == nil {
		s.log.Error("Service: Order is nil")
		return nil, domain.ErrInvalidCustomerID
	}

	if order.Amount <= 0 {
		s.log.Error("Service: Invalid amount | Amount=%.2f", order.Amount)
		return nil, domain.ErrInvalidAmount
	}

	if order.Status == "" {
		s.log.Error("Service: Status is empty")
		return nil, domain.ErrStatusEmpty
	}

	if !domain.IsValidStatus(order.Status) {
		s.log.Error("Service: Invalid status | Status=%s", order.Status)
		return nil, domain.ErrInvalidStatus
	}

	if order.CustomerID == uuid.Nil {
		s.log.Error("Service: Invalid customer ID")
		return nil, domain.ErrInvalidCustomerID
	}

	if order.ID == uuid.Nil {
		order.ID = uuid.New()
		s.log.Debug("Service: Generated new Order ID | ID=%s", order.ID)
	}
	
	if order.CreatedAt.IsZero() {
		order.CreatedAt = time.Now()
		s.log.Debug("Service: Set creation timestamp | CreatedAt=%s", order.CreatedAt)
	}

	s.log.Info("Service: Validations passed, saving order to repository | ID=%s", order.ID)
	if err := s.orderRepo.Create(ctx, order); err != nil {
		s.log.Error("Service: Failed to create order in repository | ID=%s | Error: %v", order.ID, err)
		return nil, err
	}

	s.log.Info("✅ Service: Order saved to database | ID=%s", order.ID)

	if s.eventPublisher != nil {
		s.log.Debug("Service: Publishing order.created event | ID=%s", order.ID)
		if err := s.eventPublisher.Publish(ctx, "order.created", order); err != nil {
			s.log.Error("Service: Failed to publish order.created event | ID=%s | Error: %v", order.ID, err)
		}
	}

	s.log.Info("✅ Service: Order created successfully | ID=%s | CustomerID=%s", order.ID, order.CustomerID)
	return order, nil
}

func (s *orderService) UpdateOrderStatus(ctx context.Context, id uuid.UUID, status string) error {
	s.log.Debug("Service: Updating order status | ID=%s | NewStatus=%s", id, status)
	
	if status == "" {
		s.log.Error("Service: Status is empty | ID=%s", id)
		return domain.ErrStatusEmpty
	}
	
	if !domain.IsValidStatus(status) {
		s.log.Error("Service: Invalid status | ID=%s | Status=%s", id, status)
		return domain.ErrInvalidStatus
	}
	
	err := s.orderRepo.UpdateStatus(ctx, id, status)
	if err != nil {
		s.log.Error("Service: Failed to update order status | ID=%s | Error: %v", id, err)
		return err
	}
	
	s.log.Info("✅ Service: Order status updated | ID=%s | NewStatus=%s", id, status)
	return nil
}
