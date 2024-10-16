package metrics

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	// Define a new counter metric
	requestCounter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total number of HTTP requests",
		},
		[]string{"method", "handler"},
	)

	// Define a new histogram metric
	requestDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Help:    "Duration of HTTP requests in seconds",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "handler"},
	)
)

func Init() {
	// Register the metrics with Prometheus
	prometheus.MustRegister(requestCounter)
	prometheus.MustRegister(requestDuration)
}

// Handler for exposing the metrics
func Handler() http.Handler {
	return promhttp.Handler()
}
