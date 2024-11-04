package mocks

import (
	"context"
	"guardian/internal/models/entities"

	"github.com/stretchr/testify/mock"
)

type MockPluginService struct {
	mock.Mock
}

func (m *MockPluginService) GetPluginsByTask(context context.Context, task entities.Task) ([]entities.Plugin, error) {
	args := m.Called(context, task)
    if plugins, ok := args.Get(0).([]entities.Plugin); ok {
        return plugins, args.Error(1)
    }
    return []entities.Plugin{}, args.Error(1)
}
