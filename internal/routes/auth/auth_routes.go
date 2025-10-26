package auth_routes

import (
	"github.com/Skrwbrgle/go-backend-learn/internal/domain/auth"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func RegisterAuthRoutes(rg *gin.RouterGroup, authService auth.AuthService, logger *zap.Logger) {
	h := auth.NewAuthHandler(authService, logger)

	publicGroup := rg.Group("/auth")

	{
		publicGroup.POST("/login", h.Login)
		publicGroup.POST("/register", h.Register)
		publicGroup.POST("/logout", h.Logout)
	}

}
