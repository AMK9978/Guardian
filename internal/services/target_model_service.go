package services

import (
	"context"

	"guardian/internal/models/entities"
	"guardian/internal/repository"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type TargetModelServiceInterface interface {
	GetTargetModel(ctx context.Context, modelID primitive.ObjectID) (entities.TargetModel, error)
	CreateTargetModel(ctx context.Context, model entities.TargetModel) error
}

type TargetModelService struct {
	targetModelRepo *repository.TargetModelRepository
}

func NewTargetModelService(targetModelRepo *repository.TargetModelRepository) *TargetModelService {
	return &TargetModelService{
		targetModelRepo: targetModelRepo,
	}
}

func (t *TargetModelService) GetTargetModel(ctx context.Context, modelID primitive.ObjectID) (entities.TargetModel,
	error) {
	targetModel, err := t.targetModelRepo.GetModel(ctx, modelID)
	if err != nil {
		return entities.TargetModel{}, err
	}
	return *targetModel, err
}

func (t *TargetModelService) CreateTargetModel(ctx context.Context, model entities.TargetModel) error {
	_, err := t.targetModelRepo.CreateModel(ctx, model)
	if err != nil {
		return err
	}
	return nil
}
