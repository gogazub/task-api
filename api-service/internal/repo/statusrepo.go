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
