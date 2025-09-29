package main

import (
	"fmt"
	"log"
	"os"

	amqp "github.com/rabbitmq/amqp091-go"
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Panicf("%s: %s", msg, err)
	}
}

func main() {
	user := os.Getenv("RABBITMQ_USER")
	password := os.Getenv("RABBITMQ_PASSWORD")
	adress := os.Getenv("RABBITMQ_ADRESS")
	port := os.Getenv("RABBITMQ_PORT")

	dial := fmt.Sprintf("amqp://%s:%s@%s:%s", user, password, adress, port)
	conn, err := amqp.Dial(dial)
	failOnError(err, "connection fail")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "channel error")
	defer ch.Close()

	q, err := ch.QueueDeclare("task", false, false, false, false, nil)
	failOnError(err, "declaration error")

	for msg := range q.Messages {
		fmt.Printf("msg: %v\n", msg)
	}

}
