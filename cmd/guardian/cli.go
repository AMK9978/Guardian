package guardian

import (
	"guardian/configs"
	"guardian/pkg/metrics"
	"guardian/pkg/milvus"
	"guardian/pkg/mongodb"
	"guardian/pkg/rabbitmq"
	"guardian/pkg/redis"
	"guardian/utlis/logger"
	"net/http"

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
	log := logger.GetLogger()

	metrics.Init()
	go func() {
		log.Info("Starting metrics server on :8081")
		if err := http.ListenAndServe(":8081", metrics.Handler()); err != nil {
			log.Fatalf("Failed to start metrics server: %s", err)
		}
	}()

	cfg := configs.LoadConfig()

	redisClient := redis.NewClient(cfg.RedisAddr)
	rabbitMQClient := rabbitmq.NewClient(cfg.RabbitMQURI)
	mongoDBClient := mongodb.NewClient(cfg.MongoDBURI)
	milvusClient := milvus.NewClient(cfg.MilvusURI)

	log.Info("Successfully connected to all services")
}
