package configs

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

type RedisConfig struct {
	Client *redis.Client
}

func NewRedisConfig(cfg *Config) *RedisConfig {
	timeout, cancel := context.WithTimeout(context.Background(), 5*time.Second)

	defer cancel()

	if cfg.Env == "development" {
		cfg.RedisHost = "localhost"
	}

	client := redis.NewClient(&redis.Options{
		Addr: cfg.RedisHost + ":" + cfg.RedisPort,
	})

	//  test redis connection with a simple ping
	_, err := client.Ping(timeout).Result()

	if err != nil {
		zap.L().Fatal("Error pinging redis", zap.Error(err))
	}

	zap.L().Info("Connected to redis")

	return &RedisConfig{
		Client: client,
	}
}
