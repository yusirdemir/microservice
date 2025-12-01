package server

import (
	"os"
	"time"

	"github.com/gofiber/adaptor/v2"
	"github.com/gofiber/contrib/fiberzap/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/timeout"
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

func New(cfg *config.Config, logger *zap.Logger) (*Server, error) {
	readTimeout, err := time.ParseDuration(cfg.Server.ReadTimeout)
	if err != nil {
		return nil, err
	}
	writeTimeout, err := time.ParseDuration(cfg.Server.WriteTimeout)
	if err != nil {
		return nil, err
	}
	idleTimeout, err := time.ParseDuration(cfg.Server.IdleTimeout)
	if err != nil {
		return nil, err
	}

	app := fiber.New(fiber.Config{
		DisableStartupMessage: true,
		AppName:               cfg.App.Name,
		ReadTimeout:           readTimeout,
		WriteTimeout:          writeTimeout,
		IdleTimeout:           idleTimeout,
	})

	app.Use(func(c *fiber.Ctx) error {
		h := timeout.NewWithContext(func(c *fiber.Ctx) error {
			return c.Next()
		}, writeTimeout)
		return h(c)
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
		handler.NewTimeoutHandler(),
	}

	r := router.New(app, handlers)
	r.SetupRoutes()

	return &Server{
		App:    app,
		Config: cfg,
		Logger: logger,
	}, nil
}

func (s *Server) Run() error {
	port := ":" + s.Config.App.Port
	s.Logger.Info("Initializing server...", zap.String("address", port))
	return s.App.Listen(port)
}

func (s *Server) Shutdown() error {
	return s.App.ShutdownWithTimeout(10 * time.Second)
}
