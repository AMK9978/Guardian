package services

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"guardian/configs"
	"guardian/internal/models"
	"guardian/internal/models/entities"
	"guardian/utlis/logger"
	"net/http"
	"sync"
)

type PromptServiceInterface interface {
	ProcessPrompt(req models.SendRequest, r *http.Request) (bool, error)
}

type PromptService struct {
	userService UserServiceInterface
}

func NewPromptService(userService UserServiceInterface) *PromptService {
	return &PromptService{
		userService: userService,
	}
}

func (p *PromptService) ProcessPrompt(req models.SendRequest, r *http.Request) (bool, error) {
	if req.Prompt == "" {
		return false, errors.New("empty prompt")
	}

	if p.pipeline(req, r) {
		return false, nil
	}

	return true, nil
}

func (p *PromptService) pipeline(req models.SendRequest, r *http.Request) bool {
	tasks, err := p.userService.GetUserTasksByID(req.UserID)
	if err != nil {
		logger.GetLogger().Error(err)
		return false
	}

	workerPoolSize := configs.LoadConfig().PipelineWorkerPoolSize

	taskChan := make(chan entities.Task, len(tasks))
	resultsChan := make(chan entities.TaskResult, len(tasks))
	quit := make(chan struct{})

	var wg sync.WaitGroup
	for i := 0; i < workerPoolSize; i++ {
		wg.Add(1)
		go p.worker(taskChan, resultsChan, quit, r, &wg)
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
	r *http.Request, wg *sync.WaitGroup) {
	defer wg.Done()

	for {
		select {
		case task, ok := <-taskChan:
			if !ok {
				return
			}
			taskType := task.Type

			result, err := p.forwardRequest(r, task.Address)
			if err != nil {
				resultsChan <- entities.TaskResult{TaskType: taskType, Success: false, Err: err}
				close(quit)
				return
			}

			if !result.Success {
				resultsChan <- entities.TaskResult{TaskType: taskType, Success: false}
				close(quit)
				return
			}

			resultsChan <- entities.TaskResult{TaskType: taskType, Success: true}

		case <-quit:
			return
		}
	}
}

func (p *PromptService) forwardRequest(req *http.Request, taskAddress string) (*models.SendResponse, error) {
	newReq, err := http.NewRequestWithContext(context.Background(), req.Method, taskAddress, req.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to create new request: %w", err)
	}
	defer req.Body.Close()

	for k, v := range req.Header {
		newReq.Header[k] = v
	}

	client := &http.Client{}
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
