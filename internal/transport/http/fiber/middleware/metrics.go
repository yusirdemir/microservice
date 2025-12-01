package middleware

import (
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/yusirdemir/microservice/internal/metrics"
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

	metrics.HttpRequestsTotal.WithLabelValues(method, path, statusStr).Inc()
	metrics.HttpRequestDuration.WithLabelValues(method, path).Observe(duration)

	return err
}
