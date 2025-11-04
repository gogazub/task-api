package repository

import "github.com/gogazub/consumer/model"

// IResultRepository - save result of the execution to DB.
type IResultRepository interface {
	Save(result model.Result)
}

// ResultRepository - save result of the execution to DB.
type ResultRepository struct {
	// conn to psql
}

func (repo ResultRepository) Save(result model.Result) {

}
