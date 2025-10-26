package routes

import (
	"gorm.io/gorm"

	"github.com/Skrwbrgle/go-backend-learn/internal/domain/auth"
	"github.com/Skrwbrgle/go-backend-learn/internal/domain/user"
	auth_routes "github.com/Skrwbrgle/go-backend-learn/internal/routes/auth"
	user_routes "github.com/Skrwbrgle/go-backend-learn/internal/routes/user"
	"github.com/Skrwbrgle/go-backend-learn/pkg/logger"
	"github.com/gin-gonic/gin"
)

func RegisterAuthModule(r *gin.RouterGroup, db *gorm.DB) {
	repo := auth.NewAuthRepository(db, logger.Log)
	service := auth.NewAuthService(repo, logger.Log)
	auth_routes.RegisterAuthRoutes(r, service, logger.Log)
}

func RegisterUserModule(r *gin.RouterGroup, db *gorm.DB) {
	repo := user.NewUserRepository(db, logger.Log)
	service := user.NewUserService(repo, logger.Log)
	user_routes.RegisterUserRoutes(r, service, logger.Log)
}
