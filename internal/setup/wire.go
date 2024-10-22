// wire.go
//+build wireinject

package setup

import (
	"github.com/google/wire"
	"go.mongodb.org/mongo-driver/mongo"
	"guardian/api"
	"guardian/internal/repository"
	"guardian/internal/services"
)

func NewUserTaskService(userTaskRepo *repository.UserTaskRepository) *services.UserTaskService {
    return services.NewUserTaskService(userTaskRepo)
}

var UserTaskServiceSet = wire.NewSet(
    NewUserTaskService,
    wire.Bind(new(services.UserTaskServiceInterface), new(*services.UserTaskService)),
)

var SendHandlerSet = wire.NewSet(
    api.NewSendHandlerController,
    services.NewPromptService,
    UserTaskServiceSet,
    repository.NewUserTasksRepository,
)

func InitializeSendHandlerController(db *mongo.Database) *api.SendHandlerController {
	wire.Build(SendHandlerSet)
    return nil
}

func InitializeAuthController(db *mongo.Database) *api.AuthController {
	wire.Build(repository.NewUserRepository, services.NewUserService, api.NewAuthController)
	return &api.AuthController{}
}
