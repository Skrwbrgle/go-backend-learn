package handler

import (
	"strconv"

	"github.com/Skrwbrgle/go-backend-learn/internal/model"
	"github.com/Skrwbrgle/go-backend-learn/internal/service"
	"github.com/Skrwbrgle/go-backend-learn/pkg/response"
	"github.com/Skrwbrgle/go-backend-learn/pkg/validator"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type UserHandler struct {
	service service.UserService
}

func NewUserHandler(service service.UserService) *UserHandler {
	return &UserHandler{service}
}

func (h *UserHandler) CreateUser(c *gin.Context) {
	var user model.User
	if err := c.ShouldBindJSON(&user); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	if err := validator.Validate.Struct(user); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	if err := h.service.CreateUser(&user); err != nil {
		response.InternalError(c, err.Error())
		return
	}

	response.Created(c, user, "User created successfully")
}

func (h *UserHandler) GetUsers(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	users, pagination, err := h.service.GetUsers(page, limit)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}

	response.SuccessWithMeta(c, users, pagination, "Users fetched successfully")
}

func (h *UserHandler) GetUserByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		response.BadRequest(c, "Invalid UUID format")
		return
	}

	user, err := h.service.GetUserByID(id)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}

	response.Success(c, user, "User fetched successfully")
}

func (h *UserHandler) UpdateUser(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		response.BadRequest(c, "Invalid UUID format")
		return
	}

	var input model.User
	if err := c.ShouldBindJSON(&input); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	updateData := model.User{
		Name:     input.Name,
		Email:    input.Email,
		Password: input.Password,
		Role:     input.Role,
	}

	if err := h.service.UpdateUser(id, &updateData); err != nil {
		response.InternalError(c, err.Error())
		return
	}

	response.Updated(c, updateData, "User updated successfully")
}

func (h *UserHandler) DeleteUser(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		response.BadRequest(c, "Invalid UUID format")
		return
	}

	if err := h.service.DeleteUser(id); err != nil {
		response.InternalError(c, err.Error())
		return
	}

	response.Deleted(c, "User deleted successfully")
}