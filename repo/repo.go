package repo

import (
	"fmt"
	"sync"
)

type TaskStatus int

const (
	InProgress TaskStatus = iota
	Ready
)

type StatusRepo struct {
	mu    sync.RWMutex
	tasks map[string]TaskStatus
}

func NewStatusRepo() *StatusRepo {
	return &StatusRepo{
		tasks: make(map[string]TaskStatus),
	}
}

func (r *StatusRepo) GetStatus(id string) (TaskStatus, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	status, exists := r.tasks[id]
	return status, exists
}

func (r *StatusRepo) AddTask(id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.tasks[id]; exists {
		return fmt.Errorf("task already exists: %s", id)
	}

	r.tasks[id] = InProgress
	return nil
}

func (r *StatusRepo) UpdateStatus(id string, status TaskStatus) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.tasks[id]; !exists {
		return fmt.Errorf("task not found: %s", id)
	}

	r.tasks[id] = status
	return nil
}

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

	token := fmt.Sprintf("token_%s", username)

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
