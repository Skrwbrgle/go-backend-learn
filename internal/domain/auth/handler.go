package auth

import (
	"github.com/Skrwbrgle/go-backend-learn/pkg/response"
	"github.com/Skrwbrgle/go-backend-learn/pkg/validator"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type AuthHandler struct {
	service AuthService
	logger  *zap.Logger
}

func NewAuthHandler(service AuthService, logger *zap.Logger) *AuthHandler {
	return &AuthHandler{service, logger}
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Warn("Invalid login payload", zap.Error(err))
		response.BadRequest(c, err.Error())
		return
	}

	if err := validator.Validate.Struct(req); err != nil {
		h.logger.Warn("Login validation failed", zap.Error(err))
		response.BadRequest(c, err.Error())
		return
	}

	loginResp, err := h.service.Login(req.Email, req.Password)
	if err != nil {
		h.logger.Warn("Login failed", zap.String("email", req.Email), zap.Error(err))
		response.Unauthorized(c, "Invalid email or password")
		return
	}

	h.logger.Info("User logged in", zap.String("email", req.Email))
	response.Success(c, loginResp, "User logged in successfully")
}

func (h *AuthHandler) Logout(c *gin.Context) {
	h.service.Logout(c)
	h.logger.Info("User logged out", zap.String("email", c.GetString("email")))
	response.Success(c, nil, "User logged out successfully")
}

func (h *AuthHandler) Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Warn("Invalid register payload", zap.Error(err))
		response.BadRequest(c, err.Error())
		return
	}

	if err := validator.Validate.Struct(req); err != nil {
		h.logger.Warn("Register validation failed", zap.Error(err))
		response.BadRequest(c, err.Error())
		return
	}

	registerResp, err := h.service.Register(req)
	if err != nil {
		h.logger.Warn("Register failed", zap.String("email", req.Email), zap.Error(err))
		response.InternalError(c, err.Error())
		return
	}

	h.logger.Info("User registered", zap.String("email", req.Email))
	response.Success(c, registerResp, "User registered successfully")
}
