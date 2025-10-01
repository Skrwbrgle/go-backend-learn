package routes

import (
	"github.com/Skrwbrgle/go-backend-learn/internal/handler"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetupRouter(r *gin.Engine, db *gorm.DB) {
	api := r.Group("/go-api")
	{
		api.GET("/ping", handler.HealthCheck)

		RegisterAuthModule(api, db)
		RegisterUserModule(api, db) 
	}
}
