package service

import (
	"github.com/gogazub/app/internal/repo"
)

type Service struct {
	taskService *TaskService
	userService *UserService
	producer    *Producer
}

func NewService(taskRepo *repo.StatusRepo, userService *UserService, producer *Producer) *Service {
	return &Service{
		taskService: NewTaskService(taskRepo),
		userService: userService,
		producer:    producer,
	}
}

func (s *Service) GetTaskService() *TaskService {
	return s.taskService
}

func (s *Service) GetUserService() *UserService {
	return s.userService
}

func (s *Service) GetProducer() *Producer {
	return s.producer
}

func (s *Service) GetUsernameFromToken(tokenString string) (string, error) {
	name, err := s.GetUserService().m.GetName(tokenString)
	return name, err
}
