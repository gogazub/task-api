package consumer

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/gogazub/consumer/service"
	amqp "github.com/rabbitmq/amqp091-go"
)

func StartConsumer(messageProcessor *service.MessageService) error {
	user, password, adress, port := getFullAdress()
	conn, err := amqp.Dial(fmt.Sprintf("amqp://%s:%s@%s:%s/", user, password, adress, port))
	if err != nil {
		return err
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		return err
	}
	defer ch.Close()

	q, err := ch.QueueDeclare("code_to_run", false, false, false, false, nil)
	if err != nil {
		return err
	}

	msgs, err := ch.Consume(q.Name, "code_runner", true, false, false, false, nil)
	if err != nil {
		return err
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	for {
		select {
		case msg, ok := <-msgs:
			if !ok {
				return err
			}
			messageProcessor.Accept(msg)
			log.Printf("Recivied message: %s\n", msg.Body)
		case <-quit:
			log.Print("Shutting down gracefully...")
			return err
		}
	}
}

func getFullAdress() (user, password, adress, port string) {
	user = os.Getenv("RABBITMQ_USER")
	password = os.Getenv("RABBITMQ_PASSWORD")
	adress = os.Getenv("RABBITMQ_ADRESS")
	port = os.Getenv("RABBITMQ_PORT")
	return
}
