package services

import (
	"github.com/go-redis/redis"
)

var redisClient *redis.Client

func InitRedis() *redis.Client {
	// Connect to Redis
	redisClient := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})
	return redisClient
}
