package service

import (
	"fmt"
	"time"

	"github.com/gogazub/hw1/repo"
	"github.com/gogazub/hw1/utils"
)

type TaskService struct {
	repo *repo.StatusRepo
}

func NewTaskService(repo *repo.StatusRepo) *TaskService {
	return &TaskService{repo: repo}
}

func (s *TaskService) CreateTask() (string, error) {
	id := utils.GenerateUUID()

	err := s.repo.AddTask(id)
	if err != nil {
		return "", fmt.Errorf("failed to create task: %w", err)
	}

	go s.processTask(id)

	return id, nil
}

func (s *TaskService) GetTaskStatus(id string) (repo.TaskStatus, error) {
	status, exists := s.repo.GetStatus(id)
	if !exists {
		return repo.InProgress, fmt.Errorf("task not found: %s", id)
	}
	return status, nil
}

func (s *TaskService) processTask(id string) {
	time.Sleep(5 * time.Second)

	err := s.repo.UpdateStatus(id, repo.Ready)
	if err != nil {
		fmt.Printf("Error updating task status: %v\n", err)
	}

	fmt.Printf("Task %s is ready\n", id)
}

func (s *TaskService) GetTaskResult(id string) (string, error) {
	status, exists := s.repo.GetStatus(id)
	if !exists {
		return "", fmt.Errorf("task not found: %s", id)
	}

	if status != repo.Ready {
		return "", fmt.Errorf("task not ready yet: %s", id)
	}

	return fmt.Sprintf("Result for task %s", id), nil
}

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
