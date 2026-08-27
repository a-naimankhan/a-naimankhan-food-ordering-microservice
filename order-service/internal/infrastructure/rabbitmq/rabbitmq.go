package rabbitmq

import (
	"context"
	"order-service/internal/logger"

	"github.com/goccy/go-json"
	amqp "github.com/rabbitmq/amqp091-go"
)

type RabbitPublisher struct {
	channel  *amqp.Channel
	exchange string
	log      *logger.Logger
}

func NewRabbitPublisher(ch *amqp.Channel, exchange string) *RabbitPublisher {
	log := logger.GetLogger()
	log.Debug("Initializing RabbitMQ publisher", "exchange", exchange)
	return &RabbitPublisher{
		channel:  ch,
		exchange: exchange,
		log:      log,
	}
}

func (r *RabbitPublisher) Publish(ctx context.Context, topic string, payload interface{}) error {
	r.log.Debug("Publishing event", "topic", topic, "exchange", r.exchange)
	
	body, err := json.Marshal(payload)
	if err != nil {
		r.log.Error("Failed to marshal payload", "topic", topic, "error", err.Error())
		return err
	}

	r.log.Debug("Serialized payload", "topic", topic, "size", len(body))

	err = r.channel.Publish(
		r.exchange,
		topic,
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		},
	)

	if err != nil {
		r.log.Error("Failed to publish event", "topic", topic, "error", err.Error())
		return err
	}

	r.log.Info("Event published successfully", "topic", topic, "exchange", r.exchange)
	return nil
}
