package services

import (
	"guardian/internal/repository"
)

type UserService struct {
	userRepo *repository.UserRepository
}

func NewUserRepo(userRepo *repository.UserRepository) *UserService {
	return &UserService{
		userRepo: userRepo,
	}
}
