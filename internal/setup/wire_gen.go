// Code generated by Wire. DO NOT EDIT.

//go:generate go run -mod=mod github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package setup

import (
	"github.com/google/wire"
	"go.mongodb.org/mongo-driver/mongo"
	"guardian/api"
	"guardian/internal/middleware"
	"guardian/internal/plugins"
	"guardian/internal/repository"
	"guardian/internal/services"
)

// Injectors from wire.go:

func InitializeSendHandlerController(db *mongo.Database) *api.SendHandlerController {
	userRepository := repository.NewUserRepository(db)
	taskRepository := repository.NewTaskRepository(db)
	userService := NewUserService(userRepository, taskRepository)
	client := services.NewHTTPClientProvider()
	httpClient := plugins.NewHTTPClient(client)
	pluginRepository := repository.NewPluginRepository(db)
	pluginService := services.NewPluginService(pluginRepository)
	promptService := services.NewPromptService(userService, httpClient, pluginService)
	targetModelRepository := repository.NewTargetModelRepository(db)
	targetModelService := services.NewTargetModelService(targetModelRepository)
	middlewareMiddleware := middleware.NewMiddleware()
	sendHandlerController := api.NewSendHandlerController(promptService, targetModelService, middlewareMiddleware)
	return sendHandlerController
}

func InitializeAuthController(db *mongo.Database) *api.AuthController {
	userRepository := repository.NewUserRepository(db)
	taskRepository := repository.NewTaskRepository(db)
	userService := services.NewUserService(userRepository, taskRepository)
	authController := api.NewAuthController(userService)
	return authController
}

// wire.go:

func NewUserService(userRepo *repository.UserRepository, taskRepo *repository.TaskRepository) *services.UserService {
	return services.NewUserService(userRepo, taskRepo)
}

var UserServiceSet = wire.NewSet(
	NewUserService, wire.Bind(new(services.UserServiceInterface), new(*services.UserService)),
)

var SendHandlerSet = wire.NewSet(api.NewSendHandlerController, middleware.NewMiddleware, wire.Bind(new(middleware.Interface), new(*middleware.Middleware)), repository.NewPluginRepository, wire.Bind(new(repository.PluginRepoInterface), new(*repository.PluginRepository)), services.NewPluginService, wire.Bind(new(services.PluginServiceInterface), new(*services.PluginService)), services.NewPromptService, wire.Bind(new(services.PromptServiceInterface), new(*services.PromptService)), UserServiceSet, repository.NewTaskRepository)
