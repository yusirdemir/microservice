package server

import (
	"os"
	"time"

	"github.com/gofiber/adaptor/v2"
	"github.com/gofiber/contrib/fiberzap/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/yusirdemir/microservice/internal/handler"
	"github.com/yusirdemir/microservice/internal/middleware"
	"github.com/yusirdemir/microservice/internal/router"
	"github.com/yusirdemir/microservice/pkg/config"
	"go.uber.org/zap"
)

type Server struct {
	App    *fiber.App
	Config *config.Config
	Logger *zap.Logger
}

func New(cfg *config.Config, logger *zap.Logger) *Server {
	app := fiber.New(fiber.Config{
		DisableStartupMessage: true,
		AppName:               cfg.App.Name,
	})

	app.Use(middleware.Metrics)

	app.Get("/metrics", adaptor.HTTPHandler(promhttp.Handler()))

	app.Use(fiberzap.New(fiberzap.Config{
		Logger: logger,
		Fields: []string{"latency", "status", "method", "url", "ip", "ua"},
	}))

	app.Hooks().OnListen(func(data fiber.ListenData) error {
		logger.Info("Server started successfully",
			zap.String("host", data.Host),
			zap.String("port", data.Port),
			zap.Int("pid", os.Getpid()),
		)
		return nil
	})

	handlers := []router.RouteHandler{
		handler.NewHealthHandler(),
		handler.NewUserHandler(),
	}

	r := router.New(app, handlers)
	r.SetupRoutes()

	return &Server{
		App:    app,
		Config: cfg,
		Logger: logger,
	}
}

func (s *Server) Run() error {
	port := ":" + s.Config.App.Port
	s.Logger.Info("Initializing server...", zap.String("address", port))
	return s.App.Listen(port)
}

func (s *Server) Shutdown() error {
	return s.App.ShutdownWithTimeout(10 * time.Second)
}
