package main

import (
	"log"
	"os"

	"github.com/gogazub/consumer/consumer"
	"github.com/gogazub/consumer/runner"
	"github.com/gogazub/consumer/service"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Printf("warning: no .env loaded: %v", err)
	}

	// runner.Test()
	// return

	cr, err := runner.NewCodeRunner()
	if err != nil {
		log.Printf("launch error: %v", err)
		os.Exit(1)
	}

	mp := service.NewMessageProcessor(cr)

	err = consumer.StartConsumer(mp)
	if err != nil {
		log.Printf("launch error: %v", err)
	}
}
