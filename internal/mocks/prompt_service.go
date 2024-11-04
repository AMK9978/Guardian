package mocks

import (
	"context"
	"github.com/stretchr/testify/mock"
	"guardian/internal/models"
	"net/http"
)

type MockPromptService struct {
	mock.Mock
}


func (p *MockPromptService) ProcessPrompt(_ context.Context, _ *models.PluginRequest) (bool, error) {
	args := p.Called()
	return args.Bool(0), args.Error(1)
}

func (p *MockPromptService) SendPrompt(ctx context.Context, newReq *http.Request) (*http.Response, error) {
	args := p.Called()
	return args.Get(0).(*http.Response), args.Error(1)
}