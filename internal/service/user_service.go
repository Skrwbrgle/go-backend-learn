package service

import (
	"github.com/Skrwbrgle/go-backend-learn/internal/model"
	"github.com/Skrwbrgle/go-backend-learn/internal/repository"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type UserService interface {
	CreateUser(user *model.User) error
	GetUsers(page, limit int, search string) ([]model.User, model.Pagination, error)
	GetUserByID(id uuid.UUID) (*model.User, error)
	UpdateUser(id uuid.UUID, user *model.User) error
	DeleteUser(id uuid.UUID) error
}

type userService struct {
	repo   repository.UserRepository
	logger *zap.Logger
}

func NewUserService(repo repository.UserRepository, logger *zap.Logger) UserService {
	return &userService{repo, logger}
}

func (s *userService) CreateUser(user *model.User) error {
	return s.repo.Create(user)
}

func (s *userService) GetUsers(page, limit int, search string) ([]model.User, model.Pagination, error) {
	users, totalRows, err := s.repo.FindAll(page, limit, search)
	if err != nil {
		return nil, model.Pagination{}, err
	}

	totalPages := int((totalRows + int64(limit) - 1) / int64(limit))
	pagination := model.Pagination{
		Page:       page,
		Limit:      limit,
		TotalRows:  totalRows,
		TotalPages: totalPages,
	}

	return users, pagination, nil
}

func (s *userService) GetUserByID(id uuid.UUID) (*model.User, error) {
	return s.repo.FindByID(id)
}

func (s *userService) UpdateUser(id uuid.UUID, user *model.User) error {
	existing, err := s.repo.FindByID(id)
	if err != nil {
		return err
	}

	existing.Name = user.Name
	existing.Email = user.Email
	return s.repo.Update(existing)
}

func (s *userService) DeleteUser(id uuid.UUID) error {
	return s.repo.Delete(id)
}
