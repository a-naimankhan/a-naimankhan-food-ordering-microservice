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
	s.log.Debug("Getting order", "id", id)

	order, err := s.orderRepo.GetByID(ctx, id)
	if err != nil {
		s.log.Error("Failed to get order", "id", id, "error", err.Error())
		return nil, err
	}

	if order == nil {
		s.log.Error("Order not found", "id", id)
		return nil, domain.ErrOrderNotFound
	}

	s.log.Info("Order retrieved", "id", id, "status", order.Status, "amount", order.Amount)
	return order, nil
}

func (s *orderService) CreateOrder(ctx context.Context, order *domain.Order) (*domain.Order, error) {
	s.log.Debug("Creating order", "customer_id", order.CustomerID, "amount", order.Amount, "status", order.Status)

	if order == nil {
		s.log.Error("Order is nil")
		return nil, domain.ErrInvalidCustomerID
	}

	if order.Amount <= 0 {
		s.log.Error("Invalid amount", "amount", order.Amount)
		return nil, domain.ErrInvalidAmount
	}

	if order.Status == "" {
		s.log.Error("Status is empty")
		return nil, domain.ErrStatusEmpty
	}

	if !domain.IsValidStatus(order.Status) {
		s.log.Error("Invalid status", "status", order.Status)
		return nil, domain.ErrInvalidStatus
	}

	if order.CustomerID == uuid.Nil {
		s.log.Error("Invalid customer ID")
		return nil, domain.ErrInvalidCustomerID
	}

	if order.ID == uuid.Nil {
		order.ID = uuid.New()
		s.log.Debug("Generated new Order ID", "id", order.ID)
	}

	if order.CreatedAt.IsZero() {
		order.CreatedAt = time.Now()
		s.log.Debug("Set creation timestamp", "created_at", order.CreatedAt)
	}

	s.log.Info("Validations passed, saving order", "id", order.ID)
	if err := s.orderRepo.Create(ctx, order); err != nil {
		s.log.Error("Failed to create order in repository", "id", order.ID, "error", err.Error())
		return nil, err
	}

	s.log.Info("Order saved to database", "id", order.ID)

	if s.eventPublisher != nil {
		s.log.Debug("Publishing order.created event", "id", order.ID)
		if err := s.eventPublisher.Publish(ctx, "order.created", order); err != nil {
			s.log.Error("Failed to publish order.created event", "id", order.ID, "error", err.Error())
		}
	}

	s.log.Info("Order created successfully", "id", order.ID, "customer_id", order.CustomerID)
	return order, nil
}

func (s *orderService) UpdateOrderStatus(ctx context.Context, id uuid.UUID, status string) error {
	s.log.Debug("Updating order status", "id", id, "new_status", status)

	if status == "" {
		s.log.Error("Status is empty", "id", id)
		return domain.ErrStatusEmpty
	}

	//checks if is updated status is valid
	if !(domain.IsValidStatus(status)) {
		s.log.Error("Invalid status", "id", id, "status", status)
		return domain.ErrInvalidStatus
	}

	currentOrder, err := s.orderRepo.GetByID(ctx, id)
	if err != nil {
		s.log.Error("Failed to get order", "id", id, "error", err.Error())
		return err
	}

	//checks if order really appeared
	if currentOrder == nil || currentOrder.ID == uuid.Nil {
		s.log.Error("Order not found or empty", "id", id)
		return domain.ErrOrderNotFound
	}

	//checks if new update status transaction is possible
	if !(domain.IsValidTransaction(currentOrder.Status, status)) {
		s.log.Error("Not valid status for this specific transaction", "id", id, "current_status", currentOrder.Status, "updated_status", status)
		return domain.ErrInvalidTransaction
	}

	//updates status in db
	err = s.orderRepo.UpdateStatus(ctx, id, status)
	if err != nil {
		s.log.Error("Failed to update order status", "id", id, "error", err.Error())
		return err
	}

	s.log.Info("Order status updated", "id", id, "new_status", status)
	return nil
}

//func (s *orderService) CancelOrder(ctx context.Context, id uuid.UUID, reason string) error {
//	//TODO finish this
//
//}
