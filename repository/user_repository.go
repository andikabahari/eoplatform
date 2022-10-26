package repository

import (
	"github.com/andikabahari/eoplatform/model"
	"gorm.io/gorm"
)

type UserRepository interface {
	Find(user *model.User, id string)
	FindByUsername(user *model.User, username string)
	Create(user *model.User)
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *userRepository {
	return &userRepository{db}
}

func (r *userRepository) Find(user *model.User, id string) {
	r.db.Debug().Where("id = ?", id).First(user)
}

func (r *userRepository) FindByUsername(user *model.User, username string) {
	r.db.Debug().Where("username = ?", username).First(user)
}

func (r *userRepository) Create(user *model.User) {
	r.db.Debug().Save(user)
}
