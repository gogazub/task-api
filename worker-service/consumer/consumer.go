package consumer

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	amqp "github.com/rabbitmq/amqp091-go"
)

type CodeMessage struct {
	Code string `json:"code"`
}

type MessageProcessor struct {
}

func GetCodeMessage(msg amqp.Delivery) (*CodeMessage, error) {
	var codeMessage CodeMessage
	err := json.Unmarshal(msg.Body, &codeMessage)
	if err != nil {
		return nil, err
	}
	return &codeMessage, nil
}

func (mp *MessageProcessor) Process(msg amqp.Delivery) error {
	codeMessage, err := GetCodeMessage(msg)
	if err != nil {
		return err
	}
	fmt.Printf("Recived msg: %s", codeMessage.Code)
	return nil
}

func getFullAdress() (user, password, adress, port string) {
	user = os.Getenv("RABBITMQ_USER")
	password = os.Getenv("RABBITMQ_PASSWORD")
	adress = os.Getenv("RABBITMQ_ADRESS")
	port = os.Getenv("RABBITMQ_PORT")
	return
}

func StartConsumer(messageProcessor *MessageProcessor) {

	user, password, adress, port := getFullAdress()

	conn, err := amqp.Dial(fmt.Sprint("amqp://%s:%s@%s:%es/", user, password, adress, port))
	if err != nil {
		log.Printf(err.Error())
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
			messageProcessor.Process(msg)
			log.Printf("Recivied message: %s\n", msg.Body)
		case <-quit:
			log.Print("Shutting down gracefully...")
			return
		}
	}
}
