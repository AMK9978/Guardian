package models

import (
	"context"
	"time"

	"guardian/internal/models/entities"
	"guardian/internal/mongodb"
	"guardian/utlis/logger"

	"github.com/brianvoe/gofakeit/v6"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func generateFakeSendRequest() SendRequest {
	modelID, _ := primitive.ObjectIDFromHex(gofakeit.UUID())
	chatID, _ := primitive.ObjectIDFromHex(gofakeit.UUID())
	userID, _ := primitive.ObjectIDFromHex(gofakeit.UUID())

	target := entities.TargetModel{
		ID:    modelID,
		Token: gofakeit.LetterN(10),
	}

	return SendRequest{
		UserID: userID,
		ChatID: &chatID,
		Prompt: gofakeit.Sentence(5),
		Target: []entities.TargetModel{target},
	}
}

func generateFakeGroup() entities.Group {
	groupID, _ := primitive.ObjectIDFromHex(gofakeit.UUID())
	return entities.Group{
		ID:     groupID,
		Name:   gofakeit.Company(),
		Status: gofakeit.Number(1, 3),
	}
}

func generateFakeUser() entities.User {
	userID, _ := primitive.ObjectIDFromHex(gofakeit.UUID())
	return entities.User{
		ID:       userID,
		Name:     gofakeit.Name(),
		Password: gofakeit.Password(true, false, false, false, false, 32),
		Status:   gofakeit.Number(0, 1),
		Groups:   []entities.Group{generateFakeGroup()},
	}
}

func generate() {
	gofakeit.Seed(time.Now().UnixNano())

	sendRequest := generateFakeSendRequest()
	group := generateFakeGroup()
	logger.GetLogger().Infof("Fake Group: %+v\n", group)

	collection := mongodb.Client.Database("test_db").Collection("send_requests")
	_, err := collection.InsertOne(context.Background(), bson.M{
		"user_id": sendRequest.UserID,
		"chat_id": sendRequest.ChatID,
		"prompt":  sendRequest.Prompt,
		"targets": sendRequest.Target,
	})
	if err != nil {
		logger.GetLogger().Fatal(err)
	}

	logger.GetLogger().Info("Data inserted into MongoDB")
}
