package services

import (
	"context"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"guardian/internal/mocks"
	"testing"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"guardian/internal/models/entities"
)

type MockTargetModelRepo struct {
	mock.Mock
}

var (
	GetTargetModelErr = errors.New("some GetTargetModel error")
)

func TestGetTargetModel(t *testing.T) {
	t.Parallel()

	targetModelRepo := new(mocks.MockTargetModelRepo)
	targetModelService := NewTargetModelService(targetModelRepo)

	t.Run("GetTargetModel returns err", func(t *testing.T) {
		targetModelRepo.On("GetModel", mock.Anything, mock.Anything).Return(nil,
			GetTargetModelErr)
		_, err := targetModelService.GetTargetModel(context.Background(), primitive.NewObjectID())
		assert.Equal(t, GetTargetModelErr, err)
		targetModelRepo.On("GetModel", mock.Anything, mock.Anything).Unset()
	})

	t.Run("GetTargetModel works normal", func(t *testing.T) {
		t.Parallel()

		targetModel := entities.TargetModel{
			ID:       primitive.ObjectID{},
			Provider: "",
			Name:     "",
			Address:  "",
			Status:   0,
			Token:    "",
		}
		targetModelRepo.On("GetModel", mock.Anything, mock.Anything).Return(targetModel, nil)
		result, err := targetModelService.GetTargetModel(context.Background(), primitive.NewObjectID())

		assert.Nil(t, err)
		assert.Equal(t, targetModel, *result)
	})
}

func TestCreateTargetModel(t *testing.T) {
	t.Parallel()

	targetModelRepo := new(mocks.MockTargetModelRepo)
	targetModelService := NewTargetModelService(targetModelRepo)

	t.Run("CreateModel returns err", func(t *testing.T) {
		t.Parallel()

		targetModel := entities.TargetModel{
			Provider: "",
			Name:     "",
			Address:  "",
			Status:   0,
			Token:    "",
		}
		targetModelRepo.On("CreateModel", mock.Anything, mock.Anything).Return(nil,
			GetTargetModelErr)
		err := targetModelService.CreateTargetModel(context.Background(), targetModel)
		assert.Equal(t, GetTargetModelErr, err)
		targetModelRepo.On("CreateModel", mock.Anything, mock.Anything).Unset()
	})

	t.Run("CreateModel works normal", func(t *testing.T) {
		t.Parallel()

		targetModel := entities.TargetModel{
			ID:       primitive.ObjectID{},
			Provider: "",
			Name:     "",
			Address:  "",
			Status:   0,
			Token:    "",
		}
		targetModelRepo.On("CreateModel", mock.Anything, mock.Anything).Return(targetModel, nil)
		err := targetModelService.CreateTargetModel(context.Background(), targetModel)
		assert.Nil(t, err)
	})
}
