package mocks

import (
	"context"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"guardian/internal/models/entities"
)

type MockTargetModelRepo struct {
	mock.Mock
}

func (m *MockTargetModelRepo) GetModels(_ context.Context, modelIDs []primitive.ObjectID) ([]entities.TargetModel,
	error,
) {
	args := m.Called(modelIDs)
	return []entities.TargetModel{}, args.Error(1)
}

func (m *MockTargetModelRepo) GetModel(_ context.Context, modelID primitive.ObjectID) (entities.TargetModel,
	error,
) {
	args := m.Called(modelID)
	return entities.TargetModel{}, args.Error(1)
}

func (m *MockTargetModelRepo) CreateModel(_ context.Context, model entities.TargetModel) (interface{}, error) {
	args := m.Called(model)
	return nil, args.Error(1)
}

func (m *MockTargetModelRepo) DeleteModel(_ context.Context, modelID primitive.ObjectID) (int64, error) {
	args := m.Called(modelID)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockTargetModelRepo) UpdateModel(_ context.Context, model entities.TargetModel) (int64, error) {
	args := m.Called(model)
	return args.Get(0).(int64), args.Error(1)
}


