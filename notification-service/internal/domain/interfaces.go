package domain

import (
	"context"
)

type NotificationService interface {
	HandleOrderCreated(ctx context.Context, event *OrderEvent) error
	HandleOrderCancelled(ctx context.Context, event *OrderEvent) error
}

type EventConsumer interface {
	Start(ctx context.Context) error
	Close() error
}

type Notifier interface {
	Send(ctx context.Context, notification *Notification) error
}
