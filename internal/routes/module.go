package routes

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"github.com/Skrwbrgle/go-backend-learn/internal/repository"
	"github.com/Skrwbrgle/go-backend-learn/internal/routes/auth"
	"github.com/Skrwbrgle/go-backend-learn/internal/routes/user"
	"github.com/Skrwbrgle/go-backend-learn/internal/service"
	"github.com/Skrwbrgle/go-backend-learn/pkg/logger"
)

func RegisterAuthModule(r *gin.RouterGroup, db *gorm.DB) {
	repo := repository.NewAuthRepository(db, logger.Log)
	service := service.NewAuthService(repo, logger.Log)
	auth.RegisterAuthRoutes(r, service, logger.Log)
}

func RegisterUserModule(r *gin.RouterGroup, db *gorm.DB) {
	repo := repository.NewUserRepository(db, logger.Log)
	service := service.NewUserService(repo, logger.Log)
	user.RegisterUserRoutes(r, service, logger.Log)
}
