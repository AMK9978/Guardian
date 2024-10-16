package guardian

import (
	"guardian/internal/mongodb"
	"net/http"

	"guardian/api"
	"guardian/configs"
	"guardian/internal/metrics"
	"guardian/internal/milvus"
	"guardian/internal/redis"
	"guardian/utlis/logger"

	"github.com/gorilla/mux"
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

	metrics.Init()
	go func() {
		logger.GetLogger().Info("Starting metrics server on :8081")
		if err := http.ListenAndServe(":8081", metrics.Handler()); err != nil {
			logger.GetLogger().Fatalf("Failed to start metrics server: %s", err)
		}
	}()

	redis.NewClient(configs.GlobalConfig.RedisAddr)
	//rabbitMQClient := rabbitmq.NewClient(cfg.RabbitMQURI)
	mongodb.Init()

	milvus.NewClient(configs.GlobalConfig.MilvusURI)

	startServer()

	logger.GetLogger().Info("Successfully connected to all services")
}

func startServer() {
	r := mux.NewRouter()

	r.HandleFunc("/send", api.SendHandler).Methods("POST")

	logger.GetLogger().Println("Server starting on port 8080")
	http.ListenAndServe(":8080", r)
}
