package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Pagination struct {
	Page      int `json:"page"`
	Limit     int `json:"limit"`
	TotalRows int64 `json:"totalRows"`
	TotalPages int  `json:"totalPages"`
}

type User struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	Name      string    `json:"name" validate:"required,min=3"`
	Email     string    `json:"email" gorm:"unique" validate:"required,email"`
	Password  string    `json:"-" validate:"required,min=6"`
	Role      string    `json:"role" validate:"oneof=super_admin admin viewer user"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
	if u.ID == uuid.Nil {
		u.ID = uuid.New()
	}
	return
}
