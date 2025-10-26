package auth

import (
	"errors"

	"github.com/Skrwbrgle/go-backend-learn/internal/domain/user"
	"github.com/Skrwbrgle/go-backend-learn/pkg/bcrypt"
	"github.com/Skrwbrgle/go-backend-learn/pkg/jwt"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type AuthService interface {
	Login(email, password string) (*LoginResponse, error)
	Register(req RegisterRequest) (*LoginResponse, error)
	Logout(ctx *gin.Context)
}

type authService struct {
	repo  	AuthRepository
	logger *zap.Logger
}

func NewAuthService(repo AuthRepository, logger *zap.Logger) AuthService {
	return &authService{repo, logger}
}

func (s *authService) Login(email, password string) (*LoginResponse, error) {
	user, err := s.repo.FindByEmail(email)
	if err != nil {
		s.logger.Warn("Login failed", zap.String("email", email), zap.Error(err))
		return nil, err
	}

	if !bcrypt.CheckPasswordHash(password, user.Password) {
		s.logger.Warn("Login failed", zap.String("email", email))
		return nil, errors.New("invalid email or password")
	}

	token, err := jwt.GenerateToken(user.ID.String(), user.Email, user.Role)
	if err != nil {
		s.logger.Warn("Login failed", zap.String("email", email), zap.Error(err))
		return nil, err
	}

	s.logger.Info("User logged in", zap.String("email", email))
	return &LoginResponse{
		Token: token,
		ID:    user.ID.String(),
		Name:  user.Name,
		Role:  user.Role,
	}, nil
}

func (s *authService) Register(req RegisterRequest) (*LoginResponse, error) {
	_, err := s.repo.FindByEmail(req.Email)
	if err == nil {
		s.logger.Warn("Register failed", zap.String("email", req.Email))
		return nil, errors.New("email already registered")
	}

	hashedPassword := bcrypt.HashPassword(req.Password)

	user := user.User{
		Name:     req.Name,
		Email:    req.Email,
		Password: hashedPassword,
		Role:     "user",
	}

	if err := s.repo.Register(&user); err != nil {
		s.logger.Warn("Register failed", zap.String("email", req.Email), zap.Error(err))
		return nil, err
	}

	token, err := jwt.GenerateToken(user.ID.String(), user.Email, user.Role)
	if err != nil {
		s.logger.Warn("Register failed", zap.String("email", req.Email), zap.Error(err))
		return nil, err
	}

	s.logger.Info("User registered", zap.String("email", req.Email))
	return &LoginResponse{
		Token: token,
		ID:    user.ID.String(),
		Name:  user.Name,
		Role:  user.Role,
	}, nil
}

func (s *authService) Logout(c *gin.Context) {
	s.logger.Info("User logged out", zap.String("email", c.GetString("email")))
	c.SetCookie("token", "", -1, "/", "", false, true)
}
