package middleware

import (
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	BuildInfo = promauto.NewGaugeVec(prometheus.GaugeOpts{
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

func Metrics(c *fiber.Ctx) error {
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
}
