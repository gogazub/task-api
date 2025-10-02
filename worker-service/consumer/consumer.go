package consumer

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	amqp "github.com/rabbitmq/amqp091-go"
)

func StartConsumer(messageProcessor *MessageProcessor) {
	user, password, adress, port := getFullAdress()
	conn, err := amqp.Dial(fmt.Sprintf("amqp://%s:%s@%s:%s/", user, password, adress, port))
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		log.Println(err.Error())
	}
	defer ch.Close()

	q, err := ch.QueueDeclare("code_to_run", false, false, false, false, nil)
	if err != nil {
		log.Println(err.Error())
	}

	msgs, err := ch.Consume(q.Name, "code_runner", true, false, false, false, nil)
	if err != nil {
		log.Println(err.Error())
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	for {
		select {
		case msg, ok := <-msgs:
			if !ok {
				log.Printf("Channel is closed")
				return
			}
			messageProcessor.Accept(msg)
			log.Printf("Recivied message: %s\n", msg.Body)
		case <-quit:
			log.Print("Shutting down gracefully...")
			return
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
