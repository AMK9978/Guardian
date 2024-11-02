package api

import (
	"bytes"
	"context"
	"encoding/json"
	"guardian/internal/mocks"
	"guardian/internal/models"
	"guardian/internal/models/entities"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var (
	ErrTargetModel      = errors.New("error in fetching target models")
	ErrProcessPrompt    = errors.New("error in fetching processing the prompt")
	ErrForwardingPrompt = errors.New("error in forwarding request to the target")
)

func TestSendHandler(t *testing.T) {
	t.Parallel()

	m := new(mocks.MockMiddleware)
	reqBody := models.RefereeRequest{
		UserID:   primitive.ObjectID{},
		Chat:     "Previous conversation",
		Prompt:   "Hello",
		TargetID: primitive.ObjectID{},
	}

	m.On("GetUserFromContext").Return(mock.Anything, nil)

	t.Run("target models fetch fails", func(t *testing.T) {
		t.Parallel()

		targetModelService := new(mocks.MockTargetModelService)
		promptService := new(mocks.MockPromptService)
		controller := NewSendHandlerController(promptService, targetModelService, m)

		targetModelService.On("GetTargetModel", mock.Anything, mock.Anything).
			Return(entities.TargetModel{}, ErrTargetModel)

		body, _ := json.Marshal(reqBody)
		req, _ := http.NewRequestWithContext(context.Background(), http.MethodPost, "/send", bytes.NewBuffer(body))
		rec := httptest.NewRecorder()

		controller.SendHandler(rec, req)

		assert.Equal(t, http.StatusInternalServerError, rec.Code)
	})

	t.Run("process prompt returns error", func(t *testing.T) {
		t.Parallel()

		targetModelService := new(mocks.MockTargetModelService)
		promptService := new(mocks.MockPromptService)
		controller := NewSendHandlerController(promptService, targetModelService, m)

		targetModelService.On("GetTargetModel", mock.Anything, mock.Anything).
			Return(entities.TargetModel{}, nil)
		promptService.On("ProcessPrompt").Return(false, ErrProcessPrompt)

		body, _ := json.Marshal(reqBody)
		req, _ := http.NewRequestWithContext(context.Background(), http.MethodPost, "/send", bytes.NewBuffer(body))
		rec := httptest.NewRecorder()

		controller.SendHandler(rec, req)

		assert.Equal(t, http.StatusInternalServerError, rec.Code)
	})

	t.Run("process prompt returns false", func(t *testing.T) {
		t.Parallel()

		targetModelService := new(mocks.MockTargetModelService)
		promptService := new(mocks.MockPromptService)
		controller := NewSendHandlerController(promptService, targetModelService, m)

		targetModelService.On("GetTargetModel", mock.Anything, mock.Anything).
			Return(entities.TargetModel{}, nil)

		promptService.On("ProcessPrompt").Return(false, nil)
		promptService.On("Do").Return(nil, ErrForwardingPrompt)

		body, _ := json.Marshal(reqBody)
		req, _ := http.NewRequestWithContext(context.Background(), http.MethodPost, "/send", bytes.NewBuffer(body))
		rec := httptest.NewRecorder()
		respBody := models.SendResponse{Status: false}
		respBodyJSON, _ := json.Marshal(respBody)

		controller.SendHandler(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, bytes.NewBuffer(respBodyJSON), rec.Body)
	})

	t.Run("forwarding request to the target fails", func(t *testing.T) {
		t.Parallel()

		targetModelService := new(mocks.MockTargetModelService)
		promptService := new(mocks.MockPromptService)
		controller := NewSendHandlerController(promptService, targetModelService, m)

		targetModelService.On("GetTargetModel", mock.Anything, mock.Anything).
			Return(entities.TargetModel{}, nil)
		promptService.On("ProcessPrompt").Return(true, nil)
		promptService.On("Do").Return(&http.Response{}, ErrForwardingPrompt)

		body, _ := json.Marshal(reqBody)
		req, _ := http.NewRequestWithContext(context.Background(), http.MethodPost, "/send", bytes.NewBuffer(body))
		rec := httptest.NewRecorder()

		controller.SendHandler(rec, req)

		assert.Equal(t, http.StatusInternalServerError, rec.Code)
	})
}
