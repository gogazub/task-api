package repo

import (
	"fmt"
	"sync"
)

type UserRepo struct {
	mu    sync.RWMutex
	users map[string]string // username:hash
}

func NewUserRepo() *UserRepo {
	return &UserRepo{
		users: make(map[string]string),
	}
}

func (r *UserRepo) Save(username, hash string) error {

	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.users[username]; exists {
		return fmt.Errorf("user already exists")
	}

	r.users[username] = hash
	return nil
}

func (r *UserRepo) GetHash(username string) (string, error) {
	r.mu.RLock()
	storedHash, exists := r.users[username]
	r.mu.RUnlock()
	if !exists {
		return storedHash, fmt.Errorf("user: %s doesn`t exists", username)
	}
	return storedHash, nil
}

func (r *UserRepo) Find(username string) bool {
	_, exists := r.users[username]
	return exists
}
