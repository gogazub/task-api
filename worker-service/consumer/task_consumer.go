package consumer

import amqp "github.com/rabbitmq/amqp091-go"

type ITaskConsumer interface {
	Connect() error

	StartConsuming() error
	Stop() error
	handleMessage(msg amqp.Delivery)
	AckMessage(msg amqp.Delivery)
}

type TaskConsumer struct {
}
