package services

import (
	"errors"
	"guardian/internal/models"
	"guardian/utlis/logger"
)

type PromptService struct {
	userTaskService *UserTaskService
	tasksMap        map[string]ProcessingTask
}

func NewPromptService(userTaskService *UserTaskService) *PromptService {
	return &PromptService{
		userTaskService: userTaskService,
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
	userTasks, err := p.userTaskService.GetUserTasks(req.UserID)
	if err != nil {
		logger.GetLogger().Error(err)
		return false
	}
	for _, userTask := range userTasks {
		taskType := userTask.Task.Type

		if task, exists := p.tasksMap[taskType]; exists {
			result, err := task.Process(req)
			if err != nil {
				logger.GetLogger().Error(err)
				return false
			}
			if !result {
				return false
			}
		}
	}
	
	return true
}
