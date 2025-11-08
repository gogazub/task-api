package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

var Cfg Config

type Config struct {
	// --- Server ---

	// --- RabbitMQ ---
	RABBITMQ_USER     string
	RABBITMQ_PASSWORD string
	RABBITMQ_ADDRESS  string
	RABBITMQ_PORT     string
	RABBITMQ_URL      string

	// --- Postgres ---
	POSTGRES_USER     string
	POSTGRES_PASSWORD string
	POSTGRES_HOST     string
	POSTGRES_PORT     string
	POSTGRES_DB       string
	POSTGRES_SSLMODE  string
	DATABASE_URL      string

	// --- Docker ---
}

// InitConfig loads environment variables from .env to the Config struct
func InitConfig() error {
	if err := godotenv.Load(); err != nil {
		return fmt.Errorf("init config error: %w", err)
	}

	Cfg = Config{
		// --- RabbitMQ ---
		RABBITMQ_USER:     getEnv("RABBITMQ_USER"),
		RABBITMQ_PASSWORD: getEnv("RABBITMQ_PASSWORD"),
		RABBITMQ_ADDRESS:  getEnv("RABBITMQ_ADDRESS"),
		RABBITMQ_PORT:     getEnv("RABBITMQ_PORT", "5672"),

		// --- Postgres ---
		POSTGRES_USER:     getEnv("POSTGRES_USER"),
		POSTGRES_PASSWORD: getEnv("POSTGRES_PASSWORD"),
		POSTGRES_HOST:     getEnv("POSTGRES_HOST", "localhost"),
		POSTGRES_PORT:     getEnv("POSTGRES_PORT", "5432"),
		POSTGRES_DB:       getEnv("POSTGRES_DB"),
		POSTGRES_SSLMODE:  getEnv("POSTGRES_SSLMODE", "disable"),
	}

	buildDatabaseURL(&Cfg)
	buildRabbitMQURL(&Cfg)

	return nil
}

func getEnv(key string, defaultValue ...string) string {
	value := os.Getenv(key)
	if value == "" && len(defaultValue) > 0 {
		return defaultValue[0]
	}
	return value
}

// init Database url like "postgres://postgres:password@localhost:8080/mydb?sslmode=false"
func buildDatabaseURL(c *Config) {
	if c.POSTGRES_USER != "" && c.POSTGRES_PASSWORD != "" {
		url := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",
			c.POSTGRES_USER,
			c.POSTGRES_PASSWORD,
			c.POSTGRES_HOST,
			c.POSTGRES_PORT,
			c.POSTGRES_DB,
			c.POSTGRES_SSLMODE,
		)
		c.DATABASE_URL = url
		return
	}

	url := fmt.Sprintf("postgres://%s@%s:%s/%s?sslmode=%s",
		c.POSTGRES_USER,
		c.POSTGRES_HOST,
		c.POSTGRES_PORT,
		c.POSTGRES_DB,
		c.POSTGRES_SSLMODE,
	)
	c.DATABASE_URL = url
}

// init rabbitmq url like "amqp://user:password@localhost:5672/"
func buildRabbitMQURL(c *Config) {
	user := c.RABBITMQ_USER
	password := c.RABBITMQ_PASSWORD
	address := c.RABBITMQ_ADDRESS
	port := c.RABBITMQ_PORT
	url := fmt.Sprintf("amqp://%s:%s@%s:%s/", user, password, address, port)
	c.RABBITMQ_URL = url
}
