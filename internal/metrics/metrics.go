package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	BuildInfo = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "app_build_info",
		Help: "Application build version and info",
	}, []string{"version", "app_name"})
	HttpRequestsTotal = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "http_requests_total",
		Help: "Total number of HTTP requests processed",
	}, []string{"method", "path", "status"})
	HttpRequestDuration = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Name:    "http_request_duration_seconds",
		Help:    "Duration of HTTP requests in seconds",
		Buckets: prometheus.DefBuckets,
	}, []string{"method", "path"})
)
