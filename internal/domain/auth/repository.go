package auth

import (
	"github.com/Skrwbrgle/go-backend-learn/internal/domain/user"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type AuthRepository interface {
	FindByEmail(email string) (*user.User, error)
	Register(user *user.User) error
}

type authRepo struct {
	db     *gorm.DB
	logger *zap.Logger
}

func NewAuthRepository(db *gorm.DB, logger *zap.Logger) AuthRepository {
	return &authRepo{db, logger}
}

func (r *authRepo) FindByEmail(email string) (*user.User, error) {
	var user user.User
	if err := r.db.Where("email = ?", email).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *authRepo) Register(user *user.User) error {
	return r.db.Create(user).Error
}
