package user

import (
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type UserService interface {
	CreateUser(user *User) error
	GetUsers(page, limit int, search string) ([]User, Pagination, error)
	GetUserByID(id uuid.UUID) (*User, error)
	UpdateUser(id uuid.UUID, user *User) error
	DeleteUser(id uuid.UUID) error
}

type userService struct {
	repo   UserRepository
	logger *zap.Logger
}

func NewUserService(repo UserRepository, logger *zap.Logger) UserService {
	return &userService{repo, logger}
}

func (s *userService) CreateUser(user *User) error {
	return s.repo.Create(user)
}

func (s *userService) GetUsers(page, limit int, search string) ([]User, Pagination, error) {
	users, totalRows, err := s.repo.FindAll(page, limit, search)
	if err != nil {
		return nil, Pagination{}, err
	}

	totalPages := int((totalRows + int64(limit) - 1) / int64(limit))
	pagination := Pagination{
		Page:       page,
		Limit:      limit,
		TotalRows:  totalRows,
		TotalPages: totalPages,
	}

	return users, pagination, nil
}

func (s *userService) GetUserByID(id uuid.UUID) (*User, error) {
	return s.repo.FindByID(id)
}

func (s *userService) UpdateUser(id uuid.UUID, user *User) error {
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
