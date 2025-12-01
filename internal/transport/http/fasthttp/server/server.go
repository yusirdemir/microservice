package server

import (
	"fmt"
	"time"

	"github.com/fasthttp/router"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/valyala/fasthttp"
	"github.com/valyala/fasthttp/fasthttpadaptor"
	"github.com/yusirdemir/microservice/internal/repository"
	"github.com/yusirdemir/microservice/internal/repository/couchbase"
	"github.com/yusirdemir/microservice/internal/repository/memory"
	"github.com/yusirdemir/microservice/internal/service"
	"github.com/yusirdemir/microservice/internal/transport/http/fasthttp/handler"
	"github.com/yusirdemir/microservice/internal/transport/http/fasthttp/middleware"
	localRouter "github.com/yusirdemir/microservice/internal/transport/http/fasthttp/router"
	"github.com/yusirdemir/microservice/pkg/config"
	"go.uber.org/zap"
)

type Server struct {
	Server *fasthttp.Server
	Config *config.Config
	Logger *zap.Logger
}

func New(cfg *config.Config, logger *zap.Logger) (*Server, error) {
	r := router.New()

	r.GET("/metrics", fasthttpadaptor.NewFastHTTPHandler(promhttp.Handler()))

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

	handlers := []localRouter.RouteHandler{
		handler.NewUserHandler(userService),
		handler.NewHealthHandler(),
		handler.NewTimeoutHandler(),
	}

	lr := localRouter.New(r, handlers)
	lr.SetupRoutes()

	finalHandler := middleware.Metrics(r.Handler)

	loggedHandler := func(ctx *fasthttp.RequestCtx) {
		start := time.Now()
		finalHandler(ctx)

		path := string(ctx.Path())
		if path == "/health/live" || path == "/health/ready" || path == "/metrics" {
			return
		}

		end := time.Now()
		latency := end.Sub(start)

		fields := []zap.Field{
			zap.Int("status", ctx.Response.StatusCode()),
			zap.String("method", string(ctx.Method())),
			zap.String("path", path),
			zap.String("ip", ctx.RemoteIP().String()),
			zap.Duration("latency", latency),
		}

		logger.Info(path, fields...)
	}

	return &Server{
		Server: &fasthttp.Server{
			Handler: loggedHandler,
		},
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
	return s.Server.ListenAndServe(port)
}

func (s *Server) Shutdown() error {
	return s.Server.Shutdown()
}
