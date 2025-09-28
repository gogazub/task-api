package service

import (
	"github.com/gogazub/app/internal/repo"
)

type Service struct {
	taskService *TaskService
	userService *UserService
}

func NewService(taskRepo *repo.StatusRepo, userService *UserService) *Service {
	return &Service{
		taskService: NewTaskService(taskRepo),
		userService: userService,
	}
}

func (s *Service) GetTaskService() *TaskService {
	return s.taskService
}

func (s *Service) GetUserService() *UserService {
	return s.userService
}

func (s *Service) GetUsernameFromToken(tokenString string) (string, error) {
	name, err := s.GetUserService().m.GetName(tokenString)
	return name, err
}
