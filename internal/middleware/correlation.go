package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

/*
🎯 CorrelationIDMiddleware - Setup:

1. Setiap request dapat unique ID (correlation_id)
2. ID tersimpan di context Gin
3. ID bisa di-pass ke logger untuk tracking
4. Response include correlation ID di header

Use case:
- Track request flow through multiple services
- Debug production issues dengan log ID
- Correlate logs, metrics, traces dengan request ID

Contoh flow:
Request A -> Service 1 -> Service 2 -> Service 3
Semua langkah punya ID yang sama, jadi bisa trace linearly!
*/

const CorrelationIDHeader = "X-Correlation-ID"

// CorrelationIDMiddleware - Add unique correlation ID to each request
func CorrelationIDMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Cek apakah client sudah provide correlation ID
		correlationID := c.GetHeader(CorrelationIDHeader)
		if correlationID == "" {
			// Generate baru jika belum ada
			correlationID = uuid.New().String()
		}

		// Store di context
		c.Set("correlation_id", correlationID)

		// Add ke response header
		c.Header(CorrelationIDHeader, correlationID)

		c.Next()
	}
}

// GetCorrelationID - Helper function untuk get correlation ID dari context
func GetCorrelationID(c *gin.Context) string {
	if id, exists := c.Get("correlation_id"); exists {
		if correlationID, ok := id.(string); ok {
			return correlationID
		}
	}
	return ""
}
