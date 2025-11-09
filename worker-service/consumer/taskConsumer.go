package consumer

import (
	"context"
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
	address    string
}

func NewTaskConsumer(msgService *service.MessageService) *TaskConsumer {
	return &TaskConsumer{msgService: msgService, address: ""}
}

func (tc *TaskConsumer) Connect(address string) error {
	conn, err := amqp.Dial(address)
	if err != nil {
		return fmt.Errorf("consumer connection error: %w", err)
	}

	tc.conn = conn
	tc.address = address
	return nil
}

func (tc *TaskConsumer) StartConsuming(ctx context.Context) error {
	defer tc.Close()

	if tc.conn == nil {
		return fmt.Errorf("start consuming error: no connection")
	}

	ch, err := tc.conn.Channel()
	if err != nil {
		return fmt.Errorf("start consuming error: %w", err)
	}

	q, err := ch.QueueDeclare("code_to_run", false, false, false, false, nil)
	if err != nil {
		return fmt.Errorf("start consuming error: %w", err)
	}

	msgs, err := ch.Consume(q.Name, "code_runner", true, false, false, false, nil)
	if err != nil {
		return fmt.Errorf("start consuming error: %w", err)
	}

	for {
		select {

		case msg, ok := <-msgs:
			// TODO: process !ok
			if !ok {
				continue
			}
			err := tc.msgService.Accept(msg)
			if err != nil {
				// TODO: process err
				continue
			}

		// graceful shutdown
		case <-ctx.Done():
			return tc.Close()
		}

	}
}

func (tc *TaskConsumer) Close() error {
	return tc.conn.Close()
}
