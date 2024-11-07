// wire.go
//go:build wireinject
// +build wireinject

package setup

import (
	"guardian/api"
	"guardian/internal/middleware"
	"guardian/internal/repository"
	"guardian/internal/services"
	"guardian/internal/plugins"

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
	middleware.NewMiddleware,
	wire.Bind(new(middleware.Interface), new(*middleware.Middleware)),
	repository.NewPluginRepository,
	wire.Bind(new(repository.PluginRepoInterface), new(*repository.PluginRepository)),
	services.NewPluginService,
	wire.Bind(new(services.PluginServiceInterface), new(*services.PluginService)),
	services.NewPromptService,
	wire.Bind(new(services.PromptServiceInterface), new(*services.PromptService)),
	UserServiceSet,
	repository.NewTaskRepository,
)

func InitializeSendHandlerController(db *mongo.Database) *api.SendHandlerController {
	wire.Build(
		repository.NewUserRepository,
		repository.NewTargetModelRepository,
		wire.Bind(new(repository.TargetModelRepoInterface), new(*repository.TargetModelRepository)),
		plugins.NewHTTPClient,
		wire.Bind(new(plugins.HTTPClientInterface), new(*plugins.HTTPClient)),
		SendHandlerSet,
		services.NewTargetModelService,
		wire.Bind(new(services.TargetModelServiceInterface), new(*services.TargetModelService)),
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
