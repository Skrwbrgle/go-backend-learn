package user_mock

import (
	"github.com/Skrwbrgle/go-backend-learn/internal/domain/user"
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
)

type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) Create(u *user.User) error {
	args := m.Called(u)
	return args.Error(0)
}

func (m *MockUserRepository) FindAll(page, limit int, search string) ([]user.User, int64, error) {
	args := m.Called(page, limit, search)
	return args.Get(0).([]user.User), args.Get(1).(int64), args.Error(2)
}

func (m *MockUserRepository) FindByID(id uuid.UUID) (*user.User, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*user.User), args.Error(1)
}

func (m *MockUserRepository) Update(u *user.User) error {
	args := m.Called(u)
	return args.Error(0)
}

func (m *MockUserRepository) Delete(id uuid.UUID) error {
	args := m.Called(id)
	return args.Error(0)
}