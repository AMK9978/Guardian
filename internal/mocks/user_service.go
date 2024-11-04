package mocks

import (
	"guardian/internal/models"
	"guardian/internal/models/entities"

	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type MockUserService struct {
	mock.Mock
}

func (m *MockUserService) Login(req models.LoginRequest) (string, error) {
	args := m.Called(req)
	return args.String(0), args.Error(1)
}

func (m *MockUserService) SignUp(req models.SignUpRequest) error {
	return m.Called(req).Error(0)
}

func (m *MockUserService) GetUserTasksByID(userID primitive.ObjectID) ([]entities.Task, error) {
	args := m.Called(userID)
	return args.Get(0).([]entities.Task), args.Error(1)
}

func (m *MockUserService) GetUser(_ primitive.ObjectID) (*entities.User, error) {
	return nil, nil
}

func (m *MockUserService) ActivateUser(_ models.SignUpRequest) error {
	return nil
}
