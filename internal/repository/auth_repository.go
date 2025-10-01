package repository

import (
	"github.com/Skrwbrgle/go-backend-learn/internal/model"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type AuthRepository interface {
	FindByEmail(email string) (*model.User, error)
	Register(user *model.User) error
}

type authRepo struct {
	db *gorm.DB
	logger *zap.Logger
}

func NewAuthRepository(db *gorm.DB, logger *zap.Logger) AuthRepository {
	return &authRepo{db, logger}
}

func (r *authRepo) FindByEmail(email string) (*model.User, error) {
	var user model.User
	if err := r.db.Where("email = ?", email).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *authRepo) Register(user *model.User) error {
	return r.db.Create(user).Error
}