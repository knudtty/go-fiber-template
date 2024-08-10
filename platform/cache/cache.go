package cache

import (
	"context"
	"fmt"
	"os"
	"strconv"

	"github.com/redis/go-redis/v9"
)

var client *redis.Client

func RedisClient() (*redis.Client, error) {
	if client != nil {
		return client, nil
	}

	dbNum, _ := strconv.Atoi(os.Getenv("REDIS_DB_NUMBER"))
	client = redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_HOST") + ":" + os.Getenv("REDIS_PORT"),
		Password: os.Getenv("REDIS_PASSWORD"),
		DB:       dbNum, // use default DB
	})

	status := client.Ping(context.Background())
	if _, err := status.Result(); err != nil {
		defer client.Close()
		return nil, fmt.Errorf("Error connecting to redis: %s", err)
	}

	return client, nil
}
