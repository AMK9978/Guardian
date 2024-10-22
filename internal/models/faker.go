package models

import (
	"context"
	"time"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"guardian/internal/mongodb"
	"guardian/utlis/logger"
)

func generateFakeSendRequest() SendRequest {
	chatID, _ := uuid.Parse(gofakeit.UUID())

	target := TargetModel{
		ModelID: uuid.MustParse(gofakeit.UUID()),
		Token:   gofakeit.LetterN(10),
	}

	return SendRequest{
		UserID: uuid.MustParse(gofakeit.UUID()),
		ChatID: &chatID,
		Prompt: gofakeit.Sentence(5),
		Target: []TargetModel{target},
	}
}

func generateFakeGroup() Group {
	return Group{
		ID:     uuid.MustParse(gofakeit.UUID()),
		Name:   gofakeit.Company(),
		Status: gofakeit.Number(1, 3),
		Users:  []User{generateFakeUser(), generateFakeUser()},
	}
}

func generateFakeUser() User {
	return User{
		ID:       uuid.MustParse(gofakeit.UUID()),
		Name:     gofakeit.Name(),
		Password: gofakeit.Password(true, false, false, false, false, 32),
		Status:   gofakeit.Number(0, 1),
		Groups:   []Group{generateFakeGroup()},
	}
}

func generateFakeAIModel() AIModel {
	return AIModel{
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
