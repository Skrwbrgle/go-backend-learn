package middleware

import (
	"time"

	"github.com/Skrwbrgle/go-backend-learn/pkg/metrics"
	"github.com/gin-gonic/gin"
)

/*
🎯 MetricsMiddleware berfungsi:

1. Catat waktu mulai request
2. Catat response status & size
3. Calculate durasi
4. Record ke Prometheus
5. Lanjut ke handler berikutnya

Ini ideal karena:
- Non-invasive (no need ubah existing handlers)
- Accurate (capture real durations)
- Automatic (semua endpoints tercovered)
*/

// MetricsMiddleware - Middleware untuk record HTTP metrics
func MetricsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		startTime := time.Now()

		// Increment active requests
		metrics.ActiveRequests.WithLabelValues(c.Request.Method, c.Request.URL.Path).Inc()

		// Gunakan responseWriter wrapper untuk capture response size
		writer := &responseWriter{
			ResponseWriter: c.Writer,
			size:           0,
		}
		c.Writer = writer

		// Process request
		c.Next()

		// Decrement active requests
		metrics.ActiveRequests.WithLabelValues(c.Request.Method, c.Request.URL.Path).Dec()

		// Calculate duration
		duration := time.Since(startTime).Seconds()

		// Record metrics
		metrics.HTTPRequestsTotal.WithLabelValues(
			c.Request.Method,
			c.Request.URL.Path,
			string(rune(c.Writer.Status()/100))+"xx", // Normalize: 200->2xx, 500->5xx
		).Inc()

		metrics.HTTPRequestDuration.WithLabelValues(
			c.Request.Method,
			c.Request.URL.Path,
		).Observe(duration)

		metrics.HTTPResponseSize.WithLabelValues(
			c.Request.Method,
			c.Request.URL.Path,
			string(rune(c.Writer.Status()/100))+"xx",
		).Observe(float64(writer.size))

		// Record errors (4xx, 5xx)
		if c.Writer.Status() >= 400 {
			errorType := "client"
			if c.Writer.Status() >= 500 {
				errorType = "server"
			}
			metrics.HTTPErrors.WithLabelValues(
				c.Request.Method,
				c.Request.URL.Path,
				errorType,
			).Inc()
		}
	}
}

// responseWriter wrapper untuk capture response size
type responseWriter struct {
	gin.ResponseWriter
	size int
}

func (w *responseWriter) Write(b []byte) (int, error) {
	w.size += len(b)
	return w.ResponseWriter.Write(b)
}

func (w *responseWriter) WriteString(s string) (int, error) {
	w.size += len(s)
	return w.ResponseWriter.WriteString(s)
}
