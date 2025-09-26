package service

import "github.com/gogazub/app/internal/repo"

type Service struct {
	taskService *TaskService
	userService *UserService
}

func NewService(taskRepo *repo.StatusRepo, userRepo *repo.UserRepo) *Service {
	return &Service{
		taskService: NewTaskService(taskRepo),
		userService: NewUserService(userRepo),
	}
}

func (s *Service) GetTaskService() *TaskService {
	return s.taskService
}

func (s *Service) GetUserService() *UserService {
	return s.userService
}
