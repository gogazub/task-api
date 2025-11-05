package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/gogazub/consumer/consumer"
	"github.com/gogazub/consumer/repository"
	"github.com/gogazub/consumer/runner"
	"github.com/gogazub/consumer/service"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Printf("warning: no .env loaded: %v", err)
	}

	cr, err := runner.NewCodeRunner()
	if err != nil {
		log.Printf("launch error: %v", err)
		os.Exit(1)
	}

	db := connectToDB()
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
	resService := service.NewResultService(repo)
	mp := service.NewMessageService(*cr)

	err = consumer.StartConsumer(mp)
	if err != nil {
		log.Printf("launch error: %v", err)
	}
}

func connectToDB() *sql.DB {
	connStr := fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=%s",
		os.Getenv("POSTGRES_HOST"),
		os.Getenv("POSTGRES_PORT"), 
		os.Getenv("POSTGRES_USER"),
		os.Getenv("POSTGRES_DB"),
		os.Getenv("POSTGRES_SSLMODE"),
	)

	var db *sql.DB
	var err error

	maxRetries := 5
	for i := 0; i < maxRetries; i++ {
		db, err = sql.Open("postgres", connStr)
		if err != nil {
			log.Printf("failed to open database connection (attempt %d/%d): %v", i+1, maxRetries, err)
			time.Sleep(2 * time.Second)
			continue
		}

		err = db.Ping()
		if err != nil {
			log.Printf("failed to ping database (attempt %d/%d): %v", i+1, maxRetries, err)
			db.Close()
			time.Sleep(2 * time.Second)
			continue
		}

		break
	}

	if err != nil {
		log.Printf("failed to connect to database after %d attempts: %v", maxRetries, err)
		return nil
	}

	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)
	db.SetConnMaxLifetime(5 * time.Minute)

	return db
}
