package config

import (
	"context"
	"fmt"
	"github.com/SwanHtetAungPhyo/common/pkg/logutil"
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
	logutil.GetLogger().Warn(fmt.Sprintf("Failed to load config: %v", err))
	if err != nil {
		logrus.Warn(err.Error())
	}
}

func GetRedisClient() *redis.Client {
	return redisClient
}
