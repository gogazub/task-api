package service

import (
	"encoding/json"
	"fmt"

	"github.com/gogazub/consumer/model"
	runner "github.com/gogazub/consumer/runner"
	amqp "github.com/rabbitmq/amqp091-go"
)

type MessageProcessor struct {
	processor *runner.CodeProcessor
}

func NewMessageProcessor(processor *runner.CodeProcessor) *MessageProcessor {
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

func getCodeMessage(msg amqp.Delivery) (*model.CodeMessage, error) {
	var codeMessage model.CodeMessage
	err := json.Unmarshal(msg.Body, &codeMessage)
	if err != nil {
		return nil, err
	}
	return &codeMessage, nil
}
