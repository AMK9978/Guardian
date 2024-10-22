package services

import (
	"context"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"guardian/internal/models"
	"guardian/internal/repository"
)

type UserTaskServiceInterface interface {
	GetUserTasks(userID primitive.ObjectID) ([]models.UserTask, error)
}

type UserTaskService struct {
	userTaskRepo *repository.UserTaskRepository
}

func NewUserTaskService(userTaskRepo *repository.UserTaskRepository) *UserTaskService {
	return &UserTaskService{
		userTaskRepo: userTaskRepo,
	}
}

func (u *UserTaskService) GetUserTasks(userID primitive.ObjectID) ([]models.UserTask, error) {
	return u.userTaskRepo.GetUserTasks(context.Background(), userID)
}