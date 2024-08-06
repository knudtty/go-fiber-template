package config

import (
	"context"
	"log"
	"os"
	"strconv"

	"github.com/redis/go-redis/v9"
)

var Redis *redis.Client

func InitCache() {
	dbNum, _ := strconv.Atoi(os.Getenv("REDIS_DB_NUMBER"))

	Redis = redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_HOST") + ":" + os.Getenv("REDIS_PORT"),
		Password: os.Getenv("REDIS_PASSWORD"),
		DB:       dbNum, // use default DB
	})

	status := Redis.Ping(context.Background())
	if _, err := status.Result(); err != nil {
		log.Fatalln("Error connecting to redis:", err)
	}
}
