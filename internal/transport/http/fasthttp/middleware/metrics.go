package middleware

import (
	"strconv"
	"time"

	"github.com/valyala/fasthttp"
	"github.com/yusirdemir/microservice/internal/metrics"
)

func Metrics(next fasthttp.RequestHandler) fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) {
		start := time.Now()

		next(ctx)

		duration := time.Since(start).Seconds()
		status := ctx.Response.StatusCode()
		path := string(ctx.Path())
		method := string(ctx.Method())

		metrics.HttpRequestsTotal.WithLabelValues(method, path, strconv.Itoa(status)).Inc()
		metrics.HttpRequestDuration.WithLabelValues(method, path).Observe(duration)
	}
}
