package main

import (
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/odev-swe/todoapp/configs"
	"github.com/odev-swe/todoapp/internal/ratelimiter"
	"github.com/redis/go-redis/v9"
)

type application struct {
	limiter ratelimiter.RateLimiter
	config  configs.Config
	db      *pgxpool.Pool
	redis   *redis.Client
}

func main() {
	// logger
	configs.NewLogger()

	// envs
	envConfig := configs.NewEnv()

	// database
	dbConfig := configs.NewDbConfig(envConfig)
	db := dbConfig.Db

	// redis
	redisConfig := configs.NewRedisConfig(envConfig)
	redis := redisConfig.Client

	// application
	app := &application{
		limiter: ratelimiter.NewFixedWindowLimiter(envConfig.RateLimit, time.Duration(envConfig.RateLimitWindow)),
		config:  *envConfig,
		db:      db,
		redis:   redis,
	}

	// start http server
	app.Start()

}
