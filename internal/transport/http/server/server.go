package server

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/gofiber/adaptor/v2"
	"github.com/gofiber/contrib/fiberzap/v2"
	"github.com/gofiber/contrib/otelfiber"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/timeout"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/yusirdemir/microservice/internal/repository"
	"github.com/yusirdemir/microservice/internal/repository/couchbase"
	"github.com/yusirdemir/microservice/internal/repository/memory"
	"github.com/yusirdemir/microservice/internal/service"
	"github.com/yusirdemir/microservice/internal/transport/http/handler"
	"github.com/yusirdemir/microservice/internal/transport/http/middleware"
	"github.com/yusirdemir/microservice/internal/transport/http/router"
	"github.com/yusirdemir/microservice/pkg/config"
	"github.com/yusirdemir/microservice/pkg/telemetry"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
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

	logLevel, err := zapcore.ParseLevel(cfg.Logger.Level)
	if err != nil {
		logLevel = zapcore.InfoLevel
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
		Levels: []zapcore.Level{zapcore.ErrorLevel, zapcore.WarnLevel, logLevel},
		Next: func(c *fiber.Ctx) bool {
			return c.Path() == "/health/live" || c.Path() == "/health/ready"
		},
	}))

	app.Hooks().OnListen(func(data fiber.ListenData) error {
		logger.Info("Server started successfully",
			zap.String("host", data.Host),
			zap.String("port", data.Port),
			zap.Int("pid", os.Getpid()),
		)
		return nil
	})

	var userRepo repository.UserRepository
	var productRepo repository.ProductRepository
	var errRepo error

	switch cfg.Database.Driver {
	case "couchbase":
		userRepo, errRepo = couchbase.NewUserRepository(cfg)
		if errRepo == nil {
			productRepo, errRepo = couchbase.NewProductRepository(cfg)
		}
	default:
		userRepo = memory.NewUserRepository()
		productRepo = memory.NewProductRepository()
	}

	if errRepo != nil {
		return nil, fmt.Errorf("failed to initialize repository: %w", errRepo)
	}

	userService := service.NewUserService(userRepo)
	productService := service.NewProductService(productRepo)

	handlers := []router.RouteHandler{
		handler.NewUserHandler(userService),
		handler.NewProductHandler(productService),
		handler.NewHealthHandler(),
		handler.NewTimeoutHandler(),
	}

	r := router.New(app, handlers)
	r.SetupRoutes()

	tracer, err := telemetry.InitTracer(cfg.App.Name, "1.0.0", cfg.App.Env, cfg.Trace.Endpoint)
	if err != nil {
		logger.Error("Failed to init tracer", zap.Error(err))
	} else {
		app.Hooks().OnShutdown(func() error {
			return tracer.Shutdown(context.Background())
		})

		app.Use(otelfiber.Middleware())
		logger.Info("OpenTelemetry tracer initialized and middleware added")
	}

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
