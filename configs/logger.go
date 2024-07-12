package configs

import "go.uber.org/zap"

func NewLogger() {
	logger, _ := zap.NewProduction()

	defer logger.Sync()

	logger.Info("Logger is ready")

	// set logger as global
	zap.ReplaceGlobals(logger)
}
