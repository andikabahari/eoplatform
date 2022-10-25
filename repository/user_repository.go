package repository

import (
	"github.com/andikabahari/eoplatform/model"
	"gorm.io/gorm"
)

type UserRepository interface {
	Create(user *model.User)
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *userRepository {
	return &userRepository{db}
}

func (r *userRepository) Create(user *model.User) {
	r.db.Debug().Save(user)
}
