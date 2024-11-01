package services

import (
	"context"

	"guardian/internal/models/entities"
	"guardian/internal/repository"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type TargetModelService struct {
	targetModelRepo *repository.TargetModelRepository
}

func NewTargetModelService(targetModelRepo *repository.TargetModelRepository) *TargetModelService {
	return &TargetModelService{
		targetModelRepo: targetModelRepo,
	}
}

func (t *TargetModelService) GetTargetModel(modelID primitive.ObjectID) (*entities.TargetModel, error) {
	targetModel, err := t.targetModelRepo.GetModel(context.Background(), modelID)
	if err != nil {
		return nil, err
	}
	return targetModel, err
}

func (t *TargetModelService) CreateTargetModel(model entities.TargetModel) error {
	_, err := t.targetModelRepo.CreateModel(context.Background(), model)
	if err != nil {
		return err
	}
	return nil
}
