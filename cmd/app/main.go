package main

import (
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/odev-swe/todoapp/configs"
	"github.com/odev-swe/todoapp/internal/ratelimiter"
)

type application struct {
	limiter ratelimiter.RateLimiter
	config  configs.Config
	db      *pgxpool.Pool
}

func main() {
	// logger
	configs.NewLogger()

	// envs
	envConfig := configs.NewEnv()

	// database
	dbConfig := configs.NewDbConfig(envConfig)
	db := dbConfig.Db

	// application
	app := &application{
		limiter: ratelimiter.NewFixedWindowLimiter(envConfig.RateLimit, time.Duration(envConfig.RateLimitWindow)),
		config:  *envConfig,
		db:      db,
	}

	// start http server
	app.Start()

}
