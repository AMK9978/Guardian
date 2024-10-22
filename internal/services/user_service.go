package services

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"guardian/configs"
	"guardian/internal/models"
	"guardian/internal/repository"

	"golang.org/x/crypto/bcrypt"
	"guardian/utlis/logger"
	"time"
)

type UserServiceInterface interface {
	Login(req models.LoginRequest) (string, error)
	SignUp(req models.SignUpRequest) error
}

type UserService struct {
	userRepo *repository.UserRepository
}

func NewUserService(userRepo *repository.UserRepository) *UserService {
	return &UserService{
		userRepo: userRepo,
	}
}

func (u *UserService) Login(req models.LoginRequest) (string, error) {

	user, err := u.userRepo.GetByID(context.Background(), bson.M{"user_id": req.UserID})
	if err != nil {
		return "", err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if err != nil {
		return "", err
	}

	_, tokenString, err := configs.GlobalConfig.TokenAuth.Encode(map[string]interface{}{
		"user_id": req.UserID,
		"exp":     time.Now().Add(configs.GlobalConfig.TokenExpirationTime),
	})
	if err != nil {
		logger.GetLogger().Errorf("error generating token: %v", err)
		return "", err
	}
	return tokenString, nil
}

func (u *UserService) SignUp(req models.SignUpRequest) error {
    hashedPassword, err := hashPassword(req.User.Password)
    if err != nil {
        return err
    }

    user := models.User{
        Password: hashedPassword,

    }

    err = u.userRepo.Create(context.Background(), &user)
    if err != nil {
        return err
    }

    return nil
}

func hashPassword(password string) (string, error) {
    hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
    if err != nil {
        return "", err
    }
    return string(hashedPassword), nil
}