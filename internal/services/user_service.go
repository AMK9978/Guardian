package services

import (
	"context"
	"github.com/dgrijalva/jwt-go"
	"go.mongodb.org/mongo-driver/bson"
	"guardian/configs"
	"guardian/internal/models"
	"guardian/internal/models/entities"
	"guardian/internal/repository"

	"golang.org/x/crypto/bcrypt"
	"guardian/utlis/logger"
	"time"
)

type UserServiceInterface interface {
	Login(req models.LoginRequest) (string, error)
	SignUp(req models.SignUpRequest) error
	ActivateUser(req models.SignUpRequest) error
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

func (u *UserService) ActivateUser(req models.SignUpRequest) error {
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

func sendActivationEmail(email string, activationLink string) error {
	return nil
}