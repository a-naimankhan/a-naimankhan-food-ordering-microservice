package service

import (
	"context"
	"errors"
	"order-service/internal/domain"

	"github.com/google/uuid"
)

type OrderService struct {
	orderRepo domain.OrderRepository
}

func newOrderService(orderRepo domain.OrderRepository) *OrderService {
	return &OrderService{orderRepo}
}

func (s *OrderService) GetOrder(ctx context.Context, id uuid.UUID) (*domain.Order, error) {
	order, err := s.orderRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if order == nil {
		return nil, errors.New("order not found")
	}
	return order, nil
}

func (s *OrderService) PlaceOrder(ctx context.Context, order *domain.Order) error {
	if order == nil {
		return errors.New("order is nil")
	}

	return s.orderRepo.Create(ctx, order)
}

func (s *OrderService) UpdateOrderStatus(ctx context.Context, id uuid.UUID, status string) error {
	if status == "" {
		return errors.New("status is empty")
	}

	statuses := []string{"pending", "shipped", "delivered", "cancelled"}

	for _, st := range statuses {
		if status == st {
			return s.orderRepo.UpdateStatus(ctx, id, status)
		}
	}

	return errors.New("status is invalid")

}
