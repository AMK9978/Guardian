package configs

import (
	"github.com/MicahParks/keyfunc"
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

	secretKey := getEnv("JWT_SECRET_KEY", "")
	tokenAuth := jwtauth.New("HS256", []byte(secretKey), nil)
	tokenExpTimeStr := getEnv("TOKEN_EXP_TIME", "72")
	tokenExpTime, err := strconv.Atoi(tokenExpTimeStr)

	activationSecretKey := getEnv("ACTIVATION_SECRET_KEY", "")
	activationTokenExpTimeStr := getEnv("ACTIVATION_TOKEN_EXP_TIME", "72")
	activationTokenExpTime, err := strconv.Atoi(activationTokenExpTimeStr)
	if err != nil {
		logger.GetLogger().Fatalf("couldn't convert the token expiration time to int: %s", tokenExpTimeStr)
	}

	rateLimiterStatus, err := strconv.ParseBool(getEnv("RATE_LIMITER_STATUS", "false"))
	if err != nil {
		logger.GetLogger().Fatal("couldn't convert the rate limiter status to bool")
	}

	externalAuthStatus, err := strconv.ParseBool(getEnv("EXTERNAL_AUTH_STATUS", "false"))
	if err != nil {
		logger.GetLogger().Fatal("couldn't convert the external auth status to bool")
	}

	rateInterval, requestLimit := -1, -1
	if rateLimiterStatus {
		requestLimit, err = strconv.Atoi(getEnv("REQUEST_LIMIT", "10"))
		if err != nil {
			logger.GetLogger().Fatal("couldn't convert the request limit to int")
		}
		rateInterval, err = strconv.Atoi(getEnv("RATE_INTERVAL", "1"))
		if err != nil {
			logger.GetLogger().Fatal("couldn't convert the rate interval time to int")
		}
	}

	var jwks *keyfunc.JWKS
	jwksURL := getEnv("JWKS_URL", "")

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
		ServerPort:             getEnv("SERVER_PORT", "8080"),
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
