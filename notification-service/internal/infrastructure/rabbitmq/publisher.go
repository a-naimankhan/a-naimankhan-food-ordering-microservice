package rabbitmq

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"strings"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

type Publisher struct {
	conn     *amqp.Connection
	ch       *amqp.Channel
	exchange string
	logger   *slog.Logger
}

func NewPublisher(conn *amqp.Connection, exchange string) (*Publisher, error) {
	if conn == nil {
		return nil, amqp.ErrClosed
	}

	ch, err := conn.Channel()
	if err != nil {
		return nil, err
	}

	if err := ch.ExchangeDeclare(
		exchange, // name
		"direct", // type
		true,     // durable
		false,    // auto-deleted
		false,    // internal
		false,    // no-wait
		nil,      // args
	); err != nil {
		_ = ch.Close()
		return nil, err
	}

	return &Publisher{
		conn:     conn,
		ch:       ch,
		exchange: exchange,
		logger:   slog.Default(),
	}, nil
}

func (p *Publisher) Close() error {
	var errs []error

	if p.ch != nil {
		if err := p.ch.Close(); err != nil {
			if !strings.Contains(err.Error(), "channel/connection is not open") {
				errs = append(errs, fmt.Errorf("close channel: %w", err))
			}
		}
		p.ch = nil
	}

	if p.conn != nil {
		if err := p.conn.Close(); err != nil {
			if !strings.Contains(err.Error(), "connection is not open") {
				errs = append(errs, fmt.Errorf("close connection: %w", err))
			}
		}
		p.conn = nil
	}

	if len(errs) > 0 {
		return errors.Join(errs...)
	}
	return nil
}

func (p *Publisher) PublishOrderCreated(ctx context.Context, orderID string, status string) error {
	if p.ch == nil {
		return amqp.ErrClosed
	}

	payload := map[string]string{
		"order_id": orderID,
		"status":   status,
	}
	body, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	pub := amqp.Publishing{
		ContentType: "application/json",
		Body:        body,
		Timestamp:   time.Now(),
		MessageId:   orderID,
	}

	routingKey := "order.created"

	if err := p.ch.PublishWithContext(ctx,
		p.exchange,
		routingKey,
		false, // mandatory
		false, // immediate
		pub,
	); err != nil {
		return err
	}

	p.logger.Info("Event order.created published successfully", "order_id", orderID, "exchange", p.exchange, "routing_key", routingKey)
	return nil
}
