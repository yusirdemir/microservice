package middleware

import (
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/yusirdemir/microservice/internal/metrics"
)

func Metrics(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		start := time.Now()

		err := next(c)

		duration := time.Since(start).Seconds()
		status := c.Response().Status
		path := c.Path()
		method := c.Request().Method

		metrics.HttpRequestsTotal.WithLabelValues(method, path, strconv.Itoa(status)).Inc()
		metrics.HttpRequestDuration.WithLabelValues(method, path).Observe(duration)

		return err
	}
}
