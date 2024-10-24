package ratelimit

import (
	"github.com/go-chi/jwtauth/v5"
	"github.com/redis/go-redis/v9"
	"guardian/configs"
	"net/http"
	"strconv"
)

const (
	rateLimitKeyPrefix = "rate_limiter:"
)

func RateLimiterMiddleware(redisClient *redis.Client) func(handler http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			_, claims, _ := jwtauth.FromContext(r.Context())
			userID, ok := claims["user_id"].(string)
			if !ok {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}
			redisKey := rateLimitKeyPrefix + userID

			currentCountStr, err := redisClient.Get(r.Context(), redisKey).Result()
			if err != nil && err != redis.Nil {
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				return
			}
			currentCount := 0
			if currentCountStr != "" {
				currentCount, _ = strconv.Atoi(currentCountStr)
			}

			if currentCount >= configs.GlobalConfig.RequestLimit {
				http.Error(w, "Too Many Requests", http.StatusTooManyRequests)
				return
			}

			err = redisClient.Incr(r.Context(), redisKey).Err()
			if err != nil {
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				return
			}

			if currentCount == 0 {
				redisClient.Expire(r.Context(), redisKey, configs.GlobalConfig.Interval)
			}
			next.ServeHTTP(w, r)
		})
	}
}
