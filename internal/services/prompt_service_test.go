package services

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"runtime"
	"testing"

	"guardian/configs"
	"guardian/internal/mocks"
	"guardian/internal/models"
	"guardian/internal/models/entities"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var ErrUserTasks = errors.New("error in GetUserTasks")

func TestProcessPrompt(t *testing.T) {
	t.Parallel()

	mockUserService := new(mocks.MockUserService)
	mockClient := new(mocks.MockClient)
	promptService := NewPromptService(mockUserService, mockClient)
	userID := primitive.NewObjectID()
	validReq := &models.RefereeRequest{
		UserID:   userID,
		Prompt:   "Test prompt",
		TargetID: primitive.NewObjectID(),
	}

	tests := []struct {
		name           string
		reqBody        *models.RefereeRequest
		mockTasks      []entities.Task
		mockError      error
		expectedResult bool
		expectedError  error
	}{
		{
			name:           "Valid prompt with tasks",
			reqBody:        validReq,
			mockTasks:      []entities.Task{{Type: "ExampleTask", Address: "http://example.com"}},
			mockError:      nil,
			expectedResult: true,
			expectedError:  nil,
		},
		{
			name:           "Empty prompt",
			reqBody:        &models.RefereeRequest{UserID: userID, Prompt: ""},
			expectedResult: false,
			expectedError:  nil,
		},
		{
			name:           "User tasks error",
			reqBody:        validReq,
			mockTasks:      nil,
			mockError:      ErrUserTasks,
			expectedResult: false,
			expectedError:  nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			mockUserService.On("GetUserTasksByID", tt.reqBody.UserID).Return(tt.mockTasks, tt.mockError)
			result, err := promptService.ProcessPrompt(context.Background(), tt.reqBody)

			require.Equal(t, tt.expectedResult, result)
			require.Equal(t, tt.expectedError, err)
			mockUserService.On("GetUserTasksByID", tt.reqBody.UserID).Unset()
		})
	}
}

func TestPipeline_Cases(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name         string
		userTasks    []entities.Task
		userTasksErr error
		mockResp     models.SendResponse
		mockStatus   int
		expectErr    bool
		expectRes    bool
	}{
		{
			"Normal Case",
			[]entities.Task{{Type: "ExampleTask", Address: "http://example.com"}},
			nil,
			models.SendResponse{Status: true},
			http.StatusOK, false,
			true,
		},
		{
			"GetUserTasksFails",
			[]entities.Task{},
			ErrUserTasks,
			models.SendResponse{},
			0, true, false,
		},
		{
			"ReceivesFalse",
			[]entities.Task{{Type: "ExampleTask", Address: "http://example.com"}},
			nil,
			models.SendResponse{Status: false},
			http.StatusOK, false,
			false,
		},
		{
			"ReceivesError",
			[]entities.Task{{Type: "ExampleTask", Address: "http://example.com"}},
			nil,
			models.SendResponse{Status: true},
			http.StatusInternalServerError,
			false, false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			mockUserService := new(mocks.MockUserService)
			mockClient := new(mocks.MockClient)
			promptService := NewPromptService(mockUserService, mockClient)
			configs.GlobalConfig = configs.Config{
				PipelineWorkerPoolSize: runtime.NumCPU(),
			}

			userID := primitive.NewObjectID()
			reqBody := &models.RefereeRequest{
				UserID:   userID,
				Prompt:   "Test prompt",
				TargetID: primitive.NewObjectID(),
			}

			mockUserService.On("GetUserTasksByID", userID).Return(tt.userTasks, tt.userTasksErr)

			if tt.mockStatus != 0 {
				m, _ := json.Marshal(tt.mockResp)
				mockClient.On("Do", mock.Anything).Return(&http.Response{
					StatusCode: tt.mockStatus,
					Body:       io.NopCloser(bytes.NewBuffer(m)),
				}, nil)
			}

			result, err := promptService.pipeline(context.Background(), reqBody)

			if tt.expectErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.expectRes, result)
			}
		})
	}
}
