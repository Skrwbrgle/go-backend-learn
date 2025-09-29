package user

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"github.com/Skrwbrgle/go-backend-learn/internal/repository"
	"github.com/Skrwbrgle/go-backend-learn/internal/service"
)

func RegisterUserModule(r *gin.RouterGroup, db *gorm.DB) {
	repo := repository.NewUserRepository(db)
	service := service.NewUserService(repo)
	RegisterUserRoutes(r, service)
}
