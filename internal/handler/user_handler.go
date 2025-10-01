package handler

import (
	"strconv"

	"github.com/Skrwbrgle/go-backend-learn/internal/dto"
	"github.com/Skrwbrgle/go-backend-learn/internal/model"
	"github.com/Skrwbrgle/go-backend-learn/internal/service"
	"github.com/Skrwbrgle/go-backend-learn/pkg/bcrypt"
	"github.com/Skrwbrgle/go-backend-learn/pkg/response"
	"github.com/Skrwbrgle/go-backend-learn/pkg/validator"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type UserHandler struct {
	service service.UserService
	logger  *zap.Logger
}

func NewUserHandler(service service.UserService, logger *zap.Logger) *UserHandler {
	return &UserHandler{service, logger}
}

func (h *UserHandler) CreateUser(c *gin.Context) {
	var req dto.CreateUserRequest
	
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Warn("Invalid request body", zap.Error(err))
		response.BadRequest(c, err.Error())
		return
	}

	if err := validator.Validate.Struct(req); err != nil {
		h.logger.Warn("Validation failed", zap.Error(err))
		response.BadRequest(c, err.Error())
		return
	}

	if req.Role == "super_admin" {
		h.logger.Warn("Attempt to create super_admin blocked", zap.String("email", req.Email))
		response.Forbidden(c, "You cannot create a super admin")
		return
	}

	user := model.User{
		Name:     req.Name,
		Email:    req.Email,
		Password: bcrypt.HashPassword(req.Password),
		Role:     req.Role,
	}

	if err := h.service.CreateUser(&user); err != nil {
		h.logger.Error("Failed to create user", zap.String("email", req.Email), zap.Error(err))
		response.InternalError(c, err.Error())
		return
	}

	h.logger.Info("User created successfully", zap.String("email", user.Email), zap.String("role", user.Role))
	response.Created(c, dto.UserResponse{
		ID:    user.ID,
		Name:  user.Name,
		Email: user.Email,
		Role:  user.Role,
	}, "User created successfully")
}

func (h *UserHandler) GetUsers(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	search := c.DefaultQuery("search", "")

	users, pagination, err := h.service.GetUsers(page, limit, search)
	if err != nil {
		h.logger.Error("Failed to fetch users", zap.Error(err))
		response.InternalError(c, err.Error())
		return
	}
	
	h.logger.Info("Users fetched successfully", zap.Int("page", page), zap.Int("limit", limit), zap.String("search", search))
	response.SuccessWithMeta(c, users, pagination, "Users fetched successfully")
}

func (h *UserHandler) GetUserByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		h.logger.Warn("Invalid UUID format", zap.Error(err))
		response.BadRequest(c, "Invalid UUID format")
		return
	}

	user, err := h.service.GetUserByID(id)
	if user == nil {
		h.logger.Warn("User not found", zap.String("id", idStr))
		response.NotFound(c, "User not found")
		return
	}

	if err != nil {
		h.logger.Error("Failed to fetch user", zap.String("id", idStr), zap.Error(err))
		response.InternalError(c, err.Error())
		return
	}

	h.logger.Info("User fetched successfully", zap.String("id", idStr))
	response.Success(c, user, "User fetched successfully")
}

func (h *UserHandler) UpdateUser(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		h.logger.Warn("Invalid UUID format", zap.Error(err))
		response.BadRequest(c, "Invalid UUID format")
		return
	}

	var input model.User
	if err := c.ShouldBindJSON(&input); err != nil {
		h.logger.Warn("Invalid request body", zap.Error(err))
		response.BadRequest(c, err.Error())
		return
	}

	if err := validator.Validate.Struct(input); err != nil {
		h.logger.Warn("Validation failed", zap.Error(err))
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
		h.logger.Error("Failed to update user", zap.String("id", idStr), zap.Error(err))
		response.InternalError(c, err.Error())
		return
	}

	h.logger.Info("User updated successfully", zap.String("id", idStr))
	response.Updated(c, updateData, "User updated successfully")
}

func (h *UserHandler) DeleteUser(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		h.logger.Warn("Invalid UUID format", zap.Error(err))
		response.BadRequest(c, "Invalid UUID format")
		return
	}

	if err := h.service.DeleteUser(id); err != nil {
		h.logger.Error("Failed to delete user", zap.String("id", idStr), zap.Error(err))
		response.InternalError(c, err.Error())
		return
	}

	h.logger.Info("User deleted successfully", zap.String("id", idStr))
	response.Deleted(c, "User deleted successfully")
}