package configs

import (
	"guardian/prompt_api"
	"log"
	"runtime"
	"time"

	"guardian/utlis/logger"

	"github.com/MicahParks/keyfunc"
	"github.com/go-chi/jwtauth/v5"
	"github.com/spf13/viper"
)

var GlobalConfig Config

type Collections struct {
	User         string
	Task         string
	Group        string
	TargetModel string
	Plugin      string
}

// NewCollections initializes the collection names.
func NewCollections() *Collections {
	return &Collections{
		User:        "users",
		Task:        "tasks",
		Group:       "groups",
		TargetModel: "target_models",
		Plugin:      "plugins",
	}
}

type Config struct {
	RedisAddr              string
	MongoDBURI             string
	RabbitMQURI            string
	MilvusURI              string
	ServerPort             int
	MetricServerPort       int
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
	HttpClientTimeout time.Duration
	GRPCManager       *prompt_api.ClientManager
}

func LoadConfig() Config {
	viper.AddConfigPath(".")
	viper.SetConfigName(".env")
	viper.SetConfigType("yaml")
	err := viper.ReadInConfig()
	if err != nil {
		log.Fatalf("Error while reading config file %s", err)
	}

	viper.SetDefault("PIPELINE_WORKER_POOL_SIZE", runtime.NumCPU())

	viper.SetDefault("REDIS_ADDR", "localhost:6379")
	viper.SetDefault("MONGODB_URI", "mongodb://localhost:27017")
	viper.SetDefault("RABBITMQ_URI", "amqp://guest:guest@localhost:5672/")
	viper.SetDefault("MILVUS_URI", "localhost:19530")

	viper.SetDefault("SERVER_PORT", 8080)
	viper.SetDefault("METRIC_SERVER_PORT", 8081)
	viper.SetDefault("PRIMARY_DB_NAME", "primary")

	viper.SetDefault("TOKEN_EXP_TIME", 72)
	viper.SetDefault("ACTIVATION_TOKEN_EXP_TIME", 72)
	viper.SetDefault("RATE_LIMITER_STATUS", false)
	viper.SetDefault("EXTERNAL_AUTH_STATUS", false)
	viper.SetDefault("REQUEST_LIMIT", 10)
	viper.SetDefault("RATE_INTERVAL", 1)

	viper.SetDefault("EXTERNAL_JWT_ISSUER", "")
	viper.SetDefault("EXTERNAL_JWT_AUDIENCE", "")

	viper.SetDefault("HTTP_CLIENT_TIMEOUT", 10)

	secretKey := viper.GetString("JWT_SECRET_KEY")
	tokenAuth := jwtauth.New("HS256", []byte(secretKey), nil)
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
		RedisAddr:              viper.GetString("REDIS_ADDR"),
		MongoDBURI:             viper.GetString("MONGODB_URI"),
		RabbitMQURI:            viper.GetString("RABBITMQ_URI"),
		MilvusURI:              viper.GetString("MILVUS_URI"),
		ServerPort:             viper.GetInt("SERVER_PORT"),
		MetricServerPort:       viper.GetInt("METRIC_SERVER_PORT"),
		PrimaryDBName:          viper.GetString("PRIMARY_DB_NAME"),
		TokenAuth:              tokenAuth,
		ActivationTokenKey:     activationSecretKey,
		TokenExpirationTime:    time.Hour * time.Duration(tokenExpTime),
		ActivationTokenExpTime: time.Hour * time.Duration(activationTokenExpTime),
		PipelineWorkerPoolSize: viper.GetInt("PIPELINE_WORKER_POOL_SIZE"),
		CollectionNames:        NewCollections(),
		EnableRateLimiter:      rateLimiterStatus,
		Interval:               time.Minute * time.Duration(rateInterval),
		RequestLimit:           requestLimit,
		Jwk:                    jwks,
		ExternalJwtIssuer:      viper.GetString("EXTERNAL_JWT_ISSUER"),
		ExternalJwtAudience:    viper.GetString("EXTERNAL_JWT_AUDIENCE"),
		EnableExternalAuth:     externalAuthStatus,
		HttpClientTimeout:      time.Duration(viper.GetInt("HTTP_CLIENT_TIMEOUT")) * time.Second,
		GRPCManager:            prompt_api.NewClientManager(),
	}
}
