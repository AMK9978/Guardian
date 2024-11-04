package mocks

import (
	"context"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"guardian/internal/models/entities"
)

type MockTargetModelService struct {
	mock.Mock
}


func (m *MockTargetModelService) GetTargetModel(_ context.Context, modelID primitive.ObjectID) (*entities.TargetModel,
	error) {
	args := m.Called(modelID)
	return &entities.TargetModel{}, args.Error(1)
}

func (m *MockTargetModelService) CreateTargetModel(_ context.Context, _ entities.TargetModel) error {
	args := m.Called()
	return args.Error(1)
}


