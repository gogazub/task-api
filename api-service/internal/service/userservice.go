package service

import (
	"errors"

	"github.com/gogazub/app/internal/repo"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	users *repo.UserRepo
	m     JWTManager
}

func NewUserService(users *repo.UserRepo, manager *JWTManager) *UserService {
	return &UserService{
		users: users,
		m:     *manager,
	}
}

func (s *UserService) Register(username, password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	token, err := s.m.GenerateToken(username)
	if err != nil {
		return "", err
	}

	err = s.users.Save(username, string(hash))
	if err != nil {
		return "", err
	}

	return token, nil
}

func (s *UserService) Login(username, password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	storedHash, err := s.users.GetHash(username)
	if err != nil {
		return "", err
	}

	if string(hash) != storedHash {
		return "", errors.New("invalid credentials")
	}

	token, err := s.m.GenerateToken(username)
	if err != nil {
		return "", err
	}

	return token, nil

}

func (s *UserService) FindUser(username string) bool {
	return s.users.Find(username)
}
