package service

import (
	"encoding/json"
	"fmt"

	"github.com/gogazub/consumer/model"
	runner "github.com/gogazub/consumer/runner"
	amqp "github.com/rabbitmq/amqp091-go"
)

// MessageService accepts message from kafka and send them to CodeRunner
type MessageService struct {
	processor runner.ICodeRunner
}

func NewMessageService(processor runner.CodeRunner) *MessageService {
	return &MessageService{processor: processor}
}

func (mp *MessageService) Accept(msg amqp.Delivery) error {
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
