package auth

import (
	"github.com/Skrwbrgle/go-backend-learn/internal/handler"
	"github.com/Skrwbrgle/go-backend-learn/internal/service"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func RegisterAuthRoutes(rg *gin.RouterGroup, authService service.AuthService, logger *zap.Logger) {
	h := handler.NewAuthHandler(authService, logger )

	publicGroup := rg.Group("/auth")

	{
		publicGroup.POST("/login", h.Login)
		publicGroup.POST("/register", h.Register) 
		publicGroup.POST("/logout", h.Logout)
	}

}