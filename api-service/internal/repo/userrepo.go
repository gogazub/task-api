package repo

import (
	"fmt"
	"sync"

	"github.com/gogazub/app/internal/utils"
)

type UserRepo struct {
	mu     sync.RWMutex
	users  map[string]string // username:password
	tokens map[string]string // token:username
}

func NewUserRepo() *UserRepo {
	return &UserRepo{
		users:  make(map[string]string),
		tokens: make(map[string]string),
	}
}

func (r *UserRepo) Register(username, password string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.users[username]; exists {
		return fmt.Errorf("user already exists")
	}

	r.users[username] = password
	return nil
}

func (r *UserRepo) Login(username, password string) (string, error) {
	r.mu.RLock()
	storedPassword, exists := r.users[username]
	r.mu.RUnlock()

	if !exists || storedPassword != password {
		return "", fmt.Errorf("invalid credentials")
	}

	token := utils.GenerateUUID()

	r.mu.Lock()
	r.tokens[token] = username
	r.mu.Unlock()

	return token, nil
}

func (r *UserRepo) ValidateToken(token string) (string, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	username, exists := r.tokens[token]
	return username, exists
}
