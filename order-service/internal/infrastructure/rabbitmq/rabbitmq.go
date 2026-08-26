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
	log.Debug("RabbitMQ: Initializing publisher | Exchange=%s", exchange)
	return &RabbitPublisher{
		channel:  ch,
		exchange: exchange,
		log:      log,
	}
}

func (r *RabbitPublisher) Publish(ctx context.Context, topic string, payload interface{}) error {
	r.log.Debug("RabbitMQ: Publishing event | Topic=%s | Exchange=%s", topic, r.exchange)
	
	body, err := json.Marshal(payload)
	if err != nil {
		r.log.Error("RabbitMQ: Failed to marshal payload | Topic=%s | Error: %v", topic, err)
		return err
	}

	r.log.Debug("RabbitMQ: Serialized payload | Topic=%s | Size=%d bytes", topic, len(body))

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
		r.log.Error("RabbitMQ: Failed to publish event | Topic=%s | Error: %v", topic, err)
		return err
	}

	r.log.Info("✅ RabbitMQ: Event published successfully | Topic=%s | Exchange=%s", topic, r.exchange)
	return nil
}
