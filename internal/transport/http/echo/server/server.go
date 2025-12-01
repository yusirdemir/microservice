package server

import (
	"context"
	"fmt"
	"time"

	"github.com/labstack/echo/v4"
	echoMiddleware "github.com/labstack/echo/v4/middleware"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/yusirdemir/microservice/internal/repository"
	"github.com/yusirdemir/microservice/internal/repository/couchbase"
	"github.com/yusirdemir/microservice/internal/repository/memory"
	"github.com/yusirdemir/microservice/internal/service"
	"github.com/yusirdemir/microservice/internal/transport/http/echo/handler"
	"github.com/yusirdemir/microservice/internal/transport/http/echo/middleware"
	"github.com/yusirdemir/microservice/internal/transport/http/echo/router"
	"github.com/yusirdemir/microservice/pkg/config"
	"go.uber.org/zap"
)

type Server struct {
	App    *echo.Echo
	Config *config.Config
	Logger *zap.Logger
}

func New(cfg *config.Config, logger *zap.Logger) (*Server, error) {
	app := echo.New()
	app.HideBanner = true
	app.HidePort = true

	app.Use(echoMiddleware.Recover())
	app.Use(middleware.Metrics)

	app.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			start := time.Now()
			err := next(c)

			path := c.Request().URL.Path
			if path == "/health/live" || path == "/health/ready" || path == "/metrics" {
				return err
			}

			end := time.Now()
			latency := end.Sub(start)

			fields := []zap.Field{
				zap.Int("status", c.Response().Status),
				zap.String("method", c.Request().Method),
				zap.String("path", path),
				zap.String("ip", c.RealIP()),
				zap.Duration("latency", latency),
			}

			if err != nil {
				c.Error(err)
				logger.Error(err.Error(), fields...)
			} else {
				logger.Info(path, fields...)
			}
			return nil
		}
	})

	app.GET("/metrics", echo.WrapHandler(promhttp.Handler()))

	var userRepo repository.UserRepository
	var errRepo error

	switch cfg.Database.Driver {
	case "couchbase":
		userRepo, errRepo = couchbase.NewUserRepository(cfg)
	case "memory":
		userRepo = memory.NewUserRepository()
	default:
		userRepo = memory.NewUserRepository()
	}

	if errRepo != nil {
		return nil, fmt.Errorf("failed to initialize repository: %w", errRepo)
	}

	userService := service.NewUserService(userRepo)

	handlers := []router.RouteHandler{
		handler.NewUserHandler(userService),
		handler.NewHealthHandler(),
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
	s.Logger.Info("Server started successfully",
		zap.String("host", "localhost"),
		zap.String("port", s.Config.App.Port),
	)
	return s.App.Start(port)
}

func (s *Server) Shutdown() error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	return s.App.Shutdown(ctx)
}
