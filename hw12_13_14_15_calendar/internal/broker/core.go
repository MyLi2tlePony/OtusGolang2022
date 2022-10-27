package broker

import (
	"context"

	amqp "github.com/rabbitmq/amqp091-go"
)

type Broker interface {
	Start() error
	Stop() error
	Consume(ConsumeConfig) (<-chan amqp.Delivery, error)
	PublishWithContext(context.Context, PublishConfig, []byte) error
	QueueDeclare(config QueueConfig) error
}

type ConnectionConfig interface {
	Login() string
	Password() string
	Host() string
	Port() string
}

type PublishConfig interface {
	Exchange() string
	Key() string
	ContentType() string
	Mandatory() bool
	Immediate() bool
}

type QueueConfig interface {
	Name() string
	Durable() bool
	AutoDelete() bool
	Exclusive() bool
	NoWait() bool
}

type ConsumeConfig interface {
	Queue() string
	Consumer() string
	AutoAck() bool
	Exclusive() bool
	NoLocal() bool
	NoWait() bool
}

type Logger interface {
	Fatal(string)
	Error(string)
	Warn(string)
	Info(string)
	Debug(string)
	Trace(string)
}
