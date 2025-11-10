package main

import (
	"context"
	"errors"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

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

	wg := sync.WaitGroup{}

	orderRepo := repository.NewOrderRepository(db)
	msgService := service.NewMessageService(*codeRunner, orderRepo)

	taskConsumer := consumer.NewTaskConsumer(msgService)
	taskConsumer.Connect(config.Cfg.RABBITMQ_URL)
	taskConsumerCtx, cancel := context.WithCancel(context.Background())

	errCh := make(chan error, 1)

	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := taskConsumer.StartConsuming(taskConsumerCtx); err != nil && !errors.Is(err, context.Canceled) {
			errCh <- err
		}
	}()

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	select {
	case err := <-errCh:
		log.Printf("component error: %v - initiating shutdown", err)
		cancel()
	case sig := <-sigCh:
		log.Printf("signal %v received - initiating shutdown", sig)
		cancel()
	}

	done := make(chan struct{})

	go func() {
		wg.Wait()
		close(done)
	}()

	timeout := 40 * time.Second
	select {
	case <-done:
		log.Println("all components stopped gracefully")
	case <-time.After(timeout):
		log.Printf("timeout (%s) waiting for components to stop, exiting", timeout)
	}

	log.Println("shutdown complete")

}

func must(err error) {
	if err != nil {
		log.Print(err.Error())
		os.Exit(1)
	}
}
