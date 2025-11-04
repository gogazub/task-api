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

// StartConsumer connect to kafka and start main loop
func StartConsumer(messageProcessor *service.MessageService) error {
	user, password, adress, port := getFullAdress()
	conn, err := amqp.Dial(fmt.Sprintf("amqp://%s:%s@%s:%s/", user, password, adress, port))

	if err != nil {
		return fmt.Errorf("start consumer error: %w", err)
	}
	defer func() {
		err = conn.Close()
		if err != nil {
			log.Printf("warning: %s", err.Error())
		}
	}()

	ch, err := conn.Channel()
	if err != nil {
		return err
	}

	defer func() {
		err = ch.Close()
		if err != nil {
			log.Printf("warning: %s", err.Error())
		}
	}()

	q, err := ch.QueueDeclare("code_to_run", false, false, false, false, nil)
	if err != nil {
		return err
	}
	fmt.Println("start consume...")
	msgs, err := ch.Consume(q.Name, "code_runner", true, false, false, false, nil)
	if err != nil {
		return err
	}

	// TODO: вынести в main
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	for {
		select {
		case msg, ok := <-msgs:
			if !ok {
				return err
			}
			err := messageProcessor.Accept(msg)
			if err != nil {
				return fmt.Errorf("consumer error: %w", err)
			}
			log.Printf("Recivied message: %s\n", msg.Body)
		case <-quit:
			log.Print("Shutting down gracefully...")
			return err
		}
	}
}

// getFullAdress
func getFullAdress() (user, password, adress, port string) {

	user = os.Getenv("RABBITMQ_USER")
	password = os.Getenv("RABBITMQ_PASSWORD")
	adress = os.Getenv("RABBITMQ_ADRESS")
	port = os.Getenv("RABBITMQ_PORT")
	return
}
