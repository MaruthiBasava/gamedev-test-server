package main

import "github.com/go-redis/redis"

var redisClient = InitializeRedisClient()

// InitializeRedisClient is
func InitializeRedisClient() *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr:       "localhost:6379",
		PoolSize:   100,
		MaxRetries: 2,
		Password:   "",
		DB:         0,
	})
	_, err := client.Ping().Result()

	if err != nil {
		panic(err)
	}

	return client
}