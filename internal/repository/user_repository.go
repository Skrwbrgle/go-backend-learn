package repository

import (
	"github.com/Skrwbrgle/go-backend-learn/internal/model"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type UserRepository interface {
	Create(user *model.User) error
	FindAll(page, limit int, search string) ([]model.User, int64, error)
	FindByID(id uuid.UUID) (*model.User, error)
	Update(user *model.User) error
	Delete(id uuid.UUID) error
}

type userRepo struct {
	db *gorm.DB
	logger *zap.Logger
}

func NewUserRepository(db *gorm.DB, logger *zap.Logger) UserRepository {
	return &userRepo{db, logger}
}

func (r *userRepo) Create(user *model.User) error {
	r.logger.Info("Creating user", zap.String("email", user.Email))
	return r.db.Create(user).Error
}

func (r *userRepo) FindAll(page, limit int, search string) ([]model.User, int64, error) {
	var users []model.User
	var totalRows int64

	query := r.db.Model(&model.User{})
	if search != "" {
		query = query.Where("name ILIKE  ? OR email ILIKE  ? OR role ILIKE  ?", "%"+search+"%", "%"+search+"%", "%"+search+"%")
	}

	if err := query.Count(&totalRows).Error; err != nil {
		r.logger.Error("Failed to count users", zap.Error(err))
		return nil, 0, err
	}

	offset := (page - 1) * limit
	if err := query.Limit(limit).Offset(offset).Find(&users).Error; err != nil {
		r.logger.Error("Failed to fetch users", zap.Error(err))
		return nil, 0, err
	}

	r.logger.Info("Users fetched successfully", zap.Int("page", page), zap.Int("limit", limit), zap.String("search", search))
	return users, totalRows, nil
}

func (r *userRepo) FindByID(id uuid.UUID) (*model.User, error) {
	var user model.User
	if err := r.db.First(&user, "id = ?", id).Error; err != nil {
		return nil, err
	}

	r.logger.Info("User fetched successfully", zap.String("id", id.String()))
	return &user, nil
}

func (r *userRepo) Update(user *model.User) error {
	r.logger.Info("Updating user", zap.String("email", user.Email))
	return r.db.Save(user).Error
}

func (r *userRepo) Delete(id uuid.UUID) error {
	r.logger.Info("Deleting user", zap.String("id", id.String()))
	return r.db.Delete(&model.User{}, "id = ?", id).Error
}