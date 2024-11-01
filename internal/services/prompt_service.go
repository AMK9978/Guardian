package services

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"

	"guardian/configs"
	"guardian/internal/models"
	"guardian/internal/models/entities"
	"guardian/utlis/logger"

	"github.com/pkg/errors"
)

type PromptServiceInterface interface {
	ProcessPrompt(ctx context.Context, reqBody *models.RefereeRequest) (bool, error)
}

type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

func NewHTTPClientProvider() *http.Client {
	return &http.Client{Timeout: configs.GlobalConfig.HttpClientTimeout}
}

type PromptService struct {
	userService UserServiceInterface
	client      HTTPClient
}

var (
	ErrPluginResponseFailed = errors.New("failed to receive a response")
	ErrForwardRequest       = errors.New("failed to forward request")
)

func NewPromptService(userService UserServiceInterface, client HTTPClient) *PromptService {
	return &PromptService{
		userService: userService,
		client:      client,
	}
}

func (p *PromptService) ProcessPrompt(ctx context.Context, reqBody *models.RefereeRequest) (bool, error) {
	if reqBody.Prompt == "" {
		return false, nil
	}
	result, err := p.pipeline(ctx, reqBody)
	if err != nil {
		return false, err
	}
	if !result {
		return false, nil
	}

	return true, nil
}

func (p *PromptService) pipeline(ctx context.Context, req *models.RefereeRequest) (bool, error) {
	tasks, err := p.userService.GetUserTasksByID(req.UserID)
	if err != nil {
		logger.GetLogger().Errorf("err in pipeline: %v", err)
		return false, err
	}

	workerPoolSize := configs.GlobalConfig.PipelineWorkerPoolSize

	taskChan := make(chan entities.Task, len(tasks))
	resultsChan := make(chan entities.TaskResult, len(tasks))
	quit := make(chan struct{})
	var closeQuitOnce sync.Once
	var wg sync.WaitGroup
	for i := 0; i < workerPoolSize; i++ {
		wg.Add(1)
		go p.worker(ctx, taskChan, resultsChan, quit, req, &wg, &closeQuitOnce)
	}

	for _, task := range tasks {
		taskChan <- task
	}
	close(taskChan)

	wg.Wait()
	close(resultsChan)

	for result := range resultsChan {
		if result.Err != nil {
			logger.GetLogger().Errorf("task %s faced error", result.Err)
		}

		if !result.Success {
			logger.GetLogger().Infof("task %s failed:", result.TaskType)
			return false, nil
		}
	}

	return true, nil
}

func (p *PromptService) worker(ctx context.Context, taskChan chan entities.Task, resultsChan chan entities.TaskResult,
	quit chan struct{}, reqBody *models.RefereeRequest, wg *sync.WaitGroup, closeQuitOnce *sync.Once,
) {
	defer wg.Done()

	for {
		select {
		case task, ok := <-taskChan:
			if !ok {
				return
			}
			taskType := task.Type

			result, err := p.forwardRequest(ctx, task.Address, reqBody)
			if err != nil {
				resultsChan <- entities.TaskResult{TaskType: taskType, Success: false, Err: err}
				closeQuitOnce.Do(func() {
					close(quit)
				})
				return
			}

			if !result.Status {
				resultsChan <- entities.TaskResult{TaskType: taskType, Success: false}
				closeQuitOnce.Do(func() {
					close(quit)
				})
				return
			}

			resultsChan <- entities.TaskResult{TaskType: taskType, Success: true}

		case <-quit:
			return
		}
	}
}

func (p *PromptService) forwardRequest(ctx context.Context, taskAddress string,
	reqBody *models.RefereeRequest,
) (*models.SendResponse, error) {
	marshalledBody, err := json.Marshal(&reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to marshall the request body: %w", err)
	}
	newReq, err := http.NewRequestWithContext(ctx, http.MethodPost, taskAddress, bytes.NewBuffer(marshalledBody))
	if err != nil {
		return nil, fmt.Errorf("failed to create new request: %w", err)
	}

	newReq.Header.Set("Content-Type", "application/json")

	resp, err := p.client.Do(newReq)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrForwardRequest, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%w from: %s", ErrPluginResponseFailed, taskAddress)
	}

	var sendResponse models.SendResponse
	if err := json.NewDecoder(resp.Body).Decode(&sendResponse); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &sendResponse, nil
}
