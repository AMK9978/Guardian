package repository

import (
	"context"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"guardian/configs"
	"guardian/internal/models/entities"
)

type TargetModelRepository struct {
	*MongoBaseRepository[entities.TargetModel]
}

func NewTargetModelRepository(db *mongo.Database) *TargetModelRepository {
	collection := db.Collection(configs.GlobalConfig.CollectionNames.TargetModel)
	return &TargetModelRepository{
		MongoBaseRepository: NewMongoBaseRepository[entities.TargetModel](collection),
	}
}

func (u *TargetModelRepository) GetModels(ctx context.Context, modelIDs []primitive.ObjectID) ([]entities.TargetModel,
	error) {
	var models []entities.TargetModel

	filter := bson.M{"_id": bson.M{"$in": modelIDs}}
	cursor, err := u.collection.Find(ctx, filter)
	if err != nil {
		return nil, errors.Errorf("error in GetTasks: %v", err)
	}
	err = cursor.All(ctx, &models)
	if err != nil {
		return nil, errors.Errorf("error in fetching tasks: %v", err)
	}
	return models, nil
}

func (u *TargetModelRepository) GetModel(ctx context.Context, modelID primitive.ObjectID) (*entities.TargetModel,
	error) {
	var model entities.TargetModel
	err := u.collection.FindOne(ctx, bson.D{{"_id", modelID}}).Decode(&model)
	if err != nil {
		return nil, err
	}
	return &model, err
}

func (u *TargetModelRepository) CreateModel(ctx context.Context, model entities.TargetModel) (interface{}, error) {
	cursor, err := u.collection.InsertOne(ctx, bson.D{{"name", model.Name},
		{"status", model.Status}, {"address", model.Address},
		{"provider", model.Provider}, {"token", model.Token}})
	if err != nil {
		return nil, err
	}
	return cursor.InsertedID, nil
}

func (u *TargetModelRepository) DeleteModel(ctx context.Context, modelID primitive.ObjectID) (int64, error) {
	cursor, err := u.collection.DeleteOne(ctx, bson.D{{"_id", modelID}})
	if err != nil {
		return -1, err
	}
	return cursor.DeletedCount, err
}

func (u *TargetModelRepository) UpdateModel(ctx context.Context, model entities.TargetModel) (int64, error) {
	cursor, err := u.collection.UpdateByID(ctx, model.ID, model)
	if err != nil {
		return -1, err
	}
	return cursor.ModifiedCount, err
}
