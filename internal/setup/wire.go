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

func NewUserService(userRepo *repository.UserRepository, taskRepo *repository.TaskRepository) *services.UserService {
    return services.NewUserService(userRepo, taskRepo)
}

var UserServiceSet = wire.NewSet(
    NewUserService,
    wire.Bind(new(services.UserServiceInterface), new(*services.UserService)),
)

var SendHandlerSet = wire.NewSet(
    api.NewSendHandlerController,
    services.NewPromptService,
    UserServiceSet,
    repository.NewTaskRepository,
)

func InitializeSendHandlerController(db *mongo.Database) *api.SendHandlerController {
	wire.Build(repository.NewUserRepository, repository.NewTargetModelRepository,  SendHandlerSet,
		services.NewTargetModelService)
    return nil
}

func InitializeAuthController(db *mongo.Database) *api.AuthController {
	wire.Build(repository.NewUserRepository, repository.NewTaskRepository,
		services.NewUserService, api.NewAuthController)
	return &api.AuthController{}
}
