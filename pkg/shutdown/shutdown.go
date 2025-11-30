package shutdown

import (
	"os"
	"os/signal"
	"syscall"

	"go.uber.org/zap"
)

type CleanupFunc func() error

func Wait(logger *zap.Logger, cleanup CleanupFunc) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	<-c

	logger.Info("Shutdown signal received. Shutting down gracefully...")

	if err := cleanup(); err != nil {
		logger.Error("Cleanup failed", zap.Error(err))
	} else {
		logger.Info("Cleanup completed successfully")
	}
}
