package delivery

import (
	"context"
	"log/slog"
	"notification-service/internal/domain"
)

type NotificationHandler struct {
	logger *slog.Logger
}

func NewNotificationHandler(logger *slog.Logger) *NotificationHandler {
	return &NotificationHandler{
		logger: logger,
	}
}

func (h *NotificationHandler) HandleOrderEvent(ctx context.Context, order *domain.OrderEvent) error {
	h.logger.Debug("Handling order event", "order_id", order.OrderID, "event_type", order.Status)

	switch order.Status {
	case "created":
		h.logger.Info("Sending Notification : Order created", "order_id", order.OrderID, "customer_id", order.CustomerID, "amount", order.Amount)
	case "cancelled":
		h.logger.Info("Sending Notification : Order cancelled", "order_id", order.OrderID, "customer_id", order.CustomerID, "reason", order.CancelledReason)
	default:
		h.logger.Warn("Unknown order event type", "order_id", order.OrderID, "event_type", order.Status)
	}

	//ну пока сервиса конечно нету но с контекстом получается передам туда и еще на счет statuses нужно будет еще правильно спроектировать и сделать чтобы там были все статусы которые есть в order-service
	return nil
}
