package main

import (
	"context"
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
	must(err)

	codeRunner, err := runner.NewCodeRunner()
	must(err)

	db := repository.ConnectToDB()
	must(err)
	defer db.Close()

	err = db.Ping()
	must(err)

	orderRepo := repository.NewOrderRepository(db)
	msgService := service.NewMessageService(*codeRunner, orderRepo)

	taskConsumer := consumer.NewTaskConsumer(msgService)
	taskConsumer.Connect(config.Cfg.RABBITMQ_URL)
	taskConsumerCtx, taskConsumerCancel := context.WithCancel(context.Background())
	go func() {
		taskConsumer.StartConsuming(taskConsumerCtx)
	}()

}

func must(err error) {
	if err != nil {
		log.Printf(err.Error())
		os.Exit(1)
	}
}
