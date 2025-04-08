package config

import (
	"context"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
)

var redisClient *redis.Client

func RedisConfigInit() {
	redisClient = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})
	err := redisClient.Ping(context.Background()).Err()
	if err != nil {
		logrus.Warn(err.Error())
	}
}

func GetRedisClient() *redis.Client {
	return redisClient
}
