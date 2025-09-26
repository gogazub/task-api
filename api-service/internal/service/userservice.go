package service

import "github.com/gogazub/app/internal/repo"

type UserService struct {
	repo *repo.UserRepo
}

func NewUserService(repo *repo.UserRepo) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) Register(username, password string) error {
	return s.repo.Register(username, password)
}

func (s *UserService) Login(username, password string) (string, error) {
	return s.repo.Login(username, password)
}

func (s *UserService) ValidateToken(token string) (string, bool) {
	return s.repo.ValidateToken(token)
}
