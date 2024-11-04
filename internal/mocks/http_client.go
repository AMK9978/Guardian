package mocks

import (
	"context"
	"guardian/internal/models"
	"net/http"

	"github.com/stretchr/testify/mock"
)

type MockClient struct {
	mock.Mock
}

func (m *MockClient) Do(req *http.Request) (*http.Response, error) {
	args := m.Called(req)
	return args.Get(0).(*http.Response), args.Error(1)
}

func (m *MockClient) Forward(ctx context.Context, req *models.PluginRequest) (*models.PluginResponse, error) {
	args := m.Called(ctx, req)
	return args.Get(0).(*models.PluginResponse), args.Error(1)
}