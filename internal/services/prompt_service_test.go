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
	pluginList := []primitive.ObjectID{primitive.NewObjectID()}
	mockUserService := new(mocks.MockUserService)
	mockPluginService := new(mocks.MockPluginService)
	mockClient := new(mocks.MockClient)
	pluginClient := mockClient
	promptService := NewPromptService(mockUserService, pluginClient, mockPluginService)
	userID := primitive.NewObjectID()
	validReq := &models.PluginRequest{
		UserID:   userID,
		Prompt:   "Test prompt",
		TargetID: primitive.NewObjectID(),
	}

	tests := []struct {
		name           string
		reqBody        *models.PluginRequest
		mockTasks      []entities.Task
		mockError      error
		expectedResult bool
		expectedError  error
	}{
		{
			name:           "Valid prompt with tasks",
			reqBody:        validReq,
			mockTasks:      []entities.Task{{Type: "ExampleTask", Plugins: pluginList}},
			mockError:      nil,
			expectedResult: true,
			expectedError:  nil,
		},
		{
			name:           "Empty prompt",
			reqBody:        &models.PluginRequest{UserID: userID, Prompt: ""},
			expectedResult: false,
			expectedError:  nil,
		},
		{
			name:           "User tasks error",
			reqBody:        validReq,
			mockTasks:      nil,
			mockError:      ErrUserTasks,
			expectedResult: false,
			expectedError:  ErrUserTasks,
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
	pluginListID := []primitive.ObjectID{primitive.NewObjectID()}
	tests := []struct {
		name         string
		userTasks    []entities.Task
		userTasksErr error
		mockResp     models.PluginResponse
		mockStatus   int
		expectErr    bool
		expectRes    bool
	}{
		{
			"Normal Case",
			[]entities.Task{{Type: "ExampleTask", Plugins: pluginListID}},
			nil,
			models.PluginResponse{Status: true},
			http.StatusOK, false,
			true,
		},
		{
			"GetUserTasksFails",
			[]entities.Task{},
			ErrUserTasks,
			models.PluginResponse{},
			0, true, false,
		},
		{
			"ReceivesFalse",
			[]entities.Task{{Type: "ExampleTask", Plugins: pluginListID}},
			nil,
			models.PluginResponse{Status: false},
			http.StatusOK, false,
			false,
		},
		{
			"ReceivesError",
			[]entities.Task{{Type: "ExampleTask", Plugins: pluginListID}},
			nil,
			models.PluginResponse{Status: true},
			http.StatusInternalServerError,
			false, false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			mockPluginService := new(mocks.MockPluginService)
			mockUserService := new(mocks.MockUserService)
			mockClient := new(mocks.MockClient)

			pluginList := []entities.Plugin{
				{
					ID:       primitive.ObjectID{},
					Name:     "",
					Provider: "",
					Address:  "",
					Status:   0,
					Token:    "",
					Protocol: entities.Protocol{
						ID:   primitive.ObjectID{},
						Type: entities.HTTPProtocol,
					},
				},
			}
			promptService := NewPromptService(mockUserService, mockClient, mockPluginService)
			configs.GlobalConfig = configs.Config{
				PipelineWorkerPoolSize: runtime.NumCPU(),
			}

			userID := primitive.NewObjectID()
			reqBody := &models.PluginRequest{
				UserID:   userID,
				Prompt:   "Test prompt",
				TargetID: primitive.NewObjectID(),
			}

			mockUserService.On("GetUserTasksByID", userID).Return(tt.userTasks, tt.userTasksErr)
			for _, task := range tt.userTasks {
				mockPluginService.On("GetPluginsByTask", mock.Anything, task).
					Return(pluginList, nil)
			}
			m, _ := json.Marshal(tt.mockResp)
			mockClient.On("Do", mock.Anything).Return(&http.Response{
				StatusCode: tt.mockStatus,
				Body:       io.NopCloser(bytes.NewBuffer(m)),
			}, nil)
			mockClient.On("Forward", mock.Anything, reqBody).Return(&models.PluginResponse{
				Status: tt.expectRes,
				Score: 1,
			}, nil)

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
