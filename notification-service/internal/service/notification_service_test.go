package service

import (
	"context"
	"errors"
	"testing"

	"notification-service/internal/domain"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type fakeNotifier struct {
	SendFn func(ctx context.Context, notification *domain.Notification) error
	calls  []*domain.Notification
}

func (f *fakeNotifier) Send(ctx context.Context, notification *domain.Notification) error {
	f.calls = append(f.calls, notification)
	if f.SendFn != nil {
		return f.SendFn(ctx, notification)
	}
	return nil
}

func strPtr(s string) *string {
	return &s
}

func validOrderEvent() *domain.OrderEvent {
	return &domain.OrderEvent{
		OrderID:    uuid.New(),
		CustomerID: uuid.New(),
		Status:     "pending",
		Amount:     100.0,
	}
}

func validCancelledEvent() *domain.OrderEvent {
	reason := "customer request"
	return &domain.OrderEvent{
		OrderID:         uuid.New(),
		CustomerID:      uuid.New(),
		Status:          "cancelled",
		Amount:          100.0,
		CancelledReason: &reason,
	}
}

func TestNotificationService_HandleOrderCreated(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name          string
		event         *domain.OrderEvent
		notifier      *fakeNotifier
		wantErr       error
		wantSendCalls int
		assertNotify  func(t *testing.T, n *domain.Notification)
	}{
		{
			name:  "success",
			event: validOrderEvent(),
			notifier: &fakeNotifier{
				SendFn: func(ctx context.Context, notification *domain.Notification) error {
					return nil
				},
			},
			wantErr:       nil,
			wantSendCalls: 1,
			assertNotify: func(t *testing.T, n *domain.Notification) {
				assert.NotEqual(t, uuid.Nil, n.ID)
				assert.Equal(t, domain.EventCreated, n.Type)
				assert.Contains(t, n.Message, n.OrderID.String())
				assert.Contains(t, n.Message, "has been created")
				assert.False(t, n.CreatedAt.IsZero())
			},
		},
		{
			name:          "nil event",
			event:         nil,
			notifier:      &fakeNotifier{},
			wantErr:       domain.ErrInvalidEventPayload,
			wantSendCalls: 0,
		},
		{
			name: "nil order_id",
			event: &domain.OrderEvent{
				OrderID:    uuid.Nil,
				CustomerID: uuid.New(),
			},
			notifier:      &fakeNotifier{},
			wantErr:       domain.ErrInvalidEventPayload,
			wantSendCalls: 0,
		},
		{
			name: "nil customer_id",
			event: &domain.OrderEvent{
				OrderID:    uuid.New(),
				CustomerID: uuid.Nil,
			},
			notifier:      &fakeNotifier{},
			wantErr:       domain.ErrInvalidEventPayload,
			wantSendCalls: 0,
		},
		{
			name:  "notifier returns error",
			event: validOrderEvent(),
			notifier: &fakeNotifier{
				SendFn: func(ctx context.Context, notification *domain.Notification) error {
					return errors.New("smtp unavailable")
				},
			},
			wantErr:       domain.ErrNotificationFailed,
			wantSendCalls: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := NewNotificationService(tt.notifier)

			err := svc.HandleOrderCreated(ctx, tt.event)

			if tt.wantErr != nil {
				require.Error(t, err)
				assert.ErrorIs(t, err, tt.wantErr)
			} else {
				require.NoError(t, err)
			}

			assert.Len(t, tt.notifier.calls, tt.wantSendCalls)

			if tt.wantSendCalls > 0 && tt.assertNotify != nil {
				require.NotNil(t, tt.event)
				tt.assertNotify(t, tt.notifier.calls[0])
				assert.Equal(t, tt.event.OrderID, tt.notifier.calls[0].OrderID)
				assert.Equal(t, tt.event.CustomerID, tt.notifier.calls[0].CustomerID)
			}
		})
	}
}

func TestNotificationService_HandleOrderCancelled(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name          string
		event         *domain.OrderEvent
		notifier      *fakeNotifier
		wantErr       error
		wantSendCalls int
		assertNotify  func(t *testing.T, n *domain.Notification)
	}{
		{
			name:  "success",
			event: validCancelledEvent(),
			notifier: &fakeNotifier{
				SendFn: func(ctx context.Context, notification *domain.Notification) error {
					return nil
				},
			},
			wantErr:       nil,
			wantSendCalls: 1,
			assertNotify: func(t *testing.T, n *domain.Notification) {
				assert.NotEqual(t, uuid.Nil, n.ID)
				assert.Equal(t, domain.EventOrderCancelled, n.Type)
				assert.Contains(t, n.Message, n.OrderID.String())
				assert.Contains(t, n.Message, "has been cancelled")
				assert.False(t, n.CreatedAt.IsZero())
			},
		},
		{
			name:          "nil event",
			event:         nil,
			notifier:      &fakeNotifier{},
			wantErr:       domain.ErrInvalidEventPayload,
			wantSendCalls: 0,
		},
		{
			name: "nil order_id",
			event: &domain.OrderEvent{
				OrderID:         uuid.Nil,
				CustomerID:      uuid.New(),
				CancelledReason: strPtr("reason"),
			},
			notifier:      &fakeNotifier{},
			wantErr:       domain.ErrInvalidEventPayload,
			wantSendCalls: 0,
		},
		{
			name: "nil customer_id",
			event: &domain.OrderEvent{
				OrderID:         uuid.New(),
				CustomerID:      uuid.Nil,
				CancelledReason: strPtr("reason"),
			},
			notifier:      &fakeNotifier{},
			wantErr:       domain.ErrInvalidEventPayload,
			wantSendCalls: 0,
		},
		//{
		//	name: "nil cancelled_reason pointer",
		//	event: &domain.OrderEvent{
		//		OrderID:         uuid.New(),
		//		CustomerID:      uuid.New(),
		//		CancelledReason: nil,
		//	},
		//	notifier:      &fakeNotifier{},
		//	wantErr:       domain.ErrInvalidEventPayload,
		//	wantSendCalls: 0,
		//},
		//{
		//	name: "empty cancelled_reason",
		//	event: &domain.OrderEvent{
		//		OrderID:         uuid.New(),
		//		CustomerID:      uuid.New(),
		//		CancelledReason: strPtr(""),
		//	},
		//	notifier:      &fakeNotifier{},
		//	wantErr:       domain.ErrInvalidEventPayload,
		//	wantSendCalls: 0,
		//},
		{
			name:  "notifier returns error",
			event: validCancelledEvent(),
			notifier: &fakeNotifier{
				SendFn: func(ctx context.Context, notification *domain.Notification) error {
					return errors.New("push gateway down")
				},
			},
			wantErr:       domain.ErrNotificationFailed,
			wantSendCalls: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := NewNotificationService(tt.notifier)

			err := svc.HandleOrderCancelled(ctx, tt.event)

			if tt.wantErr != nil {
				require.Error(t, err)
				assert.ErrorIs(t, err, tt.wantErr)
			} else {
				require.NoError(t, err)
			}

			assert.Len(t, tt.notifier.calls, tt.wantSendCalls)

			if tt.wantSendCalls > 0 && tt.assertNotify != nil {
				require.NotNil(t, tt.event)
				tt.assertNotify(t, tt.notifier.calls[0])
				assert.Equal(t, tt.event.OrderID, tt.notifier.calls[0].OrderID)
				assert.Equal(t, tt.event.CustomerID, tt.notifier.calls[0].CustomerID)
			}
		})
	}
}

func TestNotificationService_NotifierNotCalledOnValidationError(t *testing.T) {
	ctx := context.Background()
	notifier := &fakeNotifier{}
	svc := NewNotificationService(notifier)

	invalidCreated := &domain.OrderEvent{OrderID: uuid.Nil, CustomerID: uuid.New()}
	require.ErrorIs(t, svc.HandleOrderCreated(ctx, invalidCreated), domain.ErrInvalidEventPayload)

	invalidCancelled := &domain.OrderEvent{
		OrderID:    uuid.New(),
		CustomerID: uuid.Nil,
	}
	require.ErrorIs(t, svc.HandleOrderCancelled(ctx, invalidCancelled), domain.ErrInvalidEventPayload)

	assert.Empty(t, notifier.calls)
}
