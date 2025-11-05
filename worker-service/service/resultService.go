package service

import (
	"context"

	"github.com/gogazub/consumer/model"
	repo "github.com/gogazub/consumer/repository"
)

type ResultService struct {
	repo repo.IResultRepository
}

func NewResultService(repo repo.ResultRepository) *ResultService {
	return &ResultService{repo: repo}
}

func (svc *ResultService) Save(stdout, stderr []byte) {
	resultModel := model.Result{Error: stderr, Output: stdout}
	ctx := context.Background()
	
	svc.repo.Save(ctx, resultModel)
}
