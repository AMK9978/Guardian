package models

import (
	"context"
	"guardian/internal/models/entities"
	"time"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"guardian/internal/mongodb"
	"guardian/utlis/logger"
)

func generateFakeSendRequest() entities.SendRequest {
	chatID, _ := uuid.Parse(gofakeit.UUID())

	target := entities.TargetModel{
		ModelID: uuid.MustParse(gofakeit.UUID()),
		Token:   gofakeit.LetterN(10),
	}

	return entities.SendRequest{
		UserID: uuid.MustParse(gofakeit.UUID()),
		ChatID: &chatID,
		Prompt: gofakeit.Sentence(5),
		Target: []entities.TargetModel{target},
	}
}

func generateFakeGroup() entities.Group {
	return entities.Group{
		ID:     uuid.MustParse(gofakeit.UUID()),
		Name:   gofakeit.Company(),
		Status: gofakeit.Number(1, 3),
		Users:  []entities.User{generateFakeUser(), generateFakeUser()},
	}
}

func generateFakeUser() entities.User {
	return entities.User{
		ID:       uuid.MustParse(gofakeit.UUID()),
		Name:     gofakeit.Name(),
		Password: gofakeit.Password(true, false, false, false, false, 32),
		Status:   gofakeit.Number(0, 1),
		Groups:   []entities.Group{generateFakeGroup()},
	}
}

func generateFakeAIModel() entities.AIModel {
	return entities.AIModel{
		Name:    gofakeit.AppName(),
		Address: gofakeit.URL(),
		Status:  gofakeit.Number(1, 3),
	}
}

func generate() {
	gofakeit.Seed(time.Now().UnixNano())

	sendRequest := generateFakeSendRequest()
	group := generateFakeGroup()
	logger.GetLogger().Infof("Fake Group: %+v\n", group)
	aiModel := generateFakeAIModel()
	logger.GetLogger().Infof("Fake AIModel: %+v\n", aiModel)

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
