package middleware

import (
	"github.com/gin-gonic/gin"
)

// Configurable CORS middleware
func CORSMiddleware(allowedOrigins []string, allowedMethods []string, allowedHeaders []string) gin.HandlerFunc {
	return func(c *gin.Context) {
		origin := c.GetHeader("Origin")

		// cek apakah origin diperbolehkan
		allowed := false
		for _, o := range allowedOrigins {
			if o == "*" || o == origin {
				allowed = true
				break
			}
		}

		if allowed {
			c.Header("Access-Control-Allow-Origin", origin)
			c.Header("Access-Control-Allow-Methods", joinMethods(allowedMethods))
			c.Header("Access-Control-Allow-Headers", joinMethods(allowedHeaders))
			c.Header("Access-Control-Allow-Credentials", "true")
		}

		// Preflight request
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

func joinMethods(arr []string) string {
	if len(arr) == 0 {
		return "*"
	}
	result := ""
	for i, s := range arr {
		if i > 0 {
			result += ", "
		}
		result += s
	}
	return result
}
