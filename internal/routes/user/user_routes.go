package user_routes

import (
	"github.com/Skrwbrgle/go-backend-learn/internal/domain/user"
	"github.com/Skrwbrgle/go-backend-learn/internal/middleware"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func RegisterUserRoutes(rg *gin.RouterGroup, userService user.UserService, logger *zap.Logger) {
	h := user.NewUserHandler(userService, logger)

	privateGroup := rg.Group("/users")
	privateGroup.Use(middleware.AuthMiddleware())
	{
		privateGroup.POST("/", middleware.RoleMiddleware("admin", "super_admin"), h.CreateUser)
		privateGroup.GET("/", middleware.RoleMiddleware("admin", "super_admin", "viewer"), h.GetUsers)
		privateGroup.GET("/:id", middleware.RoleMiddleware("admin", "super_admin", "viewer"), h.GetUserByID)
		privateGroup.PUT("/:id", middleware.RoleMiddleware("admin", "super_admin"), h.UpdateUser)
		privateGroup.DELETE("/:id", middleware.RoleMiddleware("super_admin"), h.DeleteUser)
	}
}
