package server

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/yusirdemir/microservice/internal/repository"
	"github.com/yusirdemir/microservice/internal/repository/couchbase"
	"github.com/yusirdemir/microservice/internal/repository/memory"
	"github.com/yusirdemir/microservice/internal/service"
	"github.com/yusirdemir/microservice/internal/transport/http/gin/handler"
	"github.com/yusirdemir/microservice/internal/transport/http/gin/middleware"
	"github.com/yusirdemir/microservice/internal/transport/http/gin/router"
	"github.com/yusirdemir/microservice/pkg/config"
	"go.uber.org/zap"
)

type Server struct {
	App    *gin.Engine
	Server *http.Server
	Config *config.Config
	Logger *zap.Logger
}

func New(cfg *config.Config, logger *zap.Logger) (*Server, error) {
	if cfg.App.Env == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	gin.DefaultWriter = io.Discard

	app := gin.New()
	app.Use(gin.Recovery())
	app.Use(middleware.Logger(logger))
	app.Use(middleware.Metrics)

	app.GET("/metrics", gin.WrapH(promhttp.Handler()))

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
	s.Logger.Info("Initializing server...", zap.String("address", port))

	s.Server = &http.Server{
		Addr:    port,
		Handler: s.App,
	}

	s.Logger.Info("Server started successfully",
		zap.String("host", "localhost"),
		zap.String("port", s.Config.App.Port),
	)

	if err := s.Server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return err
	}
	return nil
}

func (s *Server) Shutdown() error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	return s.Server.Shutdown(ctx)
}
