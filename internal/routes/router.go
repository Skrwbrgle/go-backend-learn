package routes

import (
	"github.com/Skrwbrgle/go-backend-learn/pkg/health_check"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"gorm.io/gorm"
)

func SetupRouter(r *gin.Engine, db *gorm.DB) {
	// ============ METRICS ENDPOINT ============
	// Prometheus akan scrape ini untuk collect metrics
	// Format: http://localhost:8080/metrics
	r.GET("/metrics", gin.WrapH(promhttp.Handler()))

	// ============ HEALTH CHECK ============
	r.GET("/ping", health_check.HealthCheck)

	// ============ API ROUTES ============
	api := r.Group("/go-api")
	{
		RegisterAuthModule(api, db)
		RegisterUserModule(api, db)
	}
}
