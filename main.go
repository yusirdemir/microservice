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
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// -----------------------------------------------------------------------------
// Global Metrics Definition
// -----------------------------------------------------------------------------
var (
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
	// 1. Logger Configuration (Zap)
	// -------------------------------------------------------------------------
	config := zap.NewProductionConfig()
	config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	logger, err := config.Build()
	if err != nil {
		panic("Failed to initialize logger: " + err.Error())
	}
	defer logger.Sync()

	// -------------------------------------------------------------------------
	// 2. Fiber Application Setup
	// -------------------------------------------------------------------------
	app := fiber.New(fiber.Config{
		DisableStartupMessage: true,
		AppName:               "User Microservice v1",
	})

	// -------------------------------------------------------------------------
	// 3. Custom Metrics Middleware
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
	// 4. Metrics Endpoint
	// -------------------------------------------------------------------------
	app.Get("/metrics", adaptor.HTTPHandler(promhttp.Handler()))

	// -------------------------------------------------------------------------
	// 5. Logging Middleware (FiberZap)
	// -------------------------------------------------------------------------
	app.Use(fiberzap.New(fiberzap.Config{
		Logger: logger,
		Fields: []string{"latency", "status", "method", "url", "ip", "ua"},
	}))

	// -------------------------------------------------------------------------
	// 6. Lifecycle Hooks
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
	// 7. Application Routes
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
	// 8. Server Start
	// -------------------------------------------------------------------------
	go func() {
		logger.Info("Starting server...", zap.String("port", ":3000"))
		if err := app.Listen(":3000"); err != nil {
			logger.Fatal("Server failed to start", zap.Error(err))
		}
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c
	logger.Info("Shutting down gracefully...")

	if err := app.ShutdownWithTimeout(10 * time.Second); err != nil {
		logger.Error("Failed to shutdown", zap.Error(err))
	} else {
		logger.Info("Server shutdown successfully")
	}
	logger.Sync()
	logger.Info("Goodbye!")
}
