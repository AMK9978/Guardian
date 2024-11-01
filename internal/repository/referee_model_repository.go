package repository

import (
	"context"

	"guardian/configs"
	"guardian/internal/models/entities"

	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type RefereeModelRepository struct {
	*MongoBaseRepository[entities.RefereeModel]
}

func NewRefereeModelRepository(db *mongo.Database) *RefereeModelRepository {
	collection := db.Collection(configs.GlobalConfig.CollectionNames.RefereeModel)
	return &RefereeModelRepository{
		MongoBaseRepository: NewMongoBaseRepository[entities.RefereeModel](collection),
	}
}

func (u *TaskRepository) GetRefereeModels(ctx context.Context, modelIDs []primitive.ObjectID) ([]entities.RefereeModel,
	error,
) {
	var models []entities.RefereeModel

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

func (u *TaskRepository) GetRefereeModel(ctx context.Context, modelID primitive.ObjectID) (entities.RefereeModel,
	error,
) {
	var model entities.RefereeModel
	cursor, err := u.collection.Find(ctx, bson.D{{"_id", modelID}})
	if err != nil {
		return entities.RefereeModel{}, err
	}
	err = cursor.All(ctx, &model)
	return model, err
}

func (u *TaskRepository) CreateRefereeModel(ctx context.Context, model entities.RefereeModel) (interface{}, error) {
	cursor, err := u.collection.InsertOne(ctx, bson.D{
		{"name", model.Name},
		{"status", model.Status},
		{"address", model.Address},
		{"provider", model.Provider},
		{"token", model.Token},
	})
	if err != nil {
		return nil, err
	}
	return cursor.InsertedID, nil
}

func (u *TaskRepository) DeleteRefereeModel(ctx context.Context, modelID primitive.ObjectID) (int64, error) {
	cursor, err := u.collection.DeleteOne(ctx, bson.D{{"_id", modelID}})
	if err != nil {
		return -1, err
	}
	return cursor.DeletedCount, err
}

func (u *TaskRepository) UpdateRefereeModel(ctx context.Context, model entities.RefereeModel) (int64, error) {
	cursor, err := u.collection.UpdateByID(ctx, model.ID, model)
	if err != nil {
		return -1, err
	}
	return cursor.ModifiedCount, err
}
