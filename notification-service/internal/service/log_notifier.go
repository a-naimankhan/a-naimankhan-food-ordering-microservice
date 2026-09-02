package service

import (
	"context"
	"notification-service/internal/domain"
	"notification-service/internal/logger"
)

type LogNotifier struct {
	log *logger.Logger
}

func (n *LogNotifier) Send(ctx context.Context, notification *domain.Notification) error {
	n.log.Info("NOTIFICATION SENT", "type", notification.Type, "order_id", notification.OrderID)
	return nil
}
