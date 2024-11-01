// wire.go
//go:build wireinject
// +build wireinject

package setup

import (
	"net/http"

	"guardian/api"
	"guardian/internal/repository"
	"guardian/internal/services"

	"github.com/google/wire"
	"go.mongodb.org/mongo-driver/mongo"
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
	wire.Build(
		repository.NewUserRepository,
		wire.Bind(new(services.HTTPClient), new(*http.Client)),
		repository.NewTargetModelRepository,
		SendHandlerSet,
		services.NewTargetModelService,
		services.NewHTTPClientProvider,
	)
	return nil
}

func InitializeAuthController(db *mongo.Database) *api.AuthController {
	wire.Build(
		repository.NewUserRepository,
		repository.NewTaskRepository,
		services.NewUserService,
		wire.Bind(new(services.UserServiceInterface), new(*services.UserService)),
		api.NewAuthController,
	)
	return &api.AuthController{}
}
