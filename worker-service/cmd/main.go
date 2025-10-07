package main

import (
	"fmt"
	"log"

	"github.com/gogazub/consumer/consumer"
	"github.com/gogazub/consumer/runner"
	"github.com/gogazub/consumer/service"
	"github.com/joho/godotenv"
)

func main() {
	_ = godotenv.Load()

	runner.Test()
	return
	cr, err := runner.NewCodeRunner()
	if err != nil {
		fmt.Printf("CodeRunner creation error: %s", err.Error())
		return
	}

	mp := service.NewMessageProcessor(cr)

	err = consumer.StartConsumer(mp)
	if err != nil {
		log.Println(err.Error())
	}
}
