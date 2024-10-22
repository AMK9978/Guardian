package configs

import (
	"github.com/go-chi/jwtauth/v5"
	"guardian/utlis/logger"
	"os"
	"runtime"
	"strconv"
	"time"
)

var GlobalConfig Config

func init() {
	GlobalConfig = LoadConfig()
}

type Collections struct {
	User     string
	UserTask string
	Group    string
	AIModel  string
}

// NewCollections initializes the collection names.
func NewCollections() *Collections {
	return &Collections{
		User:     "users",
		UserTask: "user_tasks",
		Group:    "groups",
		AIModel:  "ai_models",
	}
}

type Config struct {
	RedisAddr              string
	MongoDBURI             string
	RabbitMQURI            string
	MilvusURI              string
	ServerPort             string
	PrimaryDBName          string
	CollectionNames        *Collections
	PipelineWorkerPoolSize int
	TokenAuth              *jwtauth.JWTAuth
	TokenExpirationTime    time.Duration
}

func LoadConfig() Config {
	numWorkersStr := getEnv("PIPELINE_WORKER_POOL_SIZE", strconv.Itoa(runtime.NumCPU()))
	numWorkers, err := strconv.Atoi(numWorkersStr)
	if err != nil {
		logger.GetLogger().Fatalf("coudn't convert the worker pool size to int: %s", numWorkersStr)
	}

	secretKey := getEnv("JWT_SECRET_KEY", "")
	tokenAuth := jwtauth.New("HS256", []byte(secretKey), nil)
	tokenExpTimeStr := getEnv("TOKEN_EXP_TIME", "72")
	tokenExpTime, err := strconv.Atoi(tokenExpTimeStr)
	if err != nil {
		logger.GetLogger().Fatalf("coudn't convert the token expiration time to int: %s", tokenExpTimeStr)
	}

	return Config{
		RedisAddr:              getEnv("REDIS_ADDR", "localhost:6379"),
		MongoDBURI:             getEnv("MONGODB_URI", "mongodb://localhost:27017"),
		RabbitMQURI:            getEnv("RABBITMQ_URI", "amqp://guest:guest@localhost:5672/"),
		MilvusURI:              getEnv("MILVUS_URI", "localhost:19530"),
		ServerPort:             getEnv("SERVER_PORT", "8080"),
		PrimaryDBName:          getEnv("PRIMARY_DB_NAME", "primary"),
		TokenAuth:              tokenAuth,
		TokenExpirationTime:    time.Hour * time.Duration(tokenExpTime),
		PipelineWorkerPoolSize: numWorkers,
		CollectionNames:        NewCollections(),
	}
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
