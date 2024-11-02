package mocks

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"

	"github.com/stretchr/testify/mock"
)

type MockMiddleware struct {
	mock.Mock
}

func (m *MockMiddleware) GetUserFromContext(_ *http.Request) (*primitive.ObjectID, error) {
	args := m.Called()
	return new(primitive.ObjectID), args.Error(1)
}
