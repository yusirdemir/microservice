package middleware

import (
	"net/http"
	"strconv"
	"time"

	"github.com/yusirdemir/microservice/internal/metrics"
)

type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

func Metrics(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		rw := &responseWriter{ResponseWriter: w, statusCode: http.StatusOK}
		next.ServeHTTP(rw, r)

		duration := time.Since(start).Seconds()
		path := r.URL.Path
		method := r.Method

		metrics.HttpRequestsTotal.WithLabelValues(method, path, strconv.Itoa(rw.statusCode)).Inc()
		metrics.HttpRequestDuration.WithLabelValues(method, path).Observe(duration)
	})
}
