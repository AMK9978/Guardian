package services

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"guardian/configs"
	"guardian/internal/models"
	"guardian/internal/models/entities"
	"guardian/utlis/logger"
	"net/http"
	"sync"
	"time"
)

type PromptServiceInterface interface {
	ProcessPrompt(ctx context.Context, reqBody *models.RefereeRequest) (bool, error)
}

type PromptService struct {
	userService UserServiceInterface
}

func NewPromptService(userService UserServiceInterface) *PromptService {
	return &PromptService{
		userService: userService,
	}
}

func (p *PromptService) ProcessPrompt(ctx context.Context, reqBody *models.RefereeRequest) (bool, error) {
	if reqBody.Prompt == "" {
		return false, nil
	}
	if !p.pipeline(ctx, reqBody) {
		return false, nil
	}

	return true, nil
}

func (p *PromptService) pipeline(ctx context.Context, req *models.RefereeRequest) bool {
	tasks, err := p.userService.GetUserTasksByID(req.UserID)
	if err != nil {
		logger.GetLogger().Errorf("err in pipeline: %v", err)
		return false
	}

	workerPoolSize := configs.LoadConfig().PipelineWorkerPoolSize

	taskChan := make(chan entities.Task, len(tasks))
	resultsChan := make(chan entities.TaskResult, len(tasks))
	quit := make(chan struct{})
	var closeQuitOnce sync.Once
	var wg sync.WaitGroup
	for i := 0; i < workerPoolSize; i++ {
		wg.Add(1)
		go p.worker(taskChan, resultsChan, quit, ctx, req, &wg, &closeQuitOnce)
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
			return false
		}
	}

	return true
}

func (p *PromptService) worker(taskChan chan entities.Task, resultsChan chan entities.TaskResult, quit chan struct{},
	ctx context.Context, reqBody *models.RefereeRequest, wg *sync.WaitGroup, closeQuitOnce *sync.Once) {
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
	reqBody *models.RefereeRequest) (*models.SendResponse, error) {

	marshalledBody, err := json.Marshal(&reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to marshall the request body: %w", err)
	}
	newReq, err := http.NewRequestWithContext(ctx, http.MethodPost, taskAddress, bytes.NewBuffer(marshalledBody))
	if err != nil {
		return nil, fmt.Errorf("failed to create new request: %w", err)
	}

	newReq.Header.Set("Content-Type", "application/json")

	// TODO: Use configs
	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(newReq)
	if err != nil {
		return nil, fmt.Errorf("failed to forward request: %w", err)
	}
	defer resp.Body.Close()

	var sendResponse models.SendResponse
	if err := json.NewDecoder(resp.Body).Decode(&sendResponse); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &sendResponse, nil
}
