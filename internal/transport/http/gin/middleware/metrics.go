package middleware

import (
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/yusirdemir/microservice/internal/metrics"
)

func Metrics(c *gin.Context) {
	start := time.Now()

	c.Next()

	duration := time.Since(start).Seconds()

	path := c.FullPath()
	if path == "" {
		path = "/unknown-route"
	}

	method := c.Request.Method
	statusCode := c.Writer.Status()
	statusStr := strconv.Itoa(statusCode)

	metrics.HttpRequestsTotal.WithLabelValues(method, path, statusStr).Inc()
	metrics.HttpRequestDuration.WithLabelValues(method, path).Observe(duration)
}
