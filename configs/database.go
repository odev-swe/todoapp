package configs

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

type DbConfig struct {
	Db *pgxpool.Pool
}

func NewDbConfig(cfg *Config) *DbConfig {
	timeout, cancel := context.WithTimeout(context.Background(), 5*time.Second)

	defer cancel()

	if cfg.Env == "development" {
		cfg.DbHost = "localhost"
	}

	connString := fmt.Sprintf("postgresql://%s:%s@%s:%s/%s", cfg.DbUser, cfg.DbPassword, cfg.DbHost, cfg.DbPort, cfg.DbName)

	dbCfg, err := pgxpool.ParseConfig(connString)

	if err != nil {
		zap.L().Fatal("Error parsing connection string", zap.Error(err))
	}

	conn, err := pgxpool.NewWithConfig(timeout, dbCfg)

	if err != nil {
		zap.L().Fatal("Error connecting to database", zap.Error(err))
	}

	// test database connection with a simple ping
	err = conn.Ping(timeout)

	if err != nil {
		zap.L().Fatal("Error pinging database", zap.Error(err))
	}

	zap.L().Info("Connected to database")

	return &DbConfig{
		Db: conn,
	}
}
