package redis

import (
	"context"

	"github.com/redis/go-redis/v9"
)

var ctx = context.Background()

type Client struct {
	RedisClient *redis.Client
}

func NewClient(redisAddr string) *Client {
	client := redis.NewClient(&redis.Options{
		Addr: redisAddr,
	})

	// Test the connection
	_, err := client.Ping(ctx).Result()
	if err != nil {
		panic(err) // Handle this as needed
	}

	return &Client{RedisClient: client}
}
