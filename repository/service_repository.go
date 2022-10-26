package repository

import (
	"fmt"

	"github.com/andikabahari/eoplatform/model"
	"github.com/andikabahari/eoplatform/request"
	"gorm.io/gorm"
)

type ServiceRepository interface {
	Get(services *[]model.Service, keyword string)
	Find(service *model.Service, id string)
	Create(service *model.Service)
}

type serviceRepository struct {
	db *gorm.DB
}

func NewServiceRepository(db *gorm.DB) *serviceRepository {
	return &serviceRepository{db}
}

func (r *serviceRepository) Get(services *[]model.Service, keyword string) {
	if keyword != "" {
		keyword = fmt.Sprintf("%%%s%%", keyword)
		r.db.Debug().Preload("User").Where("name LIKE ? OR description LIKE ?", keyword, keyword).Find(services)
	} else {
		r.db.Debug().Preload("User").Find(services)
	}
}

func (r *serviceRepository) Find(service *model.Service, id string) {
	r.db.Debug().Preload("User").Where("id = ?", id).Find(service)
}

func (r *serviceRepository) Create(service *model.Service) {
	r.db.Debug().Save(service)
}

func (r *serviceRepository) Update(service *model.Service, req *request.UpdateServiceRequest) {
	service.Name = req.Name
	service.Description = req.Description
	service.Cost = req.Cost

	r.db.Debug().Omit("User").Save(service)
}
