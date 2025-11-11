package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/gogazub/consumer/config"
	"github.com/gogazub/consumer/model"
)

// IResultRepository - save result of the execution to DB.
type IResultRepository interface {
	Save(ctx context.Context, result model.ExecutionResult) error
}

// ResultRepository - save result of the execution to DB.
type ResultRepository struct {
	db *sql.DB
}

func NewOrderRepository(db *sql.DB) *ResultRepository {
	return &ResultRepository{db: db}
}

func (r ResultRepository) Save(ctx context.Context, result model.ExecutionResult) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	defer func() {
		err := tx.Rollback()

		if err != nil && !errors.Is(err, sql.ErrTxDone) {
			log.Printf("Rollback error:%s", err.Error())
		}
	}()

	_, err = r.db.ExecContext(ctx, "INSERT INTO results (id, status, output, error, exit_code, duration) VALUES ($1, $2, $3, $4, $5, $6)", result.TaskID, result.Status, result.Output, result.Error, result.ExitCode, result.Duration)
	if err != nil {
		return fmt.Errorf("save result error: %v", err.Error())
	}
	return nil
}

func ConnectToDB() *sql.DB {
	connStr := config.Cfg.DATABASE_URL

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
