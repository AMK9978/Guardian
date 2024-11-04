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

type PluginRepository struct {
	*MongoBaseRepository[entities.Plugin]
}

type PluginRepoInterface interface {
	GetPluginsByTask(ctx context.Context, task entities.Task) ([]entities.Plugin, error)
	GetPlugins(ctx context.Context, modelIDs []primitive.ObjectID) ([]entities.Plugin, error)
	GetPlugin(ctx context.Context, modelID primitive.ObjectID) (entities.Plugin, error)
	CreatePlugin(ctx context.Context, model entities.Plugin) (interface{}, error)
	DeletePlugin(ctx context.Context, modelID primitive.ObjectID) (int64, error)
	UpdatePlugin(ctx context.Context, model entities.Plugin) (int64, error)
}

func NewPluginRepository(db *mongo.Database) *PluginRepository {
	collection := db.Collection(configs.GlobalConfig.CollectionNames.Plugin)
	return &PluginRepository{
		MongoBaseRepository: NewMongoBaseRepository[entities.Plugin](collection),
	}
}

func (u *PluginRepository) GetPluginsByTask(ctx context.Context, task entities.Task) ([]entities.Plugin, error) {
	var plugins []entities.Plugin

	filter := bson.M{"_id": bson.M{"$in": task.Plugins}}
	cursor, err := u.collection.Find(ctx, filter)
	if err != nil {
		return nil, errors.Errorf("error in GetPlugins: %v", err)
	}
	err = cursor.All(ctx, &plugins)
	if err != nil {
		return nil, errors.Errorf("error in fetching plugins: %v", err)
	}
	return plugins, nil
}

func (u *PluginRepository) GetPlugins(ctx context.Context, modelIDs []primitive.ObjectID) ([]entities.Plugin,
	error,
) {
	var models []entities.Plugin

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

func (u *PluginRepository) GetPlugin(ctx context.Context, modelID primitive.ObjectID) (entities.Plugin,
	error,
) {
	var model entities.Plugin
	cursor, err := u.collection.Find(ctx, bson.D{{"_id", modelID}})
	if err != nil {
		return entities.Plugin{}, err
	}
	err = cursor.All(ctx, &model)
	return model, err
}

func (u *PluginRepository) CreatePlugin(ctx context.Context, model entities.Plugin) (interface{}, error) {
	cursor, err := u.collection.InsertOne(ctx, bson.D{
		{"name", model.Name},
		{"status", model.Status},
		{"address", model.Address},
		{"provider", model.Provider},
		{"token", model.Token},
		{"protocol", model.Protocol},
	})
	if err != nil {
		return nil, err
	}
	return cursor.InsertedID, nil
}

func (u *PluginRepository) DeletePlugin(ctx context.Context, modelID primitive.ObjectID) (int64, error) {
	cursor, err := u.collection.DeleteOne(ctx, bson.D{{"_id", modelID}})
	if err != nil {
		return -1, err
	}
	return cursor.DeletedCount, err
}

func (u *PluginRepository) UpdatePlugin(ctx context.Context, model entities.Plugin) (int64, error) {
	cursor, err := u.collection.UpdateByID(ctx, model.ID, model)
	if err != nil {
		return -1, err
	}
	return cursor.ModifiedCount, err
}
