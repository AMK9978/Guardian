package redis

import (
	"context"

	"github.com/redis/go-redis/v9"
)

var ctx = context.Background()

var Client *redis.Client

func NewClient(redisAddr string) *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr: redisAddr,
	})

	// Test the connection
	_, err := client.Ping(ctx).Result()
	if err != nil {
		panic(err) // Handle this as needed
	}

	return client
}

func Init(redisAddr string) {
	Client = NewClient(redisAddr)
}
