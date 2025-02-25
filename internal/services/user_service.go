package services

import (
	"context"
	"time"

	"guardian/configs"
	"guardian/internal/models"
	"guardian/internal/models/entities"
	"guardian/internal/repository"
	"guardian/utlis/logger"

	"github.com/golang-jwt/jwt/v4"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

type UserServiceInterface interface {
	GetUserTasksByID(userID primitive.ObjectID) ([]entities.Task, error)
	GetUser(id primitive.ObjectID) (*entities.User, error)
	Login(req models.LoginRequest) (string, error)
	SignUp(req models.SignUpRequest) error
	ActivateUser(req models.SignUpRequest) error
}

type UserService struct {
	userRepo *repository.UserRepository
	taskRepo *repository.TaskRepository
}

func NewUserService(userRepo *repository.UserRepository, taskRepo *repository.TaskRepository) *UserService {
	return &UserService{
		userRepo: userRepo,
		taskRepo: taskRepo,
	}
}

func (u *UserService) GetUser(id primitive.ObjectID) (*entities.User, error) {
	user, err := u.userRepo.GetByFilter(context.Background(), bson.M{"_id": id})
	if err != nil {
		return nil, errors.Errorf("error in GetUser:%v", err)
	}
	return user, nil
}

func (u *UserService) GetUserTasksByID(userID primitive.ObjectID) ([]entities.Task, error) {
	user, err := u.GetUser(userID)
	if err != nil {
		return nil, errors.Errorf("user error:%v", userID)
	}
	if user.Tasks == nil {
		return []entities.Task{}, nil
	}
	tasks, err := u.taskRepo.GetTasks(context.Background(), user.Tasks)
	return tasks, err
}

func (u *UserService) Login(req models.LoginRequest) (string, error) {
	user, err := u.userRepo.GetByFilter(context.Background(), bson.M{"email": req.Email})
	if err != nil {
		return "", err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if err != nil {
		return "", err
	}

	_, tokenString, err := configs.GlobalConfig.TokenAuth.Encode(map[string]interface{}{
		"user_id": user.ID,
		"exp":     time.Now().Add(configs.GlobalConfig.TokenExpirationTime),
	})
	if err != nil {
		logger.GetLogger().Errorf("error generating token: %v", err)
		return "", err
	}
	return tokenString, nil
}

func (u *UserService) SignUp(req models.SignUpRequest) error {
	hashedPassword, err := hashPassword(req.Password)
	if err != nil {
		return err
	}

	user := entities.User{
		Name:     req.Name,
		Email:    req.Email,
		Password: hashedPassword,
		Status:   0,
		Groups:   nil,
	}

	err = u.userRepo.Create(context.Background(), &user)
	if err != nil {
		return err
	}

	return nil
}

func (u *UserService) ActivateUser(_ models.SignUpRequest) error {
	_, _ = generateActivationToken("")
	_ = sendActivationEmail("", "")
	return nil
}

func hashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

// TODO: Infra layer code.
func generateActivationToken(userID string) (string, error) {
	expTime := time.Now().Add(configs.GlobalConfig.ActivationTokenExpTime)

	claims := jwt.MapClaims{}
	claims["user_id"] = userID
	claims["exp"] = expTime

	token := jwt.NewWithClaims(jwt.SigningMethodES256, claims)

	activationToken, err := token.SignedString([]byte(configs.GlobalConfig.ActivationTokenKey))
	if err != nil {
		return "", err
	}
	return activationToken, nil
}

func sendActivationEmail(_ string, _ string) error {
	return nil
}
