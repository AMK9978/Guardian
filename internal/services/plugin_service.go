package services

import (
	"context"

	"guardian/internal/models/entities"
	"guardian/internal/repository"
)

type PluginServiceInterface interface {
	GetPluginsByTask(ctx context.Context, task entities.Task) ([]entities.Plugin, error)
}

type PluginService struct {
	pluginRepo repository.PluginRepoInterface
}

func NewPluginService(pluginRepo repository.PluginRepoInterface) *PluginService {
	return &PluginService{pluginRepo: pluginRepo}
}

func (t *PluginService) GetPluginsByTask(ctx context.Context, task entities.Task) ([]entities.Plugin, error) {
	plugins, err := t.pluginRepo.GetPluginsByTask(ctx, task)
	if err != nil {
		return nil, err
	}
	return plugins, err
}
