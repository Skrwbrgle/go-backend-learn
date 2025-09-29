package user

import (
	"github.com/Skrwbrgle/go-backend-learn/internal/handler"
	"github.com/Skrwbrgle/go-backend-learn/internal/service"
	"github.com/gin-gonic/gin"
)


func RegisterUserRoutes(rg *gin.RouterGroup, userService service.UserService) {
	h := handler.NewUserHandler(userService)

	userGroup := rg.Group("/users")
	{
		userGroup.POST("/", h.CreateUser)
		userGroup.GET("/", h.GetUsers)
		userGroup.GET("/:id", h.GetUserByID) 
		userGroup.PUT("/:id", h.UpdateUser)
		userGroup.DELETE("/:id", h.DeleteUser)
	}
}
