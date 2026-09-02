package domain

import "errors"

var (
	ErrInvalidEventPayload = errors.New("invalid event payload")
	ErrUnknownEventType    = errors.New("unknown event type")
	ErrNotificationFailed  = errors.New("notification failed")
)
