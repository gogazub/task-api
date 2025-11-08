package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	// --- Server ---

	// --- RabbitMQ ---
	RABBITMQ_USER     string
	RABBITMQ_PASSWORD string
	RABBITMQ_ADDRESS  string
	RABBITMQ_PORT     string

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
func InitConfig() (*Config, error) {
	if err := godotenv.Load(); err != nil {
		return nil, fmt.Errorf("init config error: %w", err)
	}

	config := &Config{
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

	config.DATABASE_URL = buildDatabaseURL(config)

	return config, nil
}

func getEnv(key string, defaultValue ...string) string {
	value := os.Getenv(key)
	if value == "" && len(defaultValue) > 0 {
		return defaultValue[0]
	}
	return value
}

func buildDatabaseURL(c *Config) string {
	if c.POSTGRES_USER != "" && c.POSTGRES_PASSWORD != "" {
		return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",
			c.POSTGRES_USER,
			c.POSTGRES_PASSWORD,
			c.POSTGRES_HOST,
			c.POSTGRES_PORT,
			c.POSTGRES_DB,
			c.POSTGRES_SSLMODE,
		)
	}

	return fmt.Sprintf("postgres://%s@%s:%s/%s?sslmode=%s",
		c.POSTGRES_USER,
		c.POSTGRES_HOST,
		c.POSTGRES_PORT,
		c.POSTGRES_DB,
		c.POSTGRES_SSLMODE,
	)
}
