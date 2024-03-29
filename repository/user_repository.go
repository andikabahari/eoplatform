package repository

import (
	"github.com/andikabahari/eoplatform/model"
	"github.com/andikabahari/eoplatform/request"
	"gorm.io/gorm"
)

type UserRepository interface {
	Find(user *model.User, id uint)
	FindByUsername(user *model.User, username string)
	Create(user *model.User)
	Update(user *model.User, req *request.UpdateUserRequest)
	ResetPassword(user *model.User, password string)
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db}
}

func (r *userRepository) Find(user *model.User, id uint) {
	r.db.Debug().Where("id = ?", id).First(user)
}

func (r *userRepository) FindByUsername(user *model.User, username string) {
	r.db.Debug().Where("username = ?", username).First(user)
}

func (r *userRepository) Create(user *model.User) {
	r.db.Debug().Save(user)
}

func (r *userRepository) Update(user *model.User, req *request.UpdateUserRequest) {
	user.Name = req.Name

	r.db.Debug().Save(user)
}

func (r *userRepository) ResetPassword(user *model.User, password string) {
	user.Password = password

	r.db.Debug().Save(user)
}
