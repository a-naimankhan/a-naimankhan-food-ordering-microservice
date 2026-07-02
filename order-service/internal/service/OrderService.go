package service

import (
	"context"
	"errors"
	"order-service/internal/domain"
	"time"

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

func (s *OrderService) CreateOrder(ctx context.Context, customerID string, amount float64, status string) (*domain.Order, error) {

	cstmId, err := uuid.Parse(customerID)
	if err != nil {
		return nil, errors.New("invalid customer id")
	}

	if amount <= 0 {
		return nil, errors.New("amount must be greater than zero")
	}

	if status == "" {
		return nil, errors.New("status is empty")
	}

	order := &domain.Order{
		ID:         uuid.New(),
		CustomerID: cstmId,
		Amount:     amount,
		Status:     status,
		CreatedAt:  time.Now(),
	}

	err = s.orderRepo.Create(ctx, order)
	if err != nil {
		return nil, err
	}

	return order, nil
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
