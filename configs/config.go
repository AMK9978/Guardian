package configs

import (
	"os"
)

type Config struct {
	RedisAddr   string
	MongoDBURI  string
	RabbitMQURI string
	MilvusURI   string
	ServerPort  string
}

func LoadConfig() Config {
	return Config{
		RedisAddr:   getEnv("REDIS_ADDR", "localhost:6379"),
		MongoDBURI:  getEnv("MONGODB_URI", "mongodb://localhost:27017"),
		RabbitMQURI: getEnv("RABBITMQ_URI", "amqp://guest:guest@localhost:5672/"),
		MilvusURI:   getEnv("MILVUS_URI", "localhost:19530"),
		ServerPort:  getEnv("SERVER_PORT", "8080"),
	}
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
