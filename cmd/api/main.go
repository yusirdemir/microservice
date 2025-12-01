package main

import (
	"os"

	"github.com/yusirdemir/microservice/internal/metrics"
	"github.com/yusirdemir/microservice/internal/server"
	"github.com/yusirdemir/microservice/pkg/config"
	"github.com/yusirdemir/microservice/pkg/logger"
	"github.com/yusirdemir/microservice/pkg/shutdown"
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
		zap.String("framework", os.Getenv("APP_FRAMEWORK")),
	)

	metrics.BuildInfo.WithLabelValues(Version, cfg.App.Name).Set(1)

	srv, err := server.New(cfg, log)
	if err != nil {
		log.Fatal("Failed to initialize server", zap.Error(err))
	}

	go func() {
		if err := srv.Run(); err != nil {
			log.Fatal("Server failed to start", zap.Error(err))
		}
	}()

	shutdown.Wait(log, func() error {
		return srv.Shutdown()
	})

	_ = log.Sync()
	log.Info("Goodbye!")
}
