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

type TaskRepository struct {
	*MongoBaseRepository[entities.Task]
}

func NewTaskRepository(db *mongo.Database) *TaskRepository {
	collection := db.Collection(configs.GlobalConfig.CollectionNames.Task)
	return &TaskRepository{
		MongoBaseRepository: NewMongoBaseRepository[entities.Task](collection),
	}
}

func (u *TaskRepository) GetTasks(ctx context.Context, taskIDs []primitive.ObjectID) ([]entities.Task, error) {
	var tasks []entities.Task

	filter := bson.M{"_id": bson.M{"$in": taskIDs}}
	cursor, err := u.collection.Find(ctx, filter)
	if err != nil {
		return nil, errors.Errorf("error in GetTasks: %v", err)
	}
	err = cursor.All(ctx, &tasks)
	if err != nil {
		return nil, errors.Errorf("error in fetching tasks: %v", err)
	}
	return tasks, nil
}

func (u *TaskRepository) GetTask(ctx context.Context, taskID primitive.ObjectID) (entities.Task, error) {
	var task entities.Task
	cursor, err := u.collection.Find(ctx, bson.D{{"_id", taskID}})
	if err != nil {
		return entities.Task{}, err
	}
	err = cursor.All(ctx, &task)
	return task, err
}

func (u *TaskRepository) CreateTask(ctx context.Context, task entities.Task) (interface{}, error) {
	cursor, err := u.collection.InsertOne(ctx, bson.D{
		{"type", task.Type},
		{"status", task.Status},
		{"address", task.Address},
	})
	if err != nil {
		return nil, err
	}
	return cursor.InsertedID, nil
}

func (u *TaskRepository) DeleteTask(ctx context.Context, taskID primitive.ObjectID) (int64, error) {
	cursor, err := u.collection.DeleteOne(ctx, bson.D{{"_id", taskID}})
	if err != nil {
		return -1, err
	}
	return cursor.DeletedCount, err
}

func (u *TaskRepository) UpdateTask(ctx context.Context, task entities.Task) (int64, error) {
	cursor, err := u.collection.UpdateByID(ctx, task.ID, task)
	if err != nil {
		return -1, err
	}
	return cursor.ModifiedCount, err
}
