package main

import (
	"log"
	"os"

	"github.com/gogazub/consumer/config"
	"github.com/gogazub/consumer/consumer"
	"github.com/gogazub/consumer/repository"
	"github.com/gogazub/consumer/runner"
	"github.com/gogazub/consumer/service"
)

func main() {
	err := config.InitConfig()
	if err != nil {
		log.Printf(err.Error())
		os.Exit(1)
	}

	cr, err := runner.NewCodeRunner()
	if err != nil {
		log.Printf("launch error: %v", err)
		os.Exit(1)
	}

	db := repository.ConnectToDB()
	if db == nil {
		log.Printf("failed to connect to database")
		os.Exit(1)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Printf("database ping error: %v", err)
		os.Exit(1)
	}

	log.Println("Successfully connected to database")

	repo := repository.NewOrderRepository(db)

	mp := service.NewMessageService(*cr, repo)

	err = consumer.StartConsumer(mp)
	if err != nil {
		log.Printf("launch error: %v", err)
	}
}
