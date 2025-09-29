package repository

import (
	"github.com/Skrwbrgle/go-backend-learn/internal/model"
	"github.com/google/uuid"
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
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepo{db}
}

func (r *userRepo) Create(user *model.User) error {
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
		return nil, 0, err
	}

	offset := (page - 1) * limit
	if err := query.Limit(limit).Offset(offset).Find(&users).Error; err != nil {
		return nil, 0, err
	}

	return users, totalRows, nil
}

func (r *userRepo) FindByID(id uuid.UUID) (*model.User, error) {
	var user model.User
	if err := r.db.First(&user, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepo) Update(user *model.User) error {
	return r.db.Save(user).Error
}

func (r *userRepo) Delete(id uuid.UUID) error {
	return r.db.Delete(&model.User{}, "id = ?", id).Error
}