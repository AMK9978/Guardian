package configs

import (
	"github.com/MicahParks/keyfunc"
	"github.com/go-chi/jwtauth/v5"
	"github.com/spf13/viper"

	"guardian/utlis/logger"
	"log"
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
	User         string
	Task         string
	Group        string
	TargetModel  string
	RefereeModel string
}

// NewCollections initializes the collection names.
func NewCollections() *Collections {
	return &Collections{
		User:         "users",
		Task:         "tasks",
		Group:        "groups",
		TargetModel:  "target_models",
		RefereeModel: "referee_models",
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
	ActivationTokenKey     string
	TokenExpirationTime    time.Duration
	ActivationTokenExpTime time.Duration
	EnableRateLimiter      bool
	RequestLimit           int
	Interval               time.Duration
	Jwk                    *keyfunc.JWKS
	ExternalJwtIssuer      string
	ExternalJwtAudience    string
	EnableExternalAuth     bool
}

func LoadConfig() Config {
	numWorkersStr := getEnv("PIPELINE_WORKER_POOL_SIZE", strconv.Itoa(runtime.NumCPU()))
	numWorkers, err := strconv.Atoi(numWorkersStr)
	if err != nil {
		logger.GetLogger().Fatalf("coudn't convert the worker pool size to int: %s", numWorkersStr)
	}

	viper.AddConfigPath(".")
	viper.SetConfigName(".env")
	viper.SetConfigType("yaml")
	err = viper.ReadInConfig()

	if err != nil {
		log.Fatalf("Error while reading config file %s", err)
	} else {

	}

	secretKey := viper.GetString("JWT_SECRET_KEY")
	tokenAuth := jwtauth.New("HS256", []byte(secretKey), nil)

	viper.SetDefault("TOKEN_EXP_TIME", 72)
	viper.SetDefault("ACTIVATION_TOKEN_EXP_TIME", 72)
	viper.SetDefault("RATE_LIMITER_STATUS", false)
	viper.SetDefault("EXTERNAL_AUTH_STATUS", false)
	viper.SetDefault("REQUEST_LIMIT", 10)
	viper.SetDefault("RATE_INTERVAL", 1)

	tokenExpTime := viper.GetInt("TOKEN_EXP_TIME")

	activationSecretKey := viper.GetString("ACTIVATION_SECRET_KEY")
	activationTokenExpTime := viper.GetInt("ACTIVATION_TOKEN_EXP_TIME")
	rateLimiterStatus := viper.GetBool("RATE_LIMITER_STATUS")

	externalAuthStatus := viper.GetBool("EXTERNAL_AUTH_STATUS")

	rateInterval, requestLimit := -1, -1
	if rateLimiterStatus {
		requestLimit = viper.GetInt("REQUEST_LIMIT")
		rateInterval = viper.GetInt("RATE_INTERVAL")
	}

	var jwks *keyfunc.JWKS
	jwksURL := viper.GetString("JWKS_URL")

	if jwksURL != "" {
		jwks, err = keyfunc.Get(jwksURL, keyfunc.Options{})
		if err != nil {
			logger.GetLogger().Fatalf("Failed to create JWKS from URL: %s\n", err)
		}
	}

	return Config{
		RedisAddr:              getEnv("REDIS_ADDR", "localhost:6379"),
		MongoDBURI:             getEnv("MONGODB_URI", "mongodb://localhost:27017"),
		RabbitMQURI:            getEnv("RABBITMQ_URI", "amqp://guest:guest@localhost:5672/"),
		MilvusURI:              getEnv("MILVUS_URI", "localhost:19530"),
		ServerPort:             getEnv("SERVER_PORT", "8081"),
		PrimaryDBName:          getEnv("PRIMARY_DB_NAME", "primary"),
		TokenAuth:              tokenAuth,
		ActivationTokenKey:     activationSecretKey,
		TokenExpirationTime:    time.Hour * time.Duration(tokenExpTime),
		ActivationTokenExpTime: time.Hour * time.Duration(activationTokenExpTime),
		PipelineWorkerPoolSize: numWorkers,
		CollectionNames:        NewCollections(),
		EnableRateLimiter:      rateLimiterStatus,
		Interval:               time.Minute * time.Duration(rateInterval),
		RequestLimit:           requestLimit,
		Jwk:                    jwks,
		ExternalJwtIssuer:      getEnv("EXTERNAL_JWT_ISSUER", ""),
		ExternalJwtAudience:    getEnv("EXTERNAL_JWT_AUDIENCE", ""),
		EnableExternalAuth:     externalAuthStatus,
	}
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
