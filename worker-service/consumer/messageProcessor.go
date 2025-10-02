package consumer

import (
	"encoding/json"
	"fmt"

	"github.com/gogazub/consumer/code"
	"github.com/gogazub/consumer/processor"
	amqp "github.com/rabbitmq/amqp091-go"
)

type MessageProcessor struct {
	processor *processor.CodeProcessor
}

func NewMessageProcessor(processor *processor.CodeProcessor) *MessageProcessor {
	return &MessageProcessor{processor: processor}
}

func (mp *MessageProcessor) Accept(msg amqp.Delivery) error {
	codeMessage, err := getCodeMessage(msg)
	if err != nil {
		return err
	}
	fmt.Printf("Recived msg: %s", codeMessage.Code)
	return nil
}

func getCodeMessage(msg amqp.Delivery) (*code.CodeMessage, error) {
	var codeMessage code.CodeMessage
	err := json.Unmarshal(msg.Body, &codeMessage)
	if err != nil {
		return nil, err
	}
	return &codeMessage, nil
}
