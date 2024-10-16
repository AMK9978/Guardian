package mongodb

import (
	"context"

	"guardian/configs"
	"guardian/utlis/logger"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	Client *mongo.Client
	Database *mongo.Database
)

func Init() {
	var err error
	Client, err = NewClient(configs.GlobalConfig.MongoDBURI)
	if err != nil {
		logger.GetLogger().Fatal(err)
	}
	Database = Client.Database(configs.GlobalConfig.PrimaryDBName)
}

func NewClient(mongoDBURI string) (*mongo.Client, error) {
	clientOptions := options.Client().ApplyURI(mongoDBURI)
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		return nil, err
	}
	err = client.Ping(context.Background(), nil)
	if err != nil {
		return nil, err
	}
	return client, err
}

func Disconnect() error {
	if err := Client.Disconnect(context.Background()); err != nil {
		return err
	}
	return nil
}
