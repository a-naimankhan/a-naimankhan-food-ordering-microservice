package rabbitmq

import (
	"context"
	"errors"
	"testing"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/stretchr/testify/assert"
)

func TestNewPublisher_NilConn(t *testing.T) {
	p, err := NewPublisher(nil, "test-exchange")
	assert.Error(t, err)
	assert.Nil(t, p)
}

func TestPublishOrderCreated_ChannelNil(t *testing.T) {
	// Create a publisher with nil channel to simulate closed/uninitialized channel
	p := &Publisher{conn: nil, ch: nil, exchange: "test", logger: nil}
	err := p.PublishOrderCreated(context.Background(), "id-123", "pending")
	assert.Error(t, err)
	// amqp.ErrClosed is returned when channel/connection is not available
	assert.True(t, errors.Is(err, amqp.ErrClosed))
}

// Basic marshalling sanity check: when ch present we cannot actually publish in unit test without rabbitmq,
// but ensure PublishOrderCreated returns an error when Publish fails. We simulate by creating a fake channel
// by embedding behavior is not possible since amqp.Channel is concrete. Thus keep tests focused on constructors and nil-ch behavior.
func TestClose_NoPanic(t *testing.T) {
	p := &Publisher{conn: nil, ch: nil, exchange: "test", logger: nil}
	// Close should not panic and should return nil when nothing to close
	err := p.Close()
	assert.NoError(t, err)
}
