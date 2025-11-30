package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/yusirdemir/microservice/internal/middleware"
	"github.com/yusirdemir/microservice/internal/server"
	"github.com/yusirdemir/microservice/pkg/config"
	"github.com/yusirdemir/microservice/pkg/logger"
	"go.uber.org/zap"
)

var Version = "dev"

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		panic("Failed to load configuration: " + err.Error())
	}

	log, err := logger.New(cfg.Logger.Level)
	if err != nil {
		panic("Failed to initialize logger: " + err.Error())
	}
	defer func() {
		_ = log.Sync()
	}()

	log.Info("Starting application",
		zap.String("app", cfg.App.Name),
		zap.String("version", Version),
		zap.String("env", cfg.App.Env),
	)

	middleware.BuildInfo.WithLabelValues(Version, cfg.App.Name).Set(1)

	srv := server.New(cfg, log)

	go func() {
		if err := srv.Run(); err != nil {
			log.Fatal("Server failed to start", zap.Error(err))
		}
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	<-c

	log.Info("Shutdown signal received. Shutting down gracefully...")
	if err := srv.Shutdown(); err != nil {
		log.Error("Server forced to shutdown", zap.Error(err))
	} else {
		log.Info("Server shutdown successfully")
	}

	_ = log.Sync()
	log.Info("Goodbye!")
}
