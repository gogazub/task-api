package service

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/gogazub/consumer/model"
	"github.com/gogazub/consumer/repository"
	runner "github.com/gogazub/consumer/runner"
	amqp "github.com/rabbitmq/amqp091-go"
)

// MessageService accepts message from kafka and send them to CodeRunner
type MessageService struct {
	processor runner.ICodeRunner
	repository repository.IResultRepository
}

func NewMessageService(processor runner.CodeRunner, repository repository.IResultRepository) *MessageService {
	return &MessageService{processor: processor, repository: repository}
}

// Accept consume messages from rabbitMQ and send them to code runner
func (mp *MessageService) Accept(msg amqp.Delivery) error {
	codeMessage, err := getCodeMessage(msg)
	if err != nil {
		return err
	}
	log.Printf("Recived msg: %s", codeMessage.Code)
	result := mp.processor.RunCode(*codeMessage)
	ctx, _ := context.WithTimeout(context.Background(),1*time.Minute)
	err = mp.SaveResult(ctx, result)

	return nil
}

// getCodeMessage converts rabbitMQ message to model.CodeMessage
func getCodeMessage(msg amqp.Delivery) (*model.CodeMessage, error) {
	var codeMessage model.CodeMessage
	err := json.Unmarshal(msg.Body, &codeMessage)
	if err != nil {
		return nil, err
	}
	return &codeMessage, nil
}


func (svc *MessageService) SaveResult(ctx context.Context,result model.Result) error {
	err := svc.repository.Save(ctx, result)
	if err != nil{
		return fmt.Errorf("save result: %v",err)
	}
	return nil
}