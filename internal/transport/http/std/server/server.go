package server

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/yusirdemir/microservice/internal/repository"
	"github.com/yusirdemir/microservice/internal/repository/couchbase"
	"github.com/yusirdemir/microservice/internal/repository/memory"
	"github.com/yusirdemir/microservice/internal/service"
	"github.com/yusirdemir/microservice/internal/transport/http/std/handler"
	"github.com/yusirdemir/microservice/internal/transport/http/std/middleware"
	"github.com/yusirdemir/microservice/internal/transport/http/std/router"
	"github.com/yusirdemir/microservice/pkg/config"
	"go.uber.org/zap"
)

type Server struct {
	Server *http.Server
	Config *config.Config
	Logger *zap.Logger
}

func New(cfg *config.Config, logger *zap.Logger) (*Server, error) {
	mux := http.NewServeMux()

	mux.Handle("GET /metrics", promhttp.Handler())

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

	r := router.New(mux, handlers)
	r.SetupRoutes()

	finalHandler := middleware.Metrics(mux)

	loggedHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		rw := &responseWriter{ResponseWriter: w, statusCode: http.StatusOK}
		finalHandler.ServeHTTP(rw, r)

		path := r.URL.Path
		if path == "/health/live" || path == "/health/ready" || path == "/metrics" {
			return
		}

		end := time.Now()
		latency := end.Sub(start)

		fields := []zap.Field{
			zap.Int("status", rw.statusCode),
			zap.String("method", r.Method),
			zap.String("path", path),
			zap.String("ip", r.RemoteAddr),
			zap.Duration("latency", latency),
		}

		logger.Info(path, fields...)
	})

	return &Server{
		Server: &http.Server{
			Addr:    ":" + cfg.App.Port,
			Handler: loggedHandler,
		},
		Config: cfg,
		Logger: logger,
	}, nil
}

type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

func (s *Server) Run() error {
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
