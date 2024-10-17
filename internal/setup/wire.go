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

func InitializeSendHandlerController(db *mongo.Database) *api.SendHandlerController {
	wire.Build(repository.NewUserTasksRepository, services.NewUserTaskService, services.NewPromptService,
		api.NewSendHandlerController)
	return &api.SendHandlerController{}
}
