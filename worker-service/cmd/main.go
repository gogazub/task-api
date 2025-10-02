package main

import (
	"log"

	"github.com/gogazub/consumer/consumer"
	"github.com/gogazub/consumer/runner"
	"github.com/gogazub/consumer/service"
	"github.com/joho/godotenv"
)

func main() {
	_ = godotenv.Load()
	cr := runner.NewCodeProcessor()

	mp := service.NewMessageProcessor(cr)

	err := consumer.StartConsumer(mp)
	if err != nil {
		log.Println(err.Error())
	}
}
