package repository

import (
	"context"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
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

func (u *UserRepository) GetUser(ctx context.Context, userID uuid.UUID) (models.User, error) {
	var user models.User
	cursor, err := u.collection.Find(ctx, bson.D{{"user_id", userID}})
	if err != nil {
		return models.User{}, err
	}
	err = cursor.All(ctx, &user)
	return user, err
}

func (u *UserRepository) CreateUser(ctx context.Context, user models.User) (interface{}, error) {
	cursor, err := u.collection.InsertOne(ctx, bson.D{{"user_id", user.ID}, {"name", user.Name},
		{"status", 1}})
	if err != nil {
		return nil, err
	}
	return cursor.InsertedID, nil
}

func (u *UserRepository) DeleteUser(ctx context.Context, userID uuid.UUID) (int64, error) {
	cursor, err := u.collection.DeleteOne(ctx, bson.D{{"user_id", userID}})
	if err != nil {
		return -1, err
	}
	return cursor.DeletedCount, err
}

func (u *UserRepository) UpdateUser(ctx context.Context, user models.User) (int64, error) {
	cursor, err := u.collection.UpdateByID(ctx, user.ID, user)
	if err != nil {
		return -1, err
	}
	return cursor.ModifiedCount, err
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

func (u *UserTaskRepository) GetUserTasks(ctx context.Context, userID primitive.ObjectID) ([]models.UserTask, error) {
	var userTasks []models.UserTask
	cursor, err := u.collection.Find(ctx, bson.D{{"user_id", userID}})
	if err != nil {
		return nil, err
	}
	err = cursor.All(ctx, &userTasks)
	return userTasks, err
}
