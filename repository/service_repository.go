package repository

import (
	"github.com/andikabahari/eoplatform/model"
	"gorm.io/gorm"
)

type ServiceRepository interface {
	Create(service *model.Service)
}

type serviceRepository struct {
	db *gorm.DB
}

func NewServiceRepository(db *gorm.DB) *serviceRepository {
	return &serviceRepository{db}
}

func (r *serviceRepository) Create(service *model.Service) {
	r.db.Debug().Save(service)
}
