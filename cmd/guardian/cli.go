package guardian

import (
	"fmt"
	"net/http"

	"guardian/configs"
	"guardian/internal/metrics"
	"guardian/internal/mongodb"
	"guardian/internal/server"
	"guardian/utlis/logger"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "guardian",
	Short: "Guardian CLI",
}

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	logger.InitLogger()

	configs.GlobalConfig = configs.LoadConfig()

	metrics.Init()
	go func() {
		logger.GetLogger().Infof("Starting metrics server on :%s", configs.GlobalConfig.ServerPort)
		err := http.ListenAndServe(fmt.Sprintf(":%s", configs.GlobalConfig.ServerPort), metrics.Handler())
		if err != nil {
			logger.GetLogger().Fatalf("Failed to start metrics server: %s", err)
		}
	}()

	// redis.Init(configs.GlobalConfig.RedisAddr)
	// rabbitMQClient := rabbitmq.NewClient(cfg.RabbitMQURI)
	mongodb.Init()

	// milvus.NewClient(configs.GlobalConfig.MilvusURI)

	startServer()

	logger.GetLogger().Info("Successfully connected to all services")
}

func startServer() {
	server.StartServer()
}
