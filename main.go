package main

import (
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/gofiber/adaptor/v2"
	"github.com/gofiber/contrib/fiberzap/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/yusirdemir/microservice/config"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// -----------------------------------------------------------------------------
// Version Variable
// -----------------------------------------------------------------------------
var Version = "dev"

// -----------------------------------------------------------------------------
// Global Metrics Definition
// -----------------------------------------------------------------------------
var (
	buildInfo = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "app_build_info",
		Help: "Application build version and info",
	}, []string{"version", "app_name"})
	httpRequestsTotal = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "http_requests_total",
		Help: "Total number of HTTP requests processed",
	}, []string{"method", "path", "status"})
	httpRequestDuration = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Name:    "http_request_duration_seconds",
		Help:    "Duration of HTTP requests in seconds",
		Buckets: prometheus.DefBuckets,
	}, []string{"method", "path"})
)

func main() {
	// -------------------------------------------------------------------------
	// Configuration Loading
	// -------------------------------------------------------------------------
	cfg, err := config.LoadConfig()
	if err != nil {
		panic("Failed to load configuration: " + err.Error())
	}

	// -------------------------------------------------------------------------
	// Logger Configuration (Zap)
	// -------------------------------------------------------------------------
	logConfig := zap.NewProductionConfig()
	logConfig.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	zapLevel, err := zapcore.ParseLevel(cfg.Logger.Level)
	if err != nil {
		panic("Failed to parse log level: " + err.Error())
	}
	logConfig.Level = zap.NewAtomicLevelAt(zapLevel)

	logger, err := logConfig.Build()
	if err != nil {
		panic("Failed to initialize logger: " + err.Error())
	}

	defer func() {
		_ = logger.Sync()
	}()

	logger.Info("Starting application",
		zap.String("app", cfg.App.Name),
		zap.String("version", Version),
		zap.String("env", cfg.App.Env),
	)

	buildInfo.WithLabelValues(Version, cfg.App.Name).Set(1)

	// -------------------------------------------------------------------------
	// Fiber Application Setup
	// -------------------------------------------------------------------------
	app := fiber.New(fiber.Config{
		DisableStartupMessage: true,
		AppName:               cfg.App.Name,
	})

	// -------------------------------------------------------------------------
	// Custom Metrics Middleware
	// -------------------------------------------------------------------------
	app.Use(func(c *fiber.Ctx) error {
		start := time.Now()

		err := c.Next()

		duration := time.Since(start).Seconds()

		path := c.Route().Path
		if path == "" {
			path = "/unknown-route"
		}

		method := c.Method()
		statusCode := c.Response().StatusCode()
		statusStr := strconv.Itoa(statusCode)

		httpRequestsTotal.WithLabelValues(method, path, statusStr).Inc()
		httpRequestDuration.WithLabelValues(method, path).Observe(duration)

		return err
	})

	// -------------------------------------------------------------------------
	// Metrics Endpoint
	// -------------------------------------------------------------------------
	app.Get("/metrics", adaptor.HTTPHandler(promhttp.Handler()))

	// -------------------------------------------------------------------------
	// Logging Middleware (FiberZap)
	// -------------------------------------------------------------------------
	app.Use(fiberzap.New(fiberzap.Config{
		Logger: logger,
		Fields: []string{"latency", "status", "method", "url", "ip", "ua"},
	}))

	// -------------------------------------------------------------------------
	// Lifecycle Hooks
	// -------------------------------------------------------------------------
	app.Hooks().OnListen(func(data fiber.ListenData) error {
		logger.Info("Server started successfully",
			zap.String("host", data.Host),
			zap.String("port", data.Port),
			zap.Int("pid", os.Getpid()),
		)
		return nil
	})

	// -------------------------------------------------------------------------
	// Application Routes
	// -------------------------------------------------------------------------
	app.Get("/", func(c *fiber.Ctx) error {
		time.Sleep(5 * time.Second)
		return c.SendString("System is running perfectly.")
	})

	app.Get("/users/:id", func(c *fiber.Ctx) error {
		id := c.Params("id")
		return c.SendString("User ID: " + id)
	})

	// -------------------------------------------------------------------------
	// Server Start (Background)
	// -------------------------------------------------------------------------
	go func() {
		port := ":" + cfg.App.Port
		logger.Info("Initializing server...", zap.String("address", port))

		if err := app.Listen(port); err != nil {
			logger.Fatal("Server failed to start", zap.Error(err))
		}
	}()

	// -------------------------------------------------------------------------
	// Graceful Shutdown Mechanism
	// -------------------------------------------------------------------------
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	<-c

	logger.Info("Shutdown signal received. Shutting down gracefully...")
	if err := app.ShutdownWithTimeout(10 * time.Second); err != nil {
		logger.Error("Server forced to shutdown", zap.Error(err))
	} else {
		logger.Info("Server shutdown successfully")
	}

	_ = logger.Sync()
	logger.Info("Goodbye!")
}
