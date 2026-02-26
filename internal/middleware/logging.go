package middleware

import (
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func LoggerMiddleware(logger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		method := c.Request.Method
		clientIP := c.ClientIP()
		userAgent := c.Request.UserAgent()
		correlationID := GetCorrelationID(c) // Get correlation ID dari context

		c.Next()

		status := c.Writer.Status()
		duration := time.Since(start)

		fields := []zap.Field{
			zap.String("correlation_id", correlationID), // Add correlation ID ke logs
			zap.String("method", method),
			zap.String("path", path),
			zap.Int("status", status),
			zap.String("ip", clientIP),
			zap.String("user-agent", userAgent),
			zap.Duration("duration", duration),
		}

		if userID, exists := c.Get("userID"); exists {
			fields = append(fields, zap.Any("userID", userID))
		}

		if status >= 500 {
			logger.Error("HTTP Request", fields...)
		} else if status >= 400 {
			logger.Warn("HTTP Request", fields...)
		} else {
			logger.Info("HTTP Request", fields...)
		}
	}
}
