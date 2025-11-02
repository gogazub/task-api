package repository

import "github.com/gogazub/consumer/model"

type IResultRepository interface {
	Save(result model.Result)
}

// Only save result
type ResultRepository struct {
	// conn to psql
}

func (repo ResultRepository) Save(result model.Result) {

}
