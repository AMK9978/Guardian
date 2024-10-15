package mongodb

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Client struct {
	MongoClient *mongo.Client
}

func NewClient(mongoDBURI string) *Client {
	client, err := mongo.NewClient(options.Client().ApplyURI(mongoDBURI))
	if err != nil {
		log.Fatalf("Failed to create MongoDB client: %s", err)
	}

	// Connect to MongoDB
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err = client.Connect(ctx); err != nil {
		log.Fatalf("Failed to connect to MongoDB: %s", err)
	}

	return &Client{MongoClient: client}
}
