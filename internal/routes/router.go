package routes

import (
	"github.com/Skrwbrgle/go-backend-learn/pkg/health_check"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetupRouter(r *gin.Engine, db *gorm.DB) {
	api := r.Group("/go-api")
	{
		api.GET("/ping", health_check.HealthCheck)

		RegisterAuthModule(api, db)
		RegisterUserModule(api, db)
	}
}
