package consumer

import (
	"fmt"

	"github.com/gogazub/consumer/service"
	amqp "github.com/rabbitmq/amqp091-go"
)

type ITaskConsumer interface {
	Connect(address string) error

	StartConsuming() error
	Stop() error
	handleMessage(msg amqp.Delivery)
	AckMessage(msg amqp.Delivery)
}

type TaskConsumer struct {
	msgService *service.MessageService
	conn       *amqp.Connection
}

func NewTaskConsumer(msgService *service.MessageService) *TaskConsumer {
	return &TaskConsumer{msgService: msgService}
}

func (tc *TaskConsumer) Connect(address string) error {
	conn, err := amqp.Dial(address)
	if err != nil {
		return fmt.Errorf("consumer connection error: %w", err)
	}
	tc.conn = conn
	return nil
}
