package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"

	"github.com/gogazub/consumer/model"
)

// IResultRepository - save result of the execution to DB.
type IResultRepository interface {
	Save(ctx context.Context,result model.Result) error
}

// ResultRepository - save result of the execution to DB.
type ResultRepository struct {
	db *sql.DB
}

func NewOrderRepository(db *sql.DB) *ResultRepository {
	return &ResultRepository{db: db}
}

func (r ResultRepository) Save(ctx context.Context,result model.Result) error{
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

	_, err = r.db.ExecContext(ctx,  "INSERT INTO results (stderr, stdout) VALUES ($1, $2) RETURNING id",result.Error,result.Output)
	if err != nil {
		return fmt.Errorf("save result error: %w",err.Error())
	}
	return nil
}
