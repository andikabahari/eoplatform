package repository

import (
	"github.com/andikabahari/eoplatform/model"
	"github.com/andikabahari/eoplatform/request"
	"gorm.io/gorm"
)

type IUserRepository interface {
	Find(user *model.User, id uint)
	FindByUsername(user *model.User, username string)
	Create(user *model.User)
	Update(user *model.User, req *request.UpdateUserRequest)
	ResetPassword(user *model.User, password string)
}

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db}
}

func (r *UserRepository) Find(user *model.User, id uint) {
	r.db.Debug().Where("id = ?", id).First(user)
}

func (r *UserRepository) FindByUsername(user *model.User, username string) {
	r.db.Debug().Where("username = ?", username).First(user)
}

func (r *UserRepository) Create(user *model.User) {
	r.db.Debug().Save(user)
}

func (r *UserRepository) Update(user *model.User, req *request.UpdateUserRequest) {
	user.Name = req.Name

	r.db.Debug().Save(user)
}

func (r *UserRepository) ResetPassword(user *model.User, password string) {
	user.Password = password

	r.db.Debug().Save(user)
}
