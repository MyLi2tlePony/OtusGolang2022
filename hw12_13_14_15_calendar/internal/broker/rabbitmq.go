package broker

import (
	"context"
	"fmt"

	amqp "github.com/rabbitmq/amqp091-go"
)

type broker struct {
	url string

	ch   *amqp.Channel
	conn *amqp.Connection
}

func New(connConfig ConnectionConfig) Broker {
	url := fmt.Sprintf("amqp://%s:%s@%s:%s",
		connConfig.Login(), connConfig.Password(), connConfig.Host(), connConfig.Port())

	return &broker{
		url: url,
	}
}

func (b *broker) Start() (err error) {
	b.conn, err = amqp.Dial(b.url)
	if err != nil {
		return fmt.Errorf("failed to connect to RabbitMQ %w", err)
	}

	b.ch, err = b.conn.Channel()
	if err != nil {
		return fmt.Errorf("failed to open a channel %w", err)
	}

	return nil
}

func (b *broker) QueueDeclare(config QueueConfig) error {
	_, err := b.ch.QueueDeclare(
		config.Name(),
		config.Durable(),
		config.AutoDelete(),
		config.Exclusive(),
		config.NoWait(),
		nil,
	)
	if err != nil {
		return fmt.Errorf("failed to declare a queue %w", err)
	}

	return nil
}

func (b *broker) Stop() (err error) {
	if err = b.ch.Close(); err != nil {
		return fmt.Errorf("failed to close a channel %w", err)
	}

	if err = b.conn.Close(); err != nil {
		return fmt.Errorf("failed to close connect to RabbitMQ %w", err)
	}

	return nil
}

func (b *broker) Consume(config ConsumeConfig) (<-chan amqp.Delivery, error) {
	delivery, err := b.ch.Consume(
		config.Queue(),
		config.Consumer(),
		config.AutoAck(),
		config.Exclusive(),
		config.NoLocal(),
		config.NoWait(),
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to register a consumer %w", err)
	}

	return delivery, nil
}

func (b *broker) PublishWithContext(ctx context.Context, config PublishConfig, body []byte) error {
	err := b.ch.PublishWithContext(ctx,
		config.Exchange(),
		config.Key(),
		config.Mandatory(),
		config.Immediate(),
		amqp.Publishing{
			ContentType: config.ContentType(),
			Body:        body,
		})
	if err != nil {
		return fmt.Errorf("failed to  publish a message %w", err)
	}

	return nil
}
