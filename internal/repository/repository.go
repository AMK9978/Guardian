package repository

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type MongoBaseRepository[T any] struct {
	collection *mongo.Collection
}

func NewMongoBaseRepository[T any](collection *mongo.Collection) *MongoBaseRepository[T] {
	return &MongoBaseRepository[T]{
		collection: collection,
	}
}

func (r *MongoBaseRepository[T]) Create(ctx context.Context, entity *T) error {
	_, err := r.collection.InsertOne(ctx, entity)
	return err
}

func (r *MongoBaseRepository[T]) Update(ctx context.Context, filter bson.M, update bson.M) error {
	_, err := r.collection.UpdateOne(ctx, filter, update)
	return err
}

func (r *MongoBaseRepository[T]) Delete(ctx context.Context, filter bson.M) error {
	_, err := r.collection.DeleteOne(ctx, filter)
	return err
}

func (r *MongoBaseRepository[T]) GetByFilter(ctx context.Context, filter bson.M) (*T, error) {
	var entity T
	err := r.collection.FindOne(ctx, filter).Decode(&entity)
	return &entity, err
}

func (r *MongoBaseRepository[T]) GetAll(ctx context.Context, filter bson.M) ([]T, error) {
	var entities []T
	cursor, err := r.collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	err = cursor.All(ctx, &entities)
	return entities, err
}
