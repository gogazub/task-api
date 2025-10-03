package service

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/gogazub/app/internal/model"
	amqp "github.com/rabbitmq/amqp091-go"
)

type Producer struct {
	conn *amqp.Connection
	ch   *amqp.Channel
	q    *amqp.Queue
}

func NewProducer() (*Producer, error) {
	user, password, adress, port := getFullAdress()
	conn, err := amqp.Dial(fmt.Sprintf("amqp://%s:%s@%s:%s/", user, password, adress, port))
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		return nil, err
	}
	defer ch.Close()

	q, err := ch.QueueDeclare("code_to_run", false, false, false, false, nil)
	if err != nil {
		return nil, err
	}

	return &Producer{
		conn: conn,
		ch:   ch,
		q:    &q,
	}, nil

}

func (p *Producer) SendMessage(task model.Task) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	body := []byte(task.Code)
	err := p.ch.PublishWithContext(ctx, "", "code_to_run", false, false, amqp.Publishing{
		ContentType: "text/plain",
		Body:        body,
	})
	if err != nil {
		return err
	}
	log.Printf("Send %s", body)
	return nil
}

func getFullAdress() (user, password, adress, port string) {

	user = os.Getenv("RABBITMQ_USER")
	password = os.Getenv("RABBITMQ_PASSWORD")
	adress = os.Getenv("RABBITMQ_ADRESS")
	port = os.Getenv("RABBITMQ_PORT")
	return
}
