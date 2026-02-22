package shutdown

import (
	"context"

	"go.uber.org/zap"
)

type CleanupFunc func() error

func Wait(ctx context.Context, logger *zap.Logger, cleanup CleanupFunc) {
	<-ctx.Done()

	logger.Info("Shutdown signal received. Shutting down gracefully...")

	if err := cleanup(); err != nil {
		logger.Error("Cleanup failed", zap.Error(err))
	} else {
		logger.Info("Cleanup completed successfully")
	}
}
