package repository

import (
	"context"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"guardian/configs"
	"guardian/internal/models/entities"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"

)

type UserRepository struct {
	*MongoBaseRepository[entities.User]
}

func NewUserRepository(db *mongo.Database) *UserRepository {
	collection := db.Collection(configs.GlobalConfig.CollectionNames.User)
	return &UserRepository{
		MongoBaseRepository: NewMongoBaseRepository[entities.User](collection),
	}
}

func (u *UserRepository) GetUser(ctx context.Context, userID primitive.ObjectID) (*entities.User, error) {
	var user entities.User
	err := u.collection.FindOne(ctx, bson.D{{"_id", userID}}).Decode(&user)
	if err != nil {
		return nil, err
	}
	return &user, err
}

func (u *UserRepository) CreateUser(ctx context.Context, user entities.User) (interface{}, error) {
	cursor, err := u.collection.InsertOne(ctx, bson.D{{"name", user.Name}, {"status", 1}})
	if err != nil {
		return nil, err
	}
	return cursor.InsertedID, nil
}

func (u *UserRepository) DeleteUser(ctx context.Context, userID uuid.UUID) (int64, error) {
	cursor, err := u.collection.DeleteOne(ctx, bson.D{{"_id", userID}})
	if err != nil {
		return -1, err
	}
	return cursor.DeletedCount, err
}

func (u *UserRepository) UpdateUser(ctx context.Context, user entities.User) (int64, error) {
	cursor, err := u.collection.UpdateByID(ctx, user.ID, user)
	if err != nil {
		return -1, err
	}
	return cursor.ModifiedCount, err
}
