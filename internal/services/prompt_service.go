package services

import (
	"errors"
	"fmt"
	"guardian/configs"
	"guardian/internal/models"
	"guardian/internal/models/entities"
	"guardian/utlis/logger"
	"sync"
)

type PromptServiceInterface interface {
	ProcessPrompt(req models.SendRequest) (string, error)
}

type PromptService struct {
	userService UserServiceInterface
	tasksMap    map[string]ProcessingTask
}

func NewPromptService(userService UserServiceInterface) *PromptService {
	return &PromptService{
		userService: userService,
		// TODO: Sample
		tasksMap: map[string]ProcessingTask{
			"external-api": &ExternalHttpServiceTask{ApiUrl: "https://google.com"},
		},
	}
}

func (p *PromptService) ProcessPrompt(req models.SendRequest) (string, error) {
	if req.Prompt == "" {
		return "", errors.New("empty prompt")
	}

	if p.pipeline(req) {
		return "malicious", nil
	}

	return "benign", nil
}

func (p *PromptService) pipeline(req models.SendRequest) bool {
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
		go p.worker(taskChan, resultsChan, quit, req, &wg)
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
	req models.SendRequest, wg *sync.WaitGroup) {
	defer wg.Done()

	for {
		select {
		case task, ok := <-taskChan:
			if !ok {
				return
			}
			taskType := task.Type

			if task, exists := p.tasksMap[taskType]; exists {
				result, err := task.Process(req)
				if err != nil {
					resultsChan <- entities.TaskResult{TaskType: taskType, Success: false, Err: err}
					close(quit)
					return
				}

				if !result {
					resultsChan <- entities.TaskResult{TaskType: taskType, Success: false}
					close(quit)
					return
				}

				resultsChan <- entities.TaskResult{TaskType: taskType, Success: true}
			} else {
				err := fmt.Errorf("task type %s not found", taskType)
				resultsChan <- entities.TaskResult{TaskType: taskType, Success: false, Err: err}
			}

		case <-quit:
			return
		}
	}
}
