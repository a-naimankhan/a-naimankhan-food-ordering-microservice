package service

import (
	"context"
	"notification-service/internal/domain"
	"notification-service/internal/logger"
	"time"

	"github.com/google/uuid"
)

type notificationService struct {
	notifier domain.Notifier
	log      *logger.Logger
}

func NewNotificationService(notifier domain.Notifier) domain.NotificationService {
	return &notificationService{
		notifier: notifier,
		log:      logger.GetLogger(),
	}
}

func (s *notificationService) HandleOrderCreated(ctx context.Context, event *domain.OrderEvent) error {
	s.log.Debug("Handling order created")

	if event == nil {
		s.log.Error("Event is nil")
		return domain.ErrInvalidEventPayload
	}

	if event.OrderID == uuid.Nil {
		s.log.Error("Event order_id is nil", "order_id", event.OrderID)
		return domain.ErrInvalidEventPayload
	}

	if event.CustomerID == uuid.Nil {
		s.log.Error("Event customer_id is nil", "customer_id", event.CustomerID)
		return domain.ErrInvalidEventPayload //это ща рефакторить буду пока что временно
	}

	// create Notification object
	notification := &domain.Notification{
		ID:         uuid.New(),
		OrderID:    event.OrderID,
		CustomerID: event.CustomerID,
		Type:       domain.EventCreated,
		Message:    "Your order " + event.OrderID.String() + " has been created.",
		CreatedAt:  time.Now(),
	}

	//call notifier.Send(ctx, notification)
	if err := s.notifier.Send(ctx, notification); err != nil {
		s.log.Error("Failed to send notification", "order_id", event.OrderID, "customer_id", event.CustomerID, "error", err.Error())
		return domain.ErrNotificationFailed //опять рефакторить поидее
	}

	//log access
	s.log.Info("Notification successfully sent", "order_id", event.OrderID, "customer_id", event.CustomerID)
	return nil
}

func (s *notificationService) HandleOrderCancelled(ctx context.Context, event *domain.OrderEvent) error {
	if event == nil {
		s.log.Error("Event is nil")
		return domain.ErrInvalidEventPayload
	}

	s.log.Debug("Handling order cancelled event", "order_id", event.OrderID, "customer_id", event.CustomerID, "reason", event.CancelledReason)

	if event.OrderID == uuid.Nil {
		s.log.Error("Event order_id is nil", "order_id", event.OrderID)
		return domain.ErrInvalidEventPayload
	}

	if event.CustomerID == uuid.Nil {
		s.log.Error("Event customer_id is nil", "customer_id", event.CustomerID)
		return domain.ErrInvalidEventPayload
	}

	//reason := "No reason provided"
	//if event.CancelledReason != nil && *event.CancelledReason != "" {
	//	reason = *event.CancelledReason
	//} TODO уберу это все нафиг
	//if event.CancelledReason == nil || *event.CancelledReason == "" {
	//	s.log.Error("Event cancelled_reason is empty", "order_id", event.OrderID)
	//	//или лучше сделать его как null guard и отдавать дефолтный ризон
	//	defaultReason := "No reason provided"
	//	event.CancelledReason = &defaultReason
	//}

	notification := &domain.Notification{
		ID:         uuid.New(),
		OrderID:    event.OrderID,
		CustomerID: event.CustomerID,
		Type:       domain.EventOrderCancelled,
		Message:    "Your order " + event.OrderID.String() + " has been cancelled.",
		Reason:     event.CancelledReason,
		CreatedAt:  time.Now(),
	}

	if err := s.notifier.Send(ctx, notification); err != nil {
		s.log.Error("Failed to send notification", "order_id", event.OrderID, "customer_id", event.CustomerID, "error", err.Error())
		return domain.ErrNotificationFailed
	}

	s.log.Info("Notification successfully sent for order cancellation", "order_id", event.OrderID, "customer_id", event.CustomerID)
	return nil
}
