package services

import (
	"context"
	"github.com/google/uuid"
	"guardian/internal/models"
	"guardian/internal/repository"
)

type UserTaskServiceInterface interface {
	GetUserTasks(userID uuid.UUID) ([]models.UserTask, error)
}

type UserTaskService struct {
	userTaskRepo *repository.UserTaskRepository
}

func NewUserTaskService(userTaskRepo *repository.UserTaskRepository) *UserTaskService {
	return &UserTaskService{
		userTaskRepo: userTaskRepo,
	}
}

func (u *UserTaskService) GetUserTasks(userID uuid.UUID) ([]models.UserTask, error) {
	return u.userTaskRepo.GetUserTasks(context.Background(), userID)
}