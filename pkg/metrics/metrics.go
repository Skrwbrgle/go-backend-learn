// Package metrics provides Prometheus metrics instrumentation for the application
package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

/*
🎯 Metrics yang kita track:

1. HTTPRequestsTotal - Total HTTP requests by method, endpoint, status
2. HTTPRequestDuration - Durasi request processing
3. HTTPResponseSize - Ukuran response body
4. ActiveRequests - Request yang sedang berjalan (gauge)
5. DatabaseQueryDuration - Durasi database queries
6. DatabaseErrors - Database error count

Ini semua adalah standard metrics untuk HTTP API.
*/

var (
	// HTTPRequestsTotal - Total HTTP requests counter
	// Labels: method, endpoint, status
	HTTPRequestsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total number of HTTP requests (partitioned by method, endpoint, status code)",
		},
		[]string{"method", "endpoint", "status"},
	)

	// HTTPRequestDuration - Request duration histogram (dalam seconds)
	// Labels: method, endpoint
	HTTPRequestDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Help:    "HTTP request latencies in seconds (partitioned by method, endpoint)",
			Buckets: prometheus.DefBuckets, // 0.005s, 0.01s, 0.025s, 0.05s, 0.1s, 0.25s, 0.5s, 1s, 2.5s, 5s, 10s
		},
		[]string{"method", "endpoint"},
	)

	// HTTPResponseSize - Response size histogram dalam bytes
	// Labels: method, endpoint
	HTTPResponseSize = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_response_size_bytes",
			Help:    "HTTP response size in bytes (partitioned by method, endpoint)",
			Buckets: prometheus.ExponentialBuckets(100, 10, 8), // 100B, 1KB, 10KB, 100KB, 1MB, etc
		},
		[]string{"method", "endpoint", "status"},
	)

	// ActiveRequests - Gauge untuk request yang currently processing
	// Labels: method, endpoint
	ActiveRequests = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "http_requests_active",
			Help: "Number of currently active HTTP requests",
		},
		[]string{"method", "endpoint"},
	)

	// HTTPErrors - Total errors by type
	// Labels: method, endpoint, error_type
	HTTPErrors = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_errors_total",
			Help: "Total number of HTTP errors (partitioned by method, endpoint, error_type)",
		},
		[]string{"method", "endpoint", "error_type"},
	)

	// DatabaseQueryDuration - Database query duration
	// Labels: query_type (select, insert, update, delete)
	DatabaseQueryDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "db_query_duration_seconds",
			Help:    "Database query latencies in seconds (partitioned by query_type)",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"query_type"},
	)

	// DatabaseErrors - Database errors counter
	// Labels: query_type, error_type
	DatabaseErrors = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "db_errors_total",
			Help: "Total number of database errors",
		},
		[]string{"query_type", "error_type"},
	)
)

// RecordHTTPMetrics - Helper function to record HTTP metrics
// Call in middleware untuk auto-record semua metrics
func RecordHTTPMetrics(method, endpoint string, statusCode int, duration float64, responseSize int) {
	HTTPRequestsTotal.WithLabelValues(method, endpoint, string(rune(statusCode/100)+'0')).Inc()
	HTTPRequestDuration.WithLabelValues(method, endpoint).Observe(duration)
	HTTPResponseSize.WithLabelValues(method, endpoint, string(rune(statusCode/100)+'0')).Observe(float64(responseSize))
}

// RecordDatabaseMetrics - Helper function untuk database metrics
func RecordDatabaseMetrics(queryType string, duration float64, err error) {
	DatabaseQueryDuration.WithLabelValues(queryType).Observe(duration)
	if err != nil {
		DatabaseErrors.WithLabelValues(queryType, err.Error()).Inc()
	}
}
