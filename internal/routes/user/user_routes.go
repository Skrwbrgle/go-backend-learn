package user

import (
	"github.com/Skrwbrgle/go-backend-learn/internal/handler"
	"github.com/Skrwbrgle/go-backend-learn/internal/middleware"
	"github.com/Skrwbrgle/go-backend-learn/internal/service"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func RegisterUserRoutes(rg *gin.RouterGroup, userService service.UserService, logger *zap.Logger) {
	h := handler.NewUserHandler(userService, logger)

	userGroup := rg.Group("/users")
	userGroup.Use(middleware.AuthMiddleware())
	{
		userGroup.POST("/", middleware.RoleMiddleware("admin", "super_admin"), h.CreateUser)
		userGroup.GET("/", middleware.RoleMiddleware("admin", "super_admin", "viewer"), h.GetUsers)
		userGroup.GET("/:id", middleware.RoleMiddleware("admin", "super_admin", "viewer"), h.GetUserByID)
		userGroup.PUT("/:id", middleware.RoleMiddleware("admin", "super_admin"), h.UpdateUser)
		userGroup.DELETE("/:id", middleware.RoleMiddleware("super_admin"), h.DeleteUser)
	}
}
