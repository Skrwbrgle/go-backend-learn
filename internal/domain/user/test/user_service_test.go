package user_test

import (
	"errors"
	"testing"

	"github.com/Skrwbrgle/go-backend-learn/internal/domain/user"
	user_mock "github.com/Skrwbrgle/go-backend-learn/internal/domain/user/mock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"
)

func TestCreateUser(t *testing.T) {
	mockRepo := new(user_mock.MockUserRepository)
	logger, _ := zap.NewDevelopment()
	service := user.NewUserService(mockRepo, logger)

	u := &user.User{
		ID:    uuid.New(),
		Name:  "John Doe",
		Email: "john@example.com",
	}

	mockRepo.On("Create", u).Return(nil)

	err := service.CreateUser(u)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestGetUsers(t *testing.T) {
	mockRepo := new(user_mock.MockUserRepository)
	logger, _ := zap.NewDevelopment()
	service := user.NewUserService(mockRepo, logger)

	expectedUsers := []user.User{
		{ID: uuid.New(), Name: "Alice"},
		{ID: uuid.New(), Name: "Bob"},
	}

	mockRepo.On("FindAll", 1, 10, "").Return(expectedUsers, int64(2), nil)

	users, pagination, err := service.GetUsers(1, 10, "")

	assert.NoError(t, err)
	assert.Equal(t, 2, len(users))
	assert.Equal(t, 1, pagination.Page)
	assert.Equal(t, 1, pagination.TotalPages)
	mockRepo.AssertExpectations(t)
}

func TestGetUserByID(t *testing.T) {
	mockRepo := new(user_mock.MockUserRepository)
	logger, _ := zap.NewDevelopment()
	service := user.NewUserService(mockRepo, logger)

	id := uuid.New()
	expectedUser := &user.User{ID: id, Name: "Test"}

	mockRepo.On("FindByID", id).Return(expectedUser, nil)

	result, err := service.GetUserByID(id)

	assert.NoError(t, err)
	assert.Equal(t, expectedUser, result)
	mockRepo.AssertExpectations(t)
}

func TestUpdateUser(t *testing.T) {
	mockRepo := new(user_mock.MockUserRepository)
	logger, _ := zap.NewDevelopment()
	service := user.NewUserService(mockRepo, logger)

	id := uuid.New()
	existing := &user.User{ID: id, Name: "Old", Email: "old@example.com"}
	updated := &user.User{Name: "New", Email: "new@example.com"}

	mockRepo.On("FindByID", id).Return(existing, nil)
	mockRepo.On("Update", mock.AnythingOfType("*user.User")).Return(nil)

	err := service.UpdateUser(id, updated)

	assert.NoError(t, err)
	assert.Equal(t, "New", existing.Name)
	assert.Equal(t, "new@example.com", existing.Email)
	mockRepo.AssertExpectations(t)
}

func TestDeleteUser(t *testing.T) {
	mockRepo := new(user_mock.MockUserRepository)
	logger, _ := zap.NewDevelopment()
	service := user.NewUserService(mockRepo, logger)

	id := uuid.New()
	mockRepo.On("Delete", id).Return(nil)

	err := service.DeleteUser(id)
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestGetUserByID_NotFound(t *testing.T) {
	mockRepo := new(user_mock.MockUserRepository)
	logger, _ := zap.NewDevelopment()
	service := user.NewUserService(mockRepo, logger)

	id := uuid.New()
	mockRepo.On("FindByID", id).Return(nil, errors.New("not found"))

	result, err := service.GetUserByID(id)

	assert.Error(t, err)
	assert.Nil(t, result)
}
