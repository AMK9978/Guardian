package services

import (
	"context"
	"fmt"
	"net/http"
	"sync"

	"guardian/configs"
	"guardian/internal/models"
	"guardian/internal/models/entities"
	"guardian/internal/plugins"
	"guardian/utlis/logger"

	"github.com/pkg/errors"
)

type PromptServiceInterface interface {
	ProcessPrompt(ctx context.Context, reqBody *models.PluginRequest) (bool, error)
	SendPrompt(ctx context.Context, newReq *http.Request) (*http.Response, error)
}

func NewHTTPClientProvider() *http.Client {
	return &http.Client{Timeout: configs.GlobalConfig.HttpClientTimeout}
}

type PromptService struct {
	userService   UserServiceInterface
	pluginService PluginServiceInterface
	client plugins.HTTPClientInterface
}

var (
	ErrForwardRequest = errors.New("failed to forward request")
)

func NewPromptService(userService UserServiceInterface, client plugins.HTTPClientInterface,
	pluginService PluginServiceInterface) *PromptService {
	return &PromptService{
		userService: userService,
		client:  client,
		pluginService: pluginService,
	}
}

func (p *PromptService) SendPrompt(ctx context.Context, newReq *http.Request) (*http.Response, error) {
	return p.client.Do(newReq)
}

func (p *PromptService) ProcessPrompt(ctx context.Context, reqBody *models.PluginRequest) (bool, error) {
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

func (p *PromptService) pipeline(ctx context.Context, req *models.PluginRequest) (bool, error) {
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
	quit chan struct{}, reqBody *models.PluginRequest, wg *sync.WaitGroup, closeQuitOnce *sync.Once,
) {
	defer wg.Done()

	for {
		select {
		case task, ok := <-taskChan:
			if !ok {
				return
			}
			taskType := task.Type
			pluginList, err := p.pluginService.GetPluginsByTask(ctx, task)
			if err != nil {
				resultsChan <- entities.TaskResult{TaskType: taskType, Success: false, Err: err}
				closeQuitOnce.Do(func() {
					close(quit)
				})
				return
			}
			result, err := p.forwardRequest(ctx, pluginList, reqBody)
			if err != nil {
				resultsChan <- entities.TaskResult{TaskType: taskType, Success: false, Err: err}
				closeQuitOnce.Do(func() {
					close(quit)
				})
				return
			}

			if !result {
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

func (p *PromptService) forwardRequest(ctx context.Context, pluginList []entities.Plugin,
	reqBody *models.PluginRequest) (bool, error) {
	for _, plugin := range pluginList {
		var client plugins.PluginClient

		switch plugin.Protocol.Type {
		case entities.HTTPProtocol:
			client = p.client

		case entities.GRPCProtocol:
			grpcConn, err := configs.GlobalConfig.GRPCManager.GetClient(plugin)
			if err != nil {
				return false, fmt.Errorf("%w: %w", ErrForwardRequest, err)
			}
			client = plugins.NewPluginGRPCClient(grpcConn)

		default:
			return false, fmt.Errorf("unsupported protocol type: %s", plugin.Protocol.Type)
		}

		reqBody.Address = plugin.Address
		result, err := client.Forward(ctx, reqBody)
		if err != nil {
			return false, err
		}
		if !result.Status {
			return false, nil
		}
	}

	return true, nil
}
