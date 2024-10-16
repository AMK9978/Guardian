package repository

import (
	"context"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"guardian/configs"
	"guardian/internal/models"
)

type UserRepository struct {
	*MongoBaseRepository[models.User]
}

func NewUserRepository(db *mongo.Database) *UserRepository {
	collection := db.Collection(configs.GlobalConfig.CollectionNames.User)
	return &UserRepository{
		MongoBaseRepository: NewMongoBaseRepository[models.User](collection),
	}
}

type UserTaskRepository struct {
	*MongoBaseRepository[models.UserTask]
}

func NewUserTasksRepository(db *mongo.Database) *UserTaskRepository {
	collection := db.Collection(configs.GlobalConfig.CollectionNames.User)
	return &UserTaskRepository{
		MongoBaseRepository: NewMongoBaseRepository[models.UserTask](collection),
	}
}

func (u *UserTaskRepository) GetUserTasks(ctx context.Context, userID uuid.UUID) ([]models.UserTask, error) {
	var userTasks []models.UserTask
	cursor, err := u.collection.Find(ctx, bson.D{{"user_id", userID}})
	if err != nil {
		return nil, err
	}
	err = cursor.All(ctx, &userTasks)
	return userTasks, err
}
